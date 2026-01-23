---
name: backend-best-practices
description: Comprehensive guidelines for Go backend development using Gin, adhering to project-specific architecture and standards.
license: Apache-2.0
---

# Backend Best Practices (Go + Gin)

This skill provides actionable standards for the **iDiary.ai** backend, ensuring consistency, scalability, and maintainability.

## 1. Project Structure & Organization

Follow the Standard Go Project Layout with DDD adaptations:

```
/
├── cmd/                    # Main applications
│   └── server/             # Entry point (main.go)
├── internal/               # Private application code
│   ├── app/                # Application Layer (HTTP handlers, routes)
│   │   ├── controller/     # Request handling logic
│   │   ├── middleware/     # Gin middleware
│   │   └── route/          # Router setup
│   ├── domain/             # Domain Layer (Interfaces, Entities/Structs)
│   │   └── [entity]/       # e.g., diary, user
│   ├── infra/              # Infrastructure Layer (Repositories, External APIs)
│   │   └── repository/     # DB implementations
│   ├── core/               # Core utilities (DB, DI, Logger)
│   └── utils/              # Shared utilities (Config, Errors)
├── pkg/                    # Library code ok to use by external apps (rare)
├── configs/                # Configuration files
└── scripts/                # Build/Maintenance scripts
```

## 2. Code Style & Quality

-   **Formatting**: MUST run `gofumpt -l -w .` (stricter than `gofmt`).
-   **Linting**: MUST pass `golangci-lint`.
    -   **Recommended Linters**: `govet`, `staticcheck`, `errcheck`, `gocyclo` (limit 10), `revive`.
-   **Comments**:
    -   Exported functions MUST have GoDoc comments.
    -   Start with the function name: `// GetUser retrieves...`
    -   Explain *why*, not just *what*.

## 3. Gin Framework Patterns

### Routing & DI
-   **Dependency Injection**: Use `internal/core/di` to inject services into controllers.
-   **Registration**: Register routes in `internal/app/route/router.go` using `SetWithContainer`.

### Request/Response
-   **Binding**: Use `ShouldBindJSON` or `ShouldBindUri`.
-   **Validation**: Use `binding:"required,email"` tags.
-   **Response**: Use `view.RenderJSON` (or similar helper) for consistent JSON structure.

```go
// BAD
c.JSON(200, gin.H{"data": data})

// GOOD
view.RenderJSON(c, data)
```

### Context
-   **Timeout**: Respect `c.Request.Context()` for DB operations.

## 4. Error Handling

-   **Type**: Use `internal/utils/errors.AppError`.
-   **Propagation**:
    -   Repository: Return `error` (wrap standard errors).
    -   Service: Map errors to `AppError` if specific business logic fails.
    -   Controller: Check `AppError` code to determine HTTP status.
-   **Wrapping**: Use `fmt.Errorf("context: %w", err)` to preserve stack/chain.

```go
if err != nil {
    return errors.NewAppError(errors.ErrDatabase, "db_error", 500)
}
```

## 5. Configuration

-   **Library**: `spf13/viper`.
-   **Pattern**: 12-Factor App. Use Environment Variables (`IDIARY_...`) for secrets.
-   **Access**: Access via `config.GetConfig()`. Do not read files directly in business logic.

## 6. Database (MySQL + sqlx)

-   **Library**: `jmoiron/sqlx`.
-   **Pattern**: Repository Pattern (defined in `internal/domain`, implemented in `internal/infra`).
-   **Queries**:
    -   Use `db.Get` for single rows, `db.Select` for slices.
    -   Prefer `NamedExec` for structs to avoid argument counting errors.
-   **Transactions**: Use `tx := db.Beginx()` and ensure `defer tx.Rollback()`.

## 7. Logging

-   **Library**: Standard `log/slog` (Go 1.21+) for structured logging + `lumberjack.v2` for rotation.
-   **Tracing**: Inject `TraceID` in middleware and include it in all logs.

## 8. Testing

-   **Unit**: Use `testify/assert`. Mock interfaces using `mockery` or manual mocks.
-   **DB**: Use `DATA-DOG/go-sqlmock` for repository unit tests.
-   **Integration**: Use `httptest` to test Gin handlers with a real (or Dockerized) DB.

```go
func TestUserHandler(t *testing.T) {
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    // ... setup and call handler
    assert.Equal(t, 200, w.Code)
}
```

## 9. Security

-   **Input**: Sanitize all user input. Rely on `html/template` auto-escaping.
-   **Secrets**: NEVER commit secrets. Use env vars.
-   **Auth**: Validate `Authorization` header in middleware.

## 10. Performance

-   **Database**: Index foreign keys and query predicates.
-   **Pooling**: Configure `SetMaxOpenConns` and `SetMaxIdleConns` in `database/connection.go`.
-   **Pprof**: Enable `pprof` in debug mode only.

## 11. Deployment

-   **Health**: Implement `/healthz` returning 200 OK.
-   **Graceful Shutdown**: Use `context.WithTimeout` on `server.Shutdown()`.
