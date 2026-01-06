# Frequently Asked Questions (FAQ)

This document addresses common questions and issues when using Uno UI.

## Table of Contents

- [Installation & Setup](#installation--setup)
- [Components](#components)
- [TypeScript](#typescript)
- [HTMX Integration](#htmx-integration)
- [Styling & Customization](#styling--customization)
- [Performance](#performance)
- [Troubleshooting](#troubleshooting)

---

## Installation & Setup

### Q: How do I install Uno UI?

**A:** There are two ways:

1. **CDN (Quick Start)**:

```html
<script src="https://unpkg.com/alpinejs@3.15.2/dist/cdn.min.js" defer></script>
<script src="https://unpkg.com/uno-ui@0.6.0/dist/uno.umd.js"></script>
```

2. **NPM (Recommended)**:

```bash
npm install uno-ui alpinejs
```

See [README.md](../README.md) for detailed setup instructions.

### Q: Do I need to install Tailwind CSS separately?

**A:** Yes. Uno UI is designed to work with Tailwind CSS. You can:

- Use Tailwind CSS CDN for quick prototyping
- Install and configure Tailwind in your build process

```html
<!-- Tailwind CDN (Browser) -->
<script src="https://cdn.jsdelivr.net/npm/@tailwindcss/browser@4"></script>
```

### Q: What are the minimum browser requirements?

**A:** Uno UI supports:

- Chrome/Edge 90+
- Firefox 88+
- Safari 14+
- Opera 76+

These versions support ES6+ features and Alpine.js requirements.

### Q: Can I use Uno UI without a build tool?

**A:** Yes! Uno UI works perfectly via CDN without any build step:

```html
<!DOCTYPE html>
<html>
  <head>
    <script src="https://cdn.jsdelivr.net/npm/@tailwindcss/browser@4"></script>
    <script src="https://unpkg.com/alpinejs@3.15.2/dist/cdn.min.js" defer></script>
    <script src="https://unpkg.com/uno-ui@0.6.0/dist/uno.umd.js"></script>
  </head>
  <body>
    <button @click="$store.uno.notify.success('Hello!')">Click Me</button>
  </body>
</html>
```

---

## Components

### Q: How do I open a modal/dialog?

**A:** Use the Box component:

```javascript
// Simple dialog
Alpine.store('uno').box.open({
  title: 'Hello',
  content: 'This is a dialog'
})

// Confirmation dialog
const confirmed = await Alpine.store('uno').box.confirm('Are you sure?')
if (confirmed) {
  // User clicked OK
}
```

See [Box API](./box-api.md) for more options.

### Q: How do I show notifications?

**A:** Use the Notify component:

```javascript
// Success notification
Alpine.store('uno').notify.success('Operation completed!')

// Error notification
Alpine.store('uno').notify.error('Something went wrong')

// Custom notification
Alpine.store('uno').notify.show({
  type: 'info',
  text: 'Custom message',
  duration: 5000
})
```

See [Notify API](./notify-api.md) for more options.

### Q: How do I create tooltips?

**A:** Use HTML attributes or JavaScript:

**HTML Directive**:

```html
<button x-tooltip="This is a tooltip">Hover me</button>

<!-- With options -->
<button x-tooltip.top="Tooltip on top">Hover me</button>
```

**JavaScript**:

```javascript
Alpine.store('uno').tooltip.show({
  triggerEl: document.querySelector('#my-button'),
  text: 'Tooltip text'
})
```

See [Tooltip API](./tooltip-api.md) for more options.

### Q: How do I validate forms?

**A:** Use the Valid component:

```html
<input type="email" x-valid:email x-model="email" />

<!-- With custom validation -->
<input type="text" x-valid:custom="value.length >= 5" x-model="username" />
```

See [Valid API](./valid-api.md) for more patterns.

### Q: Can I use multiple modals at once?

**A:** Yes, Uno UI supports stacking multiple modals:

```javascript
// First modal
Alpine.store('uno').box.open({ title: 'First' })

// Second modal (opens on top)
Alpine.store('uno').box.open({ title: 'Second' })
```

Each modal gets the correct z-index automatically.

---

## TypeScript

### Q: Does Uno UI support TypeScript?

**A:** Yes! Full TypeScript support with 1,097 lines of type definitions.

```typescript
import type { BoxOptions, NotifyOptions } from 'uno-ui'

const options: BoxOptions = {
  title: 'Hello',
  position: 'center' // ✅ Type-safe
}
```

See [TypeScript Usage Guide](./typescript-usage.md) for details.

### Q: How do I get TypeScript autocomplete?

**A:** Just import the types - no additional setup needed:

```typescript
import { Uno } from 'uno-ui'
import type { UnoInstance, BoxOptions } from 'uno-ui'

// IntelliSense works automatically
const box: BoxOptions = {
  // Press Ctrl+Space to see all options
}
```

### Q: Why am I getting TypeScript errors?

**A:** Common causes:

1. **Missing type imports**:

```typescript
// ❌ Wrong
import { BoxOptions } from 'uno-ui'

// ✅ Correct
import type { BoxOptions } from 'uno-ui'
```

2. **Incorrect option values**:

```typescript
// ❌ Wrong - 'top' is not a valid BoxPosition
position: 'top'

// ✅ Correct
position: 'center' // or 'top-center'
```

3. **Missing Alpine store types**:

```typescript
// Add Alpine types to your project
npm install -D @types/alpinejs
```

---

## HTMX Integration

### Q: How do I use Uno UI with HTMX?

**A:** Uno UI has native HTMX support:

```html
<!-- Box with HTMX loading -->
<button x-box data-box-ajax-url="/api/content" data-box-title="Dynamic Content">
  Load Content
</button>

<!-- Tab with HTMX -->
<div x-data="$store.uno.tab.create('myTabs')">
  <button x-tab:trigger="'tab1'" hx-get="/tab1" hx-target="#content">Tab 1</button>
  <div id="content" x-tab:panel="'tab1'"></div>
</div>
```

See [Tab API](./tab-api.md) for HTMX integration examples.

### Q: How do I fetch HTML fragments for modals?

**A:** Use `Uno.htmx.fragment()`:

```javascript
const html = await Uno.htmx.fragment('GET', '/api/modal-content')

Alpine.store('uno').box.open({
  title: 'Dynamic Content',
  content: html
})
```

### Q: Can I use Uno UI without HTMX?

**A:** Yes! HTMX is optional. All components work without it:

```javascript
// Pure Alpine.js + Uno
Alpine.store('uno').box.open({
  title: 'Hello',
  content: 'No HTMX needed'
})
```

---

## Styling & Customization

### Q: How do I customize component styles?

**A:** Use the `ui` configuration object:

```javascript
Uno.init({
  ui: {
    box: {
      okBtn: 'bg-green-600 text-white px-4 py-2 rounded',
      cancelBtn: 'bg-gray-200 text-gray-800 px-4 py-2 rounded'
    }
  }
})
```

See [Configuration API](./config-api.md) for all options.

### Q: Can I use different themes?

**A:** Yes, configure styles globally:

```javascript
// Dark theme
Uno.init({
  ui: {
    base: 'rounded-lg shadow-lg border border-gray-700',
    box: {
      container: 'bg-gray-900 text-white'
    },
    notify: {
      item: 'bg-gray-800 text-white'
    }
  }
})
```

### Q: How do I remove rounded corners?

**A:**

```javascript
Uno.init({
  ui: {
    base: 'rounded-none shadow-md border'
  }
})
```

### Q: Can I add custom CSS classes to components?

**A:** Yes, use the `+` prefix to append classes:

```javascript
Uno.init({
  ui: {
    box: {
      okBtn: '+my-custom-button-class'
    }
  }
})
```

---

## Performance

### Q: What is the bundle size?

**A:**

- **ESM**: 69.69 kB (gzip: 19.32 kB)
- **UMD**: 49.87 kB (gzip: 15.36 kB)

Very lightweight compared to similar libraries.

### Q: Does Uno UI support tree-shaking?

**A:** Yes, the ESM build supports tree-shaking:

```javascript
// Only imports what you use
import { Uno } from 'uno-ui'
```

### Q: How can I optimize performance?

**A:** Best practices:

1. **Use CDN for production**:

```html
<script src="https://unpkg.com/uno-ui@0.6.0/dist/uno.umd.js"></script>
```

2. **Lazy load HTMX** (if not needed immediately):

```javascript
if (needsHtmx) {
  await import('htmx.org')
}
```

3. **Limit notification duration**:

```javascript
Alpine.store('uno').notify.success('Quick message', { duration: 2000 })
```

### Q: Are there any performance benchmarks?

**A:** Yes, key metrics:

- Tooltip show/hide: ~0.19ms
- Notify creation: ~0.15ms
- Box open: ~0.25ms
- Config lookup: ~0.05ms

See `refactor/reviews/performance-baseline.md` for details.

---

## Troubleshooting

### Q: Components don't work - what should I check?

**A:** Common issues:

1. **Check initialization order**:

```javascript
// ❌ Wrong order
Alpine.start()
Uno.init()

// ✅ Correct order
Uno.init()
Alpine.start()
```

2. **Check Alpine is loaded**:

```javascript
console.log(window.Alpine) // Should not be undefined
```

3. **Check Uno store**:

```javascript
console.log(Alpine.store('uno')) // Should show all components
```

4. **Enable debug mode**:

```javascript
Uno.init({ debug: true })
```

### Q: Styles are not applied - why?

**A:** Check:

1. **Tailwind CSS is loaded**:

```html
<script src="https://cdn.jsdelivr.net/npm/@tailwindcss/browser@4"></script>
```

2. **CSP (Content Security Policy)**:

```javascript
Uno.init({
  csp: {
    nonce: 'your-nonce-value'
  }
})
```

3. **Custom styles override**:

```javascript
// Check if custom config is correct
console.log(Uno.options.ui)
```

### Q: TypeScript types are not working

**A:** Solutions:

1. **Reinstall types**:

```bash
npm install uno-ui --save
```

2. **Check tsconfig.json**:

```json
{
  "compilerOptions": {
    "moduleResolution": "node",
    "types": ["uno-ui"]
  }
}
```

3. **Restart IDE** (VS Code, WebStorm, etc.)

### Q: HTMX integration not working

**A:** Check:

1. **HTMX is loaded**:

```javascript
console.log(window.htmx) // Should be defined
```

2. **Load order**:

```html
<!-- 1. HTMX -->
<script src="https://unpkg.com/htmx.org@2.0.8"></script>
<!-- 2. Alpine -->
<script src="https://unpkg.com/alpinejs@3.15.2/dist/cdn.min.js" defer></script>
<!-- 3. Uno -->
<script src="https://unpkg.com/uno-ui@0.6.0/dist/uno.umd.js"></script>
```

3. **HTMX events**:

```javascript
document.addEventListener('htmx:responseError', (e) => {
  console.error('HTMX error:', e.detail)
})
```

### Q: How do I report bugs?

**A:**

1. **Check existing issues**: [GitHub Issues](https://github.com/idiary-ai/uno-ui/issues)
2. **Create minimal reproduction**: Use CodePen or JSFiddle
3. **Provide details**:
   - Uno UI version
   - Alpine.js version
   - Browser and OS
   - Error messages
   - Code snippet

See [CONTRIBUTING.md](../CONTRIBUTING.md) for detailed guidelines.

### Q: Where can I get help?

**A:**

- **Documentation**: Read all docs in `doc/` directory
- **Error Handling**: See [Error Handling Guide](./error-handling.md)
- **Examples**: Check `examples/` directory
- **GitHub Issues**: Ask questions or report bugs
- **Community**: Join discussions (if available)

---

## Additional Resources

- [README](../README.md) - Getting Started
- [Error Handling Guide](./error-handling.md) - Error management
- [TypeScript Usage](./typescript-usage.md) - TypeScript integration
- [Configuration API](./config-api.md) - Configuration options
- [CONTRIBUTING](../CONTRIBUTING.md) - Contribution guidelines

---

**Last Updated**: 2026-01-06  
**Version**: 0.6.0
