package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/tythe-protocol/tythe/coinbase"
	"github.com/tythe-protocol/tythe/conf"
	"github.com/tythe-protocol/tythe/dep"
	"github.com/tythe-protocol/tythe/paypal"

	"github.com/attic-labs/noms/go/d"
	homedir "github.com/mitchellh/go-homedir"
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
	cacheDir := cacheDirFlag(c.cmd)
	sandbox := sandboxFlag(c.cmd)
	totalAmount := c.cmd.Arg("totalAmount", "amount to divide amongst the dependent packages").Required().Float64()
	roots := c.cmd.Arg("package", "one or more root packages to crawl").Required().URLList()

	c.handler = func() {
		tythed := map[string]*conf.Config{}
		totalDeps := 0

		for _, r := range *roots {
			p, err := resolvePackage(r, *cacheDir)
			d.CheckErrorNoUsage(err)

			ds, err := dep.List(p)
			d.CheckErrorNoUsage(err)

			for _, dep := range ds {
				if _, ok := tythed[dep.Name]; ok {
					continue
				}

				if dep.Conf != nil {
					tythed[dep.Name] = dep.Conf
				}

				totalDeps++
			}
		}

		fmt.Printf("Found %d total deps (%d tythed)\n", totalDeps, len(tythed))

		spend := *totalAmount * float64(len(tythed)) / float64(totalDeps)
		fmt.Printf("Ready to send %.2f?\n", spend)
		confirmContinue()

		cb := map[string]coinbase.Amount{}
		pp := map[string]float64{}

		for _, cfg := range tythed {
			amt := spend / float64(len(tythed))
			if cfg.PayPal != "" {
				pp[cfg.PayPal] += amt
			} else {
				addr := cfg.USDC
				typ := "USDC"
				if addr == "" {
					addr = cfg.BTC
					typ = "BTC"
				}
				a, ok := cb[addr]
				if !ok {
					a = coinbase.Amount{Currency: typ}
				}
				a.Value = amt
				cb[addr] = a
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
	cacheDir := cacheDirFlag(c.cmd)
	url := c.cmd.Arg("package", "File path or URL of the package to list.").Required().URL()

	c.handler = func() {
		dir, err := resolvePackage(*url, *cacheDir)
		d.CheckErrorNoUsage(err)

		deps, err := dep.List(dir)
		d.CheckErrorNoUsage(err)

		for _, d := range deps {
			addr := "<no tythe>"
			if d.Conf != nil {
				addr = d.Conf.USDC
			}
			fmt.Printf("%s %s\n", d, addr)
		}
	}

	return c
}

func payOne(app *kingpin.Application) (c command) {
	c.cmd = app.Command("pay-one", "Pay a single package")
	sandbox := sandboxFlag(c.cmd)
	cacheDir := cacheDirFlag(c.cmd)
	amount := c.cmd.Arg("amount", "Amount to send to the package (in USD).").Required().Float()
	url := c.cmd.Arg("package", "File path or URL of the package to pay.").Required().URL()

	c.handler = func() {
		p, err := resolvePackage(*url, *cacheDir)
		d.CheckErrorNoUsage(err)

		config, err := conf.Read(p)
		d.CheckErrorNoUsage(err)
		if config == nil {
			fmt.Printf("no donate file for package: %s", (*url).String())
			return
		}

		fmt.Printf("Found donate file in %s:\n", (*url).String())
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		err = enc.Encode(config)
		d.CheckError(err)

		sendOneImpl(*amount, config.PayPal, config.USDC, *sandbox)
	}

	return
}

func send(app *kingpin.Application) (c command) {
	c.cmd = app.Command("send", "Sends money to the specified address (for testing/development)")
	sandbox := sandboxFlag(c.cmd)
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

func cacheDirFlag(cmd *kingpin.CmdClause) *string {
	hd, err := homedir.Dir()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	return cmd.Flag("cache-dir", "Directory to write cached repos to during crawling").
		Default(fmt.Sprintf("%s/.tythe", hd)).String()
}

func sandboxFlag(cmd *kingpin.CmdClause) *bool {
	return cmd.Flag("sandbox", "Use the sandbox Coinbase API").Default("false").Bool()
}
