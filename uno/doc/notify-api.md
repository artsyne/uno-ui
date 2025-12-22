# Notify Component

A lightweight notification component built on Alpine.js, supporting multiple positions, message types, auto-close, and dark mode.

## Features

- **Multi-Position**: `top-left`, `top-center`, `top-right`, `center-left`, `center`, `center-right`, `bottom-left`, `bottom-center`, `bottom-right`.
- **Message Types**: `info`, `success`, `warning`, `error`.
- **Dark Mode**: Perfectly adapted to Tailwind CSS Dark Mode (`dark:` variant).
- **Auto Close**: Configurable display duration (`timeout`). Set to `0` for manual close.
- **Loading State**: Supports global loading indicators (`indicator` or `progressbar`).
- **Customizable Styles**: Uses standard Tailwind CSS classes, easy to override.
- **Internationalization**: Supports `Alpine.store('i18n')` for loading text.
- **Accessibility**: Uses `role="alert"` for screen readers.

## Quick Start

### Basic Usage

Use the global store `$store.notify.show()` to trigger notifications.

```html
<!-- Trigger Buttons -->
<button @click="$store.notify.show('Operation Successful', { type: 'success' })">Success</button>
<button @click="$store.notify.show('Something went wrong', { type: 'error' })">Error</button>
```

The container is automatically injected into the DOM when needed.

## Programmatic Usage

The component exposes methods via `Alpine.store('notify')` or the global alias `Uno.notify`.

### Global Alias: `Uno.notify`

The same methods are available globally via `Uno.notify` for convenience:

```javascript
Uno.notify.success('Saved!')
Uno.notify.error('Failed!')
Uno.notify.loading(true)
```

### `show(text, options)`

Shows a notification.

| Parameter | Type     | Default | Description           |
| --------- | -------- | ------- | --------------------- |
| `text`    | `String` | -       | Notification content. |
| `options` | `Object` | `{}`    | Configuration object. |

**Options Object:**

| Option            | Type      | Default       | Description                                                                                                                                                  |
| ----------------- | --------- | ------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `type`            | `String`  | `'info'`      | Message type: `info`, `success`, `warning`, `error`.                                                                                                         |
| `timeout`         | `Number`  | `3000`        | Auto-close duration in milliseconds. `0` to disable.                                                                                                         |
| `position`        | `String`  | `'top-right'` | Container position. Options: `top-left`, `top-center`, `top-right`, `center-left`, `center`, `center-right`, `bottom-left`, `bottom-center`, `bottom-right`. |
| `showCloseButton` | `Boolean` | `true`        | Whether to show the close (x) button.                                                                                                                        |

### Shortcuts

- `success(text, options)`
- `error(text, options)`
- `warning(text, options)`
- `info(text, options)`

### `close(id)`

Manually removes a notification by ID.

## Loading Control

Controls global loading state display (useful for page-level feedback).

### `loading(isLoading)`

| Parameter   | Type      | Default | Description                                                  |
| ----------- | --------- | ------- | ------------------------------------------------------------ |
| `isLoading` | `Boolean` | `false` | `true` to show loading indicator (spinner), `false` to hide. |

### `progress(isLoading)`

| Parameter   | Type      | Default | Description                                                   |
| ----------- | --------- | ------- | ------------------------------------------------------------- |
| `isLoading` | `Boolean` | `false` | `true` to show progress bar (top of screen), `false` to hide. |

## Internationalization (i18n)

The component uses `Alpine.store('i18n')` for loading text.

- **Keys**:
  - `notify.loading`: Text displayed when using `loading(true)` (if applicable/extended).

## Styling

Component styles rely on Tailwind CSS classes:

- `notify-item`: Base style for notification items.
- `notify-success`, `notify-warning`, `notify-error`, `notify-info`: Color variants.
- **Dark Mode**: Automatically adapts to Dark Mode (e.g., dark background, light text).
