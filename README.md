# Accounting App Documentation

Complete documentation for the Accounting App project - a modern financial management system built with Go backend and React frontend.

## 📚 Documentation Overview

This directory contains comprehensive documentation for developers, architects, and stakeholders working with the Accounting App.

### 📁 Documentation Structure

```
docs/
├── README.md                           # This file - documentation index
├── SYSTEM-ARCHITECTURE.md             # Complete system architecture overview
├── FRONTEND_INTEGRATION_GUIDE.md      # Frontend integration guide
├── DEBUG_INSTRUCTIONS.md              # Debug and troubleshooting guide
└── api/
    ├── API_DOCUMENTATION.md           # Complete backend API reference
    └── QUICK_REFERENCE.md             # Fast API lookup table
```

---

## 🚀 Quick Start

### For New Developers
- **Architecture Overview**: Start with [`SYSTEM-ARCHITECTURE.md`](SYSTEM-ARCHITECTURE.md) for comprehensive system understanding
- **API Reference**: See [`api/API_DOCUMENTATION.md`](api/API_DOCUMENTATION.md) for complete endpoint documentation
- **Frontend Integration**: See [`FRONTEND_INTEGRATION_GUIDE.md`](FRONTEND_INTEGRATION_GUIDE.md) for React/TypeScript setup

### For Backend Developers
- **API Reference**: See [`api/API_DOCUMENTATION.md`](api/API_DOCUMENTATION.md) for complete endpoint documentation
- **Quick Reference**: See [`api/QUICK_REFERENCE.md`](api/QUICK_REFERENCE.md) for fast API lookups
- **Architecture**: See [`SYSTEM-ARCHITECTURE.md`](SYSTEM-ARCHITECTURE.md) for Clean Architecture implementation

### For Frontend Developers
- **Integration Guide**: See [`FRONTEND_INTEGRATION_GUIDE.md`](FRONTEND_INTEGRATION_GUIDE.md) for React/TypeScript implementation
- **API Reference**: See [`api/API_DOCUMENTATION.md`](api/API_DOCUMENTATION.md) for endpoint details
- **Debug Guide**: See [`DEBUG_INSTRUCTIONS.md`](DEBUG_INSTRUCTIONS.md) for troubleshooting

---

## 📖 Documentation Contents

### 🏗️ [System Architecture](SYSTEM-ARCHITECTURE.md)
Complete architectural overview of the multi-service accounting application.

**Contents:**
- **Clean Architecture Implementation** - 4-layer architecture with DDD principles
- **Service Architecture** - Multi-service structure and communication patterns  
- **Backend Architecture** - Go implementation with Bridge pattern
- **Frontend Architecture** - React with TypeScript and modern patterns
- **Integration Architecture** - API contracts and data flow
- **Development Workflow** - Setup, testing, and deployment strategies
- **Future Evolution** - Roadmap for CQRS, microservices, and scaling

**Target Audience:** System architects, senior developers, technical leads

---

### 🔧 [API Documentation](api/API_DOCUMENTATION.md)
Complete backend API reference for the accounting application.

**Contents:**
- **Wallet Management APIs** - Create, read, update, delete wallets
- **Transaction APIs** - Add expenses and income transactions
- **Category Management APIs** - Manage expense and income categories
- **Response Formats** - Standard success/error response patterns
- **Error Handling** - HTTP status codes and error messages
- **Currency Guidelines** - How to handle monetary amounts correctly

**Target Audience:** Backend developers, API consumers, frontend developers

---

### 🌐 [Frontend Integration Guide](FRONTEND_INTEGRATION_GUIDE.md)
Production-ready TypeScript implementation guide for frontend developers.

**Contents:**
- **API Client Setup** - Centralized HTTP client with error handling
- **TypeScript Definitions** - Complete type definitions for API data
- **Service Implementations** - Ready-to-use service classes for all APIs
- **React Components** - Example components with best practices
- **Utility Functions** - Currency, date, and error handling utilities
- **Custom Hooks** - React hooks for data fetching and state management
- **Testing Examples** - Unit tests for services and components

**Target Audience:** Frontend developers, React/TypeScript developers

---

### 🐛 [Debug Instructions](DEBUG_INSTRUCTIONS.md)
Troubleshooting guide for development and debugging.

**Contents:**
- **Frontend Debugging** - React Developer Tools, network debugging
- **Backend Debugging** - Go debugging, database troubleshooting
- **Common Issues** - Known problems and solutions
- **Testing Tools** - Debug components and utilities

**Target Audience:** All developers, QA engineers

---

## 🏗️ Architecture Overview

The Accounting App follows **Clean Architecture** principles with **Domain-Driven Design (DDD)**:

### 4-Layer Architecture
1. **Domain Layer** - Core business logic (Wallet, Money, Categories)
2. **Application Layer** - Use cases and application services
3. **Adapter Layer** - Controllers and repository implementations
4. **Infrastructure Layer** - Database connections, web framework

### Key Patterns
- **CQRS-lite** - Separate command and query operations
- **Repository Pattern** - Data access abstraction
- **Specialized Controllers** - Single-responsibility controllers
- **Bridge Pattern** - Clean layer separation

---

## 🛠️ Development Setup

### Prerequisites
- **Backend**: Go 1.19+, PostgreSQL, Docker (optional)
- **Frontend**: Node.js 16+, npm/yarn, TypeScript

### Backend Setup
```bash
# Navigate to backend directory
cd accountingApp

# Install dependencies
go mod tidy

# Start database (with Docker)
docker-compose up -d postgres

# Run the application
go run cmd/accountingApp/main.go

# Run tests
go test ./...
```

### Frontend Setup
```bash
# Navigate to frontend directory
cd frontend

# Install dependencies
npm install

# Start development server
npm run dev

# Run tests
npm run test
```

---

## 📝 API Quick Reference

### Base URL
```
http://localhost:8080/api/v1
```

### Common Endpoints
| Method | Endpoint | Description | Status |
|--------|----------|-------------|---------|
| `GET` | `/health` | Health check | ✅ Working |
| `POST` | `/wallets` | Create wallet | ✅ Working |
| `GET` | `/wallets?userID={id}` | Get user's wallets | ✅ Working |
| `GET` | `/wallets/{id}` | Get single wallet | ✅ Working |
| `PUT` | `/wallets/{id}` | Update wallet | 🚧 Planned |
| `DELETE` | `/wallets/{id}` | Delete wallet | 🚧 Planned |
| `GET` | `/wallets/{id}/balance` | Get wallet balance | ✅ Working |
| `POST` | `/expenses` | Add expense | ✅ Working |
| `POST` | `/incomes` | Add income | ✅ Working |
| `GET` | `/incomes?userID={id}` | Get income records | ✅ Working |
| `GET` | `/categories` | Get all categories | ✅ Working |
| `GET` | `/categories/expense` | Get expense categories | ✅ Working |
| `GET` | `/categories/income` | Get income categories | ✅ Working |
| `POST` | `/categories/expense` | Create expense category | ✅ Working |
| `POST` | `/categories/income` | Create income category | ✅ Working |

### Response Format
```json
{
  "success": true,
  "data": { ... }
}
```

---

## 💡 Usage Examples

### Creating a Wallet
```bash
curl -X POST http://localhost:8080/api/v1/wallets \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user-123",
    "name": "My Bank Account",
    "type": "BANK",
    "currency": "USD",
    "initialBalance": 100000
  }'
```

### Adding an Expense
```bash
curl -X POST http://localhost:8080/api/v1/expenses \
  -H "Content-Type: application/json" \
  -d '{
    "wallet_id": "wallet-123",
    "subcategory_id": "food-123",
    "amount": 1250,
    "currency": "USD",
    "description": "Lunch",
    "date": "2024-01-01T12:00:00Z"
  }'
```

---

## 🔍 Key Implementation Details

### Currency Handling
- All amounts stored as integers in smallest currency unit
- USD/EUR: stored in cents (divide by 100 for display)
- JPY: stored as yen (no conversion needed)

### Date Format
- API accepts/returns ISO 8601 format: `2024-01-01T12:00:00Z`
- Frontend should convert to/from user's timezone

### Error Handling
- Standard HTTP status codes (200, 400, 404, 500)
- Consistent error response format
- Detailed error messages for debugging

---

## 🤝 Contributing

### Documentation Updates
1. Update relevant documentation files
2. Keep examples current with API changes
3. Test all code examples before committing
4. Update this README if adding new documentation

### Code Changes
1. Update API documentation for any endpoint changes
2. Update TypeScript types in frontend guide
3. Add new examples for new features
4. Test all integration examples

---

## 📞 Support

For questions about:
- **API Usage**: See [API Documentation](api/API_DOCUMENTATION.md)
- **Frontend Integration**: See [Frontend Integration Guide](FRONTEND_INTEGRATION_GUIDE.md)
- **System Architecture**: See [System Architecture](SYSTEM-ARCHITECTURE.md)
- **Debugging**: See [Debug Instructions](DEBUG_INSTRUCTIONS.md)
- **Development Setup**: See this README

---

**Last Updated:** January 2025  
**Version:** v1.0  
**API Version:** v1