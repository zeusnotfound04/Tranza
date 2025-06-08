package controllers

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/zeusnotfound04/Tranza/models"
	"github.com/zeusnotfound04/Tranza/services"
)

type PaymentController struct {
	Razorpay *services.RazorpayService
}

func NewPaymentController(rz *services.RazorpayService) *PaymentController {
	return &PaymentController{Razorpay: rz}
}

func (pc *PaymentController) CreateOrder(ctx *gin.Context) {
	var req struct {
		Amount      float64           `json:"amount" binding:"required,gt=0"`
		Currency    string            `json:"currency"`
		Receipt     string            `json:"receipt"`
		Notes       map[string]string `json:"notes"`
		Description string            `json:"description"`
	}
	
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}
	
	// Generate receipt if not provided
	if req.Receipt == "" {
		req.Receipt = "order_" + uuid.New().String()[:8]
	}
	
	// Set default currency
	if req.Currency == "" {
		req.Currency = "INR"
	}
	
	// Add description to notes if provided
	if req.Description != "" {
		if req.Notes == nil {
			req.Notes = make(map[string]string)
		}
		req.Notes["description"] = req.Description
	}
	
	// Create context with timeout
	reqCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	order, err := pc.Razorpay.CreateOrder(reqCtx, req.Amount, req.Currency, req.Receipt, req.Notes)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create order",
			"details": err.Error(),
		})
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    order,
	})
}

func (pc *PaymentController) VerifyPayment(ctx *gin.Context) {
	var req models.PaymentVerificationRequest
	
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}
	
	if err := pc.Razorpay.VerifyPaymentSignature(req.OrderID, req.PaymentID, req.Signature); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Payment verification failed",
			"details": err.Error(),
		})
		return
	}
	
	// Fetch payment details for additional verification
	reqCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	payment, err := pc.Razorpay.FetchPayment(reqCtx, req.PaymentID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch payment details",
			"details": err.Error(),
		})
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Payment verified successfully",
		"data": gin.H{
			"payment_id": payment.ID,
			"order_id":   payment.OrderID,
			"amount":     payment.Amount,
			"status":     payment.Status,
			"method":     payment.Method,
		},
	})
}

func (pc *PaymentController) HandleWebhook(ctx *gin.Context) {
	// Get the signature from header
	signature := ctx.GetHeader("X-Razorpay-Signature")
	if signature == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Missing signature"})
		return
	}
	
	// Read the request body
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}
	
	// Verify webhook signature
	if err := pc.Razorpay.VerifyWebhookSignature(string(body), signature); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Signature verification failed"})
		return
	}
	
	// Parse webhook payload
	var webhook models.RazorpayWebhook
	if err := json.Unmarshal(body, &webhook); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})
		return
	}
	
	// Process webhook based on event type
	switch webhook.Event {
	case "payment.captured":
		// Handle successful payment
		pc.handlePaymentCaptured(webhook)
	case "payment.failed":
		// Handle failed payment
		pc.handlePaymentFailed(webhook)
	case "order.paid":
		// Handle order completion
		pc.handleOrderPaid(webhook)
	default:
		// Log unhandled events
		// You might want to use a proper logger here
		// log.Printf("Unhandled webhook event: %s", webhook.Event)
	}
	
	ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (pc *PaymentController) GetOrder(ctx *gin.Context) {
	orderID := ctx.Param("id")
	if orderID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Order ID is required"})
		return
	}
	
	reqCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	order, err := pc.Razorpay.FetchOrder(reqCtx, orderID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch order",
			"details": err.Error(),
		})
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    order,
	})
}

func (pc *PaymentController) GetPayment(ctx *gin.Context) {
	paymentID := ctx.Param("id")
	if paymentID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Payment ID is required"})
		return
	}
	
	reqCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	payment, err := pc.Razorpay.FetchPayment(reqCtx, paymentID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch payment",
			"details": err.Error(),
		})
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    payment,
	})
}

// Helper methods for webhook handling
func (pc *PaymentController) handlePaymentCaptured(webhook models.RazorpayWebhook) {
	// Implement your business logic for successful payments
	// e.g., update order status, send confirmation email, etc.
	
	// Example: Extract payment details from webhook payload
	if payment, ok := webhook.Payload["payment"].(map[string]interface{}); ok {
		paymentID := payment["id"].(string)
		amount := payment["amount"].(float64)
		// Process the captured payment
		log.Printf("Payment captured: ID=%s, Amount=%.2f", paymentID, amount/100)
	}
}

func (pc *PaymentController) handlePaymentFailed(webhook models.RazorpayWebhook) {
	// Implement your business logic for failed payments
	// e.g., update order status, notify customer, etc.
	
	if payment, ok := webhook.Payload["payment"].(map[string]interface{}); ok {
		paymentID := payment["id"].(string)
		// Process the failed payment
		log.Printf("Payment failed: ID=%s", paymentID)
	}
}

func (pc *PaymentController) handleOrderPaid(webhook models.RazorpayWebhook) {
	// Implement your business logic for completed orders
	// e.g., fulfill order, update inventory, etc.
	
	if order, ok := webhook.Payload["order"].(map[string]interface{}); ok {
		orderID := order["id"].(string)
		// Process the completed order
		log.Printf("Order completed: ID=%s", orderID)
	}
}
