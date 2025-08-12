# Dependency Flow Documentation

## Architecture Overview

The application now follows proper Clean Architecture with Dependency Inversion Principle:

```
Browser → Controller (Adapter) → UseCase Interface → UseCase Service → Domain Model → Repository Interface → Repository Implementation
```

## Layer Dependencies

### 1. Adapter Layer (Controllers)
- **Depends on**: Application Layer Interfaces (`usecase` package)
- **Does NOT depend on**: Concrete service implementations
- **Location**: `/internal/accounting/adapter/controller/`

### 2. Application Layer
- **Interfaces**: `/internal/accounting/application/usecase/interfaces.go`
- **Services**: `/internal/accounting/application/command/` and `/internal/accounting/application/query/`
- **Depends on**: Domain models and Repository interfaces
- **Does NOT depend on**: Frameworks or external libraries

### 3. Domain Layer
- **Location**: `/internal/accounting/domain/model/`
- **Depends on**: Nothing (pure business logic)

### 4. Frameworks & Drivers Layer
- **Location**: `/internal/accounting/frameworks/`
- **Depends on**: Application layer interfaces and Domain models
- **Implements**: Repository interfaces from Application layer

## Interface Definitions

### Command Use Cases
- `CreateWalletUseCase` - Interface for wallet creation
- `AddExpenseUseCase` - Interface for expense recording
- `AddIncomeUseCase` - Interface for income recording  
- `CreateExpenseCategoryUseCase` - Interface for expense category creation
- `CreateIncomeCategoryUseCase` - Interface for income category creation

### Query Use Cases
- `GetWalletBalanceUseCase` - Interface for balance queries

## Benefits

1. **Testability**: Controllers can be tested with mock implementations
2. **Flexibility**: Easy to swap implementations without changing controllers
3. **Dependency Inversion**: High-level modules don't depend on low-level modules
4. **Clean Architecture Compliance**: Proper dependency direction (inward)

## Example Usage

```go
// In main.go - Dependency Injection
var createWalletUseCase usecase.CreateWalletUseCase = createWalletService
controller := NewWalletController(createWalletUseCase, ...)

// In controller - Using interface
output := c.createWalletUseCase.Execute(input)
```