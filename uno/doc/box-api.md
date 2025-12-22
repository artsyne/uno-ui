# Box Component (Dialog)

A reusable dialog component based on Alpine.js + TailwindCSS, featuring modal overlays, flexible positioning, size control, and lazy content loading.

## Features

- **Modal & Overlay**: Supports `modal` overlay mode. Stackable dialogs are supported (z-index increases by order).
- **Flexible Positioning**: Supports 9 positions: `center`, `top-left`, `top-center`, `top-right`, `center-left`, `center-right`, `bottom-left`, `bottom-center`, `bottom-right`.
- **Size & Responsiveness**:
  - `width`/`height` numeric settings (px).
  - Custom classes via `widthClass`/`heightClass`/`containerClass`.
- **Request Caching**: Set `cache=true` for `ajaxUrl` or `hxUrl` to cache responses (GET requests).
- **Request Cancellation**: Automatically aborts pending AJAX requests when the dialog is closed.
- **Accessibility (A11y)**:
  - **Focus Trap**: Traps focus within the modal when open.
  - **ARIA Roles**: Uses `role="dialog"`, `aria-modal="true"`.
  - **Keyboard Navigation**: Close on `Escape` key.
- **Content Loading**:
  - `content`: Static text or HTML.
  - `targetId`: Inject existing DOM elements from the page.
  - `ajaxUrl`: Async content loading (with built-in Loading state).
  - `hxUrl`: HTMX content loading (uses `htmx.ajax`).
  - `iframeUrl`: Embed an Iframe.
- **Buttons & Interaction**:
  - Presets: `dialog`, `confirm` (internally also supports `message` type via `open()`).
  - Button text supports i18n fallbacks (`box.ok`, `box.cancel`).
  - Callbacks: `onOk`, `onCancel`, `onBeforeClose`, `onClose`, `onContentLoad`.

## Quick Start

### Minimal (JS)

```javascript
// Confirm
const ok = await Uno.box.confirm('Are you sure?')

// Alert
Uno.box.alert('Saved!')
```

## Programmatic Usage

You can invoke the dialog programmatically via `Uno.box` or `Alpine.store('box')`.

```javascript
// Open dialog via global API
Uno.box.dialog({
  title: 'System Notice',
  content: 'Operation Successful',
  position: 'top-center',
  width: 400,
  onClose: () => console.log('Closed'),
})

// Low-level control via store
Alpine.store('box').open('confirm', {
  content: 'Delete this item?',
  onOk: () => console.log('confirmed'),
})
```

## Configuration / Options

### Options Object (for `open()` / `Uno.box.dialog()`)

| Option | Type | Default | Description |
|---|---|---|---|
| `type` | `'dialog' \| 'confirm' \| 'message'` | `'dialog'` | Dialog type (only used when calling `open(opts)`). |
| `title` | `string` | `''` | Dialog title text. |
| `content` | `string` | `''` | Dialog content (text/HTML). |
| `trustedHTML` | `boolean` | `false` | If `false`, `content` is sanitized before render. |
| `position` | `string` | `'center'` | One of the 9 supported positions. |
| `width` | `number` | `400` | Width in px (ignored if `widthClass` set). |
| `height` | `number` | `undefined` | Height in px. |
| `widthClass` | `string` | `undefined` | Extra width class. |
| `heightClass` | `string` | `undefined` | Extra height class. |
| `containerClass` | `string` | `undefined` | Extra container class. |
| `modal` | `boolean` | `true` | Whether to show overlay. |
| `overlayOpacity` | `number` | `0.5` | Overlay opacity. |
| `showCloseButton` | `boolean` | `true` | Whether to show close button. |
| `showOkButton` | `boolean` | `false` | Show OK button for `dialog` type. |
| `showCancelButton` | `boolean` | `false` | Show Cancel button for `dialog` type. |
| `okButtonText` | `string` | i18n/`OK` | OK button text. |
| `cancelButtonText` | `string` | i18n/`Cancel` | Cancel button text. |
| `onOk` | `() => (void \| boolean \| Promise<boolean>)` | `undefined` | Return `false` to prevent close. |
| `onCancel` | `() => void` | `undefined` | Called before close when clicking cancel button. |
| `onBeforeClose` | `() => void` | `undefined` | Called before closing. |
| `onClose` | `() => void` | `undefined` | Called after closing. |
| `onContentLoad` | `(contentEl: HTMLElement \| null) => void` | `undefined` | Called after async content is inserted. |
| `cache` | `boolean` | `false` | Enable GET response cache for `ajaxUrl` / `hxUrl`. |
| `boxId` | `string` | random | Identifier of the box instance. |
| `targetId` | `string` | `undefined` | Use `innerHTML` from an existing DOM element. |
| `ajaxUrl` | `string` | `undefined` | Load content via request helper (JSON is rendered as `<pre>`). |
| `hxUrl` | `string` | `undefined` | Load HTML via request helper (prefers `htmx.ajax` if available). |
| `iframeUrl` | `string` | `undefined` | Embed an iframe (URL is escaped). |

## Security & HTMX

### Content Sanitization
All content loaded via `content`, `ajaxUrl`, or `hxUrl` is sanitized by default.
- Use `trustedHTML: true` in JS to bypass sanitization for trusted content.
- Dangerous tags (`script`, `iframe`, etc.) are stripped unless trusted.

### HTMX Integration
- **Unified Request**: `ajax-url` and `hx-url` both use internal request helpers (`requestHTML/requestData`), which prefer `htmx.ajax` if available.
- **HX-Trigger**: Responses can trigger client-side events via the `HX-Trigger` header.

## Internationalization (i18n)

The component relies on `Alpine.store('i18n')` for multi-language support.

- **Keys**:
  - `box.ok`: Confirm button text.
  - `box.cancel`: Cancel button text.

## Events & API

### Lifecycle & Button Callbacks

When using programmatic usage, you can pass these callbacks in the options object:

#### Lifecycle

- `onBeforeClose`: Fired before the dialog closes.
- `onClose`: Fired after the dialog is closed.
- `onContentLoad`: Fired after AJAX or HTMX content has been loaded and inserted into the DOM. Receives the content element as an argument. Useful for initializing scripts or HTMX on the new content.

#### Buttons

- `onOk`: Fired when the "OK" / "Confirm" button is clicked.
  - Return `false` to prevent the dialog from closing.
  - Useful for validation or async operations.
- `onCancel`: Fired when the "Cancel" button is clicked.
  - The dialog will always close after this callback (unless you manually intervene, but generally cancel means close).

### Notes

- **Height Adaptation**: Height is auto-adaptive if not set. Flex layout is used internally to accommodate header/footer when height is fixed.
- **AJAX Loading**: Uses Fetch API. JSON responses are pretty-printed (for debugging), and HTML responses are injected directly.
- **Caching**: Cache GET responses if `cache: true`.
- **Accessibility**: 
  - Focus is trapped within the modal.
  - `Escape` key closes the modal (if `modal: true`).
  - Proper ARIA roles are applied.

## Manual Testing

### 1. Callback Verification

Open the browser console and execute:

```javascript
Alpine.store('box').open({
  title: 'Callback Test',
  content: 'Close this box to check console logs',
  onBeforeClose: () => console.log('onBeforeClose triggered'),
  onClose: () => console.log('onClose triggered')
})
```

Expected:

1. Dialog opens.
2. Click Close button or Overlay.
3. Console shows "onBeforeClose triggered" then "onClose triggered".
