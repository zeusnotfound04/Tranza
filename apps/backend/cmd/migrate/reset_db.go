package main

import (
	"log"

	"github.com/zeusnotfound04/Tranza/config"
	"github.com/zeusnotfound04/Tranza/models"
)

func main() {
	// Load environment variables
	config.LoadEnv()

	// Connect to database
	db := config.ConnectDB()
	defer config.CloseDB(db)

	// Drop existing tables to fix schema conflicts
	log.Println("🗑️  Dropping existing tables...")

	// Drop tables in reverse dependency order
	db.Exec("DROP TABLE IF EXISTS ai_spending_trackers CASCADE")
	db.Exec("DROP TABLE IF EXISTS ai_spending_limits CASCADE")
	db.Exec("DROP TABLE IF EXISTS ai_payment_requests CASCADE")
	db.Exec("DROP TABLE IF EXISTS transactions CASCADE")
	db.Exec("DROP TABLE IF EXISTS linked_cards CASCADE")
	db.Exec("DROP TABLE IF EXISTS wallets CASCADE")
	db.Exec("DROP TABLE IF EXISTS api_keys CASCADE")
	db.Exec("DROP TABLE IF EXISTS email_verifications CASCADE")
	db.Exec("DROP TABLE IF EXISTS users CASCADE")

	log.Println("✅ Tables dropped successfully")

	// Auto-migrate all models with correct schema
	log.Println("🔄 Creating tables with correct schema...")

	err := db.AutoMigrate(
		&models.User{},
		&models.EmailVerification{},
		&models.Wallet{},
		&models.Transaction{},
		&models.APIKey{},
		&models.LinkedCard{},
		&models.AIPaymentRequest{},
		&models.AISpendingLimit{},
		&models.AISpendingTracker{},
	)

	if err != nil {
		log.Fatal("❌ Failed to migrate database:", err)
	}

	log.Println("✅ Database migration completed successfully!")
	log.Println("🎉 All tables created with correct UUID schema")
}
