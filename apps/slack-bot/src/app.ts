import { App } from '@slack/bolt';
import dotenv from 'dotenv';
import { registerCommands } from './commands';
import { initializeSessionManager } from './services/user-session';

// Load environment variables
dotenv.config();

// Validate required environment variables
const requiredEnvVars = {
  SLACK_BOT_TOKEN: process.env['SLACK_BOT_TOKEN'],
  SLACK_SIGNING_SECRET: process.env['SLACK_SIGNING_SECRET'],
  // SLACK_APP_TOKEN: process.env['SLACK_APP_TOKEN'], // Not needed for HTTP endpoint mode
};

const missingVars = Object.entries(requiredEnvVars)
  .filter(([, value]) => !value)
  .map(([key]) => key);

if (missingVars.length > 0) {
  console.error('❌ Missing required environment variables:');
  missingVars.forEach(varName => {
    console.error(`   - ${varName}`);
  });
  console.error('\n📝 Please check your .env file and add the missing variables.');
  console.error('   Get your tokens from: https://api.slack.com/apps');
  process.exit(1);
}

// Initialize the app
const app = new App({
  token: process.env['SLACK_BOT_TOKEN']!,
  signingSecret: process.env['SLACK_SIGNING_SECRET']!,
  // socketMode: true, // Disabled for ngrok HTTP endpoint usage
  // appToken: process.env['SLACK_APP_TOKEN']!, // Not needed when socketMode is disabled
  port: parseInt(process.env['PORT'] || '3000'),
});

// Initialize session manager with configuration
initializeSessionManager({
  sessionTimeout: 60, // 1 hour
  maxSessions: 1000,
});

// Register all commands and actions
registerCommands(app);

// Start the app
const startApp = async (): Promise<void> => {
  try {
    await app.start();
    console.log('⚡️ Tranza Slack bot is running on the port' , parseInt(process.env['PORT'] || '3000'));
    console.log('⚡️ Tranza Slack bot is running!');
    console.log(`🔗 API Base URL: ${process.env['TRANZA_API_BASE_URL'] || 'http://localhost:8080'}`);
    console.log('📱 Available commands: /auth, /fetch-balance, /send-money, /logout, /help');
  } catch (error) {
    console.error('❌ Error starting app:', error);
    process.exit(1);
  }
};

// Handle graceful shutdown
const gracefulShutdown = (): void => {
  console.log('🔄 Shutting down Tranza Slack bot...');
  try {
    // Clean up sessions
    const { destroySessionManager } = require('./services/user-session');
    destroySessionManager();
    console.log('✅ Graceful shutdown completed');
    process.exit(0);
  } catch (error) {
    console.error('❌ Error during shutdown:', error);
    process.exit(1);
  }
};

// Listen for shutdown signals
process.on('SIGINT', gracefulShutdown);
process.on('SIGTERM', gracefulShutdown);

// Start the application
startApp();
