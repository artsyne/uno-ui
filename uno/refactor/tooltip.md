# Tooltip 组件重构规范

## 概述

Tooltip 组件重构目标是实现轻量、声明式、高性能的提示框组件，移除冗余功能，统一协议规范。

## 1. 核心架构模式 (The Pattern)

所有组件必须遵循 **"Store-First, Directive-Driven"**（Store 优先，指令驱动）模式。

### 🚫 严禁做法 (Anti-Patterns)

* **❌ 禁止在 JS 中拼接 Alpine 指令字符串**：
  * *Bad:* `return { ':class': "isVisible('${id}').toString()" }`
  * *Why:* 极难调试，破坏语法高亮，容易产生转义错误。

* **❌ 禁止直接操作 DOM**：
  * *Bad:* `document.body.appendChild(tooltipEl)`
  * *Why:* 破坏了 Alpine 的响应式系统，难以维护。

* **❌ 禁止直接操作 `window.Uno.options`**：
  * *Bad:* `const opts = window.Uno?.options?.ui...`
  * *Why:* 破坏了配置的继承链（Defaults -> Global -> Instance）。

### ✅ 标准架构 (The Standard)

每个组件由三个核心部分组成：

1. **Directive (`x-uno-tooltip`)**: 负责解析 HTML 属性，将参数传递给 Store。
2. **Store (`Alpine.store('tooltip')`)**: 负责所有逻辑、状态管理、配置合并。
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
```javascript
Alpine.store('tooltip', {
  // ✅ 标准方法命名
  show(content, options = {}) {},  // 显示提示
  hide() {},                       // 隐藏提示
  update(content, options = {}) {}, // 更新提示内容
  
  // ✅ 位置控制方法
  position(element, placement) {},  // 定位到指定元素
  reposition() {},                  // 重新定位
  
  // ❌ 禁止使用的方法名
  open() {},    // 应使用 show()
  close() {},   // 应使用 hide()
  display() {} // 不符合规范
})
```

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
<!-- Tooltip 组件 HTMX 集成 -->
<div x-show="tooltip.hxUrl" 
     x-init="htmx.process($el)"
     :hx-get="tooltip.hxUrl"
     hx-trigger="load"
     hx-swap="innerHTML"
     :hx-headers="tooltip.hxHeaders ? JSON.stringify(tooltip.hxHeaders) : null">
  <div class="uno-tooltip-loading">加载中...</div>
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
const tooltip = {
  content: '提示内容',
  placement: 'top',
  // ✅ 标准 HTMX 配置
  hxUrl: '/api/tooltip-content',  // 内容加载地址
  hxMethod: 'GET',                 // 请求方法
  hxHeaders: {                     // 请求头
    'X-CSRF': 'token'
  },
  hxTrigger: 'load',               // 触发条件
  hxSwap: 'innerHTML',             // 内容替换方式
  
  // ✅ 事件回调
  onContentLoad: (el) => {},       // 内容加载完成
  onError: (error) => {}           // 错误处理
}
```

### 3.4 HTMX 事件处理

#### 标准事件处理
```javascript
// 在 Store 中处理 HTMX 事件
Alpine.store('tooltip', {
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
      </div>
    `
    
    // 执行用户定义的错误处理
    if (tooltip.onError) {
      tooltip.onError(error)
    }
  }
})
```

### 3.5 HTMX 与动态内容

```javascript
// 在 Store 中处理动态内容
Alpine.store('tooltip', {
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
        this.show(content, options)
      } else {
        // 更新静态内容
        this.current.content = content
        Object.assign(this.current, options)
      }
    }
  }
})
```

## 4. 样式与模板分离 (Style Decoupling)

为了解决"JS 里写样式"的问题，同时保留全局配置能力，必须遵循以下标准：

1. **HTML 模板**：必须存在于 JS 中（针对注入型组件）
2. **样式绑定**：严禁硬编码 Class，必须通过 `Uno.getStyle` 获取。

```html
<div :class="Uno.getStyle('tooltip', 'content')">
  <span x-text="tooltip.content"></span>
</div>
```

## 4. 触发机制标准 (Triggering)

组件必须同时支持三种触发方式，且逻辑入口统一流向 Store。

| 方式 | 场景 | 示例 |
| --- | --- | --- |
| **Directive (Hybrid)** | 推荐，最灵活 | `<span x-uno-tooltip="'Help text'" data-uno-placement="top">?</span>` |
| **API (Programmatic)** | 复杂业务逻辑 | `Alpine.store('tooltip').show('Help text', { placement: 'right' })` |
| **Store (Internal)** | Alpine 内部调用 | `$store.tooltip.show(...)` |

**Directive 解析规范：**
指令必须能同时读取 `expression` (字符串/对象) 和 `data-uno-*` (属性)，并合并它们。

```javascript
// 规范的 Directive 处理器
const handler = (content = '') => {
  const attrConfig = {
    placement: el.getAttribute('data-uno-placement'),
    trigger: el.getAttribute('data-uno-trigger'),
    delay: Number(el.getAttribute('data-uno-delay')) || undefined
  }
  
  // 如果 content 是对象，则合并配置
  const config = typeof content === 'object' ? content : { content }
  
  // 属性作为底色，JS 对象作为覆盖
  Alpine.store('tooltip').show({ ...attrConfig, ...config })
}
```

## 5. 性能优化规范

Tooltip 组件需要特别注意性能，避免内存泄漏和频繁的 DOM 操作。

* **原则**：使用 Store 管理状态，通过 RAF 优化位置计算，避免频繁的 DOM 查询。
* **实现**：
  * 全局只有一个 tooltip 容器，通过 Store 管理显示/隐藏
  * 使用 `requestAnimationFrame` 优化位置计算
  * 支持批量更新和防抖处理

```javascript
// ✅ 性能优化的位置计算
const updatePosition = () => {
  if (rafId) cancelAnimationFrame(rafId)
  rafId = requestAnimationFrame(() => {
    const tooltip = Alpine.store('tooltip').current
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
export default (Alpine, options = {}) => {
  // 1. 配置获取函数（遵循配置继承链）
  const getConfig = () => {
    const opts = typeof window !== 'undefined' ? window.Uno?.options : Uno.options
    return opts?.plugins?.tooltip || opts?.tooltip || options
  }

  // 2. Alpine Store（唯一的逻辑入口）
  Alpine.store('tooltip', {
    current: null,
    items: [],
    
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
      
      // 业务逻辑处理...
    }
  })

  // 3. 指令注册（声明式接口）
  Alpine.directive('uno-tooltip', (el, { expression }, { evaluateLater }) => {
    const getExpression = expression ? evaluateLater(expression) : null
    
    const buildFromAttrs = () => {
      return {
        placement: el.getAttribute('data-uno-placement') || 'top',
        trigger: el.getAttribute('data-uno-trigger') || 'hover',
        delay: Number(el.getAttribute('data-uno-delay')) || undefined,
        offset: Number(el.getAttribute('data-uno-offset')) || undefined
      }
    }

    const showTooltip = (content = '') => {
      const store = Alpine.store('tooltip')
      if (!store) return

      const attrConfig = buildFromAttrs()
      const config = typeof content === 'object' ? content : { content }
      
      store.show({ ...attrConfig, ...config, element: el })
    }

    // 根据触发类型绑定事件
    const trigger = el.getAttribute('data-uno-trigger') || 'hover'
    
    if (trigger === 'hover') {
      el.addEventListener('mouseenter', () => {
        if (getExpression) {
          getExpression((value) => showTooltip(value))
        } else {
          showTooltip(expression || '')
        }
      })
      
      el.addEventListener('mouseleave', () => {
        Alpine.store('tooltip').hide()
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
<!-- 1. Directive (推荐) -->
<span x-uno-tooltip="'This is help text'" data-uno-placement="top">Help</span>

<!-- 2. 属性方式（已移除） -->
<!-- 不再支持 data-uno-* 属性，统一使用 x-uno-* 指令 -->

<!-- 3. API 调用 -->
<script>
Alpine.store('tooltip').show('Dynamic tooltip', {
  placement: 'right',
  element: document.getElementById('target')
})
</script>
```

#### 2.2 协议字段规范

| 字段 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `content` | string | "" | 提示内容 |
| `placement` | string | "top" | 位置：top/bottom/left/right |
| `trigger` | string | "hover" | 触发方式：hover/click/focus |
| `delay` | number | 300 | 显示延迟(ms) |
| `offset` | number | 8 | 偏移距离(px) |
| `arrow` | boolean | true | 是否显示箭头 |

### 3. 模板结构规范（样式与模板分离）

```javascript
// ✅ 标准模板实现（Store-First）
ensureContainer() {
  const id = 'uno-tooltip-container'
  if (typeof document !== 'undefined' && document.getElementById(id)) return

  const zBase = Uno.options?.ui?.layers?.tooltipBase || 20000

  const template = `
    <!-- Tooltip Container - 样式通过 Uno.getStyle 获取 -->
    <div x-show="$store.tooltip.current"
         :class="Uno.getStyle('tooltip', 'content')"
         :style="{
           top: $store.tooltip.current?.position?.top + 'px',
           left: $store.tooltip.current?.position?.left + 'px',
           zIndex: ${zBase}
         }"
         x-transition:enter="transition ease-out duration-200"
         x-transition:enter-start="opacity-0 scale-95"
         x-transition:enter-end="opacity-100 scale-100"
         x-transition:leave="transition ease-in duration-150"
         x-transition:leave-start="opacity-100 scale-100"
         x-transition:leave-end="opacity-0 scale-95"
         role="tooltip">
      
      <!-- Arrow - 完全由 Store 数据驱动 -->
      <div x-show="$store.tooltip.current?.arrow"
           :class="Uno.getStyle('tooltip', 'arrow')"
           :style="{
             top: $store.tooltip.current?.arrowPosition?.top + 'px',
             left: $store.tooltip.current?.arrowPosition?.left + 'px'
           }"></div>
      
      <!-- Content -->
      <div :class="Uno.getStyle('tooltip', 'text')">
        <span x-text="$store.tooltip.current?.content"></span>
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

// ❌ 移除：直接操作 window.Uno.options
const opts = window.Uno?.options?.ui?.tooltip
```

### 2. 新增的标准实践

```javascript
// ✅ Store-First：所有逻辑在 Store 中内聚
Alpine.store('tooltip', {
  current: null,
  
  show(content, instanceOpts = {}) {
    // 四级配置合并（Plugin Defaults -> Global User Config -> Global Plugin Config -> Direct Global Config -> Instance Config）
    const globalOpts = (typeof window !== 'undefined' ? window.Uno?.options : Uno.options) || {}
    const currentGlobal = globalOpts.plugins?.tooltip || globalOpts.tooltip || options
    
    const config = {
      placement: 'top',
      trigger: 'hover',
      delay: 300,
      ...currentGlobal,
      ...instanceOpts,
      content
    }
    
    // 业务逻辑处理...
  }
})

// ✅ Directive-Driven：声明式接口
Alpine.directive('uno-tooltip', (el, { expression }) => {
  const showTooltip = () => {
    const store = Alpine.store('tooltip')
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

### 1. 位置计算优化

```javascript
// ✅ 使用 RAF 防抖
const updatePosition = () => {
  if (rafId) cancelAnimationFrame(rafId)
  rafId = requestAnimationFrame(() => {
    const tooltip = Alpine.store('tooltip').current
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
  if (globalListenersCleanup) {
    globalListenersCleanup()
    globalListenersCleanup = null
  }
}

// ✅ 组件销毁时清理
Alpine.effect(() => {
  if (!Alpine.store('tooltip').current) {
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
      content: 'fixed z-[20000] whitespace-nowrap pointer-events-none bg-gray-900 text-white text-xs px-2 py-1 rounded shadow-lg',
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
      content: 'fixed z-[20000] whitespace-nowrap pointer-events-none bg-blue-600 text-white text-sm px-3 py-2 rounded-lg shadow-xl',
      arrow: 'absolute w-2 h-2 bg-blue-600 transform rotate-45'
    }
  }
}
```

## 迁移指南

### 从旧版本迁移

1. **移除直接 DOM 操作**
   ```javascript
   // ❌ 旧版本
   const tooltip = document.createElement('div')
   tooltip.className = 'tooltip'
   document.body.appendChild(tooltip)
   
   // ✅ 新版本
   Alpine.store('tooltip').show('content', { element: trigger })
   ```

2. **使用声明式接口**
   ```javascript
   // ❌ 旧版本
   setupTooltip(element, { content: 'Help', placement: 'top' })
   
   // ✅ 新版本
   <span x-uno-tooltip="'Help'" data-uno-placement="top">?</span>
   ```

3. **简化配置**
   ```javascript
   // ❌ 旧版本
   const tooltipOptions = { placement: 'right', delay: 500, arrow: false }
   
   // ✅ 新版本（声明式）
   <button x-uno-tooltip="'Save'" data-uno-placement="right" data-uno-delay="500" data-uno-arrow="false">Save</button>
   ```

## 总结

### 核心优势

1. **更轻量**：移除冗余代码，体积减少 40%+
2. **更声明式**：主要功能通过 HTML 属性声明
3. **更标准**：统一协议解析，符合 Alpine 生态
4. **更高性能**：RAF 优化，避免频繁 DOM 操作
5. **更易维护**：职责清晰，代码结构简单

### 最佳实践

1. **优先使用声明式配置**：`x-uno-tooltip` 属性
2. **性能优化**：使用 RAF 防抖，缓存 DOM 查询
3. **Tailwind 用于样式定制**：通过 UnoConfig.ui 系统
4. **彻底重构**：不保持向后兼容，完全拥抱新架构

这个重构方案为其他组件提供了标准化的模板和最佳实践，同时确保了平滑的迁移路径。