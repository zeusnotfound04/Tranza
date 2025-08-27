package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/zeusnotfound04/Tranza/models/dto"
	"github.com/zeusnotfound04/Tranza/services"
	"github.com/zeusnotfound04/Tranza/utils"
)

type APIKeyController struct {
	apiKeyService   *services.APIKeyService
	usageLogService *services.APIUsageLogService
}

func NewAPIKeyController(apiKeyService *services.APIKeyService, usageLogService *services.APIUsageLogService) *APIKeyController {
	return &APIKeyController{
		apiKeyService:   apiKeyService,
		usageLogService: usageLogService,
	}
}

// CreateAPIKey creates a new universal API key for the authenticated user
// POST /api/keys
func (c *APIKeyController) CreateAPIKey(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(ctx, "User not authenticated")
		return
	}

	var req dto.CreateAPIKeyRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(ctx, "Invalid request body", err)
		return
	}

	// Validate TTL
	if req.TTLHours < 0 || req.TTLHours > 8760 { // Max 1 year
		utils.BadRequestResponse(ctx, "TTL must be between 0 and 8760 hours", nil)
		return
	}

	// Validate password
	if len(req.Password) < 6 {
		utils.BadRequestResponse(ctx, "Password must be at least 6 characters long", nil)
		return
	}

	// Default TTL if not specified
	if req.TTLHours == 0 {
		req.TTLHours = 8760 // 1 year default
	}

	// Convert userID to UUID
	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		fmt.Printf("DEBUG CreateAPIKey: Invalid user ID type: %T, value: %v\n", userID, userID)
		utils.InternalServerErrorResponse(ctx, "Invalid user ID type", nil)
		return
	}

	fmt.Printf("DEBUG CreateAPIKey: Creating API key for user %s with label '%s' and TTL %d hours\n", userUUID, req.Label, req.TTLHours)

	rawKey, err := c.apiKeyService.Generate(ctx.Request.Context(), userUUID, req.Label, req.Password, req.TTLHours)
	if err != nil {
		fmt.Printf("DEBUG CreateAPIKey: Failed to generate API key: %v\n", err)
		utils.InternalServerErrorResponse(ctx, "Failed to create API key", err)
		return
	}

	fmt.Printf("DEBUG CreateAPIKey: Successfully created API key with raw key prefix: %s...\n", rawKey[:8])

	response := dto.CreateAPIKeyResponse{
		APIKey:   rawKey,
		Label:    req.Label,
		TTLHours: req.TTLHours,
		Message:  "Universal API key created successfully. This key works with all features including Slack bot integration. Store it securely as it won't be shown again.",
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "API key created successfully", response)
}

// CreateBotAPIKey creates a new universal API key (same as CreateAPIKey for backward compatibility)
// POST /api/keys/bot
func (c *APIKeyController) CreateBotAPIKey(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(ctx, "User not authenticated")
		return
	}

	var req dto.CreateBotAPIKeyRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(ctx, "Invalid request body", err)
		return
	}

	// Default TTL for keys is 1 year
	ttl := req.TTLHours
	if ttl == 0 {
		ttl = 8760 // 1 year
	}

	// Convert userID to UUID
	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		utils.InternalServerErrorResponse(ctx, "Invalid user ID type", nil)
		return
	}

	// Use the same Generate method - all keys are now universal
	rawKey, err := c.apiKeyService.Generate(ctx.Request.Context(), userUUID, req.Label, req.Password, ttl)
	if err != nil {
		utils.InternalServerErrorResponse(ctx, "Failed to create API key", err)
		return
	}

	response := dto.CreateBotAPIKeyResponse{
		APIKey:      rawKey,
		Label:       req.Label,
		WorkspaceID: req.WorkspaceID,
		BotUserID:   req.BotUserID,
		TTLHours:    ttl,
		Scopes: []string{
			"*", // Universal access
		},
		Message: "Universal API key created successfully. This key works with all features including Slack bot integration.",
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "API key created successfully", response)
}

// GetAPIKeys lists all API keys for the authenticated user
// GET /api/keys
func (c *APIKeyController) GetAPIKeys(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(ctx, "User not authenticated")
		return
	}

	// Now user_id should already be a UUID from the middleware
	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		fmt.Printf("DEBUG GetAPIKeys: Unexpected user ID type: %T, value: %v\n", userID, userID)
		utils.InternalServerErrorResponse(ctx, "Invalid user ID type", nil)
		return
	}

	fmt.Printf("DEBUG GetAPIKeys: Fetching API keys for user %s\n", userUUID)

	// Get all API keys for the user
	keys, err := c.apiKeyService.Repo.GetByUserID(ctx.Request.Context(), userUUID)
	if err != nil {
		fmt.Printf("DEBUG GetAPIKeys: Error fetching keys: %v\n", err)
		utils.InternalServerErrorResponse(ctx, "Failed to fetch API keys", err)
		return
	}

	fmt.Printf("DEBUG GetAPIKeys: Found %d keys for user %s\n", len(keys), userUUID)

	// Convert to response format
	var keyInfos []dto.APIKeyInfo
	for _, key := range keys {
		if key.IsActive {
			keyInfo := dto.APIKeyInfo{
				ID:         key.ID,
				Label:      key.Label,
				KeyType:    key.KeyType,
				Scopes:     key.GetScopes(),
				CreatedAt:  key.CreatedAt,
				ExpiresAt:  key.ExpiresAt,
				LastUsedAt: key.LastUsedAt,
				UsageCount: key.UsageCount,
				IsActive:   key.IsActive,
			}
			keyInfos = append(keyInfos, keyInfo)
		}
	}

	response := dto.ListAPIKeysResponse{
		Keys:  keyInfos,
		Total: len(keyInfos),
	}

	utils.SuccessResponse(ctx, http.StatusOK, "API keys retrieved successfully", response)
}

// GetAPIKeyUsage gets usage statistics for a specific API key
// GET /api/keys/:id/usage
func (c *APIKeyController) GetAPIKeyUsage(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(ctx, "User not authenticated")
		return
	}

	keyIDStr := ctx.Param("id")
	keyID, err := strconv.ParseUint(keyIDStr, 10, 32)
	if err != nil {
		utils.BadRequestResponse(ctx, "Invalid key ID", err)
		return
	}

	// Convert userID to UUID
	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		utils.InternalServerErrorResponse(ctx, "Invalid user ID type", nil)
		return
	}

	stats, err := c.apiKeyService.GetUsageStats(ctx.Request.Context(), uint(keyID), userUUID)
	if err != nil {
		utils.NotFoundResponse(ctx, "API key not found or access denied")
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Usage statistics retrieved successfully", stats)
}

// RotateAPIKey rotates an existing API key
// POST /api/keys/:id/rotate
func (c *APIKeyController) RotateAPIKey(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(ctx, "User not authenticated")
		return
	}

	keyIDStr := ctx.Param("id")
	keyID, err := strconv.ParseUint(keyIDStr, 10, 32)
	if err != nil {
		utils.BadRequestResponse(ctx, "Invalid key ID", err)
		return
	}

	// Convert userID to UUID
	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		utils.InternalServerErrorResponse(ctx, "Invalid user ID type", nil)
		return
	}

	newKey, err := c.apiKeyService.RotateKey(ctx.Request.Context(), uint(keyID), userUUID)
	if err != nil {
		utils.InternalServerErrorResponse(ctx, "Failed to rotate API key", err)
		return
	}

	response := dto.RotateAPIKeyResponse{
		NewAPIKey: newKey,
		Message:   "API key rotated successfully. Update your applications with the new key.",
	}

	utils.SuccessResponse(ctx, http.StatusOK, "API key rotated successfully", response)
}

// RevokeAPIKey revokes an API key
// DELETE /api/keys/:id
func (c *APIKeyController) RevokeAPIKey(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(ctx, "User not authenticated")
		return
	}

	keyIDStr := ctx.Param("id")
	keyID, err := strconv.ParseUint(keyIDStr, 10, 32)
	if err != nil {
		utils.BadRequestResponse(ctx, "Invalid key ID", err)
		return
	}

	// Convert userID to UUID
	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		utils.InternalServerErrorResponse(ctx, "Invalid user ID type", nil)
		return
	}

	err = c.apiKeyService.Revoke(ctx.Request.Context(), uint(keyID), userUUID)
	if err != nil {
		utils.InternalServerErrorResponse(ctx, "Failed to revoke API key", err)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "API key revoked successfully", nil)
}

// ViewAPIKey allows a user to view their API key by providing the password
// POST /api/keys/:id/view
func (c *APIKeyController) ViewAPIKey(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(ctx, "User not authenticated")
		return
	}

	keyIDStr := ctx.Param("id")
	keyID, err := strconv.ParseUint(keyIDStr, 10, 32)
	if err != nil {
		utils.BadRequestResponse(ctx, "Invalid key ID", err)
		return
	}

	var req dto.ViewAPIKeyRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(ctx, "Invalid request body", err)
		return
	}

	// Convert userID to UUID
	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		utils.InternalServerErrorResponse(ctx, "Invalid user ID type", nil)
		return
	}

	apiKey, err := c.apiKeyService.ViewAPIKey(ctx.Request.Context(), uint(keyID), userUUID, req.Password)
	if err != nil {
		if err.Error() == "invalid password" {
			utils.UnauthorizedResponse(ctx, "Invalid password")
			return
		}
		utils.InternalServerErrorResponse(ctx, "Failed to retrieve API key", err)
		return
	}

	response := dto.ViewAPIKeyResponse{
		APIKey:  apiKey,
		Message: "API key retrieved successfully",
	}

	utils.SuccessResponse(ctx, http.StatusOK, "API key retrieved successfully", response)
}

// Legacy method for backwards compatibility
func (c *APIKeyController) Create(ctx *gin.Context) {
	c.CreateAPIKey(ctx)
}

type CreateAPIKeyRequest struct {
	Label    string `json:"label" binding:"required"`
	TTLHours int    `json:"ttl_hours"`
}

type RevokeRequest struct {
	KeyID uint `json:"key_id"`
}

func (c *APIKeyController) Revoke(ctx *gin.Context) {
	// Legacy method - deprecated in favor of RevokeAPIKey
	// This method uses header-based auth which doesn't align with the new UUID system
	ctx.JSON(http.StatusBadRequest, gin.H{
		"error": "This endpoint is deprecated. Please use DELETE /api/v1/keys/:id instead",
	})
}

// GetDetailedUsageStats gets comprehensive usage statistics for a specific API key
// GET /api/keys/:id/usage/detailed
func (c *APIKeyController) GetDetailedUsageStats(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(ctx, "User not authenticated")
		return
	}

	keyIDStr := ctx.Param("id")
	keyID, err := strconv.ParseUint(keyIDStr, 10, 32)
	if err != nil {
		utils.BadRequestResponse(ctx, "Invalid key ID", err)
		return
	}

	// Convert userID to UUID
	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		utils.InternalServerErrorResponse(ctx, "Invalid user ID type", nil)
		return
	}

	// Get days parameter (default 30 days)
	daysStr := ctx.DefaultQuery("days", "30")
	days, err := strconv.Atoi(daysStr)
	if err != nil || days < 1 || days > 365 {
		days = 30
	}

	summary, err := c.usageLogService.GetUsageSummary(ctx.Request.Context(), uint(keyID), userUUID, days)
	if err != nil {
		utils.NotFoundResponse(ctx, "API key not found or access denied")
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Detailed usage statistics retrieved successfully", summary)
}

// GetUsageLogs gets paginated usage logs for a specific API key
// GET /api/keys/:id/logs
func (c *APIKeyController) GetUsageLogs(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(ctx, "User not authenticated")
		return
	}

	keyIDStr := ctx.Param("id")
	keyID, err := strconv.ParseUint(keyIDStr, 10, 32)
	if err != nil {
		utils.BadRequestResponse(ctx, "Invalid key ID", err)
		return
	}

	// Convert userID to UUID
	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		utils.InternalServerErrorResponse(ctx, "Invalid user ID type", nil)
		return
	}

	// Get pagination parameters
	limitStr := ctx.DefaultQuery("limit", "50")
	offsetStr := ctx.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 1000 {
		limit = 50
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	// Verify user owns this API key
	_, err = c.apiKeyService.GetUsageStats(ctx.Request.Context(), uint(keyID), userUUID)
	if err != nil {
		utils.NotFoundResponse(ctx, "API key not found or access denied")
		return
	}

	logs, err := c.usageLogService.GetUsageLogs(ctx.Request.Context(), uint(keyID), limit, offset)
	if err != nil {
		utils.InternalServerErrorResponse(ctx, "Failed to retrieve usage logs", err)
		return
	}

	response := map[string]interface{}{
		"logs":   logs,
		"limit":  limit,
		"offset": offset,
		"total":  len(logs),
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Usage logs retrieved successfully", response)
}

// GetTimeSeriesData gets time-series usage data for charts
// GET /api/keys/:id/usage/timeseries
func (c *APIKeyController) GetTimeSeriesData(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(ctx, "User not authenticated")
		return
	}

	keyIDStr := ctx.Param("id")
	keyID, err := strconv.ParseUint(keyIDStr, 10, 32)
	if err != nil {
		utils.BadRequestResponse(ctx, "Invalid key ID", err)
		return
	}

	// Convert userID to UUID
	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		utils.InternalServerErrorResponse(ctx, "Invalid user ID type", nil)
		return
	}

	// Get days parameter
	daysStr := ctx.DefaultQuery("days", "30")
	days, err := strconv.Atoi(daysStr)
	if err != nil || days < 1 || days > 365 {
		days = 30
	}

	// Verify user owns this API key
	_, err = c.apiKeyService.GetUsageStats(ctx.Request.Context(), uint(keyID), userUUID)
	if err != nil {
		utils.NotFoundResponse(ctx, "API key not found or access denied")
		return
	}

	timeSeriesData, err := c.usageLogService.GetTimeSeriesData(ctx.Request.Context(), uint(keyID), days)
	if err != nil {
		utils.InternalServerErrorResponse(ctx, "Failed to retrieve time series data", err)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Time series data retrieved successfully", timeSeriesData)
}

// GetCommandData gets command-specific usage data
// GET /api/keys/:id/usage/commands
func (c *APIKeyController) GetCommandData(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(ctx, "User not authenticated")
		return
	}

	keyIDStr := ctx.Param("id")
	keyID, err := strconv.ParseUint(keyIDStr, 10, 32)
	if err != nil {
		utils.BadRequestResponse(ctx, "Invalid key ID", err)
		return
	}

	// Convert userID to UUID
	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		utils.InternalServerErrorResponse(ctx, "Invalid user ID type", nil)
		return
	}

	// Get days parameter
	daysStr := ctx.DefaultQuery("days", "30")
	days, err := strconv.Atoi(daysStr)
	if err != nil || days < 1 || days > 365 {
		days = 30
	}

	// Verify user owns this API key
	_, err = c.apiKeyService.GetUsageStats(ctx.Request.Context(), uint(keyID), userUUID)
	if err != nil {
		utils.NotFoundResponse(ctx, "API key not found or access denied")
		return
	}

	commandData, err := c.usageLogService.GetCommandData(ctx.Request.Context(), uint(keyID), days)
	if err != nil {
		utils.InternalServerErrorResponse(ctx, "Failed to retrieve command data", err)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Command data retrieved successfully", commandData)
}
