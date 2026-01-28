# 安全审查报告：iDiary.ai (Hebe) 工程

**日期**: 2026-01-27  
**审查对象**: iDiary.ai (Hebe) Backend & Frontend  
**审查范围**: CSRF防护、API限速、暴力破解防护、输入校验完整性

---

## 1. 执行摘要

本次安全审查对 `iDiary.ai/hebe` 工程进行了全面的代码与架构分析。审查结果显示，工程在**输入校验**方面表现良好，后端拥有健全的校验机制。然而，在**防御性安全控制**方面存在**极高风险**的缺失：

- **CSRF防护**: 未实施，用户面临账户被劫持风险。
- **API限速**: 未实施，系统面临DoS攻击与资源滥用风险。
- **暴力破解防护**: 未实施，用户凭据面临被自动化工具撞库的风险。

建议开发团队**立即**着手修复CSRF与API限速问题，并在短期内完善暴力破解防护机制。

---

## 2. 工程现状分析与风险识别

### 2.1 跨站请求伪造 (CSRF) 防护

- **现状**:
  - 后端 `internal/app/middleware` 及路由层未发现任何 CSRF 中间件。
  - 用户登录 (`/login`)、注册 (`/register`) 及敏感操作（如 `UpdateUserEmail`）直接绑定请求数据，未校验 CSRF Token。
  - 前端 `auth.js` 中 HTMX 请求未携带自定义 Token 头，仅依赖 Cookie 进行会话维持。
- **风险等级**: **高 (Critical)**
- **影响**: 攻击者可诱导已登录用户点击恶意链接，在用户不知情的情况下修改其密码、邮箱或发布内容。

### 2.2 API 接口限速与防滥用

- **现状**:
  - 全工程未发现限流中间件（如 `rate` 或 `tollbooth`）。
  - `internal/app/controller/user/auth.go` 中的 `Login` 接口无频率限制，攻击者可无限次尝试登录。
  - 短信/邮件发送接口（如 `SendResetPasswordEmail`）虽有 Task 异步处理，但入口处无 IP 级或用户级频次控制。
- **风险等级**: **高 (Critical)**
- **影响**:
  - **服务拒绝 (DoS)**: 恶意高频请求耗尽服务器资源。
  - **资源盗刷**: 短信/邮件服务被滥用，产生高额费用。

### 2.3 暴力破解防护 (Brute Force)

- **现状**:
  - `internal/domain/user/service.go` 的 `LoginWithUsernameOrEmail` 方法在密码验证失败时直接返回 `PasswordIncorrect` 错误。
  - 缺失**失败计数器**：系统未记录用户或 IP 的连续失败次数。
  - 缺失**账户锁定**：无“连续失败 N 次锁定 M 分钟”的机制。
  - 缺失**人机验证**：登录与注册接口未集成 CAPTCHA 或类似验证。
- **风险等级**: **高 (High)**
- **影响**: 攻击者可使用字典库对特定账户进行长时间、高并发的密码猜解，直至攻破账户。

### 2.4 输入校验完整性

- **现状**:
  - **后端**: 表现优秀。使用了 `middleware.Validate[T]` 结合 `gin` binding 标签（`required`, `min`, `max`）及自定义校验器（`username`, `password`）。`pkg/regexcheck.go` 提供了统一的正则逻辑。
  - **前端**: 使用 `Uno UI` valid 插件，通过 `checkUsername` 和 `checkEmail` 实现了与后端的远程校验同步。
  - **一致性**: 前后端逻辑基本一致。
    - _微小差异_: `LoginRequest` 允许密码 `min=6`，而注册时 `password` 校验器要求更严格（8-80位，含数字字母）。这属于可接受的“宽进严出”策略，不构成安全漏洞。
- **风险等级**: **低 (Low)**
- **影响**: 输入校验较为扎实，SQL 注入与基础 XSS 风险较低（依赖于 ORM 和模板引擎的自动转义）。

---

## 3. 改进方案与实施建议

### 3.1 CSRF 防护改进 (已实施)

- **方案**: 实施 **Double Submit Cookie** 模式或 **Synchronizer Token** 模式。
- **实施状态**: **已完成 (2026-01-27)**
- **实施细节**:
  1.  **中间件**: 集成了 `gorilla/csrf`，并在 `internal/app/middleware/security.go` 中实现了 `CSRF` 中间件。
  2.  **Token 生成**:
      - 密钥从 `Cookie.Secret` 获取（自动填充至 32 字节）。
      - Cookie 设置为 `HttpOnly`，`SameSite=Lax`。
      - 生产环境自动启用 `Secure` 标志。
  3.  **前端适配**:
      - `website/views/header.html` 中添加了 `<meta name="csrf-token" ...>`。
      - `website/static/src/js/core/init.js` 中配置了 `htmx` 和 `fetch` 拦截器，自动添加 `X-CSRF-Token` 请求头。
  4.  **后端验证**: 全局中间件已注册，自动校验所有非 GET/HEAD/OPTIONS 请求。

### 3.2 API 限速 (Rate Limiting) (已实施)

- **方案**: 基于 Redis 的滑动窗口限流。
- **实施状态**: **已完成 (2026-01-27)**
- **实施细节**:
  1.  **引入库**: 集成了 `github.com/didip/tollbooth/v8`。
  2.  **全局限流**: `GeneralRateLimit` 中间件已注册，默认限制 1000 req/s (burst 100)。
  3.  **敏感接口限流**: `AuthRateLimit` 中间件已应用于 `/login`, `/register`, `/findpwd`, `/resetpwd` 等接口，默认限制 5 req/s (burst 10)。
  4.  **响应处理**: 超限时返回 HTTP 429 及 JSON 错误信息，并设置 `Retry-After` 头。
  5.  **配置**: 新增 `RateLimitConfig`，支持通过配置文件调整限流阈值。

### 3.3 暴力破解防护

- **方案**: 账户锁定策略 + 渐进式延迟。
- **实施步骤**:
  1.  **数据结构**: 在 Redis 中维护 `login_fail_count:{user_id}` 和 `login_fail_count:{ip}`。
  2.  **逻辑实现**:
      - 登录失败时：计数器 +1，设置过期时间（如 10 分钟）。
      - 计数 > 5：锁定账户 15 分钟，返回“尝试次数过多，请稍后再试”。
      - 登录成功时：清除计数器。
  3.  **前端增强**: 若连续失败 3 次，强制弹出图形验证码（CAPTCHA）。

### 3.4 增强建议 (Quick Wins)

- **安全标头**: 引入 `secure` 中间件，设置 `HSTS`, `X-Content-Type-Options`, `X-Frame-Options` 等 HTTP 头。
- **Cookie 安全**: 确保所有敏感 Cookie（如 `ud`）设置 `HttpOnly` 和 `SameSite=Lax/Strict` 属性。

---

## 4. 后续行动计划 (Action Plan)

| 优先级 | 任务                                           | 预计耗时 | 责任人    |
| :----- | :--------------------------------------------- | :------- | :-------- |
| **P0** | **实施 CSRF 中间件与前端 Token 传递**          | 4h       | 后端/前端 |
| **P0** | **实现 Redis 基础限流中间件 (API Rate Limit)** | 2h       | 后端      |
| **P1** | **实现登录失败计数与锁定机制**                 | 3h       | 后端      |
| **P1** | **配置 HTTP 安全标头 (Security Headers)**      | 1h       | 后端      |
| **P2** | **集成 CAPTCHA (针对高风险 IP/用户)**          | 4h       | 全栈      |
