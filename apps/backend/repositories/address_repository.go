package repositories

import (
	"github.com/google/uuid"
	"github.com/zeusnotfound04/Tranza/models"
	"gorm.io/gorm"
)

type AddressRepository struct {
	db *gorm.DB
}

func NewAddressRepository(db *gorm.DB) *AddressRepository {
	return &AddressRepository{db: db}
}

// Create creates a new address
func (r *AddressRepository) Create(address *models.Address) (*models.Address, error) {
	if err := r.db.Create(address).Error; err != nil {
		return nil, err
	}
	return address, nil
}

// GetByID retrieves address by ID
func (r *AddressRepository) GetByID(id uuid.UUID) (*models.Address, error) {
	var address models.Address
	if err := r.db.Where("id = ?", id).First(&address).Error; err != nil {
		return nil, err
	}
	return &address, nil
}

// GetByUserID retrieves all addresses for a user
func (r *AddressRepository) GetByUserID(userID uuid.UUID) ([]models.Address, error) {
	var addresses []models.Address
	if err := r.db.Where("user_id = ?", userID).Order("is_default DESC, created_at DESC").Find(&addresses).Error; err != nil {
		return nil, err
	}
	return addresses, nil
}

// GetDefaultByUserID retrieves the default address for a user
func (r *AddressRepository) GetDefaultByUserID(userID uuid.UUID) (*models.Address, error) {
	var address models.Address
	if err := r.db.Where("user_id = ? AND is_default = ?", userID, true).First(&address).Error; err != nil {
		return nil, err
	}
	return &address, nil
}

// Update updates an existing address
func (r *AddressRepository) Update(address *models.Address) error {
	return r.db.Save(address).Error
}

// Delete deletes an address
func (r *AddressRepository) Delete(id uuid.UUID, userID uuid.UUID) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Address{}).Error
}

// SetDefault sets an address as default and unsets others
func (r *AddressRepository) SetDefault(id uuid.UUID, userID uuid.UUID) error {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Unset all other defaults for this user
	if err := tx.Model(&models.Address{}).Where("user_id = ?", userID).Update("is_default", false).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Set the specified address as default
	if err := tx.Model(&models.Address{}).Where("id = ? AND user_id = ?", id, userID).Update("is_default", true).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// Count returns the number of addresses for a user
func (r *AddressRepository) Count(userID uuid.UUID) (int64, error) {
	var count int64
	if err := r.db.Model(&models.Address{}).Where("user_id = ?", userID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
