# Valid 组件重构规范

## 概述

Valid 组件重构目标是实现轻量、声明式、HTMX 原生的表单验证组件，移除冗余功能，统一协议规范。

## 1. 核心架构模式 (The Pattern)

所有组件必须遵循 **"Store-First, Directive-Driven"**（Store 优先，指令驱动）模式。

### 🚫 严禁做法 (Anti-Patterns)

* **❌ 禁止在 JS 中拼接 Alpine 指令字符串**：
  * *Bad:* `return { ':class': "hasError('${name}')" }`
  * *Why:* 极难调试，破坏语法高亮，容易产生转义错误。

* **❌ 禁止直接操作 DOM**：
  * *Bad:* `input.classList.add('is-invalid')`
  * *Why:* 破坏了 Alpine 的响应式系统，难以维护。

* **❌ 禁止直接操作 `window.Uno.options`**：
  * *Bad:* `const opts = window.Uno?.options?.ui...`
  * *Why:* 破坏了配置的继承链（Defaults -> Global -> Instance）。

### ✅ 标准架构 (The Standard)

每个组件由三个核心部分组成：

1. **Directive (`x-uno-valid`)**: 负责解析 HTML 属性，将参数传递给 Store。
2. **Store (`Alpine.store('uno').valid`)**: 负责所有逻辑、状态管理、配置合并。
3. **HTMX 拦截器**：负责表单验证拦截，完全依赖 HTMX 生态。

### 2.1 命名空间规范
所有组件必须注册到 `uno` 命名空间，避免多插件使用冲突：
- 注册方式：`Alpine.store('uno').valid = validStore`
- 外部调用：`Alpine.store('uno').valid.validate()`
- 内部调用：`$store.uno.valid.xxx`

### 2.2 Uno 基础设施使用
所有组件重构必须使用 uno 基础设施：
- **core/uno.js**：使用 `Uno.getConfig()` 进行配置合并
- **utils/dom.js**：使用 `DomManager` 进行 DOM 操作
- **utils/log.js**：使用 `log` 进行日志记录
- **utils/validation.js**：使用验证规则工具
- **utils/position.js**：使用位置计算工具（错误显示）
- **core/layer.js**：使用 Z-Index 层级管理
- **locales/**：使用国际化支持

## 2. 统一配置合并协议 (Configuration Protocol)

所有组件必须在 `validate` 方法内部实现**四级配置合并**，逻辑必须内聚，不能散落在各处。

```javascript
// 标准实现模板（四级配置继承链）
validate(form, instanceOpts = {}) {
  // Level 1: 插件默认配置 (Plugin Defaults)
  const pluginDefaults = options || {}
  
  // Level 2: 全局用户配置 (Global User Config)
  const globalOpts = (typeof window !== 'undefined' ? window.Uno?.options : Uno.options) || {}
  
  // Level 3: 插件全局配置 (Global Plugin Config)
  const globalPluginConfig = globalOpts.plugins?.valid || {}
  
  // Level 4: 直接全局配置 (Direct Global Config)
  const directGlobalConfig = globalOpts.valid || {}
  
// 四级配置合并的标准实现
  const config = {
    realtime: true,
    showError: true,
    scrollOnError: true,
    ...pluginDefaults,           // Level 1
    ...globalPluginConfig,       // Level 2
    ...directGlobalConfig,       // Level 3
    ...instanceOpts              // Level 4 (最高优先级)
  }
}
```

## 2. Store 方法命名规范 (Store Method Naming)

### 2.1 标准命名规范

#### 验证类组件（Valid）
```javascript
Alpine.store('valid', {
  // ✅ 标准方法命名
  validate(form, options = {}) {}, // 验证表单
  clear(form) {},                  // 清除验证结果
  showError(field, message) {},    // 显示字段错误
  hideError(field) {},             // 隐藏字段错误
  
  // ✅ 批量操作方法
  validateAll() {},                // 验证所有表单
  clearAll() {},                   // 清除所有错误
  showErrors(errors) {},           // 批量显示错误
  
  // ❌ 禁止使用的方法名
  check() {},   // 应使用 validate()
  reset() {},   // 应使用 clear()
  remove() {}   // 应使用 hideError
})
```

### 2.2 命名原则

1. **动词明确**：使用清晰的动词表达操作意图
2. **语义一致**：相同类型操作使用相同动词
3. **简洁易懂**：避免过长或模糊的方法名
4. **分类明确**：按照组件类型使用对应的动词组

## 3. HTMX 集成标准 (HTMX Integration Standards)

### 3.1 核心原则

**UI 组件只负责容器状态管理，表单验证全权交给 HTMX 拦截器**

### 3.2 集成方式

#### ✅ 推荐方式：HTMX 拦截器集成
```javascript
// HTMX 拦截器自动处理表单验证
document.body.addEventListener('htmx:confirm', (evt) => {
  const el = evt.detail.elt
  const form = el.tagName === 'FORM' ? el : el.closest('form')
  
  // 只拦截带有 data-uno-valid 的表单
  if (form && form.hasAttribute('data-uno-valid')) {
    if (!form.checkValidity()) {
      evt.preventDefault()
      evt.stopPropagation()
      
      // 触发验证显示错误
      Alpine.store('uno').valid.validate(form)
      
      // 触发自定义事件
      form.dispatchEvent(new CustomEvent('uno:validation-failed', { bubbles: true }))
    }
  }
})
```

#### ❌ 禁止方式：手动 HTMX 验证
```javascript
// ❌ 禁止在组件中手动处理 HTMX 请求
form.addEventListener('submit', (e) => {
  if (form.hasAttribute('data-uno-valid')) {
    e.preventDefault()
    // 手动验证逻辑不应该在这里
  }
})
```

### 3.3 HTMX 配置规范

#### 标准配置模板
```javascript
// 表单元素配置
<form data-uno-valid 
      hx-post="/api/submit"
      hx-trigger="submit">
  <input type="text" 
         name="email" 
         required 
         data-uno-validate="email">
  <div data-uno-error="email"></div>
  
  <button type="submit">提交</button>
</form>
```

### 3.4 HTMX 事件处理

#### 标准事件处理
```javascript
// 在 Store 中处理 HTMX 验证事件
Alpine.store('uno').valid = {
  errors: {},
  
  // ✅ 表单验证
  validate(form, options = {}) {
    const formData = new FormData(form)
    const errors = {}
    
    // 遍历表单字段
    for (let [field, value] of formData.entries()) {
      const fieldElement = form.querySelector(`[name="${field}"]`)
      
      if (fieldElement && fieldElement.hasAttribute('data-uno-validate')) {
        const validateType = fieldElement.getAttribute('data-uno-validate')
        
        // 验证字段
        if (!this.validateField(field, value, validateType)) {
          errors[field] = this.getErrorMessage(validateType)
        }
      }
    }
    
    // 更新错误状态
    this.errors = errors
    
    // 返回验证结果
    return {
      isValid: Object.keys(errors).length === 0,
      errors
    }
  },
  
  // ✅ 字段验证
  validateField(field, value, validateType) {
    const validators = {
      email: (v) => /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(v),
      required: (v) => v.trim() !== '',
      minLength: (v, min) => v.length >= min,
      maxLength: (v, max) => v.length <= max
    }
    
    const validator = validators[validateType]
    return validator ? validator(value) : true
  },
  
  // ✅ 错误显示
  showError(field, message) {
    this.errors[field] = message
    
    // 更新 DOM
    const errorElement = document.querySelector(`[data-uno-error="${field}"]`)
    if (errorElement) {
      errorElement.textContent = message
      errorElement.style.display = 'block'
    }
  },
  
  // ✅ 错误隐藏
  hideError(field) {
    delete this.errors[field]
    
    // 更新 DOM
    const errorElement = document.querySelector(`[data-uno-error="${field}"]`)
    if (errorElement) {
      errorElement.style.display = 'none'
    }
  }
})
```

### 3.5 HTMX 与服务端验证

```javascript
// 处理服务端验证错误
document.body.addEventListener('htmx:afterRequest', (evt) => {
  const xhr = evt.detail.xhr
  const response = xhr.response
  
  // 检查是否是验证错误
  if (xhr.status === 422 && response.errors) {
    // 显示服务端验证错误
    Alpine.store('uno').valid.showErrors(response.errors)
    
    // 触发验证失败事件
    document.dispatchEvent(new CustomEvent('uno:server-validation-failed', {
      detail: { errors: response.errors }
    }))
  }
})
```

## 4. 样式与模板分离 (Style Decoupling)

为了解决"JS 里写样式"的问题，同时保留全局配置能力，必须遵循以下标准：

1. **样式绑定**：严禁硬编码 Class，必须通过 `Uno.getStyle` 获取。
2. **状态管理**：通过 Store 管理错误状态，避免直接 DOM 操作。

```html
<input :class="Uno.getStyle('valid', 'inputError')" x-show="$store.uno.valid.errors['fieldName']">
<div :class="Uno.getStyle('valid', 'errorText')" x-text="$store.uno.valid.errors['fieldName']"></div>
```

## 4. 触发机制标准 (Triggering)

组件必须同时支持三种触发方式，且逻辑入口统一流向 Store。

| 方式 | 场景 | 示例 |
| --- | --- | --- |
| **Directive (Hybrid)** | 推荐，最灵活 | `<input x-uno-valid="{ required: true, email: true }" name="email">` |
| **API (Programmatic)** | 复杂业务逻辑 | `Alpine.store('uno').valid.validate('email', 'test@example.com')` |
| **HTMX 拦截** | 表单提交验证 | `<form data-uno-valid="true">` |

**Directive 解析规范：**
指令必须能同时读取 `expression` (对象) 和 `data-uno-*` (属性)，并合并它们。

```javascript
// 规范的 Directive 处理器
const handler = (rules = {}) => {
  const attrConfig = {
    required: el.hasAttribute('data-uno-required'),
    email: el.hasAttribute('data-uno-email'),
    minLength: Number(el.getAttribute('data-uno-min-length')) || undefined
  }
  
  // 属性作为底色，JS 对象作为覆盖
  Alpine.store('uno').valid.addRules(el.name || el.id, { ...attrConfig, ...rules })
}
```

## 5. HTMX 集成规范

Valid 组件完全依赖 HTMX 生态，不提供自定义 AJAX 实现。

* **原则**：验证组件只负责"验证逻辑"和"错误显示"。表单提交全权交给 HTMX。
* **实现**：
  * 通过 HTMX 事件拦截器进行表单验证
  * 支持服务端验证和客户端验证
  * 错误信息通过 Store 管理

```javascript
// ✅ HTMX 拦截器实现
document.body.addEventListener('htmx:confirm', (evt) => {
  const form = evt.detail.elt.tagName === 'FORM' ? evt.detail.elt : evt.detail.elt.closest('form')
  
  // 只拦截 opt-in 的表单
  if (form && form.hasAttribute('data-uno-valid')) {
    if (!Alpine.store('uno').valid.validateForm(form)) {
      evt.preventDefault()
      form.dispatchEvent(new CustomEvent('uno:validation-failed', { bubbles: true }))
    }
  }
})
```

## Valid 组件重构规范

### 1. 标准插件结构

```javascript
// Store-First, Directive-Driven 标准结构
export default (Alpine, options = {}) => {
  // 1. 配置获取函数（遵循配置继承链）
  const getConfig = () => {
    const opts = typeof window !== 'undefined' ? window.Uno?.options : Uno.options
    return opts?.plugins?.valid || opts?.valid || options
  }

  // 2. Alpine Store（唯一的逻辑入口）
  Alpine.store('uno').valid = {
    errors: {},
    rules: {},
    
    // 四级配置合并的标准实现
    addRules(fieldName, instanceOpts = {}) {
      const globalOpts = (typeof window !== 'undefined' ? window.Uno?.options : Uno.options) || {}
      const pluginDefaults = options || {}
      const currentGlobal = globalOpts.plugins?.valid || globalOpts.valid || pluginDefaults
      
      const config = {
        realtime: true,
        showError: true,
        ...currentGlobal,
        ...instanceOpts
      }
      
      this.rules[fieldName] = config
    },
    
    validate(fieldName, value) {
      const rules = this.rules[fieldName]
      if (!rules) return true
      
      // 验证逻辑...
      const isValid = this.runValidation(value, rules)
      
      if (!isValid) {
        this.errors[fieldName] = rules.message || 'Validation failed'
      } else {
        delete this.errors[fieldName]
      }
      
      return isValid
    },
    
    validateForm(form) {
      const inputs = form.querySelectorAll('input, select, textarea')
      let isValid = true
      
      inputs.forEach(input => {
        const fieldName = input.name || input.id
        if (this.rules[fieldName]) {
          if (!this.validate(fieldName, input.value)) {
            isValid = false
          }
        }
      })
      
      return isValid
    }
  })

  // 3. 指令注册（声明式接口）
  Alpine.directive('uno-valid', (el, { expression }) => {
    const buildFromAttrs = () => {
      return {
        required: el.hasAttribute('data-uno-required'),
        email: el.hasAttribute('data-uno-email'),
        minLength: Number(el.getAttribute('data-uno-min-length')) || undefined,
        maxLength: Number(el.getAttribute('data-uno-max-length')) || undefined,
        pattern: el.getAttribute('data-uno-pattern') || undefined
      }
    }

    const rules = expression ? (typeof expression === 'object' ? expression : {}) : {}
    const attrConfig = buildFromAttrs()
    
    Alpine.store('uno').valid.addRules(el.name || el.id, { ...attrConfig, ...rules })
    
    // 实时验证
    if (Alpine.store('uno').valid.rules[el.name || el.id]?.realtime) {
      el.addEventListener('input', () => {
        Alpine.store('uno').valid.validate(el.name || el.id, el.value)
      })
    }
  })

  // 4. HTMX 拦截器
  if (typeof document !== 'undefined') {
    document.body.addEventListener('htmx:confirm', (evt) => {
      const form = evt.detail.elt.tagName === 'FORM' ? evt.detail.elt : evt.detail.elt.closest('form')
      
      if (form && form.hasAttribute('data-uno-valid')) {
        if (!Alpine.store('uno').valid.validateForm(form)) {
          evt.preventDefault()
          evt.stopPropagation()
        }
      }
    })
  }
}
```

### 2. 协议定义

#### 2.1 触发器协议（三种方式）

```html
<!-- 1. Directive (推荐) -->
<input x-uno-valid="{ required: true, email: true }" name="email" type="email">

<!-- 2. 属性方式（已移除） -->
<!-- 不再支持 data-uno-* 属性，统一使用 x-uno-* 指令和 HTMX 拦截器 -->

<!-- 3. HTMX 表单 -->
<form data-uno-valid="true" hx-post="/api/submit">
  <input name="username" x-uno-valid="{ required: true, minLength: 3 }">
  <button type="submit">Submit</button>
</form>
```

#### 2.2 协议字段规范

| 字段 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `required` | boolean | false | 是否必填 |
| `email` | boolean | false | 是否为邮箱格式 |
| `minLength` | number | undefined | 最小长度 |
| `maxLength` | number | undefined | 最大长度 |
| `pattern` | string | undefined | 正则表达式 |
| `realtime` | boolean | true | 是否实时验证 |
| `showError` | boolean | true | 是否显示错误信息 |
| `message` | string | undefined | 自定义错误信息 |

### 3. 样式绑定规范（样式与模板分离）

```html
<!-- ✅ 标准实现（Store-First） -->
<div x-data="{ fieldName: 'email' }">
  <!-- 输入框 - 样式通过 Uno.getStyle 获取 -->
  <input :name="fieldName"
         type="email"
         x-uno-valid="{ required: true, email: true }"
         :class="[
           Uno.getStyle('valid', 'input'),
           $store.uno.valid.errors[fieldName] ? Uno.getStyle('valid', 'inputError') : ''
         ]">
  
  <!-- 错误信息 - 完全由 Store 数据驱动 -->
  <div x-show="$store.uno.valid.errors[fieldName]"
       :class="Uno.getStyle('valid', 'errorText')"
       x-text="$store.uno.valid.errors[fieldName]">
  </div>
</div>

<!-- ✅ HTMX 表单实现 -->
<form data-uno-valid="true" 
      hx-post="/api/submit"
      hx-target="#result">
  <div class="space-y-4">
    <div>
      <label>Username</label>
      <input name="username"
             x-uno-valid="{ required: true, minLength: 3 }"
             :class="[
               Uno.getStyle('valid', 'input'),
               $store.uno.valid.errors['username'] ? Uno.getStyle('valid', 'inputError') : ''
             ]">
      <div x-show="$store.uno.valid.errors['username']"
           :class="Uno.getStyle('valid', 'errorText')"
           x-text="$store.uno.valid.errors['username']">
      </div>
    </div>
    
    <div>
      <label>Email</label>
      <input name="email"
             type="email"
             x-uno-valid="{ required: true, email: true }"
             :class="[
               Uno.getStyle('valid', 'input'),
               $store.uno.valid.errors['email'] ? Uno.getStyle('valid', 'inputError') : ''
             ]">
      <div x-show="$store.uno.valid.errors['email']"
           :class="Uno.getStyle('valid', 'errorText')"
           x-text="$store.uno.valid.errors['email']">
      </div>
    </div>
  </div>
  
  <button type="submit">Submit</button>
</form>

<div id="result"></div>
```

**关键原则**：

1. **样式与模板分离**：所有样式通过 `Uno.getStyle('valid', 'elementName')` 获取
2. **Store 数据驱动**：错误显示完全响应 Store 的 `errors` 对象
3. **HTMX 原生集成**：表单验证通过 HTMX 拦截器实现
4. **严禁硬编码**：不允许在模板中直接写 Tailwind 类名

## 重构实现

### 1. 移除的反模式 (Anti-Patterns)

```javascript
// ❌ 移除：直接 DOM 操作
input.classList.add('is-invalid')
errorDiv.textContent = 'This field is required'

// ❌ 移除：硬编码样式
input.className = 'border-red-500 focus:border-red-500'

// ❌ 移除：自定义 AJAX 实现
form.addEventListener('submit', async (e) => {
  e.preventDefault()
  const response = await fetch(url, { /* ... */ })
})

// ❌ 移除：直接操作 window.Uno.options
const opts = window.Uno?.options?.ui?.valid
```

### 2. 新增的标准实践

```javascript
// ✅ Store-First：所有逻辑在 Store 中内聚
Alpine.store('uno').valid = {
  errors: {},
  rules: {},
  
  validate(fieldName, value) {
    // 四级配置合并（Plugin Defaults -> Global User Config -> Global Plugin Config -> Direct Global Config -> Instance Config）
    const globalOpts = (typeof window !== 'undefined' ? window.Uno?.options : Uno.options) || {}
    const currentGlobal = globalOpts.plugins?.valid || globalOpts.valid || options
    
    // 验证逻辑处理...
  }
})

// ✅ Directive-Driven：声明式接口
Alpine.directive('uno-valid', (el, { expression }) => {
  const rules = expression || {}
  Alpine.store('uno').valid.addRules(el.name || el.id, rules)
  
  // 实时验证
  el.addEventListener('input', () => {
    Alpine.store('uno').valid.validate(el.name || el.id, el.value)
  })
})

// ✅ 样式与模板分离：通过 Uno.getStyle 获取样式
:class="[
  Uno.getStyle('valid', 'input'),
  $store.uno.valid.errors[fieldName] ? Uno.getStyle('valid', 'inputError') : ''
]"

// ✅ HTMX 原生集成：通过拦截器实现表单验证
document.body.addEventListener('htmx:confirm', (evt) => {
  if (form.hasAttribute('data-uno-valid')) {
    if (!Alpine.store('uno').valid.validateForm(form)) {
      evt.preventDefault()
    }
  }
})
```

## HTMX 集成最佳实践

### 1. 表单验证拦截

```javascript
// ✅ HTMX 验证拦截器
document.body.addEventListener('htmx:confirm', (evt) => {
  const form = evt.detail.elt.tagName === 'FORM' ? evt.detail.elt : evt.detail.elt.closest('form')
  
  if (form && form.hasAttribute('data-uno-valid')) {
    const isValid = Alpine.store('uno').valid.validateForm(form)
    
    if (!isValid) {
      evt.preventDefault()
      evt.stopPropagation()
      
      // 触发验证失败事件
      form.dispatchEvent(new CustomEvent('uno:validation-failed', { 
        bubbles: true,
        detail: { errors: Alpine.store('uno').valid.errors }
      }))
    }
  }
})
```

### 2. 服务端验证集成

```javascript
// ✅ 处理服务端验证错误
document.body.addEventListener('htmx:afterRequest', (evt) => {
  const response = evt.detail.xhr.response
  
  try {
    const data = JSON.parse(response)
    if (data.errors) {
      // 将服务端错误合并到 Store
      Object.assign(Alpine.store('uno').valid.errors, data.errors)
    }
  } catch (e) {
    // 忽略非 JSON 响应
  }
})
```

## 样式规范

### 4.1 Tailwind 配置规范

Valid 组件通过 `Uno.getStyle()` 机制获取样式配置。

```javascript
// 全局配置示例
window.UnoConfig = {
  ui: {
    valid: {
      input: 'w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent',
      inputError: 'border-red-500 focus:ring-red-500 focus:border-red-500 bg-red-50',
      errorText: 'text-red-600 text-sm mt-1',
      parentError: 'border-red-500 ring-red-500'
    }
  }
}
```

### 4.2 样式覆盖示例

```javascript
// 简单的样式自定义
window.UnoConfig = {
  ui: {
    valid: {
      input: 'w-full px-4 py-3 border-2 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500',
      inputError: 'border-rose-500 bg-rose-50 text-rose-900 placeholder-rose-300',
      errorText: 'text-rose-600 font-medium mt-2 text-sm'
    }
  }
}
```

## 迁移指南

### 从旧版本迁移

1. **移除直接 DOM 操作**
   ```javascript
   // ❌ 旧版本
   input.classList.add('is-invalid')
   errorDiv.textContent = message
   
   // ✅ 新版本
     Alpine.store('uno').valid.errors[fieldName] = message   ```

2. **使用声明式接口**
   ```javascript
   // ❌ 旧版本
   validateInput(input, { required: true, email: true })
   
   // ✅ 新版本
   <input x-uno-valid="{ required: true, email: true }">
   ```

3. **HTMX 表单集成**
   ```javascript
   // ❌ 旧版本
   form.addEventListener('submit', handleSubmit)
   
   // ✅ 新版本
   <form data-uno-valid="true" hx-post="/api/submit">
   ```

## 总结

### 核心优势

1. **更轻量**：移除冗余代码，体积减少 30%+
2. **更声明式**：主要功能通过 HTML 属性声明
3. **更标准**：统一协议解析，符合 HTMX 生态
4. **更安全**：完全依赖 HTMX 安全机制
5. **更易维护**：职责清晰，代码结构简单

### 最佳实践

1. **优先使用声明式配置**：`x-uno-valid` 属性
2. **HTMX 用于表单提交**：避免自定义 AJAX 实现
3. **Tailwind 用于样式定制**：通过 UnoConfig.ui 系统
4. **彻底重构**：不保持向后兼容，完全拥抱新架构

这个重构方案为其他组件提供了标准化的模板和最佳实践，同时确保了平滑的迁移路径。