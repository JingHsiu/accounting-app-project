package usecase

import (
	"github.com/JingHsiu/accountingApp/internal/accounting/application/command"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/common"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/mapper"
	"github.com/JingHsiu/accountingApp/internal/accounting/application/usecase"
	"github.com/JingHsiu/accountingApp/internal/accounting/domain/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type MockExpenseCategoryRepository struct {
	mock.Mock
}

func (m *MockExpenseCategoryRepository) Save(category *model.ExpenseCategory) error {
	args := m.Called(category)
	return args.Error(0)
}

func (m *MockExpenseCategoryRepository) FindByID(id string) (*model.ExpenseCategory, error) {
	args := m.Called(id)
	return args.Get(0).(*model.ExpenseCategory), args.Error(1)
}

func (m *MockExpenseCategoryRepository) FindByUserID(userID string) ([]*model.ExpenseCategory, error) {
	args := m.Called(userID)
	return args.Get(0).([]*model.ExpenseCategory), args.Error(1)
}

func (m *MockExpenseCategoryRepository) FindBySubcategoryID(subcategoryID string) (*model.ExpenseCategory, error) {
	args := m.Called(subcategoryID)
	return args.Get(0).(*model.ExpenseCategory), args.Error(1)
}

func (m *MockExpenseCategoryRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

// 實現ExpenseCategoryRepositoryPeer介面的Bridge Pattern方法
func (m *MockExpenseCategoryRepository) SaveData(data mapper.ExpenseCategoryData) error {
	args := m.Called(data)
	return args.Error(0)
}

func (m *MockExpenseCategoryRepository) FindDataByID(id string) (*mapper.ExpenseCategoryData, error) {
	args := m.Called(id)
	return args.Get(0).(*mapper.ExpenseCategoryData), args.Error(1)
}

func (m *MockExpenseCategoryRepository) FindDataBySubcategoryID(subcategoryID string) (*mapper.ExpenseCategoryData, error) {
	args := m.Called(subcategoryID)
	return args.Get(0).(*mapper.ExpenseCategoryData), args.Error(1)
}

func (m *MockExpenseCategoryRepository) DeleteData(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestCreateExpenseCategory_Success(t *testing.T) {
	repo := new(MockExpenseCategoryRepository)
	repo.On("Save", mock.AnythingOfType("*model.ExpenseCategory")).Return(nil)

	service := command.NewCreateExpenseCategoryService(repo)
	input := usecase.CreateExpenseCategoryInput{
		UserID: "user-123",
		Name:   "Food & Dining",
	}

	output := service.Execute(input)

	assert.Equal(t, common.Success, output.GetExitCode())
	assert.NotEmpty(t, output.GetID())
	repo.AssertExpectations(t)
}

func TestCreateExpenseCategory_InvalidName(t *testing.T) {
	repo := new(MockExpenseCategoryRepository)

	service := command.NewCreateExpenseCategoryService(repo)
	input := usecase.CreateExpenseCategoryInput{
		UserID: "user-123",
		Name:   "",
	}

	output := service.Execute(input)

	assert.Equal(t, common.Failure, output.GetExitCode())
	assert.Contains(t, output.GetMessage(), "Invalid category name")
	repo.AssertNotCalled(t, "Save")
}
