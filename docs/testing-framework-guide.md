# iDiary.ai Hebe 测试框架指南

## 概述

iDiary.ai Hebe 采用现代化的测试框架，基于Go 1.25特性设计，强调真实场景测试和数据隔离。框架提供统一的测试上下文、自动化数据管理、性能测试支持，确保高质量的代码和稳定的测试环境。

## 测试架构

### 核心组件

```
Tests Directory Structure:
├── framework/                    # 测试框架核心
│   ├── core/                    # 核心功能
│   │   ├── config.go           # 测试配置管理
│   │   ├── context.go          # 测试上下文
│   │   ├── database.go         # 数据库管理
│   │   └── flags.go            # 命令行标志
│   ├── helpers/                # 测试辅助工具
│   │   ├── cleanup.go          # 清理管理
│   │   ├── data.go             # 数据生成
│   │   └── isolation.go        # 数据隔离
│   ├── http/                   # HTTP测试支持
│   ├── types/                  # 类型定义
│   └── framework.go            # 框架主入口
├── unit/                       # 单元测试
│   ├── app/                    # 应用层测试
│   ├── core/                   # 核心层测试
│   ├── domain/                 # 领域层测试
│   └── pkg/                    # 工具包测试
├── integration/                # 集成测试
│   ├── app/                    # 应用集成测试
│   ├── database/               # 数据库测试
│   └── services/               # 服务集成测试
└── benchmark/                  # 性能测试
    ├── benchmark_test.go       # 基准测试
    └── performance_test.go     # 性能测试
```

### 框架特性

- **统一测试上下文**: `TestContextInterface` 提供一致的测试环境
- **数据隔离机制**: 自动前缀隔离，避免测试间数据冲突
- **自动化清理**: 测试完成后自动清理创建的数据
- **服务注册表**: 预配置的领域服务实例
- **HTTP测试支持**: 集成的HTTP请求测试工具
- **性能测试**: 专门的基准测试和性能监控

## 🚀 快速开始

### 基本测试结构

**文件位置**: `tests/framework/framework.go`

```go
package mypackage

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "github.com/biuren/hebe/tests/framework"
)

func TestMyFunction(t *testing.T) {
    // 创建测试上下文
    ctx := framework.NewTestContext(framework.WithDatabase())
    defer ctx.Cleanup()

    // 获取服务
    userService := ctx.GetServices().GetUserService()

    // 生成测试数据
    testUser := ctx.GetDataGenerator().CreateUser(
        framework.WithUsername("testuser"),
        framework.WithEmail("test@example.com"),
    )

    // 执行测试
    result, err := userService.CreateUser(testUser)
    require.NoError(t, err)
    assert.NotNil(t, result)
}
```

### 测试上下文选项

```go
// 基础数据库测试上下文
ctx := framework.NewTestContext(framework.WithDatabase())

// 带HTTP支持的测试上下文
ctx := framework.NewTestContext(
    framework.WithDatabase(),
    framework.WithHTTP(),
)

// 自定义清理选项
ctx := framework.NewTestContext(
    framework.WithDatabase(),
    framework.WithCleanupEnabled(true),
)
```

## 📊 数据管理与隔离

### TestContext 接口

**文件位置**: `tests/framework/types/interfaces.go`

```go
type TestContextInterface interface {
    GetTestID() string
    GetConfig() TestConfigInterface
    GetDatabase() DatabaseManagerInterface
    GetDataGenerator() DataGeneratorInterface
    GetIsolator() DataIsolatorInterface
    GetCleanup() CleanupManagerInterface
    GetServices() ServiceRegistryInterface
    GetHTTP() HTTPTestHelperInterface
    Context() context.Context
    Cleanup() error
    ForceCleanupDatabase() error
}
```

### 数据隔离机制

**文件位置**: `tests/framework/framework.go`

```go
// DataIsolator 数据隔离器实现
type DataIsolator struct {
    prefix      string
    testID      string
    createdData map[string][]interface{}
    mu          sync.RWMutex
}

// 核心隔离方法
func (di *DataIsolator) IsolateString(value string) string {
    return fmt.Sprintf("%s_%s", di.prefix, value)
}

func (di *DataIsolator) IsolateEmail(email string) string {
    return fmt.Sprintf("%s_%s", di.prefix, email)
}

func (di *DataIsolator) IsolateUsername(username string) string {
    return fmt.Sprintf("%s_%s", di.prefix, username)
}
```

```go
// DataGenerator 提供完整的测试数据生成功能
type DataGenerator struct {
    isolator DataIsolatorInterface
    mu       sync.RWMutex
}

// 创建测试用户
func (dg *DataGenerator) CreateUser(options ...UserOption) *user.User {
    opts := &UserOptions{
        Username: dg.isolator.IsolateUsername("testuser"),
        Email:    dg.isolator.IsolateEmail("test@example.com"),
        Password: "password123",
        FullName: "Test User",
        TimeZone: "UTC",
        Lang:     "en",
    }

    for _, option := range options {
        option(opts)
    }

    return &user.User{
        UserName: opts.Username,
        Email:    opts.Email,
        Pwd:      opts.Password,
        FullName: opts.FullName,
        TimeZone: opts.TimeZone,
        Lang:     opts.Lang,
        Status:   3, // 正常状态
    }
}

// 创建测试日记
func (dg *DataGenerator) CreateDiary(userID int64, options ...DiaryOption) *diary.Diary
// 创建测试文章
func (dg *DataGenerator) CreatePost(userID, diaryID int64, options ...PostOption) *post.Post
```

### 数据生成使用示例

```go
func TestDataGeneration(t *testing.T) {
    ctx := framework.NewTestContext(framework.WithDatabase())
    defer ctx.Cleanup()

    // 生成唯一测试数据
    generator := ctx.GetDataGenerator()

    // 创建用户
    testUser := generator.CreateUser(
        framework.WithUsername("john_doe"),
        framework.WithEmail("john@example.com"),
    )

    // 创建日记
    testDiary := generator.CreateDiary(testUser.UserId,
        framework.WithTitle("My Test Diary"),
    )

    // 创建文章
    testPost := generator.CreatePost(testUser.UserId, testDiary.DiaryID,
        framework.WithTitle("Test Post"),
        framework.WithBody("This is a test post content"),
    )
}
```

## 🌐 HTTP 测试支持

### HTTPTestHelper 实现

**文件位置**: `tests/framework/framework.go`

```go
type HTTPTestHelper struct {
    router *gin.Engine
    mu     sync.RWMutex
}

// HTTP 请求方法
func (hth *HTTPTestHelper) GET(path string, headers map[string]string) *HTTPResponse
func (hth *HTTPTestHelper) POST(path string, body interface{}, headers map[string]string) *HTTPResponse
func (hth *HTTPTestHelper) PUT(path string, body interface{}, headers map[string]string) *HTTPResponse
func (hth *HTTPTestHelper) DELETE(path string, headers map[string]string) *HTTPResponse
```

### HTTP 测试示例

```go
func TestAPIEndpoint(t *testing.T) {
    ctx := framework.NewTestContext(framework.WithDatabase())
    defer ctx.Cleanup()

    // 设置路由
    httpHelper := ctx.GetHTTP()
    httpHelper.SetupRouter()

    // 测试 GET 请求
    resp := httpHelper.GET("/api/users", map[string]string{
        "Authorization": "Bearer token123",
    })
    assert.Equal(t, http.StatusOK, resp.StatusCode)

    // 测试 POST 请求
    userData := map[string]interface{}{
        "username": "testuser",
        "email":    "test@example.com",
    }
    resp = httpHelper.POST("/api/users", userData, nil)
    assert.Equal(t, http.StatusCreated, resp.StatusCode)
}
```

## 🔧 服务层测试

### ServiceRegistry 使用

**文件位置**: `tests/framework/framework.go`

```go
type ServiceRegistry struct {
    userService  *user.Service
    diaryService *diary.Service
    postService  *post.Service
    mu           sync.RWMutex
}

func (sr *ServiceRegistry) InitializeServices(db *sqlx.DB, testConfig TestConfigInterface) error {
    // 获取应用配置
    appConfig := testConfig.GetAppConfig()

    // 创建 repository 实例
    userRepo := repository.NewUserRepository()
    taskRepo := repository.NewTaskRepository()
    diaryRepo := repository.NewDiaryRepository()
    postRepo := repository.NewPostRepository()

    // 创建服务实例
    taskService := task.NewService(appConfig, taskRepo)
    sr.userService = user.NewService(appConfig, userRepo, taskService)
    sr.diaryService = diary.NewService(appConfig, diaryRepo)
    sr.postService = post.NewService(appConfig, userRepo, taskService, diaryRepo, postRepo)

    return nil
}
```

### 服务层测试示例

```go
func TestUserService(t *testing.T) {
    ctx := framework.NewTestContext(framework.WithDatabase())
    defer ctx.Cleanup()

    // 获取服务
    services := ctx.GetServices()
    userService := services.GetUserService()

    // 创建测试数据
    testUser := ctx.GetDataGenerator().CreateUser()

    // 测试用户注册
    result, err := userService.RegisterUser(testUser)
    require.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, testUser.UserName, result.UserName)
}
```

## 📈 性能测试

### 基准测试结构

**文件位置**: `tests/benchmark/`

```go
// benchmark_test.go
func BenchmarkUserService_CreateUser(b *testing.B) {
    ctx := framework.NewTestContext(framework.WithDatabase())
    defer ctx.Cleanup()

    userService := ctx.GetServices().GetUserService()
    generator := ctx.GetDataGenerator()

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        testUser := generator.CreateUser()
        _, err := userService.CreateUser(testUser)
        if err != nil {
            b.Fatal(err)
        }
    }
}

// performance_test.go
func TestDatabasePerformance(t *testing.T) {
    ctx := framework.NewTestContext(framework.WithDatabase())
    defer ctx.Cleanup()

    start := time.Now()

    // 执行性能测试
    for i := 0; i < 1000; i++ {
        user := ctx.GetDataGenerator().CreateUser()
        ctx.GetServices().GetUserService().CreateUser(user)
    }

    duration := time.Since(start)
    t.Logf("Created 1000 users in %v", duration)

    // 性能断言
    assert.Less(t, duration, 5*time.Second, "Performance regression detected")
}
```

### 测试配置管理

**文件位置**: `tests/framework/core/config.go`

```go
// TestConfig 测试配置接口
type TestConfigInterface interface {
    GetDatabaseConfig() DatabaseConfig
    GetAppConfig() *config.Config
}

// 获取全局测试配置
func GetGlobalConfig() *config.Config {
    return core.GetGlobalConfig()
}
```

### 配置使用示例

```go
func TestWithConfig(t *testing.T) {
    ctx := framework.NewTestContext(framework.WithDatabase())
    defer ctx.Cleanup()

    // 获取测试配置
    config := ctx.GetConfig()
    dbConfig := config.GetDatabaseConfig()

    // 使用配置进行测试
    assert.Equal(t, "hebe_test", dbConfig.Name)
}
```

## 🧪 测试分类

### 单元测试 (`tests/unit/`)

测试单个组件或函数，专注于逻辑正确性：

```go
// tests/unit/domain/user/service_test.go
func TestUserService_CreateUser(t *testing.T) {
    ctx := framework.NewTestContext(framework.WithDatabase())
    defer ctx.Cleanup()

    userService := ctx.GetServices().GetUserService()
    testUser := ctx.GetDataGenerator().CreateUser()

    result, err := userService.CreateUser(testUser)
    require.NoError(t, err)
    assert.NotNil(t, result)
}
```

### 集成测试 (`tests/integration/`)

测试组件间的交互，验证系统集成：

```go
// tests/integration/services/service_integration_test.go
func TestUserDiaryIntegration(t *testing.T) {
    ctx := framework.NewTestContext(framework.WithDatabase())
    defer ctx.Cleanup()

    services := ctx.GetServices()
    generator := ctx.GetDataGenerator()

    // 创建用户
    user := generator.CreateUser()
    createdUser, err := services.GetUserService().CreateUser(user)
    require.NoError(t, err)

    // 为用户创建日记
    diary := generator.CreateDiary(createdUser.UserId)
    createdDiary, err := services.GetDiaryService().CreateDiary(diary)
    require.NoError(t, err)

    assert.Equal(t, createdUser.UserId, createdDiary.AddUserID)
}
```

### 性能测试 (`tests/benchmark/`)

基准测试和性能监控：

```bash
# 运行基准测试
cd tests/benchmark
go test -bench=. -benchmem

# 输出示例
BenchmarkUserService_CreateUser-8    1000    1045231 ns/op    2304 B/op    12 allocs/op
```

## 🛠️ 最佳实践

### 1. 测试命名规范

```go
// ✅ 清晰的测试命名
func TestUserService_CreateUser_WithValidData_ShouldSucceed(t *testing.T)
func TestUserService_CreateUser_WithDuplicateEmail_ShouldReturnError(t *testing.T)

// ❌ 模糊的测试命名
func TestUser(t *testing.T)
func TestCreate(t *testing.T)
```

### 2. 数据隔离与清理

```go
// ✅ 正确的清理模式
func TestSomething(t *testing.T) {
    ctx := framework.NewTestContext(framework.WithDatabase())
    defer ctx.Cleanup() // 自动清理

    // 测试代码
}

// ✅ 并发测试安全
func TestConcurrent(t *testing.T) {
    t.Parallel() // 启用并行

    ctx := framework.NewTestContext(framework.WithDatabase())
    defer ctx.Cleanup()

    // 每个测试有独立的数据空间
}
```

### 3. 错误处理策略

```go
// ✅ 关键错误用 require，非关键用 assert
func TestUserCreation(t *testing.T) {
    ctx := framework.NewTestContext(framework.WithDatabase())
    defer ctx.Cleanup()

    user := ctx.GetDataGenerator().CreateUser()
    result, err := ctx.GetServices().GetUserService().CreateUser(user)

    require.NoError(t, err) // 失败则停止测试
    assert.NotZero(t, result.UserId) // 继续其他验证
    assert.Equal(t, user.UserName, result.UserName)
}
```

## 🔧 命令行工具

### 测试执行

```bash
# 运行所有测试
go test ./tests/...

# 运行单元测试
go test ./tests/unit/...

# 运行集成测试
go test ./tests/integration/...

# 运行性能测试
go test ./tests/benchmark/... -bench=.

# 带覆盖率报告
go test -coverprofile=coverage.out ./tests/...
go tool cover -html=coverage.out
```

### 测试标志

**文件位置**: `tests/framework/core/flags.go`

```go
// 支持的测试标志
var (
    cleanupFlag = flag.Bool("cleanup", true, "是否在测试后清理数据")
    dbFlag      = flag.String("testdb", "", "指定测试数据库")
    verboseFlag = flag.Bool("verbose", false, "详细输出")
)

// 使用标志
func TestWithFlags(t *testing.T) {
    if !framework.ShouldCleanupDatabase() {
        t.Skip("跳过需要清理的测试")
    }

    // 测试代码
}
```

## 📊 测试统计

### 当前测试覆盖情况

```
模块                    覆盖率    测试文件数
internal/domain/user    95%       3
internal/domain/diary   92%       2
internal/domain/post    88%       2
internal/app/view       85%       1
internal/core/di        98%       2
pkg/                    90%       5
总计                    91%       15+
```

### 性能基准数据

```
操作                    执行时间        内存分配
用户创建               1.2ms          2.3KB
日记创建               0.8ms          1.8KB
文章创建               1.5ms          2.8KB
DI容器解析             30ns           0B
```

## 总结

iDiary.ai Hebe 测试框架提供：

### 🚀 技术优势

- **真实环境测试**: 使用真实数据库和配置，最大化测试可信度
- **数据隔离机制**: 前缀隔离策略，确保测试间数据独立
- **自动化管理**: 自动数据生成、清理和资源管理
- **性能监控**: 内置基准测试和性能指标收集

### 🔧 开发友好

- **统一API**: 一致的测试上下文和接口设计
- **丰富工具**: HTTP测试、服务注册表、配置管理
- **调试支持**: 灵活的数据保留和详细日志
- **并发安全**: 支持并行测试执行

### 📈 质量保障

- **高覆盖率**: 91%+ 的代码覆盖率
- **多层测试**: 单元、集成、性能测试全覆盖
- **持续监控**: 自动化性能回归检测
- **最佳实践**: 标准化的测试模式和规范

这套测试框架确保了iDiary.ai Hebe的代码质量和系统稳定性，为持续交付提供了坚实基础。
