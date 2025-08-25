package services

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/zeusnotfound04/Tranza/config"
	"github.com/zeusnotfound04/Tranza/models"
	"github.com/zeusnotfound04/Tranza/models/dto"
	"github.com/zeusnotfound04/Tranza/pkg/razorpay"
	"github.com/zeusnotfound04/Tranza/repositories"
	"github.com/zeusnotfound04/Tranza/utils"
	"gorm.io/gorm"
)

type ExternalTransferService struct {
	db                   *gorm.DB
	externalTransferRepo *repositories.ExternalTransferRepository
	walletRepo           *repositories.WalletRepository
	transactionRepo      *repositories.TransactionRepository
	razorpayClient       *razorpay.Client
	notificationService  *NotificationService
}

func NewExternalTransferService(
	db *gorm.DB,
	externalTransferRepo *repositories.ExternalTransferRepository,
	walletRepo *repositories.WalletRepository,
	transactionRepo *repositories.TransactionRepository,
	razorpayClient *razorpay.Client,
	notificationService *NotificationService,
) *ExternalTransferService {
	return &ExternalTransferService{
		db:                   db,
		externalTransferRepo: externalTransferRepo,
		walletRepo:           walletRepo,
		transactionRepo:      transactionRepo,
		razorpayClient:       razorpayClient,
		notificationService:  notificationService,
	}
}

// Constants for transfer limits and fees
const (
	MinTransferAmount    = 1.0    // ₹1
	MaxTransferAmount    = 100000 // ₹1,00,000
	DailyTransferLimit   = 50000  // ₹50,000
	MonthlyTransferLimit = 200000 // ₹2,00,000
	UPITransferFee       = 2.0    // ₹2 per transfer
	PhoneTransferFee     = 5.0    // ₹5 per transfer (if phone to UPI conversion)
)

// ValidateTransferRequest validates a transfer request before processing
func (s *ExternalTransferService) ValidateTransferRequest(userID string, req *dto.ValidateTransferRequest) (*dto.ValidateTransferResponse, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	response := &dto.ValidateTransferResponse{
		Valid:    true,
		Warnings: []string{},
		Errors:   []string{},
	}

	// Validate amount
	if req.Amount.LessThan(decimal.NewFromFloat(MinTransferAmount)) {
		response.Valid = false
		response.Errors = append(response.Errors, fmt.Sprintf("Minimum transfer amount is ₹%.0f", MinTransferAmount))
	}

	if req.Amount.GreaterThan(decimal.NewFromFloat(MaxTransferAmount)) {
		response.Valid = false
		response.Errors = append(response.Errors, fmt.Sprintf("Maximum transfer amount is ₹%.0f", float64(MaxTransferAmount)))
	}

	// Validate recipient
	if err := s.validateRecipient(req.RecipientType, req.RecipientValue); err != nil {
		response.Valid = false
		response.Errors = append(response.Errors, err.Error())
	}

	// Check wallet balance
	wallet, err := s.walletRepo.GetByUserID(uid)
	if err != nil {
		response.Valid = false
		response.Errors = append(response.Errors, "Wallet not found")
		return response, nil
	}

	// Calculate transfer fee
	transferFee := s.calculateTransferFee(req.RecipientType, req.Amount)
	totalAmount := req.Amount.Add(transferFee)

	response.TransferFee = transferFee
	response.TotalAmount = totalAmount

	// Check sufficient balance
	if wallet.Balance.LessThan(totalAmount) {
		response.Valid = false
		response.Errors = append(response.Errors, "Insufficient wallet balance")
	}

	// Check daily limit
	today := time.Now()
	dailyUsed, err := s.externalTransferRepo.GetDailyTransferLimitUsed(uid, today)
	if err == nil {
		if dailyUsed.Add(totalAmount).GreaterThan(decimal.NewFromFloat(DailyTransferLimit)) {
			response.Valid = false
			response.Errors = append(response.Errors, "Daily transfer limit exceeded")
		} else if dailyUsed.Add(totalAmount).GreaterThan(decimal.NewFromFloat(DailyTransferLimit * 0.8)) {
			response.Warnings = append(response.Warnings, "Approaching daily transfer limit")
		}
	}

	// Check monthly limit
	monthlyUsed, err := s.externalTransferRepo.GetMonthlyTransferLimitUsed(uid, today)
	if err == nil {
		if monthlyUsed.Add(totalAmount).GreaterThan(decimal.NewFromFloat(MonthlyTransferLimit)) {
			response.Valid = false
			response.Errors = append(response.Errors, "Monthly transfer limit exceeded")
		}
	}

	// Set estimated time
	response.EstimatedTime = s.getEstimatedTransferTime(req.RecipientType)

	return response, nil
}

// CreateExternalTransfer creates a new external transfer
func (s *ExternalTransferService) CreateExternalTransfer(userID string, req *dto.CreateExternalTransferRequest) (*dto.ExternalTransferResponse, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	// Validate the request first
	validateReq := &dto.ValidateTransferRequest{
		Amount:         req.Amount,
		RecipientType:  req.RecipientType,
		RecipientValue: req.RecipientValue,
	}

	validation, err := s.ValidateTransferRequest(userID, validateReq)
	if err != nil {
		return nil, err
	}

	if !validation.Valid {
		return nil, errors.New(strings.Join(validation.Errors, "; "))
	}

	// Get wallet
	wallet, err := s.walletRepo.GetByUserID(uid)
	if err != nil {
		return nil, errors.New("wallet not found")
	}

	// Calculate fees
	transferFee := s.calculateTransferFee(req.RecipientType, req.Amount)
	totalAmount := req.Amount.Add(transferFee)

	// Start database transaction
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create external transfer record
	transfer := &models.ExternalTransfer{
		UserID:         uid,
		WalletID:       wallet.ID,
		Amount:         req.Amount,
		Currency:       "INR",
		Description:    req.Description,
		RecipientType:  req.RecipientType,
		RecipientValue: req.RecipientValue,
		RecipientName:  req.RecipientName,
		Status:         models.ExternalTransferStatusPending,
		TransferMethod: models.TransferMethodRazorpayPayout,
		TransferFee:    transferFee,
		TotalAmount:    totalAmount,
		InitiatedBy:    models.InitiatedByUser,
		BalanceBefore:  wallet.Balance,
		BalanceAfter:   wallet.Balance.Sub(totalAmount),
		MaxRetries:     3,
	}

	createdTransfer, err := s.externalTransferRepo.Create(transfer)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create transfer record: %w", err)
	}

	// Create corresponding transaction record
	transaction := &models.Transaction{
		WalletID:     wallet.ID,
		UserID:       uid,
		Type:         utils.TransactionTypeExternalTransfer,
		Amount:       totalAmount,
		Currency:     "INR",
		Description:  fmt.Sprintf("External transfer to %s", req.RecipientValue),
		Status:       models.StatusPending,
		ReferenceID:  createdTransfer.ReferenceID,
		BalanceAfter: wallet.Balance.Sub(totalAmount),
	}

	createdTransaction, err := s.transactionRepo.Create(transaction)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create transaction record: %w", err)
	}

	// Update transfer with transaction ID
	createdTransfer.TransactionID = &createdTransaction.ID
	if err := s.externalTransferRepo.Update(tx, createdTransfer); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update transfer with transaction ID: %w", err)
	}

	// Deduct amount from wallet (reserve funds)
	newBalance := wallet.Balance.Sub(totalAmount)
	if err := s.walletRepo.UpdateBalance(tx, wallet.ID, newBalance); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update wallet balance: %w", err)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transfer transaction: %w", err)
	}

	// Check if we're in test mode - process synchronously for immediate completion
	cfg := config.LoadConfig()
	if cfg.Environment == "test" {
		// Process synchronously in test mode
		s.processTransferAsync(createdTransfer.ID.String())

		// Fetch updated transfer to get the new status
		updatedTransfer, err := s.externalTransferRepo.GetByID(createdTransfer.ID)
		if err == nil {
			return &dto.ExternalTransferResponse{
				ID:             updatedTransfer.ID.String(),
				ReferenceID:    updatedTransfer.ReferenceID,
				Amount:         req.Amount,
				Currency:       "INR",
				TransferFee:    transferFee,
				TotalAmount:    totalAmount,
				RecipientType:  req.RecipientType,
				RecipientValue: req.RecipientValue,
				RecipientName:  req.RecipientName,
				Status:         updatedTransfer.Status,
				CreatedAt:      updatedTransfer.CreatedAt,
				EstimatedTime:  s.getEstimatedTransferTimeForStatus(req.RecipientType, updatedTransfer.Status),
			}, nil
		}
	} else {
		// Process the transfer asynchronously in production
		go s.processTransferAsync(createdTransfer.ID.String())
	}

	// Send notification
	// go s.notificationService.SendExternalTransferInitiatedNotification(userID, req.Amount, req.RecipientValue)

	return &dto.ExternalTransferResponse{
		ID:             createdTransfer.ID.String(),
		ReferenceID:    createdTransfer.ReferenceID,
		Amount:         req.Amount,
		Currency:       "INR",
		TransferFee:    transferFee,
		TotalAmount:    totalAmount,
		RecipientType:  req.RecipientType,
		RecipientValue: req.RecipientValue,
		RecipientName:  req.RecipientName,
		Status:         models.ExternalTransferStatusPending,
		CreatedAt:      createdTransfer.CreatedAt,
		EstimatedTime:  s.getEstimatedTransferTime(req.RecipientType),
	}, nil
}

// GetExternalTransfer retrieves an external transfer by ID
func (s *ExternalTransferService) GetExternalTransfer(transferID string) (*dto.ExternalTransferResponse, error) {
	id, err := uuid.Parse(transferID)
	if err != nil {
		return nil, errors.New("invalid transfer ID")
	}

	transfer, err := s.externalTransferRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("transfer not found: %w", err)
	}

	return &dto.ExternalTransferResponse{
		ID:             transfer.ID.String(),
		ReferenceID:    transfer.ReferenceID,
		Amount:         transfer.Amount,
		Currency:       transfer.Currency,
		TransferFee:    transfer.TransferFee,
		TotalAmount:    transfer.TotalAmount,
		RecipientType:  transfer.RecipientType,
		RecipientValue: transfer.RecipientValue,
		RecipientName:  transfer.RecipientName,
		Status:         transfer.Status,
		CreatedAt:      transfer.CreatedAt,
		EstimatedTime:  s.getEstimatedTransferTimeForStatus(transfer.RecipientType, transfer.Status),
	}, nil
}

// GetExternalTransfersByUser retrieves external transfers for a user with pagination
func (s *ExternalTransferService) GetExternalTransfersByUser(userID string, page, limit int) (*dto.PaginatedExternalTransferResponse, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	offset := (page - 1) * limit
	transfers, total, err := s.externalTransferRepo.GetByUserIDWithPagination(uid, limit, offset, "", nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get transfers: %w", err)
	}

	var transferResponses []dto.ExternalTransferResponse
	for _, transfer := range transfers {
		transferResponses = append(transferResponses, dto.ExternalTransferResponse{
			ID:             transfer.ID.String(),
			ReferenceID:    transfer.ReferenceID,
			Amount:         transfer.Amount,
			Currency:       transfer.Currency,
			TransferFee:    transfer.TransferFee,
			TotalAmount:    transfer.TotalAmount,
			RecipientType:  transfer.RecipientType,
			RecipientValue: transfer.RecipientValue,
			RecipientName:  transfer.RecipientName,
			Status:         transfer.Status,
			CreatedAt:      transfer.CreatedAt,
		})
	}

	return &dto.PaginatedExternalTransferResponse{
		Transfers:   transferResponses,
		TotalCount:  total,
		CurrentPage: page,
		TotalPages:  (total + int64(limit) - 1) / int64(limit),
		HasMore:     int64(page*limit) < total,
	}, nil
}

// processTransferAsync processes the transfer with Razorpay asynchronously
func (s *ExternalTransferService) processTransferAsync(transferID string) {
	id, err := uuid.Parse(transferID)
	if err != nil {
		utils.LogError(err, map[string]interface{}{"transfer_id": transferID, "action": "parse_transfer_id"})
		return
	}

	transfer, err := s.externalTransferRepo.GetByID(id)
	if err != nil {
		utils.LogError(err, map[string]interface{}{"transfer_id": transferID, "action": "get_transfer"})
		return
	}

	// Process with Razorpay
	if err := s.processWithRazorpay(transfer); err != nil {
		utils.LogError(err, map[string]interface{}{"transfer_id": transferID, "action": "process_with_razorpay"})

		// Update status to failed
		s.externalTransferRepo.UpdateStatus(id, models.ExternalTransferStatusFailed, err.Error())

		// Refund wallet balance
		s.refundWalletBalance(transfer)

		// Send failure notification
		// go s.notificationService.SendExternalTransferFailedNotification(
		// 	transfer.UserID.String(), transfer.Amount, transfer.RecipientValue, err.Error())
	}
}

// processWithRazorpay processes the transfer using Razorpay Payouts
func (s *ExternalTransferService) processWithRazorpay(transfer *models.ExternalTransfer) error {
	// Convert amount to paise
	amountInPaise := transfer.Amount.Mul(decimal.NewFromInt(100)).IntPart()

	// Ensure recipient name is not empty (required by Razorpay)
	recipientName := transfer.RecipientName
	if recipientName == "" {
		// Use a default name for testing/when name is not provided
		recipientName = "Test User"
	}

	// Check if we're in test mode and should simulate success
	cfg := config.LoadConfig()
	if cfg.Environment == "test" {
		// Simulate successful payout for testing
		log.Printf("Simulating successful payout in test mode, transfer_id: %s, amount: %s, recipient: %s",
			transfer.ID.String(),
			transfer.Amount.String(),
			transfer.RecipientValue)

		// Generate a fake payout ID
		fakePayoutID := fmt.Sprintf("pout_test_%s", transfer.ReferenceID)

		// Update transfer with fake Razorpay payout ID
		if err := s.externalTransferRepo.UpdateRazorpayPayoutID(transfer.ID, fakePayoutID); err != nil {
			return fmt.Errorf("failed to update transfer with payout ID: %w", err)
		}

		// Immediately mark as successful
		if err := s.externalTransferRepo.UpdateStatus(transfer.ID, models.ExternalTransferStatusSuccess, ""); err != nil {
			return fmt.Errorf("failed to update transfer status: %w", err)
		}

		log.Printf("Test transfer completed successfully, transfer_id: %s, payout_id: %s",
			transfer.ID.String(),
			fakePayoutID)

		return nil
	}

	var payout *razorpay.Payout
	var err error

	switch transfer.RecipientType {
	case models.RecipientTypeUPI:
		payout, err = s.razorpayClient.CreateUPIPayout(
			transfer.RecipientValue,
			amountInPaise,
			"INR",
			razorpay.PurposePayout,
			transfer.Description,
			recipientName,
			"",
			transfer.ReferenceID,
		)
	case models.RecipientTypePhone:
		// For phone numbers, we might need to convert to UPI ID
		// For now, treat as UPI with @paytm suffix
		upiID := s.convertPhoneToUPI(transfer.RecipientValue)
		payout, err = s.razorpayClient.CreateUPIPayout(
			upiID,
			amountInPaise,
			"INR",
			razorpay.PurposePayout,
			transfer.Description,
			recipientName,
			transfer.RecipientValue,
			transfer.ReferenceID,
		)
	default:
		return errors.New("unsupported recipient type")
	}

	if err != nil {
		return fmt.Errorf("failed to create Razorpay payout: %w", err)
	}

	// Update transfer with Razorpay payout ID
	if err := s.externalTransferRepo.UpdateRazorpayPayoutID(transfer.ID, payout.ID); err != nil {
		return fmt.Errorf("failed to update transfer with payout ID: %w", err)
	}

	// Start monitoring the payout status
	go s.monitorPayoutStatus(transfer.ID.String(), payout.ID)

	return nil
}

// monitorPayoutStatus monitors the Razorpay payout status
func (s *ExternalTransferService) monitorPayoutStatus(transferID, payoutID string) {
	id, err := uuid.Parse(transferID)
	if err != nil {
		utils.LogError(err, map[string]interface{}{"transfer_id": transferID, "action": "parse_transfer_id"})
		return
	}

	// Check status every 30 seconds for up to 10 minutes
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	timeout := time.After(10 * time.Minute)

	for {
		select {
		case <-ticker.C:
			payout, err := s.razorpayClient.GetPayout(payoutID)
			if err != nil {
				utils.LogError(err, map[string]interface{}{"transfer_id": transferID, "payout_id": payoutID, "action": "get_payout_status"})
				continue
			}

			switch payout.Status {
			case razorpay.PayoutStatusProcessed:
				s.handlePayoutSuccess(id, payout)
				return
			case razorpay.PayoutStatusFailed, razorpay.PayoutStatusCancelled:
				s.handlePayoutFailure(id, payout)
				return
			case razorpay.PayoutStatusReversed:
				s.handlePayoutReversal(id, payout)
				return
			}

		case <-timeout:
			// Timeout - mark as failed
			s.externalTransferRepo.UpdateStatus(id, models.ExternalTransferStatusFailed, "Transfer timeout")
			transfer, _ := s.externalTransferRepo.GetByID(id)
			if transfer != nil {
				s.refundWalletBalance(transfer)
			}
			return
		}
	}
}

// handlePayoutSuccess handles successful payout
func (s *ExternalTransferService) handlePayoutSuccess(transferID uuid.UUID, payout *razorpay.Payout) {
	if err := s.externalTransferRepo.UpdateStatus(transferID, models.ExternalTransferStatusSuccess, ""); err != nil {
		utils.LogError(err, map[string]interface{}{"transfer_id": transferID.String(), "action": "update_success_status"})
		return
	}

	// Update corresponding transaction
	transfer, err := s.externalTransferRepo.GetByID(transferID)
	if err != nil {
		utils.LogError(err, map[string]interface{}{"transfer_id": transferID.String(), "action": "get_transfer_for_success"})
		return
	}

	if transfer.TransactionID != nil {
		s.transactionRepo.UpdateStatus(*transfer.TransactionID, utils.TransactionStatusSuccess, "")
	}

	// Send success notification
	// go s.notificationService.SendExternalTransferSuccessNotification(
	// 	transfer.UserID.String(), transfer.Amount, transfer.RecipientValue, payout.UTR)

	utils.LogInfo("External transfer completed successfully", map[string]interface{}{
		"transfer_id": transferID.String(),
		"payout_id":   payout.ID,
		"utr":         payout.UTR,
		"amount":      transfer.Amount.String(),
	})
}

// handlePayoutFailure handles failed payout
func (s *ExternalTransferService) handlePayoutFailure(transferID uuid.UUID, payout *razorpay.Payout) {
	if err := s.externalTransferRepo.UpdateStatus(transferID, models.ExternalTransferStatusFailed, payout.FailureReason); err != nil {
		utils.LogError(err, map[string]interface{}{"transfer_id": transferID.String(), "action": "update_failure_status"})
		return
	}

	// Refund wallet balance
	transfer, err := s.externalTransferRepo.GetByID(transferID)
	if err != nil {
		utils.LogError(err, map[string]interface{}{"transfer_id": transferID.String(), "action": "get_transfer_for_refund"})
		return
	}

	s.refundWalletBalance(transfer)

	// Update corresponding transaction
	if transfer.TransactionID != nil {
		s.transactionRepo.UpdateStatus(*transfer.TransactionID, utils.TransactionStatusFailed, payout.FailureReason)
	}

	// Send failure notification
	// go s.notificationService.SendExternalTransferFailedNotification(
	// 	transfer.UserID.String(), transfer.Amount, transfer.RecipientValue, payout.FailureReason)

	utils.LogWarning("External transfer failed", map[string]interface{}{
		"transfer_id":    transferID.String(),
		"payout_id":      payout.ID,
		"failure_reason": payout.FailureReason,
		"amount":         transfer.Amount.String(),
	})
}

// handlePayoutReversal handles reversed payout
func (s *ExternalTransferService) handlePayoutReversal(transferID uuid.UUID, payout *razorpay.Payout) {
	if err := s.externalTransferRepo.UpdateStatus(transferID, models.ExternalTransferStatusRefunded, "Payout reversed by bank"); err != nil {
		utils.LogError(err, map[string]interface{}{"transfer_id": transferID.String(), "action": "update_reversal_status"})
		return
	}

	// Refund wallet balance
	transfer, err := s.externalTransferRepo.GetByID(transferID)
	if err != nil {
		utils.LogError(err, map[string]interface{}{"transfer_id": transferID.String(), "action": "get_transfer_for_reversal"})
		return
	}

	s.refundWalletBalance(transfer)

	// Send reversal notification
	// go s.notificationService.SendExternalTransferReversalNotification(
	// 	transfer.UserID.String(), transfer.Amount, transfer.RecipientValue)
}

// refundWalletBalance refunds the wallet balance for failed/reversed transfers
func (s *ExternalTransferService) refundWalletBalance(transfer *models.ExternalTransfer) {
	// Get current wallet
	wallet, err := s.walletRepo.GetByID(transfer.WalletID)
	if err != nil {
		utils.LogError(err, map[string]interface{}{"wallet_id": transfer.WalletID.String(), "action": "get_wallet_for_refund"})
		return
	}

	// Refund the amount
	newBalance := wallet.Balance.Add(transfer.TotalAmount)
	if err := s.walletRepo.UpdateBalance(nil, wallet.ID, newBalance); err != nil {
		utils.LogError(err, map[string]interface{}{"wallet_id": transfer.WalletID.String(), "action": "refund_wallet_balance"})
		return
	}

	// Create refund transaction
	refundTransaction := &models.Transaction{
		WalletID:     wallet.ID,
		UserID:       transfer.UserID,
		Type:         utils.TransactionTypeRefund,
		Amount:       transfer.TotalAmount,
		Currency:     "INR",
		Description:  fmt.Sprintf("Refund for failed transfer %s", transfer.ReferenceID),
		Status:       models.StatusSuccess,
		ReferenceID:  "REFUND_" + transfer.ReferenceID,
		BalanceAfter: newBalance,
	}

	s.transactionRepo.Create(refundTransaction)

	utils.LogInfo("Wallet balance refunded", map[string]interface{}{
		"transfer_id":   transfer.ID.String(),
		"wallet_id":     wallet.ID.String(),
		"refund_amount": transfer.TotalAmount.String(),
		"new_balance":   newBalance.String(),
	})
}

// Helper functions

func (s *ExternalTransferService) validateRecipient(recipientType, recipientValue string) error {
	switch recipientType {
	case models.RecipientTypeUPI:
		return s.validateUPIID(recipientValue)
	case models.RecipientTypePhone:
		return s.validatePhoneNumber(recipientValue)
	default:
		return errors.New("invalid recipient type")
	}
}

func (s *ExternalTransferService) validateUPIID(upiID string) error {
	// UPI ID format: username@provider
	upiRegex := regexp.MustCompile(`^[a-zA-Z0-9._-]+@[a-zA-Z0-9.-]+$`)
	if !upiRegex.MatchString(upiID) {
		return errors.New("invalid UPI ID format")
	}
	return nil
}

func (s *ExternalTransferService) validatePhoneNumber(phone string) error {
	// Indian phone number format: 10 digits starting with 6-9
	phoneRegex := regexp.MustCompile(`^[6-9]\d{9}$`)
	if !phoneRegex.MatchString(phone) {
		return errors.New("invalid phone number format")
	}
	return nil
}

func (s *ExternalTransferService) calculateTransferFee(recipientType string, amount decimal.Decimal) decimal.Decimal {
	switch recipientType {
	case models.RecipientTypeUPI:
		return decimal.NewFromFloat(UPITransferFee)
	case models.RecipientTypePhone:
		return decimal.NewFromFloat(PhoneTransferFee)
	default:
		return decimal.NewFromFloat(UPITransferFee)
	}
}

func (s *ExternalTransferService) getEstimatedTransferTime(recipientType string) string {
	// In test mode, don't show estimated time for pending transfers as they complete instantly
	cfg := config.LoadConfig()
	if cfg.Environment == "test" || cfg.Environment == "development" {
		return ""
	}

	switch recipientType {
	case models.RecipientTypeUPI:
		return "Instant (within 2 minutes)"
	case models.RecipientTypePhone:
		return "2-5 minutes"
	default:
		return "2-5 minutes"
	}
}

// getEstimatedTransferTimeForStatus returns estimated time based on status
func (s *ExternalTransferService) getEstimatedTransferTimeForStatus(recipientType, status string) string {
	// If transfer is completed, failed, or cancelled, no estimated time needed
	if status == models.ExternalTransferStatusSuccess ||
		status == models.ExternalTransferStatusFailed ||
		status == models.ExternalTransferStatusCancelled {
		return ""
	}

	return s.getEstimatedTransferTime(recipientType)
}

func (s *ExternalTransferService) convertPhoneToUPI(phone string) string {
	// Simple conversion - in production, you might want to use a service
	// to determine the actual UPI ID associated with the phone number
	return phone + "@paytm"
}
