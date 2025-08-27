package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type APIKey struct {
	ID           uint      `gorm:"primaryKey"`
	UserID       uuid.UUID `json:"user_id" gorm:"type:uuid"`
	KeyHash      string    `gorm:"not null;uniqueIndex"`
	EncryptedKey string    `gorm:"not null"` // Encrypted version of the raw key for viewing
	PasswordHash string    `gorm:"not null"` // Hash of the password to view the key
	Label        string    `gorm:"type:varchar(100)"`
	Scopes       string    `gorm:"type:text"`                       // JSON array of scopes
	KeyType      string    `gorm:"type:varchar(20);default:'user'"` // 'user', 'bot', 'admin'
	UsageCount   int64     `gorm:"default:0"`
	RateLimit    int       `gorm:"default:1000"` // Requests per hour

	// Financial limits
	SpendingLimit float64 `gorm:"type:decimal(15,2);default:10000"` // Default spending limit
	SpentAmount   float64 `gorm:"type:decimal(15,2);default:0"`     // Total amount spent via this key
	Currency      string  `gorm:"type:varchar(10);default:'INR'"`

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
	// For universal keys, always return true for backward compatibility
	if k.KeyType == "universal" {
		return true
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

// CanSpend checks if the key can spend the given amount
func (k *APIKey) CanSpend(amount float64) bool {
	return k.SpentAmount+amount <= k.SpendingLimit
}

// AddSpentAmount adds to the spent amount
func (k *APIKey) AddSpentAmount(amount float64) {
	k.SpentAmount += amount
}

// GetRemainingSpendingLimit returns the remaining spending limit
func (k *APIKey) GetRemainingSpendingLimit() float64 {
	remaining := k.SpendingLimit - k.SpentAmount
	if remaining < 0 {
		return 0
	}
	return remaining
}

// IncrementUsage increments the usage counter and updates last used timestamp
func (k *APIKey) IncrementUsage() {
	k.UsageCount++
	k.LastUsedAt = time.Now()
}
