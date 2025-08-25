package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// ExternalTransfer represents a money transfer to external accounts (UPI/Phone)
type ExternalTransfer struct {
	ID       uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID   uuid.UUID `json:"user_id" gorm:"type:uuid;not null;index"`
	WalletID uuid.UUID `json:"wallet_id" gorm:"type:uuid;not null;index"`

	// Transfer Details
	Amount      decimal.Decimal `json:"amount" gorm:"type:decimal(15,2);not null"`
	Currency    string          `json:"currency" gorm:"default:'INR';not null"`
	Description string          `json:"description" gorm:"size:500"`

	// Recipient Details
	RecipientType  string `json:"recipient_type" gorm:"not null"`  // "upi" or "phone"
	RecipientValue string `json:"recipient_value" gorm:"not null"` // UPI ID or Phone Number
	RecipientName  string `json:"recipient_name,omitempty" gorm:"size:100"`

	// Transfer Status & References
	Status         string `json:"status" gorm:"default:'pending';not null"` // pending, processing, success, failed, cancelled
	TransferMethod string `json:"transfer_method" gorm:"not null"`          // "razorpay_payout", "upi_direct"

	// External Payment Gateway References
	RazorpayPayoutID  string `json:"razorpay_payout_id,omitempty" gorm:"size:50;index"`
	RazorpayContactID string `json:"razorpay_contact_id,omitempty" gorm:"size:50"`
	RazorpayFundID    string `json:"razorpay_fund_id,omitempty" gorm:"size:50"`

	// Financial Details
	TransferFee  decimal.Decimal `json:"transfer_fee" gorm:"type:decimal(10,2);default:0"`
	TotalAmount  decimal.Decimal `json:"total_amount" gorm:"type:decimal(15,2);not null"` // Amount + Fee
	ExchangeRate decimal.Decimal `json:"exchange_rate,omitempty" gorm:"type:decimal(10,6);default:1"`

	// Tracking & Audit
	ReferenceID   string     `json:"reference_id" gorm:"unique;not null;size:50"`     // Internal tracking
	TransactionID *uuid.UUID `json:"transaction_id,omitempty" gorm:"type:uuid;index"` // Link to main transaction

	// Security & Compliance
	IPAddress   string `json:"ip_address,omitempty" gorm:"size:45"`
	UserAgent   string `json:"user_agent,omitempty" gorm:"size:500"`
	InitiatedBy string `json:"initiated_by" gorm:"default:'user';not null"` // "user", "bot", "api"

	// Status Tracking
	ProcessedAt   *time.Time `json:"processed_at,omitempty"`
	CompletedAt   *time.Time `json:"completed_at,omitempty"`
	FailureReason string     `json:"failure_reason,omitempty" gorm:"size:500"`
	RetryCount    int        `json:"retry_count" gorm:"default:0"`
	MaxRetries    int        `json:"max_retries" gorm:"default:3"`

	// Wallet Balance Tracking
	BalanceBefore decimal.Decimal `json:"balance_before" gorm:"type:decimal(15,2)"`
	BalanceAfter  decimal.Decimal `json:"balance_after" gorm:"type:decimal(15,2)"`

	// Relationships
	User        User         `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Wallet      Wallet       `json:"wallet,omitempty" gorm:"foreignKey:WalletID"`
	Transaction *Transaction `json:"transaction,omitempty" gorm:"foreignKey:TransactionID"`

	// Timestamps
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName returns the table name for ExternalTransfer
func (ExternalTransfer) TableName() string {
	return "external_transfers"
}

// External Transfer Status Constants
const (
	ExternalTransferStatusPending    = "pending"
	ExternalTransferStatusProcessing = "processing"
	ExternalTransferStatusSuccess    = "success"
	ExternalTransferStatusFailed     = "failed"
	ExternalTransferStatusCancelled  = "cancelled"
	ExternalTransferStatusRefunded   = "refunded"
)

// Recipient Type Constants
const (
	RecipientTypeUPI   = "upi"
	RecipientTypePhone = "phone"
	RecipientTypeIFSC  = "ifsc" // For future bank transfer support
)

// Transfer Method Constants
const (
	TransferMethodRazorpayPayout = "razorpay_payout"
	TransferMethodUPIDirect      = "upi_direct"
	TransferMethodBankTransfer   = "bank_transfer"
)

// Initiated By Constants
const (
	InitiatedByUser = "user"
	InitiatedByBot  = "bot"
	InitiatedByAPI  = "api"
)

// BeforeCreate hook to set UUID and reference ID
func (et *ExternalTransfer) BeforeCreate(tx *gorm.DB) error {
	if et.ID == uuid.Nil {
		et.ID = uuid.New()
	}

	if et.ReferenceID == "" {
		// Generate unique reference ID
		et.ReferenceID = GenerateExternalTransferReference()
	}

	return nil
}

// IsCompleted returns true if transfer is in a final state
func (et *ExternalTransfer) IsCompleted() bool {
	return et.Status == ExternalTransferStatusSuccess ||
		et.Status == ExternalTransferStatusFailed ||
		et.Status == ExternalTransferStatusCancelled ||
		et.Status == ExternalTransferStatusRefunded
}

// CanRetry returns true if transfer can be retried
func (et *ExternalTransfer) CanRetry() bool {
	return et.Status == ExternalTransferStatusFailed &&
		et.RetryCount < et.MaxRetries
}

// GetDisplayStatus returns user-friendly status message
func (et *ExternalTransfer) GetDisplayStatus() string {
	switch et.Status {
	case ExternalTransferStatusPending:
		return "Transfer Initiated"
	case ExternalTransferStatusProcessing:
		return "Processing Transfer"
	case ExternalTransferStatusSuccess:
		return "Transfer Completed"
	case ExternalTransferStatusFailed:
		return "Transfer Failed"
	case ExternalTransferStatusCancelled:
		return "Transfer Cancelled"
	case ExternalTransferStatusRefunded:
		return "Transfer Refunded"
	default:
		return "Unknown Status"
	}
}

// GetRecipientDisplay returns formatted recipient display
func (et *ExternalTransfer) GetRecipientDisplay() string {
	switch et.RecipientType {
	case RecipientTypeUPI:
		return et.RecipientValue
	case RecipientTypePhone:
		// Mask phone number for security
		if len(et.RecipientValue) >= 10 {
			return "****" + et.RecipientValue[len(et.RecipientValue)-4:]
		}
		return et.RecipientValue
	default:
		return et.RecipientValue
	}
}
