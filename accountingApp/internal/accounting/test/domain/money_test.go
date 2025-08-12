package domain

import (
	"testing"
	"github.com/JingHsiu/accountingApp/internal/accounting/domain/model"
	"github.com/stretchr/testify/assert"
)

func TestNewMoney_Success(t *testing.T) {
	money, err := model.NewMoney(10000, "USD")
	
	assert.NoError(t, err)
	assert.Equal(t, int64(10000), money.Amount)
	assert.Equal(t, "USD", money.Currency)
}

func TestNewMoney_InvalidAmount(t *testing.T) {
	money, err := model.NewMoney(-100, "USD")
	
	assert.Error(t, err)
	assert.Nil(t, money)
	assert.Contains(t, err.Error(), "cannot be negative")
}

func TestNewMoney_InvalidCurrency(t *testing.T) {
	money, err := model.NewMoney(100, "US")
	
	assert.Error(t, err)
	assert.Nil(t, money)
	assert.Contains(t, err.Error(), "must be 3 characters")
}

func TestMoney_Add_Success(t *testing.T) {
	money1, _ := model.NewMoney(1000, "USD")
	money2, _ := model.NewMoney(500, "USD")
	
	result, err := money1.Add(*money2)
	
	assert.NoError(t, err)
	assert.Equal(t, int64(1500), result.Amount)
	assert.Equal(t, "USD", result.Currency)
}

func TestMoney_Add_DifferentCurrencies(t *testing.T) {
	money1, _ := model.NewMoney(1000, "USD")
	money2, _ := model.NewMoney(500, "EUR")
	
	result, err := money1.Add(*money2)
	
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "cannot add different currencies")
}

func TestMoney_Subtract_Success(t *testing.T) {
	money1, _ := model.NewMoney(1000, "USD")
	money2, _ := model.NewMoney(300, "USD")
	
	result, err := money1.Subtract(*money2)
	
	assert.NoError(t, err)
	assert.Equal(t, int64(700), result.Amount)
	assert.Equal(t, "USD", result.Currency)
}

func TestMoney_Subtract_NegativeResult(t *testing.T) {
	money1, _ := model.NewMoney(100, "USD")
	money2, _ := model.NewMoney(300, "USD")
	
	result, err := money1.Subtract(*money2)
	
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "cannot be negative")
}