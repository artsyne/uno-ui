# iDiary.ai Hebe 国际化 (i18n) 指南

iDiary.ai Hebe 采用现代化的国际化(i18n)系统，支持多语言界面和智能语言检测。系统基于 Go 1.25 泛型特性、预加载机制和配置驱动设计，提供高性能的多语言支持，与 ViewService 和模板函数系统完美集成。

## 1. 系统架构

### 核心组件

- **I18nManager**: 国际化管理器接口，负责语言资源加载和翻译。
- **LocalizationManager**: 具体实现，基于预加载和缓存的高性能管理器。
- **智能语言检测**: 从用户设置、Cookie、浏览器偏好中自动检测语言。
- **模板函数**: 在模板中提供便捷的翻译功能。
- **ViewService集成**: 与视图服务完美集成，支持页面级国际化。
- **配置驱动**: 通过 `config.Config` 管理支持的语言。

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

## 2. 语言资源文件

### 文件格式

语言资源文件采用 JSON 格式，支持扁平化键值对结构：

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
```

### 命名约定

- **扁平化结构**: 使用点分隔符组织层级关系。
- **功能分组**: 按模块前缀组织（common、nav、page、diary 等）。
- **描述性键名**: 键名要能清楚表达内容含义。
- **一致性**: 在所有语言文件中保持键名完全一致。

### 参数化翻译

支持使用 `{{参数名}}` 格式的占位符：

```json
{
  "diary.content_placeholder": "{{username}}，今天有什么新鲜事？",
  "greeting.message": "你好，{{username}}！欢迎使用 iDiary。"
}
```

## 3. 配置管理

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

## 4. I18nManager 接口

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

## 5. 模板函数调用标准化

```html
<!-- 标准调用格式 -->
{{getMessage "key" .Context.Lang}}

<!-- 带默认值的调用 -->
{{getMessageWithDefault "key" .Context.Lang "默认文本"}}

<!-- 带参数的调用 -->
{{getMessageWithParams "greeting.welcome" .Context.Lang (dict "username" .Context.AuthUser.UserName)}}
```

## 6. 维护指南

### 添加新的多语言文本
1. 在 `zh-cn.json` 中添加中文文本。
2. 在 `en.json` 中添加对应的英文翻译。
3. 运行 `./scripts/i18n-automation.sh validate` 验证（如果存在自动化脚本）。
4. 在模板中使用 `{{getMessage "key" .Context.Lang}}`。

### 验证多语言功能
```bash
# 验证语言文件
./scripts/i18n-automation.sh validate

# 检查模板语法
./scripts/i18n-automation.sh check-syntax

# 运行测试
./scripts/i18n-automation.sh test
```

### 扩展新语言
1. 创建新的语言文件（如 `ja.json`）。
2. 复制现有结构并翻译。
3. 更新配置文件中的支持语言列表。
4. 测试新语言的显示效果。