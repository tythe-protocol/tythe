// Package conf is for working with dot-donate files.
// See: https://github.com/aboodman/dot-donate.
package conf

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"regexp"

	"github.com/pkg/errors"
)

// PaymentType is a mechanism that tythe can use to move money between parties.
type PaymentType string

const (
	// DonateFile is the name of the dot-donate file.
	DonateFile string = ".donate"

	// PaymentTypeBTC indicates payment should happen using the Bitcoin network.
	PaymentTypeBTC string = "BTC"
	// PaymentTypePayPal indicates payment should happen using PayPal.
	PaymentTypePayPal string = "PayPal"
	// PaymentTypeUSDC indicates payment should happen using Coinbase's stablecoin, USDC.
	PaymentTypeUSDC string = "USDC"
	// PaymentTypeNone indicates no payment type.
	PaymentTypeNone string = ""
)

var (
	// ErrNoSupportedPaymentType is returned by conf.Read() when the config looks valid, but it
	// doesn't contain a payment type tythe supports.
	ErrNoSupportedPaymentType error = errors.New("No supported payment type")

	usdcAddressPattern = regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
)

// Config describes the json metadata developers add to their package to opt-in
// to receiving tythes.
type Config struct {
	BTC    string `json:"BTC,omitempty"`
	PayPal string `json:"PayPal,omitempty"`
	USDC   string `json:"USDC,omitempty"`
}

// PreferredPaymentType returns the payment type that should be used, in the case that multiple
// are specified.
func (c Config) PreferredPaymentType() string {
	switch {
	case c.BTC != "":
		return PaymentTypeBTC
	case c.USDC != "":
		return PaymentTypeUSDC
	case c.PayPal != "":
		return PaymentTypePayPal
	}
	return PaymentTypeNone
}

// AddressForType returns the address to be used for a particular payment type.
func (c Config) AddressForType(paymentType string) string {
	switch paymentType {
	case PaymentTypeBTC:
		return c.BTC
	case PaymentTypePayPal:
		return c.PayPal
	case PaymentTypeUSDC:
		return c.USDC
	}
	return ""
}

// Read loads the Config of a package if there is one. Returns nil, nil if no config.
func Read(dir string) (*Config, error) {
	w := func(err error, msg string, args ...interface{}) error {
		return errors.Wrapf(err, "Could not read config for package %s: %s", dir, fmt.Sprintf(msg, args...))
	}

	df, err := os.Open(path.Join(dir, DonateFile))
	if err != nil {
		return nil, nil
	}
	defer df.Close()

	var c Config
	err = json.NewDecoder(df).Decode(&c)
	if err != nil {
		return nil, w(err, "Could not parse donate file: %s", err.Error())
	}

	if c.BTC == "" && c.PayPal == "" && c.USDC == "" {
		return nil, ErrNoSupportedPaymentType
	}

	// TODO: Test BTC address validity? Or can we rely on Coinbase (in which case, remove below)

	if c.USDC != "" && !ValidUSDCAddress(c.USDC) {
		return nil, w(fmt.Errorf("Invalid destination address in donate file: \"%s\"", c.USDC), "")
	}

	return &c, nil
}

// ValidUSDCAddress returns true if an string is a correctly formated USDC address.
// It doesn't check whether the address actually exists.
func ValidUSDCAddress(addr string) bool {
	return usdcAddressPattern.MatchString(addr)
}
