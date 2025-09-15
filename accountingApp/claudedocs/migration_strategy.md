# Migration Strategy: From Dual-Repository to Single Aggregate Repository

## Current State Analysis

### Problems with Current Implementation
1. **Dual Repository Anti-Pattern**: Both `WalletRepository` and `IncomeRecordRepository` exist
2. **Broken Aggregate Boundaries**: Income records accessed directly, bypassing Wallet
3. **Inconsistent Persistence**: Manual coordination between repositories in `AddIncomeService`
4. **Transaction Isolation**: No single transaction boundary for aggregate operations

### Current Service Implementation Issue
```go
// AddIncomeService.Execute() - PROBLEMATIC PATTERN
func (s *AddIncomeService) Execute(input usecase.AddIncomeInput) common.Output {
    // ...
    income, err := wallet.AddIncome(*amount, input.SubcategoryID, input.Description, input.Date)
    // ...
    
    // 5. Persist wallet aggregate - GOOD
    err = s.walletRepo.Save(wallet)
    
    // 6. Separately persist IncomeRecord - BAD (violates DDD)
    if s.incomeRecordRepo != nil {
        err = s.incomeRecordRepo.Save(income)
    }
}
```

## Migration Phases

### Phase 1: Enhance WalletMapper (Low Risk)
**Goal**: Add child entity handling without changing interfaces

**Changes**:
- Extend `WalletMapper.ToData()` to include detection of loaded entities
- Add `WalletMapper.ToDomain()` support for child entities
- Keep existing interfaces intact

**Benefits**:
- Foundation for enhanced persistence
- No breaking changes
- Backward compatible

### Phase 2: Enhance WalletRepositoryPeer (Medium Risk)
**Goal**: Add aggregate persistence methods to peer layer

**Changes**:
- Add `SaveAggregate()` method to `WalletRepositoryPeer`
- Add `FindAggregateByID()` method for complete loading
- Implement transactional child entity management

**Migration Steps**:
```go
// 1. Extend WalletRepositoryPeer interface
type WalletRepositoryPeer interface {
    // Existing methods (keep for compatibility)
    Save(data mapper.WalletData) error
    FindByID(id string) (*mapper.WalletData, error)
    
    // New aggregate methods
    SaveAggregate(walletData mapper.WalletData, childEntities ChildEntitiesData) error
    FindAggregateByID(id string) (*mapper.WalletData, *ChildEntitiesData, error)
}

// 2. Update WalletRepositoryImpl.Save() to use SaveAggregate()
func (r *WalletRepositoryImpl) Save(wallet *model.Wallet) error {
    if wallet.IsFullyLoaded() {
        // Use enhanced aggregate persistence
        walletData, childEntities := r.mapper.ToAggregateData(wallet)
        return r.peer.SaveAggregate(walletData, childEntities)
    } else {
        // Use traditional wallet-only persistence
        aggregateData := r.mapper.ToData(wallet)
        return r.peer.Save(aggregateData)
    }
}
```

### Phase 3: Remove IncomeRecordRepository Dependencies (High Risk)
**Goal**: Eliminate dual repository pattern

**Changes**:
- Remove `IncomeRecordRepository` from service constructors
- Update `AddIncomeService` to use only `WalletRepository`
- Remove separate `IncomeRecordRepositoryImpl`

**Updated Service**:
```go
// Enhanced AddIncomeService - Single Repository Pattern
type AddIncomeService struct {
    walletRepo   repository.WalletRepository   // Only repository needed
    categoryRepo repository.IncomeCategoryRepository
    // incomeRecordRepo removed ✅
}

func (s *AddIncomeService) Execute(input usecase.AddIncomeInput) common.Output {
    // 1. Load complete aggregate
    wallet, err := s.walletRepo.FindByIDWithTransactions(input.WalletID)
    
    // 2. Validate category (unchanged)
    // ...
    
    // 3. Add income through aggregate root
    income, err := wallet.AddIncome(*amount, input.SubcategoryID, input.Description, input.Date)
    
    // 4. Single aggregate save - handles wallet + income record
    err = s.walletRepo.Save(wallet)  // ✅ Single transaction
    
    return common.UseCaseOutput{
        ID:       income.ID,
        ExitCode: common.Success,
    }
}
```

### Phase 4: Clean Up Legacy Code (Low Risk)
**Goal**: Remove all traces of dual repository pattern

**Changes**:
- Remove `IncomeRecordRepository` interface and implementation files
- Remove `IncomeRecordRepositoryPeer` and its PostgreSQL implementation
- Update dependency injection configuration
- Remove separate income record mappers if no longer needed

## Risk Assessment & Mitigation

### High Risk Areas
1. **Transaction Boundaries**: Multiple table operations must be atomic
2. **Performance Impact**: Loading complete aggregates may be slower
3. **Memory Usage**: Large transaction histories consume more memory

### Mitigation Strategies

#### 1. Transaction Management
```go
// Use database transactions for aggregate persistence
func (p *PGWalletRepositoryPeer) SaveAggregate(wallet WalletData, children ChildEntitiesData) error {
    return p.db.WithTransaction(func(tx Transaction) error {
        // All operations in single transaction
        if err := p.saveWallet(tx, wallet); err != nil { return err }
        if err := p.deleteOldChildren(tx, wallet.ID); err != nil { return err }
        if err := p.insertNewChildren(tx, wallet.ID, children); err != nil { return err }
        return nil
    })
}
```

#### 2. Performance Optimization
```go
// Dual loading strategy
func (r *WalletRepositoryImpl) FindByID(id string) (*model.Wallet, error) {
    // Fast path: wallet only
    return r.findWalletOnly(id)
}

func (r *WalletRepositoryImpl) FindByIDWithTransactions(id string) (*model.Wallet, error) {
    // Complete path: wallet + all child entities
    return r.findCompleteAggregate(id)
}
```

#### 3. Memory Management
```go
// Pagination for large histories
func (r *WalletRepositoryImpl) FindByIDWithRecentTransactions(id string, days int) (*model.Wallet, error) {
    // Load only recent transactions to limit memory usage
    return r.findAggregateWithDateRange(id, time.Now().AddDate(0, 0, -days), time.Now())
}
```

## Testing Strategy

### 1. Unit Tests
- Test enhanced mappers with child entities
- Test transaction boundary handling
- Test error scenarios (rollback behavior)

### 2. Integration Tests  
- Test complete aggregate persistence workflows
- Test performance with large transaction histories
- Test concurrent access scenarios

### 3. Migration Tests
- Test backward compatibility during migration
- Test data integrity across migration phases
- Test rollback procedures

## Rollback Plan

### If Migration Fails
1. **Phase 1-2**: Simple rollback, interfaces unchanged
2. **Phase 3**: Restore `IncomeRecordRepository` dependencies
3. **Phase 4**: Restore removed files from git history

### Emergency Procedures
- Keep backup of current working implementation
- Implement feature flags for new vs old repository pattern
- Monitor application performance and error rates

## Success Metrics

### Technical Metrics
- **Transaction Consistency**: All aggregate operations atomic
- **Performance**: No significant degradation in wallet operations  
- **Memory Usage**: Reasonable memory footprint for loaded aggregates
- **Code Quality**: Reduced repository coupling, cleaner service layer

### Business Metrics
- **Functionality**: All income operations work correctly
- **Data Integrity**: No lost or corrupted income records
- **User Experience**: No visible impact on application performance