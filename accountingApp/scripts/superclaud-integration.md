# SuperClaude 整合指南

## 🤖 新增的 SuperClaude 命令

### `/sc:sync-docs` - 智能文檔同步
使用 DocumentationAgent 自動分析專案並更新文檔。

```bash
# 基本使用
/sc:sync-docs

# 等同於執行
./scripts/update-docs.sh --superclaud --verbose
```

**功能**:
- 🔍 自動掃描程式碼變更
- 📊 智能分析 API 端點、Domain Models、Use Cases
- 📝 自動更新 CLAUDE.md 和 PROJECT-STATUS.md
- 📄 生成詳細分析報告

**使用時機**:
- 完成新功能實作後
- 大型重構完成後
- 定期文檔同步 (每日/每週)
- 準備新會話前

## 🔧 整合方式

### 方法 1: 在 SuperClaude 中使用 Task 工具
當用戶輸入 `/sc:sync-docs` 時，Claude 會：

```python
# SuperClaude Framework 處理
def handle_sync_docs_command():
    # 使用 Task tool 執行文檔同步
    task_description = "執行智能文檔追蹤和更新"
    task_prompt = f"""
    請執行以下命令來更新專案文檔:
    
    ./scripts/update-docs.sh --superclaud --verbose
    
    這個命令會:
    1. 分析專案中的程式碼變更
    2. 自動檢測新的 API 端點和功能
    3. 更新 CLAUDE.md 和相關文檔
    4. 生成分析報告
    """
    
    return Task(description=task_description, prompt=task_prompt, subagent_type="general-purpose")
```

### 方法 2: 直接 Bash 整合
Claude 可以直接執行更新腳本：

```bash
# 執行智能文檔更新
/Users/hsiu/Desktop/workspace/accountingApp/scripts/update-docs.sh --superclaud --verbose
```

### 方法 3: 工作流程整合
結合其他 SuperClaude 命令：

```bash
# 完整的開發後整理流程
/sc:sync-docs       # 更新文檔
/sc:build           # 驗證程式碼
/sc:test            # 執行測試
/load @CLAUDE.md    # 重新載入更新後的 context
```

## 📋 使用建議

### 何時使用 `/sc:sync-docs`
1. **功能完成後**: 使用 `/sc:implement` 完成新功能後
2. **重構完成後**: 使用 `/sc:improve` 重構程式碼後
3. **會話開始前**: 確保文檔是最新狀態
4. **定期維護**: 每日或每週定期同步

### 工作流程建議
```bash
# 標準開發流程
1. /load @CLAUDE.md           # 載入現有 context
2. /sc:implement [功能]       # 實作新功能
3. /sc:test                   # 測試驗證
4. /sc:sync-docs             # 更新文檔
5. /sc:build                 # 最終驗證

# 新會話開始
1. /sc:sync-docs             # 確保文檔最新
2. /load @CLAUDE.md          # 載入更新後的 context
3. 開始工作...
```

## 🎯 預期效果

### 自動化程度
- 📊 **95%** 的狀態更新自動化
- 📝 **80%** 的功能描述自動生成
- 🔄 **90%** 的文檔同步自動化
- 👤 **10%** 的人工審核和調整

### 智能化功能
- ✅ 自動檢測新 API 端點
- ✅ 識別 Domain Models 變更
- ✅ 追蹤 Use Case 實作狀況
- ✅ 維護完整的變更歷史
- ✅ 跨文檔一致性檢查

### 整合優勢
- 🚀 **無縫整合**: 與現有 SuperClaude 命令完美配合
- ⚡ **即時更新**: 變更後立即反映在文檔中
- 📈 **追蹤能力**: 完整的專案演進歷史
- 👥 **團隊協作**: 標準化的文檔更新流程

## 🔮 未來擴展

### 計劃中的功能
1. **智能 Commit 訊息生成**: 基於分析結果生成 git commit 訊息
2. **API 文檔自動生成**: 生成 OpenAPI/Swagger 規格
3. **測試覆蓋率追蹤**: 整合測試結果分析
4. **效能指標監控**: 追蹤程式碼品質指標

### 整合擴展
- **Git Hooks**: pre-commit 自動文檔更新
- **CI/CD 整合**: 自動化部署流程中的文檔更新
- **Slack/Teams 通知**: 文檔更新通知
- **儀表板**: 專案狀態視覺化儀表板

---

**使用提示**: 將 `/sc:sync-docs` 納入你的日常開發流程，讓文檔永遠保持最新！