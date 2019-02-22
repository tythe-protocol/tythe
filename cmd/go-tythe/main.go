package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/tythe-protocol/go-tythe/conf"
	"github.com/tythe-protocol/go-tythe/dep"

	"github.com/attic-labs/noms/go/d"
	homedir "github.com/mitchellh/go-homedir"
	gdax "github.com/preichenberger/go-gdax"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type command struct {
	cmd     *kingpin.CmdClause
	handler func()
}

func main() {
	app := kingpin.New("go-tythe", "A command-line tythe client.")

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
		tythedWeight := 0.0
		totalDeps := 0
		totalWeight := 0.0

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
					tythed[p] = dep.Conf
					tythedWeight += 1.0 // TODO: impl weight on the CLI
				}

				totalDeps++
				totalWeight += 1.0
			}
		}

		fmt.Printf("Found %d total deps (%.2f total weight) and %d tythed deps (%.2f weight)\n", totalDeps, totalWeight, len(tythed), tythedWeight)

		spend := *totalAmount * tythedWeight / totalWeight
		fmt.Printf("Ready to send %.2f?\n", spend)
		confirmContinue()

		for _, cfg := range tythed {
			const packageWeight = 1.0
			amt := spend * packageWeight / totalWeight
			sendImpl(amt, cfg.Destination.Address, *sandbox)
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
				addr = d.Conf.Destination.Address
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
	url := c.cmd.Arg("package", "File path or URL of the package to pay.").Required().URL()
	amount := c.cmd.Arg("amount", "Amount to send to the package (in USD).").Required().Float()

	c.handler = func() {
		p, err := resolvePackage(*url, *cacheDir)
		d.CheckErrorNoUsage(err)

		config, err := conf.Read(p)
		d.CheckErrorNoUsage(err)
		if config == nil {
			fmt.Printf("no tythe.json for package: %s", (*url).String())
			return
		}

		fmt.Printf("Found tythe.json in %s:\n", (*url).String())
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		err = enc.Encode(config)
		d.CheckError(err)

		sendImpl(*amount, config.Destination.Address, *sandbox)
	}

	return
}

func send(app *kingpin.Application) (c command) {
	c.cmd = app.Command("send", "Sends USDC to the specified address (for testing/development)")
	sandbox := sandboxFlag(c.cmd)
	address := c.cmd.Arg("address", "USDC address to send to.").Required().String()
	amount := c.cmd.Arg("amount", "Amount to send (in USD).").Required().Float()

	c.handler = func() {
		if !conf.ValidUSDCAddress(*address) {
			fmt.Fprintln(os.Stderr, "Invalid USDC address")
			// TODO: refactor exit handling
			os.Exit(1)
			return
		}

		fmt.Printf("Really send $%f (y/n)?\n", *amount)
		confirmContinue()

		sendImpl(*amount, *address, *sandbox)
	}

	return
}

func sendImpl(amt float64, addr string, sandbox bool) {
	key := getEnv("TYTHE_COINBASE_API_KEY")
	secret := getEnv("TYTHE_COINBASE_API_SECRET")
	passphrase := getEnv("TYTHE_COINBASE_API_PASSPHRASE")

	client := gdax.NewClient(secret, key, passphrase)
	if sandbox {
		client.BaseURL = "https://api-public.sandbox.pro.coinbase.com"
	}
	params := map[string]interface{}{
		"amount":         amt,
		"currency":       conf.USDC,
		"crypto_address": addr,
	}

	var res map[string]interface{}
	_, err := client.Request("POST", "/withdrawals/crypto", params, &res)
	d.PanicIfError(err)

	fmt.Printf("All done! Coinbase transaction ID: %s\n", res["id"])
}

func confirmContinue() {
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	d.CheckErrorNoUsage(err)

	if line != "y\n" {
		os.Exit(0)
	}
}

func getEnv(s string) string {
	v := os.Getenv(s)
	if v == "" {
		fmt.Fprintf(os.Stderr, "Could not find required environment variable: %s\n", s)
		os.Exit(1)
	}
	return v
}

func cacheDirFlag(cmd *kingpin.CmdClause) *string {
	hd, err := homedir.Dir()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	return cmd.Flag("cache-dir", "Directory to write cached repos to during crawling").
		Default(fmt.Sprintf("%s/.go-tythe", hd)).String()
}

func sandboxFlag(cmd *kingpin.CmdClause) *bool {
	return cmd.Flag("sandbox", "Use the sandbox Coinbase API").Default("false").Bool()
}
