package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/zeusnotfound04/Tranza/models"
	"gorm.io/gorm"
)

type APIKeyRepository struct {
	DB *gorm.DB
}

func NewAPIKeyRepository(db *gorm.DB) *APIKeyRepository {
	return &APIKeyRepository{
		DB: db,
	}
}

func (r *APIKeyRepository) Create(ctx context.Context, key *models.APIKey) error {
	return r.DB.WithContext(ctx).Create(key).Error
}

func (r *APIKeyRepository) FindByHash(ctx context.Context, hash string) (*models.APIKey, error) {
	var key models.APIKey
	err := r.DB.WithContext(ctx).Where("key_hash = ? AND is_active = TRUE", hash).First(&key).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &key, err

}

func (r *APIKeyRepository) UpdateUsage(ctx context.Context, keyID uint) error {
	return r.DB.WithContext(ctx).Model(&models.APIKey{}).
		Where("id = ?", keyID).
		Updates(map[string]interface{}{
			"last_used_at": time.Now(),
			"usage_count":  gorm.Expr("usage_count + 1"),
		}).Error
}

// FindByID finds an API key by its ID
func (r *APIKeyRepository) FindByID(ctx context.Context, keyID uint) (*models.APIKey, error) {
	var key models.APIKey
	err := r.DB.WithContext(ctx).Where("id = ?", keyID).First(&key).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &key, err
}

// UpdateKeyHash updates the key hash for rotation
func (r *APIKeyRepository) UpdateKeyHash(ctx context.Context, keyID uint, newHash string) error {
	return r.DB.WithContext(ctx).Model(&models.APIKey{}).
		Where("id = ?", keyID).
		Update("key_hash", newHash).Error
}

// GetByUserID returns all API keys for a user
func (r *APIKeyRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]models.APIKey, error) {
	var keys []models.APIKey
	err := r.DB.WithContext(ctx).Where("user_id = ?", userID).Find(&keys).Error
	return keys, err
}

// GetBotKeys returns all bot API keys for a workspace
func (r *APIKeyRepository) GetBotKeys(ctx context.Context, workspaceID string) ([]models.APIKey, error) {
	var keys []models.APIKey
	err := r.DB.WithContext(ctx).Where("key_type = ? AND bot_workspace = ? AND is_active = TRUE", "bot", workspaceID).Find(&keys).Error
	return keys, err
}

func (r *APIKeyRepository) RevokeByID(ctx context.Context, keyID uint, userID uuid.UUID) error {
	return r.DB.WithContext(ctx).
		Model(&models.APIKey{}).
		Where("id = ? AND user_id = ?", keyID, userID).
		Update("is_active", false).Error
}
