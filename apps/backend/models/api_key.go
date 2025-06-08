package models

import "time"


type APIKey struct {
	ID         uint      `gorm:"primaryKey"`
	UserID     uint      `gorm:"not null"`
	KeyHash    string    `gorm:"not null;uniqueIndex"`
	Label      string     `gorm:"type:varchar(100)"`
	CreatedAt  time.Time
	ExpiresAt  *time.Time
	LastUsedAt time.Time
	IsActive   bool      `gorm:"default:true"`
}


func (k *APIKey) IsExpired() bool {
	return k.ExpiresAt != nil && k.ExpiresAt.Before(time.Now())
}
