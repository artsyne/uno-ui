# Project Rules

<!-- AI_CONTEXT: Core Guidelines for AI Code Generation -->

## 1. Core Principles (Core)

- **Interaction**: Always use Chinese. Critical logic must be commented. Clear code > Clever code.
- **Development**: SpecFirst (OpenSpec) -> TDD (Red-Green-Refactor).
- **Skills**:
  - Creative Design: [Frontend Design](.trae/skills/frontend-design/SKILL.md)
  - Brand Guidelines: [Brand Guidelines](.trae/skills/brand-guidelines/SKILL.md)
  - Backend Standards: [Backend Best Practices](.trae/skills/backend-best-practices/SKILL.md)
  - Testing Strategy: [Fullstack Testing](.trae/skills/fullstack-testing/SKILL.md)
- **Debugging**: Try 3 times before asking for help.

## 2. Tech Stack (Stack)

- **Go**: v1.25, Gin. Generics, Embed. Strictly forbidden to modify `AppError` core logic.
- **Web**: HTMX ^2.0, Alpine ^3.15, Tailwind v4. Logic separation (`x-data` preferred over inline).
- **Architecture**: DDD. `internal/` (Core), `website/` (Web). Interface injection > Singleton.
- **I18n**: Dual strategy. Backend `website/locales` (Templates), Frontend `js/locales` (Interaction).

## 3. Quality Standards (Quality)

- **Build & Test**: Must compile and pass all unit tests.
- **Limits**: Functions ≤60 lines, Args ≤5. Cyclomatic complexity ≤10.
- **Checks**: Must pass golangci-lint / ESLint. Zero TODOs.
- **Refactoring**: Optimize structure/naming only. Strictly forbidden to modify business logic/API contracts.

## 5. Verification Protocol (Strict)

Before concluding any coding session or task, the following **Mandatory Checklist** must be verified:

1.  **Test Suite Execution**:
    -   Run ALL unit and integration tests (`go test ./...`).
    -   **Requirement**: 100% Pass Rate. No skipped tests without valid reason.
2.  **Build Verification**:
    -   Execute full project build/compile (`go build ./...`).
    -   **Requirement**: Zero compile errors.
3.  **Code Coverage**:
    -   Ensure critical paths (especially new logic) have sufficient test coverage.
4.  **Technical Debt Check**:
    -   Confirm no new linting errors (golangci-lint).
    -   Verify no temporary hacks or "TODO" comments left behind.

**Command Reference**:
```bash
go test ./... && go build ./...
```

