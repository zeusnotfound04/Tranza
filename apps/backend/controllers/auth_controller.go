package controllers

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/zeusnotfound04/Tranza/models"
	"github.com/zeusnotfound04/Tranza/services"
)

type AuthController struct {
	authService              *services.AuthService
	emailVerificationService *services.EmailVerificationService
}

// NewAuthController creates a new auth controller
func NewAuthController(authService *services.AuthService, emailVerificationService *services.EmailVerificationService) *AuthController {
	return &AuthController{
		authService:              authService,
		emailVerificationService: emailVerificationService,
	}
}

// Cookie helper functions for HttpOnly cookies
func (ac *AuthController) setAuthCookies(ctx *gin.Context, accessToken, refreshToken string, expiresIn int) {
	// Get domain from environment or use localhost for development
	domain := ""                            // Leave empty for localhost, set to your domain in production
	secure := gin.Mode() == gin.ReleaseMode // Only secure in production

	ctx.SetCookie("access_token", accessToken, expiresIn, "/", domain, secure, true)
	ctx.SetCookie("refresh_token", refreshToken, expiresIn*24*7, "/", domain, secure, true) // 7 days
}

func (ac *AuthController) clearAuthCookies(ctx *gin.Context) {
	domain := "" // Match the domain used in setAuthCookies
	secure := gin.Mode() == gin.ReleaseMode

	ctx.SetCookie("access_token", "", -1, "/", domain, secure, true)
	ctx.SetCookie("refresh_token", "", -1, "/", domain, secure, true)
}

func (ac *AuthController) getTokenFromCookie(ctx *gin.Context) string {
	token, err := ctx.Cookie("access_token")
	if err != nil {
		// Fallback to Authorization header
		authHeader := ctx.GetHeader("Authorization")
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			return authHeader[7:]
		}
		return ""
	}
	return token
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
		if err.Error() == "user already exists" {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}

	// Set HttpOnly cookies for immediate login after signup
	ac.setAuthCookies(ctx, response.AccessToken, response.RefreshToken, 3600)

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Signup successful",
		"user":    response.User,
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

	// Set HttpOnly cookies
	ac.setAuthCookies(ctx, response.AccessToken, response.RefreshToken, 3600) // 1 hour for access token

	ctx.JSON(http.StatusOK, gin.H{
		"message":       "Login successful",
		"user":          response.User,
		"access_token":  response.AccessToken,
		"refresh_token": response.RefreshToken,
		"expires_in":    response.ExpiresIn,
	})
}

// OAuthCallbackHandler handles OAuth callback
func (ac *AuthController) OAuthCallbackHandler(ctx *gin.Context) {
	var req models.OAuthCallbackRequest

	// Handle both GET (URL parameters) and POST (JSON) requests
	if ctx.Request.Method == "GET" {
		// Extract from URL parameters (OAuth provider redirect)
		req.Provider = ctx.Param("provider")
		if req.Provider == "" {
			// Handle direct provider callbacks like /auth/google/callback
			path := ctx.Request.URL.Path
			if strings.Contains(path, "/google/callback") {
				req.Provider = "google"
			} else if strings.Contains(path, "/github/callback") {
				req.Provider = "github"
			}
		}
		req.Code = ctx.Query("code")
		req.State = ctx.Query("state")

		// Build correct redirect URI based on the actual callback path
		if req.Provider != "" {
			req.RedirectURI = fmt.Sprintf("http://%s%s", ctx.Request.Host, ctx.Request.URL.Path)
		}
	} else {
		// Handle POST request with JSON body
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}
	}

	// Validate required fields
	if req.Provider == "" || req.Code == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Missing provider or authorization code"})
		return
	}

	response, err := ac.authService.HandleOAuthCallback(ctx.Request.Context(), req)
	if err != nil {
		// For GET requests, redirect to frontend with error
		if ctx.Request.Method == "GET" {
			frontendURL := os.Getenv("FRONTEND_URL")
			if frontendURL == "" {
				frontendURL = "http://localhost:3000"
			}
			ctx.Redirect(http.StatusTemporaryRedirect,
				fmt.Sprintf("%s/login?error=%s", frontendURL, url.QueryEscape(err.Error())))
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set HttpOnly cookies
	ac.setAuthCookies(ctx, response.AccessToken, response.RefreshToken, 3600) // 1 hour for access token

	// For GET requests, redirect to frontend with success
	if ctx.Request.Method == "GET" {
		frontendURL := os.Getenv("FRONTEND_URL")
		if frontendURL == "" {
			frontendURL = "http://localhost:3000"
		}
		ctx.Redirect(http.StatusTemporaryRedirect,
			fmt.Sprintf("%s/dashboard?oauth=success", frontendURL))
		return
	}

	// For POST requests, return JSON
	ctx.JSON(http.StatusOK, gin.H{
		"message":       "OAuth login successful",
		"user":          response.User,
		"access_token":  response.AccessToken,
		"refresh_token": response.RefreshToken,
		"expires_in":    response.ExpiresIn,
	})
}

// RefreshTokenHandler handles token refresh
func (ac *AuthController) RefreshTokenHandler(ctx *gin.Context) {
	// Try to get refresh token from cookie first
	refreshToken, err := ctx.Cookie("refresh_token")
	if err != nil {
		// Fallback to JSON body
		var req struct {
			RefreshToken string `json:"refresh_token" binding:"required"`
		}
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Refresh token required"})
			return
		}
		refreshToken = req.RefreshToken
	}

	response, err := ac.authService.RefreshToken(ctx.Request.Context(), refreshToken)
	if err != nil {
		ac.clearAuthCookies(ctx)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Set new HttpOnly cookies
	ac.setAuthCookies(ctx, response.AccessToken, response.RefreshToken, 3600)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Token refreshed successfully",
		"user":    response.User,
	})
}

// ValidateTokenHandler validates a JWT token
func (ac *AuthController) ValidateTokenHandler(ctx *gin.Context) {
	token := ac.getTokenFromCookie(ctx)
	if token == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	user, err := ac.authService.ValidateToken(ctx.Request.Context(), token)
	if err != nil {
		ac.clearAuthCookies(ctx)
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
		fmt.Printf("DEBUG: AuthMiddleware called for path: %s\n", ctx.Request.URL.Path)
		fmt.Printf("DEBUG: Request method: %s\n", ctx.Request.Method)
		fmt.Printf("DEBUG: All cookies: %+v\n", ctx.Request.Cookies())
		
		token := ac.getTokenFromCookie(ctx)
		fmt.Printf("DEBUG: Token from cookie: %s\n", token)
		
		if token == "" {
			fmt.Printf("DEBUG: No token found in cookies\n")
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			ctx.Abort()
			return
		}

		fmt.Printf("DEBUG: Validating token...\n")
		user, err := ac.authService.ValidateToken(ctx.Request.Context(), token)
		if err != nil {
			fmt.Printf("DEBUG: Token validation failed: %v\n", err)
			ac.clearAuthCookies(ctx)
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			ctx.Abort()
			return
		}

		fmt.Printf("DEBUG: Token validation successful. User: %+v\n", user)
		// Store user in context for use in subsequent handlers
		ctx.Set("user", user)
		ctx.Set("user_id", user.ID)
		ctx.Set("userID", user.ID.String()) // Add this for wallet controller compatibility
		fmt.Printf("DEBUG: Set userID in context: %s\n", user.ID.String())
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

	// Generate a state token for CSRF protection
	state := "oauth_" + provider

	// Get OAuth URL from auth service
	url, err := ac.authService.GetAuthURL(provider, state)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate OAuth URL",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"url":      url,
		"provider": provider,
		"state":    state,
	})
}

func (ac *AuthController) PreRegisterHandler(ctx *gin.Context) {
	var req models.PreRegistrationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	response, err := ac.emailVerificationService.InitiateEmailVerification(ctx.Request.Context(), req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Verification code sent successfully",
		"data":    response,
	})
}

// VerifyEmailHandler handles email verification and user creation
func (ac *AuthController) VerifyEmailHandler(ctx *gin.Context) {
	var req models.EmailVerificationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	response, err := ac.emailVerificationService.VerifyEmailCode(ctx.Request.Context(), req)
	if err != nil {
		if err.Error() == "verification code expired" ||
			err.Error() == "invalid or expired verification code" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
				"code":  "VERIFICATION_EXPIRED",
			})
		} else if err.Error() == "too many attempts" {
			ctx.JSON(http.StatusTooManyRequests, gin.H{
				"error": err.Error(),
				"code":  "TOO_MANY_ATTEMPTS",
			})
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Email verified successfully! You can now log in.",
		"data":    response,
	})
}

// ResendVerificationHandler handles resending verification codes
func (ac *AuthController) ResendVerificationHandler(ctx *gin.Context) {
	var req models.ResendVerificationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	response, err := ac.emailVerificationService.ResendVerificationCode(ctx.Request.Context(), req)
	if err != nil {
		if err.Error() == "rate limit exceeded" {
			ctx.JSON(http.StatusTooManyRequests, gin.H{
				"error": err.Error(),
				"code":  "RATE_LIMIT_EXCEEDED",
			})
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Verification code resent successfully",
		"data":    response,
	})
}

// GetCurrentUserHandler returns the current authenticated user
func (ac *AuthController) GetCurrentUserHandler(ctx *gin.Context) {
	// Get user from context (set by AuthMiddleware)
	userInterface, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	user, ok := userInterface.(*models.User)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user data"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "User retrieved successfully",
		"user":    user,
	})
}

// LogoutHandler handles user logout by clearing cookies
func (ac *AuthController) LogoutHandler(ctx *gin.Context) {
	ac.clearAuthCookies(ctx)
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Logout successful",
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
