# 测试框架使用指南

## 概述

Hebe 项目采用统一的测试框架，提供了一致的测试体验和强大的功能支持。本指南将帮助您快速上手并有效使用测试框架。

## 🚀 快速开始

### 基本测试结构

```go
package mypackage

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    framework "github.com/biuren/hebe/tests/framework/core"
)

func TestMyFunction(t *testing.T) {
    // 创建测试上下文
    ctx := framework.NewTestContext("TestMyFunction", framework.WithDatabase())
    defer ctx.Cleanup()
    
    // 您的测试代码
    // ...
}
```

### 测试上下文选项

```go
// 仅数据库支持
ctx := framework.NewTestContext("TestName", framework.WithDatabase())

// 数据库 + HTTP 支持
ctx := framework.NewTestContext("TestName", 
    framework.WithDatabase(),
    framework.WithHTTP(),
)

// 自定义配置
ctx := framework.NewTestContext("TestName",
    framework.WithDatabase(),
    framework.WithConfig(customConfig),
)
```

## 📊 数据管理

### 数据生成器

测试框架提供了强大的数据生成器，确保测试数据的唯一性：

```go
func TestUserCreation(t *testing.T) {
    ctx := framework.NewTestContext("TestUserCreation", framework.WithDatabase())
    defer ctx.Cleanup()
    
    // 生成唯一的测试数据
    testUser := &user.User{
        UserName: ctx.GetDataGenerator().GenerateUniqueString("testuser"),
        Email:    ctx.GetDataGenerator().GenerateUniqueEmail("test"),
        FullName: ctx.GetDataGenerator().GenerateUniqueString("fullname"),
        // ...
    }
    
    // 使用生成的数据进行测试
    userService := user.NewService()
    err := userService.RegisterUser(testUser)
    require.NoError(t, err)
}
```

### 数据隔离

框架使用前缀策略实现数据隔离，每个测试都有独立的数据空间：

```go
// 自动数据隔离 - 推荐方式
ctx := framework.NewTestContext("TestName", framework.WithDatabase())
defer ctx.Cleanup() // 自动清理测试数据

// 手动控制清理
ctx := framework.NewTestContext("TestName", framework.WithDatabase())
// 测试代码...
ctx.Cleanup() // 手动清理
```

## 🌐 HTTP 测试

### HTTP 客户端使用

```go
func TestAPIEndpoint(t *testing.T) {
    ctx := framework.NewTestContext("TestAPIEndpoint", 
        framework.WithDatabase(),
        framework.WithHTTP(),
    )
    defer ctx.Cleanup()
    
    httpHelper := ctx.GetHTTPHelper()
    
    // GET 请求
    resp := httpHelper.GET("/api/users", nil)
    assert.Equal(t, 200, resp.StatusCode)
    
    // POST 请求
    userData := map[string]interface{}{
        "username": "testuser",
        "email":    "test@example.com",
    }
    resp = httpHelper.POST("/api/users", userData, nil)
    assert.Equal(t, 201, resp.StatusCode)
}
```

## 🔧 配置管理

### 获取配置

```go
func TestWithConfig(t *testing.T) {
    ctx := framework.NewTestContext("TestWithConfig", framework.WithDatabase())
    defer ctx.Cleanup()
    
    // 获取测试配置
    config := ctx.GetConfig()
    dbConfig := config.GetDatabaseConfig()
    
    // 使用配置
    assert.Equal(t, "hebe_test", dbConfig.Name)
}
```

### 自定义配置

```go
func TestCustomConfig(t *testing.T) {
    customConfig := &framework.TestConfig{
        Database: framework.DatabaseConfig{
            Name: "custom_test_db",
            // 其他配置...
        },
    }
    
    ctx := framework.NewTestContext("TestCustomConfig",
        framework.WithDatabase(),
        framework.WithConfig(customConfig),
    )
    defer ctx.Cleanup()
    
    // 使用自定义配置的测试
}
```

## 📈 性能基准测试

### 编写基准测试

```go
func BenchmarkUserRegistration(b *testing.B) {
    ctx := framework.NewTestContext("BenchmarkUserRegistration", framework.WithDatabase())
    defer ctx.Cleanup()
    
    userService := user.NewService()
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        testUser := &user.User{
            UserName: ctx.GetDataGenerator().GenerateUniqueString("benchuser"),
            Email:    ctx.GetDataGenerator().GenerateUniqueEmail("bench"),
            // ...
        }
        err := userService.RegisterUser(testUser)
        if err != nil {
            b.Fatal(err)
        }
    }
}
```

### 运行基准测试

```bash
# 运行所有基准测试
./scripts/benchmark.sh

# 运行特定类型的基准测试
./scripts/benchmark.sh user
./scripts/benchmark.sh concurrent

# 查看帮助
./scripts/benchmark.sh help
```

## 🧪 测试类型

### 单元测试

位于 `tests/unit/` 目录，测试单个函数或方法：

```go
func TestUserValidation(t *testing.T) {
    ctx := framework.NewTestContext("TestUserValidation")
    defer ctx.Cleanup()
    
    // 测试用户验证逻辑
    user := &user.User{Email: "invalid-email"}
    err := user.Validate()
    assert.Error(t, err)
}
```

### 集成测试

位于 `tests/integration/` 目录，测试组件间的交互：

```go
func TestUserServiceIntegration(t *testing.T) {
    ctx := framework.NewTestContext("TestUserServiceIntegration", framework.WithDatabase())
    defer ctx.Cleanup()
    
    // 测试用户服务与数据库的集成
    userService := user.NewService()
    // 集成测试代码...
}
```

## 🛠️ 最佳实践

### 1. 测试命名

```go
// ✅ 好的命名
func TestUserService_RegisterUser_WithValidData_ShouldSucceed(t *testing.T) {}
func TestUserService_RegisterUser_WithDuplicateEmail_ShouldReturnError(t *testing.T) {}

// ❌ 避免的命名
func TestUser(t *testing.T) {}
func TestRegister(t *testing.T) {}
```

### 2. 数据清理

```go
// ✅ 推荐：使用 defer 确保清理
func TestSomething(t *testing.T) {
    ctx := framework.NewTestContext("TestSomething", framework.WithDatabase())
    defer ctx.Cleanup() // 确保清理
    
    // 测试代码...
}

// ❌ 避免：忘记清理
func TestSomething(t *testing.T) {
    ctx := framework.NewTestContext("TestSomething", framework.WithDatabase())
    // 测试代码...
    // 忘记调用 ctx.Cleanup()
}
```

### 3. 错误处理

```go
// ✅ 使用 require 处理关键错误
func TestUserCreation(t *testing.T) {
    ctx := framework.NewTestContext("TestUserCreation", framework.WithDatabase())
    defer ctx.Cleanup()
    
    userService := user.NewService()
    err := userService.RegisterUser(testUser)
    require.NoError(t, err) // 如果失败，立即停止测试
    
    // 继续其他断言
    assert.NotZero(t, testUser.UserId)
}
```

### 4. 子测试

```go
func TestUserService(t *testing.T) {
    ctx := framework.NewTestContext("TestUserService", framework.WithDatabase())
    defer ctx.Cleanup()
    
    userService := user.NewService()
    
    t.Run("注册新用户", func(t *testing.T) {
        // 子测试代码
    })
    
    t.Run("重复邮箱注册", func(t *testing.T) {
        // 子测试代码
    })
}
```

## 🔍 调试技巧

### 1. 保留测试数据

```go
func TestDebugMode(t *testing.T) {
    ctx := framework.NewTestContext("TestDebugMode", framework.WithDatabase())
    // 注释掉 defer ctx.Cleanup() 来保留数据进行调试
    // defer ctx.Cleanup()
    
    // 测试代码...
    
    // 手动清理（可选）
    // ctx.Cleanup()
}
```

### 2. 查看生成的数据

```go
func TestWithLogging(t *testing.T) {
    ctx := framework.NewTestContext("TestWithLogging", framework.WithDatabase())
    defer ctx.Cleanup()
    
    username := ctx.GetDataGenerator().GenerateUniqueString("testuser")
    t.Logf("生成的用户名: %s", username)
    
    // 测试代码...
}
```

## 📚 常见问题

### Q: 如何处理数据库连接问题？

A: 确保数据库配置正确，并且测试数据库可访问：

```go
// 检查数据库连接
func TestDatabaseConnection(t *testing.T) {
    ctx := framework.NewTestContext("TestDatabaseConnection", framework.WithDatabase())
    defer ctx.Cleanup()
    
    db := ctx.GetDatabase()
    err := db.Ping()
    require.NoError(t, err, "数据库连接失败")
}
```

### Q: 如何处理并发测试？

A: 使用数据隔离确保并发安全：

```go
func TestConcurrent(t *testing.T) {
    t.Parallel() // 启用并行测试
    
    ctx := framework.NewTestContext("TestConcurrent", framework.WithDatabase())
    defer ctx.Cleanup()
    
    // 每个测试都有独立的数据空间
    // 测试代码...
}
```

### Q: 如何测试错误情况？

A: 使用断言验证预期的错误：

```go
func TestErrorHandling(t *testing.T) {
    ctx := framework.NewTestContext("TestErrorHandling", framework.WithDatabase())
    defer ctx.Cleanup()
    
    userService := user.NewService()
    
    // 测试无效输入
    invalidUser := &user.User{Email: "invalid"}
    err := userService.RegisterUser(invalidUser)
    assert.Error(t, err, "应该返回验证错误")
    assert.Contains(t, err.Error(), "email", "错误信息应包含email相关内容")
}
```

## 🎯 总结

测试框架提供了：

- ✅ **统一的测试体验** - 一致的API和模式
- ✅ **自动数据管理** - 数据生成、隔离和清理
- ✅ **强大的工具支持** - HTTP客户端、配置管理等
- ✅ **性能测试支持** - 内置基准测试工具
- ✅ **调试友好** - 灵活的数据保留和日志

遵循本指南的最佳实践，您可以编写出高质量、可维护的测试代码。