package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"regexp"

	"github.com/attic-labs/noms/go/d"
	"github.com/pkg/errors"
	gdax "github.com/preichenberger/go-gdax"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

var (
	usdcAddressPattern = regexp.MustCompile("^0x[0-9a-f]{40}$")
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

	config, err := readConfig(*url)
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

// PaymentType is a mechanism that go-tythe can use to move money between parties.
type PaymentType string

const (
	// USDC represents the USDC stablecoin backed by Coinbase and Circle.
	USDC PaymentType = "USDC"
)

// PackageConfig describes the json metadata developers add to their package to
// opt-in to receiving tythes.
type PackageConfig struct {
	Destination struct {
		Type    PaymentType `json:"type"`
		Address string      `json:"address"`
	} `json:"destination"`
}

func readConfig(url *url.URL) (PackageConfig, error) {
	var c PackageConfig

	r, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL:   url.String(),
		Depth: 1,
	})
	if err == nil {
		h, err := r.ResolveRevision("HEAD")
		if err == nil {
			commit, err := r.CommitObject(*h)
			if err == nil {
				t, err := commit.Tree()
				if err == nil {
					f, err := t.File("tythe.json")
					if err == nil {
						r, err := f.Reader()
						if err == nil {
							d := json.NewDecoder(r)
							err = d.Decode(&c)
						}
					}
				}
			}
		}
	}

	if err != nil {
		return PackageConfig{}, errors.Wrapf(err, "Could not get config for repository: %s", url)
	}

	if c.Destination.Type != USDC {
		return PackageConfig{}, fmt.Errorf("invalid tythe.json - destination type: \"%s\" not supported", c.Destination.Type)
	}

	if !usdcAddressPattern.MatchString(c.Destination.Address) {
		return PackageConfig{}, fmt.Errorf("invalid destination address in tythe.json: \"%s\"", c.Destination.Address)
	}

	return c, nil
}

func getEnv(s string) string {
	v := os.Getenv(s)
	if v == "" {
		fmt.Fprintf(os.Stderr, "Could not find required environment variable: %s\n", s)
		os.Exit(1)
	}
	return v
}
