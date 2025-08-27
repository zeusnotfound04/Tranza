#!/usr/bin/env node

/**
 * Cleanup Script for Test-to-Live Environment Transition
 * 
 * This script prepares the database for live Razorpay environment by:
 * 1. Backing up current data
 * 2. Adding environment columns
 * 3. Marking test data
 * 4. Resetting wallet balances
 * 5. Creating audit trail
 * 
 * Usage: node cleanup_for_live_environment.js [--dry-run] [--backup-only]
 */

const { Pool } = require('pg');
const fs = require('fs');
const path = require('path');
const { execSync } = require('child_process');

// Configuration
const config = {
  database: {
    host: process.env.DB_HOST || 'localhost',
    port: process.env.DB_PORT || 5432,
    database: process.env.DB_NAME || 'tranza',
    user: process.env.DB_USER || '',
    password: process.env.DB_PASSWORD || 'password',
  },
  backupDir: './backups',
  isDryRun: process.argv.includes('--dry-run'),
  backupOnly: process.argv.includes('--backup-only'),
};

// Colors for console output
const colors = {
  reset: '\x1b[0m',
  bright: '\x1b[1m',
  red: '\x1b[31m',
  green: '\x1b[32m',
  yellow: '\x1b[33m',
  blue: '\x1b[34m',
  magenta: '\x1b[35m',
  cyan: '\x1b[36m',
};

function log(message, color = colors.reset) {
  console.log(`${color}${message}${colors.reset}`);
}

function logStep(step, message) {
  log(`${colors.cyan}[STEP ${step}]${colors.reset} ${message}`);
}

function logSuccess(message) {
  log(`${colors.green}✅ ${message}${colors.reset}`);
}

function logWarning(message) {
  log(`${colors.yellow}⚠️  ${message}${colors.reset}`);
}

function logError(message) {
  log(`${colors.red}❌ ${message}${colors.reset}`);
}

class DatabaseCleanup {
  constructor() {
    this.pool = new Pool(config.database);
    this.backupPath = null;
  }

  async connect() {
    try {
      const client = await this.pool.connect();
      log(`${colors.green}✅ Connected to database: ${config.database.database}${colors.reset}`);
      client.release();
      return true;
    } catch (error) {
      logError(`Failed to connect to database: ${error.message}`);
      return false;
    }
  }

  async createBackup() {
    logStep(1, 'Creating database backup...');
    
    if (!fs.existsSync(config.backupDir)) {
      fs.mkdirSync(config.backupDir, { recursive: true });
    }

    const timestamp = new Date().toISOString().replace(/[:.]/g, '-');
    this.backupPath = path.join(config.backupDir, `backup_before_live_${timestamp}.sql`);

    try {
      const pgDumpCommand = `pg_dump -h ${config.database.host} -p ${config.database.port} -U ${config.database.user} -d ${config.database.database} -f "${this.backupPath}"`;
      
      log(`${colors.blue}Running: ${pgDumpCommand}${colors.reset}`);
      
      // Set PGPASSWORD environment variable for pg_dump
      const env = { ...process.env, PGPASSWORD: config.database.password };
      execSync(pgDumpCommand, { env, stdio: 'inherit' });
      
      logSuccess(`Database backup created: ${this.backupPath}`);
      
      // Get backup file size
      const stats = fs.statSync(this.backupPath);
      const fileSizeInMB = (stats.size / (1024 * 1024)).toFixed(2);
      log(`Backup size: ${fileSizeInMB} MB`);
      
      return true;
    } catch (error) {
      logError(`Failed to create backup: ${error.message}`);
      return false;
    }
  }

  async getCurrentStats() {
    logStep(2, 'Gathering current database statistics...');
    
    const queries = [
      { name: 'Total Users', query: 'SELECT COUNT(*) as count FROM users' },
      { name: 'Total Wallets', query: 'SELECT COUNT(*) as count FROM wallets' },
      { name: 'Total Wallet Balance', query: 'SELECT COALESCE(SUM(balance), 0) as total FROM wallets' },
      { name: 'Total Transactions', query: 'SELECT COUNT(*) as count FROM transactions' },
      { name: 'Total Payments', query: 'SELECT COUNT(*) as count FROM payments' },
      { name: 'Total External Transfers', query: 'SELECT COUNT(*) as count FROM external_transfers' },
      { name: 'Total API Keys', query: 'SELECT COUNT(*) as count FROM api_keys' },
    ];

    const stats = {};
    
    for (const { name, query } of queries) {
      try {
        const result = await this.pool.query(query);
        stats[name] = result.rows[0].count || result.rows[0].total || 0;
        log(`  ${name}: ${colors.bright}${stats[name]}${colors.reset}`);
      } catch (error) {
        logWarning(`Failed to get ${name}: ${error.message}`);
        stats[name] = 'Error';
      }
    }

    return stats;
  }

  async addEnvironmentColumns() {
    logStep(3, 'Adding environment tracking columns...');
    
    const alterQueries = [
      {
        table: 'transactions',
        query: `ALTER TABLE transactions ADD COLUMN IF NOT EXISTS environment VARCHAR(10) DEFAULT 'live'`
      },
      {
        table: 'payments', 
        query: `ALTER TABLE payments ADD COLUMN IF NOT EXISTS environment VARCHAR(10) DEFAULT 'live'`
      },
      {
        table: 'external_transfers',
        query: `ALTER TABLE external_transfers ADD COLUMN IF NOT EXISTS environment VARCHAR(10) DEFAULT 'live'`
      }
    ];

    for (const { table, query } of alterQueries) {
      try {
        if (config.isDryRun) {
          log(`${colors.yellow}[DRY RUN]${colors.reset} Would execute: ${query}`);
        } else {
          await this.pool.query(query);
          logSuccess(`Added environment column to ${table}`);
        }
      } catch (error) {
        if (error.message.includes('already exists')) {
          log(`  Environment column already exists in ${table}`);
        } else {
          logError(`Failed to add environment column to ${table}: ${error.message}`);
        }
      }
    }
  }

  async markTestData() {
    logStep(4, 'Marking existing data as test environment...');
    
    const updateQueries = [
      {
        table: 'transactions',
        query: `UPDATE transactions SET environment = 'test' WHERE environment IS NULL OR environment = 'live'`
      },
      {
        table: 'payments',
        query: `UPDATE payments SET environment = 'test' WHERE environment IS NULL OR environment = 'live'`
      },
      {
        table: 'external_transfers', 
        query: `UPDATE external_transfers SET environment = 'test' WHERE environment IS NULL OR environment = 'live'`
      }
    ];

    for (const { table, query } of updateQueries) {
      try {
        if (config.isDryRun) {
          log(`${colors.yellow}[DRY RUN]${colors.reset} Would execute: ${query}`);
        } else {
          const result = await this.pool.query(query);
          logSuccess(`Marked ${result.rowCount} records in ${table} as test data`);
        }
      } catch (error) {
        logError(`Failed to mark test data in ${table}: ${error.message}`);
      }
    }
  }

  async resetWalletBalances() {
    logStep(5, 'Resetting wallet balances to zero...');
    
    try {
      // First, get current total balance
      const balanceResult = await this.pool.query('SELECT COALESCE(SUM(balance), 0) as total FROM wallets WHERE balance > 0');
      const totalBalance = parseFloat(balanceResult.rows[0].total);
      
      if (totalBalance > 0) {
        log(`  Current total balance across all wallets: ₹${totalBalance.toFixed(2)}`);
        
        if (config.isDryRun) {
          log(`${colors.yellow}[DRY RUN]${colors.reset} Would reset ${totalBalance.toFixed(2)} rupees across all wallets`);
        } else {
          const result = await this.pool.query('UPDATE wallets SET balance = 0.00');
          logSuccess(`Reset balances for ${result.rowCount} wallets (₹${totalBalance.toFixed(2)} total)`);
        }
      } else {
        log('  All wallet balances are already zero');
      }
    } catch (error) {
      logError(`Failed to reset wallet balances: ${error.message}`);
    }
  }

  async createAuditTrail() {
    logStep(6, 'Creating audit trail transactions...');
    
    const auditQuery = `
      INSERT INTO transactions (user_id, wallet_id, amount, transaction_type, status, description, environment, created_at, updated_at)
      SELECT 
        w.user_id, 
        w.id, 
        0, 
        'system', 
        'completed', 
        'Wallet reset for live environment transition - Previous test balance cleared', 
        'live', 
        NOW(), 
        NOW()
      FROM wallets w
      WHERE NOT EXISTS (
        SELECT 1 FROM transactions t 
        WHERE t.wallet_id = w.id 
        AND t.description = 'Wallet reset for live environment transition - Previous test balance cleared'
      )
    `;

    try {
      if (config.isDryRun) {
        log(`${colors.yellow}[DRY RUN]${colors.reset} Would create audit trail transactions`);
      } else {
        const result = await this.pool.query(auditQuery);
        logSuccess(`Created ${result.rowCount} audit trail transactions`);
      }
    } catch (error) {
      logError(`Failed to create audit trail: ${error.message}`);
    }
  }

  async resetAPIKeyUsage() {
    logStep(7, 'Resetting API key usage statistics (optional)...');
    
    try {
      const resetQuery = `
        UPDATE api_keys 
        SET usage_count = 0, 
            last_used_at = created_at 
        WHERE usage_count > 0
      `;
      
      if (config.isDryRun) {
        log(`${colors.yellow}[DRY RUN]${colors.reset} Would reset API key usage statistics`);
      } else {
        const result = await this.pool.query(resetQuery);
        logSuccess(`Reset usage statistics for ${result.rowCount} API keys`);
      }
    } catch (error) {
      logError(`Failed to reset API key usage: ${error.message}`);
    }
  }

  async verifyChanges() {
    logStep(8, 'Verifying changes...');
    
    try {
      // Check wallet balances
      const balanceCheck = await this.pool.query('SELECT COUNT(*) as count FROM wallets WHERE balance > 0');
      const walletsWithBalance = parseInt(balanceCheck.rows[0].count);
      
      if (walletsWithBalance === 0) {
        logSuccess('All wallet balances are zero');
      } else {
        logWarning(`${walletsWithBalance} wallets still have non-zero balances`);
      }

      // Check test data marking
      const testDataCheck = await this.pool.query(`
        SELECT 
          (SELECT COUNT(*) FROM transactions WHERE environment = 'test') as test_transactions,
          (SELECT COUNT(*) FROM payments WHERE environment = 'test') as test_payments,
          (SELECT COUNT(*) FROM external_transfers WHERE environment = 'test') as test_transfers
      `);
      
      const testData = testDataCheck.rows[0];
      log(`  Test transactions: ${testData.test_transactions}`);
      log(`  Test payments: ${testData.test_payments}`);
      log(`  Test transfers: ${testData.test_transfers}`);

      // Check audit trail
      const auditCheck = await this.pool.query(`
        SELECT COUNT(*) as count 
        FROM transactions 
        WHERE description LIKE '%live environment transition%'
      `);
      
      const auditRecords = parseInt(auditCheck.rows[0].count);
      log(`  Audit trail records: ${auditRecords}`);
      
    } catch (error) {
      logError(`Failed to verify changes: ${error.message}`);
    }
  }

  async generateReport() {
    log(`\n${colors.bright}=== CLEANUP REPORT ===${colors.reset}`);
    
    if (config.isDryRun) {
      logWarning('This was a DRY RUN - no actual changes were made');
    }
    
    if (this.backupPath) {
      log(`📁 Backup created: ${this.backupPath}`);
    }
    
    log(`\n${colors.bright}Next Steps:${colors.reset}`);
    log(`1. Update environment variables:`);
    log(`   ${colors.cyan}RAZORPAY_KEY_ID=rzp_live_xxxxxxxxx${colors.reset}`);
    log(`   ${colors.cyan}RAZORPAY_KEY_SECRET=your_live_secret${colors.reset}`);
    log(`   ${colors.cyan}RAZORPAY_ENV=live${colors.reset}`);
    
    log(`\n2. Update Razorpay webhook URLs in dashboard`);
    log(`\n3. Test with small amounts (₹1) first`);
    log(`\n4. Monitor logs for any issues`);
    
    log(`\n${colors.green}✅ Database is ready for live environment!${colors.reset}`);
  }

  async cleanup() {
    await this.pool.end();
  }
}

// Main execution function
async function main() {
  log(`${colors.bright}🚀 Tranza Live Environment Cleanup Script${colors.reset}`);
  log(`${colors.bright}===========================================${colors.reset}\n`);
  
  if (config.isDryRun) {
    logWarning('Running in DRY RUN mode - no changes will be made');
  }
  
  if (config.backupOnly) {
    log(`${colors.cyan}Running in BACKUP ONLY mode${colors.reset}`);
  }

  const cleanup = new DatabaseCleanup();
  
  try {
    // Connect to database
    const connected = await cleanup.connect();
    if (!connected) {
      process.exit(1);
    }

    // Get current stats
    await cleanup.getCurrentStats();
    
    // Create backup
    const backupSuccess = await cleanup.createBackup();
    if (!backupSuccess) {
      logError('Backup failed - aborting cleanup');
      process.exit(1);
    }

    if (config.backupOnly) {
      log(`\n${colors.green}✅ Backup completed successfully!${colors.reset}`);
      process.exit(0);
    }

    // Confirm before proceeding
    if (!config.isDryRun) {
      const readline = require('readline').createInterface({
        input: process.stdin,
        output: process.stdout
      });

      const answer = await new Promise(resolve => {
        readline.question(
          `\n${colors.yellow}⚠️  This will reset all wallet balances to zero and mark existing data as test data.\n` +
          `Are you sure you want to proceed? (type 'yes' to continue): ${colors.reset}`,
          resolve
        );
      });

      readline.close();

      if (answer.toLowerCase() !== 'yes') {
        log('Operation cancelled by user');
        process.exit(0);
      }
    }

    // Execute cleanup steps
    await cleanup.addEnvironmentColumns();
    await cleanup.markTestData();
    await cleanup.resetWalletBalances();
    await cleanup.createAuditTrail();
    await cleanup.resetAPIKeyUsage();
    await cleanup.verifyChanges();
    await cleanup.generateReport();

  } catch (error) {
    logError(`Cleanup failed: ${error.message}`);
    console.error(error);
    process.exit(1);
  } finally {
    await cleanup.cleanup();
  }
}

// Handle script termination
process.on('SIGINT', async () => {
  log('\n\nScript interrupted by user');
  process.exit(0);
});

process.on('unhandledRejection', (reason, promise) => {
  logError('Unhandled Rejection at:', promise, 'reason:', reason);
  process.exit(1);
});

// Run the script
if (require.main === module) {
  main().catch(console.error);
}

module.exports = { DatabaseCleanup };
