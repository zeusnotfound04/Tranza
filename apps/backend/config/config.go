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

type OAuthConfig struct {
	Google GoogleConfig `json:"google"`
	GitHub GitHubConfig `json:"github"`
}

type GoogleConfig struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectURL  string `json:"redirect_url"`
}

type GitHubConfig struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectURL  string `json:"redirect_url"`
}

type Config struct {
	OAuth OAuthConfig `json:"oauth"`
}

func LoadEnv() {
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

	fmt.Printf("DSN: %s\n", dsn)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect the DB :", err)
	}
	DB = db
	return DB
}

func CloseDB(db *gorm.DB) {
	sqlDB, err := db.DB()

	if err != nil {
		log.Fatal("DB close error", err)
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

// LoadOAuthConfig loads OAuth configuration from environment variables
func LoadOAuthConfig() *Config {
	return &Config{
		OAuth: OAuthConfig{
			Google: GoogleConfig{
				ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
				ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
				RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
			},
			GitHub: GitHubConfig{
				ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
				ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
				RedirectURL:  os.Getenv("GITHUB_REDIRECT_URL"),
			},
		},
	}
}
