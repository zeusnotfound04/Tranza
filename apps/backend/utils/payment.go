// utils/payment_utils.go
package utils

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

// Transaction Types
const (
	TransactionTypeLoadMoney  = "load_money"
	TransactionTypeAIPayment  = "ai_payment"
	TransactionTypeRefund     = "refund"
	TransactionTypeWithdrawal = "withdrawal"
)

// Transaction Status
const (
	TransactionStatusPending   = "pending"
	TransactionStatusSuccess   = "success"
	TransactionStatusFailed    = "failed"
	TransactionStatusCancelled = "cancelled"
)

// Payment Constants
const (
	// Razorpay Payment Status
	PaymentStatusCreated    = "created"
	PaymentStatusAuthorized = "authorized"
	PaymentStatusCaptured   = "captured"
	PaymentStatusRefunded   = "refunded"
	PaymentStatusFailed     = "failed"

	// Razorpay Order Status
	OrderStatusCreated   = "created"
	OrderStatusAttempted = "attempted"
	OrderStatusPaid      = "paid"

	// Payment Methods
	PaymentMethodCard       = "card"
	PaymentMethodNetbanking = "netbanking"
	PaymentMethodWallet     = "wallet"
	PaymentMethodUPI        = "upi"
	PaymentMethodEMI        = "emi"
	PaymentMethodCardless   = "cardless_emi"
	PaymentMethodPaylater   = "paylater"

	// Currency
	DefaultCurrency = "INR"

	// Amount Limits (in paise for Razorpay)
	MinAmountPaise = 100      // ₹1 in paise
	MaxAmountPaise = 10000000 // ₹1,00,000 in paise

	// Webhook Events
	WebhookPaymentCaptured   = "payment.captured"
	WebhookPaymentFailed     = "payment.failed"
	WebhookPaymentAuthorized = "payment.authorized"
	WebhookOrderPaid         = "order.paid"
	WebhookRefundCreated     = "refund.created"
	WebhookRefundProcessed   = "refund.processed"
)

// Logger instance
var Logger *logrus.Logger

// InitLogger initializes the logger if not already initialized
func InitLogger() {
	if Logger == nil {
		Logger = logrus.New()
		Logger.SetLevel(logrus.InfoLevel)
		Logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
		})
	}
}

// Logging Functions
func LogError(err error, context map[string]interface{}) {
	InitLogger()
	Logger.WithFields(context).Error(err)
}

func LogWarning(message string, context map[string]interface{}) {
	InitLogger()
	Logger.WithFields(context).Warn(message)
}

func LogInfo(message string, context map[string]interface{}) {
	InitLogger()
	Logger.WithFields(context).Info(message)
}

func LogTransaction(transactionID, userID, transactionType string, amount string, status string) {
	InitLogger()
	Logger.WithFields(logrus.Fields{
		"transaction_id":   transactionID,
		"user_id":         userID,
		"transaction_type": transactionType,
		"amount":          amount,
		"status":          status,
		"timestamp":       time.Now().UTC(),
	}).Info("Transaction processed")
}

func LogAIPayment(transactionID, userID, agentID, merchantName string, amount string, status string) {
	InitLogger()
	Logger.WithFields(logrus.Fields{
		"transaction_id": transactionID,
		"user_id":       userID,
		"agent_id":      agentID,
		"merchant_name": merchantName,
		"amount":        amount,
		"status":        status,
		"type":          "ai_payment",
		"timestamp":     time.Now().UTC(),
	}).Info("AI payment processed")
}



// Payment Validation Functions
func ValidateRazorpayPaymentID(paymentID string) error {
	if paymentID == "" {
		return errors.New("payment ID cannot be empty")
	}

	// Razorpay payment IDs start with "pay_"
	if !strings.HasPrefix(paymentID, "pay_") {
		return errors.New("invalid payment ID format")
	}

	// Should be exactly 18 characters
	if len(paymentID) != 18 {
		return errors.New("invalid payment ID length")
	}

	return nil
}

func ValidateRazorpayOrderID(orderID string) error {
	if orderID == "" {
		return errors.New("order ID cannot be empty")
	}

	// Razorpay order IDs start with "order_"
	if !strings.HasPrefix(orderID, "order_") {
		return errors.New("invalid order ID format")
	}

	// Should be exactly 20 characters
	if len(orderID) != 20 {
		return errors.New("invalid order ID length")
	}

	return nil
}

func ValidateRazorpaySignature(signature string) error {
	if signature == "" {
		return errors.New("signature cannot be empty")
	}

	// Should be a valid hex string of 64 characters (SHA256)
	if len(signature) != 64 {
		return errors.New("invalid signature length")
	}

	matched, _ := regexp.MatchString("^[a-f0-9]{64}$", signature)
	if !matched {
		return errors.New("invalid signature format")
	}

	return nil
}

func ValidatePaymentAmount(amount decimal.Decimal) error {
	if amount.LessThanOrEqual(decimal.Zero) {
		return errors.New("amount must be greater than zero")
	}

	// Convert to paise for Razorpay validation
	amountInPaise := amount.Mul(decimal.NewFromInt(100))

	if amountInPaise.LessThan(decimal.NewFromInt(MinAmountPaise)) {
		return fmt.Errorf("minimum amount is ₹%.2f", decimal.NewFromInt(MinAmountPaise).Div(decimal.NewFromInt(100)))
	}

	if amountInPaise.GreaterThan(decimal.NewFromInt(MaxAmountPaise)) {
		return fmt.Errorf("maximum amount is ₹%.2f", decimal.NewFromInt(MaxAmountPaise).Div(decimal.NewFromInt(100)))
	}

	return nil
}

// Signature Generation Functions
func GeneratePaymentSignature(orderID, paymentID, secret string) string {
	message := orderID + "|" + paymentID
	return generateHMAC(message, secret)
}

func VerifyPaymentSignature(orderID, paymentID, signature, secret string) bool {
	expectedSignature := GeneratePaymentSignature(orderID, paymentID, secret)
	return hmac.Equal([]byte(expectedSignature), []byte(signature))
}

func GenerateWebhookSignature(body []byte, secret string) string {
	return generateHMAC(string(body), secret)
}


func generateHMAC(message, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))
}

// Amount Conversion Functions
func RupeesToPaise(rupees decimal.Decimal) int64 {
	return rupees.Mul(decimal.NewFromInt(100)).IntPart()
}

func PaiseToRupees(paise int64) decimal.Decimal {
	return decimal.NewFromInt(paise).Div(decimal.NewFromInt(100))
}



// Payment Method Validation
func ValidatePaymentMethod(method string) error {
	supportedMethods := []string{
		PaymentMethodCard,
		PaymentMethodNetbanking,
		PaymentMethodWallet,
		PaymentMethodUPI,
		PaymentMethodEMI,
		PaymentMethodCardless,
		PaymentMethodPaylater,
	}

	for _, supported := range supportedMethods {
		if method == supported {
			return nil
		}
	}

	return fmt.Errorf("unsupported payment method: %s", method)
}

func GetPaymentMethodDisplayName(method string) string {
	displayNames := map[string]string{
		PaymentMethodCard:       "Credit/Debit Card",
		PaymentMethodNetbanking: "Net Banking",
		PaymentMethodWallet:     "Wallet",
		PaymentMethodUPI:        "UPI",
		PaymentMethodEMI:        "EMI",
		PaymentMethodCardless:   "Cardless EMI",
		PaymentMethodPaylater:   "Pay Later",
	}

	if displayName, exists := displayNames[method]; exists {
		return displayName
	}

	return strings.Title(method)
}

// Payment Status Functions
func IsPaymentSuccessful(status string) bool {
	return status == PaymentStatusCaptured
}

func IsPaymentPending(status string) bool {
	return status == PaymentStatusCreated || status == PaymentStatusAuthorized
}

func IsPaymentFailed(status string) bool {
	return status == PaymentStatusFailed
}

func GetPaymentStatusDisplayName(status string) string {
	displayNames := map[string]string{
		PaymentStatusCreated:    "Payment Initiated",
		PaymentStatusAuthorized: "Payment Authorized",
		PaymentStatusCaptured:   "Payment Successful",
		PaymentStatusRefunded:   "Payment Refunded",
		PaymentStatusFailed:     "Payment Failed",
	}

	if displayName, exists := displayNames[status]; exists {
		return displayName
	}

	return strings.Title(status)
}

// Reference Generation Functions
func GenerateRandomString(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		randomIndex := make([]byte, 1)
		if _, err := rand.Read(randomIndex); err != nil {
			return "", err
		}
		b[i] = charset[int(randomIndex[0])%len(charset)]
	}
	return string(b), nil
}

func GenerateTransactionReference(prefix string) string {
	timestamp := time.Now().Unix()
	randomPart, _ := GenerateRandomString(8)
	return fmt.Sprintf("%s_%d_%s", prefix, timestamp, randomPart)
}

func GenerateOrderReceipt(userID, purpose string) string {
	timestamp := time.Now().Unix()
	return fmt.Sprintf("%s_%s_%d", purpose, userID[:8], timestamp)
}

func GenerateWebhookID(eventType string) string {
	timestamp := time.Now().Unix()
	randomPart, _ := GenerateRandomString(8)
	return fmt.Sprintf("webhook_%s_%d_%s", eventType, timestamp, randomPart)
}

// Error Handling Functions
func GetUserFriendlyError(err error) string {
	if err == nil {
		return ""
	}

	errorMessage := err.Error()

	// Map common errors to user-friendly messages
	errorMappings := map[string]string{
		"invalid payment signature":       "Payment verification failed. Please try again.",
		"payment not captured":            "Payment was not processed successfully. Please try again.",
		"insufficient balance":            "Insufficient balance in your account.",
		"wallet not found":               "Wallet not found. Please contact support.",
		"transaction already processed":   "This transaction has already been processed.",
		"amount mismatch":                "Payment amount doesn't match. Please try again.",
		"invalid order ID format":        "Invalid payment details. Please try again.",
		"invalid payment ID format":      "Invalid payment details. Please try again.",
		"transaction not found":          "Transaction not found. Please contact support.",
		"failed to create payment order": "Unable to initiate payment. Please try again.",
	}

	// Check for exact matches first
	if friendlyMessage, exists := errorMappings[errorMessage]; exists {
		return friendlyMessage
	}

	// Check for partial matches
	for technicalError, friendlyMessage := range errorMappings {
		if strings.Contains(errorMessage, technicalError) {
			return friendlyMessage
		}
	}

	// Default message for unknown errors
	return "An error occurred while processing your payment. Please try again."
}

func GetErrorCode(err error) string {
	if err == nil {
		return ""
	}

	errorMessage := err.Error()

	errorCodes := map[string]string{
		"invalid payment signature":       "PAYMENT_SIGNATURE_INVALID",
		"payment not captured":            "PAYMENT_NOT_CAPTURED",
		"insufficient balance":            "INSUFFICIENT_BALANCE",
		"wallet not found":               "WALLET_NOT_FOUND",
		"transaction already processed":   "DUPLICATE_TRANSACTION",
		"amount mismatch":                "AMOUNT_MISMATCH",
		"invalid order ID format":        "INVALID_ORDER_ID",
		"invalid payment ID format":      "INVALID_PAYMENT_ID",
		"transaction not found":          "TRANSACTION_NOT_FOUND",
		"failed to create payment order": "ORDER_CREATION_FAILED",
	}

	for errorText, errorCode := range errorCodes {
		if strings.Contains(errorMessage, errorText) {
			return errorCode
		}
	}

	return "UNKNOWN_PAYMENT_ERROR"
}

// Retry Logic Functions
func ShouldRetryPayment(err error) bool {
	if err == nil {
		return false
	}

	// Errors that should not be retried
	nonRetryableErrors := []string{
		"invalid payment signature",
		"transaction already processed",
		"amount mismatch",
		"wallet not found",
		"insufficient balance",
	}

	errorMessage := err.Error()
	for _, nonRetryable := range nonRetryableErrors {
		if strings.Contains(errorMessage, nonRetryable) {
			return false
		}
	}

	// Default to allowing retry for network/temporary errors
	return true
}

func GetRetryDelay(attemptNumber int) time.Duration {
	// Exponential backoff: 1s, 2s, 4s, 8s, 16s
	delay := time.Duration(1<<attemptNumber) * time.Second
	
	// Cap at 30 seconds
	if delay > 30*time.Second {
		delay = 30 * time.Second
	}
	
	return delay
}

func GetMaxRetryAttempts() int {
	return 3
}

