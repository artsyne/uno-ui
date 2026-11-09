# Error Handling Guide

This document explains how Uno UI handles errors and provides best practices for error management in your application.

## Overview

Uno UI v0.7.0+ includes a comprehensive error handling system that:

- ✅ Catches unhandled Promise rejections
- ✅ Catches uncaught JavaScript errors
- ✅ Provides custom error reporting hooks
- ✅ Prevents duplicate error reports
- ✅ Logs errors with detailed context

## Global Error Boundary

Uno automatically sets up global error handlers during initialization. These handlers catch errors that would otherwise go unhandled.

### Unhandled Promise Rejections

When a Promise is rejected without a `.catch()` handler, Uno catches it:

```javascript
// ❌ Bad: No error handling
async function loadData() {
  const response = await fetch('/api/data')
  return response.json()
}

// ✅ Good: Explicit error handling
async function loadData() {
  try {
    const response = await fetch('/api/data')
    if (!response.ok) throw new Error('Failed to fetch')
    return response.json()
  } catch (error) {
    console.error('Load failed:', error)
    throw error // Re-throw if needed
  }
}
```

**What Uno catches:**

- Async/await without try-catch
- Promise chains without `.catch()`
- HTMX ajax errors
- Component lifecycle errors

### Uncaught JavaScript Errors

Uno catches errors that escape try-catch blocks:

```javascript
// This error will be caught by Uno's global handler
function riskyOperation() {
  throw new Error('Something went wrong')
}

// Uno logs:
// [Uno Error] Uncaught Error: {
//   message: "Something went wrong",
//   filename: "app.js",
//   lineno: 42,
//   colno: 10,
//   error: Error object
// }
```

## Custom Error Reporting

You can provide a custom error handler to send errors to your logging service:

### Configuration

```javascript
import { Uno } from '@artsyne/uno-ui'

Uno.init({
  debug: true,
  onError: (error) => {
    // Send to logging service
    console.error('Uno Error:', error)

    // Example: Send to Sentry
    if (window.Sentry) {
      window.Sentry.captureException(error)
    }

    // Example: Send to custom API
    fetch('/api/log-error', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        message: error.message,
        stack: error.stack,
        timestamp: new Date().toISOString()
      })
    })
  }
})
```

### Error Handler Details

The `onError` callback receives:

- **Error objects**: JavaScript Error instances
- **Error context**: Additional information when available
- **Component errors**: Errors from Uno components

## Component-Level Error Handling

### Box Component

```javascript
// Box loads HTMX content from hxUrl; onContentLoad runs after the response is swapped in.
Alpine.store('uno').box.open({
  title: 'Load Content',
  hxUrl: '/api/content',
  hxMethod: 'GET',
  onContentLoad: (el) => {
    // The loaded content element is available here.
  }
})
```

Box has no `onError` callback: a failed `hxUrl` request surfaces through the HTMX events below, so
report the failure with the same global listeners.

### HTMX Integration

```javascript
// HTMX error handling
Uno.htmx.onRequest({
  start: (event) => {
    console.log('Request started')
  },
  end: (event) => {
    if (event.detail.failed) {
      console.error('Request failed:', event.detail)
      Alpine.store('uno').notify.error('Request failed')
    }
  }
})
```

### HTMX Fragment Errors

```javascript
// Error handling with fragment fetch
try {
  const html = await Uno.htmx.fragment('GET', '/api/fragment')
  // Use html
} catch (error) {
  console.error('Fragment fetch failed:', error)
  // Uno already logged this error
  // Handle gracefully
  Alpine.store('uno').notify.error('Failed to load content')
}
```

## Best Practices

### 1. Always Handle Promises

```javascript
// ❌ Bad: Silent failure
async function loadUser() {
  const response = await fetch('/api/user')
  return response.json()
}

// ✅ Good: Explicit error handling
async function loadUser() {
  try {
    const response = await fetch('/api/user')
    if (!response.ok) {
      throw new Error(`HTTP ${response.status}: ${response.statusText}`)
    }
    return response.json()
  } catch (error) {
    console.error('Failed to load user:', error)
    Alpine.store('uno').notify.error('Failed to load user')
    return null
  }
}
```

### 2. Provide User Feedback

```javascript
// Show user-friendly error messages
Alpine.store('uno').box.open({
  title: 'Delete Item',
  content: 'Are you sure?',
  onOk: async () => {
    try {
      const response = await fetch('/api/delete', { method: 'DELETE' })
      if (!response.ok) throw new Error('Delete failed')

      Alpine.store('uno').notify.success('Item deleted')
      return true
    } catch (error) {
      Alpine.store('uno').notify.error('Failed to delete item')
      return false // Prevents box from closing
    }
  }
})
```

### 3. Use Debug Mode During Development

```javascript
Uno.init({
  debug: true, // Enables verbose logging
  onError: (error) => {
    console.error('Debug mode error:', error)
  }
})
```

### 4. Handle Network Errors

```javascript
async function fetchWithRetry(url, retries = 3) {
  for (let i = 0; i < retries; i++) {
    try {
      const response = await fetch(url)
      if (!response.ok) throw new Error('HTTP error')
      return response
    } catch (error) {
      if (i === retries - 1) {
        // Last retry failed
        Alpine.store('uno').notify.error('Network error. Please try again.')
        throw error
      }
      // Wait before retry
      await new Promise((resolve) => setTimeout(resolve, 1000 * (i + 1)))
    }
  }
}
```

## Common Error Scenarios

### 1. HTMX Request Failures

**Symptom**: HTMX requests fail silently

**Solution**:

```javascript
// Setup global HTMX error handling
document.addEventListener('htmx:responseError', (event) => {
  console.error('HTMX error:', event.detail)
  Alpine.store('uno').notify.error('Request failed')
})

document.addEventListener('htmx:sendError', (event) => {
  console.error('HTMX send error:', event.detail)
  Alpine.store('uno').notify.error('Network error')
})
```

### 2. Component Initialization Errors

**Symptom**: Components don't work as expected

**Solution**:

```javascript
// Ensure proper initialization order
document.addEventListener('DOMContentLoaded', () => {
  // 1. Initialize Uno
  Uno.init({ debug: true })

  // 2. Start Alpine
  window.Alpine.start()
})
```

### 3. CSP (Content Security Policy) Issues

**Symptom**: Components don't render or styles don't apply

**Solution**:

```javascript
// Provide CSP nonce
Uno.init({
  csp: {
    nonce: 'your-nonce-here'
  }
})
```

## Debugging Tips

### Enable Debug Mode

```javascript
Uno.init({ debug: true })
// Outputs: [Uno Debug] / [Uno Info] / [Uno Warn] messages to console
```

### Check Component State

```javascript
// Inspect store state
console.log(Alpine.store('uno').box)
console.log(Alpine.store('uno').notify)
console.log(Alpine.store('uno').tooltip)
```

### Monitor Errors

```javascript
// Count errors
let errorCount = 0
Uno.init({
  onError: (error) => {
    errorCount++
    console.log(`Error #${errorCount}:`, error)
  }
})
```

## Error Types

### UnoError Classes

Uno uses standard JavaScript Error objects. Common error types:

- **TypeError**: Invalid parameter types
- **ReferenceError**: Component not initialized
- **NetworkError**: HTMX/fetch failures
- **ValidationError**: Form validation failures

### Example Error Messages

```
[Uno Error] Unhandled Promise Rejection: Failed to fetch
[Uno Error] HTMX fragment fetch failed: {method: "GET", url: "/api/data", error: Error}
[Uno Error] Uncaught Error: {message: "Cannot read property 'box' of undefined", ...}
[Uno Warn] Uno.use() called after initialization. Plugin might not work correctly.
```

## See Also

- [FAQ](./faq.md) - Frequently Asked Questions
- [TypeScript Usage](./typescript-usage.md) - Type-safe error handling
- [Configuration API](./config-api.md) - Error reporting configuration
