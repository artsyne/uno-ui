# iDiary 系统启动过程详解

## 概述

iDiary采用分层启动架构，通过startup模式管理应用生命周期，实现了清晰的职责分离和优雅的启动流程。系统采用统一的错误处理机制，确保启动过程中的错误能够被正确捕获和处理。

## 启动流程图

```
main.go
├── 1. 基础日志设置 (utils.SetupBasicLogger)
├── 2. 配置初始化 (utils.config.InitConfig)
├── 3. 完整日志系统 (utils.SetLogger)
├── 4. 应用创建 (startup.NewApplication)
├── 5. 应用启动 (app.Startup)
│   ├── 5.1 核心基础设施初始化
│   │   ├── 配置验证和加载
│   │   ├── 数据库连接建立
│   │   └── DI容器创建和服务注册
│   └── 5.2 应用层组件初始化
│       ├── Cookie默认值设置
│       ├── 定时任务启动
│       └── Gin服务器配置
│           ├── 中间件设置
│           └── 路由注册
└── 6. 服务器运行 (app.Run)
```

## 详细启动过程

### 阶段1：程序入口 (main.go)

**文件位置**: `cmd/idiary/main.go`

**文件位置**: `internal/utils/logger/logger.go`

```go
func main() {
    // 1. 设置基本日志，用于应用启动阶段
    utils.SetupBasicLogger()

    // 2. 初始化配置
    if err := utils.config.InitConfig(); err != nil {
        slog.Error("初始化配置失败", "error", err)
        os.Exit(1)
    }

    // 3. 配置完整的日志系统（包含文件输出）
    logConfig := &utils.LoggerConfig{
        Path:       utils.config.GetConfig().Log.Path,
        Filename:   utils.config.GetConfig().Log.Filename,
        MaxSize:    utils.config.GetConfig().Log.MaxSize,
        MaxBackups: config.GetConfig().Log.MaxBackups,
        MaxAge:     config.GetConfig().Log.MaxAge,
        Compress:   config.GetConfig().Log.Compress,
        LocalTime:  config.GetConfig().Log.LocalTime,
        Level:      config.GetConfig().Log.Level,
        Format:     config.GetConfig().Log.Format,
    }

    if err := utils.SetLogger(logConfig); err != nil {
        slog.Error("日志系统初始化失败", "error", err)
        os.Exit(1)
    }

    // 4. 创建并启动应用
    app := startup.NewApplication()
    if err := app.Startup(); err != nil {
        slog.Error("应用启动失败", "error", err)
        os.Exit(1)
    }

    // 5. 运行应用
    if err := app.Run(); err != nil {
        slog.Error("应用运行失败", "error", err)
        os.Exit(1)
    }
}
```

**阶段特点**：

- 极简化的入口函数，只负责高层编排
- 完善的错误处理，启动失败时正确退出
- 分层的日志初始化，确保启动过程可追踪

### 阶段2：核心基础设施初始化

**文件位置**: `internal/core/startup/initializer.go`

这个阶段由`CoreInitializer`负责，包含以下步骤：

#### 2.1 配置系统初始化

```go
func (i *Initializer) initConfig() error {
    // 配置已在main.go中初始化，这里获取配置实例
    i.config = utils.config.GetConfig()

    // 验证关键配置项
    if i.config.App.Port == "" {
        return errors.New("应用端口未配置")
    }

    return nil
}
```

#### 2.2 数据库连接建立

```go
func (i *Initializer) initDatabase() error {
    db, err := database.Connect(i.config.Database)
    if err != nil {
        return fmt.Errorf("数据库连接失败: %w", err)
    }

    i.database = db
    return nil
}
```

#### 2.3 DI容器创建和服务注册

```go
func (i *Initializer) initContainer() error {
    i.container = di.New()

    slog.Info("成功初始化DI容器")
    return nil
}
```

### 阶段3：应用层组件初始化

**文件位置**: `internal/app/startup/application.go`

**相关中间件文件位置**:
- `internal/utils/middleware/recovery.go`
- `internal/utils/middleware/error_handler.go`
- `internal/app/middleware/user.go`

**文件位置**: `internal/utils/cookie/cookie.go`
#### 3.1 Cookie配置

```go
func (a *Application) initCookie() error {
    utils.SetCookieDefVal(
        a.config.Cookie.Domain,
        "/",
        a.config.Cookie.ExpireDays,
        a.config.Cookie.HttpOnly,
    )
    slog.Info("成功设置 Cookie 默认值")
    return nil
}
```

**文件位置**: `internal/domain/task/task.go`
#### 3.2 定时任务启动

```go
func (a *Application) initScheduler() error {
    domain.task.NewService().StartScheduler()
    slog.Info("成功启动定时任务")
    return nil
}
```

#### 3.3 Gin服务器配置

```go
func (a *Application) initGinServer() error {
    // 设置 Gin 模式
    gin.SetMode(a.config.App.GinMode)

    // 初始化默认gin server
    server := gin.New()

    // 初始化中间件
    server.Use(utils_middleware.PanicRecoveryMiddleware())   // panic recovery中间件
server.Use(utils_middleware.AppErrorHandlerMiddleware()) // 应用错误处理中间件
server.Use(app_middleware.SharedData())                      // 用户共享数据中间件

    // 初始化路由
    route.SetWithContainer(server, a.container)

    a.server = server
    return nil
}
```

### 阶段4：路由和中间件设置

**文件位置**: `internal/app/route/router.go`

```go
func SetWithContainer(router *gin.Engine, container *di.Container) {
    // 设置静态文件路由
    router.Static("/static", "./website/static")
    router.LoadHTMLGlob("website/views/**/*.html")

    // 从DI容器获取Handler
    userHandler := domain.user.CreateUserHandler(container)
    diaryHandler := domain.diary.CreateDiaryHandler(container)
    postHandler := domain.post.CreatePostHandler(container)

    // 注册路由组
    domain.user.SetupUserRoutes(router, userHandler)
    domain.diary.SetupDiaryRoutes(router, diaryHandler)
    domain.post.SetupPostRoutes(router, postHandler)
}
```

### 阶段5：服务器启动

```go
func (a *Application) Run() error {
    addr := ":" + a.config.App.Port
    slog.Info("Starting server", "addr", addr)

    if err := a.server.Run(addr); err != nil {
        slog.Error("服务器启动失败", "error", err, "addr", addr)
        os.Exit(1)
    }

    return nil
}
```

## 中间件加载顺序

1. **PanicRecoveryMiddleware**: 最先加载，捕获panic并记录日志
2. **AppErrorHandlerMiddleware**: 处理业务错误，统一错误响应格式
3. **SharedData**: 注入用户认证信息到上下文

## DI容器服务注册

### 核心服务

- Config: 配置实例
- Database: 数据库连接
- ViewService: 视图服务
- TemplateService: 模板服务

### 业务服务

- UserService: 用户业务逻辑
- DiaryService: 日记本业务逻辑
- PostService: 内容业务逻辑
- TaskService: 任务服务

### Handler层

- UserHandler: 用户HTTP处理器
- DiaryHandler: 日记本HTTP处理器
- PostHandler: 内容HTTP处理器
- HealthHandler: 健康检查HTTP处理器

## 配置文件结构

```yaml
# config.yaml
app:
  port: '8080'
  gin_mode: 'debug'

database:
  host: 'localhost'
  port: 3306
  name: 'idiary'
  username: 'root'
  password: 'password'

log:
  path: './logs'
  filename: 'app.log'
  level: 'info'
  format: 'json'

cookie:
  domain: 'localhost'
  expire_days: 7
  http_only: true
```

## 错误处理机制

### 启动阶段错误

- 任何启动阶段错误都会导致程序退出
- 详细的错误日志记录，便于排查问题
- 优雅的错误信息输出

### 运行时错误

- 通过中间件统一处理
- 结构化错误日志记录
- 统一的错误响应格式

## 优雅关闭（预留）

```go
func (a *Application) Shutdown() error {
    slog.Info("开始关闭应用程序...")

    // 1. 停止接收新请求
    // 2. 等待现有请求完成
    // 3. 关闭数据库连接
    // 4. 清理资源

    slog.Info("应用程序已关闭")
    return nil
}
```

## 监控和健康检查

### 启动状态监控

- 每个启动阶段都有明确的日志记录
- 启动耗时统计
- 启动失败原因追踪

### 健康检查端点

```go
// 可以添加健康检查路由
router.GET("/health", func(c *gin.Context) {
    c.JSON(200, gin.H{
        "status":    "healthy",
        "timestamp": time.Now().Unix(),
        "version":   "1.0.0",
    })
})
```

## 最佳实践

1. **分层启动**: 核心基础设施先于应用组件启动
2. **错误处理**: 每个阶段都有完善的错误处理
3. **日志记录**: 结构化日志记录启动过程
4. **依赖注入**: 通过DI容器管理组件依赖
5. **配置管理**: 集中的配置管理和验证
6. **优雅关闭**: 为未来扩展预留优雅关闭机制

这种启动架构的优势在于职责清晰、易于测试、便于扩展，为应用的稳定运行奠定了坚实基础。
