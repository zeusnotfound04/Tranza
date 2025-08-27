// pkg/razorpay/client.go
package razorpay

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// FlexibleNotes handles Razorpay's inconsistent notes field (can be object or array)
type FlexibleNotes map[string]interface{}

func (fn *FlexibleNotes) UnmarshalJSON(data []byte) error {
	// Try to unmarshal as object first
	var obj map[string]interface{}
	if err := json.Unmarshal(data, &obj); err == nil {
		*fn = FlexibleNotes(obj)
		return nil
	}

	// If it fails, it might be an array (empty notes)
	var arr []interface{}
	if err := json.Unmarshal(data, &arr); err == nil {
		// If it's an empty array, create empty map
		if len(arr) == 0 {
			*fn = make(FlexibleNotes)
			return nil
		}
		// If it's not empty, that's unexpected
		return fmt.Errorf("unexpected non-empty array for notes")
	}

	return fmt.Errorf("notes field is neither object nor array")
}

// Client represents Razorpay client
type Client struct {
	KeyID         string
	KeySecret     string
	BaseURL       string
	HTTPClient    *http.Client
	AccountNumber string // Required for payouts in live mode
}

// NewClient creates a new Razorpay client
func NewClient(keyID, keySecret string) *Client {
	accountNumber := os.Getenv("RAZORPAY_ACCOUNT_NUMBER")
	if accountNumber == "" {
		// For testing or if account number is not set, use empty string
		accountNumber = ""
	}

	return &Client{
		KeyID:         keyID,
		KeySecret:     keySecret,
		BaseURL:       "https://api.razorpay.com/v1",
		AccountNumber: accountNumber,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// NewClientWithAccount creates a new Razorpay client with explicit account number
func NewClientWithAccount(keyID, keySecret, accountNumber string) *Client {
	return &Client{
		KeyID:         keyID,
		KeySecret:     keySecret,
		BaseURL:       "https://api.razorpay.com/v1",
		AccountNumber: accountNumber,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Order represents Razorpay order
type Order struct {
	ID         string        `json:"id"`
	Entity     string        `json:"entity"`
	Amount     int64         `json:"amount"`
	AmountPaid int64         `json:"amount_paid"`
	AmountDue  int64         `json:"amount_due"`
	Currency   string        `json:"currency"`
	Receipt    string        `json:"receipt"`
	Status     string        `json:"status"`
	Attempts   int           `json:"attempts"`
	Notes      FlexibleNotes `json:"notes"`
	CreatedAt  int64         `json:"created_at"`
}

// Payment represents Razorpay payment
type Payment struct {
	ID          string        `json:"id"`
	Entity      string        `json:"entity"`
	Amount      int64         `json:"amount"`
	Currency    string        `json:"currency"`
	Status      string        `json:"status"`
	OrderID     string        `json:"order_id"`
	Method      string        `json:"method"`
	Description string        `json:"description"`
	Email       string        `json:"email"`
	Contact     string        `json:"contact"`
	Notes       FlexibleNotes `json:"notes"`
	CreatedAt   int64         `json:"created_at"`
}

// CreateOrderRequest represents order creation request
type CreateOrderRequest struct {
	Amount   int64                  `json:"amount"`
	Currency string                 `json:"currency"`
	Receipt  string                 `json:"receipt"`
	Notes    map[string]interface{} `json:"notes,omitempty"`
}

// CreateOrder creates a new order
func (c *Client) CreateOrder(amount int64, currency, receipt string) (*Order, error) {
	url := fmt.Sprintf("%s/orders", c.BaseURL)

	req := CreateOrderRequest{
		Amount:   amount,
		Currency: currency,
		Receipt:  receipt,
		Notes:    make(map[string]interface{}), // Provide empty notes map to avoid array response
	}

	var order Order
	if err := c.makeRequest("POST", url, req, &order); err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	return &order, nil
}

// GetOrder retrieves an order by ID
func (c *Client) GetOrder(orderID string) (*Order, error) {
	url := fmt.Sprintf("%s/orders/%s", c.BaseURL, orderID)

	var order Order
	if err := c.makeRequest("GET", url, nil, &order); err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	return &order, nil
}

// GetPayment retrieves a payment by ID
func (c *Client) GetPayment(paymentID string) (*Payment, error) {
	url := fmt.Sprintf("%s/payments/%s", c.BaseURL, paymentID)

	var payment Payment
	if err := c.makeRequest("GET", url, nil, &payment); err != nil {
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}

	return &payment, nil
}

// CapturePayment captures a payment
func (c *Client) CapturePayment(paymentID string, amount int64) (*Payment, error) {
	url := fmt.Sprintf("%s/payments/%s/capture", c.BaseURL, paymentID)

	req := map[string]interface{}{
		"amount": amount,
	}

	var payment Payment
	if err := c.makeRequest("POST", url, req, &payment); err != nil {
		return nil, fmt.Errorf("failed to capture payment: %w", err)
	}

	return &payment, nil
}

// makeRequest makes HTTP request to Razorpay API
func (c *Client) makeRequest(method, url string, body interface{}, result interface{}) error {
	var reqBody io.Reader

	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)

		// Add debug logging for payout-related requests
		if strings.Contains(url, "payouts") || strings.Contains(url, "contacts") || strings.Contains(url, "fund_accounts") {
			fmt.Printf("🔍 DEBUG Razorpay Request: %s %s\nPayload: %s\n", method, url, string(jsonBody))
		}
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(c.KeyID, c.KeySecret)

	// Add X-Razorpay-Account header for payouts if account number is available
	if strings.Contains(url, "payouts") || strings.Contains(url, "contacts") || strings.Contains(url, "fund_accounts") {
		if c.AccountNumber != "" {
			req.Header.Set("X-Razorpay-Account", c.AccountNumber)
			fmt.Printf("🔑 Added X-Razorpay-Account header: %s\n", c.AccountNumber)
		}
	}

	// Make request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	// Add debug logging for payout-related responses
	if strings.Contains(url, "payouts") || strings.Contains(url, "contacts") || strings.Contains(url, "fund_accounts") {
		fmt.Printf("🔍 DEBUG Razorpay Response: %d\nBody: %s\n", resp.StatusCode, string(respBody))
	}

	// Check status code
	if resp.StatusCode >= 400 {
		var errorResp RazorpayError
		if err := json.Unmarshal(respBody, &errorResp); err == nil {
			// Special handling for common payout errors
			if strings.Contains(url, "payouts") && resp.StatusCode == 404 {
				return fmt.Errorf("Razorpay Payouts API not found. This usually means: 1) Payouts not enabled on your account, 2) Wrong API endpoint, or 3) Missing account permissions. Original error: %s", errorResp.Error())
			}
			return &errorResp
		}
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	// Parse response
	if result != nil {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}

	return nil
}

// RazorpayError represents Razorpay API error
type RazorpayError struct {
	RazorpayError struct {
		Code        string `json:"code"`
		Description string `json:"description"`
		Source      string `json:"source"`
		Step        string `json:"step"`
		Reason      string `json:"reason"`
	} `json:"error"`
}

func (e *RazorpayError) Error() string {
	return fmt.Sprintf("Razorpay error [%s]: %s", e.RazorpayError.Code, e.RazorpayError.Description)
}
