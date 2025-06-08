package models

import "time"

type TransactionStatus string


const (
	StatusFailed  TransactionStatus = "FAILED"
	StatusSuccess TransactionStatus  = "SUCCESS"
	StatusPending TransactionStatus ="PENDING"
)


type Transaction struct{
	ID        uint   `gorm:"primaryKey"`
	UserID   uint    `gorm:"not null"`
	CardID    uint    `gorm:"not null"`
	Amount    float64  `gorm:"not null"`
	Description string  
	Status       TransactionStatus `gorm:"type:varchar(10);default:'PENDING'"` 
	CreatedAt    time.Time
	UpdatedAt     time.Time 
}