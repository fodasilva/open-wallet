# E2E Testing Directives

This document outlines the patterns, core utility functions, directives, and best practices for writing End-to-End (E2E) tests in the `open-wallet` platform.

## Test Setup & Utilities
To ensure true application testing capability, this project uses `testcontainers-go` to spin up live instances of PostgreSQL and Redis. Every E2E test suite has its dependencies thoroughly isolated to avoid flakiness and collisions.

### Core Helpers (from `helpers.go`)
1. **`SetupTestResources(t *testing.T)`**: Spins up Postgres and Redis containers, runs standard database migrations, and initializes singleton clients. This function ensures your environment is fresh but authentic.
2. **`SetupTestUser(t *testing.T, db *sql.DB, cfg *infra.Config)`**: Seeds a mock user into the DB and dynamically mints a valid JWT validation token required by the API. Returns a `TestUser` struct and `Token` string.
3. **`AssertTableIsEmpty(t *testing.T, db *sql.DB, tableName string)`**: A fast teardown utility checking equality counts per test table ensuring there's zero contamination across module tests.
4. **`AssertUnauthorized(t *testing.T, mux *http.ServeMux, method string, url string, body io.Reader)`**: Automatically probes an endpoint and verifies that it returns a rigid `401 Unauthorized` without a valid JWT token.

## Architectural Testing Patterns

### 1. Authentication Enforcement (Initial Route Test)
Every RESTful suite should begin with an enforcement check to catch any unprotected or erroneously public routes traversing via Gin framework definitions.
Usually achieved by looping through endpoints under an umbrella test.

```go
t.Run("Authentication Enforcement", func(t *testing.T) {
    endpoints := []struct {
        method string
        url    string
    }{
        {http.MethodPost, "/api/v1/resource"},
        {http.MethodGet, "/api/v1/resource"},
        {http.MethodDelete, "/api/v1/resource/123"},
    }

    for _, e := range endpoints {
        AssertUnauthorized(t, router, e.method, e.url, nil)
    }
})
```

### 2. Table-Driven Tests
Because payloads mutate, responses fluctuate, and side-effects matter, we adopt **table-driven tests**. This structures execution cleanly, maps intent, limits repetition and forces clear declaration over imperative messy code blocks.

**Creating a Table:**
We structure a localized array of test structs representing different conditions/inputs:
```go
type testCase struct {
    name            string
    payload         handlers.CreatePayloadRequest // Use actual requests Structs
    expectedStatus  int
    validateDB      bool
    expectEndPeriod *string
}
```

**Executing The Sequence:**
After the slice arrays are populated you simply map `t.Run()`:
```go
for _, tc := range cases {
    t.Run(tc.name, func(t *testing.T) {
        body, _ := json.Marshal(tc.payload)
        req := httptest.NewRequest(http.MethodPost, "/api/v1/resource", bytes.NewBuffer(body))
        req.Header.Set("Authorization", "Bearer "+token)
        req.Header.Set("Content-Type", "application/json")
        w := httptest.NewRecorder()

        router.ServeHTTP(w, req)

        assert.Equal(t, tc.expectedStatus, w.Code)
		
        // Add Specific Database Side-Effect Assertions if Needed
        if tc.validateDB { ... }
    })
}
```

### 3. Setup, Clearing, and Tear-Down Directives
1. **Never mock handlers**: Start your router using real instantiated `routes.SetupXRoutes()`.
2. **Mock Configurations**: You are free to seed static configuration structs where necessary (e.g., `&infra.Config{ JWTSecret: "test-secret"}`).
3. **Database Clearing**: Because PostgreSQL containers run consistently through the `TestE2E...()` wrapper func block, subsequent `t.Run()` events share state. **You must clear rows and re-seed explicitly inside each internal test block boundaries!**
   > [!WARNING]
   > Forgetting to issue `res.DB.Exec("DELETE FROM table")` across `t.Run()` blocks is the #1 cause behind test failures and data collisions!

## Structuring the Module File
All comprehensive E2E tests are stored in `tests/e2e/e2e_{resource}_test.go`.

- **Setup Factory func**: Keep a `setup{Resource}TestServer()` at the top scope.
- **Main test runner**: Provide a single overarching `TestE2e{Resource}(t *testing.T)` block per file.
- **Defer Container teardowns**: 
```go
defer func() { _ = res.PostgresContainer.Terminate(context.Background()) }()
defer func() { _ = res.RedisContainer.Terminate(context.Background()) }()
defer func() { _ = res.DB.Close() }()
```
