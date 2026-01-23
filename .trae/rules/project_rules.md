# Project Rules

<!-- AI_CONTEXT: AI 代码生成核心指导原则 -->

## 1. 核心原则 (Core)

- **交互**: 始终用中文。关键逻辑必加注释。代码清晰 > 巧妙。
- **开发**: SpecFirst (OpenSpec) -> TDD (Red-Green-Refactor)。
- **调试**: 遇到问题尝试 3 次后求助。

## 2. 技术栈 (Stack)

- **Go**: v1.25, Gin. 泛型, Embed. 严禁修改 `AppError` 核心逻辑。
- **Web**: HTMX ^2.0, Alpine ^3.15, Tailwind v4. 逻辑分离 (`x-data` 优于内联)。
- **架构**: DDD. `internal/`(核心), `website/`(Web). 接口注入 > 单例。
- **I18n**: 双层策略。后端 `website/locales` (模板), 前端 `js/locales` (交互)。

## 3. 质量规范 (Quality)

- **限制**: 函数 ≤60 行, 参数 ≤5. 圈复杂度 ≤10.
- **检查**: 必过 golangci-lint / ESLint. 零 TODOs.
- **重构**: 仅优化结构/命名，严禁修改业务逻辑/API 契约。

## 4. 红线 (Critical)

- 🚫 提交无法编译代码/暴露密钥。
- 🚫 随意引入新依赖。
- ✅ 保持 API 向后兼容。
- ✅ 增量提交，小步快跑。
