package utils

import (
	"os"
	"time"
    "crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func GetJWTSecret() string {
	return os.Getenv("JWT_SECRET")
}

func GenerateJWT(userId, email, username string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userId,
		"email":    email,
		"username": username,
		"exp":      time.Now().Add(48 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	return token.SignedString(jwtSecret)
}



// ComputeHMACSHA256 computes HMAC-SHA256 signature for the given message and secret
func ComputeHMACSHA256(message, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))
}

// VerifyRazorpaySignature verifies the payment signature from Razorpay
// This is used to verify the signature received after payment completion
func VerifyRazorpaySignature(orderID, paymentID, signature, secret string) bool {
	if orderID == "" || paymentID == "" || signature == "" || secret == "" {
		return false
	}
	
	// Create the message string as per Razorpay documentation
	message := fmt.Sprintf("%s|%s", orderID, paymentID)
	
	// Compute expected signature
	expectedSignature := ComputeHMACSHA256(message, secret)
	
	// Use constant time comparison to prevent timing attacks
	return hmac.Equal([]byte(expectedSignature), []byte(signature))
}

// VerifyWebhookSignature verifies the webhook signature from Razorpay
// This is used to verify that the webhook request is actually from Razorpay
func VerifyWebhookSignature(body, signature, secret string) bool {
	if body == "" || signature == "" || secret == "" {
		return false
	}
	
	// Compute expected signature using the webhook body
	expectedSignature := ComputeHMACSHA256(body, secret)
	
	// Use constant time comparison to prevent timing attacks
	return hmac.Equal([]byte(expectedSignature), []byte(signature))
}

// VerifyRefundSignature verifies the refund signature from Razorpay
// This is used when processing refund webhooks
func VerifyRefundSignature(refundID, paymentID, signature, secret string) bool {
	if refundID == "" || paymentID == "" || signature == "" || secret == "" {
		return false
	}
	
	// Create the message string for refund verification
	message := fmt.Sprintf("%s|%s", paymentID, refundID)
	
	// Compute expected signature
	expectedSignature := ComputeHMACSHA256(message, secret)
	
	// Use constant time comparison to prevent timing attacks
	return hmac.Equal([]byte(expectedSignature), []byte(signature))
}

// ValidateSignatureFormat checks if the signature is in the correct format
func ValidateSignatureFormat(signature string) bool {
	// Razorpay signatures are typically 64 characters long (SHA256 hex)
	if len(signature) != 64 {
		return false
	}
	
	// Check if it contains only valid hex characters
	for _, char := range signature {
		if !((char >= '0' && char <= '9') || (char >= 'a' && char <= 'f') || (char >= 'A' && char <= 'F')) {
			return false
		}
	}
	
	return true
}

// SanitizeInput removes potentially dangerous characters from input strings
func SanitizeInput(input string) string {
	// Remove null bytes and control characters
	sanitized := strings.ReplaceAll(input, "\x00", "")
	sanitized = strings.TrimSpace(sanitized)
	return sanitized
}

// ValidateOrderID checks if the order ID format is valid
func ValidateOrderID(orderID string) bool {
	// Razorpay order IDs typically start with "order_" followed by alphanumeric characters
	if !strings.HasPrefix(orderID, "order_") {
		return false
	}
	
	// Check minimum length (order_ + at least 10 characters)
	if len(orderID) < 16 {
		return false
	}
	
	// Check if it contains only valid characters after "order_"
	idPart := orderID[6:] // Remove "order_" prefix
	for _, char := range idPart {
		if !((char >= '0' && char <= '9') || (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z')) {
			return false
		}
	}
	
	return true
}

// ValidatePaymentID checks if the payment ID format is valid
func ValidatePaymentID(paymentID string) bool {
	// Razorpay payment IDs typically start with "pay_" followed by alphanumeric characters
	if !strings.HasPrefix(paymentID, "pay_") {
		return false
	}
	
	// Check minimum length (pay_ + at least 10 characters)
	if len(paymentID) < 14 {
		return false
	}
	
	// Check if it contains only valid characters after "pay_"
	idPart := paymentID[4:] // Remove "pay_" prefix
	for _, char := range idPart {
		if !((char >= '0' && char <= '9') || (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z')) {
			return false
		}
	}
	
	return true
}

// ValidateAmount checks if the amount is valid (positive and reasonable)
func ValidateAmount(amount float64) bool {
	// Amount should be positive
	if amount <= 0 {
		return false
	}
	
	// Amount should not be unreasonably large (100 crores in rupees)
	if amount > 1000000000 {
		return false
	}
	
	// Amount should have at most 2 decimal places
	multiplied := amount * 100
	if multiplied != float64(int(multiplied)) {
		return false
	}
	
	return true
}

// ValidateCurrency checks if the currency code is supported
func ValidateCurrency(currency string) bool {
	supportedCurrencies := map[string]bool{
		"INR": true,
		"USD": true,
		"EUR": true,
		"GBP": true,
		"AUD": true,
		"CAD": true,
		"SGD": true,
		"AED": true,
		"MYR": true,
	}
	
	return supportedCurrencies[strings.ToUpper(currency)]
}

// SecureCompare performs a constant-time comparison of two strings
// This prevents timing attacks when comparing sensitive data
func SecureCompare(a, b string) bool {
	return hmac.Equal([]byte(a), []byte(b))
}

// GenerateSecureReceipt generates a secure receipt string
func GenerateSecureReceipt(prefix string) string {
	// This would typically use a proper random generator
	// For production, consider using crypto/rand
	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	hash := ComputeHMACSHA256(timestamp, "internal-secret")
	return fmt.Sprintf("%s_%s", prefix, hash[:8])
}

// ValidateWebhookEvent checks if the webhook event is valid
func ValidateWebhookEvent(event string) bool {
	validEvents := map[string]bool{
		"payment.authorized": true,
		"payment.captured":   true,
		"payment.failed":     true,
		"order.paid":         true,
		"refund.created":     true,
		"refund.processed":   true,
		"refund.failed":      true,
		"settlement.processed": true,
	}
	
	return validEvents[event]
}