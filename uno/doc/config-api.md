# Uno Configuration API

Uno provides a flexible 4-level configuration system that merges settings from plugin defaults, global user config, global plugin config, and instance-specific options.

## Configuration Levels

The configuration is resolved in the following order (from lowest to highest priority):

1.  **Plugin Defaults**: Built-in defaults defined by the component.
2.  **Global User Config**: Settings provided in `window.UnoConfig.plugins.<componentName>`.
3.  **Direct Global Config**: Settings provided in `window.UnoConfig.<componentName>` (Legacy support).
4.  **Instance Config**: Options passed directly to the component method (e.g., `Alpine.store('uno').box.alert(msg, options)` or `x-uno-box='{...}'`).

### Example

```javascript
// 1. Window Config (Levels 2 & 3)
window.UnoConfig = {
  // Global defaults for all notify instances
  plugins: {
    notify: {
      duration: 5000,
      position: 'top-right'
    }
  }
}

// 2. Instance Call (Level 4)
// This overrides the global 'duration' but keeps 'position'
Alpine.store('uno').notify.show('Hello', { duration: 1000 })
```

## Initialization

You can pass configuration options when manually initializing Uno, or via the global `window.UnoConfig` object for auto-initialization.

### Auto-Initialization (Recommended)

Set `window.UnoConfig` before loading Uno.

```html
<script>
  window.UnoConfig = {
    debug: true, // Enable debug logs
    onError: (err) => {
      // Send to Sentry or other error tracking
      console.error('Uno Error:', err)
    }
  }
</script>
<script src="dist/uno.umd.js"></script>
```

### Manual Initialization

If you are bundling Uno or need manual control:

```javascript
import Uno from '@idiary/uno'
import Alpine from 'alpinejs'

Uno.init(Alpine, {
  debug: true,
  onError: (err) => console.error(err)
})

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
window.UnoConfig = {
  debug: false
}
```

## Component Options Namespace

Component behavior options are recommended to be placed under `plugins.<component>`.

```html
<script>
  window.UnoConfig = {
    plugins: {
      valid: {
        callbackDebounce: 350
      }
    }
  }
</script>
```

Backward compatibility: top-level `valid`, `tooltip`, `notify`, `box` options are still supported but `plugins.<component>` is preferred.
