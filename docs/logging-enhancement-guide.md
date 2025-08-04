# 日志系统增强功能使用指南

## 新增功能

增强后的日志系统现在包含以下新功能：

### 1. 自动日志文件测试

在初始化日志系统时，会自动测试日志文件是否可写：

```go
// 初始化时会自动测试
err := utils.logger.SetLogger(logConfig)
if err != nil {
    // 如果日志文件不可写，会立即返回错误
    slog.Error("日志系统初始化失败", "error", err)
    return err
}
```

### 2. 日志健康检查

可以随时检查日志文件是否存在并可写：

```go
// 检查日志文件健康状态
if err := utils.logger.CheckLogFileExists(logConfig); err != nil {
    slog.Error("日志文件健康检查失败", "error", err)

    // 尝试恢复日志系统
    if recoverErr := utils.logger.RecoverLogSystem(logConfig); recoverErr != nil {
        slog.Error("日志系统恢复失败", "error", recoverErr)
    }
}
```

### 3. 自动恢复功能

当检测到日志文件丢失时，可以自动恢复：

```go
// 恢复日志系统
err := utils.logger.RecoverLogSystem(logConfig)
if err != nil {
    // 恢复失败的处理
    fmt.Printf("日志系统恢复失败: %v\n", err)
} else {
    // 恢复成功，日志系统已重新初始化
    slog.Info("日志系统恢复成功")
}
```

## 建议的使用模式

### 1. 应用启动时

```go
// main.go
func main() {
    // 基础日志
    pkg.SetupBasicLogger()

    // 配置初始化
    config.InitConfig()

    // 增强的日志系统初始化
    logConfig := &pkg.LoggerConfig{
        Path:       config.GetConfig().Log.Path,
        Filename:   config.GetConfig().Log.Filename,
        // ... 其他配置
    }

    if err := pkg.SetLogger(logConfig); err != nil {
        slog.Error("日志系统初始化失败", "error", err)
        os.Exit(1)
    }

    slog.Info("应用启动成功")
}
```

### 2. 定期健康检查（可选）

```go
// 在应用中添加定期检查
func startLogHealthChecker(config *pkg.LoggerConfig) {
    ticker := time.NewTicker(5 * time.Minute)
    go func() {
        for range ticker.C {
            if err := pkg.CheckLogFileExists(config); err != nil {
                slog.Warn("检测到日志文件异常", "error", err)

                // 尝试自动恢复
                if recoverErr := pkg.RecoverLogSystem(config); recoverErr != nil {
                    slog.Error("日志系统自动恢复失败", "error", recoverErr)
                    // 可以在这里添加告警通知
                }
            }
        }
    }()
}
```

### 3. 错误处理中间件中的检查

```go
// 在关键操作前检查日志系统
func criticalOperation() {
    // 检查日志系统是否正常
    if err := pkg.CheckLogFileExists(logConfig); err != nil {
        // 尝试恢复
        pkg.RecoverLogSystem(logConfig)
    }

    // 执行关键操作
    slog.Info("执行关键操作")
    // ...
}
```

## 错误类型

新的错误消息更加详细：

- `日志文件测试失败: 无法写入日志文件 /path/to/log.log: permission denied`
- `日志文件不存在: /path/to/log.log`
- `日志文件无法写入: permission denied`
- `日志系统恢复失败: failed to create log directory`

## 优势

1. **早期发现**：在应用启动时就能发现日志问题
2. **自动恢复**：支持日志文件丢失后的自动恢复
3. **详细错误**：提供具体的错误信息，便于排查
4. **健康检查**：支持运行时的日志系统健康检查
5. **优雅降级**：即使日志文件有问题，应用也能正常运行（输出到控制台）

这些增强功能确保了日志系统的可靠性和稳定性。
