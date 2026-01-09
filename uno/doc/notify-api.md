# Uno Notify API Documentation

The Notify component provides a lightweight, toast-style notification system. It supports multiple screen positions, varied notification types (success, error, warning, info), and includes a global progress bar for HTMX integration.

## Overview

- **Core Function**: Display temporary feedback messages (toasts) and loading indicators.
- **Architecture**: Store-based (`Alpine.store('uno').notify`) with directive triggers (`x-uno-notify`).
- **Positions**: Supports 9 screen positions (e.g., top-right, bottom-center).
- **HTMX Support**: Automatic global progress bar integration during HTMX requests.

## Basic Usage

### 1. Directive Usage

Trigger notifications from HTML elements.

```html
<!-- Simple String -->
<button x-uno-notify="'Operation saved'">Save</button>

<!-- Configuration Object -->
<button x-uno-notify="{ message: 'Deleted', type: 'error', duration: 2000 }">Delete</button>
```

### 2. JavaScript Usage

Use the global store shortcuts.

```javascript
// Success
Alpine.store('uno').notify.success('Profile updated')

// Error
Alpine.store('uno').notify.error('Connection failed')

// Custom
Alpine.store('uno').notify.show('Info message', {
  position: 'bottom-right'
})
```

## API Parameters

The following options can be passed to `show()` or used in the `x-uno-notify` configuration object.

| Name       | Type      | Default       | Required | Description                                                          |
| :--------- | :-------- | :------------ | :------- | :------------------------------------------------------------------- |
| `message`  | `String`  | -             | Yes      | The text content of the notification. Supports i18n keys.            |
| `type`     | `String`  | `'info'`      | No       | Visual style. Values: `'info'`, `'success'`, `'warning'`, `'error'`. |
| `duration` | `Number`  | `3000`        | No       | Time in ms before auto-dismiss. Set `0` to disable auto-dismiss.     |
| `position` | `String`  | `'top-right'` | No       | Screen position. See supported positions below.                      |
| `closable` | `Boolean` | `true`        | No       | Whether to show the close button.                                    |
| `hxUrl`    | `String`  | -             | No       | If set, loads content via HTMX into the notification body.           |

### Supported Positions

- `top-left`, `top-center`, `top-right`
- `center-left`, `center`, `center-right`
- `bottom-left`, `bottom-center`, `bottom-right`

## Events

The Notify component integrates with HTMX global events to manage the progress bar.

| Event                | Trigger Timing                 | Behavior                                             |
| :------------------- | :----------------------------- | :--------------------------------------------------- |
| `htmx:beforeRequest` | Start of an HTMX request.      | Starts the global progress bar animation.            |
| `htmx:afterRequest`  | Completion of an HTMX request. | Stops/Completes the progress bar.                    |
| `htmx:responseError` | Server returns 4xx/5xx error.  | Stops progress bar and optionally shows error toast. |

## Methods

Accessed via `Alpine.store('uno').notify`.

### `show(message, options)`

Displays a notification.

- **Parameters**:
  - `message` (String) - Content or i18n key.
  - `options` (Object, optional) - Config object.
- **Returns**: `Number` - The notification ID.
- **Example**:
  ```javascript
  Alpine.store('uno').notify.show('Hello', { type: 'success' })
  ```

### `success(message, options)`

Shortcut for `show` with `type: 'success'`.

- **Parameters**: Same as `show`.
- **Returns**: `Number` - ID.

### `error(message, options)`

Shortcut for `show` with `type: 'error'`.

- **Parameters**: Same as `show`.
- **Returns**: `Number` - ID.

### `warning(message, options)`

Shortcut for `show` with `type: 'warning'`.

- **Parameters**: Same as `show`.
- **Returns**: `Number` - ID.

### `info(message, options)`

Shortcut for `show` with `type: 'info'`.

- **Parameters**: Same as `show`.
- **Returns**: `Number` - ID.

### `remove(id)`

Removes a specific notification immediately.

- **Parameters**: `id` (Number) - The notification ID.
- **Returns**: `void`

### `startProgress()`

Manually starts the global progress bar (loops 0-100%).

- **Parameters**: None.
- **Returns**: `void`

### `stopProgress()`

Manually stops the global progress bar (fills to 100% then fades).

- **Parameters**: None.
- **Returns**: `void`

## Examples

### Sticky Notification

A notification that stays until the user closes it manually.

```javascript
Alpine.store('uno').notify.show('Please verify your email.', {
  type: 'warning',
  duration: 0, // Disable auto-dismiss
  closable: true
})
```

### HTMX Progress Bar Control

While handled automatically, you can manually trigger it for non-HTMX async operations.

```javascript
const loadData = async () => {
  Alpine.store('uno').notify.startProgress()
  try {
    await fetch('/api/data')
    Alpine.store('uno').notify.success('Data loaded')
  } catch (e) {
    Alpine.store('uno').notify.error('Error loading data')
  } finally {
    Alpine.store('uno').notify.stopProgress()
  }
}
```

## Version History

| Version | Date       | Changes                                              |
| :------ | :--------- | :--------------------------------------------------- |
| v1.0.0  | 2025-01-09 | Initial release with Toasts and Global Progress Bar. |
