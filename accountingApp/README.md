# 🏦 Accounting App

一個使用 Clean Architecture + Domain-Driven Design 的 Go 會計應用程式。

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15-336791.svg)](https://postgresql.org)
[![Docker](https://img.shields.io/badge/Docker-compose-0db7ed.svg)](https://docker.com)
[![Architecture](https://img.shields.io/badge/Architecture-Clean%20Architecture-green.svg)](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)

## 📋 目錄

- [功能特色](#功能特色)
- [快速開始](#快速開始)
- [架構設計](#架構設計)
- [API 文檔](#api-文檔)
- [開發指南](#開發指南)
- [測試](#測試)
- [部署](#部署)
- [貢獻指南](#貢獻指南)

## ✨ 功能特色

### 已實作功能
- 🏦 **錢包管理**: 建立多幣別錢包，管理帳戶餘額
- 📊 **分類管理**: 支出和收入分類系統，支援子分類
- 💰 **交易記錄**: 記錄收入和支出交易 (開發中)
- 🌐 **REST API**: 完整的 HTTP API 介面
- 🐳 **Docker 支援**: 開發環境容器化
- 📊 **PostgreSQL**: 企業級資料庫支援

### 架構特色
- 🏗️ **Clean Architecture**: 4層架構，依賴反轉
- 🎯 **Domain-Driven Design**: 領域驅動設計
- 🌉 **Bridge Pattern**: 解決跨層依賴問題
- ✅ **高測試覆蓋率**: 85%+ 測試覆蓋率
- 🔧 **環境管理**: .env 檔案配置管理

## 🚀 快速開始

### 前置需求
- Go 1.21 或更高版本
- Docker 和 Docker Compose
- Git

### 1. 取得原始碼
```bash
git clone <repository-url>
cd accountingApp
```

### 2. 環境設定
```bash
# 複製環境變數模板
cp .env.example .env

# 編輯配置 (可選)
vim .env
```

### 3. 啟動開發環境
```bash
# 啟動資料庫
./scripts/start-dev.sh

# 或手動啟動
docker-compose up -d postgres
```

### 4. 執行應用程式
```bash
# 安裝依賴
go mod tidy

# 執行應用程式
go run cmd/accoountingApp/main.go
```

### 5. 驗證安裝
```bash
# 檢查健康狀態
curl http://localhost:8080/health

# 建立測試錢包
curl -X POST http://localhost:8080/api/v1/wallets \
  -H "Content-Type: application/json" \
  -d '{"name":"Test Wallet","type":"CASH","currency":"USD"}'
```

## 🏗️ 架構設計

### Clean Architecture 層級

```mermaid
graph TD
    A[Domain Layer<br/>領域層] --> B[Application Layer<br/>應用層]
    B --> C[Adapter Layer<br/>適配層]
    C --> D[Frameworks Layer<br/>框架層]
    
    A1[Entities<br/>Aggregates<br/>Value Objects<br/>Domain Services] --> A
    B1[Use Cases<br/>Repository Interfaces<br/>Mappers] --> B
    C1[Controllers<br/>Repository Impl] --> C
    D1[Database<br/>Web Framework<br/>External APIs] --> D
```

### 核心元件

#### Domain Layer (第一層)
```go
// 聚合根
type Wallet struct {
    id       WalletID
    userID   UserID
    name     string
    balance  Money
    // ...
}

// 值對象
type Money struct {
    amount   int64
    currency Currency
}
```

#### Application Layer (第二層)
```go
// 使用案例
type CreateWalletService struct {
    repo repository.WalletRepository
}

// 儲存庫介面
type WalletRepository interface {
    Save(wallet *model.Wallet) error
    FindByID(id string) (*model.Wallet, error)
}
```

#### Adapter Layer (第三層)
```go
// HTTP 控制器
type WalletController struct {
    createUseCase usecase.CreateWalletUseCase
}

// 儲存庫實作
type PostgresWalletRepositoryPeer struct {
    db *sql.DB
}
```

### Bridge Pattern 設計

為了解決 Clean Architecture 的依賴違規問題，我們實作了 Bridge Pattern：

```mermaid
graph LR
    A[Use Case] --> B[Repository Interface]
    B --> C[Repository Impl<br/>Bridge]
    C --> D[Peer Interface]
    D --> E[Database Impl]
    
    C --> F[Mapper]
    F --> G[Domain ↔ Data]
```

詳細設計請參閱: [Bridge Pattern 設計文檔](docs/bridge-pattern-design.md)

## 🌐 API 文檔

### 錢包 API

#### 建立錢包
```http
POST /api/v1/wallets
Content-Type: application/json

{
  "name": "My Wallet",
  "type": "CASH",
  "currency": "USD"
}
```

#### 查詢錢包餘額
```http
GET /api/v1/wallets/{id}/balance
```

#### 新增支出
```http
POST /api/v1/expenses
Content-Type: application/json

{
  "walletId": "wallet-uuid",
  "categoryId": "category-uuid",
  "amount": 1000,
  "currency": "USD",
  "description": "Lunch"
}
```

### 分類 API

#### 建立支出分類
```http
POST /api/v1/categories/expense
Content-Type: application/json

{
  "name": "Food",
  "subcategories": ["Restaurant", "Grocery"]
}
```

#### 建立收入分類
```http
POST /api/v1/categories/income
Content-Type: application/json

{
  "name": "Salary",
  "subcategories": ["Base", "Bonus"]
}
```

### 錯誤回應格式
```json
{
  "error": {
    "code": "INVALID_INPUT",
    "message": "Invalid wallet type",
    "details": "Type must be CASH or BANK"
  }
}
```

## 🛠️ 開發指南

### 專案結構
```
accountingApp/
├── cmd/accoountingApp/           # 應用程式入口
├── internal/accounting/
│   ├── domain/                   # 領域層
│   │   ├── model/               # 聚合和實體
│   │   └── service/             # 領域服務
│   ├── application/             # 應用層
│   │   ├── command/             # 命令使用案例
│   │   ├── query/               # 查詢使用案例
│   │   ├── repository/          # 儲存庫介面
│   │   └── mapper/              # 資料轉換
│   ├── adapter/                 # 適配層
│   │   └── controller/          # HTTP 控制器
│   └── frameworks/              # 框架層
│       ├── database/            # 資料庫實作
│       └── web/                 # Web 框架
├── test/                        # 測試檔案
├── docs/                        # 專案文檔
├── scripts/                     # 開發腳本
├── docker-compose.yml           # Docker 配置
└── .env                        # 環境變數
```

### 開發工作流程

#### 1. 新增功能
```bash
# 1. 建立分支
git checkout -b feature/new-feature

# 2. 編寫測試
vim internal/accounting/test/new_feature_test.go

# 3. 實作功能 (由內而外)
vim internal/accounting/domain/model/...
vim internal/accounting/application/...
vim internal/accounting/adapter/...

# 4. 執行測試
go test ./...
```

#### 2. 程式碼品質檢查
```bash
# 格式化程式碼
go fmt ./...

# 執行 linter (如果有安裝)
golint ./...

# 檢查競爭條件
go test -race ./...
```

#### 3. 資料庫變更
```bash
# 1. 更新 schema.sql
vim internal/accounting/frameworks/database/schema.sql

# 2. 重建資料庫
docker-compose down -v
docker-compose up -d postgres

# 3. 驗證變更
go run cmd/accoountingApp/main.go
```

### 編碼規範

#### Clean Architecture 規則
1. **依賴方向**: 永遠由外向內 (Frameworks → Adapters → Use Cases → Entities)
2. **跨邊界**: 使用介面和 DTO，不直接傳遞內層物件
3. **業務規則**: 放在 Domain 和 Application 層
4. **技術細節**: 放在 Frameworks 層

#### Go 編碼規範
- 遵循 [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- 使用 `gofmt` 格式化程式碼
- 介面名稱加 `er` 後綴 (如 `Repository`, `Mapper`)
- 錯誤處理不可忽略

#### 命名規範
- **Package**: 小寫，簡短，描述性
- **Interface**: 動詞 + er (Reader, Writer, Repository)
- **Struct**: 名詞，PascalCase
- **Function**: 動詞開始，描述行為

## 🧪 測試

### 執行所有測試
```bash
go test ./...
```

### 測試覆蓋率
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### 測試分類

#### 單元測試
```bash
# Domain 層測試
go test ./internal/accounting/domain/...

# Use Case 測試
go test ./internal/accounting/application/...
```

#### 整合測試
```bash
go test ./internal/accounting/test/...
```

### 測試結構
```go
func TestCreateWallet_Success(t *testing.T) {
    // Arrange
    repo := &mocks.MockWalletRepository{}
    service := command.NewCreateWalletService(repo)
    
    // Act
    result := service.Execute(request)
    
    // Assert
    assert.True(t, result.Success)
}
```

## 🐳 部署

### Docker 建置
```bash
# 建置應用程式映像
docker build -t accounting-app .

# 執行容器
docker run -p 8080:8080 --env-file .env accounting-app
```

### Docker Compose 部署
```bash
# 生產環境
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d
```

### 環境變數
| 變數名 | 說明 | 預設值 |
|--------|------|--------|
| `DATABASE_URL` | PostgreSQL 連線字串 | `postgres://postgres:password@localhost:5432/accountingdb?sslmode=disable` |
| `PORT` | HTTP 伺服器埠號 | `8080` |
| `ENV` | 執行環境 | `development` |

## 📚 相關文檔

- [CLAUDE.md](CLAUDE.md) - Claude AI Context 檔案
- [專案狀態報告](docs/PROJECT-STATUS.md) - 詳細專案狀態
- [架構決策記錄](docs/ARCHITECTURE-DECISIONS.md) - ADR 文檔
- [Bridge Pattern 設計](docs/bridge-pattern-design.md) - 架構設計詳情
- [Docker 使用指南](README-Docker.md) - Docker 環境設定

## 🤝 貢獻指南

### 提交 Issue
1. 使用清楚的標題描述問題
2. 提供重現步驟
3. 包含環境資訊 (Go 版本、OS 等)

### 提交 Pull Request
1. Fork 專案並建立功能分支
2. 編寫測試並確保通過
3. 遵循編碼規範
4. 提供清楚的 commit 訊息
5. 更新相關文檔

### 開發環境設定
```bash
# 1. Fork 和 Clone
git clone https://github.com/yourusername/accountingApp.git
cd accountingApp

# 2. 設定開發環境
./scripts/start-dev.sh

# 3. 執行測試
go test ./...
```

## 📄 授權條款

本專案採用 MIT 授權條款 - 詳見 [LICENSE](LICENSE) 檔案

## 🔗 相關連結

- [Go 官方文檔](https://golang.org/doc/)
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Domain-Driven Design](https://domainlanguage.com/ddd/)
- [PostgreSQL](https://www.postgresql.org/)

---

**維護者**: [您的名稱]  
**最後更新**: 2024年  
**專案狀態**: 🟢 積極開發中