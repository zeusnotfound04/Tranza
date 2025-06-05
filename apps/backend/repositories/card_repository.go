package repositories

import (
	"github.com/zeusnotfound04/Tranza/models"
	"gorm.io/gorm"
)

type CardRepository struct {
	DB *gorm.DB
}


func NewCardRepository(db *gorm.DB) *CardRepository {
	return &CardRepository{DB: db}
}


func (r *CardRepository) Create(card *models.LinkedCard) error {
	return r.DB.Create(card).Error
}

func (r *CardRepository) GetByUser(userID uint) ([]models.LinkedCard, error) {
	var cards []models.LinkedCard
	err := r.DB.Where("user_id = ?" , userID).Find(&cards).Error
	return cards, err
}

func (r *CardRepository) Delete(cardID, userID uint) error {
	return r.DB.Where("id = ? AND user_id = ?", cardID, userID).Delete(&models.LinkedCard{}).Error
}

func (r *CardRepository) UpdateLimit(cardID, userID uint, limit float64) error {
	return r.DB.Model(&models.LinkedCard{}).
		Where("id = ? AND user_id = ?", cardID, userID).
		Update("limit", limit).Error
}