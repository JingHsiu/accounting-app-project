package model

import (
	"errors"
	"fmt"
)

type Money struct {
	Amount   int64
	Currency string
}

// GetCurrencySubdivision returns the subdivision for different currencies
// TWD and other whole-number currencies use 1, decimal currencies use 100
func GetCurrencySubdivision(currency string) int64 {
	currencySubdivisions := map[string]int64{
		// Whole number currencies (no subdivision) - TWD as primary
		"TWD": 1,   // 1 台幣 = 1 台幣 (no cents, base unit)
		"JPY": 1,   // 1 yen = 1 yen (no sen in practice)
		"KRW": 1,   // 1 won = 1 won (no subdivision)
		"VND": 1,   // 1 dong = 1 dong (no subdivision)
		
		// Decimal currencies (1 unit = 100 smaller units)
		"USD": 100, // 1 dollar = 100 cents
		"EUR": 100, // 1 euro = 100 cents
		"GBP": 100, // 1 pound = 100 pence
		"CNY": 100, // 1 yuan = 100 fen
	}
	
	if subdivision, exists := currencySubdivisions[currency]; exists {
		return subdivision
	}
	return 1 // Default to 1 for unknown currencies (like TWD)
}

func NewMoney(amount int64, currency string) (*Money, error) {
	if amount < 0 {
		return nil, errors.New("amount cannot be negative")
	}
	if currency == "" {
		currency = "TWD" // Default to TWD
	}
	if len(currency) != 3 {
		return nil, errors.New("currency must be 3 characters (ISO 4217)")
	}
	
	return &Money{
		Amount:   amount,
		Currency: currency,
	}, nil
}

func (m Money) Add(other Money) (*Money, error) {
	if m.Currency != other.Currency {
		return nil, fmt.Errorf("cannot add different currencies: %s and %s", m.Currency, other.Currency)
	}
	return NewMoney(m.Amount+other.Amount, m.Currency)
}

func (m Money) Subtract(other Money) (*Money, error) {
	if m.Currency != other.Currency {
		return nil, fmt.Errorf("cannot subtract different currencies: %s and %s", m.Currency, other.Currency)
	}
	result := m.Amount - other.Amount
	if result < 0 {
		return nil, errors.New("result cannot be negative")
	}
	return NewMoney(result, m.Currency)
}

func (m Money) Equals(other Money) bool {
	return m.Amount == other.Amount && m.Currency == other.Currency
}

func (m Money) String() string {
	subdivision := GetCurrencySubdivision(m.Currency)
	displayAmount := float64(m.Amount) / float64(subdivision)
	
	// For whole number currencies, don't show decimal places
	if subdivision == 1 {
		return fmt.Sprintf("%.0f %s", displayAmount, m.Currency)
	}
	// For decimal currencies, show 2 decimal places
	return fmt.Sprintf("%.2f %s", displayAmount, m.Currency)
}