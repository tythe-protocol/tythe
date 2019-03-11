// Package coinbase provides utilities for working with Coinbase.
package coinbase

import (
	"fmt"
	"strconv"

	gdax "github.com/preichenberger/go-gdax"
	"github.com/tythe-protocol/tythe/utl/env"
)

// SendResult represents the results of sending one individual transaction on Coinbase.
type SendResult struct {
	TransactionID string
	Error         error
}

func (sr SendResult) String() string {
	if sr.Error != nil {
		return fmt.Sprintf("Error: %s", sr.Error.Error())
	} else {
		return fmt.Sprintf("Success - transaction ID: %s", sr.TransactionID)
	}
}

// An Amount to send with Send.
type Amount struct {
	// Currency is the type of currency to send.
	Currency string

	// Value is always denominated in USD(C). If Currency is other, then it
	// will be converted before send.
	Value float64
}

// Send sends money via Coinbase.
func Send(txs map[string]Amount, sandbox bool) (map[string]SendResult, error) {
	key := env.Must("TYTHE_COINBASE_API_KEY")
	secret := env.Must("TYTHE_COINBASE_API_SECRET")
	passphrase := env.Must("TYTHE_COINBASE_API_PASSPHRASE")

	client := gdax.NewClient(secret, key, passphrase)
	if sandbox {
		client.BaseURL = "https://api-public.sandbox.pro.coinbase.com"
	}

	btcPrice, err := btcUSD(client)
	if err != nil {
		return nil, err
	}

	ret := map[string]SendResult{}
	for addr, amt := range txs {
		v := amt.Value
		prec := 2
		if amt.Currency == "BTC" {
			v /= btcPrice
			prec = 8
		}

		params := map[string]string{
			"amount":         fmt.Sprintf("%."+strconv.Itoa(prec)+"f", v),
			"currency":       amt.Currency,
			"crypto_address": addr,
		}

		sr := SendResult{}
		var resp struct {
			ID string `json:"id"`
		}
		_, err := client.Request("POST", "/withdrawals/crypto", params, &resp)
		if err != nil {
			sr.Error = err
		} else {
			sr.TransactionID = resp.ID
		}

		ret[addr] = sr
	}

	return ret, nil
}

func btcUSD(c *gdax.Client) (float64, error) {
	t, err := c.GetTicker("BTC-USD")
	if err != nil {
		return 0, err
	}
	return strconv.ParseFloat(t.Price, 64)
}
