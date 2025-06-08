package services

import (
	"context"

	"github.com/zeusnotfound04/Tranza/models"
	"github.com/zeusnotfound04/Tranza/repositories"
)

type CardService struct {
	Repo *repositories.CardRepository
}


func NewCardService(repo *repositories.CardRepository) *CardService {
	return &CardService{Repo: repo}
}

func (s *CardService) LinkCard(ctx context.Context , card *models.LinkedCard) error {
	// Add validations if needed (e.g., max limit check)
	return s.Repo.Create(ctx , card)
}

func (s *CardService) GetCards(ctx context.Context,  userID uint) ([]models.LinkedCard, error) {
	return s.Repo.GetByUser(ctx, userID)
}

func (s *CardService) DeleteCard(ctx context.Context, cardID, userID uint) error {
	return s.Repo.Delete(ctx, cardID, userID)
}

func (s *CardService) UpdateCardLimit(ctx context.Context, cardID uint, userID uint , limit float64) error {
	return s.Repo.UpdateLimit(ctx , cardID , userID , limit)
}