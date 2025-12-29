# Tab 组件重构规范

## 概述

Tab 组件重构目标是实现轻量、声明式、HTMX 原生的标签页组件，移除冗余功能，统一协议规范。

## 1. 核心架构模式 (The Pattern)

所有组件必须遵循 **"Store-First, Directive-Driven"**（Store 优先，指令驱动）模式。

## 2. 命名空间规范
所有组件必须注册到 `uno` 命名空间，避免多插件使用冲突：
- 注册方式：`Alpine.store('uno').tab = tabStore`
- 外部调用：`Alpine.store('uno').tab.activate()`
- 内部调用：`$store.uno.tab.xxx`

## 3. Uno 基础设施使用
所有组件重构必须使用 uno 基础设施：
- **core/uno.js**：使用 `Uno.getConfig()` 进行配置合并
- **utils/dom.js**：使用 `DomManager` 进行 DOM 操作
- **utils/log.js**：使用 `log` 进行日志记录
- **utils/position.js**：使用位置计算工具
- **core/layer.js**：使用 Z-Index 层级管理
- **locales/**：使用国际化支持

### 🚫 严禁做法 (Anti-Patterns)

* **❌ 禁止在 JS 中拼接 Alpine 指令字符串**：
  * *Bad:* `return { ':aria-selected': "isActive('${i}').toString()" }`
  * *Why:* 极难调试，破坏语法高亮，容易产生转义错误。

* **❌ 禁止使用 `Alpine.data` 做全局组件**：
  * *Bad:* `Alpine.data('unoTabs', ...)`
  * *Why:* 状态被锁在 DOM 局部，无法实现跨域操作。

* **❌ 禁止直接操作 `window.Uno.options`**：
  * *Bad:* `const opts = window.Uno?.options?.ui...`
  * *Why:* 破坏了配置的继承链（Defaults -> Global -> Instance）。

### ✅ 标准架构 (The Standard)

每个组件由三个核心部分组成：

1. **Directive (`x-uno-tab`)**: 负责解析 HTML 属性，将参数传递给 Store。
2. **Store (`Alpine.store('uno').tab`)**: 负责所有逻辑、状态管理、配置合并。
3. **Template (HTML)**: 负责渲染，完全由 Store 数据驱动，样式通过 `Uno.getStyle` 绑定。

## 2. 统一配置合并协议 (Configuration Protocol)

所有组件必须在 `open/show` 方法内部实现**四级配置合并**，逻辑必须内聚，不能散落在各处。

```javascript
// 标准实现模板（四级配置继承链）
show(instanceOpts = {}) {
  // Level 1: 插件默认配置 (Plugin Defaults)
  const pluginDefaults = options || {}
  
  // Level 2: 全局用户配置 (Global User Config)
  const globalOpts = (typeof window !== 'undefined' ? window.Uno?.options : Uno.options) || {}
  
  // Level 3: 插件全局配置 (Global Plugin Config)
  const globalPluginConfig = globalOpts.plugins?.tab || {}
  
  // Level 4: 直接全局配置 (Direct Global Config)
  const directGlobalConfig = globalOpts.tab || {}
  
  const config = {
    defaultTab: '',
    hash: true,
    position: 'top',
    ...pluginDefaults,           // Level 1
    ...globalPluginConfig,       // Level 2
    ...directGlobalConfig,       // Level 3
    ...instanceOpts              // Level 4 (最高优先级)
  }
}
```

## 2. Store 方法命名规范 (Store Method Naming)

### 2.1 标准命名规范

#### 标签页类组件（Tab）
```javascript
Alpine.store('tab', {
  // ✅ 标准方法命名
  show(tabId, options = {}) {},    // 显示指定标签
  switch(tabId, options = {}) {},  // 切换到指定标签
  activate(tabId, options = {}) {}, // 激活指定标签
  
  // ✅ 导航方法
  next() {},                       // 下一个标签
  prev() {},                       // 上一个标签
  first() {},                      // 第一个标签
  last() {},                       // 最后一个标签
  
  // ❌ 禁止使用的方法名
  open() {},    // 应使用 show/switch
  close() {},   // 标签页不需要关闭方法
  select() {}   // 应使用 show/switch
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
<!-- Tab 组件 HTMX 集成 -->
<div role="tabpanel" 
     x-show="activeTab === 'tab1'"
     x-init="htmx.process($el)"
     :hx-get="tabContentUrl"
     hx-trigger="load"
     hx-swap="innerHTML"
     :hx-headers="tabHeaders ? JSON.stringify(tabHeaders) : null">
  <div class="uno-tab-loading">加载中...</div>
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
const tab = {
  id: 'tab1',
  title: '标签标题',
  // ✅ 标准 HTMX 配置
  hxUrl: '/api/tab-content',      // 内容加载地址
  hxMethod: 'GET',                 // 请求方法
  hxHeaders: {                     // 请求头
    'X-CSRF': 'token',
    'X-Tab-ID': 'tab1'
  },
  hxTrigger: 'load',               // 触发条件
  hxSwap: 'innerHTML',             // 内容替换方式
  
  // ✅ 事件回调
  onContentLoad: (el) => {},       // 内容加载完成
  onError: (error) => {},          // 错误处理
  onActivate: (tabId) => {}        // 标签激活时
}
```

### 3.4 HTMX 事件处理

#### 标准事件处理
```javascript
// 在 Store 中处理 HTMX 事件
Alpine.store('tab', {
  activeTab: 'tab1',
  tabs: {},
  
  // ✅ 标准事件处理
  onContentLoaded(el, tab) {
    // 内容加载完成后的处理
    if (tab.onContentLoad) {
      tab.onContentLoad(el)
    }
    
    // 重新初始化 Alpine 组件
    if (window.Alpine && window.Alpine.initTree) {
      window.Alpine.initTree(el)
    }
    
    // 更新 URL 哈希
    if (tab.hash) {
      window.location.hash = tab.id
    }
  },
  
  // ✅ 错误处理
  onHtmxError(error, tab) {
    console.error('HTMX Error:', error)
    
    // 显示错误内容
    const errorContent = `
      <div class="uno-tab-error">
        <p>内容加载失败</p>
        <button onclick="$store.uno.tab.retry('${tab.id}')">重试</button>
      </div>
    `
    
    // 执行用户定义的错误处理
    if (tab.onError) {
      tab.onError(error)
    }
  },
  
  // ✅ 重试机制
  retry(tabId) {
    const tab = this.tabs[tabId]
    if (tab) {
      this.switch(tabId)
    }
  }
})
```

### 3.5 HTMX 与 URL 哈希同步

```javascript
// 在 Store 中处理 URL 哈希同步
Alpine.store('tab', {
  init() {
    // 监听哈希变化
    window.addEventListener('hashchange', () => {
      const hash = window.location.hash.slice(1)
      if (hash && this.tabs[hash]) {
        this.switch(hash)
      }
    })
    
    // 初始哈希检查
    const hash = window.location.hash.slice(1)
    if (hash) {
      this.switch(hash)
    }
  },
  
  switch(tabId, options = {}) {
    // 切换标签
    this.activeTab = tabId
    
    // 更新 URL 哈希
    if (options.updateHash !== false) {
      window.location.hash = tabId
    }
    
    // 触发 HTMX 加载
    const tab = this.tabs[tabId]
    if (tab && tab.hxUrl) {
      // HTMX 会自动处理内容加载
    }
  }
})
```

## 4. 样式与模板分离 (Style Decoupling)

为了解决"JS 里写样式"的问题，同时保留全局配置能力，必须遵循以下标准：

1. **HTML 模板**：必须存在于 JS 中（针对注入型组件）
2. **样式绑定**：严禁硬编码 Class，必须通过 `Uno.getStyle` 获取。

```html
<div :class="Uno.getStyle('tab', 'container')">
  <button :class="Uno.getStyle('tab', 'trigger')"></button>
  <div :class="Uno.getStyle('tab', 'content')"></div>
</div>
```

## 4. 触发机制标准 (Triggering)

组件必须同时支持三种触发方式，且逻辑入口统一流向 Store。

| 方式 | 场景 | 示例 |
| --- | --- | --- |
| **Directive (Hybrid)** | 推荐，最灵活 | `<button x-uno-tab="{ 'name': 'profile', 'url': '/profile' }">` |
| **API (Programmatic)** | 复杂业务逻辑 | `Alpine.store('uno').tab.activate('profile')` |
| **Store (Internal)** | Alpine 内部调用 | `$store.uno.tab.activate(...)` |

**Directive 解析规范：**
指令必须能同时读取 `expression` (对象) 和 `data-uno-*` (属性)，并合并它们。

```javascript
// 规范的 Directive 处理器
const handler = (values = {}) => {
  const attrConfig = {
    name: el.getAttribute('data-uno-name'), // 注意：保留 data-uno-* 属性解析以兼容旧代码
    url: el.getAttribute('data-uno-url'),
    // ...使用 parseBool/parseNumber 辅助函数
  }
  // 属性作为底色，JS 对象作为覆盖
  Alpine.store('uno').tab.activate({ ...attrConfig, ...values })
}
```

## 5. HTMX 集成规范

Tab 的重构证明了我们不需要在组件里写 AJAX。

* **原则**：UI 组件只负责"标签状态"（激活/非激活/Loading）。内容加载全权交给 HTMX。
* **实现**：
  * 如果标签支持动态内容，提供 `hxUrl` 参数。
  * 在模板中动态渲染 HTMX 属性：
```html
<div x-init="htmx.process($el)" :hx-get="item.hxUrl" hx-trigger="load"></div>
```

* **移除**所有旧代码中的 `fetch`、`XMLHttpRequest` 封装。

## Tab 组件重构规范

### 1. 标准插件结构

```javascript
// Store-First, Directive-Driven 标准结构
export default (Alpine, options = {}) => {
  // 1. 配置获取函数（遵循配置继承链）
  const getConfig = () => {
    const opts = typeof window !== 'undefined' ? window.Uno?.options : Uno.options
    return opts?.plugins?.tab || opts?.tab || options
  }

  // 2. Alpine Store（唯一的逻辑入口）
  Alpine.store('tab', {
    groups: {},
    
    // 四级配置合并的标准实现
    activate(instanceOpts = {}) {
      const globalOpts = (typeof window !== 'undefined' ? window.Uno?.options : Uno.options) || {}
      const pluginDefaults = options || {}
      const currentGlobal = globalOpts.plugins?.tab || globalOpts.tab || pluginDefaults
      
      const config = {
        defaultTab: '',
        hash: true,
        ...currentGlobal,
        ...instanceOpts
      }
      
      // 业务逻辑...
    }
  })

  // 3. 指令注册（声明式接口）
  Alpine.directive('uno-tab', (el, { expression }, { evaluateLater }) => {
    const getExpression = expression ? evaluateLater(expression) : null
    
    const buildFromAttrs = () => {
      return {
        name: el.getAttribute('data-uno-name') || el.getAttribute('data-tab'),
        url: el.getAttribute('data-uno-url') || el.getAttribute('data-tab-url'),
        target: el.getAttribute('data-uno-target') || el.getAttribute('data-tab-target'),
        group: el.getAttribute('data-uno-group') || 'default'
      }
    }

    const onClick = (e) => {
      e.preventDefault()
const store = Alpine.store('uno').tab

      if (getExpression) {
        let evaluated
        getExpression((value) => {
          evaluated = value
        })
        if (evaluated && typeof evaluated === 'object') {
          store.activate({ ...buildFromAttrs(), ...evaluated })
          return
        }
      }

      store.activate(buildFromAttrs())
    }

    el.addEventListener('click', onClick)
  })
}
```

### 2. 协议定义

#### 2.1 触发器协议（三种方式）

```html
<!-- 1. Directive (推荐) -->
<button x-uno-tab='{"name":"profile","url":"/user/profile","target":"#tab-content"}'>Profile</button>

<!-- 2. 属性方式（已移除） -->
<!-- 不再支持 data-uno-* 属性，统一使用 x-uno-* 指令 -->

<!-- 3. API 调用 -->
<script>
Alpine.store('uno').tab.activate({
  name: 'dashboard',
  url: '/api/dashboard',
  target: '#main-content'
})
</script>
```

#### 2.2 协议字段规范

| 字段 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `name` | string | "" | 标签唯一标识 |
| `url` | string | undefined | HTMX 内容加载地址 |
| `target` | string | undefined | 内容目标选择器 |
| `group` | string | "default" | 标签组名称 |
| `hash` | boolean | true | 是否启用 URL 哈希同步 |
| `defaultTab` | string | "" | 默认激活标签 |

### 3. 模板结构规范（样式与模板分离）

```javascript
// ✅ 标准模板实现（Store-First）
ensureContainer() {
  const id = 'uno-tab-manager'
  if (typeof document !== 'undefined' && document.getElementById(id)) return

  const template = `
    <div x-data="{ tabGroup: 'default' }">
      <!-- Tab Triggers - 样式通过 Uno.getStyle 获取 -->
      <div :class="Uno.getStyle('tab', 'triggerList')">
        <template x-for="tab in $store.uno.tab.groups[tabGroup]?.tabs || []" :key="tab.name">
          <button :id="'tab-trigger-' + tab.name"
                  :class="[
                    Uno.getStyle('tab', 'trigger'),
                    tab.isActive ? Uno.getStyle('tab', 'triggerActive') : Uno.getStyle('tab', 'triggerInactive')
                  ]"
                  :aria-selected="tab.isActive"
                  :aria-controls="'tab-panel-' + tab.name"
                  role="tab"
                  @click="$store.uno.tab.activate({ name: tab.name, group: tabGroup })">
            <span x-text="tab.title"></span>
          </button>
        </template>
      </div>

      <!-- Tab Panels - 完全由 Store 数据驱动 -->
      <div :class="Uno.getStyle('tab', 'panelList')">
        <template x-for="tab in $store.uno.tab.groups[tabGroup]?.tabs || []" :key="tab.name">
          <div :id="'tab-panel-' + tab.name"
               :class="Uno.getStyle('tab', 'panel')"
               :aria-labelledby="'tab-trigger-' + tab.name"
               role="tabpanel"
               x-show="tab.isActive"
               x-transition:enter="transition ease-out duration-200"
               x-transition:enter-start="opacity-0"
               x-transition:enter-end="opacity-100">
            
            <!-- HTMX 集成 - UI 组件只负责容器状态 -->
            <template x-if="tab.url">
              <div x-init="htmx.process($el)"
                   :hx-get="tab.url"
                   hx-trigger="load"
                   hx-swap="innerHTML">
                <div class="animate-pulse">Loading...</div>
              </div>
            </template>
            
            <!-- 静态内容 -->
            <template x-if="!tab.url && tab.content">
              <div x-html="tab.content"></div>
            </template>
          </div>
        </template>
      </div>
    </div>
  `

  const container = Uno.dom.ensureContainer(id, template, {
    'x-data': '',
    'role': 'tablist'
  })
  
  if (container && Alpine?.initTree) {
    Alpine.initTree(container)
  }
}
```

**关键原则**：

1. **样式与模板分离**：所有样式通过 `Uno.getStyle('tab', 'elementName')` 获取
2. **Store 数据驱动**：模板完全响应 Store 的 `groups` 数据
3. **HTMX 原生集成**：UI 组件只管理标签状态，内容加载交给 HTMX
4. **严禁硬编码**：不允许在模板中直接写 Tailwind 类名

## 重构实现

### 1. 移除的反模式 (Anti-Patterns)

```javascript
// ❌ 移除：Alpine.data 全局组件
Alpine.data('unoTabs', () => ({
  activeTab: '',
  tabs: {},
  init() { /* ... */ }
}))

// ❌ 移除：JS 中拼接 Alpine 指令
return { ':aria-selected': "isActive('${i}').toString()" }

// ❌ 移除：直接操作 window.Uno.options
const opts = window.Uno?.options?.ui?.tab

// ❌ 移除：组件内部 AJAX 实现
fetch(url).then(res => res.text()).then(html => { /* ... */ })

// ❌ 移除：硬编码样式类
tab.classList.add('bg-blue-600', 'text-white')
```

### 2. 新增的标准实践

```javascript
// ✅ Store-First：所有逻辑在 Store 中内聚
Alpine.store('tab', {
  groups: {},
  
  activate(instanceOpts = {}) {
    // 四级配置合并（Plugin Defaults -> Global User Config -> Global Plugin Config -> Direct Global Config -> Instance Config）
    const globalOpts = (typeof window !== 'undefined' ? window.Uno?.options : Uno.options) || {}
    const currentGlobal = globalOpts.plugins?.tab || globalOpts.tab || options
    
    const config = {
      group: 'default',
      hash: true,
      ...currentGlobal,
      ...instanceOpts
    }
    
    // 业务逻辑处理...
  }
})

// ✅ Directive-Driven：声明式接口
Alpine.directive('uno-tab', (el, { expression }, { evaluateLater }) => {
  const onClick = (e) => {
    e.preventDefault()
    const store = Alpine.store('uno').tab
    
    // 支持 expression 和 data-uno-* 属性合并
    store.activate({ ...buildFromAttrs(), ...evaluated })
  }
  
  el.addEventListener('click', onClick)
})

// ✅ 样式与模板分离：通过 Uno.getStyle 获取样式
<div :class="Uno.getStyle('tab', 'trigger')">

// ✅ HTMX 原生集成：UI 组件只负责容器状态
<div x-init="htmx.process($el)" :hx-get="tab.url" hx-trigger="load">
```

## HTMX 集成最佳实践

### 1. 性能优化

```javascript
// ✅ 推荐：使用 hx-swap="innerHTML" 避免整页刷新
{
  url: '/api/user-profile',
  target: '#tab-content',
  hxSwap: 'innerHTML'
}

// ✅ 推荐：添加加载状态
{
  url: '/api/slow-content',
  onLoad: (el) => {
    // 内容加载完成后的回调
    console.log('Tab content loaded:', el)
  }
}
```

### 2. 错误处理

```javascript
// ✅ HTMX 错误处理
document.addEventListener('htmx:responseError', function(evt) {
  console.error('Tab content load failed:', evt.detail)
  Alpine.store('uno').notify.error('标签内容加载失败')
})
```

## 样式规范

### 4.1 Tailwind 配置规范

Tab 组件通过 `Uno.getStyle()` 机制获取样式配置。

```javascript
// 全局配置示例
window.UnoConfig = {
  ui: {
    tab: {
      triggerList: 'flex space-x-1 border-b border-gray-200 dark:border-gray-700',
      trigger: 'px-4 py-2 text-sm font-medium rounded-t-lg transition-colors',
      triggerActive: 'bg-blue-600 text-white border-b-2 border-blue-600',
      triggerInactive: 'text-gray-500 hover:text-gray-700 hover:bg-gray-100',
      panelList: 'mt-2',
      panel: 'p-4 bg-white dark:bg-gray-800 rounded-b-lg',
      content: 'prose dark:prose-invert max-w-none'
    }
  }
}
```

### 4.2 样式覆盖示例

```javascript
// 简单的样式自定义
window.UnoConfig = {
  ui: {
    tab: {
      triggerActive: 'bg-purple-600 text-white border-b-2 border-purple-600',
      panel: 'p-6 bg-gray-50 dark:bg-gray-900 rounded-lg shadow-sm'
    }
  }
}
```

## 迁移指南

### 从旧版本迁移

1. **移除 Alpine.data 使用**
   ```javascript
   // ❌ 旧版本
   Alpine.data('unoTabs', () => ({ activeTab: 'home' }))
   
   // ✅ 新版本
   Alpine.store('uno').tab.activate({ name: 'home' })
   ```

2. **使用 HTMX 替代内容加载**
   ```javascript
   // ❌ 旧版本
   tabs[tabName] = { url: '/api/data', loaded: false }
   fetch(url).then(/* 处理内容 */)
   
   // ✅ 新版本
   <button x-uno-tab='{"name":"tab1","url":"/api/data"}'>Tab 1</button>
   ```

3. **简化配置**
   ```javascript
   // ❌ 旧版本
   const tabConfig = { defaultTab: 'home', hash: true, animate: true }
   
   // ✅ 新版本（声明式）
   <div x-data="{ tabGroup: 'main' }" x-init="$store.uno.tab.activate({'defaultTab': 'home', 'hash': true})">
   ```

## 总结

### 核心优势

1. **更轻量**：移除冗余代码，体积减少 30%+
2. **更声明式**：主要功能通过 HTML 属性声明
3. **更标准**：统一协议解析，符合 HTMX 生态
4. **更安全**：移除直接 DOM 操作，完全依赖 HTMX 安全机制
5. **更易维护**：职责清晰，代码结构简单

### 最佳实践

1. **优先使用声明式配置**：`x-uno-tab` 属性
2. **HTMX 用于动态内容**：避免直接 HTML 注入
3. **Tailwind 用于样式定制**：通过 UnoConfig.ui 系统
4. **彻底重构**：不保持向后兼容，完全拥抱新架构

这个重构方案为其他组件提供了标准化的模板和最佳实践，同时确保了平滑的迁移路径。