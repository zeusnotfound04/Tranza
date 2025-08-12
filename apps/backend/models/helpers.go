package models

import "github.com/shopspring/decimal"

// DecimalFromFloat64 converts a float64 to decimal.Decimal
func DecimalFromFloat64(f float64) decimal.Decimal {
	return decimal.NewFromFloat(f)
}

// DecimalFromString converts a string to decimal.Decimal
func DecimalFromString(s string) (decimal.Decimal, error) {
	return decimal.NewFromString(s)
}
