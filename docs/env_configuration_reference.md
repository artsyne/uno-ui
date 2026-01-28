# 系统配置参数参考手册

本文档详细说明了 Hebe 项目的配置系统及所有可用配置参数的含义。项目使用 `viper` 和 `godotenv` 进行配置管理，支持通过环境变量或 `.env` 文件进行配置。

## 配置加载机制

系统启动时按以下优先级加载配置：
1. **环境变量**（最高优先级，通常用于生产环境注入）
2. **.env.local 文件**（用于本地开发覆盖，不应提交到版本控制）
3. **默认值**（代码中定义的兜底值）

所有环境变量统一使用 `IDIARY_` 前缀。

---

## 🔒 安全与限流配置 (Rate Limit & Security)

此部分配置直接关系到系统的稳定性和防攻击能力，请根据实际部署环境的负载能力谨慎调整。

### 限流策略 (Rate Limiting)
系统采用**令牌桶算法 (Token Bucket)** 实现限流控制。

| 配置项 | 示例值 | 说明 |
| :--- | :--- | :--- |
| `IDIARY_RATE_LIMIT_GENERAL` | `1000.0` | **通用接口限流速率**。<br>单位：请求/秒 (r/s)。<br>表示系统每秒向令牌桶中投放的令牌数量，即允许的平均请求速率。 |
| `IDIARY_RATE_LIMIT_GENERAL_BURST` | `100` | **通用接口突发容量**。<br>单位：请求数。<br>令牌桶的容量上限，允许短时间内突发的最大请求数量。 |
| `IDIARY_RATE_LIMIT_AUTH` | `5.0` | **认证接口限流速率**。<br>针对登录、注册等敏感接口的严格限流，防止撞库攻击。<br>单位：请求/秒。 |
| `IDIARY_RATE_LIMIT_AUTH_BURST` | `10` | **认证接口突发容量**。<br>认证接口允许的瞬间最大并发数。 |

### 暴力破解防护 (Brute Force Protection)
用于防止针对用户密码的暴力猜测攻击。

| 配置项 | 示例值 | 说明 |
| :--- | :--- | :--- |
| `IDIARY_BRUTE_FORCE_MAX_FAILURES` | `5` | **最大失败尝试次数**。<br>在锁定周期内，允许用户登录失败的最大次数。超过此次数账号将被临时锁定。 |
| `IDIARY_BRUTE_FORCE_LOCKOUT_MINUTES` | `15` | **账户锁定时间**。<br>单位：分钟。<br>触发暴力破解防护后，账号被禁止登录的时长。 |

---

## 🛠 基础与应用配置 (Basic & App)

### 环境标识
| 配置项 | 示例值 | 说明 |
| :--- | :--- | :--- |
| `IDIARY_ENV` | `development` | 运行环境标识。可选值：`development`, `production`, `test`。 |
| `IDIARY_GIN_MODE` | `debug` | Gin 框架运行模式。`debug` (输出详细日志) 或 `release` (高性能)。 |

### 应用参数
| 配置项 | 示例值 | 说明 |
| :--- | :--- | :--- |
| `IDIARY_DOMAIN` | `http://127.0.0.1:9090` | 应用对外服务的完整域名/地址，用于生成链接等。 |
| `IDIARY_PORT` | `9090` | HTTP 服务监听端口。 |
| `IDIARY_PAGE_SIZE` | `20` | 默认分页大小，即每页显示的数据条数。 |
| `IDIARY_MAX_DIARY_LIMIT` | `500` | 单个用户允许创建的最大日记本数量。 |
| `IDIARY_MAX_ENTRY_LENGTH` | `20000` | 单篇日记内容的最大字符数限制。 |
| `IDIARY_MAX_COMMENT_LENGTH` | `1000` | 单条评论的最大字符数限制。 |
| `IDIARY_STATIC_DIR` | `./website/static` | 静态资源文件系统路径。 |
| `IDIARY_LOCALES_DIR` | `./website/locales` | 后端多语言模板/翻译文件路径。 |
| `IDIARY_DEFAULT_LANG` | `en` | 系统默认语言代码。 |

---

## 💾 数据存储配置 (Storage)

### 数据库 (MySQL)
| 配置项 | 示例值 | 说明 |
| :--- | :--- | :--- |
| `IDIARY_DB_HOST` | `127.0.0.1` | 数据库服务器地址。 |
| `IDIARY_DB_PORT` | `3306` | 数据库端口。 |
| `IDIARY_DB_NAME` | `cd` | 数据库名称。 |
| `IDIARY_DB_USER` | `root` | 数据库用户名。 |
| `IDIARY_DB_PASSWORD` | `...` | 数据库密码。 |
| `IDIARY_DB_CONN` | `...` | **DSN 连接串模板**。<br>通常包含 `%s` 占位符供程序运行时填充上述参数。 |

### 缓存 (Redis)
| 配置项 | 示例值 | 说明 |
| :--- | :--- | :--- |
| `IDIARY_REDIS_HOST` | `127.0.0.1` | Redis 服务器地址。 |
| `IDIARY_REDIS_PORT` | `6379` | Redis 端口。 |
| `IDIARY_REDIS_PASSWORD` | `...` | Redis 认证密码。 |

---

## 🍪 会话与安全 (Session & Cookie)

| 配置项 | 示例值 | 说明 |
| :--- | :--- | :--- |
| `IDIARY_COOKIE_DOMAIN` | `127.0.0.1` | Cookie 作用域域名。生产环境应设置为顶级域名（如 `.example.com`）。 |
| `IDIARY_COOKIE_SECRET` | `...` | **Cookie 加密密钥**。<br>32位随机字符串，用于签名和加密 Cookie 数据，泄露会导致安全风险。 |
| `IDIARY_COOKIE_EXPIRE_DAYS` | `7` | "记住我" 功能的 Cookie 有效期（天）。 |
| `IDIARY_COOKIE_HTTP_ONLY` | `true` | 是否开启 HttpOnly 属性，防止 XSS 窃取 Cookie。 |

---

## 📧 邮件服务 (SMTP)

用于发送注册激活、密码重置等系统邮件。

| 配置项 | 示例值 | 说明 |
| :--- | :--- | :--- |
| `IDIARY_SMTP_SERVER` | `smtp...com:25` | SMTP 服务器地址及端口。 |
| `IDIARY_SMTP_USER` | `noreply@...` | 发件人账户/邮箱。 |
| `IDIARY_SMTP_PASSWORD` | `...` | SMTP 授权码或密码。 |
| `IDIARY_SMTP_FROM` | `noreply@...` | 邮件显示的“发件人”地址。 |
| `IDIARY_EMAIL_MAX_RETRIES`| `3` | 邮件发送失败后的最大重试次数。 |

---

## 📂 文件上传 (Upload)

| 配置项 | 示例值 | 说明 |
| :--- | :--- | :--- |
| `IDIARY_UPLOAD_PATH` | `upload/` | 上传文件的本地存储目录。 |
| `IDIARY_UPLOAD_TYPES` | `data:image...` | 允许上传的文件 MIME 类型列表（管道符 `|` 分隔）。 |
| `IDIARY_UPLOAD_MAX_SIZE` | `1048576` | 单个文件上传大小限制（字节）。例：`1048576` = 1MB。 |

---

## 📝 日志系统 (Logging)

| 配置项 | 示例值 | 说明 |
| :--- | :--- | :--- |
| `IDIARY_LOG_LEVEL` | `debug` | 日志级别。可选：`debug`, `info`, `warn`, `error`。 |
| `IDIARY_LOG_FORMAT` | `text` | 日志输出格式。可选：`text` (控制台友好), `json` (机器友好)。 |
| `IDIARY_LOG_PATH` | `./logs` | 日志文件存储目录。 |
| `IDIARY_LOG_MAX_SIZE` | `100` | 单个日志文件最大大小 (MB)，超过则轮转。 |
| `IDIARY_LOG_MAX_AGE` | `30` | 历史日志保留天数。 |
| `IDIARY_LOG_COMPRESS` | `true` | 是否压缩历史日志文件 (.gz)。 |

---

## 🌐 高级 JSON 配置

以下配置项为 JSON 数组字符串格式，用于配置复杂的数据结构。

| 配置项 | 说明 |
| :--- | :--- |
| `IDIARY_SUPPORT_LANGS` | **支持的语言列表**。<br>包含语言代码 (`code`)、显示标签 (`label`) 和本地化名称 (`local-text`)。 |
| `IDIARY_TIMEZONES` | **时区列表**。<br>字符串数组，定义了系统支持供用户选择的所有时区 ID。 |
| `IDIARY_PARAMS` | **系统参数键值对**。<br>用于定义动态文案模板，如：<br>- `reg_email_subject`: 注册邮件标题<br>- `reg_email_body`: 注册邮件内容模板 (支持 `%s` 占位符)<br>- `rest_pwd_email_...`: 重置密码邮件相关配置 |
