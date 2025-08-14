package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zeusnotfound04/Tranza/models/dto"
	"github.com/zeusnotfound04/Tranza/services"
	"github.com/zeusnotfound04/Tranza/utils"
)

type WalletHandler struct {
	walletService *services.WalletService
}

func NewWalletHandler(walletService *services.WalletService) *WalletHandler {
	return &WalletHandler{
		walletService: walletService,
	}
}

// Get wallet balance and details
func (h *WalletHandler) GetWallet(c *gin.Context) {
	fmt.Printf("DEBUG: GetWallet endpoint called\n")
	fmt.Printf("DEBUG: Request headers: %+v\n", c.Request.Header)
	fmt.Printf("DEBUG: Request URL: %s\n", c.Request.URL.Path)
	fmt.Printf("DEBUG: Request method: %s\n", c.Request.Method)
	
	userID := c.GetString("userID") // From JWT middleware
	fmt.Printf("DEBUG: UserID from JWT middleware: %s\n", userID)
	
	if userID == "" {
		fmt.Printf("DEBUG: UserID is empty - JWT middleware may not have set it\n")
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	fmt.Printf("DEBUG: Calling wallet service with userID: %s\n", userID)
	wallet, err := h.walletService.GetWalletByUserID(userID)
	if err != nil {
		fmt.Printf("DEBUG: Wallet service error: %v\n", err)
		utils.ErrorResponse(c, http.StatusNotFound, "Wallet not found", err)
		return
	}

	fmt.Printf("DEBUG: Wallet found: %+v\n", wallet)
	response := &dto.WalletResponse{
		ID:                    wallet.ID.String(),
		Balance:               wallet.Balance,
		Currency:              wallet.Currency,
		Status:                wallet.Status,
		AIAccessEnabled:       wallet.AIAccessEnabled,
		AIDailyLimit:          wallet.AIDailyLimit,
		AIPerTransactionLimit: wallet.AIPerTransactionLimit,
	}

	fmt.Printf("DEBUG: Sending successful response\n")
	utils.SuccessResponse(c, http.StatusOK, "Wallet details retrieved", response)
}

// Create order for loading money
func (h *WalletHandler) CreateLoadMoneyOrder(c *gin.Context) {
	userID := c.GetString("userID")
	
	// Debug log
	fmt.Printf("DEBUG: CreateLoadMoneyOrder called for userID: %s\n", userID)
	
	var req dto.LoadMoneyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Printf("DEBUG: JSON binding error: %v\n", err)
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request", err)
		return
	}

	// Debug log
	fmt.Printf("DEBUG: Amount received: %s\n", req.Amount.String())

	// Validate amount
	if err := utils.ValidateAmount(req.Amount); err != nil {
		fmt.Printf("DEBUG: Amount validation error: %v\n", err)
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid amount", err)
		return
	}

	order, err := h.walletService.CreateLoadMoneyOrder(userID, req.Amount)
	if err != nil {
		fmt.Printf("DEBUG: CreateLoadMoneyOrder service error: %v\n", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create order", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Order created successfully", order)
}

// Verify Razorpay payment and credit wallet
func (h *WalletHandler) VerifyPayment(c *gin.Context) {
	userID := c.GetString("userID")

	var req dto.VerifyPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request", err)
		return
	}

	result, err := h.walletService.VerifyAndCreditWallet(
		userID,
		req.RazorpayPaymentID,
		req.RazorpayOrderID,
		req.RazorpaySignature,
	)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Payment verification failed", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Payment verified and wallet credited", result)
}

// Update wallet settings
func (h *WalletHandler) UpdateWalletSettings(c *gin.Context) {
	userID := c.GetString("userID")

	var req dto.UpdateWalletSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request", err)
		return
	}

	err := h.walletService.UpdateWalletSettings(userID, &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update settings", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Wallet settings updated successfully", nil)
}
