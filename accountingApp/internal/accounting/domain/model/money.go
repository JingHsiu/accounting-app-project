package model

import (
	"errors"
	"fmt"
)

type Money struct {
	Amount   int64
	Currency string
}

func NewMoney(amount int64, currency string) (*Money, error) {
	if amount < 0 {
		return nil, errors.New("amount cannot be negative")
	}
	if currency == "" {
		return nil, errors.New("currency cannot be empty")
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
	return fmt.Sprintf("%.2f %s", float64(m.Amount)/100, m.Currency)
}