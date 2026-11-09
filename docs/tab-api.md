# Tab API

Tab manages tab groups, hash synchronization and HTMX content loading. The state lives in
`Alpine.store('uno').tab`; the `x-uno-tab` directive registers triggers, injects the tab ARIA
attributes and handles keyboard navigation. Panels are ordinary markup bound with `x-show`.

## Prerequisites

- Alpine.js 3 with Uno registered. See [Getting Started](getting-started.md).
- HTMX 2 only for `hx-get` tabs and hash-triggered loading.

## Minimal example

```html
<div x-data>
  <div role="tablist" aria-label="Account sections" class="flex">
    <button
      id="tab-profile"
      x-uno-tab="{ name: 'profile', group: 'account' }"
    >
      Profile
    </button>
    <button
      id="tab-settings"
      x-uno-tab="{ name: 'settings', group: 'account' }"
    >
      Settings
    </button>
  </div>

  <div
    id="panel-profile"
    role="tabpanel"
    aria-labelledby="tab-profile"
    x-show="$store.uno.tab.isActive('profile', 'account')"
  >
    Profile content
  </div>
  <div
    id="panel-settings"
    role="tabpanel"
    aria-labelledby="tab-settings"
    x-show="$store.uno.tab.isActive('settings', 'account')"
  >
    Settings content
  </div>
</div>
```

The directive sets `role="tab"`, `aria-controls` (`config.panel` or `panel-{name}`),
`aria-selected` and roving `tabindex`. Do not bind `:aria-selected` yourself; the directive
owns it.

## Directive

`x-uno-tab` accepts a string (name) or an object:

```html
<button x-data x-uno-tab="'profile'">Profile</button>
<!-- equivalent to { name: 'profile', group: 'default' } -->
```

| Option | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `name` | `string` | required | Tab identifier, unique in practice. |
| `group` | `string` | `'default'` | Group namespace; one active tab per group. |
| `panel` | `string` | `panel-{name}` | Element id used for `aria-controls`. |
| `hash` | `boolean` | `true` | Whether activation writes the URL hash. |

A directive expression that is a string is normalized to `{ name }`.

## Store API

```javascript
const tab = Alpine.store('uno').tab
```

### `activate(options)`

Activates a tab. It **must** receive an object; positional arguments are ignored.

```javascript
tab.activate({ name: 'settings', group: 'account' })
```

| Option | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `name` | `string` | required | Tab name. |
| `group` | `string` | `'default'` | Group. |
| `hash` | `boolean` | `true` | Write the hash. |
| `updateHash` | `boolean` | `true` | Set to `false` when activation comes from a hash change; this also enables the simulated HTMX click. |

Behaviour:

- Registers the tab if it is not in `registry`.
- Sets `groups[group] = name`.
- Updates `aria-labelledby` on shared panels (`role="tabpanel" data-tab-group="{group}"`).
- Writes `#name` with `history.replaceState()` when `hash` is true and `updateHash` is not
  `false`. No new history entry is created.
- When `updateHash === false`, schedules a click on `#tab-{name}` if it has `hx-get` and HTMX is
  loaded.

### `isActive(name, group)`

```javascript
tab.isActive('profile', 'account') // boolean
```

`group` defaults to `'default'`.

### `register(name, group)`

Registers a tab and makes it active if the group has no active tab. Called by the directive.
If the current hash equals `name`, it activates with `updateHash: false`.

### `rotate(group, current, key)`

Returns the name of the next tab for `ArrowLeft`/`ArrowRight`/`Home`/`End`, or `null`. It does
not activate anything; the directive uses it for the roving keyboard pattern.

## Hash synchronization

- The first tab registered in a group becomes active.
- On registration, a matching `location.hash` activates that tab with `updateHash: false`.
- `hashchange` and `load` listeners are installed lazily, only after the first tab registration
  or activation.
- `history.replaceState()` is used, so back/forward does not step through tab changes.
- A global `tab.hash` setting is not consulted by `activate()`: the action supplies
  `hash: true` as an instance default, which wins over global configuration. Pass `hash: false`
  on the directive or on the `activate()` call instead.

## HTMX

Use ordinary HTMX attributes on the trigger. Keep the id in the `tab-{name}` form so deep-link
activation can trigger the request:

```html
<button
  id="tab-users"
  x-uno-tab="{ name: 'users', group: 'app' }"
  hx-get="/users"
  hx-target="#panel-users"
>
  Users
</button>
<div id="panel-users" role="tabpanel" aria-labelledby="tab-users"></div>
```

A normal click activates the tab and HTMX sends the request. A visit to `#users` activates the
tab with `updateHash: false`; when HTMX is loaded, Uno calls `htmx.trigger(el, 'click')` if
`hx-get` is present.

## Accessibility

- `role="tab"` and `aria-selected` are set by the directive.
- Roving `tabindex`: active tab `0`, others `-1`.
- `aria-controls` points at `panel` or `panel-{name}`.
- ArrowLeft/ArrowRight rotate within the group; Home/End go to the first/last tab; focus follows
  and activation reuses the click path. RTL containers reverse the horizontal arrows.
- Per-tab panels declare `id="panel-{name}" role="tabpanel" aria-labelledby="tab-{name}"`.
- Shared-panel mode declares `role="tabpanel" data-tab-group="{group}"`; activation updates
  `aria-labelledby`.
- The tab list should carry an accessible label (`aria-label` or `aria-labelledby`).

The directive does not inject panel roles or `x-show`; those remain the author's markup.

## Known limitations

- `registry` is keyed by tab name globally. Two groups using the same `name` collide for
  hash-based activation; use unique names when hash sync matters.
- `history.replaceState()` means tab switches are not browser history entries.
- `hashchange` is read from `window.location`; custom routers that only use
  `history.pushState` are not observed.
- `ui.tab.active/inactive/loading/error` tokens are registered but not consumed by the current
  directive, which has no template. They are reserved for application-level styling.
- The directive's global `group`/`hash` defaults are not read; set them per tab/action.

## Version history

| Version | Date | Changes |
| :--- | :--- | :--- |
| 0.10.0 | 2026-11-09 | APG tabs pattern: roles, `aria-selected`, roving tabindex, Arrow/Home/End, per-tab and shared panels. |
| 0.8.0 | 2026-01-22 | Documentation realigned with the Store-First implementation. |
| 0.7.0 | 2026-01-20 | Hash listener starts only after first tab use. |
| 0.6.0 | 2025-12-18 | Initial release. |

## Related

- [HTMX integration](htmx.md#tab)
- [Accessibility](accessibility.md)
- [Configuration](configuration.md)
