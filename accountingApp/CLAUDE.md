# Accounting App - Claude Context

> 這是一個 Go 會計應用程式專案的 Claude AI context 檔案。
> 在新的 Claude 會話中使用 `/load @CLAUDE.md` 快速恢復專案理解。

## 🏗️ 專案概述

**專案名稱**: Accounting App  
**語言**: Go 1.21+  
**架構模式**: Clean Architecture + Domain-Driven Design (DDD)  
**資料庫**: PostgreSQL  
**容器化**: Docker Compose  

## 📋 技術棧

### 核心技術
- **Backend**: Go + net/http
- **Database**: PostgreSQL 15
- **ORM**: 原生 SQL (database/sql)
- **Architecture**: Clean Architecture (4-layer)
- **Pattern**: Bridge Pattern (解決層級依賴問題)

### 開發工具
- **Environment**: godotenv (.env 檔案管理)
- **Containerization**: Docker Compose
- **Database Admin**: pgAdmin 4 (可選)

## 🏛️ Clean Architecture 層級

```
第一層 (Domain)       │ model.Wallet, model.ExpenseCategory, model.IncomeCategory
                     │      ↑
第二層 (Application)  │ Repository Interfaces, Mapper, RepositoryImpl
                     │      ↑  
第三層 (Adapter)      │ Controllers, Peer Implementation (PostgreSQL)
                     │      ↑
第四層 (Frameworks)   │ Database Connection, Web Router
```

### 重要架構決策
- **Bridge Pattern**: 解決第四層直接依賴第一層的問題
- **層級修正**: mapper 和 repositoryimpl 在第二層，peer impl 在第三層
- **Peer Interface**: 在第二層定義，第三層實現

## ✅ 已完成功能

### 核心架構
- [x] Clean Architecture 四層架構
- [x] Bridge Pattern 實作 (解決 PostgresWalletRepository 依賴違規)  
- [x] Domain-Driven Design 實作
- [x] Repository Pattern + Dependency Inversion

### 環境管理
- [x] .env 檔案管理 (godotenv)
- [x] Docker Compose 配置
- [x] PostgreSQL 自動建表
- [x] 開發環境腳本 (`scripts/start-dev.sh`)

### Domain Models
- [x] Wallet (錢包聚合)
- [x] ExpenseCategory (支出分類聚合)
- [x] IncomeCategory (收入分類聚合)
- [x] Money Value Object
- [x] Domain Services 驗證

### Repository Implementation
- [x] WalletRepository (Bridge Pattern)
- [x] ExpenseCategoryRepository 介面
- [x] IncomeCategoryRepository 介面
- [x] PostgreSQL 資料層實作

### Use Cases
- [x] CreateWalletService
- [x] AddExpenseService / AddIncomeService  
- [x] GetWalletBalanceService
- [x] CreateExpenseCategoryService / CreateIncomeCategoryService

### Web Layer
- [x] HTTP Controllers
- [x] REST API 端點
- [x] Request/Response 處理

### Testing
- [x] Domain Model 測試
- [x] Use Case 測試
- [x] Dependency Inversion 測試
- [x] Mock Repository 實作

## 🚀 開發環境啟動

### 快速開始
```bash
# 1. 啟動資料庫
./scripts/start-dev.sh

# 2. 啟動應用
go run cmd/accoountingApp/main.go

# 3. (可選) 啟動 pgAdmin
docker-compose --profile admin up -d pgadmin
```

### 環境配置
- **Database**: `postgres://postgres:password@localhost:5432/accountingdb?sslmode=disable`
- **App Port**: `8080`
- **pgAdmin**: `http://localhost:8081` (admin@accounting.com / admin123)

### 常用指令
```bash
# 測試
go test ./...

# 編譯
go build ./cmd/accoountingApp

# Docker 管理
docker-compose up -d postgres
docker-compose down
docker-compose logs postgres
```

## 📁 專案結構

```
accountingApp/
├── cmd/accoountingApp/main.go          # 應用程式入口
├── internal/accounting/
│   ├── domain/
│   │   ├── model/                      # Domain Models & Aggregates
│   │   └── service/                    # Domain Services
│   ├── application/
│   │   ├── command/                    # Command Use Cases
│   │   ├── query/                      # Query Use Cases  
│   │   ├── repository/                 # Repository Interfaces
│   │   ├── mapper/                     # Domain ↔ Data 轉換
│   │   └── usecase/                    # Use Case Interfaces
│   ├── adapter/
│   │   └── controller/                 # HTTP Controllers
│   └── frameworks/
│       ├── database/                   # 資料庫實作
│       └── web/                        # HTTP 路由
├── test/                              # 測試檔案
├── docs/                              # 專案文檔
├── scripts/                           # 開發腳本
├── .env                              # 環境變數
├── docker-compose.yml               # Docker 配置
└── CLAUDE.md                        # 本檔案
```

## 🎯 近期完成項目

### Bridge Pattern 實作 (解決架構違規)
- **問題**: PostgresWalletRepository (第四層) 直接依賴 Domain Model (第一層)
- **解決**: 實作 Bridge Pattern，透過 Peer 介面分離依賴
- **結果**: 符合 Clean Architecture 依賴規則

### 層級架構修正
- **修正**: mapper 和 repositoryimpl 移到第二層，peer impl 在第三層
- **原因**: 使用者指正「mapper跟repositoryimpl, peer interface都在第二層，peer impl在第三層」
- **狀態**: 已完成，所有測試通過

### 環境變數管理實作
- **新增**: .env 檔案支援
- **工具**: github.com/joho/godotenv
- **配置**: DATABASE_URL, PORT, ENV
- **狀態**: 已完成並測試

### Docker 環境建立
- **服務**: PostgreSQL 15 + pgAdmin 4
- **功能**: 自動建表、資料持久化、健康檢查
- **腳本**: `scripts/start-dev.sh` 一鍵啟動
- **狀態**: 已完成並測試

## 📋 最近變更

### 2025-08-10 - Domain Model重構
**✅ 重大架構改進**:
- 🔄 **Wallet聚合增強**: 新增內部交易記錄管理 (expenseRecords, incomeRecords, transfers)
- 🔄 **Repository精簡**: 簡化為基本CRUD操作，移除複雜查詢
- 🔄 **Use Case重構**: 只依賴Repository，避免Inquiry依賴
- ➕ **新增ProcessTransferService**: 完整的轉帳處理流程
- 🎯 **Transaction統一介面**: 統一管理所有交易記錄類型

**Domain Model改進**:
- ✅ **Wallet Aggregate**: 包含完整的交易記錄聚合 (11個欄位)
- ✅ **Transaction ValueObject**: 統一交易記錄介面 (2個欄位)
- ✅ **聚合內查詢方法**: GetTransactionHistory(), GetMonthlyTotal()
- ✅ **載入狀態管理**: isFullyLoaded 標記聚合載入狀態

**Repository重構**:
- 🔄 **WalletRepository**: 基本CRUD + FindByIDWithTransactions()
- 🔄 **CategoryRepository**: 精簡為必要的Domain查詢
- ❌ **移除Inquiry依賴**: Use Case不再直接依賴複雜查詢

**Use Case改進**:
- 🔄 **AddExpenseService**: 載入完整聚合，Domain Model處理業務邏輯
- 🔄 **GetWalletBalanceService**: 只載入基本資訊，效能優化
- ➕ **ProcessTransferService**: 新增轉帳處理，支援雙錢包操作

**分析摘要**: 掃描了 22 個檔案，發現 7 個Use Cases，11個Domain Models



## 🔄 待辦事項

### 高優先級  
- [ ] 實作 WalletRepositoryImpl 的 FindByIDWithTransactions 方法
- [ ] 實作 Repository Peer 層的 PostgreSQL 實作
- [ ] 新增 ProcessTransfer API 端點
- [ ] 實作轉帳交易管理 (Transaction Manager)

### 中優先級  
- [ ] API 錯誤處理標準化
- [ ] 輸入驗證增強
- [ ] 日誌記錄系統
- [ ] API 文檔生成

### 低優先級
- [ ] 實作 Inquiry 介面 (只在真正必要時)
- [ ] Swagger API 文檔
- [ ] 效能測試
- [ ] CI/CD Pipeline
- [ ] Docker 多階段建置

## 🐛 已知問題

1. **Repository Peer層實作**: 需要實作 FindByIDWithTransactions 的資料存取邏輯
2. **交易管理**: ProcessTransfer 需要 Transaction Manager 保證 ACID
3. **錯誤處理**: API 層級的錯誤處理需要標準化  
4. **日誌**: 缺乏結構化日誌記錄

## 📚 相關文檔

- `docs/bridge-pattern-design.md` - Bridge Pattern 詳細設計
- `docs/dependency-flow.md` - 依賴流向分析
- `README-Docker.md` - Docker 使用指南
- `internal/accounting/frameworks/database/schema.sql` - 資料庫結構

## 🤖 SuperClaude 使用建議

### 開始新會話時
```
/load @CLAUDE.md - 載入此 context
/sc:analyze - 分析專案現狀  
/sc:build - 編譯並測試專案
```

### 常用命令
```
/sc:implement [feature] - 實作新功能
/sc:improve [target] - 改善程式碼品質
/sc:design [domain] - 架構設計
/sc:test - 執行測試工作流程
```

### 重要提醒
- 始終遵循 Clean Architecture 原則
- 保持 Bridge Pattern 的完整性
- 所有變更都需要通過測試驗證
- 重要架構決策需要更新此文檔

---

**最後更新**: 2025-08-10 - Domain Model重構完成  
**專案狀態**: 🟢 開發中 - Clean Architecture重構完成，Repository層實作階段  
**測試狀態**: 🟡 部分測試需要更新 (Domain Model變更)