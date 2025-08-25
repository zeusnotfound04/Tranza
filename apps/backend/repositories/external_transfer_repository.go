package repositories

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/zeusnotfound04/Tranza/models"
	"gorm.io/gorm"
)

type ExternalTransferRepository struct {
	db *gorm.DB
}

func NewExternalTransferRepository(db *gorm.DB) *ExternalTransferRepository {
	return &ExternalTransferRepository{db: db}
}

// Create creates a new external transfer
func (r *ExternalTransferRepository) Create(transfer *models.ExternalTransfer) (*models.ExternalTransfer, error) {
	if err := r.db.Create(transfer).Error; err != nil {
		return nil, err
	}
	return transfer, nil
}

// GetByID retrieves an external transfer by ID
func (r *ExternalTransferRepository) GetByID(id uuid.UUID) (*models.ExternalTransfer, error) {
	var transfer models.ExternalTransfer
	if err := r.db.Where("id = ?", id).First(&transfer).Error; err != nil {
		return nil, err
	}
	return &transfer, nil
}

// GetByReferenceID retrieves an external transfer by reference ID
func (r *ExternalTransferRepository) GetByReferenceID(referenceID string) (*models.ExternalTransfer, error) {
	var transfer models.ExternalTransfer
	if err := r.db.Where("reference_id = ?", referenceID).First(&transfer).Error; err != nil {
		return nil, err
	}
	return &transfer, nil
}

// GetByRazorpayPayoutID retrieves an external transfer by Razorpay payout ID
func (r *ExternalTransferRepository) GetByRazorpayPayoutID(payoutID string) (*models.ExternalTransfer, error) {
	var transfer models.ExternalTransfer
	if err := r.db.Where("razorpay_payout_id = ?", payoutID).First(&transfer).Error; err != nil {
		return nil, err
	}
	return &transfer, nil
}

// GetByUserID retrieves external transfers by user ID
func (r *ExternalTransferRepository) GetByUserID(userID uuid.UUID) ([]*models.ExternalTransfer, error) {
	var transfers []*models.ExternalTransfer
	if err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&transfers).Error; err != nil {
		return nil, err
	}
	return transfers, nil
}

// GetByUserIDWithPagination retrieves external transfers by user ID with pagination and filters
func (r *ExternalTransferRepository) GetByUserIDWithPagination(
	userID uuid.UUID,
	limit, offset int,
	status string,
	dateFrom, dateTo *time.Time,
) ([]*models.ExternalTransfer, int64, error) {
	var transfers []*models.ExternalTransfer
	var total int64

	query := r.db.Where("user_id = ?", userID)

	// Apply filters
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if dateFrom != nil {
		query = query.Where("created_at >= ?", *dateFrom)
	}

	if dateTo != nil {
		query = query.Where("created_at <= ?", *dateTo)
	}

	// Get total count
	if err := query.Model(&models.ExternalTransfer{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get transfers with pagination
	if err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&transfers).Error; err != nil {
		return nil, 0, err
	}

	return transfers, total, nil
}

// GetByWalletID retrieves external transfers by wallet ID
func (r *ExternalTransferRepository) GetByWalletID(walletID uuid.UUID) ([]*models.ExternalTransfer, error) {
	var transfers []*models.ExternalTransfer
	if err := r.db.Where("wallet_id = ?", walletID).Order("created_at DESC").Find(&transfers).Error; err != nil {
		return nil, err
	}
	return transfers, nil
}

// GetByStatus retrieves external transfers by status
func (r *ExternalTransferRepository) GetByStatus(status string) ([]*models.ExternalTransfer, error) {
	var transfers []*models.ExternalTransfer
	if err := r.db.Where("status = ?", status).Order("created_at DESC").Find(&transfers).Error; err != nil {
		return nil, err
	}
	return transfers, nil
}

// GetPendingTransfers retrieves all pending transfers
func (r *ExternalTransferRepository) GetPendingTransfers() ([]*models.ExternalTransfer, error) {
	var transfers []*models.ExternalTransfer
	if err := r.db.Where("status IN ?", []string{
		models.ExternalTransferStatusPending,
		models.ExternalTransferStatusProcessing,
	}).Order("created_at ASC").Find(&transfers).Error; err != nil {
		return nil, err
	}
	return transfers, nil
}

// GetFailedTransfersForRetry retrieves failed transfers that can be retried
func (r *ExternalTransferRepository) GetFailedTransfersForRetry() ([]*models.ExternalTransfer, error) {
	var transfers []*models.ExternalTransfer
	if err := r.db.Where("status = ? AND retry_count < max_retries",
		models.ExternalTransferStatusFailed).Order("created_at ASC").Find(&transfers).Error; err != nil {
		return nil, err
	}
	return transfers, nil
}

// Update updates an external transfer
func (r *ExternalTransferRepository) Update(tx *gorm.DB, transfer *models.ExternalTransfer) error {
	db := r.db
	if tx != nil {
		db = tx
	}
	return db.Save(transfer).Error
}

// UpdateStatus updates the status of an external transfer
func (r *ExternalTransferRepository) UpdateStatus(transferID uuid.UUID, status string, failureReason string) error {
	updates := map[string]interface{}{
		"status":     status,
		"updated_at": time.Now(),
	}

	if status == models.ExternalTransferStatusProcessing {
		updates["processed_at"] = time.Now()
	}

	if status == models.ExternalTransferStatusSuccess ||
		status == models.ExternalTransferStatusFailed ||
		status == models.ExternalTransferStatusCancelled {
		updates["completed_at"] = time.Now()
	}

	if failureReason != "" {
		updates["failure_reason"] = failureReason
	}

	return r.db.Model(&models.ExternalTransfer{}).Where("id = ?", transferID).Updates(updates).Error
}

// UpdateRazorpayPayoutID updates the Razorpay payout ID
func (r *ExternalTransferRepository) UpdateRazorpayPayoutID(transferID uuid.UUID, payoutID string) error {
	return r.db.Model(&models.ExternalTransfer{}).Where("id = ?", transferID).Updates(map[string]interface{}{
		"razorpay_payout_id": payoutID,
		"status":             models.ExternalTransferStatusProcessing,
		"processed_at":       time.Now(),
		"updated_at":         time.Now(),
	}).Error
}

// IncrementRetryCount increments the retry count for a transfer
func (r *ExternalTransferRepository) IncrementRetryCount(transferID uuid.UUID) error {
	return r.db.Model(&models.ExternalTransfer{}).Where("id = ?", transferID).Updates(map[string]interface{}{
		"retry_count": gorm.Expr("retry_count + 1"),
		"updated_at":  time.Now(),
	}).Error
}

// GetTransferSummary gets summary statistics for transfers
func (r *ExternalTransferRepository) GetTransferSummary(userID uuid.UUID, dateFrom, dateTo *time.Time) (*TransferSummary, error) {
	query := r.db.Model(&models.ExternalTransfer{}).Where("user_id = ?", userID)

	if dateFrom != nil {
		query = query.Where("created_at >= ?", *dateFrom)
	}
	if dateTo != nil {
		query = query.Where("created_at <= ?", *dateTo)
	}

	var summary TransferSummary

	// Get total count and amount
	if err := query.Select("COUNT(*) as total_count, COALESCE(SUM(amount), 0) as total_amount").
		Scan(&summary).Error; err != nil {
		return nil, err
	}

	// Get status-wise counts
	var statusCounts []StatusCount
	if err := query.Select("status, COUNT(*) as count").
		Group("status").Scan(&statusCounts).Error; err != nil {
		return nil, err
	}

	// Map status counts
	for _, sc := range statusCounts {
		switch sc.Status {
		case models.ExternalTransferStatusSuccess:
			summary.SuccessfulCount = sc.Count
		case models.ExternalTransferStatusFailed:
			summary.FailedCount = sc.Count
		case models.ExternalTransferStatusPending:
			summary.PendingCount = sc.Count
		case models.ExternalTransferStatusProcessing:
			summary.ProcessingCount = sc.Count
		}
	}

	// Calculate success rate
	if summary.TotalCount > 0 {
		summary.SuccessRate = float64(summary.SuccessfulCount) / float64(summary.TotalCount) * 100
	}

	return &summary, nil
}

// GetDailyTransferLimitUsed calculates the daily transfer amount used
func (r *ExternalTransferRepository) GetDailyTransferLimitUsed(userID uuid.UUID, date time.Time) (decimal.Decimal, error) {
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	var totalAmount decimal.Decimal
	row := r.db.Model(&models.ExternalTransfer{}).
		Where("user_id = ? AND created_at >= ? AND created_at < ? AND status IN ?",
			userID, startOfDay, endOfDay,
			[]string{models.ExternalTransferStatusSuccess, models.ExternalTransferStatusPending, models.ExternalTransferStatusProcessing}).
		Select("COALESCE(SUM(total_amount), 0)").Row()

	if err := row.Scan(&totalAmount); err != nil {
		return decimal.Zero, err
	}

	return totalAmount, nil
}

// GetMonthlyTransferLimitUsed calculates the monthly transfer amount used
func (r *ExternalTransferRepository) GetMonthlyTransferLimitUsed(userID uuid.UUID, date time.Time) (decimal.Decimal, error) {
	startOfMonth := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, 0)

	var totalAmount decimal.Decimal
	row := r.db.Model(&models.ExternalTransfer{}).
		Where("user_id = ? AND created_at >= ? AND created_at < ? AND status IN ?",
			userID, startOfMonth, endOfMonth,
			[]string{models.ExternalTransferStatusSuccess, models.ExternalTransferStatusPending, models.ExternalTransferStatusProcessing}).
		Select("COALESCE(SUM(total_amount), 0)").Row()

	if err := row.Scan(&totalAmount); err != nil {
		return decimal.Zero, err
	}

	return totalAmount, nil
}

// Delete soft deletes an external transfer
func (r *ExternalTransferRepository) Delete(transferID uuid.UUID) error {
	return r.db.Delete(&models.ExternalTransfer{}, transferID).Error
}

// Helper structs for summary data
type TransferSummary struct {
	TotalCount      int             `json:"total_count"`
	TotalAmount     decimal.Decimal `json:"total_amount"`
	SuccessfulCount int             `json:"successful_count"`
	FailedCount     int             `json:"failed_count"`
	PendingCount    int             `json:"pending_count"`
	ProcessingCount int             `json:"processing_count"`
	SuccessRate     float64         `json:"success_rate"`
}

type StatusCount struct {
	Status string `json:"status"`
	Count  int    `json:"count"`
}
