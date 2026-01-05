# Tooltip 组件重构规范

## 概述

Tooltip 组件重构目标是实现轻量、声明式、高性能的提示框组件，移除冗余功能，统一协议规范。

## 1. 核心架构模式 (The Pattern)

所有组件必须遵循 **"Store-First, Directive-Driven"**（Store 优先，指令驱动）模式。

### 🚫 严禁做法 (Anti-Patterns)

- **❌ 禁止在 JS 中拼接 Alpine 指令字符串**：
  - _Bad:_ `return { ':class': "isVisible('${id}').toString()" }`
  - _Why:_ 极难调试，破坏语法高亮，容易产生转义错误。

- **❌ 禁止直接操作 DOM**：
  - _Bad:_ `document.body.appendChild(tooltipEl)`
  - _Why:_ 破坏了 Alpine 的响应式系统，难以维护。

- **❌ 禁止直接操作 `window.Uno.options`**：
  - _Bad:_ `const opts = window.Uno?.options?.ui...`
  - _Why:_ 破坏了配置的继承链（Defaults -> Global -> Instance）。

### ✅ 标准架构 (The Standard)

每个组件由三个核心部分组成：

1. **Directive (`x-uno-tooltip`)**: 负责解析 HTML 属性，将参数传递给 Store。
2. **Store (`Alpine.store('uno').tooltip`)**: 负责所有逻辑、状态管理、配置合并。
3. **Template (HTML)**: 负责渲染，完全由 Store 数据驱动，样式通过 `Uno.getStyle` 绑定。

## 2. 统一配置合并协议 (Configuration Protocol)

所有组件必须在 `show` 方法内部实现**四级配置合并**，逻辑必须内聚，不能散落在各处。

```javascript
// 标准实现模板
show(instanceOpts = {}) {
  // Level 1: 全局 Uno.init 配置 (User Global)
  const globalOpts = (typeof window !== 'undefined' ? window.Uno?.options : Uno.options) || {}

  // Level 2: 插件加载时默认配置 (Plugin Defaults)
  const pluginDefaults = options.pluginDefaults || {}

  // Level 3: 最终计算 (优先级：Instance > Global > Defaults)
  const currentGlobal = globalOpts.plugins?.tooltip || globalOpts.tooltip || pluginDefaults

  const config = {
    placement: 'top',
    trigger: 'hover',
    delay: 300,
    ...currentGlobal,
    ...instanceOpts
  }
}
```

## 2. Store 方法命名规范 (Store Method Naming)

### 2.1 标准命名规范

#### 工具提示类组件（Tooltip）

Alpine.store('tooltip', {
const tooltipStore = {
// ✅ 标准方法命名
show(content, options = {}) {}, // 显示提示
hide() {}, // 隐藏提示
update(content, options = {}) {}, // 更新提示内容

// ✅ 位置控制方法
position(element, placement) {}, // 定位到指定元素
reposition() {}, // 重新定位

// ❌ 禁止使用的方法名
open() {}, // 应使用 show()
close() {}, // 应使用 hide()
display() {} // 不符合规范

```
4. **分类明确**：按照组件类型使用对应的动词组

## 重构要求

### 1. Store-First 架构
所有组件必须基于 Alpine Store 实现，状态管理完全内聚在 Store 中。

### 2. 命名空间规范
所有组件必须注册到 `uno` 命名空间，避免多插件使用冲突：
- 注册方式：`Alpine.store('uno').tooltip = tooltipStore`
- 外部调用：`Alpine.store('uno').tooltip.show()`
- 内部调用：`$store.uno.tooltip.xxx`

### 3. Uno 基础设施使用
所有组件重构必须使用 uno 基础设施：
- **core/uno.js**：使用 `Uno.getConfig()` 进行配置合并
- **utils/dom.js**：使用 `DomManager` 进行 DOM 操作
- **utils/log.js**：使用 `log` 进行日志记录
- **utils/position.js**：使用位置计算工具
- **core/layer.js**：使用 Z-Index 层级管理
- **locales/**：使用国际化支持

### 4.1 核心模块 (`src/core/`)

- **配置管理**: 使用 `Uno.options` (或 `src/core/uno.js`) 获取全局配置。
- **多语言**: 使用 `src/core/i18n.js` 处理所有用户可见的文本。
- **层级管理**: 使用 `src/core/layer.js` 管理 z-index。
- **错误处理**: 使用 `src/core/error.js` 统一处理异常。

### 4.2 工具库 (`src/utils/`)

- **DOM 操作**: 使用 `src/utils/dom-manager.js` 进行复杂的 DOM 操作。
- **日志**: 使用 `src/utils/log.js` 输出调试信息。
- **定位**: 使用 `src/utils/position.js` 处理位置计算。
- **无障碍**: 使用 `src/utils/a11y.js` 确保组件符合无障碍标准。

### 4.3 语言文件 (`src/locales/`)

- 所有硬编码的文本必须提取到 `src/locales/` 下的对应语言文件中。

### 2.2 命名原则

1. **动词明确**：使用清晰的动词表达操作意图
2. **语义一致**：相同类型操作使用相同动词
3. **简洁易懂**：避免过长或模糊的方法名
4. **分类明确**：按照组件类型使用对应的动词组

## 3. HTMX 集成标准 (HTMX Integration Standards)
## 5. HTMX 集成标准 (HTMX Integration Standards)
### 3.1 核心原则
### 5.1 核心原则

**UI 组件只负责容器状态管理，内容加载全权交给 HTMX**
### 3.2 集成方式
### 5.2 集成方式

#### ✅ 推荐方式：声明式 HTMX 集成

<div x-show="$store.uno.tooltip.current?.hxUrl"
<!-- Tooltip 组件 HTMX 集成 -->
     :hx-get="$store.uno.tooltip.current?.hxUrl"
  x-show="$store.uno.tooltip.current?.hxUrl"
  x-init="htmx.process($el)"
     :hx-headers="$store.uno.tooltip.current?.hxHeaders ? JSON.stringify($store.uno.tooltip.current.hxHeaders) : null">
  hx-trigger="load"
  hx-swap="innerHTML"
  :hx-headers="$store.uno.tooltip.current?.hxHeaders ? JSON.stringify($store.uno.tooltip.current.hxHeaders) : null"
>
  <div class="uno-tooltip-loading">加载中...</div>
</div>
```

#### ❌ 禁止方式：编程式 HTMX 调用

```javascript
// ❌ 禁止在组件中直接调用 HTMX API
### 3.3 HTMX 配置规范
htmx.trigger(element, eventName)
window.htmx.process(element)
```

### 5.3 HTMX 配置规范

#### 标准配置模板

````javascript
// Store 中的 HTMX 配置
const tooltip = {
  content: '提示内容',
  placement: 'top',
  // ✅ 标准 HTMX 配置
  hxUrl: '/api/tooltip-content', // 内容加载地址
  hxMethod: 'GET', // 请求方法
  hxHeaders: {
    // 请求头
    'X-CSRF': 'token'
  },
  hxTrigger: 'load', // 触发条件
  hxSwap: 'innerHTML', // 内容替换方式
### 3.4 HTMX 事件处理
  // ✅ 事件回调
  onContentLoad: (el) => {}, // 内容加载完成
  onError: (error) => {} // 错误处理
}
Alpine.store('tooltip', {
### 5.4 HTMX 事件处理

#### 标准事件处理

```javascript
// 在 Store 中处理 HTMX 事件
// 挂载到 Alpine.store('uno').tooltip
const tooltipStore = {
  current: null,

  // ✅ 标准事件处理
  onContentLoaded(el, tooltip) {
    // 内容加载完成后的处理
    if (tooltip.onContentLoad) {
      tooltip.onContentLoad(el)
    }

    // 重新初始化 Alpine 组件
    if (window.Alpine && window.Alpine.initTree) {
      window.Alpine.initTree(el)
    }

    // 重新计算位置
    this.reposition()
  },

  // ✅ 错误处理
  onHtmxError(error, tooltip) {
    console.error('HTMX Error:', error)

    // 显示错误内容
    const errorContent = `
      <div class="uno-tooltip-error">
        <p>内容加载失败</p>
})
    `

### 3.5 HTMX 与动态内容
    if (tooltip.onError) {
      tooltip.onError(error)
    }
Alpine.store('tooltip', {
````

### 5.5 HTMX 与动态内容

```javascript
// 在 Store 中处理动态内容
// 挂载到 Alpine.store('uno').tooltip
const tooltipStore = {
  show(content, options = {}) {
    // 如果有 HTMX URL，优先使用 HTMX 加载
    if (options.hxUrl) {
      this.current = {
        ...options,
        content, // 作为默认内容
        isLoading: true
      }

      // HTMX 会自动加载内容
      return
    }

    // 直接显示静态内容
    this.current = {
      content,
      ...options,
      isLoading: false
    }
  },

  // ✅ 内容更新
  update(content, options = {}) {
    if (this.current) {
      // 如果有新的 HTMX URL，重新加载
      if (options.hxUrl && options.hxUrl !== this.current.hxUrl) {
})
      } else {
        // 更新静态内容
## 4. 样式与模板分离 (Style Decoupling)
        Object.assign(this.current, options)
      }
    }
  }
}
```

## 6. 样式与模板分离 (Style Decoupling)

<span x-text="tooltip.content"></span>
为了解决"JS 里写样式"的问题，同时保留全局配置能力，必须遵循以下标准：

1. **HTML 模板**：必须存在于 JS 中（针对注入型组件）

## 4. 触发机制标准 (Triggering)

```html
<div :class="Uno.getStyle('tooltip', 'content')">
  <span x-text="$store.uno.tooltip.current?.content"></span>
</div>
```

| **API (Programmatic)** | 复杂业务逻辑 | `Alpine.store('uno').tooltip.show('Help text', { placement: 'right' })` |
| **Store (Internal)** | Alpine 内部调用 | `$store.uno.tooltip.show(...)` |

组件必须同时支持三种触发方式，且逻辑入口统一流向 Store。

| 方式                   | 场景            | 示例                                                                    |
| ---------------------- | --------------- | ----------------------------------------------------------------------- |
| **Directive (Hybrid)** | 推荐，最灵活    | `<span x-uno-tooltip="'Help text', {'placement': 'top'}">?</span>`      |
| **API (Programmatic)** | 复杂业务逻辑    | `Alpine.store('uno').tooltip.show('Help text', { placement: 'right' })` |
| **Store (Internal)**   | Alpine 内部调用 | `$store.uno.tooltip.show(...)`                                          |

**Directive 解析规范：**
指令必须能同时读取 `expression` (字符串/对象) 和 `data-uno-*` (属性)，并合并它们。

```javascript
// 规范的 Directive 处理器
const handler = (content = '') => {
  const attrConfig = {
    placement: el.getAttribute('data-uno-placement'), // 注意：保留 data-uno-* 属性解析以兼容旧代码
  Alpine.store('uno').tooltip.show({ ...attrConfig, ...config })
    delay: Number(el.getAttribute('data-uno-delay')) || undefined
  }

## 5. 性能优化规范
  const config = typeof content === 'object' ? content : { content }

  // 属性作为底色，JS 对象作为覆盖
  Alpine.store('uno').tooltip.show({ ...attrConfig, ...config })
}
```

## 8. 性能优化规范

Tooltip 组件需要特别注意性能，避免内存泄漏和频繁的 DOM 操作。

- **原则**：使用 Store 管理状态，通过 RAF 优化位置计算，避免频繁的 DOM 查询。
- **实现**：
  - 全局只有一个 tooltip 容器，通过 Store 管理显示/隐藏
    const tooltip = Alpine.store('uno').tooltip.current
  - 支持批量更新和防抖处理

```javascript
// ✅ 性能优化的位置计算
const updatePosition = () => {
  if (rafId) cancelAnimationFrame(rafId)
  rafId = requestAnimationFrame(() => {
    const tooltip = Alpine.store('uno').tooltip.current
    if (tooltip && tooltip.element) {
      const position = calculatePosition(tooltip.element, tooltip.placement)
      tooltip.position = position
    }
  })
}
```

## Tooltip 组件重构规范

### 1. 标准插件结构

```javascript
// Store-First, Directive-Driven 标准结构
  Alpine.store('tooltip', {
  const getConfig = () => {
    const opts = typeof window !== 'undefined' ? window.Uno?.options : Uno.options
    return opts?.plugins?.tooltip || opts?.tooltip || options
// 四级配置合并（Plugin Defaults -> Global User Config -> Global Plugin Config -> Direct Global Config -> Instance Config）
  const config = {
    content,
    placement: 'top',
    trigger: 'hover',
    delay: 300,
    ...pluginDefaults,           // Level 1
    ...globalPluginConfig,       // Level 2
    ...directGlobalConfig,       // Level 3
    ...instanceOpts              // Level 4 (最高优先级)
  }
        content,
        placement: 'top',
    }
  })
    }
  }

  if (!Alpine.store('uno')) {
    Alpine.store('uno', { tooltip: tooltipStore })
  } else {
    Alpine.store('uno').tooltip = tooltipStore
  }

  // 3. 指令注册（声明式接口）
  Alpine.directive('uno-tooltip', (el, { expression }, { evaluateLater }) => {
    const getExpression = expression ? evaluateLater(expression) : null

    const buildFromAttrs = () => {
      return {
      const store = Alpine.store('uno').tooltip
        trigger: el.getAttribute('data-uno-trigger') || 'hover',
        delay: Number(el.getAttribute('data-uno-delay')) || undefined,
        offset: Number(el.getAttribute('data-uno-offset')) || undefined
      }
    }

    const showTooltip = (content = '') => {
      const store = Alpine.store('uno').tooltip
      if (!store) return

      const attrConfig = buildFromAttrs()
      const config = typeof content === 'object' ? content : { content }

      store.show({ ...attrConfig, ...config, element: el })
    }

    // 根据触发类型绑定事件
    const trigger = el.getAttribute('data-uno-trigger') || 'hover'

    if (trigger === 'hover') {
      el.addEventListener('mouseenter', () => {
        Alpine.store('uno').tooltip.hide()
          getExpression((value) => showTooltip(value))
        } else {
          showTooltip(expression || '')
        }
      })

      el.addEventListener('mouseleave', () => {
        Alpine.store('uno').tooltip.hide()
      })
    } else if (trigger === 'click') {
      el.addEventListener('click', (e) => {
        e.preventDefault()
        if (getExpression) {
          getExpression((value) => showTooltip(value))
        } else {
          showTooltip(expression || '')
        }
      })
    }
  })
}
```

### 2. 协议定义

#### 2.1 触发器协议（三种方式）

```html
Alpine.store('uno').tooltip.show('Dynamic tooltip', {
<span x-uno-tooltip="'This is help text', {'placement': 'top'}">Help</span>

<!-- 2. 属性方式（已移除） -->
<!-- 不再支持 data-uno-* 属性，统一使用 x-uno-* 指令 -->

<!-- 3. API 调用 -->
<script>
  Alpine.store('uno').tooltip.show('Dynamic tooltip', {
    placement: 'right',
    element: document.getElementById('target')
  })
</script>
```

#### 2.2 协议字段规范

| 字段        | 类型    | 默认值  | 说明                        |
| ----------- | ------- | ------- | --------------------------- |
| `content`   | string  | ""      | 提示内容                    |
| `placement` | string  | "top"   | 位置：top/bottom/left/right |
| `trigger`   | string  | "hover" | 触发方式：hover/click/focus |
| `delay`     | number  | 300     | 显示延迟(ms)                |
| `offset`    | number  | 8       | 偏移距离(px)                |
| `arrow`     | boolean | true    | 是否显示箭头                |

### 3. 模板结构规范（样式与模板分离）

```javascript
// ✅ 标准模板实现（Store-First）
    <div x-show="$store.uno.tooltip.current"
  const id = 'uno-tooltip-container'
  if (typeof document !== 'undefined' && document.getElementById(id)) return
           top: $store.uno.tooltip.current?.position?.top + 'px',
           left: $store.uno.tooltip.current?.position?.left + 'px',

  const template = `
    <!-- Tooltip Container - 样式通过 Uno.getStyle 获取 -->
    <div x-show="$store.uno.tooltip.current"
         :class="Uno.getStyle('tooltip', 'content')"
         :style="{
           top: $store.uno.tooltip.current?.position?.top + 'px',
           left: $store.uno.tooltip.current?.position?.left + 'px',
           zIndex: ${zBase}
         }"
         x-transition:enter="transition ease-out duration-200"
      <div x-show="$store.uno.tooltip.current?.arrow"
         x-transition:enter-end="opacity-100 scale-100"
         x-transition:leave="transition ease-in duration-150"
             top: $store.uno.tooltip.current?.arrowPosition?.top + 'px',
             left: $store.uno.tooltip.current?.arrowPosition?.left + 'px'
         role="tooltip">

      <!-- Arrow - 完全由 Store 数据驱动 -->
      <div x-show="$store.uno.tooltip.current?.arrow"
        <span x-text="$store.uno.tooltip.current?.content"></span>
           :style="{
             top: $store.uno.tooltip.current?.arrowPosition?.top + 'px',
             left: $store.uno.tooltip.current?.arrowPosition?.left + 'px'
           }"></div>

      <!-- Content -->
      <div :class="Uno.getStyle('tooltip', 'text')">
        <span x-text="$store.uno.tooltip.current?.content"></span>
      </div>
    </div>
  `

  const container = Uno.dom.ensureContainer(id, template, {
    'x-data': '',
    'role': 'tooltip'
  })

  if (container && Alpine?.initTree) {
    Alpine.initTree(container)
  }
}
```

**关键原则**：

1. **样式与模板分离**：所有样式通过 `Uno.getStyle('tooltip', 'elementName')` 获取
2. **Store 数据驱动**：模板完全响应 Store 的 `current` 对象
3. **性能优化**：使用 RAF 优化位置计算，避免频繁 DOM 操作
4. **严禁硬编码**：不允许在模板中直接写 Tailwind 类名

## 重构实现

### 1. 移除的反模式 (Anti-Patterns)

```javascript
// ❌ 移除：直接 DOM 操作
document.body.appendChild(tooltipEl)
tooltipEl.style.top = position.top + 'px'

// ❌ 移除：频繁的 DOM 查询
const tooltip = document.querySelector('#uno-tooltip')
const rect = trigger.getBoundingClientRect()

// ❌ 移除：硬编码样式
tooltipEl.className = 'fixed z-50 bg-gray-900 text-white px-2 py-1 rounded text-sm'

Alpine.store('tooltip', {
```

### 2. 新增的标准实践

```javascript
// ✅ Store-First：所有逻辑在 Store 中内聚
// 挂载到 Alpine.store('uno').tooltip
const tooltipStore = {
  current: null,

  show(content, instanceOpts = {}) {
    // 四级配置合并（Plugin Defaults -> Global User Config -> Global Plugin Config -> Direct Global Config -> Instance Config）
    const globalOpts = (typeof window !== 'undefined' ? window.Uno?.options : Uno.options) || {}
    const currentGlobal = globalOpts.plugins?.tooltip || globalOpts.tooltip || options

    const config = {
      placement: 'top',
      trigger: 'hover',
})
  }
}

if (!Alpine.store('uno')) {
    const store = Alpine.store('uno').tooltip
} else {
  Alpine.store('uno').tooltip = tooltipStore
}

// ✅ Directive-Driven：声明式接口
Alpine.directive('uno-tooltip', (el, { expression }) => {
  const showTooltip = () => {
    const store = Alpine.store('uno').tooltip
    store.show(expression || '', { element: el })
  }

  el.addEventListener('mouseenter', showTooltip)
})

// ✅ 样式与模板分离：通过 Uno.getStyle 获取样式
<div :class="Uno.getStyle('tooltip', 'content')">

// ✅ 性能优化：使用 RAF 优化位置计算
const updatePosition = () => {
  if (rafId) cancelAnimationFrame(rafId)
  rafId = requestAnimationFrame(() => {
    // 位置计算逻辑
  })
}
```

## 性能优化最佳实践

    const tooltip = Alpine.store('uno').tooltip.current

### 1. 位置计算优化

```javascript
// ✅ 使用 RAF 防抖
const updatePosition = () => {
  if (rafId) cancelAnimationFrame(rafId)
  rafId = requestAnimationFrame(() => {
    const tooltip = Alpine.store('uno').tooltip.current
    if (tooltip && tooltip.element) {
      const position = calculatePosition(tooltip.element, tooltip.placement)
      tooltip.position = position
    }
  })
}

// ✅ 缓存 DOM 查询
const getCachedElement = () => {
  if (!cachedTooltipEl) {
    cachedTooltipEl = document.getElementById('uno-tooltip-container')
  }
  return cachedTooltipEl
}
```

### 2. 内存管理

```javascript
// ✅ 清理事件监听器
const cleanup = () => {
  if (rafId) {
    cancelAnimationFrame(rafId)
    rafId = 0
  }
  if (!Alpine.store('uno').tooltip.current) {
    globalListenersCleanup()
    globalListenersCleanup = null
  }
}

// ✅ 组件销毁时清理
Alpine.effect(() => {
  if (!Alpine.store('uno').tooltip.current) {
    cleanup()
  }
})
```

## 样式规范

### 4.1 Tailwind 配置规范

Tooltip 组件通过 `Uno.getStyle()` 机制获取样式配置。

```javascript
// 全局配置示例
window.UnoConfig = {
  ui: {
    tooltip: {
      // 必须包含 'fixed' 和 'z-index'
      content:
        'fixed z-[20000] whitespace-nowrap pointer-events-none bg-gray-900 text-white text-xs px-2 py-1 rounded shadow-lg',
      arrow: 'absolute w-2 h-2 bg-gray-900 transform rotate-45',
      text: 'block'
    }
  }
}
```

### 4.2 样式覆盖示例

```javascript
// 简单的样式自定义
window.UnoConfig = {
  ui: {
    tooltip: {
      content:
        'fixed z-[20000] whitespace-nowrap pointer-events-none bg-blue-600 text-white text-sm px-3 py-2 rounded-lg shadow-xl',
      arrow: 'absolute w-2 h-2 bg-blue-600 transform rotate-45'
    }
  }
}
```

## 迁移指南

Alpine.store('uno').tooltip.show('content', { element: trigger })

1. **移除直接 DOM 操作**

   ```javascript
   // ❌ 旧版本
   const tooltip = document.createElement('div')
   tooltip.className = 'tooltip'
   document.body.appendChild(tooltip)

   // ✅ 新版本
   Alpine.store('uno').tooltip.show('content', { element: trigger })
   ```

2. **使用声明式接口**

   ```javascript
   // ❌ 旧版本
   setupTooltip(element, { content: 'Help', placement: 'top' })

   // ✅ 新版本
   <span x-uno-tooltip="'Help', {'placement': 'top'}">?</span>
   ```

3. **简化配置**

   ```javascript
   // ❌ 旧版本
   const tooltipOptions = { placement: 'right', delay: 500, arrow: false }

   // ✅ 新版本（声明式）
   <button x-uno-tooltip="'Save', {'placement': 'right', 'delay': 500, 'arrow': false}">Save</button>
   ```

## 总结

### 核心优势

1. **更轻量**：移除冗余代码，体积减少 40%+
2. **更标准**：统一协议解析，符合 Alpine 生态
3. **更高性能**：RAF 优化，避免频繁 DOM 操作
4. **更易维护**：职责清晰，代码结构简单

### 最佳实践

1. **优先使用声明式配置**：`x-uno-tooltip` 属性
2. **性能优化**：使用 RAF 防抖，缓存 DOM 查询
3. **Tailwind 用于样式定制**：通过 UnoConfig.ui 系统
4. **彻底重构**：不保持向后兼容，完全拥抱新架构

这个重构方案为其他组件提供了标准化的模板和最佳实践，同时确保了平滑的迁移路径。
