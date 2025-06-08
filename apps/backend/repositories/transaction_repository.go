package repositories

import (
	"context"
	"fmt"

	"github.com/zeusnotfound04/Tranza/models"
	"gorm.io/gorm"
)



type TransactionRepository interface {
	Create (ctx context.Context, txn *models.Transaction ) error
	GetUserTransactions (ctx context.Context, userId  uint)  ([]models.Transaction , error)
	GetById (ctx context.Context, txnID  uint) (*models.Transaction , error)
	WithTxn (txn *gorm.DB) TransactionRepository
}


type transactionRepository struct {
	db  *gorm.DB
}


func NewTransactionRepository(db *gorm.DB)  TransactionRepository {
	return &transactionRepository{db: db}
}


func (r *transactionRepository ) WithTxn(tx *gorm.DB) TransactionRepository {
	return &transactionRepository{db: tx}
}

func (r *transactionRepository) Create(ctx context.Context, txn *models.Transaction) error{
	if err := r.db.WithContext(ctx).Create(txn).Error; err != nil {
		return fmt.Errorf("Failed to create a transaction: %w", err)
	}
	return nil
}
func (r *transactionRepository) GetUserTransactions(ctx context.Context , userID uint) ([]models.Transaction , error) {
	var txns []models.Transaction
	if err := r.db.WithContext(ctx).
	Where("user_id", userID).
	Order("created_at DESC").
	Find(&txns).Error; err != nil {
		return nil, fmt.Errorf("could not fetch the transaction: %w" , err)
	}
	return txns , nil
}

func (r *transactionRepository) GetById(ctx context.Context , txnID uint) (*models.Transaction , error) {
	var txn *models.Transaction
	if err := r.db.WithContext(ctx).
	First(&txn , txnID).Error; err !=nil {
		return nil , fmt.Errorf("failed to fetch transaction by ID: %w", err)
	}
	return txn , nil
}



