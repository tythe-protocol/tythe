package conf

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"regexp"

	"github.com/pkg/errors"
)

// PaymentType is a mechanism that go-tythe can use to move money between parties.
type PaymentType string

const (
	// USDC represents the USDC stablecoin backed by Coinbase and Circle.
	USDC      PaymentType = "USDC"
	TytheFile string      = ".tythe"
)

var (
	usdcAddressPattern = regexp.MustCompile("^0x[0-9a-f]{40}$")
)

// Config describes the json metadata developers add to their package to opt-in
// to receiving tythes.
type Config struct {
	Destination struct {
		Type    PaymentType `json:"type"`
		Address string      `json:"address"`
	} `json:"destination"`
}

// Read loads the Config of a package if there is one. Returns nil, nil if no config.
func Read(dir string) (*Config, error) {
	w := func(err error) error {
		return errors.Wrapf(err, "Could not read config for package: %s:", dir)
	}

	f, err := os.Open(path.Join(dir, TytheFile))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, w(err)
	}

	var c Config
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
