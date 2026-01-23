---
name: fullstack-testing
description: Testing strategy for the iDiary.ai ecosystem, covering Uno UI (Frontend Lib), Website (Consumer), and Go Backend.
license: Apache-2.0
---

# Full-Stack Testing Strategy

This skill defines the testing pyramid and tools for the **iDiary.ai** project.

## 1. Frontend: Uno UI Framework (`/uno`)

**Target**: The reusable UI library.
**Tools**: `Vitest`, `Happy-DOM`, `Testing Library`.

### Unit/Component Tests
-   **Location**: `tests/components/*.test.js`
-   **Strategy**: Test behavior, not implementation details.
-   **Requirements**:
    -   Test all props configurations.
    -   Test event emissions.
    -   Test slot rendering.
    -   Test A11y (ARIA attributes).

```javascript
// Example: Valid Plugin Test
it('should show error on invalid input', async () => {
  const { container } = render(ValidInput)
  await userEvent.type(container.querySelector('input'), 'invalid')
  expect(container.querySelector('.error-msg')).toBeVisible()
})
```

## 2. Frontend: Website Consumer (`/website`)

**Target**: The actual application code using Uno UI.
**Problem**: Logic tightly coupled in `static/src/js`.
**Tools**: `Playwright` (E2E), `Vitest` (Unit for extracted logic).

### Refactoring for Testability
-   **Action**: Extract complex logic from `modules/*.js` into pure functions in `utils/`.
-   **Test**: Write Unit Tests for these extracted functions.

### End-to-End (E2E) Testing
-   **Tool**: Playwright.
-   **Scope**: Critical User Journeys (CUJs).
    1.  **Auth**: Register -> Login -> Logout.
    2.  **Core Loop**: Write Diary -> Save -> View in List -> Edit -> Delete.
    3.  **Social**: Like -> Comment.
-   **Visual Regression**: Use Playwright snapshots for CSS changes.

```javascript
// tests/e2e/auth.spec.ts
test('user can login', async ({ page }) => {
  await page.goto('/login');
  await page.fill('input[name="email"]', 'user@example.com');
  await page.fill('input[name="password"]', 'password');
  await page.click('button[type="submit"]');
  await expect(page).toHaveURL('/dashboard');
});
```

## 3. Backend: Go API (`/internal`)

**Target**: API and Business Logic.
**Tools**: `Testify`, `Go-SQLMock`, `httptest`.

### Layered Testing
-   **Repository**: Unit test using `go-sqlmock` to verify SQL generation and error handling.
-   **Service**: Unit test using Mocks for Repositories. Test business rules.
-   **Controller**: Integration test using `httptest`.

### Mocking Strategy
-   **Interfaces**: Define interfaces for all services in `domain`.
-   **Generation**: Use `mockery` to generate mocks from interfaces.

## 4. CI/CD Integration

### Workflow
1.  **Commit**: Run `gofumpt`, `golangci-lint`, `eslint`.
2.  **PR**:
    -   Run Backend Unit Tests (`go test ./...`).
    -   Run Uno UI Tests (`npm run test`).
    -   Run E2E Smoke Tests (headless).
3.  **Merge**: Run full E2E suite and deploy to staging.

### Quality Gates
-   **Coverage**: >80% for Backend Domain Logic and Uno Core.
-   **Visual**: 0 regression differences (auto-fail if snapshots differ).

## 5. Test Data Management

-   **Fixtures**: Use SQL seeds for reference data (e.g., categories).
-   **Factories**: Use a factory pattern (Go) or helper functions (JS) to generate dynamic test data.
-   **Cleanup**:
    -   **DB**: Use transactions in tests and rollback.
    -   **E2E**: Use a dedicated test DB instance; truncate tables after suite.
