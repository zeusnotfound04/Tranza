package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// ExternalOrder represents orders made through external e-commerce websites
type ExternalOrder struct {
	ID              uuid.UUID       `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID          uuid.UUID       `json:"user_id" gorm:"type:uuid;not null;index"`
	OrderNumber     string          `json:"order_number" gorm:"unique;not null"`
	ExternalOrderID string          `json:"external_order_id" gorm:"not null"` // Order ID from external store
	Website         string          `json:"website" gorm:"not null"`           // Myntra, Amazon, Flipkart, etc.
	TotalAmount     decimal.Decimal `json:"total_amount" gorm:"type:decimal(10,2);not null"`
	Currency        string          `json:"currency" gorm:"default:'INR'"`
	Status          string          `json:"status" gorm:"default:'pending'"` // pending, confirmed, shipped, delivered, cancelled
	PaymentStatus   string          `json:"payment_status" gorm:"default:'pending'"`
	TransactionID   string          `json:"transaction_id,omitempty"`
	TrackingNumber  string          `json:"tracking_number,omitempty"`

	// Delivery details
	DeliveryAddress   Address    `json:"delivery_address" gorm:"embedded;embeddedPrefix:delivery_"`
	EstimatedDelivery *time.Time `json:"estimated_delivery,omitempty"`
	DeliveredAt       *time.Time `json:"delivered_at,omitempty"`

	// Order items from external website
	Items []ExternalOrderItem `json:"items" gorm:"serializer:json"`

	// AI-related fields
	IsAIOrder   bool   `json:"is_ai_order" gorm:"default:true"`
	AIPrompt    string `json:"ai_prompt,omitempty" gorm:"type:text"`
	AIRequestID string `json:"ai_request_id,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relations
	User User `json:"-" gorm:"foreignKey:UserID"`
}

// ExternalOrderItem represents items purchased from external websites
type ExternalOrderItem struct {
	ProductID  string          `json:"product_id"` // External product ID
	Name       string          `json:"name"`
	Brand      string          `json:"brand"`
	Category   string          `json:"category"`
	Size       string          `json:"size"`
	Color      string          `json:"color"`
	Quantity   int             `json:"quantity"`
	Price      decimal.Decimal `json:"price"`
	ImageURL   string          `json:"image_url"`
	ProductURL string          `json:"product_url"` // Link to external product page
	Website    string          `json:"website"`     // Source website
}

// DTO Types for API requests/responses

type AIClothingOrderRequest struct {
	Prompt    string  `json:"prompt" binding:"required"`
	AddressID string  `json:"address_id,omitempty"`
	Budget    float64 `json:"budget,omitempty"`
	Category  string  `json:"category,omitempty"` // shirts, pants, dresses, etc.
}

type AIClothingOrderResponse struct {
	OrderID              string                 `json:"order_id,omitempty"`
	OrderCreated         bool                   `json:"order_created"`
	RequiredInfo         []string               `json:"required_info,omitempty"`
	SuggestedItems       []ProductSuggestion    `json:"suggested_items,omitempty"`
	SuggestedProducts    []SuggestedProduct     `json:"suggested_products,omitempty"`
	TotalEstimate        decimal.Decimal        `json:"total_estimate"`
	Message              string                 `json:"message"`
	RequiresConfirmation bool                   `json:"requires_confirmation"`
	ConfirmationID       string                 `json:"confirmation_id,omitempty"`
	OrderAnalysis        *ClothingOrderAnalysis `json:"order_analysis,omitempty"`
	SelectedAddress      *Address               `json:"selected_address,omitempty"`
}

type ProductSuggestion struct {
	Product  ProductResponse `json:"product"`
	Variant  VariantResponse `json:"variant"`
	Quantity int             `json:"quantity"`
	Reason   string          `json:"reason"`
	Price    decimal.Decimal `json:"price"`
}

type ProductResponse struct {
	ID       string          `json:"id"`
	Name     string          `json:"name"`
	Brand    string          `json:"brand"`
	Category string          `json:"category"`
	Price    decimal.Decimal `json:"price"`
	ImageURL string          `json:"image_url"`
	Rating   float64         `json:"rating"`
	Website  string          `json:"website"`
	URL      string          `json:"url"`
}

type VariantResponse struct {
	ID    string `json:"id"`
	Size  string `json:"size"`
	Color string `json:"color"`
	Stock int    `json:"stock"`
}

type ClothingOrderConfirmationRequest struct {
	ConfirmationID string `json:"confirmation_id" binding:"required"`
	Confirmed      bool   `json:"confirmed" binding:"required"`
}

// Additional DTOs for external API integration
type ClothingOrderAnalysis struct {
	Category string  `json:"category"`
	Size     string  `json:"size"`
	Color    string  `json:"color"`
	Brand    string  `json:"brand"`
	Occasion string  `json:"occasion"`
	MinPrice float64 `json:"min_price"`
	MaxPrice float64 `json:"max_price"`
}

type SuggestedProduct struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Brand       string  `json:"brand"`
	Price       float64 `json:"price"`
	Currency    string  `json:"currency"`
	ImageURL    string  `json:"image_url"`
	Rating      float64 `json:"rating"`
	Description string  `json:"description"`
	URL         string  `json:"url"`
	Website     string  `json:"website"`
}

type ConfirmAIClothingOrderRequest struct {
	SelectedProducts []SelectedProduct `json:"selected_products" binding:"required"`
	AddressID        string            `json:"address_id" binding:"required"`
}

type SelectedProduct struct {
	ID       string  `json:"id" binding:"required"`
	Name     string  `json:"name" binding:"required"`
	Price    float64 `json:"price" binding:"required"`
	Quantity int     `json:"quantity" binding:"required"`
	Size     string  `json:"size"`
	Color    string  `json:"color"`
	Website  string  `json:"website" binding:"required"`
	URL      string  `json:"url" binding:"required"`
}

type ConfirmAIClothingOrderResponse struct {
	Success          bool              `json:"success"`
	OrderID          string            `json:"order_id,omitempty"`
	Message          string            `json:"message"`
	TotalAmount      float64           `json:"total_amount,omitempty"`
	RequiredAmount   float64           `json:"required_amount,omitempty"`
	CurrentBalance   float64           `json:"current_balance,omitempty"`
	RemainingBalance float64           `json:"remaining_balance,omitempty"`
	SelectedProducts []SelectedProduct `json:"selected_products,omitempty"`
}

type ClothingOrderResponse struct {
	OrderID     string              `json:"order_id"`
	OrderNumber string              `json:"order_number"`
	Status      string              `json:"status"`
	TotalAmount decimal.Decimal     `json:"total_amount"`
	PaymentInfo PaymentInfoResponse `json:"payment_info"`
	Items       []ExternalOrderItem `json:"items"`
	CreatedAt   time.Time           `json:"created_at"`
}

type PaymentInfoResponse struct {
	Status        string `json:"status"`
	TransactionID string `json:"transaction_id"`
}

// GORM hooks
func (eo *ExternalOrder) BeforeCreate(tx *gorm.DB) (err error) {
	if eo.ID == uuid.Nil {
		eo.ID = uuid.New()
	}
	if eo.OrderNumber == "" {
		eo.OrderNumber = fmt.Sprintf("TRZ%s%06d",
			time.Now().Format("20060102"),
			time.Now().Unix()%1000000)
	}
	if eo.Currency == "" {
		eo.Currency = "INR"
	}
	return
}
