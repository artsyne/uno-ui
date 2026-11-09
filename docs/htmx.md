# HTMX Integration

HTMX is an optional peer dependency (`htmx.org` 2.x). Uno uses HTMX for server-loaded content
in Box, Notify, Tab and Valid, and exposes a small helper object at `Uno.htmx`.

Uno targets **HTMX 2.x** and reads only the XHR response (`event.detail.xhr`). The Fetch API
integration is not supported.

## Load order

Load HTMX before Uno so Uno can install its `htmx:afterSwap` hook during initialization:

```html
<script src="https://cdn.jsdelivr.net/npm/htmx.org@2.0.10/dist/htmx.min.js"></script>
<script src="/dist/uno.umd.js" defer></script>
<script src="https://cdn.jsdelivr.net/npm/alpinejs@3.17.3/dist/cdn.min.js" defer></script>
```

If HTMX loads after Uno initializes, components that use `hx-*` attributes may still work
because HTMX processes the DOM itself, but Uno will not install its `afterSwap` → `initTree`
hook automatically. Call `Alpine.initTree(event.detail.elt)` from your own `htmx:afterSwap`
listener in that case.

## `Uno.htmx` helpers

```javascript
Uno.htmx.available() // boolean
Uno.htmx.ajax(method, url, options) // Promise
Uno.htmx.fragment(method, url, options) // Promise<string>
Uno.htmx.onRequest({ start, end, progress })
Uno.htmx.bindNotifyToHtmx(options)
```

| Helper | Behavior |
| :--- | :--- |
| `available()` | True when `window.htmx.ajax` is a function. |
| `ajax(method, url, options)` | Thin wrapper around `htmx.ajax`; rejects when HTMX is unavailable. |
| `fragment(method, url, options)` | Fetches HTML without inserting it. Returns the server response string. |
| `onRequest(handlers)` | Registers lifecycle handlers; listeners are installed once. |
| `bindNotifyToHtmx(options)` | Shows progress/success/error notifications for HTMX requests. Can be bound once. |

### `fragment()`

```javascript
try {
  const html = await Uno.htmx.fragment('GET', '/api/modal-content')
  Alpine.store('uno').box.open({ title: 'Details', content: html })
} catch (error) {
  Alpine.store('uno').notify.error('Could not load content')
}
```

`fragment()` creates a detached target, intercepts `htmx:beforeSwap`, and returns
`detail.serverResponse` without swapping it into the page. A failed request rejects and is also
logged through Uno's error reporter.

### `onRequest()`

```javascript
Uno.htmx.onRequest({
  start: (event) => console.log('start', event.detail),
  end: (event) => {
    if (event.detail.failed) console.warn('request failed')
  },
  progress: (event) => {
    // event.detail is the htmx:xhr:progress event
  }
})
```

- `start` runs on `htmx:beforeRequest`.
- `end` runs on `htmx:afterRequest`, `htmx:responseError`, `htmx:sendError`, and
  `htmx:timeout`.
- `progress` runs on `htmx:xhr:progress`.
- Duplicate `start`/`end` calls for the same XHR are suppressed with a `WeakSet`.
- Handlers registered with `onRequest()` are stored, but the DOM listeners are only installed
  when HTMX is already present. Call it after HTMX loads.

### `bindNotifyToHtmx()`

```javascript
Uno.htmx.bindNotifyToHtmx({
  showSuccess: true,
  showError: true,
  showProgress: true,
  successMessage: 'Operation completed',
  errorMessage: 'Operation failed'
})
```

| Option | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `showSuccess` | `boolean` | `true` | Show a success toast after a request that did not fail. |
| `showError` | `boolean` | `true` | Show an error toast on response/send/timeout errors. |
| `showProgress` | `boolean` | `true` | Drive the global progress bar while requests are active. |
| `successMessage` | `string` | `'Operation completed'` | Default success text. |
| `errorMessage` | `string` | `'Operation failed'` | Default error text. |

The binding is idempotent: only the first call installs listeners. It requires HTMX to be
available at call time and logs a warning otherwise.

Progress is driven with `notify.startProgress()` / `notify.stopProgress()` and a local request
counter. If you need concurrency-safe progress in your own code, use
`notify.beginProgress()` / `notify.endProgress()` instead.

#### Server-sent messages

| Header | Effect |
| :--- | :--- |
| `X-Notify-Message` | Overrides the success message. |
| `X-Notify-Error` | Overrides the error message. |
| `X-Notify-Suppress` | Any non-empty value suppresses the success toast. |

```http
X-Notify-Message: Item saved
```

The values are read from the XHR response headers on the `afterRequest` / error event.

## Alpine initialization after a swap

When Uno initializes and `window.htmx` exists, it installs:

```javascript
document.addEventListener('htmx:afterSwap', (event) => {
  const elt = event?.detail?.elt
  if (elt) window.Alpine.initTree(elt)
})
```

This lets Alpine directives inside swapped fragments initialize. Without it, `x-data`,
`x-on` and Uno directives in HTMX-loaded content stay inert. The listener is installed once.

## Box

```javascript
Alpine.store('uno').box.open({
  title: 'Edit user',
  hxUrl: '/users/1/edit',
  hxMethod: 'GET',
  hxHeaders: { 'X-Requested-With': 'Uno' },
  onContentLoad: (el) => {
    console.log('loaded', el)
  }
})
```

| Option | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `hxUrl` | `string` | — | URL to load. |
| `hxMethod` | `'GET' \| 'POST'` | `'GET'` | Method used by the template. Other values are not rendered. |
| `hxHeaders` | `Record<string, string>` | — | JSON-encoded into `hx-headers`. |
| `hxTrigger` | `string` | `'load'` | HTMX trigger. |
| `onContentLoad` | `(el: HTMLElement) => void` | — | Called after the swap, before focus is recalculated. |

The Box template calls `onContentLoad(el)`, `Alpine.initTree(el)` and `scheduleAutoFocus()`.
The response is inserted by HTMX as HTML; treat it as trusted server output. Box has no
`onError` callback: listen to `htmx:responseError` / `htmx:sendError` globally.

## Notify

```javascript
Alpine.store('uno').notify.show('Loading…', {
  hxUrl: '/partials/status',
  hxTrigger: 'load',
  hxSwap: 'innerHTML',
  hxHeaders: { 'X-Requested-With': 'Uno' },
  duration: 10000
})
```

| Option | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `hxUrl` | `string` | — | URL to load into the notification. |
| `hxTrigger` | `string` | `'load'` | HTMX trigger attribute. |
| `hxSwap` | `string` | `'innerHTML'` | Swap strategy. |
| `hxHeaders` | `Record<string, string>` | — | Serialized into `hx-headers`. |
| `onContentLoad` | `(el: HTMLElement) => void` | — | Called after the response is copied into the item. |
| `onError` | `(detail: any) => void` | — | Called on response/send/timeout errors. |

The response is copied into `item.content` and rendered with `x-html`; it is not sanitized.
Notify drives `beginProgress()`/`endProgress()` around the request.

## Tab

HTMX tabs use ordinary HTMX attributes on the trigger. The trigger should be
`id="tab-{name}"`:

```html
<button
  id="tab-users"
  x-uno-tab="{ name: 'users', group: 'app' }"
  hx-get="/users"
  hx-target="#panel-users"
>
  Users
</button>
```

- A user click activates the tab and lets HTMX issue the request natively.
- A deep link (`#users`) activates the tab with `updateHash: false`. If the trigger has
  `hx-get`, Uno calls `htmx.trigger(el, 'click')` so the content still loads.
- The hash listener starts only after the first tab is registered or activated.

## Valid

The Valid plugin intercepts `htmx:confirm` on `document.body` when the form carries
`x-uno-valid` or contains a field with `x-uno-valid`:

```html
<form x-uno-valid hx-post="/save" hx-swap="none">
  <input name="email" x-uno-valid="['required', 'email']" />
  <button type="submit">Save</button>
</form>
```

If validation fails, the request is not issued. If validation passes,
`evt.detail.issueRequest(true)` continues the HTMX request.

Server-side validation errors are accepted on HTTP `422` with a JSON body:

```json
{
  "errors": {
    "email": "Email is already registered",
    "name": "Name is required"
  }
}
```

Uno calls `valid.showErrors(errors)`, which updates inline errors and `aria-invalid`. Non-JSON
422 bodies are ignored. An exception during validation fails closed: the request is blocked and
the error is logged.

## Server contract

- Return HTML fragments for `hxUrl` / `hx-get`, not a full document.
- Return `422` with `{ "errors": { ... } }` for server-side validation failures.
- Treat every HTML fragment as trusted. Uno does not sanitize HTMX responses.
- For notifications, use the `X-Notify-*` headers described above.

## Limitations

- HTMX 2 XHR only; the Fetch API extension is not supported.
- `hxMethod` in Box is limited to `GET` and `POST` by the template.
- Uno does not install the `afterSwap` hook if HTMX is loaded after Uno initializes.
- `bindNotifyToHtmx()` binds once per page; call `Uno.htmx.onRequest()` if you need multiple
  independent handlers.
- HTMX content is HTML; sanitize untrusted server output before returning it.

## Related

- [Notify API](notify-api.md)
- [Box API](box-api.md)
- [Tab API](tab-api.md)
- [Valid API](valid-api.md)
- [Error handling](error-handling.md)
