# 模板函数使用指南

## 概述

iDiary提供了丰富的模板函数，用于在HTML模板中访问用户信息、格式化数据、国际化消息等。这些函数通过ViewService统一管理，确保模板层的数据访问安全和一致性。

## 模板函数分类

### 1. 用户认证相关函数

#### `isLogin .`

检查当前用户是否已登录。

```html
{{if isLogin .}}
<div class="user-panel">欢迎回来！</div>
{{else}}
<a href="/login">请登录</a>
{{end}}
```

#### `getUserID .`

获取当前登录用户的ID。

```html
{{if isLogin .}}
<a href="/u/{{getUserID .}}">我的主页</a>
{{end}}
```

#### `getUserName .`

获取当前登录用户的用户名。

```html
{{if isLogin .}}
<span>欢迎，{{getUserName .}}</span>
{{end}}
```

#### `getLoginName .`

获取当前登录用户的登录名。

```html
{{if isLogin .}}
<input type="hidden" name="login_name" value="{{getLoginName .}}" />
{{end}}
```

#### `getFullName .`

获取当前登录用户的全名。

```html
{{if isLogin .}}
<h1>{{getFullName .}}的日记空间</h1>
{{end}}
```

#### `getEmail .`

获取当前登录用户的邮箱。

```html
{{if isLogin .}}
<span class="email">{{getEmail .}}</span>
{{end}}
```

#### `getUserGroup .`

获取当前登录用户的用户组。

```html
{{if isLogin .}} {{if eq (getUserGroup .) 1}}
<div class="admin-badge">管理员</div>
{{end}} {{end}}
```

#### `getFaceUrl .`

获取当前登录用户的头像URL。

```html
{{if isLogin .}}
<img src="{{getFaceUrl .}}" alt="用户头像" />
{{end}}
```

#### `getUserAvatar .`

获取用户头像（带默认值处理）。

```html
{{if isLogin .}}
<img src="{{getUserAvatar .}}" alt="用户头像" class="avatar" />
{{end}}
```

#### `getLoginUTCTime .`

获取当前登录用户的登录时间。

```html
{{if isLogin .}}
<span>上次登录：{{formatTime (getLoginUTCTime .) "2006-01-02 15:04"}}</span>
{{end}}
```

#### `getTimeZone .`

获取当前登录用户的时区设置。

```html
{{if isLogin .}}
<span>时区：{{getTimeZone .}}</span>
{{end}}
```

#### `getLang .`

获取当前登录用户的语言设置。

```html
{{if isLogin .}}
<span>语言：{{getLang .}}</span>
{{end}}
```

#### `getStatus .`

获取当前登录用户的状态。

```html
{{if isLogin .}} {{if eq (getStatus .) 1}}
<span class="status-active">正常</span>
{{else}}
<span class="status-disabled">已禁用</span>
{{end}} {{end}}
```

### 2. 时间格式化函数

#### `formatTime`

格式化时间戳，支持自定义格式和时区。

```html
<!-- 基础用法 -->
{{formatTime .Data.Post.CreateTime "2006-01-02 15:04"}}

<!-- 使用用户时区 -->
{{formatTime .Data.Post.CreateTime "2006-01-02 15:04" (getTimeZone .)}}

<!-- 常用格式示例 -->
{{formatTime .Data.Post.CreateTime "2006-01-02"}}
<!-- 2024-01-15 -->
{{formatTime .Data.Post.CreateTime "2006-01-02 15:04"}}
<!-- 2024-01-15 14:30 -->
{{formatTime .Data.Post.CreateTime "2006-01-02 15:04:05"}}
<!-- 2024-01-15 14:30:45 -->
{{formatTime .Data.Post.CreateTime "01-02 15:04"}}
<!-- 01-15 14:30 -->
```

### 3. 国际化函数

#### `getMessage`

获取国际化消息。

```html
<!-- 基础消息 -->
<h1>{{getMessage "diary.title"}}</h1>

<!-- 带参数的消息 -->
<p>{{getMessage "welcome_user" (getUserName .)}}</p>

<!-- 错误消息 -->
<div class="error">{{getMessage "validation_error"}}</div>
```

### 4. 导航和面包屑函数

#### `renderBreadcrumb`

渲染面包屑导航。

```html
{{if .Breadcrumb}}
<nav class="breadcrumb">{{renderBreadcrumb .Breadcrumb}}</nav>
{{end}}
```

### 5. 数据处理函数

#### `safeHTML`

安全输出HTML内容。

```html
<div class="content">{{safeHTML .Data.Post.Content}}</div>
```

#### `truncateText`

截断文本到指定长度。

```html
<p class="summary">{{truncateText .Data.Post.Body 150 "..."}}</p>
```

#### `formatFileSize`

格式化文件大小。

```html
<span class="file-size"> 文件大小：{{formatFileSize .Data.File.Size}} </span>
```

## 实际使用示例

### 1. 用户导航栏

```html
<!-- website/views/navbar.html -->
<nav class="navbar">
  <div class="navbar-brand">
    <a href="/">iDiary</a>
  </div>

  <div class="navbar-nav">
    {{if isLogin .}}
    <!-- 已登录用户 -->
    <div class="user-menu">
      <img src="{{getUserAvatar .}}" alt="头像" class="avatar-sm" />
      <span>{{getUserName .}}</span>
      <div class="dropdown">
        <a href="/u/{{getUserID .}}">我的空间</a>
        <a href="/settings">设置</a>
        <a href="/logout">退出</a>
      </div>
    </div>
    {{else}}
    <!-- 未登录用户 -->
    <a href="/login" class="btn">登录</a>
    <a href="/register" class="btn btn-primary">注册</a>
    {{end}}
  </div>
</nav>
```

### 2. 日记列表页面

```html
<!-- website/views/diary/index.html -->
<div class="diary-container">
  <!-- 页面标题 -->
  <h1>{{.Title}}</h1>

  <!-- 面包屑导航 -->
  {{if .Breadcrumb}} {{renderBreadcrumb .Breadcrumb}} {{end}}

  <!-- 日记本信息 -->
  {{with .Data.Diary}}
  <div class="diary-header">
    <h2>{{.Title}}</h2>
    <p class="diary-meta">
      创建者：{{.AddUserName}} | 创建时间：{{formatTime .CreateTime
      "2006-01-02"}} | {{.InfoNum}}篇日记
    </p>

    {{if $.Data.IsOwner}}
    <div class="owner-actions">
      <button class="btn" onclick="editDiary()">编辑</button>
      <button class="btn" onclick="addPost()">写日记</button>
    </div>
    {{end}}
  </div>
  {{end}}

  <!-- 日记列表 -->
  {{range .Data.Posts}}
  <article class="post-item">
    <h3><a href="/post/{{.PostID}}">{{.Title}}</a></h3>
    <div class="post-meta">
      <time>{{formatTime .CreateTime "2006-01-02 15:04" (getTimeZone $)}}</time>
      {{if eq .InfoType 1}}
      <span class="post-type">图片</span>
      {{else if eq .InfoType 2}}
      <span class="post-type">视频</span>
      {{else}}
      <span class="post-type">文字</span>
      {{end}}
    </div>
    <div class="post-excerpt">{{truncateText .Body 200 "..."}}</div>
  </article>
  {{end}}
</div>

<script>
  {{if isLogin .}}
    // JavaScript中使用模板函数
    var currentUserId = {{getUserID .}};
    var currentUserName = '{{getUserName .}}';
    var userTimeZone = '{{getTimeZone .}}';

    console.log('当前用户:', currentUserName, '时区:', userTimeZone);
  {{end}}
</script>
```

### 3. 用户设置页面

```html
<!-- website/views/member/settings.html -->
<div class="settings-container">
  <h1>{{getMessage "user.settings"}}</h1>

  {{if isLogin .}}
  <form class="settings-form" method="post">
    <div class="form-group">
      <label>{{getMessage "user.username"}}</label>
      <input type="text" name="username" value="{{getUserName .}}" readonly />
    </div>

    <div class="form-group">
      <label>{{getMessage "user.email"}}</label>
      <input type="email" name="email" value="{{getEmail .}}" />
    </div>

    <div class="form-group">
      <label>{{getMessage "user.full_name"}}</label>
      <input type="text" name="full_name" value="{{getFullName .}}" />
    </div>

    <div class="form-group">
      <label>{{getMessage "user.timezone"}}</label>
      <select name="timezone">
        {{range .Data.TimeZones}}
        <option value="{{.}}" {{if eq . (getTimeZone $)}}selected{{end}}>
          {{.}}
        </option>
        {{end}}
      </select>
    </div>

    <div class="form-group">
      <label>{{getMessage "user.language"}}</label>
      <select name="language">
        {{range .Data.Languages}}
        <option value="{{.Code}}" {{if eq .Code (getLang $)}}selected{{end}}>
          {{.Name}}
        </option>
        {{end}}
      </select>
    </div>

    <button type="submit" class="btn btn-primary">{{getMessage "save"}}</button>
  </form>
  {{else}}
  <p>{{getMessage "login_required"}}</p>
  {{end}}
</div>
```

### 4. 错误页面

```html
<!-- website/views/error.html -->
<div class="error-container">
  <h1>{{getMessage "error"}}</h1>

  <div class="error-details">
    <p>{{.Data.ErrMsg}}</p>
  </div>

  {{if isLogin .}}
  <div class="error-actions">
    <a href="/u/{{getUserID .}}" class="btn">返回我的主页</a>
    <a href="/" class="btn btn-primary">返回首页</a>
  </div>
  {{else}}
  <div class="error-actions">
    <a href="/" class="btn btn-primary">返回首页</a>
    <a href="/login" class="btn">登录</a>
  </div>
  {{end}}
</div>
```

## 高级用法

### 1. 条件判断组合

```html
<!-- 管理员功能 -->
{{if and (isLogin .) (eq (getUserGroup .) 1)}}
<div class="admin-panel">
  <h3>管理员功能</h3>
  <a href="/admin/users">用户管理</a>
  <a href="/admin/settings">系统设置</a>
</div>
{{end}}

<!-- 用户状态检查 -->
{{if isLogin .}} {{if eq (getStatus .) 1}}
<div class="user-active">
  <p>账户状态：正常</p>
</div>
{{else}}
<div class="user-disabled">
  <p>您的账户已被禁用，请联系管理员</p>
</div>
{{end}} {{end}}
```

### 2. 循环中使用模板函数

```html
<!-- 在循环中访问当前用户信息 -->
{{range .Data.Comments}}
<div class="comment">
  <div class="comment-meta">
    <span>{{.UserName}}</span>
    <time>{{formatTime .CreateTime "2006-01-02 15:04" (getTimeZone $)}}</time>

    {{if and (isLogin $) (eq .UserId (getUserID $))}}
    <a href="#" class="delete-comment">删除</a>
    {{end}}
  </div>
  <div class="comment-content">{{.Content}}</div>
</div>
{{end}}
```

### 3. JavaScript集成

```html
<script>
  // 用户信息注入
  window.USER_INFO = {
    {{if isLogin .}}
      isLogin: true,
      userId: {{getUserID .}},
      userName: '{{getUserName .}}',
      userGroup: {{getUserGroup .}},
      timeZone: '{{getTimeZone .}}',
      language: '{{getLang .}}'
    {{else}}
      isLogin: false
    {{end}}
  };

  // 国际化消息注入
  window.MESSAGES = {
    'save': '{{getMessage "save"}}',
    'cancel': '{{getMessage "cancel"}}',
    'confirm': '{{getMessage "confirm"}}',
    'error': '{{getMessage "error"}}'
  };

  // 时间格式化函数
  function formatTime(timestamp) {
    {{if isLogin .}}
      return moment.unix(timestamp).tz('{{getTimeZone .}}').format('YYYY-MM-DD HH:mm');
    {{else}}
      return moment.unix(timestamp).format('YYYY-MM-DD HH:mm');
    {{end}}
  }
</script>
```

## 最佳实践

### 1. 安全性

- 始终使用模板函数而不是直接访问原始数据
- 对用户输入使用`safeHTML`函数进行安全处理
- 在JavaScript中注入数据时注意XSS防护

### 2. 性能优化

- 在模板顶部缓存常用的模板函数结果
- 避免在循环中重复调用相同的模板函数
- 使用条件判断减少不必要的函数调用

```html
<!-- 好的做法：缓存结果 -->
{{$isLoggedIn := isLogin .}} {{$currentUserId := getUserID .}} {{if
$isLoggedIn}} {{range .Data.Items}} {{if eq .UserId $currentUserId}}
<!-- 当前用户的项目 -->
{{end}} {{end}} {{end}}
```

### 3. 错误处理

- 使用条件判断检查函数返回值
- 为模板函数提供默认值
- 在关键位置添加错误提示

```html
{{$userName := getUserName .}} {{if $userName}}
<span>欢迎，{{$userName}}</span>
{{else}}
<span>游客</span>
{{end}}
```

### 4. 国际化

- 所有显示文本都使用`getMessage`函数
- 时间格式化考虑用户时区
- 数字和货币格式化遵循用户语言设置

```html
<p>
  {{getMessage "created_at"}}: {{formatTime .CreateTime "2006-01-02"
  (getTimeZone .)}}
</p>
```

这套模板函数体系提供了完整的用户信息访问、数据格式化和国际化支持，确保了模板层的安全性和一致性。
