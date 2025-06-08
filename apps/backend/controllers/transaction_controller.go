package controllers

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zeusnotfound04/Tranza/services"
)



type TransactionController struct {
	PaymentService *services.PaymentService
}


func NewTransactionController(ps *services.PaymentService) *TransactionController {
	return &TransactionController{PaymentService: ps}
}


type MakePaymentRequest struct {
	CardID   uint  `json:"card_id" binding:"requried"`
	Amount   float64 `json:"amount" binding:"requried"`
	Description string `json:"description" binding:"required"`
}


func (c *TransactionController) MakePayment(ctx *gin.Context) {
	var req MakePaymentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error" : "invalid input:" + err.Error()})
	}

	userIDStr := ctx.GetHeader("X-User-ID")
	userID , err := strconv.Atoi(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error" : "Invalid user"})
	}

	timeCtx , cancel := context.WithTimeout(ctx.Request.Context(), 5*time.Second)
	defer cancel()
	txn , err := c.PaymentService.MakePayment(timeCtx , uint(userID) , req.CardID , req.Amount , req.Description)

	if err != nil {
		if errors.Is(err , services.ErrUnauthorizedCardAccess) {
			ctx.JSON(http.StatusForbidden , gin.H{"error" : "Unauthorized card access"})
		} else if errors.Is(err , services.ErrPaymentLimitExeceed) {
			ctx.JSON(http.StatusBadRequest , gin.H{"error" :"Limit exceeded"})
		} else {
			ctx.JSON(http.StatusInternalServerError , gin.H{"error" : "transaction failed"})
		}
		return

	}

	ctx.JSON(http.StatusCreated , txn)

}



func (c *TransactionController) GetUserTransactions(ctx *gin.Context) {
	userIDStr := ctx.GetHeader("X-User-ID")
	userID , err := strconv.Atoi(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized , gin.H{"error" : "Invalid user"})
		return
	}
	timeCtx , cancel := context.WithTimeout(ctx.Request.Context() , 3*time.Second)
	defer cancel()

	txns , err := c.PaymentService.TxnRepo.GetUserTransactions(timeCtx , uint(userID))

	if err != nil {
		ctx.JSON(http.StatusInternalServerError  , gin.H{"error" : "could not retrieve transactions"})
		return
	}
	
	ctx.JSON(http.StatusOK , txns)

}
