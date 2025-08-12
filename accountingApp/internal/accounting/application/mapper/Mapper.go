package mapper

// Mapper 通用資料轉換介面
// 目的：在Use Case層統一處理Domain Model與Persistence Data的轉換
// 避免Domain層與外部技術細節耦合
type Mapper[TDomain any, TData any] interface {
	// ToData 將Domain Model轉換為Persistence Data
	ToData(domain TDomain) TData
	
	// ToDomain 將Persistence Data轉換為Domain Model
	ToDomain(data TData) (TDomain, error)
}

// AggregateData 通用聚合資料結構介面
// 所有持久化資料結構都應實現此介面
type AggregateData interface {
	GetID() string
}

// MapperRegistry Mapper註冊中心
// 用於管理所有Mapper實例
type MapperRegistry struct {
	mappers map[string]interface{}
}

func NewMapperRegistry() *MapperRegistry {
	return &MapperRegistry{
		mappers: make(map[string]interface{}),
	}
}

func (r *MapperRegistry) Register(name string, mapper interface{}) {
	r.mappers[name] = mapper
}

func (r *MapperRegistry) Get(name string) interface{} {
	return r.mappers[name]
}