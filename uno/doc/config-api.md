# Uno Configuration API

Uno can be configured during initialization to control global behavior, debugging, and error handling.

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
            console.error('Uno Error:', err);
        }
    };
</script>
<script src="dist/uno.umd.js"></script>
```

### Manual Initialization

If you are bundling Uno or need manual control:

```javascript
import Uno from '@idiary/uno';
import Alpine from 'alpinejs';

Uno.init(Alpine, {
    debug: true,
    onError: (err) => console.error(err)
});

Alpine.start();
```

## Configuration Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `debug` | `boolean` | `true` | Enable/Disable internal logs. Set to `false` for production. |
| `onError` | `function` | `null` | Custom error handler callback. Receives the error object. |
| `csp` | `object` | `undefined` | Content Security Policy settings (e.g., `{ nonce: '...' }`). |

## Debug Mode

When `debug` is `true` (default), Uno prints lifecycle events and internal status to the console with `[Uno Info]` or `[Uno Warn]` prefixes.

To disable logs in production:

```javascript
window.UnoConfig = {
    debug: false
};
```
