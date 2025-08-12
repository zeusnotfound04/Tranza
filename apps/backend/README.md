# Tranza Backend API

A comprehensive financial transaction platform built with Go, Gin, and GORM.

## 🚀 Quick Start

### Prerequisites
- Go 1.19 or higher
- PostgreSQL database
- Razorpay account (for payments)

### Setup

1. **Clone and Navigate**
   ```bash
   cd apps/backend
   ```

2. **Environment Configuration**
   ```bash
   # Copy example environment file
   cp .env.example .env
   
   # Update .env with your configuration
   nano .env  # or use your preferred editor
   ```

3. **Install Dependencies**
   ```bash
   go mod tidy
   ```

4. **Run the Server**
   
   **Linux/macOS:**
   ```bash
   chmod +x start.sh
   ./start.sh
   ```
   
   **Windows:**
   ```cmd
   start.bat
   ```
   
   **Or directly:**
   ```bash
   go run cmd/server/main.go
   ```

## 📊 API Features

### 🔐 Authentication System
- Email verification flow
- JWT cookie-based authentication
- OAuth integration (Google, GitHub)
- API key management for external access

### 💰 Financial Operations
- **Wallet Management**: Balance, settings, load money
- **Card Management**: Link cards, set limits, manage payment methods
- **Transaction Processing**: Create, track, analyze transactions
- **Payment Integration**: Razorpay payment gateway

### 📈 Analytics & Reporting
- Transaction statistics and trends
- Monthly/daily summaries
- Export functionality (CSV/PDF)
- Real-time analytics

### 🔑 API Access
- RESTful API design
- JWT authentication for web apps
- API key authentication for external integrations
- Comprehensive webhook support

## 🛠️ Available Endpoints

### Public Routes
- `GET /ping` - Health check
- `GET /health` - Detailed health status

### Authentication (`/auth`)
- `POST /auth/pre-register` - Start email verification
- `POST /auth/verify-email` - Complete registration
- `POST /auth/login` - User login
- `POST /auth/logout` - User logout
- `GET /auth/me` - Get current user
- OAuth routes for Google/GitHub

### Protected API (`/api/v1`) 🔒
- `/api/v1/profile` - User profile management
- `/api/v1/wallet` - Wallet operations
- `/api/v1/cards` - Card management
- `/api/v1/transactions` - Transaction operations
- `/api/v1/payments` - Payment processing
- `/api/v1/api-keys` - API key management

### External API (`/api/external`) 🔑
- API key authenticated endpoints for integrations

### Webhooks
- `POST /webhooks/razorpay` - Razorpay webhook handler

## 📁 Project Structure

```
apps/backend/
├── cmd/
│   └── server/          # Application entry point
├── config/              # Configuration management
├── controllers/         # HTTP handlers
├── middleware/          # Custom middleware
├── models/             # Database models
├── repositories/       # Data access layer
├── routes/             # Route definitions
├── services/           # Business logic
├── utils/              # Utility functions
├── pkg/                # External packages
└── .env.example        # Environment template
```

## 🔧 Configuration

Key environment variables:

```env
# Server
PORT=8080

# Database
DB_HOST=localhost
DB_NAME=tranza_db

# Security
JWT_SECRET=your-secret-key

# Payments
RAZORPAY_KEY_ID=your-key
RAZORPAY_KEY_SECRET=your-secret

# Frontend
FRONTEND_URL=http://localhost:3000
```

## 🐳 Docker Support

```bash
# Build image
docker build -t tranza-api .

# Run container
docker run -p 8080:8080 --env-file .env tranza-api
```

## 📚 Documentation

- **API Routes**: See [API_ROUTES_DOCUMENTATION.md](../../API_ROUTES_DOCUMENTATION.md)
- **Postman Collection**: Import from `/docs/postman/`
- **OpenAPI Spec**: Available at `/docs/swagger.json`

## 🧪 Testing

```bash
# Run tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package
go test ./services/...
```

## 🚦 Health Monitoring

- Health endpoint: `GET /health`
- Metrics endpoint: `GET /metrics` (if enabled)
- Application logs: Structured JSON logging

## 🔒 Security Features

- CORS configuration
- Rate limiting
- Input validation
- SQL injection prevention
- XSS protection
- CSRF protection via cookies

## 📦 Deployment

### Production Checklist
- [ ] Update environment variables
- [ ] Enable HTTPS
- [ ] Configure database connection pooling
- [ ] Set up monitoring and logging
- [ ] Configure backup strategy
- [ ] Set up CI/CD pipeline

## 🤝 Contributing

1. Fork the repository
2. Create feature branch
3. Make changes
4. Run tests
5. Submit pull request

## 📄 License

This project is licensed under the MIT License.

---

**🌟 Happy Coding with Tranza API!**
