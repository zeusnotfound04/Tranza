package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type TransactionStatus string

const (
	StatusFailed  TransactionStatus = "FAILED"
	StatusSuccess TransactionStatus = "SUCCESS"
	StatusPending TransactionStatus = "PENDING"
)

type Transaction struct {
	ID       uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	WalletID uuid.UUID `gorm:"type:uuid;not null" json:"wallet_id"`
	UserID   uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	CardID   *uint     `gorm:"null" json:"card_id,omitempty"`

	Type         string            `gorm:"type:varchar(50);not null" json:"type"`
	Amount       decimal.Decimal   `gorm:"type:decimal(15,2);not null" json:"amount"`
	BalanceAfter decimal.Decimal   `gorm:"type:decimal(15,2)" json:"balance_after"`
	Currency     string            `gorm:"type:varchar(3);default:'INR'" json:"currency"`
	Description  string            `gorm:"type:text" json:"description"`
	Status       TransactionStatus `gorm:"type:varchar(10);default:'PENDING'" json:"status"`

	// Payment related fields
	PaymentMethod     string `gorm:"type:varchar(50)" json:"payment_method"`
	RazorpayOrderID   string `gorm:"type:varchar(255)" json:"razorpay_order_id"`
	RazorpayPaymentID string `gorm:"type:varchar(255)" json:"razorpay_payment_id"`
	ReferenceID       string `gorm:"type:varchar(255);unique" json:"reference_id"`

	// AI related fields
	AIAgentID     string `gorm:"type:varchar(255)" json:"ai_agent_id"`
	MerchantName  string `gorm:"type:varchar(255)" json:"merchant_name"`
	MerchantUPIID string `gorm:"type:varchar(255)" json:"merchant_upi_id"`

	// Failure information
	FailureReason string `gorm:"type:text" json:"failure_reason"`

	// Audit fields
	IPAddress string `gorm:"type:varchar(45)" json:"ip_address"`
	UserAgent string `gorm:"type:text" json:"user_agent"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
