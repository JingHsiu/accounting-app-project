#!/bin/bash

# 智能文檔更新腳本 - SuperClaude 整合版本
# 使用 DocumentationAgent 自動分析專案並更新文檔

set -e

# 顏色定義
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 專案根目錄
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
AGENT_PATH="$PROJECT_ROOT/scripts/doc-agent.go"
DOCS_DIR="$PROJECT_ROOT/docs"

# 函數定義
log_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

log_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

log_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

log_error() {
    echo -e "${RED}❌ $1${NC}"
}

# 顯示使用說明
show_help() {
    cat << EOF
🤖 智能文檔更新工具 - SuperClaude 整合版本

使用方式:
  $0 [選項]

選項:
  -a, --analyze-only    只執行分析，不更新文檔
  -f, --files PATTERN   只分析符合模式的檔案
  -v, --verbose         顯示詳細輸出
  -h, --help           顯示此說明

範例:
  $0                    # 完整分析並更新文檔
  $0 --analyze-only     # 只分析不更新
  $0 --files "*.go"     # 只分析 Go 檔案
  $0 --verbose          # 顯示詳細輸出

SuperClaude 整合:
  /sc:sync-docs         # SuperClaude 命令 (相當於執行此腳本)

EOF
}

# 檢查前置條件
check_prerequisites() {
    log_info "檢查前置條件..."
    
    # 檢查 Go 是否安裝
    if ! command -v go &> /dev/null; then
        log_error "Go 未安裝，請先安裝 Go"
        exit 1
    fi

    # 檢查專案結構
    if [[ ! -d "$PROJECT_ROOT/internal/accounting" ]]; then
        log_error "專案結構不正確，請確認在正確的專案目錄中執行"
        exit 1
    fi

    # 檢查 Agent 檔案
    if [[ ! -f "$AGENT_PATH" ]]; then
        log_error "DocumentationAgent 不存在: $AGENT_PATH"
        exit 1
    fi

    log_success "前置條件檢查完成"
}

# 執行 DocumentationAgent
run_analysis() {
    log_info "啟動智能文檔分析..."
    
    # 編譯並執行 DocumentationAgent
    cd "$PROJECT_ROOT/scripts"
    
    if [[ "$VERBOSE" == "true" ]]; then
        log_info "執行命令: go run doc-agent.go \"$PROJECT_ROOT\""
        go run doc-agent.go "$PROJECT_ROOT"
    else
        go run doc-agent.go "$PROJECT_ROOT" 2>/dev/null
    fi
    
    if [[ $? -eq 0 ]]; then
        log_success "文檔分析完成"
    else
        log_error "文檔分析失敗"
        exit 1
    fi
}

# 顯示分析結果
show_results() {
    local analysis_file="$DOCS_DIR/analysis-result.json"
    
    if [[ -f "$analysis_file" ]]; then
        log_info "分析結果摘要:"
        
        # 提取關鍵資訊 (需要 jq，如果沒有則顯示原始檔案)
        if command -v jq &> /dev/null; then
            local timestamp=$(jq -r '.timestamp' "$analysis_file" 2>/dev/null)
            local files_scanned=$(jq -r '.files_scanned' "$analysis_file" 2>/dev/null)
            local api_count=$(jq -r '.api_endpoints | length' "$analysis_file" 2>/dev/null)
            local model_count=$(jq -r '.domain_models | length' "$analysis_file" 2>/dev/null)
            local usecase_count=$(jq -r '.use_cases | length' "$analysis_file" 2>/dev/null)
            
            echo "  📊 分析時間: $timestamp"
            echo "  📁 掃描檔案: $files_scanned 個"
            echo "  🌐 API 端點: $api_count 個"
            echo "  🏗️ Domain Models: $model_count 個" 
            echo "  ⚙️ Use Cases: $usecase_count 個"
            
            # 顯示摘要
            local summary=$(jq -r '.summary' "$analysis_file" 2>/dev/null)
            if [[ "$summary" != "null" && "$VERBOSE" == "true" ]]; then
                echo ""
                echo "$summary"
            fi
        else
            log_warning "未安裝 jq，無法解析 JSON 結果。安裝 jq 以獲得更好的體驗。"
            echo "分析結果已儲存到: $analysis_file"
        fi
    else
        log_warning "未找到分析結果檔案"
    fi
}

# 驗證文檔更新
verify_updates() {
    log_info "驗證文檔更新..."
    
    local updated_files=()
    
    # 檢查各個文檔檔案的更新時間
    local current_time=$(date +%s)
    local threshold=300  # 5分鐘內的更新視為最新
    
    for doc_file in "$PROJECT_ROOT/CLAUDE.md" "$DOCS_DIR/PROJECT-STATUS.md"; do
        if [[ -f "$doc_file" ]]; then
            local file_time=$(stat -f %m "$doc_file" 2>/dev/null || stat -c %Y "$doc_file" 2>/dev/null)
            local time_diff=$((current_time - file_time))
            
            if [[ $time_diff -lt $threshold ]]; then
                updated_files+=("$(basename "$doc_file")")
            fi
        fi
    done
    
    if [[ ${#updated_files[@]} -gt 0 ]]; then
        log_success "以下文檔已更新: ${updated_files[*]}"
    else
        log_warning "文檔可能未更新，或更新時間超過 5 分鐘"
    fi
}

# 生成 SuperClaude 使用建議
generate_superclaud_suggestions() {
    if [[ "$VERBOSE" == "true" ]]; then
        cat << EOF

🤖 SuperClaude 使用建議:

下次開始新會話時:
  /load @CLAUDE.md                # 載入更新後的 context
  
相關命令:
  /sc:analyze --comprehensive     # 深度分析專案
  /sc:build                       # 驗證更新後的程式碼
  /sc:test                        # 執行測試驗證

如果發現新問題:
  /sc:improve [問題領域]          # 改善程式碼品質
  /sc:implement [功能描述]        # 實作新功能

EOF
    fi
}

# 主要執行流程
main() {
    # 預設值
    ANALYZE_ONLY=false
    VERBOSE=false
    FILE_PATTERN=""

    # 解析命令行參數
    while [[ $# -gt 0 ]]; do
        case $1 in
            -a|--analyze-only)
                ANALYZE_ONLY=true
                shift
                ;;
            -f|--files)
                FILE_PATTERN="$2"
                shift 2
                ;;
            -v|--verbose)
                VERBOSE=true
                shift
                ;;
            -h|--help)
                show_help
                exit 0
                ;;
            *)
                log_error "未知參數: $1"
                show_help
                exit 1
                ;;
        esac
    done

    # 顯示標題
    echo ""
    log_info "🤖 智能文檔更新工具 - SuperClaude 整合版本"
    echo ""

    # 執行主要流程
    check_prerequisites
    
    # 確保 docs 目錄存在
    mkdir -p "$DOCS_DIR"
    
    # 執行分析
    run_analysis
    
    # 顯示結果
    show_results
    
    if [[ "$ANALYZE_ONLY" == "false" ]]; then
        verify_updates
    fi
    
    # 生成使用建議
    generate_superclaud_suggestions
    
    echo ""
    log_success "🎉 文檔更新完成！"
    
    if [[ -f "$DOCS_DIR/analysis-result.json" ]]; then
        echo "📄 詳細分析結果: $DOCS_DIR/analysis-result.json"
    fi
}

# SuperClaude 命令別名檢測
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    # 檢查是否由 SuperClaude 呼叫
    if [[ "$1" == "--superclaud" ]]; then
        shift  # 移除 --superclaud 參數
        VERBOSE=true  # SuperClaude 模式預設顯示詳細輸出
    fi
    
    # 執行主程式
    main "$@"
fi