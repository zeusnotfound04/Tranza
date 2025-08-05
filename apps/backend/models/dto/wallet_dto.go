package dto

import (
	"github.com/shopspring/decimal"
)

// Load Money Request
type LoadMoneyRequest struct {
	Amount decimal.Decimal `json:"amount" validate:"required,gt=0" binding:"required"`
}

// Load Money Response
type LoadMoneyResponse struct {
	OrderID       string          `json:"order_id"`
	Amount        decimal.Decimal `json:"amount"`
	Currency      string          `json:"currency"`
	TransactionID string          `json:"transaction_id"`
	RazorpayKeyID string          `json:"razorpay_key_id"`
}

// Verify Payment Request
type VerifyPaymentRequest struct {
	RazorpayPaymentID string `json:"razorpay_payment_id" validate:"required" binding:"required"`
	RazorpayOrderID   string `json:"razorpay_order_id" validate:"required" binding:"required"`
	RazorpaySignature string `json:"razorpay_signature" validate:"required" binding:"required"`
}

// Payment Verification Response
type PaymentVerificationResponse struct {
	Success       bool            `json:"success"`
	NewBalance    decimal.Decimal `json:"new_balance"`
	TransactionID string          `json:"transaction_id"`
	Message       string          `json:"message"`
	Amount        decimal.Decimal `json:"amount"`
}

// Wallet Response
type WalletResponse struct {
	ID                    string          `json:"id"`
	Balance               decimal.Decimal `json:"balance"`
	Currency              string          `json:"currency"`
	Status                string          `json:"status"`
	DailyLimit            decimal.Decimal `json:"daily_limit"`
	MonthlyLimit          decimal.Decimal `json:"monthly_limit"`
	AIAccessEnabled       bool            `json:"ai_access_enabled"`
	AIDailyLimit          decimal.Decimal `json:"ai_daily_limit"`
	AIPerTransactionLimit decimal.Decimal `json:"ai_per_transaction_limit"`
	CreatedAt             string          `json:"created_at"`
	UpdatedAt             string          `json:"updated_at"`
}

// Update Wallet Settings Request
type UpdateWalletSettingsRequest struct {
	AIDailyLimit          *decimal.Decimal `json:"ai_daily_limit" validate:"omitempty,gt=0"`
	AIPerTransactionLimit *decimal.Decimal `json:"ai_per_transaction_limit" validate:"omitempty,gt=0"`
	AIAccessEnabled       *bool            `json:"ai_access_enabled"`
	DailyLimit            *decimal.Decimal `json:"daily_limit" validate:"omitempty,gt=0"`
	MonthlyLimit          *decimal.Decimal `json:"monthly_limit" validate:"omitempty,gt=0"`
}

// Wallet Settings Response
type WalletSettingsResponse struct {
	EmailNotifications             bool            `json:"email_notifications"`
	SMSNotifications               bool            `json:"sms_notifications"`
	PushNotifications              bool            `json:"push_notifications"`
	RequireOTPForLargeTransactions bool            `json:"require_otp_large_transactions"`
	LargeTransactionThreshold      decimal.Decimal `json:"large_transaction_threshold"`
	AutoBlockSuspiciousActivity    bool            `json:"auto_block_suspicious"`
	AISpendingNotifications        bool            `json:"ai_spending_notifications"`
	AIApprovalRequired             bool            `json:"ai_approval_required"`
	AIAllowedMerchants             []string        `json:"ai_allowed_merchants"`
	AIBlockedCategories            []string        `json:"ai_blocked_categories"`
	AIWeeklySpendingLimit          decimal.Decimal `json:"ai_weekly_spending_limit"`
}

// Update Wallet Settings Advanced Request
type UpdateWalletSettingsAdvancedRequest struct {
	EmailNotifications             *bool            `json:"email_notifications"`
	SMSNotifications               *bool            `json:"sms_notifications"`
	PushNotifications              *bool            `json:"push_notifications"`
	RequireOTPForLargeTransactions *bool            `json:"require_otp_large_transactions"`
	LargeTransactionThreshold      *decimal.Decimal `json:"large_transaction_threshold"`
	AutoBlockSuspiciousActivity    *bool            `json:"auto_block_suspicious"`
	AISpendingNotifications        *bool            `json:"ai_spending_notifications"`
	AIApprovalRequired             *bool            `json:"ai_approval_required"`
	AIAllowedMerchants             []string         `json:"ai_allowed_merchants"`
	AIBlockedCategories            []string         `json:"ai_blocked_categories"`
	AIWeeklySpendingLimit          *decimal.Decimal `json:"ai_weekly_spending_limit"`
}