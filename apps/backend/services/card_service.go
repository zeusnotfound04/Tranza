package services

import (
	"github.com/zeusnotfound04/Tranza/models"
	"github.com/zeusnotfound04/Tranza/repositories"
)

type CardService struct {
	Repo *repositories.CardRepository
}


func NewCardService(repo *repositories.CardRepository) *CardService {
	return &CardService{Repo: repo}
}

func (s *CardService) LinkCard(card *models.LinkedCard) error {
	// Add validations if needed (e.g., max limit check)
	return s.Repo.Create(card)
}

func (s *CardService) GetCards(userID uint) ([]models.LinkedCard, error) {
	return s.Repo.GetByUser(userID)
}

func (s *CardService) DeleteCard(cardID, userID uint) error {
	return s.Repo.Delete(cardID, userID)
}

func (s *CardService) UpdateCardLimit(cardID, userID uint, limit float64) error {
	return s.Repo.UpdateLimit(cardID, userID, limit)
}