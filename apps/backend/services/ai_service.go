package services

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/zeusnotfound04/Tranza/models"
	"gorm.io/gorm"
)

type AIService struct {
	db           *gorm.DB
	geminiAPIKey string
}

func NewAIService(db *gorm.DB, geminiAPIKey string) *AIService {
	return &AIService{
		db:           db,
		geminiAPIKey: geminiAPIKey,
	}
}

// ProcessPaymentRequest analyzes natural language and creates payment request
func (s *AIService) ProcessPaymentRequest(userID uuid.UUID, req models.AIPaymentRequestDTO) (*models.AIPaymentResponse, error) {
	// Check if user has AI access enabled
	limits, err := s.GetSpendingLimits(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get spending limits: %v", err)
	}

	if !limits.AIAccessEnabled {
		return nil, fmt.Errorf("AI payment processing is not enabled for this user")
	}

	// Parse the natural language prompt using simple pattern matching
	// Note: Gemini AI integration can be added later
	analysisResult := s.analyzePaymentPrompt(req.Prompt)

	// Override with explicit values if provided
	if req.Amount > 0 {
		analysisResult.Amount = req.Amount
	}
	if req.Merchant != "" {
		analysisResult.Merchant = req.Merchant
	}
	if req.Description != "" {
		analysisResult.Description = req.Description
	}

	// Validate against spending limits
	if err := s.validateSpendingLimits(userID, analysisResult.Amount); err != nil {
		return nil, err
	}

	// Assess risk level
	riskLevel := s.assessRiskLevel(analysisResult.Amount, analysisResult.Merchant)

	// Create AI payment request record
	paymentRequest := &models.AIPaymentRequest{
		UserID:       userID,
		Amount:       analysisResult.Amount,
		Description:  analysisResult.Description,
		MerchantName: analysisResult.Merchant,
		AIPrompt:     req.Prompt,
		Status:       "pending",
		AIResponse:   analysisResult.AIReasoning,
		RiskLevel:    riskLevel,
		Confidence:   analysisResult.Confidence,
	}

	if err := s.db.Create(paymentRequest).Error; err != nil {
		return nil, fmt.Errorf("failed to create payment request: %v", err)
	}

	// Get user's wallet balance
	var wallet models.Wallet
	if err := s.db.Where("user_id = ?", userID).First(&wallet).Error; err != nil {
		return nil, fmt.Errorf("failed to get wallet balance: %v", err)
	}

	// Calculate remaining daily limit
	remainingLimit, err := s.calculateRemainingDailyLimit(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate remaining limit: %v", err)
	}

	// Determine if confirmation is required
	requiresConfirmation := analysisResult.Amount >= limits.ConfirmationThreshold ||
		limits.RequireConfirmation ||
		riskLevel == "high"

	// Convert decimal to float64 for response
	walletBalance, _ := wallet.Balance.Float64()

	response := &models.AIPaymentResponse{
		ID:                   paymentRequest.ID.String(),
		Amount:               analysisResult.Amount,
		Merchant:             analysisResult.Merchant,
		Description:          analysisResult.Description,
		Confidence:           analysisResult.Confidence,
		RiskLevel:            riskLevel,
		RequiresConfirmation: requiresConfirmation,
		AIReasoning:          analysisResult.AIReasoning,
		WalletBalance:        walletBalance,
		RemainingLimit:       remainingLimit,
		Suggestions:          s.generateSuggestions(analysisResult, walletBalance, remainingLimit),
	}

	return response, nil
}

// ConfirmPayment processes confirmed payment requests
func (s *AIService) ConfirmPayment(userID, paymentID uuid.UUID, confirmed bool) (interface{}, error) {
	// Get the payment request
	var paymentRequest models.AIPaymentRequest
	if err := s.db.Where("id = ? AND user_id = ?", paymentID, userID).First(&paymentRequest).Error; err != nil {
		return nil, fmt.Errorf("payment request not found: %v", err)
	}

	if paymentRequest.Status != "pending" {
		return nil, fmt.Errorf("payment request is not in pending status")
	}

	if !confirmed {
		// Mark as cancelled
		paymentRequest.Status = "cancelled"
		if err := s.db.Save(&paymentRequest).Error; err != nil {
			return nil, fmt.Errorf("failed to cancel payment: %v", err)
		}
		return gin.H{"status": "cancelled", "message": "Payment request cancelled"}, nil
	}

	// Get wallet for transaction
	var wallet models.Wallet
	if err := s.db.Where("user_id = ?", userID).First(&wallet).Error; err != nil {
		return nil, fmt.Errorf("failed to get wallet: %v", err)
	}

	// Process the payment - create actual transaction
	transaction := &models.Transaction{
		WalletID:     wallet.ID,
		UserID:       userID,
		Type:         "debit",
		Amount:       models.DecimalFromFloat64(paymentRequest.Amount),
		Description:  fmt.Sprintf("AI Payment: %s", paymentRequest.Description),
		Status:       models.StatusSuccess,
		BalanceAfter: wallet.Balance.Sub(models.DecimalFromFloat64(paymentRequest.Amount)),
	}

	// Start database transaction
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create transaction record
	if err := tx.Create(transaction).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create transaction: %v", err)
	}

	// Update wallet balance
	wallet.Balance = wallet.Balance.Sub(models.DecimalFromFloat64(paymentRequest.Amount))
	if err := tx.Save(&wallet).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update wallet: %v", err)
	}

	// Update payment request
	paymentRequest.Status = "processed"
	now := time.Now()
	paymentRequest.ProcessedAt = &now
	paymentRequest.TransactionID = transaction.ID.String()

	if err := tx.Save(&paymentRequest).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update payment request: %v", err)
	}

	// Update spending tracker
	if err := s.updateSpendingTracker(tx, userID, paymentRequest.Amount); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update spending tracker: %v", err)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return gin.H{
		"status":         "processed",
		"transaction_id": transaction.ID.String(),
		"amount":         paymentRequest.Amount,
		"merchant":       paymentRequest.MerchantName,
		"message":        "Payment processed successfully",
	}, nil
}

// GetPaymentHistory retrieves AI payment history for a user
func (s *AIService) GetPaymentHistory(userID uuid.UUID, page, limit int, status string) ([]models.AIPaymentRequest, int64, error) {
	var payments []models.AIPaymentRequest
	var total int64

	query := s.db.Where("user_id = ?", userID)
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Count total records
	if err := query.Model(&models.AIPaymentRequest{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count payments: %v", err)
	}

	// Get paginated results
	offset := (page - 1) * limit
	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&payments).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get payment history: %v", err)
	}

	return payments, total, nil
}

// GetSpendingLimits retrieves user's spending limits
func (s *AIService) GetSpendingLimits(userID uuid.UUID) (*models.AISpendingLimit, error) {
	var limits models.AISpendingLimit
	err := s.db.Where("user_id = ?", userID).First(&limits).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Create default limits for new user
			limits = models.AISpendingLimit{
				UserID:                userID,
				DailyLimit:            10000,
				TransactionLimit:      2000,
				MonthlyLimit:          100000,
				AIAccessEnabled:       false,
				RequireConfirmation:   true,
				ConfirmationThreshold: 1000,
			}
			if err := s.db.Create(&limits).Error; err != nil {
				return nil, fmt.Errorf("failed to create default limits: %v", err)
			}
		} else {
			return nil, fmt.Errorf("failed to get spending limits: %v", err)
		}
	}
	return &limits, nil
}

// UpdateSpendingLimits updates user's spending limits
func (s *AIService) UpdateSpendingLimits(userID uuid.UUID, req models.AISpendingLimitsDTO) (*models.AISpendingLimit, error) {
	limits, err := s.GetSpendingLimits(userID)
	if err != nil {
		return nil, err
	}

	// Update limits
	limits.DailyLimit = req.DailyLimit
	limits.TransactionLimit = req.TransactionLimit
	limits.MonthlyLimit = req.MonthlyLimit
	limits.AIAccessEnabled = req.AIAccessEnabled
	limits.RequireConfirmation = req.RequireConfirmation
	limits.ConfirmationThreshold = req.ConfirmationThreshold

	if err := s.db.Save(limits).Error; err != nil {
		return nil, fmt.Errorf("failed to update spending limits: %v", err)
	}

	return limits, nil
}

// GetSpendingAnalytics provides spending insights
func (s *AIService) GetSpendingAnalytics(userID uuid.UUID, period string) (interface{}, error) {
	var startDate time.Time
	now := time.Now()

	switch period {
	case "day":
		startDate = now.Truncate(24 * time.Hour)
	case "week":
		startDate = now.AddDate(0, 0, -7)
	case "month":
		startDate = now.AddDate(0, -1, 0)
	case "year":
		startDate = now.AddDate(-1, 0, 0)
	default:
		startDate = now.AddDate(0, -1, 0) // Default to month
	}

	// Get spending data
	var totalSpent float64
	var transactionCount int64

	if err := s.db.Model(&models.AIPaymentRequest{}).
		Where("user_id = ? AND status = 'processed' AND created_at >= ?", userID, startDate).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalSpent).Error; err != nil {
		return nil, fmt.Errorf("failed to calculate total spent: %v", err)
	}

	if err := s.db.Model(&models.AIPaymentRequest{}).
		Where("user_id = ? AND status = 'processed' AND created_at >= ?", userID, startDate).
		Count(&transactionCount).Error; err != nil {
		return nil, fmt.Errorf("failed to count transactions: %v", err)
	}

	// Get top merchants
	var topMerchants []struct {
		MerchantName string  `json:"merchant_name"`
		TotalSpent   float64 `json:"total_spent"`
		Count        int     `json:"count"`
	}

	if err := s.db.Model(&models.AIPaymentRequest{}).
		Where("user_id = ? AND status = 'processed' AND created_at >= ?", userID, startDate).
		Select("merchant_name, SUM(amount) as total_spent, COUNT(*) as count").
		Group("merchant_name").
		Order("total_spent DESC").
		Limit(10).
		Scan(&topMerchants).Error; err != nil {
		return nil, fmt.Errorf("failed to get top merchants: %v", err)
	}

	return gin.H{
		"period":            period,
		"total_spent":       totalSpent,
		"transaction_count": transactionCount,
		"average_amount": func() float64 {
			if transactionCount > 0 {
				return totalSpent / float64(transactionCount)
			}
			return 0
		}(),
		"top_merchants": topMerchants,
		"insights": []string{
			fmt.Sprintf("You've spent %.2f in the last %s", totalSpent, period),
			fmt.Sprintf("Average transaction amount: %.2f", func() float64 {
				if transactionCount > 0 {
					return totalSpent / float64(transactionCount)
				}
				return 0
			}()),
		},
	}, nil
}

// GetPaymentRequest retrieves specific payment request details
func (s *AIService) GetPaymentRequest(userID, paymentID uuid.UUID) (*models.AIPaymentRequest, error) {
	var payment models.AIPaymentRequest
	if err := s.db.Where("id = ? AND user_id = ?", paymentID, userID).First(&payment).Error; err != nil {
		return nil, fmt.Errorf("payment request not found: %v", err)
	}
	return &payment, nil
}

// CancelPaymentRequest cancels a pending payment request
func (s *AIService) CancelPaymentRequest(userID, paymentID uuid.UUID) error {
	var payment models.AIPaymentRequest
	if err := s.db.Where("id = ? AND user_id = ?", paymentID, userID).First(&payment).Error; err != nil {
		return fmt.Errorf("payment request not found: %v", err)
	}

	if payment.Status != "pending" {
		return fmt.Errorf("can only cancel pending payment requests")
	}

	payment.Status = "cancelled"
	if err := s.db.Save(&payment).Error; err != nil {
		return fmt.Errorf("failed to cancel payment: %v", err)
	}

	return nil
}

// Private helper methods

type PaymentAnalysis struct {
	Amount      float64
	Merchant    string
	Description string
	Confidence  float64
	AIReasoning string
}

func (s *AIService) analyzePaymentPrompt(prompt string) *PaymentAnalysis {
	// Enhanced pattern matching - fallback when Gemini is not available
	prompt = strings.ToLower(prompt)

	analysis := &PaymentAnalysis{
		Confidence:  0.6, // Start with lower confidence
		AIReasoning: "Pattern matching analysis",
	}

	// Enhanced amount extraction with multiple patterns
	amountPatterns := []string{
		`(?:rs\.?|rupees?|₹)\s*(\d+(?:,\d{3})*(?:\.\d{2})?)`, // Rs. 1,000.00
		`(\d+(?:,\d{3})*(?:\.\d{2})?)\s*(?:rs\.?|rupees?|₹)`, // 1000 rs
		`(\d+(?:,\d{3})*(?:\.\d{2})?)\s*(?:only|/-)?`,        // 1000 only
	}

	for _, pattern := range amountPatterns {
		amountRegex := regexp.MustCompile(pattern)
		if matches := amountRegex.FindStringSubmatch(prompt); len(matches) > 1 {
			amountStr := strings.ReplaceAll(matches[1], ",", "") // Remove commas
			if amount, err := strconv.ParseFloat(amountStr, 64); err == nil && amount > 0 {
				analysis.Amount = amount
				analysis.Confidence += 0.3
				break
			}
		}
	}

	// Enhanced merchant extraction with better patterns
	merchantPatterns := []string{
		`(?:pay|send|transfer)\s+(?:to\s+)?([a-z][a-z0-9\s]{1,30})(?:\s+for|\s+rs|\s*₹|\s*$)`,
		`(?:for|at|from)\s+([a-z][a-z0-9\s]{1,30})(?:\s+for|\s+rs|\s*₹|\s*$)`,
		`(?:to|towards)\s+([a-z][a-z0-9\s]{1,30})(?:\s+for|\s+rs|\s*₹|\s*$)`,
	}

	for _, pattern := range merchantPatterns {
		merchantRegex := regexp.MustCompile(pattern)
		if matches := merchantRegex.FindStringSubmatch(prompt); len(matches) > 1 {
			merchant := strings.TrimSpace(matches[1])
			// Clean up common words
			merchant = regexp.MustCompile(`\b(the|and|for|of|in|on|at|to|from)\b`).ReplaceAllString(merchant, "")
			merchant = strings.TrimSpace(merchant)
			if len(merchant) > 2 {
				analysis.Merchant = strings.Title(merchant)
				analysis.Confidence += 0.2
				break
			}
		}
	}

	// Extract purpose/description
	purposePatterns := []string{
		`for\s+([a-z\s]{3,30})(?:\s|$)`,
		`(?:buying|purchasing)\s+([a-z\s]{3,30})(?:\s|$)`,
		`(?:payment for|bill for)\s+([a-z\s]{3,30})(?:\s|$)`,
	}

	var purpose string
	for _, pattern := range purposePatterns {
		purposeRegex := regexp.MustCompile(pattern)
		if matches := purposeRegex.FindStringSubmatch(prompt); len(matches) > 1 {
			purpose = strings.TrimSpace(matches[1])
			if len(purpose) > 2 {
				analysis.Confidence += 0.1
				break
			}
		}
	}

	// Generate intelligent description
	if analysis.Merchant != "" && analysis.Amount > 0 {
		if purpose != "" {
			analysis.Description = fmt.Sprintf("Payment of Rs. %.2f to %s for %s", analysis.Amount, analysis.Merchant, purpose)
		} else {
			analysis.Description = fmt.Sprintf("Payment of Rs. %.2f to %s", analysis.Amount, analysis.Merchant)
		}
	} else if analysis.Amount > 0 {
		analysis.Description = fmt.Sprintf("Payment of Rs. %.2f", analysis.Amount)
	} else if analysis.Merchant != "" {
		analysis.Description = fmt.Sprintf("Payment to %s", analysis.Merchant)
	} else {
		analysis.Description = "Payment request"
		analysis.Confidence = 0.3 // Low confidence for unclear requests
	}

	// Update reasoning based on what was extracted
	reasoningParts := []string{}
	if analysis.Amount > 0 {
		reasoningParts = append(reasoningParts, fmt.Sprintf("extracted amount: Rs. %.2f", analysis.Amount))
	}
	if analysis.Merchant != "" {
		reasoningParts = append(reasoningParts, fmt.Sprintf("identified merchant: %s", analysis.Merchant))
	}
	if purpose != "" {
		reasoningParts = append(reasoningParts, fmt.Sprintf("payment purpose: %s", purpose))
	}

	if len(reasoningParts) > 0 {
		analysis.AIReasoning = "Successfully " + strings.Join(reasoningParts, ", ")
	} else {
		analysis.AIReasoning = "Could not extract clear payment details from prompt"
	}

	return analysis
}

func (s *AIService) validateSpendingLimits(userID uuid.UUID, amount float64) error {
	limits, err := s.GetSpendingLimits(userID)
	if err != nil {
		return err
	}

	// Check transaction limit
	if amount > limits.TransactionLimit {
		return fmt.Errorf("amount exceeds transaction limit of Rs. %.2f", limits.TransactionLimit)
	}

	// Check daily limit
	remainingDaily, err := s.calculateRemainingDailyLimit(userID)
	if err != nil {
		return err
	}
	if amount > remainingDaily {
		return fmt.Errorf("amount exceeds remaining daily limit of Rs. %.2f", remainingDaily)
	}

	return nil
}

func (s *AIService) assessRiskLevel(amount float64, merchant string) string {
	if amount > 5000 {
		return "high"
	}
	if amount > 1000 {
		return "medium"
	}
	return "low"
}

func (s *AIService) calculateRemainingDailyLimit(userID uuid.UUID) (float64, error) {
	limits, err := s.GetSpendingLimits(userID)
	if err != nil {
		return 0, err
	}

	today := time.Now().Truncate(24 * time.Hour)
	var dailySpent float64

	if err := s.db.Model(&models.AIPaymentRequest{}).
		Where("user_id = ? AND status = 'processed' AND created_at >= ?", userID, today).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&dailySpent).Error; err != nil {
		return 0, fmt.Errorf("failed to calculate daily spending: %v", err)
	}

	return limits.DailyLimit - dailySpent, nil
}

func (s *AIService) updateSpendingTracker(tx *gorm.DB, userID uuid.UUID, amount float64) error {
	today := time.Now().Truncate(24 * time.Hour)
	firstOfMonth := time.Date(today.Year(), today.Month(), 1, 0, 0, 0, 0, today.Location())

	var tracker models.AISpendingTracker
	err := tx.Where("user_id = ? AND date = ?", userID, today).First(&tracker).Error

	if err == gorm.ErrRecordNotFound {
		// Calculate monthly spending properly
		var monthlySpent float64
		tx.Model(&models.AIPaymentRequest{}).
			Where("user_id = ? AND status = 'processed' AND created_at >= ? AND created_at < ?",
				userID, firstOfMonth, firstOfMonth.AddDate(0, 1, 0)).
			Select("COALESCE(SUM(amount), 0)").
			Scan(&monthlySpent)

		// Create new tracker for today
		tracker = models.AISpendingTracker{
			UserID:           userID,
			Date:             today,
			DailySpent:       amount,
			TransactionCount: 1,
			MonthlySpent:     monthlySpent + amount,
		}
		return tx.Create(&tracker).Error
	} else if err != nil {
		return err
	}

	// Update existing tracker
	tracker.DailySpent += amount
	tracker.TransactionCount++

	// Recalculate monthly spending properly
	var monthlySpent float64
	tx.Model(&models.AIPaymentRequest{}).
		Where("user_id = ? AND status = 'processed' AND created_at >= ? AND created_at < ?",
			userID, firstOfMonth, firstOfMonth.AddDate(0, 1, 0)).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&monthlySpent)
	tracker.MonthlySpent = monthlySpent + amount

	return tx.Save(&tracker).Error
}

func (s *AIService) generateSuggestions(analysis *PaymentAnalysis, walletBalance, remainingLimit float64) string {
	suggestions := []string{}

	if analysis.Amount > walletBalance {
		suggestions = append(suggestions, "Insufficient wallet balance. Consider adding funds.")
	}

	if analysis.Amount > remainingLimit*0.8 {
		suggestions = append(suggestions, "This transaction will use most of your daily limit.")
	}

	if analysis.Confidence < 0.8 {
		suggestions = append(suggestions, "Please verify the payment details as AI confidence is low.")
	}

	if len(suggestions) == 0 {
		suggestions = append(suggestions, "Payment details look good. Proceed with confidence.")
	}

	return strings.Join(suggestions, " ")
}
