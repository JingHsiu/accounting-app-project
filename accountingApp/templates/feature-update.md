# 功能更新模板

## 新功能: {{.FeatureName}}

**實作日期**: {{.Date}}  
**類型**: {{.Type}} (API端點/Domain Model/Use Case)  
**狀態**: ✅ 完成

### 📋 功能描述
{{.Description}}

### 🏗️ 架構影響
- **Domain Layer**: {{.DomainChanges}}
- **Application Layer**: {{.ApplicationChanges}}
- **Adapter Layer**: {{.AdapterChanges}}
- **Frameworks Layer**: {{.FrameworksChanges}}

### 📁 相關檔案
{{range .RelatedFiles}}
- `{{.}}`
{{end}}

### 🌐 API 端點
{{if .APIEndpoints}}
{{range .APIEndpoints}}
- `{{.Method}} {{.Path}}` - {{.Description}}
{{end}}
{{else}}
無新的 API 端點
{{end}}

### 🏗️ Domain Models
{{if .DomainModels}}
{{range .DomainModels}}
- **{{.Name}}** ({{.Type}}) - {{.Description}}
  - 欄位數量: {{len .Fields}}
  - 檔案位置: `{{.File}}`
{{end}}
{{else}}
無新的 Domain Models
{{end}}

### ⚙️ Use Cases
{{if .UseCases}}
{{range .UseCases}}
- **{{.Name}}** ({{.Type}}) - {{.Description}}
  - 服務: {{.Service}}
  - 檔案位置: `{{.File}}`
{{end}}
{{else}}
無新的 Use Cases
{{end}}

### 🧪 測試狀態
- [ ] 單元測試
- [ ] 整合測試
- [ ] API 測試
- [ ] 效能測試

### 📝 待辦事項
- [ ] 完善錯誤處理
- [ ] 新增日誌記錄
- [ ] 效能優化
- [ ] 文檔更新

### 🔗 相關連結
- 相關 Issue: #{{.IssueNumber}}
- Pull Request: #{{.PRNumber}}
- 設計文件: [連結]({{.DesignDocLink}})

---
*🤖 此更新由 DocumentationAgent 自動生成*