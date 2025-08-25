# API Keys Cleanup Script

This script removes all existing API keys from the database to resolve the UUID/uint mismatch issue after updating the APIKey model to use UUID instead of integer user IDs.

## ⚠️ Important Warning

**This script will DELETE ALL existing API keys from your database.** This is necessary because:

1. The old API keys were created with integer `user_id` values
2. The new system expects UUID `user_id` values  
3. There's no clean way to migrate the data without knowing the mapping

Users will need to create new API keys after running this cleanup.

## 🚀 How to Run

### Step 1: Install Dependencies
```bash
cd f:\VISHESH\VS code\JavaScript\MERN\Project\tranza\apps\backend\scripts
npm install
```

### Step 2: Configure Database Connection
```bash
# Copy the example environment file
cp .env.example .env

# Edit .env with your actual database credentials
# Use the same values from your main backend .env file
```

### Step 3: Run the Cleanup Script
```bash
npm run cleanup:api-keys
```

Or directly:
```bash
node cleanup_api_keys.js
```

## 📋 What the Script Does

1. **Connects** to your PostgreSQL database
2. **Counts** existing API keys
3. **Shows** a sample of existing keys for reference
4. **Deletes** all API keys from the `api_keys` table
5. **Verifies** the cleanup was successful
6. **Resets** the ID sequence (optional)

## 🔧 Database Configuration

The script expects these environment variables:

```env
DB_HOST=localhost
DB_PORT=5432
DB_NAME=tranza
DB_USER=postgres
DB_PASSWORD=your_password
```

## 🐛 Troubleshooting

### "Connection Refused" Error
- Ensure PostgreSQL is running
- Check database host and port
- Verify credentials are correct

### "Table does not exist" Error  
- Run database migrations first
- Ensure the `api_keys` table exists

### Permission Errors
- Ensure the database user has DELETE permissions
- Check if the user can connect to the specified database

## 📝 After Running the Script

1. **Restart** your Go backend server
2. **Test** API key creation through your frontend
3. **Inform users** they need to create new API keys
4. **Verify** the UUID system is working correctly

## 🔄 Next Steps

After cleanup:
- The `api_keys` table will be empty
- New API keys will use UUID `user_id` format
- The UUID/uint type mismatch will be resolved
- Users can create new API keys through the web interface

## 🚨 Rollback

There's no automatic rollback for this script. If you need to restore API keys:
- Use database backups
- Or recreate API keys manually through the application
