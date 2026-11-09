# Error Handling

Uno installs global listeners for unhandled promise rejections and uncaught errors, and routes
its own internal failures through a small logger. This page explains what is captured, what is
not, and how to connect your own reporting.

## What Uno does

During initialization, once per page, Uno adds:

- `window.addEventListener('unhandledrejection', ...)`
- `window.addEventListener('error', ...)`

Both call `log.error()`. If `onError` is configured, it is called with the same arguments. If
`debug` is true, the message is also printed to the console with an `[Uno Error]` prefix.
Duplicate `Error` objects are suppressed with a `WeakSet`.

The listeners do not call `preventDefault()`; the browser still reports the error in its own
way.

Uno also catches failures while loading a plugin and logs `Failed to load plugin` with the
error.

## What is not captured

- Errors that you catch yourself.
- HTMX request failures as promise rejections. HTMX reports them as DOM events
  (`htmx:responseError`, `htmx:sendError`, `htmx:timeout`), not rejected promises. Listen to
  those events or use `Uno.htmx.onRequest({ end })`.
- Rejections that happen before Uno initializes.
- Errors in another JavaScript realm (iframe, worker).
- Anything when `onError` is unset and `debug` is false: `log.error()` has no visible output in
  that configuration.

## Configuration

```javascript
Uno.setConfig({
  debug: true,
  onError: (...args) => {
    console.error('[app]', ...args)
    // window.Sentry?.captureException(args[0])
  }
})
```

- `onError` receives the arguments passed to `log.error(...)`. That may be an `Error`, an event
  payload, a string plus an object, or several values.
- `onError` and `debug` are applied during `Uno.init()`. Changing them later with `setConfig()`
  does not update the logger.
- An exception thrown by `onError` is swallowed; in debug mode it prints `Reporter failed`.

## HTMX errors

```javascript
document.addEventListener('htmx:responseError', (event) => {
  console.error('HTMX response error', event.detail.xhr.status)
  Alpine.store('uno').notify.error('Request failed')
})

document.addEventListener('htmx:sendError', () => {
  Alpine.store('uno').notify.error('Network error')
})

document.addEventListener('htmx:timeout', () => {
  Alpine.store('uno').notify.error('Request timed out')
})
```

Or with the helper:

```javascript
Uno.htmx.onRequest({
  end: (event) => {
    if (event.detail.failed) {
      console.warn('HTMX request failed', event.detail)
    }
  }
})
```

`Uno.htmx.fragment()` rejects on failure, so a `try/catch` around it is enough:

```javascript
try {
  const html = await Uno.htmx.fragment('GET', '/api/fragment')
  Alpine.store('uno').box.open({ title: 'Loaded', content: html })
} catch (error) {
  Alpine.store('uno').notify.error('Unable to load content')
}
```

Box has no `onError` callback. A failed `hxUrl` request surfaces through the HTMX events above.
Notify has an `onError` item callback for its own HTMX content.

## Component-level handling

### Box

`onOk` / `onCancel` / `onClose` are the available hooks. A rejected promise inside `onOk` is
not caught by Uno; catch it yourself and return `false` to keep the dialog open.

```javascript
Alpine.store('uno').box.open({
  title: 'Delete item',
  content: 'This cannot be undone.',
  showOkButton: true,
  showCancelButton: true,
  onOk: async () => {
    try {
      const response = await fetch('/api/delete', { method: 'DELETE' })
      if (!response.ok) throw new Error(`HTTP ${response.status}`)
      Alpine.store('uno').notify.success('Deleted')
      return true
    } catch (error) {
      Alpine.store('uno').notify.error('Delete failed')
      return false
    }
  }
})
```

`confirm()` / `prompt()` compose their internal Promise resolver with your callbacks. See the
[Box API](box-api.md#store-api) for the exact close semantics.

### Notify

`onContentLoad(el)` and `onError(detail)` are available for `hxUrl` notifications. A failed
content load also shows a generic `Content load failed` toast.

### Valid

Custom validators should catch their own network errors and return a message. A thrown
exception is treated as a failed validation with the message `Validation error`. During HTMX
form submission, unexpected validation failures are logged and the request is blocked (fail
closed).

## Best practices

1. **Handle promises you create.** Uno's global handler is a safety net, not a replacement for
   `try/catch`.
2. **Return a user-facing message** from validators and callbacks instead of only logging.
3. **Use `debug: true` during development**, and either `onError` or your own reporting in
   production.
4. **Do not assume `onError` receives an `Error`.** Inspect the arguments.
5. **Listen to HTMX events** rather than expecting a rejected promise.
6. **Never report secrets** to a third-party service from the client.

## Error types

Uno uses standard JavaScript `Error` objects. It does not define `NetworkError`,
`ValidationError` or a custom `UnoError` class. Common messages include:

```text
[Uno Error] Unhandled Promise Rejection: <reason>
[Uno Error] Uncaught Error: { message, filename, lineno, colno, error }
[Uno Error] HTMX fragment fetch failed: { method, url, error }
[Uno Warn] Uno.use() called after initialization. Plugin might not work correctly.
[Uno Warn] setConfig called with invalid argument, expected object
```

Warnings are only printed when debug mode is enabled.

## Related

- [Troubleshooting](troubleshooting.md)
- [HTMX integration](htmx.md)
- [Configuration](configuration.md)
