# Proposal: Generic Validation Middleware & Unified Response

**Status**: Implemented
**Date**: 2025-12-18
**Author**: Assistant

## 1. Background

To modernize the Go backend and align with the `GOLANG_OPTIMIZATION_PLAN.md` (Phase 2), we need to replace manual request binding and validation in controllers with a generic, middleware-driven approach. This reduces boilerplate, ensures type safety, and enforces consistent API error handling.

## 2. Changes

### 2.1 Unified Response Structure

Refactored `internal/app/view/response.go` to use generics:

```go
type APIResponse[T any] struct {
    Success   bool   `json:"success"`
    Message   string `json:"message,omitempty"`
    Data      T      `json:"data,omitempty"`
    RequestID string `json:"request_id"`
}
```

### 2.2 Generic Validation Middleware

Introduced `internal/app/middleware/validator.go`:

- `Validate[T any]()`: Middleware that binds and validates the request struct `T`.
- `MustGetValidatedRequest[T any](c)`: Type-safe helper to retrieve the validated struct in controllers.

### 2.3 Controller Refactoring

Refactored the following controllers to use the new middleware and response format:

- **User Controller** (`internal/app/controller/user`): Login, Register, Profile, Password.
- **Diary Controller** (`internal/app/controller/diary`): Create, Update, Like.
- **Post Controller** (`internal/app/controller/post`): Create, Update, Delete, Like.
- **Comment Controller** (`internal/app/controller/comment`): Create, Delete.

**Example Change (Diary Controller):**

*Before:*
```go
var req CreateDiaryRequest
if err := c.ShouldBind(&req); err != nil { ... }
```

*After:*
```go
// Route
authAPI.POST("/diary", middleware.Validate[CreateDiaryRequest](), ctrl.CreateDiary)

// Handler
req := middleware.MustGetValidatedRequest[CreateDiaryRequest](c)
```

## 3. Frontend Synchronization

Updated frontend modules to align with RESTful routes (PUT/DELETE) and unified response structure:

- `diary.js`: Updated `_handleDiarySubmit` (PUT) and `_handlePostSubmit`.
- `member.js`: Updated `submitDiaryForm` (POST).
- `auth.js`: Merged validation logic and exposed `ValidForm` component.
- `uno-shim.js`: Implemented `Box`, `Notify`, `Valid` shims for Alpine store access.

## 4. Verification

- Validated that `Box` and `Notify` work via `uno-shim.js`.
- Confirmed `emptyquit` validation rule usage in templates.
- Verified RESTful route alignment (e.g., `PUT /diary/:id`).
