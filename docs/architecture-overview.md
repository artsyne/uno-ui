# iDiary.ai Hebe 系统架构概览

## 概述

iDiary.ai Hebe是一个基于Go 1.25和Gin框架的现代化智能日记管理系统，采用DDD（领域驱动设计）的分层架构，通过Go 1.25泛型依赖注入容器实现松耦合设计，具备高性能、高可维护性和可扩展性。系统支持多种内容类型（文字、图片、视频、链接）和多云盘存储。

## 分层架构

```
┌─────────────────────────────────────────────────────────────┐
│                    Website Layer                            │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐           │
│  │Templates│ │ Static  │ │ Kiko.js │ │Tailwind │           │
│  │(模板文件)│ │(静态资源)│ │(前端框架)│ │CSS 4.x  │           │
│  └─────────┘ └─────────┘ └─────────┘ └─────────┘           │
└─────────────────────────────────────────────────────────────┘
                              │
┌─────────────────────────────────────────────────────────────┐
│                App Layer (应用服务层)                      │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐           │
│  │Startup  │ │ Route   │ │  View   │ │Controller│           │
│  │(启动器) │ │(路由)   │ │(视图)   │ │(控制器) │           │
│  └─────────┘ └─────────┘ └─────────┘ └─────────┘           │
│  ┌─────────┐ ┌─────────┐                                   │
│  │Middleware│ │Templates│                                   │
│  │(中间件) │ │(模板)   │                                   │
│  └─────────┘ └─────────┘                                   │
└─────────────────────────────────────────────────────────────┘
                              │
┌─────────────────────────────────────────────────────────────┐
│                Domain Layer (业务领域层)                   │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐           │
│  │  User   │ │  Diary  │ │  Post   │ │  Task   │           │
│  │(用户领域)│ │(日记领域)│ │(内容领域)│ │(任务领域)│           │
│  │R-S-T    │ │R-S-T    │ │R-S-T    │ │Service  │           │
│  └─────────┘ └─────────┘ └─────────┘ └─────────┘           │
└─────────────────────────────────────────────────────────────┘
                              │
┌─────────────────────────────────────────────────────────────┐
│                Infra Layer (基础设施层)                    │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐                       │
│  │Repository│ │Database │ │ Cache   │                       │
│  │(仓储实现)│ │(数据库) │ │(缓存)   │                       │
│  └─────────┘ └─────────┘ └─────────┘                       │
└─────────────────────────────────────────────────────────────┘
                              │
┌─────────────────────────────────────────────────────────────┐
│                Core Layer (核心基础层)                     │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐           │
│  │   DI    │ │Database │ │Initializer│ │ Config │           │
│  │(依赖注入)│ │(数据库) │ │(初始化) │ │(配置)   │           │
│  └─────────┘ └─────────┘ └─────────┘ └─────────┘           │
│       │           │           │           │                │
│       │           │           └───────────┘                │
│       │           │              配置注册到DI容器           │
│       └───────────┴─────────────────┐                      │
│              统一依赖管理             │                      │
└─────────────────────────────────────────────────────────────┘
                              │
┌─────────────────────────────────────────────────────────────┐
│            Utils & Contracts Layer (工具和契约层)          │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐           │
│  │ Utils   │ │ Error   │ │  I18n   │ │Contracts│           │
│  │(工具)   │ │(错误)   │ │(国际化) │ │(契约)   │           │
│  └─────────┘ └─────────┘ └─────────┘ └─────────┘           │
│  ┌─────────┐ ┌─────────┐                                   │
│  │Middleware│ │ Config  │                                   │
│  │(中间件) │ │(配置)   │                                   │
│  └─────────┘ └─────────┘                                   │
└─────────────────────────────────────────────────────────────┘
```

**依赖方向**: Website → App → Domain → Infra → Core → Utils/Contracts (严禁反向依赖)

## 分层职责详解

### Website Layer (网站层)

**核心职责**: 前端展示和用户交互

- **Templates**: HTML模板文件，使用Go template语法
- **Static**: 静态资源文件（CSS、JS、图片等）
- **Kiko.js**: 自研前端框架，提供组件化和国际化支持
- **Tailwind CSS 4.x**: 现代化CSS框架，提供响应式设计

### App Layer (应用服务层)

**核心职责**: 应用编排、HTTP处理、渲染服务

#### Startup 模块

- **职责**: 应用生命周期管理
- **位置**: `internal/app/startup/`
- **功能**:
  - 应用启动和关闭流程编排
  - 组件初始化顺序管理
  - 错误处理和恢复机制

#### Route 模块

- **职责**: HTTP路由配置和管理
- **位置**: `internal/app/route/`
- **功能**:
  - 路由注册和分组
  - 中间件配置
  - 静态资源路由

#### View 模块

- **职责**: 统一视图服务（整合视图渲染、数据绑定和模板函数）
- **位置**: `internal/app/view/`
- **功能**:
  - 视图渲染和数据绑定
  - 模板函数管理
  - 国际化支持
  - 错误页面渲染

#### Controller 模块

- **职责**: 基础页面控制器
- **位置**: `internal/app/controller/`
- **功能**:
  - 首页等基础页面处理
  - 静态页面渲染

#### Middleware 模块

- **职责**: 应用级中间件
- **位置**: `internal/app/middleware/`
- **功能**:
  - 认证中间件
  - 权限检查中间件
  - 业务逻辑中间件

#### Templates 模块

- **职责**: 模板管理和渲染
- **位置**: `internal/app/templates/`
- **功能**:
  - 模板文件管理
  - 模板函数注册
  - 模板缓存机制

### Domain Layer (业务领域层)

**核心职责**: 业务逻辑实现，采用Repository-Service-Types模式

#### User 模块

- **职责**: 用户管理和认证
- **位置**: `internal/domain/user/`
- **组件**:
  - **Repository**: 用户数据访问接口定义
  - **Service**: 用户业务逻辑层 (用户认证、密码加密等)
  - **Types**: 用户相关数据结构定义

#### Diary 模块

- **职责**: 日记本管理
- **位置**: `internal/domain/diary/`
- **组件**:
  - **Repository**: 日记本数据访问接口定义
  - **Service**: 日记本业务逻辑
  - **Types**: 日记本相关数据结构

#### Post 模块

- **职责**: 日记内容管理
- **位置**: `internal/domain/post/`
- **组件**:
  - **Repository**: 内容数据访问接口定义
  - **Service**: 内容业务逻辑
  - **Types**: 内容相关数据结构

#### Task 模块

- **职责**: 后台任务处理和异步任务管理
- **位置**: `internal/domain/task/`
- **组件**:
  - **Repository**: 任务数据访问接口定义
  - **Service**: 定时任务服务（基于gocron的邮件任务调度）
  - **Types**: 任务相关数据结构（JobQueue、EmailTask等）

### Infra Layer (基础设施层)

**核心职责**: 数据访问和外部服务集成

#### Repository 模块

- **职责**: 数据访问层实现
- **位置**: `internal/infra/repository/`
- **功能**:
  - 实现Domain层定义的Repository接口
  - 数据库操作封装
  - 缓存策略实现

### Core Layer (核心基础层)

**核心职责**: 核心基础设施服务，为上层提供技术支撑

#### DI 模块

- **职责**: Go 1.25泛型依赖注入容器
- **位置**: `internal/core/di/`
- **功能**:
  - 零反射泛型设计，完全消除反射开销
  - 三种生命周期管理 (瞬态Bind/单例Singleton/实例Value)
  - 编译时类型安全，使用泛型确保类型检查
  - 高并发性能，基于读写锁和sync.Once优化
  - 重复注册检测，防止意外覆盖服务
  - 明确覆盖支持，提供Rebind方法用于测试场景

#### Database 模块

- **职责**: 数据库连接和管理
- **位置**: `internal/core/database/`
- **功能**:
  - 数据库连接池管理
  - SQL执行辅助函数
  - 事务管理
  - SQLx集成

#### Initializer 模块

- **职责**: 系统初始化管理
- **位置**: `internal/core/initializer/`
- **功能**:
  - 组件初始化顺序管理
  - 依赖关系解析
  - 启动流程编排
  - **配置系统集成**: 将配置对象注册到DI容器中管理

#### 配置系统与DI容器集成

项目采用混合架构模式管理配置，同时支持全局配置访问和DI容器注入：

**配置生命周期管理**:
1. **初始化阶段**: 配置系统初始化，创建全局配置实例
2. **注册阶段**: DI容器初始化，将配置注册为预创建实例
3. **注入阶段**: 应用层服务通过构造函数注入配置依赖

**分层使用策略**:
- **应用层/领域层**: 优先使用DI注入配置，提升测试性和松耦合性
- **基础设施层**: 保留全局配置访问以保证向后兼容性
- **模板/工具层**: 使用全局配置函数提供便利性

### Utils & Contracts Layer (工具和契约层)

**核心职责**: 通用工具和系统契约定义

#### Utils 模块

- **职责**: 通用工具函数
- **位置**: `internal/utils/`
- **子模块**:
  - **Config**: 配置管理和解析
  - **其他工具**: 文本处理、时间处理等

#### Error 模块

- **职责**: 错误定义和处理
- **位置**: `internal/utils/error/`
- **功能**:
  - 统一错误码体系 (4位错误码)
  - AppError结构定义
  - 错误分类和处理策略

#### I18n 模块

- **职责**: 国际化支持
- **位置**: `internal/utils/i18n/`
- **功能**:
  - 多语言消息管理
  - 语言切换机制
  - 消息格式化

#### Contracts 模块

- **职责**: 系统契约和接口定义
- **位置**: `internal/contracts/`
- **功能**:
  - 共享接口定义
  - 类型安全保证
  - 系统契约规范

#### Middleware 模块

- **职责**: 通用中间件
- **位置**: `internal/utils/middleware/`
- **功能**:
  - Panic恢复中间件
  - 错误处理中间件
  - 结构化日志记录

#### Config 模块

- **职责**: 配置管理
- **位置**: `internal/utils/config/`
- **功能**:
  - 配置文件加载和解析
  - 环境变量处理
  - 配置验证
  - Viper集成

## 数据流向

### 请求处理流程

```
HTTP Request
    ↓
Middleware (Panic Recovery, Error Handler, Auth, etc.)
    ↓
Route (路由匹配)
    ↓
Controller/Handler (HTTP处理层)
    ↓
Domain Service (业务逻辑层)
    ↓
Domain Repository Interface (领域仓储接口)
    ↓
Infra Repository Implementation (基础设施仓储实现)
    ↓
Database (数据库)
```

### 响应处理流程

```
Database
    ↓
Infra Repository (数据组装)
    ↓
Domain Service (业务逻辑处理)
    ↓
Controller (调用ViewService)
    ↓
ViewService (页面数据组装、模板渲染)
    ↓
HTTP Response (JSON/HTML)
```

### 依赖注入流程

```
Application Startup
    ↓
DI Container Initialization
    ↓
Service Registration (按依赖顺序)
    ↓
Service Resolution (泛型类型安全)
    ↓
Component Injection
    ↓
Application Ready
```

## 关键设计模式

### 1. Go 1.25泛型依赖注入容器

```go
// 创建容器
container := di.New()

// 瞬态服务注册 - 每次获取都创建新实例
err := di.Bind(container, func() *user.Service {
    return user.NewService(config, userRepo, taskService)
})

// 单例服务注册 - 全局唯一实例，延迟初始化
err := di.Singleton(container, func() *view.ViewService {
    return view.NewViewService(config)
})

// 实例服务注册 - 预创建实例
err := di.Value(container, config)

// 泛型服务解析 - 零反射，高性能
userService := di.Get[*user.Service](container)

// 安全获取服务，返回错误而非panic
viewService, err := di.TryGet[*view.ViewService](container)
```

### 2. 仓储模式 (Repository Pattern)

```go
// Domain层接口定义
type UserRepository interface {
    GetByID(ctx context.Context, id int) (*User, error)
    Create(ctx context.Context, user *User) error
}

// Infra层实现
type userRepositoryImpl struct {
    db *sqlx.DB
}

func (r *userRepositoryImpl) GetByID(ctx context.Context, id int) (*User, error) {
    // 数据库操作实现
}
```

### 3. 领域服务模式 (Domain Service)

```go
// 用户服务实现
type Service struct {
    config      *config.Config
    repo        RepositoryInterface
    regexCheck  *pkg.RegexCheck
    taskService *task.Service
}

func NewService(config *config.Config, repo RepositoryInterface, taskService *task.Service) *Service {
    return &Service{
        config:      config,
        repo:        repo,
        regexCheck:  &pkg.RegexCheck{},
        taskService: taskService,
    }
}

func (s *Service) LoginWithUsernameOrEmail(usernameOrEmail, password string) (*User, error) {
    var user *User
    var err error

    // 判断输入是否为邮箱地址
    if s.regexCheck.IsEmail(usernameOrEmail) {
        user, err = s.repo.GetUserByEmail(usernameOrEmail)
    } else {
        user, err = s.repo.GetUserByUserName(usernameOrEmail)
    }

    if err != nil {
        return nil, err
    }

    // 验证密码
    if !s.verifyPassword(password, user.Pwd) {
        return nil, appErrors.PasswordIncorrect.WithDetails("密码不正确")
    }

    return user, nil
}
```

### 4. 统一错误处理模式

```go
// 错误定义结构体
type ErrorDef struct {
    Code     int    `json:"code"`      // 错误数字码
    Key      string `json:"key"`       // 错误标识符
    HttpCode int    `json:"http_code"` // HTTP状态码
}

// 应用错误
type AppError struct {
    ErrorDef        // 嵌入错误定义
    Details  string `json:"-"` // 错误详情
}

// 预定义错误
var (
    UserNotFound     = ErrorDef{1301, "user_not_found", 404}
    PasswordIncorrect = ErrorDef{2201, "password_incorrect", 400}
    DatabaseError    = ErrorDef{6001, "database_error", 500}
)

// 创建错误
func (def ErrorDef) WithDetails(details string) *AppError {
    return &AppError{
        ErrorDef: def,
        Details:  details,
    }
}

// 错误处理中间件
func AppErrorHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        defer func() {
            if len(c.Errors) > 0 {
                lastError := c.Errors.Last()
                if appErr, ok := lastError.Err.(*AppError); ok {
                    c.JSON(appErr.HttpCode, gin.H{
                        "code":    appErr.Code,
                        "message": appErr.Key,
                        "data":    nil,
                    })
                }
            }
        }()
        c.Next()
    }
}
```

### 5. 视图服务模式 (View Service)

```go
// ViewService 视图服务实现
type ViewService struct {
    config      *config.Config
    i18nManager i18n.I18nManager
}

func NewViewService(cfg *config.Config) *ViewService {
    i18nManager := i18n.NewLocalizationManager(cfg)
    return &ViewService{
        config:      cfg,
        i18nManager: i18nManager,
    }
}

// 创建页面数据（泛型支持）
func (s *ViewService) CreatePageData(c *gin.Context, titleKey string) *contracts.PageData[map[string]interface{}] {
    ctx := s.BuildRenderContext(c)
    title := s.generateTitle(c, titleKey)

    return &contracts.PageData[map[string]interface{}]{
        Title:   title,
        Context: ctx,
        Data:    make(map[string]interface{}),
    }
}

// JSON响应
func (s *ViewService) JSONSuccess(c *gin.Context, data interface{}) {
    response := JSONResponse{
        Success: true,
        Message: "success",
        Code:    0,
        Data:    data,
    }
    c.JSON(http.StatusOK, response)
}

// HTML渲染
func (s *ViewService) HTML(c *gin.Context, name string, data interface{}) {
    c.HTML(http.StatusOK, name, data)
}
```

### 6. 中间件链模式

```go
// 在 initGinServer 中设置中间件
func (a *Application) initGinServer() error {
    gin.SetMode(a.config.App.GinMode)
    server := gin.New()

    // 初始化中间件链
    server.Use(utilsMiddleware.PanicRecovery())   // panic recovery中间件
    server.Use(utilsMiddleware.AppErrorHandler()) // 应用错误处理中间件
    server.Use(middleware.SharedData())           // 用户共享数据中间件

    // 设置路由
    route.SetWithContainer(server, a.container)

    a.server = server
    return nil
}

// PanicRecovery中间件实现
func PanicRecovery() gin.HandlerFunc {
    return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
        if recovered != nil {
            slog.Error("panic recovered",
                "panic", recovered,
                "path", c.Request.URL.Path,
                "method", c.Request.Method,
                "stack", string(debug.Stack()))

            internalErr := errors.InternalError
            c.JSON(http.StatusInternalServerError, gin.H{
                "code":    internalErr.Code,
                "message": internalErr.Key,
                "data":    nil,
            })
            c.Abort()
        }
    })
}
```

## 配置管理

### 配置层级 (基于Viper)

1. **默认配置**: 代码中的默认值
2. **配置文件**: `config/config.yaml`
3. **环境变量**: 覆盖配置文件
4. **命令行参数**: 最高优先级

### 配置结构

```go
type Config struct {
    App      AppConfig      `yaml:"app" mapstructure:"app"`
    Database DatabaseConfig `yaml:"database" mapstructure:"database"`
    Log      LogConfig      `yaml:"log" mapstructure:"log"`
    Cookie   CookieConfig   `yaml:"cookie" mapstructure:"cookie"`
    Redis    RedisConfig    `yaml:"redis" mapstructure:"redis"`
    I18n     I18nConfig     `yaml:"i18n" mapstructure:"i18n"`
}

type AppConfig struct {
    Name    string `yaml:"name" mapstructure:"name"`
    Version string `yaml:"version" mapstructure:"version"`
    Port    int    `yaml:"port" mapstructure:"port"`
    Debug   bool   `yaml:"debug" mapstructure:"debug"`
}
```

### 配置加载

```go
// 使用Viper加载配置
viper.SetConfigName("config")
viper.SetConfigType("yaml")
viper.AddConfigPath("./config")
viper.AutomaticEnv()

if err := viper.ReadInConfig(); err != nil {
    return nil, err
}

var config Config
if err := viper.Unmarshal(&config); err != nil {
    return nil, err
}
```

## 日志系统 (基于slog)

### 日志级别

- **DEBUG**: 详细调试信息
- **INFO**: 一般信息记录
- **WARN**: 警告信息
- **ERROR**: 错误信息

### 结构化日志特性

- **零分配**: 高性能日志记录
- **类型安全**: 编译时类型检查
- **上下文感知**: 自动携带请求上下文
- **多输出**: 文件、控制台、远程日志服务

### 日志配置

```go
// 配置slog
logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
    Level: slog.LevelInfo,
    AddSource: true,
}))
slog.SetDefault(logger)
```

### 日志示例

```go
// 结构化日志记录
slog.Info("用户登录成功",
    slog.Int("user_id", 123),
    slog.String("username", "john_doe"),
    slog.String("ip", "192.168.1.100"),
    slog.Duration("duration", 150*time.Millisecond))

// 带上下文的日志
slog.InfoContext(ctx, "处理请求",
    slog.String("method", "POST"),
    slog.String("path", "/api/users"))
```

## 错误处理体系

### 错误码分类

- **1xxx**: 通用错误 (系统错误、验证错误等)
- **2xxx**: 用户模块错误
- **3xxx**: 日记模块错误
- **4xxx**: 内容模块错误
- **5xxx**: 网络和外部服务错误

### 错误处理流程

1. **业务层**: 创建AppError
2. **Handler层**: 调用ViewService处理错误
3. **中间件层**: 统一错误格式和日志记录
4. **客户端**: 接收统一格式的错误响应

## 性能特征

### Go 1.25泛型DI容器性能

- **解析速度**: ~57纳秒 (零反射设计)
- **内存分配**: 最小化内存分配
- **并发安全**: 读写锁+sync.Once优化
- **类型安全**: 编译时类型检查，运行时零错误
- **生命周期**: 瞬态(Bind)/单例(Singleton)/实例(Value)三种模式
- **错误检测**: 重复注册检查，防止意外覆盖
- **高级功能**: 安全获取(TryGet)、服务检查(Has)、测试覆盖(Rebind)

### 前端性能

- **Kiko框架**: 轻量级组件化框架
- **Tailwind CSS 4.x**: 现代化CSS框架
- **ES2023**: 现代JavaScript特性
- **模块化**: 按需加载和缓存

### 渲染性能

- **模板缓存**: Go template预编译缓存
- **页面数据**: 结构化数据组装
- **国际化**: 多语言消息缓存机制
- **视图服务**: 统一渲染管道

### 数据库性能 (SQLx + MySQL)

- **连接池**: 高效数据库连接池管理
- **预编译语句**: SQL语句预编译和缓存
- **事务管理**: 自动事务管理和回滚
- **查询优化**: 结构化查询和映射

### 日志性能 (slog)

- **零分配**: 高性能结构化日志
- **异步写入**: 非阻塞日志记录
- **类型安全**: 编译时类型检查
- **上下文传播**: 自动上下文携带

## 安全机制

### 认证和授权

- **Session管理**: Redis会话存储
- **Cookie安全**: HttpOnly、Secure、SameSite设置
- **权限控制**: 基于角色的权限管理
- **中间件保护**: 认证和授权中间件

### 数据安全

- **密码加密**: 安全哈希算法
- **数据传输**: HTTPS加密传输
- **输入验证**: Gin validator参数验证
- **SQL注入防护**: SQLx预编译语句

### 错误处理安全

- **敏感信息**: 不在错误响应中暴露敏感信息
- **日志脱敏**: 敏感数据脱敏记录
- **错误码**: 抽象的4位错误码，不暴露内部实现
- **Panic恢复**: 全局Panic恢复中间件

### 前端安全

- **XSS防护**: 模板自动转义
- **CSRF防护**: CSRF Token验证
- **内容安全**: CSP头部设置
- **现代浏览器**: 支持现代安全特性

## 部署架构

### 开发环境

```
iDiary Application (Go 1.25)
    ↓
MySQL 8.0+ (本地)
    ↓
Redis 6.0+ (本地)
```

### 生产环境

```
Load Balancer (Nginx)
    ↓
iDiary Application Cluster
    ↓
MySQL 8.0+ (主从复制)
    ↓
Redis Cluster (高可用)
```

### 容器化部署

```
Docker Compose:
- iDiary (Go 1.25 Application)
- MySQL 8.0+ (Database)
- Redis 6.0+ (Cache)
- Nginx (Static Files & Reverse Proxy)
```

### 云原生部署

```
Kubernetes:
- iDiary Deployment (多副本)
- MySQL StatefulSet
- Redis Cluster
- Ingress Controller
- ConfigMap & Secret
```

## 扩展点

### 水平扩展

- **无状态设计**: 应用层无状态，支持多实例部署
- **数据库扩展**: 支持MySQL读写分离和分库分表
- **缓存扩展**: Redis集群和分片支持
- **负载均衡**: 支持多种负载均衡策略

### 功能扩展

- **模块化设计**: 基于DDD的模块化架构
- **API扩展**: RESTful API和GraphQL支持
- **插件系统**: 预留插件接口和扩展点
- **多租户**: 支持SaaS多租户架构

### 技术扩展

- **微服务**: 可拆分为微服务架构
- **消息队列**: 支持异步消息处理
- **搜索引擎**: 集成Elasticsearch等搜索引擎
- **监控告警**: 集成Prometheus、Grafana等监控系统

### 存储扩展

- **云存储**: 支持AWS S3、阿里云OSS等云存储
- **文件系统**: 本地和分布式文件存储
- **数据备份**: 自动备份和灾难恢复机制
- **数据同步**: 支持多数据中心同步

## 总结

iDiary.ai Hebe系统采用现代化的Go 1.25技术栈，通过DDD分层架构、泛型依赖注入容器、统一错误处理等设计模式，构建了一个高性能、高可维护性和高可扩展性的智能日记管理系统。

### 核心优势

- **现代化技术栈**: Go 1.25、Gin 1.10.0、SQLx、Redis、gocron等
- **高性能DI容器**: 零反射设计、57纳秒解析速度、并发安全
- **类型安全**: 编译时类型检查，减少运行时错误
- **分层清晰**: 严格的依赖方向控制，遵循DDD最佳实践
- **统一错误处理**: 4位错误码体系，结构化错误定义
- **结构化日志**: 基于slog的高性能日志系统
- **前端现代化**: Kiko框架 + Tailwind CSS 4.x + ES2023
- **智能任务**: 基于gocron的异步邮件任务系统
- **多云存储**: 预留多云盘存储接口支持

### 设计理念

- **领域驱动**: 业务逻辑内聚，技术关注点分离
- **依赖倒置**: 接口与实现分离，便于测试和扩展
- **单一职责**: 每个模块职责明确，便于维护
- **开闭原则**: 对扩展开放，对修改封闭

这种架构设计确保了系统的高可维护性、可扩展性和稳定性，为未来的功能扩展和性能优化奠定了坚实的基础。
