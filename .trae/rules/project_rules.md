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

## 4. Red Lines (Critical)

- 🚫 Do not commit code that fails to compile or exposes secrets.
- 🚫 Do not introduce new dependencies arbitrarily.
- ✅ Maintain API backward compatibility.
- ✅ Incremental commits, small steps.
