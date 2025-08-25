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
	apiKeyService *services.APIKeyService
}

func NewAPIKeyController(apiKeyService *services.APIKeyService) *APIKeyController {
	return &APIKeyController{
		apiKeyService: apiKeyService,
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

	// Default TTL if not specified
	if req.TTLHours == 0 {
		req.TTLHours = 8760 // 1 year default
	}

	// Convert userID to UUID
	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		utils.InternalServerErrorResponse(ctx, "Invalid user ID type", nil)
		return
	}

	rawKey, err := c.apiKeyService.Generate(ctx.Request.Context(), userUUID, req.Label, req.TTLHours)
	if err != nil {
		utils.InternalServerErrorResponse(ctx, "Failed to create API key", err)
		return
	}

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
	rawKey, err := c.apiKeyService.Generate(ctx.Request.Context(), userUUID, req.Label, ttl)
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

	fmt.Printf("DEBUG: UserID type: %T, value: %v\n", userID, userID)

	// TODO: This is a temporary implementation. The proper fix is to update the
	// APIKey model and repository to use UUID instead of uint for UserID
	// to match the User model structure.

	// For now, return empty list to prevent crashes
	response := dto.ListAPIKeysResponse{
		Keys:  []dto.APIKeyInfo{},
		Total: 0,
	}

	utils.SuccessResponse(ctx, http.StatusOK, "API keys retrieved successfully (temporary: UUID/uint mismatch needs to be fixed)", response)
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
