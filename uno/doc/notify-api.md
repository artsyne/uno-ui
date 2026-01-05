# Notify 组件 API 文档

一个基于 Alpine.js 构建的轻量级通知组件，支持多位置、多类型、自动关闭和暗色模式。

## 目录

- [功能特性](#功能特性)
- [快速开始](#快速开始)
- [核心 API](#核心-api)
- [Progress Bar（进度条）](#progress-bar进度条)
- [HTMX 全局集成](#htmx-全局集成)
- [配置选项](#配置选项)
- [样式定制](#样式定制)
- [国际化支持](#国际化支持)
- [使用示例](#使用示例)

## 功能特性

- **多位置显示**：支持 9 个位置 - `top-left`, `top-center`, `top-right`, `center-left`, `center`, `center-right`, `bottom-left`, `bottom-center`, `bottom-right`
- **消息类型**：`info`, `success`, `warning`, `error`
- **暗色模式**：完美适配 Tailwind CSS Dark Mode（`dark:` 变体）
- **自动关闭**：可配置显示时长（`duration`），设为 `0` 表示手动关闭
- **全局进度条**：支持全局 progress bar，显示在页面顶部
- **HTMX 集成**：自动绑定 HTMX 全局事件，显示请求进度
- **样式可定制**：使用标准 Tailwind CSS 类，易于覆盖
- **国际化**：支持 `Alpine.store('i18n')` 进行文本翻译
- **无障碍访问**：使用 `role="alert"` 支持屏幕阅读器

## 快速开始

### 基础使用

使用全局 store `Alpine.store('uno').notify.show()` 触发通知。

```html
<!-- 触发按钮 -->
<button @click="$store.uno.notify.show('操作成功', { type: 'success' })">成功</button>
<button @click="$store.uno.notify.show('发生错误', { type: 'error' })">错误</button>
<button @click="$store.uno.notify.show('警告信息', { type: 'warning' })">警告</button>
<button @click="$store.uno.notify.show('提示信息', { type: 'info' })">提示</button>
```

容器会在需要时自动注入到 DOM 中。

### 使用指令

```html
<!-- 使用 x-uno-notify 指令 -->
<button x-uno-notify="'保存成功!'" data-uno-type="success">保存</button>

<!-- 使用数组语法传递配置 -->
<button x-uno-notify="['操作完成', { type: 'success', duration: 5000 }]">提交</button>
```

## 核心 API

组件通过 `Alpine.store('uno').notify` 暴露方法。

### 全局访问

```javascript
// 显示成功消息
Alpine.store('uno').notify.success('保存成功!')

// 显示错误消息
Alpine.store('uno').notify.error('操作失败!')

// 显示警告消息
Alpine.store('uno').notify.warning('请注意!')

// 显示提示消息
Alpine.store('uno').notify.info('温馨提示')
```

### `show(message, options)`

显示一个通知。

| 参数      | 类型     | 默认值 | 说明     |
| --------- | -------- | ------ | -------- |
| `message` | `String` | -      | 通知内容 |
| `options` | `Object` | `{}`   | 配置对象 |

**配置选项：**

| 选项       | 类型      | 默认值        | 说明                                                                                                                                           |
| ---------- | --------- | ------------- | ---------------------------------------------------------------------------------------------------------------------------------------------- |
| `type`     | `String`  | `'info'`      | 消息类型：`info`, `success`, `warning`, `error`                                                                                                |
| `duration` | `Number`  | `3000`        | 自动关闭时长（毫秒）。设为 `0` 禁用自动关闭                                                                                                    |
| `position` | `String`  | `'top-right'` | 容器位置。选项：`top-left`, `top-center`, `top-right`, `center-left`, `center`, `center-right`, `bottom-left`, `bottom-center`, `bottom-right` |
| `closable` | `Boolean` | `true`        | 是否显示关闭按钮                                                                                                                               |
| `icon`     | `String`  | 根据类型自动  | 图标类名（可选）                                                                                                                               |

**使用示例：**

```javascript
// 基础使用
Alpine.store('uno').notify.show('操作成功', { type: 'success' })

// 自定义持续时间和位置
Alpine.store('uno').notify.show('请注意', {
  type: 'warning',
  duration: 5000,
  position: 'top-center'
})

// 手动关闭（不自动消失）
Alpine.store('uno').notify.show('重要通知', {
  type: 'info',
  duration: 0, // 不自动关闭
  closable: true
})
```

### 快捷方法

为常用类型提供的便捷方法：

- `success(message, options)` - 显示成功消息
- `error(message, options)` - 显示错误消息
- `warning(message, options)` - 显示警告消息
- `info(message, options)` - 显示提示消息

```javascript
// 等价于 show(message, { type: 'success', ...options })
Alpine.store('uno').notify.success('保存成功')

Alpine.store('uno').notify.error('网络错误')

Alpine.store('uno').notify.warning('磁盘空间不足', { duration: 5000 })

Alpine.store('uno').notify.info('系统将在5分钟后维护', { position: 'top-center' })
```

### `remove(id)`

手动移除指定 ID 的通知。

```javascript
// show() 方法返回通知 ID
const notifyId = Alpine.store('uno').notify.show('处理中...')

// 稍后手动移除
setTimeout(() => {
  Alpine.store('uno').notify.remove(notifyId)
}, 3000)
```

## Progress Bar（进度条）

组件提供全局 progress bar 功能，显示在页面顶部，用于显示加载状态。

### `startProgress()`

启动全局进度条。进度条会从 0% 快速循环加载到 100%，然后立即重置到 0% 开始下一轮循环，直到调用 `stopProgress()`。

```javascript
// 启动进度条
Alpine.store('uno').notify.startProgress()
```

**特性：**

- 进度条显示在浏览器窗口顶部
- 从 0% 开始，每 100ms 增加 10%
- 达到 100% 后立即重置到 0%，形成连续循环动画
- 使用渐变色（蓝色到靳蓝色），带有微妙的阴影效果
- 支持暗色模式

### `stopProgress()`

停止全局进度条。进度条会先完成到 100%，然后渐隐。

```javascript
// 停止进度条
Alpine.store('uno').notify.stopProgress()
```

**工作流程：**

1. 立即设置进度为 100%
2. 等待 300ms（显示完成效果）
3. 隐藏进度条
4. 重置进度为 0%

### 使用示例

```javascript
// 基础使用
Alpine.store('uno').notify.startProgress()

// 模拟异步操作
setTimeout(() => {
  Alpine.store('uno').notify.stopProgress()
}, 3000)

// 在 Alpine 组件中使用
<button
  @click="
    $store.uno.notify.startProgress();
    // 执行一些操作...
    setTimeout(() => $store.uno.notify.stopProgress(), 2000)
  ">
  开始处理
</button>
```

### 样式定制

进度条样式可通过 Uno 配置系统进行定制：

```javascript
window.Uno.setConfig({
  ui: {
    notify: {
      // 进度条容器样式
      progressWrap:
        'fixed top-0 left-0 right-0 h-1 bg-gray-200 dark:bg-gray-700 pointer-events-none',

      // 进度条样式
      progressBar:
        'h-full bg-gradient-to-r from-blue-500 to-indigo-600 shadow-sm shadow-blue-500/50'
    }
  }
})
```

## HTMX 全局集成

Notify 组件提供了与 HTMX 的全局集成，自动在所有 HTMX 请求时显示进度条。

### 启用 HTMX 集成

在页面初始化时调用 `bindNotifyToHtmx()` 方法：

```html
<script>
  document.addEventListener('alpine:init', () => {
    if (window.Uno) {
      // 绑定 notify 功能到所有 HTMX 请求
      window.Uno.htmx.bindNotifyToHtmx()
    }
  })
</script>
```

### 工作原理

一旦启用，组件会自动监听以下 HTMX 事件：

| HTMX 事件            | 触发时机        | Notify 操作                      |
| -------------------- | --------------- | -------------------------------- |
| `htmx:beforeRequest` | HTMX 请求开始前 | 启动 progress bar                |
| `htmx:afterRequest`  | HTMX 请求完成后 | 停止 progress bar                |
| `htmx:responseError` | HTMX 响应错误   | 停止 progress bar + 显示错误通知 |
| `htmx:sendError`     | HTMX 发送错误   | 停止 progress bar + 显示错误通知 |
| `htmx:timeout`       | HTMX 请求超时   | 停止 progress bar + 显示超时警告 |

### 使用示例

```html
<!DOCTYPE html>
<html>
  <head>
    <script src="https://cdn.jsdelivr.net/npm/htmx.org@2.0.4/dist/htmx.min.js"></script>
    <script src="uno.umd.js"></script>
    <script src="https://cdn.jsdelivr.net/npm/alpinejs@3.15.2/dist/cdn.min.js" defer></script>
    <script>
      // 启用 HTMX 全局进度条
      document.addEventListener('alpine:init', () => {
        if (window.Uno) {
          window.Uno.htmx.bindNotifyToHtmx()
        }
      })
    </script>
  </head>
  <body>
    <!-- 任何 HTMX 请求都会自动显示进度条 -->
    <button hx-get="/api/data" hx-target="#result">加载数据</button>

    <div id="result"></div>
  </body>
</html>
```

### 特性

- **自动化**：一次设置，全局生效，无需为每个 HTMX 请求单独配置
- **智能错误处理**：自动捕获并显示 HTMX 请求错误
- **视觉反馈**：用户在每次 HTMX 请求时都能看到清晰的加载进度
- **零侵入**：不影响现有的 HTMX 代码，只是添加视觉反馈

## 配置选项

### 四级配置继承

Notify 组件支持四级配置继承：

1. **插件默认配置**（Plugin Defaults）
2. **全局插件配置**（Global Plugin Config）
3. **直接全局配置**（Direct Global Config）
4. **实例配置**（Instance Config） - 最高优先级

```javascript
// 方式 1：全局插件配置
window.Uno.setConfig({
  plugins: {
    notify: {
      duration: 5000,
      position: 'top-center'
    }
  }
})

// 方式 2：直接全局配置
window.Uno.setConfig({
  notify: {
    duration: 5000,
    position: 'top-center'
  }
})

// 方式 3：实例配置（优先级最高）
Alpine.store('uno').notify.show('消息', {
  duration: 3000, // 覆盖全局配置
  position: 'bottom-right'
})
```

### 默认配置

```javascript
{
  type: 'info',
  duration: 3000,
  position: 'top-right',
  closable: true,
  icon: true
}
```

## 样式定制

组件样式基于 Tailwind CSS，通过 `Uno.getStyle()` 机制获取样式配置。

### 可定制的样式类

| 样式类         | 说明         | 默认值                                                                                     |
| -------------- | ------------ | ------------------------------------------------------------------------------------------ |
| `container`    | 通知容器     | `fixed flex gap-2 p-4 pointer-events-none`                                                 |
| `item`         | 通知项       | `pointer-events-auto max-w-[calc(100vw-2rem)] overflow-hidden rounded-lg shadow-md border` |
| `content`      | 内容区域     | `p-4`                                                                                      |
| `row`          | 横向布局     | `flex items-start`                                                                         |
| `iconWrap`     | 图标容器     | `flex-shrink-0`                                                                            |
| `textWrap`     | 文本容器     | `ml-3 min-w-0 flex-1 pt-0.5`                                                               |
| `text`         | 文本内容     | `text-sm font-medium text-gray-900 dark:text-slate-100 break-words`                        |
| `closeWrap`    | 关闭按钮容器 | `ml-4 flex flex-shrink-0`                                                                  |
| `closeBtn`     | 关闭按钮     | `inline-flex rounded-md bg-white dark:bg-slate-900 text-gray-400 hover:text-gray-500`      |
| `closeIcon`    | 关闭图标     | `h-5 w-5`                                                                                  |
| `progressWrap` | 进度条容器   | `fixed top-0 left-0 right-0 h-1 bg-gray-200 dark:bg-gray-700 pointer-events-none`          |
| `progressBar`  | 进度条       | `h-full bg-gradient-to-r from-blue-500 to-indigo-600 shadow-sm shadow-blue-500/50`         |

### 类型样式

| 类型      | 样式                               |
| --------- | ---------------------------------- |
| `success` | `bg-green-100 dark:bg-green-950`   |
| `error`   | `bg-red-100 dark:bg-red-950`       |
| `warning` | `bg-yellow-100 dark:bg-yellow-950` |
| `info`    | `bg-blue-100 dark:bg-blue-950`     |

### 图标样式

| 类型      | 图标样式                               |
| --------- | -------------------------------------- |
| `success` | `text-green-500 dark:text-green-300`   |
| `error`   | `text-red-500 dark:text-red-300`       |
| `warning` | `text-yellow-500 dark:text-yellow-300` |
| `info`    | `text-blue-500 dark:text-blue-300`     |

### 自定义示例

```javascript
// 自定义通知项样式
window.Uno.setConfig({
  ui: {
    notify: {
      item: 'pointer-events-auto max-w-md rounded-xl shadow-xl border-2',
      types: {
        success:
          'bg-emerald-100 text-emerald-900 border-emerald-300 dark:bg-emerald-950 dark:text-emerald-100',
        error: 'bg-rose-100 text-rose-900 border-rose-300 dark:bg-rose-950 dark:text-rose-100'
      },
      icons: {
        success: 'text-emerald-600 dark:text-emerald-400',
        error: 'text-rose-600 dark:text-rose-400'
      }
    }
  }
})
```

### 进度条自定义

```javascript
window.Uno.setConfig({
  ui: {
    notify: {
      // 自定义进度条高度和颜色
      progressWrap:
        'fixed top-0 left-0 right-0 h-2 bg-gray-300 dark:bg-gray-800 pointer-events-none',
      progressBar: 'h-full bg-gradient-to-r from-purple-500 to-pink-600'
    }
  }
})
```

## 国际化支持

组件支持 `Alpine.store('i18n')` 进行文本翻译。

```javascript
// 如果存在 i18n store，消息会自动翻译
Alpine.store('uno').notify.show('notify.success', { type: 'success' })

// i18n 配置示例
Alpine.store('i18n', {
  messages: {
    'zh-CN': {
      'notify.success': '操作成功',
      'notify.error': '操作失败',
      'notify.loading': '加载中...',
      'notify.requestFailed': '请求失败',
      'notify.timeout': '请求超时'
    },
    en: {
      'notify.success': 'Operation successful',
      'notify.error': 'Operation failed',
      'notify.loading': 'Loading...',
      'notify.requestFailed': 'Request failed',
      'notify.timeout': 'Request timeout'
    }
  }
})
```

## 使用示例

### 基础通知

```html
<!-- 基础使用 -->
<button @click="$store.uno.notify.success('保存成功')">保存</button>

<!-- 自定义持续时间 -->
<button @click="$store.uno.notify.warning('警告', { duration: 5000 })">警告</button>

<!-- 不自动关闭 -->
<button @click="$store.uno.notify.info('重要信息', { duration: 0 })">提示</button>
```

### 多位置显示

```html
<div class="grid grid-cols-3 gap-4">
  <button @click="$store.uno.notify.show('左上', { position: 'top-left' })">左上</button>
  <button @click="$store.uno.notify.show('中上', { position: 'top-center' })">中上</button>
  <button @click="$store.uno.notify.show('右上', { position: 'top-right' })">右上</button>

  <button @click="$store.uno.notify.show('左中', { position: 'center-left' })">左中</button>
  <button @click="$store.uno.notify.show('居中', { position: 'center' })">居中</button>
  <button @click="$store.uno.notify.show('右中', { position: 'center-right' })">右中</button>

  <button @click="$store.uno.notify.show('左下', { position: 'bottom-left' })">左下</button>
  <button @click="$store.uno.notify.show('中下', { position: 'bottom-center' })">中下</button>
  <button @click="$store.uno.notify.show('右下', { position: 'bottom-right' })">右下</button>
</div>
```

### 使用指令

```html
<!-- 简单使用 -->
<button x-uno-notify="'操作成功'" data-uno-type="success">提交</button>

<!-- 数组语法 -->
<button x-uno-notify="['删除成功', { type: 'success', duration: 2000 }]">删除</button>

<!-- 结合 Alpine 数据 -->
<div x-data="{ message: '保存成功' }">
  <button x-uno-notify="[message, { type: 'success' }]">保存</button>
</div>
```

### 异步操作与进度条

```html
<div
  x-data="{
  async saveData() {
    this.$store.uno.notify.startProgress()
    try {
      await fetch('/api/save', { method: 'POST' })
      this.$store.uno.notify.success('保存成功')
    } catch (error) {
      this.$store.uno.notify.error('保存失败')
    } finally {
      this.$store.uno.notify.stopProgress()
    }
  }
}"
>
  <button @click="saveData()">保存数据</button>
</div>
```

### HTMX 集成示例

```html
<!DOCTYPE html>
<html>
  <head>
    <script src="https://cdn.jsdelivr.net/npm/htmx.org@2.0.4/dist/htmx.min.js"></script>
    <script src="uno.umd.js"></script>
    <script src="https://cdn.jsdelivr.net/npm/alpinejs@3.15.2/dist/cdn.min.js" defer></script>
    <script>
      document.addEventListener('alpine:init', () => {
        // 启用 HTMX 全局进度条
        window.Uno.htmx.bindNotifyToHtmx()
      })
    </script>
  </head>
  <body>
    <!-- HTMX 请求会自动显示进度条 -->
    <button hx-get="/api/users" hx-target="#users">加载用户</button>
    <button hx-post="/api/users" hx-target="#result">创建用户</button>
    <button hx-delete="/api/users/1" hx-target="#result">删除用户</button>

    <div id="users"></div>
    <div id="result"></div>
  </body>
</html>
```

## 最佳实践

1. **优先使用快捷方法**：`success()`, `error()`, `warning()`, `info()` 比 `show()` 更简洁
2. **合理设置持续时间**：
   - 成功消息：2-3 秒
   - 错误消息：4-5 秒或手动关闭
   - 警告消息：3-4 秒
3. **使用 HTMX 全局集成**：为所有 HTMX 请求提供统一的视觉反馈
4. **避免通知滥用**：不要在短时间内显示过多通知
5. **使用国际化**：如果项目支持多语言，使用 i18n 键值而不是硬编码文本
6. **利用配置继承**：在全局配置中设置通用选项，在实例中只覆盖特殊选项

---

更多信息，请参考：

- [Uno 框架文档](../README.md)
- [Notify 组件重构规范](../refactor/notify.md)
- [HTMX 集成指南](https://htmx.org/)
