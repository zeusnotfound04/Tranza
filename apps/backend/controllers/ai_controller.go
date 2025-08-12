package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/zeusnotfound04/Tranza/models"
	"github.com/zeusnotfound04/Tranza/services"
	"github.com/zeusnotfound04/Tranza/utils"
)

type AIController struct {
	aiService      *services.AIService
	walletService  *services.WalletService
	paymentService *services.PaymentService
}

func NewAIController(
	aiService *services.AIService,
	walletService *services.WalletService,
	paymentService *services.PaymentService,
) *AIController {
	return &AIController{
		aiService:      aiService,
		walletService:  walletService,
		paymentService: paymentService,
	}
}

// ProcessPaymentRequest handles natural language payment requests
func (ac *AIController) ProcessPaymentRequest(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req models.AIPaymentRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert userID to UUID
	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Process the AI payment request
	response, err := ac.aiService.ProcessPaymentRequest(userUUID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Payment request processed", response)
}

// ConfirmPayment handles payment confirmation
func (ac *AIController) ConfirmPayment(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req models.AIPaymentConfirmationDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert userID to UUID
	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Parse payment ID
	paymentUUID, err := uuid.Parse(req.PaymentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment ID"})
		return
	}

	// Process confirmation
	result, err := ac.aiService.ConfirmPayment(userUUID, paymentUUID, req.Confirmed)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if req.Confirmed {
		utils.SuccessResponse(c, http.StatusOK, "Payment confirmed and processed", result)
	} else {
		utils.SuccessResponse(c, http.StatusOK, "Payment cancelled", result)
	}
}

// GetPaymentHistory returns AI payment request history
func (ac *AIController) GetPaymentHistory(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Convert userID to UUID
	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Parse query parameters
	page := 1
	limit := 20
	status := c.Query("status")

	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	// Get payment history
	history, total, err := ac.aiService.GetPaymentHistory(userUUID, page, limit, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := gin.H{
		"history": history,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	}

	utils.SuccessResponse(c, http.StatusOK, "Payment history retrieved", response)
}

// GetSpendingLimits returns user's AI spending limits
func (ac *AIController) GetSpendingLimits(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Convert userID to UUID
	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Get spending limits
	limits, err := ac.aiService.GetSpendingLimits(userUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Spending limits retrieved", limits)
}

// UpdateSpendingLimits updates user's AI spending limits
func (ac *AIController) UpdateSpendingLimits(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req models.AISpendingLimitsDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert userID to UUID
	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Update spending limits
	limits, err := ac.aiService.UpdateSpendingLimits(userUUID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Spending limits updated", limits)
}

// GetSpendingAnalytics returns spending analytics and insights
func (ac *AIController) GetSpendingAnalytics(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Convert userID to UUID
	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Parse query parameters for date range
	period := c.DefaultQuery("period", "month") // day, week, month, year

	// Get analytics
	analytics, err := ac.aiService.GetSpendingAnalytics(userUUID, period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Spending analytics retrieved", analytics)
}

// GetPaymentRequest returns details of a specific payment request
func (ac *AIController) GetPaymentRequest(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	paymentID := c.Param("id")
	if paymentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Payment ID is required"})
		return
	}

	// Convert IDs to UUIDs
	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	paymentUUID, err := uuid.Parse(paymentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment ID"})
		return
	}

	// Get payment request details
	payment, err := ac.aiService.GetPaymentRequest(userUUID, paymentUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Payment request details retrieved", payment)
}

// CancelPaymentRequest cancels a pending payment request
func (ac *AIController) CancelPaymentRequest(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	paymentID := c.Param("id")
	if paymentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Payment ID is required"})
		return
	}

	// Convert IDs to UUIDs
	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	paymentUUID, err := uuid.Parse(paymentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment ID"})
		return
	}

	// Cancel payment request
	err = ac.aiService.CancelPaymentRequest(userUUID, paymentUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Payment request cancelled", nil)
}
