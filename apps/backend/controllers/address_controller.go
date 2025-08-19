package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zeusnotfound04/Tranza/models"
	"github.com/zeusnotfound04/Tranza/services"
	"github.com/zeusnotfound04/Tranza/utils"
)

type AddressController struct {
	addressService *services.AddressService
}

func NewAddressController(addressService *services.AddressService) *AddressController {
	return &AddressController{
		addressService: addressService,
	}
}

// CreateAddress creates a new address for the user
func (ac *AddressController) CreateAddress(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req models.AddressCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	address, err := ac.addressService.CreateAddress(userID.(string), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Address created successfully", address)
}

// GetAddresses retrieves all addresses for the user
func (ac *AddressController) GetAddresses(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	addresses, err := ac.addressService.GetAddresses(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Addresses retrieved successfully", addresses)
}

// GetAddress retrieves a specific address
func (ac *AddressController) GetAddress(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	addressID := c.Param("id")
	if addressID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Address ID is required"})
		return
	}

	address, err := ac.addressService.GetAddress(userID.(string), addressID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Address retrieved successfully", address)
}

// GetDefaultAddress retrieves the default address for the user
func (ac *AddressController) GetDefaultAddress(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	address, err := ac.addressService.GetDefaultAddress(userID.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No default address found"})
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Default address retrieved successfully", address)
}

// UpdateAddress updates an existing address
func (ac *AddressController) UpdateAddress(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	addressID := c.Param("id")
	if addressID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Address ID is required"})
		return
	}

	var req models.AddressUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	address, err := ac.addressService.UpdateAddress(userID.(string), addressID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Address updated successfully", address)
}

// DeleteAddress deletes an address
func (ac *AddressController) DeleteAddress(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	addressID := c.Param("id")
	if addressID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Address ID is required"})
		return
	}

	err := ac.addressService.DeleteAddress(userID.(string), addressID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Address deleted successfully", nil)
}

// SetDefaultAddress sets an address as default
func (ac *AddressController) SetDefaultAddress(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	addressID := c.Param("id")
	if addressID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Address ID is required"})
		return
	}

	err := ac.addressService.SetDefaultAddress(userID.(string), addressID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Default address updated successfully", nil)
}
