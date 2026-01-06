# TypeScript Usage Guide

Uno UI v0.6.0+ includes full TypeScript type definitions for enhanced developer experience. This guide shows you how to use Uno with TypeScript.

## Installation

```bash
npm install uno-ui alpinejs htmx.org
```

No additional type packages needed! Type definitions are included in the package.

## Basic Usage

### Import with Types

```typescript
import { Uno, type UnoOptions } from 'uno-ui'
import Alpine from 'alpinejs'

// Initialize with type-safe configuration
const config: UnoOptions = {
  debug: true,
  csp: {
    nonce: 'your-nonce-here'
  }
}

window.Alpine = Alpine
Uno.init(config)
Alpine.start()
```

### Using Box Component

```typescript
import type { BoxOptions, BoxItem } from 'uno-ui'

// Type-safe box options
const options: BoxOptions = {
  title: 'Confirmation',
  content: 'Are you sure?',
  position: 'center',
  showOkButton: true,
  showCancelButton: true,
  onOk: () => {
    console.log('Confirmed')
  },
  onCancel: () => {
    console.log('Cancelled')
  }
}

// Get box item with full type information
const box: BoxItem = Alpine.store('uno').box.open(options)

// Promise-based confirm dialog
Alpine.store('uno')
  .box.confirm('Delete this item?', {
    title: 'Confirm Delete'
  })
  .then((confirmed: boolean) => {
    if (confirmed) {
      console.log('Item deleted')
    }
  })
```

### Using Notify Component

```typescript
import type { NotifyOptions, NotifyType } from 'uno-ui'

// Type-safe notify options
const options: NotifyOptions = {
  type: 'success',
  position: 'top-right',
  duration: 3000,
  icon: 'fas fa-check'
}

// Show notification
const id: number = Alpine.store('uno').notify.show('Success!', options)

// Type-safe shortcuts
Alpine.store('uno').notify.success('Operation completed')
Alpine.store('uno').notify.error('Something went wrong')
Alpine.store('uno').notify.warning('Please be careful')
Alpine.store('uno').notify.info('FYI: Update available')
```

### Using Tooltip Component

```typescript
import type { TooltipOptions, TooltipPlacement } from 'uno-ui'

// Type-safe placement
const placement: TooltipPlacement = 'top-start'

const options: TooltipOptions = {
  content: 'This is a tooltip',
  placement: placement,
  trigger: 'hover',
  arrow: true,
  offset: 10
}

Alpine.store('uno').tooltip.show(options)
```

### Using Tab Component

```typescript
import type { TabOptions } from 'uno-ui'

// Type-safe tab activation
const options: TabOptions = {
  name: 'profile',
  group: 'settings',
  hash: true
}

Alpine.store('uno').tab.activate(options)

// Check if tab is active
const isActive: boolean = Alpine.store('uno').tab.isActive('profile', 'settings')
```

## Advanced Configuration

### Global UI Configuration

```typescript
import type { UnoOptions, UiConfig, BoxUiConfig } from 'uno-ui'

const boxUi: BoxUiConfig = {
  container: 'rounded-xl bg-white dark:bg-gray-900 shadow-2xl',
  okBtn: 'bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded',
  cancelBtn: 'bg-gray-200 hover:bg-gray-300 text-gray-800 px-4 py-2 rounded'
}

const uiConfig: UiConfig = {
  base: 'rounded-lg shadow-md',
  layers: {
    box: 1400,
    tooltip: 1700,
    notify: 1600
  },
  box: boxUi
}

const config: UnoOptions = {
  debug: false,
  ui: uiConfig
}

Uno.init(config)
```

### Plugin Configuration

```typescript
import type { BoxPluginDefaults, NotifyPluginDefaults } from 'uno-ui'

const boxDefaults: BoxPluginDefaults = {
  defaultWidth: 500,
  defaultPosition: 'center',
  okText: 'Confirm',
  cancelText: 'Cancel',
  closeOnEsc: true,
  closeOnOverlayClick: false
}

const notifyDefaults: NotifyPluginDefaults = {
  duration: 5000,
  position: 'top-right'
}

const config: UnoOptions = {
  plugins: {
    box: boxDefaults,
    notify: notifyDefaults
  }
}

Uno.init(config)
```

## HTMX Integration with Types

```typescript
import type { BoxOptions } from 'uno-ui'

// HTMX content loading with type safety
const htmxOptions: BoxOptions = {
  title: 'User Profile',
  hxUrl: '/api/user/profile',
  hxMethod: 'GET',
  hxHeaders: {
    'X-CSRF-Token': getCsrfToken()
  },
  onContentLoad: (el: HTMLElement) => {
    console.log('Content loaded:', el)
  }
}

Alpine.store('uno').box.open(htmxOptions)
```

## Custom Types

### Extending Store Types

If you need to extend the store types:

```typescript
import 'uno-ui'

declare module 'alpinejs' {
  interface Stores {
    // Uno stores are already defined
    // You can add your custom stores here
    myCustomStore: {
      myMethod(): void
    }
  }
}
```

### Using Generic Types

```typescript
import type { BoxOptions, NotifyOptions } from 'uno-ui'

// Create a type-safe wrapper function
function showDialog<T extends BoxOptions>(options: T): Promise<boolean> {
  return Alpine.store('uno').box.confirm('Proceed?', options)
}

// Use with full type checking
showDialog({
  title: 'Confirm Action',
  position: 'center',
  closeOnEsc: true
})
```

## Type Utilities

### Position Types

```typescript
import type { BoxPosition, NotifyPosition, TooltipPlacement } from 'uno-ui'

// Box positions (9-point grid)
const boxPos: BoxPosition = 'top-center' // ✅ Valid
// const invalid: BoxPosition = 'invalid' // ❌ Type error

// Notify positions
const notifyPos: NotifyPosition = 'bottom-right' // ✅ Valid

// Tooltip placements (12 directions)
const tooltipPlace: TooltipPlacement = 'top-start' // ✅ Valid
```

### Type Guards

```typescript
import type { BoxType, NotifyType } from 'uno-ui'

function isBoxType(type: string): type is BoxType {
  return type === 'dialog' || type === 'confirm'
}

function isNotifyType(type: string): type is NotifyType {
  return ['success', 'error', 'warning', 'info'].includes(type)
}
```

## IntelliSense Benefits

With TypeScript support, you get:

- ✅ **Autocomplete**: Full autocomplete for all options and methods
- ✅ **Type Checking**: Catch errors at compile time
- ✅ **Documentation**: Inline documentation in your IDE
- ✅ **Refactoring**: Safe refactoring with type guarantees
- ✅ **Error Prevention**: Prevent typos and invalid configurations

## Migration from JavaScript

If you're migrating from JavaScript:

1. **Rename files**: `.js` → `.ts`
2. **Add type imports**: `import type { ... } from 'uno-ui'`
3. **Add type annotations**: Start with function parameters and return types
4. **Fix type errors**: TypeScript will guide you

Example migration:

```javascript
// Before (JavaScript)
const options = {
  title: 'Hello',
  content: 'World',
  position: 'center'
}
Alpine.store('uno').box.open(options)
```

```typescript
// After (TypeScript)
import type { BoxOptions } from 'uno-ui'

const options: BoxOptions = {
  title: 'Hello',
  content: 'World',
  position: 'center'
}
Alpine.store('uno').box.open(options)
```

## Common Issues

### Alpine Not Found

If you see `Cannot find module 'alpinejs'`:

```bash
npm install --save-dev @types/alpinejs
```

### Store Types Not Working

Make sure to import the types:

```typescript
import 'uno-ui' // This declares the store types
import Alpine from 'alpinejs'

// Now stores are typed
Alpine.store('uno').box.open({ title: 'Test' })
```

## Resources

- [TypeScript Documentation](https://www.typescriptlang.org/docs/)
- [Alpine.js with TypeScript](https://alpinejs.dev/advanced/typescript)
- [Uno UI API Reference](./box-api.md)

## Support

For TypeScript-related issues, please check:

1. Ensure you're using TypeScript 4.5+
2. Check your `tsconfig.json` includes type definitions
3. Verify type imports are correct
4. Report issues on GitHub with type examples
