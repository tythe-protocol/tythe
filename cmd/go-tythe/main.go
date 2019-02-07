package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aboodman/go-tythe/packageconfig"
	"github.com/attic-labs/noms/go/d"
	gdax "github.com/preichenberger/go-gdax"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	app := kingpin.New("go-tythe", "A command-line tythe client.")

	pay := app.Command("pay", "Pay a single package")
	url := pay.Arg("package-url", "URL of the package to pay.").Required().URL()
	amount := pay.Arg("amount", "Amount to send to the package (in USD).").Required().Float()

	key := getEnv("TYTHE_COINBASE_API_KEY")
	secret := getEnv("TYTHE_COINBASE_API_SECRET")
	passphrase := getEnv("TYTHE_COINBASE_API_PASSPHRASE")

	kingpin.MustParse(app.Parse(os.Args[1:]))

	config, err := packageconfig.Read(*url)
	d.CheckErrorNoUsage(err)

	fmt.Printf("Found tythe.json in %s:\n", (*url).String())
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	err = enc.Encode(config)
	d.CheckError(err)
	fmt.Printf("Really send $%f (y/n)?\n", *amount)

	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	d.CheckErrorNoUsage(err)

	if line != "y\n" {
		return
	}

	client := gdax.NewClient(secret, key, passphrase)
	params := map[string]interface{}{
		"amount":         *amount,
		"currency":       config.Destination.Type,
		"crypto_address": config.Destination.Address,
	}

	var res map[string]interface{}
	_, err = client.Request("POST", "/withdrawals/crypto", params, &res)
	d.PanicIfError(err)

	fmt.Printf("All done! Coinbase transaction ID: %s\n", res["id"])
}

func getEnv(s string) string {
	v := os.Getenv(s)
	if v == "" {
		fmt.Fprintf(os.Stderr, "Could not find required environment variable: %s\n", s)
		os.Exit(1)
	}
	return v
}
