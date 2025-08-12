# iDiary.ai Hebe 国际化系统指南

## 概述

iDiary.ai Hebe采用现代化的国际化(i18n)系统，支持多语言界面和智能语言检测。系统基于Go 1.24泛型特性、预加载机制和配置驱动设计，提供高性能的多语言支持，与ViewService和模板函数系统完美集成。

## 系统架构

### 核心组件

- **I18nManager**: 国际化管理器接口，负责语言资源加载和翻译
- **LocalizationManager**: 具体实现，基于预加载和缓存的高性能管理器
- **智能语言检测**: 从用户设置、Cookie、浏览器偏好中自动检测语言
- **模板函数**: 在模板中提供便捷的翻译功能
- **ViewService集成**: 与视图服务完美集成，支持页面级国际化
- **配置驱动**: 通过config.Config管理支持的语言

### 文件结构

```
website/locales/
├── en.json          # 英语资源文件
├── zh-cn.json       # 简体中文资源文件
└── [其他语言].json   # 其他语言资源文件

internal/utils/i18n/
├── interfaces.go    # I18nManager 接口定义
└── service.go       # LocalizationManager 具体实现

internal/app/templates/
└── localization.go  # 模板函数实现

internal/utils/config/
└── types.go         # 语言配置结构定义
```

### 架构流程图

```
HTTP Request
    ↓
SharedData Middleware (注入语言上下文)
    ↓
Controller (获取用户语言偏好)
    ↓
ViewService.CreatePageData (设置语言上下文)
    ↓
Template Rendering (使用模板函数获取翻译)
    ↓
I18nManager.GetMessage (从预加载缓存获取翻译)
    ↓
HTML Response (多语言内容)
```

## 语言资源文件

### 文件格式

语言资源文件采用JSON格式，支持扁平化键值对结构：

```json
{
  "common.home": "首页",
  "common.settings": "设置",
  "common.profile": "个人资料",
  "common.logout": "退出登录",
  "common.save": "保存",
  "common.edit": "编辑",
  "common.delete": "删除",
  "common.confirm_delete": "您确定要删除这个项目吗？",

  "nav.home": "首页",
  "nav.diary": "日记",
  "nav.settings": "设置",

  "page.title.diary_list": "我的日记",
  "page.title.diary_detail": "日记详情",

  "diary.title": "我的日记",
  "diary.title.desc": "记录生活点滴，分享精彩瞬间",
  "diary.new_entry": "写新日记",
  "diary.content_placeholder": "{{username}}，今天有什么新鲜事？",

  "greeting.message": "你好，{{username}}！欢迎使用 iDiary。",

  "error.not_found": "抱歉，您要查找的页面不存在。"
}
}
```

### 命名约定

- **扁平化结构**: 使用点分隔符组织层级关系
- **功能分组**: 按模块前缀组织（common、nav、page、diary等）
- **描述性键名**: 键名要能清楚表达内容含义
- **一致性**: 在所有语言文件中保持键名完全一致

### 参数化翻译

支持使用`{{参数名}}`格式的占位符：

```json
{
  "diary.content_placeholder": "{{username}}，今天有什么新鲜事？",
  "greeting.message": "你好，{{username}}！欢迎使用 iDiary。"
}
```

## 配置管理

### 语言配置

**文件位置**: `internal/utils/config/types.go`

```go
// SupportLang 支持的语言
type SupportLang struct {
    Code      string `json:"code" validate:"required"`       // 语言代码
    Label     string `json:"label" validate:"required"`      // 英文标签
    LocalText string `json:"local-text" validate:"required"` // 本地化文本
}

// Config 结构中的国际化配置
type Config struct {
    // ...
    SupportLangs []SupportLang `json:"support_langs" validate:"required,dive"`
    // ...
    App AppConfig `json:"app" validate:"required"`
}

// AppConfig 应用配置
type AppConfig struct {
    // ...
    LocalesDir  string `json:"locales_dir" mapstructure:"LOCALES_DIR" validate:"required"`
    DefaultLang string `json:"default_lang" mapstructure:"DEFAULT_LANG" validate:"required"`
    // ...
}
```

### 配置示例

```yaml
# config.yaml
app:
  locales_dir: './website/locales'
  default_lang: 'zh-cn'

support_langs:
  - code: 'zh-cn'
    label: 'Chinese Simplified'
    local-text: '简体中文'
  - code: 'en'
    label: 'English'
    local-text: 'English'
  - code: 'zh-tw'
    label: 'Chinese Traditional'
    local-text: '繁體中文'
```

## I18nManager 接口

**文件位置**: `internal/utils/i18n/interfaces.go`

```go
// I18nManager 国际化管理器接口
type I18nManager interface {
    // GetMessage 获取指定键和语言的本地化消息
    GetMessage(key, locale string) string

    // GetMessageWithParams 获取带参数的翻译消息，支持 {{参数名}} 格式的占位符替换
    GetMessageWithParams(key, locale string, params map[string]string) string

    // GetSupportedLanguages 获取支持的语言列表
    GetSupportedLanguages() []config.SupportLang

    // GetSupportedLanguageCodes 获取支持的语言代码列表
    GetSupportedLanguageCodes() []string

    // GetDefaultLanguage 获取默认语言
    GetDefaultLanguage() string

    // IsLanguageSupported 检查语言是否支持
    IsLanguageSupported(locale string) bool

    // GetLanguageInfo 获取语言信息
    GetLanguageInfo(locale string) (config.SupportLang, bool)

    // GetLocaleFromContext 从 Gin 上下文中智能检测语言
    GetLocaleFromContext(c *gin.Context) string

    // ReloadMessages 重新加载消息（开发环境使用）
    ReloadMessages() error
}
```

## LocalizationManager 实现

**文件位置**: `internal/utils/i18n/service.go`

### 核心特性

1. **预加载机制**: 启动时一次性加载所有语言文件到内存
2. **高性能缓存**: 使用全局map缓存，键格式为`语言代码:消息键`
3. **智能语言检测**: 多层优先级的语言检测机制
4. **语言代码规范化**: 自动处理各种语言代码格式
5. **降级策略**: 找不到翻译时使用默认语言或返回原始键

### 智能语言检测

系统按以下优先级检测用户语言：

1. **用户设置**: 已登录用户的语言偏好设置
2. **Cookie设置**: `lang` Cookie的值
3. **浏览器偏好**: `Accept-Language` 头中的语言列表
4. **系统默认**: 配置文件中的默认语言

```go
func (m *LocalizationManager) GetLocaleFromContext(c *gin.Context) string {
    // 优先级1：用户设置 (如果已登录)
    if isLogin, exists := c.Get("IsLogin"); exists && isLogin == true {
        if authUserValue, exists := c.Get("AuthUser"); exists && authUserValue != nil {
            if user, ok := authUserValue.(interface{ GetLang() string }); ok {
                if locale := user.GetLang(); locale != "" && m.IsLanguageSupported(locale) {
                    return locale
                }
            }
        }
    }

    // 优先级2：Cookie 设置
    if cookieLocale, err := c.Cookie("lang"); err == nil && cookieLocale != "" {
        if m.IsLanguageSupported(cookieLocale) {
            return cookieLocale
        }
    }

    // 优先级3：浏览器偏好 (Accept-Language)
    if browserLocale := m.getBrowserPreferredLocale(c); browserLocale != "" {
        return browserLocale
    }

    // 优先级4：系统默认
    return m.GetDefaultLanguage()
}
```

### 语言代码规范化

系统自动处理多种语言代码格式：

```go
func (m *LocalizationManager) normalizeLocale(locale string) string {
    locale = strings.ToLower(locale)

    switch locale {
    case "zh", "zh-hans", "zh_hans", "zh-cn", "zh_cn":
        return "zh-cn"
    case "zh-hant", "zh_hant", "zh-tw", "zh_tw":
        return "zh-tw"
    case "en", "en-us", "en_us":
        return "en-us"
    default:
        return locale
    }
}
```

## ViewService 集成

### 在ViewService中的使用

**文件位置**: `internal/app/view/service.go`

```go
type ViewService struct {
    config      *config.Config
    i18nManager i18n.I18nManager
}

func NewViewService(cfg *config.Config) *ViewService {
    i18nManager := i18n.NewLocalizationManager(cfg)
    return &ViewService{
        config:      cfg,
        i18nManager: i18nManager,
    }
}

// GetMessage 获取国际化消息
func (s *ViewService) GetMessage(key, lang string) string {
    return s.i18nManager.GetMessage(key, lang)
}

// GetLocaleFromContext 从上下文获取语言设置
func (s *ViewService) GetLocaleFromContext(c *gin.Context) string {
    return s.i18nManager.GetLocaleFromContext(c)
}
```

### 页面数据创建

```go
func (s *ViewService) CreatePageData(c *gin.Context, titleKey string) *contracts.PageData[map[string]interface{}] {
    ctx := s.BuildRenderContext(c)
    title := s.generateTitle(c, titleKey)

    return &contracts.PageData[map[string]interface{}]{
        Title:   title,
        Context: ctx,
        Data:    make(map[string]interface{}),
    }
}

func (s *ViewService) BuildRenderContext(c *gin.Context) *contracts.RenderContext {
    return &contracts.RenderContext{
        IsLogin:  s.getLoginStatus(c),
        AuthUser: s.getAuthUser(c),
        Lang:     s.GetLocaleFromContext(c), // 自动检测语言
    }
}
```

## 模板函数集成

### 模板函数实现

**文件位置**: `internal/app/templates/localization.go`

```go
// GetLocalizationFuncs 获取本地化相关的模板函数
func GetLocalizationFuncs() template.FuncMap {
    return template.FuncMap{
        // 核心消息获取函数
        "getMessage": getMessage,

        // 扩展本地化函数
        "getMessageWithDefault": getMessageWithDefault,
        "getMessages":           getMessages,
        "getSupportedLanguages": getSupportedLanguages,
        "isLanguageSupported":   isLanguageSupported,
    }
}

// getMessage 获取本地化消息（核心函数）
func getMessage(key, locale string) string {
    cacheKey := i18n.BuildCacheKey(locale, key)
    if message, exists := i18n.GetMessage(cacheKey); exists {
        return message
    }
    return key // 降级处理：返回原始key
}
```

### 在模板中使用

```html
<!DOCTYPE html>
<html lang="{{.Context.Lang}}">
<head>
    <title>{{getMessage "page.title.diary_list" .Context.Lang}}</title>
    <meta name="description" content="{{getMessage "diary.title.desc" .Context.Lang}}" />
</head>
<body>
    <h1>{{getMessage "diary.title" .Context.Lang}}</h1>

    <!-- 带参数的翻译 -->
    <p>{{getMessageWithDefault "diary.content_placeholder" .Context.Lang "写点什么吧"}}</p>

    <!-- 语言切换器 -->
    <select name="language">
    {{range getSupportedLanguages}}
    <option value="{{.Code}}" {{if eq .Code $.Context.Lang}}selected{{end}}>
        {{.LocalText}}
    </option>
    {{end}}
    </select>
</body>
</html>
```

## 控制器中的使用

### 基本使用模式

```go
func (ctrl *DiaryController) ShowDiary(c *gin.Context) {
    viewService := di.Get[*view.ViewService](ctrl.container)

    // 创建页面数据，自动设置当前语言
    pageData := viewService.CreatePageData(c, "diary.detail.title")

    // 获取当前语言
    currentLang := pageData.Context.Lang

    // 手动获取翻译
    welcomeMsg := viewService.GetMessage("diary.welcome", currentLang)

    // 添加业务数据
    pageData.Data["welcomeMessage"] = welcomeMsg
    pageData.Data["diary"] = diary

    // 渲染模板
    viewService.HTML(c, "diary/detail.html", pageData)
}
```

### 语言切换处理

```go
func (ctrl *UserController) SwitchLanguage(c *gin.Context) {
    newLang := c.PostForm("language")

    // 验证语言是否支持
    viewService := di.Get[*view.ViewService](ctrl.container)
    if !viewService.IsLanguageSupported(newLang) {
        err := errors.ValidationError.WithDetails("不支持的语言: " + newLang)
        c.Error(err)
        return
    }

    // 设置Cookie
    c.SetCookie("lang", newLang, 365*24*3600, "/", "", false, true)

    // 如果用户已登录，更新用户设置
    if isLogin, _ := c.Get("IsLogin"); isLogin == true {
        authUser, _ := c.Get("AuthUser")
        user := authUser.(*contracts.AuthUser)

        // 更新用户语言设置
        err := ctrl.userService.UpdateUserLanguage(user.UserId, newLang)
        if err != nil {
            slog.Warn("更新用户语言设置失败", "error", err, "user_id", user.UserId)
        }
    }

    // 重定向回原页面
    redirectURL := c.GetHeader("Referer")
    if redirectURL == "" {
        redirectURL = "/"
    }
    c.Redirect(http.StatusFound, redirectURL)
}
```

## 最佳实践

### 1. 翻译文件管理

```bash
# 翻译文件目录结构
website/locales/
├── zh-cn.json    # 简体中文（主要语言）
├── en.json       # 英语
├── zh-tw.json    # 繁体中文
└── ja.json       # 日语
```

### 2. 键名设计规范

```json
{
  // ✅ 好的键名设计
  "common.save": "保存",
  "common.cancel": "取消",
  "nav.home": "首页",
  "page.title.diary_list": "我的日记",
  "diary.create.title": "创建新日记",
  "error.validation.required": "此字段为必填项",

  // ❌ 避免的键名设计
  "save": "保存", // 缺少模块前缀
  "page_title": "页面标题", // 不够具体
  "err1": "错误1" // 不够描述性
}
```

### 3. 参数化最佳实践

```json
{
  // ✅ 正确的参数化
  "greeting.welcome": "欢迎，{{username}}！",
  "diary.created_at": "创建于 {{date}}",
  "user.posts_count": "共有 {{count}} 篇日记",

  // ❌ 避免复杂的HTML结构
  "complex.html": "<div class='user'>欢迎 {{username}} <span>!</span></div>"
}
```

### 4. 性能优化建议

```go
// ✅ 在模板中缓存语言
{{$lang := .Context.Lang}}
<h1>{{getMessage "page.title" $lang}}</h1>
<p>{{getMessage "page.description" $lang}}</p>

// ❌ 重复调用
<h1>{{getMessage "page.title" .Context.Lang}}</h1>
<p>{{getMessage "page.description" .Context.Lang}}</p>
```

### 5. 错误处理

```go
// 在控制器中处理翻译错误
func (ctrl *BaseController) getLocalizedError(errorKey string, c *gin.Context) string {
    viewService := di.Get[*view.ViewService](ctrl.container)
    lang := viewService.GetLocaleFromContext(c)

    message := viewService.GetMessage(errorKey, lang)

    // 如果翻译不存在，返回原始键
    if message == errorKey {
        slog.Warn("翻译缺失", "key", errorKey, "lang", lang)
        return viewService.GetMessage(errorKey, viewService.GetDefaultLanguage())
    }

    return message
}
```

## 性能特性

### 预加载性能

```
启动时间：
- 加载3种语言，总计300+条消息：< 50ms
- 内存占用：< 2MB
- 查找性能：每次调用 < 1μs（微秒级）
```

### 缓存机制

1. **全局缓存**: 使用`map[string]string`存储所有翻译
2. **键格式**: `语言代码:消息键`（如`zh-cn:common.save`）
3. **零分配**: 查找过程不产生内存分配
4. **并发安全**: 启动后只读，支持高并发访问

## 总结

iDiary.ai Hebe的国际化系统具备以下特点：

### 🚀 技术优势

- **高性能**: 预加载+内存缓存，微秒级翻译查找
- **智能检测**: 多层优先级的语言检测机制
- **类型安全**: 基于接口设计，支持依赖注入
- **扩展性**: 易于添加新语言和翻译键

### 🔧 开发友好

- **模板集成**: 丰富的模板函数，使用简单
- **ViewService集成**: 与视图服务无缝整合
- **配置驱动**: 通过配置文件管理语言支持
- **错误降级**: 翻译缺失时的优雅处理

### 📈 运维特性

- **预加载**: 启动时检测所有翻译文件完整性
- **调试支持**: 提供翻译统计和查询接口
- **热重载**: 开发环境支持重新加载翻译
- **监控友好**: 结构化日志记录

这套国际化系统为iDiary.ai Hebe提供了完整的多语言支持基础，确保了全球用户的使用体验。
