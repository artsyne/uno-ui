# iDiary 国际化系统指南

## 概述

iDiary采用现代化的国际化(i18n)系统，支持多语言界面和智能语言检测。系统基于预加载机制和配置驱动设计，提供高性能的多语言支持。

## 系统架构

### 核心组件

- **I18nManager**: 国际化管理器，负责语言资源加载和翻译
- **语言检测**: 智能检测用户语言偏好
- **模板函数**: 在模板中提供便捷的翻译功能
- **配置驱动**: 通过配置文件管理支持的语言

### 文件结构

```
website/locales/
├── en.json          # 英语资源文件
├── zh-cn.json       # 简体中文资源文件
└── [其他语言].json   # 其他语言资源文件

internal/utils/i18n/
└── i18n.go  # 国际化管理器实现
```

## 语言资源文件

### 文件格式

语言资源文件采用JSON格式，支持嵌套结构：

```json
{
  "common": {
    "save": "Save",
    "cancel": "Cancel",
    "delete": "Delete"
  },
  "user": {
    "login": "Login",
    "logout": "Logout",
    "profile": "Profile"
  },
  "diary": {
    "title": "My Diary",
    "create": "Create New Entry",
    "edit": "Edit Entry"
  }
}
```

### 命名约定

- 使用小写字母和下划线
- 按功能模块组织键名
- 保持键名在所有语言文件中一致

## I18nManager API

### 核心方法

```go
// 创建国际化管理器
func NewI18nManager(config *config.Config) *I18nManager

// 获取翻译文本
func (im *I18nManager) GetText(locale, key string) string

// 智能语言检测
func (im *I18nManager) GetLocaleFromContext(c *gin.Context) string

// 验证语言是否支持
func (im *I18nManager) isValidLocale(locale string) bool
```

### 语言检测机制

系统按以下优先级检测用户语言：

1. **URL参数**: `?lang=zh-cn`
2. **Cookie**: `lang=zh-cn`
3. **浏览器偏好**: `Accept-Language` 头
4. **默认语言**: 配置文件中的默认语言

```go
func (lm *I18nManager) GetLocaleFromContext(c *gin.Context) string {
    // 1. 检查URL参数
    if lang := c.Query("lang"); lang != "" && lm.isValidLocale(lang) {
        return lang
    }
    
    // 2. 检查Cookie
    if lang, err := c.Cookie("lang"); err == nil && lm.isValidLocale(lang) {
        return lang
    }
    
    // 3. 检查浏览器偏好
    if lang := lm.getBrowserPreferredLocale(c); lang != "" {
        return lang
    }
    
    // 4. 返回默认语言
    return lm.defaultLocale
}
```

## 模板中的使用

### 翻译函数

在HTML模板中使用翻译函数：

```html
<!-- 基本翻译 -->
<h1>{{t "diary.title"}}</h1>

<!-- 带参数的翻译 -->
<p>{{tf "user.welcome" .Username}}</p>

<!-- 获取当前语言 -->
<html lang="{{.Lang}}">
```

### PageData结构

```go
type PageData struct {
    Lang string                    // 当前语言
    Data map[string]interface{}   // 页面数据
    // 其他字段...
}

// 便捷方法
func (pd *PageData) GetLang() string {
    return pd.Lang
}
```

## 配置管理

### 配置文件设置

```yaml
# config.yaml
app:
  default_locale: "en"
  supported_locales:
    - "en"
    - "zh-cn"
  locale_path: "website/locales"
```

### 配置结构

```go
type Config struct {
    App struct {
        DefaultLocale     string   `yaml:"default_locale"`
        SupportedLocales  []string `yaml:"supported_locales"`
        LocalePath        string   `yaml:"locale_path"`
    } `yaml:"app"`
}
```

## 性能优化

### 预加载机制

- 所有语言文件在应用启动时加载到内存
- 避免运行时文件I/O操作
- 提供快速的翻译查找

### 缓存策略

```go
type I18nManager struct {
    messages      map[string]map[string]string // [locale][key]value
    defaultLocale string
    config        *config.Config
}
```

## 最佳实践

### 1. 键名设计

- 使用层次化的键名结构
- 避免过深的嵌套（建议不超过3层）
- 使用描述性的键名

```json
{
  "user": {
    "form": {
      "email": "Email Address",
      "password": "Password"
    },
    "message": {
      "login_success": "Login successful",
      "login_failed": "Login failed"
    }
  }
}
```

### 2. 翻译管理

- 保持所有语言文件的键结构一致
- 定期检查缺失的翻译
- 使用占位符处理动态内容

### 3. 语言切换

- 提供用户友好的语言切换界面
- 保存用户的语言偏好
- 支持URL参数临时切换语言

## 扩展新语言

### 添加新语言支持

1. **创建语言文件**:
   ```bash
   cp website/locales/en.json website/locales/fr.json
   ```

2. **翻译内容**:
   ```json
   {
     "common": {
       "save": "Enregistrer",
       "cancel": "Annuler"
     }
   }
   ```

3. **更新配置**:
   ```yaml
   app:
     supported_locales:
       - "en"
       - "zh-cn"
       - "fr"  # 新增法语
   ```

4. **重启应用**:
   新语言将在应用重启后生效。

## 故障排除

### 常见问题

1. **翻译不显示**:
   - 检查键名是否正确
   - 确认语言文件格式正确
   - 验证语言代码是否在支持列表中

2. **语言检测失败**:
   - 检查浏览器Accept-Language头
   - 确认Cookie设置正确
   - 验证默认语言配置

3. **性能问题**:
   - 确认使用预加载机制
   - 避免在模板中进行复杂的翻译逻辑
   - 检查语言文件大小

### 调试技巧

```go
// 启用调试日志
slog.Debug("Language detected", "locale", locale, "source", "browser")

// 检查翻译键是否存在
if text := im.GetText(locale, key); text == key {
    slog.Warn("Translation missing", "locale", locale, "key", key)
}
```

## 重构历史

### 主要变更

- **v1.0**: 基础MessageService实现
- **v2.0**: 重构为I18nManager
- **v2.1**: 添加智能语言检测
- **v2.2**: 优化性能和预加载机制

### 架构演进

- 删除了兼容性接口MessageService
- 采用配置驱动的设计
- 引入智能语言检测机制
- 优化了模板函数接口
- 统一了国际化模块的命名

---

**文档版本**: v2.2  
**最后更新**: 2024年  
**维护团队**: iDiary开发团队