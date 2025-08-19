package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Address struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID      uuid.UUID `json:"user_id" gorm:"type:uuid;not null;index"`
	Name        string    `json:"name" gorm:"not null"`
	Phone       string    `json:"phone" gorm:"not null"`
	AddressLine string    `json:"address_line" gorm:"not null"`
	City        string    `json:"city" gorm:"not null"`
	State       string    `json:"state" gorm:"not null"`
	PinCode     string    `json:"pin_code" gorm:"not null"`
	Country     string    `json:"country" gorm:"default:'India'"`
	Landmark    string    `json:"landmark,omitempty"`
	IsDefault   bool      `json:"is_default" gorm:"default:false"`
	AddressType string    `json:"address_type" gorm:"default:'home'"` // home, office, other
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relations
	User User `json:"-" gorm:"foreignKey:UserID"`
}

// DTO for address operations
type AddressCreateRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=100"`
	Phone       string `json:"phone" binding:"required,min=10,max=15"`
	AddressLine string `json:"address_line" binding:"required,min=10,max=500"`
	City        string `json:"city" binding:"required,min=2,max=100"`
	State       string `json:"state" binding:"required,min=2,max=100"`
	PinCode     string `json:"pin_code" binding:"required,min=6,max=6"`
	Country     string `json:"country"`
	Landmark    string `json:"landmark,omitempty"`
	IsDefault   bool   `json:"is_default"`
	AddressType string `json:"address_type"`
}

type AddressUpdateRequest struct {
	Name        *string `json:"name,omitempty" binding:"omitempty,min=2,max=100"`
	Phone       *string `json:"phone,omitempty" binding:"omitempty,min=10,max=15"`
	AddressLine *string `json:"address_line,omitempty" binding:"omitempty,min=10,max=500"`
	City        *string `json:"city,omitempty" binding:"omitempty,min=2,max=100"`
	State       *string `json:"state,omitempty" binding:"omitempty,min=2,max=100"`
	PinCode     *string `json:"pin_code,omitempty" binding:"omitempty,min=6,max=6"`
	Country     *string `json:"country,omitempty"`
	Landmark    *string `json:"landmark,omitempty"`
	IsDefault   *bool   `json:"is_default,omitempty"`
	AddressType *string `json:"address_type,omitempty"`
}

type AddressResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Phone       string `json:"phone"`
	AddressLine string `json:"address_line"`
	City        string `json:"city"`
	State       string `json:"state"`
	PinCode     string `json:"pin_code"`
	Country     string `json:"country"`
	Landmark    string `json:"landmark,omitempty"`
	IsDefault   bool   `json:"is_default"`
	AddressType string `json:"address_type"`
	CreatedAt   string `json:"created_at"`
}

// BeforeCreate hook
func (a *Address) BeforeCreate(tx *gorm.DB) (err error) {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	if a.Country == "" {
		a.Country = "India"
	}
	if a.AddressType == "" {
		a.AddressType = "home"
	}
	return
}

// BeforeUpdate hook - ensure only one default address per user
func (a *Address) BeforeUpdate(tx *gorm.DB) (err error) {
	if a.IsDefault {
		// Set all other addresses as non-default for this user
		err = tx.Model(&Address{}).Where("user_id = ? AND id != ?", a.UserID, a.ID).Update("is_default", false).Error
	}
	return
}

// BeforeCreate hook - ensure only one default address per user
func (a *Address) AfterCreate(tx *gorm.DB) (err error) {
	if a.IsDefault {
		// Set all other addresses as non-default for this user
		err = tx.Model(&Address{}).Where("user_id = ? AND id != ?", a.UserID, a.ID).Update("is_default", false).Error
	}
	return
}

// Helper method to format address
func (a *Address) GetFormattedAddress() string {
	formatted := a.AddressLine + ", " + a.City + ", " + a.State + " - " + a.PinCode
	if a.Landmark != "" {
		formatted = a.AddressLine + ", " + a.Landmark + ", " + a.City + ", " + a.State + " - " + a.PinCode
	}
	return formatted
}
