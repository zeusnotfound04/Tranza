package models

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/shopspring/decimal"
)

// DecimalFromFloat64 converts a float64 to decimal.Decimal
func DecimalFromFloat64(f float64) decimal.Decimal {
	return decimal.NewFromFloat(f)
}

// DecimalFromString converts a string to decimal.Decimal
func DecimalFromString(s string) (decimal.Decimal, error) {
	return decimal.NewFromString(s)
}

// GenerateExternalTransferReference generates a unique reference ID for external transfers
func GenerateExternalTransferReference() string {
	// Format: EXT_YYYYMMDD_HHMMSS_RAND
	now := time.Now()
	randNum := rand.Intn(9999)
	return fmt.Sprintf("EXT_%s_%04d",
		now.Format("20060102_150405"),
		randNum)
}
