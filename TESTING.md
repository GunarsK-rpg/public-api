# Testing Guide

## Quick Commands

```bash
task test              # Run all tests
task test:coverage     # Run with coverage
task ci:all            # Full CI suite (format, lint, vet, test)
```

## Test Environment

- **Framework**: Go `testing` package with subtests (`t.Run`)
- **HTTP**: `httptest` + Gin, shared `performRequest` / `withAuth` helpers
- **Redis**: `miniredis` for in-memory Redis (no real server needed)
- **Mocks**: Function-field mocks in `mock_repo_test.go`
- **CI**: GitHub Actions via `.github/workflows/ci.yml`

## Test Files

### Cache Tests

| File            | Coverage                                         |
| --------------- | ------------------------------------------------ |
| `cache_test.go` | Get/set, HasFlag/SetFlag, TTL expiry, cache miss |

### Handler Tests

| File                  | Coverage                                           |
| --------------------- | -------------------------------------------------- |
| `mock_repo_test.go`   | Shared function-field mock for all repositories    |
| `helpers_test.go`     | Auth context, pgx error mapping, handler helpers   |
| `classifiers_test.go` | Cache hit/miss, Redis fallback, scoped lookups     |
| `heroes_test.go`      | Hero CRUD and matrix coverage of all wrappers      |
| `campaigns_test.go`   | Campaign CRUD, join-by-code, remove hero           |
| `combat_test.go`      | NPCs, combats, instances, multi-param getters      |
| `avatar_test.go`      | Hero/NPC avatar set/delete with validation         |

### Middleware Tests

| File                | Coverage                                      |
| ------------------- | --------------------------------------------- |
| `bodylimit_test.go` | Body size under/at/over limit, nil body       |
| `usersync_test.go`  | Missing auth, cache hit, DB sync, failure     |

### Route Tests

| File             | Coverage                                     |
| ---------------- | -------------------------------------------- |
| `routes_test.go` | Hero sub-resource route registration         |

## What's Not Tested

- **Repository layer**: Calls PostgreSQL functions directly; tested via
  Flyway migration CI against a real database.
- **UserSync DB path**: The `pool.Exec` write path needs a real or
  mocked pgxpool; cache-hit and no-auth paths are tested.
- **Route middleware chain**: JWT and RBAC come from the common library
  and are tested there.
- **Config loading**: Verified implicitly by CI service startup.

## Contributing Tests

1. Add new hero wrapper handlers to the matching matrix table in
   `heroes_test.go` (GetByID, POST, or DELETE).
2. Add full individual coverage for custom/multi-param handlers:
   success, no auth, invalid params, not found, repo error.
3. Add the mock function field in `mock_repo_test.go`.
4. Run `task ci:all` before committing.
