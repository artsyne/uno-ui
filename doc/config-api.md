# Uno Configuration API

Uno provides a flexible 4-level configuration system that merges plugin defaults, global user config, direct global config, and instance-specific options.

## Configuration Levels

The configuration is resolved in the following order (from lowest to highest priority):

1.  **Plugin Defaults**: Built-in defaults registered by the component.
2.  **Global User Config**: Settings provided via `Uno.setConfig()` in `plugins.<componentName>`.
3.  **Direct Global Config**: Settings provided via `Uno.setConfig()` under the component name itself (`<componentName>`), a shorthand for global overrides.
4.  **Instance Config**: Options passed directly to the component method (e.g., `Alpine.store('uno').box.alert(msg, options)` or `x-uno-box='{...}'`).

### Example

```javascript
// 1. Global Config (Level 2)
import Uno from '@artsyne/uno-ui'

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
<script defer src="https://unpkg.com/alpinejs@3.17.3/dist/cdn.min.js"></script>
<script defer src="https://unpkg.com/@artsyne/uno-ui@0.10.0/dist/uno.umd.js"></script>
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
    import Uno from 'https://unpkg.com/@artsyne/uno-ui@0.10.0/dist/uno.es.js'

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
  <script defer src="https://unpkg.com/alpinejs@3.17.3/dist/cdn.min.js"></script>
</head>
```

#### Module Usage

For NPM/bundler usage:

```javascript
import Uno from '@artsyne/uno-ui'
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
import Uno from '@artsyne/uno-ui'

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

## CSP (Content Security Policy) Configuration

Uno UI provides built-in support for Content Security Policy (CSP). When using CSP with inline styles or scripts, you can configure Uno to automatically inject the nonce attribute.

### What is CSP?

Content Security Policy (CSP) is an added layer of security that helps to detect and mitigate certain types of attacks, including Cross Site Scripting (XSS) and data injection attacks.

### Basic CSP Configuration

#### 1. Server-Side CSP Header

Configure your server to send a CSP header with a nonce:

```http
Content-Security-Policy: default-src 'self'; script-src 'self' 'nonce-abc123'; style-src 'self' 'nonce-abc123'
```

#### 2. Configure Uno with Nonce

Pass the nonce to Uno during initialization:

```javascript
// CDN Usage
<script type="module">
  import Uno from 'https://unpkg.com/@artsyne/uno-ui@0.10.0/dist/uno.es.js'

  Uno.setConfig({
    csp: {
      nonce: 'abc123' // Must match the nonce in your CSP header
    }
  })
</script>
```

```javascript
// NPM/Module Usage
import Uno from '@artsyne/uno-ui'

Uno.setConfig({
  csp: {
    nonce: 'abc123'
  }
})
```

### Dynamic Nonce (Recommended for Production)

In production, generate a unique nonce for each request:

```javascript
// Server-side (Node.js/Express example)
const crypto = require('crypto')

app.use((req, res, next) => {
  const nonce = crypto.randomBytes(16).toString('base64')
  res.locals.nonce = nonce
  res.setHeader(
    'Content-Security-Policy',
    `default-src 'self'; script-src 'self' 'nonce-${nonce}'; style-src 'self' 'nonce-${nonce}'`
  )
  next()
})

// Pass nonce to template
app.get('/', (req, res) => {
  res.render('index.html', { nonce: res.locals.nonce })
})
```

```html
<!-- Template (EJS example) -->
<script type="module">
  import Uno from '/dist/uno.es.js'

  Uno.setConfig({
    csp: {
      nonce: '<%= nonce %>'
    }
  })
</script>
```

### CSP Configuration Options

| Option  | Type     | Default | Description                                                                                        |
| ------- | -------- | ------- | -------------------------------------------------------------------------------------------------- |
| `nonce` | `string` | `null`  | The CSP nonce value. If provided, Uno will automatically inject it into inline styles and scripts. |

### What Uno Does with CSP

When CSP nonce is configured, Uno automatically:

1. **Injects nonce into inline styles**: All dynamically created `<style>` elements will have the `nonce` attribute
2. **Injects nonce into inline scripts**: Any dynamically created `<script>` elements will have the `nonce` attribute
3. **Maintains security**: Ensures all inline content complies with your CSP policy

### Example: Complete CSP Setup

```html
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <title>Uno UI with CSP</title>

    <!-- CSP Header (set by server) -->
    <!-- Content-Security-Policy: default-src 'self'; script-src 'self' 'nonce-abc123'; style-src 'self' 'nonce-abc123' -->

    <!-- Tailwind CSS -->
    <script src="https://cdn.jsdelivr.net/npm/@tailwindcss/browser@4"></script>

    <!-- Configure Uno with CSP nonce -->
    <script type="module">
      import Uno from 'https://unpkg.com/@artsyne/uno-ui@0.10.0/dist/uno.es.js'

      Uno.setConfig({
        debug: false,
        csp: {
          nonce: 'abc123' // Must match server-generated nonce
        },
        ui: {
          base: 'rounded-lg shadow-md'
        }
      })
    </script>

    <!-- Load Alpine -->
    <script defer src="https://cdn.jsdelivr.net/npm/alpinejs@3.17.3/dist/cdn.min.js"></script>
  </head>
  <body>
    <button x-data @click="$store.uno.box.alert('Hello!')">Click Me</button>
  </body>
</html>
```

### Troubleshooting CSP Issues

| Issue                     | Cause                                     | Solution                                                    |
| ------------------------- | ----------------------------------------- | ----------------------------------------------------------- |
| **Styles not applied**    | CSP blocking inline styles                | Ensure `style-src` includes the nonce or `'unsafe-inline'`  |
| **Scripts not executing** | CSP blocking inline scripts               | Ensure `script-src` includes the nonce or `'unsafe-inline'` |
| **Mixed Content errors**  | Loading resources from HTTP on HTTPS page | Ensure all resources use HTTPS                              |
| **Nonce mismatch**        | Nonce in header doesn't match config      | Verify nonce is correctly passed from server to client      |

### Best Practices

1. **Always use unique nonces**: Generate a new nonce for each request
2. **Don't use 'unsafe-inline'**: Use nonce instead of `'unsafe-inline'` for better security
3. **Test in development**: Enable CSP in development to catch issues early
4. **Monitor CSP violations**: Use the `Content-Security-Policy-Report-Only` header to test without blocking

### Learn More

- [MDN: Content Security Policy](https://developer.mozilla.org/en-US/docs/Web/HTTP/CSP)
- [CSP Evaluator](https://csp-evaluator.withgoogle.com/)
- [OWASP: Content Security Policy Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Content_Security_Policy_Cheat_Sheet.html)
