package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// ParseUUID safely parses UUID string
func ParseUUID(uuidStr string) (uuid.UUID, error) {
	if uuidStr == "" {
		return uuid.Nil, fmt.Errorf("UUID string is empty")
	}
	return uuid.Parse(uuidStr)
}

// GetUserIDFromContext extracts and converts user_id from gin context to UUID
func GetUserIDFromContext(ctx *gin.Context) (uuid.UUID, error) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		return uuid.Nil, fmt.Errorf("user not authenticated")
	}

	// Handle both UUID and string types for backward compatibility
	switch v := userID.(type) {
	case uuid.UUID:
		return v, nil
	case string:
		return uuid.Parse(v)
	default:
		return uuid.Nil, fmt.Errorf("invalid user ID type: %T", userID)
	}
}

// GetUserIDStringFromContext extracts user_id from gin context as string
func GetUserIDStringFromContext(ctx *gin.Context) (string, error) {
	userUUID, err := GetUserIDFromContext(ctx)
	if err != nil {
		return "", err
	}
	return userUUID.String(), nil
}

// FormatAmount formats decimal amount to string with 2 decimal places
func FormatAmount(amount decimal.Decimal) string {
	return fmt.Sprintf("₹%.2f", amount)
}

// ParseAmount parses string to decimal amount
func ParseAmount(amountStr string) (decimal.Decimal, error) {
	// Remove currency symbols and spaces
	amountStr = strings.ReplaceAll(amountStr, "₹", "")
	amountStr = strings.ReplaceAll(amountStr, ",", "")
	amountStr = strings.TrimSpace(amountStr)

	return decimal.NewFromString(amountStr)
}

// GetPaginationParams extracts pagination parameters from query
func GetPaginationParams(c *gin.Context) (page, limit int) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "50")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err = strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 50
	}

	return page, limit
}

// CalculateOffset calculates database offset from page and limit
func CalculateOffset(page, limit int) int {
	return (page - 1) * limit
}

// FormatPhoneNumber formats phone number to standard format
func FormatPhoneNumber(phone string) string {
	// Remove all non-digit characters
	phone = strings.ReplaceAll(phone, " ", "")
	phone = strings.ReplaceAll(phone, "-", "")
	phone = strings.ReplaceAll(phone, "+91", "")
	phone = strings.ReplaceAll(phone, "(", "")
	phone = strings.ReplaceAll(phone, ")", "")

	// Add +91 prefix if not present
	if len(phone) == 10 && !strings.HasPrefix(phone, "+91") {
		phone = "+91" + phone
	}

	return phone
}

// MaskPhoneNumber masks phone number for display (shows last 4 digits)
func MaskPhoneNumber(phone string) string {
	if len(phone) < 4 {
		return phone
	}

	masked := strings.Repeat("*", len(phone)-4) + phone[len(phone)-4:]
	return masked
}

// MaskEmail masks email for display
func MaskEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return email
	}

	username := parts[0]
	domain := parts[1]

	if len(username) <= 2 {
		return email
	}

	maskedUsername := string(username[0]) + strings.Repeat("*", len(username)-2) + string(username[len(username)-1])
	return maskedUsername + "@" + domain
}

// IsValidTimeRange checks if time range is valid
func IsValidTimeRange(startTime, endTime time.Time) bool {
	return startTime.Before(endTime) && !startTime.IsZero() && !endTime.IsZero()
}

// ConvertToIST converts UTC time to IST
func ConvertToIST(utcTime time.Time) time.Time {
	ist, _ := time.LoadLocation("Asia/Kolkata")
	return utcTime.In(ist)
}

// GetStartAndEndOfDay returns start and end time of a given date
func GetStartAndEndOfDay(date time.Time) (time.Time, time.Time) {
	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	end := start.Add(24 * time.Hour).Add(-time.Nanosecond)
	return start, end
}

// GetClientIP extracts client IP from gin context
func GetClientIP(c *gin.Context) string {
	// Check X-Forwarded-For header first
	forwarded := c.GetHeader("X-Forwarded-For")
	if forwarded != "" {
		// X-Forwarded-For can contain multiple IPs, get the first one
		ips := strings.Split(forwarded, ",")
		return strings.TrimSpace(ips[0])
	}

	// Check X-Real-IP header
	realIP := c.GetHeader("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	// Fall back to RemoteAddr
	return c.ClientIP()
}

// SanitizeString removes potentially harmful characters from string
func SanitizeString(input string) string {
	// Remove null bytes and control characters
	sanitized := strings.Map(func(r rune) rune {
		if r == 0 || (r < 32 && r != 9 && r != 10 && r != 13) {
			return -1
		}
		return r
	}, input)

	return strings.TrimSpace(sanitized)
}

// ParseDecimal safely parses string to decimal
func ParseDecimal(value string) decimal.Decimal {
	result, err := decimal.NewFromString(value)
	if err != nil {
		return decimal.Zero
	}
	return result
}

// GetCurrentTimestamp returns current Unix timestamp
func GetCurrentTimestamp() int64 {
	return time.Now().Unix()
}
