# Frontend Integration Guide

Complete guide for integrating with the Accounting App backend API from frontend applications.

> **⚠️ Updated**: This guide now reflects the actual working API implementation as of January 2025.

## 🚀 Quick Start

### 1. API Client Setup

Create a centralized API client for consistent error handling and configuration:

```typescript
// src/services/apiClient.ts
const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080';

export class ApiClient {
  private baseUrl: string;
  
  constructor(baseUrl = API_BASE_URL) {
    this.baseUrl = baseUrl;
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<T> {
    const url = `${this.baseUrl}${endpoint}`;
    
    const config: RequestInit = {
      headers: {
        'Content-Type': 'application/json',
        ...options.headers,
      },
      ...options,
    };

    try {
      const response = await fetch(url, config);
      const data = await response.json();

      if (!data.success) {
        throw new ApiError(data.error || 'Unknown API error', response.status);
      }

      return data.data || data;
    } catch (error) {
      if (error instanceof ApiError) {
        throw error;
      }
      throw new ApiError('Network error', 0);
    }
  }

  async get<T>(endpoint: string): Promise<T> {
    return this.request<T>(endpoint, { method: 'GET' });
  }

  async post<T>(endpoint: string, data: any): Promise<T> {
    return this.request<T>(endpoint, {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  async put<T>(endpoint: string, data: any): Promise<T> {
    return this.request<T>(endpoint, {
      method: 'PUT',
      body: JSON.stringify(data),
    });
  }

  async delete<T>(endpoint: string): Promise<T> {
    return this.request<T>(endpoint, { method: 'DELETE' });
  }
}

// Custom error class for API errors
export class ApiError extends Error {
  constructor(
    message: string,
    public status: number,
    public originalError?: Error
  ) {
    super(message);
    this.name = 'ApiError';
  }
}

// Export singleton instance
export const apiClient = new ApiClient();
```

### 2. Type Definitions

Define TypeScript interfaces for API data:

```typescript
// src/types/api.ts
export interface Wallet {
  id: string;
  user_id: string;
  name: string;
  type: 'CASH' | 'BANK' | 'CREDIT' | 'INVESTMENT';  // UPPERCASE only
  currency: string;
  balance: {
    amount: number;
    currency: string;
  };
  created_at: string;
  updated_at: string;
  is_fully_loaded?: boolean;
  transactions?: Transaction[];
}

export interface Transaction {
  id: string;
  wallet_id: string;
  type: 'expense' | 'income';
  amount: number;
  currency: string;
  description?: string;
  date: string;
  category_id: string;
}

export interface Category {
  id: string;
  user_id: string;
  name: string;
  type: 'expense' | 'income';
  created_at: string;
}

// Request/Response types
export interface CreateWalletRequest {
  user_id: string;
  name: string;
  type: 'CASH' | 'BANK' | 'CREDIT' | 'INVESTMENT';  // UPPERCASE only
  currency: string;
  initialBalance?: number;
}

export interface UpdateWalletRequest {
  name?: string;
  type?: string;
  currency?: string;
}

export interface AddExpenseRequest {
  wallet_id: string;
  subcategory_id: string;
  amount: number;
  currency: string;
  description?: string;
  date: string;
}

export interface AddIncomeRequest {
  wallet_id: string;
  subcategory_id: string;
  amount: number;
  currency: string;
  description?: string;
  date: string;
}

export interface CreateCategoryRequest {
  user_id: string;
  name: string;
}
```

## 🏦 Wallet Service Implementation

Complete service for wallet operations:

```typescript
// src/services/walletService.ts
import { apiClient } from './apiClient';
import type {
  Wallet,
  CreateWalletRequest,
  UpdateWalletRequest,
} from '../types/api';

export class WalletService {
  async createWallet(walletData: CreateWalletRequest): Promise<{ id: string; message: string }> {
    return apiClient.post('/api/v1/wallets', walletData);
  }

  async getUserWallets(userID: string): Promise<Wallet[]> {
    return apiClient.get(`/api/v1/wallets?userID=${encodeURIComponent(userID)}`);
  }

  async getWallet(walletID: string, includeTransactions = false): Promise<{ data: Wallet }> {
    const queryParam = includeTransactions ? '?includeTransactions=true' : '';
    return apiClient.get(`/api/v1/wallets/${encodeURIComponent(walletID)}${queryParam}`);
  }

  async updateWallet(walletID: string, updates: UpdateWalletRequest): Promise<{ message: string }> {
    return apiClient.put(`/api/v1/wallets/${encodeURIComponent(walletID)}`, updates);
  }

  async deleteWallet(walletID: string): Promise<{ message: string }> {
    return apiClient.delete(`/api/v1/wallets/${encodeURIComponent(walletID)}`);
  }

  async getWalletBalance(walletID: string): Promise<{
    walletId: string;
    balance: number;
    currency: string;
    message: string;
  }> {
    return apiClient.get(`/api/v1/wallets/${encodeURIComponent(walletID)}/balance`);
  }
}

export const walletService = new WalletService();
```

## 💸 Transaction Service Implementation

Service for managing transactions:

```typescript
// src/services/transactionService.ts
import { apiClient } from './apiClient';
import type { AddExpenseRequest, AddIncomeRequest } from '../types/api';

export class TransactionService {
  async addExpense(expenseData: AddExpenseRequest): Promise<{ id: string; message: string }> {
    return apiClient.post('/api/v1/expenses', expenseData);
  }

  async addIncome(incomeData: AddIncomeRequest): Promise<{ id: string; message: string }> {
    return apiClient.post('/api/v1/incomes', incomeData);
  }
}

export const transactionService = new TransactionService();
```

## 🏆 Income Service Implementation

Service for managing income queries:

```typescript
// src/services/incomeService.ts
import { apiClient } from './apiClient';
import type { IncomeRecord } from '../types/api';

export interface GetIncomesFilters {
  walletID?: string;
  categoryID?: string;
  startDate?: string;
  endDate?: string;
  minAmount?: number;
  maxAmount?: number;
  description?: string;
}

export class IncomeService {
  async getIncomes(userID: string, filters?: GetIncomesFilters): Promise<{
    data: IncomeRecord[];
    count: number;
    message: string;
  }> {
    const params = new URLSearchParams({ userID });
    
    if (filters) {
      Object.entries(filters).forEach(([key, value]) => {
        if (value !== undefined) {
          params.append(key, value.toString());
        }
      });
    }
    
    return apiClient.get(`/api/v1/incomes?${params}`);
  }
}

export const incomeService = new IncomeService();
```

## 🏷️ Category Service Implementation

Service for category management:

```typescript
// src/services/categoryService.ts
import { apiClient } from './apiClient';
import type { CreateCategoryRequest } from '../types/api';

export class CategoryService {
  async getAllCategories(): Promise<Category[]> {
    return apiClient.get('/api/v1/categories');
  }

  async getExpenseCategories(): Promise<Category[]> {
    return apiClient.get('/api/v1/categories/expense');
  }

  async getIncomeCategories(): Promise<Category[]> {
    return apiClient.get('/api/v1/categories/income');
  }

  async createExpenseCategory(categoryData: CreateCategoryRequest): Promise<{ id: string; message: string }> {
    return apiClient.post('/api/v1/categories/expense', categoryData);
  }

  async createIncomeCategory(categoryData: CreateCategoryRequest): Promise<{ id: string; message: string }> {
    return apiClient.post('/api/v1/categories/income', categoryData);
  }
}

export const categoryService = new CategoryService();
```

## 💰 Currency Utilities

Handle currency conversion and formatting:

```typescript
// src/utils/currency.ts
const CENTS_BASED_CURRENCIES = ['USD', 'EUR', 'GBP', 'CAD', 'AUD', 'TWD'];

export const CurrencyUtils = {
  // Convert display amount to API format (smallest currency unit)
  toApiAmount(displayAmount: number, currency: string): number {
    if (CENTS_BASED_CURRENCIES.includes(currency.toUpperCase())) {
      return Math.round(displayAmount * 100);
    }
    return displayAmount;
  },

  // Convert API amount to display format
  toDisplayAmount(apiAmount: number, currency: string): number {
    if (CENTS_BASED_CURRENCIES.includes(currency.toUpperCase())) {
      return apiAmount / 100;
    }
    return apiAmount;
  },

  // Format amount for display
  formatAmount(amount: number, currency: string): string {
    const displayAmount = this.toDisplayAmount(amount, currency);
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: currency,
    }).format(displayAmount);
  },

  // Parse user input to API format
  parseUserInput(input: string, currency: string): number {
    const numericValue = parseFloat(input.replace(/[^\d.-]/g, ''));
    if (isNaN(numericValue)) return 0;
    return this.toApiAmount(numericValue, currency);
  },
};
```

## 📅 Date Utilities

Handle date formatting:

```typescript
// src/utils/date.ts
export const DateUtils = {
  // Convert Date object to API format
  toApiFormat(date: Date): string {
    return date.toISOString();
  },

  // Convert API date string to Date object
  fromApiFormat(dateString: string): Date {
    return new Date(dateString);
  },

  // Format date for display
  formatForDisplay(date: Date | string): string {
    const dateObj = typeof date === 'string' ? new Date(date) : date;
    return new Intl.DateTimeFormat('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
    }).format(dateObj);
  },

  // Get current date in API format
  now(): string {
    return this.toApiFormat(new Date());
  },
};
```

## ⚛️ React Hooks

Custom hooks for common operations:

```typescript
// src/hooks/useWallets.ts
import { useState, useEffect } from 'react';
import { walletService } from '../services/walletService';
import type { Wallet } from '../types/api';

export const useWallets = (userID: string) => {
  const [wallets, setWallets] = useState<Wallet[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchWallets = async () => {
    try {
      setLoading(true);
      setError(null);
      const result = await walletService.getUserWallets(userID);
      setWallets(result.data);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch wallets');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (userID) {
      fetchWallets();
    }
  }, [userID]);

  return { wallets, loading, error, refetch: fetchWallets };
};

// src/hooks/useWallet.ts
import { useState, useEffect } from 'react';
import { walletService } from '../services/walletService';
import type { Wallet } from '../types/api';

export const useWallet = (walletID: string, includeTransactions = false) => {
  const [wallet, setWallet] = useState<Wallet | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchWallet = async () => {
    try {
      setLoading(true);
      setError(null);
      const result = await walletService.getWallet(walletID, includeTransactions);
      setWallet(result.data);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch wallet');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (walletID) {
      fetchWallet();
    }
  }, [walletID, includeTransactions]);

  return { wallet, loading, error, refetch: fetchWallet };
};
```

## 🎯 React Component Examples

### Wallet List Component

```tsx
// src/components/WalletList.tsx
import React from 'react';
import { useWallets } from '../hooks/useWallets';
import { CurrencyUtils } from '../utils/currency';

interface WalletListProps {
  userID: string;
  onWalletSelect?: (walletID: string) => void;
}

export const WalletList: React.FC<WalletListProps> = ({ userID, onWalletSelect }) => {
  const { wallets, loading, error } = useWallets(userID);

  if (loading) return <div>Loading wallets...</div>;
  if (error) return <div>Error: {error}</div>;
  if (wallets.length === 0) return <div>No wallets found</div>;

  return (
    <div className="wallet-list">
      <h2>Your Wallets</h2>
      {wallets.map((wallet) => (
        <div 
          key={wallet.id} 
          className="wallet-card"
          onClick={() => onWalletSelect?.(wallet.id)}
        >
          <h3>{wallet.name}</h3>
          <p>Type: {wallet.type}</p>
          <p>Balance: {CurrencyUtils.formatAmount(wallet.balance.amount, wallet.balance.currency)}</p>
          <p>Created: {new Date(wallet.created_at).toLocaleDateString()}</p>
        </div>
      ))}
    </div>
  );
};
```

### Create Wallet Form

```tsx
// src/components/CreateWalletForm.tsx
import React, { useState } from 'react';
import { walletService } from '../services/walletService';
import { CurrencyUtils } from '../utils/currency';
import type { CreateWalletRequest } from '../types/api';

interface CreateWalletFormProps {
  userID: string;
  onSuccess?: (walletID: string) => void;
}

export const CreateWalletForm: React.FC<CreateWalletFormProps> = ({ userID, onSuccess }) => {
  const [formData, setFormData] = useState({
    name: '',
    type: 'BANK' as 'CASH' | 'BANK' | 'CREDIT' | 'INVESTMENT',
    currency: 'USD',
    initialBalance: '',
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    try {
      setLoading(true);
      setError(null);

      const requestData: CreateWalletRequest = {
        user_id: userID,
        name: formData.name,
        type: formData.type,
        currency: formData.currency,
      };

      if (formData.initialBalance) {
        requestData.initialBalance = CurrencyUtils.parseUserInput(formData.initialBalance, formData.currency);
      }

      const result = await walletService.createWallet(requestData);
      onSuccess?.(result.id);
      
      // Reset form
      setFormData({ name: '', type: 'BANK', currency: 'USD', initialBalance: '' });
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to create wallet');
    } finally {
      setLoading(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="create-wallet-form">
      <h2>Create New Wallet</h2>
      
      {error && <div className="error">{error}</div>}
      
      <div className="form-group">
        <label htmlFor="name">Wallet Name *</label>
        <input
          type="text"
          id="name"
          value={formData.name}
          onChange={(e) => setFormData(prev => ({ ...prev, name: e.target.value }))}
          required
        />
      </div>

      <div className="form-group">
        <label htmlFor="type">Wallet Type *</label>
        <select
          id="type"
          value={formData.type}
          onChange={(e) => setFormData(prev => ({ ...prev, type: e.target.value as any }))}
          required
        >
          <option value="BANK">Bank Account</option>
          <option value="CASH">Cash</option>
          <option value="CREDIT">Credit Card</option>
          <option value="INVESTMENT">Investment</option>
        </select>
      </div>

      <div className="form-group">
        <label htmlFor="currency">Currency *</label>
        <select
          id="currency"
          value={formData.currency}
          onChange={(e) => setFormData(prev => ({ ...prev, currency: e.target.value }))}
          required
        >
          <option value="USD">USD</option>
          <option value="EUR">EUR</option>
          <option value="TWD">TWD</option>
          <option value="JPY">JPY</option>
        </select>
      </div>

      <div className="form-group">
        <label htmlFor="initialBalance">Initial Balance</label>
        <input
          type="number"
          step="0.01"
          id="initialBalance"
          value={formData.initialBalance}
          onChange={(e) => setFormData(prev => ({ ...prev, initialBalance: e.target.value }))}
          placeholder="0.00"
        />
      </div>

      <button type="submit" disabled={loading}>
        {loading ? 'Creating...' : 'Create Wallet'}
      </button>
    </form>
  );
};
```

### Add Expense Form

```tsx
// src/components/AddExpenseForm.tsx
import React, { useState } from 'react';
import { transactionService } from '../services/transactionService';
import { CurrencyUtils, DateUtils } from '../utils';
import type { AddExpenseRequest } from '../types/api';

interface AddExpenseFormProps {
  walletID: string;
  currency: string;
  onSuccess?: (expenseID: string) => void;
}

export const AddExpenseForm: React.FC<AddExpenseFormProps> = ({ walletID, currency, onSuccess }) => {
  const [formData, setFormData] = useState({
    subcategory_id: '',
    amount: '',
    description: '',
    date: new Date().toISOString().split('T')[0], // YYYY-MM-DD format
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    try {
      setLoading(true);
      setError(null);

      const requestData: AddExpenseRequest = {
        wallet_id: walletID,
        subcategory_id: formData.subcategory_id,
        amount: CurrencyUtils.parseUserInput(formData.amount, currency),
        currency: currency,
        description: formData.description || undefined,
        date: DateUtils.toApiFormat(new Date(formData.date)),
      };

      const result = await transactionService.addExpense(requestData);
      onSuccess?.(result.id);
      
      // Reset form
      setFormData({
        subcategory_id: '',
        amount: '',
        description: '',
        date: new Date().toISOString().split('T')[0],
      });
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to add expense');
    } finally {
      setLoading(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="add-expense-form">
      <h3>Add Expense</h3>
      
      {error && <div className="error">{error}</div>}
      
      <div className="form-group">
        <label htmlFor="subcategory_id">Category *</label>
        <select
          id="subcategory_id"
          value={formData.subcategory_id}
          onChange={(e) => setFormData(prev => ({ ...prev, subcategory_id: e.target.value }))}
          required
        >
          <option value="">Select a category</option>
          {/* You would populate this with actual categories */}
        </select>
      </div>

      <div className="form-group">
        <label htmlFor="amount">Amount ({currency}) *</label>
        <input
          type="number"
          step="0.01"
          id="amount"
          value={formData.amount}
          onChange={(e) => setFormData(prev => ({ ...prev, amount: e.target.value }))}
          required
          placeholder="0.00"
        />
      </div>

      <div className="form-group">
        <label htmlFor="description">Description</label>
        <input
          type="text"
          id="description"
          value={formData.description}
          onChange={(e) => setFormData(prev => ({ ...prev, description: e.target.value }))}
          placeholder="Optional description"
        />
      </div>

      <div className="form-group">
        <label htmlFor="date">Date *</label>
        <input
          type="date"
          id="date"
          value={formData.date}
          onChange={(e) => setFormData(prev => ({ ...prev, date: e.target.value }))}
          required
        />
      </div>

      <button type="submit" disabled={loading}>
        {loading ? 'Adding...' : 'Add Expense'}
      </button>
    </form>
  );
};
```

## 🚨 Error Handling Best Practices

### Global Error Handler

```typescript
// src/utils/errorHandler.ts
import { ApiError } from '../services/apiClient';

export const handleApiError = (error: unknown): string => {
  if (error instanceof ApiError) {
    switch (error.status) {
      case 400:
        return `Invalid request: ${error.message}`;
      case 404:
        return 'Resource not found';
      case 500:
        return 'Server error. Please try again later.';
      default:
        return error.message;
    }
  }
  
  if (error instanceof Error) {
    return error.message;
  }
  
  return 'An unexpected error occurred';
};
```

### Error Boundary Component

```tsx
// src/components/ErrorBoundary.tsx
import React from 'react';

interface ErrorBoundaryState {
  hasError: boolean;
  error?: Error;
}

export class ErrorBoundary extends React.Component<
  React.PropsWithChildren<{}>,
  ErrorBoundaryState
> {
  constructor(props: React.PropsWithChildren<{}>) {
    super(props);
    this.state = { hasError: false };
  }

  static getDerivedStateFromError(error: Error): ErrorBoundaryState {
    return { hasError: true, error };
  }

  componentDidCatch(error: Error, errorInfo: React.ErrorInfo) {
    console.error('Error caught by boundary:', error, errorInfo);
  }

  render() {
    if (this.state.hasError) {
      return (
        <div className="error-boundary">
          <h2>Something went wrong.</h2>
          <p>{this.state.error?.message}</p>
          <button onClick={() => this.setState({ hasError: false })}>
            Try again
          </button>
        </div>
      );
    }

    return this.props.children;
  }
}
```

## 🧪 Testing Examples

### Service Tests

```typescript
// src/services/__tests__/walletService.test.ts
import { walletService } from '../walletService';
import { apiClient } from '../apiClient';

// Mock the API client
jest.mock('../apiClient');
const mockedApiClient = apiClient as jest.Mocked<typeof apiClient>;

describe('WalletService', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('should create a wallet', async () => {
    const mockResponse = { id: 'wallet-123', message: 'Wallet created' };
    mockedApiClient.post.mockResolvedValue(mockResponse);

    const walletData = {
      user_id: 'user-123',
      name: 'Test Wallet',
      type: 'checking',
      currency: 'USD',
    };

    const result = await walletService.createWallet(walletData);

    expect(mockedApiClient.post).toHaveBeenCalledWith('/api/v1/wallets', walletData);
    expect(result).toEqual(mockResponse);
  });

  it('should get user wallets', async () => {
    const mockResponse = {
      data: [{ id: 'wallet-123', name: 'Test Wallet' }],
      count: 1,
    };
    mockedApiClient.get.mockResolvedValue(mockResponse);

    const result = await walletService.getUserWallets('user-123');

    expect(mockedApiClient.get).toHaveBeenCalledWith('/api/v1/wallets?userID=user-123');
    expect(result).toEqual(mockResponse);
  });
});
```

### Component Tests

```tsx
// src/components/__tests__/CreateWalletForm.test.tsx
import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { CreateWalletForm } from '../CreateWalletForm';
import { walletService } from '../../services/walletService';

jest.mock('../../services/walletService');
const mockedWalletService = walletService as jest.Mocked<typeof walletService>;

describe('CreateWalletForm', () => {
  const mockOnSuccess = jest.fn();
  const defaultProps = {
    userID: 'user-123',
    onSuccess: mockOnSuccess,
  };

  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('should render form fields', () => {
    render(<CreateWalletForm {...defaultProps} />);
    
    expect(screen.getByLabelText(/wallet name/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/wallet type/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/currency/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/initial balance/i)).toBeInTheDocument();
  });

  it('should submit form with valid data', async () => {
    const mockResponse = { id: 'wallet-123', message: 'Success' };
    mockedWalletService.createWallet.mockResolvedValue(mockResponse);

    render(<CreateWalletForm {...defaultProps} />);
    
    fireEvent.change(screen.getByLabelText(/wallet name/i), {
      target: { value: 'My Checking' }
    });
    
    fireEvent.click(screen.getByRole('button', { name: /create wallet/i }));

    await waitFor(() => {
      expect(mockedWalletService.createWallet).toHaveBeenCalled();
      expect(mockOnSuccess).toHaveBeenCalledWith('wallet-123');
    });
  });
});
```

## 🏗️ Project Structure Recommendation

```
src/
├── components/          # React components
│   ├── common/         # Reusable components
│   ├── wallet/         # Wallet-specific components
│   ├── transaction/    # Transaction components
│   └── category/       # Category components
├── hooks/              # Custom React hooks
├── services/           # API services
├── types/              # TypeScript type definitions
├── utils/              # Utility functions
│   ├── currency.ts
│   ├── date.ts
│   └── errorHandler.ts
├── contexts/           # React contexts
└── __tests__/          # Test files
```

This guide provides a complete foundation for integrating with the accounting app backend API, with proper error handling, TypeScript support, and React best practices.