# Notify API

Notify shows toast-style messages in nine screen positions, with optional HTMX content and a
global progress bar. It is store-first (`Alpine.store('uno').notify`) and has a directive
(`x-uno-notify`) for click-driven use.

## Prerequisites

- Alpine.js 3 with Uno registered. See [Getting Started](getting-started.md).
- Tailwind CSS 4 classes generated for Uno's templates. See [Theming](theming.md).
- HTMX 2 only if you use `hxUrl` or the progress helpers with HTMX.

## Minimal example

```html
<button x-data @click="$store.uno.notify.show('Saved', { type: 'success' })">
  Show notification
</button>
```

Directive form:

```html
<button x-data x-uno-notify="{ message: 'Saved', type: 'success' }">Show notification</button>
```

## Directive

`x-uno-notify` accepts a string (the message) or an object with a `message` key plus any
option:

```html
<button x-data x-uno-notify="'Saved'">Show</button>
<button x-data x-uno-notify="{ message: 'Failed', type: 'error', duration: 5000 }">Show</button>
```

The directive listens for `click` and calls `notify.show(message, config)`.

## Store API

```javascript
const notify = Alpine.store('uno').notify
```

### `show(message, options)`

Displays a notification and returns its numeric id.

```javascript
const id = notify.show('Saved', { type: 'success', duration: 2000 })
notify.remove(id)
```

`message` is passed through `i18n.t(message)`. A normal sentence that is not a translation key
comes back unchanged; a key is translated.

### Shortcuts

| Method | Equivalent |
| :--- | :--- |
| `success(message, options)` | `show(message, { type: 'success', ...options })` |
| `error(message, options)` | `show(message, { type: 'error', ...options })` |
| `warning(message, options)` | `show(message, { type: 'warning', ...options })` |
| `info(message, options)` | `show(message, { type: 'info', ...options })` |

Each returns the numeric id.

### `remove(id)`

Hides the notification (`isVisible = false`) and removes it after a 300 ms transition window.
The id may belong to any position group.

### Progress bar: counted

Use this pair when several operations can overlap:

```javascript
notify.beginProgress()
try {
  await loadData()
} finally {
  notify.endProgress()
}
```

| Method | Behavior |
| :--- | :--- |
| `beginProgress()` | Increments the active count; starts the bar on the first request. |
| `endProgress()` | Decrements the count; stops the bar when it reaches zero. Safe to call extra times. |

`progressRequests` is exposed on the store if you need to inspect the count.

### Progress bar: direct

```javascript
notify.startProgress()
// ...
notify.stopProgress()
```

| Method | Behavior |
| :--- | :--- |
| `startProgress()` | Shows the bar and advances it by 10% every 100 ms, resetting at 100%. |
| `stopProgress()` | Stops the timer, fills to 100%, hides after 300 ms and resets the percentage. |

These do not count concurrency. Prefer `beginProgress()` / `endProgress()` when requests can
overlap.

## Options

| Name | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `message` | `string` | required | Text or i18n key. Rendered with `x-text`. |
| `type` | `'info' \| 'success' \| 'warning' \| 'error'` | none | Controls the background, icon colour and announcement politeness. There is no default; pass a type for a styled notification. |
| `duration` | `number` | `3000` | Auto-dismiss after this many milliseconds. `0` disables auto-dismiss. |
| `position` | `NotifyPosition` | `'top-right'` | Screen position. Unknown values fall back to the `top-right` group. |
| `closable` | `boolean` | `true` | Show the close button. |
| `icon` | `string` | — | CSS class string for an `<i>` element, for example `fas fa-check`. The icon is only rendered when set. |
| `hxUrl` | `string` | — | Load content with HTMX into the notification. |
| `hxTrigger` | `string` | `'load'` | HTMX trigger attribute. |
| `hxSwap` | `string` | `'innerHTML'` | HTMX swap strategy. |
| `hxHeaders` | `Record<string, string>` | — | Serialized into `hx-headers`. |
| `onContentLoad` | `(el: HTMLElement) => void` | — | Called after HTMX content is copied into the item. |
| `onError` | `(detail: any) => void` | — | Called on HTMX response/send/timeout errors. |

`NotifyPosition`: `top-left`, `top-center`, `top-right`, `center-left`, `center`, `center-right`,
`bottom-left`, `bottom-center`, `bottom-right`.

```javascript
notify.show('Please verify your email.', {
  type: 'warning',
  duration: 0,
  closable: true
})
```

## HTMX content

When `hxUrl` is set, the notification body is replaced by an HTMX request. The response is
copied into `item.content` and rendered with `x-html`; it is not sanitized. Treat the server
response as trusted HTML.

Notify calls `beginProgress()` / `endProgress()` around each HTMX notification request. A
failed load calls your `onError(detail)` and also shows a generic `Content load failed` toast.

See [HTMX integration](htmx.md#notify).

## Automatic HTMX notifications

`bindNotifyToHtmx()` lives on `Uno.htmx`, not on the notify store:

```javascript
Uno.htmx.bindNotifyToHtmx({
  showSuccess: true,
  showError: true,
  showProgress: true
})
```

It binds HTMX request events to progress/success/error notifications and can only be bound once.
See [HTMX integration](htmx.md#bindnotifytohtmx) for options and server headers.

## Accessibility

| Type | Role | `aria-live` |
| :--- | :--- | :--- |
| `error`, `warning` | `alert` | `assertive` |
| `success`, `info` | `status` | `polite` |

The role and live-region attributes are bound per item in the template. The close button is an
icon-only button with a visually hidden label from the `box.close` i18n key. Notifications do
not steal focus.

The global progress bar is visual only; it has no `progressbar` role in the current template.

## Notes

- The notification container is created lazily on the first `show()`.
- The container's z-index comes from `ui.layers.notify` (default 1600); the progress bar uses
  9999.
- `duration` is measured from insertion; the leave transition adds roughly 300 ms before the
  element is removed.
- The progress bar container is created by `show()`. If you call `startProgress()` (or
  `bindNotifyToHtmx()` progress) before any notification has been shown, call
  `notify.ensureContainer()` once or show a notification first.
- `message` is text, not HTML. Use `hxUrl` for server-rendered markup.

## Version history

| Version | Date | Changes |
| :--- | :--- | :--- |
| 0.10.0 | 2026-11-09 | Announcement politeness graded by type; close-button name uses Uno i18n. |
| 0.9.0 | 2026-08-26 | Added `beginProgress()` / `endProgress()`; fixed `remove()` hiding via `isVisible`. |
| 0.8.0 | 2026-01-22 | Added `bindNotifyToHtmx()` on `Uno.htmx`. |
| 0.7.0 | 2026-01-20 | Lazy initialization; containers on first use. |
| 0.6.0 | 2025-12-18 | Initial release. |

## Related

- [HTMX integration](htmx.md)
- [Configuration](configuration.md)
- [Theming](theming.md)
- [Accessibility](accessibility.md)
