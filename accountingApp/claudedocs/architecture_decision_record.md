# Architecture Decision Record (ADR): Wallet Aggregate Persistence

## Decision Status
**RECOMMENDED** - Enhanced Aggregate Store Pattern (Option 1)

## Context
The current system violates DDD aggregate boundaries by implementing separate `WalletRepository` and `IncomeRecordRepository`, breaking the aggregate root principle where all operations should go through the Wallet aggregate.

## Decision Drivers
1. **DDD Compliance**: Maintain proper aggregate boundaries
2. **Transaction Consistency**: Ensure atomic operations across aggregate + entities  
3. **Performance**: Minimize database queries and memory usage
4. **Maintainability**: Reduce repository coupling and complexity
5. **Migration Safety**: Low-risk transition from current implementation

## Considered Options

### Option 1: Enhanced WalletMapper (RECOMMENDED) ⭐
**Strategy**: Extend existing mapper to handle child entities

**Pros**:
- ✅ Minimal code changes required
- ✅ Leverages existing mapper pattern
- ✅ Maintains Clean Architecture layer separation
- ✅ Low migration risk
- ✅ Backward compatible during transition
- ✅ Single transaction boundary for aggregate operations

**Cons**:
- ⚠️ Mapper becomes more complex
- ⚠️ Requires enhanced peer interface

**Implementation Effort**: Medium (2-3 weeks)
**Risk Level**: Low
**Maintainability**: High

### Option 2: Event-Driven Persistence
**Strategy**: Use domain events for child entity persistence

**Pros**:
- ✅ Excellent decoupling
- ✅ Audit trail capabilities
- ✅ Scalable architecture

**Cons**:
- ❌ High complexity for current needs
- ❌ Eventual consistency challenges
- ❌ Significant architectural changes required
- ❌ Learning curve for team

**Implementation Effort**: High (6-8 weeks)
**Risk Level**: High
**Maintainability**: Medium

### Option 3: Enhanced Aggregate Store  
**Strategy**: Generic pattern for aggregate + entity persistence

**Pros**:
- ✅ Reusable across other aggregates
- ✅ Future-proof design
- ✅ Clean abstraction

**Cons**:
- ❌ Over-engineering for current requirements
- ❌ Complex generic implementation  
- ❌ High development effort
- ❌ May not be needed for other aggregates

**Implementation Effort**: High (4-6 weeks)
**Risk Level**: Medium
**Maintainability**: Medium

## Decision Outcome

**Selected**: Option 1 - Enhanced WalletMapper

### Rationale
1. **Lowest Risk**: Builds on existing patterns and infrastructure
2. **Fastest Delivery**: Can be implemented incrementally
3. **DDD Compliant**: Achieves proper aggregate boundaries
4. **Performance Efficient**: Single transaction, optimized loading strategies
5. **Team Familiarity**: Uses well-understood mapper pattern

## Implementation Plan

### Phase 1: Enhanced Mapper (Week 1)
- Extend `WalletMapper` to handle child entities
- Add `ChildEntitiesData` structure
- Implement `ToAggregateData()` and `ToDomain()` methods
- Unit tests for enhanced mapper

### Phase 2: Enhanced Peer Interface (Week 2)  
- Add `SaveAggregate()` and `FindAggregateByID()` to `WalletRepositoryPeer`
- Implement PostgreSQL peer with transaction management
- Integration tests for aggregate persistence

### Phase 3: Service Layer Updates (Week 3)
- Remove `IncomeRecordRepository` dependency from services
- Update `AddIncomeService` to use single repository pattern
- End-to-end testing of income operations

### Phase 4: Cleanup & Monitoring (Week 4)
- Remove legacy `IncomeRecordRepository` implementation
- Add performance monitoring and health checks
- Documentation and team training

## Success Metrics

### Technical Metrics
- **Transaction Atomicity**: 100% of aggregate operations in single transaction
- **Performance**: No degradation in wallet operation response times
- **Memory Efficiency**: Aggregate loading within 10MB limit per wallet
- **Code Quality**: Reduced cyclomatic complexity in service layer

### Business Metrics  
- **Data Consistency**: Zero data integrity issues post-migration
- **Feature Completeness**: All income operations working correctly
- **System Reliability**: No increase in error rates or downtime

## Consequences

### Positive Consequences
- ✅ **Proper DDD Implementation**: Aggregate boundaries respected
- ✅ **Simplified Service Layer**: Single repository per aggregate
- ✅ **Improved Data Consistency**: Atomic aggregate operations
- ✅ **Better Testability**: Clear aggregate testing boundaries
- ✅ **Reduced Coupling**: Services depend on fewer repositories

### Negative Consequences
- ⚠️ **Increased Mapper Complexity**: More logic in mapping layer
- ⚠️ **Memory Usage**: Complete aggregates consume more memory
- ⚠️ **Query Complexity**: More complex database queries for aggregate loading

### Mitigation Strategies
- **Mapper Complexity**: Comprehensive unit testing and clear documentation
- **Memory Usage**: Implement partial loading strategies and monitoring
- **Query Complexity**: Database query optimization and performance testing

## Monitoring Plan

### Performance Monitoring
```go
// Key metrics to track
- wallet_aggregate_save_duration_seconds (histogram)
- wallet_aggregate_load_duration_seconds (histogram)  
- wallet_aggregate_size_bytes (histogram)
- wallet_transaction_count_per_aggregate (histogram)
```

### Health Checks
- Repository connectivity tests
- Aggregate consistency validation
- Transaction rollback testing
- Memory usage monitoring

## Rollback Plan

### If Issues Arise
1. **Performance Issues**: Enable partial loading by default
2. **Memory Issues**: Implement streaming for large aggregates
3. **Data Integrity Issues**: Rollback to dual repository pattern with feature flag
4. **Critical Failures**: Full rollback using git history and database backups

### Emergency Procedures
- Feature flag to switch between old/new repository patterns
- Database transaction monitoring and alerting
- Automated rollback on critical error thresholds

## Review Schedule

### Implementation Reviews
- **Week 1**: Mapper implementation review
- **Week 2**: Peer interface and transaction handling review  
- **Week 3**: Service layer integration review
- **Week 4**: Performance and monitoring review

### Post-Implementation Reviews
- **1 Month**: Performance impact assessment
- **3 Months**: Long-term stability evaluation
- **6 Months**: Architecture pattern effectiveness review