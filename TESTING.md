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
go test -v -run TestGetHeroes_WithCampaignFilter ./internal/handlers/

# Run all tests in a package
go test -v ./internal/handlers/
go test -v ./internal/middleware/
go test -v ./internal/routes/
```

## Test Structure

```text
internal/
  cache/
    cache_test.go          Redis get/set/flag operations with miniredis
  handlers/
    mock_repo_test.go      Shared mock implementing repository.Repository (64 methods)
    helpers_test.go        Auth context, error mapping, generic handler helpers
    classifiers_test.go    Classifier caching (cache hit/miss/fallback with miniredis)
    heroes_test.go         Hero CRUD, sub-resources, patches, equipment mods, favorites
    campaigns_test.go      Campaign CRUD, join-by-code, RemoveHeroFromCampaign
    combat_test.go         NPC templates, combat encounters, NPC instances, companions
    avatar_test.go         Hero avatar, NPC avatar (set/delete with JSON binding)
  middleware/
    bodylimit_test.go      Request body size limiting
    usersync_test.go       User sync middleware (auth rejection, cache-hit skip)
  routes/
    routes_test.go         registerHeroSubResource route registration
```

## Test Files

### `internal/cache/cache_test.go` -- 6 tests

- Cache miss returns nil
- Set then get round-trip
- HasFlag miss, SetFlag then HasFlag
- TTL expiry for Set and SetFlag

### `internal/handlers/mock_repo_test.go` -- shared mock

Implements all 64 methods across `ClassifierRepository`, `HeroRepository`,
`CampaignRepository`, and `CombatRepository`. Each method delegates to a
function field; returns "not implemented" if nil. Used by all handler test files.

### `internal/handlers/helpers_test.go` -- 43 tests

- `GetAuthContext`: success, missing user_id, missing username, wrong types
- `HandlePgxError`: no rows, 7 PG error codes (table-driven),
  unknown code, generic error
- `handleGet`: success, no auth, repo error
- `handleGetByID`: success, no auth, invalid ID (table-driven),
  null result, nil result, repo error
- `handleGetByString`: success, no auth, null result, repo error
- `handlePost`: success, no auth, invalid JSON (table-driven), repo error
- `handleDelete`: success, no auth, invalid ID, not found, repo error
- `getPathParamInt64`: valid, non-integer

### `internal/handlers/classifiers_test.go` -- 4 tests

- Cache hit returns cached data, skips DB
- Cache miss hits DB, stores in cache, verifies TTL
- Second call hits cache (only 1 DB call total)
- Redis failure falls back to DB

### `internal/handlers/heroes_test.go` -- 27 tests

- `GetHeroes`: success (no filter), campaign_id filter,
  invalid campaign_id, no auth, repo error
- `GetHero`: success, not found (pgx.ErrNoRows), no auth, invalid ID
- `GetHeroSheet`: success with ID verification
- `CreateHero`: success (body forwarded), no auth
- `UpdateHero`: success
- `DeleteHero`: success, not found
- Sub-resources (attributes): get success, upsert success, delete success
- Resource patches (health): success, no auth
- Equipment modifications: add success, remove success
- Favorite actions: add success, remove success

### `internal/handlers/campaigns_test.go` -- 15 tests

- `GetCampaigns`: success, no auth, repo error
- `GetCampaign`: success (ID verification), not found
- `GetCampaignByCode`: success (code verification), not found (null result)
- `CreateCampaign`: success
- `DeleteCampaign`: success (ID verification)
- `RemoveHeroFromCampaign`: success (both params verified), no auth,
  invalid hero ID, invalid campaign ID, not found, repo error

### `internal/handlers/combat_test.go` -- 25 tests

- `GetNpcOptions`: success (campaign ID verified)
- `GetNpcLibrary`: success (campaign ID verified)
- `GetNpc` (multi-param): success, no auth, invalid NPC ID,
  invalid campaign ID, null result, repo error
- `GetNpcByID`: success
- `CreateNpc`: success
- `DeleteNpc` (multi-param): success, no auth, invalid NPC ID, not found
- `GetCombat` (multi-param): success, no auth, invalid combat ID, null result
- `DeleteCombat` (multi-param): success, not found
- `GetCombats`: success
- `EndCombatRound`: success
- NPC instances: get, create, patch, delete -- all success
- Companions: GetHeroCompanions, GetCompanionNpcOptions -- success

### `internal/handlers/avatar_test.go` -- 13 tests

- `SetHeroAvatar`: success (ID + avatarKey verified), no auth,
  invalid ID, missing avatarKey, repo error
- `DeleteHeroAvatar`: success, no auth, repo error
- `SetNpcAvatar`: success (npcID + campaignID + avatarKey verified),
  invalid campaign ID, invalid NPC ID
- `DeleteNpcAvatar`: success, no auth

### `internal/middleware/bodylimit_test.go` -- 4 tests

- Body under limit passes through
- Body over limit returns 413
- Body at exact limit passes through
- Nil body (GET request) passes through

### `internal/middleware/usersync_test.go` -- 2 tests

- Missing auth context returns 401
- Cache hit skips DB call, handler proceeds

### `internal/routes/routes_test.go` -- 2 tests

- `registerHeroSubResource` registers GET/POST/DELETE with correct paths
- Multiple resources don't conflict

## Testing Patterns

### Mock Repository

All handler tests share `mockRepo` from `mock_repo_test.go`. Set the function
field for the method your test needs; unset methods return "not implemented".

```go
mock := &mockRepo{}
mock.getHeroFunc = func(
    _ context.Context, _ repository.AuthContext, id int64,
) (json.RawMessage, error) {
    return json.RawMessage(`{"id":1}`), nil
}
handler := New(mock, nil)
```

### HTTP Handler Tests

Use `performRequest` helper from `helpers_test.go` with `withAuth` for auth context:

```go
router := gin.New()
router.GET("/heroes/:id", func(c *gin.Context) {
    withAuth(c)
    handler.GetHero(c)
})
w := performRequest(t, router, "GET", "/heroes/42", nil)
```

### Redis Tests

Use `miniredis` for in-memory Redis (no real server needed):

```go
mr := miniredis.RunT(t)
client := redis.NewClient(&redis.Options{Addr: mr.Addr()})
appCache := cache.New(client)
```

### Table-Driven Tests

For multiple scenarios on the same handler:

```go
tests := []struct {
    name   string
    id     string
    status int
}{
    {"alphabetic", "abc", http.StatusBadRequest},
    {"float", "1.5", http.StatusBadRequest},
}
for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        w := performRequest(t, router, "GET", "/items/"+tt.id, nil)
        if w.Code != tt.status {
            t.Errorf("status = %d, want %d", w.Code, tt.status)
        }
    })
}
```

## What's Not Tested

- **Repository layer**: Calls PostgreSQL functions directly. Tested via Flyway
  migration CI (real PostgreSQL in GitHub Actions).
- **UserSync DB path**: The `pool.Exec` call in usersync middleware requires
  a real or mocked pgxpool. Cache-hit and no-auth paths are tested.
- **Route middleware chain**: JWT validation and RBAC permission middleware
  come from the common library and are tested there.
- **Config loading**: Reads environment variables with validation. Tested
  implicitly via CI (service starts with valid config).

## Adding New Tests

1. If adding a new handler that delegates to `handleGetByID`/`handlePost`/`handleDelete`,
   add a success test verifying correct param name and repo method wiring.
2. If adding a custom handler (multi-param, special logic), add full coverage:
   success, no auth, invalid params, not found, repo error.
3. Add the mock function field to `mockRepo` in `mock_repo_test.go` and
   implement the delegation method.
4. Run `task ci:all` before committing.
