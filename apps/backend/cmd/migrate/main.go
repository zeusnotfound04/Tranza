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

	log.Println("🚀 Running database migrations...")

	// Auto-migrate all models
	err := db.AutoMigrate(
		&models.User{},
		&models.Wallet{},
		&models.Transaction{},
		&models.APIKey{},
		&models.ExternalTransfer{},
		&models.LinkedCard{},
		&models.Address{},
		&models.EmailVerification{},
	)

	if err != nil {
		log.Fatalf("❌ Migration failed: %v", err)
	}

	log.Println("✅ Database migrations completed successfully!")

	// Fix foreign key constraint for external_transfers
	log.Println("🔧 Fixing foreign key constraints...")

	// Drop the problematic foreign key constraint if it exists
	db.Exec("ALTER TABLE external_transfers DROP CONSTRAINT IF EXISTS fk_external_transfers_transaction;")

	// Recreate it as nullable
	db.Exec("ALTER TABLE external_transfers ALTER COLUMN transaction_id DROP NOT NULL;")

	// Add back the foreign key constraint that allows nulls
	db.Exec("ALTER TABLE external_transfers ADD CONSTRAINT fk_external_transfers_transaction FOREIGN KEY (transaction_id) REFERENCES transactions(id) ON DELETE SET NULL;")

	log.Println("✅ Foreign key constraints fixed!")
	log.Println("🎉 All migrations completed successfully!")
}
