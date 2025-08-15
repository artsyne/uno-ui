# Project Rules

本文档定义了 iDiary.ai Hebe 项目的开发规范和配置要求。

<!-- AI_CONTEXT: 本文档为 AI 代码生成提供核心指导原则 -->

## 1. 核心开发原则

<!-- AI_RULE: MANDATORY -->
### 基本信念

- **渐进式改进胜过大爆炸式变更** - 小的变更能够编译并通过测试
- **从现有代码中学习** - 在实施前先研究和规划
- **实用主义胜过教条主义** - 适应项目实际情况
- **清晰意图胜过巧妙代码** - 选择无聊且明显的解决方案
<!-- /AI_RULE -->

<!-- AI_PATTERN: PREFERRED -->
### 简洁性实践

- ✅ 每个函数/类单一职责
- ✅ 使用标准库函数而非自定义实现
- ✅ 优先选择已验证的设计模式
- ❌ 避免过早抽象
- ❌ 不使用巧妙技巧和难以理解的语言特性
- ❌ 如果需要解释，说明过于复杂
<!-- /AI_PATTERN -->

## 2. AI 代码生成决策流程

<!-- AI_CHECKLIST: MANDATORY -->
### 生成前检查

#### 步骤 1: 技术选型验证
- [ ] 是否符合技术栈版本要求？
- [ ] 是否遵循现有架构模式？
- [ ] 是否需要国际化支持？
- [ ] 是否使用项目现有的库和工具？

#### 步骤 2: 代码质量验证
- [ ] 函数长度 ≤ 60 行（Go）/ 50 行（JS）？
- [ ] 参数数量 ≤ 5 个（Go）/ 4 个（JS）？
- [ ] 是否包含相应测试？
- [ ] 是否遵循项目命名约定？

#### 步骤 3: 架构原则验证
- [ ] 是否使用依赖注入而非单例？
- [ ] 是否正确处理错误？
- [ ] 是否包含必要的注释？
- [ ] 是否符合 golangci-lint 规范？
<!-- /AI_CHECKLIST -->

### 优先级排序

1. **安全性** - 永远不暴露敏感信息
2. **兼容性** - 保持向后兼容
3. **可测试性** - 确保代码可测试
4. **可读性** - 代码自文档化
5. **性能** - 在满足功能前提下优化

## 3. 目录结构规范

### 文档目录
- **路径**: `docs/`
- **用途**: 存放项目文档，包括架构设计、开发指南、API 文档等
- **格式**: Markdown
- **要求**: 必须包含 README、API 文档和变更日志

### 测试目录
- **路径**: `tests/`
- **用途**: 存放所有测试代码，包括单元测试、集成测试等
- **结构**: 按模块和测试类型组织

## 4. 技术栈规范

<!-- AI_TECH_STACK: MANDATORY -->

### Go 语言规范

<!-- AI_VERSION: 1.25 -->
- **版本**: 1.25
- **特性要求**:
  - ✅ 使用 Go Modules
  - ✅ 使用 Go Work
  - ✅ 使用泛型 (Generics)
  - ✅ 使用 embed 包

<!-- AI_GO_PATTERNS -->
**代码模式示例**:
```go
// ✅ 正确：使用泛型
func Process[T any](items []T) []T { ... }

// ✅ 正确：使用 embed
//go:embed static/*
var staticFiles embed.FS

// ❌ 错误：避免 interface{}
func Process(items interface{}) interface{} { ... }
```
<!-- /AI_GO_PATTERNS -->

### Gin 框架规范

<!-- AI_VERSION: 1.10.0 -->
- **版本**: 1.10.0
- **使用规范**:
  - ✅ 使用 Context
  - ✅ 使用中间件 (Middleware)
  - ✅ 使用路由组 (Group)
  - ✅ 使用数据绑定 (Binding)

<!-- AI_GIN_PATTERNS -->
**代码模式示例**:
```go
// ✅ 正确：路由组和中间件
api := r.Group("/api/v1")
api.Use(middleware.Auth())
api.POST("/users", handler.CreateUser)

// ✅ 正确：数据绑定
var req CreateUserRequest
if err := c.ShouldBindJSON(&req); err != nil {
    c.JSON(400, gin.H{"error": err.Error()})
    return
}
```
<!-- /AI_GIN_PATTERNS -->

### JavaScript 规范

<!-- AI_VERSION: ES2023 -->
- **版本**: ES2023
- **浏览器支持**:
  - Chrome 110+
  - Edge 110+
  - Firefox 128 ESR+

<!-- AI_JS_CHECKLIST -->
**编码规范检查清单**:
- [ ] 使用严格模式 (`'use strict'`)
- [ ] 使用模块化 (`import/export`)
- [ ] 使用 `async/await` 而非 Promise.then()
- [ ] 使用箭头函数 (`() => {}`)
- [ ] 使用 `const/let` 而非 `var`
- [ ] 使用模板字符串 (`` `${variable}` ``)
- [ ] 使用解构赋值 (`const {a, b} = obj`)
- [ ] 使用展开运算符 (`...array`)
- [ ] 使用可选链 (`obj?.prop`)
- [ ] 使用空值合并运算符 (`value ?? default`)
<!-- /AI_JS_CHECKLIST -->

<!-- AI_JS_PATTERNS -->
**代码模式示例**:
```javascript
// ✅ 正确：现代 JavaScript
const fetchData = async (url) => {
    try {
        const response = await fetch(url);
        const data = await response.json();
        return data?.results ?? [];
    } catch (error) {
        console.error(`Failed to fetch: ${error.message}`);
        return [];
    }
};

// ❌ 错误：旧式写法
function fetchData(url) {
    return fetch(url).then(function(response) {
        return response.json();
    }).then(function(data) {
        return data.results || [];
    });
}
```
<!-- /AI_JS_PATTERNS -->

### jQuery 规范

<!-- AI_VERSION: 4.x -->
- **版本**: 4.x
- **使用规范**:
  - ✅ 使用现代 API
  - ❌ 避免已弃用的方法

### Tailwind CSS 规范

<!-- AI_VERSION: 4.x -->
- **版本**: 4.x
- **使用规范**:
  - ✅ 使用现代类名
  - ✅ 优先使用 Utility Classes
  - ❌ 避免自定义 CSS（除非必要）

## 5. 国际化 (i18n) 规范

<!-- AI_I18N: MANDATORY -->

### Kiko 插件国际化
- **类型**: 插件
- **路径**: `kiko/src/i18n`
- **策略**:
  - 默认语言: `en`
  - 回退语言: `zh-CN`
  - 命名空间: `kiko`

### Website 应用国际化
- **类型**: 应用程序
- **路径**: `website/static/src/js/i18n`
- **策略**:
  - 默认语言: `en`
  - 回退语言: `zh-CN`
  - 命名空间: `website.js`
  - 作用域: `website/static/src/js`

<!-- AI_I18N_CHECKLIST -->
**国际化检查清单**:
- [ ] 所有用户可见文本是否已国际化？
- [ ] 是否使用正确的命名空间？
- [ ] 是否提供英文和中文翻译？
- [ ] 日期、数字格式是否本地化？
<!-- /AI_I18N_CHECKLIST -->

## 6. 代码质量要求

<!-- AI_QUALITY: CRITICAL -->

### 重构核心原则

<!-- AI_REFACTOR_RULES: IMMUTABLE -->
**核心重构原则：业务逻辑和外部接口契约不变性**

⚠️ **严格遵守以下规则进行 Go 语言代码重构，任何违反将被视为无效**：

#### 1. 业务逻辑（Behavioral Logic）必须保持完全一致

**🚫 禁止修改**:
- 算法、计算公式、数据转换规则
- 条件判断（`if/else`, `switch`）、循环（`for`）、流程控制
- 对数据状态的读写操作（数据库操作、缓存操作、文件 I/O）
        - 验证规则、业务决策点
- 任何会改变最终输出结果或产生副作用的行为

#### 2. 外部接口契约必须保持不变

**🚫 禁止修改**:
- 公共函数/方法的签名（名称、参数类型、返回类型）
- 公共结构体的导出字段
- HTTP API 的路径、方法、请求/响应格式
- 数据库表结构和字段名称

**✅ 允许修改**:
- 代码组织和结构优化
- 变量和函数命名改进
- 注释和文档完善
- 性能优化（不改变行为）
- 错误处理改进
- 代码风格统一
<!-- /AI_REFACTOR_RULES -->

### 代码复杂度控制

<!-- AI_COMPLEXITY_LIMITS -->
#### Go 语言限制
- **函数长度**: ≤ 60 行
- **参数数量**: ≤ 5 个
- **嵌套层级**: ≤ 4 层
- **圈复杂度**: ≤ 10

#### JavaScript 限制
- **函数长度**: ≤ 50 行
- **参数数量**: ≤ 4 个
- **嵌套层级**: ≤ 3 层
- **圈复杂度**: ≤ 8
<!-- /AI_COMPLEXITY_LIMITS -->

### 代码检查工具

<!-- AI_LINTING_TOOLS -->
#### Go 语言工具
- **golangci-lint**: 必须通过所有检查
- **gofmt**: 代码格式化
- **go vet**: 静态分析
- **go test**: 单元测试覆盖率 ≥ 80%

#### JavaScript 工具
- **ESLint**: 必须通过所有检查
- **Prettier**: 代码格式化
- **JSHint**: 静态分析
<!-- /AI_LINTING_TOOLS -->

## 7. 开发指导原则

<!-- AI_DEVELOPMENT_PRINCIPLES: CORE -->

### 核心理念

#### 渐进式改进
- 🎯 **目标**: 每次变更都能编译并通过测试
- 📋 **实践**: 小步快跑，频繁集成
- ⚡ **优势**: 降低风险，快速反馈

#### 从现有代码学习
- 🔍 **研究**: 在实施前先分析现有模式
- 📚 **学习**: 理解项目的架构决策
- 🎨 **模仿**: 保持代码风格一致性

#### 实用主义优先
- 🛠️ **适应**: 根据项目实际情况调整
- ⚖️ **平衡**: 在理想与现实间找到平衡
- 🎯 **聚焦**: 解决实际问题而非追求完美

#### 清晰胜过巧妙
- 📖 **可读**: 代码应该像散文一样易读
- 🔧 **简单**: 选择最直接的解决方案
- 🚫 **避免**: 过度工程和炫技
<!-- /AI_DEVELOPMENT_PRINCIPLES -->
    ### 开发流程

<!-- AI_WORKFLOW: TDD_DDD -->

#### 1. 规划与分阶段 (Planning & Staging)

将复杂工作分解为 3-5 个阶段。在 `IMPLEMENTATION_PLAN.md` 中记录：

```markdown
## Stage N: [Name]
**Goal**: [Specific deliverable]
**Success Criteria**: [Testable outcomes]
**Tests**: [Specific test cases]
**Status**: [Not Started|In Progress|Complete]
```

- 随着进展更新状态
- 所有阶段完成后删除文件

#### 2. TDD 实施流程 (Test-Driven Development)

**🔄 严格按照 TDD 红-绿-重构循环**：

1. **🔍 理解 (Understand)** - 研究代码库中的现有模式
2. **🔴 测试 (Test)** - 先写测试 (红色)
3. **🟢 实现 (Implement)** - 最少代码通过测试 (绿色)
4. **🔧 重构 (Refactor)** - 在测试通过的情况下清理代码
5. **📝 提交 (Commit)** - 使用清晰的消息链接到计划

#### 3. DDD 架构原则 (Domain-Driven Design)

**🏗️ 领域驱动设计模式**：
- **高内聚按领域组织** - 相关业务逻辑聚合在同一领域
- **组合胜过继承** - 使用依赖注入
- **接口胜过单例** - 启用测试和灵活性
- **显式胜过隐式** - 清晰的数据流和依赖关系
- **尽可能测试驱动** - 永远不要禁用测试，而是修复它们

#### 4. 遇到困难时 (After 3 Attempts)

**⚠️ 关键**: 每个问题最多尝试3次，然后停止。

1. **📋 记录失败内容**:
   - 你尝试了什么
   - 具体错误消息
   - 你认为失败的原因

2. **🔍 研究替代方案**:
   - 找到 2-3 个类似的实现
   - 注意使用的不同方法

3. **❓ 质疑基本假设**:
   - 这是正确的抽象层次吗？
   - 能否分解为更小的问题？
   - 是否有完全更简单的方法？

4. **🔄 尝试不同角度**:
   - 不同的库/框架特性？
   - 不同的架构模式？
   - 移除抽象而不是添加？

<!-- /AI_WORKFLOW -->

## 8. 技术标准

<!-- AI_TECHNICAL_STANDARDS -->

### 架构原则

- **组合胜过继承** - 使用依赖注入
- **接口胜过单例** - 启用测试和灵活性
- **显式胜过隐式** - 清晰的数据流和依赖关系
- **尽可能测试驱动** - 永远不要禁用测试，而是修复它们

### 代码质量

**每次提交必须**:
- 编译成功
- 通过所有现有测试
- 包含新功能的测试
- 遵循项目格式化/linting

**提交前**:
- 运行格式化工具/linters
- 自我审查变更
- 确保提交消息解释"为什么"

### 错误处理

- 快速失败并提供描述性消息
- 包含调试上下文
- 在适当级别处理错误
- 永远不要静默吞噬异常

<!-- /AI_TECHNICAL_STANDARDS -->

## 9. 决策框架

<!-- AI_DECISION_FRAMEWORK -->

当存在多个有效方法时，基于以下标准选择：

1. **可测试性** - 我能轻松测试这个吗？
2. **可读性** - 6个月后有人能理解这个吗？
3. **一致性** - 这符合项目模式吗？
4. **简洁性** - 这是有效的最简单解决方案吗？
5. **可逆性** - 以后更改有多困难？

<!-- /AI_DECISION_FRAMEWORK -->

## 10. 项目集成

<!-- AI_PROJECT_INTEGRATION -->

### 学习代码库

- 找到 3 个类似的功能/组件
- 识别常见模式和约定
- 尽可能使用相同的库/工具
- 遵循现有测试模式

### 工具使用

- 使用项目现有的构建系统
- 使用项目的测试框架
- 使用项目的格式化工具/linter 设置
- 没有强有力的理由不要引入新工具

<!-- /AI_PROJECT_INTEGRATION -->

## 11. 质量门控

<!-- AI_QUALITY_GATES -->

### 完成定义 (Definition of Done)

- [ ] 测试已编写并通过
- [ ] 代码遵循项目约定
- [ ] 无 linter/formatter 警告
- [ ] 提交消息清晰
- [ ] 实现符合计划
- [ ] 无 TODOs（除非有问题编号）

### 测试指南

- 测试行为，而非实现
- 尽可能每个测试一个断言
- 清晰的测试名称描述场景
- 使用现有测试工具/助手
- 测试应该是确定性的

<!-- /AI_QUALITY_GATES -->

## 12. 重要提醒

<!-- AI_IMPORTANT_REMINDERS -->

**🚫 永远不要**:
- 使用 `--no-verify` 绕过提交钩子
- 禁用测试而不是修复它们
- 提交不能编译的代码
- 做假设 - 用现有代码验证

**✅ 总是**:
- 增量提交工作代码
- 随时更新计划文档
- 从现有实现中学习
- 3次失败尝试后停止并重新评估

<!-- /AI_IMPORTANT_REMINDERS -->

## 13. AI 代码生成特别指导

<!-- AI_SPECIFIC_GUIDANCE: CRITICAL -->

### 错误处理规范

**🚫 禁止修改现有错误处理逻辑**:
- 不得改变现有错误返回的条件、类型、顺序和内容
- 特别是 `github.com/biuren/hebe/internal/utils/error` 包定义的 `*errors.AppError` 类型
- 保持 `Code` 和 `Key` 字段的一致性

### 代码生成最佳实践

#### 1. 上下文感知
```go
// ✅ 正确：了解项目结构后生成
func (s *UserService) CreateUser(ctx context.Context, req *CreateUserRequest) (*User, error) {
    // 使用项目现有的验证器
    if err := s.validator.Validate(req); err != nil {
        return nil, errors.NewValidationError("invalid_request", err.Error())
    }
    // ...
}

// ❌ 错误：忽略项目约定
func CreateUser(request map[string]interface{}) interface{} {
    // 不符合项目风格
}
```

#### 2. 依赖注入模式
```go
// ✅ 正确：使用依赖注入
type UserService struct {
    repo UserRepository
    logger Logger
}

func NewUserService(repo UserRepository, logger Logger) *UserService {
    return &UserService{repo: repo, logger: logger}
}

// ❌ 错误：使用全局变量或单例
var globalDB *sql.DB
func CreateUser() { /* 使用 globalDB */ }
```

#### 3. 测试友好设计
```go
// ✅ 正确：接口抽象，便于测试
type UserRepository interface {
    Create(ctx context.Context, user *User) error
    FindByID(ctx context.Context, id string) (*User, error)
}

// ❌ 错误：直接依赖具体实现
type UserService struct {
    db *sql.DB // 难以测试
}
```

### 重要提醒

<!-- AI_REMINDERS: CRITICAL -->
**🔥 关键注意事项**:

1. **安全第一**: 永远不要在代码中暴露或记录密钥和敏感信息
2. **向后兼容**: 保持现有 API 和接口的兼容性
3. **测试覆盖**: 为所有新功能编写相应的测试
4. **文档同步**: 更新相关的文档和注释
5. **性能考虑**: 避免引入性能瓶颈
6. **错误处理**: 提供清晰、有用的错误信息
7. **国际化**: 所有用户可见文本都要支持多语言
8. **代码风格**: 严格遵循项目的编码规范
<!-- /AI_REMINDERS -->

---

<!-- AI_FOOTER -->
**📝 文档版本**: v2.0  
**🔄 最后更新**: 2024年  
**🎯 适用范围**: iDiary.ai Hebe 项目所有 AI 代码生成任务  
**⚡ 优先级**: 本文档中的所有 `MANDATORY`、`CRITICAL` 标记的规则具有最高优先级
<!-- /AI_FOOTER -->