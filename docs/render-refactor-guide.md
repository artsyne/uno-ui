# Render 模块渐进式重构指南 (更新版)

## 概述

当前 Render 模块是一个集成了多种职责的单体服务，包含页面数据组装、HTTP响应处理、模板函数管理、本地化等功能。本文档提供渐进式重构方案，逐步将其拆分为松耦合、可测试的组件，充分利用 Go 1.24 新特性并拥抱依赖注入 (DI) 理念，最终实现高性能、高可维护性的模块。

**核心角色演变**: 本次重构将保持 `RenderServiceInterface` 对外接口的完全兼容，确保调用方代码无需任何修改。然而，`RenderService` 结构体本身的角色将发生根本性转变：它将从一个集成了所有实现细节的“万能瑞士军刀”，演变为一个非常“薄”的“项目经理”。重构后，`Service` 实例将不再包含任何业务逻辑，其所有方法都将直接委托给通过依赖注入传入的、高度解耦的专业组件来完成。

## 当前架构分析

### 现状概览

```
internal/app/render/
├── service.go              # 主服务入口（91行）- 所有功能集成点
├── interfaces.go           # 服务接口定义（35行）- 包含所有方法
├── types.go               # 类型定义（208行）- RenderUser接口, PageData等
├── context_builder.go     # 上下文构建器（47行）- 处理登录状态和语言检测
├── page_builder.go        # 页面构建器（123行）- 生成页面数据
├── response_handler.go    # 响应处理器（约80行）- JSON/HTML响应
├── template_funcs.go      # 模板函数（约200行）- 模板辅助函数
├── localization.go        # 本地化管理（约150行）- 多语言支持
└── utils.go              # 工具类（约180行）- 各种辅助方法
```

### 依赖关系分析 (更新)

- **⚠️ 存在循环依赖**:
  - 当前 `RenderService` (在 `service.go` 中) 依赖 `ResponseHandler`, `TemplateFuncs`, `RenderUtils`。
  - 同时，`response_handler.go`, `template_funcs.go`, `utils.go` 中的对应结构体（`ResponseHandler`, `TemplateFuncs`, `RenderUtils`）的构造函数都接收 `*RenderService` 作为参数，并在各自结构体中持有 `*RenderService` 的引用。
  - 这导致了 `RenderService <-> ResponseHandler`、`RenderService <-> TemplateFuncs` 和 `RenderService <-> RenderUtils` 的**相互循环依赖**。
  - **重构目标**：本方案的后续阶段将明确解决这些循环依赖，实现单向依赖。
- **⚠️ 组件耦合**: 各组件间存在直接依赖，测试困难。
- **⚠️ 职责混合**: `service.go` 承担了过多协调责任。

### Go 1.24 新特性应用机会

- **泛型约束**: 用于类型安全的数据处理。
- **结构化日志**: `slog` 包的性能优化。
- **context包增强**: 更好的请求上下文管理。
- **性能优化**: 新的内存布局和GC改进。

## 重构原则

1.  **保持功能不变**: 重构过程中所有对外接口保持兼容。
2.  **解决依赖问题**: 通过接口和**依赖注入 (DI)** 解决循环依赖和组件耦合。
3.  **拥抱 DI 框架**: 充分利用现有 DI 框架的能力，将组件的生命周期和依赖管理交给框架。
4.  **简单设计**: 避免过度设计，采用简单直接的架构。
5.  **测试驱动**: 每个步骤都要确保测试覆盖。
6.  **Go 1.24优化**: 充分利用新特性提升性能和可维护性。

## 渐进式重构计划 (更新版)

### 阶段1: 接口分离与类型安全强化 (1-2天)

**目标**: 为各组件定义清晰接口，利用 Go 1.24 泛型提升类型安全，为后续的依赖注入打下基础。

#### 1.1 创建组件接口定义

```go
// internal/app/render/contracts/interfaces.go
package contracts

import (
    "html/template"
    "github.com/gin-gonic/gin"
    // 假设RenderContext和RenderUser类型提升到render包或共享包
    "github.com/biuren/hebe/internal/app/render"
)

// 使用Go 1.24泛型约束提升类型安全 (示例)
type DataConstraint interface {
    map[string]any | map[string]interface{} | any
}

// ContextBuilder 接口
type ContextBuilderInterface interface {
    BuildContext(c *gin.Context) *render.RenderContext // 返回render包下的RenderContext
}

// PageBuilder 接口
type PageBuilderInterface interface {
    BuildPageData(c *gin.Context, titleKey string) *render.PageData // 返回render包下的PageData
}

// ResponseHandler 接口
type ResponseHandlerInterface interface {
    JSON(c *gin.Context, code int, data any)
    JSONSuccess(c *gin.Context, data any)
    JSONError(c *gin.Context, err error, data ...any)
    HTML(c *gin.Context, name string, data any)
    HTMLError(c *gin.Context, title, message string)
    Redirect(c *gin.Context, url string)
}

// LocalizationManager 接口
type LocalizationManagerInterface interface {
    GetMessage(key, locale string) string
    GetMessageWithParams(key, locale string, params map[string]string) string // 如果有
    GetLocaleFromContext(c *gin.Context) string
    GetUserLocale(user render.RenderUser) string // 接收render包下的RenderUser接口
}

// TemplateManager 接口 (根据重构后的模板管理职责定义)
type TemplateManagerInterface interface {
    RegisterFunc(name string, fn any) // 注册单个模板函数
    GetFuncMap() template.FuncMap    // 获取所有注册的模板函数映射
}

// RenderUtilsInterface 接口 (我的建议：如果utils功能无法完全变为纯函数，则为其定义接口)
type RenderUtilsInterface interface {
    SafeHTML(content string) template.HTML
    FormatFileSize(size int64) string
    TruncateString(str string, maxLength int) string
    FormatDate(timestamp int64, format string) string
    HtmlUnescape(str string) string
    GetCurrentTimestamp() int64
    GetCurrentTimestampMs() int64
    ParseTimestamp(timestamp int64) time.Time // 需要导入time包
    ParseTimestampMs(timestamp int64) time.Time
    GenerateRandomString(length int) string
    // ... 其他原utils.go中的方法 ...
}
```

#### 1.2 改进类型定义 (部分已在原文档中呈现)

```go
// internal/app/render/types.go 优化
package render

import "context"

// 使用Go 1.24的类型约束
type DataMap = map[string]any

// 增强的上下文结构，支持Go 1.24的context包新特性
type RenderContext struct {
    IsLogin  bool       `json:"is_login"`
    AuthUser RenderUser `json:"auth_user,omitempty"`
    Lang     string     `json:"lang"`

    // Go 1.24 新增：结构化上下文支持，用于在请求链中传递取消信号、截止日期和值
    ctx context.Context `json:"-"`
}

// WithContext 支持Go 1.24的增强context
func (rc *RenderContext) WithContext(ctx context.Context) *RenderContext {
    rc.ctx = ctx
    return rc
}

func (rc *RenderContext) Context() context.Context {
    if rc.ctx == nil {
        return context.Background() // 默认背景上下文
    }
    return rc.ctx
}

// PageData 结构体需要确保内部引用的是 contracts 中的接口或 render 包中的共享类型
type PageData struct {
    Title               string                  `json:"title"`
    Description         string                  `json:"description"`
    Data                DataMap                 `json:"data"`
    Breadcrumb          []BreadcrumbItem        `json:"breadcrumb"`
    Context             *RenderContext          `json:"context"`
    // 注入LocalizationManagerInterface，而不是具体的LocalizationManager
    localizationManager contracts.LocalizationManagerInterface `json:"-"`
}

// ErrorResponse 和 JSONResponse 结构体也应保持清晰，最好直接定义在 response/handler.go 或一个公共 errors/types.go 中
```

#### 1.3 执行步骤

1.  创建 `contracts/` 目录和所有组件的接口定义，确保接口方法签名与原方法兼容。
2.  更新 `internal/app/render/types.go`，确保 `PageData` 和 `RenderContext` 等关键类型能够适配新的接口和 Go 1.24 特性。
3.  修改现有组件实现文件（如 `context_builder.go`, `page_builder.go`, `response_handler.go`, `template_funcs.go`, `localization.go`, `utils.go`）以实现对应的接口。暂时保持其内部依赖 `*RenderService` 的现状。
4.  修改 `service.go`，使其开始通过接口（而非具体类型）引用其子组件，但这不涉及 `NewService` 的签名改变。
5.  运行测试确保无功能回归。

**验收标准**:

- [x] 所有组件都有明确接口定义。
- [x] 编译通过，无类型错误。
- [x] 现有功能测试全部通过。

---

### 阶段2: Context 组件独立化 (1天)

**目标**: 将上下文构建逻辑完全独立，并通过接口注入 `LocalizationManager`，消除其对 `RenderService` 的间接依赖。

#### 2.1 增强 ContextBuilder

```go
// internal/app/render/context/builder.go
package context

import (
    "context"
    "log/slog" // Go 1.24 新增的结构化日志
    "github.com/gin-gonic/gin"
    "github.com/biuren/hebe/internal/app/render/contracts"
    "github.com/biuren/hebe/internal/app/render" // 引用共享的render包类型
)

type Builder struct {
    localization contracts.LocalizationManagerInterface // 注入接口
    logger       *slog.Logger
}

// NewBuilder 构造函数：通过参数接收所有依赖接口，而不是在内部创建
func NewBuilder(loc contracts.LocalizationManagerInterface) *Builder {
    return &Builder{
        localization: loc,
        logger:       slog.Default().With("component", "context_builder"),
    }
}

// BuildContext 构建增强的渲染上下文
func (b *Builder) BuildContext(c *gin.Context) *render.RenderContext {
    ctx := c.Request.Context() // Go 1.24: 获取请求上下文

    renderCtx := &render.RenderContext{}
    renderCtx.WithContext(ctx)

    // 现有逻辑保持不变，但确保类型转换正确
    if isLogin, exists := c.Get("IsLogin"); exists {
        if login, ok := isLogin.(bool); ok {
            renderCtx.IsLogin = login
        }
    }

    if authUserRaw, exists := c.Get("AuthUser"); exists && authUserRaw != nil {
        if renderUser, ok := authUserRaw.(render.RenderUser); ok { // 使用render.RenderUser接口
            renderCtx.AuthUser = renderUser
        }
    }

    renderCtx.Lang = b.localization.GetLocaleFromContext(c)

    b.logger.DebugContext(ctx, "context built",
        "is_login", renderCtx.IsLogin,
        "lang", renderCtx.Lang)

    return renderCtx
}
```

#### 2.2 执行步骤

1.  创建独立的 `context/` 包，并将 `context_builder.go` 移动到其中。
2.  更新 `context/builder.go` 中的 `NewBuilder` 方法，使其通过参数接收 `contracts.LocalizationManagerInterface`，而不是直接依赖 `*LocalizationManager`。
3.  更新 `service.go` 中的 `NewRenderService`，修改 `contextBuilder` 的初始化方式，使其使用新的 `context.NewBuilder` 并传入 `localizationManager`。
4.  添加单元测试覆盖 `context/builder.go`，使用 Mock `LocalizationManagerInterface` 进行测试。
5.  验证功能完整性，确保上下文构建逻辑正确。

**验收标准**:

- [x] Context 构建逻辑完全独立。
- [x] `ContextBuilder` 仅依赖接口，消除对具体实现的耦合。
- [x] 单元测试覆盖率 \> 80%。
- [x] 集成测试通过。

---

### 阶段3: Page 组件模块化 (1-2天)

**目标**: 将页面构建逻辑模块化，并通过接口注入其依赖，提升扩展性和测试性，消除对 `RenderService` 的间接依赖。

#### 3.1 重构 PageBuilder

```go
// internal/app/render/page/builder.go
package page

import (
    "context"
    "fmt"
    "strings"
    "github.com/gin-gonic/gin"
    "github.com/biuren/hebe/internal/app/render/contracts"
    "github.com/biuren/hebe/internal/core/config"
    "github.com/biuren/hebe/internal/app/render" // 引用共享的render包类型
)

type Builder struct {
    config       *config.Config
    context      contracts.ContextBuilderInterface      // 注入接口
    localization contracts.LocalizationManagerInterface // 注入接口
    processors   []DataProcessor[any]                   // Go 1.24: 泛型支持的数据处理器
}

// NewBuilder 构造函数：通过参数接收所有依赖接口
func NewBuilder(
    cfg *config.Config,
    contextBuilder contracts.ContextBuilderInterface,
    localizationManager contracts.LocalizationManagerInterface) *Builder {
    return &Builder{
        config:       cfg,
        context:      contextBuilder,
        localization: localizationManager,
        processors:   make([]DataProcessor[any], 0),
    }
}

// DataProcessor 使用Go 1.24泛型的数据处理器接口
type DataProcessor[T any] interface {
    Process(ctx context.Context, data T) T
}

// RegisterProcessor 注册数据处理器
func (pb *Builder) RegisterProcessor(processor DataProcessor[any]) {
    pb.processors = append(pb.processors, processor)
}

// MetadataProcessor 元数据处理器示例 (可在独立的Go文件中实现)
type MetadataProcessor struct {
    localization contracts.LocalizationManagerInterface
}

func NewMetadataProcessor(loc contracts.LocalizationManagerInterface) *MetadataProcessor {
    return &MetadataProcessor{localization: loc}
}

func (mp *MetadataProcessor) Process(ctx context.Context, data *render.PageData) *render.PageData {
    // SEO 元数据处理逻辑
    if data.Description == "" { // 如果没有设置描述，则从本地化获取默认
        data.Description = mp.localization.GetMessage("common.default_page_desc", data.Context.Lang)
    }
    return data
}

// BuildPageData 构建页面数据（核心逻辑保持不变，但内部调用使用注入的接口）
func (pb *Builder) BuildPageData(c *gin.Context, titleKey string) *render.PageData {
    pageData := &render.PageData{
        Title:       pb.generatePageTitle(c, titleKey),
        Description: pb.generateDescription(c, titleKey),
        Data:        make(render.DataMap),
        Breadcrumb:  pb.generateBreadcrumb(c),
        Context:     pb.context.BuildContext(c),
        // 直接将 localizationManager 注入到 PageData 中，供模板函数直接使用
        localizationManager: pb.localization,
    }

    // Go 1.24: 使用泛型处理器链
    ctx := c.Request.Context()
    for _, processor := range pb.processors {
        if p, ok := processor.(DataProcessor[*render.PageData]); ok {
            pageData = p.Process(ctx, pageData)
        }
    }

    return pageData
}

// generatePageTitle, generateDescription, generateBreadcrumb 等方法保持不变，但确保调用的是注入的接口
func (pb *Builder) generatePageTitle(c *gin.Context, titleKey string) string {
    locale := pb.localization.GetLocaleFromContext(c)
    title := pb.localization.GetMessage(titleKey, locale)
    siteName := pb.config.App.SiteName
    return fmt.Sprintf("%s - %s", title, siteName)
}

func (pb *Builder) generateDescription(c *gin.Context, titleKey string) string {
    locale := pb.localization.GetLocaleFromContext(c)
    descKey := titleKey + ".desc"
    desc := pb.localization.GetMessage(descKey, locale)
    if desc == descKey { // 如果没有找到对应的描述，使用通用默认描述
        return pb.localization.GetMessage("common.default_page_desc", locale)
    }
    return desc
}

func (pb *Builder) generateBreadcrumb(c *gin.Context) []render.BreadcrumbItem {
    // 现有逻辑保持不变，但确保引用render包下的类型
    path := c.Request.URL.Path
    breadcrumb := []render.BreadcrumbItem{
        {Name: "首页", URL: "/"},
    }
    // ... 根据路径模式智能生成面包屑逻辑 ...
    return breadcrumb
}
```

#### 3.2 执行步骤

1.  创建独立的 `page/` 包，并将 `page_builder.go` 移动到其中。
2.  更新 `page/builder.go` 中的 `NewBuilder` 方法，使其通过参数接收 `contracts.ContextBuilderInterface` 和 `contracts.LocalizationManagerInterface`。
3.  在 `service.go` 的 `NewRenderService` 中，修改 `pageBuilder` 的初始化方式，传入正确的依赖。
4.  添加单元测试覆盖 `page/builder.go`，使用 Mock 依赖进行测试。
5.  验证功能完整性，确保所有页面正常渲染。

**验收标准**:

- [x] 页面构建逻辑模块化。
- [x] `PageBuilder` 仅依赖接口，消除对具体实现的耦合。
- [x] 支持可扩展的数据处理器。
- [x] 单元测试覆盖率 \> 80%。
- [x] 性能无明显下降。

---

### 阶段4: Response 组件优化 (1天)

**目标**: 优化响应处理，提升性能和可维护性，并消除对 `RenderService` 的依赖。

#### 4.1 响应处理器增强

```go
// internal/app/render/response/handler.go
package response

import (
    "context"
    "log/slog"
    "net/http"
    "sync"
    "github.com/gin-gonic/gin"
    errors "github.com/biuren/hebe/internal/core/error" // 假设错误类型在core/error
    "github.com/biuren/hebe/internal/app/render"        // 引用共享的render包类型
)

// 使用Go 1.24的泛型约束，确保JSON响应的数据类型
type ResponseData interface {
    ~map[string]any | ~struct{} | ~string | ~int | ~bool | any // 允许更广泛的类型，但建议限制
}

// JSONResponse 和 ErrorResponse 结构体可以在这里定义，或者在共享的render/types.go中
type JSONResponse struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data"`
}

type ErrorResponse struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}

type Handler struct {
    jsonPool sync.Pool // Go 1.24: 预分配的响应池（针对map或struct的复用）
    logger   *slog.Logger
}

// NewHandler 构造函数：不再接收RenderService，实现彻底解耦
func NewHandler() *Handler {
    return &Handler{
        jsonPool: sync.Pool{
            New: func() interface{} {
                return make(map[string]any, 8) // 预分配常用大小
            },
        },
        logger: slog.Default().With("component", "response_handler"),
    }
}

// JSON 使用Go 1.24性能优化的JSON响应 (泛型)
func (h *Handler) JSON[T ResponseData](c *gin.Context, code int, data T) {
    c.Header("Content-Type", "application/json; charset=utf-8")
    // 使用统一的JSONResponse结构
    c.JSON(http.StatusOK, JSONResponse{Code: code, Message: "success", Data: data})

    h.logger.DebugContext(c.Request.Context(), "json response",
        "status_code", code,
        "path", c.Request.URL.Path)
}

// JSONSuccess, JSONError, HTML, HTMLError, Redirect 方法调整为不依赖RenderService
func (h *Handler) JSONSuccess(c *gin.Context, data any) {
    h.JSON(c, 0, data) // 通常0代表成功
}

func (h *Handler) JSONError(c *gin.Context, err error, data ...any) {
    var errorData any
    if len(data) > 0 {
        errorData = data[0]
    }

    if appErr, ok := err.(*errors.AppError); ok {
        c.JSON(appErr.HttpCode, ErrorResponse{ // 使用 ErrorResponse 结构
            Code:    appErr.Code,
            Message: appErr.Key,
            Data:    errorData,
        })
        h.logAppError(c, appErr)
        return
    }

    // 处理其他错误，转换为内部错误
    internalErr := errors.InternalError.Err()
    response := ErrorResponse{
        Code:    internalErr.Code,
        Message: internalErr.Key,
        Data:    errorData,
    }
    c.JSON(internalErr.HttpCode, response)
}

func (rh *Handler) HTML(c *gin.Context, name string, data any) {
    c.HTML(http.StatusOK, name, data)
}

func (rh *Handler) HTMLError(c *gin.Context, title, message string) {
    pageData := &render.PageData{ // 使用 render 包下的 PageData
        Title: title,
        Data: map[string]any{
            "error_title":   title,
            "error_message": message,
        },
    }
    c.HTML(http.StatusOK, "error.html", pageData)
}

func (rh *Handler) Redirect(c *gin.Context, url string) {
    c.Redirect(http.StatusFound, url)
}

// logAppError 辅助函数
func (rh *Handler) logAppError(c *gin.Context, appErr *errors.AppError) {
    h.logger.ErrorContext(c.Request.Context(), "app error",
        "code", appErr.Code,
        "message", appErr.Message,
        "http_code", appErr.HttpCode,
        "path", c.Request.URL.Path,
        "original_error", appErr.Err.Error())
}
```

#### 4.2 执行步骤

1.  创建独立的 `response/` 包，并将 `response_handler.go` 移动到其中。
2.  更新 `response/handler.go` 中的 `NewHandler` 方法，使其**不再接收 `*RenderService`**，实现彻底解耦。
3.  修改 `service.go` 中 `responseHandler` 的初始化方式，直接调用 `response.NewHandler()`。
4.  添加性能优化（如 `sync.Pool` 的使用）。
5.  对响应处理进行单元测试和性能测试对比。

**验收标准**:

- [x] 响应处理逻辑完全独立，无循环依赖。
- [x] `ResponseHandler` 无需依赖 `RenderService`。
- [x] 性能提升 10%+。
- [x] 内存使用优化。
- [x] API 响应正常。

---

### 阶段5: Template 组件重构 (1-2天)

**目标**: 优化模板函数管理，提升模板渲染性能，并消除对 `RenderService` 的依赖。

#### 5.1 模板函数管理器

```go
// internal/app/render/template/manager.go
package template

import (
    "html/template"
    "log/slog"
    "strings"
    "sync"
    "time" // 用于FormatTime
    "github.com/biuren/hebe/internal/app/render/contracts"
    "github.com/biuren/hebe/internal/app/render" // 引用共享的render包类型
)

// Manager 模板函数管理器，使用Go 1.24的并发安全优化
type Manager struct {
    // 函数注册表，使用读写锁保护，支持并发读写
    funcs        sync.Map
    compiled     template.FuncMap // 预编译的函数映射，用于提高GetFuncMap的性能
    mu           sync.RWMutex     // 保护compiled字段
    localization contracts.LocalizationManagerInterface // 注入LocalizationManagerInterface
    logger       *slog.Logger
}

// NewManager 构造函数：接收必要的依赖，不再接收RenderService
func NewManager(loc contracts.LocalizationManagerInterface) *Manager {
    m := &Manager{
        localization: loc,
        logger:       slog.Default().With("component", "template_manager"),
    }
    // 注册核心模板函数
    m.RegisterCommonFuncs()
    return m
}

// RegisterFunc 注册单个模板函数 (Go 1.24泛型支持，如果需要)
func (tm *Manager) RegisterFunc(name string, fn any) {
    tm.funcs.Store(name, fn)
    tm.invalidateCompiled() // 函数注册后，需要使缓存失效
}

// GetFuncMap 获取编译后的函数映射 (懒加载+缓存)
func (tm *Manager) GetFuncMap() template.FuncMap {
    tm.mu.RLock()
    if tm.compiled != nil {
        defer tm.mu.RUnlock()
        return tm.compiled
    }
    tm.mu.RUnlock()

    return tm.buildFuncMap()
}

func (tm *Manager) buildFuncMap() template.FuncMap {
    tm.mu.Lock()
    defer tm.mu.Unlock()

    // 双重检查锁定
    if tm.compiled != nil {
        return tm.compiled
    }

    funcMap := make(template.FuncMap)
    tm.funcs.Range(func(key, value interface{}) bool {
        if name, ok := key.(string); ok {
            funcMap[name] = value
        }
        return true
    })

    tm.compiled = funcMap
    return funcMap
}

func (tm *Manager) invalidateCompiled() {
    tm.mu.Lock()
    defer tm.mu.Unlock()
    tm.compiled = nil
}

// RegisterCommonFuncs 注册通用的模板函数 (从原template_funcs.go迁移)
func (tm *Manager) RegisterCommonFuncs() {
    tm.RegisterFunc("renderBreadcrumb", tm.RenderBreadcrumb)
    tm.RegisterFunc("getMessage", tm.GetMessageFunc) // 包装，使用注入的LocalizationManager
    tm.RegisterFunc("formatTime", tm.FormatTime)
    tm.RegisterFunc("isLogin", func(data *render.PageData) bool { return data.IsAuthenticated() })
    tm.RegisterFunc("getLoginName", func(data *render.PageData) string {
        if user := data.GetAuthUser(); user != nil { return user.GetLoginName() }
        return ""
    })
    tm.RegisterFunc("getAccessUrl", func(data *render.PageData) string {
        if user := data.GetAuthUser(); user != nil { return user.GetAccessUrl() }
        return ""
    })
    tm.RegisterFunc("getUserID", func(data *render.PageData) int64 {
        if user := data.GetAuthUser(); user != nil { return user.GetUserID() }
        return 0
    })
    tm.RegisterFunc("getUserName", func(data *render.PageData) string {
        if user := data.GetAuthUser(); user != nil { return user.GetUserName() }
        return ""
    })
    tm.RegisterFunc("getFullName", func(data *render.PageData) string {
        if user := data.GetAuthUser(); user != nil { return user.GetFullName() }
        return ""
    })
    tm.RegisterFunc("getEmail", func(data *render.PageData) string {
        if user := data.GetAuthUser(); user != nil { return user.GetEmail() }
        return ""
    })
    tm.RegisterFunc("getUserGroup", func(data *render.PageData) int8 {
        if user := data.GetAuthUser(); user != nil { return user.GetUserGroup() }
        return 0
    })
    tm.RegisterFunc("getFaceUrl", func(data *render.PageData) string {
        if user := data.GetAuthUser(); user != nil { return user.GetFaceUrl() }
        return ""
    })
    tm.RegisterFunc("getLoginUTCTime", func(data *render.PageData) int64 {
        if user := data.GetAuthUser(); user != nil { return user.GetLoginUTCTime() }
        return 0
    })
    tm.RegisterFunc("getTimeZone", func(data *render.PageData) string {
        if user := data.GetAuthUser(); user != nil { return user.GetTimeZone() }
        return ""
    })
    tm.RegisterFunc("getUserLang", func(data *render.PageData) string {
        if user := data.GetAuthUser(); user != nil { return user.GetLang() }
        return ""
    })
    tm.RegisterFunc("getUserStatus", func(data *render.PageData) int8 {
        if user := data.GetAuthUser(); user != nil { return user.GetStatus() }
        return 0
    })
    tm.RegisterFunc("getIP", func(data *render.PageData) string {
        // 假设IP地址可能在Context中
        // if ip, ok := data.Context.Value("ip").(string); ok { return ip }
        return ""
    })
    tm.RegisterFunc("isMobile", func(data *render.PageData) bool {
        // 假设移动设备标识可能在Context中
        // if isMob, ok := data.Context.Value("is_mobile").(bool); ok { return isMob }
        return false
    })
    tm.RegisterFunc("formatFileSize", render.FormatFileSize) // 假设FormatFileSize成为纯函数或独立工具
    tm.RegisterFunc("truncateString", render.TruncateString) // 假设TruncateString成为纯函数或独立工具
    tm.RegisterFunc("tWithParams", tm.TWithParamsFunc) // 包装，用于处理带参数的本地化消息
}

// GetMessageFunc 包装器，用于模板函数中调用LocalizationManager
func (tm *Manager) GetMessageFunc(key, locale string) string {
    return tm.localization.GetMessage(key, locale)
}

// TWithParamsFunc 包装器，用于模板函数中调用带参数的本地化消息
func (tm *Manager) TWithParamsFunc(key string, params map[string]string, data *render.PageData) string {
    if data == nil || data.Context == nil { // 确保data和data.Context不为空
        return key // 无法获取语言信息时返回key
    }
    return tm.localization.GetMessageWithParams(key, data.Context.Lang, params)
}

// RenderBreadcrumb 示例（需要根据PageData结构调整，并确保引用render包下的类型）
func (tm *Manager) RenderBreadcrumb(breadcrumb []render.BreadcrumbItem) template.HTML {
    if len(breadcrumb) == 0 {
        return ""
    }
    var parts []string
    for i, item := range breadcrumb {
        if i == len(breadcrumb)-1 {
            parts = append(parts, fmt.Sprintf(`<span>%s</span>`, item.Name))
        } else if item.URL != "" {
            parts = append(parts, fmt.Sprintf(`<a href="%s">%s</a>`, item.URL, item.Name))
        } else {
            parts = append(parts, fmt.Sprintf(`<span>%s</span>`, item.Name))
        }
    }
    html := fmt.Sprintf(`<nav class="breadcrumb">%s</nav>`, strings.Join(parts, ` <span class="separator">/</span> `))
    return template.HTML(html)
}

// FormatTime 格式化时间 (从原template_funcs.go迁移)
func (tm *Manager) FormatTime(timestamp int64, format string, timezone ...string) string {
    t := time.Unix(timestamp/1000, (timestamp%1000)*1000000)
    if len(timezone) > 0 && timezone[0] != "" {
        if loc, err := time.LoadLocation(timezone[0]); err == nil {
            t = t.In(loc)
        }
    }
    return t.Format(format)
}
```

#### 5.2 执行步骤

1.  创建独立的 `template/` 包，并将 `template_funcs.go` 移动到其中。
2.  重命名 `template_funcs.go` 为 `manager.go`。
3.  更新 `template/manager.go` 中的 `NewManager` 方法，使其通过参数接收 `contracts.LocalizationManagerInterface`，并**不再接收 `*RenderService`**，消除循环依赖。
4.  实现并发安全的函数管理（`sync.Map` 和 `sync.RWMutex`）。
5.  优化模板编译性能（缓存 `template.FuncMap`）。
6.  将原 `template_funcs.go` 中的所有模板辅助函数迁移到 `manager.go` 或独立的 `functions.go` 中，并确保它们使用注入的依赖（如 `localization`）。
7.  修改 `service.go` 中 `templateFuncs` 的初始化方式，直接调用 `template.NewManager()`。
8.  添加模板渲染性能测试和单元测试。

**验收标准**:

- [x] 模板函数管理逻辑完全独立，无循环依赖。
- [x] `TemplateManager` 无需依赖 `RenderService`。
- [x] 模板渲染性能提升 15%+。
- [x] 并发安全。
- [x] 所有模板函数正常工作。

---

### 阶段6: 服务集成与全面 DI 优化 (1天)

**目标**: 整合所有组件，确保系统稳定运行，并通过 DI 框架进行组件的创建和依赖注入。

#### 6.1 服务集成器 (核心变更)

```go
// internal/app/render/service.go (最终优化版)
package render

import (
    "html/template" // 用于RegisterFuncs的返回类型
    "github.com/gin-gonic/gin"
    "github.com/biuren/hebe/internal/core/config"
    "github.com/biuren/hebe/internal/app/render/contracts" // 引入所有组件的contracts接口
    // 确保PageData等共享类型在当前包可访问或从共享包导入
)

// RenderService 最终重构后的渲染服务
type Service struct {
    config       *config.Config // 如果Service自身还需要访问config
    context      contracts.ContextBuilderInterface
    page         contracts.PageBuilderInterface
    response     contracts.ResponseHandlerInterface
    template     contracts.TemplateManagerInterface
    localization contracts.LocalizationManagerInterface
    // 如果utils的功能被抽取为独立纯函数，则Service不再需要持有其引用
    // 如果部分utils功能被封装为接口并需要注入，则在此处添加
    // utils        contracts.RenderUtilsInterface
}

// NewService 通过DI框架注入所有依赖接口
// 这是最核心的优化：RenderService 不再手动创建其依赖，而是声明它需要哪些接口。
// 所有的具体实现和依赖关系由 DI 框架在应用启动时负责创建和组装。
func NewService(
    cfg *config.Config,
    contextBuilder contracts.ContextBuilderInterface,
    pageBuilder contracts.PageBuilderInterface,
    responseHandler contracts.ResponseHandlerInterface,
    templateManager contracts.TemplateManagerInterface,
    localizationManager contracts.LocalizationManagerInterface,
    // 如果有，则添加 utilsInterface contracts.RenderUtilsInterface,
) RenderServiceInterface {
    return &Service{
        config:       cfg,
        context:      contextBuilder,
        page:         pageBuilder,
        response:     responseHandler,
        template:     templateManager,
        localization: localizationManager,
        // utils:        utilsInterface, // 如果有
    }
}

// 保持所有现有RenderServiceInterface的方法不变，它们将委托给内部注入的组件接口
func (s *Service) PageData(c *gin.Context, titleKey string) *PageData {
    return s.page.BuildPageData(c, titleKey)
}

func (s *Service) JSON(c *gin.Context, code int, data interface{}) {
    s.response.JSON(c, code, data)
}

func (s *Service) JSONSuccess(c *gin.Context, data interface{}) {
    s.response.JSONSuccess(c, data)
}

func (s *Service) JSONError(c *gin.Context, err error, data ...interface{}) {
    s.response.JSONError(c, err, data...)
}

func (s *Service) HTML(c *gin.Context, name string, data interface{}) {
    s.response.HTML(c, name, data)
}

func (s *Service) HTMLError(c *gin.Context, title, message string) {
    s.response.HTMLError(c, title, message)
}

func (s *Service) Redirect(c *gin.Context, url string) {
    s.response.Redirect(c, url)
}

// RegisterFuncs 不再直接在Service中注册，而是通过TemplateManager接口提供
func (s *Service) RegisterFuncs() template.FuncMap {
    return s.template.GetFuncMap()
}

// GetLocaleFromContext 方法不再由Service直接实现，而是委托给LocalizationManager
func (s *Service) GetLocaleFromContext(c *gin.Context) string {
    return s.localization.GetLocaleFromContext(c)
}

// GetUserLocale 方法不再由Service直接实现，而是委托给LocalizationManager
func (s *Service) GetUserLocale(user RenderUser) string {
    return s.localization.GetUserLocale(user)
}
// ... 其他方法保持不变，或根据实际情况委托给内部组件 ...
```

#### 6.2 DI 容器注册示例 (关键步骤)

在应用启动的DI容器配置部分，注册所有 `render` 相关的组件。这展示了所有组件是如何被“装配”起来的。

```go
// 伪代码：在你的DI容器初始化部分 (例如: internal/app/bootstrap/dependencies.go)
container := di.New()

// 注册 Config (假设已存在)
di.Singleton(container, func() *config.Config { return appConfig })

// 1. 注册 LocalizationManager (注意错误处理)
di.Singleton(container, func() (contracts.LocalizationManagerInterface, error) {
    cfg := di.Get[*config.Config](container)
    // 严格错误处理：如果本地化文件加载失败，应中断启动
    return localization.NewManager(cfg)
})

// 2. 注册 ContextBuilder
di.Singleton(container, func() contracts.ContextBuilderInterface {
    loc := di.Get[contracts.LocalizationManagerInterface](container)
    return context.NewBuilder(loc)
})

// 3. 注册 PageBuilder
di.Singleton(container, func() contracts.PageBuilderInterface {
    cfg := di.Get[*config.Config](container)
    ctxBuilder := di.Get[contracts.ContextBuilderInterface](container)
    loc := di.Get[contracts.LocalizationManagerInterface](container)
    return page.NewBuilder(cfg, ctxBuilder, loc)
})

// 4. 注册 ResponseHandler
di.Singleton(container, func() contracts.ResponseHandlerInterface {
    return response.NewHandler()
})

// 5. 注册 TemplateManager
di.Singleton(container, func() contracts.TemplateManagerInterface {
    loc := di.Get[contracts.LocalizationManagerInterface](container)
    return template.NewManager(loc)
})

// 6. 注册最终的 RenderService
di.Singleton(container, func() RenderServiceInterface {
    return NewService(
        di.Get[*config.Config](container),
        di.Get[contracts.ContextBuilderInterface](container),
        di.Get[contracts.PageBuilderInterface](container),
        di.Get[contracts.ResponseHandlerInterface](container),
        di.Get[contracts.TemplateManagerInterface](container),
        di.Get[contracts.LocalizationManagerInterface](container),
    )
})
```

#### 6.3 `internal/app/render/localization/manager.go` (考虑错误处理)

```go
// internal/app/render/localization/manager.go (核心逻辑不变，但思考NewManager的错误处理策略)
package localization

import (
    "encoding/json"
    "fmt"
    "log/slog" // Go 1.24 结构化日志
    "os"
    "path/filepath"
    "strings"

    "github.com/biuren/hebe/internal/core/config"
    "github.com/gin-gonic/gin"
    // 假设RenderUser接口在render包中
    "github.com/biuren/hebe/internal/app/render"
)

type Manager struct {
    cfg      *config.Config
    messages map[string]map[string]string
    logger   *slog.Logger
}

// NewManager 创建并初始化本地化管理器。
// 优化考虑：对于关键的默认语言加载失败，应该返回错误，阻止应用启动。
// 对于非默认语言，可以记录警告并继续。
func NewManager(cfg *config.Config) (*Manager, error) {
    manager := &Manager{
        cfg:      cfg,
        messages: make(map[string]map[string]string),
        logger:   slog.Default().With("component", "localization_manager"),
    }

    var loadedCount int
    // 首先尝试加载默认语言，如果失败则直接返回错误
    if err := manager.loadLocale(cfg.App.DefaultLang); err != nil {
        return nil, fmt.Errorf("无法加载默认语言 (%s)：%w，请检查语言文件是否存在", cfg.App.DefaultLang, err)
    }
    loadedCount++
    manager.logger.Info("默认语言加载成功", "locale", cfg.App.DefaultLang)

    // 然后加载其他支持的语言，如果失败则记录警告
    for _, supportedLocale := range cfg.SupportLangs {
        if supportedLocale.Code == cfg.App.DefaultLang { // 默认语言已加载
            continue
        }
        if err := manager.loadLocale(supportedLocale.Code); err != nil {
            manager.logger.Warn("加载语言文件失败", "locale", supportedLocale.Code, "error", err)
            continue // 继续加载其他语言
        }
        loadedCount++
    }

    if loadedCount == 0 { // 理论上不会发生，因为默认语言已经确保加载
        return nil, fmt.Errorf("没有成功加载任何语言文件，这不应该发生")
    }

    return manager, nil
}

// loadLocale ... (保持原逻辑不变)
func (m *Manager) loadLocale(localeCode string) error {
    filePath := filepath.Join(m.cfg.App.LangsDir, fmt.Sprintf("%s.json", localeCode))
    data, err := os.ReadFile(filePath)
    if err != nil {
        return fmt.Errorf("读取语言文件失败: %w", err)
    }

    var messages map[string]string
    if err := json.Unmarshal(data, &messages); err != nil {
        return fmt.Errorf("解析语言文件失败: %w", err)
    }

    m.messages[localeCode] = messages
    return nil
}

// GetMessage ... (保持原逻辑不变)
func (m *Manager) GetMessage(key, locale string) string {
    if langMessages, ok := m.messages[locale]; ok {
        if msg, ok := langMessages[key]; ok {
            return msg
        }
    }
    // Fallback to default language if message not found in specified locale
    if defaultLangMessages, ok := m.messages[m.cfg.App.DefaultLang]; ok {
        if msg, ok := defaultLangMessages[key]; ok {
            return msg
        }
    }
    return key // Return key if no message found
}

// GetMessageWithParams 处理带参数的本地化消息
// 示例：message "Hello, {{name}}!" and params {"name": "World"} -> "Hello, World!"
func (m *Manager) GetMessageWithParams(key, locale string, params map[string]string) string {
    msg := m.GetMessage(key, locale)
    for param, value := range params {
        msg = strings.ReplaceAll(msg, fmt.Sprintf("{{%s}}", param), value)
    }
    return msg
}


// GetLocaleFromContext ... (保持原逻辑不变)
func (m *Manager) GetLocaleFromContext(c *gin.Context) string {
    // 优先级1：从URL查询参数获取语言 (e.g., ?lang=en)
    if lang := c.Query("lang"); lang != "" && m.isValidLocale(lang) {
        return lang
    }

    // 优先级2：从Gin Context获取（中间件注入）
    if lang, exists := c.Get("language"); exists {
        if langStr, ok := lang.(string); ok && langStr != "" {
            return langStr
        }
    }

    // 优先级3：从用户会话获取 (假设AuthUser实现了GetLang())
    if authUserRaw, exists := c.Get("AuthUser"); exists && authUserRaw != nil {
        if user, ok := authUserRaw.(render.RenderUser); ok {
            if userLang := user.GetLang(); userLang != "" && m.isValidLocale(userLang) {
                return userLang
            }
        }
    }

    // 优先级4：浏览器偏好 (Accept-Language)
    if browserLocale := m.getBrowserPreferredLocale(c); browserLocale != "" {
        return browserLocale
    }

    // 优先级5：系统默认
    return m.cfg.App.DefaultLang
}

// GetUserLocale 从 RenderUser 获取用户语言
func (m *Manager) GetUserLocale(user render.RenderUser) string {
    if user != nil {
        if userLang := user.GetLang(); userLang != "" && m.isValidLocale(userLang) {
            return userLang
        }
    }
    return m.cfg.App.DefaultLang
}

// isValidLocale ... (保持原逻辑不变)
func (m *Manager) isValidLocale(locale string) bool {
    _, ok := m.messages[locale]
    return ok
}

// getBrowserPreferredLocale ... (保持原逻辑不变)
func (m *Manager) getBrowserPreferredLocale(c *gin.Context) string {
    acceptLang := c.GetHeader("Accept-Language")
    if acceptLang == "" {
        return ""
    }

    languages := m.parseAcceptLanguage(acceptLang)

    for _, lang := range languages {
        if m.isValidLocale(lang) { // 直接匹配
            return lang
        }
        if normalized := m.normalizeLocale(lang); normalized != "" && m.isValidLocale(normalized) { // 尝试标准化匹配 (如 en-US -> en)
            return normalized
        }
    }
    return ""
}

// normalizeLocale ... (保持原逻辑不变)
func (m *Manager) normalizeLocale(locale string) string {
    parts := strings.Split(strings.ToLower(locale), "-")
    if len(parts) == 1 {
        return parts[0]
    }
    if len(parts) >= 2 {
        // 处理特殊情况：zh-CN -> zh-cn
        return parts[0] + "-" + parts[1]
    }
    return ""
}

// parseAcceptLanguage ... (保持原逻辑不变)
func (m *Manager) parseAcceptLanguage(acceptLang string) []string {
    // 假设这里有一个更完善的解析逻辑，例如处理q值
    return strings.Split(acceptLang, ",")
}
```

#### 6.4 `internal/app/render/utils/utils.go` (优化考虑：纯函数化或接口化)

原 `utils.go` 中的许多函数（如 `FormatFileSize`, `TruncateString`, `SafeHTML`, 时间戳处理）都是**纯函数**或**几乎无状态**的工具函数。

**优化建议**：

1.  **纯函数化**：将这些函数从 `RenderUtils` 结构体中解耦出来，变为普通的包级别函数（例如 `render.FormatFileSize`）。它们不需要任何依赖。
2.  **接口化**：对于少量确实需要依赖（如 `slog`）或需要保持状态（如果存在）的工具，可以为其定义接口 `contracts.RenderUtilsInterface`，并创建 `utils/manager.go` 实现该接口，然后通过 DI 注入。

**示例（纯函数化）：**

```go
// internal/app/render/utils.go (建议直接放在render包下，或一个共享的util包)
package render // 或 package internal/utils

import (
    "bytes"
    "encoding/json"
    "fmt"
    "html/template"
    "net/http"
    "path/filepath"
    "strings"
    "time"
)

// SafeHTML 安全的HTML输出
func SafeHTML(content string) template.HTML {
    return template.HTML(content)
}

// FormatFileSize 格式化文件大小
func FormatFileSize(size int64) string {
    const unit = 1024
    if size < unit {
        return fmt.Sprintf("%d B", size)
    }
    div, exp := int64(unit), 0
    for n := size / unit; n >= unit; n /= unit {
        div *= unit
        exp++
    }
    return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}

// TruncateString 截断字符串
func TruncateString(str string, maxLength int) string {
    if len(str) <= maxLength {
        return str
    }
    return str[:maxLength] + "..."
}

// FormatDate 格式化日期
func FormatDate(timestamp int64, format string) string {
    return time.Unix(timestamp, 0).Format(format)
}

// HtmlUnescape HTML解码
func HtmlUnescape(str string) string {
    str = strings.ReplaceAll(str, "&amp;", "&")
    str = strings.ReplaceAll(str, "&lt;", "<")
    str = strings.ReplaceAll(str, "&gt;", ">")
    str = strings.ReplaceAll(str, "&quot;", "\"")
    str = strings.ReplaceAll(str, "&#39;", "'")
    return str
}

// GetCurrentTimestamp 获取当前时间戳 (秒)
func GetCurrentTimestamp() int64 {
    return time.Now().Unix()
}

// GetCurrentTimestampMs 获取当前时间戳（毫秒）
func GetCurrentTimestampMs() int64 {
    return time.Now().UnixNano() / int64(time.Millisecond)
}

// ParseTimestamp 解析时间戳 (秒)
func ParseTimestamp(timestamp int64) time.Time {
    return time.Unix(timestamp, 0)
}

// ParseTimestampMs 解析时间戳（毫秒）
func ParseTimestampMs(timestamp int64) time.Time {
    return time.Unix(timestamp/1000, (timestamp%1000)*1000000)
}

// GenerateRandomString 生成随机字符串
func GenerateRandomString(length int) string {
    const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    b := make([]byte, length)
    for i := range b {
        // NOTE: time.Now().UnixNano() 作为伪随机源不够安全，仅用于简单场景
        // 生产环境应使用 crypto/rand
        b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
    }
    return string(b)
}

// HttpGetJSON 用于从指定URL进行GET请求并解析JSON响应
// 注意：这个函数可能需要依赖 HTTP 客户端（如 http.Client），
// 如果要保持 Render 模块的纯粹性，这个功能可能更适合放在一个独立的HTTP客户端包中，
// 或者作为 RenderUtilsInterface 的一个方法，并注入 http.Client。
func HttpGetJSON(url string, target interface{}) error {
    resp, err := http.Get(url)
    if err != nil {
        return fmt.Errorf("HTTP GET request failed: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("HTTP GET request returned status code: %d", resp.StatusCode)
    }

    decoder := json.NewDecoder(resp.Body)
    if err := decoder.Decode(target); err != nil {
        return fmt.Errorf("Failed to decode JSON response: %w", err)
    }
    return nil
}

// IsFileExists 检查文件是否存在
func IsFileExists(filePath string) bool {
    _, err := os.Stat(filePath)
    return !os.IsNotExist(err)
}

// GetFileExtension 获取文件扩展名
func GetFileExtension(fileName string) string {
    return strings.TrimPrefix(filepath.Ext(fileName), ".")
}

// BytesToString 将字节切片转换为字符串
func BytesToString(data []byte) string {
    return string(data)
}
```

#### 6.4 执行步骤

1.  **DI 框架集成**：在应用程序的启动入口（通常是 `main` 函数或 DI 容器的配置层），负责创建所有组件的具体实现，并通过 DI 框架将它们注入到 `RenderService` 的 `NewService` 构造函数中。
2.  **`NewRenderService` 签名变更**：修改 `internal/app/render/service.go` 中的 `NewService` 函数签名，使其接收 `contracts` 包中定义的接口，而不是手动创建。
3.  **`utils.go` 拆分**：将原 `utils.go` 中的纯函数抽取出来，作为包级别函数，放置在 `internal/app/render/utils.go` 或更通用的 `internal/utils/` 包中。对于需要依赖的函数（如 `HttpGetJSON`），考虑为其创建独立的接口和实现，并通过 DI 注入。
4.  **测试与验证**：
    - 进行全面的**集成测试**，确保所有功能在新的架构下协同工作。
    - 运行**性能基准测试**，确认达到预期的性能提升。
    - 验证**内存使用**是否降低。
    - 更新**文档**，反映最终的架构和实现。

**验收标准**:

- [x] 所有功能正常，无功能回归。
- [x] 整个 `Render` 模块通过 DI 框架进行组装，实现高度解耦。
- [x] 性能提升 20%+，内存使用降低 15%+。
- [x] 测试覆盖率 \> 90%。
- [x] 代码结构清晰，符合最终目标架构。

---

## 最终目标架构 (更新版)

```
internal/app/render/
├── service.go              # 轻量级服务集成器 (~50行，通过DI注入依赖)
├── contracts/              # 接口定义
│   └── interfaces.go
├── context/                # 上下文构建
│   ├── builder.go
│   └── builder_test.go
├── page/                   # 页面构建
│   ├── builder.go
│   ├── processors.go       # 独立的数据处理器实现
│   └── builder_test.go
├── response/               # 响应处理
│   ├── handler.go
│   └── handler_test.go
├── template/               # 模板管理
│   ├── manager.go
│   ├── functions.go (可选，如果通用函数过多可独立)
│   └── manager_test.go
├── localization/           # 本地化管理
│   ├── manager.go
│   └── manager_test.go
├── types.go               # 共享类型定义 (PageData, RenderContext, RenderUser等)
└── utils.go               # 纯函数工具类 (已从RenderService解耦)
```

## 预期收益

1.  **性能显著提升**: 整体性能提升20%+，内存使用降低15%+，得益于组件优化、Go 1.24特性和更高效的资源管理。
2.  **极致可维护性**: 组件解耦，单一职责，依赖清晰。每个组件独立可测，降低了维护成本和引入 Bug 的风险。
3.  **更强的可测试性**: 所有组件都通过接口进行依赖，使得单元测试和集成测试的编写变得简单和高效，可以轻松使用 Mock 对象。
4.  **类型安全与现代化**: 充分利用 Go 1.24 泛型，提升代码的类型安全和表达力。
5.  **开发效率提升**: 模块化设计便于团队并行开发，组件的独立性减少了相互影响。
6.  **代码质量优化**: 通过强制接口依赖和 DI，鼓励更好的设计实践，测试覆盖率可达90%+。

---
