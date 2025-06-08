package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/zeusnotfound04/Tranza/models"
	"gorm.io/gorm"
)



type APIKeyRepository struct {
	DB *gorm.DB
}



func NewAPIKeyRepository(db *gorm.DB)  *APIKeyRepository {
	return &APIKeyRepository{
		DB: db ,
	}
}


func (r *APIKeyRepository) Create (ctx context.Context , key *models.APIKey) error {
	return r.DB.WithContext(ctx).Create(key).Error
}


func (r *APIKeyRepository) FindByHash (ctx context.Context , hash string)  (*models.APIKey , error ) {
	var key models.APIKey
	err := r.DB.WithContext(ctx).Where("key_hash = ? AND is_active = TRUE", hash).First(&key).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &key, err

}


func (r *APIKeyRepository) UpdateUsage(ctx context.Context ,keyID uint)  error {
	return r.DB.WithContext(ctx).Model(&models.APIKey{}).
		Where("id = ?", keyID).
		Update("last_used_at", time.Now()).Error
}


func (r *APIKeyRepository) RevokeByID(ctx context.Context, keyID uint, userID uint) error {
   	return r.DB.WithContext(ctx).
		Model(&models.APIKey{}).
		Where("id = ? AND user_id = ?", keyID, userID).
		Update("is_active", false).Error	
}