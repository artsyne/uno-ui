# iDiary.ai Hebe 模板翻译功能指南

## 概述

基于 **iDiary.ai Hebe** 的智能国际化系统，模板支持通过 `getMessage` 系列函数实现多语言翻译。系统采用预加载缓存设计，提供微秒级翻译查找性能。

> 📖 详细的国际化架构请参考 [国际化系统指南](./internationalization-guide.md)

## 模板翻译函数

### 1. 基础翻译方法 - getMessage

**文件位置**: `internal/app/templates/localization.go`

```html
<!-- 基础翻译 -->
<h1>{{getMessage "page.title.home" .Context.Lang}}</h1>
<p>{{getMessage "page.description.welcome" .Context.Lang}}</p>

<!-- 导航菜单 -->
<nav>
  <a href="/">{{getMessage "nav.home" .Context.Lang}}</a>
  <a href="/diary">{{getMessage "nav.diary" .Context.Lang}}</a>
  <a href="/settings">{{getMessage "nav.settings" .Context.Lang}}</a>
</nav>

<!-- 按钮和表单 -->
<button>{{getMessage "button.submit" .Context.Lang}}</button>
<button>{{getMessage "button.cancel" .Context.Lang}}</button>
<label>{{getMessage "form.username" .Context.Lang}}</label>
<label>{{getMessage "form.password" .Context.Lang}}</label>
```

### 2. 带默认值的翻译方法 - getMessageWithDefault

```html
<!-- 带默认值的翻译（如果key不存在，显示默认值） -->
<span>{{getMessageWithDefault "user.status" .Context.Lang "未知状态"}}</span>
<p>{{getMessageWithDefault "new.feature" .Context.Lang "新功能"}}</p>
```

### 3. 批量翻译方法 - getMessages

```html
<!-- 批量获取多个翻译 -->
{{$navKeys := slice "nav.home" "nav.diary" "nav.settings"}} {{$navMessages :=
getMessages $navKeys .Context.Lang}} {{range $key, $value := $navMessages}}
<a href="/{{$key}}" class="nav-link">{{$value}}</a>
{{end}}
```

## 语言文件结构

### 英文 (website/locales/en.json)

```json
{
  "page.title.home": "Home",
  "page.title.diary": "My Diary",
  "page.title.settings": "Settings",
  "nav.home": "Home",
  "nav.diary": "Diary",
  "nav.settings": "Settings",
  "button.submit": "Submit",
  "button.cancel": "Cancel",
  "form.username": "Username",
  "form.password": "Password",
  "greeting.welcome_user": "Welcome back, {{username}}!",
  "stats.item_count": "You have {{count}} {{type}}",
  "common.default_page_desc": "Your personal diary space"
}
```

### 中文 (website/locales/zh-cn.json)

```json
{
  "page.title.home": "首页",
  "page.title.diary": "我的日记",
  "page.title.settings": "设置",
  "nav.home": "首页",
  "nav.diary": "日记",
  "nav.settings": "设置",
  "button.submit": "提交",
  "button.cancel": "取消",
  "form.username": "用户名",
  "form.password": "密码",
  "greeting.welcome_user": "欢迎回来，{{username}}！",
  "stats.item_count": "您有 {{count}} 个{{type}}",
  "common.default_page_desc": "您的个人日记空间"
}
```

## 语言检测优先级

系统会按以下优先级自动检测用户语言：

1. **用户设置** - 已登录用户的语言偏好设置
2. **Cookie设置** - 浏览器 Cookie 中的 `lang` 值
3. **浏览器偏好** - `Accept-Language` 请求头
4. **系统默认** - 配置文件中的默认语言

## 语言检测和切换

### 智能语言检测

**优先级顺序**:

1. **用户设置**: 已登录用户的语言偏好
2. **Cookie**: 浏览器 Cookie 中的 `lang` 值
3. **浏览器**: `Accept-Language` 请求头
4. **默认**: 系统配置的默认语言

### 语言切换实现

```html
<!-- 语言切换器 -->
<div class="language-switcher">
  {{$currentLang := .Context.Lang}} {{$supportedLangs := getSupportedLanguages}}
  {{range $lang := $supportedLangs}}
  <a
    href="/lang/{{$lang}}"
    class="lang-link {{if eq $lang $currentLang}}active{{end}}"
  >
    {{getMessageWithDefault (printf "lang.%s" $lang) $currentLang $lang}}
  </a>
  {{end}}
</div>
```

## 完整模板示例

```html
<!DOCTYPE html>
<html lang="{{.Context.Lang}}">
  <head>
    <title>{{.Title}} - {{getMessage "site.name" .Context.Lang}}</title>
    <meta name="description" content="{{.Description}}" />
  </head>
  <body>
    <!-- 导航栏 -->
    <nav class="navbar">
      <a href="/">{{getMessage "nav.home" .Context.Lang}}</a>
      <a href="/diary">{{getMessage "nav.diary" .Context.Lang}}</a>

      {{if isLogin .}}
      <div class="user-menu">
        <span
          >{{getMessage "greeting.hello" .Context.Lang}} {{getUserName .}}</span
        >
        <a href="/settings">{{getMessage "nav.settings" .Context.Lang}}</a>
        <a href="/logout">{{getMessage "nav.logout" .Context.Lang}}</a>
      </div>
      {{else}}
      <div class="auth-menu">
        <a href="/login">{{getMessage "nav.login" .Context.Lang}}</a>
        <a href="/register">{{getMessage "nav.register" .Context.Lang}}</a>
      </div>
      {{end}}
    </nav>

    <!-- 主要内容 -->
    <main>
      <h1>{{getMessage "page.diary.title" .Context.Lang}}</h1>

      <!-- 日记列表 -->
      {{if .Data.diaries}} {{range .Data.diaries}}
      <article class="diary-item">
        <h2>{{.Title}}</h2>
        <p>{{.Description}}</p>
        <time>{{formatTime .CreatedTime "2006-01-02" (getTimeZone .)}}</time>
      </article>
      {{end}} {{else}}
      <p>{{getMessage "diary.empty" .Context.Lang}}</p>
      {{end}}
    </main>
  </body>
</html>
```

## 最佳实践

### 1. 性能优化

- **缓存语言**: 在模板开始时获取 `{{$lang := .Context.Lang}}`
- **批量查询**: 使用 `getMessages` 获取多个翻译
- **合理默认值**: 为新功能提供合适的默认翻译

### 2. 翻译键命名

- **层次化**: `page.diary.title`、`nav.home`
- **语义化**: `button.submit`、`message.success`
- **一致性**: 同类功能使用相同前缀

### 3. 错误处理

- **容错设计**: 使用 `getMessageWithDefault` 提供回退
- **键验证**: 确保翻译键在所有语言文件中存在
- **测试覆盖**: 验证所有语言的翻译正确性

## 翻译键的命名约定

建议使用以下命名约定：

- `page.title.*` - 页面标题
- `page.description.*` - 页面描述
- `nav.*` - 导航相关
- `button.*` - 按钮文本
- `form.*` - 表单标签
- `message.*` - 提示消息
- `error.*` - 错误信息
- `common.*` - 通用文本

## 系统特性

### 🚀 高性能特性

- **预加载缓存**: 系统启动时加载所有语言文件到内存
- **微秒级查找**: 高效缓存map，翻译查找<1μs
- **智能回退**: 自动降级到默认语言或键名

### 🛡️ 容错机制

- **键不存在**: 返回键本身或默认值
- **语言不存在**: 自动降级到默认语言
- **文件缺失**: 不会导致应用崩溃
- **类型安全**: 编译时检查模板函数参数

### 🔄 动态特性

- **热切换**: 支持运行时语言切换
- **智能检测**: 自动检测用户首选语言
- **Cookie持久化**: 记住用户语言选择

## 总结

iDiary.ai Hebe 的模板翻译系统提供了：

- **🎯 简单易用**: 通过 `getMessage` 系列函数实现翻译
- **⚡ 高性能**: 预加载缓存，微秒级查找速度
- **🔒 类型安全**: Go 1.25 泛型支持，编译时类型检查
- **🛡️ 容错机制**: 完善的错误处理和回退策略
- **🌍 国际化**: 支持多语言和RTL语言方向
