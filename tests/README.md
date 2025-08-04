# Hebe 项目测试框架文档

## 📋 目录

- [项目概况](#项目概况)
- [测试框架设计](#测试框架设计)
- [目录结构](#目录结构)
- [快速开始](#快速开始)
- [测试类型](#测试类型)
- [框架功能](#框架功能)
- [最佳实践](#最佳实践)
- [故障排除](#故障排除)

## 📊 项目概况

**Hebe** 是一个基于 Go 和 Gin 框架的日记应用系统，采用分层架构设计。本测试框架为项目提供了统一、高效的测试解决方案。

### 技术栈

- **Web框架**: Gin
- **数据库**: MySQL (使用 sqlx)
- **测试数据库**: hebe_test (默认配置)
- **依赖注入**: 自实现的泛型DI容器
- **配置管理**: Viper + 环境变量
- **日志**: slog
- **模板引擎**: Go原生template
- **国际化**: 自实现的i18n系统

### 架构分层

```
├── cmd/idiary/           # 应用入口
├── internal/
│   ├── app/             # 应用层 (controller, middleware, route, view)
│   ├── core/            # 核心层 (database, di, initializer)
│   ├── domain/          # 领域层 (diary, post, task, user)
│   └── utils/           # 工具层 (config, error, http, i18n, validation)
├── pkg/                 # 公共包
└── website/             # 前端资源
```

## 🎯 测试框架设计

### 设计目标

新测试框架旨在解决以下核心问题：

1. **统一测试模式** - 消除不同测试初始化方式的差异
2. **简化工具链** - 合并重复的helper文件为统一工具集
3. **提升性能** - 去除事务依赖，实现并发测试执行
4. **改善数据管理** - 用唯一前缀替代事务，实现灵活的数据隔离
5. **优化结构** - 建立清晰的测试组织结构

### 清理机制

- **数据隔离**: 基于唯一前缀的数据隔离机制
- **动态数据生成**: 智能的测试数据生成器
- **HTTP测试助手**: 完整的HTTP测试工具链
- **性能监控**: 内置的性能基准测试
- **清理机制**: 自动化的测试数据清理
- **强制清理**: 无论配置如何都可强制清理数据库

## 📁 目录结构

```
tests/
├── framework/                    # 统一测试框架
│   ├── core/                    # 核心组件
│   │   ├── context.go           # 测试上下文管理
│   │   ├── config.go            # 统一配置管理
│   │   └── database.go          # 数据库连接池管理
│   ├── helpers/                 # 辅助工具
│   │   ├── data.go              # 测试数据生成器
│   │   ├── isolation.go         # 数据隔离器
│   │   └── cleanup.go           # 清理管理器
│   ├── http/                    # HTTP测试工具
│   │   └── helper.go            # HTTP测试助手
│   └── types/                   # 类型定义
│       ├── interfaces.go        # 测试接口
│       └── structs.go           # 测试结构体
├── unit/                        # 单元测试
│   ├── core/                    # 核心层测试
│   ├── app/                     # 应用层测试
│   ├── domain/                  # 领域层测试
│   └── pkg/                     # 工具包测试
├── integration/                 # 集成测试
│   ├── app/                     # 应用集成测试
│   ├── database/                # 数据库集成测试
│   └── services/                # 服务集成测试
├── benchmark/                   # 性能基准测试
│   ├── performance_test.go      # 性能测试用例
│   └── website/                 # 前端资源测试
└── framework_validation_test.go # 框架验证测试
```

## 🚀 快速开始

### 1. 基本测试用例

```go
package tests

import (
    "testing"
    framework "github.com/biuren/hebe/tests/framework/core"
    "github.com/stretchr/testify/assert"
)

func TestExample(t *testing.T) {
    // 创建测试上下文
    ctx := framework.NewTestContext("example_test")
    defer ctx.Cleanup()
    
    // 获取数据生成器
    generator := ctx.GetDataGenerator()
    
    // 生成测试用户
    user := generator.CreateUser()
    assert.NotNil(t, user)
    assert.NotEmpty(t, user.UserName)
}
```

### 2. 数据隔离测试

```go
func TestDataIsolation(t *testing.T) {
    ctx := framework.NewTestContext("isolation_test")
    defer ctx.Cleanup()
    
    isolator := ctx.GetIsolator()
    
    // 隔离字符串
    isolated := isolator.IsolateString("test")
    assert.Contains(t, isolated, "test")
    assert.NotEqual(t, "test", isolated)
}
```

### 3. 数据库测试

```go
func TestDatabase(t *testing.T) {
    ctx := framework.NewTestContext("db_test")
    defer ctx.Cleanup()
    
    db := ctx.GetDatabase().GetDB()
    
    // 执行数据库操作
    var count int
    err := db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
    assert.NoError(t, err)
}
```

## 🧪 测试类型

### 单元测试 (Unit Tests)

- **位置**: `tests/unit/`
- **目的**: 测试单个函数或方法
- **特点**: 快速、独立、无外部依赖

### 集成测试 (Integration Tests)

- **位置**: `tests/integration/`
- **目的**: 测试组件间的交互
- **特点**: 涉及数据库、HTTP请求等

### 性能基准测试 (Benchmark Tests)

- **位置**: `tests/benchmark/`
- **目的**: 测量性能指标
- **运行**: `go test -bench=. ./tests/benchmark/...`

### 框架验证测试 (Framework Validation)

- **位置**: `tests/framework_validation_test.go`
- **目的**: 验证测试框架本身的功能
- **特点**: 确保测试工具的可靠性

## ⚙️ 框架功能

### TestContext (测试上下文)

```go
type TestContext interface {
    GetTestID() string
    GetConfig() ConfigManagerInterface
    GetDatabase() DatabaseManagerInterface
    GetIsolator() DataIsolatorInterface
    GetDataGenerator() DataGeneratorInterface
    Cleanup() error
}
```

### DataGenerator (数据生成器)

```go
type DataGeneratorInterface interface {
    CreateUser() *structs.User
    CreateDiary(userID int64) *structs.Diary
    CreatePost(diaryID, userID int64) *structs.Post
}
```

### DataIsolator (数据隔离器)

```go
type DataIsolatorInterface interface {
    IsolateString(input string) string
    IsolateEmail(email string) string
    IsolateUsername(username string) string
    IsolateTitle(title string) string
}
```

## 📝 最佳实践

### 1. 测试命名

- 使用描述性的测试名称
- 遵循 `TestFunctionName_Scenario_ExpectedResult` 格式
- 使用中文注释说明测试目的

### 2. 数据管理

- 始终使用 `defer ctx.Cleanup()` 清理测试数据
- 使用数据隔离器避免测试间的数据冲突
- 优先使用数据生成器创建测试数据

### 3. 并发测试

- 在适当的地方使用 `t.Parallel()`
- 确保测试间的数据隔离
- 避免共享可变状态

### 4. 错误处理

- 使用 `assert` 包进行断言
- 提供有意义的错误消息
- 测试正常和异常情况

## 🧹 数据库清理功能

### 清理方式

#### 1. 自动清理（推荐）
```go
func TestExample(t *testing.T) {
    ctx := framework.NewTestContext("test_example")
    defer ctx.Cleanup() // 测试结束时自动清理
    
    // 测试代码...
}
```

#### 2. 强制清理数据库
```go
func TestForceCleanup(t *testing.T) {
    ctx := framework.NewTestContext("test_force_cleanup")
    
    // 创建测试数据
    user := ctx.GetDataGenerator().CreateUser()
    
    // 强制清理数据库（无论配置如何）
    err := ctx.ForceCleanupDatabase()
    assert.NoError(t, err)
    
    // 继续测试...
}
```

#### 3. 手动控制清理时机
```go
func TestManualCleanup(t *testing.T) {
    ctx := framework.NewTestContext("test_manual")
    defer ctx.Cleanup()
    
    // 第一阶段测试
    user1 := ctx.GetDataGenerator().CreateUser()
    
    // 中间清理
    ctx.ForceCleanupDatabase()
    
    // 第二阶段测试
    user2 := ctx.GetDataGenerator().CreateUser()
}
```

#### 4. 禁用自动清理（用于调试）
```go
func TestWithoutAutoCleanup(t *testing.T) {
    // 禁用自动清理，数据保留用于调试
    ctx := framework.NewTestContext("debug_test", 
        framework.WithCleanupEnabled(false))
    
    // 测试代码...
    
    // 需要时手动清理
    ctx.ForceCleanupDatabase()
}
```

### 清理策略

- **prefix**: 按前缀清理（默认）
- **tracked**: 只清理跟踪的数据
- **immediate**: 立即清理
 
### 命令行参数控制

测试框架支持标准的 `go test` 语法传递参数：

```bash
# 使用标准 go test 语法，需要在测试包中添加 TestMain
go test ./tests/framework/examples/ -v -cleandb
go test ./tests/framework/examples/ -v -cleandb -timeout=60
go test ./tests/... -v -cleandb -parallel=false
```

#### 支持的参数：
- `-cleandb` - 强制清理数据库
- `-verbose` - 详细输出
- `-parallel=false` - 禁用并行执行
- `-timeout=60` - 设置超时时间(秒)

### 性能基准测试

```bash
# 运行性能基准测试
./tests/benchmark.sh

# 运行特定类型的基准测试
./tests/benchmark.sh user      # 用户相关测试
./tests/benchmark.sh database  # 数据库操作测试
./tests/benchmark.sh concurrent # 并发测试
```

### 命令行参数控制

测试框架现在支持以下自定义参数：

- `-cleardb`: 强制清理数据库
- `-verbose`: 详细输出
- `-parallel=false`: 禁用并行执行
- `-timeout=60`: 设置超时时间(秒)

使用方法：
```bash
# 编译测试二进制文件并传递自定义参数
go test -c ./tests/framework/examples/ -o test-binary
./test-binary -test.v -test.run TestCommandLineFlagsExample -- -cleardb -verbose -parallel=false -timeout=60
rm -f test-binary

# 运行特定测试并强制清理
go test -c ./tests/integration/ -o integration-test
./integration-test -test.v -- -cleardb
rm -f integration-test

# 运行所有测试并强制清理
go test -c ./tests/... -o all-tests
./all-tests -test.v -- -cleardb
rm -f all-tests
```

## 🗄️ 数据库配置

### 测试数据库设置

测试框架使用独立的测试数据库，默认配置如下：

```bash
# 测试数据库配置（在 .env.local 或环境变量中设置）
IDIARY_DB_HOST=127.0.0.1
IDIARY_DB_PORT=3306
IDIARY_DB_NAME=hebe_test        # 测试专用数据库
IDIARY_DB_USER=root
IDIARY_DB_PASSWORD=123123
```

### 数据库初始化

测试框架会自动：
1. 检查测试数据库是否存在
2. 如果不存在，自动创建 `hebe_test` 数据库
3. 加载数据库架构（从 `sql/schema.sql`）
4. 为每个测试提供独立的数据隔离

### 数据库权限要求

确保数据库用户具有以下权限：
- CREATE DATABASE（创建测试数据库）
- DROP DATABASE（清理测试数据库）
- CREATE TABLE（创建测试表）
- INSERT, UPDATE, DELETE, SELECT（数据操作）

## 🔧 故障排除

### 常见问题

1. **数据库连接失败**
   - 检查 `.env.local` 配置
   - 确认数据库服务运行状态
   - 验证数据库权限
   - 确保测试数据库 `hebe_test` 存在或用户有创建权限

2. **测试数据冲突**
   - 确保使用数据隔离器
   - 检查清理逻辑是否正确执行
   - 避免硬编码测试数据

3. **性能测试不稳定**
   - 在系统空闲时运行
   - 多次运行取平均值
   - 检查系统资源使用情况

4. **清理失败**
   - 检查数据库权限
   - 查看清理日志
   - 使用 `ForceCleanupDatabase()` 强制清理

### 调试技巧

- 使用 `t.Logf()` 输出调试信息
- 启用详细测试输出: `go test -v`
- 运行单个测试: `go test -run TestName`
- 查看覆盖率: `go test -cover`
- 禁用自动清理保留数据用于调试

## 📊 测试命令

### 基本测试命令

```bash
# 运行所有测试
go test ./tests/...

# 运行单元测试
go test ./tests/unit/...

# 运行集成测试
go test ./tests/integration/...

# 启用详细输出
go test ./tests/... -v
```

### 数据库清理测试

```bash
# 使用标准语法
go test ./tests/framework/examples/ -v -cleandb
go test ./tests/integration/... -v -cleandb
go test ./tests/... -v -cleandb -timeout=60
```

### 性能和覆盖率测试

```bash
# 运行性能基准测试
./tests/benchmark.sh

# 生成覆盖率报告
go test -coverprofile=coverage.out ./tests/...
go tool cover -html=coverage.out

# 运行特定测试
go test ./tests/unit/core/... -v -run TestConfig
go test ./tests/integration/... -v -cleandb -run TestUserService
```

### 高级用法

```bash
# 禁用并行测试
go test ./tests/... -v -cleandb -parallel=false

# 设置自定义超时
go test ./tests/... -v -cleandb -timeout=120

# 调试模式（保留测试数据）
go test ./tests/integration/... -v  # 不使用 -cleandb
```

---

*本文档整合了项目测试分析和框架设计方案，为开发者提供完整的测试指南。*