# Uno UI - CSP (Content Security Policy) 配置指南

## 概述

Content Security Policy (CSP) 是一个额外的安全层，用于检测和缓解某些类型的攻击，包括跨站脚本攻击 (XSS) 和数据注入攻击。Uno UI 提供了内置的 CSP 支持，可以自动处理内联样式和脚本的 nonce 注入。

## 为什么需要 CSP？

在现代 Web 应用中，XSS 攻击是一个严重的安全威胁。CSP 通过限制浏览器可以加载和执行的资源来源，提供了有效的防护机制。

## Uno UI 的 CSP 支持

Uno UI 自动处理以下 CSP 相关任务：

1. **自动注入 nonce 到内联样式**：所有动态创建的 `<style>` 元素都会自动添加 `nonce` 属性
2. **自动注入 nonce 到内联脚本**：所有动态创建的 `<script>` 元素都会自动添加 `nonce` 属性
3. **保持安全性**：确保所有内联内容符合您的 CSP 策略

## 基础配置

### 1. 服务器端 CSP 头部

首先，您需要在服务器上配置 CSP 头部。以下是一些常见服务器的配置示例：

#### Nginx

```nginx
server {
    listen 80;
    server_name example.com;

    location / {
        # 为每个请求生成唯一的 nonce
        add_header Content-Security-Policy "default-src 'self'; script-src 'self' 'nonce-$csp_nonce'; style-src 'self' 'nonce-$csp_nonce'" always;
        root /var/www/html;
        index index.html;
    }
}
```

#### Apache

```apache
<VirtualHost *:80>
    ServerName example.com
    DocumentRoot /var/www/html

    <Directory /var/www/html>
        # 配置 CSP 头部
        Header always set Content-Security-Policy "default-src 'self'; script-src 'self' 'nonce-%{CSP_NONCE}e'; style-src 'self' 'nonce-%{CSP_NONCE}e'"
    </Directory>
</VirtualHost>
```

#### Node.js (Express)

```javascript
const express = require('express')
const crypto = require('crypto')
const app = express()

// 为每个请求生成唯一的 nonce
app.use((req, res, next) => {
  const nonce = crypto.randomBytes(16).toString('base64')
  res.locals.nonce = nonce

  // 设置 CSP 头部
  res.setHeader(
    'Content-Security-Policy',
    `default-src 'self'; script-src 'self' 'nonce-${nonce}'; style-src 'self' 'nonce-${nonce}'`
  )

  next()
})

app.get('/', (req, res) => {
  // 将 nonce 传递给模板
  res.render('index.html', { nonce: res.locals.nonce })
})

app.listen(3000)
```

### 2. 配置 Uno UI

在客户端，您需要将 nonce 传递给 Uno UI：

#### CDN 使用方式

```html
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <title>Uno UI with CSP</title>

    <!-- Tailwind CSS -->
    <script src="https://cdn.jsdelivr.net/npm/@tailwindcss/browser@4"></script>

    <!-- 配置 Uno UI，传入 CSP nonce -->
    <script type="module">
      import Uno from 'https://unpkg.com/uno-ui@0.9.0/dist/uno.es.js'

      Uno.setConfig({
        debug: false,
        csp: {
          nonce: '<%= nonce %>' // 从服务器传入的 nonce
        },
        ui: {
          base: 'rounded-lg shadow-md'
        }
      })
    </script>

    <!-- 加载 Alpine.js -->
    <script defer src="https://cdn.jsdelivr.net/npm/alpinejs@3.15.2/dist/cdn.min.js"></script>
  </head>
  <body>
    <button x-data @click="$store.uno.box.alert('Hello!')">点击我</button>
  </body>
</html>
```

#### NPM/模块化使用方式

```javascript
// app.js
import Uno from 'uno-ui'
import Alpine from 'alpinejs'

// 从服务器获取 nonce（通过全局变量或 API）
const nonce = window.__CSP_NONCE__ || null

// 配置 Uno UI
Uno.setConfig({
  debug: false,
  csp: {
    nonce: nonce
  },
  ui: {
    base: 'rounded-lg shadow-md'
  }
})

// 注册 Alpine 组件
Alpine.data('myComponent', () => ({
  // ...
}))

// 启动 Alpine
window.Alpine = Alpine
Alpine.start()
```

## 完整示例

### Express + EJS 示例

```javascript
// server.js
const express = require('express')
const crypto = require('crypto')
const app = express()

app.set('view engine', 'ejs')

// CSP 中间件
app.use((req, res, next) => {
  const nonce = crypto.randomBytes(16).toString('base64')
  res.locals.nonce = nonce
  res.setHeader(
    'Content-Security-Policy',
    `default-src 'self'; script-src 'self' 'nonce-${nonce}' 'unsafe-inline'; style-src 'self' 'nonce-${nonce}' 'unsafe-inline'`
  )
  next()
})

app.get('/', (req, res) => {
  res.render('index', { nonce: res.locals.nonce })
})

app.listen(3000, () => {
  console.log('Server running on http://localhost:3000')
})
```

```html
<!-- views/index.ejs -->
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <title>Uno UI CSP Example</title>

    <!-- Tailwind CSS -->
    <script src="https://cdn.jsdelivr.net/npm/@tailwindcss/browser@4"></script>

    <!-- Uno UI 配置 -->
    <script type="module">
      import Uno from '/dist/uno.es.js'

      Uno.setConfig({
        debug: false,
        csp: {
          nonce: '<%= nonce %>'
        },
        ui: {
          base: 'rounded-lg shadow-md'
        }
      })
    </script>

    <!-- Alpine.js -->
    <script defer src="https://cdn.jsdelivr.net/npm/alpinejs@3.15.2/dist/cdn.min.js"></script>
  </head>
  <body class="bg-gray-100 min-h-screen p-8">
    <div class="max-w-md mx-auto">
      <h1 class="text-2xl font-bold mb-4">Uno UI CSP 示例</h1>

      <!-- Box 组件示例 -->
      <button
        x-data
        @click="$store.uno.box.alert('这是一个受 CSP 保护的弹窗！')"
        class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700"
      >
        显示弹窗
      </button>

      <!-- Notify 组件示例 -->
      <button
        x-data
        @click="$store.uno.notify.success('操作成功！')"
        class="px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 ml-2"
      >
        显示通知
      </button>
    </div>
  </body>
</html>
```

## CSP 配置选项

Uno UI 的 CSP 配置通过 `Uno.setConfig()` 方法进行：

```javascript
Uno.setConfig({
  csp: {
    nonce: 'abc123' // CSP nonce 值
  }
})
```

| 选项    | 类型     | 默认值 | 描述                                                                |
| ------- | -------- | ------ | ------------------------------------------------------------------- |
| `nonce` | `string` | `null` | CSP nonce 值。如果提供，Uno UI 将自动将其注入到所有内联样式和脚本中 |

## 常见问题排查

### 问题 1：样式未应用

**症状**：页面样式混乱，组件没有正确的样式。

**原因**：CSP 阻止了内联样式。

**解决方案**：

- 确保 CSP 头部中的 `style-src` 包含 nonce 或 `'unsafe-inline'`
- 检查 nonce 是否正确传递给 Uno UI

```http
Content-Security-Policy: style-src 'self' 'nonce-abc123'
```

### 问题 2：脚本未执行

**症状**：JavaScript 代码没有执行，控制台显示 CSP 错误。

**原因**：CSP 阻止了内联脚本。

**解决方案**：

- 确保 CSP 头部中的 `script-src` 包含 nonce 或 `'unsafe-inline'`
- 验证 nonce 的一致性

```http
Content-Security-Policy: script-src 'self' 'nonce-abc123'
```

### 问题 3：混合内容错误

**症状**：浏览器阻止加载 HTTP 资源。

**原因**：在 HTTPS 页面上加载 HTTP 资源。

**解决方案**：确保所有资源都使用 HTTPS。

### 问题 4：Nonce 不匹配

**症状**：控制台显示 CSP nonce 错误。

**原因**：服务器端 nonce 和客户端 nonce 不一致。

**解决方案**：

- 确保服务器为每个请求生成唯一的 nonce
- 确保 nonce 正确传递到客户端
- 检查模板渲染是否正确

## 最佳实践

### 1. 始终使用唯一的 Nonce

为每个请求生成唯一的 nonce，不要重用或硬编码：

```javascript
// ❌ 错误：硬编码 nonce
const nonce = 'abc123'

// ✅ 正确：为每个请求生成唯一 nonce
const nonce = crypto.randomBytes(16).toString('base64')
```

### 2. 避免使用 'unsafe-inline'

尽可能使用 nonce 而不是 `'unsafe-inline'`，以获得更好的安全性：

```http
# ❌ 不推荐：使用 unsafe-inline
Content-Security-Policy: script-src 'self' 'unsafe-inline'

# ✅ 推荐：使用 nonce
Content-Security-Policy: script-src 'self' 'nonce-abc123'
```

### 3. 在开发环境中测试 CSP

在开发环境中启用 CSP，尽早发现和解决问题：

```javascript
if (process.env.NODE_ENV === 'development') {
  // 开发环境：使用较宽松的 CSP
  res.setHeader(
    'Content-Security-Policy',
    "default-src 'self'; script-src 'self' 'unsafe-inline' 'unsafe-eval'; style-src 'self' 'unsafe-inline'"
  )
} else {
  // 生产环境：使用严格的 CSP
  res.setHeader(
    'Content-Security-Policy',
    `default-src 'self'; script-src 'self' 'nonce-${nonce}'; style-src 'self' 'nonce-${nonce}'`
  )
}
```

### 4. 监控 CSP 违规

使用 `Content-Security-Policy-Report-Only` 头部来测试 CSP 策略，而不会实际阻止资源加载：

```javascript
// 报告模式：只报告不阻止
res.setHeader(
  'Content-Security-Policy-Report-Only',
  `default-src 'self'; script-src 'self' 'nonce-${nonce}'; style-src 'self' 'nonce-${nonce}'; report-uri /csp-violation-report-endpoint`
)
```

### 5. 使用 CSP 评估工具

使用在线工具评估您的 CSP 策略：

- [CSP Evaluator](https://csp-evaluator.withgoogle.com/)
- [Mozilla Observatory](https://observatory.mozilla.org/)

## 安全性提示

1. **不要在客户端生成 nonce**：nonce 应该在服务器端生成，以防止绕过 CSP
2. **使用 HTTPS**：CSP 在 HTTPS 下最有效
3. **定期审查 CSP 策略**：随着应用的发展，定期更新和审查您的 CSP 策略
4. **最小权限原则**：只允许必要的资源来源，不要过度放宽限制

## 相关资源

- [MDN: Content Security Policy](https://developer.mozilla.org/en-US/docs/Web/HTTP/CSP)
- [CSP Evaluator](https://csp-evaluator.withgoogle.com/)
- [OWASP: Content Security Policy Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Content_Security_Policy_Cheat_Sheet.html)
- [Google CSP Best Practices](https://web.dev/csp/)

## 总结

Uno UI 的 CSP 支持使得在严格的安全策略下使用组件库变得简单。通过正确配置 CSP，您可以：

- ✅ 防止 XSS 攻击
- ✅ 保护用户数据
- ✅ 符合安全合规要求
- ✅ 提高应用的整体安全性

记住：安全性是一个持续的过程，定期审查和更新您的 CSP 策略是保持应用安全的关键。
