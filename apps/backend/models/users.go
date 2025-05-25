package models

type User struct {
	ID        string `json:"id" gorm:"primaryKey"`
	Email     string `json:"email" gorm:"unique;not null"`
	Username  string `json:"username" gorm:"unique;not null"`
	Password  string `json:"password,omitempty"`
	CreatedAt int64  `json:"created_at"`
}
