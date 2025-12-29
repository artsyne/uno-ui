# Box Component (Dialog)

A lightweight, declarative, HTMX-native modal component based on Alpine.js + TailwindCSS. Follows the "Store-First, Directive-Driven" architecture pattern.

## Architecture

**Store-First, Directive-Driven Pattern:**

- **Store**: All logic and state managed via `Alpine.store('uno').box`
- **Directive**: Declarative interface via `x-uno-box` directive
- **Namespace**: Registered under `uno` namespace to avoid conflicts

## Features

- **Modal & Overlay**: Supports `modal` overlay mode with click-to-close and ESC key support
- **9-Point Positioning**: `center`, `top-left`, `top-center`, `top-right`, `center-left`, `center-right`, `bottom-left`, `bottom-center`, `bottom-right`
- **Size Control**: Numeric width/height (px) or CSS classes via `widthClass`/`heightClass`
- **XSS Protection**: All content is automatically sanitized (trustedHTML removed for security)
- **HTMX Integration**: Native HTMX content loading via `hxUrl`
- **Accessibility (A11y)**:
  - Focus trap within modal
  - ARIA roles (`role="dialog"`, `aria-modal`, `aria-labelledby`)
  - Keyboard navigation (Escape to close)
- **Callbacks**: `onOk`, `onCancel`, `onBeforeClose`, `onClose`, `onContentLoad`
- **i18n Support**: Button text via Alpine i18n store

## Quick Start

### Programmatic API

```javascript
// Alert dialog
Uno.box.alert('Operation completed!')

// Confirm dialog (returns Promise)
const ok = await Uno.box.confirm('Are you sure?')

// Prompt dialog (returns Promise)
const name = await Uno.box.prompt('Enter your name:', 'John')

// Custom dialog
Uno.box.dialog({
  title: 'Notice',
  content: 'Hello World',
  position: 'top-center'
})
```

### Alpine Store API

```javascript
// Access via Alpine store
Alpine.store('uno').box.open({
  title: 'System Notice',
  content: 'Operation Successful',
  position: 'top-center',
  onClose: () => console.log('Closed')
})

// Confirm with promise
Alpine.store('uno')
  .box.confirm('Delete this item?', {
    title: 'Confirm Delete'
  })
  .then((result) => {
    if (result) console.log('Confirmed')
  })
```

### Directive Usage (x-uno-box)

```html
<!-- Simple dialog -->
<button x-uno-box='{"title": "Hello", "content": "World"}'>Open Dialog</button>

<!-- With HTMX content -->
<button x-uno-box='{"title": "Dynamic", "hxUrl": "/api/content"}'>Load Content</button>
```

## Configuration Options

### Options Object (for `open()`, `dialog()`, etc.)

| Option             | Type                    | Default     | Description                               |
| ------------------ | ----------------------- | ----------- | ----------------------------------------- |
| `type`             | `'dialog' \| 'confirm'` | `'dialog'`  | Dialog type                               |
| `title`            | `string`                | `''`        | Dialog title                              |
| `content`          | `string`                | `''`        | Dialog content (auto-sanitized)           |
| `position`         | `string`                | `'center'`  | One of 9 positions                        |
| `width`            | `number`                | `400`       | Width in px (ignored if `widthClass` set) |
| `height`           | `number`                | `undefined` | Height in px                              |
| `widthClass`       | `string`                | `undefined` | CSS class for width                       |
| `heightClass`      | `string`                | `undefined` | CSS class for height                      |
| `containerClass`   | `string`                | `undefined` | Additional container class                |
| `modal`            | `boolean`               | `true`      | Show overlay, enable ESC close            |
| `showCloseButton`  | `boolean`               | `true`      | Show close button                         |
| `showOkButton`     | `boolean`               | `false`     | Show OK button                            |
| `showCancelButton` | `boolean`               | `false`     | Show Cancel button                        |
| `okText`           | `string`                | `'OK'`      | OK button text                            |
| `cancelText`       | `string`                | `'Cancel'`  | Cancel button text                        |
| `boxId`            | `string`                | random      | Custom box instance ID                    |
| `hxUrl`            | `string`                | `undefined` | HTMX content URL                          |
| `hxMethod`         | `string`                | `'GET'`     | HTMX request method                       |
| `hxHeaders`        | `object`                | `undefined` | HTMX request headers                      |
| `hxTrigger`        | `string`                | `'load'`    | HTMX trigger condition                    |

### Callback Options

| Callback        | Signature                   | Description                                         |
| --------------- | --------------------------- | --------------------------------------------------- |
| `onOk`          | `() => void \| false`       | Called on OK click. Return `false` to prevent close |
| `onCancel`      | `() => void`                | Called on Cancel click                              |
| `onBeforeClose` | `() => void`                | Called before dialog closes                         |
| `onClose`       | `() => void`                | Called after dialog closes                          |
| `onContentLoad` | `(el: HTMLElement) => void` | Called after HTMX content loads                     |

## HTMX Integration

### Loading Content

```javascript
Uno.box.dialog({
  title: 'User Profile',
  hxUrl: '/api/user/profile',
  hxMethod: 'GET',
  hxHeaders: { 'X-CSRF-Token': 'token123' },
  onContentLoad: (el) => {
    console.log('Content loaded into:', el)
  }
})
```

### Server Response Format

Return plain HTML fragments:

```html
<div class="p-4">
  <h2>User Info</h2>
  <p>Name: John Doe</p>
  <!-- Nested HTMX is supported -->
  <button hx-post="/api/update" hx-target="#result">Update</button>
  <div id="result"></div>
</div>
```

## Security

### Automatic XSS Protection

All content passed via `content` option is automatically sanitized:

- Script tags removed
- Event handlers (`onclick`, `onerror`, etc.) removed
- `javascript:` URLs removed
- Image tags removed (to prevent `onerror` XSS)

**Note:** `trustedHTML` option has been removed for security. Use HTMX for dynamic content.

## Internationalization (i18n)

Button text is resolved from Alpine i18n store:

```javascript
// Set up translations
Alpine.store('i18n').setLang('zh-CN')

// Button text uses i18n keys:
// - box.ok -> "确定"
// - box.cancel -> "取消"
```

## Global Configuration

Configure defaults via `Uno.init()` or `window.UnoConfig`:

```javascript
// Via Uno.init
Uno.init({
  plugins: {
    box: {
      defaultWidth: 500,
      defaultPosition: 'top-center',
      okText: 'Confirm',
      cancelText: 'Dismiss'
    }
  }
})

// Or via UnoConfig
window.UnoConfig = {
  ui: {
    box: {
      container: 'bg-white dark:bg-gray-900 rounded-xl shadow-2xl',
      okBtn: 'bg-green-600 hover:bg-green-700 text-white px-4 py-2 rounded'
    }
  }
}
```

## Store Methods

| Method                              | Description                           |
| ----------------------------------- | ------------------------------------- |
| `open(type, options)`               | Open a dialog                         |
| `dialog(options)`                   | Alias for `open('dialog', options)`   |
| `alert(message, options)`           | Show alert dialog                     |
| `confirm(message, options)`         | Show confirm dialog (returns Promise) |
| `prompt(message, default, options)` | Show prompt dialog (returns Promise)  |
| `close(id)`                         | Close specific dialog by ID           |
| `closeAll()`                        | Close all dialogs                     |

## Migration from v0.5.x

### Breaking Changes

1. **Namespace Change**: Store moved from `Alpine.store('box')` to `Alpine.store('uno').box`

```javascript
// Old
Alpine.store('box').open({...})

// New
Alpine.store('uno').box.open({...})
// Or use global API
Uno.box.dialog({...})
```

2. **trustedHTML Removed**: For security, all content is auto-sanitized. Use HTMX for dynamic content:

```javascript
// Old (no longer supported)
store.open({ content: htmlString, trustedHTML: true })

// New - use HTMX
store.open({ hxUrl: '/api/content' })
```

3. **Global API**: Use `Uno.box` instead of accessing store directly for convenience

## Example

```javascript
// Full example with all features
Uno.box.dialog({
  title: 'Edit Profile',
  hxUrl: '/api/profile/edit',
  position: 'center',
  width: 600,
  showOkButton: true,
  showCancelButton: true,
  okText: 'Save',
  cancelText: 'Discard',
  onOk: () => {
    // Validate before close
    const form = document.querySelector('#profile-form')
    if (!form.checkValidity()) return false
    form.submit()
  },
  onBeforeClose: () => console.log('About to close'),
  onClose: () => console.log('Closed'),
  onContentLoad: (el) => {
    // Initialize form validation
    Alpine.initTree(el)
  }
})
```
