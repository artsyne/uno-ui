# 中间件系统详解

## 概述

iDiary采用分层中间件架构，通过Gin框架的中间件机制实现请求预处理、错误处理、用户认证等功能。中间件按功能分为核心中间件（Core层）和业务中间件（Module层），确保关注点分离和职责清晰。

## 中间件分层架构

```
Gin Request Pipeline
    ↓
┌─────────────────────────────┐
│    Core Middleware Layer    │  (核心中间件层)
│  - PanicRecoveryMiddleware  │  ← 最优先级，捕获panic
│  - AppErrorHandlerMiddleware│  ← 统一错误处理
└─────────────────────────────┘
    ↓
┌─────────────────────────────┐
│  Business Middleware Layer  │  (业务中间件层)
│  - SharedData              │  ← 用户认证信息注入
│  - AuthRequired            │  ← 认证检查
│  - Permission              │  ← 权限验证
└─────────────────────────────┘
    ↓
Handler Processing
```

## 核心中间件详解

### 1. PanicRecoveryMiddleware

**文件位置**: `internal/core/middleware/error.go`

**职责**: 捕获和处理panic，防止应用崩溃

```go
func PanicRecoveryMiddleware() gin.HandlerFunc {
    return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
        // 记录详细的panic信息
        slog.Error("panic recovered",
            "panic", recovered,
            "path", c.Request.URL.Path,
            "method", c.Request.Method,
            "client_ip", c.ClientIP(),
            "user_agent", c.Request.UserAgent(),
            "stack", string(debug.Stack()))

        // 返回统一的错误响应
        errorResponse := gin.H{
            "code":    1000,
            "message": "系统内部错误",
            "data":    nil,
        }

        c.JSON(http.StatusInternalServerError, errorResponse)
    })
}
```

**特点**:

- 捕获所有未处理的panic
- 记录完整的调用堆栈
- 返回用户友好的错误信息
- 防止敏感信息泄露

### 2. AppErrorHandlerMiddleware

**文件位置**: `internal/core/middleware/error.go`

**职责**: 处理业务错误，统一错误响应格式

```go
func AppErrorHandlerMiddleware() gin.HandlerFunc {
    return gin.HandlerFunc(func(c *gin.Context) {
        c.Next()

        // 处理c.Error()收集的错误
        if len(c.Errors) > 0 {
            err := c.Errors.Last()

            if appErr, ok := err.Err.(*errors.AppError); ok {
                // 处理AppError类型错误
                handleAppError(c, appErr)
            } else {
                // 处理其他类型错误
                handleGenericError(c, err.Err)
            }
        }
    })
}

func handleAppError(c *gin.Context, appErr *errors.AppError) {
    // 记录错误日志
    slog.Error("业务错误",
        "code", appErr.Code,
        "message", appErr.MsgKey,
        "http_code", appErr.HttpCode,
        "path", c.Request.URL.Path,
        "method", c.Request.Method)

    // 返回错误响应
    c.JSON(appErr.HttpCode, gin.H{
        "code":    appErr.Code,
        "message": appErr.MsgKey,
        "data":    nil,
    })
}
```

**特点**:

- 统一的错误响应格式
- 结构化错误日志记录
- 区分业务错误和系统错误
- 支持国际化错误消息

## 业务中间件详解

### 1. SharedData 中间件

**文件位置**: `internal/module/user/middlewares.go`

**职责**: 注入用户认证信息到请求上下文

```go
func SharedData() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 从session或token中获取用户信息
        authUser := getUserFromSession(c)

        if authUser != nil {
            // 注入认证用户信息到上下文
            c.Set("auth_user", authUser)
            c.Set("user_id", authUser.UserId)
            c.Set("is_login", true)
        } else {
            c.Set("is_login", false)
        }

        c.Next()
    }
}

func getUserFromSession(c *gin.Context) *AuthUser {
    // 1. 尝试从JWT token获取
    if token := c.GetHeader("Authorization"); token != "" {
        if user := validateJWTToken(token); user != nil {
            return user
        }
    }

    // 2. 尝试从session获取
    if sessionID := getSessionID(c); sessionID != "" {
        if user := getSessionUser(sessionID); user != nil {
            return user
        }
    }

    // 3. 尝试从cookie获取记住我信息
    if rememberToken := c.GetCookie("remember_token"); rememberToken != "" {
        if user := validateRememberToken(rememberToken); user != nil {
            return user
        }
    }

    return nil
}
```

**特点**:

- 支持多种认证方式（JWT、Session、Remember Token）
- 自动注入用户信息到上下文
- 为模板函数提供数据支持
- 不阻断请求，只是注入信息

### 2. AuthRequired 中间件

**文件位置**: `internal/module/user/middlewares.go`

**职责**: 检查用户是否已认证，未认证则拦截请求

```go
func AuthRequired() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 检查用户是否已登录
        isLogin, exists := c.Get("is_login")
        if !exists || !isLogin.(bool) {
            // 根据请求类型返回不同响应
            if isAPIRequest(c) {
                // API请求返回JSON错误
                c.JSON(http.StatusUnauthorized, gin.H{
                    "code":    2001,
                    "message": "未认证",
                    "data":    nil,
                })
            } else {
                // 页面请求重定向到登录页
                redirectURL := c.Request.URL.Path
                if c.Request.URL.RawQuery != "" {
                    redirectURL += "?" + c.Request.URL.RawQuery
                }

                loginURL := "/login?redirect=" + url.QueryEscape(redirectURL)
                c.Redirect(http.StatusTemporaryRedirect, loginURL)
            }

            c.Abort()
            return
        }

        c.Next()
    }
}

func isAPIRequest(c *gin.Context) bool {
    // 检查请求头或路径来判断是否为API请求
    contentType := c.GetHeader("Content-Type")
    accept := c.GetHeader("Accept")
    path := c.Request.URL.Path

    return strings.Contains(contentType, "application/json") ||
           strings.Contains(accept, "application/json") ||
           strings.HasPrefix(path, "/api/")
}
```

**特点**:

- 区分API请求和页面请求
- API请求返回JSON错误
- 页面请求重定向到登录页
- 支持登录后跳转回原页面

### 3. Permission 中间件

**文件位置**: `internal/module/user/middlewares.go`

**职责**: 检查用户权限

```go
func Permission(requiredGroup int) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 先检查是否已认证
        isLogin, _ := c.Get("is_login")
        if !isLogin.(bool) {
            c.JSON(http.StatusUnauthorized, gin.H{
                "code":    2001,
                "message": "未认证",
            })
            c.Abort()
            return
        }

        // 获取用户信息
        authUser, exists := c.Get("auth_user")
        if !exists {
            c.JSON(http.StatusInternalServerError, gin.H{
                "code":    1000,
                "message": "系统错误",
            })
            c.Abort()
            return
        }

        user := authUser.(*AuthUser)

        // 检查用户组权限
        if user.UserGroup < requiredGroup {
            c.JSON(http.StatusForbidden, gin.H{
                "code":    2003,
                "message": "权限不足",
            })
            c.Abort()
            return
        }

        c.Next()
    }
}

// 使用示例
func AdminRequired() gin.HandlerFunc {
    return Permission(1) // 需要管理员权限
}

func ModeratorRequired() gin.HandlerFunc {
    return Permission(2) // 需要版主权限
}
```

## 中间件注册和使用

### 1. 全局中间件注册

**文件位置**: `internal/app/bootstrap/application.go`

```go
func (a *Application) initGinServer() error {
    server := gin.New()

    // 注册核心中间件（按顺序）
    server.Use(middleware.PanicRecoveryMiddleware())   // 1. Panic恢复
    server.Use(middleware.AppErrorHandlerMiddleware()) // 2. 错误处理
    server.Use(user.SharedData())                      // 3. 用户数据注入

    // 设置路由
    route.SetWithContainer(server, a.container)

    a.server = server
    return nil
}
```

### 2. 路由组中间件使用

**文件位置**: `internal/app/route/router.go`

```go
func setupUserRoutes(router *gin.Engine, userHandler *user.Handler) {
    // 公开路由（无需认证）
    public := router.Group("/")
    {
        public.GET("/login", userHandler.ShowLoginForm)
        public.POST("/login", userHandler.HandleLogin)
        public.GET("/register", userHandler.ShowRegisterForm)
        public.POST("/register", userHandler.HandleRegister)
    }

    // 需要认证的路由
    auth := router.Group("/")
    auth.Use(user.AuthRequired())
    {
        auth.GET("/u/:user_id", userHandler.ShowUserProfile)
        auth.GET("/settings", userHandler.ShowSettings)
        auth.POST("/settings", userHandler.UpdateSettings)
        auth.POST("/logout", userHandler.HandleLogout)
    }

    // 管理员路由
    admin := router.Group("/admin")
    admin.Use(user.AuthRequired())
    admin.Use(user.AdminRequired())
    {
        admin.GET("/users", userHandler.ListUsers)
        admin.POST("/users/:id/disable", userHandler.DisableUser)
        admin.POST("/users/:id/enable", userHandler.EnableUser)
    }
}
```

### 3. 单个路由中间件使用

```go
// 特定路由使用中间件
router.GET("/sensitive-data",
    user.AuthRequired(),
    user.Permission(1),
    handler.GetSensitiveData)

// 链式中间件
router.POST("/admin/action",
    user.AuthRequired(),
    user.AdminRequired(),
    customLoggingMiddleware(),
    handler.AdminAction)
```

## 自定义中间件开发

### 1. 请求日志中间件

```go
func RequestLoggingMiddleware() gin.HandlerFunc {
    return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
        var statusColor, methodColor, resetColor string
        if param.IsOutputColor() {
            statusColor = param.StatusCodeColor()
            methodColor = param.MethodColor()
            resetColor = param.ResetColor()
        }

        return fmt.Sprintf("%s[%s] %s \"%s %s %s %d %s\" %s\n%s",
            statusColor, param.TimeStamp.Format("2006/01/02 - 15:04:05"),
            methodColor, param.Method, resetColor,
            param.Path,
            param.StatusCode,
            param.Latency,
            param.ClientIP,
            param.ErrorMessage,
        )
    })
}
```

### 2. CORS中间件

```go
func CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        origin := c.GetHeader("Origin")

        // 检查是否允许该origin
        if isAllowedOrigin(origin) {
            c.Header("Access-Control-Allow-Origin", origin)
        }

        c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
        c.Header("Access-Control-Allow-Credentials", "true")
        c.Header("Access-Control-Max-Age", "86400")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(http.StatusNoContent)
            return
        }

        c.Next()
    }
}
```

### 3. 限流中间件

```go
func RateLimitMiddleware(limit int, window time.Duration) gin.HandlerFunc {
    limiter := make(map[string]*rate.Limiter)
    mu := sync.RWMutex{}

    return func(c *gin.Context) {
        ip := c.ClientIP()

        mu.RLock()
        l, exists := limiter[ip]
        mu.RUnlock()

        if !exists {
            mu.Lock()
            l = rate.NewLimiter(rate.Every(window), limit)
            limiter[ip] = l
            mu.Unlock()
        }

        if !l.Allow() {
            c.JSON(http.StatusTooManyRequests, gin.H{
                "code":    4001,
                "message": "请求频率过高",
            })
            c.Abort()
            return
        }

        c.Next()
    }
}
```

## 中间件最佳实践

### 1. 执行顺序

中间件的执行顺序很重要，遵循以下原则：

1. **Panic恢复** - 最先执行，确保能捕获所有panic
2. **请求日志** - 记录所有请求
3. **CORS处理** - 处理跨域请求
4. **限流控制** - 防止滥用
5. **用户认证** - 识别用户身份
6. **权限检查** - 验证用户权限
7. **业务中间件** - 具体业务逻辑

### 2. 错误处理

```go
func MyMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 中间件逻辑
        if err := doSomething(); err != nil {
            // 使用统一错误处理
            appErr := errors.NewAppError(1001, "middleware_error", http.StatusInternalServerError)
            c.Error(appErr) // 传递给错误处理中间件
            c.Abort()
            return
        }

        c.Next()
    }
}
```

### 3. 上下文数据传递

```go
func DataInjectionMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 注入数据到上下文
        c.Set("request_id", generateRequestID())
        c.Set("start_time", time.Now())

        c.Next()

        // 请求结束后处理
        duration := time.Since(c.GetTime("start_time"))
        slog.Info("请求完成",
            "request_id", c.GetString("request_id"),
            "duration", duration)
    }
}
```

### 4. 性能优化

```go
func OptimizedMiddleware() gin.HandlerFunc {
    // 预编译正则表达式
    pathRegex := regexp.MustCompile(`^/api/v\d+/`)

    return func(c *gin.Context) {
        // 快速路径匹配
        if !pathRegex.MatchString(c.Request.URL.Path) {
            c.Next()
            return
        }

        // 复杂逻辑只在匹配时执行
        if err := complexLogic(c); err != nil {
            c.Error(err)
            c.Abort()
            return
        }

        c.Next()
    }
}
```

## 调试和监控

### 1. 中间件执行追踪

```go
func TracingMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        startTime := time.Now()
        middlewareName := "TracingMiddleware"

        slog.Debug("中间件开始",
            "middleware", middlewareName,
            "path", c.Request.URL.Path)

        c.Next()

        duration := time.Since(startTime)
        slog.Debug("中间件结束",
            "middleware", middlewareName,
            "duration", duration,
            "status", c.Writer.Status())
    }
}
```

### 2. 性能监控

```go
func PerformanceMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()

        c.Next()

        duration := time.Since(start)

        // 记录慢请求
        if duration > 1*time.Second {
            slog.Warn("慢请求检测",
                "path", c.Request.URL.Path,
                "method", c.Request.Method,
                "duration", duration,
                "status", c.Writer.Status())
        }

        // 发送到监控系统
        sendMetrics(c.Request.URL.Path, c.Request.Method, duration, c.Writer.Status())
    }
}
```

这套中间件系统为iDiary提供了完整的请求处理流程控制，确保了安全性、可维护性和性能监控。
