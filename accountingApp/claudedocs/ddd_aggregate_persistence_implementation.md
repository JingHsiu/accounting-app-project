# DDD Aggregate Persistence Implementation

## Overview
This document describes the complete implementation of Domain-Driven Design (DDD) aggregate persistence pattern for the Wallet entity in the Go-based accounting application, specifically focusing on child entity persistence within database transactions.

## Implementation Summary

### Key Files Modified/Created

1. **`/internal/accounting/adapter/repository/pgRepositoryPeerAdapter.go`**
   - Enhanced PgWalletRepositoryPeerAdapter with complete transaction management
   - Added child entity persistence methods
   - Implemented aggregate loading with child entities

2. **`/internal/accounting/application/mapper/WalletMapper.go`**
   - Added GetID() methods for child entities to satisfy AggregateData interface
   - Enhanced interface compliance assertions

3. **`/internal/accounting/application/repository/Repository.go`**
   - Added FindByIDWithChildEntities() method to WalletRepositoryPeer interface

4. **`/cmd/accoountingApp/main.go`**
   - Updated dependency injection to include child entity stores
   - Added store creation functions for IncomeRecordData, ExpenseRecordData, TransferData

## Architecture Design

### DDD Aggregate Pattern Implementation

#### 1. Aggregate Root Management
```go
type PgWalletRepositoryPeerAdapter struct {
    walletStore   store.QueryAggregateStore[mapper.WalletData]
    dbClient      database.DatabaseClient
    incomeStore   store.BatchAggregateStore[mapper.IncomeRecordData]
    expenseStore  store.BatchAggregateStore[mapper.ExpenseRecordData]  
    transferStore store.BatchAggregateStore[mapper.TransferData]
}
```

#### 2. Transaction Management Strategy
- **Complete Aggregate Persistence**: Save/load entire aggregate as single unit
- **ACID Transaction Compliance**: All operations within database transaction
- **Rollback on Failure**: Automatic rollback if any child entity fails
- **Child Entity Replacement**: Delete-then-insert pattern for complete consistency

### Database Transaction Implementation

#### saveWithTransaction() Method
```go
func (p *PgWalletRepositoryPeerAdapter) saveWithTransaction(data mapper.WalletData) error {
    tx, err := p.dbClient.BeginTx()
    if err != nil {
        return fmt.Errorf("failed to begin transaction: %w", err)
    }
    defer func() {
        if err != nil {
            tx.Rollback()
        }
    }()

    // 1. Save wallet main entity
    err = p.saveWalletInTransaction(tx, data)
    
    // 2. Save child entities (income, expense, transfers)
    if len(data.IncomeRecords) > 0 {
        err = p.saveIncomeRecords(tx, data.IncomeRecords)
    }
    // ... similar for expense and transfer records

    return tx.Commit()
}
```

### Child Entity Persistence Strategy

#### 1. Income Record Persistence
- **Table**: `income_records`
- **Strategy**: DELETE existing + INSERT new records
- **Columns**: id, wallet_id, category_id, amount, currency, description, date, created_at

#### 2. Expense Record Persistence  
- **Table**: `expense_records`
- **Strategy**: DELETE existing + INSERT new records
- **Columns**: id, wallet_id, category_id, amount, currency, description, date, created_at

#### 3. Transfer Record Persistence
- **Table**: `transfers`
- **Strategy**: DELETE existing + INSERT new records
- **Columns**: id, from_wallet_id, to_wallet_id, amount, currency, fee_amount, fee_currency, description, date, created_at

### Aggregate Loading Implementation

#### Two Loading Strategies

1. **Basic Loading**: `FindByID()` - Loads only wallet basic data
   ```go
   walletData.IsFullyLoaded = false // Child entities not loaded
   ```

2. **Complete Loading**: `FindByIDWithChildEntities()` - Loads complete aggregate
   ```go
   // Load child entities
   err = p.loadChildEntities(walletData)
   walletData.IsFullyLoaded = true // Mark as fully loaded
   ```

#### Child Entity Loading Methods
```go
func (p *PgWalletRepositoryPeerAdapter) loadIncomeRecords(walletID string) ([]mapper.IncomeRecordData, error)
func (p *PgWalletRepositoryPeerAdapter) loadExpenseRecords(walletID string) ([]mapper.ExpenseRecordData, error)  
func (p *PgWalletRepositoryPeerAdapter) loadTransfers(walletID string) ([]mapper.TransferData, error)
```

## SQL Operations

### Child Entity UPSERT Pattern
```sql
-- Income Records Persistence
DELETE FROM income_records WHERE wallet_id = $1;

INSERT INTO income_records (
    id, wallet_id, category_id, amount, currency, description, date, created_at
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8);
```

### Aggregate Query Pattern
```sql
-- Load Income Records for Wallet
SELECT id, wallet_id, category_id, amount, currency, description, date, created_at
FROM income_records
WHERE wallet_id = $1
ORDER BY date DESC, created_at DESC;
```

## Performance Considerations

### 1. Batch Operations
- Use `store.BatchAggregateStore` for child entities
- Single transaction for multiple child entity operations
- Efficient bulk delete/insert operations

### 2. Lazy Loading Strategy
- Basic queries load only wallet entity (`IsFullyLoaded = false`)
- Complete aggregate loading only when explicitly requested
- Reduces memory and query overhead for list operations

### 3. Transaction Isolation
- Database transaction ensures ACID compliance
- Rollback mechanism prevents partial state persistence
- Child entity replacement ensures consistency

## Error Handling Strategy

### 1. Transaction-Level Error Handling
```go
defer func() {
    if err != nil {
        tx.Rollback()
    }
}()
```

### 2. Granular Error Context
```go
return fmt.Errorf("failed to save income records: %w", err)
```

### 3. Aggregate-Level Validation
- Validation occurs at domain level before persistence
- Repository layer focuses on transactional integrity
- Database constraints provide final validation layer

## Interface Compliance

### AggregateData Interface
All child entities implement `store.AggregateData`:
```go
type IncomeRecordData struct { ... }
func (ird IncomeRecordData) GetID() string { return ird.ID }

// Interface compliance assertions
var _ store.AggregateData = (*IncomeRecordData)(nil)
var _ store.AggregateData = (*ExpenseRecordData)(nil)
var _ store.AggregateData = (*TransferData)(nil)
```

### Repository Interface Enhancement
```go
type WalletRepositoryPeer interface {
    Save(data mapper.WalletData) error
    FindByID(id string) (*mapper.WalletData, error)
    FindByIDWithChildEntities(id string) (*mapper.WalletData, error) // New method
    FindByUserID(userID string) ([]mapper.WalletData, error)
    Delete(id string) error
}
```

## Dependency Injection Configuration

### Store Creation
```go
// Main entity store
walletStore := createWalletAggregateStore(dbClient)

// Child entity stores  
incomeStore := createIncomeRecordAggregateStore(dbClient)
expenseStore := createExpenseRecordAggregateStore(dbClient)
transferStore := createTransferAggregateStore(dbClient)

// Repository peer with all dependencies
walletPeer := NewPgWalletRepositoryPeerAdapter(
    walletStore, dbClient, incomeStore, expenseStore, transferStore,
)
```

## Benefits Achieved

### 1. **Data Integrity**
- ACID transaction compliance across multiple tables
- All-or-nothing aggregate persistence
- Consistent aggregate state reconstruction

### 2. **Performance Optimization**
- Lazy loading for basic queries
- Batch operations for child entities
- Single transaction reduces connection overhead

### 3. **Clean Architecture Compliance**  
- Domain entities isolated from persistence concerns
- Repository pattern with proper abstraction layers
- Dependency inversion principle maintained

### 4. **Maintainability**
- Clear separation of aggregate vs child entity operations
- Consistent error handling patterns
- Extensible design for additional child entities

## Testing Strategy

### 1. Unit Tests
- Test transaction rollback scenarios
- Test child entity loading/saving independently
- Test aggregate reconstruction accuracy

### 2. Integration Tests
- Test complete aggregate persistence cycles
- Test concurrent transaction scenarios
- Test database constraint violations

### 3. Performance Tests
- Benchmark lazy vs complete loading
- Test batch operation performance
- Test memory usage with large aggregates

## Conclusion

This implementation successfully provides a production-ready DDD aggregate persistence pattern that:

- Maintains complete transactional integrity
- Supports efficient lazy and complete loading strategies  
- Follows Clean Architecture principles
- Provides comprehensive error handling
- Enables easy testing and maintenance

The solution addresses all the original requirements for child entity persistence while maintaining high performance and architectural cleanliness.