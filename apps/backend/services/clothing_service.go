package services

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/zeusnotfound04/Tranza/models"
	"github.com/zeusnotfound04/Tranza/repositories"
)

type ClothingService struct {
	addressRepo     *repositories.AddressRepository
	userRepo        *repositories.UserRepository
	transactionRepo *repositories.TransactionRepository
	walletService   *WalletService
}

func NewClothingService(
	addressRepo *repositories.AddressRepository,
	userRepo *repositories.UserRepository,
	transactionRepo *repositories.TransactionRepository,
	walletService *WalletService,
) *ClothingService {
	return &ClothingService{
		addressRepo:     addressRepo,
		userRepo:        userRepo,
		transactionRepo: transactionRepo,
		walletService:   walletService,
	}
}

// ProcessAIClothingOrder processes natural language clothing requests and searches external APIs
func (cs *ClothingService) ProcessAIClothingOrder(userID string, req *models.AIClothingOrderRequest) (*models.AIClothingOrderResponse, error) {
	// Parse the AI request to extract clothing preferences
	analysis, err := cs.analyzeClothingRequest(req.Prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze clothing request: %w", err)
	}

	// Search external e-commerce APIs for matching products
	suggestedProducts, err := cs.searchExternalAPIs(analysis)
	if err != nil {
		return nil, fmt.Errorf("failed to search external APIs: %w", err)
	}

	// Check if user has delivery address
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	addresses, err := cs.addressRepo.GetByUserID(userUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user addresses: %w", err)
	}

	var selectedAddress *models.Address
	if len(addresses) > 0 {
		// Use first address as default or find default address
		for _, addr := range addresses {
			if addr.IsDefault {
				selectedAddress = &addr
				break
			}
		}
		if selectedAddress == nil {
			selectedAddress = &addresses[0]
		}
	}

	// Calculate total estimate
	var totalEstimate decimal.Decimal
	for _, product := range suggestedProducts {
		totalEstimate = totalEstimate.Add(decimal.NewFromFloat(product.Price))
	}

	// Generate confirmation ID for this request
	confirmationID := uuid.New().String()

	response := &models.AIClothingOrderResponse{
		OrderCreated:         false,
		SuggestedProducts:    suggestedProducts,
		TotalEstimate:        totalEstimate,
		Message:              cs.generateAIResponse(analysis, suggestedProducts),
		RequiresConfirmation: true,
		ConfirmationID:       confirmationID,
		OrderAnalysis:        analysis,
		SelectedAddress:      selectedAddress,
	}

	if len(addresses) == 0 {
		response.RequiredInfo = []string{"delivery_address"}
		response.Message = "I found some great clothing options for you! However, I need a delivery address to proceed with the order. Please add your address first."
	}

	return response, nil
}

// ConfirmAIClothingOrder confirms the order and processes payment
func (cs *ClothingService) ConfirmAIClothingOrder(userID string, req *models.ConfirmAIClothingOrderRequest) (*models.ConfirmAIClothingOrderResponse, error) {
	// Validate the selected product
	if req.SelectedProduct == nil {
		return nil, fmt.Errorf("no product selected")
	}

	// Validate address
	var deliveryAddress *models.Address
	if req.AddressID != nil {
		address, err := cs.addressRepo.GetByID(*req.AddressID)
		if err != nil {
			return nil, fmt.Errorf("failed to get address: %w", err)
		}
		deliveryAddress = address
	} else if req.NewAddress != nil {
		// Create new address
		newAddr := &models.Address{
			UserID:      userID,
			FullName:    req.NewAddress.FullName,
			PhoneNumber: req.NewAddress.PhoneNumber,
			AddressLine: req.NewAddress.AddressLine,
			City:        req.NewAddress.City,
			State:       req.NewAddress.State,
			PostalCode:  req.NewAddress.PostalCode,
			Country:     req.NewAddress.Country,
			IsDefault:   len(req.NewAddress.FullName) > 0, // Set as default if it's the first address
		}

		createdAddr, err := cs.addressRepo.Create(newAddr)
		if err != nil {
			return nil, fmt.Errorf("failed to create address: %w", err)
		}
		deliveryAddress = createdAddr
	} else {
		return nil, fmt.Errorf("delivery address is required")
	}

	// Calculate total amount
	totalAmount := req.SelectedProduct.Price * float64(req.Quantity)

	// Check wallet balance
	user, err := cs.userRepo.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	wallet, err := cs.walletService.GetWalletByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get wallet: %w", err)
	}

	if wallet.Balance < totalAmount {
		return nil, fmt.Errorf("insufficient wallet balance. Required: %.2f, Available: %.2f", totalAmount, wallet.Balance)
	}

	// Create transaction record
	transaction := &models.Transaction{
		UserID:        userID,
		Amount:        totalAmount,
		Type:          "debit",
		Status:        "completed",
		Description:   fmt.Sprintf("AI Clothing Order: %s", req.SelectedProduct.Name),
		ReferenceID:   uuid.New().String(),
		PaymentMethod: "wallet",
		Metadata: map[string]interface{}{
			"product_name":     req.SelectedProduct.Name,
			"product_url":      req.SelectedProduct.URL,
			"quantity":         req.Quantity,
			"delivery_address": deliveryAddress,
			"order_type":       "ai_clothing",
		},
	}

	// Process payment through wallet service
	_, err = cs.walletService.DebitWallet(userID, totalAmount, fmt.Sprintf("AI Clothing Order: %s", req.SelectedProduct.Name))
	if err != nil {
		return nil, fmt.Errorf("failed to process payment: %w", err)
	}

	// Save transaction
	createdTransaction, err := cs.transactionRepo.Create(transaction)
	if err != nil {
		// Rollback wallet transaction if database save fails
		cs.walletService.CreditWallet(userID, totalAmount, "Rollback for failed order")
		return nil, fmt.Errorf("failed to save transaction: %w", err)
	}

	// Create order confirmation
	orderConfirmation := &models.OrderConfirmation{
		OrderID:           createdTransaction.ReferenceID,
		ProductName:       req.SelectedProduct.Name,
		ProductURL:        req.SelectedProduct.URL,
		Quantity:          req.Quantity,
		TotalAmount:       totalAmount,
		DeliveryAddress:   *deliveryAddress,
		PaymentStatus:     "completed",
		OrderStatus:       "confirmed",
		EstimatedDelivery: time.Now().AddDate(0, 0, 7), // 7 days from now
		TrackingInfo:      fmt.Sprintf("TRZ-%s", createdTransaction.ReferenceID[:8]),
	}

	response := &models.ConfirmAIClothingOrderResponse{
		Success:           true,
		Message:           "Order confirmed successfully! Payment has been processed from your Tranza wallet.",
		OrderConfirmation: orderConfirmation,
		TransactionID:     createdTransaction.ReferenceID,
		WalletBalance:     wallet.Balance - totalAmount,
	}

	return response, nil
}

// parseClothingRequest extracts clothing preferences from natural language
func (cs *ClothingService) parseClothingRequest(message string) (*models.ClothingPreferences, error) {
	message = strings.ToLower(message)
	
	preferences := &models.ClothingPreferences{
		Category: "clothing",
		Gender:   "unisex",
		Color:    "any",
		Size:     "M", // default size
		Budget:   5000, // default budget in INR
	}

	// Extract category
	categories := map[string]string{
		"shirt":     "shirts",
		"t-shirt":   "t-shirts",
		"tshirt":    "t-shirts",
		"pant":      "pants",
		"trouser":   "pants",
		"jeans":     "jeans",
		"dress":     "dresses",
		"shoe":      "shoes",
		"sneaker":   "shoes",
		"jacket":    "jackets",
		"sweater":   "sweaters",
		"hoodie":    "hoodies",
		"skirt":     "skirts",
		"shorts":    "shorts",
	}

	for keyword, category := range categories {
		if strings.Contains(message, keyword) {
			preferences.Category = category
			break
		}
	}

	// Extract gender
	if strings.Contains(message, "men") || strings.Contains(message, "male") {
		preferences.Gender = "men"
	} else if strings.Contains(message, "women") || strings.Contains(message, "female") || strings.Contains(message, "ladies") {
		preferences.Gender = "women"
	}

	// Extract size
	sizes := []string{"xs", "s", "m", "l", "xl", "xxl", "xxxl"}
	for _, size := range sizes {
		if strings.Contains(message, "size "+size) || strings.Contains(message, " "+size+" ") {
			preferences.Size = strings.ToUpper(size)
			break
		}
	}

	// Extract color
	colors := []string{"black", "white", "red", "blue", "green", "yellow", "pink", "purple", "orange", "brown", "gray", "grey"}
	for _, color := range colors {
		if strings.Contains(message, color) {
			preferences.Color = color
			break
		}
	}

	// Extract budget (look for numbers followed by currency indicators)
	if strings.Contains(message, "₹") || strings.Contains(message, "rs") || strings.Contains(message, "rupees") {
		// Simple regex-like extraction for budget
		words := strings.Fields(message)
		for i, word := range words {
			if strings.Contains(word, "₹") || (i > 0 && (words[i-1] == "rs" || words[i-1] == "rupees")) {
				// Extract number
				numStr := strings.Trim(word, "₹rs ")
				if budget, err := strconv.ParseFloat(numStr, 64); err == nil {
					preferences.Budget = budget
				}
			}
		}
	}

	return preferences, nil
}

// searchExternalAPIs searches external e-commerce APIs for products
func (cs *ClothingService) searchExternalAPIs(preferences *models.ClothingPreferences) ([]*models.ExternalProduct, error) {
	var allProducts []*models.ExternalProduct

	// Search Myntra API (mock implementation)
	myntraProducts, err := cs.searchMyntraAPI(preferences)
	if err == nil {
		allProducts = append(allProducts, myntraProducts...)
	}

	// Search Flipkart API (mock implementation)
	flipkartProducts, err := cs.searchFlipkartAPI(preferences)
	if err == nil {
		allProducts = append(allProducts, flipkartProducts...)
	}

	// Search Amazon API (mock implementation)
	amazonProducts, err := cs.searchAmazonAPI(preferences)
	if err == nil {
		allProducts = append(allProducts, amazonProducts...)
	}

	// If no products found, return mock data for demonstration
	if len(allProducts) == 0 {
		allProducts = cs.getMockProducts(preferences)
	}

	return allProducts, nil
}

// searchMyntraAPI searches Myntra for products (mock implementation)
func (cs *ClothingService) searchMyntraAPI(preferences *models.ClothingPreferences) ([]*models.ExternalProduct, error) {
	// In a real implementation, you would call Myntra's API here
	// For now, return mock data
	
	searchQuery := fmt.Sprintf("%s %s %s", preferences.Gender, preferences.Color, preferences.Category)
	
	// Mock Myntra API response
	products := []*models.ExternalProduct{
		{
			ID:          "myntra_001",
			Name:        fmt.Sprintf("Myntra %s %s %s", preferences.Gender, preferences.Color, preferences.Category),
			Price:       preferences.Budget * 0.8, // 80% of budget
			OriginalPrice: preferences.Budget,
			Discount:    20,
			URL:         fmt.Sprintf("https://myntra.com/search?q=%s", url.QueryEscape(searchQuery)),
			ImageURL:    "https://via.placeholder.com/300x400?text=Myntra+Product",
			Brand:       "Myntra Brand",
			Rating:      4.2,
			Reviews:     1250,
			Sizes:       []string{preferences.Size, "L", "XL"},
			Colors:      []string{preferences.Color, "black", "white"},
			Description: fmt.Sprintf("High-quality %s from Myntra with premium fabric and comfortable fit", preferences.Category),
			Platform:    "Myntra",
			InStock:     true,
			Delivery:    "3-5 days",
		},
	}

	return products, nil
}

// searchFlipkartAPI searches Flipkart for products (mock implementation)
func (cs *ClothingService) searchFlipkartAPI(preferences *models.ClothingPreferences) ([]*models.ExternalProduct, error) {
	searchQuery := fmt.Sprintf("%s %s %s", preferences.Gender, preferences.Color, preferences.Category)
	
	products := []*models.ExternalProduct{
		{
			ID:          "flipkart_001",
			Name:        fmt.Sprintf("Flipkart %s %s %s", preferences.Gender, preferences.Color, preferences.Category),
			Price:       preferences.Budget * 0.7, // 70% of budget
			OriginalPrice: preferences.Budget * 0.9,
			Discount:    22,
			URL:         fmt.Sprintf("https://flipkart.com/search?q=%s", url.QueryEscape(searchQuery)),
			ImageURL:    "https://via.placeholder.com/300x400?text=Flipkart+Product",
			Brand:       "Flipkart Brand",
			Rating:      4.0,
			Reviews:     890,
			Sizes:       []string{preferences.Size, "M", "L"},
			Colors:      []string{preferences.Color, "navy", "gray"},
			Description: fmt.Sprintf("Trendy %s from Flipkart with modern design and great value", preferences.Category),
			Platform:    "Flipkart",
			InStock:     true,
			Delivery:    "2-4 days",
		},
	}

	return products, nil
}

// searchAmazonAPI searches Amazon for products (mock implementation)
func (cs *ClothingService) searchAmazonAPI(preferences *models.ClothingPreferences) ([]*models.ExternalProduct, error) {
	searchQuery := fmt.Sprintf("%s %s %s", preferences.Gender, preferences.Color, preferences.Category)
	
	products := []*models.ExternalProduct{
		{
			ID:          "amazon_001",
			Name:        fmt.Sprintf("Amazon %s %s %s", preferences.Gender, preferences.Color, preferences.Category),
			Price:       preferences.Budget * 0.9, // 90% of budget
			OriginalPrice: preferences.Budget * 1.1,
			Discount:    18,
			URL:         fmt.Sprintf("https://amazon.in/s?k=%s", url.QueryEscape(searchQuery)),
			ImageURL:    "https://via.placeholder.com/300x400?text=Amazon+Product",
			Brand:       "Amazon Brand",
			Rating:      4.3,
			Reviews:     2100,
			Sizes:       []string{preferences.Size, "S", "M", "L"},
			Colors:      []string{preferences.Color, "charcoal", "beige"},
			Description: fmt.Sprintf("Premium %s from Amazon with excellent quality and fast delivery", preferences.Category),
			Platform:    "Amazon",
			InStock:     true,
			Delivery:    "1-2 days",
		},
	}

	return products, nil
}

// getMockProducts returns mock products when API calls fail
func (cs *ClothingService) getMockProducts(preferences *models.ClothingPreferences) []*models.ExternalProduct {
	return []*models.ExternalProduct{
		{
			ID:          "mock_001",
			Name:        fmt.Sprintf("Classic %s %s %s", preferences.Gender, preferences.Color, preferences.Category),
			Price:       preferences.Budget * 0.75,
			OriginalPrice: preferences.Budget,
			Discount:    25,
			URL:         "https://example.com/product/1",
			ImageURL:    "https://via.placeholder.com/300x400?text=Product+Image",
			Brand:       "Generic Brand",
			Rating:      4.1,
			Reviews:     567,
			Sizes:       []string{preferences.Size, "M", "L", "XL"},
			Colors:      []string{preferences.Color, "black", "white", "gray"},
			Description: fmt.Sprintf("High-quality %s with comfortable fit and durable material", preferences.Category),
			Platform:    "Demo Store",
			InStock:     true,
			Delivery:    "3-5 days",
		},
		{
			ID:          "mock_002",
			Name:        fmt.Sprintf("Premium %s %s %s", preferences.Gender, preferences.Color, preferences.Category),
			Price:       preferences.Budget * 0.85,
			OriginalPrice: preferences.Budget * 1.2,
			Discount:    29,
			URL:         "https://example.com/product/2",
			ImageURL:    "https://via.placeholder.com/300x400?text=Premium+Product",
			Brand:       "Premium Brand",
			Rating:      4.5,
			Reviews:     1234,
			Sizes:       []string{preferences.Size, "S", "M", "L"},
			Colors:      []string{preferences.Color, "navy", "maroon"},
			Description: fmt.Sprintf("Premium quality %s with luxury finish and elegant design", preferences.Category),
			Platform:    "Premium Store",
			InStock:     true,
			Delivery:    "2-3 days",
		},
	}
}

// generateAIResponse creates a natural language response for the user
func (cs *ClothingService) generateAIResponse(preferences *models.ClothingPreferences, products []*models.ExternalProduct) string {
	if len(products) == 0 {
		return "I couldn't find any products matching your requirements. Please try with different preferences."
	}

	response := fmt.Sprintf("Great! I found %d products for you based on your request for %s %s %s", 
		len(products), preferences.Gender, preferences.Color, preferences.Category)

	if preferences.Budget > 0 {
		response += fmt.Sprintf(" within your budget of ₹%.0f", preferences.Budget)
	}

	response += ". Here are the best options I found from different platforms. You can select any product to proceed with the order using your Tranza wallet!"

	return response
}
