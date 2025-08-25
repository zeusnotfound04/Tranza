/**
 * Drop API Keys Table Script
 * 
 * This script drops the existing api_keys table to resolve the UUID migration issue.
 * After running this, restart your Go server and GORM will recreate the table
 * with the correct UUID schema.
 */

const { Client } = require('pg');
require('dotenv').config();

// Database configuration
const dbConfig = {
  host: process.env.DB_HOST || 'localhost',
  port: process.env.DB_PORT || 5432,
  database: process.env.DB_NAME || 'postgres',
  user: process.env.DB_USER || 'postgres',
  password: process.env.DB_PASSWORD || 'tranza@47323',
};

async function dropAPIKeysTable() {
  const client = new Client(dbConfig);
  
  try {
    console.log('🔌 Connecting to database...');
    await client.connect();
    console.log('✅ Connected to database successfully');

    // Check if table exists
    const checkTableQuery = `
      SELECT EXISTS (
        SELECT FROM information_schema.tables 
        WHERE table_schema = 'public' 
        AND table_name = 'api_keys'
      );
    `;
    
    const tableExists = await client.query(checkTableQuery);
    const exists = tableExists.rows[0].exists;
    
    if (!exists) {
      console.log('ℹ️  Table "api_keys" does not exist. Nothing to drop.');
      return;
    }

    console.log('📋 Table "api_keys" found. Checking current schema...');
    
    // Show current table schema
    const schemaQuery = `
      SELECT column_name, data_type, is_nullable
      FROM information_schema.columns 
      WHERE table_name = 'api_keys' 
      AND table_schema = 'public'
      ORDER BY ordinal_position;
    `;
    
    const schemaResult = await client.query(schemaQuery);
    console.log('\n📊 Current table schema:');
    schemaResult.rows.forEach(col => {
      console.log(`  - ${col.column_name}: ${col.data_type} (nullable: ${col.is_nullable})`);
    });

    // Check for any existing data
    const countQuery = 'SELECT COUNT(*) FROM api_keys';
    const countResult = await client.query(countQuery);
    const rowCount = parseInt(countResult.rows[0].count);
    
    console.log(`\n📊 Current row count: ${rowCount}`);

    if (rowCount > 0) {
      console.log('⚠️  WARNING: This table contains data that will be permanently lost!');
      console.log('   Make sure you want to proceed with dropping the table.');
    }

    // Drop the table with CASCADE to handle any dependencies
    console.log('\n🗑️  Dropping table "api_keys"...');
    await client.query('DROP TABLE IF EXISTS api_keys CASCADE');
    
    console.log('✅ Table "api_keys" dropped successfully!');

    // Verify the table is gone
    const verifyResult = await client.query(checkTableQuery);
    const stillExists = verifyResult.rows[0].exists;
    
    if (!stillExists) {
      console.log('✨ Verification successful: Table has been completely removed.');
    } else {
      console.log('⚠️  Warning: Table still exists after drop command.');
    }

  } catch (error) {
    console.error('❌ Error during table drop:', error.message);
    
    if (error.code === 'ECONNREFUSED') {
      console.error('\n💡 Database connection failed. Please check:');
      console.error('   - PostgreSQL is running');
      console.error('   - Connection details are correct');
      console.error('   - Database exists');
    } else if (error.code === '42P01') {
      console.error('\n💡 Table does not exist - this is actually good for our use case!');
    }
    
    throw error;
  } finally {
    await client.end();
    console.log('🔌 Database connection closed');
  }
}

// Run the drop table script
if (require.main === module) {
  console.log('🗑️  Starting API Keys Table Drop Script...\n');
  
  dropAPIKeysTable()
    .then(() => {
      console.log('\n🎉 Table drop completed successfully!');
      console.log('\n📝 Next steps:');
      console.log('   1. Start your Go backend server: go run .\\cmd\\server\\main.go');
      console.log('   2. GORM will automatically create the table with UUID schema');
      console.log('   3. Test API key creation through your frontend');
      console.log('   4. The UUID/uint mismatch error should be resolved');
      process.exit(0);
    })
    .catch((error) => {
      console.error('\n💥 Table drop failed:', error.message);
      console.log('\n🛠️  Manual alternative:');
      console.log('   You can also drop the table manually using pgAdmin or any PostgreSQL client');
      console.log('   Run this SQL command: DROP TABLE IF EXISTS api_keys CASCADE;');
      process.exit(1);
    });
}

module.exports = { dropAPIKeysTable };
