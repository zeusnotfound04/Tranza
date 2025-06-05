package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zeusnotfound04/Tranza/models"
	"github.com/zeusnotfound04/Tranza/services"
)

type CardController struct {
	Service *services.CardService
}

func NewCardController(service *services.CardService) *CardController {
	return &CardController{Service: service}
}

func (cc *CardController) LinkCard(c *gin.Context) {
	var card models.LinkedCard
	if err := c.ShouldBindJSON(&card); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	userID := c.GetUint("user_id")
	card.UserId = userID

	if err := cc.Service.LinkCard(&card); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to link card"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Card linked"})
}

func (cc *CardController) GetCards(c *gin.Context) {
	userID := c.GetUint("user_id")

	cards, err := cc.Service.GetCards(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch cards"})
		return
	}

	c.JSON(http.StatusOK, cards)
}

func (cc *CardController) DeleteCard(c *gin.Context) {
	cardID, _ := strconv.Atoi(c.Param("id"))
	userID := c.GetUint("user_id")

	err := cc.Service.DeleteCard(uint(cardID), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete card"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Card deleted"})
}

func (cc *CardController) UpdateLimit(c *gin.Context) {
	cardID, _ := strconv.Atoi(c.Param("id"))
	userID := c.GetUint("user_id")

	var payload struct {
		Limit float64 `json:"limit"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	err := cc.Service.UpdateCardLimit(uint(cardID), userID, payload.Limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update limit"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Limit updated"})
}
