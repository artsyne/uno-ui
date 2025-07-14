# 模板翻译功能使用指南

## 概述

现在 `PageData` 已经集成了翻译功能，您可以在模板中直接使用 `{{.T "key"}}` 和 `{{.TWithParams "key" .params}}` 来获取多语言文本。

## 可用的翻译方法

### 1. 基础翻译方法 - T

在模板中使用简单翻译：

```html
<!-- 基础使用 -->
<h1>{{.T "page.title.home"}}</h1>
<p>{{.T "page.description.welcome"}}</p>

<!-- 导航菜单 -->
<nav>
  <a href="/">{{.T "nav.home"}}</a>
  <a href="/diary">{{.T "nav.diary"}}</a>
  <a href="/settings">{{.T "nav.settings"}}</a>
</nav>

<!-- 按钮和表单 -->
<button>{{.T "button.submit"}}</button>
<button>{{.T "button.cancel"}}</button>
<label>{{.T "form.username"}}</label>
<label>{{.T "form.password"}}</label>
```

### 2. 带参数的翻译方法 - TWithParams

```html
<!-- 带参数的翻译 -->
{{$params := dict "username" .Context.AuthUser.GetUserName}}
<p>{{.TWithParams "greeting.welcome_user" $params}}</p>

<!-- 显示统计信息 -->
{{$stats := dict "count" "5" "type" "articles"}}
<span>{{.TWithParams "stats.item_count" $stats}}</span>
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

## 实际模板使用示例

```html
<!DOCTYPE html>
<html lang="{{.GetLang}}">
  <head>
    <title>{{.Title}}</title>
    <meta name="description" content="{{.Description}}" />
  </head>
  <body>
    <!-- 主导航 -->
    <nav class="navbar">
      <a href="/" class="nav-link">{{.T "nav.home"}}</a>
      <a href="/diary" class="nav-link">{{.T "nav.diary"}}</a>
      {{if .IsAuthenticated}}
      <a href="/settings" class="nav-link">{{.T "nav.settings"}}</a>
      {{$params := dict "username" .GetAuthUser.GetUserName}}
      <span class="welcome"
        >{{.TWithParams "greeting.welcome_user" $params}}</span
      >
      {{else}}
      <a href="/login" class="nav-link">{{.T "nav.login"}}</a>
      {{end}}
    </nav>

    <!-- 页面内容 -->
    <main>
      <h1>{{.T "page.title.home"}}</h1>
      <p>{{.T "page.description.welcome"}}</p>

      <!-- 表单示例 -->
      <form>
        <label>{{.T "form.username"}}</label>
        <input type="text" name="username" />

        <label>{{.T "form.password"}}</label>
        <input type="password" name="password" />

        <button type="submit">{{.T "button.submit"}}</button>
        <button type="button">{{.T "button.cancel"}}</button>
      </form>
    </main>
  </body>
</html>
```

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

## 注意事项

1. 如果翻译键不存在，方法会返回键本身作为回退
2. 参数格式使用 `{{参数名}}` 语法
3. 语言检测是自动的，无需手动处理
4. 系统已配置容错机制，缺失的语言文件不会导致应用崩溃
