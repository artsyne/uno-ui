# iDiary.ai Hebe 中间件系统详解

## 概述

iDiary.ai Hebe采用分层中间件架构，通过Gin框架的中间件机制实现请求预处理、错误处理、用户认证等功能。中间件按功能分为系统级中间件(Utils层)和业务级中间件(App层)，确保关注点分离和职责清晰，与DDD分层架构完美融合。

## 中间件分层架构

```
HTTP Request
    ↓
┌─────────────────────────────────────────────────┐
│          系统级中间件层 (Utils Layer)             │
│   - utilsMiddleware.PanicRecovery()             │  ← 最高优先级，panic恢复
│   - utilsMiddleware.AppErrorHandler()           │  ← 统一错误处理
└─────────────────────────────────────────────────┘
    ↓
┌─────────────────────────────────────────────────┐
│          业务级中间件层 (App Layer)               │
│   - middleware.SharedData()                    │  ← 用户认证信息注入
│   - middleware.AuthRequired()                  │  ← 认证检查（按需）
└─────────────────────────────────────────────────┘
    ↓
┌─────────────────────────────────────────────────┐
│              控制器层 (Controller)                │
│   - UserController, DiaryController, etc.      │
└─────────────────────────────────────────────────┘
    ↓
┌─────────────────────────────────────────────────┐
│               服务层 (Service)                   │
│   - UserService, DiaryService, etc.            │
└─────────────────────────────────────────────────┘
    ↓
HTTP Response (JSON/HTML)
```

## 系统级中间件详解

### 1. PanicRecovery 中间件

**文件位置**: `internal/utils/middleware/panic.go`

**职责**: 捕获和处理运行时panic，防止应用崩溃

**优先级**: 最高（第一个加载）

```go
// PanicRecovery panic 恢复中间件
// 捕获并处理运行时 panic，记录详细的错误信息和堆栈跟踪
// 使用统一的错误响应格式
func PanicRecovery() gin.HandlerFunc {
    return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
        if recovered != nil {
            // 记录 panic 详细信息
            slog.Error("panic recovered",
                "panic", recovered,
                "path", c.Request.URL.Path,
                "method", c.Request.Method,
                "client_ip", c.ClientIP(),
                "user_agent", c.Request.UserAgent(),
                "stack", string(debug.Stack()))

            // 使用统一的错误响应格式
            internalErr := errors.InternalError
            response := gin.H{
                "code":    internalErr.Code,
                "message": internalErr.Key,
                "data":    nil,
            }

            c.JSON(http.StatusInternalServerError, response)
            c.Abort()
        }
    })
}
```

**特点**:

- **详细日志记录**: 记录panic内容、请求路径、客户端信息、完整堆栈
- **统一错误格式**: 使用预定义的`InternalError`错误码
- **优雅降级**: 防止整个应用崩溃，返回友好的错误响应

### 2. AppErrorHandler 中间件

**文件位置**: `internal/utils/middleware/apperror.go`

**职责**: 处理通过`c.Error()`添加的业务错误和应用错误

**优先级**: 系统级第二优先级

```go
// AppErrorHandler 应用错误处理中间件
// 处理通过 c.Error() 添加的业务错误和应用错误
// 使用统一的错误响应格式
func AppErrorHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        defer func() {
            if len(c.Errors) > 0 {
                lastError := c.Errors.Last()

                // 处理 AppError 类型的业务错误
                if appErr, ok := lastError.Err.(*errors.AppError); ok {
                    logAppError(c, appErr)
                    c.JSON(appErr.HttpCode, gin.H{
                        "code":    appErr.Code,
                        "message": appErr.Key,
                        "data":    nil,
                    })
                    return
                }

                // 处理其他类型的错误
                slog.Error("未知应用错误",
                    "error", lastError.Err.Error(),
                    "path", c.Request.URL.Path,
                    "method", c.Request.Method,
                    "user_agent", c.Request.UserAgent(),
                    "client_ip", c.ClientIP())

                internalErr := errors.InternalError
                c.JSON(http.StatusInternalServerError, gin.H{
                    "code":    internalErr.Code,
                    "message": internalErr.Key,
                    "data":    nil,
                })
            }
        }()

        c.Next()
    }
}

// logAppError 记录应用错误的详细信息
func logAppError(c *gin.Context, err *errors.AppError) {
    slog.Error("应用业务错误",
        "error_code", err.Code,
        "error_key", err.Key,
        "http_code", err.HttpCode,
        "method", c.Request.Method,
        "path", c.Request.URL.Path,
        "user_agent", c.Request.UserAgent(),
        "client_ip", c.ClientIP(),
        "details", err.Details)
}
```

**特点**:

- **统一错误格式**: 支持AppError的4位错误码体系
- **详细错误日志**: 记录错误码、HTTP状态码、请求信息、错误详情
- **错误分类处理**: 区分AppError和普通错误，分别处理
- **错误传播**: 不阻止错误向上传播，在defer中统一处理

## 业务级中间件详解

### 1. SharedData 中间件

**文件位置**: `internal/app/middleware/user.go`

**职责**: 解析用户登录信息并注入到上下文

**优先级**: 业务级第一优先级

```go
// SharedData 共享数据中间件
func SharedData() gin.HandlerFunc {
    return func(c *gin.Context) {
        auth := user.NewAuthHelper(config.GetConfig())
        // 解析用户登录信息
        authUser := auth.ParseUserFromCookie(c)

        // 注入登录状态和用户数据到上下文
        if authUser != nil {
            c.Set("IsLogin", true)
            c.Set("AuthUser", authUser)
        } else {
            c.Set("IsLogin", false)
        }

        c.Next()
    }
}
```

**特点**:

- **用户状态注入**: 自动解析Cookie中的用户信息
- **上下文共享**: 将用户数据注入到Gin上下文，供后续处理器使用
- **状态标识**: 提供`IsLogin`布尔值和`AuthUser`对象
- **无侵入性**: 不阻止请求继续，仅提供用户信息

### 2. AuthRequired 中间件

**文件位置**: `internal/app/middleware/user.go`

**职责**: 检查用户是否已登录，未登录则拒绝访问

**使用场景**: 保护需要登录的路由

```go
// AuthRequired 需要登录访问的中间件
func AuthRequired() gin.HandlerFunc {
    auth := user.NewAuthHelper(config.GetConfig())
    return func(c *gin.Context) {
        if isAuth, _ := auth.IsAuthenticated(c); isAuth {
            c.Next()
        } else {
            slog.Warn("User not authorized", "uri", c.Request.RequestURI)
            c.HTML(http.StatusForbidden, "error", nil)
            c.Abort()
        }
    }
}
```

**特点**:

- **认证检查**: 使用AuthHelper验证用户登录状态
- **访问控制**: 未登录用户无法访问受保护的路由
- **友好错误页**: 返回HTML错误页面而非JSON响应
- **日志记录**: 记录未授权访问尝试

## 中间件注册和加载顺序

### 在应用启动中的配置

**文件位置**: `internal/app/startup/application.go`

```go
func (a *Application) initGinServer() error {
    gin.SetMode(a.config.App.GinMode)
    server := gin.New()

    // 配置中间件链 - 严格按顺序加载
    server.Use(utilsMiddleware.PanicRecovery())   // 1. panic恢复中间件
    server.Use(utilsMiddleware.AppErrorHandler()) // 2. 应用错误处理中间件
    server.Use(middleware.SharedData())           // 3. 共享数据中间件

    // 配置路由 - 通过DI容器创建控制器
    route.SetWithContainer(server, a.container)

    a.server = server
    slog.Info("成功初始化Gin服务器")
    return nil
}
```

### 中间件执行流程

```
HTTP Request
    ↓
1. utilsMiddleware.PanicRecovery()
   └── 捕获panic，防止应用崩溃
    ↓
2. utilsMiddleware.AppErrorHandler()
   └── 处理业务错误，统一错误响应
    ↓
3. middleware.SharedData()
   └── 解析用户信息，注入上下文
    ↓
4. middleware.AuthRequired() (可选，在特定路由使用)
   └── 检查用户认证状态
    ↓
5. Controller Handler
   └── 业务逻辑处理
    ↓
6. ViewService Response
   └── 统一响应格式
```

**执行顺序解释**:

1. **PanicRecovery**: 最先执行，确保任何panic都能被捕获，防止应用崩溃
2. **AppErrorHandler**: 第二执行，处理业务层抛出的错误，统一错误响应格式
3. **SharedData**: 解析用户认证信息，为后续中间件和控制器提供用户上下文
4. **AuthRequired**: 按需使用，在需要认证的路由上检查用户状态

## 中间件最佳实践

### 1. 错误处理模式

```go
// 在控制器中使用中间件错误处理
func (ctrl *UserController) UpdateProfile(c *gin.Context) {
    user, err := ctrl.userService.GetCurrentUser(c)
    if err != nil {
        // 将错误传递给AppErrorHandler中间件处理
        _ = c.Error(err)
        return
    }

    // 业务逻辑...
}
```

### 2. 用户信息获取模式

```go
// 在控制器中获取SharedData注入的用户信息
func (ctrl *DiaryController) CreateDiary(c *gin.Context) {
    // 从上下文获取用户信息
    isLogin, _ := c.Get("IsLogin")
    if !isLogin.(bool) {
        err := errors.AuthRequired.Err()
        _ = c.Error(err)
        return
    }

    authUser, _ := c.Get("AuthUser")
    user := authUser.(*user.AuthUser)

    // 使用用户信息进行业务处理...
}
```

### 3. 路由级中间件使用

```go
// 在路由注册时使用AuthRequired中间件
func (ctrl *UserController) RegisterRoutes(router *gin.Engine) {
    userGroup := router.Group("/user")

    // 公开路由
    userGroup.POST("/login", ctrl.Login)
    userGroup.POST("/register", ctrl.Register)

    // 需要认证的路由
    authGroup := userGroup.Group("").Use(middleware.AuthRequired())
    authGroup.GET("/profile", ctrl.GetProfile)
    authGroup.PUT("/profile", ctrl.UpdateProfile)
}
```

### 4. 中间件开发规范

```go
// 标准中间件函数签名
func CustomMiddleware(params ...interface{}) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 前置处理逻辑

        // 调用下一个中间件/处理器
        c.Next()

        // 后置处理逻辑（可选）
    }
}
```

## 性能考量

### 1. 中间件性能优化

- **减少内存分配**: 复用对象，避免频繁创建
- **懒加载**: 只在需要时才进行复杂计算
- **缓存策略**: 缓存频繁访问的数据（如用户信息）

### 2. 执行顺序优化

- **系统级中间件优先**: 确保错误处理和panic恢复最先执行
- **业务中间件后置**: 用户认证等业务逻辑在系统级中间件之后
- **按需使用**: 只在需要的路由上使用AuthRequired等中间件

### 3. 日志性能

- **结构化日志**: 使用slog的零分配特性
- **日志级别控制**: 生产环境减少DEBUG级别日志
- **异步日志**: 避免日志记录阻塞请求处理

## 总结

iDiary.ai Hebe的中间件系统展现了以下特点：

### 🚀 架构优势

- **分层清晰**: 系统级中间件处理技术关注点，业务级中间件处理业务逻辑
- **统一错误处理**: AppError体系与中间件完美集成
- **高性能设计**: 基于slog的零分配日志，最小化中间件开销

### 🔧 开发友好

- **简洁API**: 清晰的中间件注册和使用方式
- **上下文共享**: 自动注入用户信息，简化控制器开发
- **类型安全**: 结合DI容器提供编译时类型安全

### 📈 可扩展性

- **可插拔设计**: 中间件可以根据需要灵活添加或移除
- **标准化接口**: 遵循Gin中间件规范，易于集成第三方中间件
- **模块化架构**: 与DDD分层架构完美融合

这套中间件系统为iDiary.ai Hebe提供了坚实的技术基础，确保了系统的稳定性、可维护性和可扩展性。
