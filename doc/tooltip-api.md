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
  placement: 'right'
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

| Name           | Type      | Default   | Required | Description                                                                   |
| :------------- | :-------- | :-------- | :------- | :---------------------------------------------------------------------------- |
| `content`      | `String`  | `''`      | No\*     | Text/HTML content. If starting with `#` (e.g., `'#id'`), it acts as `target`. |
| `target`       | `String`  | -         | No\*     | CSS Selector to an element whose `innerHTML` will be used as content.         |
| `placement`    | `String`  | `'top'`   | No       | Preferred position. Values: `'top'`, `'bottom'`, `'left'`, `'right'`.         |
| `trigger`      | `String`  | `'hover'` | No       | Trigger mode. Values: `'hover'`, `'click'`, `'focus'`.                        |
| `arrow`        | `Boolean` | `true`    | No       | Whether to show the directional arrow.                                        |
| `offset`       | `Number`  | `8`       | No       | Distance in pixels between trigger and tooltip.                               |
| `clickOutside` | `Boolean` | `true`    | No       | If `true`, clicking outside closes the tooltip (for `click` mode).            |
| `hxUrl`        | `String`  | -         | No       | URL to fetch content from via HTMX.                                           |
| `hxTrigger`    | `String`  | `'load'`  | No       | When to trigger HTMX fetch. Default `'load'` means fetch when tooltip shows.  |

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

### `cancelHide()`

Cancels a scheduled hide operation. Use this when you want to prevent the tooltip from hiding after `scheduleHide()` was called.

- **Parameters**: None.
- **Returns**: `void`
- **Example**:
  ```javascript
  Alpine.store('uno').tooltip.cancelHide()
  ```

### `scheduleHide(delay = 300)`

Schedules the tooltip to hide after a delay. This is useful for creating a delay before hiding when the mouse leaves the trigger.

- **Parameters**:
  - `delay` (Number, default: `300`) - Delay in milliseconds before hiding.
- **Returns**: `void`
- **Example**:
  ```javascript
  Alpine.store('uno').tooltip.scheduleHide(500) // Hide after 500ms
  ```

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
  placement: 'bottom',
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

## Troubleshooting

### Common Issues

| Issue                   | Possible Cause                                  | Solution                                                                                                      |
| :---------------------- | :---------------------------------------------- | :------------------------------------------------------------------------------------------------------------ |
| **Tooltip not showing** | `z-index` too low or overflow hidden on parent. | Tooltips are appended to body. Override the layer via `Uno.setConfig({ ui: { layers: { tooltip: 3000 } } })`. |
| **Wrong position**      | Element near edge or moving.                    | Call `Alpine.store('uno').tooltip.updatePosition()` if anchor moves.                                          |
| **Flickering**          | Hover logic conflict.                           | Ensure `offset` is sufficient to prevent mouse leaving trigger area.                                          |

## Version History

| Version | Date       | Changes                                                                                                                                                                                                                          |
| :------ | :--------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| v0.9.0  | 2026-08-26 | `offset` option is now honored by `calculatePosition()` (previously ignored); the arrow element is rendered from the returned anchor. `placement` documented/typed as the four supported values (`top`/`bottom`/`left`/`right`). |
| v0.8.0  | 2026-02-28 | Added `bindNotifyToHtmx()` method to HTMX integration for automatic notifications.                                                                                                                                               |
| v0.7.0  | 2026-01-20 | Refactored for **Lazy Initialization**.                                                                                                                                                                                          |
| v0.6.0  | 2025-12-18 | Initial stable release.                                                                                                                                                                                                          |
