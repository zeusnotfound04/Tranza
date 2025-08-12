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

type WalletRepository struct {
	db *gorm.DB
}

func NewWalletRepository(db *gorm.DB) *WalletRepository {
	return &WalletRepository{
		db: db,
	}
}

// Create creates a new wallet for a user
func (r *WalletRepository) Create(wallet *models.Wallet) (*models.Wallet, error) {
	// Set default values
	if wallet.Balance.IsZero() {
		wallet.Balance = decimal.Zero
	}
	if wallet.Currency == "" {
		wallet.Currency = "INR"
	}
	if wallet.Status == "" {
		wallet.Status = "active"
	}
	if wallet.DailyLimit.IsZero() {
		wallet.DailyLimit = decimal.NewFromInt(10000) // ₹10,000
	}
	if wallet.MonthlyLimit.IsZero() {
		wallet.MonthlyLimit = decimal.NewFromInt(100000) // ₹1,00,000
	}
	if wallet.AIDailyLimit.IsZero() {
		wallet.AIDailyLimit = decimal.NewFromInt(1000) // ₹1,000
	}
	if wallet.AIPerTransactionLimit.IsZero() {
		wallet.AIPerTransactionLimit = decimal.NewFromInt(500) // ₹500
	}

	if err := r.db.Create(wallet).Error; err != nil {
		return nil, fmt.Errorf("failed to create wallet: %w", err)
	}
	return wallet, nil
}

// GetByID retrieves a wallet by ID
func (r *WalletRepository) GetByID(id uuid.UUID) (*models.Wallet, error) {
	var wallet models.Wallet
	if err := r.db.Where("id = ?", id).First(&wallet).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("wallet not found")
		}
		return nil, fmt.Errorf("failed to get wallet: %w", err)
	}
	return &wallet, nil
}

// GetByUserID retrieves a wallet by user ID
func (r *WalletRepository) GetByUserID(userID uuid.UUID) (*models.Wallet, error) {
	var wallet models.Wallet
	if err := r.db.Where("user_id = ?", userID).First(&wallet).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("wallet not found")
		}
		return nil, fmt.Errorf("failed to get wallet by user ID: %w", err)
	}
	return &wallet, nil
}

// GetByUserIDWithUser retrieves wallet with user details
func (r *WalletRepository) GetByUserIDWithUser(userID uuid.UUID) (*models.Wallet, error) {
	var wallet models.Wallet
	if err := r.db.Preload("User").Where("user_id = ?", userID).First(&wallet).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("wallet not found")
		}
		return nil, fmt.Errorf("failed to get wallet with user: %w", err)
	}
	return &wallet, nil
}

// Update updates a wallet
func (r *WalletRepository) Update(wallet *models.Wallet) error {
	wallet.UpdatedAt = time.Now()
	if err := r.db.Save(wallet).Error; err != nil {
		return fmt.Errorf("failed to update wallet: %w", err)
	}
	return nil
}

// UpdateBalance updates wallet balance (thread-safe with transaction)
func (r *WalletRepository) UpdateBalance(tx *gorm.DB, walletID uuid.UUID, newBalance decimal.Decimal) error {
	db := r.db
	if tx != nil {
		db = tx
	}

	result := db.Model(&models.Wallet{}).
		Where("id = ?", walletID).
		Updates(map[string]interface{}{
			"balance":    newBalance,
			"updated_at": time.Now(),
		})

	if result.Error != nil {
		return fmt.Errorf("failed to update wallet balance: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.New("wallet not found or balance not updated")
	}

	return nil
}

// IncrementBalance adds amount to wallet balance (atomic operation)
func (r *WalletRepository) IncrementBalance(tx *gorm.DB, walletID uuid.UUID, amount decimal.Decimal) (*models.Wallet, error) {
	db := r.db
	if tx != nil {
		db = tx
	}

	// Use raw SQL for atomic increment
	result := db.Model(&models.Wallet{}).
		Where("id = ?", walletID).
		Updates(map[string]interface{}{
			"balance":    gorm.Expr("balance + ?", amount),
			"updated_at": time.Now(),
		})

	if result.Error != nil {
		return nil, fmt.Errorf("failed to increment balance: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return nil, errors.New("wallet not found")
	}

	// Return updated wallet
	return r.GetByID(walletID)
}

// DecrementBalance subtracts amount from wallet balance (atomic operation)
func (r *WalletRepository) DecrementBalance(tx *gorm.DB, walletID uuid.UUID, amount decimal.Decimal) (*models.Wallet, error) {
	db := r.db
	if tx != nil {
		db = tx
	}

	// Check sufficient balance first
	var currentWallet models.Wallet
	if err := db.Where("id = ?", walletID).First(&currentWallet).Error; err != nil {
		return nil, fmt.Errorf("failed to get wallet: %w", err)
	}

	if currentWallet.Balance.LessThan(amount) {
		return nil, errors.New("insufficient balance")
	}

	// Use raw SQL for atomic decrement
	result := db.Model(&models.Wallet{}).
		Where("id = ? AND balance >= ?", walletID, amount).
		Updates(map[string]interface{}{
			"balance":    gorm.Expr("balance - ?", amount),
			"updated_at": time.Now(),
		})

	if result.Error != nil {
		return nil, fmt.Errorf("failed to decrement balance: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return nil, errors.New("insufficient balance or wallet not found")
	}

	// Return updated wallet
	return r.GetByID(walletID)
}

// UpdateAISettings updates AI-related settings
func (r *WalletRepository) UpdateAISettings(walletID uuid.UUID, settings *AISettings) error {
	updates := map[string]interface{}{
		"updated_at": time.Now(),
	}

	if settings.AIAccessEnabled != nil {
		updates["ai_access_enabled"] = *settings.AIAccessEnabled
	}
	if settings.AIDailyLimit != nil {
		updates["ai_daily_limit"] = *settings.AIDailyLimit
	}
	if settings.AIPerTransactionLimit != nil {
		updates["ai_per_transaction_limit"] = *settings.AIPerTransactionLimit
	}

	if err := r.db.Model(&models.Wallet{}).
		Where("id = ?", walletID).
		Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to update AI settings: %w", err)
	}

	return nil
}

// UpdateStatus updates wallet status
func (r *WalletRepository) UpdateStatus(walletID uuid.UUID, status string) error {
	if err := r.db.Model(&models.Wallet{}).
		Where("id = ?", walletID).
		Updates(map[string]interface{}{
			"status":     status,
			"updated_at": time.Now(),
		}).Error; err != nil {
		return fmt.Errorf("failed to update wallet status: %w", err)
	}
	return nil
}

// UpdateRazorpayCustomerID updates Razorpay customer ID
func (r *WalletRepository) UpdateRazorpayCustomerID(walletID uuid.UUID, customerID string) error {
	if err := r.db.Model(&models.Wallet{}).
		Where("id = ?", walletID).
		Updates(map[string]interface{}{
			"razorpay_customer_id": customerID,
			"updated_at":          time.Now(),
		}).Error; err != nil {
		return fmt.Errorf("failed to update Razorpay customer ID: %w", err)
	}
	return nil
}

// GetActiveWallets retrieves all active wallets
func (r *WalletRepository) GetActiveWallets() ([]*models.Wallet, error) {
	var wallets []*models.Wallet
	if err := r.db.Where("status = ?", "active").Find(&wallets).Error; err != nil {
		return nil, fmt.Errorf("failed to get active wallets: %w", err)
	}
	return wallets, nil
}

// GetWalletsByStatus retrieves wallets by status
func (r *WalletRepository) GetWalletsByStatus(status string) ([]*models.Wallet, error) {
	var wallets []*models.Wallet
	if err := r.db.Where("status = ?", status).Find(&wallets).Error; err != nil {
		return nil, fmt.Errorf("failed to get wallets by status: %w", err)
	}
	return wallets, nil
}

// GetWalletsWithBalance retrieves wallets with balance greater than specified amount
func (r *WalletRepository) GetWalletsWithBalance(minBalance decimal.Decimal) ([]*models.Wallet, error) {
	var wallets []*models.Wallet
	if err := r.db.Where("balance > ? AND status = ?", minBalance, "active").Find(&wallets).Error; err != nil {
		return nil, fmt.Errorf("failed to get wallets with balance: %w", err)
	}
	return wallets, nil
}

// GetWalletsWithAIEnabled retrieves wallets with AI access enabled
func (r *WalletRepository) GetWalletsWithAIEnabled() ([]*models.Wallet, error) {
	var wallets []*models.Wallet
	if err := r.db.Where("ai_access_enabled = ? AND status = ?", true, "active").Find(&wallets).Error; err != nil {
		return nil, fmt.Errorf("failed to get AI-enabled wallets: %w", err)
	}
	return wallets, nil
}

// GetTotalBalance calculates total balance across all active wallets
func (r *WalletRepository) GetTotalBalance() (decimal.Decimal, error) {
	var result struct {
		Total decimal.Decimal
	}

	if err := r.db.Model(&models.Wallet{}).
		Select("COALESCE(SUM(balance), 0) as total").
		Where("status = ?", "active").
		Scan(&result).Error; err != nil {
		return decimal.Zero, fmt.Errorf("failed to calculate total balance: %w", err)
	}

	return result.Total, nil
}

// GetWalletStatistics returns comprehensive wallet statistics
func (r *WalletRepository) GetWalletStatistics() (*WalletStatistics, error) {
	var stats WalletStatistics

	// Total wallets
	if err := r.db.Model(&models.Wallet{}).Count(&stats.TotalWallets).Error; err != nil {
		return nil, fmt.Errorf("failed to count total wallets: %w", err)
	}

	// Active wallets
	if err := r.db.Model(&models.Wallet{}).
		Where("status = ?", "active").
		Count(&stats.ActiveWallets).Error; err != nil {
		return nil, fmt.Errorf("failed to count active wallets: %w", err)
	}

	// Frozen wallets
	if err := r.db.Model(&models.Wallet{}).
		Where("status = ?", "frozen").
		Count(&stats.FrozenWallets).Error; err != nil {
		return nil, fmt.Errorf("failed to count frozen wallets: %w", err)
	}

	// AI enabled wallets
	if err := r.db.Model(&models.Wallet{}).
		Where("ai_access_enabled = ? AND status = ?", true, "active").
		Count(&stats.AIEnabledWallets).Error; err != nil {
		return nil, fmt.Errorf("failed to count AI enabled wallets: %w", err)
	}

	// Total balance
	var balanceResult struct {
		Total decimal.Decimal
	}
	if err := r.db.Model(&models.Wallet{}).
		Select("COALESCE(SUM(balance), 0) as total").
		Where("status = ?", "active").
		Scan(&balanceResult).Error; err != nil {
		return nil, fmt.Errorf("failed to calculate total balance: %w", err)
	}
	stats.TotalBalance = balanceResult.Total

	// Average balance
	if stats.ActiveWallets > 0 {
		stats.AverageBalance = stats.TotalBalance.Div(decimal.NewFromInt(stats.ActiveWallets))
	}

	return &stats, nil
}

// IsWalletExists checks if wallet exists for user
func (r *WalletRepository) IsWalletExists(userID uuid.UUID) (bool, error) {
	var count int64
	if err := r.db.Model(&models.Wallet{}).
		Where("user_id = ?", userID).
		Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check wallet existence: %w", err)
	}
	return count > 0, nil
}

// GetWalletWithSettings retrieves wallet with settings
func (r *WalletRepository) GetWalletWithSettings(userID uuid.UUID) (*models.Wallet, error) {
	var wallet models.Wallet
	if err := r.db.Preload("WalletSettings").
		Where("user_id = ?", userID).
		First(&wallet).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("wallet not found")
		}
		return nil, fmt.Errorf("failed to get wallet with settings: %w", err)
	}
	return &wallet, nil
}

// BulkUpdateStatus updates status for multiple wallets
func (r *WalletRepository) BulkUpdateStatus(walletIDs []uuid.UUID, status string) error {
	if err := r.db.Model(&models.Wallet{}).
		Where("id IN ?", walletIDs).
		Updates(map[string]interface{}{
			"status":     status,
			"updated_at": time.Now(),
		}).Error; err != nil {
		return fmt.Errorf("failed to bulk update wallet status: %w", err)
	}
	return nil
}

// GetWalletsCreatedInDateRange retrieves wallets created within date range
func (r *WalletRepository) GetWalletsCreatedInDateRange(startDate, endDate time.Time) ([]*models.Wallet, error) {
	var wallets []*models.Wallet
	if err := r.db.Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Order("created_at DESC").
		Find(&wallets).Error; err != nil {
		return nil, fmt.Errorf("failed to get wallets by date range: %w", err)
	}
	return wallets, nil
}

// GetLowBalanceWallets retrieves wallets with balance below threshold
func (r *WalletRepository) GetLowBalanceWallets(threshold decimal.Decimal) ([]*models.Wallet, error) {
	var wallets []*models.Wallet
	if err := r.db.Preload("User").
		Where("balance < ? AND status = ?", threshold, "active").
		Find(&wallets).Error; err != nil {
		return nil, fmt.Errorf("failed to get low balance wallets: %w", err)
	}
	return wallets, nil
}

// GetHighValueWallets retrieves wallets with balance above threshold
func (r *WalletRepository) GetHighValueWallets(threshold decimal.Decimal) ([]*models.Wallet, error) {
	var wallets []*models.Wallet
	if err := r.db.Preload("User").
		Where("balance > ? AND status = ?", threshold, "active").
		Order("balance DESC").
		Find(&wallets).Error; err != nil {
		return nil, fmt.Errorf("failed to get high value wallets: %w", err)
	}
	return wallets, nil
}

// GetInactiveWallets retrieves wallets with no transactions in specified duration
func (r *WalletRepository) GetInactiveWallets(inactiveDuration time.Duration) ([]*models.Wallet, error) {
	cutoffTime := time.Now().Add(-inactiveDuration)
	
	var wallets []*models.Wallet
	// This is a complex query - wallets that don't have any transactions after cutoff time
	if err := r.db.Preload("User").
		Where(`id NOT IN (
			SELECT DISTINCT wallet_id 
			FROM transactions 
			WHERE created_at > ? AND status = 'success'
		) AND status = 'active'`, cutoffTime).
		Find(&wallets).Error; err != nil {
		return nil, fmt.Errorf("failed to get inactive wallets: %w", err)
	}
	return wallets, nil
}

// SoftDelete soft deletes a wallet (sets status to closed)
func (r *WalletRepository) SoftDelete(walletID uuid.UUID) error {
	if err := r.db.Model(&models.Wallet{}).
		Where("id = ?", walletID).
		Updates(map[string]interface{}{
			"status":     "closed",
			"updated_at": time.Now(),
		}).Error; err != nil {
		return fmt.Errorf("failed to soft delete wallet: %w", err)
	}
	return nil
}

// RestoreWallet restores a soft deleted wallet
func (r *WalletRepository) RestoreWallet(walletID uuid.UUID) error {
	if err := r.db.Model(&models.Wallet{}).
		Where("id = ?", walletID).
		Updates(map[string]interface{}{
			"status":     "active",
			"updated_at": time.Now(),
		}).Error; err != nil {
		return fmt.Errorf("failed to restore wallet: %w", err)
	}
	return nil
}

// GetWalletsByBalanceRange retrieves wallets within balance range
func (r *WalletRepository) GetWalletsByBalanceRange(minBalance, maxBalance decimal.Decimal, limit int) ([]*models.Wallet, error) {
	var wallets []*models.Wallet
	query := r.db.Where("status = ?", "active")
	
	if !minBalance.IsZero() {
		query = query.Where("balance >= ?", minBalance)
	}
	if !maxBalance.IsZero() {
		query = query.Where("balance <= ?", maxBalance)
	}
	
	if err := query.Order("balance DESC").Limit(limit).Find(&wallets).Error; err != nil {
		return nil, fmt.Errorf("failed to get wallets by balance range: %w", err)
	}
	return wallets, nil
}

// Helper structs
type AISettings struct {
	AIAccessEnabled       *bool            `json:"ai_access_enabled"`
	AIDailyLimit          *decimal.Decimal `json:"ai_daily_limit"`
	AIPerTransactionLimit *decimal.Decimal `json:"ai_per_transaction_limit"`
}

type WalletStatistics struct {
	TotalWallets      int64           `json:"total_wallets"`
	ActiveWallets     int64           `json:"active_wallets"`
	FrozenWallets     int64           `json:"frozen_wallets"`
	AIEnabledWallets  int64           `json:"ai_enabled_wallets"`
	TotalBalance      decimal.Decimal `json:"total_balance"`
	AverageBalance    decimal.Decimal `json:"average_balance"`
}