package dto

import (
	"time"

	"github.com/shopspring/decimal"
)

// Pagination represents pagination metadata (imported from utils)
type Pagination struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
	HasNext    bool  `json:"has_next"`
	HasPrev    bool  `json:"has_prev"`
}

// CreateExternalTransferRequest represents the request to create an external transfer
type CreateExternalTransferRequest struct {
	Amount         decimal.Decimal `json:"amount" binding:"required,gt=0" example:"100.00"`
	Currency       string          `json:"currency,omitempty" example:"INR"`
	Description    string          `json:"description,omitempty" binding:"max=500" example:"Payment for services"`
	RecipientType  string          `json:"recipient_type" binding:"required,oneof=upi phone" example:"upi"`
	RecipientValue string          `json:"recipient_value" binding:"required" example:"user@paytm"`
	RecipientName  string          `json:"recipient_name,omitempty" binding:"max=100" example:"John Doe"`
}

// ExternalTransferResponse represents the response after creating an external transfer
type ExternalTransferResponse struct {
	ID             string          `json:"id"`
	ReferenceID    string          `json:"reference_id"`
	Amount         decimal.Decimal `json:"amount"`
	Currency       string          `json:"currency"`
	TransferFee    decimal.Decimal `json:"transfer_fee"`
	TotalAmount    decimal.Decimal `json:"total_amount"`
	RecipientType  string          `json:"recipient_type"`
	RecipientValue string          `json:"recipient_value"`
	RecipientName  string          `json:"recipient_name,omitempty"`
	Status         string          `json:"status"`
	CreatedAt      time.Time       `json:"created_at"`
	EstimatedTime  string          `json:"estimated_time,omitempty"`
}

// ExternalTransferStatusResponse represents transfer status information
type ExternalTransferStatusResponse struct {
	ID               string          `json:"id"`
	ReferenceID      string          `json:"reference_id"`
	Status           string          `json:"status"`
	DisplayStatus    string          `json:"display_status"`
	Amount           decimal.Decimal `json:"amount"`
	TransferFee      decimal.Decimal `json:"transfer_fee"`
	TotalAmount      decimal.Decimal `json:"total_amount"`
	RecipientDisplay string          `json:"recipient_display"`
	ProcessedAt      *time.Time      `json:"processed_at,omitempty"`
	CompletedAt      *time.Time      `json:"completed_at,omitempty"`
	FailureReason    string          `json:"failure_reason,omitempty"`
	CanRetry         bool            `json:"can_retry"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
}

// ExternalTransferHistoryRequest represents request parameters for transfer history
type ExternalTransferHistoryRequest struct {
	Page     int    `form:"page,default=1" binding:"min=1"`
	Limit    int    `form:"limit,default=20" binding:"min=1,max=100"`
	Status   string `form:"status,omitempty" binding:"omitempty,oneof=pending processing success failed cancelled"`
	DateFrom string `form:"date_from,omitempty" binding:"omitempty,datetime=2006-01-02"`
	DateTo   string `form:"date_to,omitempty" binding:"omitempty,datetime=2006-01-02"`
}

// ExternalTransferHistoryResponse represents the transfer history response
type ExternalTransferHistoryResponse struct {
	Transfers  []ExternalTransferStatusResponse `json:"transfers"`
	Pagination Pagination                       `json:"pagination"`
	Summary    ExternalTransferSummary          `json:"summary"`
}

// ExternalTransferSummary provides summary statistics
type ExternalTransferSummary struct {
	TotalTransfers  int             `json:"total_transfers"`
	TotalAmount     decimal.Decimal `json:"total_amount"`
	SuccessfulCount int             `json:"successful_count"`
	FailedCount     int             `json:"failed_count"`
	PendingCount    int             `json:"pending_count"`
	SuccessRate     float64         `json:"success_rate"`
}

// ValidateTransferRequest represents request to validate transfer before processing
type ValidateTransferRequest struct {
	Amount         decimal.Decimal `json:"amount" binding:"required,gt=0"`
	RecipientType  string          `json:"recipient_type" binding:"required,oneof=upi phone"`
	RecipientValue string          `json:"recipient_value" binding:"required"`
}

// ValidateTransferResponse represents transfer validation response
type ValidateTransferResponse struct {
	Valid         bool            `json:"valid"`
	TransferFee   decimal.Decimal `json:"transfer_fee"`
	TotalAmount   decimal.Decimal `json:"total_amount"`
	EstimatedTime string          `json:"estimated_time"`
	RecipientName string          `json:"recipient_name,omitempty"`
	Warnings      []string        `json:"warnings,omitempty"`
	Errors        []string        `json:"errors,omitempty"`
}

// CancelTransferRequest represents request to cancel a transfer
type CancelTransferRequest struct {
	TransferID string `json:"transfer_id" binding:"required"`
	Reason     string `json:"reason,omitempty" binding:"max=500"`
}

// CancelTransferResponse represents cancel transfer response
type CancelTransferResponse struct {
	Success  bool   `json:"success"`
	Message  string `json:"message"`
	RefundID string `json:"refund_id,omitempty"`
}

// RetryTransferRequest represents request to retry a failed transfer
type RetryTransferRequest struct {
	TransferID string `json:"transfer_id" binding:"required"`
}

// RetryTransferResponse represents retry transfer response
type RetryTransferResponse struct {
	Success    bool   `json:"success"`
	Message    string `json:"message"`
	NewStatus  string `json:"new_status"`
	RetryCount int    `json:"retry_count"`
}

// BotTransferRequest represents transfer request from Slack bot
type BotTransferRequest struct {
	Amount         decimal.Decimal `json:"amount" binding:"required,gt=0"`
	RecipientType  string          `json:"recipient_type" binding:"required,oneof=upi phone"`
	RecipientValue string          `json:"recipient_value" binding:"required"`
	Description    string          `json:"description,omitempty"`
	SlackUserID    string          `json:"slack_user_id,omitempty"`
	SlackChannelID string          `json:"slack_channel_id,omitempty"`
}

// BotTransferResponse represents transfer response for Slack bot
type BotTransferResponse struct {
	Success       bool            `json:"success"`
	TransferID    string          `json:"transfer_id"`
	ReferenceID   string          `json:"reference_id"`
	Amount        decimal.Decimal `json:"amount"`
	TotalAmount   decimal.Decimal `json:"total_amount"`
	TransferFee   decimal.Decimal `json:"transfer_fee"`
	Status        string          `json:"status"`
	Message       string          `json:"message"`
	Recipient     string          `json:"recipient"`
	EstimatedTime string          `json:"estimated_time,omitempty"`
}

// TransferFeesResponse represents transfer fees information
type TransferFeesResponse struct {
	UPIFee       decimal.Decimal `json:"upi_fee"`
	PhoneFee     decimal.Decimal `json:"phone_fee"`
	MinAmount    decimal.Decimal `json:"min_amount"`
	MaxAmount    decimal.Decimal `json:"max_amount"`
	DailyLimit   decimal.Decimal `json:"daily_limit"`
	MonthlyLimit decimal.Decimal `json:"monthly_limit"`
	FeeStructure []FeeRange      `json:"fee_structure"`
}

// FeeRange represents fee structure for different amount ranges
type FeeRange struct {
	MinAmount decimal.Decimal `json:"min_amount"`
	MaxAmount decimal.Decimal `json:"max_amount"`
	Fee       decimal.Decimal `json:"fee"`
	FeeType   string          `json:"fee_type"` // "fixed" or "percentage"
}

// PaginatedExternalTransferResponse represents paginated external transfer response
type PaginatedExternalTransferResponse struct {
	Transfers   []ExternalTransferResponse `json:"transfers"`
	TotalCount  int64                      `json:"total_count"`
	CurrentPage int                        `json:"current_page"`
	TotalPages  int64                      `json:"total_pages"`
	HasMore     bool                       `json:"has_more"`
}

// Bot-specific DTOs for Slack integration
type BotValidateTransferRequest struct {
	Amount         decimal.Decimal `json:"amount" binding:"required"`
	RecipientType  string          `json:"recipient_type" binding:"required"`  // "upi" or "phone"
	RecipientValue string          `json:"recipient_value" binding:"required"` // UPI ID or phone number
}

type BotValidateTransferResponse struct {
	Valid         bool            `json:"valid"`
	Errors        []string        `json:"errors,omitempty"`
	Warnings      []string        `json:"warnings,omitempty"`
	TransferFee   decimal.Decimal `json:"transfer_fee"`
	TotalAmount   decimal.Decimal `json:"total_amount"`
	EstimatedTime string          `json:"estimated_time"`
}

type BotCreateTransferRequest struct {
	Amount         decimal.Decimal `json:"amount" binding:"required"`
	RecipientType  string          `json:"recipient_type" binding:"required"`  // "upi" or "phone"
	RecipientValue string          `json:"recipient_value" binding:"required"` // UPI ID or phone number
	RecipientName  string          `json:"recipient_name,omitempty"`
	Description    string          `json:"description,omitempty"`
}

type BotTransferStatusResponse struct {
	TransferID    string          `json:"transfer_id"`
	ReferenceID   string          `json:"reference_id"`
	Status        string          `json:"status"`
	Amount        decimal.Decimal `json:"amount"`
	Recipient     string          `json:"recipient"`
	EstimatedTime string          `json:"estimated_time,omitempty"`
}
