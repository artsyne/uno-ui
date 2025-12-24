# Uno UI 基础设施重构规范

## 概述

Uno UI 基础设施重构目标是建立统一、高性能、可扩展的核心架构，支撑所有组件的重构和未来发展。

## 1. Core (核心) 重构规范

### 1.1 配置管理系统

#### 当前问题
- 配置获取逻辑分散在各个组件中
- 缺乏统一的配置继承链实现
- 样式解析机制不够灵活

#### 重构目标

```javascript
// core/uno.js - 增强配置管理
class UnoInstance {
  constructor() {
    this.version = '0.6.0'
    this.pluginDefaults = {} // 新增：插件默认配置存储
    this.performance = new PerformanceManager() // 新增：性能管理器
  }

  /**
   * ✅ 统一配置获取方法
   * 支持四级配置合并：Plugin Defaults -> Global User Config -> Global Plugin Config -> Direct Global Config -> Instance Config
   */
  getConfig(componentName, instanceOpts = {}) {
    const globalOpts = this.options || {}
    
    return {
      // Level 1: 插件默认配置
      ...this.pluginDefaults[componentName],
      // Level 2: 全局用户配置  
      ...globalOpts.plugins?.[componentName],
      // Level 3: 直接配置
      ...globalOpts[componentName],
      // Level 4: 实例配置（最高优先级）
      ...instanceOpts
    }
  }

  /**
   * ✅ 增强样式解析系统
   * 支持嵌套配置、动态令牌、继承等高级特性
   */
  getStyle(comp, elem, defaultCls) {
    const config = this._getUiValue(comp, elem)
    
    // 支持嵌套配置，如 types.success
    if (elem.includes('.')) {
      const [parent, child] = elem.split('.')
      const parentConfig = this._getUiValue(comp, parent)
      const childConfig = parentConfig?.[child]
      
      if (childConfig) {
        return this._resolveUiTokens(childConfig)
      }
    }
    
    return this._resolveUiTokens(config || defaultCls)
  }

  /**
   * ✅ 插件默认配置注册
   */
  registerPluginDefaults(name, defaults) {
    this.pluginDefaults[name] = defaults
  }
}
```

### 1.2 插件系统标准化

```javascript
// core/plugin.js - 新增插件基类
export class BasePlugin {
  constructor(name, Alpine, options = {}) {
    this.name = name
    this.Alpine = Alpine
    this.options = options
    this.store = null
    this.directive = null
  }

  /**
   * ✅ 标准插件初始化流程
   */
  init() {
    // 1. 注册默认配置
    this.registerDefaults()
    
    // 2. 初始化 Store
    this.initStore()
    
    // 3. 注册指令
    this.initDirective()
    
    // 4. 初始化拦截器（如 HTMX）
    this.initInterceptors()
  }

  /**
   * ✅ 注册插件默认配置
   */
  registerDefaults() {
    const defaults = this.getDefaults()
    Uno.registerPluginDefaults(this.name, defaults)
  }

  /**
   * ✅ 统一 Store 创建模式
   */
  initStore() {
    this.store = this.Alpine.store(this.name, this.createStore())
  }

  /**
   * ✅ 统一指令注册模式
   */
  initDirective() {
    const directiveName = `uno-${this.name}`
    this.Alpine.directive(directiveName, this.createDirective())
  }

  /**
   * ✅ 子类必须实现的方法
   */
  getDefaults() {
    throw new Error('Plugin must implement getDefaults()')
  }

  createStore() {
    throw new Error('Plugin must implement createStore()')
  }

  createDirective() {
    throw new Error('Plugin must implement createDirective()')
  }

  initInterceptors() {
    // 可选实现
  }
}
```

## 2. Locales (多语言) 重构规范

### 2.1 多语言核心系统

#### 当前问题
- 多语言支持不够灵活
- 缺乏动态语言切换能力
- 组件内部硬编码文本

#### 重构目标

```javascript
// locales/i18n.js - 完全重构
export default (Alpine, options = {}) => {
  Alpine.store('i18n', {
    currentLang: options.defaultLang || 'zh-CN',
    fallbackLang: options.fallbackLang || 'en',
    messages: {},
    loading: false,
    
    /**
     * ✅ 异步语言加载
     */
    async loadLang(lang) {
      if (this.messages[lang]) {
        this.currentLang = lang
        return
      }

      this.loading = true
      try {
        const messages = await import(`./locales/${lang}.js`)
        this.messages[lang] = messages.default
        this.currentLang = lang
        
        // 触发语言切换事件
        window.dispatchEvent(new CustomEvent('uno:lang-changed', { 
          detail: { lang: this.currentLang }
        }))
      } catch (e) {
        console.warn(`Failed to load language: ${lang}`)
        if (lang !== this.fallbackLang) {
          await this.loadLang(this.fallbackLang)
        }
      } finally {
        this.loading = false
      }
    },
    
    /**
     * ✅ 统一翻译方法
     */
    t(key, params = {}) {
      const message = this.getNestedValue(this.messages[this.currentLang], key)
      
      if (!message) {
        // 回退到默认语言
        const fallbackMessage = this.getNestedValue(this.messages[this.fallbackLang], key)
        return this.interpolate(fallbackMessage || key, params)
      }
      
      return this.interpolate(message, params)
    },
    
    /**
     * ✅ 支持嵌套键名
     */
    getNestedValue(obj, key) {
      return key.split('.').reduce((o, k) => o?.[k], obj)
    },
    
    /**
     * ✅ 参数插值
     */
    interpolate(message, params) {
      if (typeof message !== 'string') return message
      
      return message.replace(/\{\{(\w+)\}\}/g, (match, key) => {
        return params[key] !== undefined ? params[key] : match
      })
    },
    
    /**
     * ✅ 批量翻译
     */
    tBatch(keys, params = {}) {
      return keys.reduce((result, key) => {
        result[key] = this.t(key, params[key] || {})
        return result
      }, {})
    }
  })

  /**
   * ✅ 全局翻译指令
   */
  Alpine.directive('t', (el, { expression }) => {
    const updateText = () => {
      el.textContent = Alpine.store('i18n').t(expression)
    }
    
    // 初始翻译
    updateText()
    
    // 监听语言切换
    window.addEventListener('uno:lang-changed', updateText)
    
    // 清理函数
    Alpine.effect(() => {
      return () => {
        window.removeEventListener('uno:lang-changed', updateText)
      }
    })
  })

  /**
   * ✅ 属性翻译指令
   */
  Alpine.directive('t-attr', (el, { expression }) => {
    const [attr, key] = expression.split(',')
    
    const updateAttr = () => {
      el.setAttribute(attr, Alpine.store('i18n').t(key.trim()))
    }
    
    updateAttr()
    window.addEventListener('uno:lang-changed', updateAttr)
  })
}
```

### 2.2 语言文件标准化

```javascript
// locales/zh-CN.js - 标准化结构
export default {
  // 通用文本
  common: {
    ok: '确定',
    cancel: '取消',
    close: '关闭',
    save: '保存',
    delete: '删除',
    edit: '编辑',
    loading: '加载中...',
    error: '错误',
    success: '成功',
    warning: '警告',
    info: '信息'
  },
  
  // 组件特定文本
  box: {
    confirmDelete: '确认删除此项目？',
    operationSuccess: '操作成功',
    operationFailed: '操作失败'
  },
  
  tab: {
    loading: '标签加载中...',
    error: '标签加载失败'
  },
  
  tooltip: {
    clickToShow: '点击显示提示',
    hoverToShow: '悬停显示提示'
  },
  
  notify: {
    autoClose: '{{seconds}} 秒后自动关闭',
    closeManually: '点击手动关闭'
  },
  
  valid: {
    required: '此字段为必填项',
    email: '请输入有效的邮箱地址',
    minLength: '最少需要 {{min}} 个字符',
    maxLength: '最多允许 {{max}} 个字符',
    pattern: '格式不正确'
  }
}
```

## 3. Z-Index 管理系统

### 3.1 Z-Index 管理器

#### 问题分析
- 现状：文档中零星提到了层级，但未见统一的 Z-Index Manager
- 需求：避免 Modal 盖不住 Tooltip 或 Notify 被遮挡的问题

#### 重构目标

```javascript
// utils/z-index.js - 新增 Z-Index 管理器
export class ZIndexManager {
  constructor() {
    // 定义层级标准
    this.layers = {
      base: 1000,        // 基础层级
      dropdown: 2000,     // 下拉菜单
      tooltip: 3000,      // 提示框
      notify: 4000,       // 通知
      modal: 5000,        // 模态框
      drawer: 6000,       // 抽屉
      loading: 7000,      // 加载遮罩
      debug: 9999         // 调试工具
    }
    
    // 当前使用计数
    this.counters = {}
    Object.keys(this.layers).forEach(layer => {
      this.counters[layer] = 0
    })
  }
  
  /**
   * ✅ 获取指定层级的 Z-Index
   * @param {string} layer - 层级名称
   * @param {number} offset - 偏移量
   * @returns {number} Z-Index 值
   */
  get(layer, offset = 0) {
    const base = this.layers[layer] || this.layers.base
    const count = this.counters[layer] || 0
    
    return base + count + offset
  }
  
  /**
   * ✅ 增加层级计数（用于同层级多个元素）
   * @param {string} layer - 层级名称
   * @returns {number} 当前 Z-Index
   */
  increment(layer) {
    this.counters[layer] = (this.counters[layer] || 0) + 1
    return this.get(layer)
  }
  
  /**
   * ✅ 减少层级计数
   * @param {string} layer - 层级名称
   */
  decrement(layer) {
    if (this.counters[layer] > 0) {
      this.counters[layer]--
    }
  }
  
  /**
   * ✅ 重置指定层级计数
   * @param {string} layer - 层级名称
   */
  reset(layer) {
    this.counters[layer] = 0
  }
  
  /**
   * ✅ 获取层级信息
   * @param {string} layer - 层级名称
   * @returns {Object} 层级信息
   */
  getLayerInfo(layer) {
    return {
      base: this.layers[layer] || this.layers.base,
      current: this.get(layer),
      count: this.counters[layer] || 0
    }
  }
}

// 全局 Z-Index 管理器实例
export const zIndexManager = new ZIndexManager()
```

#### 集成到 Uno Core

```javascript
// core/uno.js - 集成 Z-Index 管理
import { zIndexManager } from '../utils/z-index.js'

class UnoInstance {
  constructor() {
    // ... 其他初始化
    this.zIndex = zIndexManager
  }
  
  /**
   * ✅ 获取组件 Z-Index
   */
  getComponentZIndex(componentName, instanceId = null) {
    const layerMap = {
      'box': 'modal',
      'tooltip': 'tooltip', 
      'notify': 'notify',
      'tab': 'base',
      'valid': 'base'
    }
    
    const layer = layerMap[componentName] || 'base'
    
    if (instanceId) {
      return this.zIndex.increment(layer)
    }
    
    return this.zIndex.get(layer)
  }
}
```

## 4. Utils (工具函数) 重构规范

### 4.1 性能管理系统

```javascript
// utils/performance.js - 新增性能管理器
export class PerformanceManager {
  constructor() {
    this.rafId = 0
    this.timeouts = new Map()
    this.intervals = new Map()
    this.observers = new Set()
  }
  
  /**
   * ✅ RAF 防抖
   */
  rafDebounce(fn) {
    if (this.rafId) {
      cancelAnimationFrame(this.rafId)
    }
    this.rafId = requestAnimationFrame(fn)
  }
  
  /**
   * ✅ 超时管理（支持命名和清理）
   */
  setTimeout(id, fn, delay) {
    this.clearTimeout(id)
    const timeoutId = setTimeout(() => {
      fn()
      this.timeouts.delete(id)
    }, delay)
    this.timeouts.set(id, timeoutId)
  }
  
  clearTimeout(id) {
    const timeoutId = this.timeouts.get(id)
    if (timeoutId) {
      clearTimeout(timeoutId)
      this.timeouts.delete(id)
    }
  }
  
  /**
   * ✅ 间隔管理
   */
  setInterval(id, fn, interval) {
    this.clearInterval(id)
    const intervalId = setInterval(() => {
      fn()
    }, interval)
    this.intervals.set(id, intervalId)
  }
  
  clearInterval(id) {
    const intervalId = this.intervals.get(id)
    if (intervalId) {
      clearInterval(intervalId)
      this.intervals.delete(id)
    }
  }
  
  /**
   * ✅ Intersection Observer 管理
   */
  observe(id, element, callback, options = {}) {
    this.unobserve(id)
    
    const observer = new IntersectionObserver(callback, options)
    observer.observe(element)
    
    this.observers.add({
      id,
      observer,
      element
    })
  }
  
  unobserve(id) {
    const item = Array.from(this.observers).find(item => item.id === id)
    if (item) {
      item.observer.unobserve(item.element)
      this.observers.delete(item)
    }
  }
  
  /**
   * ✅ 组件销毁时清理
   */
  cleanup() {
    if (this.rafId) {
      cancelAnimationFrame(this.rafId)
      this.rafId = 0
    }
    
    this.timeouts.forEach(timeoutId => clearTimeout(timeoutId))
    this.intervals.forEach(intervalId => clearInterval(intervalId))
    this.observers.forEach(item => item.observer.disconnect())
    
    this.timeouts.clear()
    this.intervals.clear()
    this.observers.clear()
  }
}
```

### 3.2 DOM 工具标准化

```javascript
// utils/dom.js - 标准化 DOM 管理
export class DomManager {
  /**
   * ✅ 统一的容器创建
   */
  static ensureContainer(id, template, attributes = {}) {
    let container = document.getElementById(id)
    if (!container) {
      container = document.createElement('div')
      container.id = id
      container.innerHTML = template
      
      Object.entries(attributes).forEach(([key, value]) => {
        if (key.startsWith('x-')) {
          // Alpine 指令
          container.setAttribute(key, value)
        } else {
          // 普通属性
          container.setAttribute(key, value)
        }
      })
      
      document.body.appendChild(container)
    }
    return container
  }
  
  /**
   * ✅ 安全的 HTML 设置
   */
  static setHTML(el, html, trusted = false) {
    if (trusted) {
      el.innerHTML = html
    } else {
      el.textContent = html
    }
  }
  
  /**
   * ✅ 性能优化的查询
   */
  static querySelector(selector, parent = document) {
    return parent.querySelector(selector)
  }
  
  static querySelectorAll(selector, parent = document) {
    return Array.from(parent.querySelectorAll(selector))
  }
  
  /**
   * ✅ 批量 DOM 操作
   */
  static batch(operations) {
    // 使用 DocumentFragment 优化性能
    const fragment = document.createDocumentFragment()
    
    operations.forEach(op => {
      if (op.type === 'create') {
        const el = document.createElement(op.tag)
        Object.entries(op.attributes || {}).forEach(([key, value]) => {
          el.setAttribute(key, value)
        })
        if (op.html) el.innerHTML = op.html
        fragment.appendChild(el)
      }
    })
    
    return fragment
  }
  
  /**
   * ✅ 事件委托
   */
  static delegate(parent, selector, event, handler) {
    parent.addEventListener(event, (e) => {
      const target = e.target.closest(selector)
      if (target) {
        handler.call(target, e)
      }
    })
  }
}
```

### 3.3 验证系统重构

```javascript
// utils/validation.js - 增强验证系统
export const ValidationRules = {
  required: (value) => value && value.trim() !== '',
  email: (value) => /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(value),
  minLength: (value, min) => value.length >= min,
  maxLength: (value, max) => value.length <= max,
  pattern: (value, regex) => new RegExp(regex).test(value),
  url: (value) => {
    try {
      new URL(value)
      return true
    } catch {
      return false
    }
  },
  phone: (value) => /^1[3-9]\d{9}$/.test(value),
  idCard: (value) => /(^\d{15}$)|(^\d{18}$)|(^\d{17}(\d|X|x)$)/.test(value)
}

export class Validator {
  constructor(rules = {}) {
    this.rules = rules
    this.i18n = null
  }
  
  /**
   * ✅ 设置多语言实例
   */
  setI18n(i18nStore) {
    this.i18n = i18nStore
  }
  
  /**
   * ✅ 统一验证接口
   */
  validate(fieldName, value, fieldRules = {}) {
    const rules = { ...this.rules[fieldName], ...fieldRules }
    const errors = []
    
    Object.entries(rules).forEach(([rule, params]) => {
      let isValid = false
      let ruleParams = {}
      
      if (typeof params === 'boolean' && params) {
        isValid = ValidationRules[rule](value)
      } else if (typeof params === 'object') {
        isValid = ValidationRules[rule](value, params.value)
        ruleParams = params
      } else if (typeof params !== 'boolean') {
        isValid = ValidationRules[rule](value, params)
        ruleParams = { value: params }
      }
      
      if (!isValid) {
        errors.push(this.getErrorMessage(rule, ruleParams))
      }
    })
    
    return {
      isValid: errors.length === 0,
      errors,
      field: fieldName
    }
  }
  
  /**
   * ✅ 批量验证
   */
  validateBatch(data) {
    const results = {}
    let overallValid = true
    
    Object.entries(data).forEach(([fieldName, value]) => {
      const result = this.validate(fieldName, value)
      results[fieldName] = result
      if (!result.isValid) {
        overallValid = false
      }
    })
    
    return {
      isValid: overallValid,
      results
    }
  }
  
  /**
   * ✅ 多语言错误消息
   */
  getErrorMessage(rule, params) {
    if (this.i18n) {
      return this.i18n.t(`validation.${rule}`, params)
    }
    
    // 回退到英文错误消息
    const messages = {
      required: 'This field is required',
      email: 'Please enter a valid email address',
      minLength: 'Must be at least {{min}} characters',
      maxLength: 'Must be no more than {{max}} characters',
      pattern: 'Invalid format'
    }
    
    return messages[rule] || `Validation failed: ${rule}`
  }
}
```

## 9. 建议的目录结构

```
src/
├── core/
│   ├── uno.js (增强配置管理)
│   ├── plugin.js (新增：插件基类)
│   ├── constants.js (新增：常量定义)
│   ├── events.js (新增：事件管理)
│   └── error-boundary.js (新增：错误边界)
├── locales/
│   ├── i18n.js (重构：多语言核心)
│   ├── zh-CN.js (标准化结构)
│   ├── zh-TW.js (标准化结构)
│   └── en.js (标准化结构)
├── utils/
│   ├── dom.js (标准化)
│   ├── performance.js (新增)
│   ├── validation.js (增强)
│   ├── position.js (优化)
│   ├── log.js (保持)
│   ├── a11y.js (增强：无障碍访问)
│   ├── z-index.js (新增：Z-Index 管理)
│   └── error-boundary.js (新增：错误边界)
├── components/
│   ├── base/ (新增：组件基类)
│   │   ├── base-component.js
│   │   └── base-plugin.js
│   ├── box.js
│   ├── tab.js
│   ├── tooltip.js
│   ├── notify.js
│   └── valid.js
└── index.js (更新：统一导出)
```

## 10. JSDoc 注释标准 (JSDoc Standards)

### 10.1 核心原则

**所有代码注释必须采用 JSDoc 标准，统一使用英文注释，遵循开源项目最佳实践。**

#### 语言要求
- **必须使用英文注释**：所有注释、文档、说明都必须使用英文
- **遵循 JSDoc 3.x 标准**：使用最新的 JSDoc 语法
- **保持简洁明了**：注释应该解释"为什么"而不是"做什么"

#### 注释原则
- **注释解释意图**：解释代码的目的和设计决策
- **避免显而易见**：不要注释明显的代码
- **保持更新**：代码变更时必须同步更新注释

### 10.2 函数注释标准

#### 基本函数注释
```javascript
/**
 * Opens a dialog box with the specified configuration.
 * 
 * @param {Object} config - The dialog configuration object
 * @param {string} [config.type='dialog'] - The type of dialog ('dialog', 'confirm', 'alert')
 * @param {string} [config.title=''] - The dialog title
 * @param {string} [config.content=''] - The dialog content
 * @param {string} [config.hxUrl] - URL for dynamic content loading via HTMX
 * @param {Function} [config.onOk] - Callback function when OK button is clicked
 * @param {Function} [config.onCancel] - Callback function when Cancel button is clicked
 * @param {Function} [config.onClose] - Callback function when dialog is closed
 * @returns {string} The unique ID of the created dialog
 * 
 * @example
 * // Open a simple dialog
 * const dialogId = Alpine.store('box').open({
 *   type: 'dialog',
 *   title: 'Information',
 *   content: 'Operation completed successfully'
 * });
 * 
 * @since 2.0.0
 */
open(config = {}) {
  // Implementation...
}
```

#### 异步函数注释
```javascript
/**
 * Loads content from the specified URL using HTMX.
 * 
 * @param {string} url - The URL to load content from
 * @param {Object} [options] - Additional options for the request
 * @param {string} [options.method='GET'] - HTTP method to use
 * @param {Object} [options.headers={}] - Additional headers to send
 * @param {number} [options.timeout=5000] - Request timeout in milliseconds
 * @returns {Promise<string>} The loaded HTML content
 * 
 * @throws {Error} When the request fails or times out
 * 
 * @since 2.0.0
 */
async function loadContent(url, options = {}) {
  // Implementation...
}
```

### 10.3 类注释标准

#### 类定义注释
```javascript
/**
 * Box component store for managing modal dialogs.
 * 
 * This store manages the lifecycle of modal dialogs, including creation,
 * positioning, and destruction. It supports multiple dialog types and
 * provides a unified API for all dialog operations.
 * 
 * @class BoxStore
 * @since 2.0.0
 * 
 * @example
 * // Access the box store
 * const boxStore = Alpine.store('box');
 * 
 * // Open a dialog
 * const dialogId = boxStore.open({
 *   type: 'dialog',
 *   title: 'Hello World',
 *   content: 'This is a dialog'
 * });
 */
Alpine.store('box', {
  // Store implementation...
});
```

### 10.4 配置对象注释

#### 配置接口注释
```javascript
/**
 * Configuration object for dialog boxes.
 * 
 * @typedef {Object} DialogConfig
 * @property {string} [type='dialog'] - The type of dialog
 * @property {string} [title=''] - The dialog title
 * @property {string} [content=''] - The dialog content
 * @property {string} [hxUrl] - URL for dynamic content loading
 * @property {number} [width] - Dialog width in pixels
 * @property {number} [height] - Dialog height in pixels
 * @property {string} [position='center'] - Dialog position
 * @property {boolean} [modal=true] - Whether the dialog is modal
 * @property {boolean} [closable=true] - Whether the dialog can be closed
 * @property {Function} [onOk] - Callback for OK button
 * @property {Function} [onCancel] - Callback for Cancel button
 * @property {Function} [onClose] - Callback for close event
 * @since 2.0.0
 */
```

### 10.5 注释最佳实践

#### ✅ 推荐的注释
```javascript
/**
 * Calculates the optimal position for a tooltip relative to its target element.
 * 
 * This algorithm considers viewport boundaries, element dimensions, and
 * available space to determine the best placement. It prioritizes keeping
 * the tooltip visible within the viewport.
 * 
 * @param {HTMLElement} target - The element the tooltip is attached to
 * @param {string} preferredPlacement - The preferred placement ('top', 'bottom', 'left', 'right')
 * @returns {Object} The calculated position with coordinates
 */
function calculateTooltipPosition(target, preferredPlacement) {
  // Complex logic implementation...
}
```

#### ❌ 避免的注释
```javascript
// ❌ 避免这样的注释
const x = 5; // Set x to 5

// ❌ 避免过度注释
// Function to open a dialog
function openDialog(config) {
  // Create a new dialog
  const dialog = createDialogElement();
  // Set the title
  dialog.title = config.title;
}
```

### 10.6 工具和验证

#### 推荐的 VS Code 扩展
- **JSDoc Comments** - 提供 JSDoc 语法高亮和自动补全
- **Document This** - 自动生成 JSDoc 注释模板
- **Better Comments** - 改善注释的可读性

#### 注释检查清单
- [ ] 所有公共函数都有 JSDoc 注释
- [ ] 所有类都有完整的类文档
- [ ] 所有复杂逻辑都有解释性注释
- [ ] 所有注释都使用英文
- [ ] 注释与代码保持同步
- [ ] 没有 TODO 注释（除非有明确的计划）
- [ ] 没有过时的注释
- [ ] 所有示例代码都是可执行的

## 11. 组件基类标准化

```javascript
// components/base/base-component.js - 新增
import { PerformanceManager } from '../../utils/performance.js'
import { DomManager } from '../../utils/dom.js'
import { Validator } from '../../utils/validation.js'
import { a11yManager } from '../../utils/a11y.js'
import { zIndexManager } from '../../utils/z-index.js'
import { errorBoundary } from '../../utils/error-boundary.js'

/**
 * Base component class that provides common functionality for all Uno UI components.
 * 
 * This class encapsulates shared utilities like performance management, DOM operations,
 * validation, accessibility, z-index management, and error handling. All components should
 * extend this base class to ensure consistent behavior and maintainability.
 * 
 * @class BaseComponent
 * @since 2.0.0
 * 
 * @example
 * // Create a custom component
 * class MyComponent extends BaseComponent {
 *   constructor(Alpine, options = {}) {
 *     super('myComponent', Alpine, options);
 *   }
 *   
 *   render() {
 *     // Custom rendering logic
 *   }
 * }
 */
export class BaseComponent {
  /**
   * Creates a new BaseComponent instance.
   * 
   * @param {string} name - The component name
   * @param {Object} Alpine - The Alpine.js instance
   * @param {Object} [options={}] - Component-specific options
   * @since 2.0.0
   */
  constructor(name, Alpine, options = {}) {
    this.name = name
    this.Alpine = Alpine
    this.options = options
    this.performance = new PerformanceManager()
    this.dom = new DomManager()
    this.validator = new Validator()
    this.a11y = a11yManager
    this.zIndex = zIndexManager
    this.errorBoundary = errorBoundary
    
    // Register error handlers
    this.errorBoundary.addHandler(this.handleError.bind(this))
  }
  
  /**
   * Gets the component configuration with proper inheritance chain.
   * 
   * @param {Object} [instanceOpts={}] - Instance-specific options
   * @returns {Object} The merged configuration object
   * @since 2.0.0
   */
  getConfig(instanceOpts = {}) {
    return Uno.getConfig(this.name, instanceOpts)
  }
  
  /**
   * Gets the style class for a component element.
   * 
   * @param {string} element - The element name
   * @param {string} [defaultCls=''] - Default class if no style is configured
   * @returns {string} The resolved style class
   * @since 2.0.0
   */
  getStyle(element, defaultCls = '') {
    return Uno.getStyle(this.name, element, defaultCls)
  }
  
  /**
   * Gets a translated string for the component.
   * 
   * @param {string} key - The translation key
   * @param {Object} [params={}] - Parameters for interpolation
   * @returns {string} The translated string
   * @since 2.0.0
   */
  t(key, params = {}) {
    const i18n = this.Alpine.store('i18n')
    return i18n?.t(`${this.name}.${key}`, params) || i18n?.t(key, params) || key
  }
  
  /**
   * Validates a field with the specified rules.
   * 
   * @param {string} fieldName - The field name
   * @param {*} value - The value to validate
   * @param {Object} rules - Validation rules
   * @returns {Object} Validation result with isValid and errors
   * @since 2.0.0
   */
  validate(fieldName, value, rules) {
    this.validator.setI18n(this.Alpine.store('i18n'))
    return this.validator.validate(fieldName, value, rules)
  }
  
  /**
   * Gets the z-index for the component layer.
   * 
   * @param {number} [offset=0] - Offset from base layer
   * @returns {number} The calculated z-index
   * @since 2.0.0
   */
  getZIndex(offset = 0) {
    const layerMap = {
      'box': 'modal',
      'tooltip': 'tooltip', 
      'notify': 'notify',
      'tab': 'base',
      'valid': 'base'
    }
    
    const layer = layerMap[this.name] || 'base'
    return this.zIndex.get(layer, offset)
  }
  
  /**
   * Sets accessibility attributes on an element.
   * 
   * @param {HTMLElement} element - The element to set attributes on
   * @param {Object} config - Accessibility configuration
   * @since 2.0.0
   */
  setA11y(element, config) {
    this.a11y.setA11yAttributes(element, config)
  }
  
  /**
   * Traps focus within a container element.
   * 
   * @param {HTMLElement} container - The container element
   * @returns {Function} Cleanup function to remove focus trap
   * @since 2.0.0
   */
  trapFocus(container) {
    return this.a11y.trapFocus(container)
  }
  
  /**
   * Announces a message to screen readers.
   * 
   * @param {string} message - The message to announce
   * @param {string} [priority='polite'] - The announcement priority
   * @since 2.0.0
   */
  announce(message, priority = 'polite') {
    this.a11y.announce(message, priority)
  }
  
  /**
   * Handles component-specific errors.
   * 
   * @param {Object} errorInfo - Error information object
   * @since 2.0.0
   */
  handleError(errorInfo) {
    if (errorInfo.componentName === this.name) {
      // Component-specific error handling
      console.error(`Component ${this.name} error:`, errorInfo)
    }
  }
  
  /**
   * Wraps a method with error boundary protection.
   * 
   * @param {Function} method - The method to wrap
   @param {...*} args - Arguments to pass to the method
   @returns {*} The method result or null if error occurs
   * @since 2.0.0
   */
  safeWrap(method, ...args) {
    return this.errorBoundary.wrapComponent(this.name, method)(...args)
  }
  
  /**
   * Cleans up component resources.
   * 
   * @since 2.0.0
   */
  destroy() {
    this.performance.cleanup()
    this.errorBoundary.removeHandler(this.handleError.bind(this))
    this.validator = null
  }
  
  /**
   * Called when component is initialized.
   * 
   * @since 2.0.0
   */
  onInit() {
    // Subclass can override
  }
  
  /**
   * Called when component is destroyed.
   * 
   * @since 2.0.0
   */
  onDestroy() {
    // Subclass can override
  }
}
```

## 6. 无障碍访问 (A11y) 系统

### 6.1 A11y 管理器

#### 问题分析
- 现状：提到较少，缺乏系统的无障碍支持
- 需求：现代 UI 库必须支持 aria-* 属性和键盘导航

#### 重构目标

```javascript
// utils/a11y.js - 增强 A11y 支持
export class A11yManager {
  constructor() {
    this.focusableSelectors = [
      'button:not([disabled])',
      'input:not([disabled])',
      'select:not([disabled])',
      'textarea:not([disabled])',
      'a[href]',
      '[tabindex]:not([tabindex="-1"])',
      '[contenteditable="true"]'
    ].join(', ')
    
    this.announcer = null
    this.initAnnouncer()
  }
  
  /**
   * ✅ 初始化屏幕阅读器通知器
   */
  initAnnouncer() {
    this.announcer = document.createElement('div')
    this.announcer.setAttribute('aria-live', 'polite')
    this.announcer.setAttribute('aria-atomic', 'true')
    this.announcer.style.position = 'absolute'
    this.announcer.style.left = '-10000px'
    this.announcer.style.width = '1px'
    this.announcer.style.height = '1px'
    this.announcer.style.overflow = 'hidden'
    document.body.appendChild(this.announcer)
  }
  
  /**
   * ✅ 为元素添加标准 A11y 属性
   */
  setA11yAttributes(element, config) {
    const {
      role,
      label,
      describedBy,
      expanded,
      selected,
      hidden,
      busy
    } = config
    
    if (role) element.setAttribute('role', role)
    if (label) element.setAttribute('aria-label', label)
    if (describedBy) element.setAttribute('aria-describedby', describedBy)
    if (expanded !== undefined) element.setAttribute('aria-expanded', expanded)
    if (selected !== undefined) element.setAttribute('aria-selected', selected)
    if (hidden !== undefined) element.setAttribute('aria-hidden', hidden)
    if (busy !== undefined) element.setAttribute('aria-busy', busy)
  }
  
  /**
   * ✅ 屏幕阅读器通知
   */
  announce(message, priority = 'polite') {
    if (this.announcer) {
      this.announcer.setAttribute('aria-live', priority)
      this.announcer.textContent = message
      
      // 清理通知
      setTimeout(() => {
        this.announcer.textContent = ''
      }, 1000)
    }
  }
  
  /**
   * ✅ 焦点管理
   */
  trapFocus(container) {
    const focusableElements = container.querySelectorAll(this.focusableSelectors)
    const firstElement = focusableElements[0]
    const lastElement = focusableElements[focusableElements.length - 1]
    
    const handleKeydown = (e) => {
      if (e.key === 'Tab') {
        if (e.shiftKey) {
          if (document.activeElement === firstElement) {
            lastElement.focus()
            e.preventDefault()
          }
        } else {
          if (document.activeElement === lastElement) {
            firstElement.focus()
            e.preventDefault()
          }
        }
      }
    }
    
    container.addEventListener('keydown', handleKeydown)
    firstElement?.focus()
    
    return () => {
      container.removeEventListener('keydown', handleKeydown)
    }
  }
  
  /**
   * ✅ 获取焦点元素列表
   */
  getFocusableElements(container) {
    return Array.from(container.querySelectorAll(this.focusableSelectors))
      .filter(el => {
        const style = window.getComputedStyle(el)
        return style.display !== 'none' && style.visibility !== 'hidden'
      })
  }
}

// 全局 A11y 管理器实例
export const a11yManager = new A11yManager()
```

### 6.2 组件 A11y 标准

```javascript
// 组件必须遵循的 A11y 标准
const A11yStandards = {
  // Box 组件
  box: {
    container: {
      role: 'dialog',
      'aria-modal': 'true',
      'aria-labelledby': 'box-title',
      'aria-describedby': 'box-content'
    },
    closeButton: {
      'aria-label': 'Close dialog'
    }
  },
  
  // Tooltip 组件
  tooltip: {
    container: {
      role: 'tooltip',
      'aria-hidden': 'false'
    }
  },
  
  // Tab 组件
  tab: {
    trigger: {
      role: 'tab',
      'aria-selected': 'false',
      'aria-controls': 'tab-panel-{id}'
    },
    panel: {
      role: 'tabpanel',
      'aria-labelledby': 'tab-trigger-{id}',
      'aria-hidden': 'true'
    }
  },
  
  // Notify 组件
  notify: {
    container: {
      role: 'alert',
      'aria-live': 'polite'
    }
  }
}
```

## 7. 错误边界 (Error Boundary) 系统

### 7.1 错误管理器

#### 问题分析
- 现状：虽然有 HTMX 错误处理，但缺乏全局错误捕获
- 需求：防止某个组件报错导致整个页面 Alpine 挂掉

#### 重构目标

```javascript
// utils/error-boundary.js - 新增错误边界系统
export class ErrorBoundary {
  constructor() {
    this.errorHandlers = new Set()
    this.setupGlobalHandlers()
  }
  
  /**
   * ✅ 设置全局错误处理器
   */
  setupGlobalHandlers() {
    // 全局 JavaScript 错误
    window.addEventListener('error', (event) => {
      this.handleError({
        type: 'javascript',
        message: event.message,
        filename: event.filename,
        lineno: event.lineno,
        colno: event.colno,
        error: event.error
      })
    })
    
    // Promise 错误
    window.addEventListener('unhandledrejection', (event) => {
      this.handleError({
        type: 'promise',
        message: event.reason?.message || 'Unhandled Promise Rejection',
        error: event.reason
      })
    })
    
    // Alpine 错误处理
    if (typeof window !== 'undefined' && window.Alpine) {
      window.Alpine.error = (error, instance, expression) => {
        this.handleError({
          type: 'alpine',
          message: error.message,
          expression,
          instance,
          error
        })
      }
    }
  }
  
  /**
   * ✅ 注册错误处理器
   */
  addHandler(handler) {
    this.errorHandlers.add(handler)
  }
  
  /**
   * ✅ 移除错误处理器
   */
  removeHandler(handler) {
    this.errorHandlers.delete(handler)
  }
  
  /**
   * ✅ 统一错误处理
   */
  handleError(errorInfo) {
    console.group('🚨 Uno UI Error Boundary')
    console.error('Error Type:', errorInfo.type)
    console.error('Message:', errorInfo.message)
    console.error('Error:', errorInfo.error)
    console.groupEnd()
    
    // 通知所有错误处理器
    this.errorHandlers.forEach(handler => {
      try {
        handler(errorInfo)
      } catch (e) {
        console.error('Error in error handler:', e)
      }
    })
    
    // 默认错误处理
    this.defaultErrorHandler(errorInfo)
  }
  
  /**
   * ✅ 默认错误处理器
   */
  defaultErrorHandler(errorInfo) {
    // 显示用户友好的错误提示
    if (errorInfo.type === 'alpine') {
      this.showUserError('页面交互出现问题，请刷新重试')
    } else if (errorInfo.type === 'javascript') {
      this.showUserError('页面执行出错，部分功能可能不可用')
    }
    
    // 发送错误到监控系统（如果配置了）
    this.reportError(errorInfo)
  }
  
  /**
   * ✅ 显示用户错误提示
   */
  showUserError(message) {
    // 创建简单的错误提示
    const errorEl = document.createElement('div')
    errorEl.textContent = message
    errorEl.style.cssText = `
      position: fixed;
      top: 20px;
      right: 20px;
      background: #ef4444;
      color: white;
      padding: 12px 16px;
      border-radius: 6px;
      z-index: 99999;
      font-family: system-ui, sans-serif;
      font-size: 14px;
      box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
    `
    
    document.body.appendChild(errorEl)
    
    // 3秒后自动移除
    setTimeout(() => {
      if (errorEl.parentNode) {
        errorEl.parentNode.removeChild(errorEl)
      }
    }, 3000)
  }
  
  /**
   * ✅ 错误上报（可选）
   */
  reportError(errorInfo) {
    // 这里可以集成错误监控服务
    if (window.Uno?.options?.errorReporting?.enabled) {
      // 发送错误到监控服务
      console.log('Error reported:', errorInfo)
    }
  }
  
  /**
   * ✅ 组件错误包装器
   */
  wrapComponent(componentName, componentFn) {
    return (...args) => {
      try {
        return componentFn(...args)
      } catch (error) {
        this.handleError({
          type: 'component',
          componentName,
          args,
          error
        })
        return null
      }
    }
  }
}

// 全局错误边界实例
export const errorBoundary = new ErrorBoundary()
```

## 8. 事件系统标准化

```javascript
// core/events.js - 新增统一事件管理
export class EventManager {
  constructor() {
    this.events = new Map()
  }
  
  /**
   * ✅ 注册事件监听器
   */
  on(event, handler, options = {}) {
    if (!this.events.has(event)) {
      this.events.set(event, new Set())
    }
    
    const listener = { handler, options }
    this.events.get(event).add(listener)
    
    // 如果是 DOM 事件，同时注册到 document
    if (options.dom) {
      document.addEventListener(event, handler, options)
    }
  }
  
  /**
   * ✅ 触发事件
   */
  emit(event, data) {
    const listeners = this.events.get(event)
    if (listeners) {
      listeners.forEach(({ handler }) => {
        try {
          handler(data)
        } catch (e) {
          console.error(`Event handler error for ${event}:`, e)
        }
      })
    }
    
    // 同时触发 DOM 事件
    window.dispatchEvent(new CustomEvent(`uno:${event}`, { detail: data }))
  }
  
  /**
   * ✅ 移除事件监听器
   */
  off(event, handler) {
    const listeners = this.events.get(event)
    if (listeners) {
      listeners.forEach(listener => {
        if (listener.handler === handler) {
          listeners.delete(listener)
        }
      })
    }
  }
  
  /**
   * ✅ 清理所有监听器
   */
  cleanup() {
    this.events.clear()
  }
}

// 全局事件管理器实例
export const eventManager = new EventManager()
```

## 总结

### 重构收益

1. **统一架构**：所有组件遵循相同的开发模式
2. **性能优化**：统一的性能管理和内存清理
3. **多语言支持**：完整的国际化解决方案
4. **开发效率**：标准化的工具函数和基类
5. **可维护性**：清晰的职责分离和代码组织

### 实施优先级

1. **高优先级**：多语言系统、Core 配置管理
2. **中优先级**：性能工具、验证系统
3. **低优先级**：事件系统、组件基类

这个重构方案将为 Uno UI 提供坚实的基础设施，支撑所有组件的现代化重构。