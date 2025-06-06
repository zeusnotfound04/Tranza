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

func (s *CardService) LinkCard(ctx context.Context , card *models.LinkedCard) erro,r {
	// Add validations if needed (e.g., max limit check)
	return s.Repo.Create(card)
}

func (s *CardService) GetCards(ctx context.Context,  userID uint) ([]models.LinkedCard, error) {
	return s.Repo.GetByUser(userID)
}

func (s *CardService) DeleteCard(ctx context.Context, cardID, userID uint) error {
	return s.Repo.Delete(cardID, userID)
}

func (s *CardService) UpdateCardLimit(ctx context.Context, card *models.LinkedCard) error {
	return s.Repo.UpdateLimit(ctx )
}