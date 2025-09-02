# Accounting App Backend API Documentation

**Version:** v1  
**Base URL:** `http://localhost:8080/api/v1`  
**Content-Type:** `application/json`

## Overview

This API follows RESTful principles and Clean Architecture patterns. All responses include a `success` boolean and either `data` or `error` field.

### Standard Response Format

**Success Response:**
```json
{
  "success": true,
  "data": { ... }
}
```

**Error Response:**
```json
{
  "success": false,
  "error": "Error message"
}
```

### Common HTTP Status Codes

- `200` - Success
- `400` - Bad Request (validation errors)
- `404` - Not Found
- `405` - Method Not Allowed
- `500` - Internal Server Error

---

## 🏦 Wallet Management APIs

### Create Wallet
Create a new wallet for a user.

**Endpoint:** `POST /api/v1/wallets`

**Request Body:**
```json
{
  "user_id": "string",          // Required: User ID
  "name": "string",             // Required: Wallet name
  "type": "string",             // Required: Wallet type (checking, savings, credit, etc.)
  "currency": "string",         // Required: Currency code (USD, TWD, etc.)
  "initialBalance": 0           // Optional: Initial balance in smallest currency unit (cents)
}
```

**Response:**
```json
{
  "id": "wallet-uuid",
  "success": true,
  "message": "Wallet created successfully"
}
```

**Frontend Example:**
```typescript
const createWallet = async (walletData: {
  user_id: string;
  name: string;
  type: string;
  currency: string;
  initialBalance?: number;
}) => {
  const response = await fetch('/api/v1/wallets', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(walletData)
  });
  return await response.json();
};
```

---

### Get User's Wallets
Retrieve all wallets for a specific user.

**Endpoint:** `GET /api/v1/wallets?userID={userID}`

**Query Parameters:**
- `userID` (required): User ID to fetch wallets for

**Response:**
```json
{
  "success": true,
  "data": {
    "data": [
      {
        "id": "wallet-uuid",
        "user_id": "user-uuid",
        "name": "My Checking Account",
        "type": "checking",
        "currency": "USD",
        "balance": {
          "amount": 150000,        // Amount in cents (1500.00)
          "currency": "USD"
        },
        "created_at": "2024-01-01T12:00:00Z",
        "updated_at": "2024-01-01T12:00:00Z",
        "is_fully_loaded": false
      }
    ],
    "count": 1
  }
}
```

**Frontend Example:**
```typescript
const getUserWallets = async (userID: string) => {
  const response = await fetch(`/api/v1/wallets?userID=${encodeURIComponent(userID)}`);
  return await response.json();
};
```

---

### Get Single Wallet
Retrieve detailed information about a specific wallet.

**Endpoint:** `GET /api/v1/wallets/{walletID}`

**Path Parameters:**
- `walletID`: Wallet ID to retrieve

**Query Parameters:**
- `includeTransactions` (optional): Set to "true" to include transaction history

**Response:**
```json
{
  "success": true,
  "data": {
    "data": {
      "id": "wallet-uuid",
      "user_id": "user-uuid",
      "name": "My Checking Account",
      "type": "checking",
      "currency": "USD",
      "balance": {
        "amount": 150000,
        "currency": "USD"
      },
      "created_at": "2024-01-01T12:00:00Z",
      "updated_at": "2024-01-01T12:00:00Z",
      "is_fully_loaded": true,
      "transactions": []           // Only if includeTransactions=true
    }
  }
}
```

**Frontend Example:**
```typescript
const getWallet = async (walletID: string, includeTransactions = false) => {
  const url = `/api/v1/wallets/${encodeURIComponent(walletID)}${includeTransactions ? '?includeTransactions=true' : ''}`;
  const response = await fetch(url);
  return await response.json();
};
```

---

### Update Wallet
Update wallet properties (partial updates supported).

**Endpoint:** `PUT /api/v1/wallets/{walletID}`

**Path Parameters:**
- `walletID`: Wallet ID to update

**Request Body:** (all fields optional)
```json
{
  "name": "Updated Wallet Name",     // Optional: New wallet name
  "type": "savings",                 // Optional: New wallet type
  "currency": "EUR"                  // Optional: New currency
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "message": "Wallet updated successfully"
  }
}
```

**Frontend Example:**
```typescript
const updateWallet = async (walletID: string, updates: {
  name?: string;
  type?: string;
  currency?: string;
}) => {
  const response = await fetch(`/api/v1/wallets/${encodeURIComponent(walletID)}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(updates)
  });
  return await response.json();
};
```

---

### Delete Wallet
Delete a wallet permanently.

**Endpoint:** `DELETE /api/v1/wallets/{walletID}`

**Path Parameters:**
- `walletID`: Wallet ID to delete

**Response:**
```json
{
  "success": true,
  "data": {
    "message": "Wallet deleted successfully"
  }
}
```

**Frontend Example:**
```typescript
const deleteWallet = async (walletID: string) => {
  const response = await fetch(`/api/v1/wallets/${encodeURIComponent(walletID)}`, {
    method: 'DELETE'
  });
  return await response.json();
};
```

---

### Get Wallet Balance
Get the current balance of a specific wallet.

**Endpoint:** `GET /api/v1/wallets/{walletID}/balance`

**Path Parameters:**
- `walletID`: Wallet ID to get balance for

**Response:**
```json
{
  "walletId": "wallet-uuid",
  "balance": 150000,              // Amount in smallest currency unit
  "currency": "USD",
  "success": true,
  "message": "Balance retrieved successfully"
}
```

**Frontend Example:**
```typescript
const getWalletBalance = async (walletID: string) => {
  const response = await fetch(`/api/v1/wallets/${encodeURIComponent(walletID)}/balance`);
  return await response.json();
};
```

---

## 💸 Transaction Management APIs

### Add Expense
Record a new expense transaction.

**Endpoint:** `POST /api/v1/expenses`

**Request Body:**
```json
{
  "wallet_id": "string",        // Required: Target wallet ID
  "subcategory_id": "string",   // Required: Expense subcategory ID
  "amount": 5000,               // Required: Amount in smallest currency unit (positive)
  "currency": "USD",            // Required: Currency code
  "description": "Coffee",      // Optional: Transaction description
  "date": "2024-01-01T12:00:00Z" // Required: Transaction date
}
```

**Response:**
```json
{
  "id": "expense-uuid",
  "success": true,
  "message": "Expense added successfully"
}
```

**Frontend Example:**
```typescript
const addExpense = async (expenseData: {
  wallet_id: string;
  subcategory_id: string;
  amount: number;
  currency: string;
  description?: string;
  date: string; // ISO format
}) => {
  const response = await fetch('/api/v1/expenses', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(expenseData)
  });
  return await response.json();
};
```

---

### Add Income
Record a new income transaction.

**Endpoint:** `POST /api/v1/incomes`

**Request Body:**
```json
{
  "wallet_id": "string",        // Required: Target wallet ID
  "subcategory_id": "string",   // Required: Income subcategory ID
  "amount": 500000,             // Required: Amount in smallest currency unit (positive)
  "currency": "USD",            // Required: Currency code
  "description": "Salary",     // Optional: Transaction description
  "date": "2024-01-01T12:00:00Z" // Required: Transaction date
}
```

**Response:**
```json
{
  "id": "income-uuid",
  "success": true,
  "message": "Income added successfully"
}
```

**Frontend Example:**
```typescript
const addIncome = async (incomeData: {
  wallet_id: string;
  subcategory_id: string;
  amount: number;
  currency: string;
  description?: string;
  date: string; // ISO format
}) => {
  const response = await fetch('/api/v1/incomes', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(incomeData)
  });
  return await response.json();
};
```

---

## 🏷️ Category Management APIs

### Create Expense Category
Create a new expense category.

**Endpoint:** `POST /api/v1/categories/expense`

**Request Body:**
```json
{
  "user_id": "string",          // Required: User ID
  "name": "string"              // Required: Category name
}
```

**Response:**
```json
{
  "id": "category-uuid",
  "success": true,
  "message": "Expense category created successfully"
}
```

**Frontend Example:**
```typescript
const createExpenseCategory = async (categoryData: {
  user_id: string;
  name: string;
}) => {
  const response = await fetch('/api/v1/categories/expense', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(categoryData)
  });
  return await response.json();
};
```

---

### Create Income Category
Create a new income category.

**Endpoint:** `POST /api/v1/categories/income`

**Request Body:**
```json
{
  "user_id": "string",          // Required: User ID
  "name": "string"              // Required: Category name
}
```

**Response:**
```json
{
  "id": "category-uuid",
  "success": true,
  "message": "Income category created successfully"
}
```

**Frontend Example:**
```typescript
const createIncomeCategory = async (categoryData: {
  user_id: string;
  name: string;
}) => {
  const response = await fetch('/api/v1/categories/income', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(categoryData)
  });
  return await response.json();
};
```

---

## 🔧 Utility APIs

### Health Check
Check API health status.

**Endpoint:** `GET /health`

**Response:**
```json
{
  "status": "healthy"
}
```

**Frontend Example:**
```typescript
const checkHealth = async () => {
  const response = await fetch('/health');
  return await response.json();
};
```

---

## 🚧 Not Yet Implemented

The following endpoints are planned but not yet implemented:

### Category Collection Operations
- `GET /api/v1/categories` - Get all categories with type filter
- `POST /api/v1/categories` - Generic create category
- `GET /api/v1/categories/{categoryID}` - Get category by ID
- `PUT /api/v1/categories/{categoryID}` - Update category
- `DELETE /api/v1/categories/{categoryID}` - Delete category
- `GET /api/v1/categories/tree` - Get category tree structure

---

## 🛠️ Frontend Integration Tips

### Currency Handling
All monetary amounts are stored and transmitted as integers in the smallest currency unit:
- USD: cents (divide by 100 for display)
- EUR: cents (divide by 100 for display)
- JPY: yen (no division needed)

```typescript
// Convert from display to API format
const displayToAPI = (amount: number, currency: string): number => {
  const centsBasedCurrencies = ['USD', 'EUR', 'TWD'];
  return centsBasedCurrencies.includes(currency) ? Math.round(amount * 100) : amount;
};

// Convert from API to display format
const apiToDisplay = (amount: number, currency: string): number => {
  const centsBasedCurrencies = ['USD', 'EUR', 'TWD'];
  return centsBasedCurrencies.includes(currency) ? amount / 100 : amount;
};
```

### Error Handling
Always check the `success` field in responses:

```typescript
const handleApiResponse = async (response: Response) => {
  const data = await response.json();
  if (!data.success) {
    throw new Error(data.error || 'Unknown API error');
  }
  return data.data || data;
};
```

### Date Format
Use ISO 8601 format for dates:

```typescript
const formatDateForAPI = (date: Date): string => {
  return date.toISOString();
};
```

### URL Encoding
Always encode path parameters to handle special characters:

```typescript
const buildApiUrl = (path: string, ...params: string[]): string => {
  const encodedParams = params.map(param => encodeURIComponent(param));
  return `/api/v1/${path}/${encodedParams.join('/')}`;
};
```

---

## 📚 Architecture Notes

This API follows Clean Architecture principles with:

- **Domain Layer**: Business logic and entities (Wallet, Money, Categories)
- **Application Layer**: Use cases and application services
- **Adapter Layer**: Controllers that handle HTTP requests/responses
- **Infrastructure Layer**: Database connections and external services

Each controller is specialized for single operations following the Single Responsibility Principle.