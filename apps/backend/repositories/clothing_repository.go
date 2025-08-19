package repositories

import (
	"github.com/google/uuid"
	"github.com/zeusnotfound04/Tranza/models"
	"gorm.io/gorm"
)

type ClothingRepository struct {
	db *gorm.DB
}

func NewClothingRepository(db *gorm.DB) *ClothingRepository {
	return &ClothingRepository{db: db}
}

// Product operations
func (r *ClothingRepository) CreateProduct(product *models.Product) (*models.Product, error) {
	if err := r.db.Create(product).Error; err != nil {
		return nil, err
	}
	return product, nil
}

func (r *ClothingRepository) GetProductByID(id uuid.UUID) (*models.Product, error) {
	var product models.Product
	if err := r.db.Preload("Variants").Where("id = ? AND is_active = ?", id, true).First(&product).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ClothingRepository) GetProducts(category, subCategory string, limit, offset int) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	query := r.db.Model(&models.Product{}).Where("is_active = ?", true)
	
	if category != "" {
		query = query.Where("category = ?", category)
	}
	if subCategory != "" {
		query = query.Where("sub_category = ?", subCategory)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results with variants
	if err := query.Preload("Variants").Limit(limit).Offset(offset).Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (r *ClothingRepository) SearchProducts(searchTerm, category string, minPrice, maxPrice float64, limit, offset int) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	query := r.db.Model(&models.Product{}).Where("is_active = ?", true)
	
	if searchTerm != "" {
		query = query.Where("name ILIKE ? OR description ILIKE ? OR brand ILIKE ?", 
			"%"+searchTerm+"%", "%"+searchTerm+"%", "%"+searchTerm+"%")
	}
	if category != "" {
		query = query.Where("category = ?", category)
	}
	if minPrice > 0 {
		query = query.Where("price >= ?", minPrice)
	}
	if maxPrice > 0 {
		query = query.Where("price <= ?", maxPrice)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	if err := query.Preload("Variants").Limit(limit).Offset(offset).Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

// Product Variant operations
func (r *ClothingRepository) GetVariantByID(id uuid.UUID) (*models.ProductVariant, error) {
	var variant models.ProductVariant
	if err := r.db.Preload("Product").Where("id = ? AND is_active = ?", id, true).First(&variant).Error; err != nil {
		return nil, err
	}
	return &variant, nil
}

func (r *ClothingRepository) UpdateVariantStock(variantID uuid.UUID, quantity int) error {
	return r.db.Model(&models.ProductVariant{}).Where("id = ?", variantID).Update("stock", gorm.Expr("stock - ?", quantity)).Error
}

// Order operations
func (r *ClothingRepository) CreateOrder(order *models.Order) (*models.Order, error) {
	if err := r.db.Create(order).Error; err != nil {
		return nil, err
	}
	return order, nil
}

func (r *ClothingRepository) GetOrderByID(id uuid.UUID) (*models.Order, error) {
	var order models.Order
	if err := r.db.Preload("OrderItems.Product").Preload("OrderItems.Variant").Where("id = ?", id).First(&order).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *ClothingRepository) GetOrdersByUserID(userID uuid.UUID, limit, offset int) ([]models.Order, int64, error) {
	var orders []models.Order
	var total int64

	// Count total
	if err := r.db.Model(&models.Order{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	if err := r.db.Preload("OrderItems.Product").Preload("OrderItems.Variant").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func (r *ClothingRepository) UpdateOrderStatus(orderID uuid.UUID, status string) error {
	return r.db.Model(&models.Order{}).Where("id = ?", orderID).Update("status", status).Error
}

func (r *ClothingRepository) UpdateOrderPaymentStatus(orderID uuid.UUID, status, transactionID string) error {
	updates := map[string]interface{}{
		"payment_status": status,
	}
	if transactionID != "" {
		updates["transaction_id"] = transactionID
	}
	return r.db.Model(&models.Order{}).Where("id = ?", orderID).Updates(updates).Error
}

// Order Item operations
func (r *ClothingRepository) CreateOrderItems(orderItems []models.OrderItem) error {
	return r.db.Create(&orderItems).Error
}

func (r *ClothingRepository) GetOrderItemsByOrderID(orderID uuid.UUID) ([]models.OrderItem, error) {
	var orderItems []models.OrderItem
	if err := r.db.Preload("Product").Preload("Variant").Where("order_id = ?", orderID).Find(&orderItems).Error; err != nil {
		return nil, err
	}
	return orderItems, nil
}

// AI-specific operations
func (r *ClothingRepository) GetAIRecommendations(category, subCategory string, minPrice, maxPrice float64, limit int) ([]models.Product, error) {
	var products []models.Product
	
	query := r.db.Model(&models.Product{}).Where("is_active = ?", true)
	
	if category != "" {
		query = query.Where("category = ?", category)
	}
	if subCategory != "" {
		query = query.Where("sub_category = ?", subCategory)
	}
	if minPrice > 0 {
		query = query.Where("price >= ?", minPrice)
	}
	if maxPrice > 0 {
		query = query.Where("price <= ?", maxPrice)
	}

	// Order by rating and review count for better recommendations
	if err := query.Preload("Variants").
		Order("rating DESC, review_count DESC, created_at DESC").
		Limit(limit).
		Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

func (r *ClothingRepository) GetCategories() ([]string, error) {
	var categories []string
	if err := r.db.Model(&models.Product{}).Distinct("category").Where("is_active = ?", true).Pluck("category", &categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *ClothingRepository) GetSizes() ([]string, error) {
	var sizes []string
	if err := r.db.Model(&models.ProductVariant{}).Distinct("size").Where("is_active = ?", true).Pluck("size", &sizes).Error; err != nil {
		return nil, err
	}
	return sizes, nil
}

func (r *ClothingRepository) GetColors() ([]string, error) {
	var colors []string
	if err := r.db.Model(&models.ProductVariant{}).Distinct("color").Where("is_active = ?", true).Pluck("color", &colors).Error; err != nil {
		return nil, err
	}
	return colors, nil
}
