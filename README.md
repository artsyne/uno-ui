# Uno UI

[![npm version](https://img.shields.io/npm/v/uno-ui.svg?style=flat-square)](https://www.npmjs.com/package/uno-ui)
[![npm downloads](https://img.shields.io/npm/dm/uno-ui.svg?style=flat-square)](https://www.npmjs.com/package/uno-ui)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg?style=flat-square)](https://opensource.org/licenses/MIT)
[![TypeScript](https://img.shields.io/badge/TypeScript-Ready-blue.svg?style=flat-square)](https://www.typescriptlang.org/)
[![Bundle Size](https://img.shields.io/bundlephobia/minzip/uno-ui?style=flat-square)](https://bundlephobia.com/package/uno-ui)

Uno is a lightweight, modern UI component library built on **Alpine.js** and **Tailwind CSS**, designed for simplicity, performance, and ease of integration.

It provides a curated set of essential UI components that require no complex build steps to use, while also supporting modern ESM-based development workflows.

## Features

- ⚡ **Lightweight**: Built on a micro-kernel architecture relying only on Alpine.js, ensuring minimal runtime overhead.
- 🎨 **Modern Stack**: Deeply integrated with Tailwind CSS 4.0, offering flexible and customizable styling.
- 📦 **Zero Build Config**: Can be used directly via CDN without any build tools, or imported as an ESM module.
- 🔌 **HTMX Ready**: Natively supports HTMX's server-driven UI patterns, with components designed for seamless AJAX interactions.
- 🌍 **I18n Support**: Built-in internationalization mechanism for easy multi-language switching.
- ♿ **Accessible**: Follows WAI-ARIA best practices to ensure good accessibility.
- 📘 **TypeScript Support**: Full TypeScript type definitions included for enhanced developer experience.

## Getting Started

### 1. Using CDN (Fastest)

Simply include the compiled CSS and JS files in your HTML.

**Dependencies**: Alpine.js (Required), HTMX (Optional, Recommended).

**Best Practice**: Use `defer` for all scripts to ensure non-blocking loading and correct execution order.

```html
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <title>Uno App</title>

    <!-- 1. Styles (Tailwind CSS 4.x) -->
    <script src="https://cdn.jsdelivr.net/npm/@tailwindcss/browser@4"></script>

    <!-- 2. HTMX (Optional) -->
    <script src="https://unpkg.com/htmx.org@2.0.10"></script>

    <!-- 3. Uno Framework (Must load before your app code if you use global Uno) -->
    <script defer src="https://unpkg.com/uno-ui@0.9.0/dist/uno.umd.js"></script>

    <!-- 4. Your App Code (Defensive Initialization) -->
    <script defer src="/js/app.js"></script>

    <!-- 5. Alpine.js (Base Framework) -->
    <script defer src="https://unpkg.com/alpinejs@3.17.3/dist/cdn.min.js"></script>
  </head>
  <body>
    <!-- Component Usage -->
    <button x-data @click="$store.uno.box.alert('Hello World')">Click Me</button>
  </body>
</html>
```

#### Defensive Initialization (Recommended)

When using `defer` scripts, the exact execution timing can vary. To ensure your configuration is applied correctly, use this defensive pattern in your `app.js`:

```javascript
// app.js

const initApp = () => {
  // 1. Configure Uno
  if (window.Uno) {
    window.Uno.setConfig({
      debug: true,
      ui: { base: 'rounded-lg shadow-md' }
    })
  }

  // 2. Register Alpine Components
  Alpine.data('myComponent', () => ({
    init() {
      console.log('Initialized')
    }
  }))
}

// Defensive Start: Run immediately if Alpine is ready, otherwise wait
if (window.Alpine) {
  initApp()
} else {
  document.addEventListener('alpine:init', initApp)
}
```

#### Optional: Custom Configuration (Module)

For NPM users, configure Uno using `setConfig()` BEFORE Alpine loads:

````html
<head>
  <!-- 1. Configure Uno (before Alpine loads) -->
  <script type="module">
    import Uno from 'https://unpkg.com/uno-ui@0.9.0/dist/uno.es.js'

    Uno.setConfig({
      debug: true,
      ui: {
        base: 'rounded-lg shadow-md',
        box: {
          container: 'custom-box-styles'
          // ... other customizations
        }
      }
    })
  </script>

  <!-- 2. Load Alpine (Uno auto-initializes) -->
  <script defer src="https://unpkg.com/alpinejs@3.17.3/dist/cdn.min.js"></script>
</head>

#### Optional: Register Components To register custom Alpine components: ```html
<script defer>
  document.addEventListener('alpine:init', () => {
    Alpine.data('myComponent', () => ({
      message: 'Hello World'
    }))
  })
</script>
````

### 2. Using NPM (Module)

```bash
npm install uno-ui alpinejs htmx.org
```

Import into your project's entry file:

```javascript
import Alpine from 'alpinejs'
import Uno from 'uno-ui'

// Configure Uno BEFORE Alpine loads
Uno.setConfig({
  debug: true,
  ui: {
    base: 'rounded-lg shadow-md'
    // Your custom styles
  }
})

// Register your Alpine components
Alpine.data('myComponent', () => ({
  // ...
}))

// Start Alpine (Uno auto-initializes)
window.Alpine = Alpine
Alpine.start()
```

### 3. TypeScript Support

Uno UI includes full TypeScript type definitions. No additional setup required!

```typescript
import { Uno, type BoxOptions, type NotifyOptions } from 'uno-ui'
import Alpine from 'alpinejs'

// Full type checking and IntelliSense
const options: BoxOptions = {
  title: 'Hello',
  content: 'World',
  position: 'center'
}

Alpine.store('uno').box.open(options)
```

**Type-safe configuration:**

```typescript
import type { UnoOptions } from 'uno-ui'
import Uno from 'uno-ui'

const config: UnoOptions = {
  debug: true,
  plugins: {
    box: {
      defaultWidth: 500,
      closeOnEsc: true
    }
  }
}

Uno.setConfig(config)
```

## Components

Uno provides a set of essential components. Detailed documentation is available in the `doc/` directory.

| Component   | Description                                                                            | Documentation                     |
| ----------- | -------------------------------------------------------------------------------------- | --------------------------------- |
| **Box**     | Modals, dialogs, slide-overs with flexible positioning and lazy loading (AJAX/HTMX).   | [Box API](doc/box-api.md)         |
| **Notify**  | Global notifications and toast alerts, supporting multiple positions and states.       | [Notify API](doc/notify-api.md)   |
| **Tab**     | Tab switching with support for AJAX/HTMX content loading and URL Hash synchronization. | [Tab API](doc/tab-api.md)         |
| **Tooltip** | Smart positioning tooltips supporting Hover/Click triggers.                            | [Tooltip API](doc/tooltip-api.md) |
| **Valid**   | Powerful form validation engine with declarative rules and async support.              | [Valid API](doc/valid-api.md)     |

### Security

Uno UI provides built-in support for Content Security Policy (CSP) to help you secure your applications against XSS attacks.

- [CSP Configuration Guide](doc/csp-guide.md) - Learn how to configure CSP with Uno UI

## TypeScript Support

Uno UI v0.7.0+ includes full TypeScript type definitions. See the [TypeScript Usage Guide](doc/typescript-usage.md) for detailed information.

**Key Features:**

- ✅ Full type definitions for all components
- ✅ IntelliSense and autocomplete in your IDE
- ✅ Compile-time type checking
- ✅ No additional setup required

**Quick Example:**

```typescript
import type { BoxOptions } from 'uno-ui'

const options: BoxOptions = {
  title: 'Hello',
  position: 'center' // ✅ Type-safe
  // position: 'invalid' // ❌ Type error
}

Alpine.store('uno').box.open(options)
```

## Customization

Uno UI provides a flexible styling system that allows you to customize the look and feel of components globally. You can override default Tailwind CSS classes via `Uno.setConfig()` — callable both before initialization (recommended) and at runtime.

### Configuration Structure

The `ui` configuration object is structured as follows:

```javascript
Uno.setConfig({
  ui: {
    // Base style string, substituted wherever a component style contains `{base}`
    base: 'rounded-lg shadow-md border border-gray-200 dark:border-gray-700',
    // Z-index layers managed by the core layer system
    layers: { box: 1300, tooltip: 1700, notify: 1600 },
    // Component-specific overrides
    box: {
      overlay: '...', // Override overlay style
      container: '...', // Override container style
      okBtn: '...', // Override OK button style
      cancelBtn: '...' // Override Cancel button style
    },
    notify: {
      types: {
        success: 'bg-green-100 ...',
        error: 'bg-red-100 ...'
      }
    },
    tooltip: {
      content: '...' // Override tooltip content style
    },
    tab: {
      active: 'border-b-2 border-blue-500',
      inactive: 'text-slate-500 hover:text-slate-700'
    },
    valid: {
      inputError: '!border-red-500',
      inputValid: '!border-green-500',
      helperText: 'text-red-600 text-xs'
    }
  }
})
```

> **Tip:** Prefix a style value with `+` to append it to the default instead of replacing it, e.g. `okBtn: '+min-w-24'`.

### Examples

#### 1. Flat & Square Theme (No Rounded Corners)

```javascript
Uno.setConfig({
  ui: {
    base: 'rounded-none shadow-none border border-gray-200 dark:border-gray-700'
  }
})
```

#### 2. Brand Color Theme

```javascript
Uno.setConfig({
  ui: {
    tab: {
      active: 'text-indigo-600 border-b-2 border-indigo-600 font-bold'
    },
    box: {
      okBtn:
        'inline-flex w-full justify-center rounded-md bg-indigo-600 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 sm:ml-3 sm:w-auto'
    }
  }
})
```

## Internationalization (i18n)

Uno components include built-in support for internationalization.

### Supported Languages

The library comes with the following built-in locales:

- `en` (English) - Default
- `zh-CN` (Simplified Chinese)
- `zh-TW` (Traditional Chinese)

### Basic Usage

You can switch languages globally using the `i18n` store. This should be done after Alpine initializes.

```javascript
document.addEventListener('alpine:init', () => {
  // Switch to Simplified Chinese
  Alpine.store('i18n').setLang('zh-CN')
})
```

### Extending Locales

You can add support for new languages or override existing translations using the `extend` method. This is useful for customizing component text or adding translations for your own application.

```javascript
document.addEventListener('alpine:init', () => {
  Alpine.store('i18n').extend({
    // Override existing English text
    en: {
      box: {
        ok: 'Yes',
        cancel: 'No'
      }
    },
    // Add a new language (French)
    'fr-FR': {
      box: {
        ok: 'Confirmer',
        cancel: 'Annuler',
        loading: 'Chargement...',
        close: 'Fermer'
      },
      valid: {
        required: 'Ce champ est requis',
        email: 'Veuillez entrer une adresse email valide'
      }
    }
  })

  // Switch to the new language
  Alpine.store('i18n').setLang('fr-FR')
})
```

## Development

### Prerequisites

- Node.js 18+
- NPM 9+

### Setup

```bash
git clone https://github.com/your-org/uno-ui.git
cd uno-ui
npm install
```

### Scripts

- `npm run dev`: Start the Vite development server.
- `npm run build`: Build for production (output to `dist/`).
- `npm run preview`: Preview the production build.

### Testing

The project uses Vitest for unit testing with comprehensive coverage of all components and utilities.

```bash
# Run unit tests once (this is the CI entry point)
npm run test:run

# Watch mode during development
npm run test -- --watch
```

**Coverage:**

- Unit tests for the components, the core modules and the utilities, plus component-level axe scans
- Behaviour and accessibility contracts live in `tests/`; counts change per release, so run the suite rather than reading a fixed number here
- See [Test Reports](doc/test-reports.md) for the layout and what the component scans do and do not cover

You can also verify component functionality manually using the interactive demo:

```bash
npm run build
npm run dev
# Open http://localhost:5176/examples/
```

## Project Structure

```
uno/
├── dist/                  # Build Output (ESM & UMD bundles)
├── types/                 # TypeScript Definitions
├── src/
│   ├── components/       # UI Components (box, notify, tooltip, tab, valid)
│   ├── core/             # Core System (uno, i18n, layer, plugin-registry)
│   ├── utils/            # Utilities (dom, position, validation, a11y, log)
│   ├── locales/          # I18n Language Files (en, zh-CN, zh-TW)
│   └── index.js          # Entry point
├── tests/                # Unit tests (components, core, utils) + axe helper
├── doc/                  # API Documentation
├── examples/             # Usage Examples & Interactive Demo
├── CONTRIBUTING.md       # Contribution Guidelines
├── CODE_OF_CONDUCT.md    # Community Code of Conduct
└── package.json          # Package Configuration
```

## License

MIT License

## HTMX version

Uno and Hebe use the same HTMX 2 version resolved by `uno/package-lock.json`. Hebe copies the runtime from Uno dependencies. Only XHR event details are supported; upgrades update both projects together.
