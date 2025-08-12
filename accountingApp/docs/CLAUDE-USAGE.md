# SuperClaude 使用指南

> 如何在此專案中有效使用 Claude AI 和 SuperClaude 框架

## 🚀 快速開始新會話

### 1. 載入專案 Context
```
/load @CLAUDE.md
```

這會載入完整的專案 context，包含：
- 專案概述和技術棧
- 架構設計和已完成功能
- 開發環境和常用指令
- 待辦事項和已知問題

### 2. 確認專案狀態
```
/sc:analyze
```

分析當前專案狀態，檢查：
- 程式碼品質
- 測試狀態
- 架構合規性
- 潛在問題

### 3. 驗證環境
```
/sc:build
```

編譯並測試專案，確保一切正常運作。

## 🛠️ 常用 SuperClaude 命令

### 開發相關

#### 實作新功能
```bash
# 實作完整功能
/sc:implement "實作支出記錄的 CRUD 操作" --type feature --with-tests

# 實作 API 端點
/sc:implement "新增錢包查詢 API" --type api --framework http

# 實作 UI 元件
/sc:implement "錢包餘額顯示元件" --type component --framework react
```

#### 改善程式碼品質
```bash
# 改善特定檔案
/sc:improve internal/accounting/adapter/controller/walletController.go

# 效能優化
/sc:improve --focus performance

# 安全性加強
/sc:improve --focus security

# 程式碼重構
/sc:improve --focus quality
```

#### 分析和調試
```bash
# 全面分析
/sc:analyze --comprehensive

# 效能分析
/sc:analyze --focus performance

# 架構分析
/sc:analyze --focus architecture

# 安全性分析
/sc:analyze --focus security
```

### 設計相關

#### 架構設計
```bash
# 新功能架構設計
/sc:design "使用者認證系統架構"

# 效能優化設計
/sc:design "快取層架構" --focus performance

# 資料庫設計
/sc:design "交易記錄表結構" --focus database
```

### 測試相關

#### 測試執行
```bash
# 執行所有測試
/sc:test

# 執行特定類型測試
/sc:test --type unit

# 效能測試
/sc:test --benchmark

# 整合測試
/sc:test --type integration
```

### 文檔相關

#### 文檔生成
```bash
# 生成 API 文檔
/sc:document --type api

# 更新 README
/sc:document README.md

# 生成架構文檔
/sc:document --type architecture
```

## 🎯 工作流程建議

### 1. 每日開發流程
```bash
# 1. 開始新會話
/load @CLAUDE.md

# 2. 檢查專案狀態
/sc:analyze

# 3. 查看待辦事項
查看 CLAUDE.md 中的待辦事項

# 4. 開始開發
/sc:implement [功能描述]

# 5. 測試驗證
/sc:test

# 6. 程式碼改善
/sc:improve [目標檔案或功能]
```

### 2. 新功能開發
```bash
# 1. 設計階段
/sc:design "新功能架構設計"

# 2. 實作階段
/sc:implement "功能實作" --with-tests

# 3. 測試階段
/sc:test --type unit

# 4. 整合階段
/sc:build

# 5. 文檔階段
/sc:document --type feature
```

### 3. 問題解決
```bash
# 1. 問題分析
/sc:analyze --focus [問題領域]

# 2. 解決方案設計
/sc:design "問題解決方案"

# 3. 實作修復
/sc:implement "修復實作"

# 4. 驗證修復
/sc:test
```

## 🎨 SuperClaude 特色功能

### Auto-Activation (自動啟用)

SuperClaude 會根據專案內容自動啟用適合的 persona：

- **Backend 開發**: 自動啟用 backend persona，使用 Context7 和 Sequential
- **Frontend 開發**: 自動啟用 frontend persona，使用 Magic 和 Playwright
- **架構設計**: 自動啟用 architect persona，使用 Sequential 和 Context7
- **效能優化**: 自動啟用 performance persona，使用 Playwright
- **安全分析**: 自動啟用 security persona，使用 Sequential

### Flag 系統

#### 思考深度控制
```bash
--think          # 中等深度分析 (~4K tokens)
--think-hard     # 深度架構分析 (~10K tokens)  
--ultrathink     # 最大深度分析 (~32K tokens)
```

#### 效率優化
```bash
--uc             # 極簡輸出模式 (30-50% token 節省)
--validate       # 操作前驗證風險
--safe-mode      # 最大驗證模式
```

#### MCP Server 控制
```bash
--c7             # 啟用 Context7 (文檔查詢)
--seq            # 啟用 Sequential (複雜分析)
--magic          # 啟用 Magic (UI 生成)
--play           # 啟用 Playwright (測試)
--all-mcp        # 啟用所有 MCP servers
```

### Wave 系統 (高級功能)

對於複雜的多階段操作，可以使用 Wave 系統：

```bash
# 自動 Wave 模式
/sc:improve --comprehensive  # 自動判斷是否需要 Wave

# 強制 Wave 模式
/sc:analyze --wave-mode force

# Wave 策略選擇
/sc:implement --wave-strategy progressive  # 漸進式
/sc:improve --wave-strategy systematic     # 系統化
/sc:design --wave-strategy adaptive        # 適應式
```

## 📝 最佳實踐

### Context 管理
1. **每次開始新會話都要載入 CLAUDE.md**
2. **重要變更後更新 CLAUDE.md**
3. **保持文檔與程式碼同步**

### 命令使用
1. **從分析開始**: 先了解現狀再行動
2. **測試驅動**: 實作後立即測試驗證
3. **漸進改善**: 小步快跑，持續改善

### 效率提升
1. **使用 --uc 模式節省 token**
2. **善用自動啟用功能**
3. **批次處理相關任務**

## 🔍 故障排除

### 常見問題

#### Claude 忘記 context
**解決方案**: 重新載入 CLAUDE.md
```
/load @CLAUDE.md
```

#### 命令執行失敗
**檢查項目**:
1. 確認檔案路徑正確
2. 檢查權限設定
3. 驗證依賴套件
4. 查看錯誤訊息

#### 性能問題
**優化策略**:
1. 使用 `--uc` 節省 token
2. 限制分析範圍
3. 分批處理大型操作

### Debug 模式
```bash
# 詳細輸出模式
/sc:analyze --verbose

# 內省模式 (查看 AI 思考過程)
/sc:improve --introspect
```

## 🎓 進階技巧

### 自定義工作流程
```bash
# 建立自定義命令組合
alias sc-daily="/load @CLAUDE.md && /sc:analyze && /sc:build"
```

### 專案特定配置
根據專案需求調整 SuperClaude 行為：
- 在 CLAUDE.md 中記錄常用 flag 組合
- 建立專案特定的 prompt templates
- 維護專案特定的最佳實踐列表

### 團隊協作
1. **統一 Context**: 確保團隊成員都使用相同的 CLAUDE.md
2. **文檔同步**: 重要變更後立即更新文檔
3. **知識分享**: 分享有效的 SuperClaude 使用模式

## 📚 參考資料

### SuperClaude 文檔
- SuperClaude Framework 官方文檔
- Persona 系統說明
- MCP Server 指南
- Wave 系統詳解

### 專案文檔
- [CLAUDE.md](../CLAUDE.md) - 主要 context 檔案
- [PROJECT-STATUS.md](PROJECT-STATUS.md) - 專案狀態
- [ARCHITECTURE-DECISIONS.md](ARCHITECTURE-DECISIONS.md) - 架構決策

---

**提示**: 善用 SuperClaude 的自動化功能，專注於創意和決策，讓 AI 處理重複性任務！