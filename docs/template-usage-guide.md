# 模板使用说明

## 基于新架构的模板函数使用示例

### 1. 认证用户信息访问

新架构中，认证用户信息通过模板函数访问：

```html
<!-- 检查是否登录 -->
{{ if isLogin . }}
<div class="user-info">
  <!-- 获取用户基本信息 -->
  <p>登录名: {{ getLoginName . }}</p>
  <p>用户名: {{ getUserName . }}</p>
  <p>全名: {{ getFullName . }}</p>
  <p>邮箱: {{ getEmail . }}</p>

  <!-- 获取用户头像 -->
  <img
    src="{{ getUserAvatar . }}"
    alt="{{ getFullName . }}的头像"
    class="avatar"
  />

  <!-- 获取用户设置 -->
  <p>用户组: {{ getUserGroup . }}</p>
  <p>语言: {{ getLang . }}</p>
  <p>时区: {{ getTimeZone . }}</p>
  <p>状态: {{ getStatus . }}</p>

  <!-- 获取时间信息 -->
  <p>登录时间: {{ getLoginUTCTime . }}</p>
</div>
{{ else }}
<a href="/login">请登录</a>
{{ end }}
```

### 2. 页面所有者判断

```html
<!-- 判断是否为页面所有者 -->
{{ if .Data.IsOwner }}
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

## 架构优势

1. **无循环依赖**: render 模块只依赖接口，user 模块实现接口
2. **类型安全**: 编译时检查接口实现
3. **职责清晰**:
   - 认证用户信息: 通过模板函数访问
   - 页面用户信息: 通过 Data 传递
4. **易于维护**: 接口抽象使代码更清晰
5. **高性能**: 预计算上下文，避免重复计算
