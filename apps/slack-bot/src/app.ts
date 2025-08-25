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

const app = new App({
  token: process.env['SLACK_BOT_TOKEN']!,
  signingSecret: process.env['SLACK_SIGNING_SECRET']!,
  port: parseInt(process.env['PORT'] || '3000'),
});

// Initialize session manager with configuration
initializeSessionManager({
  sessionTimeout: 60, // 1 hour
  maxSessions: 1000,
});

// Register all commands and actions
registerCommands(app);

app.use(async ({ body, next }) => {
  console.log('📨 Incoming Slack request:', {
    type: (body as any).type || 'unknown',
    user: (body as any).user?.id || 'unknown',
    timestamp: new Date().toISOString()
  });
  await next();
});

// Add error handling for unhandled requests
app.error(async (error) => {
  console.error('❌ Slack app error:', error);
});

// Start the app
const startApp = async (): Promise<void> => {
  try {
    await app.start();
    const port = parseInt(process.env['PORT'] || '3000');
    console.log('⚡️ Tranza Slack bot is running!');
    console.log(`🌐 Server running on port: ${port}`);
    console.log(`🔗 API Base URL: ${process.env['TRANZA_API_BASE_URL'] || 'http://localhost:8080'}`);
    console.log('📱 Available commands: /auth, /fetch-balance, /send-money, /logout, /help');
    console.log('\n🔧 Slack App Configuration URLs:');
    console.log(`   Event Subscriptions → Request URL: https://your-domain.com/slack/events`);
    console.log(`   Interactivity & Shortcuts → Request URL: https://your-domain.com/slack/actions`);
    console.log(`   Slash Commands → Request URL: https://your-domain.com/slack/commands`);
    console.log('\n📝 Replace "your-domain.com" with your actual ngrok or server domain');
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
