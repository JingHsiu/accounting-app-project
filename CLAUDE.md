# Accounting App Project Standards & Context

This file provides guidance to Claude Code (claude.ai/code) when working with code in this multi-service personal accounting application.

> 📚 **Service-Specific Context**: See service directories for specialized development guidance
> 
> - **Backend (Go)**: [`accountingApp/CLAUDE.md`](accountingApp/CLAUDE.md) - Clean Architecture + DDD specifics
> - **Frontend (React)**: [`frontend/CLAUDE.md`](frontend/CLAUDE.md) - React/TypeScript patterns
> 
> 🏗️ **Complete Documentation**: See [`README.md`](README.md) and [`docs/`](docs/) for comprehensive system documentation

## Project Architecture Overview

This is a **multi-service personal accounting application** implementing **Clean Architecture** with **Domain-Driven Design (DDD)** principles.

### System Components

**Backend Service** (`accountingApp/`)
- **Language**: Go 1.19+
- **Architecture**: Clean Architecture (4-layer) with DDD
- **Patterns**: CQRS-lite, Repository with Bridge pattern, Specialized Controllers
- **Database**: PostgreSQL with ACID transactions
- **API**: RESTful HTTP API (port 8080)

**Frontend Service** (`frontend/`)
- **Framework**: React 18 + TypeScript
- **Build Tool**: Vite
- **Styling**: Tailwind CSS with glass-morphism design
- **State Management**: React Query
- **Development**: Hot reload (port 3000)

## Development Standards

### Git Workflow & Commit Conventions
```bash
# Branch naming convention
feature/wallet-management
bugfix/currency-conversion-issue
hotfix/critical-security-patch

# Commit message format
type(scope): brief description

feat(wallet): implement multi-currency support
fix(api): resolve 400 error in expense creation  
docs(readme): update development setup instructions
refactor(domain): extract Money value object patterns
```

### Code Quality Standards

**Validation Pipeline**
1. **Syntax Check**: Language parsers (Go compiler, TypeScript)
2. **Type Safety**: Go type system, TypeScript strict mode
3. **Linting**: Backend (go vet, golint), Frontend (ESLint)
4. **Testing**: Unit tests (100%), Integration tests (100%)
5. **Security**: Static analysis, dependency scanning
6. **Performance**: Load testing, memory profiling

**Quality Gates**
- All tests must pass before merge
- No critical security vulnerabilities
- Code coverage thresholds met
- Documentation updated for API changes

### API Contract Standards

**HTTP Response Format**
```json
{
  "success": true,
  "data": { ... },
  "error": "error message (when success: false)"
}
```

**Naming Conventions**
- **Backend**: snake_case for all API fields
- **Frontend**: camelCase for internal TypeScript types
- **Database**: snake_case for table/column names
- **Transformation**: Automatic conversion at API boundary

**Currency Handling**
- **Default Currency**: TWD (Taiwan Dollar) as primary base unit
- Store amounts as integers in base units (TWD=1, USD/EUR=100)
- Separate currency field (ISO 4217 3-character codes)
- Currency subdivision mapping: TWD/JPY/KRW/VND=1, USD/EUR/GBP/CNY=100

**Date/Time Format**
- API: ISO 8601 format (`2024-01-01T12:00:00Z`)
- Database: PostgreSQL TIMESTAMP WITH TIME ZONE
- Frontend: Local timezone conversion for display

## Cross-Service Integration Guidelines

### API Communication
```bash
# Base URLs
Backend API: http://localhost:8080/api/v1
Frontend Dev: http://localhost:3000

# Common Headers
Content-Type: application/json
Accept: application/json
```

### Data Transformation Strategy
**Shared Transformation Utilities** (`frontend/src/utils/apiTransform.ts`)
- Bidirectional snake_case ↔ camelCase conversion
- Special field mappings for irregular conversions
- Currency amount conversion (subdivision-aware)
- Date format standardization

### Error Handling Strategy
**Backend Error Responses**
```json
{
  "success": false,
  "error": "wallet_id is required",
  "details": { "field": "wallet_id", "code": "REQUIRED_FIELD" }
}
```

**Frontend Error Management**
- Enhanced React Query hooks with error classification
- User-friendly error messages with retry capabilities
- Error state components with recovery actions
- Automatic retry for server errors (500-level)

### Development Environment
```bash
# Environment Variables
DATABASE_URL="postgres://postgres:password@localhost:5432/accountingdb?sslmode=disable"
PORT="8080"
NODE_ENV="development"

# Docker Services
docker-compose up -d postgres    # PostgreSQL database
docker-compose up -d --profile admin  # With pgAdmin
```

## Security Standards & Compliance

### Data Protection
- **No sensitive data logging**: Never log passwords, tokens, or PII
- **Input validation**: All user inputs validated at API boundary
- **SQL injection prevention**: Parameterized queries only
- **XSS protection**: Proper input sanitization and output encoding

### Authentication & Authorization
- **Future implementation**: JWT-based authentication planned
- **Access control**: User-scoped data access (userID filtering)
- **Session security**: Secure session management patterns

### Dependency Management
```bash
# Backend dependency updates
go mod tidy && go mod download

# Frontend dependency updates  
npm audit && npm update

# Security scanning
go list -m all | nancy sleuth  # Backend
npm audit --audit-level moderate  # Frontend
```

## Performance Standards & Monitoring

### Response Time Targets
- **API Endpoints**: <200ms average response time
- **Database Queries**: <50ms for simple queries, <200ms for complex
- **Frontend Load**: <3 seconds initial load, <1 second navigation
- **Bundle Size**: <500KB initial JavaScript, <2MB total assets

### Scalability Considerations
- **Database**: Proper indexing, query optimization
- **API**: Stateless design, horizontal scaling ready
- **Frontend**: Code splitting, lazy loading, caching strategies

### Resource Management
```bash
# Database connection pooling
max_connections = 100
pool_size = 20

# Memory monitoring
go tool pprof -http=:8081 http://localhost:8080/debug/pprof/heap

# Frontend bundle analysis
npm run build && npm run analyze
```

## Context Management Rules

### Code Preservation During Compaction
When executing `/compact` or auto-compact operations:
- **ALWAYS preserve code changes**: Keep all recent code modifications, implementations, and fixes
- **Retain architectural decisions**: Maintain records of design choices and implementation patterns  
- **Preserve troubleshooting solutions**: Keep diagnostic information and resolution steps
- **Maintain development context**: Retain information about current development progress and active features
- **Cross-service integration**: Preserve API contract changes and integration patterns

This ensures continuity across sessions and prevents loss of important development work during context compression.

### Documentation Synchronization
- **API changes**: Update both backend and frontend documentation
- **Schema modifications**: Update all affected service documentation
- **Integration patterns**: Document cross-service communication updates
- **Performance changes**: Update benchmarks and monitoring guidelines

## Future Evolution Roadmap

### Planned Enhancements
- **Microservices Architecture**: Break backend into domain-specific services
- **Event Sourcing**: Implement event-driven architecture for audit trails
- **Advanced CQRS**: Separate read/write databases for optimal performance
- **Authentication Service**: Dedicated user management and authorization
- **Notification Service**: Real-time updates and alerting system
- **Reporting Service**: Advanced analytics and financial insights

### Scaling Considerations
- **Service Mesh**: Inter-service communication and observability
- **Container Orchestration**: Kubernetes deployment strategy
- **Database Sharding**: Data partitioning for large-scale deployments
- **CDN Integration**: Static asset optimization and global distribution

---

**Project Type**: Multi-Service Personal Accounting Application  
**Architecture**: Clean Architecture + DDD + React Frontend  
**Development Stage**: Core functionality complete, scaling preparation  
**Last Updated**: January 2025