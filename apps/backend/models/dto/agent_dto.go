package dto

import (
	"github.com/shopspring/decimal"
)

// AI Payment Request
type AIPaymentRequest struct {
	AgentID       string          `json:"agent_id" validate:"required,min=3,max=100" binding:"required"`
	Amount        decimal.Decimal `json:"amount" validate:"required,gt=0" binding:"required"`
	MerchantName  string          `json:"merchant_name" validate:"required,min=2,max=255" binding:"required"`
	MerchantUPIID string          `json:"merchant_upi_id" validate:"omitempty"`
	Description   string          `json:"description" validate:"required,min=5,max=500" binding:"required"`
	Category      string          `json:"category" validate:"omitempty,max=100"`
	Metadata      map[string]interface{} `json:"metadata"`
}

// AI Payment Response
type AIPaymentResponse struct {
	Success       bool            `json:"success"`
	TransactionID string          `json:"transaction_id"`
	NewBalance    decimal.Decimal `json:"new_balance"`
	Message       string          `json:"message"`
	Amount        decimal.Decimal `json:"amount"`
	MerchantName  string          `json:"merchant_name"`
	AgentID       string          `json:"agent_id"`
	Timestamp     string          `json:"timestamp"`
}

// AI Spending Limits Response
type AISpendingLimitsResponse struct {
	Balance                   decimal.Decimal `json:"balance"`
	AIAccessEnabled           bool            `json:"ai_access_enabled"`
	AIDailyLimit              decimal.Decimal `json:"ai_daily_limit"`
	AIPerTransactionLimit     decimal.Decimal `json:"ai_per_transaction_limit"`
	AIWeeklyLimit             decimal.Decimal `json:"ai_weekly_limit"`
	DailySpent                decimal.Decimal `json:"daily_spent"`
	WeeklySpent               decimal.Decimal `json:"weekly_spent"`
	DailyRemaining            decimal.Decimal `json:"daily_remaining"`
	WeeklyRemaining           decimal.Decimal `json:"weekly_remaining"`
	TransactionsToday         int             `json:"transactions_today"`
	TransactionsThisWeek      int             `json:"transactions_this_week"`
	SpendingDate              string          `json:"spending_date"`
	LastTransactionTime       string          `json:"last_transaction_time,omitempty"`
	RecommendedTransactionLimit decimal.Decimal `json:"recommended_transaction_limit"`
}

// AI Agent Configuration Request
type AIAgentConfigRequest struct {
	AgentID                string            `json:"agent_id" validate:"required"`
	Enabled                *bool             `json:"enabled"`
	DailyLimit             *decimal.Decimal  `json:"daily_limit" validate:"omitempty,gt=0"`
	PerTransactionLimit    *decimal.Decimal  `json:"per_transaction_limit" validate:"omitempty,gt=0"`
	AllowedMerchants       []string          `json:"allowed_merchants"`
	BlockedMerchants       []string          `json:"blocked_merchants"`
	AllowedCategories      []string          `json:"allowed_categories"`
	BlockedCategories      []string          `json:"blocked_categories"`
	RequireApproval        *bool             `json:"require_approval"`
	AutoApprovalThreshold  *decimal.Decimal  `json:"auto_approval_threshold"`
	NotificationSettings   map[string]bool   `json:"notification_settings"`
}

// AI Agent Configuration Response
type AIAgentConfigResponse struct {
	AgentID                string            `json:"agent_id"`
	Enabled                bool              `json:"enabled"`
	DailyLimit             decimal.Decimal   `json:"daily_limit"`
	PerTransactionLimit    decimal.Decimal   `json:"per_transaction_limit"`
	AllowedMerchants       []string          `json:"allowed_merchants"`
	BlockedMerchants       []string          `json:"blocked_merchants"`
	AllowedCategories      []string          `json:"allowed_categories"`
	BlockedCategories      []string          `json:"blocked_categories"`
	RequireApproval        bool              `json:"require_approval"`
	AutoApprovalThreshold  decimal.Decimal   `json:"auto_approval_threshold"`
	NotificationSettings   map[string]bool   `json:"notification_settings"`
	CreatedAt              string            `json:"created_at"`
	UpdatedAt              string            `json:"updated_at"`
}

// AI Analytics Response
type AIAnalyticsResponse struct {
	TotalAITransactions     int64             `json:"total_ai_transactions"`
	TotalAISpending         decimal.Decimal   `json:"total_ai_spending"`
	AverageTransactionAmount decimal.Decimal  `json:"average_transaction_amount"`
	MostActiveAgent         string            `json:"most_active_agent"`
	MostUsedMerchant        string            `json:"most_used_merchant"`
	SpendingByCategory      map[string]decimal.Decimal `json:"spending_by_category"`
	SpendingByAgent         map[string]decimal.Decimal `json:"spending_by_agent"`
	DailySpending           map[string]decimal.Decimal `json:"daily_spending"`
	WeeklyTrend             []AISpendingTrend `json:"weekly_trend"`
	SavingsRecommendations  []string          `json:"savings_recommendations"`
}

// AI Spending Trend
type AISpendingTrend struct {
	Date         string          `json:"date"`
	Amount       decimal.Decimal `json:"amount"`
	Transactions int             `json:"transactions"`
	TopMerchant  string          `json:"top_merchant"`
}
// AI Agent Response