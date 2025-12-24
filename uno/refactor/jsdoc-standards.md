# JSDoc 注释标准

## 📋 概述

Uno UI v2.0 所有代码注释必须采用 JSDoc 标准，统一使用英文注释，遵循开源项目最佳实践。

---

## 🎯 核心原则

### 1. **语言要求**
- **必须使用英文注释**：所有注释、文档、说明都必须使用英文
- **遵循 JSDoc 3.x 标准**：使用最新的 JSDoc 语法
- **保持简洁明了**：注释应该解释"为什么"而不是"做什么"

### 2. **注释原则**
- **注释解释意图**：解释代码的目的和设计决策
- **避免显而易见**：不要注释明显的代码
- **保持更新**：代码变更时必须同步更新注释

---

## 📝 JSDoc 注释标准

### 1. 函数注释

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
 * @example
 * // Open a confirm dialog with callbacks
 * const dialogId = Alpine.store('box').open({
 *   type: 'confirm',
 *   title: 'Confirm Action',
 *   content: 'Are you sure you want to delete this item?',
 *   onOk: () => deleteItem(),
 *   onCancel: () => console.log('Cancelled')
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
 * @example
 * try {
 *   const content = await loadContent('/api/user-info');
 *   console.log('Content loaded:', content);
 * } catch (error) {
 *   console.error('Failed to load content:', error.message);
 * }
 * 
 * @since 2.0.0
 */
async function loadContent(url, options = {}) {
  // Implementation...
}
```

### 2. 类注释

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

#### 类方法注释
```javascript
/**
 * Validates a form field using the specified validation rules.
 * 
 * @memberof BoxStore
 * @param {HTMLInputElement} field - The input element to validate
 * @param {string} rule - The validation rule to apply
 * @param {*} [value] - The value to validate (defaults to field value)
 * @returns {boolean} True if the field is valid, false otherwise
 * 
 * @example
 * const isValid = this.validateField(emailInput, 'email');
 * if (!isValid) {
 *   this.showError(emailInput, 'Please enter a valid email');
 * }
 * 
 * @since 2.0.0
 */
validateField(field, rule, value = field.value) {
  // Implementation...
}
```

### 3. 常量注释

#### 常量定义注释
```javascript
/**
 * Default z-index values for different UI layers.
 * 
 * @readonly
 * @enum {number}
 * @since 2.0.0
 */
const Z_INDEX = {
  /** Base layer for components */
  BASE: 1000,
  /** Modal dialogs and overlays */
  MODAL: 10000,
  /** Tooltips and popovers */
  TOOLTIP: 20000,
  /** Notifications and alerts */
  NOTIFY: 30000
};
```

### 4. 事件注释

#### 事件处理器注释
```javascript
/**
 * Handles the click event on dialog triggers.
 * 
 * @param {MouseEvent} event - The click event object
 * @param {HTMLElement} target - The clicked element
 * @returns {void}
 * 
 * @listens click
 * @since 2.0.0
 */
const handleDialogClick = (event, target) => {
  // Implementation...
};
```

### 5. 配置对象注释

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

---

## 📏 注释最佳实践

### ✅ 推荐的注释

#### 1. 解释复杂逻辑
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

#### 2. 记录设计决策
```javascript
/**
 * Uses requestAnimationFrame for smooth animations instead of CSS transitions.
 * 
 * This approach provides better performance for frequently updating animations
 * and allows for more precise timing control. It also avoids layout thrashing
 * that can occur with CSS transitions during rapid updates.
 * 
 * @param {HTMLElement} element - The element to animate
 * @param {Object} fromStyles - Starting styles
 * @param {Object} toStyles - Ending styles
 * @param {number} [duration=300] - Animation duration in milliseconds
 */
function animateElement(element, fromStyles, toStyles, duration = 300) {
  // Animation implementation...
}
```

#### 3. API 文档
```javascript
/**
 * Creates a new notification with the specified type and message.
 * 
 * This method supports all notification types (success, error, warning, info)
 * and automatically handles positioning and auto-dismissal based on the
 * global configuration.
 * 
 * @param {string} message - The message to display
 * @param {Object} [options={}] - Additional options for the notification
 * @param {string} [options.type='info'] - The notification type
 * @param {number} [options.duration] - Auto-dismiss duration in milliseconds
 * @param {string} [options.position] - Notification position
 * @returns {string} The unique ID of the created notification
 * 
 * @example
 * // Create a success notification
 * const id = Alpine.store('notify').success('Operation completed!');
 * 
 * @example
 * // Create a persistent error notification
 * const id = Alpine.store('notify').error('Failed to save data', {
 *   duration: 0 // Don't auto-dismiss
 * });
 * 
 * @since 2.0.0
 */
```

### ❌ 避免的注释

#### 1. 注释显而易见的代码
```javascript
// ❌ 避免这样的注释
const x = 5; // Set x to 5

// ✅ 更好的做法（如果需要注释）
const x = 5; // Initial counter value
```

#### 2. 过度注释
```javascript
// ❌ 避免这样的注释
// Function to open a dialog
function openDialog(config) {
  // Create a new dialog
  const dialog = createDialogElement();
  
  // Set the title
  dialog.title = config.title;
  
  // Set the content
  dialog.content = config.content;
  
  // Show the dialog
  showDialog(dialog);
}

// ✅ 更好的做法
function openDialog(config) {
  const dialog = createDialogElement();
  Object.assign(dialog, config);
  showDialog(dialog);
}
```

#### 3. 过时的注释
```javascript
// ❌ 避免这样的注释
// TODO: Remove this legacy code in v2.1
function legacyFunction() {
  // This should be replaced with the new implementation
}

// ✅ 更好的做法
/**
 * @deprecated Use openDialog() instead. Will be removed in v2.1.
 * @since 1.0.0
 * @deprecated since 2.0.0
 */
function legacyFunction() {
  // Implementation...
}
```

---

## 🔧 工具和验证

### JSDoc 验证工具
```bash
# 安装 JSDoc 验证工具
npm install -D jsdoc

# 验证注释语法
jsdoc src/components/box.js

# 生成文档
jsdoc -c jsdoc.conf.json src/
```

### 推荐的 VS Code 扩展
- **JSDoc Comments** - 提供 JSDoc 语法高亮和自动补全
- **Document This** - 自动生成 JSDoc 注释模板
- **Better Comments** - 改善注释的可读性

---

## 📚 注释检查清单

在提交代码前，请确保：

- [ ] 所有公共函数都有 JSDoc 注释
- [ ] 所有类都有完整的类文档
- [ ] 所有复杂逻辑都有解释性注释
- [ ] 所有注释都使用英文
- [ ] 注释与代码保持同步
- [ ] 没有 TODO 注释（除非有明确的计划）
- [ ] 没有过时的注释
- [ ] 所有示例代码都是可执行的

---

*遵循这些标准将确保 Uno UI 代码库的文档质量和可维护性。*