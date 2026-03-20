# Testing Guide

## Quick Commands

```bash
# Run all tests
task test

# Or directly with Go
go test -v ./...

# Run with coverage
task test:coverage

# Run specific test
go test -v -run TestGetAllClassifiers_CacheHit ./internal/handlers/

# Run all tests in a package
go test -v ./internal/handlers/
```

## Test Files

**`internal/cache/cache_test.go`** - 3 tests

- Cache miss returns nil (1)
- Set then get returns value (1)
- TTL expiry removes entry (1)

**`internal/handlers/classifiers_test.go`** - 7 tests

- GetAllClassifiers cache hit (1)
- GetAllClassifiers cache miss (1)
- GetAllClassifiers second call hits cache (1)
- GetAllClassifiers Redis fallback to DB (1)
- Auth context extraction (4)

**`internal/handlers/helpers_test.go`** - 43 tests

- HandlePgxError: no rows, PG error codes, unknown code,
  generic error (10)
- HandleGet: success, no auth, repo error (3)
- HandleGetByID: success, no auth, invalid ID,
  null result, nil result, repo error (8)
- HandleGetByString: success, no auth, null, error (4)
- HandlePost: success, no auth, invalid JSON, error (6)
- HandleDelete: success, no auth, invalid ID,
  not found, repo error (5)
- GetPathParamInt64: valid, non-integer (2)

## Key Testing Patterns

**Table-driven tests**: Multiple scenarios with `tests := []struct{...}`

```go
tests := []struct {
    name       string
    paramValue string
    wantCode   int
}{
    {"alphabetic", "abc", http.StatusBadRequest},
    {"float", "1.5", http.StatusBadRequest},
    {"special chars", "!@#", http.StatusBadRequest},
}

for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        // test logic
    })
}
```

**Gin test context**: HTTP handler tests use `httptest` with Gin

```go
w := httptest.NewRecorder()
c, _ := gin.CreateTestContext(w)
c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
```

**Mock repositories**: Interface-based mocks for database layer

```go
type mockRepo struct {
    result json.RawMessage
    err    error
}

func (m *mockRepo) GetAll(ctx context.Context) (json.RawMessage, error) {
    return m.result, m.err
}
```

## Contributing Tests

1. Follow naming: `Test<FunctionName>_<Scenario>`
2. Use table-driven tests for multiple scenarios
3. Use `httptest.NewRecorder()` and `gin.CreateTestContext()` for handler tests
4. Mock the repository interface, not the database
5. Verify: `task ci:all`
