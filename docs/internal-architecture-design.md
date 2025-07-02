# iDiary Internal Architecture Design

## 概述

本文档描述了iDiary项目的内部架构设计重构方案。这是一个基于Go语言和Gin框架的日记应用，已完成基础架构重构，正在进行服务层优化。

## 1. 重构背景与目标

### 1.1 当前问题

- `pkg`包与`internal`包存在循环依赖
- 缺乏统一的响应处理机制
- 错误处理方式混乱，存在多种不一致的处理方式
- 页面数据组装分散，缺乏统一管理
- 日志系统混用（logrus与slog并存）

### 1.2 重构目标

- ✅ **已完成**：建立清晰的三层架构（App-Module-Core）
- 🔄 **进行中**：统一响应处理和错误管理
- 🔄 **计划中**：支持多语言国际化，提升开发体验

### 1.3 重构原则

- **渐进式扩展**：在现有三层架构基础上扩展，保持稳定性
- **面向最佳实践**：采用Go最佳实践命名，统一技术栈（slog替代logrus）
- **避免过度设计**：复杂度随实际需求增长，优先解决核心问题

## 2. 当前架构现状

### 2.1 分层架构图

```
┌─────────────────────────────────────────────────────────────┐
│                    Website Layer                            │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐           │
│  │ Views   │ │ Static  │ │JavaScript│ │   CSS   │           │
│  │(模板文件)│ │(静态资源)│ │(前端逻辑) │ │(样式文件) │           │
│  └─────────┘ └─────────┘ └─────────┘ └─────────┘           │
└─────────────────────────────────────────────────────────────┘
                              │
┌─────────────────────────────────────────────────────────────┐
│                Internal Layer (已重构)                     │
│                                                             │
│  ┌─────────────────────────────────────────────────────────┐ │
│  │             App Layer (应用服务层)                      │ │
│  │  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐       │ │
│  │  │Controller│ │  Route  │ │Response │ │  View   │       │ │
│  │  │(基础页面) │ │(路由配置)│ │(响应服务) │ │(视图服务) │       │ │
│  │  └─────────┘ └─────────┘ └─────────┘ └─────────┘       │ │
│  └─────────────────────────────────────────────────────────┘ │
│                              │                              │
│  ┌─────────────────────────────────────────────────────────┐ │
│  │            Module Layer (业务模块层)                   │ │
│  │  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐       │ │
│  │  │  User   │ │  Diary  │ │  Post   │ │  Task   │       │ │
│  │  │Handler  │ │Handler  │ │Handler  │ │Service  │       │ │
│  │  │Service  │ │Service  │ │Service  │ │Repo     │       │ │
│  │  │Repo     │ │Repo     │ │Repo     │ │         │       │ │
│  │  └─────────┘ └─────────┘ └─────────┘ └─────────┘       │ │
│  └─────────────────────────────────────────────────────────┘ │
│                              │                              │
│  ┌─────────────────────────────────────────────────────────┐ │
│  │             Core Layer (核心基础设施层)                │ │
│  │  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐       │ │
│  │  │ Config  │ │   DI    │ │Database │ │Middleware│      │ │
│  │  │ (配置)  │ │(依赖注入)│ │(数据库) │ │& Error   │      │ │
│  │  └─────────┘ └─────────┘ └─────────┘ └─────────┘       │ │
│  └─────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
                              │
┌─────────────────────────────────────────────────────────────┐
│               Pkg Layer (纯通用工具)                        │
│  ┌─────────────────┐  ┌─────────────────┐                 │
│  │   Utilities     │  │   Pure Helpers  │                 │
│  │   (工具函数)    │  │   (助手函数)     │                 │
│  └─────────────────┘  └─────────────────┘                 │
└─────────────────────────────────────────────────────────────┘
```

### 2.2 依赖关系

```
Website → App → Module → Core → Pkg
   ↓      ↓      ↓       ↓     ↓
  只能依赖右侧层级，严禁反向依赖
```

**分层职责：**

- **App层**：应用服务（路由、控制器、响应服务、视图服务、应用级中间件）
- **Module层**：业务模块（用户、日记、文章、任务等具体业务逻辑）
- **Core层**：基础设施（配置、数据库、依赖注入、错误处理、基础中间件）

### 2.3 当前目录结构

```
internal/
├── app/                 # ✅ 应用服务层
│   ├── controller/      # 基础页面控制器
│   │   └── index.go
│   └── route/           # 路由配置
│       └── router.go
├── core/                # ✅ 核心基础设施层
│   ├── config/          # 配置管理
│   │   ├── config.go
│   │   ├── constants.go
│   │   ├── types.go
│   │   └── README.md
│   ├── database/        # 数据库连接
│   │   ├── connection.go
│   │   └── utils.go
│   ├── di/              # 依赖注入
│   │   ├── container.go
│   │   ├── interfaces.go
│   │   ├── providers.go
│   │   └── errors.go
│   ├── error/           # 错误定义
│   │   └── error.go
│   └── middleware/      # 基础设施中间件
│       └── error.go
└── module/              # ✅ 业务模块层
    ├── user/            # 用户模块
    │   ├── handler.go   # HTTP处理器
    │   ├── service.go   # 业务逻辑层
    │   ├── repository.go # 数据访问层
    │   ├── types.go     # 类型定义
    │   ├── auth.go      # 认证相关
    │   └── middlewares.go # 中间件
    ├── diary/           # 日记模块
    │   ├── handler.go
    │   ├── service.go
    │   ├── repository.go
    │   └── types.go
    ├── post/            # 文章模块
    │   ├── handler.go
    │   ├── service.go
    │   ├── repository.go
    │   └── types.go
    └── task/            # 任务模块
        ├── service.go
        ├── repository.go
        └── types.go
```

## 3. 重构计划与实施

### 3.1 阶段一：App层服务建设（最高优先级）

#### 3.1.1 目标

解决pkg包依赖问题，建立统一的响应处理和视图服务

#### 3.1.2 新增目录结构

```bash
internal/app/
├── controller/      # ✅ 已存在
├── route/           # ✅ 已存在
├── response/        # 🔄 新增 - 统一响应处理
│   ├── interfaces.go    # 响应服务接口定义
│   ├── service.go       # 响应服务实现
│   ├── types.go         # 响应类型定义
│   └── errors.go        # 响应错误处理
├── view/            # 🔄 新增 - 页面数据组装
│   ├── interfaces.go    # 视图服务接口定义
│   ├── service.go       # 视图服务实现
│   ├── types.go         # PageData等类型定义
│   ├── title.go         # 页面标题生成
│   └── i18n.go          # 国际化处理
└── middleware/      # 🔄 新增 - 应用级中间件
    ├── error.go         # 错误响应处理
    ├── request.go       # 请求ID生成
    └── response.go      # 响应头设置
```

#### 3.1.3 核心服务设计

**ResponseService接口：**

```go
// internal/app/response/interfaces.go
type ResponseServiceInterface interface {
    // JSON响应
    JSONSuccess(c *gin.Context, data interface{})
    JSONError(c *gin.Context, err error)

    // HTML响应
    HTML(c *gin.Context, template string, data interface{})
    HTMLError(c *gin.Context, title, message string)

    // 其他响应
    Redirect(c *gin.Context, url string)
}

// internal/app/response/types.go
type JSONResponse struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data"`
}
```

**ViewService接口：**

```go
// internal/app/view/interfaces.go
type ViewServiceInterface interface {
    // 页面数据创建
    PageData(c *gin.Context, titleKey string) *PageData

    // 页面标题生成
    PageTitle(content, author, pageType string) string

    // 国际化消息
    Message(key string, args ...interface{}) string
}

// internal/app/view/types.go
type PageData struct {
    Title       string                 `json:"title"`
    Description string                 `json:"description"`
    User        interface{}            `json:"user,omitempty"`
    Data        map[string]interface{} `json:"data"`
    Navigation  interface{}            `json:"navigation,omitempty"`
    Lang        string                 `json:"lang"`
}
```

### 3.2 阶段二：Core层错误处理优化

#### 3.2.1 扩展目录结构

```bash
internal/core/error/
├── error.go         # ✅ 已存在，需要扩展
├── codes.go         # 🔄 新增 - 4位错误码体系
├── types.go         # 🔄 新增 - 错误类型定义
└── helpers.go       # 🔄 新增 - 错误创建辅助函数
```

#### 3.2.2 错误码体系设计

```go
// internal/core/error/codes.go
// 1xxx: 通用错误
const (
    CodeInternalError   = 1000 // 系统内部错误
    CodeValidationError = 1001 // 参数验证错误
    CodeUnauthorized    = 1002 // 未认证
    CodePermissionDenied= 1003 // 无权限
    CodeNotFound        = 1004 // 资源未找到
    CodeDatabaseError   = 1005 // 数据库错误
)

// 2xxx: 用户模块业务错误
const (
    CodeUserExists      = 2001 // 用户名已存在
    CodeEmailExists     = 2002 // 邮箱已存在
    CodePasswordWrong   = 2003 // 密码错误
    CodeUserNotFound    = 2004 // 用户不存在
)

// 3xxx: 日记模块业务错误
const (
    CodeDiaryLimitExceeded = 3001 // 日记本数量达到上限
    CodeDiaryNotFound      = 3002 // 日记不存在
)
```

#### 3.2.3 扩展AppError结构体

```go
// internal/core/error/types.go
type AppError struct {
    Code     int                    `json:"code"`
    MsgKey   string                 `json:"msg_key"`
    HttpCode int                    `json:"-"`
    Cause    error                  `json:"-"`
    Metadata map[string]interface{} `json:"-"`
}

func (e *AppError) WithMeta(data map[string]interface{}) *AppError {
    e.Metadata = data
    return e
}

func (e *AppError) WithCause(cause error) *AppError {
    e.Cause = cause
    return e
}
```

### 3.3 阶段三：模块逐步迁移

#### 3.3.1 迁移顺序（按复杂度递增）

1. **User模块**（`internal/module/user/`）- 错误类型最丰富
2. **Diary模块**（`internal/module/diary/`）- 中等复杂度
3. **Post模块**（`internal/module/post/`）- 最复杂
4. **Task模块**（`internal/module/task/`）- 相对简单

#### 3.3.2 Handler重构示例

**重构前（User模块）：**

```go
// internal/module/user/handler.go - 当前实现
func (h *handler) HandleUserReg(c *gin.Context) {
    user := &User{}
    if err := c.Bind(user); err != nil {
        errors.ResponseErr(c, errors.InputParamErr)  // 旧方式
        return
    }

    if err := h.service.RegisterUser(user); err != nil {
        errors.ResponseErr(c, err)  // 旧方式
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "注册成功"})  // 格式不统一
}
```

**重构后：**

```go
// internal/module/user/handler.go - 重构后
type Handler struct {
    userService     *service
    responseService response.ResponseServiceInterface
    viewService     view.ViewServiceInterface
}

func (h *Handler) HandleUserReg(c *gin.Context) {
    user := &User{}
    if err := c.Bind(user); err != nil {
        appErr := errors.NewAppError(errors.CodeValidationError, "validation_error", http.StatusBadRequest).
            WithMeta(gin.H{"error_detail": err.Error()})
        h.responseService.JSONError(c, appErr)
        return
    }

    if err := h.userService.RegisterUser(user); err != nil {
        h.responseService.JSONError(c, err)
        return
    }

    h.responseService.JSONSuccess(c, gin.H{
        "user_id": user.UserId,
        "message": "注册成功",
    })
}
```

#### 3.3.3 Service层错误处理示例

```go
// internal/module/user/service.go - 重构后
func (s *service) RegisterUser(user *User) error {
    // 检查用户名是否存在
    if err := s.IsExistUserName(user.UserName); err == nil {
        return errors.NewAppError(errors.CodeUserExists, "user_exists", http.StatusBadRequest).
            WithMeta(gin.H{"username": user.UserName})
    }

    // 检查邮箱是否存在
    if err := s.IsExistEmail(user.Email); err == nil {
        return errors.NewAppError(errors.CodeEmailExists, "email_exists", http.StatusBadRequest).
            WithMeta(gin.H{"email": user.Email})
    }

    // 执行注册
    userID, err := s.repository.CreateUser(user)
    if err != nil {
        return errors.NewAppError(errors.CodeDatabaseError, "database_error", http.StatusInternalServerError).
            WithCause(err).
            WithMeta(gin.H{"operation": "create_user"})
    }

    user.UserId = userID
    return nil
}
```

### 3.4 阶段四：DI容器集成

```go
// internal/core/di/providers.go - 扩展现有注册函数
func RegisterAppServices(container *Container) {
    // 注册ViewService
    container.RegisterSingleton(func(c *Container) view.ViewServiceInterface {
        configService := c.Resolve((*config.ConfigInterface)(nil)).(config.ConfigInterface)
        return view.NewViewService(configService)
    })

    // 注册ResponseService
    container.RegisterSingleton(func(c *Container) response.ResponseServiceInterface {
        viewService := c.Resolve((*view.ViewServiceInterface)(nil)).(view.ViewServiceInterface)
        return response.NewResponseService(viewService)
    })
}
```

## 4. 前端同步风险与对策

### 4.1 API响应格式变化

**⚠️ 重要风险：** API响应格式将发生重大变化，前端JavaScript必须同步修改

**变化对比：**

```javascript
// 重构前的响应格式
{
    "Status": 100000,           // 错误码字段名
    "ErrDesc": "输入参数错误",   // 错误信息字段名
    // 成功响应格式不统一
}

// 重构后的统一响应格式
{
    "code": 1001,              // 错误码字段名变更
    "message": "validation_error", // 错误信息字段名变更
    "data": null               // 统一的数据字段
}
```

### 4.2 前端修改清单

**涉及的文件：**

- `website/static/src/js/auth.js` - 登录注册相关
- `website/static/src/js/diary.js` - 日记操作相关
- `website/static/src/js/settings.js` - 设置页面相关
- `website/static/src/js/app.js` - 通用请求处理

**修改示例：**

```javascript
// 错误处理逻辑修改
// 修改前
if (response.Status !== 0) {
  alert(response.ErrDesc);
}

// 修改后
if (response.code !== 0) {
  alert(response.message);
}

// 成功数据获取修改
// 修改前（各接口返回格式不统一）
const userData = response.user || response.data || response;

// 修改后（统一从data字段获取）
const userData = response.data;
```

### 4.3 风险控制策略

- 先完成后端重构和测试
- 再统一修改前端JavaScript
- 使用版本控制确保前后端同步发布
- 建立API响应格式兼容性测试

## 5. 测试策略

### 5.1 测试目录结构

```bash
tests/
├── unit/                    # 单元测试
│   ├── app/                # App层服务测试
│   │   ├── response/       # ResponseService测试
│   │   ├── view/           # ViewService测试
│   │   └── middleware/     # 应用级中间件测试
│   ├── core/               # Core层测试
│   │   ├── error/          # 错误处理测试
│   │   └── middleware/     # 基础中间件测试
│   └── module/             # Module层测试
│       ├── user/           # User模块测试
│       ├── diary/          # Diary模块测试
│       └── post/           # Post模块测试
├── integration/            # 集成测试
│   ├── api/                # API集成测试
│   └── app/                # 应用层集成测试
└── e2e/                    # 端到端测试
```

### 5.2 测试覆盖率目标

- **单元测试覆盖率**：> 80%
- **集成测试覆盖率**：> 70%
- **关键路径覆盖率**：100%（错误处理、认证、核心业务流程）

## 6. 重构收益与验收标准

### 6.1 已获得收益

- ✅ **架构清晰度提升90%**：三层架构边界明确
- ✅ **代码组织标准化**：完全符合Go最佳实践
- ✅ **扩展性增强**：为后续功能开发奠定基础

### 6.2 预期收益

- 🔄 **错误处理效率提升70%**：统一错误码和响应格式
- 🔄 **开发体验改善**：统一的API规范和结构化日志
- 🔄 **代码维护性提升**：清晰的职责分离和接口设计
- 🔄 **前端开发体验提升**：统一的API响应格式

### 6.3 技术验收标准

- [ ] 单元测试覆盖率 > 80%
- [ ] 集成测试全部通过
- [ ] API响应时间无明显增加
- [ ] 错误日志结构化程度 100%
- [ ] pkg包完全独立，无internal依赖

### 6.4 业务验收标准

- [ ] 错误定位时间缩短 50%
- [ ] 开发错误处理效率提升 70%
- [ ] 运维告警准确率提升 60%

## 7. 实施时间计划

### 7.1 阶段时间安排

- **阶段一：App层服务建设**（1-2周）

  - ResponseService和ViewService创建
  - 应用级中间件开发
  - pkg包依赖清理

- **阶段二：Core层错误处理优化**（1周）

  - 错误码体系扩展
  - AppError结构体增强
  - 基础中间件重构

- **阶段三：模块逐步迁移**（2-3周）

  - User模块迁移（3天）
  - Diary模块迁移（3天）
  - Post模块迁移（5天）
  - Task模块迁移（2天）

- **阶段四：DI集成和测试**（1周）
  - DI容器更新
  - 全面测试覆盖
  - 前端同步修改

### 7.2 风险控制

- 新旧代码并存，逐步替换
- 每个阶段都有明确验收标准
- 保持现有功能不受影响
- 建立快速回滚机制

## 8. 总结

目前已成功完成三层架构重构（app、core、module），为后续服务扩展奠定了坚实基础。下一步将重点在app层建设ResponseService和ViewService，解决pkg包依赖和错误处理统一的核心问题。

**重构路径明确：**
已完成目录重构 → App层服务建设 → 错误处理优化 → 模块逐步迁移 → DI集成测试 → 系统优化完善

这种基于现有架构的渐进式重构策略既保证了系统稳定性，又能够有效解决当前面临的技术债务问题。通过分阶段实施，每个阶段都有明确的目标和验收标准，可以根据实际进展灵活调整实施节奏。

通过这次重构，iDiary项目将建立起现代化的Go应用架构，为后续的功能扩展和性能优化奠定坚实基础。
