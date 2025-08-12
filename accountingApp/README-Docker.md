# Docker 環境設定

## 快速啟動

### 1. 啟動資料庫
```bash
# 啟動 PostgreSQL 資料庫
docker-compose up -d postgres

# 檢查資料庫健康狀態
docker-compose ps
```

### 2. 啟動應用程式
```bash
# 啟動 Go 應用程式 (本機運行)
go run cmd/accoountingApp/main.go
```

### 3. (可選) 啟動 pgAdmin 管理介面
```bash
# 啟動 pgAdmin (可選的資料庫管理介面)
docker-compose --profile admin up -d pgadmin

# 訪問 pgAdmin: http://localhost:8081
# 帳號: admin@accounting.com
# 密碼: admin123
```

## 詳細說明

### 服務配置

**PostgreSQL 資料庫**:
- 埠號: `5432`
- 資料庫名稱: `accountingdb`
- 使用者名稱: `postgres`
- 密碼: `password`
- 自動載入: `schema.sql`

**pgAdmin 管理介面** (可選):
- 埠號: `8081`
- 帳號: `admin@accounting.com`
- 密碼: `admin123`

### 資料持久化

資料庫資料會持久化儲存在 Docker Volume `postgres_data` 中，即使停止容器資料也不會遺失。

### 常用指令

```bash
# 啟動所有服務
docker-compose up -d

# 啟動特定服務
docker-compose up -d postgres

# 啟動包含 pgAdmin
docker-compose --profile admin up -d

# 檢視日誌
docker-compose logs postgres

# 停止服務
docker-compose down

# 完全清除 (包含資料)
docker-compose down -v
```

### 連線資訊

應用程式使用的連線字串:
```
DATABASE_URL=postgres://postgres:password@localhost:5432/accountingdb?sslmode=disable
```

### 初始化

當第一次啟動 PostgreSQL 容器時，會自動執行 `schema.sql` 建立所需的資料表:
- wallets
- expense_categories
- expense_subcategories  
- income_categories
- income_subcategories
- expense_records
- income_records

### 故障排除

1. **埠號衝突**: 如果本機已有 PostgreSQL 運行，修改 docker-compose.yml 中的埠號映射
2. **權限問題**: 確認 schema.sql 檔案可讀
3. **健康檢查**: 等待 PostgreSQL 完全啟動後再連線應用程式å