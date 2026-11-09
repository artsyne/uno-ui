# Box API

Box provides modal dialogs, alerts, confirms and prompts. It is store-first
(`Alpine.store('uno').box`) with a directive (`x-uno-box`) for click-driven use. Dialogs can
load HTMX content and stack; each dialog has a focus trap and restores focus to its opener.

## Prerequisites

- Alpine.js 3 with Uno registered. See [Getting Started](getting-started.md).
- Tailwind CSS 4 classes generated for Uno's templates. See [Theming](theming.md).
- HTMX 2 only if you use `hxUrl`.

## Minimal example

```html
<button
  x-data
  @click="$store.uno.box.open({ title: 'Profile', content: '<p>Hello</p>' })"
>
  Open dialog
</button>
```

Equivalent directive form:

```html
<button x-data x-uno-box="{ title: 'Profile', content: '<p>Hello</p>' }">Open dialog</button>
```

## Directive

`x-uno-box` evaluates its expression to an object and calls `box.open(object)` on click.

```html
<button
  x-data
  x-uno-box="{
    type: 'confirm',
    title: 'Delete item?',
    content: 'This action cannot be undone.',
    okText: 'Delete',
    okButtonClass: '!bg-red-600 hover:!bg-red-700',
    onOk: () => deleteItem()
  }"
>
  Delete
</button>
```

The directive calls `preventDefault()` and logs a warning if the expression is not an object.
The expression is re-evaluated on each click, so it can read Alpine state.

## Store API

```javascript
const box = Alpine.store('uno').box
```

### `open(options)` / `open(type, options)`

Opens a dialog and returns the reactive item.

```javascript
const item = box.open({ title: 'Info', content: 'Text' })
box.open('confirm', { title: 'Continue?', content: 'Proceed?' })
```

The item contains a `close()` method and an `isOpen` flag. Use `close(id)` or the returned item.

### `alert(message, options)`

Opens a `dialog` with an OK button and the message as `content`. If `options` is a string, it is
treated as the title.

```javascript
box.alert('Saved')
box.alert('Saved', { title: 'Success' })
box.alert('Saved', 'Success') // string shorthand for title
```

Returns the item.

### `confirm(message, options)`

Returns a Promise that resolves to `true` on OK and `false` on Cancel or close.

```javascript
const ok = await box.confirm('Delete this item?', { title: 'Confirm' })
if (ok) {
  // ...
}
```

User callbacks are composed with the resolver:

- `onOk` returning `false` (or a Promise resolving to `false`) keeps the dialog open and the
  Promise pending.
- `onCancel` runs first, then the Promise resolves `false`.
- `onClose` runs on close and resolves `false` if the Promise is still pending.

### `prompt(message, defaultValue, options)`

Returns a Promise resolving to the input value, or `null` on Cancel.

```javascript
const name = await box.prompt('Your name:', 'Guest', { title: 'Welcome' })
```

The message and default value are passed through `Uno.dom.sanitize()` before being inserted into
the prompt markup. `prompt()` does **not** compose an `onClose` resolver: closing through the
close button, Escape or the overlay (when enabled) can leave the Promise pending. Use the Cancel
button or call the returned item's `close()` only after handling that state yourself.

### `dialog(options)`

Alias for `open('dialog', options)`.

### `handleOk(item)` / `handleExtraButton(item, btn)`

Internal handlers used by the template. They run the callback, await a Promise, and close the
dialog unless the result is `false`. They are public on the store but most applications do not
need to call them.

### `close(id)`

Closes the dialog with that id.

```javascript
const item = box.open({ title: 'Info' })
box.close(item.id)
```

### `closeAll()`

Calls `onBeforeClose` and `onClose` for every item, then clears the item array immediately. It
does not animate or restore focus. Prefer closing individual items when focus management
matters.

### `getPositionClasses(position)`

Returns the flex alignment classes for a 9-point position. Mainly used by the template.

## Options

| Name | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `type` | `'dialog' \| 'confirm'` | `'dialog'` | Dialog type. Other strings are accepted but do not add button defaults. |
| `title` | `string` | `''` | Header text. Also used as the dialog label when present. |
| `content` | `string` | `''` | HTML body. Passed through `Uno.dom.sanitize()`. |
| `width` | `number` | `400` | Width in pixels. The 400 default applies only when both `width` and `widthClass` are absent. |
| `widthClass` | `string` | — | Tailwind classes for width, applied alongside any inline width. |
| `height` | `number` | — | Height in pixels. No default. |
| `heightClass` | `string` | — | Tailwind classes for height, applied alongside any inline height. |
| `containerClass` | `string` | — | Extra classes appended to the container. |
| `boxId` | `string` | random | Explicit item id. |
| `fill` | `boolean` | `false` | Fill mode; see [Fill mode](#fill-mode). |
| `position` | `BoxPosition` | `'center'` | One of the 9 positions below. |
| `modal` | `boolean` | `true` | Sets the `aria-modal` attribute. The backdrop is rendered whenever an item exists; use `closeOnOverlayClick` for dismissal. |
| `showCloseButton` | `boolean` | `true` | Show the floating close button. |
| `showOkButton` | `boolean` | `false` (`true` for `confirm`) | Show the OK button. |
| `showCancelButton` | `boolean` | `false` (`true` for `confirm`) | Show the Cancel button. |
| `okText` | `string` | `'OK'` | OK label. Rendered literally, not translated. |
| `cancelText` | `string` | `'Cancel'` | Cancel label. Rendered literally, not translated. |
| `okButtonClass` | `string` | — | Extra classes for the OK button. |
| `cancelButtonClass` | `string` | — | Extra classes for the Cancel button. |
| `extraButtons` | `BoxExtraButton[]` | `[]` | Additional buttons; see below. |
| `closeOnEsc` | `boolean` | `false` | Close on Escape. |
| `closeOnOverlayClick` | `boolean` | `false` | Close on backdrop/dialog-wrapper click. |
| `autoFocus` | `boolean` | `true` | Focus the first inputable element. The neutral target is still focused when disabled. |
| `hxUrl` | `string` | — | Load the body with HTMX. |
| `hxMethod` | `'GET' \| 'POST'` | `'GET'` | HTMX method; other values are not rendered by the template. |
| `hxHeaders` | `Record<string, string>` | — | Headers serialized into `hx-headers`. |
| `hxTrigger` | `string` | `'load'` | HTMX trigger. |
| `onOk` | `() => void \| false \| Promise<void \| false>` | — | OK callback. `false` keeps the dialog open. |
| `onCancel` | `() => void \| false` | — | Cancel callback. Returning `false` also keeps the dialog open. |
| `onBeforeClose` | `() => void` | — | Runs before the close animation; returning `false` does not stop it. |
| `onClose` | `() => void` | — | Runs after the close animation/removal. |
| `onContentLoad` | `(el: HTMLElement) => void` | — | Runs after HTMX content is swapped in. |

`BoxPosition` values: `center`, `top-left`, `top-center`, `top-right`, `center-left`,
`center-right`, `bottom-left`, `bottom-center`, `bottom-right`.

`defaultWidth` and `defaultPosition` are plugin-level defaults, not per-instance options; they are
read when `width` / `position` are absent. See [Configuration](configuration.md#plugin-defaults).

## Extra buttons

| Name | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `text` | `string` | required | Visible label. |
| `position` | `'before' \| 'after'` | `'before'` | `before` renders to the visual left of OK; `after` renders after Cancel. |
| `buttonClass` | `string` | — | Replaces the default `extraBtn` classes. |
| `onClick` | `(item: BoxItem) => void \| false \| Promise<void \| false>` | — | Callback. `false` keeps the dialog open. |

```javascript
const item = box.open({
  title: 'New entry',
  showOkButton: true,
  okText: 'Publish',
  extraButtons: [
    {
      text: 'Save draft',
      position: 'before',
      onClick: async () => {
        await saveDraft()
        return false // keep open
      }
    }
  ],
  onOk: async () => {
    await publish()
  }
})

// Replace the array later; Alpine's proxy keeps the template reactive.
item.extraButtons = [{ text: 'Archive', position: 'after', onClick: () => archive() }]
```

## Events and callbacks

| Callback | Timing | Arguments |
| :--- | :--- | :--- |
| `onOk` | OK button click | none |
| `onCancel` | Cancel button click | none |
| `onBeforeClose` | Before closing, before the leave transition | none |
| `onClose` | After the leave transition and removal | none |
| `onContentLoad` | After an HTMX swap, before autofocus | loaded element |

`onOk` and `onCancel` are only invoked by their buttons. Closing through the close button,
Escape or the overlay calls `onClose` (and `onBeforeClose`), not `onCancel`.

## Accessibility

- `role="dialog"`, `:aria-modal`, `:aria-labelledby` (title) and `:aria-describedby`
  (content container).
- Initial focus goes to the first visible input/textarea/select when `autoFocus` is `true`;
  otherwise to the title or container. Hidden, disabled and readonly inputs are skipped, as are
  buttons.
- `x-uno-trap` traps Tab/Shift+Tab and wraps in DOM order.
- Closing restores focus to the opener when focus is still inside the dialog or on `body`.
- Underlay dialogs receive `inert` and `aria-hidden="true"`.
- The floating close button gets its accessible name from `box.close`.

See [Accessibility](accessibility.md) for test scope and limitations.

## Fill mode

`fill: true` makes the dialog body a flex column so loaded content can take the remaining
height:

```javascript
Alpine.store('uno').box.open({
  title: 'Editor',
  hxUrl: '/editor',
  fill: true,
  widthClass: 'w-full lg:w-[70%]'
})
```

- The content area, body and HTMX container receive `contentFill` / `bodyFill`.
- Without `height`/`heightClass`, `containerFill` supplies the height.
- The loaded content must continue the chain: each wrapper down to the growing element needs
  `flex-1 min-h-0 flex flex-col`, and the growing element needs a minimum height.
- If the minimum height exceeds the dialog, the content area scrolls as one container.

## HTMX

See [HTMX integration](htmx.md#box). The response is inserted as HTML by HTMX and is not
sanitized; treat it as trusted server output. Box has no `onError`; use HTMX events.

## Notes

- `content` is sanitized with a regex-based sanitizer; it is not a general-purpose HTML
  sanitizer. Do not pass untrusted HTML.
- The backdrop is always rendered while an item exists. `modal: false` only changes
  `aria-modal`; it does not remove the backdrop or make the page interactive.
- Stacked dialogs add `index * 10` to `ui.layers.box`.
- `closeAll()` is immediate and does not run the leave transition.

## Related

- [Configuration](configuration.md)
- [Theming](theming.md)
- [HTMX integration](htmx.md)
- [Accessibility](accessibility.md)
- [Notify API](notify-api.md)
