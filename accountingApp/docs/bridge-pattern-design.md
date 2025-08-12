# Bridge Pattern Repository架構設計

## 問題背景

原本的設計違反了Clean Architecture的依賴規則：
- 第四層 (Frameworks & Drivers) 的 `PostgresWalletRepository` 直接依賴第一層 (Domain) 的 `model.Wallet`
- 這種直接依賴違反了"外層不應該依賴內層"的原則

## Bridge Pattern解決方案

### 架構層級

```
第一層 (Domain)          │ model.Wallet
                        │      ↑
第二層 (Application)     │ WalletRepository interface
                        │ WalletRepositoryPeer interface  
                        │      ↑
第三層 (Adapter)         │ WalletRepositoryImpl
                        │      ↑ (Bridge)
第四層 (Frameworks)      │ PostgresWalletRepositoryPeer
```

### 介面設計

#### 第二層 - Repository Interfaces

```go
// WalletRepositoryPeer 第四層儲存實現的橋接介面
// 只處理資料結構，不接觸Domain Model
type WalletRepositoryPeer interface {
    SaveData(data mapper.WalletData) error
    FindDataByID(id string) (*mapper.WalletData, error)
    DeleteData(id string) error
}

// WalletRepository 錢包專用儲存庫介面
// 繼承Peer介面，同時提供Domain Model介面
type WalletRepository interface {
    WalletRepositoryPeer // 橋接到第四層
    
    // Domain Model介面
    Save(wallet *model.Wallet) error
    FindByID(id string) (*model.Wallet, error)
    Delete(id string) error
}
```

#### 第三層 - Repository Implementation (Bridge)

```go
type WalletRepositoryImpl struct {
    peer   repository.WalletRepositoryPeer // 橋接到第四層
    mapper *mapper.WalletMapper            // Domain/Data轉換
}

func (r *WalletRepositoryImpl) Save(wallet *model.Wallet) error {
    // Domain Model → Data結構
    data := r.mapper.ToData(wallet)
    
    // 透過peer介面橋接到第四層
    return r.peer.SaveData(data)
}
```

#### 第四層 - Database Implementation

```go
type PostgresWalletRepositoryPeer struct {
    db *sql.DB
}

// 只實現Peer介面，不接觸Domain Model
func (r *PostgresWalletRepositoryPeer) SaveData(data mapper.WalletData) error {
    // 直接操作資料結構，不依賴Domain Model
    query := `INSERT INTO wallets (...) VALUES (...)`
    _, err := r.db.Exec(query, data.ID, data.UserID, ...)
    return err
}
```

### 依賴注入更新

```go
// main.go
// 第四層 - Peer實現
walletPeer := database.NewPostgresWalletRepositoryPeer(dbConn.GetDB())

// 第三層 - Bridge實現  
var walletRepo repository.WalletRepository = adapterRepo.NewWalletRepositoryImpl(walletPeer)
```

## 架構優點

### ✅ 符合Clean Architecture依賴規則
- 第四層只依賴第二層的Peer介面和mapper資料結構
- 不再直接依賴Domain Model
- 依賴方向：外層→內層

### ✅ 責任分離
- **第四層**: 純資料存取，處理SQL和資料結構
- **第三層**: Bridge橋接，處理Domain/Data轉換
- **第二層**: 介面定義，定義契約

### ✅ 可測試性
- 第四層可以獨立測試資料存取邏輯
- 第三層可以用Mock Peer測試
- Domain邏輯與存取邏輯完全分離

### ✅ 可擴展性
- 可以輕鬆替換不同的資料庫實現
- Peer介面支援多種存取方式
- Bridge Pattern支援複雜的轉換邏輯

## 適用範圍

這個Bridge Pattern設計已經應用到：
- ✅ WalletRepository
- ✅ ExpenseCategoryRepository  
- ✅ IncomeCategoryRepository

所有Repository都遵循相同的Bridge Pattern架構。

## 測試驗證

- ✅ 所有現有測試通過
- ✅ Mock repositories已更新支援Peer介面
- ✅ 編譯和運行時都正常工作
- ✅ Clean Architecture依賴規則得到遵守

## 結論

Bridge Pattern成功解決了Clean Architecture違規問題，同時保持了：
- 代碼的可維護性
- 測試的便利性  
- 架構的清晰性
- 未來的擴展性

這是一個符合SOLID原則和Clean Architecture的優雅解決方案。