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
	"gorm.io/gorm"
)

type WalletService struct {
	walletRepo      *repositories.WalletRepository
	transactionRepo *repositories.TransactionRepository
	razorpayClient  *razorpay.Client
	notificationSvc *NotificationService
	db              *gorm.DB
}

func NewWalletService(
	walletRepo *repositories.WalletRepository,
	transactionRepo *repositories.TransactionRepository,
	razorpayClient *razorpay.Client,
	notificationSvc *NotificationService,
	db *gorm.DB,
) *WalletService {
	return &WalletService{
		walletRepo:      walletRepo,
		transactionRepo: transactionRepo,
		razorpayClient:  razorpayClient,
		notificationSvc: notificationSvc,
		db:              db,
	}
}

// Get wallet by user ID
func (s *WalletService) GetWalletByUserID(userID string) (*models.Wallet, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	wallet, err := s.walletRepo.GetByUserID(uid)
	if err != nil {
		return nil, err
	}

	return wallet, nil
}

// Create Razorpay order for loading money
func (s *WalletService) CreateLoadMoneyOrder(userID string, amount decimal.Decimal) (*dto.LoadMoneyResponse, error) {
	// Validate amount (₹10 - ₹50,000)
	if amount.LessThan(decimal.NewFromInt(10)) || amount.GreaterThan(decimal.NewFromInt(50000)) {
		return nil, errors.New("amount must be between ₹10 and ₹50,000")
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

	// Create Razorpay order
	amountInPaise := amount.Mul(decimal.NewFromInt(100)).IntPart()
	order, err := s.razorpayClient.CreateOrder(amountInPaise, "INR", fmt.Sprintf("wallet_load_%s_%d", userID, time.Now().Unix()))
	if err != nil {
		return nil, fmt.Errorf("failed to create Razorpay order: %v", err)
	}

	// Create pending transaction
	transaction := &models.Transaction{
		WalletID:        wallet.ID,
		UserID:          uid,
		Type:            "load_money",
		Amount:          amount,
		Currency:        "INR",
		Description:     "Loading money to wallet",
		RazorpayOrderID: order.ID,
		Status:          "pending",
		ReferenceID:     fmt.Sprintf("LOAD_%s_%d", userID, time.Now().Unix()),
	}

	createdTxn, err := s.transactionRepo.Create(transaction)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %v", err)
	}

	return &dto.LoadMoneyResponse{
		OrderID:       order.ID,
		Amount:        amount,
		Currency:      "INR",
		TransactionID: createdTxn.ID.String(),
	}, nil
}

// Verify Razorpay payment and credit wallet
func (s *WalletService) VerifyAndCreditWallet(userID, paymentID, orderID, signature string) (*dto.PaymentVerificationResponse, error) {
	// Verify Razorpay signature
	if !s.verifyRazorpaySignature(orderID, paymentID, signature) {
		return nil, errors.New("invalid payment signature")
	}

	// Get payment details from Razorpay
	payment, err := s.razorpayClient.GetPayment(paymentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment details: %v", err)
	}

	if payment.Status != "captured" {
		return nil, errors.New("payment not captured")
	}

	amount := decimal.NewFromInt(payment.Amount).Div(decimal.NewFromInt(100)) // Convert paise to rupees

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

	if transaction.Status != "pending" {
		tx.Rollback()
		return nil, errors.New("transaction already processed")
	}

	// Update wallet balance
	wallet, err := s.walletRepo.GetByID(transaction.WalletID)
	if err != nil {
		tx.Rollback()
		return nil, errors.New("wallet not found")
	}

	newBalance := wallet.Balance.Add(amount)
	err = s.walletRepo.UpdateBalance(tx, wallet.ID, newBalance)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update wallet balance: %v", err)
	}

	// Update transaction
	transaction.Status = "success"
	transaction.RazorpayPaymentID = paymentID
	transaction.BalanceAfter = newBalance
	transaction.PaymentMethod = payment.Method

	err = s.transactionRepo.Update(tx, transaction)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update transaction: %v", err)
	}

	// Commit transaction
	if err = tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	// Send notification
	go s.notificationSvc.SendWalletCreditNotification(userID, amount, newBalance)

	return &dto.PaymentVerificationResponse{
		Success:       true,
		NewBalance:    newBalance,
		TransactionID: transaction.ID.String(),
		Message:       fmt.Sprintf("₹%s added to wallet successfully", amount.StringFixed(2)),
	}, nil
}

// Update wallet settings
func (s *WalletService) UpdateWalletSettings(userID string, req *dto.UpdateWalletSettingsRequest) error {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return errors.New("invalid user ID")
	}

	wallet, err := s.walletRepo.GetByUserID(uid)
	if err != nil {
		return errors.New("wallet not found")
	}

	// Update fields if provided
	if req.AIDailyLimit != nil {
		wallet.AIDailyLimit = *req.AIDailyLimit
	}
	if req.AIPerTransactionLimit != nil {
		wallet.AIPerTransactionLimit = *req.AIPerTransactionLimit
	}
	if req.AIAccessEnabled != nil {
		wallet.AIAccessEnabled = *req.AIAccessEnabled
	}

	return s.walletRepo.Update(wallet)
}

// Verify Razorpay signature
func (s *WalletService) verifyRazorpaySignature(orderID, paymentID, signature string) bool {
	body := orderID + "|" + paymentID
	expectedSignature := s.generateSignature(body, s.razorpayClient.KeySecret)
	return hmac.Equal([]byte(expectedSignature), []byte(signature))
}

func (s *WalletService) generateSignature(body, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(body))
	return hex.EncodeToString(h.Sum(nil))
}
