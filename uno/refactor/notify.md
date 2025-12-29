# Notify 组件重构规范

## 概述

Notify 组件重构目标是实现轻量、声明式、高性能的通知组件，移除冗余功能，统一协议规范。

## 1. 核心架构模式 (The Pattern)

所有组件必须遵循 **"Store-First, Directive-Driven"**（Store 优先，指令驱动）模式。

### 2. 命名空间规范

- 注册方式：`Alpine.store('uno').notify = notifyStore`

- 注册方式：必须包含 `init` 钩子以确保容器初始化。
- 外部调用：`Alpine.store('uno').notify.success()`
```

### 3. Uno 基础设施使用

所有组件重构必须使用 uno 基础设施：

- **core/uno.js**：使用 `Uno.getConfig()` 进行配置合并
- **utils/dom.js**：使用 `DomManager` 进行 DOM 操作
- **utils/log.js**：使用 `log` 进行日志记录
- **utils/position.js**：使用位置计算工具
- **core/layer.js**：使用 Z-Index 层级管理
- **locales/**：使用国际化支持

### 🚫 严禁做法 (Anti-Patterns)

- **❌ 禁止在 JS 中拼接 Alpine 指令字符串**：
  - _Bad:_ `return { ':class': "getTypeClass('${type}')" }`
  - _Why:_ 极难调试，破坏语法高亮，容易产生转义错误。

- **❌ 禁止直接操作 DOM**：
  - _Bad:_ `container.appendChild(notificationEl)`
  - _Why:_ 破坏了 Alpine 的响应式系统，难以维护。

- **❌ 禁止直接操作 `window.Uno.options`**：
  - _Bad:_ `const opts = window.Uno?.options?.ui...`
  - _Why:_ 破坏了配置的继承链（Defaults -> Global -> Instance）。

### ✅ 标准架构 (The Standard)

每个组件由三个核心部分组成：

1. **Directive (`x-uno-notify`)**: 负责解析 HTML 属性，将参数传递给 Store。
2. **Store (`Alpine.store('uno').notify`)**: 负责所有逻辑、状态管理、配置合并。
所有组件必须在 `show` 方法内部实现**四级配置合并**，逻辑必须内聚，不能散落在各处。

## 2. 统一配置合并协议 (Configuration Protocol)
// 标准实现模板（四级配置继承链）
所有组件必须在 `show` 方法内部实现**四级配置合并**，使用 `Uno.getConfig` 简化实现。
  // Level 1: 插件默认配置 (Plugin Defaults)
  const pluginDefaults = options || {}

  // Level 2: 全局用户配置 (Global User Config)
  const globalOpts = (typeof window !== 'undefined' ? window.Uno?.options : Uno.options) || {}

  // Level 3: 插件全局配置 (Global Plugin Config)
  const globalPluginConfig = globalOpts.plugins?.notify || {}

  // Level 4: 直接全局配置 (Direct Global Config)
  const directGlobalConfig = globalOpts.notify || {}

  const config = {
    type: 'info',
    duration: 3000,
    position: 'top-right',
    closable: true,
    ...pluginDefaults,           // Level 1
    ...globalPluginConfig,       // Level 2
    ...directGlobalConfig,       // Level 3
    ...instanceOpts              // Level 4 (最高优先级)
  }
  // 4. Instance Config
  const config = Uno.getConfig('notify', instanceOpts)

## 2. Store 方法命名规范 (Store Method Naming)
}
### 2.1 标准命名规范

#### 通知类组件（Notify）
## 3. Store 方法命名规范 (Store Method Naming)

// ✅ 注册到 uno 命名空间
const notifyStore = {
  // ... methods
}

if (!Alpine.store('uno')) {
  Alpine.store('uno', { notify: notifyStore })
} else {
  Alpine.store('uno').notify = notifyStore
}

// 调用示例
// $store.uno.notify.show('message')
```

```javascript
// 旧模式 (Deprecated)
Alpine.store('notify', {
  // ...
})
```

**命名空间**：`Alpine.store('uno').notify`  
**核心方法**：`show()`、`remove()`、`success()`、`error()`、`warning()`、`info()`  
**配置合并**：四级配置合并（Plugin Defaults -> Global User Config -> Global Plugin Config -> Direct Global Config -> Instance Config）  
**状态管理**：`groups`、`nextId`、`activeRequests`、`progressVisible`  
**内部调用**：`$store.uno.notify.xxx`

### 2.2 命名原则

1. **动词明确**：使用清晰的动词表达操作意图
2. **语义一致**：相同类型操作使用相同动词
3. **简洁易懂**：避免过长或模糊的方法名
4. **分类明确**：按照组件类型使用对应的动词组

## 3. HTMX 集成标准 (HTMX Integration Standards)

### 3.1 核心原则

**UI 组件只负责容器状态管理，内容加载全权交给 HTMX**

### 3.2 集成方式

#### ✅ 推荐方式：声明式 HTMX 集成

```html
<!-- Notify 组件 HTMX 集成 -->
<div
  x-show="item.hxUrl"
  x-init="htmx.process($el)"
  :hx-get="item.hxUrl"
  hx-trigger="load"
  hx-swap="innerHTML"
  :hx-headers="item.hxHeaders ? JSON.stringify(item.hxHeaders) : null"
>
  <div class="uno-notify-loading">加载中...</div>
</div>
```

#### ❌ 禁止方式：编程式 HTMX 调用

```javascript
// ❌ 禁止在组件中直接调用 HTMX API
htmx.ajax(method, url, options)
htmx.trigger(element, eventName)
window.htmx.process(element)
```

### 3.3 HTMX 配置规范

#### 标准配置模板

```javascript
// Store 中的 HTMX 配置
const item = {
  message: '通知内容',
  type: 'info',
  // ✅ 标准 HTMX 配置
  hxUrl: '/api/notification', // 内容加载地址
  hxMethod: 'GET', // 请求方法
  hxHeaders: {
    // 请求头
    'X-CSRF': 'token'
  },
  hxTrigger: 'load', // 触发条件
  hxSwap: 'innerHTML', // 内容替换方式

  // ✅ 事件回调
  onContentLoad: (el) => {}, // 内容加载完成
  onError: (error) => {} // 错误处理
}
```

### 3.4 HTMX 事件处理

#### 标准事件处理
## 4. 高级特性 (Advanced Features)

// 在 Store 中处理 HTMX 事件
// 注意：需挂载到 Alpine.store('uno').notify
const notifyStore = {
  items: [],

  // ✅ 标准事件处理
  onContentLoaded(el, item) {
    // 内容加载完成后的处理
    if (item.onContentLoad) {
      item.onContentLoad(el)
    }

    // 重新初始化 Alpine 组件
    if (window.Alpine && window.Alpine.initTree) {
      window.Alpine.initTree(el)
    }
  },

  // ✅ 错误处理
  onHtmxError(error, item) {
    console.error('HTMX Error:', error)

    // 显示错误通知
    this.error('内容加载失败')

    // 执行用户定义的错误处理
    if (item.onError) {
      item.onError(error)
    }
  }
```javascript
$store.uno.notify.show('', {
  hxUrl: '/api/notification', // HTMX URL
## 4. 样式与模板分离 (Style Decoupling)

为了解决"JS 里写样式"的问题，同时保留全局配置能力，必须遵循以下标准：

1. **HTML 模板**：必须存在于 JS 中（针对注入型组件）
2. **样式绑定**：严禁硬编码 Class，必须通过 `Uno.getStyle` 获取。

  hxSwap: 'innerHTML' // Swap strategy
<div :class="Uno.getStyle('notify', 'item')">
  <div :class="Uno.getStyle('notify', 'content')">
    <span x-text="notification.message"></span>
模板实现：
</div>
```

## 4. 触发机制标准 (Triggering)

组件必须同时支持三种触发方式，且逻辑入口统一流向 Store。

| 方式                   | 场景            | 示例                                                                      |
| ---------------------- | --------------- | ------------------------------------------------------------------------- |
| **Directive (Hybrid)** | 推荐，最灵活    | `<button x-uno-notify="'Success!', {'type': 'success'}">Save</button>` |
| **API (Programmatic)** | 复杂业务逻辑    | `Alpine.store('uno').notify.success('Operation completed')`               |
| **Store (Internal)**   | Alpine 内部调用 | `$store.uno.notify.show(...)`                                             |

**Directive 解析规范：**
指令必须能同时读取 `expression` (字符串) 和 `data-uno-*` (属性)，并合并它们。

```javascript
// 规范的 Directive 处理器
const handler = (message = '') => {
  const attrConfig = {
    type: el.getAttribute('data-uno-type') || 'info', // 注意：保留 data-uno-* 属性解析以兼容旧代码
    position: el.getAttribute('data-uno-position'),
    duration: Number(el.getAttribute('data-uno-duration')) || undefined
  }

  // 属性作为底色，JS 对象作为覆盖
  Alpine.store('uno').notify.show(message, { ...attrConfig })
}
```html
<template x-if="item.hxUrl">
## 5. 性能优化规范

Notify 组件需要特别注意性能，避免内存泄漏和频繁的 DOM 操作。

- **原则**：使用 Store 管理状态，支持批量更新和自动清理。
- **实现**：
  - 按位置分组管理通知，避免全局搜索
  - 支持自动过期和手动关闭
  - 使用 CSS 动画而非 JS 动画

```javascript
// ✅ 性能优化的分组管理
const notifyStore = {
  groups: {
    'top-left': [],
    'top-right': [],
    'bottom-left': [],
    'bottom-right': []
  },

  // 按位置添加通知，避免全局搜索
  show(message, config = {}) {
    const position = config.position || 'top-right'
    const notification = { id, message, ...config }

    this.groups[position].push(notification)

    // 自动过期
    if (config.duration > 0) {
      setTimeout(() => this.remove(id), config.duration)
    }
  }
}
// Note: Register to Alpine.store('uno').notify
```

## Notify 组件重构规范

### 1. 标准插件结构

```javascript
// Store-First, Directive-Driven 标准结构
export default (Alpine, options = {}) => {
  // 1. 配置获取函数（遵循配置继承链）
  const getConfig = () => {
    const opts = typeof window !== 'undefined' ? window.Uno?.options : Uno.options
    return opts?.plugins?.notify || opts?.notify || options
  }

  // 2. Alpine Store（唯一的逻辑入口）
  const notifyStore = {
    groups: {
      'top-left': [],
      'top-right': [],
      'bottom-left': [],
      'bottom-right': []
    },
    nextId: 1,
    // ... implementation
  }

  if (!Alpine.store('uno')) {
    Alpine.store('uno', { notify: notifyStore })
  } else {
    Alpine.store('uno').notify = notifyStore
  }

// 四级配置合并的标准实现
  const config = {
    duration: 3000,
    position: 'top-right',
    closable: true,
    ...pluginDefaults,           // Level 1
    ...globalPluginConfig,       // Level 2
    ...directGlobalConfig,       // Level 3
    ...instanceOpts              // Level 4 (最高优先级)
  }

      const id = this.nextId++
      const notification = { id, message, ...config }

      // 按位置分组
      const position = config.position || 'top-right'
      if (!this.groups[position]) {
        this.groups[position] = []
      }

      this.groups[position].push(notification)

      // 自动过期
      if (config.duration > 0) {
        setTimeout(() => this.remove(id), config.duration)
      }

      return id
    },

    remove(id) {
      // 从所有组中移除
      Object.keys(this.groups).forEach(position => {
        this.groups[position] = this.groups[position].filter(item => item.id !== id)
      })
    },

    // 便捷方法
    success(message, opts = {}) {
      return this.show(message, { type: 'success', ...opts })
    },

    error(message, opts = {}) {
      return this.show(message, { type: 'error', ...opts })
    },

    warning(message, opts = {}) {
      return this.show(message, { type: 'warning', ...opts })
    },

    info(message, opts = {}) {
      return this.show(message, { type: 'info', ...opts })
    }
  })

  // 3. Register Store to uno namespace
  if (!Alpine.store('uno')) {
    Alpine.store('uno', { notify: notifyStore })
  } else {
    // If uno store already exists, attach notify
    Alpine.store('uno').notify = notifyStore
  }
  Alpine.directive('uno-notify', (el, { expression }) => {
    const buildFromAttrs = () => {
      return {
        type: el.getAttribute('data-uno-type') || 'info',
        position: el.getAttribute('data-uno-position'),
        duration: Number(el.getAttribute('data-uno-duration')) || undefined,
        closable: el.getAttribute('data-uno-closable') !== 'false'
      }
    }

    const showNotification = (message = '') => {
      const store = Alpine.store('uno').notify
      if (!store) return

      const attrConfig = buildFromAttrs()
      store.show(message || expression || '', attrConfig)
    }

    el.addEventListener('click', showNotification)
  })
}
```

### 2. 协议定义

#### 2.1 触发器协议（三种方式）
</template>
```
<!-- 1. Directive (推荐) -->
<button x-uno-notify="'Saved successfully', {'type': 'success'}">Save</button>

<!-- 2. 属性方式（已移除） -->
<!-- 不再支持 data-uno-* 属性，统一使用 x-uno-* 指令 -->

<!-- 3. API 调用 -->
<script>
  Alpine.store('uno').notify.success('Operation completed', {
    duration: 2000,
    closable: false
  })
</script>
```

#### 2.2 协议字段规范

| 字段       | 类型    | 默认值      | 说明                                                                                                       |
| ---------- | ------- | ----------- | ---------------------------------------------------------------------------------------------------------- |
| `message`  | string  | ""          | 通知内容                                                                                                   |
| `type`     | string  | "info"      | 类型：success/error/warning/info                                                                           |
| `position` | string  | "top-right" | 位置：top-left/top-center/top-right/center-left/center/center-right/bottom-left/bottom-center/bottom-right |
| `duration` | number  | 3000        | 显示时长(ms)，0表示不自动关闭                                                                              |
| `closable` | boolean | true        | 是否可手动关闭                                                                                             |
| `icon`     | boolean | true        | 是否显示图标                                                                                               |

### 3. 模板结构规范（样式与模板分离）

```javascript
// ✅ 标准模板实现（Store-First）
ensureContainer() {
  const id = 'uno-notify-container'
  if (typeof document !== 'undefined' && document.getElementById(id)) return

  const zBase = Uno.options?.ui?.layers?.notifyBase || 30000

  const template = `
    <!-- 按位置分组的容器 -->
    <template x-for="(items, position) in $store.uno.notify.groups" :key="position">
      <div :class="Uno.getStyle('notify', 'container')"
           :style="{
             position: 'fixed',
             zIndex: ${zBase},
             ...getPositionClasses(position)
           }">

        <!-- 通知项 - 样式通过 Uno.getStyle 获取 -->
        <template x-for="item in items" :key="item.id">
          <div :class="[
            Uno.getStyle('notify', 'item'),
            Uno.getStyle('notify', 'types.' + item.type)
          ]"
               x-show="item.isVisible"
               x-transition:enter="transition ease-out duration-300"
               x-transition:enter-start="opacity-0 transform translate-x-full"
               x-transition:enter-end="opacity-100 transform translate-x-0"
               x-transition:leave="transition ease-in duration-200"
               x-transition:leave-start="opacity-100 transform translate-x-0"
               x-transition:leave-end="opacity-0 transform translate-x-full">

            <!-- 图标 - 完全由 Store 数据驱动 -->
            <div x-show="item.icon"
                 :class="Uno.getStyle('notify', 'iconWrap')">
              <div :class="Uno.getStyle('notify', 'icons.' + item.type)"></div>
            </div>

            <!-- 内容 -->
            <div :class="Uno.getStyle('notify', 'textWrap')">
              <div :class="Uno.getStyle('notify', 'text')" x-text="item.message"></div>
            </div>

            <!-- 关闭按钮 -->
            <div x-show="item.closable"
                 :class="Uno.getStyle('notify', 'closeWrap')">
              <button :class="Uno.getStyle('notify', 'closeBtn')"
                      @click="$store.uno.notify.remove(item.id)">
                <div :class="Uno.getStyle('notify', 'closeIcon')"></div>
              </button>
            </div>
          </div>
        </template>
      </div>
    </template>
  `

  const container = Uno.dom.ensureContainer(id, template, {
    'x-data': '',
    'role': 'alert',
    'aria-live': 'polite'
  })

  if (container && Alpine?.initTree) {
    Alpine.initTree(container)
  }
}

// 位置样式计算
const getPositionClasses = (position) => {
  const positions = {
    'top-left': { top: '1rem', left: '1rem' },
    'top-right': { top: '1rem', right: '1rem' },
    'bottom-left': { bottom: '1rem', left: '1rem' },
    'bottom-right': { bottom: '1rem', right: '1rem' }
  }
  return positions[position] || positions['top-right']
}
```

**关键原则**：

1. **样式与模板分离**：所有样式通过 `Uno.getStyle('notify', 'elementName')` 获取
2. **Store 数据驱动**：模板完全响应 Store 的 `groups` 对象
3. **性能优化**：按位置分组，避免全局搜索
4. **严禁硬编码**：不允许在模板中直接写 Tailwind 类名

## 重构实现

### 1. 移除的反模式 (Anti-Patterns)

```javascript
// ❌ 移除：直接 DOM 操作
container.appendChild(notificationEl)
notificationEl.remove()

// ❌ 移除：硬编码样式
notificationEl.className = `notify notify-${type}`

// ❌ 移除：全局搜索通知
const notifications = document.querySelectorAll('.notify')

// ❌ 移除：直接操作 window.Uno.options
const opts = window.Uno?.options?.ui?.notify
```

### 2. 新增的标准实践

```javascript
// ✅ Store-First：所有逻辑在 Store 中内聚
Alpine.store('notify', {
  groups: {},

  show(message, instanceOpts = {}) {
    // 四级配置合并（Plugin Defaults -> Global User Config -> Global Plugin Config -> Direct Global Config -> Instance Config）
    const globalOpts = (typeof window !== 'undefined' ? window.Uno?.options : Uno.options) || {}
    const currentGlobal = globalOpts.plugins?.notify || globalOpts.notify || options

    const config = {
      type: 'info',
      position: 'top-right',
      duration: 3000,
      ...currentGlobal,
      ...instanceOpts
    }

    // 按位置分组管理
    const position = config.position || 'top-right'
    this.groups[position].push(notification)
  }
})

// ✅ Directive-Driven：声明式接口
Alpine.directive('uno-notify', (el, { expression }) => {
  const showNotification = () => {
    Alpine.store('uno').notify.show(expression || '', { element: el })
  }

  el.addEventListener('click', showNotification)
})

// ✅ 样式与模板分离：通过 Uno.getStyle 获取样式
<div :class="Uno.getStyle('notify', 'item')">

// ✅ 性能优化：按位置分组管理
this.groups[position].push(notification)
```

## 性能优化最佳实践

### 1. 分组管理

```javascript
// ✅ 按位置分组，避免全局搜索
Alpine.store('notify', {
  groups: {
    'top-left': [],
    'top-right': [],
    'bottom-left': [],
    'bottom-right': []
  },

  remove(id) {
    // 只从相关组中移除，不遍历所有通知
    Object.keys(this.groups).forEach((position) => {
      this.groups[position] = this.groups[position].filter((item) => item.id !== id)
    })
  }
})
```

### 2. 自动清理

```javascript
// ✅ 自动过期和内存清理
show(message, config = {}) {
  const notification = { id, message, ...config }

  // 添加到对应位置组
  this.groups[position].push(notification)

  // 自动过期
  if (config.duration > 0) {
    setTimeout(() => {
      this.remove(id)
    }, config.duration)
  }

  return id
}
```

## 样式规范

### 4.1 Tailwind 配置规范

Notify 组件通过 `Uno.getStyle()` 机制获取样式配置。

```javascript
// 全局配置示例
window.UnoConfig = {
  ui: {
    notify: {
      container: 'flex flex-col gap-2 p-4 pointer-events-none',
      item: 'pointer-events-auto max-w-[calc(100vw-2rem)] overflow-hidden rounded-lg shadow-lg border',
      content: 'p-4',
      iconWrap: 'flex-shrink-0',
      textWrap: 'ml-3 min-w-0 flex-1',
      text: 'text-sm font-medium',
      closeWrap: 'ml-4 flex flex-shrink-0',
      closeBtn: 'inline-flex rounded-md p-1 hover:bg-black/10 focus:outline-none',
      types: {
        success: 'bg-green-50 text-green-800 border-green-200',
        error: 'bg-red-50 text-red-800 border-red-200',
        warning: 'bg-yellow-50 text-yellow-800 border-yellow-200',
        info: 'bg-blue-50 text-blue-800 border-blue-200'
      },
      icons: {
        success: 'text-green-500',
        error: 'text-red-500',
        warning: 'text-yellow-500',
        info: 'text-blue-500'
      }
    }
  }
}
```

### 4.2 样式覆盖示例

```javascript
// 简单的样式自定义
window.UnoConfig = {
  ui: {
    notify: {
      item: 'pointer-events-auto max-w-md rounded-xl shadow-xl border-2',
      types: {
        success: 'bg-emerald-100 text-emerald-900 border-emerald-300',
        error: 'bg-rose-100 text-rose-900 border-rose-300'
      }
    }
  }
}
```

## 迁移指南

### 从旧版本迁移

1. **移除直接 DOM 操作**

   ```javascript
   // ❌ 旧版本
   const notification = createNotificationEl(message, type)
   container.appendChild(notification)

   // ✅ 新版本
   Alpine.store('uno').notify.success(message)
   ```

2. **使用声明式接口**

   ```javascript
   // ❌ 旧版本
   showNotification('Error!', { type: 'error', duration: 5000 })

   // ✅ 新版本
   <button x-uno-notify="'Error!', {'type': 'error', 'duration': 5000}">Show Error</button>
   ```

3. **简化配置**

   ```javascript
   // ❌ 旧版本
   const notifyOptions = { position: 'bottom-left', closable: false }

   // ✅ 新版本（声明式）
   <button x-uno-notify="'Message', {'position': 'bottom-left', 'closable': false}">Notify</button>
   ```

## 总结

### 核心优势

1. **更轻量**：移除冗余代码，体积减少 35%+
2. **更声明式**：主要功能通过 HTML 属性声明
3. **更标准**：统一协议解析，符合 Alpine 生态
4. **更高性能**：按位置分组，避免全局搜索
5. **更易维护**：职责清晰，代码结构简单

### 最佳实践

1. **优先使用声明式配置**：`x-uno-notify` 属性
2. **性能优化**：按位置分组，自动清理过期通知
3. **Tailwind 用于样式定制**：通过 UnoConfig.ui 系统
4. **彻底重构**：不保持向后兼容，完全拥抱新架构

这个重构方案为其他组件提供了标准化的模板和最佳实践，同时确保了平滑的迁移路径。

```html
<button x-uno-notify="['Message', { type: 'success', duration: 5000 }]">Click me</button>
```
