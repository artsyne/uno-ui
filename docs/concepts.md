# Core Concepts

Uno is small enough to read in one sitting, but five ideas explain most of its behaviour:
**Store**, **Directive**, **Template**, **Layer** and **Action**. The rest of this page ties
them to the component lifecycle and the configuration model.

## Store-first

Every component is an Alpine store under the `uno` namespace:

```javascript
Alpine.store('uno').box
Alpine.store('uno').notify
Alpine.store('uno').tooltip
Alpine.store('uno').tab
Alpine.store('uno').valid
```

The `uno` store is assembled by the plugin registry. Each plugin registers its object, and the
registry re-registers the whole store so Alpine keeps the nested properties reactive. A
compatibility alias also exposes the plugin on its own name:

```javascript
Alpine.store('box') // same object as Alpine.store('uno').box
```

Programmatic control goes through the store. This is the API you use from application code,
from tests and from the browser console.

## Directive-driven

Directives are the declarative surface in HTML:

| Directive | Component | Purpose |
| :--- | :--- | :--- |
| `x-uno-box` | Box | Open a dialog from an element click. |
| `x-uno-notify` | Notify | Show a notification on click. |
| `x-uno-tooltip` | Tooltip | Attach a hover, click or focus popup. |
| `x-uno-tab` | Tab | Register a tab trigger and wire keyboard/ARIA. |
| `x-uno-valid` | Valid | Attach validation rules to a field or mark a form. |
| `x-uno-trap` | Box | Internal focus-trap helper used by the dialog template. |

Directives read their expression, attach native DOM listeners, call the store, and clean up
through Alpine's `cleanup` hook. They do not emit custom DOM events. When an API needs to react
to state, read the store or pass a callback.

The `x-uno-tab` directive is the one exception to "user owns ARIA": it injects `role="tab"`,
`id="tab-{name}"` (when no id exists), `aria-controls`, `aria-selected` and the roving
`tabindex`, because those values are fully derived from the tab configuration.

## Templates

Component markup is stored in JavaScript template strings (`*.template.js`) and injected into
containers on first use. This keeps the published bundle self-contained and avoids shipping a
CSS file or a separate HTML partial for each component.

Consequences:

- The injected templates use Alpine expressions and the global `Uno` object. `window.Uno` is
  set by `src/index.js`, and the package marks that entry as having side effects.
- Containers are created lazily: no dialog, toast or tooltip DOM exists until the component is
  used for the first time. This is the "lazy initialization" introduced in 0.7.0.
- Content inserted by HTMX is handed back to Alpine with `Alpine.initTree()`. Uno registers an
  `htmx:afterSwap` listener during initialization when `window.htmx` is already present.
- The container IDs are stable and can be used for styling or tests:
  `uno-box-manager`, `uno-notify-container`, `uno-global-progress`,
  `uno-tooltip-manager`, `uno-tooltip-content`.

Templates are an implementation detail. For styling, use the `ui` configuration rather than
targeting internal classes; see [Theming](theming.md).

## Actions

An **Action** is a public store method that performs a state transition or opens UI. Examples:

```javascript
Alpine.store('uno').box.open({ title: 'Profile' })
Alpine.store('uno').notify.show('Saved', { type: 'success' })
Alpine.store('uno').tooltip.show({ element, content: 'Help' })
Alpine.store('uno').tab.activate({ name: 'settings', group: 'main' })
Alpine.store('uno').valid.validate('email', value, 'blur')
```

Callbacks (`onOk`, `onCancel`, `onClose`, `onContentLoad`, `onClick`) are not actions; they are
hooks invoked by an action or by a user interaction. Some actions return a reactive item
(`box.open`) or a Promise (`box.confirm`, `box.prompt`); the individual API pages document the
return value.

## Layers

Uno keeps z-index values in one place, `src/core/layer.js`:

| Layer | Value | Used by |
| :--- | ---: | :--- |
| `BOX_BACKDROP` | 1300 | Reserved backdrop band. |
| `BOX` | 1400 | Box manager container; each stacked dialog adds `index * 10`. |
| `NOTIFY` | 1600 | Notification container. |
| `TOOLTIP` | 1700 | Tooltip container. |
| `MAX` | 9999 | Global progress bar. |

Override the manager values globally:

```javascript
Uno.setConfig({
  ui: {
    layers: { box: 2000, notify: 2100, tooltip: 2200 }
  }
})
```

Layers do not defeat a parent stacking context. If a component is rendered inside a
`transform`, `filter`, `contain` or `position: fixed` ancestor that creates a stacking context,
the layer value only orders siblings inside that context. This is a CSS constraint, not a Uno
setting.

## Lifecycle

The order is:

1. **Script evaluation.** `src/index.js` attaches `window.Uno`, registers the built-in plugins
   (`i18n`, `tooltip`, `notify`, `box`, `valid`, `tab`), and decides how to reach Alpine.
2. **Alpine readiness.** If `window.Alpine` is already present, `Uno.init(window.Alpine)` runs
   immediately. Otherwise Uno waits for the one-time `alpine:init` event.
3. **Plugin registration.** Stores and directives are registered. If `window.htmx` exists,
   Uno adds the `htmx:afterSwap` listener that calls `Alpine.initTree()`.
4. **Global error handlers.** `unhandledrejection` and `error` listeners are installed once.
5. **First use.** A directive or store action creates the component container, calls
   `Alpine.initTree()` on it, and schedules focus/animation work.
6. **Interaction.** Alpine reactivity updates the injected template. Closing a Box waits for
   `transitionend` (with a 350 ms fallback) before removing the item.
7. **HTMX content.** After a swap, Box/Notify call their `onContentLoad` callback, re-run
   `Alpine.initTree()` on the loaded element, and (Box) schedule initial focus again.

`Uno.init()` is public but normally unnecessary. Calling it after Uno has already initialized
logs a warning and returns the Alpine instance. Use `Uno.setConfig()` for configuration and let
the automatic registration run.

## Configuration model

Configuration is resolved at four levels, lowest priority first:

1. **Plugin defaults** registered by the component (`Uno.registerPluginDefaults`).
2. **`plugins.<component>`** in `Uno.setConfig()`.
3. **Top-level `<component>`** in `Uno.setConfig()` (shorthand).
4. **Instance options** passed to an Action or directive expression.

The merge is deep for plain objects. Two string conventions apply:

- `{base}` is replaced with `ui.base` when a style is resolved.
- `+` at the start of a style string appends to the current value instead of replacing it.

See [Configuration](configuration.md) for timing, defaults and the full `ui` map.

## Sanitization boundary

Uno's built-in sanitizer is a small regex pass, not an HTML parser. It removes `<script>`,
`<img>`, `on*` attributes and `javascript:` URLs. It is intended for semi-trusted content and
can be bypassed. The boundary is:

| Content | Handling |
| :--- | :--- |
| Box `content` | Passed through `Uno.dom.sanitize()`. |
| Box `prompt` message and default value | Passed through `Uno.dom.sanitize()` before interpolation. |
| Tooltip `content` and `target` innerHTML | Passed through `sanitizeHTML()`. |
| Notify `message` | Rendered with `x-text`; it is not interpreted as HTML. |
| Notify HTMX response | Inserted with `x-html`; treat the server response as trusted. |
| Tooltip HTMX response | Inserted by HTMX; treat the server response as trusted. |
| Valid error messages | Rendered with `x-text`; plain text. |

For untrusted HTML, use a dedicated sanitizer such as DOMPurify before passing content to Uno.
See [Security](https://github.com/artsyne/uno-ui/blob/main/SECURITY.md) and the [CSP guide](csp.md).

## Next

- [Configuration](configuration.md)
- [Theming](theming.md)
- [HTMX integration](htmx.md)
- [Accessibility](accessibility.md)
