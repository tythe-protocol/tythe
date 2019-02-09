package pcfg

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"regexp"

	"github.com/pkg/errors"
	"github.com/tythe-protocol/go-tythe/git"
)

// PaymentType is a mechanism that go-tythe can use to move money between parties.
type PaymentType string

const (
	// USDC represents the USDC stablecoin backed by Coinbase and Circle.
	USDC PaymentType = "USDC"
)

var (
	usdcAddressPattern = regexp.MustCompile("^0x[0-9a-f]{40}$")
)

// PackageConfig describes the json metadata developers add to their package to
// opt-in to receiving tythes.
type PackageConfig struct {
	Destination struct {
		Type    PaymentType `json:"type"`
		Address string      `json:"address"`
	} `json:"destination"`
}

// Read loads the PackageConfig of a package if there is one. Returns nil, nil if no config.
func Read(url *url.URL) (*PackageConfig, error) {
	w := func(err error) error {
		return errors.Wrapf(err, "Could not read config for repo: %s:", url.String())
	}

	dir, err := ioutil.TempDir("", "go-tythe")
	if err != nil {
		return nil, w(err)
	}

	p, err := git.Clone(url, dir)
	if err != nil {
		return nil, w(err)
	}

	f, err := os.Open(path.Join(p, "tythe.json"))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, w(err)
	}

	var c PackageConfig
	d := json.NewDecoder(f)
	err = d.Decode(&c)

	if err != nil {
		return nil, w(err)
	}

	if c.Destination.Type != USDC {
		return nil, fmt.Errorf("invalid tythe.json - destination type: \"%s\" not supported", c.Destination.Type)
	}

	if !ValidUSDCAddress(c.Destination.Address) {
		return nil, fmt.Errorf("invalid destination address in tythe.json: \"%s\"", c.Destination.Address)
	}

	return &c, nil
}

// ValidUSDCAddress returns true if an string is a correctly formated USDC address.
// It doesn't check whether the address actually exists.
func ValidUSDCAddress(addr string) bool {
	return usdcAddressPattern.MatchString(addr)
}
