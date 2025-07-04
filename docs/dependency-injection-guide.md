# 依赖注入容器使用指南

## 概述

iDiary采用自研的轻量级依赖注入（DI）容器，提供类型安全的泛型接口和高性能的服务解析能力。该DI容器是整个应用架构的核心基础设施，实现了19.89纳秒的解析速度和零内存分配的优异性能。

## 核心概念

### 1. 容器架构

```go
type Container struct {
    mu          sync.RWMutex                 // 读写锁，保护并发访问
    providers   map[reflect.Type]*provider   // 服务提供者映射
    singletons  map[reflect.Type]interface{} // 单例实例缓存
    initialized bool                         // 初始化标志
}

type provider struct {
    factory   reflect.Value     // 工厂函数
    lifecycle ServiceLifecycle  // 生命周期管理
    deps      []reflect.Type    // 依赖类型列表
}
```

### 2. 服务生命周期

#### Singleton（单例）

- 全应用生命周期内只创建一次
- 适用于无状态服务、配置、数据库连接等
- 线程安全，支持并发访问

#### Transient（瞬态）

- 每次请求都创建新实例
- 适用于有状态服务、请求特定的服务
- 无缓存，每次都调用工厂函数

```go
type ServiceLifecycle int

const (
    Singleton ServiceLifecycle = iota  // 单例模式
    Transient                          // 瞬态模式
)
```

## 基础使用

### 1. 服务注册

#### 注册单例服务

```go
// 注册数据库服务（单例）
container.Register(func() *database.DB {
    db, err := database.Connect(config.DatabaseURL)
    if err != nil {
        slog.Error("数据库连接失败", "error", err)
        panic(err)
    }
    return db
}, di.Singleton)

// 注册用户Repository（单例）
container.Register(func(db *database.DB) user.RepositoryInterface {
    return user.NewRepository(db)
}, di.Singleton)

// 注册用户Service（单例）
container.Register(func(repo user.RepositoryInterface) user.ServiceInterface {
    return user.NewService(repo)
}, di.Singleton)
```

#### 注册瞬态服务

```go
// 注册HTTP客户端（瞬态）
container.Register(func() *http.Client {
    return &http.Client{
        Timeout: 30 * time.Second,
        Transport: &http.Transport{
            MaxIdleConns:        100,
            MaxIdleConnsPerHost: 10,
        },
    }
}, di.Transient)

// 注册请求处理器（瞬态）
container.Register(func(userService user.ServiceInterface, render *render.Service) *user.Handler {
    return user.NewHandler(userService, render)
}, di.Transient)
```

### 2. 服务解析

#### 基本解析

```go
// 解析单例服务
userService := di.Resolve[user.ServiceInterface](container)

// 解析瞬态服务
handler := di.Resolve[*user.Handler](container)
```

#### 在Handler中使用

```go
func SetupRoutes(router *gin.Engine, container *di.Container) {
    // 解析Handler实例
    userHandler := di.Resolve[*user.Handler](container)
    diaryHandler := di.Resolve[*diary.Handler](container)
    postHandler := di.Resolve[*post.Handler](container)

    // 注册路由
    router.GET("/users", userHandler.ListUsers)
    router.GET("/diary", diaryHandler.ShowDiary)
    router.POST("/posts", postHandler.CreatePost)
}
```

## 高级用法

### 1. 条件注册

```go
// 根据环境注册不同的服务实现
if config.Environment == "production" {
    // 生产环境使用Redis缓存
    container.Register(func() cache.Interface {
        return cache.NewRedisCache(config.RedisURL)
    }, di.Singleton)
} else {
    // 开发环境使用内存缓存
    container.Register(func() cache.Interface {
        return cache.NewMemoryCache()
    }, di.Singleton)
}
```

### 2. 接口与实现注册

```go
// 定义接口
type EmailService interface {
    SendEmail(to, subject, body string) error
}

// 注册实现
container.Register(func() EmailService {
    if config.EmailProvider == "smtp" {
        return email.NewSMTPService(config.SMTPConfig)
    } else {
        return email.NewSendGridService(config.SendGridAPIKey)
    }
}, di.Singleton)

// 在其他服务中使用
container.Register(func(emailService EmailService) *user.Service {
    return user.NewService(emailService)
}, di.Singleton)
```

### 3. 配置注入

```go
// 注册配置
container.Register(func() *config.Config {
    return config.Load()
}, di.Singleton)

// 服务使用配置
container.Register(func(cfg *config.Config) *logger.Service {
    return logger.NewService(cfg.LogLevel, cfg.LogFile)
}, di.Singleton)
```

## 在iDiary中的实际应用

### 1. 应用启动时的服务注册

**文件位置**: `internal/app/bootstrap/application.go`

```go
func (a *Application) registerServices() error {
    // 1. 注册核心基础设施
    a.registerCoreServices()

    // 2. 注册数据访问层
    a.registerRepositories()

    // 3. 注册业务服务层
    a.registerServices()

    // 4. 注册应用服务层
    a.registerAppServices()

    return nil
}

func (a *Application) registerCoreServices() {
    // 数据库连接
    a.container.Register(func() *database.DB {
        db, err := database.NewConnection()
        if err != nil {
            panic(err)
        }
        return db
    }, di.Singleton)

    // 渲染服务
    a.container.Register(func() *render.Service {
        return render.NewService()
    }, di.Singleton)
}

func (a *Application) registerRepositories() {
    // 用户Repository
    a.container.Register(func(db *database.DB) user.RepositoryInterface {
        return user.NewRepository(db)
    }, di.Singleton)

    // 日记Repository
    a.container.Register(func(db *database.DB) diary.RepositoryInterface {
        return diary.NewRepository(db)
    }, di.Singleton)

    // 文章Repository
    a.container.Register(func(db *database.DB) post.RepositoryInterface {
        return post.NewRepository(db)
    }, di.Singleton)
}
```

### 2. 模块级别的服务组织

```go
// User模块服务注册
func RegisterUserServices(container *di.Container) {
    // Repository
    container.Register(func(db *database.DB) user.RepositoryInterface {
        return user.NewRepository(db)
    }, di.Singleton)

    // Service
    container.Register(func(repo user.RepositoryInterface) user.ServiceInterface {
        return user.NewService(repo)
    }, di.Singleton)

    // Handler
    container.Register(func(
        service user.ServiceInterface,
        render *render.Service,
    ) *user.Handler {
        return user.NewHandler(service, render)
    }, di.Transient)
}

// Diary模块服务注册
func RegisterDiaryServices(container *di.Container) {
    container.Register(func(db *database.DB) diary.RepositoryInterface {
        return diary.NewRepository(db)
    }, di.Singleton)

    container.Register(func(repo diary.RepositoryInterface) diary.ServiceInterface {
        return diary.NewService(repo)
    }, di.Singleton)

    container.Register(func(
        service diary.ServiceInterface,
        userService user.ServiceInterface,
        render *render.Service,
    ) *diary.Handler {
        return diary.NewHandler(service, userService, render)
    }, di.Transient)
}
```

### 3. 路由中的依赖解析

**文件位置**: `internal/app/route/router.go`

```go
func SetWithContainer(router *gin.Engine, container *di.Container) {
    // 解析Handler实例
    userHandler := di.Resolve[*user.Handler](container)
    diaryHandler := di.Resolve[*diary.Handler](container)
    postHandler := di.Resolve[*post.Handler](container)
    taskHandler := di.Resolve[*task.Handler](container)

    // 用户相关路由
    setupUserRoutes(router, userHandler)

    // 日记相关路由
    setupDiaryRoutes(router, diaryHandler)

    // 文章相关路由
    setupPostRoutes(router, postHandler)

    // 任务相关路由
    setupTaskRoutes(router, taskHandler)
}

func setupUserRoutes(router *gin.Engine, handler *user.Handler) {
    // 公开路由
    router.GET("/login", handler.ShowLoginForm)
    router.POST("/login", handler.HandleLogin)
    router.GET("/register", handler.ShowRegisterForm)
    router.POST("/register", handler.HandleRegister)

    // 需要认证的路由
    auth := router.Group("/")
    auth.Use(user.AuthRequired())
    {
        auth.GET("/settings", handler.ShowSettings)
        auth.POST("/settings", handler.UpdateSettings)
        auth.POST("/logout", handler.HandleLogout)
    }
}
```

## 性能优化

### 1. 预热策略

```go
func (a *Application) warmupServices() error {
    // 预热关键单例服务
    criticalServices := []interface{}{
        (*database.DB)(nil),
        (*render.Service)(nil),
        (user.ServiceInterface)(nil),
        (diary.ServiceInterface)(nil),
    }

    for _, serviceType := range criticalServices {
        start := time.Now()

        // 通过反射获取类型并解析
        t := reflect.TypeOf(serviceType).Elem()
        service := a.container.resolve(t)

        duration := time.Since(start)
        slog.Debug("服务预热完成",
            "service", t.String(),
            "duration", duration)

        if service == nil {
            return fmt.Errorf("预热服务失败: %s", t.String())
        }
    }

    return nil
}
```

### 2. 性能监控

```go
func (c *Container) resolveWithMetrics(serviceType reflect.Type) interface{} {
    start := time.Now()
    defer func() {
        duration := time.Since(start)

        // 记录解析性能
        slog.Debug("服务解析性能",
            "service", serviceType.String(),
            "duration", duration)

        // 发送到监控系统
        metrics.RecordServiceResolution(serviceType.String(), duration)
    }()

    return c.resolve(serviceType)
}
```

### 3. 内存优化

```go
func (c *Container) optimizeMemory() {
    c.mu.Lock()
    defer c.mu.Unlock()

    // 清理未使用的瞬态服务提供者
    for serviceType, provider := range c.providers {
        if provider.lifecycle == Transient {
            // 检查是否有其他服务依赖此服务
            if !c.hasDependent(serviceType) {
                delete(c.providers, serviceType)
                slog.Debug("清理未使用的瞬态服务", "service", serviceType.String())
            }
        }
    }
}
```

## 错误处理和调试

### 1. 依赖循环检测

```go
func (c *Container) detectCircularDependency(serviceType reflect.Type, visiting map[reflect.Type]bool) error {
    if visiting[serviceType] {
        return fmt.Errorf("检测到循环依赖: %s", serviceType.String())
    }

    visiting[serviceType] = true
    defer delete(visiting, serviceType)

    provider, exists := c.providers[serviceType]
    if !exists {
        return nil
    }

    for _, depType := range provider.deps {
        if err := c.detectCircularDependency(depType, visiting); err != nil {
            return err
        }
    }

    return nil
}
```

### 2. 服务注册验证

```go
func (c *Container) validateRegistration(factory interface{}, lifecycle ServiceLifecycle) error {
    factoryType := reflect.TypeOf(factory)

    // 检查是否为函数
    if factoryType.Kind() != reflect.Func {
        return errors.New("工厂必须是函数类型")
    }

    // 检查返回值
    if factoryType.NumOut() != 1 {
        return errors.New("工厂函数必须返回一个值")
    }

    // 检查参数类型
    for i := 0; i < factoryType.NumIn(); i++ {
        paramType := factoryType.In(i)
        if !c.canResolve(paramType) {
            return fmt.Errorf("无法解析参数类型: %s", paramType.String())
        }
    }

    return nil
}
```

### 3. 调试工具

```go
func (c *Container) PrintDependencyGraph() {
    c.mu.RLock()
    defer c.mu.RUnlock()

    fmt.Println("=== 依赖注入容器服务图 ===")

    for serviceType, provider := range c.providers {
        fmt.Printf("服务: %s\n", serviceType.String())
        fmt.Printf("  生命周期: %s\n", provider.lifecycle)
        fmt.Printf("  依赖: ")

        if len(provider.deps) == 0 {
            fmt.Println("无")
        } else {
            for i, dep := range provider.deps {
                if i > 0 {
                    fmt.Print(", ")
                }
                fmt.Print(dep.String())
            }
            fmt.Println()
        }
        fmt.Println()
    }
}

func (c *Container) GetServiceStats() map[string]interface{} {
    c.mu.RLock()
    defer c.mu.RUnlock()

    stats := map[string]interface{}{
        "total_services":     len(c.providers),
        "singleton_services": 0,
        "transient_services": 0,
        "initialized_singletons": len(c.singletons),
    }

    for _, provider := range c.providers {
        if provider.lifecycle == Singleton {
            stats["singleton_services"] = stats["singleton_services"].(int) + 1
        } else {
            stats["transient_services"] = stats["transient_services"].(int) + 1
        }
    }

    return stats
}
```

## 最佳实践

### 1. 服务设计原则

```go
// ✅ 好的做法：接口依赖
type UserService struct {
    repo   user.RepositoryInterface  // 依赖接口
    email  EmailServiceInterface    // 依赖接口
    logger LoggerInterface          // 依赖接口
}

// ❌ 不好的做法：具体类型依赖
type UserService struct {
    repo   *user.Repository  // 直接依赖具体类型
    email  *SMTPService      // 具体实现
    logger *FileLogger      // 具体实现
}
```

### 2. 生命周期选择

```go
// 单例服务：无状态、共享资源
container.Register(func() *database.DB { /*...*/ }, di.Singleton)           // ✅ 数据库连接
container.Register(func() *config.Config { /*...*/ }, di.Singleton)         // ✅ 配置
container.Register(func() user.RepositoryInterface { /*...*/ }, di.Singleton) // ✅ Repository

// 瞬态服务：有状态、请求特定
container.Register(func() *user.Handler { /*...*/ }, di.Transient)          // ✅ HTTP Handler
container.Register(func() *http.Client { /*...*/ }, di.Transient)           // ✅ HTTP客户端
container.Register(func() *session.Context { /*...*/ }, di.Transient)       // ✅ 会话上下文
```

### 3. 错误处理

```go
// 工厂函数中的错误处理
container.Register(func() (*database.DB, error) {
    db, err := database.Connect(config.DatabaseURL)
    if err != nil {
        slog.Error("数据库连接失败", "error", err)
        return nil, err
    }
    return db, nil
}, di.Singleton)

// 服务解析时的错误处理
userService, err := di.TryResolve[user.ServiceInterface](container)
if err != nil {
    slog.Error("解析用户服务失败", "error", err)
    return err
}
```

### 4. 模块化注册

```go
// 按模块组织服务注册
func RegisterAllServices(container *di.Container) error {
    // 注册核心服务
    if err := core.RegisterServices(container); err != nil {
        return err
    }

    // 注册业务模块
    user.RegisterServices(container)
    diary.RegisterServices(container)
    post.RegisterServices(container)
    task.RegisterServices(container)

    return nil
}
```

这个DI容器为iDiary项目提供了强大而高效的依赖管理能力，是整个架构的核心基础设施。通过合理使用，可以大大提高代码的可维护性和可测试性。
