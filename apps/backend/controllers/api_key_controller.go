package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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

// CreateAPIKey creates a new API key for the authenticated user
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

	rawKey, err := c.apiKeyService.Generate(ctx.Request.Context(), userID.(uint), req.Label, req.TTLHours)
	if err != nil {
		utils.InternalServerErrorResponse(ctx, "Failed to create API key", err)
		return
	}

	response := dto.CreateAPIKeyResponse{
		APIKey:   rawKey,
		Label:    req.Label,
		TTLHours: req.TTLHours,
		Message:  "API key created successfully. Store it securely as it won't be shown again.",
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "API key created successfully", response)
}

// CreateBotAPIKey creates a new bot API key
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

	// Validate required fields for bot keys
	if req.WorkspaceID == "" {
		utils.BadRequestResponse(ctx, "WorkspaceID is required for bot keys", nil)
		return
	}

	if req.BotUserID == "" {
		utils.BadRequestResponse(ctx, "BotUserID is required for bot keys", nil)
		return
	}

	// Default TTL for bot keys is 1 year
	ttl := req.TTLHours
	if ttl == 0 {
		ttl = 8760 // 1 year
	}

	rawKey, err := c.apiKeyService.GenerateBotKey(
		ctx.Request.Context(),
		userID.(uint),
		req.Label,
		req.WorkspaceID,
		req.BotUserID,
		ttl,
	)
	if err != nil {
		utils.InternalServerErrorResponse(ctx, "Failed to create bot API key", err)
		return
	}

	response := dto.CreateBotAPIKeyResponse{
		APIKey:      rawKey,
		Label:       req.Label,
		WorkspaceID: req.WorkspaceID,
		BotUserID:   req.BotUserID,
		TTLHours:    ttl,
		Scopes: []string{
			"bot:transfer:validate",
			"bot:transfer:create",
			"bot:transfer:status",
			"bot:wallet:balance",
			"bot:user:info",
		},
		Message: "Bot API key created successfully. Configure your Slack bot with this key.",
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "Bot API key created successfully", response)
}

// GetAPIKeys lists all API keys for the authenticated user
// GET /api/keys
func (c *APIKeyController) GetAPIKeys(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(ctx, "User not authenticated")
		return
	}

	keys, err := c.apiKeyService.Repo.GetByUserID(ctx.Request.Context(), userID.(uint))
	if err != nil {
		utils.InternalServerErrorResponse(ctx, "Failed to retrieve API keys", err)
		return
	}

	var keyList []dto.APIKeyInfo
	for _, key := range keys {
		keyInfo := dto.APIKeyInfo{
			ID:         key.ID,
			Label:      key.Label,
			KeyType:    key.KeyType,
			Scopes:     key.GetScopes(),
			UsageCount: key.UsageCount,
			RateLimit:  key.RateLimit,
			IsActive:   key.IsActive,
			CreatedAt:  key.CreatedAt,
			ExpiresAt:  key.ExpiresAt,
			LastUsedAt: key.LastUsedAt,
		}

		// Include bot-specific fields if it's a bot key
		if key.IsBot() {
			keyInfo.BotWorkspace = &key.BotWorkspace
			keyInfo.BotUserID = &key.BotUserID
		}

		keyList = append(keyList, keyInfo)
	}

	response := dto.ListAPIKeysResponse{
		Keys:  keyList,
		Total: len(keyList),
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

	stats, err := c.apiKeyService.GetUsageStats(ctx.Request.Context(), uint(keyID), userID.(uint))
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

	newKey, err := c.apiKeyService.RotateKey(ctx.Request.Context(), uint(keyID), userID.(uint))
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

	err = c.apiKeyService.Revoke(ctx.Request.Context(), uint(keyID), userID.(uint))
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
	userIDStr := ctx.GetHeader("X-User-ID")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user"})
		return
	}

	var req RevokeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = c.apiKeyService.Revoke(ctx.Request.Context(), req.KeyID, uint(userID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not revoke key"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "API key revoked"})
}
