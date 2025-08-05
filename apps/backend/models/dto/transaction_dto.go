package dto

import (
	"github.com/shopspring/decimal"
)

// Transaction Response
type TransactionResponse struct {
	ID            string          `json:"id"`
	Type          string          `json:"type"`
	Amount        decimal.Decimal `json:"amount"`
	BalanceAfter  decimal.Decimal `json:"balance_after"`
	Currency      string          `json:"currency"`
	Description   string          `json:"description"`
	Status        string          `json:"status"`
	PaymentMethod string          `json:"payment_method,omitempty"`
	AIAgentID     string          `json:"ai_agent_id,omitempty"`
	MerchantName  string          `json:"merchant_name,omitempty"`
	MerchantUPIID string          `json:"merchant_upi_id,omitempty"`
	ReferenceID   string          `json:"reference_id"`
	FailureReason string          `json:"failure_reason,omitempty"`
	CreatedAt     string          `json:"created_at"`
	UpdatedAt     string          `json:"updated_at"`
}

// Transaction History Request
type TransactionHistoryRequest struct {
	Page            int    `json:"page" form:"page" validate:"min=1"`
	Limit           int    `json:"limit" form:"limit" validate:"min=1,max=100"`
	TransactionType string `json:"type" form:"type"`
	Status          string `json:"status" form:"status"`
	StartDate       string `json:"start_date" form:"start_date"`
	EndDate         string `json:"end_date" form:"end_date"`
	MinAmount       string `json:"min_amount" form:"min_amount"`
	MaxAmount       string `json:"max_amount" form:"max_amount"`
}

// Transaction Statistics Response
type TransactionStatsResponse struct {
	TotalTransactions    int64           `json:"total_transactions"`
	TotalAmount          decimal.Decimal `json:"total_amount"`
	TodayTransactions    int64           `json:"today_transactions"`
	TodayAmount          decimal.Decimal `json:"today_amount"`
	MonthTransactions    int64           `json:"month_transactions"`
	MonthAmount          decimal.Decimal `json:"month_amount"`
	AITransactions       int64           `json:"ai_transactions"`
	AIAmount             decimal.Decimal `json:"ai_amount"`
	LastTransactionDate  string          `json:"last_transaction_date"`
	MostUsedPaymentMethod string         `json:"most_used_payment_method"`
}