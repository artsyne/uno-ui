# Tooltip Component

A lightweight, smart positioning Tooltip component built on Alpine.js.

## Features

- **Smart Positioning**: Automatically detects viewport edges and adjusts position (Top/Bottom/Left/Right).
- **Interaction Modes**: Supports `hover` (mouse hover, default) and `click` (click trigger).
- **Themes**: Built-in `dark` (default) and `light` themes.
- **Rich Content**: Supports plain text or HTML content (referenced via DOM element).
- **High Performance**: Uses `requestAnimationFrame` for rendering optimization to prevent jitter.
- **Accessibility**:
  - Adds `role="tooltip"` to content.
  - Automatically manages `aria-describedby` on trigger element.

## Quick Start

### Basic Structure

Attach the directive `x-tooltip` on any trigger element.

```html
<button class="btn" x-tooltip="'Hello World'">Hover Me</button>
```

> **Note**: Tooltip uses a global singleton container injected on-demand by `Alpine.store('uno').tooltip.show()`.

## Programmatic Usage

### Store API

The component exposes a global API for manual control.

#### `Alpine.store('uno').tooltip.show(el, content, config)`

Manually show a tooltip anchored to an element.

```javascript
Alpine.store('uno').tooltip.show(document.getElementById('my-btn'), 'Manual Tooltip', { placement: 'bottom' })
```

#### `Alpine.store('uno').tooltip.hide()`

Hides the current tooltip.

#### `Alpine.store('uno').tooltip.updatePosition()`

Recalculates position (useful if trigger moves).

## Configuration / Options

### Options Object

| Option        | Type                                     | Default   | Description                                                                  |
| ------------- | ---------------------------------------- | --------- | ---------------------------------------------------------------------------- |
| `placement`   | `'top' \| 'bottom' \| 'left' \| 'right'` | `'top'`   | Preferred placement. Flips automatically if space is insufficient.           |
| `trigger`     | `'hover' \| 'click' \| 'focus'`          | `'hover'` | Trigger mode. In `click` mode, clicking outside closes it.                   |
| `theme`       | `'dark' \| 'light'`                      | `'dark'`  | Visual theme.                                                                |
| `timeout`     | `Number`                                 | `100`     | Delay in ms before closing on mouse leave (hover mode).                      |
| `target`      | `String`                                 | `null`    | CSS selector of the DOM element to use as content source.                    |
| `onShow`      | `Function`                               | `null`    | Callback function triggered when tooltip is shown. Receives content element. |
| `trustedHTML` | `Boolean`                                | `false`   | If true, disables HTML sanitization. Use with caution.                       |

## Security & Sanitization

By default, all HTML content (whether from `content` prop or `target` DOM) is sanitized using `Uno.dom.sanitizeHTML`.

- **Dangerous Tags**: `<script>`, `<style>`, `<iframe>`, `<object>`, `<embed>` are completely stripped.
- **Attributes**: Only safe attributes (e.g. `class`, `id`) are allowed. `href` in `<a>` tags is validated (no `javascript:`).
- **Trusted Mode**: Set `trustedHTML: true` to bypass sanitization if you trust the content source.

## Examples

### 1. Specific Position (Bottom)

```html
<button x-tooltip.bottom="'Bottom Tip'">Hover Me</button>
```

### 2. Click Trigger + Light Theme (Popover Style)

```html
<button x-tooltip.click="{ content: 'Light Theme Popover', theme: 'light' }">Click Me</button>
```

### 3. HTML Content (Target Mode)

```html
<!-- Hidden Content Source -->
<div id="tip-content" style="display: none;">
  <strong>Bold Title</strong><br />
  <span class="text-blue-500">Supports HTML styles</span>
</div>

<!-- Tooltip Instance -->
<button x-tooltip="{ target: '#tip-content' }">Hover to see HTML</button>
```

## Developer Notes

- **Core Positioning**: Positioning logic reuses `src/utils/position.js`, consistent with Box/Notify components.
- **Debugging**: Enable debug mode in `src/utils/log.js` to view internal logs in the console.
- **A11y**: Ensure your trigger element is focusable (e.g., `<button>` or `tabindex="0"`) for keyboard users.
