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

type RazorpayConfig struct {
	RazorpayKeyID     string
	RazorpayKeySecret string
	WebhookSecret     string
	Environment       string // "test" or "live"
}

func LoadEnv(){
	err := godotenv.Load()

	if err != nil {
		log.Println("No .env file found")
	}
}


func ConnectDB() *gorm.DB {

    dsn := fmt.Sprintf(
        "host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
        os.Getenv("DB_HOST"),
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_NAME"),
        os.Getenv("DB_PORT"),
        os.Getenv("DB_SSL"),
    )


	fmt.Printf(dsn)
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
func LoadConfig() *RazorpayConfig {
	keyID := os.Getenv("RAZORPAY_KEY_ID")
	keySecret := os.Getenv("RAZORPAY_KEY_SECRET")
	webhookSecret := os.Getenv("RAZORPAY_WEBHOOK_SECRET")
	
	if keyID == "" || keySecret == "" {
		log.Fatal("RAZORPAY_KEY_ID and RAZORPAY_KEY_SECRET must be set")
	}
	
	env := os.Getenv("RAZORPAY_ENV")
	if env == "" {
		env = "test"
	}
	
	return &RazorpayConfig{
		RazorpayKeyID:     keyID,
		RazorpayKeySecret: keySecret,
		WebhookSecret:     webhookSecret,
		Environment:       env,
	}
}

func (c *RazorpayConfig) GetBaseURL() string {
	return "https://api.razorpay.com/v1"
}

func (c *RazorpayConfig) IsProduction() bool {
	return c.Environment == "live"
}
func (c *RazorpayConfig) Validate() error {
	if c.RazorpayKeyID == "" {
		return fmt.Errorf("RAZORPAY_KEY_ID is required")
	}
	if c.RazorpayKeySecret == "" {
		return fmt.Errorf("RAZORPAY_KEY_SECRET is required")
	}
	return nil
}