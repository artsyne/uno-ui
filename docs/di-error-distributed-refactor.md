# DI错误处理重构总结 - 分散错误函数到相关文件

## 重构目标

1. 将DI模块的独立错误处理机制统一到 `core/error` 包
2. **移除独立的 `errors.go` 文件**，将错误函数分散到相关代码文件中
3. 实现统一的错误处理和更好的代码组织

## 重构内容

### 1. 新增DI专用错误码

在 `internal/core/error/error.go` 中新增DI相关错误码：

```go
// 6250-6269: DI容器错误 (400/500)
DIServiceAlreadyRegistered = 6251 // DI服务已注册
DIServiceNotRegistered     = 6252 // DI服务未注册
DIInvalidFactory           = 6253 // DI工厂函数无效
DIContainerNotInitialized  = 6254 // DI容器未初始化
DITypeConversionError      = 6255 // DI类型转换错误
```

### 2. 错误函数分散到相关文件

**优化前**：

- 独立的 `internal/core/di/errors.go` 文件
- 统一的错误创建函数

**优化后**：

- **删除了 `errors.go` 文件**
- 错误函数分散到具体使用的文件中：
  - `container.go`: `newAlreadyRegisteredError`, `newNotRegisteredError`
  - `providers.go`: `newInvalidFactoryError`
  - 其他错误直接内联使用

### 3. 内联错误处理

对于使用频率低的错误（如初始化错误），直接在使用处内联：

```go
// 原来
return NewAlreadyInitializedError()

// 现在
return errors.DIContainerNotInitializedErr().WithDetails("container is already initialized")
```

## 文件结构变化

### 优化前

```
internal/core/di/
├── container.go
├── errors.go       # 独立错误文件 ❌
├── interfaces.go
└── providers.go
```

### 优化后 ✅

```
internal/core/di/
├── container.go    # 包含注册相关错误函数
├── interfaces.go   # 包含内联类型转换错误
└── providers.go    # 包含工厂函数相关错误函数
```

## 代码分布策略

### container.go

```go
// newAlreadyRegisteredError 创建已注册错误
func newAlreadyRegisteredError(serviceType reflect.Type) error {
    return errors.DIServiceAlreadyRegisteredErr().WithDetails(...)
}

// newNotRegisteredError 创建未注册错误
func newNotRegisteredError(serviceType reflect.Type) error {
    return errors.DIServiceNotRegisteredErr().WithDetails(...)
}

// 内联初始化相关错误
return errors.DIContainerNotInitializedErr().WithDetails("container is already initialized")
```

### providers.go

```go
// newInvalidFactoryError 创建无效工厂函数错误
func newInvalidFactoryError(factoryType reflect.Type) error {
    return errors.DIInvalidFactoryErr().WithDetails(...)
}
```

### interfaces.go

```go
// 内联类型转换错误
return errors.DITypeConversionErr().WithDetails(
    fmt.Sprintf("type conversion failed: expected %v, got %v", serviceType, reflect.TypeOf(instance)),
)
```

## 改进效果

### ✅ 代码组织

- **就近原则**：错误函数放在最相关的文件中
- **减少文件数**：移除了独立的errors.go文件
- **更好的内聚性**：相关代码聚集在一起

### ✅ 维护性

- **容易定位**：错误函数在使用它们的地方
- **减少跳转**：不需要在多个文件间跳转查看错误实现
- **统一标准**：使用统一的错误处理机制

### ✅ 性能和一致性

- **统一错误格式**：所有错误使用相同的AppError结构
- **统一错误码**：便于错误分类和国际化
- **减少依赖**：移除了不必要的文件依赖

## 验证结果

- ✅ 所有DI模块测试通过
- ✅ 应用编译成功
- ✅ 错误信息格式正确
- ✅ 代码组织更加合理
- ✅ 统一的错误处理机制

## 最佳实践总结

1. **错误函数就近放置**：将错误创建函数放在最相关的文件中
2. **低频错误内联**：对于使用频率低的错误，直接内联使用
3. **统一错误标准**：所有模块使用相同的错误处理机制
4. **避免过度抽象**：不为了重用而创建不必要的独立文件

这次重构成功实现了更好的代码组织和统一的错误处理，提高了代码的可维护性和一致性。通过将错误函数分散到相关文件中，消除了不必要的文件依赖，让代码结构更加清晰。
