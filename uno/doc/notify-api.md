# Notify Component

A lightweight notification component built on Alpine.js, supporting multiple positions, message types, auto-close, and dark mode.

## Features

- **Multi-Position**: `top-left`, `top-right`, `bottom-left`, `bottom-right`.
- **Message Types**: `info`, `success`, `warning`, `error`.
- **Dark Mode**: Perfectly adapted to Tailwind CSS Dark Mode (`dark:` variant).
- **Auto Close**: Configurable display duration (`timeout`). Set to `0` for manual close.
- **Loading State**: Supports global loading indicators (`indicator` or `progressbar`).
- **Customizable Styles**: Uses standard Tailwind CSS classes, easy to override.
- **Internationalization**: Supports `Alpine.store('i18n')` for loading text.

## Quick Start

### Basic Usage

Initialize the component on a container and use the `show` method to display notifications.

```html
<div x-data="notify()">
  <!-- Trigger Buttons -->
  <button @click="show('Operation Successful', { type: 'success' })" class="btn">Success</button>
  <button @click="show('Something went wrong', { type: 'error' })" class="btn">Error</button>
  
  <!-- Notification Container (Managed by component) -->
  <div :class="containerClasses" :style="containerStyle">
    <template x-for="item in items" :key="item.id">
      <div :class="itemClasses(item.type)" @click="remove(item.id)">
        <span x-text="item.text"></span>
      </div>
    </template>
  </div>
</div>
```

## Programmatic Usage

The component exposes methods via the Alpine data object or can be accessed if you make the store global (depending on implementation). Within the `x-data` scope:

### `notify({ position })`
Initializes the component.

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `position` | `String` | `'top-right'` | Container position. Options: `top-left`, `top-right`, `bottom-left`, `bottom-right`. |

### `show(text, options)`
Shows a notification.

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `text` | `String` | - | Notification content. |
| `options` | `Object` | `{}` | Configuration object. |

**Options Object:**

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `type` | `String` | `'info'` | Message type: `info`, `success`, `warning`, `error`. |
| `duration` | `Number` | `3000` | Auto-close duration in milliseconds. `0` to disable. |

### `remove(id)`
Manually removes a notification by ID.

## Loading Control

Controls global loading state display (useful for page-level feedback).

### `loading(isLoading, mode)`

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `isLoading` | `Boolean` | `false` | `true` to show, `false` to hide. |
| `mode` | `String` | `'indicator'` | Loading style: `'indicator'` (spinner) or `'progressbar'`. |

```html
<!-- Spinner -->
<button @click="loading(true, 'indicator')">Start Loading</button>
<button @click="loading(false)">Stop Loading</button>

<!-- Progress Bar -->
<button @click="loading(true, 'progressbar')">Start Progress</button>
<button @click="loading(false)">Stop Progress</button>
```

## Internationalization (i18n)

The component uses `Alpine.store('i18n')` for loading text.

- **Keys**:
  - `notify.loading`: Text displayed when using `loading(true)` (if applicable/extended).

## Styling

Component styles rely on Tailwind CSS classes:

- `notify-item`: Base style for notification items.
- `notify-success`, `notify-warning`, `notify-error`, `notify-info`: Color variants.
- **Dark Mode**: Automatically adapts to Dark Mode (e.g., dark background, light text).
