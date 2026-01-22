# Scenario Results: expert

**Status:** SUCCESS
**Duration:** 44.518s

## Score

**Total:** 9/12 (Success)

| Dimension | Rating | Notes |
|-----------|--------|-------|
| Completeness | 0/3 | Resource count validation failed |
| Lint Quality | 3/3 | Deferred to domain tools |
| Output Validity | 3/3 | 7 files generated |
| Question Efficiency | 3/3 | 0 questions asked |

## Generated Files

- [README.md](README.md)
- [cmd/main.go](cmd/main.go)
- [go.mod](go.mod)
- [helpers.go](helpers.go)
- [jobs.go](jobs.go)
- [triggers.go](triggers.go)
- [workflows.go](workflows.go)

## Validation

**Status:** ❌ FAILED

### Resource Counts

| Domain | Type | Found | Constraint | Status |
|--------|------|-------|------------|--------|
| github | workflows | 0 | min: 1 | ❌ |

## Conversation

See [conversation.txt](conversation.txt) for the full prompt and response.
