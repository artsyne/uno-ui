# Hebe配置系统

本文档描述了Hebe应用的配置管理系统。

## 概述

Hebe使用基于环境变量的配置系统，遵循12-Factor应用设计原则。配置系统具有以下特点：

1. **基于环境变量** - 所有配置项都通过环境变量设置
2. **类型安全** - 使用强类型结构体表示配置
3. **配置验证** - 内置配置验证机制确保配置正确
4. **合理默认值** - 为所有配置项提供合理的默认值
5. **灵活扩展** - 易于添加新的配置项

## 使用方法

### 1. 设置环境变量

所有配置项都以`IDIARY_`为前缀，例如：

```bash
export IDIARY_PORT=8080
export IDIARY_DB_HOST=localhost
```

### 2. 在应用中使用配置

```go
import "github.com/biuren/hebe/internal/configs"

func main() {
    // 初始化配置
    if err := configs.InitConfig(); err != nil {
        log.Fatalf("配置初始化失败: %v", err)
    }

    // 获取配置
    config := configs.GetConfig()

    // 使用配置
    log.Printf("应用端口: %s", config.App.Port)
}
```

### 3. 获取特定配置

提供了多个辅助函数获取特定配置：

```go
// 获取参数配置
params := configs.GetParams()

// 获取时区配置
timezones := configs.GetTimeZone()

// 获取支持的语言
langs := configs.GetSupportLangs()

// 获取静态文件路径
staticPath := configs.GetStaticPath()
```

## 配置验证

配置系统使用[validator](https://github.com/go-playground/validator)库进行验证，确保所有必需的配置项都存在且有效。验证规则定义在配置结构体的标签中。

## 示例

下面是一个完整的环境变量配置示例：

```bash
# 环境配置
export IDIARY_ENV=development
export IDIARY_GIN_MODE=debug

# 应用配置
export IDIARY_DOMAIN=http://127.0.0.1:8080
export IDIARY_PORT=9090
export IDIARY_SUMMARY_MAX=500
export IDIARY_HASH_SALT=some-secret-salt

# 数据库配置
export IDIARY_DB_HOST=127.0.0.1
export IDIARY_DB_PORT=3306
export IDIARY_DB_NAME=mydatabase
export IDIARY_DB_USER=root
export IDIARY_DB_PASSWORD=mypassword

# 更多配置...
```

## 测试

可以使用提供的`run_with_env.sh`脚本来测试配置系统：

```bash
./run_with_env.sh
```

或者使用测试工具：

```bash
go run cmd/configtest/main.go
```

## 添加新配置

要添加新的配置项，请按照以下步骤操作：

1. 在`internal/configs/types.go`中的相应结构体中添加字段
2. 为新字段添加适当的验证规则
3. 在`internal/configs/config.go`的`setDefaults`函数中添加默认值

## 最佳实践

1. 保持环境变量的命名一致性，使用`IDIARY_`前缀
2. 为敏感信息（如密码、密钥）使用环境变量，避免硬编码
3. 为所有配置提供合理的默认值
4. 使用配置验证确保应用程序使用有效的配置
5. 在不同环境（开发、测试、生产）使用不同的环境变量值
