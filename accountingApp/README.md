# Accounting App Backend

A Clean Architecture implementation of a personal accounting application backend built with Go, featuring Domain-Driven Design principles and CQRS patterns.

## 🏗️ Architecture Overview

This backend implements **Clean Architecture** with **4 distinct layers**, ensuring strict dependency inversion and maintainable separation of concerns.

### Layer Structure

```
cmd/accountingApp/           # Application Entry Point
├── main.go                  # Dependency injection & server setup

internal/accounting/         # Core Application Modules
├── domain/                  # Layer 1: Business Logic (Core)
├── application/             # Layer 2: Use Cases & Application Services
├── adapter/                 # Layer 3: Interface Adapters
├── frameworks/              # Layer 4: Infrastructure & External Systems
└── test/                    # Comprehensive Test Suite
```

## 📋 Component Overview

### 🎯 Domain Layer (Layer 1)
**Pure business logic with no external dependencies**

**Domain Models** (`domain/model/`)
- `wallet.go` - Wallet aggregate root with transaction history
- `money.go` - Money value object with currency validation
- `expenseCategory.go` / `incomeCategory.go` - Hierarchical category system
- `expenseRecord.go` / `incomeRecord.go` - Transaction entities

**Domain Services** (`domain/service/`)
- `CategoryValidationService.go` - Business rule validation for categories

### ⚙️ Application Layer (Layer 2)
**Orchestrates business workflows and defines contracts**

**Use Case Interfaces** (`application/usecase/`)
- `usecases.go` - Contract definitions for all business operations
- Input/Output models with standardized response patterns

**Command Services** (`application/command/`) - Write Operations
- `CreateWalletService.go` - Wallet creation with optional initial balance
- `UpdateWalletService.go` - Wallet property modifications
- `DeleteWalletService.go` - Safe wallet deletion
- `AddExpenseService.go` / `AddIncomeService.go` - Transaction recording
- `CreateExpenseCategoryService.go` / `CreateIncomeCategoryService.go` - Category management
- `ProcessTransferService.go` - Inter-wallet transfers with fees

**Query Services** (`application/query/`) - Read Operations
- `GetWalletsService.go` - User wallet listing
- `GetWalletService.go` - Single wallet retrieval with optional transactions
- `GetWalletBalanceService.go` - Wallet balance queries

**Repository Layer** (`application/repository/`)
- `Repository.go` - Generic repository interfaces
- `WalletRepositoryImpl.go` - Wallet repository implementation using Bridge pattern

**Data Mapping** (`application/mapper/`)
- `WalletMapper.go` - Domain ↔ Data transformation
- `CategoryMapper.go` - Category data mapping
- `Mapper.go` - Base mapping interfaces

### 🔌 Adapter Layer (Layer 3)  
**Interface adapters connecting use cases to external systems**

**Specialized Controllers** (`adapter/controller/`)
- `createWalletController.go` - POST /api/v1/wallets
- `updateWalletController.go` - PUT /api/v1/wallets/{id}
- `deleteWalletController.go` - DELETE /api/v1/wallets/{id}
- `getWalletController.go` - GET /api/v1/wallets, GET /api/v1/wallets/{id}
- `getWalletBalanceController.go` - GET /api/v1/wallets/{id}/balance
- `addExpenseController.go` - POST /api/v1/expenses
- `addIncomeController.go` - POST /api/v1/incomes
- `categoryController.go` - Category management endpoints

**Repository Adapters** (`adapter/repository/`)
- `pgRepositoryPeerAdapter.go` - PostgreSQL repository bridge implementation

**Storage Abstractions** (`adapter/store/`)
- `AggregateStore.go` - Generic aggregate persistence interfaces

### 🖥️ Frameworks Layer (Layer 4)
**Infrastructure, databases, and external system integrations**

**Database Systems** (`frameworks/database/`)
- `PostgreSQLClient.go` - PostgreSQL database client implementation
- `PgAggregateStoreAdapter.go` - PostgreSQL-specific aggregate store
- `DatabaseClient.go` - Database client interface
- `connection.go` - Connection management
- `schema.sql` - Database schema with indexes

**Web Framework** (`frameworks/web/`)
- `router.go` - HTTP routing with RESTful API design

## 🚀 API Endpoints

### Wallet Management
```http
GET    /api/v1/wallets?userID={id}     # List user wallets
POST   /api/v1/wallets                 # Create new wallet
GET    /api/v1/wallets/{id}            # Get wallet details
PUT    /api/v1/wallets/{id}            # Update wallet
DELETE /api/v1/wallets/{id}            # Delete wallet
GET    /api/v1/wallets/{id}/balance    # Get wallet balance
```

### Transaction Management  
```http
POST   /api/v1/expenses                # Record expense
POST   /api/v1/incomes                 # Record income
```

### Category Management
```http
POST   /api/v1/categories/expense      # Create expense category
POST   /api/v1/categories/income       # Create income category
```

### Health Check
```http
GET    /health                         # Service health status
```

## 🏛️ Architectural Patterns

### Clean Architecture Principles
- **Dependency Inversion**: All dependencies point inward toward the domain
- **Interface Segregation**: Small, focused interfaces for each responsibility
- **Single Responsibility**: Each component handles one specific concern

### Domain-Driven Design
- **Aggregate Pattern**: Wallet as aggregate root managing contained entities
- **Value Objects**: Money type ensuring currency consistency and precision
- **Domain Services**: Complex business logic separated from entities

### CQRS-Lite Implementation
- **Command/Query Separation**: Distinct handlers for read and write operations
- **Specialized Controllers**: Single-purpose controllers per operation type
- **Standardized Outputs**: Common response patterns across all operations

### Repository Pattern with Bridge
- **Interface Abstraction**: Repository contracts defined in application layer
- **Bridge Implementation**: Peer adapters connecting to infrastructure layer
- **Generic Type Safety**: Go generics for type-safe aggregate operations

## 💾 Data Management

### Money Value Object
- Integer-based amounts (avoiding floating-point precision issues)
- Currency validation (ISO 4217 3-character codes)
- Currency-aware arithmetic operations

### Database Schema
- **PostgreSQL** with proper indexing for performance
- **ACID Transactions** ensuring data consistency  
- **Foreign Key Constraints** maintaining referential integrity
- **Hierarchical Categories** supporting main/sub-category structures

### Aggregate Persistence
- **Generic Aggregate Store** supporting any aggregate type
- **Optimistic Concurrency** through version tracking
- **Query Optimization** with strategic indexes

## 🧪 Testing Strategy

### Test Organization (`test/`)
```
test/
├── domain/          # Unit tests for domain logic
├── usecase/         # Integration tests for business workflows  
├── controller/      # HTTP integration tests
├── repository/      # Database integration tests
└── fake_wallet_repo.go  # Test doubles for isolation
```

### Testing Principles
- **Real Implementations**: Tests use actual services (not mocks) for integration validation
- **Layer-Specific Testing**: Each layer tested independently and in integration
- **Domain Logic Focus**: Comprehensive coverage of business rules and invariants

## 🛠️ Development Commands

```bash
# Build & Run
go build -o build/accountingApp cmd/accoountingApp/main.go
go run cmd/accoountingApp/main.go

# Testing
go test ./...                           # All tests
go test ./internal/accounting/test/...  # Application tests
go test -run TestWalletCreation ./internal/accounting/test/domain/

# Database
docker-compose up -d postgres           # Start PostgreSQL
docker-compose --profile admin up -d    # With pgAdmin
docker-compose down -v                  # Reset database
```

## 🔧 Configuration

### Environment Variables
- `DATABASE_URL` - PostgreSQL connection string
  - Default: `postgres://postgres:password@localhost:5432/accountingdb?sslmode=disable`
- `PORT` - HTTP server port (default: `8080`)

### Database Connection
- **Connection Pooling**: Managed by Go's database/sql package
- **Health Checks**: Built-in database connectivity validation
- **Migration**: Automatic schema setup via Docker initialization

## 🎯 Key Design Decisions

### Why Clean Architecture?
- **Long-term Maintainability**: Business logic isolated from infrastructure changes
- **Testability**: Each layer can be tested independently
- **Flexibility**: Easy to swap infrastructure components without affecting business logic

### Why Specialized Controllers?
- **Single Responsibility**: Each controller handles one HTTP operation
- **Easier Testing**: Focused testing scope per controller
- **Better Error Handling**: Tailored error responses per operation type

### Why Integer-based Money?
- **Precision**: Avoids floating-point arithmetic errors
- **Currency Safety**: Enforces currency matching in operations
- **Domain Modeling**: Reflects real-world monetary precision requirements

This backend provides a robust, scalable foundation for personal accounting applications with clear separation of concerns and comprehensive business logic encapsulation.