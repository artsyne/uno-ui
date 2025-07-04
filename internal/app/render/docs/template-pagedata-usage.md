# PageData 模板访问更新指南

## 概述

基于最新的PageData结构重构，用户认证信息现在通过`RenderContext`和模板函数访问，而不是直接在`Data`字段中传递。

## 模板访问方式更新

### 旧的访问方式 (需要更新)

```html
<!-- 检查登录状态 -->
{{if .IsLogin}} {{if not .IsLogin}}

<!-- 访问用户信息 -->
{{.Data.AuthUser.UserName}} {{.Data.AuthUser.UserId}} {{.Data.AuthUser.FaceUrl}}

<!-- JavaScript中的用户信息 -->
var login_user_name = '{{.Data.AuthUser.UserName}}'; var login_user_id =
{{.Data.AuthUser.UserId}};
```

### 新的访问方式 (推荐)

```html
<!-- 检查登录状态 - 使用模板函数 -->
{{if isLogin .}} {{if not (isLogin .)}}

<!-- 访问用户信息 - 使用模板函数 -->
{{getUserName .}} {{getUserID .}} {{getFaceUrl .}} {{getUserAvatar .}}
<!-- 带默认头像处理 -->

<!-- JavaScript中的用户信息 -->
{{if isLogin .}}
<script>
  var login_user_name = '{{getUserName .}}';
  var login_user_id = {{getUserID .}};
</script>
{{end}}
```

### 完整的模板函数列表

- `isLogin .` - 检查是否登录
- `getLoginName .` - 获取登录名
- `getAccessUrl .` - 获取访问URL
- `getUserID .` - 获取用户ID
- `getUserName .` - 获取用户名
- `getFullName .` - 获取全名
- `getEmail .` - 获取邮箱
- `getUserGroup .` - 获取用户组
- `getFaceUrl .` - 获取头像URL
- `getUserAvatar .` - 获取用户头像（带默认值处理）
- `getLoginUTCTime .` - 获取登录时间
- `getTimeZone .` - 获取时区
- `getLang .` - 获取语言
- `getStatus .` - 获取状态

## 需要更新的模板文件

### 1. website/views/navbar.html

```html
<!-- 更新前 -->
{{if not .IsLogin }}
<a href="/register">免费注册</a>
<a href="/login">登录</a>
{{else}}
<img src="{{.Data.AuthUser.FaceUrl}}" />
<a href="/u/{{.Data.AuthUser.UserId }}">我的主页</a>

<!-- 更新后 -->
{{if not (isLogin .)}}
<a href="/register">免费注册</a>
<a href="/login">登录</a>
{{else}}
<img src="{{getUserAvatar .}}" />
<a href="/u/{{getUserID .}}">我的主页</a>
{{end}}
```

### 2. website/views/header.html

```html
<!-- 更新前 -->
{{if .IsLogin }}
<script>
  var login_user_name = '{{.Data.AuthUser.UserName}}';
  var login_user_id = {{.Data.AuthUser.UserId }};
</script>
{{end}}

<!-- 更新后 -->
{{if isLogin .}}
<script>
  var login_user_name = '{{getUserName .}}';
  var login_user_id = {{getUserID .}};
</script>
{{end}}
```

### 3. 其他需要检查的模板文件

- `website/views/member/index.html`
- `website/views/member/settings.html`
- `website/views/diary/index.html`
- `website/views/form/diary.html`

## Handler更新说明

Handler中无需在`data.Data`中设置用户信息，这些信息会自动通过`RenderContext`提供：

```go
// 错误的做法 (不要这样做)
data.Data["AuthUser"] = authUser
data.Data["UserData"] = loginUser
data.Data["IsLogin"] = true

// 正确的做法 (推荐)
// 用户信息会自动通过 ContextBuilder 设置到 data.Context 中
// 业务数据通过 data.Data 传递
data.Data["Post"] = post
data.Data["Diary"] = diary
data.Data["IsOwner"] = isOwner
```

## 页面用户信息传递

对于页面显示的其他用户信息（非当前登录用户），继续通过`data.Data`传递：

```go
// 页面用户信息（如用户主页显示的用户）
data.Data["PageUser"] = pageUser
```

在模板中访问：

```html
<!-- 页面用户信息 -->
{{.Data.PageUser.UserName}}

<!-- 当前登录用户信息 -->
{{getUserName .}}
```

## 注意事项

1. **向后兼容**: 更新模板后需要确保功能正常
2. **默认值处理**: 使用`getUserAvatar`函数会自动处理默认头像
3. **空值检查**: 模板函数内部已经处理了空值情况
4. **性能**: 新的模板函数比直接访问更安全，且性能相当
