# Box Component (Dialog)

A reusable dialog component based on Alpine.js + TailwindCSS, featuring modal overlays, flexible positioning, size control, and lazy content loading.

## Features

- **Modal & Overlay**: Supports `modal` mode, automatically locks body scroll when open, and restores it on close. Stackable modals are supported.
- **Flexible Positioning**: Supports 9 positions: `center`, `top`, `bottom`, `top-right`, `bottom-right`, `top-left`, `bottom-left`, `top-center`, `bottom-center`.
- **Size & Responsiveness**:
  - Direct `width`/`height` settings.
  - Custom classes via `widthClass`/`heightClass`/`containerClass`.
  - `contentOverflow` controls scroll behavior within the content area.
- **Smart Caching**: Set `cache=true` + `boxId` to reuse DOM elements. Closing only hides the element. Supports manual cache clearing.
- **Content Loading**:
  - `content`: Static text or HTML.
  - `targetId`: Inject existing DOM elements from the page.
  - `ajaxUrl`: Async content loading (with built-in Loading state).
  - `hxUrl`: HTMX content loading (uses `htmx.ajax`).
  - `iframeUrl`: Embed an Iframe.
- **Buttons & Interaction**:
  - Three presets: `dialog`, `confirm`, `message`.
  - Button text supports i18n.
  - Lifecycle events: `onBeforeOpen`, `onOpen`, `onBeforeClose`, `onClose`.
  - Button alignment: `buttonAlign` (left/center/right).
- **Declarative Usage**: Quickly use via `data-box-type` and `data-box-*` attributes without writing JS.
- **Internationalization**: Supports `Alpine.store('i18n')` for text (`box.ok`, `box.cancel`, `box.loading`, `box.close`).

## Quick Start

### Declarative Configuration (HTML Attributes)

The simplest way is to use `data-box-*` attributes directly on a button to trigger the dialog.

```html
<button
  data-box-type="confirm"
  data-box-title="Confirm Action"
  data-box-content="Are you sure you want to proceed?"
  data-box-position="center"
  data-box-button-align="right"
  class="btn btn-primary"
>
  Open Confirm
</button>
```

## Programmatic Usage

You can invoke the dialog programmatically via `Alpine.store('box')`.

```javascript
// Open Dialog
Alpine.store('box').dialog({
  title: 'System Notice',
  content: 'Operation Successful',
  position: 'top-center',
  width: '400px',
  // Callback events
  onClose: () => console.log('Closed')
})

// Clear Cache
Alpine.store('box').clearCache()
```

## Configuration / Options

### Data Attributes (Declarative)

| Attribute              | Description               | Default / Options                                                                           |
| ---------------------- | ------------------------- | ------------------------------------------------------------------------------------------- |
| `data-box-type`        | Dialog type               | `dialog`, `confirm`, `message`                                                              |
| `data-box-title`       | Title text                | -                                                                                           |
| `data-box-content`     | Content text/HTML         | -                                                                                           |
| `data-box-modal`       | Show overlay              | `true`                                                                                      |
| `data-box-position`    | Screen position           | `center` (default), `top`, `bottom`, `top-right`, `bottom-right`, `top-left`, `bottom-left` |
| `data-box-width`       | Width (px/rem)            | -                                                                                           |
| `data-box-height`      | Height (px/rem)           | Auto-adaptive                                                                               |
| `data-box-ajax-url`    | AJAX URL                  | -                                                                                           |
| `data-box-hx-url`      | HTMX URL (uses htmx.ajax) | -                                                                                           |
| `data-box-iframe-url`  | Iframe URL                | -                                                                                           |
| `data-box-target-id`   | Content Source DOM ID     | -                                                                                           |
| `data-box-ok-text`     | Confirm button text       | -                                                                                           |
| `data-box-cancel-text` | Cancel button text        | -                                                                                           |
| `data-box-cache`       | Cache DOM                 | `false` (Requires `data-box-id`)                                                            |
| `data-box-id`          | Unique Identifier         | -                                                                                           |

## Internationalization (i18n)

The component relies on `Alpine.store('i18n')` for multi-language support.

- **Keys**:
  - `box.ok`: Confirm button text.
  - `box.cancel`: Cancel button text.
  - `box.loading`: Loading text.
  - `box.close`: Close button aria-label.

## Events & API

### Lifecycle Events

When using programmatic usage, you can pass these callbacks in the options object:

- `onBeforeOpen`: Fired before the dialog opens.
- `onOpen`: Fired after the dialog is fully open.
- `onBeforeClose`: Fired before the dialog closes.
- `onClose`: Fired after the dialog is closed.

### Notes

- **Height Adaptation**: Height is auto-adaptive if not set. Flex layout is used internally to accommodate header/footer when height is fixed.
- **AJAX Loading**: Uses Fetch API. JSON responses are pretty-printed (for debugging), and HTML responses are injected directly.
