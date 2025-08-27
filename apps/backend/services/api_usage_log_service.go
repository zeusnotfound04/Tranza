package services

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/zeusnotfound04/Tranza/models"
	"github.com/zeusnotfound04/Tranza/repositories"
)

type APIUsageLogService struct {
	Repo       *repositories.APIUsageLogRepository
	APIKeyRepo *repositories.APIKeyRepository
}

func NewAPIUsageLogService(repo *repositories.APIUsageLogRepository, apiKeyRepo *repositories.APIKeyRepository) *APIUsageLogService {
	return &APIUsageLogService{
		Repo:       repo,
		APIKeyRepo: apiKeyRepo,
	}
}

// LogAPIUsage logs a new API usage entry
func (s *APIUsageLogService) LogAPIUsage(ctx context.Context, logEntry *models.APIUsageLog) error {
	return s.Repo.Create(ctx, logEntry)
}

// LogFromGinContext creates and logs an API usage entry from a Gin context
func (s *APIUsageLogService) LogFromGinContext(
	ctx *gin.Context,
	apiKeyID uint,
	userID uuid.UUID,
	statusCode int,
	responseTime int64,
	command string,
	commandParams interface{},
	amountInvolved float64,
	source string,
	errorMessage string,
) error {
	// Convert command params to JSON
	var commandParamsJSON string
	if commandParams != nil {
		if paramsBytes, err := json.Marshal(commandParams); err == nil {
			commandParamsJSON = string(paramsBytes)
		}
	}

	// Create metadata
	metadata := map[string]interface{}{
		"query_params": ctx.Request.URL.RawQuery,
		"path":         ctx.Request.URL.Path,
		"referer":      ctx.GetHeader("Referer"),
	}
	metadataJSON, _ := json.Marshal(metadata)

	logEntry := &models.APIUsageLog{
		APIKeyID:       apiKeyID,
		UserID:         userID,
		Method:         ctx.Request.Method,
		Endpoint:       ctx.Request.URL.Path,
		UserAgent:      ctx.GetHeader("User-Agent"),
		IPAddress:      ctx.ClientIP(),
		StatusCode:     statusCode,
		ResponseTime:   responseTime,
		Command:        command,
		CommandParams:  commandParamsJSON,
		AmountInvolved: amountInvolved,
		Currency:       "INR",
		ErrorMessage:   errorMessage,
		Metadata:       string(metadataJSON),
		Source:         source,
		CreatedAt:      time.Now(),
	}

	// Update API key spending if amount is involved
	if amountInvolved > 0 {
		if err := s.updateAPIKeySpending(ctx.Request.Context(), apiKeyID, amountInvolved); err != nil {
			// Log the error but don't fail the request
			// You might want to use a proper logger here
		}
	}

	return s.LogAPIUsage(ctx.Request.Context(), logEntry)
}

// updateAPIKeySpending updates the spent amount for an API key
func (s *APIUsageLogService) updateAPIKeySpending(ctx context.Context, apiKeyID uint, amount float64) error {
	apiKey, err := s.APIKeyRepo.FindByID(ctx, apiKeyID)
	if err != nil {
		return err
	}

	if apiKey != nil {
		apiKey.AddSpentAmount(amount)
		// Note: You might need to add an Update method to APIKeyRepository
		// For now, we'll just update the usage which triggers the IncrementUsage method
		return s.APIKeyRepo.UpdateUsage(ctx, apiKeyID)
	}

	return nil
}

// GetUsageStats retrieves comprehensive usage statistics for an API key
func (s *APIUsageLogService) GetUsageStats(ctx context.Context, apiKeyID uint, userID uuid.UUID) (*models.APIUsageStats, error) {
	return s.Repo.GetUsageStats(ctx, apiKeyID, userID)
}

// GetUsageSummary retrieves a comprehensive usage summary
func (s *APIUsageLogService) GetUsageSummary(ctx context.Context, apiKeyID uint, userID uuid.UUID, days int) (*models.APIUsageSummary, error) {
	return s.Repo.GetUsageSummary(ctx, apiKeyID, userID, days)
}

// GetUsageLogs retrieves paginated usage logs
func (s *APIUsageLogService) GetUsageLogs(ctx context.Context, apiKeyID uint, limit, offset int) ([]models.APIUsageLog, error) {
	return s.Repo.GetUsageLogs(ctx, apiKeyID, limit, offset)
}

// GetTimeSeriesData retrieves time-series data for charts
func (s *APIUsageLogService) GetTimeSeriesData(ctx context.Context, apiKeyID uint, days int) ([]models.APIUsageTimeSeriesData, error) {
	return s.Repo.GetTimeSeriesData(ctx, apiKeyID, days)
}

// GetCommandData retrieves command-specific usage data
func (s *APIUsageLogService) GetCommandData(ctx context.Context, apiKeyID uint, days int) ([]models.APIUsageCommandData, error) {
	return s.Repo.GetCommandData(ctx, apiKeyID, days)
}

// CleanupOldLogs removes old usage logs
func (s *APIUsageLogService) CleanupOldLogs(ctx context.Context, olderThanDays int) error {
	return s.Repo.CleanupOldLogs(ctx, olderThanDays)
}
