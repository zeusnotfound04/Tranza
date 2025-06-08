package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/zeusnotfound04/Tranza/config"
	"github.com/zeusnotfound04/Tranza/models"
	"github.com/zeusnotfound04/Tranza/utils"
)

type RazorpayService struct {
	KeyID     string
	KeySecret string
	BaseURL   string
	Client    *http.Client
}

func NewRazorpayService() *RazorpayService {
	cfg := config.LoadConfig()
	
	baseURL := "https://api.razorpay.com/v1"
	if cfg.Environment == "test" {
		// Razorpay uses the same URL for both test and live, 
		// the environment is determined by the API keys
		baseURL = "https://api.razorpay.com/v1"
	}
	
	return &RazorpayService{
		KeyID:     cfg.RazorpayKeyID,
		KeySecret: cfg.RazorpayKeySecret,
		BaseURL:   baseURL,
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (r *RazorpayService) CreateOrder(ctx context.Context, amount float64, currency, receipt string, notes map[string]string) (*models.RazorpayOrderResponse, error) {
	if amount <= 0 {
		return nil, fmt.Errorf("amount must be positive")
	}
	
	if currency == "" {
		currency = "INR"
	}
	
	payload := models.RazorpayOrderRequest{
		Amount:   int(amount * 100), // Convert to paise
		Currency: currency,
		Receipt:  receipt,
		Notes:    notes,
	}
	
	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	
	req, err := http.NewRequestWithContext(ctx, "POST", r.BaseURL+"/orders", bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	req.SetBasicAuth(r.KeyID, r.KeySecret)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Razorpay-Go/1.0")
	
	resp, err := r.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	
	if resp.StatusCode >= 400 {
		var razorpayErr models.RazorpayError
		if err := json.Unmarshal(body, &razorpayErr); err == nil {
			return nil, fmt.Errorf("razorpay error: %s - %s", razorpayErr.Error.Code, razorpayErr.Error.Description)
		}
		return nil, fmt.Errorf("razorpay error: %s", string(body))
	}
	
	var orderResp models.RazorpayOrderResponse
	if err := json.Unmarshal(body, &orderResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	
	return &orderResp, nil
}

func (r *RazorpayService) FetchOrder(ctx context.Context, orderID string) (*models.RazorpayOrderResponse, error) {
	if orderID == "" {
		return nil, fmt.Errorf("order ID cannot be empty")
	}
	
	req, err := http.NewRequestWithContext(ctx, "GET", r.BaseURL+"/orders/"+orderID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	req.SetBasicAuth(r.KeyID, r.KeySecret)
	req.Header.Set("User-Agent", "Razorpay-Go/1.0")
	
	resp, err := r.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	
	if resp.StatusCode >= 400 {
		var razorpayErr models.RazorpayError
		if err := json.Unmarshal(body, &razorpayErr); err == nil {
			return nil, fmt.Errorf("razorpay error: %s - %s", razorpayErr.Error.Code, razorpayErr.Error.Description)
		}
		return nil, fmt.Errorf("razorpay error: %s", string(body))
	}
	
	var orderResp models.RazorpayOrderResponse
	if err := json.Unmarshal(body, &orderResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	
	return &orderResp, nil
}

func (r *RazorpayService) FetchPayment(ctx context.Context, paymentID string) (*models.RazorpayPaymentResponse, error) {
	if paymentID == "" {
		return nil, fmt.Errorf("payment ID cannot be empty")
	}
	
	req, err := http.NewRequestWithContext(ctx, "GET", r.BaseURL+"/payments/"+paymentID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	req.SetBasicAuth(r.KeyID, r.KeySecret)
	req.Header.Set("User-Agent", "Razorpay-Go/1.0")
	
	resp, err := r.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	
	if resp.StatusCode >= 400 {
		var razorpayErr models.RazorpayError
		if err := json.Unmarshal(body, &razorpayErr); err == nil {
			return nil, fmt.Errorf("razorpay error: %s - %s", razorpayErr.Error.Code, razorpayErr.Error.Description)
		}
		return nil, fmt.Errorf("razorpay error: %s", string(body))
	}
	
	var paymentResp models.RazorpayPaymentResponse
	if err := json.Unmarshal(body, &paymentResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	
	return &paymentResp, nil
}

func (r *RazorpayService) VerifyPaymentSignature(orderID, paymentID, signature string) error {
	if orderID == "" || paymentID == "" || signature == "" {
		return fmt.Errorf("order ID, payment ID, and signature are required")
	}
	
	if !utils.VerifyRazorpaySignature(orderID, paymentID, signature, r.KeySecret) {
		return fmt.Errorf("signature verification failed")
	}
	
	return nil
}

func (r *RazorpayService) VerifyWebhookSignature(body, signature string) error {
	cfg := config.LoadConfig()
	if cfg.WebhookSecret == "" {
		return fmt.Errorf("webhook secret not configured")
	}
	
	if !utils.VerifyWebhookSignature(body, signature, cfg.WebhookSecret) {
		return fmt.Errorf("webhook signature verification failed")
	}
	
	return nil
}