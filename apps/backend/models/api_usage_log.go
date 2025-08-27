package models

import (
	"time"

	"github.com/google/uuid"
)

// APIUsageLog represents a detailed log entry for each API request
type APIUsageLog struct {
	ID       uint      `gorm:"primaryKey"`
	APIKeyID uint      `gorm:"index;not null"`
	UserID   uuid.UUID `gorm:"type:uuid;index;not null"`

	// Request Details
	Method    string `gorm:"type:varchar(10);not null"`
	Endpoint  string `gorm:"type:varchar(255);not null"`
	UserAgent string `gorm:"type:text"`
	IPAddress string `gorm:"type:varchar(45)"`

	// Response Details
	StatusCode   int   `gorm:"not null"`
	ResponseTime int64 `gorm:"not null"` // Response time in milliseconds

	// Command/Action Details (for Slack bot and other integrations)
	Command       string `gorm:"type:varchar(100)"` // e.g., "/send-money", "/fetch-balance"
	CommandParams string `gorm:"type:text"`         // JSON encoded command parameters

	// Financial Details
	AmountInvolved float64 `gorm:"type:decimal(15,2);default:0"` // Amount in the transaction if any
	Currency       string  `gorm:"type:varchar(10);default:'INR'"`

	// Error Details
	ErrorMessage string `gorm:"type:text"` // Error message if request failed

	// Metadata
	Metadata string `gorm:"type:text"`                      // JSON encoded additional metadata
	Source   string `gorm:"type:varchar(50);default:'api'"` // 'api', 'slack-bot', 'web-dashboard'

	// Timestamps
	CreatedAt time.Time
	UpdatedAt time.Time

	// Relationships
	APIKey APIKey `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

// APIUsageStats represents aggregated usage statistics for an API key
type APIUsageStats struct {
	KeyID   uint   `json:"key_id"`
	Label   string `json:"label"`
	KeyType string `json:"key_type"`

	// Basic Usage Stats
	TotalRequests      int64     `json:"total_requests"`
	SuccessfulRequests int64     `json:"successful_requests"`
	FailedRequests     int64     `json:"failed_requests"`
	LastUsedAt         time.Time `json:"last_used_at"`

	// Financial Stats
	TotalAmountSpent float64 `json:"total_amount_spent"`
	Currency         string  `json:"currency"`
	SpendingLimit    float64 `json:"spending_limit"` // From API key settings
	RemainingLimit   float64 `json:"remaining_limit"`

	// Performance Stats
	AverageResponseTime float64 `json:"average_response_time"` // in milliseconds

	// Command Stats (for bot usage)
	TopCommands []CommandUsage `json:"top_commands"`

	// Time-based Stats
	RequestsLast24Hours int64 `json:"requests_last_24_hours"`
	RequestsLast7Days   int64 `json:"requests_last_7_days"`
	RequestsLast30Days  int64 `json:"requests_last_30_days"`

	// Rate Limiting
	RateLimit      int     `json:"rate_limit"`
	RateLimitUsage float64 `json:"rate_limit_usage"` // percentage of rate limit used
}

// CommandUsage represents usage statistics for a specific command
type CommandUsage struct {
	Command     string    `json:"command"`
	Count       int64     `json:"count"`
	LastUsed    time.Time `json:"last_used"`
	SuccessRate float64   `json:"success_rate"`
}

// APIUsageTimeSeriesData represents time-series data for usage analytics
type APIUsageTimeSeriesData struct {
	Date            string  `json:"date"`
	TotalRequests   int64   `json:"total_requests"`
	SuccessRequests int64   `json:"success_requests"`
	FailedRequests  int64   `json:"failed_requests"`
	TotalAmount     float64 `json:"total_amount"`
	AvgResponseTime float64 `json:"avg_response_time"`
}

// APIUsageCommandData represents command-specific usage data
type APIUsageCommandData struct {
	Command         string    `json:"command"`
	TotalRequests   int64     `json:"total_requests"`
	SuccessRequests int64     `json:"success_requests"`
	FailedRequests  int64     `json:"failed_requests"`
	TotalAmount     float64   `json:"total_amount"`
	AvgResponseTime float64   `json:"avg_response_time"`
	LastUsed        time.Time `json:"last_used"`
}

// APIUsageSummary represents a comprehensive summary of API usage
type APIUsageSummary struct {
	Stats          APIUsageStats            `json:"stats"`
	TimeSeriesData []APIUsageTimeSeriesData `json:"time_series_data"`
	CommandData    []APIUsageCommandData    `json:"command_data"`
	RecentLogs     []APIUsageLog            `json:"recent_logs"`
}
