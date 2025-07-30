package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zeusnotfound04/Tranza/models"
	"github.com/zeusnotfound04/Tranza/services"
)

/*
Example usage in routes:

	// Initialize dependencies
	db := config.ConnectDB()
	userRepo := repositories.NewUserRepository(db)
	jwtService := utils.NewJWTService(os.Getenv("JWT_SECRET"))
	oauthService := services.NewOAuthServiceFromEnv()
	authService := services.NewAuthService(userRepo, jwtService, oauthService)
	authController := controllers.NewAuthController(authService)

	// Setup routes
	auth := router.Group("/auth")
	{
		auth.POST("/signup", authController.SignupHandler)
		auth.POST("/login", authController.LoginHandler)
		auth.POST("/refresh", authController.RefreshTokenHandler)
		auth.GET("/validate", authController.ValidateTokenHandler)
		auth.POST("/oauth/callback", authController.OAuthCallbackHandler)
		auth.GET("/oauth/:provider", authController.GetOAuthURLHandler)
	}

	// Protected routes example
	protected := router.Group("/api")
	protected.Use(authController.AuthMiddleware())
	{
		protected.GET("/profile", getUserProfile)
		protected.PUT("/profile", updateUserProfile)
	}
*/

// AuthController handles authentication related requests
type AuthController struct {
	authService *services.AuthService
}

// NewAuthController creates a new auth controller
func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

// SignupHandler handles user registration
func (ac *AuthController) SignupHandler(ctx *gin.Context) {
	var req models.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	response, err := ac.authService.Register(ctx.Request.Context(), req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Signup successful",
		"data":    response,
	})
}

// LoginHandler handles user login
func (ac *AuthController) LoginHandler(ctx *gin.Context) {
	var req models.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	response, err := ac.authService.Login(ctx.Request.Context(), req)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"data":    response,
	})
}

// OAuthCallbackHandler handles OAuth callback
func (ac *AuthController) OAuthCallbackHandler(ctx *gin.Context) {
	var req models.OAuthCallbackRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	response, err := ac.authService.HandleOAuthCallback(ctx.Request.Context(), req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "OAuth login successful",
		"data":    response,
	})
}

// RefreshTokenHandler handles token refresh
func (ac *AuthController) RefreshTokenHandler(ctx *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	response, err := ac.authService.RefreshToken(ctx.Request.Context(), req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Token refreshed successfully",
		"data":    response,
	})
}

// ValidateTokenHandler validates a JWT token
func (ac *AuthController) ValidateTokenHandler(ctx *gin.Context) {
	token := ctx.GetHeader("Authorization")
	if token == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
		return
	}

	// Remove "Bearer " prefix if present
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	user, err := ac.authService.ValidateToken(ctx.Request.Context(), token)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Token is valid",
		"user":    user,
	})
}

// AuthMiddleware provides JWT authentication middleware
func (ac *AuthController) AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")
		if token == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			ctx.Abort()
			return
		}

		// Remove "Bearer " prefix if present
		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}

		user, err := ac.authService.ValidateToken(ctx.Request.Context(), token)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			ctx.Abort()
			return
		}

		// Store user in context for use in subsequent handlers
		ctx.Set("user", user)
		ctx.Set("user_id", user.ID)
		ctx.Next()
	}
}

// GetOAuthURLHandler generates OAuth URLs for different providers
func (ac *AuthController) GetOAuthURLHandler(ctx *gin.Context) {
	provider := ctx.Param("provider")
	if provider != "google" && provider != "github" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported provider. Use 'google' or 'github'"})
		return
	}

	state := ctx.Query("state")
	if state == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "State parameter is required"})
		return
	}

	// This would require adding GetAuthURL method to auth service
	// For now, return a placeholder response
	ctx.JSON(http.StatusNotImplemented, gin.H{
		"error":    "OAuth URL generation not implemented",
		"message":  "Please implement GetAuthURL in auth service",
		"provider": provider,
		"state":    state,
	})
}

// Legacy functions for backward compatibility
func SignupHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusNotImplemented, gin.H{
		"error":   "This endpoint is deprecated",
		"message": "Please use the new AuthController endpoints",
	})
}

func LoginHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusNotImplemented, gin.H{
		"error":   "This endpoint is deprecated",
		"message": "Please use the new AuthController endpoints",
	})
}
