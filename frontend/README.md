# 記帳管理系統 Frontend

使用 React + TypeScript 開發的理財管理系統前端，採用淺粉紫色主題設計。

## 🚀 技術棧

- **React 18** - 用戶界面框架
- **TypeScript** - 靜態類型檢查
- **Vite** - 構建工具
- **Tailwind CSS** - 樣式框架
- **React Query** - 數據獲取和狀態管理
- **React Router** - 路由管理
- **Axios** - HTTP 客戶端
- **Lucide React** - 圖標庫
- **Date-fns** - 日期處理

## 🎨 設計系統

### 色彩主題
- **主色調**: 淺粉紫色 (`primary`)
- **輔助色**: 翠綠色 (`secondary`) 用於收入
- **強調色**: 珊瑚色 (`accent`) 用於支出
- **中性色**: 灰色系列 (`neutral`)

### 漸層效果
- `gradient-primary`: 主要漸層背景
- `gradient-secondary`: 收入相關漸層
- `gradient-accent`: 支出相關漸層

### 玻璃態效果
使用 `glass-card` 類別實現毛玻璃效果，搭配背景模糊和邊框。

## 📱 功能特色

### 核心功能
1. **儀表板** - 財務概況總覽
2. **錢包管理** - 多錢包支持（現金、銀行、信用卡、投資）
3. **交易記錄** - 收入、支出和轉帳管理
4. **類別管理** - 自定義收支分類

### UI/UX 特色
- 響應式設計，支持手機和桌面端
- 毛玻璃風格界面
- 流暢的動畫效果
- 直觀的色彩區分（收入綠色、支出紅色）
- 側邊欄導航

## 🛠️ 開發指南

### 安裝依賴
```bash
npm install
```

### 開發模式
```bash
npm run dev
```

### 構建生產版本
```bash
npm run build
```

### 代碼檢查
```bash
npm run lint
```

## 📁 項目結構

```
src/
├── components/          # 共用組件
│   ├── ui/             # 基礎 UI 組件
│   │   ├── Card.tsx    # 卡片組件
│   │   ├── Button.tsx  # 按鈕組件
│   │   ├── Input.tsx   # 輸入框組件
│   │   └── Modal.tsx   # 彈窗組件
│   └── Layout.tsx      # 主佈局
├── pages/              # 頁面組件
│   ├── Dashboard.tsx   # 儀表板
│   ├── Wallets.tsx     # 錢包管理
│   ├── Transactions.tsx # 交易記錄
│   └── Categories.tsx  # 類別管理
├── services/           # API 服務
│   ├── api.ts          # Axios 配置
│   ├── walletService.ts    # 錢包 API
│   ├── transactionService.ts # 交易 API
│   ├── categoryService.ts   # 類別 API
│   └── dashboardService.ts  # 儀表板 API
├── types/              # TypeScript 類型
│   └── index.ts        # 類型定義
├── utils/              # 工具函數
│   └── format.ts       # 格式化函數
├── App.tsx             # 主應用
└── main.tsx            # 應用入口
```

## 🔧 配置說明

### Vite 配置
- 代理 API 請求到後端 (localhost:8080)
- 路徑別名 `@` 指向 `src` 目錄
- 開發服務器端口: 3000

### TypeScript 配置
- 嚴格模式啟用
- JSX 運行時: `react-jsx`
- 路徑映射支持

### Tailwind CSS
- 自定義顏色主題
- 玻璃態效果工具類
- 動畫和漸變配置

## 🌐 API 集成

前端通過 Axios 與後端 Go 服務通信：

- **錢包 API**: `/api/wallets`
- **交易 API**: `/api/expenses`, `/api/incomes`, `/api/transfers`
- **類別 API**: `/api/categories`
- **儀表板 API**: `/api/dashboard`

### 請求攔截器
- 自動添加認證 token
- 統一錯誤處理
- 401 自動跳轉登入

## 📝 組件說明

### UI 組件
- **Card**: 支持玻璃態、懸停效果
- **Button**: 多種變體（primary, secondary, outline, ghost, danger）
- **Input**: 包含標籤、錯誤提示、圖標支持
- **Modal**: 支持鍵盤操作、背景關閉

### 業務組件
- **Layout**: 響應式側邊欄導航
- **Dashboard**: 財務概況卡片和圖表
- **WalletCard**: 錢包信息展示
- **TransactionItem**: 交易記錄項目

## 🎯 未來規劃

- [ ] 圖表和數據可視化
- [ ] 深色模式支持
- [ ] PWA 離線支持
- [ ] 多語言國際化
- [ ] 單元測試覆蓋
- [ ] 性能優化和懶加載
- [ ] 導入/導出功能

## 🐛 已知問題

- 部分 npm 依賴版本過舊的警告
- PostCSS 配置警告（不影響功能）
- 待實現圖表組件整合

## 📄 許可證

MIT License - 詳見 LICENSE 文件