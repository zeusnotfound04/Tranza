package dto

import (
	"github.com/shopspring/decimal"
)

// Wallet Analytics Request
type WalletAnalyticsRequest struct {
	StartDate string `json:"start_date" form:"start_date" validate:"omitempty"`
	EndDate   string `json:"end_date" form:"end_date" validate:"omitempty"`
	Period    string `json:"period" form:"period" validate:"omitempty,oneof=day week month year"`
}

// Wallet Analytics Response
type WalletAnalyticsResponse struct {
	Period               string                      `json:"period"`
	TotalInflow          decimal.Decimal             `json:"total_inflow"`
	TotalOutflow         decimal.Decimal             `json:"total_outflow"`
	NetFlow              decimal.Decimal             `json:"net_flow"`
	TransactionCount     int64                       `json:"transaction_count"`
	AverageTransaction   decimal.Decimal             `json:"average_transaction"`
	LargestTransaction   decimal.Decimal             `json:"largest_transaction"`
	SmallestTransaction  decimal.Decimal             `json:"smallest_transaction"`
	InflowByMethod       map[string]decimal.Decimal  `json:"inflow_by_method"`
	OutflowByCategory    map[string]decimal.Decimal  `json:"outflow_by_category"`
	DailyBreakdown       []DailyAnalytics            `json:"daily_breakdown"`
	MonthlyComparison    []MonthlyAnalytics          `json:"monthly_comparison"`
	SpendingPattern      SpendingPattern             `json:"spending_pattern"`
	Recommendations      []string                    `json:"recommendations"`
}

// Daily Analytics
type DailyAnalytics struct {
	Date         string          `json:"date"`
	Inflow       decimal.Decimal `json:"inflow"`
	Outflow      decimal.Decimal `json:"outflow"`
	NetFlow      decimal.Decimal `json:"net_flow"`
	Transactions int             `json:"transactions"`
	Balance      decimal.Decimal `json:"balance"`
}

// Monthly Analytics
type MonthlyAnalytics struct {
	Month        string          `json:"month"`
	Year         int             `json:"year"`
	Inflow       decimal.Decimal `json:"inflow"`
	Outflow      decimal.Decimal `json:"outflow"`
	NetFlow      decimal.Decimal `json:"net_flow"`
	Transactions int             `json:"transactions"`
	AvgBalance   decimal.Decimal `json:"avg_balance"`
}

// Spending Pattern
type SpendingPattern struct {
	PeakSpendingHour    int             `json:"peak_spending_hour"`
	PeakSpendingDay     string          `json:"peak_spending_day"`
	AverageDaily        decimal.Decimal `json:"average_daily"`
	AverageWeekly       decimal.Decimal `json:"average_weekly"`
	AverageMonthly      decimal.Decimal `json:"average_monthly"`
	SpendingVolatility  string          `json:"spending_volatility"` // low, medium, high
	SavingsRate         decimal.Decimal `json:"savings_rate"`
	TopMerchants        []MerchantSpending `json:"top_merchants"`
	TopCategories       []CategorySpending `json:"top_categories"`
}

// Merchant Spending
type MerchantSpending struct {
	MerchantName string          `json:"merchant_name"`
	Amount       decimal.Decimal `json:"amount"`
	Transactions int             `json:"transactions"`
	Percentage   decimal.Decimal `json:"percentage"`
}

// Category Spending
type CategorySpending struct {
	Category     string          `json:"category"`
	Amount       decimal.Decimal `json:"amount"`
	Transactions int             `json:"transactions"`
	Percentage   decimal.Decimal `json:"percentage"`
}