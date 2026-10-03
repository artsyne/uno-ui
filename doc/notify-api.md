# Uno Notify API Documentation

The Notify component provides a lightweight, toast-style notification system. It supports multiple screen positions, varied notification types (success, error, warning, info), and includes a global progress bar for HTMX integration.

## Overview

- **Core Function**: Display temporary feedback messages (toasts) and loading indicators.
- **Architecture**: Store-based (`Alpine.store('uno').notify`) with directive triggers (`x-uno-notify`).
- **Positions**: Supports 9 screen positions (e.g., top-right, bottom-center).
- **HTMX Support**: Automatic global progress bar integration during HTMX requests.
- **A11y**: The close button's screen-reader label is localized via the `box.close` i18n key.

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

### `bindNotifyToHtmx(options)`

Binds the Notify component to HTMX events for automatic notification display during HTMX requests. This provides a seamless integration where success/error notifications and progress bar are automatically shown based on HTMX request lifecycle.

- **Parameters**:
  - `options` (Object, optional) - Configuration options:
    - `showSuccess` (Boolean, default: `true`) - Show success notification after successful requests.
    - `showError` (Boolean, default: `true`) - Show error notification on failed requests.
    - `showProgress` (Boolean, default: `true`) - Show progress bar during active requests.
    - `successMessage` (String, default: `'Operation completed'`) - Default success message.
    - `errorMessage` (String, default: `'Operation failed'`) - Default error message.
- **Returns**: `void`
- **Example**:

  ```javascript
  // Initialize with defaults
  Alpine.store('uno').htmx.bindNotifyToHtmx()

  // Custom messages
  Alpine.store('uno').htmx.bindNotifyToHtmx({
    showSuccess: true,
    showError: true,
    showProgress: true,
    successMessage: 'Saved successfully!',
    errorMessage: 'Something went wrong'
  })

  // Only show progress bar, disable notifications
  Alpine.store('uno').htmx.bindNotifyToHtmx({
    showSuccess: false,
    showError: false,
    showProgress: true
  })
  ```

## HTMX Integration

### Automatic Setup

The `bindNotifyToHtmx()` method binds to the following HTMX events:

| Event                | Trigger Timing                 | Behavior                                                                 |
| :------------------- | :----------------------------- | :----------------------------------------------------------------------- |
| `htmx:beforeRequest` | Start of an HTMX request.      | Starts the global progress bar animation.                                |
| `htmx:afterRequest`  | Completion of an HTMX request. | Stops/Completes the progress bar, shows success toast if response is OK. |
| `htmx:responseError` | Server returns 4xx/5xx error.  | Stops progress bar and shows error toast.                                |
| `htmx:sendError`     | HTMX request failed to send.   | Stops progress bar and shows error toast.                                |
| `htmx:timeout`       | HTMX request timed out.        | Stops progress bar and shows error toast.                                |

### Server-Sent Custom Messages

The integration supports custom messages sent from the server via response headers:

| Header              | Description                                       |
| :------------------ | :------------------------------------------------ |
| `X-Notify-Message`  | Custom success message (overrides default).       |
| `X-Notify-Error`    | Custom error message (overrides default).         |
| `X-Notify-Suppress` | Set to `'true'` to suppress success notification. |

Example server response (any language):

```python
# Python/Flask example
return jsonify({'result': 'ok'}, headers={
    'X-Notify-Message': 'Item saved successfully!'
})
```

## Accessibility

Notifications are announced by screen readers with politeness graded by type (B7, `docs/a11y-design.md`):

| Type               | `role`   | `aria-live` | Behavior                                                    |
| :----------------- | :------- | :---------- | :---------------------------------------------------------- |
| `error`, `warning` | `alert`  | `assertive` | Interrupts the current utterance — urgent feedback.         |
| `success`, `info`  | `status` | `polite`    | Queued after current speech — avoids breaking user context. |

Both attributes are bound dynamically per item in the template; no configuration is needed.

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

## Troubleshooting

### Common Issues

| Issue                          | Possible Cause                          | Solution                                                  |
| :----------------------------- | :-------------------------------------- | :-------------------------------------------------------- |
| **Notification not appearing** | `z-index` conflict or container hidden. | Check CSS. Notify container uses high z-index.            |
| **Progress bar stuck**         | HTMX request failed or didn't complete. | Manually call `stopProgress()` in `catch/finally` block.  |
| **Missing styles**             | Tailwind classes purged or not loaded.  | Ensure `safelist` includes Uno classes if using PurgeCSS. |

## Version History

| Version | Date       | Changes                                                                                                                                                               |
| :------ | :--------- | :-------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| v0.8.0  | 2026-02-28 | Added `bindNotifyToHtmx()` method for automatic HTMX integration with progress bar, success/error notifications, and server-sent custom messages support via headers. |
| v0.7.0  | 2026-01-20 | Refactored for **Lazy Initialization**. Containers injected on first usage.                                                                                           |
| v0.6.0  | 2025-12-18 | Initial stable release.                                                                                                                                               |
