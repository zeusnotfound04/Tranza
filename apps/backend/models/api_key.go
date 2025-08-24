package models

import (
	"encoding/json"
	"time"
)

type APIKey struct {
	ID         uint   `gorm:"primaryKey"`
	UserID     uint   `gorm:"not null"`
	KeyHash    string `gorm:"not null;uniqueIndex"`
	Label      string `gorm:"type:varchar(100)"`
	Scopes     string `gorm:"type:text"`                       // JSON array of scopes
	KeyType    string `gorm:"type:varchar(20);default:'user'"` // 'user', 'bot', 'admin'
	UsageCount int64  `gorm:"default:0"`
	RateLimit  int    `gorm:"default:1000"` // Requests per hour
	CreatedAt  time.Time
	ExpiresAt  *time.Time
	LastUsedAt time.Time
	IsActive   bool `gorm:"default:true"`

	// Additional bot-specific fields
	BotWorkspace string `gorm:"type:varchar(100)"` // Slack workspace ID for bot keys
	BotUserID    string `gorm:"type:varchar(100)"` // Bot's associated user ID
}

// GetScopes returns the scopes as a slice of strings
func (k *APIKey) GetScopes() []string {
	if k.Scopes == "" {
		return []string{}
	}

	var scopes []string
	json.Unmarshal([]byte(k.Scopes), &scopes)
	return scopes
}

// SetScopes sets the scopes from a slice of strings
func (k *APIKey) SetScopes(scopes []string) error {
	scopesJSON, err := json.Marshal(scopes)
	if err != nil {
		return err
	}
	k.Scopes = string(scopesJSON)
	return nil
}

// HasScope checks if the API key has a specific scope
func (k *APIKey) HasScope(scope string) bool {
	scopes := k.GetScopes()
	for _, s := range scopes {
		if s == scope || s == "*" { // "*" means all permissions
			return true
		}
	}
	return false
}

// IsExpired checks if the API key has expired
func (k *APIKey) IsExpired() bool {
	return k.ExpiresAt != nil && k.ExpiresAt.Before(time.Now())
}

// IsBot checks if this is a bot API key
func (k *APIKey) IsBot() bool {
	return k.KeyType == "bot"
}

// CanMakeRequest checks if the key can make requests (not expired, active, etc.)
func (k *APIKey) CanMakeRequest() bool {
	return k.IsActive && !k.IsExpired()
}

// IncrementUsage increments the usage counter and updates last used timestamp
func (k *APIKey) IncrementUsage() {
	k.UsageCount++
	k.LastUsedAt = time.Now()
}
