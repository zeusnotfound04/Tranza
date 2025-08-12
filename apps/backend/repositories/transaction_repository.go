package repositories

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/zeusnotfound04/Tranza/models"
	"gorm.io/gorm"
)

type TransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{
		db: db,
	}
}

// Create creates a new transaction
func (r *TransactionRepository) Create(transaction *models.Transaction) (*models.Transaction, error) {
	if err := r.db.Create(transaction).Error; err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}
	return transaction, nil
}

// CreateWithTx creates a new transaction within a database transaction
func (r *TransactionRepository) CreateWithTx(tx *gorm.DB, transaction *models.Transaction) (*models.Transaction, error) {
	if err := tx.Create(transaction).Error; err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}
	return transaction, nil
}

// GetByID retrieves a transaction by ID
func (r *TransactionRepository) GetByID(id uuid.UUID) (*models.Transaction, error) {
	var transaction models.Transaction
	if err := r.db.Where("id = ?", id).First(&transaction).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("transaction not found")
		}
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}
	return &transaction, nil
}

// GetByIDAndUserID retrieves a transaction by ID and user ID (for security)
func (r *TransactionRepository) GetByIDAndUserID(id, userID uuid.UUID) (*models.Transaction, error) {
	var transaction models.Transaction
	if err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&transaction).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("transaction not found")
		}
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}
	return &transaction, nil
}

// GetByOrderID retrieves a transaction by Razorpay order ID
func (r *TransactionRepository) GetByOrderID(orderID string) (*models.Transaction, error) {
	var transaction models.Transaction
	if err := r.db.Where("razorpay_order_id = ?", orderID).First(&transaction).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("transaction not found")
		}
		return nil, fmt.Errorf("failed to get transaction by order ID: %w", err)
	}
	return &transaction, nil
}

// GetByPaymentID retrieves a transaction by Razorpay payment ID
func (r *TransactionRepository) GetByPaymentID(paymentID string) (*models.Transaction, error) {
	var transaction models.Transaction
	if err := r.db.Where("razorpay_payment_id = ?", paymentID).First(&transaction).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("transaction not found")
		}
		return nil, fmt.Errorf("failed to get transaction by payment ID: %w", err)
	}
	return &transaction, nil
}

// GetByReferenceID retrieves a transaction by reference ID
func (r *TransactionRepository) GetByReferenceID(referenceID string) (*models.Transaction, error) {
	var transaction models.Transaction
	if err := r.db.Where("reference_id = ?", referenceID).First(&transaction).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("transaction not found")
		}
		return nil, fmt.Errorf("failed to get transaction by reference ID: %w", err)
	}
	return &transaction, nil
}

// Update updates a transaction
func (r *TransactionRepository) Update(tx *gorm.DB, transaction *models.Transaction) error {
	db := r.db
	if tx != nil {
		db = tx
	}

	if err := db.Save(transaction).Error; err != nil {
		return fmt.Errorf("failed to update transaction: %w", err)
	}
	return nil
}

// GetByUserIDWithPagination retrieves transactions by user ID with pagination and filters
func (r *TransactionRepository) GetByUserIDWithPagination(userID uuid.UUID, limit, offset int, transactionType string) ([]*models.Transaction, int64, error) {
	var transactions []*models.Transaction
	var total int64

	query := r.db.Where("user_id = ?", userID)

	// Apply transaction type filter
	if transactionType != "" {
		query = query.Where("type = ?", transactionType)
	}

	// Get total count
	if err := query.Model(&models.Transaction{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count transactions: %w", err)
	}

	// Get transactions with pagination
	if err := query.Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&transactions).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get transactions: %w", err)
	}

	return transactions, total, nil
}

// GetTransactionsWithFilters retrieves transactions with advanced filters
func (r *TransactionRepository) GetTransactionsWithFilters(
	userID uuid.UUID,
	limit, offset int,
	transactionType, status string,
	startDate, endDate *time.Time,
	minAmount, maxAmount *decimal.Decimal,
) ([]*models.Transaction, int64, error) {
	var transactions []*models.Transaction
	var total int64

	query := r.db.Where("user_id = ?", userID)

	// Apply filters
	if transactionType != "" {
		query = query.Where("type = ?", transactionType)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if startDate != nil {
		query = query.Where("created_at >= ?", *startDate)
	}
	if endDate != nil {
		query = query.Where("created_at <= ?", *endDate)
	}
	if minAmount != nil {
		query = query.Where("amount >= ?", *minAmount)
	}
	if maxAmount != nil {
		query = query.Where("amount <= ?", *maxAmount)
	}

	// Get total count
	if err := query.Model(&models.Transaction{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count transactions: %w", err)
	}

	// Get transactions with pagination
	if err := query.Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&transactions).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get transactions: %w", err)
	}

	return transactions, total, nil
}

// GetByWalletID retrieves all transactions for a wallet
func (r *TransactionRepository) GetByWalletID(walletID uuid.UUID) ([]*models.Transaction, error) {
	var transactions []*models.Transaction
	if err := r.db.Where("wallet_id = ?", walletID).
		Order("created_at DESC").
		Find(&transactions).Error; err != nil {
		return nil, fmt.Errorf("failed to get transactions by wallet ID: %w", err)
	}
	return transactions, nil
}

// GetByWalletIDWithDateRange retrieves transactions within a date range
func (r *TransactionRepository) GetByWalletIDWithDateRange(walletID uuid.UUID, startDate, endDate time.Time) ([]*models.Transaction, error) {
	var transactions []*models.Transaction
	if err := r.db.Where("wallet_id = ? AND created_at BETWEEN ? AND ?", walletID, startDate, endDate).
		Order("created_at DESC").
		Find(&transactions).Error; err != nil {
		return nil, fmt.Errorf("failed to get transactions by date range: %w", err)
	}
	return transactions, nil
}

// GetSuccessfulTransactionsByWalletID retrieves only successful transactions
func (r *TransactionRepository) GetSuccessfulTransactionsByWalletID(walletID uuid.UUID) ([]*models.Transaction, error) {
	var transactions []*models.Transaction
	if err := r.db.Where("wallet_id = ? AND status = ?", walletID, "success").
		Order("created_at DESC").
		Find(&transactions).Error; err != nil {
		return nil, fmt.Errorf("failed to get successful transactions: %w", err)
	}
	return transactions, nil
}

// GetAITransactionsByWalletID retrieves only AI transactions
func (r *TransactionRepository) GetAITransactionsByWalletID(walletID uuid.UUID) ([]*models.Transaction, error) {
	var transactions []*models.Transaction
	if err := r.db.Where("wallet_id = ? AND type = ? AND status = ?", walletID, "ai_payment", "success").
		Order("created_at DESC").
		Find(&transactions).Error; err != nil {
		return nil, fmt.Errorf("failed to get AI transactions: %w", err)
	}
	return transactions, nil
}

// GetAIDailySpending calculates total AI spending for a specific date
func (r *TransactionRepository) GetAIDailySpending(walletID uuid.UUID, date string) (decimal.Decimal, error) {
	// Parse date
	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return decimal.Zero, errors.New("invalid date format")
	}

	// Calculate start and end of day
	startOfDay := parsedDate
	endOfDay := parsedDate.Add(24 * time.Hour).Add(-time.Nanosecond)

	// Query AI transactions for the day
	var result struct {
		Total decimal.Decimal
	}

	if err := r.db.Model(&models.Transaction{}).
		Select("COALESCE(SUM(amount), 0) as total").
		Where("wallet_id = ? AND type = ? AND status = ? AND created_at BETWEEN ? AND ?",
			walletID, "ai_payment", "success", startOfDay, endOfDay).
		Scan(&result).Error; err != nil {
		return decimal.Zero, fmt.Errorf("failed to calculate daily AI spending: %w", err)
	}

	return result.Total, nil
}

// GetAIWeeklySpending calculates total AI spending for current week
func (r *TransactionRepository) GetAIWeeklySpending(walletID uuid.UUID) (decimal.Decimal, error) {
	now := time.Now()
	startOfWeek := now.AddDate(0, 0, -int(now.Weekday()))
	startOfWeek = time.Date(startOfWeek.Year(), startOfWeek.Month(), startOfWeek.Day(), 0, 0, 0, 0, startOfWeek.Location())

	var result struct {
		Total decimal.Decimal
	}

	if err := r.db.Model(&models.Transaction{}).
		Select("COALESCE(SUM(amount), 0) as total").
		Where("wallet_id = ? AND type = ? AND status = ? AND created_at >= ?",
			walletID, "ai_payment", "success", startOfWeek).
		Scan(&result).Error; err != nil {
		return decimal.Zero, fmt.Errorf("failed to calculate weekly AI spending: %w", err)
	}

	return result.Total, nil
}

// GetAITransactionCount counts AI transactions for today
func (r *TransactionRepository) GetAITransactionCount(walletID uuid.UUID, date string) (int64, error) {
	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return 0, errors.New("invalid date format")
	}

	startOfDay := parsedDate
	endOfDay := parsedDate.Add(24 * time.Hour).Add(-time.Nanosecond)

	var count int64
	if err := r.db.Model(&models.Transaction{}).
		Where("wallet_id = ? AND type = ? AND status = ? AND created_at BETWEEN ? AND ?",
			walletID, "ai_payment", "success", startOfDay, endOfDay).
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count AI transactions: %w", err)
	}

	return count, nil
}

// GetTransactionStats retrieves transaction statistics for a user
func (r *TransactionRepository) GetTransactionStats(userID uuid.UUID) (*TransactionStats, error) {
	var stats TransactionStats

	// Total transactions
	if err := r.db.Model(&models.Transaction{}).
		Where("user_id = ? AND status = ?", userID, "success").
		Count(&stats.TotalTransactions).Error; err != nil {
		return nil, fmt.Errorf("failed to count total transactions: %w", err)
	}

	// Total amount
	var totalResult struct {
		Total decimal.Decimal
	}
	if err := r.db.Model(&models.Transaction{}).
		Select("COALESCE(SUM(amount), 0) as total").
		Where("user_id = ? AND status = ? AND type != ?", userID, "success", "load_money").
		Scan(&totalResult).Error; err != nil {
		return nil, fmt.Errorf("failed to calculate total amount: %w", err)
	}
	stats.TotalAmount = totalResult.Total

	// Today's transactions
	today := time.Now()
	startOfDay := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())
	endOfDay := startOfDay.Add(24 * time.Hour).Add(-time.Nanosecond)

	if err := r.db.Model(&models.Transaction{}).
		Where("user_id = ? AND status = ? AND created_at BETWEEN ? AND ?", userID, "success", startOfDay, endOfDay).
		Count(&stats.TodayTransactions).Error; err != nil {
		return nil, fmt.Errorf("failed to count today's transactions: %w", err)
	}

	// Today's amount
	var todayResult struct {
		Total decimal.Decimal
	}
	if err := r.db.Model(&models.Transaction{}).
		Select("COALESCE(SUM(amount), 0) as total").
		Where("user_id = ? AND status = ? AND type != ? AND created_at BETWEEN ? AND ?",
			userID, "success", "load_money", startOfDay, endOfDay).
		Scan(&todayResult).Error; err != nil {
		return nil, fmt.Errorf("failed to calculate today's amount: %w", err)
	}
	stats.TodayAmount = todayResult.Total

	// AI transactions
	if err := r.db.Model(&models.Transaction{}).
		Where("user_id = ? AND status = ? AND type = ?", userID, "success", "ai_payment").
		Count(&stats.AITransactions).Error; err != nil {
		return nil, fmt.Errorf("failed to count AI transactions: %w", err)
	}

	// AI amount
	var aiResult struct {
		Total decimal.Decimal
	}
	if err := r.db.Model(&models.Transaction{}).
		Select("COALESCE(SUM(amount), 0) as total").
		Where("user_id = ? AND status = ? AND type = ?", userID, "success", "ai_payment").
		Scan(&aiResult).Error; err != nil {
		return nil, fmt.Errorf("failed to calculate AI amount: %w", err)
	}
	stats.AIAmount = aiResult.Total

	// Last transaction date
	var lastTransaction models.Transaction
	if err := r.db.Where("user_id = ? AND status = ?", userID, "success").
		Order("created_at DESC").
		First(&lastTransaction).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("failed to get last transaction: %w", err)
		}
	} else {
		stats.LastTransactionDate = lastTransaction.CreatedAt
	}

	return &stats, nil
}

// GetTransactionsByStatus retrieves transactions by status
func (r *TransactionRepository) GetTransactionsByStatus(status string, limit int) ([]*models.Transaction, error) {
	var transactions []*models.Transaction
	if err := r.db.Where("status = ?", status).
		Order("created_at DESC").
		Limit(limit).
		Find(&transactions).Error; err != nil {
		return nil, fmt.Errorf("failed to get transactions by status: %w", err)
	}
	return transactions, nil
}

// GetPendingTransactions retrieves pending transactions older than specified duration
func (r *TransactionRepository) GetPendingTransactions(olderThan time.Duration) ([]*models.Transaction, error) {
	cutoffTime := time.Now().Add(-olderThan)
	var transactions []*models.Transaction

	if err := r.db.Where("status = ? AND created_at < ?", "pending", cutoffTime).
		Find(&transactions).Error; err != nil {
		return nil, fmt.Errorf("failed to get pending transactions: %w", err)
	}
	return transactions, nil
}

// UpdateStatus updates transaction status
func (r *TransactionRepository) UpdateStatus(id uuid.UUID, status string, failureReason string) error {
	updates := map[string]interface{}{
		"status":     status,
		"updated_at": time.Now(),
	}

	if failureReason != "" {
		updates["failure_reason"] = failureReason
	}

	if err := r.db.Model(&models.Transaction{}).
		Where("id = ?", id).
		Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to update transaction status: %w", err)
	}
	return nil
}

// BulkUpdateStatus updates status for multiple transactions
func (r *TransactionRepository) BulkUpdateStatus(ids []uuid.UUID, status string) error {
	if err := r.db.Model(&models.Transaction{}).
		Where("id IN ?", ids).
		Updates(map[string]interface{}{
			"status":     status,
			"updated_at": time.Now(),
		}).Error; err != nil {
		return fmt.Errorf("failed to bulk update transaction status: %w", err)
	}
	return nil
}

// DeleteOldTransactions deletes transactions older than specified duration (for cleanup)
func (r *TransactionRepository) DeleteOldTransactions(olderThan time.Duration) (int64, error) {
	cutoffTime := time.Now().Add(-olderThan)

	result := r.db.Where("created_at < ? AND status IN ?", cutoffTime, []string{"failed", "cancelled"}).
		Delete(&models.Transaction{})

	if result.Error != nil {
		return 0, fmt.Errorf("failed to delete old transactions: %w", result.Error)
	}

	return result.RowsAffected, nil
}

// GetTransactionSummaryByDateRange returns summary of transactions in date range
func (r *TransactionRepository) GetTransactionSummaryByDateRange(walletID uuid.UUID, startDate, endDate time.Time) (*TransactionSummary, error) {
	var summary TransactionSummary

	// Total transactions
	if err := r.db.Model(&models.Transaction{}).
		Where("wallet_id = ? AND status = ? AND created_at BETWEEN ? AND ?", walletID, "success", startDate, endDate).
		Count(&summary.TotalTransactions).Error; err != nil {
		return nil, fmt.Errorf("failed to count transactions: %w", err)
	}

	// Total inflow (load_money, refund)
	var inflowResult struct {
		Total decimal.Decimal
	}
	if err := r.db.Model(&models.Transaction{}).
		Select("COALESCE(SUM(amount), 0) as total").
		Where("wallet_id = ? AND status = ? AND type IN ? AND created_at BETWEEN ? AND ?",
			walletID, "success", []string{"load_money", "refund"}, startDate, endDate).
		Scan(&inflowResult).Error; err != nil {
		return nil, fmt.Errorf("failed to calculate inflow: %w", err)
	}
	summary.TotalInflow = inflowResult.Total

	// Total outflow (ai_payment, withdrawal)
	var outflowResult struct {
		Total decimal.Decimal
	}
	if err := r.db.Model(&models.Transaction{}).
		Select("COALESCE(SUM(amount), 0) as total").
		Where("wallet_id = ? AND status = ? AND type IN ? AND created_at BETWEEN ? AND ?",
			walletID, "success", []string{"ai_payment", "withdrawal"}, startDate, endDate).
		Scan(&outflowResult).Error; err != nil {
		return nil, fmt.Errorf("failed to calculate outflow: %w", err)
	}
	summary.TotalOutflow = outflowResult.Total

	// Calculate net flow
	summary.NetFlow = summary.TotalInflow.Sub(summary.TotalOutflow)

	// Average transaction amount
	if summary.TotalTransactions > 0 {
		summary.AverageAmount = summary.TotalInflow.Add(summary.TotalOutflow).Div(decimal.NewFromInt(summary.TotalTransactions))
	}

	return &summary, nil
}

// Helper structs for statistics
type TransactionStats struct {
	TotalTransactions   int64           `json:"total_transactions"`
	TotalAmount         decimal.Decimal `json:"total_amount"`
	TodayTransactions   int64           `json:"today_transactions"`
	TodayAmount         decimal.Decimal `json:"today_amount"`
	AITransactions      int64           `json:"ai_transactions"`
	AIAmount            decimal.Decimal `json:"ai_amount"`
	LastTransactionDate time.Time       `json:"last_transaction_date"`
}

type TransactionSummary struct {
	TotalTransactions int64           `json:"total_transactions"`
	TotalInflow       decimal.Decimal `json:"total_inflow"`
	TotalOutflow      decimal.Decimal `json:"total_outflow"`
	NetFlow           decimal.Decimal `json:"net_flow"`
	AverageAmount     decimal.Decimal `json:"average_amount"`
}
