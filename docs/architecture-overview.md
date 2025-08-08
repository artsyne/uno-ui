# iDiary 系统架构概览

## 概述

iDiary是一个基于Go语言和Gin框架的智能日记管理系统，采用DDD（领域驱动设计）的分层架构，通过依赖注入实现松耦合设计，具备良好的可维护性和可扩展性。

## 分层架构

```
┌─────────────────────────────────────────────────────────────┐
│                    Website Layer                            │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐           │
│  │ Views   │ │ Static  │ │JavaScript│ │   CSS   │           │
│  │(模板文件)│ │(静态资源)│ │(前端逻辑) │ │(样式文件) │           │
│  └─────────┘ └─────────┘ └─────────┘ └─────────┘           │
└─────────────────────────────────────────────────────────────┘
                              │
┌─────────────────────────────────────────────────────────────┐
│                App Layer (应用服务层)                      │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐           │
│  │Startup  │ │ Route   │ │  View   │ │Controller│           │
│  │(启动器) │ │(路由)   │ │(视图)   │ │(控制器) │           │
│  └─────────┘ └─────────┘ └─────────┘ └─────────┘           │
└─────────────────────────────────────────────────────────────┘
                              │
┌─────────────────────────────────────────────────────────────┐
│                Domain Layer (业务领域层)                   │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐           │
│  │  User   │ │  Diary  │ │  Post   │ │  Task   │           │
│  │(用户领域)│ │(日记领域)│ │(内容领域)│ │(任务领域)│           │
│  │H-S-R    │ │H-S-R    │ │H-S-R    │ │Service  │           │
│  └─────────┘ └─────────┘ └─────────┘ └─────────┘           │
└─────────────────────────────────────────────────────────────┘
                              │
┌─────────────────────────────────────────────────────────────┐
│                Core Layer (核心基础设施层)                 │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐           │
│  │ Config  │ │   DI    │ │Database │ │Error &  │           │
│  │(配置)   │ │(依赖注入)│ │(数据库) │ │Middleware│          │
│  └─────────┘ └─────────┘ └─────────┘ └─────────┘           │
└─────────────────────────────────────────────────────────────┘
                              │
┌─────────────────────────────────────────────────────────────┐
│                Pkg Layer (工具层)                          │
│  ┌─────────────────┐  ┌─────────────────┐                 │
│  │ Internal Utilities│  │ External Utilities│                 │
│  │   (内部工具)    │  │   (外部工具)     │                 │
│  └─────────────────┘  └─────────────────┘                 │
└─────────────────────────────────────────────────────────────┘
                              │

```

**依赖方向**: Website → App → Module → Core → Pkg (严禁反向依赖)

## 分层职责详解

### App Layer (应用服务层)

**核心职责**: 应用编排、HTTP处理、渲染服务

#### Startup 模块

- **职责**: 应用生命周期管理
- **位置**: `internal/app/startup/`
- **功能**:
  - 应用启动和关闭流程编排
  - 组件初始化顺序管理
  - 错误处理和恢复机制

#### Route 模块

- **职责**: HTTP路由配置和管理
- **位置**: `internal/app/route/`
- **功能**:
  - 路由注册和分组
  - 中间件配置
  - 静态资源路由

#### View 模块

- **职责**: 统一视图服务（整合视图渲染、数据绑定和模板函数）
- **位置**: `internal/app/view/`
- **功能**:
  - 视图渲染和数据绑定
  - 模板函数管理
  - 国际化支持
  - 错误页面渲染

#### Controller 模块

- **职责**: 基础页面控制器
- **位置**: `internal/app/controller/`
- **功能**:
  - 首页等基础页面处理
  - 静态页面渲染

### Module Layer (业务模块层)

**核心职责**: 业务逻辑实现，采用Handler-Service-Repository模式

#### User 模块

- **职责**: 用户管理和认证
- **位置**: `internal/module/user/`
- **组件**:
  - **Handler**: HTTP请求处理 (注册、登录、设置等)
  - **Service**: 业务逻辑层 (用户认证、密码加密等)
  - **Repository**: 数据访问层 (用户CRUD操作)
  - **Auth**: 认证助手和中间件
  - **Types**: 数据结构定义

#### Diary 模块

- **职责**: 日记本管理
- **位置**: `internal/module/diary/`
- **组件**:
  - **Handler**: 日记本HTTP处理
  - **Service**: 日记本业务逻辑
  - **Repository**: 日记本数据访问
  - **Types**: 日记本相关数据结构

#### Post 模块

- **职责**: 日记内容管理
- **位置**: `internal/module/post/`
- **组件**:
  - **Handler**: 内容HTTP处理
  - **Service**: 内容业务逻辑
  - **Repository**: 内容数据访问
  - **Types**: 内容相关数据结构

#### Task 模块

- **职责**: 后台任务处理
- **位置**: `internal/module/task/`
- **组件**:
  - **Service**: 定时任务服务
  - **Repository**: 任务数据访问

### Core Layer (核心基础设施层)

**核心职责**: 基础设施服务，为上层提供技术支撑

#### Config 模块

- **职责**: 配置管理
- **位置**: `internal/core/config/`
- **功能**:
  - 配置文件加载和解析
  - 环境变量处理
  - 配置验证
  - 热更新支持

#### DI 模块

- **职责**: Go 1.24泛型依赖注入容器
- **位置**: `internal/core/di/`
- **功能**:
  - 零反射泛型设计
  - 三种生命周期 (瞬态/单例/实例)
  - 编译时类型安全
  - 高并发性能 (~57纳秒解析)

#### Database 模块

- **职责**: 数据库连接和管理
- **位置**: `internal/core/database/`
- **功能**:
  - 数据库连接池管理
  - SQL执行辅助函数
  - 事务管理

#### Error 模块

- **职责**: 错误定义和处理
- **位置**: `internal/core/error/`
- **功能**:
  - 统一错误码体系 (4位错误码)
  - AppError结构定义
  - 错误分类和处理策略

#### Middleware 模块

- **职责**: 核心中间件
- **位置**: `internal/core/middleware/`
- **功能**:
  - Panic恢复中间件
  - 错误处理中间件
  - 结构化日志记录

### Pkg Layer (工具层)

**核心职责**: 纯函数工具库，无依赖，可复用

**位置**: `pkg/`

**工具模块**:

- `logging.go`: slog日志配置和管理
- `crypto.go`: 加密解密工具 (AES、3DES、SHA256)
- `text.go`: 文本处理和智能截断
- `time.go`: 时间处理工具
- `image.go`: 图片上传和处理
- `cookie.go`: Cookie操作工具
- `helpers.go`: 通用助手函数
- `worker.go`: 多线程工作容器
- `session.go`: Redis会话管理

## 数据流向

### 请求处理流程

```
HTTP Request
    ↓
Middleware (Panic Recovery, Error Handler, SharedData)
    ↓
Route (路由匹配)
    ↓
Handler (HTTP处理层)
    ↓
Service (业务逻辑层)
    ↓
Repository (数据访问层)
    ↓
Database (数据库)
```

### 响应处理流程

```
Database
    ↓
Repository (数据组装)
    ↓
Service (业务逻辑处理)
    ↓
Handler (调用ViewService)
    ↓
ViewService (页面数据组装、模板渲染)
    ↓
HTTP Response (JSON/HTML)
```

## 关键设计模式

### 1. 依赖注入 (DI)

```go
// 服务注册
container.RegisterSingleton(func(c di.ContainerInterface) *view.ViewService {
    config := di.Resolve[*config.Config](c)
    return view.NewViewService(config)
})

// 服务解析
viewService := di.Resolve[*view.ViewService](container)
```

### 2. 策略模式 (错误处理)

```go
type AppError struct {
    Code     int    `json:"code"`
    MsgKey   string `json:"msg_key"`
    HttpCode int    `json:"-"`
    Details  string `json:"-"`
}

// 不同错误类型的处理策略
func (e *AppError) Handle(c *gin.Context) {
    switch e.HttpCode {
    case 400:
        c.JSON(400, gin.H{"code": e.Code, "message": e.MsgKey})
    case 500:
        c.JSON(500, gin.H{"code": e.Code, "message": "internal_error"})
    }
}
```

### 3. 模板方法模式 (Handler结构)

```go
type Handler struct {
    service     Service
    viewSvc     *view.ViewService
}

// 统一的处理模板
func (h *Handler) handleRequest(c *gin.Context, businessLogic func() (interface{}, error)) {
    data, err := businessLogic()
    if err != nil {
        h.viewSvc.JSONError(c, err)
        return
    }
    h.viewSvc.JSONSuccess(c, data)
}
```

### 4. 装饰器模式 (中间件)

```go
func PanicRecoveryMiddleware() gin.HandlerFunc {
    return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
        slog.Error("panic recovered", "panic", recovered)
        c.JSON(500, gin.H{"code": 1000, "message": "internal_error"})
    })
}
```

## 配置管理

### 配置层级

1. **默认配置**: 代码中的默认值
2. **配置文件**: `config/config.yaml`
3. **环境变量**: 覆盖配置文件
4. **命令行参数**: 最高优先级

### 配置结构

```go
type Config struct {
    App      AppConfig      `yaml:"app"`
    Database DatabaseConfig `yaml:"database"`
    Log      LogConfig      `yaml:"log"`
    Cookie   CookieConfig   `yaml:"cookie"`
    Redis    RedisConfig    `yaml:"redis"`
}
```

## 日志系统

### 日志级别

- **DEBUG**: 详细调试信息
- **INFO**: 一般信息记录
- **WARN**: 警告信息
- **ERROR**: 错误信息

### 日志格式

- **结构化日志**: 使用slog库，JSON格式输出
- **文件轮转**: 基于大小和时间的日志轮转
- **多输出**: 同时输出到文件和控制台

### 日志示例

```go
slog.Info("用户登录",
    "user_id", 123,
    "username", "john_doe",
    "ip", "192.168.1.100",
    "duration", "150ms")
```

## 错误处理体系

### 错误码分类

- **1xxx**: 通用错误 (系统错误、验证错误等)
- **2xxx**: 用户模块错误
- **3xxx**: 日记模块错误
- **4xxx**: 内容模块错误
- **5xxx**: 网络和外部服务错误

### 错误处理流程

1. **业务层**: 创建AppError
2. **Handler层**: 调用ViewService处理错误
3. **中间件层**: 统一错误格式和日志记录
4. **客户端**: 接收统一格式的错误响应

## 性能特征

### DI容器性能

- **解析速度**: 19.89纳秒
- **内存分配**: 零分配
- **并发安全**: 读写锁保护

### 渲染性能

- **模板缓存**: 预编译模板缓存
- **页面数据**: 结构化数据组装
- **国际化**: 消息缓存机制

### 数据库性能

- **连接池**: 数据库连接池管理
- **预编译语句**: SQL语句预编译
- **事务管理**: 自动事务管理

## 安全机制

### 认证和授权

- **JWT Token**: 无状态认证
- **Session管理**: Redis会话存储
- **权限控制**: 基于用户组的权限管理

### 数据安全

- **密码加密**: SHA256哈希
- **数据传输**: HTTPS加密传输
- **输入验证**: 参数验证和过滤

### 错误处理安全

- **敏感信息**: 不在错误响应中暴露敏感信息
- **日志脱敏**: 敏感数据脱敏记录
- **错误码**: 抽象的错误码，不暴露内部实现

## 部署架构

### 单体部署

```
Load Balancer
    ↓
iDiary Application (Go Binary)
    ↓
MySQL Database
    ↓
Redis Cache
```

### 容器化部署

```
Docker Compose:
- iDiary (Go Application)
- MySQL (Database)
- Redis (Cache)
- Nginx (Static Files & Reverse Proxy)
```

## 扩展点

### 水平扩展

- **无状态设计**: 应用层无状态，支持多实例部署
- **数据库连接池**: 支持数据库读写分离
- **缓存层**: Redis集群支持

### 功能扩展

- **插件系统**: 预留插件接口
- **API开放**: RESTful API支持
- **多租户**: 支持多租户架构

### 存储扩展

- **云存储**: 支持多云盘存储
- **文件系统**: 本地和云端文件存储
- **数据备份**: 自动备份和恢复机制

这种架构设计确保了系统的高可维护性、可扩展性和稳定性，为未来的功能扩展和性能优化奠定了坚实的基础。
