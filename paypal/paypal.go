package paypal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/tythe-protocol/go-tythe/env"

	"github.com/pkg/errors"
)

func url(sandbox bool, path string) string {
	var ss string
	if sandbox {
		ss = "sandbox."
	}
	return fmt.Sprintf("https://api.%spaypal.com/v1/%s", ss, path)
}

func getToken(sandbox bool) (string, error) {
	clientID := env.Must("PAYPAL_CLIENT_ID")
	secret := env.Must("PAYPAL_SECRET")

	w := func(err error) error {
		return errors.Wrap(err, "Failed to get auth token")
	}

	url := url(sandbox, "oauth2/token")
	req, err := http.NewRequest("POST", url, strings.NewReader("grant_type=client_credentials"))
	if err != nil {
		return "", w(err)
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Accept-Language", "en_US")
	req.Header.Add("Accept-Language", "en_US")
	req.SetBasicAuth(clientID, secret)
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return "", w(err)
	}
	if resp.StatusCode != http.StatusOK {
		buf := &bytes.Buffer{}
		io.Copy(buf, resp.Body)
		resp.Body.Close()
		return "", w(fmt.Errorf("Response: %d %s - %s", resp.StatusCode, resp.Status, string(buf.Bytes())))
	}

	var s struct {
		Token string `json:"access_token"`
	}
	err = json.NewDecoder(resp.Body).Decode(&s)
	if err != nil {
		return "", w(err)
	}

	if s.Token == "" {
		return "", w(errors.New("No access token in response"))
	}

	return s.Token, nil
}

// Send sends money via PayPal.
func Send(txs map[string]float64, sandbox bool) (batchID, status string, err error) {
	// Remove when out of testing
	sandbox = true

	token, err := getToken(sandbox)
	if err != nil {
		return "", "", err
	}

	batch := payoutBatch{
		SenderBatchHeader: payoutHeader{
			SenderBatchID: fmt.Sprintf("%d", time.Now().UnixNano()),
			EmailSubject:  "subject",
			EmailMessage:  "message",
		},
	}

	for addr, amt := range txs {
		batch.Items = append(batch.Items, payoutItem{
			RecipientType: "EMAIL",
			Amount: payoutAmount{
				Currency: "USD",
				Value:    fmt.Sprintf("%.2f", amt),
			},
			Note:         "note",
			SenderItemID: fmt.Sprintf("%s.%d", addr, time.Now().UnixNano()),
			Receiver:     addr,
		})
	}

	data, err := json.Marshal(batch)
	if err != nil {
		return "", "", err
	}

	req, err := http.NewRequest("POST", url(sandbox, "payments/payouts"), bytes.NewReader(data))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return "", "", err
	}

	if resp.StatusCode != http.StatusCreated {
		buf := &bytes.Buffer{}
		buf.Write([]byte(fmt.Sprintf("Error %d: %s\n", resp.StatusCode, resp.Status)))
		io.Copy(buf, resp.Body)
		buf.Write([]byte("\nRequest:"))
		buf.Write(data)
		return "", "", errors.New(string(buf.Bytes()))
	}

	var pr payoutBatchResponse
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&pr)
	if err != nil {
		return "", "", err
	}

	return pr.BatchHeader.PayoutBatchID, pr.BatchHeader.BatchStatus, nil
}
