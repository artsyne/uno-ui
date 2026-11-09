# Tooltip API

Tooltip attaches a positioned popup to an element. It supports hover, click and focus triggers,
auto-flipping and viewport clamping, an optional arrow, HTMX content, and a keyboard contract
for click popups.

## Prerequisites

- Alpine.js 3 with Uno registered. See [Getting Started](getting-started.md).
- Tailwind CSS 4 classes generated for Uno's templates. See [Theming](theming.md).
- HTMX 2 only if you use `hxUrl`.

## Minimal example

```html
<button x-data x-uno-tooltip="'Save changes'">Save</button>
```

## Directive

`x-uno-tooltip` accepts a string (content) or an object:

```html
<button
  x-data
  x-uno-tooltip="{
    content: '<strong>Bold</strong> info',
    placement: 'right'
  }"
>
  Info
</button>
```

The directive attaches `mouseenter`, `mouseleave`, `click`, `focus` and `blur` listeners and
filters them by the resolved `trigger`. Modifiers are not supported; a modifier logs a
deprecation warning and is ignored.

### DOM element as content

```html
<div id="my-tip" class="hidden">Detailed info...</div>
<button x-data x-uno-tooltip="{ target: '#my-tip', trigger: 'click' }">Details</button>
```

`target` uses the element's `innerHTML`. A `content` string that starts with `#` is treated as
a target when no `target` is set and a matching element exists.

## Store API

```javascript
const tooltip = Alpine.store('uno').tooltip
```

### `show(options)`

Shows a tooltip programmatically. `element` is required when called from the store; the
directive supplies it.

```javascript
tooltip.show({
  element: document.querySelector('#help'),
  content: 'Help text',
  placement: 'top'
})
```

Returns `void`. The store is a singleton: showing a new tooltip replaces the current one.

### `hide()`

Hides the current tooltip and restores focus for click popups when focus was inside the popup
or on `body`.

### `updatePosition()`

Recalculates position and arrow placement. Positioning already runs on open, scroll, resize and
popup size changes; call this after moving the trigger yourself.

### `cancelHide()`

Cancels a scheduled hide.

### `scheduleHide(delay = 300)`

Schedules `hide()` after `delay` milliseconds. Mouseleave on the trigger uses the default; the
popup's own `mouseenter` cancels it.

## Options

| Name | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `content` | `string` | `''` | Text/HTML content. If it starts with `#` and no `target` is set, it may be used as a target selector. |
| `target` | `string` | — | CSS selector whose `innerHTML` becomes the content. |
| `placement` | `'top' \| 'bottom' \| 'left' \| 'right'` | `'top'` | Preferred side; the popup flips when there is no room. Other strings resolve to `bottom`. |
| `trigger` | `'hover' \| 'click' \| 'focus'` | `'hover'` | Interaction mode. |
| `arrow` | `boolean` | `true` | Render the 8 px rotated arrow. |
| `offset` | `number \| (ctx) => number` | `8` | Gap in pixels. The function form is re-evaluated on every positioning pass. |
| `clickOutside` | `boolean` | `true` | Close a click popup on an outside click. |
| `hxUrl` | `string` | — | Load content with HTMX. |
| `hxHeaders` | `Record<string, string>` | — | Serialized into `hx-headers`. |
| `hxTrigger` | `string` | `'load'` | HTMX trigger attribute. |
| `element` | `HTMLElement` | — | Anchor element for `show()` from the store. |

### `offset` as a function

```javascript
const offset = ({ triggerEl, tooltipEl, placement }) => {
  const rect = triggerEl.getBoundingClientRect()
  return placement === 'bottom' ? rect.height / 2 + 8 : 8
}

Alpine.store('uno').tooltip.show({
  element: document.querySelector('#anchor'),
  content: 'Live gap',
  offset
})
```

The function receives `{ triggerEl, tooltipEl, placement }` and must return a finite number.
A thrown error or a non-numeric result falls back to the default 8 px gap. It runs on open,
scroll, resize and popup size changes.

## Events

| Behavior | Native events | Notes |
| :--- | :--- | :--- |
| `hover` | `mouseenter`, `mouseleave` | Popup enter cancels hide; popup leave schedules hide. |
| `click` | `click` | Toggles. `preventDefault()` is called. Outside click closes when `clickOutside` is true. |
| `focus` | `focus`, `blur` | Shows on focus and hides on blur. |

Escape closes the popup while it is visible.

### Keyboard contract for click popups

With `trigger: 'click'`:

- Trigger gets `aria-haspopup="true"`, `aria-expanded="true"` and
  `aria-describedby="uno-tooltip-content"`.
- After render, focus moves to the first focusable element in the popup when one exists.
- Tab is trapped inside the popup.
- ArrowDown/ArrowUp/Home/End move between focusable items **only when the content contains an
  element with `role="menu"`**.
- Escape, outside click, an action, or Space on an anchor trigger closes it. `aria-expanded`
  becomes `"false"` and `aria-describedby` is removed.
- Focus returns to the trigger when it was inside the popup (or on `body`).

Hover and focus modes remain hints: no focus movement, no trap, no ARIA state changes.

Menu content should declare APG semantics:

```html
<div id="menu-content" class="hidden">
  <ul role="menu">
    <li role="none"><button role="menuitem" type="button">Edit</button></li>
    <li role="none"><button role="menuitem" type="button">Delete</button></li>
  </ul>
</div>
```

Known limitation: the popup keeps `role="tooltip"` in click mode, even when it behaves like a
disclosure or menu. This is not a full APG menu pattern.

## Positioning

- Preferred placement is used when the popup fits; otherwise it flips to the opposite side.
- The box is clamped at least 10 px inside the viewport on both axes. A popup wider than the
  viewport stays at the left safe margin instead of producing a negative coordinate.
- The popup size is observed with `ResizeObserver` after opening; late content, fonts and
  language switches re-run positioning.
- `left`/`right` placements are mirrored for RTL anchors based on the computed `direction` of
  the trigger.

The default content class uses `whitespace-nowrap`. For menus or paragraphs, override
`ui.tooltip.content` and add a max width so a wide box still clamps correctly:

```javascript
Uno.setConfig({
  ui: {
    tooltip: {
      content: '+max-w-[calc(100vw-1rem)] whitespace-normal'
    }
  }
})
```

## HTMX content

```html
<span
  x-data
  x-uno-tooltip="{
    hxUrl: '/api/user/123/card',
    hxTrigger: 'load',
    placement: 'bottom'
  }"
>
  @username
</span>
```

HTMX inserts the response into the popup. It is not sanitized; treat the server response as
trusted. `hxTrigger` defaults to `'load'`, so the request happens when the popup renders.

## Accessibility

- `role="tooltip"` on the popup.
- Click mode adds the ARIA state and focus behaviour described above.
- Hover/focus mode does not add `aria-describedby`; screen readers may not announce a
  hover-only tooltip. Provide the same information in visible text when it is essential.

## Notes

- The tooltip store is a singleton; opening a second tooltip hides the first.
- `ui.tooltip.content` contains `{base}` and defaults to `whitespace-nowrap`.
- `scheduleHide()` defaults to 300 ms. The registered plugin default `delay` is not currently
  read by the directive.
- The popup container is appended to the document body.

## Version history

| Version | Date | Changes |
| :--- | :--- | :--- |
| 0.10.0 | 2026-11-09 | Click popup keyboard contract, popup size observation, non-negative horizontal clamping, function `offset`. |
| 0.9.0 | 2026-08-26 | `offset` honored by positioning; arrow rendered; `updatePosition()` reads `position`; content sanitized; placement narrowed to four values. |
| 0.7.0 | 2026-01-20 | Lazy initialization. |
| 0.6.0 | 2025-12-18 | Initial release. |

## Related

- [Configuration](configuration.md)
- [Theming](theming.md)
- [HTMX integration](htmx.md)
- [Accessibility](accessibility.md)
