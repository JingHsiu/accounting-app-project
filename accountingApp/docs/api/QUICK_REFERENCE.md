# API Quick Reference

Fast reference for the Accounting App API endpoints.

**Base URL:** `http://localhost:8080/api/v1`

## 🏦 Wallet Endpoints

| Method | Endpoint | Description | Required Fields |
|--------|----------|-------------|-----------------|
| `POST` | `/wallets` | Create wallet | `user_id`, `name`, `type`, `currency` |
| `GET` | `/wallets?userID={id}` | Get user's wallets | `userID` (query) |
| `GET` | `/wallets/{id}` | Get single wallet | `id` (path) |
| `PUT` | `/wallets/{id}` | Update wallet | `id` (path) |
| `DELETE` | `/wallets/{id}` | Delete wallet | `id` (path) |
| `GET` | `/wallets/{id}/balance` | Get wallet balance | `id` (path) |

## 💸 Transaction Endpoints

| Method | Endpoint | Description | Required Fields |
|--------|----------|-------------|-----------------|
| `POST` | `/expenses` | Add expense | `wallet_id`, `subcategory_id`, `amount`, `currency`, `date` |
| `POST` | `/incomes` | Add income | `wallet_id`, `subcategory_id`, `amount`, `currency`, `date` |

## 🏷️ Category Endpoints

| Method | Endpoint | Description | Required Fields |
|--------|----------|-------------|-----------------|
| `POST` | `/categories/expense` | Create expense category | `user_id`, `name` |
| `POST` | `/categories/income` | Create income category | `user_id`, `name` |

## 🔧 Utility Endpoints

| Method | Endpoint | Description | Required Fields |
|--------|----------|-------------|-----------------|
| `GET` | `/health` | Health check | None |

## 📄 Response Formats

### Success Response
```json
{
  "success": true,
  "data": { ... }
}
```

### Error Response
```json
{
  "success": false,
  "error": "Error message"
}
```

## 🚨 Common Status Codes

- `200` - Success
- `400` - Bad Request (validation errors)
- `404` - Not Found
- `405` - Method Not Allowed
- `500` - Internal Server Error

## 💰 Currency Notes

- Amounts are integers in smallest unit (cents for USD/EUR)
- USD $12.34 → `1234`
- EUR €12.34 → `1234`
- JPY ¥1234 → `1234`

## 📅 Date Format

Use ISO 8601 format: `2024-01-01T12:00:00Z`

## 🔗 Full Documentation

See [API_DOCUMENTATION.md](API_DOCUMENTATION.md) for complete details, examples, and integration guides.