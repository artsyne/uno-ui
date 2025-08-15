# iDiary.ai Hebe 模板开发最佳实践

## 概述

本文档基于 iDiary.ai Hebe 的 ViewService 架构和 Go 1.25 泛型 PageData 结构，提供模板开发的最佳实践指南。结合现代化的模板函数体系，帮助开发者构建高效、类型安全的模板系统。

> 📖 详细的模板函数说明请参考 [模板函数使用指南](./template-functions-guide.md)

## PageData 泛型结构

### 数据结构定义

**文件位置**: `internal/contracts/types.go`

```go
type PageData[T any] struct {
    Title       string         `json:"title"`       // 页面标题
    Description string         `json:"description"` // 页面描述
    Keywords    string         `json:"keywords"`    // 页面关键词
    Context     *RenderContext `json:"context"`     // 渲染上下文
    Data        T              `json:"data"`        // 业务数据（泛型）
}

type RenderContext struct {
    IsLogin  bool      `json:"is_login"`  // 登录状态
    AuthUser *AuthUser `json:"auth_user"` // 认证用户信息
    Lang     string    `json:"lang"`      // 当前语言
}
```

### 控制器中创建 PageData

```go
func (ctrl *DiaryController) ShowDiary(c *gin.Context) {
    viewService := di.Get[*view.ViewService](ctrl.container)

    // 创建泛型 PageData，自动设置标题、上下文等
    pageData := viewService.CreatePageData(c, "diary.detail.title")

    // 添加业务数据
    diary, _ := ctrl.diaryService.GetDiary(diaryID)
    pageData.Data["diary"] = diary
    pageData.Data["posts"] = diary.Posts

    // 渲染模板
    viewService.HTML(c, "diary/detail.html", pageData)
}
```

## 模板函数使用指南

### 1. 认证用户信息访问

通过模板函数安全访问当前登录用户信息：

```html
<!-- 检查登录状态 -->
{{if isLogin .}}
<div class="user-info">
  <!-- 基本用户信息 -->
  <h3>{{getUserName .}}</h3>
  <p>{{getFullName .}}</p>
  <p>{{getEmail .}}</p>

  <!-- 用户头像 -->
  <img src="{{getUserAvatar .}}" alt="{{getUserName .}}" class="avatar" />

  <!-- 用户设置 -->
  <span class="user-group">{{getUserGroup .}}</span>
  <span class="timezone">{{getTimeZone .}}</span>
  <span class="language">{{getLang .}}</span>

  <!-- 特殊信息 -->
  <span class="first-char">{{getFirstCharOfName .}}</span>
</div>
{{else}}
<a href="/login" class="login-link">请登录</a>
{{end}}
```

### 2. 页面所有者判断

```html
<!-- 判断是否为页面所有者 -->
{{ if .IsOwner }}
<button class="btn" data-action="showDiaryBox" data-form-action="diary.add">
  <i class="mdi mdi-plus"></i> 新增日记本
</button>
<button class="btn">编辑设置</button>
{{ end }}
```

### 3. 页面用户信息访问

页面用户信息（如用户主页的用户信息）通过 `.Data.PageUser` 访问：

```html
<!-- 页面用户信息 -->
{{ with .Data.PageUser }}
<div class="page-user-info">
  <h1>{{ .UserName }}的空间</h1>
  <p>全名: {{ .FullName }}</p>
  <p>邮箱: {{ .Email }}</p>
  <img src="{{ .FaceUrl }}" alt="用户头像" class="page-avatar" />
  <p>注册时间: {{ formatTime .CreatedUTCTime "2006-01-02" }}</p>
  <p>最后登录: {{ formatTime .LastLoginUTCTime "2006-01-02 15:04" }}</p>
  <p>用户组: {{ .UserGroup }}</p>
  <p>状态: {{ .Status }}</p>
  <p>简介: {{ .Remark }}</p>
</div>
{{ end }}
```

### 4. 时间格式化

使用用户的时区设置进行时间格式化：

```html
<!-- 使用当前登录用户的时区 -->
{{ $timeZone := getTimeZone . }} {{ range .Data.Diaries }}
<div class="diary-item">
  <h3>{{ .Title }}</h3>
  <p>
    最后更新: {{ formatTime .LastActionUTCTime "2006-01-02 15:04" $timeZone }}
  </p>
  <p>{{ .InfoNum }}篇日记</p>
</div>
{{ end }}
```

### 5. 导航和面包屑

```html
<!-- 用户导航 -->
<nav class="user-nav">
  {{ if isLogin . }}
  <a href="/u/{{ getUserID . }}">我的空间</a>
  <a href="/settings">设置</a>
  <a href="/logout">退出</a>
  {{ else }}
  <a href="/login">登录</a>
  <a href="/register">注册</a>
  {{ end }}
</nav>

<!-- 面包屑导航 -->
{{ renderBreadcrumb .Breadcrumb }}
```

### 6. 条件显示

```html
<!-- 根据用户权限显示不同内容 -->
{{ if isLogin . }} {{ if eq (getUserGroup .) 1 }}
<!-- 管理员功能 -->
<div class="admin-panel">
  <h3>管理员面板</h3>
  <button>用户管理</button>
  <button>系统设置</button>
</div>
{{ end }} {{ if eq (getStatus .) 1 }}
<!-- 正常用户功能 -->
<div class="user-actions">
  <button>写日记</button>
  <button>创建日记本</button>
</div>
{{ else }}
<!-- 用户被禁用 -->
<div class="user-disabled">
  <p>您的账户已被禁用，请联系管理员</p>
</div>
{{ end }} {{ end }}
```

### 2. 国际化翻译

```html
<!-- 基础翻译 -->
<h1>{{getMessage "page.title.diary_list" .Context.Lang}}</h1>
<p>{{getMessage "diary.description" .Context.Lang}}</p>

<!-- 带默认值的翻译 -->
<span>{{getMessageWithDefault "user.status" .Context.Lang "未知状态"}}</span>

<!-- 批量获取翻译 -->
{{$navMessages := getMessages (slice "nav.home" "nav.diary" "nav.settings")
.Context.Lang}} {{range $key, $value := $navMessages}}
<a href="/{{$key}}">{{$value}}</a>
{{end}}
```

### 3. 时间格式化

```html
<!-- 使用用户时区格式化时间 -->
{{$timezone := getTimeZone .}} {{range .Data.Posts}}
<div class="post-item">
  <h3>{{.Title}}</h3>
  <time>{{formatTime .CreatedTime "2006-01-02 15:04" $timezone}}</time>
</div>
{{end}}

<!-- 相对时间显示 -->
<span
  >创建于: {{formatTime .Data.CreateTime "2006-01-02" (getTimeZone .)}}</span
>
```

### 4. UI 组件渲染

```html
<!-- 面包屑导航 -->
{{if .Data.breadcrumb}} {{renderBreadcrumb .Data.breadcrumb}} {{end}}

<!-- 动态版本号 -->
<link rel="stylesheet" href="/static/css/app.css?v={{getRandomVersion .}}" />
<script src="/static/js/app.js?v={{getRandomVersion .}}"></script>
```

## 最佳实践模式

### 1. 页面数据访问模式

```html
<!-- ✅ 正确: 通过模板函数访问用户信息 -->
{{if isLogin .}} {{$user := authUser .}}
<div>用户名: {{$user.UserName}}</div>
<div>邮箱: {{$user.Email}}</div>
{{end}}

<!-- ✅ 正确: 通过 .Data 访问业务数据 -->
{{if .Data.diary}}
<h1>{{.Data.diary.Title}}</h1>
<p>{{.Data.diary.Description}}</p>
{{end}}

<!-- ❌ 错误: 不要直接访问 Context -->
<!-- {{.Context.AuthUser.UserName}} -->
```

### 2. 性能优化模式

```html
<!-- ✅ 正确: 缓存模板函数结果 -->
{{$lang := .Context.Lang}} {{$timezone := getTimeZone .}} {{range .Data.Items}}
<div>
  <h3>{{getMessage "item.title" $lang}}</h3>
  <time>{{formatTime .CreatedTime "2006-01-02" $timezone}}</time>
</div>
{{end}}

<!-- ❌ 错误: 重复调用模板函数 -->
{{range .Data.Items}}
<div>
  <h3>{{getMessage "item.title" .Context.Lang}}</h3>
  <time>{{formatTime .CreatedTime "2006-01-02" (getTimeZone .)}}</time>
</div>
{{end}}
```

### 3. 条件渲染模式

```html
<!-- 权限检查 -->
{{if isLogin .}} {{if eq (getUserGroup .) 1}}
<!-- 管理员功能 -->
<div class="admin-panel">
  <button>用户管理</button>
  <button>系统设置</button>
</div>
{{end}}

<!-- 普通用户功能 -->
<div class="user-actions">
  <button>写日记</button>
  <button>创建日记本</button>
</div>
{{else}}
<!-- 游客功能 -->
<div class="guest-actions">
  <a href="/login">登录</a>
  <a href="/register">注册</a>
</div>
{{end}}
```

## 架构优势

### 🏗️ ViewService 架构优势

1. **类型安全**: Go 1.25 泛型 PageData，编译时类型检查
2. **职责分离**:
   - **认证信息**: 通过模板函数访问，自动处理登录状态
   - **业务数据**: 通过 .Data 传递，类型安全
   - **国际化**: 自动语言检测和翻译
3. **高性能**: 预计算上下文，避免重复计算
4. **易维护**: 统一的 ViewService 接口，标准化的数据流

### 🔒 安全特性

- **XSS防护**: 模板函数自动转义输出
- **权限控制**: 基于用户组和状态的条件渲染
- **数据隔离**: 认证信息与业务数据分离访问

### ⚡ 性能特性

- **缓存优化**: 智能缓存模板函数结果
- **按需加载**: 只在需要时获取用户信息
- **批量操作**: 支持批量获取翻译等操作
