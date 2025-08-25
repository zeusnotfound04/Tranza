package main

import (
	"log"

	"github.com/zeusnotfound04/Tranza/config"
)

func main() {
	// Load environment variables
	config.LoadEnv()

	// Connect to database
	db := config.ConnectDB()
	defer config.CloseDB(db)

	log.Println("🚀 Adding performance indexes...")

	// Add indexes for external_transfers table
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_external_transfers_user_id ON external_transfers(user_id);",
		"CREATE INDEX IF NOT EXISTS idx_external_transfers_status ON external_transfers(status);",
		"CREATE INDEX IF NOT EXISTS idx_external_transfers_created_at ON external_transfers(created_at);",
		"CREATE INDEX IF NOT EXISTS idx_external_transfers_user_status ON external_transfers(user_id, status);",
		"CREATE INDEX IF NOT EXISTS idx_external_transfers_reference_id ON external_transfers(reference_id);",

		// Add indexes for api_keys table
		"CREATE INDEX IF NOT EXISTS idx_api_keys_user_id ON api_keys(user_id);",
		"CREATE INDEX IF NOT EXISTS idx_api_keys_key_type ON api_keys(key_type);",
		"CREATE INDEX IF NOT EXISTS idx_api_keys_is_active ON api_keys(is_active);",
		"CREATE INDEX IF NOT EXISTS idx_api_keys_workspace ON api_keys(bot_workspace) WHERE bot_workspace IS NOT NULL;",
		"CREATE INDEX IF NOT EXISTS idx_api_keys_expires_at ON api_keys(expires_at) WHERE expires_at IS NOT NULL;",
		"CREATE INDEX IF NOT EXISTS idx_api_keys_user_active ON api_keys(user_id, is_active);",

		// Add indexes for transactions table
		"CREATE INDEX IF NOT EXISTS idx_transactions_user_id ON transactions(user_id);",
		"CREATE INDEX IF NOT EXISTS idx_transactions_status ON transactions(status);",
		"CREATE INDEX IF NOT EXISTS idx_transactions_created_at ON transactions(created_at);",
		"CREATE INDEX IF NOT EXISTS idx_transactions_type ON transactions(type);",

		// Add indexes for wallets table
		"CREATE INDEX IF NOT EXISTS idx_wallets_user_id ON wallets(user_id);",
		"CREATE INDEX IF NOT EXISTS idx_wallets_updated_at ON wallets(updated_at);",

		// Add indexes for users table
		"CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);",
		"CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at);",
	}

	// Execute each index creation
	for _, indexSQL := range indexes {
		if err := db.Exec(indexSQL).Error; err != nil {
			log.Printf("⚠️  Warning: Failed to create index: %v", err)
		}
	}

	log.Println("✅ Performance indexes added successfully!")
	log.Println("🎉 Database optimization completed!")
}
