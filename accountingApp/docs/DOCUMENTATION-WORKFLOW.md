# 📝 Documentation Update Workflow

> **Systematic approach to maintaining accurate project documentation**  
> Ensures CLAUDE.md and context files remain current and useful

## 🎯 Documentation Philosophy

### Core Principles
- **Backend Focus**: CLAUDE.md covers only Go backend implementation
- **Frontend Separation**: Frontend features documented separately
- **Accuracy First**: Documentation reflects actual implementation state
- **Developer Efficiency**: Quick context restoration for development sessions

### Documentation Hierarchy
```
CLAUDE.md                 # Main context file (backend only)
├── CONTEXT-TRACKER.md    # Development status tracking
├── PROJECT-STATUS.md     # Detailed implementation matrix
├── API-REFERENCE.md      # REST API documentation
└── DEVELOPER-GUIDE.md    # Architecture and workflow guide
```

## 🔄 Update Triggers

### Automatic Updates Required
- [ ] **New Domain Models**: Added entities, aggregates, or value objects
- [ ] **Use Case Changes**: New services or modified business logic
- [ ] **Repository Implementation**: New peer implementations or bridge changes
- [ ] **API Endpoint Changes**: New routes or modified request/response formats
- [ ] **Architecture Decisions**: Clean Architecture or pattern changes
- [ ] **Test Coverage Changes**: New test packages or significant coverage shifts

### Manual Update Indicators
- [ ] **Priority Task Completion**: When 🔴 Critical tasks are finished
- [ ] **Development Session Start**: Before beginning new feature work
- [ ] **Before Committing**: Major implementation milestones
- [ ] **Weekly Maintenance**: Regular accuracy verification

## 📋 Update Checklist

### CLAUDE.md Updates
```markdown
## Implementation Status Section
- [ ] Update Domain Layer status (✅/🟡/🔴)
- [ ] Update Application Layer completion percentages
- [ ] Update Adapter Layer implementation state
- [ ] Update Framework Layer configuration status

## Metrics Section  
- [ ] Count Go files: `find internal/accounting -name "*.go" | wc -l`
- [ ] Count test files: `find internal/accounting -name "*test*.go" | wc -l`
- [ ] Verify test status: `go test ./... -short`
- [ ] Update API endpoint count

## Priority Tasks Section
- [ ] Move completed tasks to "✅ Complete" 
- [ ] Add new critical/high priority items
- [ ] Update impact analysis and estimates
- [ ] Refresh next session goals

## Known Issues Section
- [ ] Remove resolved issues
- [ ] Add newly discovered bugs or limitations
- [ ] Update issue severity and impact assessment
```

### CONTEXT-TRACKER.md Updates
```markdown
## Implementation Matrix
- [ ] Update component status (✅/🟡/🔴)
- [ ] Refresh file counts and coverage percentages
- [ ] Update development priorities and impact analysis

## Test Coverage Section
- [ ] Run test suite and update timing metrics
- [ ] Verify all test packages still passing
- [ ] Update coverage percentages if significantly changed

## Development Metrics
- [ ] Auto-update file counts using commands
- [ ] Update performance benchmarks
- [ ] Refresh architecture compliance percentage
```

## 🛠️ Update Commands

### Backend Analysis Commands
```bash
# File counting
find internal/accounting -name "*.go" | wc -l          # Backend Go files
find internal/accounting -name "*test*.go" | wc -l     # Test files
find internal/accounting -type d | wc -l               # Directory structure

# Test verification
go test ./... -short                                   # Quick test run
go test ./... -cover                                   # With coverage
go test ./... -v | grep -E "(PASS|FAIL|SKIP)"        # Test status

# Code quality checks
go vet ./...                                           # Static analysis
go fmt ./... -l                                        # Format check
go mod tidy                                            # Dependency cleanup
```

### Implementation Status Verification
```bash
# Check for TODO comments (incomplete implementation)
grep -r "TODO" internal/accounting --include="*.go"

# Count interface implementations
grep -r "interface {" internal/accounting --include="*.go" | wc -l

# Find missing implementations
grep -r "// TODO:" internal/accounting --include="*.go" -A 2

# Check for error handling completeness
grep -r "if err != nil" internal/accounting --include="*.go" | wc -l
```

## 📅 Update Schedule

### Weekly Maintenance (Every Friday)
```markdown
## Weekly Documentation Review
- [ ] Run all update commands and refresh metrics
- [ ] Verify test suite still passes completely
- [ ] Review and update priority task lists
- [ ] Check for new architectural decisions or patterns
- [ ] Update last modified dates in all documentation files
```

### Session-Based Updates
```markdown
## Before Development Session
- [ ] Load current context: `/load @CLAUDE.md`
- [ ] Review CONTEXT-TRACKER.md priority tasks
- [ ] Verify development environment (database, tests)

## After Development Session  
- [ ] Update implementation status based on completed work
- [ ] Add new issues or technical debt discovered
- [ ] Refresh priority tasks for next session
- [ ] Commit documentation changes with implementation
```

### Milestone Updates
```markdown
## Major Feature Completion
- [ ] Comprehensive update of all documentation files
- [ ] Architectural decision recording if patterns changed
- [ ] API documentation updates if endpoints modified
- [ ] Test coverage analysis and reporting

## Architecture Changes
- [ ] Update Clean Architecture compliance status
- [ ] Document Bridge Pattern modifications
- [ ] Refresh dependency flow diagrams
- [ ] Update development workflow guides
```

## 🚀 Automation Opportunities

### Automated Metrics Collection
```bash
#!/bin/bash
# metrics-collector.sh - Auto-update project metrics

echo "## Auto-Generated Metrics ($(date))"
echo "Backend Go Files: $(find internal/accounting -name "*.go" | wc -l)"
echo "Test Files: $(find internal/accounting -name "*test*.go" | wc -l)"
echo "API Endpoints: $(grep -r "HandleFunc" internal/accounting --include="*.go" | wc -l)"
echo "Test Status: $(go test ./... -short | grep -c "ok")/6 packages passing"
```

### Git Hooks Integration
```bash
# .git/hooks/pre-commit
#!/bin/bash
# Ensure documentation stays current

# Check if backend files changed
if git diff --cached --name-only | grep -q "internal/accounting.*\.go$"; then
    echo "Backend changes detected. Consider updating CLAUDE.md"
    echo "Run: /document CLAUDE.md --update"
fi

# Verify tests still pass
if ! go test ./... -short; then
    echo "Tests failing. Update documentation after fixing."
    exit 1
fi
```

## 📊 Quality Metrics

### Documentation Accuracy Indicators
- **Stale Metrics**: File counts >10% off from actual
- **Outdated Status**: Implementation status doesn't match code
- **Missing Features**: New endpoints/models not documented
- **Broken References**: Links to non-existent files or sections

### Update Success Criteria
- [ ] All file counts match actual implementation
- [ ] All test packages show current pass/fail status
- [ ] Priority tasks reflect actual development needs
- [ ] Known issues list is current and actionable
- [ ] Developer can restore full context in <2 minutes

## 🎯 Best Practices

### Content Guidelines
- **Be Specific**: "🔴 Missing PostgreSQL peer implementation" vs "needs work"
- **Include Impact**: Why each priority task matters for the project
- **Use Consistent Status**: ✅ Complete, 🟡 Partial, 🔴 Missing
- **Quantify When Possible**: "85% coverage" vs "good coverage"

### Maintenance Habits
- **Update Immediately**: Don't batch documentation updates
- **Test Documentation**: Verify someone else could use the context
- **Keep Focused**: Backend documentation only covers backend
- **Prune Regularly**: Remove outdated information and stale tasks

---

**Workflow Version**: v1.0  
**Established**: 2025-08-16  
**Next Review**: Weekly Friday maintenance