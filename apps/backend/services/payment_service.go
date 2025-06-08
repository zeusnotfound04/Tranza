package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/zeusnotfound04/Tranza/models"
	"github.com/zeusnotfound04/Tranza/repositories"
)



var ErrPaymentLimitExeceed = errors.New("payment amount exceeds card limit")
var ErrUnauthorizedCardAccess = errors.New("authorised access to card")

type PaymentService struct {
	CardRepo *repositories.CardRepository
	TxnRepo  repositories.TransactionRepository
}


func NewPaymentService(cardRepo *repositories.CardRepository , txnRepo repositories.TransactionRepository) *PaymentService {
	return &PaymentService{
		CardRepo: cardRepo,
		TxnRepo: txnRepo,
	}
}
func (s *PaymentService) MakePayment(ctx context.Context , userID uint , CardID uint, amount float64, desc string ) (*models.Transaction , error){
	card , err := s.CardRepo.FindByID(ctx , userID )
	if err != nil {
		return nil , fmt.Errorf("Card not found: %w" , err)
	}

	if card.UserId != userID {
		return nil , ErrUnauthorizedCardAccess
	}

	if card.Limit < amount {
		return nil , ErrPaymentLimitExeceed
	}

	card.Limit -= amount
	if err := s.CardRepo.UpdateLimit(ctx , card.ID , card.UserId , card.Limit); err != nil {
		return nil , fmt.Errorf("Failed to update the limit: %w", err)
	}	

	txn := &models.Transaction{
        UserID:      userID,
		CardID:      CardID, 
		Amount:      amount,
		Description: desc,
		Status:      models.StatusSuccess,
	}

	if err := s.TxnRepo.Create(ctx , txn); err != nil {
		return nil, fmt.Errorf("failed to record transaction: %w", err)
	}
    
	return txn, nil
}

