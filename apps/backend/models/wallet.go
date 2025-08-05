package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Wallet struct {
	ID     uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID uuid.UUID `gorm:"type:uuid;unique;not null" json:"user_id"`

	Balance  decimal.Decimal `gorm:"type:decimal(15,2);default:0.00" json:"balance"`
	Currency string          `gorm:"size:3;default:'INR'" json:"currency"`
	Status   string          `gorm:"type:varchar(20);default:'active'" json:"status"` // active, frozen, closed

	DailyLimit   decimal.Decimal `gorm:"type:decimal(15,2);default:10000.00" json:"daily_limit"`
	MonthlyLimit decimal.Decimal `gorm:"type:decimal(15,2);default:100000.00" json:"monthly_limit"`

	AIAccessEnabled       bool            `gorm:"default:true" json:"ai_access_enabled"`
	AIDailyLimit          decimal.Decimal `gorm:"type:decimal(15,2);default:1000.00" json:"ai_daily_limit"`
	AIPerTransactionLimit decimal.Decimal `gorm:"type:decimal(15,2);default:500.00" json:"ai_per_transaction_limit"`

	RazorpayCustomerID string `gorm:"size:255" json:"razorpay_customer_id"` // Store Razorpay customer ID

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relations will be loaded separately to avoid circular dependencies
}

// HasSufficientBalance checks if wallet has sufficient balance for the given amount
func (w *Wallet) HasSufficientBalance(amount decimal.Decimal) bool {
	return w.Balance.GreaterThanOrEqual(amount)
}
