# Uno 插件延迟加载（Lazy Initialization）重构方案

## 1. 目标

优化 Uno 插件（Box, Notify, Tab, Tooltip, Valid）的初始化机制，实现**按需加载（Lazy Loading）**。
目前，部分插件在应用启动时立即注入 HTML 模版或绑定全局事件，无论页面是否实际使用了这些功能。重构的目标是将这些操作推迟到插件真正被使用时（例如打开对话框、显示通知或解析到相关指令时）执行。

## 2. 现状分析

### 2.1 Box 组件 (对话框)

- **入口 (`src/components/box/index.js`)**: 注册 Store。
- **Store (`src/components/box/box.store.js`)**: `init()` 生命周期方法调用了 `this.ensureContainer()`，导致应用启动时立即在 DOM 中注入 Box 的 HTML 结构（遮罩层、容器等）。
- **使用方式**: `open()`, `alert()`, `confirm()` 等方法内部已经调用了 `ensureContainer()`。

### 2.2 Notify 组件 (通知)

- **入口 (`src/components/notify/index.js`)**: 显式调用了 `notifyStore.init()`。
- **Store (`src/components/notify/notify.store.js`)**: `init()` 方法调用 `this.ensureContainer()`，注入通知容器和全局进度条。
- **使用方式**:
  - `show()` 会调用 `ensureContainer()`。
  - `startProgress()` (全局 Fetch/HTMX 进度条) **目前未调用** `ensureContainer()`，依赖于初始化时的注入。

### 2.3 Tab 组件 (选项卡)

- **入口 (`src/components/tab/index.js`)**: 显式调用了 `Alpine.store('uno').tab.init()`。
- **Store (`src/components/tab/tab.store.js`)**: `init()` 方法绑定了 `hashchange` 和 `load` 事件，用于 URL Hash 同步。
- **问题**: 即使页面没有 Tab 组件，也会监听 Hash 变化。

### 2.4 Tooltip 组件 (提示框)

- **入口 (`src/components/tooltip/index.js`)**: 未显式调用 `init()`。
- **Store (`src/components/tooltip/tooltip.store.js`)**: `init()` 方法包含了 `ensureContainer()`。虽然目前入口未调用，但为了防止未来自动调用导致 Eager 加载，应予以清理。
- **使用方式**: `show()` 方法已经正确处理了容器注入和事件绑定。

### 2.5 Valid 组件 (验证)

- **分析**: 该组件是纯逻辑/指令驱动的，没有 `init()` 方法，也不注入全局 DOM 元素。
- **结论**: **无需重构**，已经是按需执行。

## 3. 重构方案

### 3.1 Box 组件重构

**文件:** `src/components/box/box.store.js`

- **操作**: 移除 `init()` 方法中的 `this.ensureContainer()` 调用。
- **原因**: `open()` 及其衍生方法会在调用时自动注入容器。

```javascript
// src/components/box/box.store.js
init() {
  // 移除 this.ensureContainer()
},
```

### 3.2 Notify 组件重构

**文件:** `src/components/notify/index.js`

- **操作**: 移除显式的 `notifyStore.init()` 调用。

**文件:** `src/components/notify/notify.store.js`

- **操作 1**: 移除 `init()` 方法中的 `this.ensureContainer()`。
- **操作 2**: 在 `startProgress()` 方法开头添加 `this.ensureContainer()`。
- **原因**: 确保只有在触发通知或发起网络请求（触发进度条）时才注入 DOM。

```javascript
// src/components/notify/notify.store.js
init() {
  // 移除 this.ensureContainer()
},

startProgress() {
  this.ensureContainer() // 新增: 确保进度条容器存在
  // ... 原有逻辑
},
```

### 3.3 Tab 组件重构

**文件:** `src/components/tab/index.js`

- **操作**: 移除显式的 `Alpine.store('uno').tab.init()` 调用。

**文件:** `src/components/tab/tab.store.js`

- **操作 1**: 移除 `init()` 方法。
- **操作 2**: 新增 `initialized` 标志位和 `ensureInit()` 方法。
- **操作 3**: 将原 `init()` 的逻辑（绑定事件）移动到 `ensureInit()` 中，并确保只执行一次。
- **操作 4**: 在 `register()` (由 `x-uno-tab` 指令调用) 和 `activate()` 方法开头调用 `ensureInit()`。
- **原因**: 只有当页面存在 Tab 指令或通过 JS 激活 Tab 时，才开始监听 Hash 变化。

```javascript
// src/components/tab/tab.store.js
// 新增状态
initialized: false,

// 新增初始化方法
ensureInit() {
  if (this.initialized) return
  this.initialized = true

  // 原 init() 的逻辑：绑定 hashchange 等
  const checkHash = () => { ... }
  checkHash()
  window.addEventListener('hashchange', checkHash)
  window.addEventListener('load', checkHash)
},

register(name, group = 'default') {
  this.ensureInit() // 新增调用
  // ...
},

activate(instanceOpts = {}) {
  this.ensureInit() // 新增调用
  // ...
}
```

### 3.4 Tooltip 组件重构

**文件:** `src/components/tooltip/tooltip.store.js`

- **操作**: 移除 `init()` 方法中的 `this.ensureContainer()`。
- **原因**: 防御性编程，确保容器只由 `show()` 方法触发注入。

```javascript
// src/components/tooltip/tooltip.store.js
init() {
  // 移除 this.ensureContainer()
},
```

## 4. 验证策略

1.  **初始加载检查**:
    - 启动应用，检查 DOM。
    - 验证 **不应存在**: `#uno-box-manager`, `#uno-notify-container`, `#uno-global-progress`, `#uno-tooltip-manager`。
    - 验证 **不应存在**: Tab 组件的 `hashchange` 监听器（难以直接通过 DOM 验证，可通过 Console Log 验证 `ensureInit` 是否未执行）。

2.  **Box 功能验证**:
    - 执行 `uno.box.alert('test')` -> 验证容器注入，弹窗显示。

3.  **Notify 功能验证**:
    - 执行 `uno.notify.success('test')` -> 验证容器注入，Toast 显示。
    - 发起 HTMX 请求 -> 验证进度条容器注入，进度条动画显示。

4.  **Tab 功能验证**:
    - 访问带 Tab 的页面（存在 `x-uno-tab` 指令）。
    - 验证 `ensureInit` 被调用，Hash 同步功能正常（修改 URL hash，Tab 切换）。

5.  **Tooltip 功能验证**:
    - 触发 Tooltip（鼠标悬停或 `uno.tooltip.show()`）。
    - 验证容器注入，Tooltip 显示。

## 5. 收益

- **性能提升**: 减少页面初始 DOM 节点数量，降低浏览器解析和渲染开销。
- **资源优化**: 内存和事件监听器仅在功能被实际使用时占用。
- **代码整洁**: 逻辑更加模块化，副作用（DOM 注入、全局事件）与组件生命周期绑定更紧密。
