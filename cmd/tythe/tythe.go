package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/tythe-protocol/tythe/cmd/flags"
	"github.com/tythe-protocol/tythe/coinbase"
	"github.com/tythe-protocol/tythe/conf"
	"github.com/tythe-protocol/tythe/dep"
	"github.com/tythe-protocol/tythe/dep/crawl"
	"github.com/tythe-protocol/tythe/git"
	"github.com/tythe-protocol/tythe/paypal"
	"github.com/tythe-protocol/tythe/utl/status"

	"github.com/attic-labs/noms/go/d"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type command struct {
	cmd     *kingpin.CmdClause
	handler func()
}

func main() {
	app := kingpin.New("tythe", "A command-line tythe client.")

	commands := []command{
		payAll(app),
		payOne(app),
		send(app),
		list(app),
	}

	selected := kingpin.MustParse(app.Parse(os.Args[1:]))
	for _, c := range commands {
		if selected == c.cmd.FullCommand() {
			c.handler()
			break
		}
	}
}

func payAll(app *kingpin.Application) (c command) {
	c.cmd = app.Command("pay-all", "Pay tythes for listed packages and their transitive dependencies")
	cacheDir := flags.CacheDir(c.cmd)
	sandbox := flags.Sandbox(c.cmd)
	totalAmount := c.cmd.Arg("totalAmount", "amount to divide amongst the dependent packages").Required().Float64()
	roots := c.cmd.Arg("package", "one or more root packages to crawl").Required().URLList()

	c.handler = func() {
		tythed := map[dep.ID]*conf.Config{}
		totalDeps := 0

		l := log.New(status.Writer{}, "", 0)
		for _, r := range *roots {
			p, err := git.Resolve(r, *cacheDir, l)
			d.CheckErrorNoUsage(err)

			rs := crawl.Crawl(p, *cacheDir, l)

			for r := range rs {
				if r.Dep == nil {
					continue
				}
				if _, ok := tythed[r.Dep.ID]; ok {
					continue
				}

				if r.Dep.Conf != nil {
					tythed[r.Dep.ID] = r.Dep.Conf
				}

				totalDeps++
			}
		}

		status.Clear()
		fmt.Printf("Found %d total deps (%d tythed)\n", totalDeps, len(tythed))

		spend := *totalAmount * float64(len(tythed)) / float64(totalDeps)
		fmt.Printf("Ready to send %.2f?\n", spend)
		confirmContinue()

		cb := map[string]coinbase.Amount{}
		pp := map[string]float64{}

		for _, cfg := range tythed {
			value := spend / float64(len(tythed))
			if cfg.PayPal != "" {
				pp[cfg.PayPal] += value
			} else {
				pt := cfg.PreferredPaymentType()
				addr := cfg.AddressForType(pt)
				amt, ok := cb[addr]
				if !ok {
					amt = coinbase.Amount{Currency: pt}
				}
				amt.Value = value
				cb[addr] = amt
			}
		}

		fmt.Println()

		if len(pp) > 0 {
			batchID, status, err := paypal.Send(pp, *sandbox)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error from PayPal: %s", err)
			} else {
				fmt.Printf("Sent %d PayPal transactions. BatchID: %s, Status: %s:\n", len(pp), batchID, status)
				for addr, amt := range pp {
					fmt.Printf("%s (%.2f)\n", addr, amt)
				}
			}
			fmt.Println()
		}

		if len(cb) > 0 {
			fmt.Println()
			srs, err := coinbase.Send(cb, *sandbox)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error from Coinbase: %s", err)
			}
			fmt.Printf("Sent %d Coinbase transactions:\n", len(srs))
			for addr, sr := range srs {
				fmt.Printf("%s (%.2f): %s\n", addr, cb[addr].Value, sr.String())
			}
			fmt.Println()
		}

	}

	return c
}

func list(app *kingpin.Application) (c command) {
	c.cmd = app.Command("list", "List transitive dependencies of a package")
	cacheDir := flags.CacheDir(c.cmd)
	url := c.cmd.Arg("package", "File path or URL of the package to list.").Required().URL()

	c.handler = func() {
		w := status.Writer{}
		l := log.New(w, "", 0)
		dir, err := git.Resolve(*url, *cacheDir, l)
		d.CheckErrorNoUsage(err)

		rs := crawl.Crawl(dir, *cacheDir, l)
		for r := range rs {
			if r.Dep == nil {
				continue
			}
			addr := "<no tythe>"
			if r.Dep.Conf != nil {
				addr = r.Dep.Conf.USDC
			}
			status.Clear()
			status.Printf("%s %s", r.Dep, addr)
			status.Enter()
		}
	}

	return c
}

func payOne(app *kingpin.Application) (c command) {
	c.cmd = app.Command("pay-one", "Pay a single package")
	sandbox := flags.Sandbox(c.cmd)
	cacheDir := flags.CacheDir(c.cmd)
	amount := c.cmd.Arg("amount", "Amount to send to the package (in USD).").Required().Float()
	url := c.cmd.Arg("package", "File path or URL of the package to pay.").Required().URL()

	c.handler = func() {
		l := log.New(status.Writer{}, "", 0)
		p, err := git.Resolve(*url, *cacheDir, l)
		d.CheckErrorNoUsage(err)
		status.Clear()

		config, err := conf.Read(p)
		d.CheckErrorNoUsage(err)
		if config == nil {
			fmt.Printf("no donate file for package: %s\n", (*url).String())
			return
		}

		fmt.Printf("Found donate file in %s:\n", (*url).String())
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		err = enc.Encode(config)
		d.CheckError(err)

		pt := config.PreferredPaymentType()
		addr := config.AddressForType(pt)
		sendOneImpl(*amount, pt, addr, *sandbox)
	}

	return
}

func send(app *kingpin.Application) (c command) {
	c.cmd = app.Command("send", "Sends money to the specified address (for testing/development)")
	sandbox := flags.Sandbox(c.cmd)
	paymentType := c.cmd.Arg("type", "Payment type to use").Required().HintOptions("BTC", "PayPal", "USDC").String()
	address := c.cmd.Arg("address", "Address to send to.").Required().String()
	amount := c.cmd.Arg("amount", "Amount to send (in USD).").Required().Float()

	c.handler = func() {
		// TODO: validate paypal, bitcoin addresses

		if *paymentType == "USDC" && !conf.ValidUSDCAddress(*address) {
			fmt.Fprintln(os.Stderr, "Invalid USDC address")
			// TODO: refactor exit handling
			os.Exit(1)
			return
		}

		fmt.Printf("Really send $%.2f (y/n)?\n", *amount)
		confirmContinue()

		sendOneImpl(*amount, *paymentType, *address, *sandbox)
	}

	return
}

func sendOneImpl(amt float64, paymentType string, address string, sandbox bool) {
	if paymentType == "PayPal" {
		batchID, status, err := paypal.Send(map[string]float64{address: amt}, sandbox)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failure: %s\n", err.Error())
			return
		}
		fmt.Printf("Success. PayPal Batch ID: %s, Status: %s\n", batchID, status)
	} else {
		srs, err := coinbase.Send(
			map[string]coinbase.Amount{address: coinbase.Amount{Currency: paymentType, Value: amt}},
			sandbox)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failure: %s\n", err.Error())
			return
		}
		sr := srs[address]
		if sr.Error != nil {
			fmt.Fprintf(os.Stderr, "Failure: %s\n", sr.Error.Error())
			return
		}
		fmt.Printf("Success. Coinbase Transaction ID: %s\n", sr.TransactionID)
	}
}

func confirmContinue() {
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	d.CheckErrorNoUsage(err)

	if line != "y\n" {
		os.Exit(0)
	}
}
