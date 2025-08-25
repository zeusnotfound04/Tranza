package razorpay

import (
	"fmt"
)

// Payout represents Razorpay payout
type Payout struct {
	ID            string                 `json:"id"`
	Entity        string                 `json:"entity"`
	Amount        int64                  `json:"amount"`
	Currency      string                 `json:"currency"`
	Status        string                 `json:"status"`
	Purpose       string                 `json:"purpose"`
	UTR           string                 `json:"utr"`
	Mode          string                 `json:"mode"`
	Reference     string                 `json:"reference_id"`
	Narration     string                 `json:"narration"`
	Notes         map[string]interface{} `json:"notes"`
	Fees          int64                  `json:"fees"`
	Tax           int64                  `json:"tax"`
	FailureReason string                 `json:"failure_reason,omitempty"`
	CreatedAt     int64                  `json:"created_at"`
	ProcessedAt   int64                  `json:"processed_at,omitempty"`
	FundAccountID string                 `json:"fund_account_id"`
}

// PayoutRequest represents payout creation request
type PayoutRequest struct {
	AccountNumber     string                 `json:"account_number"`
	Amount            int64                  `json:"amount"`
	Currency          string                 `json:"currency"`
	Mode              string                 `json:"mode"`
	Purpose           string                 `json:"purpose"`
	FundAccount       PayoutFundAccount      `json:"fund_account"`
	QueueIfLowBalance bool                   `json:"queue_if_low_balance,omitempty"`
	Reference         string                 `json:"reference_id,omitempty"`
	Narration         string                 `json:"narration,omitempty"`
	Notes             map[string]interface{} `json:"notes,omitempty"`
}

// PayoutFundAccount represents fund account for payout
type PayoutFundAccount struct {
	AccountType string             `json:"account_type"` // vpa, bank_account
	VPA         *PayoutVPA         `json:"vpa,omitempty"`
	BankAccount *PayoutBankAccount `json:"bank_account,omitempty"`
	Contact     PayoutContact      `json:"contact"`
}

// PayoutVPA represents UPI VPA details
type PayoutVPA struct {
	Address string `json:"address"` // UPI ID
}

// PayoutBankAccount represents bank account details
type PayoutBankAccount struct {
	Name          string `json:"name"`
	IFSC          string `json:"ifsc"`
	AccountNumber string `json:"account_number"`
}

// PayoutContact represents contact details
type PayoutContact struct {
	Name        string                 `json:"name"`
	Email       string                 `json:"email,omitempty"`
	Contact     string                 `json:"contact,omitempty"` // Phone number
	Type        string                 `json:"type"`              // customer, employee, vendor
	ReferenceID string                 `json:"reference_id,omitempty"`
	Notes       map[string]interface{} `json:"notes,omitempty"`
}

// ContactResponse represents contact creation response
type ContactResponse struct {
	ID          string                 `json:"id"`
	Entity      string                 `json:"entity"`
	Name        string                 `json:"name"`
	Contact     string                 `json:"contact"`
	Email       string                 `json:"email"`
	Type        string                 `json:"type"`
	ReferenceID string                 `json:"reference_id"`
	BatchID     string                 `json:"batch_id"`
	Active      bool                   `json:"active"`
	Notes       map[string]interface{} `json:"notes"`
	CreatedAt   int64                  `json:"created_at"`
}

// FundAccountRequest represents fund account creation request
type FundAccountRequest struct {
	ContactID   string             `json:"contact_id"`
	AccountType string             `json:"account_type"`
	VPA         *PayoutVPA         `json:"vpa,omitempty"`
	BankAccount *PayoutBankAccount `json:"bank_account,omitempty"`
}

// FundAccountResponse represents fund account creation response
type FundAccountResponse struct {
	ID          string             `json:"id"`
	Entity      string             `json:"entity"`
	ContactID   string             `json:"contact_id"`
	AccountType string             `json:"account_type"`
	VPA         *PayoutVPA         `json:"vpa,omitempty"`
	BankAccount *PayoutBankAccount `json:"bank_account,omitempty"`
	Active      bool               `json:"active"`
	CreatedAt   int64              `json:"created_at"`
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

// CreateContact creates a contact for payouts
func (c *Client) CreateContact(contact *PayoutContact) (*ContactResponse, error) {
	url := fmt.Sprintf("%s/contacts", c.BaseURL)

	var result ContactResponse
	if err := c.makeRequest("POST", url, contact, &result); err != nil {
		return nil, fmt.Errorf("failed to create contact: %w", err)
	}

	return &result, nil
}

// CreateFundAccount creates a fund account
func (c *Client) CreateFundAccount(fundAccount *FundAccountRequest) (*FundAccountResponse, error) {
	url := fmt.Sprintf("%s/fund_accounts", c.BaseURL)

	var result FundAccountResponse
	if err := c.makeRequest("POST", url, fundAccount, &result); err != nil {
		return nil, fmt.Errorf("failed to create fund account: %w", err)
	}

	return &result, nil
}

// CreateUPIPayout creates a UPI payout (convenience method)
func (c *Client) CreateUPIPayout(upiID string, amount int64, currency, purpose, narration, contactName, phone, referenceID string) (*Payout, error) {
	// Create contact first
	contact := &PayoutContact{
		Name:        contactName,
		Contact:     phone,
		Type:        ContactTypeCustomer,
		ReferenceID: referenceID + "_contact",
	}

	contactResp, err := c.CreateContact(contact)
	if err != nil {
		return nil, fmt.Errorf("failed to create contact: %w", err)
	}

	// Create fund account
	fundAccount := &FundAccountRequest{
		ContactID:   contactResp.ID,
		AccountType: AccountTypeVPA,
		VPA: &PayoutVPA{
			Address: upiID,
		},
	}

	_, err = c.CreateFundAccount(fundAccount)
	if err != nil {
		return nil, fmt.Errorf("failed to create fund account: %w", err)
	}

	// Create payout request
	payoutReq := &PayoutRequest{
		Amount:    amount,
		Currency:  currency,
		Mode:      PayoutModeUPI,
		Purpose:   purpose,
		Reference: referenceID,
		Narration: narration,
		FundAccount: PayoutFundAccount{
			AccountType: AccountTypeVPA,
			VPA: &PayoutVPA{
				Address: upiID,
			},
			Contact: *contact,
		},
		QueueIfLowBalance: true,
	}

	return c.CreatePayout(payoutReq)
}

// Payout Status Constants
const (
	PayoutStatusQueued     = "queued"
	PayoutStatusPending    = "pending"
	PayoutStatusProcessing = "processing"
	PayoutStatusProcessed  = "processed"
	PayoutStatusReversed   = "reversed"
	PayoutStatusCancelled  = "cancelled"
	PayoutStatusFailed     = "failed"
)

// Payout Mode Constants
const (
	PayoutModeUPI  = "UPI"
	PayoutModeIMPS = "IMPS"
	PayoutModeNEFT = "NEFT"
	PayoutModeRTGS = "RTGS"
)

// Contact Type Constants
const (
	ContactTypeCustomer = "customer"
	ContactTypeEmployee = "employee"
	ContactTypeVendor   = "vendor"
)

// Account Type Constants
const (
	AccountTypeVPA  = "vpa"
	AccountTypeBank = "bank_account"
)

// Purpose Constants
const (
	PurposeRefund   = "refund"
	PurposeCashback = "cashback"
	PurposePayout   = "payout"
	PurposeSalary   = "salary"
	PurposeUtility  = "utility_bill"
)
