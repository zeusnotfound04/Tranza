package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/zeusnotfound04/Tranza/models"
	"gorm.io/gorm"
)

var db *gorm.DB

// UserRepository interface defines the contract for user data operations
type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	FindByProviderID(ctx context.Context, provider, providerID string) (*models.User, error)
	Create(ctx context.Context, user *models.User) error
}

// userRepository struct implements UserRepository interface
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new UserRepository instance
func NewUserRepository(database *gorm.DB) UserRepository {
	return &userRepository{db: database}
}

func InitRepo(database *gorm.DB) {
	db = database
}

// Legacy functions - keeping for backward compatibility
func GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func CreateUser(user *models.User) error {
	return db.Create(user).Error
}

// UserRepository interface implementations
func (r *userRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByProviderID(ctx context.Context, provider, providerID string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Where("provider = ? AND provider_id = ?", provider, providerID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}
