package routes

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/zeusnotfound04/Tranza/controllers"
	middlewares "github.com/zeusnotfound04/Tranza/middleware"
	"github.com/zeusnotfound04/Tranza/pkg/razorpay"
	"github.com/zeusnotfound04/Tranza/repositories"
	"github.com/zeusnotfound04/Tranza/services"
	"github.com/zeusnotfound04/Tranza/utils"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	emailVerificationRepo := repositories.NewEmailVerificationRepository(db)
	cardRepo := repositories.NewCardRepository(db)
	txnRepo := repositories.NewTransactionRepository(db)
	walletRepo := repositories.NewWalletRepository(db)
	apiKeyRepo := repositories.NewAPIKeyRepository(db)

	// Initialize external clients
	razorpayClient := razorpay.NewClient(
		os.Getenv("RAZORPAY_KEY_ID"),
		os.Getenv("RAZORPAY_KEY_SECRET"),
	)

	// Initialize services
	jwtService := utils.NewJWTService(os.Getenv("JWT_SECRET"))
	emailService := services.NewEmailService()
	oauthService := services.NewOAuthServiceFromEnv()
	notificationService := services.NewNotificationService()

	// Initialize main services
	walletService := services.NewWalletService(walletRepo, txnRepo, razorpayClient, notificationService, db)
	authService := services.NewAuthService(userRepo, jwtService, oauthService, walletService)
	emailVerificationService := services.NewEmailVerificationService(emailVerificationRepo, userRepo, emailService)
	cardService := services.NewCardService(cardRepo)
	paymentService := services.NewPaymentService(razorpayClient, walletRepo, txnRepo, notificationService, db, os.Getenv("RAZORPAY_WEBHOOK_SECRET"))
	transactionService := services.NewTransactionService(txnRepo, walletRepo, paymentService)
	razorpayService := services.NewRazorpayService()
	apiKeyService := services.NewAPIKeyService(apiKeyRepo)
	aiService := services.NewAIService(db, os.Getenv("GEMINI_API_KEY"))

	// Initialize controllers
	authController := controllers.NewAuthController(authService, emailVerificationService)
	cardController := controllers.NewCardController(cardService)
	walletController := controllers.NewWalletHandler(walletService)
	transactionController := controllers.NewTransactionController(transactionService, paymentService)
	paymentController := controllers.NewPaymentController(razorpayService)
	apiKeyController := controllers.NewAPIKeyController(apiKeyService)
	aiController := controllers.NewAIController(aiService, walletService, paymentService)

	// ======================
	// Public Routes (No Auth Required)
	// ======================

	// Health check
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "pong"})
	})
	r.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"status": "healthy", "service": "tranza-api"})
	})

	// ======================
	// Authentication Routes
	// ======================
	auth := r.Group("/auth")
	{
		// Email verification flow (new registration system)
		auth.POST("/pre-register", authController.PreRegisterHandler)
		auth.POST("/verify-email", authController.VerifyEmailHandler)
		auth.POST("/resend-verification", authController.ResendVerificationHandler)

		// Authentication
		auth.POST("/login", authController.LoginHandler)
		auth.POST("/logout", authController.LogoutHandler)
		auth.POST("/refresh", authController.RefreshTokenHandler)
		auth.GET("/validate", authController.ValidateTokenHandler)
		auth.GET("/me", authController.AuthMiddleware(), authController.GetCurrentUserHandler)

		// OAuth routes
		auth.GET("/oauth/:provider", authController.GetOAuthURLHandler)            // Get OAuth URL
		auth.POST("/oauth/callback", authController.OAuthCallbackHandler)          // Handle OAuth callback via POST
		auth.GET("/oauth/:provider/callback", authController.OAuthCallbackHandler) // Handle OAuth callback via GET (for redirects)

		// Backward compatibility for direct provider callbacks (e.g., /auth/google/callback)
		auth.GET("/google/callback", authController.OAuthCallbackHandler)
		auth.GET("/github/callback", authController.OAuthCallbackHandler)

		// Legacy routes (deprecated but kept for backward compatibility)
		auth.POST("/register", controllers.SignupHandler)
		auth.POST("/signup", controllers.LoginHandler)
	}

	// ======================
	// API v1 Routes (Protected)
	// ======================
	api := r.Group("/api/v1")
	api.Use(authController.AuthMiddleware()) // Apply JWT auth to all API routes

	// ======================
	// User Profile Routes
	// ======================
	profile := api.Group("/profile")
	{
		profile.GET("", func(ctx *gin.Context) {
			// Get current user profile
			user, exists := ctx.Get("user")
			if !exists {
				ctx.JSON(500, gin.H{"error": "User not found in context"})
				return
			}
			ctx.JSON(200, gin.H{"user": user})
		})
		profile.PUT("", func(ctx *gin.Context) {
			// Update user profile
			ctx.JSON(200, gin.H{"message": "Profile updated successfully"})
		})
	}

	// ======================
	// Wallet Management Routes
	// ======================
	wallet := api.Group("/wallet")
	{
		wallet.GET("", walletController.GetWallet)                     // Get wallet details
		wallet.PUT("/settings", walletController.UpdateWalletSettings) // Update wallet settings
		wallet.POST("/load", walletController.CreateLoadMoneyOrder)    // Create load money order
		wallet.POST("/verify-payment", walletController.VerifyPayment) // Verify payment and credit wallet
	}

	// ======================
	// Card Management Routes
	// ======================
	cards := api.Group("/cards")
	{
		cards.POST("", cardController.LinkCard)             // Link a new card
		cards.GET("", cardController.GetCards)              // Get all user cards
		cards.DELETE("/:id", cardController.DeleteCard)     // Delete a card
		cards.PUT("/:id/limit", cardController.UpdateLimit) // Update card limit
	}

	// ======================
	// Transaction Routes
	// ======================
	transactions := api.Group("/transactions")
	{
		// Basic transaction operations
		transactions.GET("", transactionController.GetTransactionHistory)             // Get transaction history with pagination
		transactions.GET("/:id", transactionController.GetTransaction)                // Get specific transaction
		transactions.GET("/search", transactionController.SearchTransactions)         // Search transactions
		transactions.GET("/type/:type", transactionController.GetTransactionsByType)  // Get transactions by type
		transactions.GET("/:id/receipt", transactionController.GetTransactionReceipt) // Get transaction receipt

		// Transaction analytics and reporting
		transactions.GET("/stats", transactionController.GetTransactionStats)                    // Get transaction statistics
		transactions.GET("/analytics", transactionController.GetTransactionAnalytics)            // Get transaction analytics
		transactions.GET("/summary/monthly", transactionController.GetMonthlyTransactionSummary) // Monthly summary
		transactions.GET("/summary/daily", transactionController.GetDailyTransactionSummary)     // Daily summary
		transactions.GET("/trends", transactionController.GetTransactionTrends)                  // Transaction trends

		// Transaction export and admin functions
		transactions.GET("/export", transactionController.ExportTransactions)         // Export transactions
		transactions.POST("/:id/validate", transactionController.ValidateTransaction) // Validate transaction (admin)
		transactions.POST("/:id/retry", transactionController.RetryFailedTransaction) // Retry failed transaction
	}

	// ======================
	// Payment Routes (Razorpay Integration)
	// ======================
	payments := api.Group("/payments")
	{
		payments.POST("/orders", paymentController.CreateOrder)     // Create Razorpay order
		payments.POST("/verify", paymentController.VerifyPayment)   // Verify payment
		payments.GET("/orders/:id", paymentController.GetOrder)     // Get order details
		payments.GET("/payments/:id", paymentController.GetPayment) // Get payment details
	}

	// Payment webhooks (public, but secured with signature verification)
	webhooks := r.Group("/webhooks")
	{
		webhooks.POST("/razorpay", paymentController.HandleWebhook) // Razorpay webhook handler
	}

	// ======================
	// API Key Management Routes
	// ======================
	apiKeys := api.Group("/api-keys")
	{
		apiKeys.POST("", apiKeyController.Create)   // Create new API key
		apiKeys.DELETE("", apiKeyController.Revoke) // Revoke API key
	}

	// ======================
	// Protected Routes with API Key Authentication
	// ======================
	apiKeyRoutes := r.Group("/api/external")
	apiKeyRoutes.Use(middlewares.APIKeyAuthMiddleware(apiKeyService)) // Use API key auth instead of JWT
	{
		// External API endpoints that use API key authentication
		// These would be used by third-party integrations
		apiKeyRoutes.GET("/transactions", transactionController.GetTransactionHistory)
		apiKeyRoutes.GET("/wallet/balance", walletController.GetWallet)
		apiKeyRoutes.POST("/payments/create", paymentController.CreateOrder)
	}

	// ======================
	// AI Payment Routes
	// ======================
	ai := api.Group("/ai")
	{
		// AI Payment Processing
		ai.POST("/payment/request", aiController.ProcessPaymentRequest) // Process natural language payment request
		ai.POST("/payment/confirm", aiController.ConfirmPayment)        // Confirm or cancel payment request
		ai.GET("/payment/:id", aiController.GetPaymentRequest)          // Get specific payment request details
		ai.DELETE("/payment/:id", aiController.CancelPaymentRequest)    // Cancel pending payment request

		// AI Payment History and Analytics
		ai.GET("/payments", aiController.GetPaymentHistory)     // Get AI payment history with pagination
		ai.GET("/analytics", aiController.GetSpendingAnalytics) // Get AI spending analytics and insights

		// AI Spending Limits Management
		ai.GET("/limits", aiController.GetSpendingLimits)    // Get user's AI spending limits
		ai.PUT("/limits", aiController.UpdateSpendingLimits) // Update user's AI spending limits
	}

	// ======================
	// Admin Routes (Future Implementation)
	// ======================
	admin := api.Group("/admin")
	// admin.Use(middlewares.AdminAuthMiddleware()) // Future: Admin-only middleware
	{
		admin.GET("/users", func(ctx *gin.Context) {
			ctx.JSON(501, gin.H{"message": "Admin routes not implemented yet"})
		})
		admin.GET("/transactions", func(ctx *gin.Context) {
			ctx.JSON(501, gin.H{"message": "Admin routes not implemented yet"})
		})
		admin.GET("/analytics", func(ctx *gin.Context) {
			ctx.JSON(501, gin.H{"message": "Admin routes not implemented yet"})
		})
	}
}
