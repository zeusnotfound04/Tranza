# Tranza Slack Bot

A powerful Slack bot that enables secure financial operations directly from your Slack workspace. Users can authenticate with universal API keys, check wallet balances, and send money through UPI or phone numbers.

## Features

### 🔐 Authentication
- Secure API key authentication
- Session management with automatic expiration
- User session tracking and cleanup

### 💰 Wallet Operations
- Check wallet balance
- Real-time balance updates
- Secure API communication

### 💸 Money Transfers
- Send money to UPI IDs
- Send money to phone numbers
- Transfer validation before execution
- Real-time transfer status tracking
- Fee calculation and warnings

### 🤖 Bot Commands
- `/auth <api-key>` - Authenticate with your API key
- `/logout` - Log out and clear session
- `/fetch-balance` - Check your wallet balance
- `/send-money <amount> <type> <recipient> [name]` - Send money
- `/transfer-status <transfer-id>` - Check transfer status
- `/bot-status` - Show bot and session status
- `/help` - Show available commands

## Architecture

The bot follows a modern functional approach with clean separation of concerns:

### 📁 Project Structure
```
src/
├── app.ts                    # Main Slack bot application
├── clients/
│   └── tranza-api.ts        # API client for backend communication
└── services/
    └── user-session.ts      # User session management
```

### 🔧 Core Components

#### API Client (`clients/tranza-api.ts`)
- Functional approach using pure functions
- Axios-based HTTP client with interceptors
- Comprehensive error handling
- Type-safe request/response interfaces

#### Session Management (`services/user-session.ts`)
- Functional session management
- In-memory session storage
- Automatic session cleanup
- Session timeout handling

#### Bot Application (`app.ts`)
- Command handlers for all bot operations
- Event handlers for mentions and messages
- Error handling and logging
- Functional programming patterns

## Getting Started

### Prerequisites
- Node.js 18 or higher
- npm or pnpm
- Access to Tranza backend API
- Slack app with bot token

### Installation

1. **Install dependencies:**
   ```bash
   npm install
   ```

2. **Configure environment:**
   ```bash
   cp .env.example .env
   ```
   
   Edit `.env` with your configuration:
   ```env
   SLACK_BOT_TOKEN=xoxb-your-bot-token-here
   SLACK_SIGNING_SECRET=your-signing-secret-here
   TRANZA_API_BASE_URL=http://localhost:8080
   PORT=3000
   ```

3. **Build the project:**
   ```bash
   npm run build
   ```

4. **Start the bot:**
   ```bash
   npm start
   ```

   For development:
   ```bash
   npm run dev
   ```

## Usage Examples

### Authentication
```
/auth your-api-key-12345
```

### Check Balance
```
/fetch-balance
```

### Send Money
```
# Send to UPI ID
/send-money 100 upi john@paytm

# Send to phone number
/send-money 500 phone 9876543210

# Send with recipient name
/send-money 1000 upi merchant@paytm "John Doe"
```

### Check Transfer Status
```
/transfer-status txn_123456789
```

## Security Features

- **API Key Authentication**: Secure authentication using backend API keys
- **Session Management**: Time-limited sessions with automatic cleanup
- **Rate Limiting**: Built-in rate limiting in API client
- **Error Handling**: Comprehensive error handling and logging
- **Input Validation**: Strict validation of all user inputs

## Error Handling

The bot provides user-friendly error messages for common scenarios:
- Invalid API keys
- Expired sessions
- Network connectivity issues
- Invalid transfer parameters
- Backend service errors

## Development

### Available Scripts
- `npm run build` - Build TypeScript to JavaScript
- `npm run start` - Start the production bot
- `npm run dev` - Start in development mode with hot reload
- `npm run type-check` - Run TypeScript type checking
- `npm run lint` - Run ESLint
- `npm run format` - Format code with Prettier
- `npm test` - Run tests

### Code Style
- **Functional Programming**: Uses functional approach instead of classes
- **TypeScript**: Fully typed for better development experience
- **Modern ES6+**: Uses modern JavaScript features
- **Clean Code**: Well-structured with clear separation of concerns

## API Integration

The bot integrates with the Tranza backend through these endpoints:

### Bot-Specific Endpoints
- `POST /api/bot/transfers/validate` - Validate transfer
- `POST /api/bot/transfers` - Create transfer
- `GET /api/bot/transfers/:id/status` - Get transfer status
- `GET /api/bot/wallet/balance` - Get wallet balance

### Authentication
Uses `X-API-Key` header for authentication with universal API keys. You can create a universal API key from the Tranza web dashboard that will work with all features including this Slack bot.

## Configuration

### Environment Variables
- `SLACK_BOT_TOKEN` - Your Slack bot token
- `SLACK_SIGNING_SECRET` - Slack app signing secret
- `TRANZA_API_BASE_URL` - Backend API URL
- `PORT` - Server port (default: 3000)
- `SESSION_TIMEOUT_MINUTES` - Session timeout (default: 60)
- `MAX_SESSIONS` - Maximum concurrent sessions (default: 1000)

### Session Configuration
The session manager can be configured with:
- Session timeout duration
- Maximum number of concurrent sessions
- Cleanup interval frequency

## Monitoring

The bot includes comprehensive logging:
- Authentication attempts
- API calls with timing
- Error tracking
- Session lifecycle events
- Transfer operations

## Contributing

1. Follow the functional programming approach
2. Maintain TypeScript types
3. Add appropriate error handling
4. Include logging for important operations
5. Update tests for new features

## License

ISC License - see package.json for details.

## Support

For issues and questions:
1. Check the logs for error details
2. Verify API key permissions
3. Ensure backend connectivity
4. Contact the Tranza team for API-related issues
