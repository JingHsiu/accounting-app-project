package repository

import (
	"testing"
	"time"

	"github.com/JingHsiu/accountingApp/internal/accounting/application/mapper"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/repository"
	"github.com/JingHsiu/accountingApp/internal/accounting/domain/model"
)

// Mock peer implementation for testing
type MockWalletRepositoryPeer struct {
	data      map[string]mapper.WalletData
	userData  map[string][]mapper.WalletData
	saveFunc  func(data mapper.WalletData) error
	findFunc  func(id string) (*mapper.WalletData, error)
	deleteFunc func(id string) error
}

func NewMockWalletRepositoryPeer() *MockWalletRepositoryPeer {
	return &MockWalletRepositoryPeer{
		data:     make(map[string]mapper.WalletData),
		userData: make(map[string][]mapper.WalletData),
	}
}

func (m *MockWalletRepositoryPeer) Save(data mapper.WalletData) error {
	if m.saveFunc != nil {
		return m.saveFunc(data)
	}
	
	m.data[data.ID] = data
	
	// Update user data index
	userWallets := m.userData[data.UserID]
	found := false
	for i, wallet := range userWallets {
		if wallet.ID == data.ID {
			userWallets[i] = data
			found = true
			break
		}
	}
	if !found {
		userWallets = append(userWallets, data)
	}
	m.userData[data.UserID] = userWallets
	
	return nil
}

func (m *MockWalletRepositoryPeer) FindByID(id string) (*mapper.WalletData, error) {
	if m.findFunc != nil {
		return m.findFunc(id)
	}
	
	if data, exists := m.data[id]; exists {
		return &data, nil
	}
	return nil, nil
}

func (m *MockWalletRepositoryPeer) FindByUserID(userID string) ([]mapper.WalletData, error) {
	if wallets, exists := m.userData[userID]; exists {
		return wallets, nil
	}
	return []mapper.WalletData{}, nil
}

func (m *MockWalletRepositoryPeer) Delete(id string) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(id)
	}
	
	if data, exists := m.data[id]; exists {
		// Remove from data
		delete(m.data, id)
		
		// Remove from user data index
		userWallets := m.userData[data.UserID]
		for i, wallet := range userWallets {
			if wallet.ID == id {
				m.userData[data.UserID] = append(userWallets[:i], userWallets[i+1:]...)
				break
			}
		}
	}
	
	return nil
}

func TestWalletRepositoryImpl_Save(t *testing.T) {
	// Arrange
	mockPeer := NewMockWalletRepositoryPeer()
	repo := repository.NewWalletRepositoryImpl(mockPeer)
	
	wallet, err := model.NewWallet("test-user", "Test Wallet", model.WalletTypeCash, "USD")
	if err != nil {
		t.Fatalf("Failed to create test wallet: %v", err)
	}

	// Act
	err = repo.Save(wallet)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	// Verify data was saved
	if len(mockPeer.data) != 1 {
		t.Errorf("Expected 1 wallet in data, got %d", len(mockPeer.data))
	}
	
	savedData, exists := mockPeer.data[wallet.ID]
	if !exists {
		t.Error("Expected wallet to be saved in data")
	}
	
	if savedData.Name != "Test Wallet" {
		t.Errorf("Expected saved wallet name to be 'Test Wallet', got %s", savedData.Name)
	}
}

func TestWalletRepositoryImpl_FindByID(t *testing.T) {
	// Arrange
	mockPeer := NewMockWalletRepositoryPeer()
	repo := repository.NewWalletRepositoryImpl(mockPeer)
	
	// Create test data
	now := time.Now()
	testData := mapper.WalletData{
		ID:              "test-wallet-id",
		UserID:          "test-user",
		Name:            "Test Wallet",
		Type:            "CASH",
		Currency:        "USD",
		BalanceAmount:   10000, // $100.00
		BalanceCurrency: "USD",
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	mockPeer.data["test-wallet-id"] = testData

	// Act
	wallet, err := repo.FindByID("test-wallet-id")

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	if wallet == nil {
		t.Error("Expected wallet to be found, got nil")
		return
	}
	
	if wallet.ID != "test-wallet-id" {
		t.Errorf("Expected wallet ID to be 'test-wallet-id', got %s", wallet.ID)
	}
	
	if wallet.Name != "Test Wallet" {
		t.Errorf("Expected wallet name to be 'Test Wallet', got %s", wallet.Name)
	}
}

func TestWalletRepositoryImpl_FindByID_NotFound(t *testing.T) {
	// Arrange
	mockPeer := NewMockWalletRepositoryPeer()
	repo := repository.NewWalletRepositoryImpl(mockPeer)

	// Act
	wallet, err := repo.FindByID("non-existent-id")

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	if wallet != nil {
		t.Error("Expected wallet to be nil, got wallet")
	}
}

func TestWalletRepositoryImpl_FindByUserID(t *testing.T) {
	// Arrange
	mockPeer := NewMockWalletRepositoryPeer()
	repo := repository.NewWalletRepositoryImpl(mockPeer)
	
	// Create test data
	now := time.Now()
	testData1 := mapper.WalletData{
		ID:              "wallet-1",
		UserID:          "test-user",
		Name:            "Wallet 1",
		Type:            "CASH",
		Currency:        "USD",
		BalanceAmount:   10000,
		BalanceCurrency: "USD",
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	testData2 := mapper.WalletData{
		ID:              "wallet-2",
		UserID:          "test-user",
		Name:            "Wallet 2",
		Type:            "BANK",
		Currency:        "USD",
		BalanceAmount:   20000,
		BalanceCurrency: "USD",
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	
	mockPeer.userData["test-user"] = []mapper.WalletData{testData1, testData2}

	// Act
	wallets, err := repo.FindByUserID("test-user")

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	if len(wallets) != 2 {
		t.Errorf("Expected 2 wallets, got %d", len(wallets))
		return
	}
	
	// Verify wallet details
	if wallets[0].Name != "Wallet 1" {
		t.Errorf("Expected first wallet name to be 'Wallet 1', got %s", wallets[0].Name)
	}
	
	if wallets[1].Name != "Wallet 2" {
		t.Errorf("Expected second wallet name to be 'Wallet 2', got %s", wallets[1].Name)
	}
}

func TestWalletRepositoryImpl_FindByUserID_NoWallets(t *testing.T) {
	// Arrange
	mockPeer := NewMockWalletRepositoryPeer()
	repo := repository.NewWalletRepositoryImpl(mockPeer)

	// Act
	wallets, err := repo.FindByUserID("non-existent-user")

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	if len(wallets) != 0 {
		t.Errorf("Expected 0 wallets, got %d", len(wallets))
	}
}

func TestWalletRepositoryImpl_FindByIDWithTransactions(t *testing.T) {
	// Arrange
	mockPeer := NewMockWalletRepositoryPeer()
	repo := repository.NewWalletRepositoryImpl(mockPeer)
	
	// Create test data
	now := time.Now()
	testData := mapper.WalletData{
		ID:              "test-wallet-id",
		UserID:          "test-user",
		Name:            "Test Wallet",
		Type:            "CASH",
		Currency:        "USD",
		BalanceAmount:   10000,
		BalanceCurrency: "USD",
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	mockPeer.data["test-wallet-id"] = testData

	// Act
	wallet, err := repo.FindByIDWithTransactions("test-wallet-id")

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	if wallet == nil {
		t.Error("Expected wallet to be found, got nil")
		return
	}
	
	// Verify wallet is marked as fully loaded
	if !wallet.IsFullyLoaded() {
		t.Error("Expected wallet to be marked as fully loaded")
	}
	
	if wallet.ID != "test-wallet-id" {
		t.Errorf("Expected wallet ID to be 'test-wallet-id', got %s", wallet.ID)
	}
}

func TestWalletRepositoryImpl_Delete(t *testing.T) {
	// Arrange
	mockPeer := NewMockWalletRepositoryPeer()
	repo := repository.NewWalletRepositoryImpl(mockPeer)
	
	// Create test data
	now := time.Now()
	testData := mapper.WalletData{
		ID:              "test-wallet-id",
		UserID:          "test-user",
		Name:            "Test Wallet",
		Type:            "CASH",
		Currency:        "USD",
		BalanceAmount:   10000,
		BalanceCurrency: "USD",
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	mockPeer.data["test-wallet-id"] = testData
	mockPeer.userData["test-user"] = []mapper.WalletData{testData}

	// Act
	err := repo.Delete("test-wallet-id")

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	// Verify data was deleted
	if len(mockPeer.data) != 0 {
		t.Errorf("Expected 0 wallets in data, got %d", len(mockPeer.data))
	}
	
	// Verify user data was updated
	if len(mockPeer.userData["test-user"]) != 0 {
		t.Errorf("Expected 0 wallets for user, got %d", len(mockPeer.userData["test-user"]))
	}
}