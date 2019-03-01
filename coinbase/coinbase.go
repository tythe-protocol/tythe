package coinbase

import (
	"fmt"

	gdax "github.com/preichenberger/go-gdax"
	"github.com/tythe-protocol/go-tythe/env"
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

// Send sends money via Coinbase.
func Send(txs map[string]float64, sandbox bool) map[string]SendResult {
	key := env.Must("TYTHE_COINBASE_API_KEY")
	secret := env.Must("TYTHE_COINBASE_API_SECRET")
	passphrase := env.Must("TYTHE_COINBASE_API_PASSPHRASE")

	client := gdax.NewClient(secret, key, passphrase)
	if sandbox {
		client.BaseURL = "https://api-public.sandbox.pro.coinbase.com"
	}

	ret := map[string]SendResult{}
	for addr, amt := range txs {
		params := map[string]interface{}{
			"amount":         amt,
			"currency":       "USDC",
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

	return ret
}
