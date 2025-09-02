# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

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

### Frontend (React/TypeScript)
Navigate to `../frontend` first, then:
```bash
npm install      # Install dependencies
npm run dev      # Start development server
npm run build    # Build for production
npm run lint     # Run linting
npm run preview  # Preview production build
```

### Database
```bash
# Start PostgreSQL with Docker
docker-compose up -d postgres

# Start with pgAdmin management interface
docker-compose --profile admin up -d

# Reset database (removes all data)
docker-compose down -v
```

## Architecture Overview

This is a **Clean Architecture** implementation with **Domain-Driven Design (DDD)** principles for a personal accounting application.

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

**Money Value Object**: All monetary values use the `Money` type with integer amounts (to avoid floating-point precision issues) and currency codes.

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

All monetary amounts stored as BIGINT (cents/smallest currency unit) with separate currency fields.

### Dependency Injection

The `main.go` demonstrates proper dependency injection flow:
1. Create infrastructure layer components (database, stores)
2. Create adapter layer components (repository peers)  
3. Create application layer components (repositories, use cases)
4. Create adapter layer components (controllers)
5. Wire everything together through constructor injection

Environment variables:
- `DATABASE_URL` - PostgreSQL connection (default: postgres://postgres:password@localhost:5432/accountingdb?sslmode=disable)
- `PORT` - HTTP server port (default: 8080)