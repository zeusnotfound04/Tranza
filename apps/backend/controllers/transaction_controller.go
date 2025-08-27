// controllers/transaction_controller.go
package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/zeusnotfound04/Tranza/models/dto"
	"github.com/zeusnotfound04/Tranza/services"
	"github.com/zeusnotfound04/Tranza/utils"
)

type TransactionController struct {
	transactionService *services.TransactionService
	paymentService     *services.PaymentService
}

func NewTransactionController(
	transactionService *services.TransactionService,
	paymentService *services.PaymentService,
) *TransactionController {
	return &TransactionController{
		transactionService: transactionService,
		paymentService:     paymentService,
	}
}

// GetTransactionHistory retrieves paginated transaction history with filters
func (c *TransactionController) GetTransactionHistory(ctx *gin.Context) {
	fmt.Printf("DEBUG: GetTransactionHistory endpoint called\n")
	fmt.Printf("DEBUG: Request URL: %s\n", ctx.Request.URL.String())
	fmt.Printf("DEBUG: Request method: %s\n", ctx.Request.Method)
	fmt.Printf("DEBUG: Request headers: %+v\n", ctx.Request.Header)

	userID, exists := ctx.Get("user_id") // From JWT middleware
	if !exists {
		fmt.Printf("DEBUG: UserID not found in context - authentication failed\n")
		utils.UnauthorizedResponse(ctx, "User not authenticated")
		return
	}

	// Convert to UUID
	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		fmt.Printf("DEBUG: UserID type assertion failed, got: %T\n", userID)
		utils.InternalServerErrorResponse(ctx, "Invalid user ID type", nil)
		return
	}

	fmt.Printf("DEBUG: UserID from context: %s\n", userUUID)

	// Parse query parameters
	var req dto.TransactionHistoryRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		fmt.Printf("DEBUG: Query parameter binding error: %v\n", err)
		utils.ValidationErrorResponse(ctx, map[string]string{
			"error": "Invalid query parameters",
		})
		return
	}

	fmt.Printf("DEBUG: Parsed request: %+v\n", req)

	// Set defaults
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 || req.Limit > 100 {
		req.Limit = 50
	}

	// Calculate offset
	offset := (req.Page - 1) * req.Limit

	// Get transaction history
	transactions, total, err := c.transactionService.GetTransactionHistory(
		userUUID.String(),
		req.Limit,
		offset,
		req.TransactionType,
	)
	if err != nil {
		utils.LogError(err, map[string]interface{}{
			"user_id": userUUID.String(),
			"action":  "get_transaction_history",
		})
		utils.InternalServerErrorResponse(ctx, "Failed to get transaction history", err)
		return
	}

	// Create paginated response
	utils.PaginatedSuccessResponse(
		ctx,
		"Transaction history retrieved successfully",
		transactions,
		req.Page,
		req.Limit,
		total,
	)
}

// GetTransaction retrieves a specific transaction by ID
func (c *TransactionController) GetTransaction(ctx *gin.Context) {
	userID, err := utils.GetUserIDStringFromContext(ctx)
	if err != nil {
		utils.UnauthorizedResponse(ctx, err.Error())
		return
	}

	transactionID := ctx.Param("id")
	if transactionID == "" {
		utils.BadRequestResponse(ctx, "Transaction ID is required", nil)
		return
	}

	// Get transaction details
	transaction, err := c.transactionService.GetTransactionByID(userID, transactionID)
	if err != nil {
		utils.LogError(err, map[string]interface{}{
			"user_id":        userID,
			"transaction_id": transactionID,
			"action":         "get_transaction_by_id",
		})
		utils.NotFoundResponse(ctx, "Transaction not found")
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Transaction details retrieved successfully", transaction)
}

// GetTransactionStats retrieves transaction statistics
func (c *TransactionController) GetTransactionStats(ctx *gin.Context) {
	userID, err := utils.GetUserIDStringFromContext(ctx)
	if err != nil {
		utils.UnauthorizedResponse(ctx, "User authentication failed")
		return
	}

	// Get transaction statistics
	stats, err := c.transactionService.GetTransactionStats(userID)
	if err != nil {
		utils.LogError(err, map[string]interface{}{
			"user_id": userID,
			"action":  "get_transaction_stats",
		})
		utils.InternalServerErrorResponse(ctx, "Failed to get transaction statistics", err)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Transaction statistics retrieved successfully", stats)
}

// GetTransactionAnalytics retrieves transaction analytics
func (c *TransactionController) GetTransactionAnalytics(ctx *gin.Context) {
	userID, err := utils.GetUserIDStringFromContext(ctx)
	if err != nil {
		utils.UnauthorizedResponse(ctx, "User authentication failed")
		return
	}

	// Parse days parameter
	daysStr := ctx.DefaultQuery("days", "30")
	days, err := strconv.Atoi(daysStr)
	if err != nil || days < 1 || days > 365 {
		days = 30
	}

	// Get payment analytics
	analytics, err := c.paymentService.GetPaymentAnalytics(userID, days)
	if err != nil {
		utils.LogError(err, map[string]interface{}{
			"user_id": userID,
			"days":    days,
			"action":  "get_payment_analytics",
		})
		utils.InternalServerErrorResponse(ctx, "Failed to get transaction analytics", err)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Transaction analytics retrieved successfully", analytics)
}

// SearchTransactions searches transactions with advanced filters
func (c *TransactionController) SearchTransactions(ctx *gin.Context) {
	userID := ctx.GetString("user_id")
	if userID == "" {
		utils.UnauthorizedResponse(ctx, "User not authenticated")
		return
	}

	var req dto.TransactionHistoryRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		utils.ValidationErrorResponse(ctx, map[string]string{
			"error": "Invalid search parameters",
		})
		return
	}

	// Set defaults
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 || req.Limit > 100 {
		req.Limit = 50
	}

	// Get filtered transactions
	transactions, total, err := c.transactionService.SearchTransactions(userID, &req)
	if err != nil {
		utils.LogError(err, map[string]interface{}{
			"user_id": userID,
			"filters": req,
			"action":  "search_transactions",
		})
		utils.InternalServerErrorResponse(ctx, "Failed to search transactions", err)
		return
	}

	utils.PaginatedSuccessResponse(
		ctx,
		"Transactions found successfully",
		transactions,
		req.Page,
		req.Limit,
		total,
	)
}

// GetTransactionsByType retrieves transactions filtered by type
func (c *TransactionController) GetTransactionsByType(ctx *gin.Context) {
	userID := ctx.GetString("user_id")
	if userID == "" {
		utils.UnauthorizedResponse(ctx, "User not authenticated")
		return
	}

	transactionType := ctx.Param("type")
	if transactionType == "" {
		utils.BadRequestResponse(ctx, "Transaction type is required", nil)
		return
	}

	// Validate transaction type
	validTypes := []string{
		utils.TransactionTypeLoadMoney,
		utils.TransactionTypeAIPayment,
		utils.TransactionTypeRefund,
		utils.TransactionTypeWithdrawal,
	}

	isValidType := false
	for _, validType := range validTypes {
		if transactionType == validType {
			isValidType = true
			break
		}
	}

	if !isValidType {
		utils.BadRequestResponse(ctx, "Invalid transaction type", nil)
		return
	}

	// Parse pagination parameters
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "50"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 50
	}

	offset := (page - 1) * limit

	// Get transactions by type
	transactions, total, err := c.transactionService.GetTransactionHistory(
		userID,
		limit,
		offset,
		transactionType,
	)
	if err != nil {
		utils.LogError(err, map[string]interface{}{
			"user_id":          userID,
			"transaction_type": transactionType,
			"action":           "get_transactions_by_type",
		})
		utils.InternalServerErrorResponse(ctx, "Failed to get transactions", err)
		return
	}

	utils.PaginatedSuccessResponse(
		ctx,
		fmt.Sprintf("%s transactions retrieved successfully", transactionType),
		transactions,
		page,
		limit,
		total,
	)
}

// GetTransactionReceipt generates and returns transaction receipt
func (c *TransactionController) GetTransactionReceipt(ctx *gin.Context) {
	userID := ctx.GetString("user_id")
	if userID == "" {
		utils.UnauthorizedResponse(ctx, "User not authenticated")
		return
	}

	transactionID := ctx.Param("id")
	if transactionID == "" {
		utils.BadRequestResponse(ctx, "Transaction ID is required", nil)
		return
	}

	// Get transaction details
	transaction, err := c.transactionService.GetTransactionByID(userID, transactionID)
	if err != nil {
		utils.LogError(err, map[string]interface{}{
			"user_id":        userID,
			"transaction_id": transactionID,
			"action":         "get_transaction_receipt",
		})
		utils.NotFoundResponse(ctx, "Transaction not found")
		return
	}

	// Generate receipt
	receipt, err := c.transactionService.GenerateTransactionReceipt(transaction)
	if err != nil {
		utils.LogError(err, map[string]interface{}{
			"user_id":        userID,
			"transaction_id": transactionID,
			"action":         "generate_receipt",
		})
		utils.InternalServerErrorResponse(ctx, "Failed to generate receipt", err)
		return
	}

	// Log receipt generation
	utils.LogInfo("Transaction receipt generated", map[string]interface{}{
		"user_id":        userID,
		"transaction_id": transactionID,
	})

	utils.SuccessResponse(ctx, http.StatusOK, "Transaction receipt generated successfully", receipt)
}

// ExportTransactions exports transaction history as CSV
func (c *TransactionController) ExportTransactions(ctx *gin.Context) {
	userID := ctx.GetString("user_id")
	if userID == "" {
		utils.UnauthorizedResponse(ctx, "User not authenticated")
		return
	}

	// Parse date range parameters
	startDate := ctx.Query("start_date")
	endDate := ctx.Query("end_date")
	format := ctx.DefaultQuery("format", "csv")

	// Validate format
	if format != "csv" && format != "pdf" {
		utils.BadRequestResponse(ctx, "Invalid export format. Supported formats: csv, pdf", nil)
		return
	}

	// Export transactions
	exportData, err := c.transactionService.ExportTransactions(userID, startDate, endDate, format)
	if err != nil {
		utils.LogError(err, map[string]interface{}{
			"user_id":    userID,
			"start_date": startDate,
			"end_date":   endDate,
			"format":     format,
			"action":     "export_transactions",
		})
		utils.InternalServerErrorResponse(ctx, "Failed to export transactions", err)
		return
	}

	// Set appropriate headers for file download
	filename := fmt.Sprintf("transactions_%s.%s", time.Now().Format("2006-01-02"), format)

	if format == "csv" {
		ctx.Header("Content-Type", "text/csv")
	} else {
		ctx.Header("Content-Type", "application/pdf")
	}

	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

	// Log export
	utils.LogInfo("Transactions exported", map[string]interface{}{
		"user_id":  userID,
		"format":   format,
		"filename": filename,
	})

	ctx.Data(http.StatusOK, ctx.GetHeader("Content-Type"), exportData)
}

// GetMonthlyTransactionSummary retrieves monthly transaction summary
func (c *TransactionController) GetMonthlyTransactionSummary(ctx *gin.Context) {
	userID := ctx.GetString("user_id")
	if userID == "" {
		utils.UnauthorizedResponse(ctx, "User not authenticated")
		return
	}

	// Parse month and year parameters
	monthStr := ctx.DefaultQuery("month", strconv.Itoa(int(time.Now().Month())))
	yearStr := ctx.DefaultQuery("year", strconv.Itoa(time.Now().Year()))

	month, err := strconv.Atoi(monthStr)
	if err != nil || month < 1 || month > 12 {
		utils.BadRequestResponse(ctx, "Invalid month parameter", nil)
		return
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil || year < 2020 || year > time.Now().Year()+1 {
		utils.BadRequestResponse(ctx, "Invalid year parameter", nil)
		return
	}

	// Get monthly summary
	summary, err := c.transactionService.GetMonthlyTransactionSummary(userID, month, year)
	if err != nil {
		utils.LogError(err, map[string]interface{}{
			"user_id": userID,
			"month":   month,
			"year":    year,
			"action":  "get_monthly_summary",
		})
		utils.InternalServerErrorResponse(ctx, "Failed to get monthly summary", err)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Monthly transaction summary retrieved successfully", summary)
}

// GetDailyTransactionSummary retrieves daily transaction summary
func (c *TransactionController) GetDailyTransactionSummary(ctx *gin.Context) {
	userID := ctx.GetString("user_id")
	if userID == "" {
		utils.UnauthorizedResponse(ctx, "User not authenticated")
		return
	}

	// Parse date parameter
	dateStr := ctx.DefaultQuery("date", time.Now().Format("2006-01-02"))

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		utils.BadRequestResponse(ctx, "Invalid date format. Use YYYY-MM-DD", nil)
		return
	}

	// Get daily summary
	summary, err := c.transactionService.GetDailyTransactionSummary(userID, date)
	if err != nil {
		utils.LogError(err, map[string]interface{}{
			"user_id": userID,
			"date":    dateStr,
			"action":  "get_daily_summary",
		})
		utils.InternalServerErrorResponse(ctx, "Failed to get daily summary", err)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Daily transaction summary retrieved successfully", summary)
}

// GetTransactionTrends retrieves transaction trends and patterns
func (c *TransactionController) GetTransactionTrends(ctx *gin.Context) {
	userID := ctx.GetString("user_id")
	if userID == "" {
		utils.UnauthorizedResponse(ctx, "User not authenticated")
		return
	}

	// Parse period parameter
	period := ctx.DefaultQuery("period", "month") // day, week, month, year

	validPeriods := []string{"day", "week", "month", "year"}
	isValidPeriod := false
	for _, validPeriod := range validPeriods {
		if period == validPeriod {
			isValidPeriod = true
			break
		}
	}

	if !isValidPeriod {
		utils.BadRequestResponse(ctx, "Invalid period. Valid options: day, week, month, year", nil)
		return
	}

	// Get transaction trends
	trends, err := c.transactionService.GetTransactionTrends(userID, period)
	if err != nil {
		utils.LogError(err, map[string]interface{}{
			"user_id": userID,
			"period":  period,
			"action":  "get_transaction_trends",
		})
		utils.InternalServerErrorResponse(ctx, "Failed to get transaction trends", err)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Transaction trends retrieved successfully", trends)
}

// ValidateTransaction validates a transaction (for admin use)
func (c *TransactionController) ValidateTransaction(ctx *gin.Context) {
	// This would typically require admin privileges
	userRole := ctx.GetString("user_role")
	if userRole != "admin" {
		utils.ForbiddenResponse(ctx, "Insufficient privileges")
		return
	}

	transactionID := ctx.Param("id")
	if transactionID == "" {
		utils.BadRequestResponse(ctx, "Transaction ID is required", nil)
		return
	}

	// Validate transaction
	result, err := c.transactionService.ValidateTransaction(transactionID)
	if err != nil {
		utils.LogError(err, map[string]interface{}{
			"transaction_id": transactionID,
			"action":         "validate_transaction",
		})
		utils.InternalServerErrorResponse(ctx, "Failed to validate transaction", err)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Transaction validation completed", result)
}

// RetryFailedTransaction retries a failed transaction
func (c *TransactionController) RetryFailedTransaction(ctx *gin.Context) {
	userID := ctx.GetString("user_id")
	if userID == "" {
		utils.UnauthorizedResponse(ctx, "User not authenticated")
		return
	}

	transactionID := ctx.Param("id")
	if transactionID == "" {
		utils.BadRequestResponse(ctx, "Transaction ID is required", nil)
		return
	}

	// Retry failed transaction
	result, err := c.transactionService.RetryFailedTransaction(userID, transactionID)
	if err != nil {
		utils.LogError(err, map[string]interface{}{
			"user_id":        userID,
			"transaction_id": transactionID,
			"action":         "retry_failed_transaction",
		})
		utils.BadRequestResponse(ctx, utils.GetUserFriendlyError(err), err)
		return
	}

	utils.LogInfo("Transaction retry initiated", map[string]interface{}{
		"user_id":            userID,
		"transaction_id":     transactionID,
		"new_transaction_id": result.ID,
	})

	utils.SuccessResponse(ctx, http.StatusOK, "Transaction retry initiated successfully", result)
}
