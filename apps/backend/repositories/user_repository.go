package repositories

import (
	"github.com/zeusnotfound04/Tranza/models"
	"gorm.io/gorm"
)

var db *gorm.DB


func InitRepo(database *gorm.DB){
	db = database
}

func GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil , err
	}

	return &user , nil
}

func CreateUser(user *models.User) error {
	return db.Create(user).Error
}

