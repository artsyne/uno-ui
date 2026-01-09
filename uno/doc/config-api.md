# Uno Configuration API

Uno provides a flexible 4-level configuration system that merges settings from plugin defaults, global user config, global plugin config, and instance-specific options.

## Configuration Levels

The configuration is resolved in the following order (from lowest to highest priority):

1.  **Plugin Defaults**: Built-in defaults defined by the component.
2.  **Global User Config**: Settings provided via `Uno.setConfig()` in `plugins.<componentName>`.
3.  **Instance Config**: Options passed directly to the component method (e.g., `Alpine.store('uno').box.alert(msg, options)` or `x-uno-box='{...}'`).

### Example

```javascript
// 1. Global Config (Level 2)
import Uno from 'uno-ui'

Uno.setConfig({
  // Global defaults for all notify instances
  plugins: {
    notify: {
      duration: 5000,
      position: 'top-right'
    }
  }
})

// 2. Instance Call (Level 3)
// This overrides the global 'duration' but keeps 'position'
Alpine.store('uno').notify.show('Hello', { duration: 1000 })
```

## Initialization

**Simplified**: Uno v0.6.0+ uses `setConfig()` for all configuration. No more `window.UnoConfig`!

### Automatic Initialization

Uno automatically initializes when Alpine is ready. Just include the scripts:

#### CDN Usage (No Config)

```html
<script defer src="https://unpkg.com/alpinejs@3.15.2/dist/cdn.min.js"></script>
<script defer src="https://unpkg.com/uno-ui@0.6.0/dist/uno.umd.js"></script>
```

**That's it!** Uno will:

1. Wait for the `alpine:init` event
2. Initialize all components automatically
3. Use default configuration

#### CDN Usage (With Config)

For CDN users who need custom configuration:

```html
<head>
  <!-- Configure Uno using ES module import -->
  <script type="module">
    import Uno from 'https://unpkg.com/uno-ui@0.6.0/dist/uno.es.js'

    Uno.setConfig({
      debug: true,
      onError: (err) => console.error('Uno Error:', err),
      ui: {
        base: 'rounded-lg shadow-md'
        // Your customizations
      }
    })
  </script>

  <!-- Load Alpine (Uno auto-initializes) -->
  <script defer src="https://unpkg.com/alpinejs@3.15.2/dist/cdn.min.js"></script>
</head>
```

#### Module Usage

For NPM/bundler usage:

```javascript
import Uno from 'uno-ui'
import Alpine from 'alpinejs'

// Configure Uno BEFORE Alpine starts
Uno.setConfig({
  debug: true,
  onError: (err) => console.error(err),
  ui: {
    // Your customizations
  }
})

// Register Alpine components
Alpine.data('myComponent', () => ({
  /* ... */
}))

// Start Alpine (Uno auto-initializes)
window.Alpine = Alpine
Alpine.start()
```

## Configuration Options

| Option    | Type       | Default           | Description                                                  |
| --------- | ---------- | ----------------- | ------------------------------------------------------------ |
| `debug`   | `boolean`  | `false`           | Enable/Disable internal logs. Set to `false` for production. |
| `onError` | `function` | `null`            | Custom error handler callback. Receives the error object.    |
| `csp`     | `object`   | `{ nonce: null }` | Content Security Policy settings (e.g., `{ nonce: '...' }`). |
| `plugins` | `object`   | `{}`              | Recommended namespace for component behavior options.        |
| `ui`      | `object`   | `{...}`           | Style configuration for `Uno.getStyle(comp, key, fallback)`. |

## Debug Mode

When `debug` is `true`, Uno prints lifecycle events and internal status to the console with `[Uno Info]` or `[Uno Warn]` prefixes.

To disable logs in production:

```javascript
Uno.setConfig({
  debug: false
})
```

## Component Options Namespace

Component behavior options are recommended to be placed under `plugins.<component>`.

```javascript
import Uno from 'uno-ui'

Uno.setConfig({
  plugins: {
    valid: {
      callbackDebounce: 350
    },
    box: {
      closeOnEsc: true
    }
  }
})
```
