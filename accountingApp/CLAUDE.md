# Backend Development Context

This file provides Go-specific development guidance for the accounting application backend service.

> 📚 **Project Standards**: See [../CLAUDE.md](../CLAUDE.md) for shared project guidelines and cross-service standards
> 🏗️ **System Architecture**: See [../docs/SYSTEM-ARCHITECTURE.md](../docs/SYSTEM-ARCHITECTURE.md) for complete system documentation
> 🌐 **API Integration**: See [../docs/FRONTEND_INTEGRATION_GUIDE.md](../docs/FRONTEND_INTEGRATION_GUIDE.md) for frontend integration patterns

## Development Commands

### Backend (Go)
```bash
# Build the application
go build -o build/accountingApp cmd/accoountingApp/main.go

# Run the application
go run cmd/accoountingApp/main.go

# Run tests
go test ./...

# Run tests for specific layers
go test ./internal/accounting/test/...
go test ./internal/accounting/test/domain/...
go test ./internal/accounting/test/usecase/...

# Run tests with coverage
go test -cover ./...

# Run a specific test function
go test -run TestWalletCreation ./internal/accounting/test/domain/

# Install/update dependencies
go mod tidy
go mod download
```

### Database Management
```bash
# Start PostgreSQL with Docker
docker-compose up -d postgres

# Start with pgAdmin management interface
docker-compose --profile admin up -d

# Reset database (removes all data)
docker-compose down -v

# Database migrations and schema
psql $DATABASE_URL -f frameworks/database/schema.sql
```

### Go-Specific Development
```bash
# Code formatting and linting
go fmt ./...
go vet ./...

# Dependency vulnerability checking
go list -m all | nancy sleuth

# Performance profiling
go tool pprof -http=:8081 http://localhost:8080/debug/pprof/heap
go tool pprof -http=:8081 http://localhost:8080/debug/pprof/profile?seconds=30
```

## Go Backend Architecture Overview

This backend service implements **Clean Architecture** with **Domain-Driven Design (DDD)** principles, focusing on the accounting domain.

> 📋 **Complete Architecture**: See [../CLAUDE.md](../CLAUDE.md) for system-wide architectural decisions and patterns

### 4-Layer Architecture

**Layer 1: Domain (Core)**
- `domain/model/` - Aggregate roots, entities, value objects (Wallet, Money, Categories)
- `domain/service/` - Domain services implementing complex business rules

**Layer 2: Application (Use Cases)**
- `application/usecase/` - Use case interface contracts
- `application/command/` - Command handlers (write operations)
- `application/query/` - Query handlers (read operations)
- `application/repository/` - Repository implementations using Bridge pattern
- `application/mapper/` - Aggregate mappers (Domain ↔ Data conversion)

**Layer 3: Adapter (Interface Adapters)**
- `adapter/controller/` - Specialized HTTP controllers (single responsibility)
- `adapter/repository/` - Repository peer adapters (bridge to Layer 4)
- `adapter/store/` - Abstract aggregate store interfaces

**Layer 4: Frameworks (Infrastructure)**
- `frameworks/database/` - PostgreSQL clients, connections, schema
- `frameworks/web/` - HTTP routing, middleware

### Architectural Patterns

**Dependency Inversion**: All dependencies point inward. Interfaces defined in inner layers, implementations in outer layers.

**Bridge Pattern**: Repository layer uses "peer" interfaces to bridge between Application layer repository implementations and Infrastructure layer aggregate stores.

**Specialized Controllers**: Each controller handles a single operation type:
- `CreateWalletController`
- `UpdateWalletController` 
- `DeleteWalletController`
- `QueryWalletController`
- `GetWalletBalanceController`

**CQRS-lite**: Clear separation of Command (write) and Query (read) operations with dedicated handlers.

**Aggregate Pattern**: Wallet is the primary aggregate root containing:
- Basic wallet properties (ID, name, type, balance)
- Contained transaction records (expenses, incomes, transfers)
- Domain invariants enforcement

### Key Implementation Details

**Money Value Object**: All monetary values use the `Money` type with integer amounts stored in base units. TWD (default currency) uses 1:1 ratio, while USD/EUR use 100:1 (cents). Currency defaults to TWD when not specified.

**Generic Aggregate Store**: Uses Go generics for type-safe aggregate persistence:
```go
type AggregateStore[T AggregateData] interface {
    Save(data T) error
    FindByID(id string) (*T, error)
    Delete(id string) error
}
```

**Use Case Contracts**: All use cases implement a common interface pattern:
```go
type CreateWalletUseCase interface {
    Execute(input CreateWalletInput) common.Output
}
```

**Testing Strategy**: Tests use real implementations (not mocks) to validate layer integration:
- Domain tests focus on business logic correctness
- Use case tests verify application workflow
- Repository tests validate data persistence
- Controller tests check HTTP integration

### Database Schema

Key relationships:
- `wallets` stores aggregate state with balance tracking
- `expense_categories`/`income_categories` implement hierarchical categorization
- `expense_records`/`income_records` link to wallets and categories
- `transfers` handle wallet-to-wallet movements with fee support

All monetary amounts stored as BIGINT in base units (TWD=1, USD=100) with separate currency fields. TWD is the default currency for new records.

### Dependency Injection

The `main.go` demonstrates proper dependency injection flow:
1. Create infrastructure layer components (database, stores)
2. Create adapter layer components (repository peers)  
3. Create application layer components (repositories, use cases)
4. Create adapter layer components (controllers)
5. Wire everything together through constructor injection

### Go Development Best Practices

**Error Handling Patterns**
```go
// Domain errors with context
func (s *CreateWalletService) Execute(input CreateWalletInput) common.Output {
    if err := input.Validate(); err != nil {
        return common.NewFailureOutput("validation failed", err)
    }
    // Business logic...
}

// Repository error handling
func (r *WalletRepositoryImpl) Save(wallet *model.Wallet) error {
    if err := r.store.Save(wallet.ToData()); err != nil {
        return fmt.Errorf("failed to save wallet %s: %w", wallet.ID, err)
    }
    return nil
}
```

**Dependency Injection Patterns**
```go
// Constructor injection for services
func NewCreateWalletService(repo repository.WalletRepository) usecase.CreateWalletUseCase {
    return &CreateWalletService{walletRepo: repo}
}

// Interface-based dependency inversion
type WalletRepository interface {
    Save(wallet *model.Wallet) error
    FindByID(id string) (*model.Wallet, error)
    FindByUserID(userID string) ([]*model.Wallet, error)
}
```

**Testing Patterns with Real Implementations**
```go
// Integration tests use real services
func TestCreateWalletIntegration(t *testing.T) {
    db := testSetupDatabase()
    store := database.NewPgAggregateStoreAdapter[model.WalletData](db)
    repo := repository.NewWalletRepositoryImpl(store)
    service := command.NewCreateWalletService(repo)
    
    // Test with actual implementations
    output := service.Execute(createValidWalletInput())
    assert.True(t, output.Success)
}
```

### Environment Configuration
- `DATABASE_URL` - PostgreSQL connection (default: postgres://postgres:password@localhost:5432/accountingdb?sslmode=disable)
- `PORT` - HTTP server port (default: 8080)
- `GO_ENV` - Environment mode (development, production)
- `LOG_LEVEL` - Logging level (debug, info, warn, error)

## Backend Development Context

### Domain Modeling Guidelines

**Aggregate Root Patterns** (Wallet example)
- Single entry point for business operations
- Encapsulate business invariants and rules
- Handle complex domain operations internally
- Emit domain events for cross-aggregate communication

**Value Object Patterns** (Money example)
- Immutable objects representing descriptive aspects
- Validate construction parameters
- Implement equals and hash methods
- Encapsulate domain-specific operations

**Domain Service Patterns**
- Operations spanning multiple aggregates
- Complex business rule implementations
- Stateless service objects
- Domain-specific validation logic

### Performance Optimization

**Database Query Optimization**
```go
// Use indexes for frequent queries
func (r *WalletRepositoryImpl) FindByUserIDOptimized(userID string) ([]*model.Wallet, error) {
    // Leverages index on user_id column
    query := `SELECT * FROM wallets WHERE user_id = $1 ORDER BY created_at DESC`
    // Implementation...
}

// Batch operations for better performance
func (r *WalletRepositoryImpl) SaveBatch(wallets []*model.Wallet) error {
    // Use database transactions for batch operations
}
```

**Memory Management**
- Use connection pooling for database connections
- Implement proper error handling to prevent resource leaks
- Profile memory usage in production environments
- Use Go's built-in garbage collector efficiently

### Security Considerations

**Input Validation**
```go
// Validate at service boundary
func (input CreateWalletInput) Validate() error {
    if input.UserID == "" {
        return errors.New("user_id is required")
    }
    if input.Currency != "USD" && input.Currency != "EUR" && input.Currency != "TWD" {
        return errors.New("unsupported currency")
    }
    return nil
}
```

**SQL Injection Prevention**
- Always use parameterized queries
- Never concatenate user input into SQL strings
- Validate input types and ranges
- Use prepared statements for repeated queries

---

**Backend Service Context**: Go Clean Architecture + DDD Implementation  
**Layer Focus**: Domain modeling, application services, infrastructure patterns  
**Integration**: RESTful API serving React frontend via standardized contracts