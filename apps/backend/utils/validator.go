package utils

import (
	"errors"
	"regexp"
	"strings"

	"github.com/shopspring/decimal"
)

// ValidateAmount validates transaction amount
func ValidateAmount(amount decimal.Decimal) error {
	if amount.LessThanOrEqual(decimal.Zero) {
		return errors.New("amount must be greater than 0")
	}

	if amount.LessThan(decimal.NewFromInt(1)) {
		return errors.New("minimum amount is ₹1")
	}

	if amount.GreaterThan(decimal.NewFromInt(100000)) {
		return errors.New("maximum amount is ₹1,00,000")
	}

	return nil
}

// ValidateLoadAmount validates wallet loading amount
func ValidateLoadAmount(amount decimal.Decimal) error {
	if amount.LessThanOrEqual(decimal.Zero) {
		return errors.New("amount must be greater than 0")
	}

	if amount.LessThan(decimal.NewFromInt(10)) {
		return errors.New("minimum load amount is ₹10")
	}

	if amount.GreaterThan(decimal.NewFromInt(50000)) {
		return errors.New("maximum load amount is ₹50,000")
	}

	return nil
}

// ValidatePhoneNumber validates Indian phone number
func ValidatePhoneNumber(phone string) error {
	// Remove spaces and special characters
	phone = strings.ReplaceAll(phone, " ", "")
	phone = strings.ReplaceAll(phone, "-", "")
	phone = strings.ReplaceAll(phone, "+91", "")

	// Check if it's exactly 10 digits
	matched, _ := regexp.MatchString(`^[6-9]\d{9}$`, phone)
	if !matched {
		return errors.New("invalid phone number format")
	}

	return nil
}

// ValidateEmail validates email format
func ValidateEmail(email string) error {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if !emailRegex.MatchString(strings.ToLower(email)) {
		return errors.New("invalid email format")
	}
	return nil
}

// ValidateUPIID validates UPI ID format
func ValidateUPIID(upiID string) error {
	if upiID == "" {
		return nil // UPI ID is optional
	}

	upiRegex := regexp.MustCompile(`^[\w\.\-_]{2,256}@[a-zA-Z]{2,64}$`)
	if !upiRegex.MatchString(upiID) {
		return errors.New("invalid UPI ID format")
	}

	return nil
}

// ValidateAgentID validates AI agent ID
func ValidateAgentID(agentID string) error {
	if len(agentID) < 3 || len(agentID) > 100 {
		return errors.New("agent ID must be between 3 and 100 characters")
	}

	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_\-]+$`, agentID)
	if !matched {
		return errors.New("agent ID can only contain letters, numbers, underscores, and hyphens")
	}

	return nil
}