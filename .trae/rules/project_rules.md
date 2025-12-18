# Project Rules

本文档定义了 iDiary.ai Hebe 项目的开发规范和配置要求。

<!-- AI_CONTEXT: 本文档为 AI 代码生成提供核心指导原则 -->

## 1. 核心协议与原则 (Core Protocols & Principles)

<!-- AI_RULE: MANDATORY -->

### AI 交互规范

- **语言偏好**: 请始终用中文回复。
- **注释要求**: 提供代码时，必须为关键点和复杂逻辑添加中文注释。
- **代码粒度**: 当生成的代码超过 20 行时，应主动评估粒度是否合适，并考虑拆分或整合。

### 基本信念

- **渐进式改进胜过大爆炸式变更** - 小的变更能够编译并通过测试
- **从现有代码中学习** - 在实施前先研究和规划
- **实用主义胜过教条主义** - 适应项目实际情况
- **清晰意图胜过巧妙代码** - 选择无聊且明显的解决方案

### 简洁性实践

- ✅ 每个函数/类单一职责
- ✅ 使用标准库函数而非自定义实现
- ✅ 优先选择已验证的设计模式
- ❌ 避免过早抽象
- ❌ 不使用巧妙技巧和难以理解的语言特性
- ❌ 如果需要解释，说明过于复杂

<!-- /AI_RULE -->

## 2. 决策与验证流程 (Decision & Workflow)

<!-- AI_CHECKLIST: MANDATORY -->

### 生成前检查清单 (Pre-Generation Checklist)

#### 1. 技术选型与架构

- [ ] 是否符合技术栈版本要求 (Go 1.25, HTMX 2.0.8, Alpine 3.15+)?
- [ ] 是否遵循现有架构模式 (DDD, Interface-based)?
- [ ] 是否使用项目现有的库和工具 (避免引入新依赖)?
- [ ] **国际化**: 是否需要 i18n 支持？(后端模板优先 > 前端 JS)

#### 2. 代码质量与规范

- [ ] 函数长度 ≤ 60 行（Go）/ 50 行（JS）？
- [ ] 参数数量 ≤ 5 个（Go）/ 4 个（JS）？
- [ ] 是否包含相应测试？
- [ ] 是否遵循项目命名约定？
- [ ] **规范性**: 是否符合 golangci-lint/ESLint 规范？

#### 3. 鲁棒性与安全

- [ ] 是否使用依赖注入而非单例？
- [ ] 是否正确处理错误 (不吞噬异常)?
- [ ] **异步反馈**: 异步操作是否使用了 `hx-indicator`？
- [ ] **安全性**: 是否处理了 CSRF？是否避免了敏感信息泄露？

<!-- /AI_CHECKLIST -->

<!-- AI_WORKFLOW: CORE -->

### 开发流程规范

#### 1. SpecFirst 开发模式

- **核心原则**: 在编写代码前锁定需求 (Lock intent before implementation)。
- **标准**: 采用 [Fission-AI/OpenSpec](https://github.com/Fission-AI/OpenSpec) 规范。
- **流程**: Draft (Proposal) -> Review -> Implement -> Archive。

#### 2. TDD 实施流程 (Red-Green-Refactor)

1. **🔍 理解 (Understand)** - 研究代码库中的现有模式
2. **🔴 测试 (Test)** - 先写测试
3. **🟢 实现 (Implement)** - 最少代码通过测试
4. **🔧 重构 (Refactor)** - 在测试通过的情况下清理代码
5. **📝 提交 (Commit)** - 使用清晰的消息链接到计划

#### 3. 决策框架

当存在多个有效方法时，基于以下标准选择：

1. **兼容性**: 保持向后兼容 (API/DB/Interface)
2. **可测试性**: 确保代码可测试 (Interface over Struct)
3. **可读性**: 代码自文档化
4. **简洁性**: 有效的最简单解决方案

#### 4. 遇到困难时 (3-Attempt Rule)

每个问题最多尝试 3 次。如果失败：

1. 记录失败内容和错误消息。
2. 研究 2-3 个类似的现有实现。
3. 质疑基本假设 (抽象层次是否正确?)。
4. 寻求帮助或尝试完全不同的思路。

<!-- /AI_WORKFLOW -->

## 3. 架构与目录规范 (Architecture & Structure)

<!-- AI_ARCHITECTURE: DDD -->

### DDD 架构原则

- **高内聚按领域组织**: 相关业务逻辑聚合在同一领域。
- **组合胜过继承**: 使用依赖注入。
- **接口胜过单例**: 启用测试和灵活性。
- **显式胜过隐式**: 清晰的数据流和依赖关系。

### 目录结构规范

- **`docs/`**: 项目文档 (架构、指南、API)。必须包含 README 和变更日志。
- **`tests/`**: 测试代码 (按模块和类型组织)。
- **`website/`**: 前端代码 (模板、静态资源)。
- **`internal/`**: 核心业务逻辑 (Go)。

<!-- /AI_ARCHITECTURE -->

## 4. 技术栈规范 (Tech Stack)

<!-- AI_TECH_STACK: MANDATORY -->

### 4.1 后端规范 (Go/Gin)

**基本要求**:

- **版本**: Go 1.25
- **特性**: Go Modules, Go Work, Generics, Embed
- **工具**: golangci-lint, gofmt, go vet

**代码模式 (Patterns)**:

```go
// ✅ 正确：使用泛型与 Embed
func Process[T any](items []T) []T { ... }
//go:embed static/*
var staticFiles embed.FS

// ✅ 正确：Gin 路由组和中间件
api := r.Group("/api/v1")
api.Use(middleware.Auth())
api.POST("/users", handler.CreateUser)

// ✅ 正确：依赖注入 (Interface over Struct)
type UserService struct {
    repo UserRepository
}
func NewUserService(repo UserRepository) *UserService { ... }
```

**错误处理规范 (CRITICAL)**:

- **禁止修改**: `github.com/biuren/hebe/internal/utils/error` 定义的 `*errors.AppError` 类型。
- **一致性**: 保持 `Code` 和 `Key` 字段的一致性。
- **原则**: 快速失败，提供描述性消息，不吞噬异常。

### 4.2 前端规范 (JS/HTMX/UI)

#### JavaScript (Alpine.js)

- **版本**: ES2023 / Alpine.js ^3.15.2
- **原则**:
  - ✅ **逻辑分离**: 业务逻辑位于 `/website/static/src/js/modules/`
  - ✅ **声明式交互**: 优先使用 `x-data`, `x-on`
  - ❌ **禁止**: 内联 JS (`onclick`), jQuery
- **最佳实践**:
  ```javascript
  // ✅ Alpine 组件
  export default () => ({
    isOpen: false,
    toggle() {
      this.isOpen = !this.isOpen;
    },
  });
  ```

#### HTMX (核心交互)

- **版本**: ^2.0.8
- **核心交互模式**:
  - **服务端驱动**: 返回 HTML 片段 (Fragment)。
  - **副作用 (Side Effects)**: 使用 `HX-Trigger` 处理 UI 反馈 (如 `{"show-toast": "Success"}`)。
  - **CSRF 保护**: 所有 `hx-post/put/delete` 必须包含 Token (通过 `htmx:configRequest` 注入)。
  - **错误处理**: Form 验证失败返回 HTTP 200 (带错误片段)；系统错误返回 4xx/5xx (前端监听 `htmx:responseError`)。

**场景示例**:

```html
<!-- 局部刷新 -->
<div id="list" hx-get="/list?p=2" hx-select="#list" hx-swap="outerHTML">
  ...
</div>

<!-- 表单提交与反馈 -->
<form
  hx-post="/api"
  hx-on::after-request="if(event.detail.successful) this.reset()"
>
  ...
</form>

<!-- 搜索防抖 -->
<input
  hx-get="/search"
  hx-trigger="keyup changed delay:500ms"
  hx-target="#results"
/>

<!-- 加载指示器 -->
<button hx-post="/save" hx-indicator="#loading">Save</button>
```

#### UI 框架 (Tailwind & Uno)

- **Tailwind CSS**: v4.x, 使用 Utility Classes, `hx-boost="true"` 实现类 SPA 导航。
- **Uno UI**: 使用编译版 (`uno.umd.js`), 插件优先 (Dialog, Toast)。

<!-- /AI_TECH_STACK -->

## 5. 国际化 (i18n) 规范

<!-- AI_I18N: MANDATORY -->

### 后端国际化 (Go/Gin) - **Primary**

- **路径**: `website/locales` (JSON)
- **用途**: 服务端模板渲染 (`{{t "key"}}`)，API 错误消息。
- **策略**: 默认 `en`, 回退 `zh-cn`。

### 前端国际化 (Website JS) - **Secondary**

- **路径**: `website/static/src/js/locales` (JS Module)
- **用途**: Alpine.js 组件动态交互文本 (`$store.i18n.t('key')`)。
- **策略**: 页面加载时由服务端注入或异步加载。

### 检查清单

- [ ] 所有用户可见文本是否已国际化？
- [ ] 是否使用正确的命名空间？
- [ ] 是否提供英文和中文翻译？
- [ ] 日期、数字格式是否本地化？

<!-- /AI_I18N -->

## 6. 质量与重构体系 (Quality & Refactoring)

<!-- AI_REFACTORING: CRITICAL -->

### 重构核心原则 (Immutable Rules)

⚠️ **严格遵守以下规则，任何违反将被视为无效**：

1. **业务逻辑不变**: 禁止修改算法、流程控制、数据状态读写、验证规则。
2. **外部契约不变**: 禁止修改公共函数签名、API 路径/格式、DB 表结构。
3. **允许修改**: 代码结构优化、命名改进、注释完善、性能优化、错误处理改进。

### 代码坏味道识别 (Code Smells)

基于 Martin Fowler 指南：

1. **神秘命名** -> 重命名
2. **重复代码** -> 提取函数
3. **过长函数** -> 提取函数
4. **过长参数 (>4)** -> 引入参数对象
5. **基本类型偏执** -> 以对象取代
6. **散弹式修改** -> 搬移函数/字段

### 复杂度控制与工具

- **Go**: 函数 ≤ 60 行, 参数 ≤ 5, 圈复杂度 ≤ 10. (必过 `golangci-lint`)
- **JS**: 函数 ≤ 50 行, 参数 ≤ 4, 圈复杂度 ≤ 8. (必过 `ESLint`)

### 质量门控 (Definition of Done)

- [ ] 测试已编写并通过
- [ ] 代码遵循项目约定 (无 Linter 警告)
- [ ] 提交消息清晰 (解释"为什么")
- [ ] 实现符合计划
- [ ] 无 TODOs

<!-- /AI_REFACTORING -->

## 7. 禁忌与红线 (Critical Rules)

<!-- AI_CRITICAL: ALWAYS_APPLY -->

**🚫 永远不要 (Never)**:

- 使用 `--no-verify` 绕过提交钩子。
- 禁用测试而不是修复它们。
- 提交无法编译的代码。
- 在代码中暴露密钥/敏感信息。
- 随意引入新的外部依赖 (除非必要)。
- 修改现有错误处理的核心逻辑 (`AppError`).

**✅ 总是 (Always)**:

- 增量提交，小步快跑。
- 保持 API 向后兼容。
- 遇到问题尝试 3 次后求助/换思路。
- 更新相关文档。
- 确保所有新功能都有测试覆盖。

<!-- /AI_CRITICAL -->

---

<!-- AI_FOOTER -->

**📝 文档版本**: v3.1 (Refactored)
**🔄 最后更新**: 2025-12-18
**🎯 适用范围**: iDiary.ai Hebe 项目
**⚡ 优先级**: 本文档所有规则对 AI 代码生成具有强制约束力。

<!-- /AI_FOOTER -->
