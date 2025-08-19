package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zeusnotfound04/Tranza/models"
	"github.com/zeusnotfound04/Tranza/services"
	"github.com/zeusnotfound04/Tranza/utils"
)

type ClothingController struct {
	clothingService *services.ClothingService
}

func NewClothingController(clothingService *services.ClothingService) *ClothingController {
	return &ClothingController{
		clothingService: clothingService,
	}
}

// ProcessAIClothingOrder handles AI-powered clothing order requests
// This endpoint processes natural language requests and searches external e-commerce sites
func (cc *ClothingController) ProcessAIClothingOrder(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	var req models.AIClothingOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request format", err)
		return
	}

	response, err := cc.clothingService.ProcessAIClothingOrder(userID.(string), &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to process AI clothing order", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "AI clothing search completed", response)
}

// ConfirmAIClothingOrder handles order confirmation and payment processing
// This confirms the order and deducts payment from the Tranza wallet
func (cc *ClothingController) ConfirmAIClothingOrder(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	var req models.ConfirmAIClothingOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request format", err)
		return
	}

	response, err := cc.clothingService.ConfirmAIClothingOrder(userID.(string), &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to confirm order", err)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Order confirmed and payment processed", response)
}
