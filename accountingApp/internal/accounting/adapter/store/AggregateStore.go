package store

// AggregateData represents generic aggregate state data for persistence
// This is the base interface that all aggregate data structures must implement
type AggregateData interface {
	GetID() string
}

// AggregateStore defines generic operations for aggregate state persistence
// This interface sits in Layer 3 (Infrastructure Interface) and abstracts
// persistence operations for any aggregate type using state sourcing pattern
type AggregateStore[T AggregateData] interface {
	// Save persists aggregate state (upsert operation)
	Save(data T) error

	// FindByID retrieves aggregate state by ID
	FindByID(id string) (*T, error)

	// Delete removes aggregate state by ID
	Delete(id string) error
}

// AggregateMapper defines the contract for converting between domain aggregates and data structures
// This interface enables clean separation between domain models and persistence data
type AggregateMapper[DomainType any, DataType AggregateData] interface {
	// ToData converts domain aggregate to data structure for persistence
	ToData(domain DomainType) DataType

	// ToDomain converts data structure to domain aggregate
	ToDomain(data DataType) (DomainType, error)
}

// BatchAggregateStore extends AggregateStore with batch operations for performance optimization
type BatchAggregateStore[T AggregateData] interface {
	AggregateStore[T]
	
	// SaveBatch persists multiple aggregates in a single operation
	SaveBatch(data []T) error
	
	// FindBatch retrieves multiple aggregates by their IDs
	FindBatch(ids []string) ([]T, error)
}

// QueryAggregateStore extends AggregateStore with query capabilities
type QueryAggregateStore[T AggregateData] interface {
	AggregateStore[T]
	
	// FindBy retrieves aggregates matching given criteria
	FindBy(criteria map[string]interface{}) ([]T, error)
	
	// Count returns the number of aggregates matching criteria
	Count(criteria map[string]interface{}) (int64, error)
}