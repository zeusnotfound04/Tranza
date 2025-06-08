package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zeusnotfound04/Tranza/services"
)

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

		ctx.Set("userID", apiKey.UserID)
		ctx.Next()
	}
}