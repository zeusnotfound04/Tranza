package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AIPaymentRequest struct {
	ID             uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID         uuid.UUID  `json:"user_id" gorm:"type:uuid;not null;index"`
	Amount         float64    `json:"amount" gorm:"not null"`
	Description    string     `json:"description" gorm:"type:text"`
	MerchantName   string     `json:"merchant_name,omitempty"`
	AIPrompt       string     `json:"ai_prompt" gorm:"type:text;not null"`
	ConfirmationID string     `json:"confirmation_id,omitempty"`
	Status         string     `json:"status" gorm:"default:'pending'"` // pending, confirmed, processed, failed, cancelled
	ProcessedAt    *time.Time `json:"processed_at,omitempty"`
	TransactionID  string     `json:"transaction_id,omitempty"`
	AIResponse     string     `json:"ai_response,omitempty" gorm:"type:text"`
	RiskLevel      string     `json:"risk_level" gorm:"default:'low'"` // low, medium, high
	Confidence     float64    `json:"confidence" gorm:"default:0.0"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type AISpendingLimit struct {
	ID                    uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID                uuid.UUID `json:"user_id" gorm:"type:uuid;not null;unique;index"`
	DailyLimit            float64   `json:"daily_limit" gorm:"default:10000"`
	TransactionLimit      float64   `json:"transaction_limit" gorm:"default:2000"`
	MonthlyLimit          float64   `json:"monthly_limit" gorm:"default:100000"`
	AIAccessEnabled       bool      `json:"ai_access_enabled" gorm:"default:false"`
	RequireConfirmation   bool      `json:"require_confirmation" gorm:"default:true"`
	ConfirmationThreshold float64   `json:"confirmation_threshold" gorm:"default:1000"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

type AISpendingTracker struct {
	ID               uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID           uuid.UUID `json:"user_id" gorm:"type:uuid;not null;index"`
	Date             time.Time `json:"date" gorm:"type:date;not null;index"`
	DailySpent       float64   `json:"daily_spent" gorm:"default:0"`
	TransactionCount int       `json:"transaction_count" gorm:"default:0"`
	MonthlySpent     float64   `json:"monthly_spent" gorm:"default:0"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// DTO Types for API requests/responses
type AIPaymentRequestDTO struct {
	Prompt      string  `json:"prompt" binding:"required"`
	Amount      float64 `json:"amount,omitempty"`
	Merchant    string  `json:"merchant,omitempty"`
	Description string  `json:"description,omitempty"`
}

type AIPaymentResponse struct {
	ID                   string  `json:"id"`
	Amount               float64 `json:"amount"`
	Merchant             string  `json:"merchant"`
	Description          string  `json:"description"`
	Confidence           float64 `json:"confidence"`
	RiskLevel            string  `json:"risk_level"`
	RequiresConfirmation bool    `json:"requires_confirmation"`
	AIReasoning          string  `json:"ai_reasoning"`
	WalletBalance        float64 `json:"wallet_balance"`
	RemainingLimit       float64 `json:"remaining_limit"`
	Suggestions          string  `json:"suggestions,omitempty"`
}

type AIPaymentConfirmationDTO struct {
	PaymentID string `json:"payment_id" binding:"required"`
	Confirmed bool   `json:"confirmed" binding:"required"`
}

type AISpendingLimitsDTO struct {
	DailyLimit            float64 `json:"daily_limit" binding:"min=0,max=50000"`
	TransactionLimit      float64 `json:"transaction_limit" binding:"min=0,max=10000"`
	MonthlyLimit          float64 `json:"monthly_limit" binding:"min=0,max=500000"`
	AIAccessEnabled       bool    `json:"ai_access_enabled"`
	RequireConfirmation   bool    `json:"require_confirmation"`
	ConfirmationThreshold float64 `json:"confirmation_threshold" binding:"min=0"`
}

// BeforeCreate hook
func (a *AIPaymentRequest) BeforeCreate(tx *gorm.DB) (err error) {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	if a.ConfirmationID == "" {
		a.ConfirmationID = uuid.New().String()
	}
	return
}

func (a *AISpendingLimit) BeforeCreate(tx *gorm.DB) (err error) {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return
}

func (a *AISpendingTracker) BeforeCreate(tx *gorm.DB) (err error) {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return
}
