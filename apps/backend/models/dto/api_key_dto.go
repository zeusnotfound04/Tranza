package dto

import "time"

// CreateAPIKeyRequest represents the request to create a new API key
type CreateAPIKeyRequest struct {
	Label    string `json:"label" binding:"required,max=100"`
	Password string `json:"password" binding:"required,min=6,max=50"` // Password to protect the API key
	TTLHours int    `json:"ttl_hours" binding:"min=0,max=8760"`       // Max 1 year
}

// CreateAPIKeyResponse represents the response after creating an API key
type CreateAPIKeyResponse struct {
	APIKey   string `json:"api_key"`
	Label    string `json:"label"`
	TTLHours int    `json:"ttl_hours"`
	Message  string `json:"message"`
}

// CreateBotAPIKeyRequest represents the request to create a bot API key
type CreateBotAPIKeyRequest struct {
	Label       string `json:"label" binding:"required,max=100"`
	Password    string `json:"password" binding:"required,min=6,max=50"` // Password to protect the API key
	WorkspaceID string `json:"workspace_id" binding:"required"`
	BotUserID   string `json:"bot_user_id" binding:"required"`
	TTLHours    int    `json:"ttl_hours" binding:"min=0,max=8760"`
}

// CreateBotAPIKeyResponse represents the response after creating a bot API key
type CreateBotAPIKeyResponse struct {
	APIKey      string   `json:"api_key"`
	Label       string   `json:"label"`
	WorkspaceID string   `json:"workspace_id"`
	BotUserID   string   `json:"bot_user_id"`
	TTLHours    int      `json:"ttl_hours"`
	Scopes      []string `json:"scopes"`
	Message     string   `json:"message"`
}

// APIKeyInfo represents basic information about an API key
type APIKeyInfo struct {
	ID           uint       `json:"id"`
	Label        string     `json:"label"`
	KeyType      string     `json:"key_type"`
	Scopes       []string   `json:"scopes"`
	UsageCount   int64      `json:"usage_count"`
	RateLimit    int        `json:"rate_limit"`
	IsActive     bool       `json:"is_active"`
	CreatedAt    time.Time  `json:"created_at"`
	ExpiresAt    *time.Time `json:"expires_at,omitempty"`
	LastUsedAt   time.Time  `json:"last_used_at"`
	BotWorkspace *string    `json:"bot_workspace,omitempty"`
	BotUserID    *string    `json:"bot_user_id,omitempty"`
}

// ListAPIKeysResponse represents the response when listing API keys
type ListAPIKeysResponse struct {
	Keys  []APIKeyInfo `json:"keys"`
	Total int          `json:"total"`
}

// RotateAPIKeyResponse represents the response after rotating an API key
type RotateAPIKeyResponse struct {
	NewAPIKey string `json:"new_api_key"`
	Message   string `json:"message"`
}

// ViewAPIKeyRequest represents the request to view an API key with password
type ViewAPIKeyRequest struct {
	Password string `json:"password" binding:"required"`
}

// ViewAPIKeyResponse represents the response when viewing an API key
type ViewAPIKeyResponse struct {
	APIKey  string `json:"api_key"`
	Message string `json:"message"`
}

// API Usage DTOs

// APIUsageStatsResponse represents the response for API usage statistics
type APIUsageStatsResponse struct {
	KeyID               uint              `json:"key_id"`
	Label               string            `json:"label"`
	KeyType             string            `json:"key_type"`
	TotalRequests       int64             `json:"total_requests"`
	SuccessfulRequests  int64             `json:"successful_requests"`
	FailedRequests      int64             `json:"failed_requests"`
	LastUsedAt          time.Time         `json:"last_used_at"`
	TotalAmountSpent    float64           `json:"total_amount_spent"`
	Currency            string            `json:"currency"`
	SpendingLimit       float64           `json:"spending_limit"`
	RemainingLimit      float64           `json:"remaining_limit"`
	AverageResponseTime float64           `json:"average_response_time"`
	RequestsLast24Hours int64             `json:"requests_last_24_hours"`
	RequestsLast7Days   int64             `json:"requests_last_7_days"`
	RequestsLast30Days  int64             `json:"requests_last_30_days"`
	RateLimit           int               `json:"rate_limit"`
	RateLimitUsage      float64           `json:"rate_limit_usage"`
	TopCommands         []CommandUsageDTO `json:"top_commands"`
}

// CommandUsageDTO represents command usage statistics
type CommandUsageDTO struct {
	Command     string    `json:"command"`
	Count       int64     `json:"count"`
	LastUsed    time.Time `json:"last_used"`
	SuccessRate float64   `json:"success_rate"`
}

// APIUsageLogDTO represents a single API usage log entry
type APIUsageLogDTO struct {
	ID             uint      `json:"id"`
	Method         string    `json:"method"`
	Endpoint       string    `json:"endpoint"`
	StatusCode     int       `json:"status_code"`
	ResponseTime   int64     `json:"response_time"`
	Command        string    `json:"command,omitempty"`
	AmountInvolved float64   `json:"amount_involved,omitempty"`
	Currency       string    `json:"currency,omitempty"`
	Source         string    `json:"source"`
	CreatedAt      time.Time `json:"created_at"`
	ErrorMessage   string    `json:"error_message,omitempty"`
}

// TimeSeriesDataDTO represents time-series data for charts
type TimeSeriesDataDTO struct {
	Date            string  `json:"date"`
	TotalRequests   int64   `json:"total_requests"`
	SuccessRequests int64   `json:"success_requests"`
	FailedRequests  int64   `json:"failed_requests"`
	TotalAmount     float64 `json:"total_amount"`
	AvgResponseTime float64 `json:"avg_response_time"`
}

// CommandDataDTO represents command-specific usage data
type CommandDataDTO struct {
	Command         string    `json:"command"`
	TotalRequests   int64     `json:"total_requests"`
	SuccessRequests int64     `json:"success_requests"`
	FailedRequests  int64     `json:"failed_requests"`
	TotalAmount     float64   `json:"total_amount"`
	AvgResponseTime float64   `json:"avg_response_time"`
	LastUsed        time.Time `json:"last_used"`
}

// UsageSummaryResponse represents comprehensive usage data
type UsageSummaryResponse struct {
	Stats          APIUsageStatsResponse `json:"stats"`
	TimeSeriesData []TimeSeriesDataDTO   `json:"time_series_data"`
	CommandData    []CommandDataDTO      `json:"command_data"`
	RecentLogs     []APIUsageLogDTO      `json:"recent_logs"`
}
