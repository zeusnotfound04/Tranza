package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zeusnotfound04/Tranza/services"
)

type APIKeyController struct {
	Service *services.APIKeyService
}

func NewAPIKeyController(s *services.APIKeyService) *APIKeyController {
	return &APIKeyController{Service: s}
}

type CreateAPIKeyRequest struct {
	Label    string `json:"label"`
	TTLHours int    `json:"ttl_hours"`
}

func (c *APIKeyController) Create(ctx *gin.Context) {
	userIDStr := ctx.GetHeader("X-User-ID")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user"})
		return
	}

	var req CreateAPIKeyRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	key, err := c.Service.Generate(ctx.Request.Context(), uint(userID), req.Label, req.TTLHours)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate API key"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"api_key": key})
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

	err = c.Service.Revoke(ctx.Request.Context(), req.KeyID, uint(userID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not revoke key"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "API key revoked"})
}