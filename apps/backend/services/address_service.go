package services

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/zeusnotfound04/Tranza/models"
	"github.com/zeusnotfound04/Tranza/repositories"
)

type AddressService struct {
	addressRepo *repositories.AddressRepository
}

func NewAddressService(addressRepo *repositories.AddressRepository) *AddressService {
	return &AddressService{
		addressRepo: addressRepo,
	}
}

// CreateAddress creates a new address for the user
func (s *AddressService) CreateAddress(userID string, req *models.AddressCreateRequest) (*models.AddressResponse, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	// Check if this is the first address and set as default
	addressCount, err := s.addressRepo.Count(uid)
	if err != nil {
		return nil, fmt.Errorf("failed to count addresses: %v", err)
	}

	isDefault := req.IsDefault || addressCount == 0

	address := &models.Address{
		UserID:      uid,
		Name:        req.Name,
		Phone:       req.Phone,
		AddressLine: req.AddressLine,
		City:        req.City,
		State:       req.State,
		PinCode:     req.PinCode,
		Country:     req.Country,
		Landmark:    req.Landmark,
		IsDefault:   isDefault,
		AddressType: req.AddressType,
	}

	if address.Country == "" {
		address.Country = "India"
	}
	if address.AddressType == "" {
		address.AddressType = "home"
	}

	createdAddress, err := s.addressRepo.Create(address)
	if err != nil {
		return nil, fmt.Errorf("failed to create address: %v", err)
	}

	return s.convertToResponse(createdAddress), nil
}

// GetAddresses retrieves all addresses for the user
func (s *AddressService) GetAddresses(userID string) ([]models.AddressResponse, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	addresses, err := s.addressRepo.GetByUserID(uid)
	if err != nil {
		return nil, fmt.Errorf("failed to get addresses: %v", err)
	}

	responses := make([]models.AddressResponse, len(addresses))
	for i, addr := range addresses {
		responses[i] = *s.convertToResponse(&addr)
	}

	return responses, nil
}

// GetAddress retrieves a specific address
func (s *AddressService) GetAddress(userID, addressID string) (*models.AddressResponse, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	aid, err := uuid.Parse(addressID)
	if err != nil {
		return nil, errors.New("invalid address ID")
	}

	address, err := s.addressRepo.GetByID(aid)
	if err != nil {
		return nil, fmt.Errorf("address not found: %v", err)
	}

	// Verify ownership
	if address.UserID != uid {
		return nil, errors.New("address not found")
	}

	return s.convertToResponse(address), nil
}

// GetDefaultAddress retrieves the default address for the user
func (s *AddressService) GetDefaultAddress(userID string) (*models.AddressResponse, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	address, err := s.addressRepo.GetDefaultByUserID(uid)
	if err != nil {
		return nil, fmt.Errorf("no default address found: %v", err)
	}

	return s.convertToResponse(address), nil
}

// UpdateAddress updates an existing address
func (s *AddressService) UpdateAddress(userID, addressID string, req *models.AddressUpdateRequest) (*models.AddressResponse, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	aid, err := uuid.Parse(addressID)
	if err != nil {
		return nil, errors.New("invalid address ID")
	}

	address, err := s.addressRepo.GetByID(aid)
	if err != nil {
		return nil, fmt.Errorf("address not found: %v", err)
	}

	// Verify ownership
	if address.UserID != uid {
		return nil, errors.New("address not found")
	}

	// Update fields if provided
	if req.Name != nil {
		address.Name = *req.Name
	}
	if req.Phone != nil {
		address.Phone = *req.Phone
	}
	if req.AddressLine != nil {
		address.AddressLine = *req.AddressLine
	}
	if req.City != nil {
		address.City = *req.City
	}
	if req.State != nil {
		address.State = *req.State
	}
	if req.PinCode != nil {
		address.PinCode = *req.PinCode
	}
	if req.Country != nil {
		address.Country = *req.Country
	}
	if req.Landmark != nil {
		address.Landmark = *req.Landmark
	}
	if req.AddressType != nil {
		address.AddressType = *req.AddressType
	}
	if req.IsDefault != nil {
		address.IsDefault = *req.IsDefault
	}

	err = s.addressRepo.Update(address)
	if err != nil {
		return nil, fmt.Errorf("failed to update address: %v", err)
	}

	return s.convertToResponse(address), nil
}

// DeleteAddress deletes an address
func (s *AddressService) DeleteAddress(userID, addressID string) error {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return errors.New("invalid user ID")
	}

	aid, err := uuid.Parse(addressID)
	if err != nil {
		return errors.New("invalid address ID")
	}

	// Check if address exists and belongs to user
	address, err := s.addressRepo.GetByID(aid)
	if err != nil {
		return fmt.Errorf("address not found: %v", err)
	}

	if address.UserID != uid {
		return errors.New("address not found")
	}

	// Don't allow deletion of default address if there are other addresses
	if address.IsDefault {
		count, err := s.addressRepo.Count(uid)
		if err != nil {
			return fmt.Errorf("failed to count addresses: %v", err)
		}
		if count > 1 {
			return errors.New("cannot delete default address. Please set another address as default first")
		}
	}

	return s.addressRepo.Delete(aid, uid)
}

// SetDefaultAddress sets an address as default
func (s *AddressService) SetDefaultAddress(userID, addressID string) error {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return errors.New("invalid user ID")
	}

	aid, err := uuid.Parse(addressID)
	if err != nil {
		return errors.New("invalid address ID")
	}

	// Check if address exists and belongs to user
	address, err := s.addressRepo.GetByID(aid)
	if err != nil {
		return fmt.Errorf("address not found: %v", err)
	}

	if address.UserID != uid {
		return errors.New("address not found")
	}

	return s.addressRepo.SetDefault(aid, uid)
}

// Helper method to convert address to response
func (s *AddressService) convertToResponse(address *models.Address) *models.AddressResponse {
	return &models.AddressResponse{
		ID:          address.ID.String(),
		Name:        address.Name,
		Phone:       address.Phone,
		AddressLine: address.AddressLine,
		City:        address.City,
		State:       address.State,
		PinCode:     address.PinCode,
		Country:     address.Country,
		Landmark:    address.Landmark,
		IsDefault:   address.IsDefault,
		AddressType: address.AddressType,
		CreatedAt:   address.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}
