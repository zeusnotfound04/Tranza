/**
 * Cleanup Script: Remove Existing API Keys
 * 
 * This script removes all existing API keys from the database to resolve
 * the UUID/uint mismatch issue after updating the APIKey model to use UUID.
 * 
 * Run this script before starting the server with the updated UUID-based model.
 */

const { Client } = require('pg');
require('dotenv').config();

// Database configuration
const dbConfig = {
  host: process.env.DB_HOST || 'localhost',
  port: process.env.DB_PORT || 5432,
  database: process.env.DB_NAME || 'tranza',
  user: process.env.DB_USER || 'postgres',
  password: process.env.DB_PASSWORD || 'password',
};

async function cleanupAPIKeys() {
  const client = new Client(dbConfig);
  
  try {
    console.log('🔌 Connecting to database...');
    await client.connect();
    console.log('✅ Connected to database successfully');

    // First, check how many API keys exist
    const countResult = await client.query('SELECT COUNT(*) FROM api_keys');
    const existingCount = parseInt(countResult.rows[0].count);
    
    console.log(`📊 Found ${existingCount} existing API keys`);

    if (existingCount === 0) {
      console.log('✨ No API keys to clean up. Database is already clean!');
      return;
    }

    // Show existing API keys (optional - for debugging)
    const existingKeys = await client.query(`
      SELECT id, user_id, label, key_type, created_at 
      FROM api_keys 
      ORDER BY created_at DESC 
      LIMIT 10
    `);
    
    console.log('\n📋 Sample existing API keys:');
    existingKeys.rows.forEach((key, index) => {
      console.log(`  ${index + 1}. ID: ${key.id}, UserID: ${key.user_id}, Label: ${key.label}, Type: ${key.key_type}`);
    });

    // Ask for confirmation (in a real scenario, you might want to add readline for interactive confirmation)
    console.log('\n⚠️  WARNING: This will DELETE ALL existing API keys!');
    console.log('   This is necessary because the user_id column type has changed from integer to UUID.');
    console.log('   Users will need to create new API keys after this cleanup.\n');

    // Delete all API keys
    console.log('🗑️  Deleting all existing API keys...');
    const deleteResult = await client.query('DELETE FROM api_keys');
    
    console.log(`✅ Successfully deleted ${deleteResult.rowCount} API keys`);

    // Verify cleanup
    const verifyResult = await client.query('SELECT COUNT(*) FROM api_keys');
    const remainingCount = parseInt(verifyResult.rows[0].count);
    
    if (remainingCount === 0) {
      console.log('✨ Cleanup completed successfully! Database is now ready for UUID-based API keys.');
    } else {
      console.log(`⚠️  Warning: ${remainingCount} API keys still remain in database.`);
    }

    // Optional: Reset the auto-increment sequence for the ID column
    try {
      await client.query('ALTER SEQUENCE api_keys_id_seq RESTART WITH 1');
      console.log('🔄 Reset API key ID sequence to start from 1');
    } catch (sequenceError) {
      console.log('ℹ️  Note: Could not reset ID sequence (this is okay)');
    }

  } catch (error) {
    console.error('❌ Error during cleanup:', error.message);
    
    if (error.code === 'ECONNREFUSED') {
      console.error('\n💡 Database connection failed. Please check:');
      console.error('   - Database is running');
      console.error('   - Connection details in .env file are correct');
      console.error('   - Database name exists');
    } else if (error.code === '42P01') {
      console.error('\n💡 Table "api_keys" does not exist. This might mean:');
      console.error('   - Database migration hasn\'t been run yet');
      console.error('   - Database schema is not set up');
    }
    
    throw error;
  } finally {
    await client.end();
    console.log('🔌 Database connection closed');
  }
}

// Run the cleanup
if (require.main === module) {
  console.log('🧹 Starting API Keys Cleanup Script...\n');
  
  cleanupAPIKeys()
    .then(() => {
      console.log('\n🎉 Cleanup script completed successfully!');
      console.log('\n📝 Next steps:');
      console.log('   1. Start your Go backend server');
      console.log('   2. Test API key creation with the new UUID system');
      console.log('   3. Users will need to create new API keys');
      process.exit(0);
    })
    .catch((error) => {
      console.error('\n💥 Cleanup script failed:', error.message);
      console.log('\n🛠️  Troubleshooting:');
      console.log('   1. Ensure database is running');
      console.log('   2. Check database connection details');
      console.log('   3. Verify .env file is properly configured');
      console.log('   4. Make sure api_keys table exists');
      process.exit(1);
    });
}

module.exports = { cleanupAPIKeys };
