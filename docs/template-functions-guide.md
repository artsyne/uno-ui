# iDiary.ai Hebe 模板函数使用指南

## 概述

iDiary.ai Hebe提供了丰富的模板函数，用于在HTML模板中访问用户信息、格式化数据、国际化消息等。这些函数通过ViewService统一管理，基于Go 1.25泛型特性的PageData结构确保模板层的数据访问安全和一致性。

## 模板数据结构

### PageData 泛型结构

iDiary.ai Hebe使用基于Go 1.25泛型的PageData结构来传递模板数据：

```go
type PageData[T any] struct {
    Title       string         `json:"title"`       // 页面标题
    Description string         `json:"description"` // 页面描述
    Keywords    string         `json:"keywords"`    // 页面关键词
    Context     *RenderContext `json:"context"`     // 渲染上下文
    Data        T              `json:"data"`        // 业务数据（泛型）
}

type RenderContext struct {
    IsLogin  bool      `json:"is_login"`              // 用户是否登录
    AuthUser *AuthUser `json:"auth_user,omitempty"`   // 认证用户信息
    Lang     string    `json:"lang"`                  // 当前语言
    Theme    string    `json:"theme,omitempty"`       // 主题
    Layout   string    `json:"layout,omitempty"`      // 布局
}
```

### 在控制器中创建PageData

```go
func (ctrl *DiaryController) ShowDiary(c *gin.Context) {
    viewService := di.Get[*view.ViewService](ctrl.container)

    // 创建页面数据
    pageData := viewService.CreatePageData(c, "diary.detail_title")

    // 添加业务数据
    diary, _ := ctrl.diaryService.GetDiary(diaryID)
    pageData.Data["diary"] = diary
    pageData.Data["posts"] = diary.Posts

    // 渲染模板
    viewService.HTML(c, "diary/detail.html", pageData)
}
```

## 模板函数分类

### 1. 用户认证相关函数

**文件位置**: `internal/app/templates/auth.go`

#### `isLogin pageData`

检查当前用户是否已登录。

```html
{{if isLogin .}}
<div class="user-panel">欢迎回来！</div>
{{else}}
<a href="/login">请登录</a>
{{end}}
```

#### `authUser pageData`

获取认证用户对象，返回`*AuthUser`类型。

```html
{{if isLogin .}} {{$user := authUser .}}
<span>用户ID: {{$user.UserId}}</span>
{{end}}
```

#### `getUserID pageData`

获取当前登录用户的ID。

```html
{{if isLogin .}}
<a href="/u/{{getUserID .}}">我的主页</a>
{{end}}
```

#### `getUserName pageData`

获取当前登录用户的用户名。

```html
{{if isLogin .}}
<span>欢迎，{{getUserName .}}</span>
{{end}}
```

#### `getLoginName pageData`

获取当前登录用户的登录名。

```html
{{if isLogin .}}
<input type="hidden" name="login_name" value="{{getLoginName .}}" />
{{end}}
```

#### `getEmail pageData`

获取当前用户的邮箱地址。

```html
{{if isLogin .}}
<span>邮箱: {{getEmail .}}</span>
{{end}}
```

#### `getFaceUrl pageData`

获取用户头像URL。

```html
{{if isLogin .}}
<img src="{{getFaceUrl .}}" alt="用户头像" />
{{end}}
```

#### `getUserAvatar pageData`

获取用户头像URL，如果没有则返回默认头像。

```html
{{if isLogin .}}
<img src="{{getUserAvatar .}}" alt="用户头像" />
{{end}}
```

#### `getFullName pageData`

获取用户全名。

```html
{{if isLogin .}}
<span>全名: {{getFullName .}}</span>
{{end}}
```

#### `getUserGroup pageData`

获取用户组级别。

```html
{{if isLogin .}} {{$group := getUserGroup .}} {{if eq $group 1}}
<span class="admin-badge">管理员</span>
{{end}} {{end}}
```

#### `getFirstCharOfName pageData`

获取用户姓名首字符（大写），用于头像占位符。

```html
{{if isLogin .}}
<div class="avatar-placeholder">{{getFirstCharOfName .}}</div>
{{end}}
```

#### `getRandomVersion pageData`

获取随机版本号，用于资源防缓存。

```html
<link rel="stylesheet" href="/static/css/app.css?v={{getRandomVersion .}}" />
```

### 2. 国际化相关函数

**文件位置**: `internal/app/templates/localization.go`

#### `getMessage key locale`

获取国际化消息，是最核心的本地化函数。

```html
<h1>{{getMessage "page.title" "zh-CN"}}</h1>
<p>{{getMessage "page.description" "en"}}</p>
```

#### `getMessageWithDefault key locale defaultMsg`

带默认值的消息获取，当找不到翻译时使用默认值。

```html
<span>{{getMessageWithDefault "user.welcome" "zh-CN" "欢迎"}}</span>
```

#### `getMessages keys locale`

批量获取多个消息。

```html
{{$messages := getMessages (slice "nav.home" "nav.about" "nav.contact")
"zh-CN"}} {{range $key, $value := $messages}}
<li><a href="/{{$key}}">{{$value}}</a></li>
{{end}}
```

#### `getSupportedLanguages`

获取系统支持的语言列表。

```html
<select name="language">
  {{range getSupportedLanguages}}
  <option value="{{.Code}}">{{.Name}}</option>
  {{end}}
</select>
```

#### `isLanguageSupported locale`

检查指定语言是否被支持。

```html
{{if isLanguageSupported "zh-CN"}}
<div class="lang-support">支持中文</div>
{{end}}
```

### 3. 时间格式化函数

**文件位置**: `internal/app/templates/time.go`

#### `formatTime timestamp format timezone...`

格式化时间戳，支持自定义格式和时区。

`formatLocalTime` 函数用于将 UTC 时间戳和时区格式化为当地习惯的日期时间字符串，无需指定格式。

**参数**:

- `timestamp`: 时间戳（毫秒）
- `format`: Go时间格式字符串
- `timezone`: 可选，时区名称

```html
<!-- 基本用法 -->
<span>{{formatTime 1672531200000 "2006-01-02 15:04:05"}}</span>

<!-- 带时区转换 -->
<span>{{formatTime 1672531200000 "2006-01-02 15:04:05" "Asia/Shanghai"}}</span>

<!-- 用户时区 -->
{{if isLogin .}}
<span>{{formatTime 1672531200000 "2006-01-02 15:04:05" (getTimeZone .)}}</span>
{{end}}

<!-- 常用格式示例 -->
<div>年月日: {{formatTime 1672531200000 "2006-01-02"}}</div>
<div>时分秒: {{formatTime 1672531200000 "15:04:05"}}</div>
<div>中文格式: {{formatTime 1672531200000 "2006年01月02日"}}</div>
```

### 4. UI相关函数

**文件位置**: `internal/app/templates/ui.go`

#### `renderBreadcrumb breadcrumbItems`

渲染面包屑导航。

**参数**: `[]BreadcrumbItem` 面包屑项目数组

```go
type BreadcrumbItem struct {
    Name string `json:"name"` // 显示名称
    URL  string `json:"url"`  // 链接地址
}
```

**模板使用**:

```html
<!-- 在控制器中准备数据 -->
{{/* 控制器代码: breadcrumb := []templates.BreadcrumbItem{ {Name: "首页", URL:
"/"}, {Name: "日记", URL: "/diary"}, {Name: "当前日记", URL: ""}, //
最后一项不需要链接 } pageData.Data["breadcrumb"] = breadcrumb */}}

<!-- 在模板中渲染 -->
{{if .Data.breadcrumb}} {{renderBreadcrumb .Data.breadcrumb}} {{end}}

<!-- 输出HTML -->
<nav class="breadcrumb">
  <a href="/">首页</a> <span class="separator">/</span>
  <a href="/diary">日记</a> <span class="separator">/</span>
  <span>当前日记</span>
</nav>
```

## 模板函数注册和管理

### 函数注册机制

**文件位置**: `internal/app/templates/registry.go`

```go
// RegisterAllFuncs 注册所有模板函数
func RegisterAllFuncs() template.FuncMap {
    funcs := make(template.FuncMap)

    // 合并所有分类的模板函数
    for name, fn := range GetLocalizationFuncs() {
        funcs[name] = fn
    }

    for name, fn := range GetAuthFuncs() {
        funcs[name] = fn
    }

    for name, fn := range GetUIFuncs() {
        funcs[name] = fn
    }

    for name, fn := range GetTimeFuncs() {
        funcs[name] = fn
    }

    return funcs
}
```

### 在ViewService中的使用

模板函数通过ViewService进行统一管理，确保类型安全和上下文正确性：

```go
// ViewService 中的模板渲染
func (s *ViewService) HTML(c *gin.Context, name string, data interface{}) {
    // 模板函数已在初始化时注册
    c.HTML(http.StatusOK, name, data)
}
```

## 实际使用示例

### 1. 完整页面模板示例

```html
<!DOCTYPE html>
<html lang="{{.Context.Lang}}">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>{{.Title}}</title>
    <meta name="description" content="{{.Description}}" />
    <link
      rel="stylesheet"
      href="/static/css/app.css?v={{getRandomVersion .}}"
    />
  </head>
  <body>
    <header class="header">
      {{if isLogin .}}
      <div class="user-info">
        <img src="{{getUserAvatar .}}" alt="{{getUserName .}}" class="avatar" />
        <span class="username">{{getUserName .}}</span>
        {{if eq (getUserGroup .) 1}}
        <span class="admin-badge"
          >{{getMessage "user.admin" .Context.Lang}}</span
        >
        {{end}}
      </div>
      {{else}}
      <div class="auth-links">
        <a href="/login">{{getMessage "auth.login" .Context.Lang}}</a>
        <a href="/register">{{getMessage "auth.register" .Context.Lang}}</a>
      </div>
      {{end}}
    </header>

    <main class="main">
      <!-- 面包屑导航 -->
      {{if .Data.breadcrumb}} {{renderBreadcrumb .Data.breadcrumb}} {{end}}

      <!-- 页面内容 -->
      <h1>{{getMessage "page.diary.title" .Context.Lang}}</h1>

      {{if .Data.diary}}
      <article class="diary">
        <h2>{{.Data.diary.Title}}</h2>
        <div class="meta">
          <span class="author">{{getUserName .}}</span>
          <span class="date"
            >{{formatTime .Data.diary.CreateTime "2006-01-02 15:04" (getTimeZone
            .)}}</span
          >
        </div>
        <div class="content">{{.Data.diary.Content}}</div>
      </article>
      {{end}}
    </main>

    <footer class="footer">
      <div class="lang-switcher">
        <select name="language">
          {{range getSupportedLanguages}}
          <option
            value="{{.Code}}"
            {{if
            eq
            .Code
            $.Context.Lang}}selected{{end}}
          >
            {{.Name}}
          </option>
          {{end}}
        </select>
      </div>
    </footer>
  </body>
</html>
```

### 2. 组件化模板示例

```html
<!-- 用户卡片组件 -->
{{define "user-card"}}
<div class="user-card">
  {{if isLogin .}}
  <div class="avatar">
    {{if getFaceUrl .}}
    <img src="{{getFaceUrl .}}" alt="{{getUserName .}}" />
    {{else}}
    <div class="avatar-placeholder">{{getFirstCharOfName .}}</div>
    {{end}}
  </div>
  <div class="info">
    <h3>{{getFullName .}}</h3>
    <p>@{{getUserName .}}</p>
    <span class="email">{{getEmail .}}</span>
  </div>
  {{else}}
  <div class="guest-info">
    <p>{{getMessage "user.guest" .Context.Lang}}</p>
    <a href="/login">{{getMessage "auth.login" .Context.Lang}}</a>
  </div>
  {{end}}
</div>
{{end}}

<!-- 使用组件 -->
{{template "user-card" .}}
```

### 3. 列表渲染示例

```html
<!-- 日记列表 -->
{{if .Data.diaries}}
<div class="diary-list">
  {{range .Data.diaries}}
  <article class="diary-item">
    <h3><a href="/diary/{{.ID}}">{{.Title}}</a></h3>
    <div class="meta">
      <span class="date"
        >{{formatTime .CreateTime "01-02" (getTimeZone $)}}</span
      >
      <span class="author">{{.AuthorName}}</span>
    </div>
    <p class="excerpt">{{.Excerpt}}</p>
  </article>
  {{end}}
</div>
{{else}}
<div class="empty-state">
  <p>{{getMessage "diary.empty" .Context.Lang}}</p>
  {{if isLogin .}}
  <a href="/diary/new" class="btn"
    >{{getMessage "diary.create" .Context.Lang}}</a
  >
  {{end}}
</div>
{{end}}
```

## 最佳实践

### 1. 类型安全使用

```html
<!-- ✅ 正确：检查登录状态 -->
{{if isLogin .}}
<span>{{getUserName .}}</span>
{{end}}

<!-- ❌ 错误：未检查登录状态 -->
<span>{{getUserName .}}</span>
```

### 2. 国际化最佳实践

```html
<!-- ✅ 正确：使用上下文语言 -->
<h1>{{getMessage "page.title" .Context.Lang}}</h1>

<!-- ✅ 正确：带默认值 -->
<span>{{getMessageWithDefault "custom.message" .Context.Lang "默认文本"}}</span>

<!-- ❌ 避免：硬编码语言 -->
<h1>{{getMessage "page.title" "zh-CN"}}</h1>
```

### 3. 时间格式化最佳实践

```html
<!-- ✅ 正确：使用用户时区 -->
{{if isLogin .}}
<time>{{formatTime .Data.timestamp "2006-01-02 15:04" (getTimeZone .)}}</time>
{{else}}
<time>{{formatTime .Data.timestamp "2006-01-02 15:04"}}</time>
{{end}}

<!-- ✅ 正确：相对时间显示 -->
<span class="time-ago">{{formatTime .Data.createTime "01-02 15:04"}}</span>
```

### 4. 性能优化

```html
<!-- ✅ 正确：缓存用户信息 -->
{{if isLogin .}} {{$user := authUser .}}
<div>用户名: {{$user.UserName}}</div>
<div>邮箱: {{$user.Email}}</div>
<div>头像: <img src="{{$user.FaceUrl}}" alt="{{$user.UserName}}" /></div>
{{end}}

<!-- ❌ 避免：重复调用函数 -->
{{if isLogin .}}
<div>用户名: {{getUserName .}}</div>
<div>邮箱: {{getEmail .}}</div>
<div>头像: <img src="{{getFaceUrl .}}" alt="{{getUserName .}}" /></div>
{{end}}
```

## 总结

iDiary.ai Hebe的模板函数系统提供了以下优势：

### 🚀 技术特性

- **类型安全**: 基于Go 1.25泛型的PageData确保编译时类型检查
- **模块化设计**: 函数按功能分类，便于维护和扩展
- **国际化集成**: 完整的多语言支持，智能语言检测
- **用户上下文**: 自动注入用户认证信息，安全访问

### 🔧 开发友好

- **简洁API**: 直观的函数命名和参数设计
- **安全访问**: 自动处理未登录状态，避免空指针错误
- **灵活配置**: 支持时区、语言等个性化设置
- **组件化**: 支持模板组件化开发，代码复用

### 📈 可扩展性

- **插件化**: 模块化的函数注册机制，易于添加新功能
- **统一管理**: 通过ViewService统一管理，确保一致性
- **性能优化**: 合理的缓存策略和函数设计

这套模板函数系统为iDiary.ai Hebe的前端开发提供了强大而灵活的工具，确保了开发效率和代码质量。
