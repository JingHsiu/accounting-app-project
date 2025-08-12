package main

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// DocumentationAgent 智能文檔追蹤 AI Agent
type DocumentationAgent struct {
	projectRoot string
	config      *AgentConfig
	scanner     *CodeScanner
	analyzer    *FeatureAnalyzer
	updater     *DocUpdater
}

// AgentConfig Agent 配置
type AgentConfig struct {
	MonitorPaths []string `json:"monitor_paths"`
	IgnorePaths  []string `json:"ignore_paths"`
	OutputFormat string   `json:"output_format"`
	UpdateDocs   bool     `json:"update_docs"`
}

// CodeScanner 程式碼掃描器
type CodeScanner struct {
	fileSet *token.FileSet
}

// FeatureAnalyzer 功能分析器
type FeatureAnalyzer struct {
	apiEndpoints  []APIEndpoint
	domainModels  []DomainModel
	useCases      []UseCase
	dbChanges     []DBChange
}

// DocUpdater 文檔更新器
type DocUpdater struct {
	claudeFile   string
	statusFile   string
	readmeFile   string
}

// APIEndpoint API 端點資訊
type APIEndpoint struct {
	Method      string    `json:"method"`
	Path        string    `json:"path"`
	Handler     string    `json:"handler"`
	Controller  string    `json:"controller"`
	Description string    `json:"description"`
	File        string    `json:"file"`
	Line        int       `json:"line"`
	CreatedAt   time.Time `json:"created_at"`
}

// DomainModel 領域模型資訊
type DomainModel struct {
	Name        string            `json:"name"`
	Type        string            `json:"type"` // Aggregate, Entity, ValueObject
	Fields      []ModelField      `json:"fields"`
	Methods     []string          `json:"methods"`
	File        string            `json:"file"`
	Description string            `json:"description"`
	CreatedAt   time.Time         `json:"created_at"`
}

// ModelField 模型欄位
type ModelField struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Tags string `json:"tags"`
}

// UseCase 使用案例資訊
type UseCase struct {
	Name        string    `json:"name"`
	Type        string    `json:"type"` // Command, Query
	Service     string    `json:"service"`
	Methods     []string  `json:"methods"`
	File        string    `json:"file"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

// DBChange 資料庫變更資訊
type DBChange struct {
	Type        string    `json:"type"` // CREATE_TABLE, ALTER_TABLE, etc
	Table       string    `json:"table"`
	Changes     string    `json:"changes"`
	File        string    `json:"file"`
	CreatedAt   time.Time `json:"created_at"`
}

// AnalysisResult 分析結果
type AnalysisResult struct {
	Timestamp    time.Time     `json:"timestamp"`
	FilesScanned int           `json:"files_scanned"`
	APIEndpoints []APIEndpoint `json:"api_endpoints"`
	DomainModels []DomainModel `json:"domain_models"`
	UseCases     []UseCase     `json:"use_cases"`
	DBChanges    []DBChange    `json:"db_changes"`
	Summary      string        `json:"summary"`
}

// NewDocumentationAgent 建立新的 Documentation Agent
func NewDocumentationAgent(projectRoot string) *DocumentationAgent {
	config := &AgentConfig{
		MonitorPaths: []string{
			"internal/accounting/adapter/controller",
			"internal/accounting/application/command",
			"internal/accounting/application/query",
			"internal/accounting/domain/model",
			"internal/accounting/frameworks/database",
		},
		IgnorePaths: []string{
			"*_test.go",
			"mock_*.go",
			"vendor/",
		},
		OutputFormat: "json",
		UpdateDocs:   true,
	}

	return &DocumentationAgent{
		projectRoot: projectRoot,
		config:      config,
		scanner:     &CodeScanner{fileSet: token.NewFileSet()},
		analyzer:    &FeatureAnalyzer{},
		updater: &DocUpdater{
			claudeFile: filepath.Join(projectRoot, "CLAUDE.md"),
			statusFile: filepath.Join(projectRoot, "docs", "PROJECT-STATUS.md"),
			readmeFile: filepath.Join(projectRoot, "README.md"),
		},
	}
}

// AnalyzeProject 分析整個專案
func (agent *DocumentationAgent) AnalyzeProject() (*AnalysisResult, error) {
	fmt.Println("🔍 開始分析專案...")

	result := &AnalysisResult{
		Timestamp: time.Now(),
	}

	// 重置分析器
	agent.analyzer = &FeatureAnalyzer{}

	// 掃描所有監控路徑
	for _, monitorPath := range agent.config.MonitorPaths {
		fullPath := filepath.Join(agent.projectRoot, monitorPath)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			continue
		}

		err := filepath.WalkDir(fullPath, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if d.IsDir() || !strings.HasSuffix(path, ".go") {
				return nil
			}

			// 檢查是否要忽略
			if agent.shouldIgnoreFile(path) {
				return nil
			}

			return agent.analyzeFile(path)
		})

		if err != nil {
			fmt.Printf("❌ 掃描路徑 %s 時發生錯誤: %v\n", fullPath, err)
		}
	}

	// 整合分析結果
	result.APIEndpoints = agent.analyzer.apiEndpoints
	result.DomainModels = agent.analyzer.domainModels
	result.UseCases = agent.analyzer.useCases
	result.DBChanges = agent.analyzer.dbChanges
	result.FilesScanned = len(agent.analyzer.apiEndpoints) + len(agent.analyzer.domainModels) + len(agent.analyzer.useCases)

	// 生成摘要
	result.Summary = agent.generateSummary(result)

	fmt.Printf("✅ 分析完成! 掃描了 %d 個檔案\n", result.FilesScanned)
	fmt.Printf("   - 發現 %d 個 API 端點\n", len(result.APIEndpoints))
	fmt.Printf("   - 發現 %d 個 Domain Models\n", len(result.DomainModels))
	fmt.Printf("   - 發現 %d 個 Use Cases\n", len(result.UseCases))

	return result, nil
}

// shouldIgnoreFile 檢查是否應該忽略檔案
func (agent *DocumentationAgent) shouldIgnoreFile(path string) bool {
	for _, ignorePattern := range agent.config.IgnorePaths {
		if matched, _ := filepath.Match(ignorePattern, filepath.Base(path)); matched {
			return true
		}
	}
	return false
}

// analyzeFile 分析單個檔案
func (agent *DocumentationAgent) analyzeFile(filePath string) error {
	// 解析 Go 檔案
	src, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	file, err := parser.ParseFile(agent.scanner.fileSet, filePath, src, parser.ParseComments)
	if err != nil {
		return err
	}

	// 根據檔案類型進行不同的分析
	relativePath := strings.TrimPrefix(filePath, agent.projectRoot+"/")
	
	switch {
	case strings.Contains(relativePath, "controller"):
		agent.analyzeController(file, relativePath)
	case strings.Contains(relativePath, "model"):
		agent.analyzeDomainModel(file, relativePath)
	case strings.Contains(relativePath, "command") || strings.Contains(relativePath, "query"):
		agent.analyzeUseCase(file, relativePath)
	case strings.Contains(relativePath, "schema.sql"):
		agent.analyzeDBSchema(filePath)
	}

	return nil
}

// analyzeController 分析控制器檔案
func (agent *DocumentationAgent) analyzeController(file *ast.File, filePath string) {
	ast.Inspect(file, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncDecl:
			// 查找 HTTP handler 方法
			if x.Recv != nil && x.Name.IsExported() {
				endpoint := agent.extractAPIEndpoint(x, filePath)
				if endpoint != nil {
					agent.analyzer.apiEndpoints = append(agent.analyzer.apiEndpoints, *endpoint)
				}
			}
		}
		return true
	})
}

// extractAPIEndpoint 提取 API 端點資訊
func (agent *DocumentationAgent) extractAPIEndpoint(funcDecl *ast.FuncDecl, filePath string) *APIEndpoint {
	// 簡單的啟發式方法識別 HTTP handler
	funcName := funcDecl.Name.Name
	
	// 常見的 HTTP handler 模式
	httpMethods := map[string]string{
		"Create": "POST",
		"Get":    "GET", 
		"Update": "PUT",
		"Delete": "DELETE",
		"List":   "GET",
	}

	var method string
	var path string
	var description string

	// 根據函數名稱推測 HTTP 方法
	for prefix, httpMethod := range httpMethods {
		if strings.HasPrefix(funcName, prefix) {
			method = httpMethod
			break
		}
	}

	// 推測 API 路徑
	if strings.Contains(filePath, "wallet") {
		path = "/api/v1/wallets"
		description = "錢包相關 API"
	} else if strings.Contains(filePath, "category") {
		path = "/api/v1/categories"
		description = "分類相關 API"
	} else {
		path = "/api/v1/" + strings.ToLower(strings.TrimSuffix(funcName, "Handler"))
	}

	if method == "" {
		return nil
	}

	return &APIEndpoint{
		Method:      method,
		Path:        path,
		Handler:     funcName,
		Controller:  filepath.Base(filePath),
		Description: description,
		File:        filePath,
		Line:        agent.scanner.fileSet.Position(funcDecl.Pos()).Line,
		CreatedAt:   time.Now(),
	}
}

// analyzeDomainModel 分析領域模型
func (agent *DocumentationAgent) analyzeDomainModel(file *ast.File, filePath string) {
	ast.Inspect(file, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.TypeSpec:
			if structType, ok := x.Type.(*ast.StructType); ok {
				model := agent.extractDomainModel(x.Name.Name, structType, filePath)
				agent.analyzer.domainModels = append(agent.analyzer.domainModels, *model)
			}
		}
		return true
	})
}

// extractDomainModel 提取領域模型資訊
func (agent *DocumentationAgent) extractDomainModel(name string, structType *ast.StructType, filePath string) *DomainModel {
	model := &DomainModel{
		Name:      name,
		File:      filePath,
		CreatedAt: time.Now(),
	}

	// 判斷模型類型
	if strings.Contains(strings.ToLower(name), "aggregate") || 
	   name == "Wallet" || name == "ExpenseCategory" || name == "IncomeCategory" {
		model.Type = "Aggregate"
	} else if strings.Contains(strings.ToLower(name), "record") {
		model.Type = "Entity"
	} else {
		model.Type = "ValueObject"
	}

	// 提取欄位
	for _, field := range structType.Fields.List {
		for _, name := range field.Names {
			fieldInfo := ModelField{
				Name: name.Name,
				Type: agent.typeToString(field.Type),
			}
			
			if field.Tag != nil {
				fieldInfo.Tags = field.Tag.Value
			}
			
			model.Fields = append(model.Fields, fieldInfo)
		}
	}

	// 生成描述
	model.Description = fmt.Sprintf("%s 類型的 %s", model.Type, name)

	return model
}

// analyzeUseCase 分析使用案例
func (agent *DocumentationAgent) analyzeUseCase(file *ast.File, filePath string) {
	ast.Inspect(file, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.TypeSpec:
			if _, ok := x.Type.(*ast.StructType); ok && strings.HasSuffix(x.Name.Name, "Service") {
				useCase := &UseCase{
					Name:      x.Name.Name,
					Service:   x.Name.Name,
					File:      filePath,
					CreatedAt: time.Now(),
				}

				// 判斷類型
				if strings.Contains(filePath, "command") {
					useCase.Type = "Command"
					useCase.Description = "執行業務操作的命令處理器"
				} else {
					useCase.Type = "Query" 
					useCase.Description = "查詢資料的查詢處理器"
				}

				agent.analyzer.useCases = append(agent.analyzer.useCases, *useCase)
			}
		}
		return true
	})
}

// analyzeDBSchema 分析資料庫 Schema
func (agent *DocumentationAgent) analyzeDBSchema(filePath string) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return
	}

	// 使用正則表達式查找 CREATE TABLE 語句
	createTableRegex := regexp.MustCompile(`CREATE TABLE IF NOT EXISTS (\w+) \(`)
	matches := createTableRegex.FindAllStringSubmatch(string(content), -1)

	for _, match := range matches {
		if len(match) >= 2 {
			dbChange := DBChange{
				Type:      "CREATE_TABLE",
				Table:     match[1],
				Changes:   fmt.Sprintf("建立資料表 %s", match[1]),
				File:      filePath,
				CreatedAt: time.Now(),
			}
			agent.analyzer.dbChanges = append(agent.analyzer.dbChanges, dbChange)
		}
	}
}

// typeToString 將 AST 類型轉換為字串
func (agent *DocumentationAgent) typeToString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return "*" + agent.typeToString(t.X)
	case *ast.ArrayType:
		return "[]" + agent.typeToString(t.Elt)
	case *ast.SelectorExpr:
		return agent.typeToString(t.X) + "." + t.Sel.Name
	default:
		return "interface{}"
	}
}

// generateSummary 生成分析摘要
func (agent *DocumentationAgent) generateSummary(result *AnalysisResult) string {
	var summary strings.Builder
	
	summary.WriteString(fmt.Sprintf("📊 專案分析摘要 (%s)\n\n", result.Timestamp.Format("2006-01-02 15:04:05")))
	
	if len(result.APIEndpoints) > 0 {
		summary.WriteString("🌐 新發現的 API 端點:\n")
		for _, endpoint := range result.APIEndpoints {
			summary.WriteString(fmt.Sprintf("  - %s %s (%s)\n", endpoint.Method, endpoint.Path, endpoint.Handler))
		}
		summary.WriteString("\n")
	}

	if len(result.DomainModels) > 0 {
		summary.WriteString("🏗️ Domain Models:\n")
		for _, model := range result.DomainModels {
			summary.WriteString(fmt.Sprintf("  - %s (%s) - %d 欄位\n", model.Name, model.Type, len(model.Fields)))
		}
		summary.WriteString("\n")
	}

	if len(result.UseCases) > 0 {
		summary.WriteString("⚙️ Use Cases:\n")
		for _, useCase := range result.UseCases {
			summary.WriteString(fmt.Sprintf("  - %s (%s)\n", useCase.Name, useCase.Type))
		}
		summary.WriteString("\n")
	}

	return summary.String()
}

// UpdateDocumentation 更新文檔
func (agent *DocumentationAgent) UpdateDocumentation(result *AnalysisResult) error {
	if !agent.config.UpdateDocs {
		return nil
	}

	fmt.Println("📝 開始更新文檔...")

	// 更新 CLAUDE.md
	if err := agent.updateClaudeFile(result); err != nil {
		return fmt.Errorf("更新 CLAUDE.md 失敗: %w", err)
	}

	// 更新 PROJECT-STATUS.md
	if err := agent.updateStatusFile(result); err != nil {
		return fmt.Errorf("更新 PROJECT-STATUS.md 失敗: %w", err)
	}

	fmt.Println("✅ 文檔更新完成!")
	return nil
}

// updateClaudeFile 更新 CLAUDE.md 檔案
func (agent *DocumentationAgent) updateClaudeFile(result *AnalysisResult) error {
	// 讀取現有檔案
	content, err := os.ReadFile(agent.updater.claudeFile)
	if err != nil {
		return err
	}

	contentStr := string(content)
	
	// 新增最近變更區段
	changesSection := agent.generateChangesSection(result)
	
	// 尋找插入點或建立新區段
	if strings.Contains(contentStr, "## 📋 最近變更") {
		// 替換現有區段
		re := regexp.MustCompile(`## 📋 最近變更.*?(?=## |\z)`)
		contentStr = re.ReplaceAllString(contentStr, changesSection)
	} else {
		// 在適當位置插入新區段
		insertPoint := strings.Index(contentStr, "## 🔄 待辦事項")
		if insertPoint == -1 {
			contentStr += "\n\n" + changesSection
		} else {
			contentStr = contentStr[:insertPoint] + changesSection + "\n\n" + contentStr[insertPoint:]
		}
	}

	// 寫回檔案
	return os.WriteFile(agent.updater.claudeFile, []byte(contentStr), 0644)
}

// updateStatusFile 更新 PROJECT-STATUS.md 檔案
func (agent *DocumentationAgent) updateStatusFile(result *AnalysisResult) error {
	// 讀取現有檔案
	content, err := os.ReadFile(agent.updater.statusFile)
	if err != nil {
		return err
	}

	contentStr := string(content)
	
	// 更新最後更新時間
	now := time.Now().Format("2006-01-02 15:04:05")
	re := regexp.MustCompile(`> 最後更新:.*`)
	contentStr = re.ReplaceAllString(contentStr, fmt.Sprintf("> 最後更新: %s", now))

	// 寫回檔案
	return os.WriteFile(agent.updater.statusFile, []byte(contentStr), 0644)
}

// generateChangesSection 生成變更區段內容
func (agent *DocumentationAgent) generateChangesSection(result *AnalysisResult) string {
	var section strings.Builder
	
	section.WriteString("## 📋 最近變更\n\n")
	section.WriteString(fmt.Sprintf("### %s\n", result.Timestamp.Format("2006-01-02")))
	
	if len(result.APIEndpoints) > 0 {
		section.WriteString("**API 端點**:\n")
		for _, endpoint := range result.APIEndpoints {
			section.WriteString(fmt.Sprintf("- ✅ %s `%s %s` - %s\n", endpoint.Handler, endpoint.Method, endpoint.Path, endpoint.Description))
		}
		section.WriteString("\n")
	}

	if len(result.DomainModels) > 0 {
		section.WriteString("**Domain Models**:\n")
		for _, model := range result.DomainModels {
			section.WriteString(fmt.Sprintf("- ✅ %s (%s) - %s\n", model.Name, model.Type, model.Description))
		}
		section.WriteString("\n")
	}

	if len(result.UseCases) > 0 {
		section.WriteString("**Use Cases**:\n")
		for _, useCase := range result.UseCases {
			section.WriteString(fmt.Sprintf("- ✅ %s (%s) - %s\n", useCase.Name, useCase.Type, useCase.Description))
		}
		section.WriteString("\n")
	}

	section.WriteString(fmt.Sprintf("**分析摘要**: 掃描了 %d 個檔案，發現 %d 個新功能\n\n", 
		result.FilesScanned, len(result.APIEndpoints)+len(result.DomainModels)+len(result.UseCases)))

	return section.String()
}

// SaveAnalysisResult 儲存分析結果
func (agent *DocumentationAgent) SaveAnalysisResult(result *AnalysisResult) error {
	outputFile := filepath.Join(agent.projectRoot, "docs", "analysis-result.json")
	
	// 確保目錄存在
	if err := os.MkdirAll(filepath.Dir(outputFile), 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(outputFile, data, 0644)
}

// main 函數
func main() {
	if len(os.Args) < 2 {
		fmt.Println("使用方式: go run doc-agent.go <project-root>")
		os.Exit(1)
	}

	projectRoot := os.Args[1]
	agent := NewDocumentationAgent(projectRoot)

	// 執行專案分析
	result, err := agent.AnalyzeProject()
	if err != nil {
		fmt.Printf("❌ 分析失敗: %v\n", err)
		os.Exit(1)
	}

	// 儲存分析結果
	if err := agent.SaveAnalysisResult(result); err != nil {
		fmt.Printf("⚠️ 儲存分析結果失敗: %v\n", err)
	}

	// 更新文檔
	if err := agent.UpdateDocumentation(result); err != nil {
		fmt.Printf("❌ 更新文檔失敗: %v\n", err)
		os.Exit(1)
	}

	// 輸出摘要
	fmt.Println("\n" + result.Summary)
}