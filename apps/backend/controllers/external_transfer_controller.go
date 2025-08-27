package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/zeusnotfound04/Tranza/models/dto"
	"github.com/zeusnotfound04/Tranza/services"
	"github.com/zeusnotfound04/Tranza/utils"
)

type ExternalTransferController struct {
	externalTransferService *services.ExternalTransferService
	walletService           *services.WalletService
}

func NewExternalTransferController(externalTransferService *services.ExternalTransferService, walletService *services.WalletService) *ExternalTransferController {
	return &ExternalTransferController{
		externalTransferService: externalTransferService,
		walletService:           walletService,
	}
}

// ValidateTransfer validates a transfer request before processing
// POST /api/transfers/validate
func (c *ExternalTransferController) ValidateTransfer(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(ctx, "User not authenticated")
		return
	}

	// Convert userID to string safely
	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		utils.InternalServerErrorResponse(ctx, "Invalid user ID type", nil)
		return
	}

	var req dto.ValidateTransferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(ctx, "Invalid request body", err)
		return
	}

	response, err := c.externalTransferService.ValidateTransferRequest(userUUID.String(), &req)
	if err != nil {
		utils.InternalServerErrorResponse(ctx, "Failed to validate transfer", err)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Transfer validation completed", response)
}

// CreateTransfer initiates a new external transfer
// POST /api/transfers
func (c *ExternalTransferController) CreateTransfer(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(ctx, "User not authenticated")
		return
	}

	// Convert userID to string safely
	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		utils.InternalServerErrorResponse(ctx, "Invalid user ID type", nil)
		return
	}

	var req dto.CreateExternalTransferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(ctx, "Invalid request body", err)
		return
	}

	response, err := c.externalTransferService.CreateExternalTransfer(userUUID.String(), &req)
	if err != nil {
		utils.BadRequestResponse(ctx, "Failed to create transfer", err)
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "Transfer initiated successfully", response)
}

// GetTransfer retrieves a specific transfer by ID
// GET /api/transfers/:id
func (c *ExternalTransferController) GetTransfer(ctx *gin.Context) {
	transferID := ctx.Param("id")
	if transferID == "" {
		utils.BadRequestResponse(ctx, "Transfer ID is required", nil)
		return
	}

	response, err := c.externalTransferService.GetExternalTransfer(transferID)
	if err != nil {
		utils.NotFoundResponse(ctx, "Transfer not found")
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Transfer retrieved successfully", response)
}

// GetUserTransfers retrieves transfers for the authenticated user with pagination
// GET /api/transfers
func (c *ExternalTransferController) GetUserTransfers(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(ctx, "User not authenticated")
		return
	}

	// Convert userID to string safely
	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		utils.InternalServerErrorResponse(ctx, "Invalid user ID type", nil)
		return
	}

	// Get pagination parameters
	page, limit := utils.GetPaginationParams(ctx)

	response, err := c.externalTransferService.GetExternalTransfersByUser(userUUID.String(), page, limit)
	if err != nil {
		utils.InternalServerErrorResponse(ctx, "Failed to retrieve transfers", err)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Transfers retrieved successfully", response)
}

// Bot-specific endpoints for Slack integration

// BotValidateTransfer validates a transfer request for bot users
// POST /api/bot/transfers/validate
func (c *ExternalTransferController) BotValidateTransfer(ctx *gin.Context) {
	// Bot authentication should be handled by API key middleware
	userID, exists := ctx.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(ctx, "Bot user not authenticated")
		return
	}

	// Convert userID to string safely
	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		utils.InternalServerErrorResponse(ctx, "Invalid user ID type", nil)
		return
	}

	var req dto.BotValidateTransferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(ctx, "Invalid request body", err)
		return
	}

	// Convert to standard validation request
	validateReq := &dto.ValidateTransferRequest{
		Amount:         req.Amount,
		RecipientType:  req.RecipientType,
		RecipientValue: req.RecipientValue,
	}

	response, err := c.externalTransferService.ValidateTransferRequest(userUUID.String(), validateReq)
	if err != nil {
		utils.InternalServerErrorResponse(ctx, "Failed to validate transfer", err)
		return
	}

	// Convert to bot response format
	botResponse := &dto.BotValidateTransferResponse{
		Valid:         response.Valid,
		Errors:        response.Errors,
		Warnings:      response.Warnings,
		TransferFee:   response.TransferFee,
		TotalAmount:   response.TotalAmount,
		EstimatedTime: response.EstimatedTime,
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Transfer validation completed", botResponse)
}

// BotCreateTransfer initiates a new external transfer for bot users
// POST /api/bot/transfers
func (c *ExternalTransferController) BotCreateTransfer(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(ctx, "Bot user not authenticated")
		return
	}

	// Convert userID to string safely
	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		utils.InternalServerErrorResponse(ctx, "Invalid user ID type", nil)
		return
	}

	var req dto.BotCreateTransferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(ctx, "Invalid request body", err)
		return
	}

	// Convert to standard transfer request
	transferReq := &dto.CreateExternalTransferRequest{
		Amount:         req.Amount,
		RecipientType:  req.RecipientType,
		RecipientValue: req.RecipientValue,
		RecipientName:  req.RecipientName,
		Description:    req.Description,
	}

	response, err := c.externalTransferService.CreateExternalTransfer(userUUID.String(), transferReq)
	if err != nil {
		utils.BadRequestResponse(ctx, "Failed to create transfer", err)
		return
	}

	// Convert to bot response format
	botResponse := &dto.BotTransferResponse{
		TransferID:    response.ID,
		ReferenceID:   response.ReferenceID,
		Amount:        response.Amount,
		TransferFee:   response.TransferFee,
		TotalAmount:   response.TotalAmount,
		Status:        response.Status,
		Recipient:     response.RecipientValue,
		EstimatedTime: response.EstimatedTime,
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "Transfer initiated successfully", botResponse)
}

// BotGetTransferStatus retrieves transfer status for bot users
// GET /api/bot/transfers/:id/status
func (c *ExternalTransferController) BotGetTransferStatus(ctx *gin.Context) {
	transferID := ctx.Param("id")
	if transferID == "" {
		utils.BadRequestResponse(ctx, "Transfer ID is required", nil)
		return
	}

	response, err := c.externalTransferService.GetExternalTransfer(transferID)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, "Transfer not found", err)
		return
	}

	// Convert to bot status response format
	botResponse := &dto.BotTransferStatusResponse{
		TransferID:    response.ID,
		ReferenceID:   response.ReferenceID,
		Status:        response.Status,
		Amount:        response.Amount,
		Recipient:     response.RecipientValue,
		EstimatedTime: response.EstimatedTime,
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Transfer status retrieved successfully", botResponse)
}


func (c *ExternalTransferController) BotGetWalletBalance(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(ctx, "Bot user not authenticated")
		return
	}

	// Convert userID to string (it should be UUID)
	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		utils.InternalServerErrorResponse(ctx, "Invalid user ID type", nil)
		return
	}

	// Get wallet using wallet service
	wallet, err := c.walletService.GetWalletByUserID(userUUID.String())
	if err != nil {
		utils.NotFoundResponse(ctx, "Wallet not found")
		return
	}

	// Create response structure for bot API
	response := map[string]interface{}{
		"user_id":     userUUID.String(),
		"balance":     wallet.Balance.String(),
		"currency":    wallet.Currency,
		"status":      wallet.Status,
		"message":     fmt.Sprintf("Current balance: ₹%s", wallet.Balance.String()),
		"balance_inr": wallet.Balance, // Numeric value for calculations
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Wallet balance retrieved successfully", response)
}

func (c *ExternalTransferController) GetTransferFees(ctx *gin.Context) {
	response := &dto.TransferFeesResponse{
		UPIFee:       utils.ParseDecimal("2.00"),   // ₹2 for UPI transfers
		PhoneFee:     utils.ParseDecimal("5.00"),   // ₹5 for phone transfers
		MinAmount:    utils.ParseDecimal("1.00"),   // ₹1 minimum
		MaxAmount:    utils.ParseDecimal("100000"), // ₹1,00,000 maximum
		DailyLimit:   utils.ParseDecimal("50000"),  // ₹50,000 daily limit
		MonthlyLimit: utils.ParseDecimal("200000"), // ₹2,00,000 monthly limit
		FeeStructure: []dto.FeeRange{
			{
				MinAmount: utils.ParseDecimal("1"),
				MaxAmount: utils.ParseDecimal("100000"),
				Fee:       utils.ParseDecimal("2"),
				FeeType:   "fixed",
			},
		},
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Transfer fees retrieved successfully", response)
}

// Health check endpoint for external transfer service
// GET /api/transfers/health
func (c *ExternalTransferController) HealthCheck(ctx *gin.Context) {
	response := map[string]interface{}{
		"service":   "external_transfer",
		"status":    "healthy",
		"timestamp": utils.GetCurrentTimestamp(),
		"version":   "1.0.0",
	}

	utils.SuccessResponse(ctx, http.StatusOK, "External transfer service is healthy", response)
}
