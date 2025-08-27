package middlewares

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/zeusnotfound04/Tranza/models"
	"github.com/zeusnotfound04/Tranza/services"
)

// APIUsageLoggingMiddleware creates middleware that logs detailed API usage
func APIUsageLoggingMiddleware(usageLogService *services.APIUsageLogService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Record start time
		startTime := time.Now()

		// Create a custom response writer to capture response data
		writer := &responseWriter{
			ResponseWriter: ctx.Writer,
			body:           &strings.Builder{},
		}
		ctx.Writer = writer

		// Process the request
		ctx.Next()

		// Calculate response time
		responseTime := time.Since(startTime).Milliseconds()

		// Get API key info from context (set by API key auth middleware)
		apiKeyInterface, exists := ctx.Get("api_key")
		if !exists {
			return // No API key, skip logging
		}

		apiKey, ok := apiKeyInterface.(*models.APIKey)
		if !ok {
			return // Invalid API key type
		}

		// Get user ID
		userID, exists := ctx.Get("user_id")
		if !exists {
			return
		}

		userUUID, ok := userID.(uuid.UUID)
		if !ok {
			return
		}

		// Extract command and amount information based on the endpoint
		command, commandParams, amountInvolved := extractCommandInfo(ctx)

		// Determine source
		source := determineSource(ctx)

		// Get error message if any
		errorMessage := ""
		if ctx.Writer.Status() >= 400 {
			errorMessage = writer.body.String()
		}

		// Log the usage asynchronously to avoid impacting response time
		go func() {
			err := usageLogService.LogFromGinContext(
				ctx,
				apiKey.ID,
				userUUID,
				ctx.Writer.Status(),
				responseTime,
				command,
				commandParams,
				amountInvolved,
				source,
				errorMessage,
			)
			if err != nil {
				// Log the error (you might want to use a proper logger here)
				// For now, we'll just print it
				// log.Printf("Failed to log API usage: %v", err)
			}
		}()
	}
}

// responseWriter wraps gin.ResponseWriter to capture response body
type responseWriter struct {
	gin.ResponseWriter
	body *strings.Builder
}

func (w *responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// extractCommandInfo extracts command and transaction information from the request
func extractCommandInfo(ctx *gin.Context) (string, interface{}, float64) {
	path := ctx.Request.URL.Path
	method := ctx.Request.Method

	var command string
	var commandParams interface{}
	var amountInvolved float64

	// Map endpoints to commands
	switch {
	case strings.Contains(path, "/wallet/balance"):
		command = "fetch-balance"
	case strings.Contains(path, "/external-transfer"):
		command = "send-money"
		// Try to extract amount from request body
		if amount := extractAmountFromRequest(ctx); amount > 0 {
			amountInvolved = amount
		}
	case strings.Contains(path, "/payment/create"):
		command = "create-payment"
		if amount := extractAmountFromRequest(ctx); amount > 0 {
			amountInvolved = amount
		}
	case strings.Contains(path, "/keys") && method == "POST":
		command = "create-api-key"
	case strings.Contains(path, "/keys") && method == "GET":
		command = "list-api-keys"
	case strings.Contains(path, "/keys") && method == "DELETE":
		command = "revoke-api-key"
	case strings.Contains(path, "/auth"):
		command = "authenticate"
	default:
		command = method + " " + path
	}

	// Extract query parameters as command params
	if len(ctx.Request.URL.RawQuery) > 0 {
		commandParams = map[string]interface{}{
			"query": ctx.Request.URL.RawQuery,
			"path":  path,
		}
	}

	return command, commandParams, amountInvolved
}

// extractAmountFromRequest tries to extract amount from various request formats
func extractAmountFromRequest(ctx *gin.Context) float64 {
	// This is a simplified implementation
	// You might need to adjust based on your actual request structures

	// Try to get from JSON body
	var body map[string]interface{}
	if err := ctx.ShouldBindJSON(&body); err == nil {
		if amount, exists := body["amount"]; exists {
			switch v := amount.(type) {
			case float64:
				return v
			case string:
				// Parse string amount if needed
				// You might want to use strconv.ParseFloat here
			}
		}
	}

	// Try to get from query parameters
	if amountStr := ctx.Query("amount"); amountStr != "" {
		// Parse amount from query
		// You might want to use strconv.ParseFloat here
	}

	return 0
}

// determineSource determines the source of the API request
func determineSource(ctx *gin.Context) string {
	userAgent := ctx.GetHeader("User-Agent")
	referer := ctx.GetHeader("Referer")

	switch {
	case strings.Contains(userAgent, "Slack"):
		return "slack-bot"
	case strings.Contains(referer, "dashboard"):
		return "web-dashboard"
	case strings.Contains(userAgent, "curl") || strings.Contains(userAgent, "Postman"):
		return "api-client"
	default:
		return "api"
	}
}

// SlackBotUsageLoggingMiddleware specifically for Slack bot requests
func SlackBotUsageLoggingMiddleware(usageLogService *services.APIUsageLogService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		startTime := time.Now()

		// Process the request
		ctx.Next()

		responseTime := time.Since(startTime).Milliseconds()

		// Get API key info from context
		apiKeyInterface, exists := ctx.Get("api_key")
		if !exists {
			return
		}

		apiKey, ok := apiKeyInterface.(*models.APIKey)
		if !ok {
			return
		}

		userID, exists := ctx.Get("user_id")
		if !exists {
			return
		}

		userUUID, ok := userID.(uuid.UUID)
		if !ok {
			return
		}

		// Extract Slack-specific command info
		command, commandParams, amountInvolved := extractSlackCommandInfo(ctx)

		// Log asynchronously
		go func() {
			err := usageLogService.LogFromGinContext(
				ctx,
				apiKey.ID,
				userUUID,
				ctx.Writer.Status(),
				responseTime,
				command,
				commandParams,
				amountInvolved,
				"slack-bot",
				"",
			)
			if err != nil {
				// Log error
			}
		}()
	}
}

// extractSlackCommandInfo extracts Slack-specific command information
func extractSlackCommandInfo(ctx *gin.Context) (string, interface{}, float64) {
	// This would be called from your Slack bot endpoints
	// You can customize this based on how you structure your Slack bot API

	path := ctx.Request.URL.Path
	var command string
	var amountInvolved float64

	// Map Slack bot endpoints to commands
	switch {
	case strings.Contains(path, "/bot/balance"):
		command = "/fetch-balance"
	case strings.Contains(path, "/bot/transfer"):
		command = "/send-money"
		// Extract amount from Slack command
	case strings.Contains(path, "/bot/auth"):
		command = "/auth"
	default:
		command = "slack-" + path
	}

	commandParams := map[string]interface{}{
		"workspace_id": ctx.GetHeader("X-Slack-Workspace-ID"),
		"user_id":      ctx.GetHeader("X-Slack-User-ID"),
		"channel_id":   ctx.GetHeader("X-Slack-Channel-ID"),
	}

	return command, commandParams, amountInvolved
}
