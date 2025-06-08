package repositories

import (
	"context"
	"fmt"

	"github.com/zeusnotfound04/Tranza/models"
	"gorm.io/gorm"
)

type CardRepository struct {
	DB *gorm.DB
}


func NewCardRepository(db *gorm.DB) *CardRepository {
	return &CardRepository{DB: db}
}

func (r *CardRepository) FindByID(ctx context.Context, cardID uint) (models.LinkedCard , error){
	var card *models.LinkedCard
	if err := r.DB.WithContext(ctx).
	Where("id = ?", cardID).
		First(&card).Error ; err != nil {
		 return models.LinkedCard{} , fmt.Errorf("Cannot find the Linked Card with this Id: %w",err)
	}

	return *card , nil
}


func (r *CardRepository) Create(ctx context.Context,card *models.LinkedCard) error {
	return r.DB.WithContext(ctx).Create(card).Error
}

func (r *CardRepository) GetByUser(ctx context.Context,userID uint) ([]models.LinkedCard, error) {
	var cards []models.LinkedCard
	err := r.DB.WithContext(ctx).Where("user_id = ?" , userID).Find(&cards).Error
	return cards, err
}

func (r *CardRepository) Delete(ctx context.Context,cardID, userID uint) error {
	return r.DB.WithContext(ctx).Where("id = ? AND user_id = ?", cardID, userID).Delete(&models.LinkedCard{}).Error
}

// func (r *CardRepository) UpdateLimit(ctx context.Context , card *models.LinkedCard) error {
// 	return r.DB.WithContext(ctx).Model(&models.LinkedCard{}).
// 		Where("id = ? AND user_id = ?", card.ID, card.UserId).
// 		Update("limit", card.Limit).Error
// }


func (r *CardRepository) UpdateLimit(ctx context.Context ,cardID uint, userID uint , limit float64) error {
	return r.DB.WithContext(ctx).Model(&models.LinkedCard{}).
		Where("id = ? AND user_id = ?", cardID, userID).
		Update("limit", limit).Error
}