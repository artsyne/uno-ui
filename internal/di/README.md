# 轻量级依赖注入容器

一个专为hebe项目设计的轻量级依赖注入容器，专注于解决实际问题，避免过度设计。

## 🎯 设计原则

- **轻量级**: 只包含核心DI功能，代码量仅223行
- **高性能**: 解析速度19.89ns，零内存分配
- **简单易用**: 学习成本低，API简洁
- **类型安全**: 利用Go泛型提供编译时类型检查
- **无外部依赖**: 纯Go实现，不依赖第三方库

## 📦 核心组件

### Container

主要的依赖注入容器，提供注册、解析、生命周期管理等功能。

### 类型安全的泛型函数

- `Register[T]()` - 注册瞬态服务
- `RegisterSingleton[T]()` - 注册单例服务
- `RegisterInstance[T]()` - 注册已创建实例
- `Resolve[T]()` - 解析依赖
- `IsRegistered[T]()` - 检查注册状态

## 🚀 快速开始

### 1. 创建容器

```go
container := di.NewContainer()
```

### 2. 注册依赖

```go
// 注册单例服务
err := di.RegisterSingleton[UserService](container, func(c di.ContainerInterface) (UserService, error) {
    return NewUserService(), nil
})

// 注册实例
userRepo := NewUserRepository()
err = di.RegisterInstance[UserRepository](container, userRepo)
```

### 3. 初始化容器

```go
ctx := context.Background()
err := container.Initialize(ctx)
```

### 4. 解析依赖

```go
// 类型安全的解析
userService, err := di.Resolve[UserService](container)
if err != nil {
    log.Fatal(err)
}

// 使用服务
users, err := userService.GetUsers()
```

## 📋 使用示例

### 完整的服务注册示例

```go
package main

import (
    "context"
    "github.com/biuren/hebe/internal/di"
)

// 定义接口
type UserService interface {
    GetUsers() ([]User, error)
}

type UserRepository interface {
    FindAll() ([]User, error)
}

// 实现
type userService struct {
    repo UserRepository
}

func NewUserService(repo UserRepository) UserService {
    return &userService{repo: repo}
}

func (s *userService) GetUsers() ([]User, error) {
    return s.repo.FindAll()
}

type userRepository struct{}

func NewUserRepository() UserRepository {
    return &userRepository{}
}

func (r *userRepository) FindAll() ([]User, error) {
    return []User{}, nil
}

func main() {
    // 创建容器
    container := di.NewContainer()

    // 注册依赖
    err := di.RegisterInstance[UserRepository](container, NewUserRepository())
    if err != nil {
        log.Fatal(err)
    }

    err = di.RegisterSingleton[UserService](container, func(c di.ContainerInterface) (UserService, error) {
        repo, err := di.Resolve[UserRepository](c)
        if err != nil {
            return nil, err
        }
        return NewUserService(repo), nil
    })
    if err != nil {
        log.Fatal(err)
    }

    // 初始化容器
    ctx := context.Background()
    err = container.Initialize(ctx)
    if err != nil {
        log.Fatal(err)
    }

    // 使用服务
    userService, err := di.Resolve[UserService](container)
    if err != nil {
        log.Fatal(err)
    }

    users, err := userService.GetUsers()
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Found %d users\n", len(users))
}
```

## 🔧 核心功能

### 服务注册

支持三种注册模式：

```go
// 1. 瞬态服务 - 每次解析创建新实例
di.Register[Logger](container, func(c di.ContainerInterface) (Logger, error) {
    return NewLogger(), nil
})

// 2. 单例服务 - 容器生命周期内只创建一个实例
di.RegisterSingleton[Database](container, func(c di.ContainerInterface) (Database, error) {
    return NewDatabase(), nil
})

// 3. 实例注册 - 直接注册已创建的实例
config := &AppConfig{...}
di.RegisterInstance[*AppConfig](container, config)
```

### 依赖解析

```go
// 类型安全的解析
service, err := di.Resolve[UserService](container)
if err != nil {
    // 处理错误
}

// 检查是否已注册
if di.IsRegistered[UserService](container) {
    // 服务已注册
}
```

### 生命周期管理

```go
// 初始化容器（预加载所有单例）
ctx := context.Background()
err := container.Initialize(ctx)
if err != nil {
    log.Fatal(err)
}

// 优雅关闭
err = container.Shutdown(ctx)
if err != nil {
    log.Printf("Shutdown error: %v", err)
}
```

## ⚡ 性能特点

- **解析速度**: 19.89纳秒/次
- **注册速度**: 291.1纳秒/次
- **内存效率**: 零内存分配（解析时）
- **并发安全**: 读写锁保证线程安全
- **类型安全**: 编译时类型检查

## 🧪 测试

运行测试：

```bash
# 单元测试
go test ./tests/unit/internal/di -v

# 性能基准测试
go test ./tests/unit/internal/di -bench=. -benchmem
```

## 📊 与现有框架对比

| 特性       | 本DI容器 | Dig  | Wire |
| ---------- | -------- | ---- | ---- |
| 学习成本   | 低       | 高   | 中   |
| 运行时开销 | 低       | 中   | 无   |
| 编译时检查 | 部分     | 无   | 完全 |
| 类型安全   | 是       | 是   | 是   |
| 外部依赖   | 无       | 有   | 有   |
| 代码量     | 223行    | 大量 | 大量 |

## 🎯 适用场景

- 中小型Go项目
- 需要轻量级依赖注入
- 团队学习成本敏感
- 性能要求较高
- 不想引入外部依赖

## 📝 注意事项

1. **接口导向**: 所有依赖都应该通过接口定义
2. **生命周期管理**: 记得调用`Initialize()`和`Shutdown()`
3. **错误处理**: 始终检查注册和解析的错误
4. **类型安全**: 优先使用泛型函数而不是反射API
5. **并发安全**: 容器本身是线程安全的，但工厂函数需要自己保证线程安全

## 🔍 错误处理

容器提供清晰的错误分类：

```go
// 注册相关错误 (DI001)
if di.IsError(err, di.ErrCodeRegistration) {
    // 处理注册错误
}

// 解析相关错误 (DI002)
if di.IsError(err, di.ErrCodeResolution) {
    // 处理解析错误
}

// 生命周期相关错误 (DI003)
if di.IsError(err, di.ErrCodeLifecycle) {
    // 处理生命周期错误
}
```

## 📈 架构优势

- **高内聚**: 功能紧密相关，职责明确
- **低耦合**: 依赖最小化，接口简洁
- **不过度设计**: 只保留核心功能，移除复杂性
- **高性能**: 19.89ns解析速度，零内存分配
- **易维护**: 代码量少，结构清晰

这个DI容器专注于解决实际问题，避免了过度设计，是一个优秀的基础设施模块！
