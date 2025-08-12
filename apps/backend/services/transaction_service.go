package services

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/zeusnotfound04/Tranza/models"
	"github.com/zeusnotfound04/Tranza/models/dto"
	"github.com/zeusnotfound04/Tranza/repositories"
	"github.com/zeusnotfound04/Tranza/utils"
)

type TransactionService struct {
	transactionRepo *repositories.TransactionRepository
	walletRepo      *repositories.WalletRepository
	paymentService  *PaymentService
}

func NewTransactionService(
	transactionRepo *repositories.TransactionRepository,
	walletRepo *repositories.WalletRepository,
	paymentService *PaymentService,
) *TransactionService {
	return &TransactionService{
		transactionRepo: transactionRepo,
		walletRepo:      walletRepo,
		paymentService:  paymentService,
	}
}

// GetTransactionHistory retrieves transaction history with pagination and filters
func (s *TransactionService) GetTransactionHistory(userID string, limit, offset int, transactionType string) ([]*dto.TransactionResponse, int64, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, 0, errors.New("invalid user ID")
	}

	transactions, total, err := s.transactionRepo.GetByUserIDWithPagination(uid, limit, offset, transactionType)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get transaction history: %w", err)
	}

	// Convert to response DTOs
	response := make([]*dto.TransactionResponse, len(transactions))
	for i, txn := range transactions {
		response[i] = s.convertToTransactionResponse(txn)
	}

	return response, total, nil
}

// GetTransactionByID retrieves a specific transaction by ID
func (s *TransactionService) GetTransactionByID(userID, transactionID string) (*dto.TransactionResponse, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	tid, err := uuid.Parse(transactionID)
	if err != nil {
		return nil, errors.New("invalid transaction ID")
	}

	transaction, err := s.transactionRepo.GetByIDAndUserID(tid, uid)
	if err != nil {
		return nil, fmt.Errorf("transaction not found: %w", err)
	}

	return s.convertToTransactionResponse(transaction), nil
}

// GetTransactionStats retrieves comprehensive transaction statistics
func (s *TransactionService) GetTransactionStats(userID string) (*dto.TransactionStatsResponse, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	stats, err := s.transactionRepo.GetTransactionStats(uid)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction stats: %w", err)
	}

	// Get most used payment method
	mostUsedMethod, err := s.getMostUsedPaymentMethod(uid)
	if err != nil {
		mostUsedMethod = "N/A"
	}

	return &dto.TransactionStatsResponse{
		TotalTransactions:     stats.TotalTransactions,
		TotalAmount:           stats.TotalAmount,
		TodayTransactions:     stats.TodayTransactions,
		TodayAmount:           stats.TodayAmount,
		MonthTransactions:     s.getMonthTransactions(uid),
		MonthAmount:           s.getMonthAmount(uid),
		AITransactions:        stats.AITransactions,
		AIAmount:              stats.AIAmount,
		LastTransactionDate:   stats.LastTransactionDate.Format("2006-01-02 15:04:05"),
		MostUsedPaymentMethod: mostUsedMethod,
	}, nil
}

// SearchTransactions searches transactions with advanced filters
func (s *TransactionService) SearchTransactions(userID string, req *dto.TransactionHistoryRequest) ([]*dto.TransactionResponse, int64, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, 0, errors.New("invalid user ID")
	}

	// Parse date filters if provided
	var startDate, endDate *time.Time
	if req.StartDate != "" {
		if parsed, err := time.Parse("2006-01-02", req.StartDate); err == nil {
			startDate = &parsed
		}
	}
	if req.EndDate != "" {
		if parsed, err := time.Parse("2006-01-02", req.EndDate); err == nil {
			endDate = &parsed
		}
	}

	// Parse amount filters
	var minAmount, maxAmount *decimal.Decimal
	if req.MinAmount != "" {
		if parsed, err := decimal.NewFromString(req.MinAmount); err == nil {
			minAmount = &parsed
		}
	}
	if req.MaxAmount != "" {
		if parsed, err := decimal.NewFromString(req.MaxAmount); err == nil {
			maxAmount = &parsed
		}
	}

	offset := (req.Page - 1) * req.Limit

	// Get filtered transactions
	transactions, total, err := s.transactionRepo.GetTransactionsWithFilters(
		uid, req.Limit, offset, req.TransactionType, req.Status,
		startDate, endDate, minAmount, maxAmount,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search transactions: %w", err)
	}

	// Convert to response DTOs
	response := make([]*dto.TransactionResponse, len(transactions))
	for i, txn := range transactions {
		response[i] = s.convertToTransactionResponse(txn)
	}

	return response, total, nil
}

// GenerateTransactionReceipt generates a receipt for a transaction
func (s *TransactionService) GenerateTransactionReceipt(transaction *dto.TransactionResponse) (*TransactionReceipt, error) {
	receipt := &TransactionReceipt{
		TransactionID: transaction.ID,
		ReferenceID:   transaction.ReferenceID,
		Date:          transaction.CreatedAt,
		Type:          transaction.Type,
		Amount:        transaction.Amount,
		Currency:      transaction.Currency,
		Status:        transaction.Status,
		Description:   transaction.Description,
		PaymentMethod: transaction.PaymentMethod,
		MerchantName:  transaction.MerchantName,
		BalanceAfter:  transaction.BalanceAfter,
		GeneratedAt:   time.Now().Format("2006-01-02 15:04:05"),
	}

	return receipt, nil
}

// ExportTransactions exports transaction history as CSV or PDF
func (s *TransactionService) ExportTransactions(userID, startDate, endDate, format string) ([]byte, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	// Parse date range
	var start, end *time.Time
	if startDate != "" {
		if parsed, err := time.Parse("2006-01-02", startDate); err == nil {
			start = &parsed
		}
	}
	if endDate != "" {
		if parsed, err := time.Parse("2006-01-02", endDate); err == nil {
			end = &parsed
		}
	}

	// Get transactions for export
	transactions, err := s.getTransactionsForExport(uid, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions for export: %w", err)
	}

	switch format {
	case "csv":
		return s.exportTransactionsAsCSV(transactions)
	case "pdf":
		return s.exportTransactionsAsPDF(transactions)
	default:
		return nil, errors.New("unsupported export format")
	}
}

// GetMonthlyTransactionSummary retrieves monthly transaction summary
func (s *TransactionService) GetMonthlyTransactionSummary(userID string, month, year int) (*MonthlyTransactionSummary, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	// Get wallet
	wallet, err := s.walletRepo.GetByUserID(uid)
	if err != nil {
		return nil, errors.New("wallet not found")
	}

	// Calculate date range for the month
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, -1).Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	// Get transaction summary
	summary, err := s.transactionRepo.GetTransactionSummaryByDateRange(wallet.ID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get monthly summary: %w", err)
	}

	// Get transactions by type
	typeBreakdown, err := s.getTransactionsByTypeInRange(wallet.ID, startDate, endDate)
	if err != nil {
		typeBreakdown = make(map[string]TransactionTypeBreakdown)
	}

	return &MonthlyTransactionSummary{
		Month:             month,
		Year:              year,
		TotalTransactions: summary.TotalTransactions,
		TotalInflow:       summary.TotalInflow,
		TotalOutflow:      summary.TotalOutflow,
		NetFlow:           summary.NetFlow,
		AverageAmount:     summary.AverageAmount,
		TypeBreakdown:     typeBreakdown,
		StartingBalance:   s.getBalanceAtDate(wallet.ID, startDate),
		EndingBalance:     s.getBalanceAtDate(wallet.ID, endDate),
	}, nil
}

// GetDailyTransactionSummary retrieves daily transaction summary
func (s *TransactionService) GetDailyTransactionSummary(userID string, date time.Time) (*DailyTransactionSummary, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	// Get wallet
	wallet, err := s.walletRepo.GetByUserID(uid)
	if err != nil {
		return nil, errors.New("wallet not found")
	}

	// Calculate day range
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour).Add(-time.Nanosecond)

	// Get transaction summary
	summary, err := s.transactionRepo.GetTransactionSummaryByDateRange(wallet.ID, startOfDay, endOfDay)
	if err != nil {
		return nil, fmt.Errorf("failed to get daily summary: %w", err)
	}

	// Get hourly breakdown
	hourlyBreakdown, err := s.getHourlyTransactionBreakdown(wallet.ID, date)
	if err != nil {
		hourlyBreakdown = make([]HourlyTransaction, 0)
	}

	return &DailyTransactionSummary{
		Date:              date.Format("2006-01-02"),
		TotalTransactions: summary.TotalTransactions,
		TotalInflow:       summary.TotalInflow,
		TotalOutflow:      summary.TotalOutflow,
		NetFlow:           summary.NetFlow,
		AverageAmount:     summary.AverageAmount,
		HourlyBreakdown:   hourlyBreakdown,
		StartingBalance:   s.getBalanceAtDate(wallet.ID, startOfDay),
		EndingBalance:     s.getBalanceAtDate(wallet.ID, endOfDay),
	}, nil
}

// GetTransactionTrends retrieves transaction trends and patterns
func (s *TransactionService) GetTransactionTrends(userID, period string) (*TransactionTrends, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	// Get wallet
	wallet, err := s.walletRepo.GetByUserID(uid)
	if err != nil {
		return nil, errors.New("wallet not found")
	}

	var trends *TransactionTrends

	switch period {
	case "day":
		trends, err = s.getDailyTrends(wallet.ID)
	case "week":
		trends, err = s.getWeeklyTrends(wallet.ID)
	case "month":
		trends, err = s.getMonthlyTrends(wallet.ID)
	case "year":
		trends, err = s.getYearlyTrends(wallet.ID)
	default:
		return nil, errors.New("invalid period")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get transaction trends: %w", err)
	}

	return trends, nil
}

// ValidateTransaction validates a transaction (admin function)
func (s *TransactionService) ValidateTransaction(transactionID string) (*TransactionValidationResult, error) {
	tid, err := uuid.Parse(transactionID)
	if err != nil {
		return nil, errors.New("invalid transaction ID")
	}

	transaction, err := s.transactionRepo.GetByID(tid)
	if err != nil {
		return nil, fmt.Errorf("transaction not found: %w", err)
	}

	result := &TransactionValidationResult{
		TransactionID: transactionID,
		IsValid:       true,
		Issues:        make([]string, 0),
		Checks:        make(map[string]bool),
	}

	// Perform various validation checks
	s.validateTransactionAmount(transaction, result)
	s.validateTransactionStatus(transaction, result)
	s.validateTransactionBalance(transaction, result)
	s.validateRazorpayData(transaction, result)

	// Overall validation result
	result.IsValid = len(result.Issues) == 0

	return result, nil
}

// RetryFailedTransaction retries a failed transaction
func (s *TransactionService) RetryFailedTransaction(userID, transactionID string) (*dto.TransactionResponse, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	tid, err := uuid.Parse(transactionID)
	if err != nil {
		return nil, errors.New("invalid transaction ID")
	}

	// Get the original transaction
	originalTxn, err := s.transactionRepo.GetByIDAndUserID(tid, uid)
	if err != nil {
		return nil, errors.New("transaction not found")
	}

	// Validate that transaction can be retried
	if originalTxn.Status != utils.TransactionStatusFailed {
		return nil, errors.New("only failed transactions can be retried")
	}

	if originalTxn.Type != "load_money" {
		return nil, errors.New("only load money transactions can be retried")
	}

	// Create new order for retry
	response, err := s.paymentService.CreateLoadMoneyOrder(userID, originalTxn.Amount)
	if err != nil {
		return nil, fmt.Errorf("failed to create retry order: %w", err)
	}

	// Get the new transaction
	newTxn, err := s.transactionRepo.GetByID(uuid.MustParse(response.TransactionID))
	if err != nil {
		return nil, fmt.Errorf("failed to get new transaction: %w", err)
	}

	return s.convertToTransactionResponse(newTxn), nil
}

// Helper methods

func (s *TransactionService) convertToTransactionResponse(txn *models.Transaction) *dto.TransactionResponse {
	return &dto.TransactionResponse{
		ID:            txn.ID.String(),
		Type:          txn.Type,
		Amount:        decimal.NewFromFloat(0).Add(txn.Amount), // Convert to ensure it's decimal.Decimal
		BalanceAfter:  txn.BalanceAfter,
		Currency:      txn.Currency,
		Description:   txn.Description,
		Status:        string(txn.Status),
		PaymentMethod: txn.PaymentMethod,
		AIAgentID:     txn.AIAgentID,
		MerchantName:  txn.MerchantName,
		MerchantUPIID: txn.MerchantUPIID,
		ReferenceID:   txn.ReferenceID,
		FailureReason: txn.FailureReason,
		CreatedAt:     txn.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:     txn.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func (s *TransactionService) getMostUsedPaymentMethod(userID uuid.UUID) (string, error) {
	// Implementation to get most used payment method
	return "UPI", nil // Placeholder
}

func (s *TransactionService) getMonthTransactions(userID uuid.UUID) int64 {
	// Implementation to get current month transactions count
	return 0 // Placeholder
}

func (s *TransactionService) getMonthAmount(userID uuid.UUID) decimal.Decimal {
	// Implementation to get current month amount
	return decimal.Zero // Placeholder
}

func (s *TransactionService) getTransactionsForExport(userID uuid.UUID, startDate, endDate *time.Time) ([]*models.Transaction, error) {
	// Get wallet
	wallet, err := s.walletRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	if startDate != nil && endDate != nil {
		return s.transactionRepo.GetByWalletIDWithDateRange(wallet.ID, *startDate, *endDate)
	}

	return s.transactionRepo.GetByWalletID(wallet.ID)
}

func (s *TransactionService) exportTransactionsAsCSV(transactions []*models.Transaction) ([]byte, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Write header
	header := []string{
		"Transaction ID", "Date", "Type", "Amount", "Currency",
		"Status", "Description", "Payment Method", "Merchant",
		"Balance After", "Reference ID",
	}
	if err := writer.Write(header); err != nil {
		return nil, err
	}

	// Write data
	for _, txn := range transactions {
		record := []string{
			txn.ID.String(),
			txn.CreatedAt.Format("2006-01-02 15:04:05"),
			txn.Type,
			txn.Amount.String(),
			txn.Currency,
			string(txn.Status),
			txn.Description,
			txn.PaymentMethod,
			txn.MerchantName,
			txn.BalanceAfter.String(),
			txn.ReferenceID,
		}
		if err := writer.Write(record); err != nil {
			return nil, err
		}
	}

	writer.Flush()
	return buf.Bytes(), nil
}

func (s *TransactionService) exportTransactionsAsPDF(transactions []*models.Transaction) ([]byte, error) {
	// Placeholder for PDF generation
	// You would use a library like gofpdf or wkhtmltopdf
	return []byte("PDF export not implemented yet"), nil
}

func (s *TransactionService) getTransactionsByTypeInRange(walletID uuid.UUID, startDate, endDate time.Time) (map[string]TransactionTypeBreakdown, error) {
	// Implementation to get transactions by type in date range
	return make(map[string]TransactionTypeBreakdown), nil
}

func (s *TransactionService) getBalanceAtDate(walletID uuid.UUID, date time.Time) decimal.Decimal {
	// Implementation to get balance at specific date
	return decimal.Zero
}

func (s *TransactionService) getHourlyTransactionBreakdown(walletID uuid.UUID, date time.Time) ([]HourlyTransaction, error) {
	// Implementation to get hourly breakdown
	return make([]HourlyTransaction, 0), nil
}

func (s *TransactionService) getDailyTrends(walletID uuid.UUID) (*TransactionTrends, error) {
	// Implementation for daily trends
	return &TransactionTrends{
		Period:     "day",
		DataPoints: make([]TrendDataPoint, 0),
	}, nil
}

func (s *TransactionService) getWeeklyTrends(walletID uuid.UUID) (*TransactionTrends, error) {
	// Implementation for weekly trends
	return &TransactionTrends{
		Period:     "week",
		DataPoints: make([]TrendDataPoint, 0),
	}, nil
}

func (s *TransactionService) getMonthlyTrends(walletID uuid.UUID) (*TransactionTrends, error) {
	// Implementation for monthly trends
	return &TransactionTrends{
		Period:     "month",
		DataPoints: make([]TrendDataPoint, 0),
	}, nil
}

func (s *TransactionService) getYearlyTrends(walletID uuid.UUID) (*TransactionTrends, error) {
	// Implementation for yearly trends
	return &TransactionTrends{
		Period:     "year",
		DataPoints: make([]TrendDataPoint, 0),
	}, nil
}

// Validation helper methods
func (s *TransactionService) validateTransactionAmount(txn *models.Transaction, result *TransactionValidationResult) {
	if txn.Amount.LessThanOrEqual(decimal.Zero) {
		result.Issues = append(result.Issues, "Invalid transaction amount")
		result.Checks["amount_valid"] = false
	} else {
		result.Checks["amount_valid"] = true
	}
}

func (s *TransactionService) validateTransactionStatus(txn *models.Transaction, result *TransactionValidationResult) {
	validStatuses := []string{
		utils.TransactionStatusPending,
		utils.TransactionStatusSuccess,
		utils.TransactionStatusFailed,
		utils.TransactionStatusCancelled,
	}

	isValid := false
	txnStatus := string(txn.Status)
	for _, status := range validStatuses {
		if txnStatus == status {
			isValid = true
			break
		}
	}

	if !isValid {
		result.Issues = append(result.Issues, "Invalid transaction status")
		result.Checks["status_valid"] = false
	} else {
		result.Checks["status_valid"] = true
	}
}

func (s *TransactionService) validateTransactionBalance(txn *models.Transaction, result *TransactionValidationResult) {
	if txn.BalanceAfter.LessThan(decimal.Zero) {
		result.Issues = append(result.Issues, "Invalid balance after transaction")
		result.Checks["balance_valid"] = false
	} else {
		result.Checks["balance_valid"] = true
	}
}

func (s *TransactionService) validateRazorpayData(txn *models.Transaction, result *TransactionValidationResult) {
	if txn.Type == "load_money" {
		if txn.RazorpayOrderID == "" {
			result.Issues = append(result.Issues, "Missing Razorpay order ID")
			result.Checks["razorpay_data_valid"] = false
		} else if string(txn.Status) == "success" && txn.RazorpayPaymentID == "" {
			result.Issues = append(result.Issues, "Missing Razorpay payment ID for successful transaction")
			result.Checks["razorpay_data_valid"] = false
		} else {
			result.Checks["razorpay_data_valid"] = true
		}
	} else {
		result.Checks["razorpay_data_valid"] = true
	}
}

// Data structures for responses

type TransactionReceipt struct {
	TransactionID string          `json:"transaction_id"`
	ReferenceID   string          `json:"reference_id"`
	Date          string          `json:"date"`
	Type          string          `json:"type"`
	Amount        decimal.Decimal `json:"amount"`
	Currency      string          `json:"currency"`
	Status        string          `json:"status"`
	Description   string          `json:"description"`
	PaymentMethod string          `json:"payment_method"`
	MerchantName  string          `json:"merchant_name"`
	BalanceAfter  decimal.Decimal `json:"balance_after"`
	GeneratedAt   string          `json:"generated_at"`
}

type MonthlyTransactionSummary struct {
	Month             int                                 `json:"month"`
	Year              int                                 `json:"year"`
	TotalTransactions int64                               `json:"total_transactions"`
	TotalInflow       decimal.Decimal                     `json:"total_inflow"`
	TotalOutflow      decimal.Decimal                     `json:"total_outflow"`
	NetFlow           decimal.Decimal                     `json:"net_flow"`
	AverageAmount     decimal.Decimal                     `json:"average_amount"`
	TypeBreakdown     map[string]TransactionTypeBreakdown `json:"type_breakdown"`
	StartingBalance   decimal.Decimal                     `json:"starting_balance"`
	EndingBalance     decimal.Decimal                     `json:"ending_balance"`
}

type DailyTransactionSummary struct {
	Date              string              `json:"date"`
	TotalTransactions int64               `json:"total_transactions"`
	TotalInflow       decimal.Decimal     `json:"total_inflow"`
	TotalOutflow      decimal.Decimal     `json:"total_outflow"`
	NetFlow           decimal.Decimal     `json:"net_flow"`
	AverageAmount     decimal.Decimal     `json:"average_amount"`
	HourlyBreakdown   []HourlyTransaction `json:"hourly_breakdown"`
	StartingBalance   decimal.Decimal     `json:"starting_balance"`
	EndingBalance     decimal.Decimal     `json:"ending_balance"`
}

type TransactionTypeBreakdown struct {
	Count  int64           `json:"count"`
	Amount decimal.Decimal `json:"amount"`
}

type HourlyTransaction struct {
	Hour   int             `json:"hour"`
	Count  int64           `json:"count"`
	Amount decimal.Decimal `json:"amount"`
}

type TransactionTrends struct {
	Period     string           `json:"period"`
	DataPoints []TrendDataPoint `json:"data_points"`
}

type TrendDataPoint struct {
	Date   string          `json:"date"`
	Amount decimal.Decimal `json:"amount"`
	Count  int64           `json:"count"`
}

type TransactionValidationResult struct {
	TransactionID string          `json:"transaction_id"`
	IsValid       bool            `json:"is_valid"`
	Issues        []string        `json:"issues"`
	Checks        map[string]bool `json:"checks"`
}
