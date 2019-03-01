package paypal

type payoutBatch struct {
	SenderBatchHeader payoutHeader `json:"sender_batch_header"`
	Items             []payoutItem `json:"items"`
}

type payoutHeader struct {
	SenderBatchID string `json:"sender_batch_id"`
	EmailSubject  string `json:"email_subject"`
	EmailMessage  string `json:"email_message"`
}

type payoutItem struct {
	RecipientType string       `json:"recipient_type"`
	Amount        payoutAmount `json:"amount"`
	Note          string       `json:"note"`
	SenderItemID  string       `json:"sender_item_id"`
	Receiver      string       `json:"receiver"`
}

type payoutAmount struct {
	Value    string `json:"value"`
	Currency string `json:"currency"`
}

type payoutBatchResponse struct {
	BatchHeader struct {
		PayoutBatchID string `json:"payout_batch_id"`
		BatchStatus   string `json:"batch_status"`
	} `json:"batch_header"`
}
