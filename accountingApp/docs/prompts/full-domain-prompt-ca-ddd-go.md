<!--
Title: Generate Accounting Domain + Usecase + Adapter in Go
Tags: Clean Architecture, DDD, Go, PostgreSQL, Bridge Pattern
Last-Updated: 2025-07-28
-->

我正在開發一個帳務管理應用，請使用 Go 語言實作，架構遵循 Clean Architecture 與 Domain-Driven Design，資料庫使用 PostgreSQL。請根據以下 Domain Model 建立完整的分層代碼：

---

📦 架構要求：

- 遵循 Clean Architecture 分層：
    - domain layer：Aggregate、Entity、ValueObject、Repository Interface
    - usecase layer：Application Service 與 Input/Output 結構
    - adapter layer：repository 實作（Postgres），需使用 Mapper 轉換 Aggregate <-> AggregateData
- DDD
- Event Sourcing 與 CQRS (Command Query Responsibility Segregation)

- 實作 Repository 時使用 **Bridge Pattern**，不改變原本的 repository interface，但不能將domain 資料外流到第三層以外
- 所有 Aggregate 的 ID 應由 Domain constructor 自行產生
- UseCase Output 為統一格式：`ExitCode`, `Message`, `ID`
- 請為每個 Aggregate 各生成至少 1 組 UseCase 與測試範例（如 CreateXXX）

---

🧩 Domain Model：

### 1️⃣ User BC

- Entity：`User`
- 備註：此領域為認證授權使用，可與會計邏輯解耦，不需實作帳務邏輯。

---

### 2️⃣ Project BC

- Entity：`Project`
- 用途：作為帳務集合，例如旅遊、婚禮、計畫等，讓 Wallet 與 Record 可歸屬某個 Project。

---

### 3️⃣ Account BC

#### Aggregate 1：Wallet Aggregate

- Aggregate Root：`Wallet`
- Entities：
    - `ExpenseRecord`
    - `IncomeRecord`
    - `Transfer`
- Value Objects：
    - `Money`（Amount + Currency）
    - `WalletType`（CASH, BANK）
- 屬性：
    - `Wallet.ID`（由 domain 自產）
    - `UserID`
    - `Name`
    - `Type`
    - `Currency`
    - `CreatedAt`
- 備註：
    - `ExpenseRecord`, `IncomeRecord`, `Transfer` 為 Wallet 下的資金紀錄行為，應封裝於 Wallet 中。
    - `Records` 不需持有實體，可透過 repository 查詢。

#### Aggregate 2：Expense Category Aggregate

- Aggregate Root：`ExpenseCategory`
- Entity：
    - `ExpenseSubcategory`
- Value Objects：
    - `CategoryName`
- 屬性：
    - `ExpenseCategory.ID`
    - `UserID`
    - `Name`
    - `Subcategories []ExpenseSubcategory`

#### Aggregate 3：Income Category Aggregate

- Aggregate Root：`IncomeCategory`
- Entity：
    - `IncomeSubcategory`
- 屬性：
    - `IncomeCategory.ID`
    - `UserID`
    - `Name`
    - `Subcategories []IncomeSubcategory`

---

📌 附加指引：

- 所有 Aggregate 應有建構邏輯驗證（例如：金額不可為負、WalletType 合法性）
- Domain 層不能依賴 Infrastructure（Postgres、Mapper）
- Repository 不能直接儲存 Aggregate，必須經由 Mapper 轉為 Persistence Model

---

請為以上需求產生：
- domain code（entities, aggregates, value objects）
- usecase（CreateWallet 等典型操作，含 input/output struct）
- adapter（repository 實作含 Mapper 與 DB model）
- 簡單範例測試與用例（如單元測試與 main 範例）
- 簡單文青的前端介面 (React)