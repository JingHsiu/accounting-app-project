# Implementation Guidelines: Clean Architecture Aggregate Persistence

## Clean Architecture Compliance

### Layer Separation Principles
```
Domain Layer (Core)
├── Wallet (Aggregate Root) ✅
├── IncomeRecord (Entity) ✅  
├── ExpenseRecord (Entity) ✅
└── Transfer (Entity) ✅

Application Layer (Use Cases)
├── WalletRepository Interface ✅
├── Enhanced WalletMapper ✅
└── AddIncomeService (simplified) ✅

Adapter Layer (Interface Adapters)  
├── WalletRepositoryPeer Interface ✅
├── PostgreSQL Peer Implementation ✅
└── Database Transaction Management ✅

Infrastructure Layer (Frameworks & Drivers)
├── PostgreSQL Database ✅
├── Multiple Tables (wallets, income_records) ✅
└── Connection Pool Management ✅
```

### Dependency Inversion
```go
// ✅ CORRECT: Application depends on abstractions
type AddIncomeService struct {
    walletRepo   repository.WalletRepository      // Interface (abstraction)
    categoryRepo repository.IncomeCategoryRepository // Interface
}

// ❌ WRONG: Application depends on concrete implementations  
type AddIncomeService struct {
    walletRepo   *repository.WalletRepositoryImpl    // Concrete class
    categoryRepo *repository.IncomeCategoryRepoImpl  // Concrete class
}
```

## Performance Considerations

### Loading Strategies

#### 1. Partial Loading (Fast Path)
```go
// For balance queries, category validation, quick operations
wallet, err := walletRepo.FindByID(walletID)
// Loads: Wallet metadata only (~100 bytes)
// Performance: ~1-2ms database query
// Memory: Minimal footprint
```

#### 2. Complete Loading (Full Path) 
```go
// For income operations, comprehensive analysis
wallet, err := walletRepo.FindByIDWithTransactions(walletID)
// Loads: Wallet + all IncomeRecords + ExpenseRecords + Transfers
// Performance: ~10-50ms depending on transaction count  
// Memory: Proportional to transaction history size
```

#### 3. Selective Loading (Optimized Path)
```go
// Future enhancement for date range queries
wallet, err := walletRepo.FindByIDWithRecentTransactions(walletID, 30) // Last 30 days
// Loads: Wallet + recent transactions only
// Performance: ~5-15ms optimized query
// Memory: Bounded by date range
```

### Memory Management

#### Aggregate Size Guidelines
```
Small Aggregate:   < 100 transactions  (~10KB memory)
Medium Aggregate:  100-1000 transactions (~100KB memory)  
Large Aggregate:   1000+ transactions (>1MB memory)

Recommendation: Monitor and implement pagination for large aggregates
```

#### Memory Optimization Patterns
```go
// 1. Lazy Loading Pattern
type Wallet struct {
    // Core fields loaded immediately
    ID, UserID, Name, Type, Balance // ~200 bytes
    
    // Child entities loaded on demand
    incomeRecords  []IncomeRecord // Loaded only when requested
    expenseRecords []ExpenseRecord
    transfers      []Transfer
    
    isFullyLoaded bool // Track loading state
}

// 2. Streaming Pattern (Future Enhancement)
func (r *WalletRepository) StreamTransactions(walletID string, handler func(Transaction) error) error {
    // Stream large transaction sets without loading all into memory
}
```

### Database Performance

#### Query Optimization
```sql
-- Optimized aggregate loading query
SELECT 
    w.*,
    ir.id as income_id, ir.amount as income_amount, ir.description as income_desc,
    er.id as expense_id, er.amount as expense_amount, er.description as expense_desc,
    t.id as transfer_id, t.amount as transfer_amount, t.description as transfer_desc
FROM wallets w
LEFT JOIN income_records ir ON w.id = ir.wallet_id  
LEFT JOIN expense_records er ON w.id = er.wallet_id
LEFT JOIN transfers t ON w.id = t.from_wallet_id OR w.id = t.to_wallet_id
WHERE w.id = $1
ORDER BY ir.date DESC, er.date DESC, t.date DESC;

-- Index strategy
CREATE INDEX idx_income_records_wallet_date ON income_records(wallet_id, date DESC);
CREATE INDEX idx_expense_records_wallet_date ON expense_records(wallet_id, date DESC); 
CREATE INDEX idx_transfers_wallet_date ON transfers(from_wallet_id, date DESC);
```

#### Transaction Boundaries
```go
// Single aggregate save = single database transaction
func (p *PostgreSQLPeer) SaveAggregate(wallet WalletData, children ChildEntitiesData) error {
    return p.db.WithTransaction(func(tx Transaction) error {
        // 1. Update wallet (balance, metadata)
        if err := p.updateWallet(tx, wallet); err != nil {
            return err // Automatic rollback
        }
        
        // 2. Handle child entity changes
        if err := p.syncChildEntities(tx, wallet.ID, children); err != nil {
            return err // Automatic rollback  
        }
        
        return nil // Commit all changes
    })
}
```

## Testing Strategy

### 1. Unit Testing

#### Domain Layer Tests
```go
func TestWallet_AddIncome_UpdatesBalanceAndRecords(t *testing.T) {
    // Test aggregate behavior
    wallet := createTestWallet()
    income, err := wallet.AddIncome(money, categoryID, "test income", time.Now())
    
    assert.NoError(t, err)
    assert.Equal(t, expectedBalance, wallet.Balance)
    assert.Len(t, wallet.GetIncomeRecords(), 1)
    assert.Equal(t, income.ID, wallet.GetIncomeRecords()[0].ID)
}
```

#### Mapper Tests  
```go
func TestWalletMapper_ToAggregateData_IncludesChildEntities(t *testing.T) {
    // Test enhanced mapper
    wallet := createWalletWithIncomeRecords()
    
    walletData, childEntities := mapper.ToAggregateData(wallet)
    
    assert.Equal(t, wallet.ID, walletData.ID)
    assert.Len(t, childEntities.IncomeRecords, 2)
}
```

### 2. Integration Testing

#### Repository Integration Tests
```go
func TestWalletRepository_SaveAndLoad_PreservesAggregateState(t *testing.T) {
    // Test complete aggregate persistence
    wallet := createWalletWithTransactions()
    
    // Save complete aggregate
    err := repo.Save(wallet)
    assert.NoError(t, err)
    
    // Load complete aggregate  
    loaded, err := repo.FindByIDWithTransactions(wallet.ID)
    assert.NoError(t, err)
    assert.Equal(t, len(wallet.GetIncomeRecords()), len(loaded.GetIncomeRecords()))
    assert.True(t, loaded.IsFullyLoaded())
}
```

#### Service Integration Tests
```go
func TestAddIncomeService_Execute_PersistsToWalletOnly(t *testing.T) {
    // Test single repository pattern
    service := NewAddIncomeService(walletRepo, categoryRepo) // No incomeRecordRepo
    
    result := service.Execute(addIncomeInput)
    
    assert.Equal(t, common.Success, result.ExitCode)
    
    // Verify persistence through wallet repository only
    wallet, _ := walletRepo.FindByIDWithTransactions(walletID)
    assert.Len(t, wallet.GetIncomeRecords(), 1)
}
```

### 3. Performance Testing

#### Load Testing
```go
func BenchmarkWalletRepository_SaveLargeAggregate(b *testing.B) {
    wallet := createWalletWith1000Transactions()
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        err := repo.Save(wallet)
        if err != nil {
            b.Fatal(err)
        }
    }
}
```

#### Memory Profiling
```go
func TestWalletRepository_MemoryUsage_WithinLimits(t *testing.T) {
    var memStart, memEnd runtime.MemStats
    runtime.GC()
    runtime.ReadMemStats(&memStart)
    
    // Load large aggregate
    wallet, err := repo.FindByIDWithTransactions(largeWalletID)
    runtime.ReadMemStats(&memEnd)
    
    memoryUsed := memEnd.Alloc - memStart.Alloc
    assert.Less(t, memoryUsed, uint64(10*1024*1024)) // Max 10MB per aggregate
}
```

## Error Handling & Resilience

### Transaction Error Handling
```go
func (p *PostgreSQLPeer) SaveAggregate(wallet WalletData, children ChildEntitiesData) error {
    return p.db.WithTransaction(func(tx Transaction) error {
        defer func() {
            if r := recover(); r != nil {
                // Log panic and convert to error
                log.Errorf("Panic during aggregate save: %v", r)
                // Transaction automatically rolled back
            }
        }()
        
        // Critical section - all operations must succeed
        if err := p.validateWalletData(wallet); err != nil {
            return fmt.Errorf("invalid wallet data: %w", err)
        }
        
        if err := p.saveWalletData(tx, wallet); err != nil {
            return fmt.Errorf("failed to save wallet: %w", err)
        }
        
        if err := p.saveChildEntities(tx, wallet.ID, children); err != nil {
            return fmt.Errorf("failed to save child entities: %w", err)
        }
        
        return nil
    })
}
```

### Consistency Validation
```go
func (r *WalletRepositoryImpl) Save(wallet *model.Wallet) error {
    // Pre-save validation
    if err := r.validateAggregateConsistency(wallet); err != nil {
        return fmt.Errorf("aggregate consistency check failed: %w", err)
    }
    
    // Attempt save
    err := r.peer.SaveAggregate(walletData, childEntities)
    if err != nil {
        // Log error and return wrapped error
        log.Errorf("Failed to save wallet aggregate %s: %v", wallet.ID, err)
        return fmt.Errorf("repository save failed: %w", err)
    }
    
    return nil
}

func (r *WalletRepositoryImpl) validateAggregateConsistency(wallet *model.Wallet) error {
    // Validate balance consistency
    calculatedBalance := r.calculateBalanceFromTransactions(wallet)
    if !calculatedBalance.Equals(wallet.Balance) {
        return errors.New("wallet balance inconsistent with transaction history")
    }
    
    // Validate entity relationships  
    for _, income := range wallet.GetIncomeRecords() {
        if income.WalletID != wallet.ID {
            return errors.New("income record references wrong wallet")
        }
    }
    
    return nil
}
```

## Monitoring & Observability

### Performance Metrics
```go
// Repository layer metrics
var (
    aggregateSaveTime = prometheus.NewHistogramVec(prometheus.HistogramOpts{
        Name: "wallet_aggregate_save_duration_seconds",
        Help: "Time taken to save wallet aggregate",
    }, []string{"operation", "fully_loaded"})
    
    aggregateSize = prometheus.NewHistogramVec(prometheus.HistogramOpts{
        Name: "wallet_aggregate_size_bytes",
        Help: "Size of wallet aggregate in memory",  
    }, []string{"transaction_count_range"})
)

func (r *WalletRepositoryImpl) Save(wallet *model.Wallet) error {
    start := time.Now()
    defer func() {
        aggregateSaveTime.WithLabelValues("save", 
            strconv.FormatBool(wallet.IsFullyLoaded())).Observe(time.Since(start).Seconds())
    }()
    
    // ... save logic
}
```

### Health Checks
```go
func (r *WalletRepositoryImpl) HealthCheck() error {
    // Test basic connectivity
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    testWallet, err := r.FindByID("health-check-wallet")
    if err != nil && !errors.Is(err, ErrWalletNotFound) {
        return fmt.Errorf("wallet repository health check failed: %w", err)
    }
    
    return nil
}
```