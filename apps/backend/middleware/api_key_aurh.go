package middlewares

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zeusnotfound04/Tranza/services"
)

// RateLimiter implements a simple in-memory rate limiter
type RateLimiter struct {
	mu       sync.Mutex
	requests map[string][]time.Time
}

func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]time.Time),
	}
}

func (rl *RateLimiter) Allow(keyID string, limit int) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-time.Hour)

	// Clean old requests
	requests := rl.requests[keyID]
	validRequests := make([]time.Time, 0)
	for _, req := range requests {
		if req.After(windowStart) {
			validRequests = append(validRequests, req)
		}
	}

	// Check if we can make another request
	if len(validRequests) >= limit {
		return false
	}

	// Add current request
	validRequests = append(validRequests, now)
	rl.requests[keyID] = validRequests

	return true
}

var globalRateLimiter = NewRateLimiter()

// APIKeyAuthMiddleware provides basic API key authentication
func APIKeyAuthMiddleware(s *services.APIKeyService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		rawKey := ctx.GetHeader("X-API-Key")
		if rawKey == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "API key required"})
			return
		}

		apiKey, err := s.Validate(ctx.Request.Context(), rawKey)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid API key"})
			return
		}

		// Rate limiting
		keyIDStr := string(rune(apiKey.ID))
		if !globalRateLimiter.Allow(keyIDStr, apiKey.RateLimit) {
			ctx.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "rate limit exceeded",
				"limit": apiKey.RateLimit,
			})
			return
		}

		ctx.Set("user_id", apiKey.UserID)
		ctx.Set("api_key", apiKey)
		ctx.Next()
	}
}

// APIKeyAuthWithScopeMiddleware provides API key authentication with scope validation
func APIKeyAuthWithScopeMiddleware(s *services.APIKeyService, requiredScope string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		rawKey := ctx.GetHeader("X-API-Key")
		if rawKey == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "API key required"})
			return
		}

		apiKey, err := s.ValidateWithScope(ctx.Request.Context(), rawKey, requiredScope)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}

		// Rate limiting
		keyIDStr := string(rune(apiKey.ID))
		if !globalRateLimiter.Allow(keyIDStr, apiKey.RateLimit) {
			ctx.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "rate limit exceeded",
				"limit": apiKey.RateLimit,
			})
			return
		}

		ctx.Set("user_id", apiKey.UserID)
		ctx.Set("api_key", apiKey)
		ctx.Next()
	}
}

// BotAPIKeyAuthMiddleware accepts both bot and universal API keys for bot endpoints
func BotAPIKeyAuthMiddleware(s *services.APIKeyService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		rawKey := ctx.GetHeader("X-API-Key")
		if rawKey == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "API key required for bot operations"})
			return
		}

		apiKey, err := s.Validate(ctx.Request.Context(), rawKey)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid API key"})
			return
		}

		// Accept both bot keys and universal keys
		if !apiKey.IsBot() && apiKey.KeyType != "universal" {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "bot or universal API key required"})
			return
		}

		// Rate limiting
		keyIDStr := string(rune(apiKey.ID))
		if !globalRateLimiter.Allow(keyIDStr, apiKey.RateLimit) {
			ctx.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "rate limit exceeded",
				"limit": apiKey.RateLimit,
			})
			return
		}

		ctx.Set("user_id", apiKey.UserID)
		ctx.Set("api_key", apiKey)
		ctx.Set("bot_workspace", apiKey.BotWorkspace)
		ctx.Set("bot_user_id", apiKey.BotUserID)
		ctx.Next()
	}
}
