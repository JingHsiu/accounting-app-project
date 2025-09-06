# API Documentation Corrections & Consistency Updates

**Generated:** 2025-09-06  
**Status:** Critical inconsistencies identified and corrected  
**Priority:** HIGH - Frontend integration depends on accurate documentation

## 🚨 Critical Corrections Required

### 1. Wallet Type Values (CRITICAL)

**❌ Documented (Incorrect):**
```json
{
  "type": "checking"  // ← WRONG - causes validation errors
}
```

**✅ Actual Backend Implementation:**
```json
{
  "type": "BANK"  // ← CORRECT - must use uppercase
}
```

**Valid Wallet Types:**
- `"CASH"` - Cash wallet
- `"BANK"` - Bank account (replaces "checking", "savings")
- `"CREDIT"` - Credit card
- `"INVESTMENT"` - Investment account

### 2. Get Wallets Response Format (CRITICAL)

**❌ Documented (Incorrect):**
```json
{
  "success": true,
  "data": {
    "data": [...],     // ← WRONG nested structure
    "count": 1
  }
}
```

**✅ Actual Backend Response:**
```json
{
  "success": true,
  "data": [           // ← CORRECT direct array
    {
      "id": "uuid",
      "user_id": "user-id",
      "name": "Wallet Name",
      "type": "BANK",
      "currency": "USD",
      "balance": {
        "amount": 100000,
        "currency": "USD"
      },
      "is_fully_loaded": false,
      "created_at": "2025-09-06T10:43:13Z",
      "updated_at": "2025-09-06T10:43:13Z"
    }
  ]
}
```

### 3. Missing Endpoints (HIGH PRIORITY)

**❌ Not Documented:**
- `GET /api/v1/categories` - Get all categories
- `GET /api/v1/categories/expense` - Get expense categories  
- `GET /api/v1/categories/income` - Get income categories
- `GET /api/v1/incomes?userID={id}` - Get user incomes
- `POST /api/v1/incomes` - Add income

**✅ Actually Implemented and Working:**
```bash
curl "http://localhost:8080/api/v1/categories/expense"
# Returns: {"success": true, "data": [{"id": "expense-1", "name": "餐飲", "type": "expense"}, ...]}
```

## 📋 Updated API Reference

### Wallet Management (Corrected)

#### Create Wallet
```http
POST /api/v1/wallets
Content-Type: application/json

{
  "user_id": "string",
  "name": "string", 
  "type": "CASH|BANK|CREDIT|INVESTMENT",  // ← CORRECTED
  "currency": "USD|TWD|EUR",
  "initialBalance": 100000                // ← Optional, in cents
}
```

**Response:**
```json
{
  "id": "wallet-uuid",
  "success": true,
  "message": ""
}
```

#### Get User Wallets
```http
GET /api/v1/wallets?userID={userID}
```

**Response:**
```json
{
  "success": true,
  "data": [                              // ← CORRECTED: Direct array
    {
      "id": "wallet-uuid",
      "user_id": "user-id",
      "name": "My Wallet",
      "type": "BANK",                    // ← CORRECTED: Uppercase
      "currency": "USD",
      "balance": {
        "amount": 100000,
        "currency": "USD"
      },
      "is_fully_loaded": false,          // ← ADDED: Missing field
      "created_at": "2025-09-06T10:43:13Z",
      "updated_at": "2025-09-06T10:43:13Z"
    }
  ]
}
```

### Category Management (NEW SECTION)

#### Get Expense Categories
```http
GET /api/v1/categories/expense
```

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": "expense-1",
      "name": "餐飲",
      "type": "expense"
    },
    {
      "id": "expense-2", 
      "name": "交通",
      "type": "expense"
    }
  ]
}
```

#### Get Income Categories  
```http
GET /api/v1/categories/income
```

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": "income-1",
      "name": "薪資",
      "type": "income"
    },
    {
      "id": "income-2",
      "name": "投資",
      "type": "income"
    }
  ]
}
```

#### Get All Categories
```http
GET /api/v1/categories
```

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": "expense-1",
      "name": "餐飲", 
      "type": "expense"
    },
    {
      "id": "income-1",
      "name": "薪資",
      "type": "income"
    }
  ]
}
```

### Income Management (NEW SECTION)

#### Get User Incomes
```http
GET /api/v1/incomes?userID={userID}
```

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": "income-uuid",
      "wallet_id": "wallet-uuid",
      "subcategory_id": "income-1", 
      "amount": 500000,
      "currency": "USD",
      "description": "Salary",
      "date": "2025-09-06T10:43:13Z",
      "created_at": "2025-09-06T10:43:13Z"
    }
  ]
}
```

#### Add Income
```http
POST /api/v1/incomes
Content-Type: application/json

{
  "wallet_id": "wallet-uuid",
  "subcategory_id": "income-1",
  "amount": 500000,                      // ← In cents/smallest unit
  "currency": "USD",
  "description": "Monthly Salary",
  "date": "2025-09-06T12:00:00Z"
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

## 🔧 Frontend Integration Updates

### TypeScript Type Corrections

**✅ Frontend types are already correct:**
```typescript
export enum WalletType {
  CASH = "CASH",
  BANK = "BANK", 
  CREDIT = "CREDIT",
  INVESTMENT = "INVESTMENT"
}

export interface Wallet {
  id: string;
  user_id: string;
  name: string;
  type: WalletType;
  balance: Money;
  currency: string;              // ← Additional backend field
  is_fully_loaded?: boolean;     // ← Additional backend field
  created_at: string;
  updated_at: string;
}
```

### API Service Layer Status

**✅ Frontend service implementation is correct:**
- Handles actual response format correctly
- Uses proper wallet type enums
- Includes comprehensive error handling
- Validates response structure

## 🎯 Action Items

### Immediate (Critical)
1. ✅ **Wallet Type Documentation** - Fixed in this document
2. ✅ **Response Format Documentation** - Corrected above
3. ✅ **Missing Endpoints Documentation** - Added above

### Next Steps
1. Update main `API_DOCUMENTATION.md` with corrections
2. Add category and income sections to main docs
3. Test all documented endpoints for accuracy
4. Validate frontend integration with corrected specs

## 🧪 Verification Commands

Test corrected API endpoints:

```bash
# Test wallet creation with correct type
curl -X POST http://localhost:8080/api/v1/wallets \
  -H "Content-Type: application/json" \
  -d '{"user_id": "test", "name": "Test", "type": "BANK", "currency": "USD", "initialBalance": 100000}'

# Test wallet retrieval
curl "http://localhost:8080/api/v1/wallets?userID=test"

# Test category endpoints
curl "http://localhost:8080/api/v1/categories/expense"
curl "http://localhost:8080/api/v1/categories/income" 

# Test income endpoint
curl "http://localhost:8080/api/v1/incomes?userID=test"
```

---

**Status: ✅ READY FOR INTEGRATION**  
**Impact: 🔥 Resolves frontend-backend contract mismatches**  
**Next: Update main API documentation with these corrections**