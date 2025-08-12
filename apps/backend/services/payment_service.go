package services

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/zeusnotfound04/Tranza/models"
	"github.com/zeusnotfound04/Tranza/models/dto"
	"github.com/zeusnotfound04/Tranza/pkg/razorpay"
	"github.com/zeusnotfound04/Tranza/repositories"
	"github.com/zeusnotfound04/Tranza/utils"
	"gorm.io/gorm"
)

type PaymentService struct {
	razorpayClient  *razorpay.Client
	walletRepo      *repositories.WalletRepository
	transactionRepo *repositories.TransactionRepository
	notificationSvc *NotificationService
	db              *gorm.DB
	webhookSecret   string
}

func NewPaymentService(
	razorpayClient *razorpay.Client,
	walletRepo *repositories.WalletRepository,
	transactionRepo *repositories.TransactionRepository,
	notificationSvc *NotificationService,
	db *gorm.DB,
	webhookSecret string,
) *PaymentService {
	return &PaymentService{
		razorpayClient:  razorpayClient,
		walletRepo:      walletRepo,
		transactionRepo: transactionRepo,
		notificationSvc: notificationSvc,
		db:              db,
		webhookSecret:   webhookSecret,
	}
}

// CreateLoadMoneyOrder creates Razorpay order for loading money
func (s *PaymentService) CreateLoadMoneyOrder(userID string, amount decimal.Decimal) (*dto.LoadMoneyResponse, error) {
	// Validate amount
	if err := utils.ValidateLoadAmount(amount); err != nil {
		return nil, err
	}

	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	// Get wallet
	wallet, err := s.walletRepo.GetByUserID(uid)
	if err != nil {
		return nil, errors.New("wallet not found")
	}

	if wallet.Status != "active" {
		return nil, errors.New("wallet is not active")
	}

	// Create Razorpay order
	amountInPaise := amount.Mul(decimal.NewFromInt(100)).IntPart()
	receipt := utils.GenerateTransactionReference("WALLET_LOAD")

	order, err := s.razorpayClient.CreateOrder(amountInPaise, "INR", receipt)
	if err != nil {
		utils.LogError(err, map[string]interface{}{
			"user_id": userID,
			"amount":  amount.String(),
			"action":  "create_razorpay_order",
		})
		return nil, fmt.Errorf("failed to create payment order: %w", err)
	}

	// Create pending transaction
	transaction := &models.Transaction{
		WalletID:        wallet.ID,
		UserID:          uid,
		Type:            utils.TransactionTypeLoadMoney,
		Amount:          amount,
		Currency:        "INR",
		Description:     "Loading money to wallet",
		RazorpayOrderID: order.ID,
		Status:          utils.TransactionStatusPending,
		ReferenceID:     receipt,
	}

	createdTxn, err := s.transactionRepo.Create(transaction)
	if err != nil {
		utils.LogError(err, map[string]interface{}{
			"user_id":  userID,
			"order_id": order.ID,
			"action":   "create_pending_transaction",
		})
		return nil, fmt.Errorf("failed to create transaction record: %w", err)
	}

	utils.LogInfo("Load money order created", map[string]interface{}{
		"user_id":        userID,
		"order_id":       order.ID,
		"transaction_id": createdTxn.ID.String(),
		"amount":         amount.String(),
	})

	return &dto.LoadMoneyResponse{
		OrderID:       order.ID,
		Amount:        amount,
		Currency:      "INR",
		TransactionID: createdTxn.ID.String(),
		RazorpayKeyID: s.razorpayClient.KeyID,
	}, nil
}

// VerifyAndProcessPayment verifies Razorpay payment and processes wallet credit
func (s *PaymentService) VerifyAndProcessPayment(
	userID, paymentID, orderID, signature string,
) (*dto.PaymentVerificationResponse, error) {
	// Verify Razorpay signature
	if !s.verifyPaymentSignature(orderID, paymentID, signature) {
		utils.LogWarning("Invalid payment signature", map[string]interface{}{
			"user_id":    userID,
			"payment_id": paymentID,
			"order_id":   orderID,
		})
		return nil, errors.New("invalid payment signature")
	}

	// Get payment details from Razorpay
	payment, err := s.razorpayClient.GetPayment(paymentID)
	if err != nil {
		utils.LogError(err, map[string]interface{}{
			"payment_id": paymentID,
			"action":     "get_payment_details",
		})
		return nil, fmt.Errorf("failed to get payment details: %w", err)
	}

	// Validate payment status
	if payment.Status != "captured" {
		return nil, fmt.Errorf("payment not captured, status: %s", payment.Status)
	}

	// Convert amount from paise to rupees
	amount := decimal.NewFromInt(payment.Amount).Div(decimal.NewFromInt(100))

	// Process payment in database transaction
	result, err := s.processPaymentTransaction(userID, paymentID, orderID, amount, payment.Method)
	if err != nil {
		return nil, err
	}

	// Send success notification
	go s.notificationSvc.SendWalletCreditNotification(userID, amount, result.NewBalance)

	utils.LogTransaction(
		result.TransactionID,
		userID,
		utils.TransactionTypeLoadMoney,
		amount.String(),
		utils.TransactionStatusSuccess,
	)

	return result, nil
}

// processPaymentTransaction handles the database transaction for payment processing
func (s *PaymentService) processPaymentTransaction(
	userID, paymentID, orderID string,
	amount decimal.Decimal,
	paymentMethod string,
) (*dto.PaymentVerificationResponse, error) {
	// Start database transaction
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Get pending transaction
	transaction, err := s.transactionRepo.GetByOrderID(orderID)
	if err != nil {
		tx.Rollback()
		return nil, errors.New("transaction not found")
	}

	// Validate transaction
	if transaction.Status != utils.TransactionStatusPending {
		tx.Rollback()
		return nil, errors.New("transaction already processed")
	}

	if transaction.UserID.String() != userID {
		tx.Rollback()
		return nil, errors.New("unauthorized access to transaction")
	}

	// Validate amount matches
	if !transaction.Amount.Equal(amount) {
		tx.Rollback()
		return nil, errors.New("amount mismatch")
	}

	// Get wallet
	wallet, err := s.walletRepo.GetByID(transaction.WalletID)
	if err != nil {
		tx.Rollback()
		return nil, errors.New("wallet not found")
	}

	// Calculate new balance
	newBalance := wallet.Balance.Add(amount)

	// Update wallet balance
	err = s.walletRepo.UpdateBalance(tx, wallet.ID, newBalance)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update wallet balance: %w", err)
	}

	// Update transaction
	transaction.Status = utils.TransactionStatusSuccess
	transaction.RazorpayPaymentID = paymentID
	transaction.BalanceAfter = newBalance
	transaction.PaymentMethod = paymentMethod

	err = s.transactionRepo.Update(tx, transaction)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update transaction: %w", err)
	}

	// Commit transaction
	if err = tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit payment transaction: %w", err)
	}

	return &dto.PaymentVerificationResponse{
		Success:       true,
		NewBalance:    newBalance,
		TransactionID: transaction.ID.String(),
		Message:       fmt.Sprintf("₹%s added to wallet successfully", amount.StringFixed(2)),
		Amount:        amount,
	}, nil
}

// ProcessWebhookEvent processes Razorpay webhook events
func (s *PaymentService) ProcessWebhookEvent(body []byte, signature string) error {
	// Verify webhook signature
	if !s.verifyWebhookSignature(body, signature) {
		utils.LogWarning("Invalid webhook signature", map[string]interface{}{
			"signature": signature,
		})
		return errors.New("invalid webhook signature")
	}

	// Parse webhook event
	event, err := s.razorpayClient.ParseWebhookEvent(body)
	if err != nil {
		utils.LogError(err, map[string]interface{}{
			"action": "parse_webhook_event",
		})
		return fmt.Errorf("failed to parse webhook event: %w", err)
	}

	utils.LogInfo("Webhook event received", map[string]interface{}{
		"event":   event.Event,
		"account": event.Account,
	})

	// Process based on event type
	switch event.Event {
	case "payment.captured":
		return s.handlePaymentCaptured(event)
	case "payment.failed":
		return s.handlePaymentFailed(event)
	case "payment.authorized":
		return s.handlePaymentAuthorized(event)
	case "order.paid":
		return s.handleOrderPaid(event)
	default:
		utils.LogInfo("Unhandled webhook event", map[string]interface{}{
			"event": event.Event,
		})
		return nil // Don't return error for unhandled events
	}
}

// handlePaymentCaptured handles payment.captured webhook
func (s *PaymentService) handlePaymentCaptured(event *razorpay.WebhookEvent) error {
	paymentData, ok := event.Payload["payment"].(map[string]interface{})
	if !ok {
		return errors.New("invalid payment data in webhook")
	}

	paymentID, ok := paymentData["id"].(string)
	if !ok {
		return errors.New("payment ID not found in webhook")
	}

	orderID, ok := paymentData["order_id"].(string)
	if !ok {
		return errors.New("order ID not found in webhook")
	}

	// Check if we already processed this payment
	_, err := s.transactionRepo.GetByPaymentID(paymentID)
	if err == nil {
		// Payment already processed
		return nil
	}

	// Get the pending transaction
	transaction, err := s.transactionRepo.GetByOrderID(orderID)
	if err != nil {
		return fmt.Errorf("transaction not found for order %s: %w", orderID, err)
	}

	if transaction.Status != utils.TransactionStatusPending {
		// Already processed
		return nil
	}

	// Update transaction status
	err = s.transactionRepo.UpdateStatus(
		transaction.ID,
		utils.TransactionStatusSuccess,
		"",
	)
	if err != nil {
		return fmt.Errorf("failed to update transaction status: %w", err)
	}

	utils.LogInfo("Payment captured via webhook", map[string]interface{}{
		"payment_id":     paymentID,
		"order_id":       orderID,
		"transaction_id": transaction.ID.String(),
	})

	return nil
}

// handlePaymentFailed handles payment.failed webhook
func (s *PaymentService) handlePaymentFailed(event *razorpay.WebhookEvent) error {
	paymentData, ok := event.Payload["payment"].(map[string]interface{})
	if !ok {
		return errors.New("invalid payment data in webhook")
	}

	orderID, ok := paymentData["order_id"].(string)
	if !ok {
		return errors.New("order ID not found in webhook")
	}

	// Get the pending transaction
	transaction, err := s.transactionRepo.GetByOrderID(orderID)
	if err != nil {
		return fmt.Errorf("transaction not found for order %s: %w", orderID, err)
	}

	if transaction.Status != utils.TransactionStatusPending {
		// Already processed
		return nil
	}

	// Get failure reason
	failureReason := "Payment failed"
	if errorData, ok := paymentData["error_description"].(string); ok {
		failureReason = errorData
	}

	// Update transaction status
	err = s.transactionRepo.UpdateStatus(
		transaction.ID,
		utils.TransactionStatusFailed,
		failureReason,
	)
	if err != nil {
		return fmt.Errorf("failed to update transaction status: %w", err)
	}

	utils.LogInfo("Payment failed via webhook", map[string]interface{}{
		"order_id":       orderID,
		"transaction_id": transaction.ID.String(),
		"reason":         failureReason,
	})

	return nil
}

// handlePaymentAuthorized handles payment.authorized webhook
func (s *PaymentService) handlePaymentAuthorized(event *razorpay.WebhookEvent) error {
	// For automatic capture, this might not be needed
	// But useful for manual capture scenarios
	utils.LogInfo("Payment authorized", map[string]interface{}{
		"event": event.Event,
	})
	return nil
}

// handleOrderPaid handles order.paid webhook
func (s *PaymentService) handleOrderPaid(event *razorpay.WebhookEvent) error {
	orderData, ok := event.Payload["order"].(map[string]interface{})
	if !ok {
		return errors.New("invalid order data in webhook")
	}

	orderID, ok := orderData["id"].(string)
	if !ok {
		return errors.New("order ID not found in webhook")
	}

	utils.LogInfo("Order paid", map[string]interface{}{
		"order_id": orderID,
	})

	return nil
}

// GetPaymentStatus gets payment status from Razorpay
func (s *PaymentService) GetPaymentStatus(paymentID string) (*razorpay.Payment, error) {
	payment, err := s.razorpayClient.GetPayment(paymentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment status: %w", err)
	}
	return payment, nil
}

// GetOrderStatus gets order status from Razorpay
func (s *PaymentService) GetOrderStatus(orderID string) (*razorpay.Order, error) {
	order, err := s.razorpayClient.GetOrder(orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order status: %w", err)
	}
	return order, nil
}

// RefundPayment processes a refund (for future use)
func (s *PaymentService) RefundPayment(paymentID string, amount decimal.Decimal, reason string) error {
	// This would use Razorpay's refund API
	// Implementation depends on your refund requirements
	utils.LogInfo("Refund requested", map[string]interface{}{
		"payment_id": paymentID,
		"amount":     amount.String(),
		"reason":     reason,
	})

	// TODO: Implement refund logic using Razorpay refund API
	return errors.New("refund functionality not implemented yet")
}

// CleanupExpiredOrders marks expired pending orders as failed
func (s *PaymentService) CleanupExpiredOrders() error {
	// Get pending transactions older than 1 hour
	expiredTransactions, err := s.transactionRepo.GetPendingTransactions(1 * time.Hour)
	if err != nil {
		return fmt.Errorf("failed to get expired transactions: %w", err)
	}

	if len(expiredTransactions) == 0 {
		return nil
	}

	// Update expired transactions to failed
	expiredIDs := make([]uuid.UUID, len(expiredTransactions))
	for i, txn := range expiredTransactions {
		expiredIDs[i] = txn.ID
	}

	err = s.transactionRepo.BulkUpdateStatus(expiredIDs, utils.TransactionStatusFailed)
	if err != nil {
		return fmt.Errorf("failed to update expired transactions: %w", err)
	}

	utils.LogInfo("Expired orders cleaned up", map[string]interface{}{
		"count": len(expiredTransactions),
	})

	return nil
}

// verifyPaymentSignature verifies Razorpay payment signature
func (s *PaymentService) verifyPaymentSignature(orderID, paymentID, signature string) bool {
	body := orderID + "|" + paymentID
	expectedSignature := s.generateSignature(body, s.razorpayClient.KeySecret)
	return hmac.Equal([]byte(expectedSignature), []byte(signature))
}

// verifyWebhookSignature verifies Razorpay webhook signature
func (s *PaymentService) verifyWebhookSignature(body []byte, signature string) bool {
	expectedSignature := s.generateSignature(string(body), s.webhookSecret)
	return hmac.Equal([]byte(expectedSignature), []byte(signature))
}

// generateSignature generates HMAC-SHA256 signature
func (s *PaymentService) generateSignature(body, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(body))
	return hex.EncodeToString(h.Sum(nil))
}

// GetPaymentAnalytics returns payment analytics
func (s *PaymentService) GetPaymentAnalytics(userID string, days int) (*PaymentAnalytics, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	// Get date range
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -days)

	// Get wallet
	wallet, err := s.walletRepo.GetByUserID(uid)
	if err != nil {
		return nil, errors.New("wallet not found")
	}

	// Get transaction summary
	summary, err := s.transactionRepo.GetTransactionSummaryByDateRange(wallet.ID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction summary: %w", err)
	}

	// Get transaction stats
	stats, err := s.transactionRepo.GetTransactionStats(uid)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction stats: %w", err)
	}

	return &PaymentAnalytics{
		DateRange:         fmt.Sprintf("%s to %s", startDate.Format("2006-01-02"), endDate.Format("2006-01-02")),
		TotalTransactions: summary.TotalTransactions,
		TotalInflow:       summary.TotalInflow,
		TotalOutflow:      summary.TotalOutflow,
		NetFlow:           summary.NetFlow,
		AverageAmount:     summary.AverageAmount,
		TodayTransactions: stats.TodayTransactions,
		TodayAmount:       stats.TodayAmount,
		AITransactions:    stats.AITransactions,
		AIAmount:          stats.AIAmount,
	}, nil
}

// PaymentAnalytics holds payment analytics data
type PaymentAnalytics struct {
	DateRange         string          `json:"date_range"`
	TotalTransactions int64           `json:"total_transactions"`
	TotalInflow       decimal.Decimal `json:"total_inflow"`
	TotalOutflow      decimal.Decimal `json:"total_outflow"`
	NetFlow           decimal.Decimal `json:"net_flow"`
	AverageAmount     decimal.Decimal `json:"average_amount"`
	TodayTransactions int64           `json:"today_transactions"`
	TodayAmount       decimal.Decimal `json:"today_amount"`
	AITransactions    int64           `json:"ai_transactions"`
	AIAmount          decimal.Decimal `json:"ai_amount"`
}
