# Tooltip Component

A lightweight, smart positioning Tooltip component built on Alpine.js.

## Features

- **Smart Positioning**: Automatically detects viewport edges and adjusts position (Top/Bottom/Left/Right).
- **Interaction Modes**: Supports `hover` (mouse hover, default) and `click` (click trigger).
- **Themes**: Built-in `dark` (default) and `light` themes.
- **Rich Content**: Supports plain text or HTML content (referenced via DOM element).
- **High Performance**: Uses `requestAnimationFrame` for rendering optimization to prevent jitter.

## Quick Start

### Basic Structure

Initialize the component using `x-data="tooltip(...)"` and bind the trigger and content elements.

```html
<div x-data="tooltip('Hello World')">
    <!-- Trigger Element -->
    <button data-trigger @mouseenter="onMouseEnter" @mouseleave="onMouseLeave" class="btn">
        Hover Me
    </button>

    <!-- Tooltip Body -->
    <div 
        x-ref="tooltipRef" 
        x-show="open" 
        x-transition 
        :class="classes"
        :style="{ top: top + 'px', left: left + 'px', visibility: positioned ? 'visible' : 'hidden' }"
        style="display: none;"
    >
        <span x-html="content"></span>
        <div :class="arrowClasses"></div>
    </div>
</div>
```

> **Note**: `style="display: none;"` prevents Flash of Unstyled Content (FOUC) before Alpine initialization.

## Programmatic Usage

### `tooltip(initialContent, options)`

Initializes the Tooltip data object.

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `initialContent` | `String` | `''` | Tooltip text. Ignored if `target` option is used. |
| `options` | `Object` | `{}` | Configuration object. |

## Configuration / Options

### Options Object

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `position` | `'top' \| 'bottom' \| 'left' \| 'right'` | `'top'` | Preferred position. Flips automatically if space is insufficient. |
| `trigger` | `'hover' \| 'click'` | `'hover'` | Trigger mode. In `click` mode, clicking outside closes it. |
| `theme` | `'dark' \| 'light'` | `'dark'` | Visual theme. |
| `timeout` | `Number` | `100` | Delay in ms before closing on mouse leave (hover mode). |
| `target` | `String` | `null` | ID of the DOM element to use as content source. |
| `onShow` | `Function` | `null` | Callback function triggered when tooltip is shown. Receives content element. |

## Examples

### 1. Specific Position (Bottom)

```html
<div x-data="tooltip('Bottom Tip', { position: 'bottom' })">
   <!-- ... -->
</div>
```

### 2. Click Trigger + Light Theme (Popover Style)

```html
<div x-data="tooltip('Light Theme Popover', { trigger: 'click', theme: 'light' })">
    <button data-trigger @click="onClick" @click.outside="onClickOutside">Click Me</button>
    <!-- ... -->
</div>
```

### 3. HTML Content (Target Mode)

```html
<!-- Hidden Content Source -->
<div id="tip-content" style="display: none;">
    <strong>Bold Title</strong><br>
    <span class="text-blue-500">Supports HTML styles</span>
</div>

<!-- Tooltip Instance -->
<div x-data="tooltip('', { target: 'tip-content' })">
    <button data-trigger ...>Hover to see HTML</button>
    <!-- ... -->
</div>
```

## Developer Notes

- **Core Positioning**: Positioning logic reuses `src/core/positioning.js`, consistent with Box/Notify components.
- **Debugging**: Enable debug mode in `src/core/log.js` to view positioning calculations in the console.
