# iDiary.ai Hebe 日志系统指南

## 概述

iDiary.ai Hebe 采用基于 Go 1.24 `slog` 标准库的现代化结构化日志系统，提供高性能、类型安全的日志记录功能。

## 日志架构

### 核心特性

- **结构化日志**: 基于 `slog` 标准库，支持 JSON 和文本格式
- **分级日志**: Debug、Info、Warn、Error 四个级别
- **高性能**: 零分配设计，支持高并发日志记录
- **配置驱动**: 通过 Viper 配置管理日志行为
- **自动轮转**: 支持按大小和时间自动轮转日志文件
- **统一格式**: 项目中所有日志都使用统一的结构化格式

### 日志配置

**文件位置**: `internal/utils/config/types.go`

```go
// LogConfig 日志配置
type LogConfig struct {
    Level      string `json:"level" mapstructure:"LOG_LEVEL" validate:"required,oneof=debug info warn error"`
    Format     string `json:"format" mapstructure:"LOG_FORMAT" validate:"required,oneof=text json"`
    Path       string `json:"path" mapstructure:"LOG_PATH" validate:"required"`
    Filename   string `json:"filename" mapstructure:"LOG_FILENAME" validate:"required"`
    MaxSize    int    `json:"max_size" mapstructure:"LOG_MAX_SIZE" validate:"required,min=1,max=500"`
    MaxAge     int    `json:"max_age" mapstructure:"LOG_MAX_AGE" validate:"required,min=1"`
    MaxBackups int    `json:"max_backups" mapstructure:"LOG_MAX_BACKUPS" validate:"required,min=1"`
    Compress   bool   `json:"compress" mapstructure:"LOG_COMPRESS"`
    LocalTime  bool   `json:"local_time" mapstructure:"LOG_LOCAL_TIME"`
}
```

### 日志初始化

**文件位置**: `pkg/logging.go`

```go
// SetLogger 设置 slog 默认日志器
func SetLogger(config config.LogConfig) error {
    logFile := filepath.Join(config.Path, config.Filename)

    // 创建日志目录
    if err := os.MkdirAll(config.Path, 0755); err != nil {
        return fmt.Errorf("创建日志目录失败: %w", err)
    }

    // 配置日志轮转
    writer := &lumberjack.Logger{
        Filename:   logFile,
        MaxSize:    config.MaxSize,
        MaxBackups: config.MaxBackups,
        MaxAge:     config.MaxAge,
        Compress:   config.Compress,
        LocalTime:  config.LocalTime,
    }

    // 设置日志级别和格式
    var handler slog.Handler
    if config.Format == "json" {
        handler = slog.NewJSONHandler(writer, &slog.HandlerOptions{
            Level: getLogLevel(config.Level),
        })
    } else {
        handler = slog.NewTextHandler(writer, &slog.HandlerOptions{
            Level: getLogLevel(config.Level),
        })
    }

    slog.SetDefault(slog.New(handler))
    return nil
}
```

## 实际使用示例

### 在控制器中使用

```go
// internal/app/controller/user.go
func (ctrl *UserController) Login(c *gin.Context) {
    username := c.PostForm("username")

    slog.Info("用户登录尝试",
        "username", username,
        "ip", c.ClientIP(),
        "user_agent", c.GetHeader("User-Agent"),
    )

    user, err := ctrl.userService.Login(username, password)
    if err != nil {
        slog.Warn("用户登录失败",
            "username", username,
            "error", err,
            "ip", c.ClientIP(),
        )
        // 错误处理...
        return
    }

    slog.Info("用户登录成功",
        "user_id", user.UserId,
        "username", username,
        "ip", c.ClientIP(),
    )
}
```

### 在服务层使用

```go
// internal/domain/user/service.go
func (s *Service) CreateUser(user *User) (*User, error) {
    slog.Debug("开始创建用户",
        "username", user.UserName,
        "email", user.Email,
    )

    // 验证用户数据
    if err := s.validateUser(user); err != nil {
        slog.Error("用户数据验证失败",
            "username", user.UserName,
            "error", err,
        )
        return nil, err
    }

    // 创建用户
    result, err := s.repo.Create(user)
    if err != nil {
        slog.Error("用户创建失败",
            "username", user.UserName,
            "error", err,
        )
        return nil, err
    }

    slog.Info("用户创建成功",
        "user_id", result.UserId,
        "username", result.UserName,
    )

    return result, nil
}
```

### 在中间件中使用

```go
// internal/utils/middleware/apperror.go
func AppErrorHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        defer func() {
            if len(c.Errors) > 0 {
                err := c.Errors.Last()

                slog.Error("请求处理错误",
                    "method", c.Request.Method,
                    "path", c.Request.URL.Path,
                    "ip", c.ClientIP(),
                    "error", err.Error(),
                )

                // 错误处理逻辑...
            }
        }()

        c.Next()
    }
}
```

## 日志格式和输出

### JSON 格式输出

```json
{
  "time": "2024-01-15T10:30:45.123Z",
  "level": "INFO",
  "msg": "用户登录成功",
  "user_id": 123,
  "username": "johndoe",
  "ip": "192.168.1.100"
}
```

### 文本格式输出

```
time=2024-01-15T10:30:45.123Z level=INFO msg="用户登录成功" user_id=123 username=johndoe ip=192.168.1.100
```

## 最佳实践

### 1. 使用结构化字段

```go
// ✅ 好的做法
slog.Info("用户操作",
    "action", "create_diary",
    "user_id", userID,
    "diary_id", diaryID,
    "duration_ms", time.Since(start).Milliseconds(),
)

// ❌ 避免字符串拼接
slog.Info(fmt.Sprintf("用户 %d 创建了日记 %d", userID, diaryID))
```

### 2. 合理使用日志级别

```go
// Debug: 详细的开发调试信息
slog.Debug("SQL 查询", "sql", query, "args", args)

// Info: 一般的业务流程信息
slog.Info("用户注册", "username", username)

// Warn: 警告但不影响正常流程
slog.Warn("缓存未命中", "key", cacheKey)

// Error: 错误信息，需要关注
slog.Error("数据库连接失败", "error", err)
```

### 3. 错误处理中的日志

```go
func (s *Service) ProcessData(data *Data) error {
    if err := s.validate(data); err != nil {
        slog.Error("数据验证失败",
            "data_id", data.ID,
            "error", err,
            "validation_rules", s.getRules(),
        )
        return fmt.Errorf("数据验证失败: %w", err)
    }

    return nil
}
```

## 性能考虑

### 高性能日志记录

- **结构化字段**: slog 使用高效的键值对结构
- **延迟计算**: 只在实际需要时格式化日志内容
- **零分配**: 在热路径中避免不必要的内存分配
- **异步写入**: 使用缓冲 I/O 提高写入性能

### 性能监控

```go
func (s *Service) PerformanceLog(operation string, duration time.Duration) {
    slog.Info("性能监控",
        "operation", operation,
        "duration_ms", duration.Milliseconds(),
        "slow_query", duration > 100*time.Millisecond,
    )
}
```

## 总结

iDiary.ai Hebe 的日志系统基于 Go 1.24 slog 标准库，提供了现代化、高性能的结构化日志功能。通过统一的配置管理、合理的日志分级和结构化字段设计，确保了系统的可观测性和调试便利性。

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
