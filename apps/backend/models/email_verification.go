package models

import (
	"time"
)


type EmailVerification struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Email      string    `gorm:"unique;not null;index" json:"email"`
	Username   string    `gorm:"not null" json:"username"`
	Password   string    `gorm:"not null" json:"-"` // Hashed password
	Code       string    `gorm:"not null" json:"-"` // Verification code (hashed)
	ExpiresAt  time.Time `gorm:"not null;index" json:"expires_at"`
	Attempts   int       `gorm:"default:0" json:"attempts"`
	IsVerified bool      `gorm:"default:false" json:"is_verified"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type PreRegistrationRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=8"`
}

type EmailVerificationRequest struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required,len=6"`
}

type ResendVerificationRequest struct {
	Email string `json:"email" binding:"required,email"`
}


type PreRegistrationResponse struct {
	Message   string    `json:"message"`
	Email     string    `json:"email"`
	ExpiresAt time.Time `json:"expires_at"`
}

// EmailVerificationResponse after successful verification
type EmailVerificationResponse struct {
	Message string `json:"message"`
	User    *User  `json:"user"`
}

// TableName returns the table name for EmailVerification
func (EmailVerification) TableName() string {
	return "email_verifications"
}