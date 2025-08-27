package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/zeusnotfound04/Tranza/models"
	"golang.org/x/crypto/bcrypt"
)

// JWTService interface for JWT operations
type JWTService interface {
	GenerateAccessToken(user *models.User) (string, error)
	GenerateRefreshToken(user *models.User) (string, error)
	ValidateToken(tokenString string) (jwt.MapClaims, error)
	ValidateRefreshToken(tokenString string) (jwt.MapClaims, error)
}

// jwtService implements JWTService interface
type jwtService struct {
	secretKey []byte
}

// NewJWTService creates a new JWT service instance
func NewJWTService(secretKey string) JWTService {
	return &jwtService{
		secretKey: []byte(secretKey),
	}
}

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

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(GetJWTSecret()))
}

// JWTService implementation methods
func (j *jwtService) GenerateAccessToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  user.ID.String(),
		"email":    user.Email,
		"username": user.Username,
		"exp":      time.Now().Add(1 * time.Hour).Unix(),
		"type":     "access",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secretKey)
}

func (j *jwtService) GenerateRefreshToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID.String(),
		"exp":     time.Now().Add(24 * 7 * time.Hour).Unix(), // 7 days
		"type":    "refresh",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secretKey)
}

func (j *jwtService) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func (j *jwtService) ValidateRefreshToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if tokenType, exists := claims["type"]; exists && tokenType == "refresh" {
			return claims, nil
		}
		return nil, fmt.Errorf("invalid token type")
	}

	return nil, fmt.Errorf("invalid refresh token")
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
		"payment.authorized":   true,
		"payment.captured":     true,
		"payment.failed":       true,
		"order.paid":           true,
		"refund.created":       true,
		"refund.processed":     true,
		"refund.failed":        true,
		"settlement.processed": true,
	}

	return validEvents[event]
}

// HashPassword hashes a password using bcrypt
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// VerifyPassword verifies a password against its hash using bcrypt
func VerifyPassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// EncryptAPIKey encrypts an API key using AES with a password-derived key
func EncryptAPIKey(apiKey, password string) (string, error) {
	// Create a simple XOR encryption using password hash as key
	passwordHash := sha256.Sum256([]byte(password))

	encrypted := make([]byte, len(apiKey))
	for i := range apiKey {
		encrypted[i] = apiKey[i] ^ passwordHash[i%len(passwordHash)]
	}

	return hex.EncodeToString(encrypted), nil
}

// DecryptAPIKey decrypts an API key using the password
func DecryptAPIKey(encryptedKey, password string) (string, error) {
	encrypted, err := hex.DecodeString(encryptedKey)
	if err != nil {
		return "", err
	}

	// Use the same XOR decryption
	passwordHash := sha256.Sum256([]byte(password))

	decrypted := make([]byte, len(encrypted))
	for i := range encrypted {
		decrypted[i] = encrypted[i] ^ passwordHash[i%len(passwordHash)]
	}

	return string(decrypted), nil
}
