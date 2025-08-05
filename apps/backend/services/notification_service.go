package services

import (
	"fmt"
	"log"

	"github.com/shopspring/decimal"
)

type NotificationService struct {
	// Add SMS, Email, Push notification clients here
}

func NewNotificationService() *NotificationService {
	return &NotificationService{}
}

// Send wallet credit notification
func (s *NotificationService) SendWalletCreditNotification(userID string, amount, newBalance decimal.Decimal) {
	message := fmt.Sprintf("₹%.2f added to your wallet. New balance: ₹%.2f", amount, newBalance)
	log.Printf("Notification to user %s: %s", userID, message)
	// Implement actual notification logic (SMS, Email, Push)
}

// Send AI payment notification
func (s *NotificationService) SendAIPaymentNotification(userID, agentID string, amount decimal.Decimal, merchantName string, newBalance decimal.Decimal) {
	message := fmt.Sprintf("AI Agent %s made payment of ₹%.2f to %s. Balance: ₹%.2f", agentID, amount, merchantName, newBalance)
	log.Printf("Notification to user %s: %s", userID, message)
	// Implement actual notification logic (SMS, Email, Push)
}