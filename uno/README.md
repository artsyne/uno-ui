# Uno UI

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

```html
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <title>Uno App</title>

    <!-- 1. Styles (Tailwind CSS 4.x) -->
    <script src="https://cdn.jsdelivr.net/npm/@tailwindcss/browser@4"></script>

    <!-- 2. HTMX (Optional) -->
    <script src="https://unpkg.com/htmx.org@2.0.8"></script>

    <!-- 3. Alpine.js (Required) -->
    <script src="https://unpkg.com/alpinejs@3.15.2/dist/cdn.min.js" defer></script>

    <!-- 4. Uno Framework -->
    <script src="https://unpkg.com/uno-ui@0.1.0/dist/uno.umd.js"></script>
  </head>
  <body>
    <!-- Component Usage -->
    <button data-box-type="confirm" data-box-title="Hello World">Click Me</button>
  </body>
</html>
```

### 2. Using NPM (Module)

```bash
npm install uno-ui alpinejs htmx.org
```

Import into your project's entry file:

```javascript
import Alpine from 'alpinejs'
import 'uno-ui' // Registers Uno plugins automatically

// Initialize Alpine
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

const config: UnoOptions = {
  debug: true,
  plugins: {
    box: {
      defaultWidth: 500,
      closeOnEsc: true
    }
  }
}

Uno.init(config)
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

## TypeScript Support

Uno UI v0.6.0+ includes full TypeScript type definitions. See the [TypeScript Usage Guide](doc/typescript-usage.md) for detailed information.

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

Uno UI provides a flexible styling system that allows you to customize the look and feel of components globally. You can override default Tailwind CSS classes by setting the `window.UnoConfig.ui` object before initialization.

### Configuration Structure

The `ui` configuration object is structured as follows:

```javascript
window.UnoConfig = {
  ui: {
    // Base styles used across multiple components
    base: {
      rounded: 'rounded-lg', // Default: rounded-lg. Set to 'rounded-none' for square look.
      primary: 'bg-blue-600', // Default: bg-blue-600.
      shadow: 'shadow-md' // Default: shadow-md.
    },
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
      triggerActive: 'border-b-2 border-blue-500',
      triggerInactive: 'text-gray-500'
    },
    valid: {
      inputError: '!border-red-500',
      inputValid: '!border-green-500',
      helperText: 'text-red-600 text-xs'
    }
  }
}
```

### Examples

#### 1. Flat & Square Theme (No Rounded Corners)

```javascript
window.UnoConfig = {
  ui: {
    base: {
      rounded: 'rounded-none',
      shadow: 'shadow-none'
    }
  }
}
```

#### 2. Brand Color Theme

```javascript
window.UnoConfig = {
  ui: {
    base: {
      primary: 'bg-indigo-600 hover:bg-indigo-700'
    },
    tab: {
      triggerActive: 'text-indigo-600 border-b-2 border-indigo-600 font-bold'
    }
  }
}
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

The project includes a manual test suite located in the `examples/` directory. After building, you can verify component functionality by opening `examples/index.html`.

```bash
npm run build
open examples/index.html
```

## Project Structure

```
uno/
├── dist/               # Build artifacts (uno.umd.js)
├── doc/                # API Documentation
├── examples/           # Manual test cases
│   ├── components/     # Individual component examples
│   └── index.html      # Test dashboard
├── src/
│   ├── components/     # Component source code
│   │   ├── box/        # Box component (store, directive, template)
│   │   ├── notify/     # Notify component
│   │   ├── tab/        # Tab component
│   │   ├── tooltip/    # Tooltip component
│   │   └── valid/      # Valid component
│   ├── core/           # Core utilities (i18n, positioning, etc.)
│   └── index.js        # Entry point
├── package.json
└── vite.config.js
```

## License

MIT License
