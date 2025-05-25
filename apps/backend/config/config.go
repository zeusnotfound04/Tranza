package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func LoadEnv(){
	err := godotenv.Load()

	if err != nil {
		log.Println("No .env file found")
	}
}


func ConnectDB() *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable" ,
			os.Getenv("DB_HOST"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_NAME"),
			os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn) , &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect the DB :" , err)
	}
	DB = db
	return DB
}



func CloseDB (db *gorm.DB){
	sqlDB , err := db.DB()

	if err != nil {
		log.Fatal("DB close error" , err)
	}

	sqlDB.Close()
}