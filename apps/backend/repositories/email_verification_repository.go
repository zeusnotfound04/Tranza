package repositories

import (
	"context"
	"time"

	"github.com/zeusnotfound04/Tranza/models"
	"gorm.io/gorm"
)

// EmailVerificationRepository handles email verification database operations
type EmailVerificationRepository struct {
	db *gorm.DB
}

// NewEmailVerificationRepository creates a new email verification repository
func NewEmailVerificationRepository(db *gorm.DB) *EmailVerificationRepository {
	return &EmailVerificationRepository{db: db}
}

// CreateVerification creates a new email verification record
func (evr *EmailVerificationRepository) CreateVerification(ctx context.Context, verification *models.EmailVerification) error {
	return evr.db.WithContext(ctx).Create(verification).Error
}

// GetVerificationByEmail retrieves verification record by email
func (evr *EmailVerificationRepository) GetVerificationByEmail(ctx context.Context, email string) (*models.EmailVerification, error) {
	var verification models.EmailVerification
	err := evr.db.WithContext(ctx).Where("email = ? AND is_verified = false", email).First(&verification).Error
	if err != nil {
		return nil, err
	}
	return &verification, nil
}

// UpdateVerification updates a verification record
func (evr *EmailVerificationRepository) UpdateVerification(ctx context.Context, verification *models.EmailVerification) error {
	return evr.db.WithContext(ctx).Save(verification).Error
}

// DeleteVerification deletes a verification record
func (evr *EmailVerificationRepository) DeleteVerification(ctx context.Context, email string) error {
	return evr.db.WithContext(ctx).Where("email = ?", email).Delete(&models.EmailVerification{}).Error
}

// DeleteExpiredVerifications deletes expired verification records
func (evr *EmailVerificationRepository) DeleteExpiredVerifications(ctx context.Context) error {
	return evr.db.WithContext(ctx).Where("expires_at < ?", time.Now()).Delete(&models.EmailVerification{}).Error
}

// IncrementAttempts increments the verification attempts counter
func (evr *EmailVerificationRepository) IncrementAttempts(ctx context.Context, email string) error {
	return evr.db.WithContext(ctx).Model(&models.EmailVerification{}).
		Where("email = ?", email).
		UpdateColumn("attempts", gorm.Expr("attempts + 1")).Error
}

// CheckEmailExists checks if email already exists in users table
func (evr *EmailVerificationRepository) CheckEmailExists(ctx context.Context, email string) (bool, error) {
	var count int64
	err := evr.db.WithContext(ctx).Model(&models.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

// CheckUsernameExists checks if username already exists in users table
func (evr *EmailVerificationRepository) CheckUsernameExists(ctx context.Context, username string) (bool, error) {
	var count int64
	err := evr.db.WithContext(ctx).Model(&models.User{}).Where("username = ?", username).Count(&count).Error
	return count > 0, err
}