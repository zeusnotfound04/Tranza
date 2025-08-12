package razorpay

import (
	"fmt"
)

// Payout represents Razorpay payout
type Payout struct {
	ID          string                 `json:"id"`
	Entity      string                 `json:"entity"`
	Amount      int64                  `json:"amount"`
	Currency    string                 `json:"currency"`
	Status      string                 `json:"status"`
	Purpose     string                 `json:"purpose"`
	UTR         string                 `json:"utr"`
	Mode        string                 `json:"mode"`
	Reference   string                 `json:"reference_id"`
	Narration   string                 `json:"narration"`
	Notes       map[string]interface{} `json:"notes"`
	CreatedAt   int64                  `json:"created_at"`
}

// PayoutRequest represents payout creation request
type PayoutRequest struct {
	AccountNumber string                 `json:"account_number"`
	Amount        int64                  `json:"amount"`
	Currency      string                 `json:"currency"`
	Mode          string                 `json:"mode"`
	Purpose       string                 `json:"purpose"`
	FundAccount   FundAccount            `json:"fund_account"`
	QueueIfLowBalance bool               `json:"queue_if_low_balance"`
	Reference     string                 `json:"reference_id"`
	Narration     string                 `json:"narration"`
	Notes         map[string]interface{} `json:"notes,omitempty"`
}

// FundAccount represents fund account for payout
type FundAccount struct {
	ID      string  `json:"id,omitempty"`
	Entity  string  `json:"entity"`
	Account Account `json:"account"`
}

// Account represents account details
type Account struct {
	Name   string `json:"name"`
	IFSC   string `json:"ifsc"`
	Number string `json:"number"`
}

// CreatePayout creates a new payout
func (c *Client) CreatePayout(req *PayoutRequest) (*Payout, error) {
	url := fmt.Sprintf("%s/payouts", c.BaseURL)
	
	var payout Payout
	if err := c.makeRequest("POST", url, req, &payout); err != nil {
		return nil, fmt.Errorf("failed to create payout: %w", err)
	}

	return &payout, nil
}

// GetPayout retrieves a payout by ID
func (c *Client) GetPayout(payoutID string) (*Payout, error) {
	url := fmt.Sprintf("%s/payouts/%s", c.BaseURL, payoutID)
	
	var payout Payout
	if err := c.makeRequest("GET", url, nil, &payout); err != nil {
		return nil, fmt.Errorf("failed to get payout: %w", err)
	}

	return &payout, nil
}

