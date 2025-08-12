# 架構決策記錄 (Architecture Decision Record)

> 記錄專案中重要的架構決策，包含背景、決策內容、後果及替代方案。

## ADR-001: 採用 Clean Architecture

**日期**: 專案初期  
**狀態**: ✅ 已採用  
**決策者**: 開發團隊  

### 背景
需要建立一個可維護、可測試、可擴展的會計應用程式架構。

### 決策
採用 Robert C. Martin 的 Clean Architecture 模式，分為四層：
1. Domain Layer (Entities)
2. Application Layer (Use Cases)  
3. Adapter Layer (Interface Adapters)
4. Frameworks Layer (Frameworks & Drivers)

### 理由
- **依賴反轉**: 核心業務邏輯不依賴外部框架
- **可測試性**: 各層可獨立測試
- **可維護性**: 層級分離使程式碼更易維護
- **可擴展性**: 新功能可按層級組織

### 後果
**正面**:
- 業務邏輯與技術實作分離
- 高測試覆蓋率
- 易於理解和維護

**負面**:
- 初期開發速度較慢
- 檔案結構較複雜
- 學習曲線較陡

### 替代方案
1. **MVC 模式**: 較簡單但職責混合
2. **分層架構**: 缺乏依賴反轉
3. **微服務**: 過度複雜，不適合單體應用

---

## ADR-002: 採用 Domain-Driven Design (DDD)

**日期**: 專案初期  
**狀態**: ✅ 已採用  
**決策者**: 開發團隊  

### 背景
會計領域有複雜的業務規則，需要清楚表達領域概念。

### 決策
採用 DDD 戰術模式：
- **Aggregate**: Wallet, ExpenseCategory, IncomeCategory
- **Entity**: Wallet, Category, Subcategory
- **Value Object**: Money, CategoryName
- **Domain Service**: CategoryValidationService
- **Repository Pattern**: 資料存取抽象

### 理由
- **領域表達**: 程式碼直接反映業務概念
- **業務規則集中**: 避免業務邏輯散落各處
- **聚合邊界**: 確保資料一致性
- **通用語言**: 開發者與業務專家使用相同術語

### 後果
**正面**:
- 業務邏輯清晰表達
- 資料一致性保證
- 易於溝通和理解

**負面**:
- 增加抽象層次
- 需要深入理解 DDD 概念

### 替代方案
1. **貧血模型**: 簡單但業務邏輯散落
2. **活動記錄**: 與資料庫耦合度高

---

## ADR-003: 解決依賴違規 - 實作 Bridge Pattern

**日期**: 架構審查期間  
**狀態**: ✅ 已實作  
**決策者**: 開發團隊  

### 背景
發現 `PostgresWalletRepository` (第四層) 直接依賴 `model.Wallet` (第一層)，違反 Clean Architecture 的依賴規則。

### 問題
```go
// 違規: 第四層直接依賴第一層
func (r *PostgresWalletRepository) Save(wallet *model.Wallet) error {
    // 直接使用 Domain Model
}
```

### 決策
實作 Bridge Pattern 解決依賴問題：

1. **在第二層定義 Peer 介面**:
   ```go
   type WalletRepositoryPeer interface {
       SaveData(data mapper.WalletData) error
       FindDataByID(id string) (*mapper.WalletData, error)
       DeleteData(id string) error
   }
   ```

2. **在第二層實作 Bridge**:
   ```go
   type WalletRepositoryImpl struct {
       peer   WalletRepositoryPeer
       mapper *mapper.WalletMapper
   }
   ```

3. **在第三層實作 Peer**:
   ```go
   type PostgresWalletRepositoryPeer struct {
       db *sql.DB
   }
   ```

### 理由
- **符合依賴規則**: 第四層只依賴第二層介面
- **職責分離**: 各層職責更明確
- **可測試性**: 可獨立測試各層實作

### 後果
**正面**:
- 完全符合 Clean Architecture
- 層級職責更清楚
- 更易於單元測試

**負面**:
- 增加程式碼複雜度
- 需要額外的 Mapper 轉換

### 替代方案
1. **忽略違規**: 簡單但不符合架構原則
2. **DTO 模式**: 類似 Bridge 但較不靈活

---

## ADR-004: 層級配置修正

**日期**: 架構審查後期  
**狀態**: ✅ 已修正  
**決策者**: 使用者指正  

### 背景
初期實作將 `RepositoryImpl` 放在第三層 (Adapter)，但使用者指正應該在第二層。

### 原始配置 (錯誤)
```
第二層: Repository Interface
第三層: RepositoryImpl, Peer Interface  
第四層: Peer Implementation
```

### 修正後配置
```
第二層: Repository Interface, RepositoryImpl, Peer Interface, Mapper
第三層: Peer Implementation
第四層: Database Connection
```

### 決策
根據使用者指正「mapper跟repositoryimpl, peer interface都在第二層，peer impl在第三層」進行調整。

### 理由
- **層級職責更清楚**: Application Layer 包含所有應用邏輯
- **依賴關係簡化**: 減少跨層依賴
- **符合 Clean Architecture 精神**: 應用層統一管理業務邏輯

### 後果
**正面**:
- 層級職責更明確
- 依賴關係更清晰
- 更符合 Clean Architecture 原則

**負面**:
- 需要重新組織檔案結構
- 需要更新所有引用

---

## ADR-005: 環境變數管理

**日期**: 開發環境設定期間  
**狀態**: ✅ 已實作  
**決策者**: 開發團隊  

### 背景
需要管理不同環境的配置，特別是資料庫連線字串和伺服器埠號。

### 決策
採用 `.env` 檔案 + `github.com/joho/godotenv` 套件：

```go
// 載入 .env 檔案
if err := godotenv.Load(); err != nil {
    log.Printf("Warning: Error loading .env file: %v", err)
}

// 使用環境變數
config := &AppConfig{
    DatabaseURL: getEnv("DATABASE_URL", defaultURL),
    Port:        getEnv("PORT", "8080"),
}
```

### 理由
- **12-Factor App**: 符合現代應用程式原則
- **安全性**: 敏感資訊不進入版本控制
- **環境分離**: 開發/測試/生產環境獨立配置
- **簡單易用**: .env 檔案格式直觀

### 後果
**正面**:
- 配置管理標準化
- 敏感資訊保護
- 環境切換容易

**負面**:
- 需要維護 .env.example
- 部署時需要額外設定

### 替代方案
1. **命令行參數**: 不適合複雜配置
2. **配置檔案**: 需要額外的解析邏輯
3. **環境變數**: 沒有預設值和檔案管理

---

## ADR-006: Docker 開發環境

**日期**: 開發環境設定期間  
**狀態**: ✅ 已實作  
**決策者**: 開發團隊  

### 背景
需要統一開發環境，簡化資料庫設定和管理。

### 決策
使用 Docker Compose 管理開發環境：

```yaml
services:
  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: accountingdb
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: password
    volumes:
      - ./schema.sql:/docker-entrypoint-initdb.d/01-schema.sql:ro
      - postgres_data:/var/lib/postgresql/data
```

### 理由
- **環境一致性**: 所有開發者使用相同的資料庫版本
- **自動初始化**: 容器啟動時自動建立表格
- **資料持久化**: 使用 Volume 保存資料
- **管理便利**: 一鍵啟動完整開發環境

### 後果
**正面**:
- 快速環境設定
- 團隊開發一致性
- 自動化資料庫管理

**負面**:
- 需要安裝 Docker
- 初次啟動較慢

### 替代方案
1. **本機 PostgreSQL**: 環境差異大
2. **雲端資料庫**: 成本高，需要網路連線
3. **SQLite**: 與生產環境差異大

---

## ADR-007: Repository Pattern 實作

**日期**: 資料層設計期間  
**狀態**: ✅ 已實作  
**決策者**: 開發團隊  

### 背景
需要抽象化資料存取，支援不同的儲存實作。

### 決策
實作完整的 Repository Pattern：

```go
// Domain 介面
type WalletRepository interface {
    Save(wallet *model.Wallet) error
    FindByID(id string) (*model.Wallet, error)
    Delete(id string) error
}

// Bridge 實作
type WalletRepositoryImpl struct {
    peer   WalletRepositoryPeer
    mapper *mapper.WalletMapper
}

// Peer 介面
type WalletRepositoryPeer interface {
    SaveData(data mapper.WalletData) error
    FindDataByID(id string) (*mapper.WalletData, error)
    DeleteData(id string) error
}
```

### 理由
- **依賴反轉**: Use Case 依賴抽象而非具體實作
- **可測試性**: 可使用 Mock Repository
- **可替換性**: 可輕鬆更換儲存方式
- **Bridge Pattern**: 解決依賴違規問題

### 後果
**正面**:
- 高度可測試
- 儲存層可替換
- 符合 Clean Architecture

**負面**:
- 程式碼量增加
- 抽象層次較高

### 替代方案
1. **Active Record**: 簡單但耦合度高
2. **Data Mapper**: 類似但沒有 Bridge Pattern
3. **直接 SQL**: 沒有抽象但簡單

---

## 決策狀態說明

- ✅ **已採用**: 決策已實作並在使用中
- 🟡 **進行中**: 決策正在實作過程中
- ❌ **已廢棄**: 決策已被其他方案取代
- 🔄 **重新評估**: 決策需要重新考慮

---

## 未來決策考慮

### 計劃中的 ADR
1. **ADR-008: 錯誤處理策略** - 統一錯誤處理機制
2. **ADR-009: 日誌記錄策略** - 結構化日誌實作
3. **ADR-010: API 認證機制** - 使用者認證方案
4. **ADR-011: 快取策略** - Redis 快取實作
5. **ADR-012: 事件溯源** - Event Sourcing 導入評估

### 可能的架構演進
- **微服務拆分**: 當功能複雜度增加時
- **CQRS 模式**: 讀寫分離優化
- **Event-Driven Architecture**: 提高系統響應性
- **GraphQL API**: 替代 REST API