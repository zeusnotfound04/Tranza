package razorpay

import (
	"fmt"
)

// Payout represents Razorpay payout
type Payout struct {
	ID            string        `json:"id"`
	Entity        string        `json:"entity"`
	Amount        int64         `json:"amount"`
	Currency      string        `json:"currency"`
	Status        string        `json:"status"`
	Purpose       string        `json:"purpose"`
	UTR           string        `json:"utr"`
	Mode          string        `json:"mode"`
	Reference     string        `json:"reference_id"`
	Narration     string        `json:"narration"`
	Notes         FlexibleNotes `json:"notes"` // Handles both map and array from Razorpay
	Fees          int64         `json:"fees"`
	Tax           int64         `json:"tax"`
	FailureReason string        `json:"failure_reason,omitempty"`
	CreatedAt     int64         `json:"created_at"`
	ProcessedAt   int64         `json:"processed_at,omitempty"`
	FundAccountID string        `json:"fund_account_id"`
}

// PayoutRequest represents payout creation request
type PayoutRequest struct {
	AccountNumber     string            `json:"account_number"`
	Amount            int64             `json:"amount"`
	Currency          string            `json:"currency"`
	Mode              string            `json:"mode"`
	Purpose           string            `json:"purpose"`
	FundAccount       PayoutFundAccount `json:"fund_account"`
	QueueIfLowBalance bool              `json:"queue_if_low_balance,omitempty"`
	Reference         string            `json:"reference_id,omitempty"`
	Narration         string            `json:"narration,omitempty"`
	Notes             FlexibleNotes     `json:"notes,omitempty"` // Handles both map and array from Razorpay
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
	Name        string        `json:"name"`
	Email       string        `json:"email,omitempty"`
	Contact     string        `json:"contact,omitempty"` // Phone number
	Type        string        `json:"type"`              // customer, employee, vendor
	ReferenceID string        `json:"reference_id,omitempty"`
	Notes       FlexibleNotes `json:"notes,omitempty"` // Handles both map and array from Razorpay
}

// ContactResponse represents contact creation response
type ContactResponse struct {
	ID          string        `json:"id"`
	Entity      string        `json:"entity"`
	Name        string        `json:"name"`
	Contact     string        `json:"contact"`
	Email       string        `json:"email"`
	Type        string        `json:"type"`
	ReferenceID string        `json:"reference_id"`
	BatchID     string        `json:"batch_id"`
	Active      bool          `json:"active"`
	Notes       FlexibleNotes `json:"notes"` // Handles both map and array from Razorpay
	CreatedAt   int64         `json:"created_at"`
}

// FundAccountRequest represents fund account creation request
type FundAccountRequest struct {
	ContactID   string             `json:"contact_id"`
	AccountType string             `json:"account_type"`
	VPA         *PayoutVPA         `json:"vpa,omitempty"`
	BankAccount *PayoutBankAccount `json:"bank_account,omitempty"`
}

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
	fmt.Printf("🚀 Starting UPI Payout: UPI=%s, Amount=%d, Contact=%s\n", upiID, amount, contactName)

	// Create contact first
	contact := &PayoutContact{
		Name:        contactName,
		Contact:     phone,
		Type:        ContactTypeCustomer,
		ReferenceID: referenceID + "_contact",
	}

	fmt.Printf("📞 Creating contact for: %s\n", contactName)
	contactResp, err := c.CreateContact(contact)
	if err != nil {
		return nil, fmt.Errorf("failed to create contact: %w", err)
	}
	fmt.Printf("✅ Contact created with ID: %s\n", contactResp.ID)

	// Create fund account
	fundAccount := &FundAccountRequest{
		ContactID:   contactResp.ID,
		AccountType: AccountTypeVPA,
		VPA: &PayoutVPA{
			Address: upiID,
		},
	}

	fmt.Printf("💰 Creating fund account for UPI: %s\n", upiID)
	fundAccountResp, err := c.CreateFundAccount(fundAccount)
	if err != nil {
		return nil, fmt.Errorf("failed to create fund account: %w", err)
	}
	fmt.Printf("✅ Fund account created with ID: %s\n", fundAccountResp.ID)

	// Create payout using fund_account_id (correct Razorpay API approach)
	payoutData := map[string]interface{}{
		"fund_account_id":      fundAccountResp.ID,
		"amount":               amount,
		"currency":             currency,
		"mode":                 PayoutModeUPI,
		"purpose":              purpose,
		"queue_if_low_balance": true,
		"reference_id":         referenceID,
		"narration":            narration,
	}

	// Only add account number if available (optional for most accounts)
	if c.AccountNumber != "" {
		payoutData["account_number"] = c.AccountNumber
		fmt.Printf("💳 Using account number: %s\n", c.AccountNumber)
	} else {
		fmt.Printf("💳 No account number set - proceeding without it\n")
	}

	fmt.Printf("💸 Creating payout with fund_account_id: %s\n", fundAccountResp.ID)
	url := fmt.Sprintf("%s/payouts", c.BaseURL)
	fmt.Printf("🌐 Payout URL: %s\n", url)

	var payout Payout
	if err := c.makeRequest("POST", url, payoutData, &payout); err != nil {
		return nil, fmt.Errorf("failed to create payout: %w", err)
	}

	fmt.Printf("🎉 Payout created successfully with ID: %s\n", payout.ID)
	return &payout, nil
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
