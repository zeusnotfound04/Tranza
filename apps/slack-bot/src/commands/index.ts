import { App, SlackCommandMiddlewareArgs, SlackActionMiddlewareArgs } from '@slack/bolt';
import { createAPIClient, validateTransfer, createTransfer, getWalletBalance } from '../clients/tranza-api';
import { getUserSession, authenticateUser, logoutUser } from '../services/user-session';

// Command: /auth - Authenticate user with API key
export const handleAuthCommand = async ({ command, ack, respond }: SlackCommandMiddlewareArgs) => {
  await ack();

  const apiKey = command.text?.trim();
  
  if (!apiKey) {
    await respond({
      text: "❌ Please provide your API key: `/auth your-api-key-here`",
    });
    return;
  }

  try {
    // Authenticate user with backend
    const authResult = await authenticateUser(command.user_id, apiKey);
    
    if (!authResult.success) {
      await respond({
        text: `❌ ${authResult.message}`,
      });
      return;
    }
    
    await respond({
      text: "✅ Authentication successful! You can now use `/fetch-balance` and `/send-money` commands.",
    });
  } catch (error) {
    console.error('Auth command error:', error);
    await respond({
      text: "❌ Authentication failed. Please try again or contact support.",
    });
  }
};

// Command: /fetch-balance - Get wallet balance
export const handleFetchBalanceCommand = async ({ command, ack, respond }: SlackCommandMiddlewareArgs) => {
  await ack();

  try {
    // Check if user is authenticated
    const session = getUserSession(command.user_id);
    if (!session) {
      await respond({
        text: "❌ You need to authenticate first. Use `/auth your-api-key` to get started.",
      });
      return;
    }

    // Get wallet balance
    const client = createAPIClient({
      baseURL: process.env['TRANZA_API_BASE_URL'] || 'http://localhost:8080',
      apiKey: session.apiKey,
    });

    const balance = await getWalletBalance(client);
    
    await respond({
      text: `💰 *Wallet Balance*\n${balance.message}`,
    });
  } catch (error) {
    console.error('Fetch balance error:', error);
    await respond({
      text: "❌ Failed to fetch balance. Please try again or check your authentication.",
    });
  }
};
export const handlePing = async ({ command, ack, respond }: SlackCommandMiddlewareArgs) => {

    console.log("Ping command received");
  await ack();
  await respond({
    text: "🏓 Pong!",
  });
};

export const handleSendMoneyCommand = async ({ command, ack, respond }: SlackCommandMiddlewareArgs) => {
  await ack();

  try {
    // Check if user is authenticated
    const session = getUserSession(command.user_id);
    if (!session) {
      await respond({
        text: "❌ You need to authenticate first. Use `/auth your-api-key` to get started.",
      });
      return;
    }

    // Parse command parameters
    const params = command.text?.trim().split(' ') || [];
    if (params.length < 3) {
      await respond({
        text: "❌ Invalid format. Use: `/send-money <amount> <upi|phone> <recipient>`\n" +
              "Example: `/send-money 100 upi user@paytm` or `/send-money 50 phone 9876543210`",
      });
      return;
    }

    const [amount, recipientType, recipient] = params;
    
    if (!['upi', 'phone'].includes(recipientType)) {
      await respond({
        text: "❌ Recipient type must be either 'upi' or 'phone'",
      });
      return;
    }

    // Create API client
    const client = createAPIClient({
      baseURL: process.env['TRANZA_API_BASE_URL'] || 'http://localhost:8080',
      apiKey: session.apiKey,
    });

    // Validate transfer first
    const validation = await validateTransfer(client, {
      amount,
      recipient_type: recipientType as 'upi' | 'phone',
      recipient_value: recipient,
    });

    if (!validation.valid) {
      await respond({
        text: `❌ Transfer validation failed:\n${validation.errors.join('\n')}`,
      });
      return;
    }

    // Show confirmation with buttons
    await respond({
      text: `💸 *Transfer Confirmation*\n` +
            `Amount: ₹${amount}\n` +
            `Recipient: ${recipient}\n` +
            `Fee: ₹${validation.transfer_fee}\n` +
            `Total: ₹${validation.total_amount}\n` +
            `Estimated Time: ${validation.estimated_time}`,
      blocks: [
        {
          type: "section",
          text: {
            type: "mrkdwn",
            text: `💸 *Transfer Confirmation*\n` +
                  `Amount: ₹${amount}\n` +
                  `Recipient: ${recipient}\n` +
                  `Fee: ₹${validation.transfer_fee}\n` +
                  `Total: ₹${validation.total_amount}\n` +
                  `Estimated Time: ${validation.estimated_time}`
          }
        },
        {
          type: "actions",
          elements: [
            {
              type: "button",
              text: {
                type: "plain_text",
                text: "✅ Confirm Transfer"
              },
              style: "primary",
              action_id: "confirm_transfer",
              value: JSON.stringify({ amount, recipientType, recipient })
            },
            {
              type: "button",
              text: {
                type: "plain_text",
                text: "❌ Cancel"
              },
              style: "danger",
              action_id: "cancel_transfer"
            }
          ]
        }
      ]
    });

  } catch (error) {
    console.error('Send money command error:', error);
    await respond({
      text: "❌ Failed to process transfer request. Please try again or check your authentication.",
    });
  }
};


export const handleConfirmTransfer = async ({ ack, body, respond }: SlackActionMiddlewareArgs) => {
  await ack();

  try {
    const userId = body.user.id;
    const session = getUserSession(userId);
    
    if (!session) {
      await respond({
        text: "❌ Session expired. Please authenticate again with `/auth your-api-key`",
      });
      return;
    }

    // Get transfer details from button value
    const buttonValue = (body as any).actions[0].value;
    const { amount, recipientType, recipient } = JSON.parse(buttonValue);

    // Create API client
    const client = createAPIClient({
      baseURL: process.env['TRANZA_API_BASE_URL'] || 'http://localhost:8080',
      apiKey: session.apiKey,
    });

    // Create transfer
    const transfer = await createTransfer(client, {
      amount,
      recipient_type: recipientType,
      recipient_value: recipient,
      description: 'Transfer via Slack Bot',
    });

    await respond({
      text: `✅ *Transfer Initiated Successfully!*\n` +
            `Transfer ID: ${transfer.transfer_id}\n` +
            `Reference: ${transfer.reference_id}\n` +
            `Status: ${transfer.status}\n` +
            `Amount: ₹${transfer.amount}\n` +
            `Total: ₹${transfer.total_amount}\n` +
            `Estimated Time: ${transfer.estimated_time}`,
      replace_original: true
    });

  } catch (error) {
    console.error('Confirm transfer error:', error);
    await respond({
      text: "❌ Transfer failed. Please try again or contact support.",
      replace_original: true
    });
  }
};

// Button Action: Cancel Transfer
export const handleCancelTransfer = async ({ ack, respond }: SlackActionMiddlewareArgs) => {
  await ack();
  
  await respond({
    text: "❌ Transfer cancelled.",
    replace_original: true
  });
};

// Command: /logout - Clear user session
export const handleLogoutCommand = async ({ command, ack, respond }: SlackCommandMiddlewareArgs) => {
  await ack();

  try {
    const logoutSuccess = logoutUser(command.user_id);
    if (logoutSuccess) {
      await respond({
        text: "✅ You have been logged out successfully. Use `/auth your-api-key` to authenticate again.",
      });
    } else {
      await respond({
        text: "❌ No active session found.",
      });
    }
  } catch (error) {
    console.error('Logout command error:', error);
    await respond({
      text: "❌ Logout failed. Please try again.",
    });
  }
};

// Command: /help - Show available commands
export const handleHelpCommand = async ({ ack, respond }: SlackCommandMiddlewareArgs) => {
  await ack();

  const helpText = `🤖 *Tranza Bot Commands*

*Authentication:*
• \`/auth <api-key>\` - Authenticate with your API key
• \`/logout\` - Clear your session

*Wallet Operations:*
• \`/fetch-balance\` - Get your wallet balance
• \`/send-money <amount> <upi|phone> <recipient>\` - Send money
  Example: \`/send-money 100 upi user@paytm\`
  Example: \`/send-money 50 phone 9876543210\`

*Help:*
• \`/help\` - Show this help message

*Getting Started:*
1. Get your API key from the Tranza dashboard
2. Authenticate: \`/auth your-api-key\`
3. Start using wallet commands!

Need help? Contact our support team.`;

  await respond({
    text: helpText,
  });
};

// Register all commands and actions
export const registerCommands = (app: App) => {
  // Slash commands
  app.command('/auth', handleAuthCommand);
  app.command('/fetch-balance', handleFetchBalanceCommand);
  app.command('/send-money', handleSendMoneyCommand);
  app.command('/logout', handleLogoutCommand);
  app.command('/help', handleHelpCommand);
  app.command('/ping', handlePing);

  // Button actions
  app.action('confirm_transfer', handleConfirmTransfer);
  app.action('cancel_transfer', handleCancelTransfer);

  console.log('✅ All Slack bot commands registered successfully');
};
