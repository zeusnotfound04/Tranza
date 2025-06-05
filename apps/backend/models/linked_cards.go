package models

import "time"

type LinkedCard struct {
	ID         uint    `gorm:"primaryKey" json:"id"`
	UserId     uint    `json:"user_id"`
	CardNumber string  `json:"card_number"`
	CardType   string  `json:"card_type"`
	Limit      float64 `json:"card_type"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}