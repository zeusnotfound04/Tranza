package dto

import "time"

// CreateAPIKeyRequest represents the request to create a new API key
type CreateAPIKeyRequest struct {
	Label    string `json:"label" binding:"required,max=100"`
	TTLHours int    `json:"ttl_hours" binding:"min=0,max=8760"` // Max 1 year
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
