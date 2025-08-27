package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/zeusnotfound04/Tranza/models"
	"gorm.io/gorm"
)

type APIUsageLogRepository struct {
	DB *gorm.DB
}

func NewAPIUsageLogRepository(db *gorm.DB) *APIUsageLogRepository {
	return &APIUsageLogRepository{
		DB: db,
	}
}

// Create logs a new API usage entry
func (r *APIUsageLogRepository) Create(ctx context.Context, log *models.APIUsageLog) error {
	return r.DB.WithContext(ctx).Create(log).Error
}

// GetUsageLogs retrieves usage logs for a specific API key with pagination
func (r *APIUsageLogRepository) GetUsageLogs(ctx context.Context, apiKeyID uint, limit, offset int) ([]models.APIUsageLog, error) {
	var logs []models.APIUsageLog
	err := r.DB.WithContext(ctx).
		Where("api_key_id = ?", apiKeyID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&logs).Error
	return logs, err
}

// GetUsageStats calculates comprehensive usage statistics for an API key
func (r *APIUsageLogRepository) GetUsageStats(ctx context.Context, apiKeyID uint, userID uuid.UUID) (*models.APIUsageStats, error) {
	var stats models.APIUsageStats

	// Get basic API key info
	var apiKey models.APIKey
	if err := r.DB.WithContext(ctx).Where("id = ? AND user_id = ?", apiKeyID, userID).First(&apiKey).Error; err != nil {
		return nil, err
	}

	stats.KeyID = apiKey.ID
	stats.Label = apiKey.Label
	stats.KeyType = apiKey.KeyType
	stats.RateLimit = apiKey.RateLimit
	stats.LastUsedAt = apiKey.LastUsedAt
	stats.SpendingLimit = apiKey.SpendingLimit
	stats.TotalAmountSpent = apiKey.SpentAmount
	stats.Currency = apiKey.Currency
	stats.RemainingLimit = apiKey.GetRemainingSpendingLimit()

	// Calculate request statistics
	var requestStats struct {
		TotalRequests      int64
		SuccessfulRequests int64
		FailedRequests     int64
		AvgResponseTime    float64
	}

	err := r.DB.WithContext(ctx).Model(&models.APIUsageLog{}).
		Select(`
			COUNT(*) as total_requests,
			COUNT(CASE WHEN status_code >= 200 AND status_code < 300 THEN 1 END) as successful_requests,
			COUNT(CASE WHEN status_code >= 400 THEN 1 END) as failed_requests,
			AVG(response_time) as avg_response_time
		`).
		Where("api_key_id = ?", apiKeyID).
		Scan(&requestStats).Error

	if err != nil {
		return nil, err
	}

	stats.TotalRequests = requestStats.TotalRequests
	stats.SuccessfulRequests = requestStats.SuccessfulRequests
	stats.FailedRequests = requestStats.FailedRequests
	stats.AverageResponseTime = requestStats.AvgResponseTime

	// Calculate time-based request counts
	now := time.Now()

	// Last 24 hours
	r.DB.WithContext(ctx).Model(&models.APIUsageLog{}).
		Where("api_key_id = ? AND created_at >= ?", apiKeyID, now.Add(-24*time.Hour)).
		Count(&stats.RequestsLast24Hours)

	// Last 7 days
	r.DB.WithContext(ctx).Model(&models.APIUsageLog{}).
		Where("api_key_id = ? AND created_at >= ?", apiKeyID, now.Add(-7*24*time.Hour)).
		Count(&stats.RequestsLast7Days)

	// Last 30 days
	r.DB.WithContext(ctx).Model(&models.APIUsageLog{}).
		Where("api_key_id = ? AND created_at >= ?", apiKeyID, now.Add(-30*24*time.Hour)).
		Count(&stats.RequestsLast30Days)

	// Get top commands
	topCommands, err := r.getTopCommands(ctx, apiKeyID, 10)
	if err != nil {
		return nil, err
	}
	stats.TopCommands = topCommands

	// Calculate rate limit usage (requests in the last hour)
	var requestsLastHour int64
	r.DB.WithContext(ctx).Model(&models.APIUsageLog{}).
		Where("api_key_id = ? AND created_at >= ?", apiKeyID, now.Add(-time.Hour)).
		Count(&requestsLastHour)

	if apiKey.RateLimit > 0 {
		stats.RateLimitUsage = float64(requestsLastHour) / float64(apiKey.RateLimit) * 100
	}

	return &stats, nil
}

// getTopCommands retrieves the most used commands for an API key
func (r *APIUsageLogRepository) getTopCommands(ctx context.Context, apiKeyID uint, limit int) ([]models.CommandUsage, error) {
	var commandUsages []models.CommandUsage

	rows, err := r.DB.WithContext(ctx).Raw(`
		SELECT 
			command,
			COUNT(*) as count,
			MAX(created_at) as last_used,
			(COUNT(CASE WHEN status_code >= 200 AND status_code < 300 THEN 1 END) * 100.0 / COUNT(*)) as success_rate
		FROM api_usage_logs 
		WHERE api_key_id = ? AND command != '' AND command IS NOT NULL
		GROUP BY command 
		ORDER BY count DESC 
		LIMIT ?
	`, apiKeyID, limit).Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var usage models.CommandUsage
		err := rows.Scan(&usage.Command, &usage.Count, &usage.LastUsed, &usage.SuccessRate)
		if err != nil {
			return nil, err
		}
		commandUsages = append(commandUsages, usage)
	}

	return commandUsages, nil
}

// GetTimeSeriesData retrieves time-series usage data for charts
func (r *APIUsageLogRepository) GetTimeSeriesData(ctx context.Context, apiKeyID uint, days int) ([]models.APIUsageTimeSeriesData, error) {
	var timeSeriesData []models.APIUsageTimeSeriesData

	query := `
		SELECT 
			DATE(created_at) as date,
			COUNT(*) as total_requests,
			COUNT(CASE WHEN status_code >= 200 AND status_code < 300 THEN 1 END) as success_requests,
			COUNT(CASE WHEN status_code >= 400 THEN 1 END) as failed_requests,
			COALESCE(SUM(amount_involved), 0) as total_amount,
			AVG(response_time) as avg_response_time
		FROM api_usage_logs 
		WHERE api_key_id = ? AND created_at >= DATE_SUB(NOW(), INTERVAL ? DAY)
		GROUP BY DATE(created_at) 
		ORDER BY date DESC
	`

	rows, err := r.DB.WithContext(ctx).Raw(query, apiKeyID, days).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var data models.APIUsageTimeSeriesData
		err := rows.Scan(&data.Date, &data.TotalRequests, &data.SuccessRequests,
			&data.FailedRequests, &data.TotalAmount, &data.AvgResponseTime)
		if err != nil {
			return nil, err
		}
		timeSeriesData = append(timeSeriesData, data)
	}

	return timeSeriesData, nil
}

// GetCommandData retrieves command-specific usage data
func (r *APIUsageLogRepository) GetCommandData(ctx context.Context, apiKeyID uint, days int) ([]models.APIUsageCommandData, error) {
	var commandData []models.APIUsageCommandData

	query := `
		SELECT 
			command,
			COUNT(*) as total_requests,
			COUNT(CASE WHEN status_code >= 200 AND status_code < 300 THEN 1 END) as success_requests,
			COUNT(CASE WHEN status_code >= 400 THEN 1 END) as failed_requests,
			COALESCE(SUM(amount_involved), 0) as total_amount,
			AVG(response_time) as avg_response_time,
			MAX(created_at) as last_used
		FROM api_usage_logs 
		WHERE api_key_id = ? AND command != '' AND command IS NOT NULL 
		AND created_at >= DATE_SUB(NOW(), INTERVAL ? DAY)
		GROUP BY command 
		ORDER BY total_requests DESC
	`

	rows, err := r.DB.WithContext(ctx).Raw(query, apiKeyID, days).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var data models.APIUsageCommandData
		err := rows.Scan(&data.Command, &data.TotalRequests, &data.SuccessRequests,
			&data.FailedRequests, &data.TotalAmount, &data.AvgResponseTime, &data.LastUsed)
		if err != nil {
			return nil, err
		}
		commandData = append(commandData, data)
	}

	return commandData, nil
}

// GetUsageSummary retrieves a comprehensive usage summary
func (r *APIUsageLogRepository) GetUsageSummary(ctx context.Context, apiKeyID uint, userID uuid.UUID, days int) (*models.APIUsageSummary, error) {
	// Get stats
	stats, err := r.GetUsageStats(ctx, apiKeyID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get usage stats: %w", err)
	}

	// Get time series data
	timeSeriesData, err := r.GetTimeSeriesData(ctx, apiKeyID, days)
	if err != nil {
		return nil, fmt.Errorf("failed to get time series data: %w", err)
	}

	// Get command data
	commandData, err := r.GetCommandData(ctx, apiKeyID, days)
	if err != nil {
		return nil, fmt.Errorf("failed to get command data: %w", err)
	}

	// Get recent logs
	recentLogs, err := r.GetUsageLogs(ctx, apiKeyID, 50, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent logs: %w", err)
	}

	return &models.APIUsageSummary{
		Stats:          *stats,
		TimeSeriesData: timeSeriesData,
		CommandData:    commandData,
		RecentLogs:     recentLogs,
	}, nil
}

// CleanupOldLogs removes logs older than the specified number of days
func (r *APIUsageLogRepository) CleanupOldLogs(ctx context.Context, olderThanDays int) error {
	cutoffDate := time.Now().AddDate(0, 0, -olderThanDays)
	return r.DB.WithContext(ctx).
		Where("created_at < ?", cutoffDate).
		Delete(&models.APIUsageLog{}).Error
}
