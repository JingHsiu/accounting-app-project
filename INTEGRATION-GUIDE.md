# Integration Guide

> Technical integration guide for the Accounting App multi-service architecture, focusing on frontend-backend integration patterns, API contracts, and development workflows.

**Target Audience**: Developers working on cross-service integration  
**Scope**: Frontend ↔ Backend integration patterns  
**Last Updated**: 2025-08-12  

## 📋 Table of Contents

- [Integration Overview](#integration-overview)
- [API Contract Specifications](#api-contract-specifications)
- [Frontend Integration Patterns](#frontend-integration-patterns)
- [Backend Integration Patterns](#backend-integration-patterns)
- [Data Transformation Layer](#data-transformation-layer)
- [Error Handling Strategy](#error-handling-strategy)
- [Authentication & Authorization](#authentication--authorization)
- [Testing Integration Points](#testing-integration-points)
- [Development Workflow](#development-workflow)
- [Troubleshooting Guide](#troubleshooting-guide)

## Integration Overview

### Service Communication Architecture

```mermaid
graph TD
    subgraph "Frontend Service (Port 3000)"
        UI[React Components] 
        Services[API Services]
        Types[TypeScript Types]
        State[React Query Cache]
    end
    
    subgraph "Integration Layer"
        Proxy[Vite Dev Proxy]
        HTTP[HTTP/JSON]
        Interceptors[Axios Interceptors]
    end
    
    subgraph "Backend Service (Port 8080)" 
        Controllers[HTTP Controllers]
        UseCases[Use Cases]
        Domain[Domain Models]
        DB[(PostgreSQL)]
    end
    
    UI --> Services
    Services --> Interceptors
    Interceptors --> Proxy
    Proxy --> Controllers
    Controllers --> UseCases
    UseCases --> Domain
    Domain --> DB
```

### Integration Principles

1. **Contract-First Development**: API contracts defined before implementation
2. **Type Safety**: Shared type definitions between frontend and backend
3. **Error Consistency**: Unified error format across all endpoints
4. **Async Communication**: Non-blocking operations with proper loading states
5. **Data Validation**: Input validation on both client and server sides

### Technology Stack Integration

| Layer | Frontend | Integration | Backend |
|-------|----------|-------------|---------|
| **Language** | TypeScript | HTTP/JSON | Go |
| **Framework** | React 18 | REST API | Clean Architecture |
| **State Management** | React Query | Axios | Repository Pattern |
| **Validation** | Zod (planned) | JSON Schema | Go validation |
| **Error Handling** | Error Boundaries | HTTP Status Codes | Custom Error Types |

## API Contract Specifications

### Base Configuration

#### Development Environment
```typescript
// Frontend: vite.config.ts
export default defineConfig({
  server: {
    port: 3000,
    proxy: {
      '/api/v1': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
})
```

```go
// Backend: main.go
server := &http.Server{
    Addr:         ":8080",
    Handler:      app.router.SetupRoutes(),
    ReadTimeout:  10 * time.Second,
    WriteTimeout: 10 * time.Second,
}
```

#### API Base Configuration
```typescript
// Frontend: services/api.ts
export const api = axios.create({
  baseURL: '/api/v1',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
})
```

### Endpoint Specifications

#### Wallet Management

**Create Wallet**
```yaml
Endpoint: POST /api/v1/wallets
Frontend Service: walletService.create()
Backend Handler: WalletController.CreateWallet()

Request Schema:
  user_id: string (required)
  name: string (required, 1-100 chars)
  type: "CASH" | "BANK" (required)
  currency: string (required, 3 chars ISO 4217)

Success Response (200):
  id: string (UUID)
  success: boolean (true)
  message: string

Error Response (400):
  success: boolean (false)
  error:
    code: string
    message: string
    details?: string
```

**Get Wallet Balance**
```yaml
Endpoint: GET /api/v1/wallets/{id}/balance
Frontend Service: walletService.getBalance(id)
Backend Handler: WalletController.GetWalletBalance()

Path Parameters:
  id: string (UUID, required)

Success Response (200):
  walletId: string
  balance: number (cents)
  currency: string
  success: boolean (true)
  message: string

Error Response (404):
  success: boolean (false)
  error:
    code: "WALLET_NOT_FOUND"
    message: "Wallet not found"
```

#### Transaction Management

**Add Expense**
```yaml
Endpoint: POST /api/v1/expenses
Frontend Service: transactionService.addExpense()
Backend Handler: WalletController.AddExpense()

Request Schema:
  wallet_id: string (UUID, required)
  category_id: string (UUID, required)
  amount: number (cents, positive integer)
  currency: string (3 chars, required)
  description: string (optional, max 500 chars)
  date?: string (YYYY-MM-DD format, optional)

Success Response (200):
  id: string (expense record UUID)
  success: boolean (true)
  message: string

Error Response (400):
  success: boolean (false)
  error:
    code: "INSUFFICIENT_BALANCE" | "INVALID_INPUT" | "WALLET_NOT_FOUND"
    message: string
    details?: string
```

**Add Income**
```yaml
Endpoint: POST /api/v1/incomes
Frontend Service: transactionService.addIncome()  
Backend Handler: WalletController.AddIncome()

Request Schema: (Same as Add Expense)
Response Schema: (Same as Add Expense, with income record UUID)
```

#### Category Management

**Create Expense Category**
```yaml
Endpoint: POST /api/v1/categories/expense
Frontend Service: categoryService.createExpense()
Backend Handler: CategoryController.CreateExpenseCategory()

Request Schema:
  name: string (required, 1-50 chars)
  subcategories?: string[] (optional, max 20 items)

Success Response (200):
  id: string (category UUID)
  success: boolean (true)
  message: string
```

**Create Income Category**
```yaml
Endpoint: POST /api/v1/categories/income
Frontend Service: categoryService.createIncome()
Backend Handler: CategoryController.CreateIncomeCategory()

Request Schema: (Same as expense category)
Response Schema: (Same as expense category)
```

### Type Definitions

#### Shared Types
```typescript
// Frontend: types/index.ts
export interface ApiResponse<T = any> {
  success: boolean
  data?: T
  error?: ApiError
}

export interface ApiError {
  code: string
  message: string
  details?: string
}

export interface Wallet {
  id: string
  user_id: string
  name: string
  type: 'CASH' | 'BANK'
  balance: Money
  created_at: string
  updated_at: string
}

export interface Money {
  amount: number  // in cents
  currency: string // ISO 4217
}

export interface Transaction {
  id: string
  wallet_id: string
  category_id: string
  amount: Money
  description: string
  date: string
  type: 'expense' | 'income'
}

export interface Category {
  id: string
  name: string
  type: 'expense' | 'income'
  subcategories: Subcategory[]
}

export interface Subcategory {
  id: string
  category_id: string
  name: string
}
```

#### Backend Go Types
```go
// Backend: domain/model
type Wallet struct {
    ID        string    `json:"id"`
    UserID    string    `json:"user_id"`
    Name      string    `json:"name"`
    Type      WalletType `json:"type"`
    Balance   Money     `json:"balance"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

type Money struct {
    Amount   int64  `json:"amount"`    // cents
    Currency string `json:"currency"`  // ISO 4217
}

// Controller DTOs
type CreateWalletRequest struct {
    UserID   string `json:"user_id"`
    Name     string `json:"name"`
    Type     string `json:"type"`
    Currency string `json:"currency"`
}

type CreateWalletResponse struct {
    ID      string `json:"id"`
    Success bool   `json:"success"`
    Message string `json:"message"`
}
```

## Frontend Integration Patterns

### Service Layer Architecture

#### API Service Base
```typescript
// services/api.ts - Base API configuration
import axios from 'axios'
import type { ApiResponse } from '@/types'

export const api = axios.create({
  baseURL: '/api/v1',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Generic API request wrapper
export const apiRequest = {
  get: async <T>(url: string): Promise<ApiResponse<T>> => {
    try {
      const response = await api.get(url)
      return { success: true, data: response.data }
    } catch (error: any) {
      return {
        success: false,
        error: {
          code: error.response?.data?.error?.code || 'NETWORK_ERROR',
          message: error.response?.data?.error?.message || error.message,
        }
      }
    }
  },
  // ... post, put, delete methods
}
```

#### Domain Service Implementation
```typescript
// services/walletService.ts
import { apiRequest } from './api'
import type { Wallet, Money, ApiResponse } from '@/types'

export interface CreateWalletData {
  user_id: string
  name: string
  type: 'CASH' | 'BANK'
  currency: string
}

export interface WalletBalance {
  walletId: string
  balance: number
  currency: string
}

export const walletService = {
  create: async (data: CreateWalletData): Promise<ApiResponse<{id: string}>> => {
    return apiRequest.post('/wallets', data)
  },

  getBalance: async (id: string): Promise<ApiResponse<WalletBalance>> => {
    return apiRequest.get(`/wallets/${id}/balance`)
  },

  getAll: async (): Promise<ApiResponse<Wallet[]>> => {
    return apiRequest.get('/wallets')
  },

  getById: async (id: string): Promise<ApiResponse<Wallet>> => {
    return apiRequest.get(`/wallets/${id}`)
  },

  update: async (id: string, data: Partial<Wallet>): Promise<ApiResponse<Wallet>> => {
    return apiRequest.put(`/wallets/${id}`, data)
  },

  delete: async (id: string): Promise<ApiResponse<void>> => {
    return apiRequest.delete(`/wallets/${id}`)
  }
}
```

### React Query Integration

#### Query Configuration
```typescript
// hooks/useWallets.ts
import { useQuery, useMutation, useQueryClient } from 'react-query'
import { walletService, CreateWalletData } from '@/services'
import type { Wallet } from '@/types'

export const useWallets = () => {
  return useQuery<Wallet[], Error>(
    'wallets',
    async () => {
      const response = await walletService.getAll()
      if (!response.success) {
        throw new Error(response.error?.message || 'Failed to fetch wallets')
      }
      return response.data!
    },
    {
      staleTime: 5 * 60 * 1000, // 5 minutes
      cacheTime: 10 * 60 * 1000, // 10 minutes
    }
  )
}

export const useWallet = (id: string) => {
  return useQuery<Wallet, Error>(
    ['wallet', id],
    async () => {
      const response = await walletService.getById(id)
      if (!response.success) {
        throw new Error(response.error?.message || 'Wallet not found')
      }
      return response.data!
    },
    {
      enabled: !!id,
    }
  )
}

export const useCreateWallet = () => {
  const queryClient = useQueryClient()
  
  return useMutation<{id: string}, Error, CreateWalletData>(
    async (data) => {
      const response = await walletService.create(data)
      if (!response.success) {
        throw new Error(response.error?.message || 'Failed to create wallet')
      }
      return response.data!
    },
    {
      onSuccess: () => {
        // Invalidate and refetch wallets list
        queryClient.invalidateQueries('wallets')
      },
      onError: (error) => {
        console.error('Failed to create wallet:', error)
      }
    }
  )
}
```

#### Component Integration
```typescript
// components/WalletForm.tsx
import React, { useState } from 'react'
import { useCreateWallet } from '@/hooks/useWallets'
import { Button, Input, Card } from '@/components/ui'

const WalletForm: React.FC = () => {
  const [formData, setFormData] = useState({
    user_id: 'user-1', // TODO: Get from auth context
    name: '',
    type: 'CASH' as const,
    currency: 'USD'
  })
  
  const createWallet = useCreateWallet()

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      await createWallet.mutateAsync(formData)
      setFormData({ ...formData, name: '' }) // Reset form
      // Success notification handled by mutation
    } catch (error) {
      // Error handled by mutation onError
    }
  }

  return (
    <Card className="p-6">
      <h2 className="text-xl font-semibold mb-4">Create Wallet</h2>
      
      <form onSubmit={handleSubmit} className="space-y-4">
        <Input
          label="Wallet Name"
          value={formData.name}
          onChange={(e) => setFormData({...formData, name: e.target.value})}
          required
          disabled={createWallet.isLoading}
        />
        
        <div>
          <label className="block text-sm font-medium mb-2">Type</label>
          <select 
            value={formData.type}
            onChange={(e) => setFormData({...formData, type: e.target.value as 'CASH' | 'BANK'})}
            className="w-full px-3 py-2 border rounded-md"
            disabled={createWallet.isLoading}
          >
            <option value="CASH">Cash</option>
            <option value="BANK">Bank Account</option>
          </select>
        </div>
        
        <Input
          label="Currency"
          value={formData.currency}
          onChange={(e) => setFormData({...formData, currency: e.target.value.toUpperCase()})}
          maxLength={3}
          pattern="[A-Z]{3}"
          required
          disabled={createWallet.isLoading}
        />
        
        <Button 
          type="submit" 
          loading={createWallet.isLoading}
          disabled={!formData.name.trim()}
        >
          Create Wallet
        </Button>
        
        {createWallet.error && (
          <div className="text-red-600 text-sm">
            {createWallet.error.message}
          </div>
        )}
      </form>
    </Card>
  )
}

export default WalletForm
```

### State Management Strategy

#### React Query for Server State
```typescript
// Query Keys Strategy
export const queryKeys = {
  wallets: ['wallets'] as const,
  wallet: (id: string) => ['wallet', id] as const,
  walletBalance: (id: string) => ['wallet', id, 'balance'] as const,
  transactions: ['transactions'] as const,
  transaction: (id: string) => ['transaction', id] as const,
  categories: ['categories'] as const,
  category: (id: string) => ['category', id] as const,
}

// Cache Invalidation Strategy
export const invalidateQueries = {
  afterWalletCreate: () => ['wallets'],
  afterWalletUpdate: (id: string) => ['wallets', ['wallet', id]],
  afterTransaction: (walletId: string) => [
    'wallets',
    ['wallet', walletId],
    ['wallet', walletId, 'balance'],
    'transactions'
  ],
}
```

#### Local State for UI
```typescript
// Context for global UI state
export const AppContext = createContext<{
  theme: 'light' | 'dark'
  sidebarOpen: boolean
  notifications: Notification[]
}>()

// Component-level state for forms and UI interactions
const [isModalOpen, setIsModalOpen] = useState(false)
const [selectedWallet, setSelectedWallet] = useState<string | null>(null)
```

## Backend Integration Patterns

### Controller Layer Design

#### HTTP Request Handling
```go
// adapter/controller/walletController.go
func (c *WalletController) CreateWallet(w http.ResponseWriter, r *http.Request) {
    // 1. Validate HTTP method
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    // 2. Parse and validate request body
    var req struct {
        UserID   string `json:"user_id"`
        Name     string `json:"name"`
        Type     string `json:"type"`
        Currency string `json:"currency"`
    }

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }

    // 3. Transform to use case input
    input := command.CreateWalletInput{
        UserID:   req.UserID,
        Name:     req.Name,
        Type:     req.Type,
        Currency: req.Currency,
    }

    // 4. Execute use case
    output := c.createWalletUseCase.Execute(input)
    
    // 5. Transform response
    w.Header().Set("Content-Type", "application/json")
    if output.GetExitCode() != 0 {
        w.WriteHeader(http.StatusBadRequest)
    }
    
    response := map[string]interface{}{
        "id":      output.GetID(),
        "success": output.GetExitCode() == 0,
        "message": output.GetMessage(),
    }
    
    json.NewEncoder(w).Encode(response)
}
```

#### Input Validation Strategy
```go
// application/command/CreateWalletService.go
func (s *CreateWalletService) Execute(input CreateWalletInput) common.UseCaseOutput {
    // 1. Input validation
    if err := s.validateInput(input); err != nil {
        return &CreateWalletOutput{
            BaseOutput: common.BaseOutput{
                ExitCode: 1,
                Message:  err.Error(),
            },
        }
    }

    // 2. Business logic execution
    walletType, err := model.ParseWalletType(input.Type)
    if err != nil {
        return &CreateWalletOutput{
            BaseOutput: common.BaseOutput{
                ExitCode: 1,
                Message:  "Invalid wallet type: " + input.Type,
            },
        }
    }

    // 3. Domain model creation
    wallet, err := model.NewWallet(input.UserID, input.Name, walletType, input.Currency)
    if err != nil {
        return &CreateWalletOutput{
            BaseOutput: common.BaseOutput{
                ExitCode: 1,
                Message:  err.Error(),
            },
        }
    }

    // 4. Persistence
    if err := s.walletRepo.Save(wallet); err != nil {
        return &CreateWalletOutput{
            BaseOutput: common.BaseOutput{
                ExitCode: 1,
                Message:  "Failed to save wallet: " + err.Error(),
            },
        }
    }

    // 5. Success response
    return &CreateWalletOutput{
        BaseOutput: common.BaseOutput{
            ExitCode: 0,
            Message:  "Wallet created successfully",
        },
        WalletID: wallet.ID,
    }
}

func (s *CreateWalletService) validateInput(input CreateWalletInput) error {
    if input.UserID == "" {
        return errors.New("user ID is required")
    }
    if strings.TrimSpace(input.Name) == "" {
        return errors.New("wallet name is required")
    }
    if input.Currency == "" || len(input.Currency) != 3 {
        return errors.New("currency must be 3 characters (ISO 4217)")
    }
    return nil
}
```

### Use Case Pattern Implementation

#### Command Use Cases
```go
// application/usecase/usecases.go - Interface definitions
type CreateWalletUseCase interface {
    Execute(input command.CreateWalletInput) common.UseCaseOutput
}

type AddExpenseUseCase interface {
    Execute(input command.AddExpenseInput) common.UseCaseOutput
}

// application/command/CreateWalletService.go - Implementation
type CreateWalletService struct {
    walletRepo repository.WalletRepository
}

func NewCreateWalletService(walletRepo repository.WalletRepository) *CreateWalletService {
    return &CreateWalletService{
        walletRepo: walletRepo,
    }
}

func (s *CreateWalletService) Execute(input CreateWalletInput) common.UseCaseOutput {
    // Use case implementation...
}
```

#### Query Use Cases
```go
// application/query/GetWalletBalanceService.go
type GetWalletBalanceService struct {
    walletRepo repository.WalletRepository
}

func (s *GetWalletBalanceService) Execute(input GetWalletBalanceInput) *GetWalletBalanceOutput {
    // 1. Find wallet
    wallet, err := s.walletRepo.FindByID(input.WalletID)
    if err != nil {
        return &GetWalletBalanceOutput{
            BaseOutput: common.BaseOutput{
                ExitCode: 1,
                Message:  "Failed to find wallet: " + err.Error(),
            },
        }
    }

    if wallet == nil {
        return &GetWalletBalanceOutput{
            BaseOutput: common.BaseOutput{
                ExitCode: 1,  
                Message:  "Wallet not found",
            },
        }
    }

    // 2. Return balance information
    return &GetWalletBalanceOutput{
        BaseOutput: common.BaseOutput{
            ExitCode: 0,
            Message:  "Balance retrieved successfully",
        },
        WalletID: wallet.ID,
        Balance:  wallet.Balance.Amount,
        Currency: wallet.Balance.Currency,
    }
}
```

### Repository Bridge Pattern

#### Application Layer Bridge
```go
// application/repository/WalletRepositoryImpl.go
type WalletRepositoryImpl struct {
    peer   WalletRepositoryPeer    // Interface to layer 3
    mapper *mapper.WalletMapper    // Data transformation
}

func (r *WalletRepositoryImpl) Save(wallet *model.Wallet) error {
    // 1. Transform domain model to data structure
    data := r.mapper.ToData(wallet)
    
    // 2. Use peer interface to persist (layer 3)
    return r.peer.SaveData(data)
}

func (r *WalletRepositoryImpl) FindByID(id string) (*model.Wallet, error) {
    // 1. Use peer interface to retrieve data (layer 3)
    data, err := r.peer.FindDataByID(id)
    if err != nil {
        return nil, err
    }
    
    if data == nil {
        return nil, nil
    }
    
    // 2. Transform data structure to domain model
    return r.mapper.ToDomain(*data)
}
```

#### Peer Interface Definition
```go
// application/repository/Repository.go
type WalletRepositoryPeer interface {
    SaveData(data mapper.WalletData) error
    FindDataByID(id string) (*mapper.WalletData, error)
    DeleteData(id string) error
}
```

#### Infrastructure Implementation
```go
// frameworks/database/postgresWalletRepository.go
type PostgresWalletRepositoryPeer struct {
    db *sql.DB
}

func (r *PostgresWalletRepositoryPeer) SaveData(data mapper.WalletData) error {
    query := `
        INSERT INTO wallets (id, user_id, name, type, balance_amount, balance_currency, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        ON CONFLICT (id) 
        DO UPDATE SET 
            name = EXCLUDED.name,
            type = EXCLUDED.type,
            balance_amount = EXCLUDED.balance_amount,
            balance_currency = EXCLUDED.balance_currency,
            updated_at = EXCLUDED.updated_at
    `
    
    _, err := r.db.Exec(query,
        data.ID,
        data.UserID,
        data.Name,
        data.Type,
        data.BalanceAmount,
        data.BalanceCurrency,
        data.CreatedAt,
        data.UpdatedAt,
    )
    
    return err
}

func (r *PostgresWalletRepositoryPeer) FindDataByID(id string) (*mapper.WalletData, error) {
    query := `
        SELECT id, user_id, name, type, balance_amount, balance_currency, created_at, updated_at
        FROM wallets 
        WHERE id = $1
    `
    
    var data mapper.WalletData
    err := r.db.QueryRow(query, id).Scan(
        &data.ID,
        &data.UserID,
        &data.Name,
        &data.Type,
        &data.BalanceAmount,
        &data.BalanceCurrency,
        &data.CreatedAt,
        &data.UpdatedAt,
    )
    
    if err == sql.ErrNoRows {
        return nil, nil
    }
    
    if err != nil {
        return nil, err
    }
    
    return &data, nil
}
```

## Data Transformation Layer

### Mapper Implementation

#### Domain to Data Transformation
```go
// application/mapper/WalletMapper.go
type WalletData struct {
    ID               string
    UserID           string
    Name             string
    Type             string
    BalanceAmount    int64
    BalanceCurrency  string
    CreatedAt        time.Time
    UpdatedAt        time.Time
}

type WalletMapper struct{}

func NewWalletMapper() *WalletMapper {
    return &WalletMapper{}
}

func (m *WalletMapper) ToData(wallet *model.Wallet) WalletData {
    return WalletData{
        ID:               wallet.ID,
        UserID:           wallet.UserID,
        Name:             wallet.Name,
        Type:             string(wallet.Type),
        BalanceAmount:    wallet.Balance.Amount,
        BalanceCurrency:  wallet.Balance.Currency,
        CreatedAt:        wallet.CreatedAt,
        UpdatedAt:        wallet.UpdatedAt,
    }
}

func (m *WalletMapper) ToDomain(data WalletData) (*model.Wallet, error) {
    // 1. Parse wallet type
    walletType, err := model.ParseWalletType(data.Type)
    if err != nil {
        return nil, err
    }

    // 2. Create Money value object
    balance, err := model.NewMoney(data.BalanceAmount, data.BalanceCurrency)
    if err != nil {
        return nil, err
    }

    // 3. Create Wallet domain model
    wallet := &model.Wallet{
        ID:        data.ID,
        UserID:    data.UserID,
        Name:      data.Name,
        Type:      walletType,
        Balance:   *balance,
        CreatedAt: data.CreatedAt,
        UpdatedAt: data.UpdatedAt,
    }

    return wallet, nil
}
```

### Frontend-Backend Type Alignment

#### Money Value Handling
```typescript
// Frontend: utils/format.ts
export const formatMoney = (amount: number, currency: string): string => {
  // amount comes from backend in cents
  const dollars = amount / 100
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: currency,
  }).format(dollars)
}

export const parseMoney = (value: string): number => {
  // Convert user input to cents for backend
  const cleaned = value.replace(/[^0-9.-]/g, '')
  const dollars = parseFloat(cleaned)
  return Math.round(dollars * 100) // Convert to cents
}
```

```go
// Backend: domain/model/money.go
func NewMoney(amountInCents int64, currency string) (*Money, error) {
    if amountInCents < 0 {
        return nil, errors.New("money amount cannot be negative")
    }
    
    if len(currency) != 3 {
        return nil, errors.New("currency must be 3 characters (ISO 4217)")
    }
    
    return &Money{
        Amount:   amountInCents,
        Currency: strings.ToUpper(currency),
    }, nil
}

func (m Money) ToDollars() float64 {
    return float64(m.Amount) / 100
}
```

#### Date/Time Handling
```typescript
// Frontend: Date handling
export const formatDate = (dateString: string): string => {
  return new Date(dateString).toLocaleDateString()
}

export const parseDate = (date: Date): string => {
  return date.toISOString().split('T')[0] // YYYY-MM-DD format for backend
}
```

```go
// Backend: Date handling
func ParseDate(dateStr string) (time.Time, error) {
    if dateStr == "" {
        return time.Now(), nil
    }
    
    return time.Parse("2006-01-02", dateStr)
}
```

## Error Handling Strategy

### Backend Error Patterns

#### Domain Error Types
```go
// domain/errors/errors.go
type DomainError struct {
    Code    string
    Message string
    Details string
}

func (e DomainError) Error() string {
    return e.Message
}

var (
    ErrInsufficientBalance = DomainError{
        Code:    "INSUFFICIENT_BALANCE",
        Message: "Insufficient wallet balance",
    }
    
    ErrInvalidCurrency = DomainError{
        Code:    "INVALID_CURRENCY", 
        Message: "Currency mismatch or invalid format",
    }
    
    ErrWalletNotFound = DomainError{
        Code:    "WALLET_NOT_FOUND",
        Message: "Wallet not found",
    }
)
```

#### Use Case Error Handling
```go
// application/command/AddExpenseService.go
func (s *AddExpenseService) Execute(input AddExpenseInput) common.UseCaseOutput {
    // 1. Find wallet
    wallet, err := s.walletRepo.FindByID(input.WalletID)
    if err != nil {
        return s.createErrorOutput("DATABASE_ERROR", err.Error())
    }
    
    if wallet == nil {
        return s.createErrorOutput("WALLET_NOT_FOUND", "Wallet not found")
    }

    // 2. Create money object
    amount, err := model.NewMoney(input.Amount, input.Currency)
    if err != nil {
        return s.createErrorOutput("INVALID_CURRENCY", err.Error())
    }

    // 3. Add expense (domain logic)
    expense, err := wallet.AddExpense(*amount, input.SubcategoryID, input.Description, input.Date)
    if err != nil {
        // Domain errors are already properly typed
        if strings.Contains(err.Error(), "insufficient balance") {
            return s.createErrorOutput("INSUFFICIENT_BALANCE", err.Error())
        }
        return s.createErrorOutput("DOMAIN_ERROR", err.Error())
    }

    // 4. Save wallet
    if err := s.walletRepo.Save(wallet); err != nil {
        return s.createErrorOutput("DATABASE_ERROR", "Failed to save wallet: " + err.Error())
    }

    return &AddExpenseOutput{
        BaseOutput: common.BaseOutput{
            ExitCode: 0,
            Message:  "Expense added successfully",
        },
        ExpenseID: expense.ID,
    }
}

func (s *AddExpenseService) createErrorOutput(code, message string) *AddExpenseOutput {
    return &AddExpenseOutput{
        BaseOutput: common.BaseOutput{
            ExitCode: 1,
            Message:  message,
        },
        Error: &errors.DomainError{
            Code:    code,
            Message: message,
        },
    }
}
```

### Frontend Error Handling

#### Service Layer Error Processing
```typescript
// services/api.ts - Error processing
export const apiRequest = {
  post: async <T>(url: string, data?: any): Promise<ApiResponse<T>> => {
    try {
      const response = await api.post(url, data)
      return {
        success: true,
        data: response.data,
      }
    } catch (error: any) {
      // Process different error types
      if (error.response) {
        // Backend returned error response
        const backendError = error.response.data?.error
        return {
          success: false,
          error: {
            code: backendError?.code || 'BACKEND_ERROR',
            message: backendError?.message || error.response.data?.message || 'Server error',
            details: backendError?.details,
          }
        }
      } else if (error.request) {
        // Network error
        return {
          success: false,
          error: {
            code: 'NETWORK_ERROR',
            message: 'Unable to connect to server',
          }
        }
      } else {
        // Request setup error
        return {
          success: false,
          error: {
            code: 'REQUEST_ERROR',
            message: error.message,
          }
        }
      }
    }
  },
}
```

#### Component Error Display
```typescript
// hooks/useErrorHandler.ts
import { useState, useCallback } from 'react'
import type { ApiError } from '@/types'

export const useErrorHandler = () => {
  const [error, setError] = useState<ApiError | null>(null)

  const handleError = useCallback((error: ApiError) => {
    setError(error)
    
    // Log error for debugging
    console.error('API Error:', error)
    
    // Could also send to error tracking service
    // errorTracker.captureException(error)
  }, [])

  const clearError = useCallback(() => {
    setError(null)
  }, [])

  const getErrorMessage = useCallback((error: ApiError): string => {
    switch (error.code) {
      case 'INSUFFICIENT_BALANCE':
        return 'Insufficient funds in wallet'
      case 'WALLET_NOT_FOUND':
        return 'Wallet not found'
      case 'INVALID_CURRENCY':
        return 'Invalid currency format'
      case 'NETWORK_ERROR':
        return 'Connection error. Please check your internet connection.'
      default:
        return error.message || 'An unexpected error occurred'
    }
  }, [])

  return {
    error,
    handleError,
    clearError,
    getErrorMessage,
  }
}
```

#### Error Boundary Implementation
```typescript
// components/ErrorBoundary.tsx
import React, { Component, ErrorInfo, ReactNode } from 'react'

interface Props {
  children: ReactNode
  fallback?: ReactNode
}

interface State {
  hasError: boolean
  error?: Error
}

class ErrorBoundary extends Component<Props, State> {
  public state: State = {
    hasError: false
  }

  public static getDerivedStateFromError(error: Error): State {
    return { hasError: true, error }
  }

  public componentDidCatch(error: Error, errorInfo: ErrorInfo) {
    console.error('Uncaught error:', error, errorInfo)
    
    // Could send to error tracking service
    // errorTracker.captureException(error, { extra: errorInfo })
  }

  public render() {
    if (this.state.hasError) {
      return this.props.fallback || (
        <div className="min-h-screen flex items-center justify-center">
          <div className="text-center">
            <h2 className="text-xl font-semibold text-red-600">Something went wrong</h2>
            <p className="text-gray-600 mt-2">
              {this.state.error?.message || 'An unexpected error occurred'}
            </p>
            <button 
              className="mt-4 px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
              onClick={() => this.setState({ hasError: false })}
            >
              Try again
            </button>
          </div>
        </div>
      )
    }

    return this.props.children
  }
}

export default ErrorBoundary
```

## Authentication & Authorization

### Current State: Development Mode
Currently, the system operates without authentication for development purposes:

```typescript
// Frontend: Hardcoded user context
const getCurrentUser = () => ({
  id: 'user-1',
  name: 'Development User',
  email: 'dev@example.com'
})
```

```go
// Backend: No authentication middleware
// All endpoints accessible without tokens
```

### Planned Authentication Architecture

#### JWT-Based Authentication
```typescript
// Frontend: Auth context and token management
export const AuthContext = createContext<{
  user: User | null
  token: string | null
  login: (email: string, password: string) => Promise<void>
  logout: () => void
  isAuthenticated: boolean
}>()

export const useAuth = () => {
  const context = useContext(AuthContext)
  if (!context) {
    throw new Error('useAuth must be used within AuthProvider')
  }
  return context
}

// Token management
export const tokenService = {
  getToken: (): string | null => localStorage.getItem('authToken'),
  setToken: (token: string): void => localStorage.setItem('authToken', token),
  removeToken: (): void => localStorage.removeItem('authToken'),
  isTokenExpired: (token: string): boolean => {
    try {
      const payload = JSON.parse(atob(token.split('.')[1]))
      return payload.exp * 1000 < Date.now()
    } catch {
      return true
    }
  }
}
```

```go
// Backend: JWT middleware (planned)
type JWTMiddleware struct {
    secretKey []byte
}

func (m *JWTMiddleware) ValidateToken(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Extract token from Authorization header
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, "Missing authorization header", http.StatusUnauthorized)
            return
        }
        
        // Validate JWT token
        token := strings.TrimPrefix(authHeader, "Bearer ")
        claims, err := m.parseToken(token)
        if err != nil {
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }
        
        // Add user context to request
        ctx := context.WithValue(r.Context(), "userID", claims.UserID)
        next.ServeHTTP(w, r.WithContext(ctx))
    }
}
```

## Testing Integration Points

### Frontend Testing Strategy

#### API Service Testing
```typescript
// __tests__/services/walletService.test.ts
import { rest } from 'msw'
import { setupServer } from 'msw/node'
import { walletService } from '@/services/walletService'

const server = setupServer(
  rest.post('/api/v1/wallets', (req, res, ctx) => {
    return res(ctx.json({
      id: 'wallet-123',
      success: true,
      message: 'Wallet created successfully'
    }))
  }),
  
  rest.get('/api/v1/wallets/:id/balance', (req, res, ctx) => {
    const { id } = req.params
    return res(ctx.json({
      walletId: id,
      balance: 10000, // $100.00 in cents
      currency: 'USD',
      success: true,
      message: 'Balance retrieved'
    }))
  })
)

beforeAll(() => server.listen())
afterEach(() => server.resetHandlers())
afterAll(() => server.close())

describe('WalletService', () => {
  it('creates wallet successfully', async () => {
    const response = await walletService.create({
      user_id: 'user-1',
      name: 'Test Wallet',
      type: 'CASH',
      currency: 'USD'
    })
    
    expect(response.success).toBe(true)
    expect(response.data?.id).toBe('wallet-123')
  })
  
  it('gets wallet balance successfully', async () => {
    const response = await walletService.getBalance('wallet-123')
    
    expect(response.success).toBe(true)
    expect(response.data?.balance).toBe(10000)
    expect(response.data?.currency).toBe('USD')
  })
})
```

#### Component Integration Testing
```typescript
// __tests__/components/WalletForm.test.tsx
import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import { QueryClient, QueryClientProvider } from 'react-query'
import { rest } from 'msw'
import { setupServer } from 'msw/node'
import WalletForm from '@/components/WalletForm'

const server = setupServer(
  rest.post('/api/v1/wallets', (req, res, ctx) => {
    return res(ctx.json({
      id: 'new-wallet-id',
      success: true,
      message: 'Wallet created successfully'
    }))
  })
)

beforeAll(() => server.listen())
afterEach(() => server.resetHandlers())
afterAll(() => server.close())

const renderWithQuery = (component: React.ReactElement) => {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: { retry: false },
      mutations: { retry: false },
    },
  })
  
  return render(
    <QueryClientProvider client={queryClient}>
      {component}
    </QueryClientProvider>
  )
}

describe('WalletForm', () => {
  it('creates wallet when form is submitted', async () => {
    renderWithQuery(<WalletForm />)
    
    // Fill form
    fireEvent.change(screen.getByLabelText(/wallet name/i), {
      target: { value: 'My Test Wallet' }
    })
    
    fireEvent.change(screen.getByLabelText(/currency/i), {
      target: { value: 'USD' }
    })
    
    // Submit form
    fireEvent.click(screen.getByRole('button', { name: /create wallet/i }))
    
    // Wait for success
    await waitFor(() => {
      expect(screen.getByDisplayValue('')).toBeInTheDocument() // Form reset
    })
  })
})
```

### Backend Testing Strategy

#### Use Case Testing
```go
// test/usecase/create_wallet_test.go
func TestCreateWallet_Success(t *testing.T) {
    // Arrange
    mockRepo := &mocks.MockWalletRepository{}
    service := command.NewCreateWalletService(mockRepo)
    
    input := command.CreateWalletInput{
        UserID:   "user-1",
        Name:     "Test Wallet",
        Type:     "CASH",
        Currency: "USD",
    }
    
    mockRepo.On("Save", mock.AnythingOfType("*model.Wallet")).Return(nil)
    
    // Act
    output := service.Execute(input)
    
    // Assert
    assert.Equal(t, 0, output.GetExitCode())
    assert.Equal(t, "Wallet created successfully", output.GetMessage())
    assert.NotEmpty(t, output.GetID())
    mockRepo.AssertExpectations(t)
}

func TestCreateWallet_InvalidInput(t *testing.T) {
    // Arrange
    mockRepo := &mocks.MockWalletRepository{}
    service := command.NewCreateWalletService(mockRepo)
    
    input := command.CreateWalletInput{
        UserID:   "",  // Invalid: empty user ID
        Name:     "Test Wallet",
        Type:     "CASH",
        Currency: "USD",
    }
    
    // Act
    output := service.Execute(input)
    
    // Assert
    assert.Equal(t, 1, output.GetExitCode())
    assert.Contains(t, output.GetMessage(), "user ID is required")
    mockRepo.AssertNotCalled(t, "Save")
}
```

#### Integration Testing with Test Database
```go
// test/integration/wallet_integration_test.go
func TestWalletIntegration(t *testing.T) {
    // Setup test database
    testDB := setupTestDatabase()
    defer cleanupTestDatabase(testDB)
    
    // Create real implementations
    peer := database.NewPostgresWalletRepositoryPeer(testDB)
    repo := repository.NewWalletRepositoryImpl(peer)
    
    // Test create wallet
    wallet, err := model.NewWallet("user-1", "Integration Test Wallet", model.WalletTypeCash, "USD")
    require.NoError(t, err)
    
    err = repo.Save(wallet)
    require.NoError(t, err)
    
    // Test find wallet
    foundWallet, err := repo.FindByID(wallet.ID)
    require.NoError(t, err)
    require.NotNil(t, foundWallet)
    
    assert.Equal(t, wallet.ID, foundWallet.ID)
    assert.Equal(t, wallet.Name, foundWallet.Name)
    assert.Equal(t, wallet.Type, foundWallet.Type)
    assert.Equal(t, wallet.Balance.Amount, foundWallet.Balance.Amount)
    assert.Equal(t, wallet.Balance.Currency, foundWallet.Balance.Currency)
}

func setupTestDatabase() *sql.DB {
    // Connect to test database
    db, err := sql.Open("postgres", os.Getenv("TEST_DATABASE_URL"))
    if err != nil {
        panic(err)
    }
    
    // Run schema
    schema, err := ioutil.ReadFile("../internal/accounting/frameworks/database/schema.sql")
    if err != nil {
        panic(err)
    }
    
    _, err = db.Exec(string(schema))
    if err != nil {
        panic(err)
    }
    
    return db
}
```

### End-to-End Testing

#### API Contract Testing
```typescript
// e2e/api-contract.test.ts
describe('API Contract Tests', () => {
  const baseURL = 'http://localhost:8080/api/v1'
  
  it('POST /wallets follows contract', async () => {
    const response = await fetch(`${baseURL}/wallets`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        user_id: 'user-1',
        name: 'Contract Test Wallet',
        type: 'CASH',
        currency: 'USD'
      })
    })
    
    expect(response.status).toBe(200)
    
    const data = await response.json()
    expect(data).toHaveProperty('id')
    expect(data).toHaveProperty('success', true)
    expect(data).toHaveProperty('message')
    expect(typeof data.id).toBe('string')
    expect(data.id).toMatch(/^[0-9a-f-]{36}$/) // UUID format
  })
  
  it('GET /wallets/{id}/balance follows contract', async () => {
    // First create wallet
    const createResponse = await fetch(`${baseURL}/wallets`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        user_id: 'user-1',
        name: 'Balance Test Wallet',
        type: 'CASH',
        currency: 'USD'
      })
    })
    
    const createData = await createResponse.json()
    
    // Then get balance
    const balanceResponse = await fetch(`${baseURL}/wallets/${createData.id}/balance`)
    expect(balanceResponse.status).toBe(200)
    
    const balanceData = await balanceResponse.json()
    expect(balanceData).toHaveProperty('walletId', createData.id)
    expect(balanceData).toHaveProperty('balance', 0)
    expect(balanceData).toHaveProperty('currency', 'USD')
    expect(balanceData).toHaveProperty('success', true)
  })
})
```

## Development Workflow

### Local Development Setup

#### Complete Environment Setup
```bash
# 1. Repository setup
git clone <repository-url>
cd accounting-app-project

# 2. Backend development environment
cd accountingApp
cp .env.example .env

# Edit .env if needed:
# DATABASE_URL=postgres://postgres:password@localhost:5432/accountingdb?sslmode=disable
# PORT=8080

# Start database
./scripts/start-dev.sh
# Or manually: docker-compose up -d postgres

# Install Go dependencies
go mod tidy

# Run backend tests
go test ./...

# Start backend server
go run cmd/accoountingApp/main.go

# 3. Frontend development environment (new terminal)
cd frontend

# Install Node dependencies  
npm install

# Start frontend dev server
npm run dev

# 4. Verification
curl http://localhost:8080/health     # Backend health
open http://localhost:3000            # Frontend UI
```

#### Development Scripts
```bash
# Backend scripts (accountingApp/)
./scripts/start-dev.sh              # Start PostgreSQL container
./scripts/reset-db.sh               # Reset database (planned)
./scripts/test.sh                   # Run all tests (planned)

# Frontend scripts (frontend/)
npm run dev                         # Development server with HMR
npm run build                       # Production build
npm run preview                     # Preview production build
npm run lint                        # ESLint code checking
npm run type-check                  # TypeScript checking
```

### Development Workflow Patterns

#### Feature Development Flow
```mermaid
graph TD
    A[Create Feature Branch] --> B[Backend Development]
    B --> C[Write Domain Tests]
    C --> D[Implement Domain Logic]
    D --> E[Write Use Case Tests]
    E --> F[Implement Use Cases]
    F --> G[Write Controller Tests]
    G --> H[Implement Controllers]
    H --> I[Frontend Development]
    I --> J[Create TypeScript Types]
    J --> K[Implement API Service]
    K --> L[Write Component Tests]
    L --> M[Implement Components]
    M --> N[Integration Testing]
    N --> O[Manual Testing]
    O --> P[Code Review]
    P --> Q[Merge to Main]
```

#### Daily Development Cycle
```yaml
Morning Setup:
  - git pull origin main
  - docker-compose up -d postgres  
  - go run cmd/accoountingApp/main.go (terminal 1)
  - npm run dev (terminal 2)
  - Verify both services running

Development Loop:
  - Write tests (TDD approach)
  - Implement feature (backend first)
  - Test with curl/Postman
  - Implement frontend integration
  - Test in browser
  - Fix issues and iterate

End of Day:
  - Run full test suite
  - git add . && git commit
  - git push origin feature-branch
  - docker-compose down (optional)
```

### Cross-Service Development

#### API-First Development
```yaml
1. Define API Contract:
   - Document endpoint in this guide
   - Create TypeScript types  
   - Create Go request/response structs

2. Backend Implementation:
   - Write use case tests
   - Implement domain logic
   - Implement use case
   - Implement controller
   - Test with curl

3. Frontend Implementation:
   - Create API service method
   - Write service tests
   - Implement React hooks
   - Create components
   - Test integration

4. Integration Testing:
   - End-to-end user flow testing
   - Error scenario testing
   - Performance testing
```

#### Debugging Integration Issues

**Common Issues & Solutions**:

```yaml
CORS Issues:
  Symptoms: Browser blocks requests
  Solution: Vite proxy handles CORS in development
  Check: Ensure proxy configuration correct

Type Mismatches:
  Symptoms: Frontend/backend data format issues
  Solution: Verify mapper implementations
  Check: Money handling (cents vs dollars)

Network Errors:
  Symptoms: Connection refused errors
  Solution: Check both servers running
  Check: Backend on :8080, frontend on :3000

Database Connection:
  Symptoms: Backend database errors
  Solution: Ensure PostgreSQL container running
  Check: docker-compose ps, connection string

Authentication Errors (Future):
  Symptoms: 401 responses
  Solution: Check token handling
  Check: Authorization header format
```

#### Performance Monitoring During Development

```typescript
// Frontend: Request timing
api.interceptors.request.use(config => {
  config.metadata = { startTime: new Date() }
  return config
})

api.interceptors.response.use(response => {
  const endTime = new Date()
  const duration = endTime.getTime() - response.config.metadata.startTime.getTime()
  console.log(`API Request: ${response.config.url} took ${duration}ms`)
  return response
})
```

```go
// Backend: Request logging middleware
func LoggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        
        // Call next handler
        next.ServeHTTP(w, r)
        
        // Log request details
        log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
    }
}
```

## Troubleshooting Guide

### Common Integration Issues

#### 1. Service Communication Problems

**Problem**: Frontend cannot reach backend API
```
Error: Network Error
fetch failed to http://localhost:3000/api/v1/wallets
```

**Diagnosis**:
```bash
# Check if backend is running
curl http://localhost:8080/health

# Check if frontend proxy is working  
curl http://localhost:3000/api/v1/health

# Check Vite proxy configuration
cat frontend/vite.config.ts
```

**Solutions**:
- Ensure backend running on port 8080
- Verify Vite proxy configuration
- Check firewall settings
- Restart frontend dev server

#### 2. Database Connection Issues

**Problem**: Backend cannot connect to PostgreSQL
```
Error: dial tcp 127.0.0.1:5432: connect: connection refused
```

**Diagnosis**:
```bash
# Check if PostgreSQL container is running
docker-compose ps

# Check PostgreSQL logs
docker-compose logs postgres

# Test connection manually
psql -h localhost -p 5432 -U postgres -d accountingdb
```

**Solutions**:
- Start PostgreSQL: `docker-compose up -d postgres`
- Reset database: `docker-compose down -v && docker-compose up -d postgres`
- Check DATABASE_URL in .env file
- Verify schema initialization

#### 3. Type Safety Issues

**Problem**: Frontend/backend type mismatches
```typescript
// Frontend expects
interface Money {
  amount: number    // dollars
  currency: string
}

// Backend returns  
{
  "amount": 1000,   // cents
  "currency": "USD"
}
```

**Solutions**:
- Use consistent units (cents) across both services
- Implement proper transformation functions
- Add runtime validation
- Update TypeScript types

#### 4. State Synchronization Issues

**Problem**: UI shows stale data after mutations
```typescript
// User creates wallet but list doesn't update
const createWallet = useCreateWallet()
await createWallet.mutateAsync(data)
// Wallet list still shows old data
```

**Solutions**:
```typescript
// Ensure proper cache invalidation
export const useCreateWallet = () => {
  const queryClient = useQueryClient()
  
  return useMutation(walletService.create, {
    onSuccess: () => {
      queryClient.invalidateQueries('wallets')  // This is crucial
    }
  })
}
```

### Performance Troubleshooting

#### Slow API Responses

**Diagnosis Tools**:
```bash
# Backend performance
time curl -X POST http://localhost:8080/api/v1/wallets \
  -H "Content-Type: application/json" \
  -d '{"user_id":"user1","name":"Test","type":"CASH","currency":"USD"}'

# Database query performance
docker exec -it accounting-postgres psql -U postgres -d accountingdb -c "EXPLAIN ANALYZE SELECT * FROM wallets;"
```

**Common Causes & Solutions**:
- **Missing database indexes**: Add indexes on frequently queried columns
- **N+1 query problems**: Implement proper eager loading
- **Large response payloads**: Implement pagination
- **Blocking operations**: Use async patterns

#### Memory Leaks

**Frontend Monitoring**:
```typescript
// Monitor component unmounting
useEffect(() => {
  return () => {
    // Cleanup subscriptions, timers, etc.
    console.log('Component unmounted, cleaning up...')
  }
}, [])
```

**Backend Monitoring**:
```go
// Monitor goroutine leaks
import _ "net/http/pprof"

func main() {
    go func() {
        log.Println(http.ListenAndServe("localhost:6060", nil))
    }()
    
    // Your application code
}
```

### Security Troubleshooting

#### CORS Issues (Production)

**Problem**: CORS errors in production environment
```
Access to fetch at 'https://api.example.com' from origin 'https://app.example.com' has been blocked by CORS policy
```

**Solutions**:
```go
// Backend: Add CORS middleware
func CORSMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        
        next.ServeHTTP(w, r)
    }
}
```

#### Authentication Token Issues (Planned)

**Problem**: Token expiration handling
```typescript
// Token expired but frontend doesn't handle it properly
axios.interceptors.response.use(
  response => response,
  error => {
    if (error.response?.status === 401) {
      // Redirect to login
      tokenService.removeToken()
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)
```

### Development Environment Issues

#### Port Conflicts

**Problem**: Ports already in use
```
Error: listen EADDRINUSE: address already in use :::8080
Error: listen EADDRINUSE: address already in use :::3000
```

**Solutions**:
```bash
# Find process using port
lsof -ti :8080
lsof -ti :3000

# Kill process
kill -9 $(lsof -ti :8080)

# Use different ports
export PORT=8081  # Backend
npm run dev -- --port 3001  # Frontend
```

#### Docker Issues

**Problem**: Container startup failures
```bash
# Check container status
docker-compose ps

# View logs
docker-compose logs postgres
docker-compose logs pgadmin

# Reset containers
docker-compose down -v
docker system prune -f
docker-compose up -d
```

### Error Monitoring Setup (Production)

#### Frontend Error Tracking
```typescript
// Sentry integration (example)
import * as Sentry from '@sentry/react'

Sentry.init({
  dsn: process.env.VITE_SENTRY_DSN,
  environment: process.env.NODE_ENV,
})

// Error boundary integration
const SentryErrorBoundary = Sentry.withErrorBoundary(App, {
  fallback: ErrorFallback,
})
```

#### Backend Error Tracking
```go
// Structured logging with context
import (
    "github.com/sirupsen/logrus"
    "context"
)

func logError(ctx context.Context, err error, message string) {
    entry := logrus.WithFields(logrus.Fields{
        "error": err.Error(),
        "user_id": getUserIDFromContext(ctx),
        "request_id": getRequestIDFromContext(ctx),
    })
    entry.Error(message)
}
```

This comprehensive integration guide provides the technical foundation for developing and maintaining the accounting app's multi-service architecture. Following these patterns ensures consistent, reliable integration between the React frontend and Go backend services.