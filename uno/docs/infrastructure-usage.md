# Infrastructure Usage Guide

This guide provides examples and best practices for using the Uno infrastructure components.

## 1. Configuration

Uno supports a 4-level configuration merge strategy.

```javascript
// 1. Initialize with Global Config
Uno.init(Alpine, {
  debug: true,
  ui: {
    // Override base classes
    base: 'rounded-xl shadow-lg',
    // Component specific config
    notify: {
      position: 'top-right'
    }
  }
})

// 2. Component Usage with Instance Config
// <div x-data="notify({ timeout: 5000 })">
```

## 2. Internationalization (i18n)

Use the `x-t` directive for declarative translations.

```html
<!-- Simple Key -->
<span x-t="welcome.message"></span>

<!-- With Attribute Modifier -->
<input x-t.placeholder="forms.email_placeholder" type="email">

<!-- Dynamic Language Switching -->
<button @click="$store.i18n.setLang('zh-CN')">中文</button>
```

## 3. Z-Index Management

Avoid magic numbers by using the Z-Index manager.

```javascript
import { getZIndex } from 'uno/core/layer'

const modalZ = getZIndex('modal') // Returns 1100
const toastZ = getZIndex('notify') // Returns 3000
```

## 4. Accessibility (A11y)

Use the `trapFocus` utility for modals and dialogs.

```javascript
import { trapFocus } from 'uno/utils/a11y'

// Trap focus within an element
const cleanup = trapFocus(modalElement)

// Release focus when done
cleanup()
```

## 5. Error Handling

Uno automatically captures errors, but you can add custom handlers.

```javascript
Uno.init(Alpine, {
  onError: (error, context) => {
    // Send to Sentry or logging service
    console.error('Captured:', error)
  }
})
```

## 6. Performance

Use built-in utilities for expensive operations.

```javascript
import { debounce, throttle } from 'uno/utils/performance'

const onResize = throttle(() => {
  // Handle resize
}, 200)

const onSearch = debounce((query) => {
  // API call
}, 300)
```
