package paypal

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

func url(sandbox bool, path string) string {
	var ss string
	if sandbox {
		ss = "sandbox."
	}
	return fmt.Sprintf("https://api.%spaypal.com/v1/%s", ss, path)
}

func getToken(clientID, secret string, sandbox bool) (string, error) {
	//clientid := getEnv("PAYPAL_CLIENT_ID")
	//secret := getEnv("PAYPAL_SECRET")

	url := url(sandbox, "oauth2/token")
	req, err := http.NewRequest("POST", url, strings.NewReader("grant_type=client_credentials"))
	if err != nil {
		return "", err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Accept-Language", "en_US")
	req.Header.Add("Accept-Language", "en_US")
	req.SetBasicAuth(clientID, secret)
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Response: %d: %s", resp.StatusCode, resp.Status)
	}

	var s struct {
		Token string `json:"access_token"`
	}
	err = json.NewDecoder(resp.Body).Decode(&s)
	if err != nil {
		return "", err
	}

	if s.Token == "" {
		return "", errors.New("No access token in response")
	}

	return s.Token, nil
}

// Send sends money via PayPal.
func Send(clientID, secret string, amt float64, addr string, sandbox bool) (batchID, status string, err error) {
	token, err := getToken(clientID, secret, sandbox)
	if err != nil {
		return "", "", err
	}

	batch := payoutBatch{
		SenderBatchHeader: payoutHeader{
			SenderBatchID: fmt.Sprintf("%d", time.Now().UnixNano()),
			EmailSubject:  "subject",
			EmailMessage:  "message",
		},
		Items: []payoutItem{
			{
				RecipientType: "EMAIL",
				Amount: payoutAmount{
					Currency: "USD",
					Value:    fmt.Sprintf("%.2f", amt),
				},
				Note:         "note",
				SenderItemID: fmt.Sprintf("%d", time.Now().UnixNano()),
				Receiver:     addr,
			},
		},
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
