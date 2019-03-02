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
	// DonateFile is the name of the dot-donate file (see https://github.com/aboodman/dot-donate).
	DonateFile string = ".donate"

	PaymentTypeBTC    string = "BTC"
	PaymentTypePayPal string = "PayPal"
	PaymentTypeUSDC   string = "USDC"
	PaymentTypeNone   string = ""
)

var (
	// ErrNoSupportedPaymentType is returned by conf.Read() when the config looks valid, but it
	// doesn't contain a payment type tythe supports.
	ErrNoSupportedPaymentType error = errors.New("No supported payment type")

	usdcAddressPattern = regexp.MustCompile("^0x[0-9a-f]{40}$")
)

// Config describes the json metadata developers add to their package to opt-in
// to receiving tythes.
type Config struct {
	BTC    string `json:"BTC,omitempty"`
	PayPal string `json:"PayPal,omitempty"`
	USDC   string `json:"USDC,omitempty"`
}

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
	w := func(err error) error {
		return errors.Wrapf(err, "Could not read config for package: %s:", dir)
	}

	f, err := os.Open(path.Join(dir, DonateFile))
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
		return nil, w(errors.Wrap(err, "donate file is not valid JSON"))
	}

	if c.BTC == "" && c.PayPal == "" && c.USDC == "" {
		return nil, ErrNoSupportedPaymentType
	}

	// TODO: Test BTC address validity? Or can we rely on Coinbase (in which case, remove below)

	if c.USDC != "" && !ValidUSDCAddress(c.USDC) {
		return nil, fmt.Errorf("invalid destination address in donate file: \"%s\"", c.USDC)
	}

	return &c, nil
}

// ValidUSDCAddress returns true if an string is a correctly formated USDC address.
// It doesn't check whether the address actually exists.
func ValidUSDCAddress(addr string) bool {
	return usdcAddressPattern.MatchString(addr)
}
