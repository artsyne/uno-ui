# Uno Tooltip API Documentation

The Tooltip component provides smart, positioning-aware tooltips. It supports auto-flipping, different trigger modes (hover, click, focus), and can load content dynamically via HTMX.

## Overview

- **Core Function**: Display contextual information attached to an element.
- **Architecture**: Directive (`x-uno-tooltip`) registers triggers; Store (`Alpine.store('uno').tooltip`) manages visibility and position.
- **Positioning**: Uses `requestAnimationFrame` for high-performance, collision-aware positioning.
- **Interaction**: Supports Hover (default), Click (Popover), and Focus triggers.

## Basic Usage

### 1. Simple Text Tooltip

```html
<button x-uno-tooltip="'Save changes'">Save</button>
```

### 2. Configuration Object

```html
<button
  x-uno-tooltip="{ 
  content: '<strong>Bold</strong> info', 
  placement: 'right',
  theme: 'light' 
}"
>
  Info
</button>
```

### 3. DOM Element as Content

```html
<div id="my-tip" class="hidden">Detailed info...</div>
<button x-uno-tooltip="{ target: '#my-tip', trigger: 'click' }">Show Details</button>
```

## API Parameters

These options are passed in the configuration object to `x-uno-tooltip`.

| Name           | Type      | Default   | Required | Description                                                                  |
| :------------- | :-------- | :-------- | :------- | :--------------------------------------------------------------------------- |
| `content`      | `String`  | `''`      | No\*     | Text/HTML content. If starting with `#` (e.g., `'#id'`), it acts as `target`.|
| `target`       | `String`  | -         | No\*     | CSS Selector to an element whose `innerHTML` will be used as content.        |
| `placement`    | `String`  | `'top'`   | No       | Preferred position. Values: `'top'`, `'bottom'`, `'left'`, `'right'`.        |
| `trigger`      | `String`  | `'hover'` | No       | Trigger mode. Values: `'hover'`, `'click'`, `'focus'`.                       |
| `arrow`        | `Boolean` | `true`    | No       | Whether to show the directional arrow.                                       |
| `offset`       | `Number`  | `8`       | No       | Distance in pixels between trigger and tooltip.                              |
| `clickOutside` | `Boolean` | `true`    | No       | If `true`, clicking outside closes the tooltip (for `click` mode).           |
| `hxUrl`        | `String`  | -         | No       | URL to fetch content from via HTMX.                                          |
| `hxTrigger`    | `String`  | `'load'`  | No       | When to trigger HTMX fetch. Default `'load'` means fetch when tooltip shows. |

## Events

The directive handles native DOM events on the trigger element based on the `trigger` mode.

| Mode    | Events Listened            | Behavior                                     |
| :------ | :------------------------- | :------------------------------------------- |
| `hover` | `mouseenter`, `mouseleave` | Shows on enter, hides on leave (with delay). |
| `click` | `click`                    | Toggles visibility.                          |
| `focus` | `focus`, `blur`            | Shows on focus, hides on blur.               |

## Methods

Accessed via `Alpine.store('uno').tooltip`.

### `show(options)`

Manually shows a tooltip.

- **Parameters**: `options` (Object)
  - `element` (HTMLElement, Required): The anchor element.
  - `content` (String): Content to show.
  - `placement` (String): Position preference.
- **Returns**: `void`

### `hide()`

Hides the currently visible tooltip.

- **Parameters**: None.
- **Returns**: `void`

### `updatePosition()`

Force recalculates the tooltip position. Useful if the anchor element moves or resizes.

- **Parameters**: None.
- **Returns**: `void`

## Examples

### HTMX Loaded Tooltip

Fetch user details when hovering over a username.

```html
<span
  x-uno-tooltip="{ 
  hxUrl: '/api/user/123/card', 
  hxTrigger: 'mouseenter',
  placement: 'bottom'
}"
>
  @username
</span>
```

### Click-Triggered Popover (Menu)

```html
<button
  x-uno-tooltip="{ 
  target: '#menu-content', 
  trigger: 'click', 
  placement: 'bottom-start',
  arrow: false,
  offset: 4
}"
>
  Options
</button>

<div id="menu-content" class="hidden">
  <ul class="menu">
    <li>Edit</li>
    <li>Delete</li>
  </ul>
</div>
```

## Version History

| Version | Date       | Changes                                                  |
| :------ | :--------- | :------------------------------------------------------- |
| v1.0.0  | 2025-01-09 | Initial release with Smart Positioning and HTMX support. |
