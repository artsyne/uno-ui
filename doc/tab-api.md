# Tab API

The Tab component provides a flexible tab interface with support for deep linking (URL hash synchronization) and lazy loading of remote content via HTMX.

## Overview

- **Core Function**: Manages tab switching logic and visibility states per group.
- **Architecture**: Store-First — all state lives in `Alpine.store('uno').tab`; the `x-uno-tab` directive only registers triggers and forwards click events.
- **Deep Linking**: The active tab is synchronized with the URL hash (`#tab-name`). Landing directly on a hash URL activates the matching tab automatically.
- **HTMX Support**: When a tab is activated via hash (not by user click) and its trigger element (`#tab-{name}`) carries `hx-get`, a simulated click is fired so remote content loads without user interaction.
- **Panel Visibility**: Panels are plain `x-show` bindings — there is no `panel()` helper. You bind each panel to `isActive(name, group)` yourself, which keeps markup fully under your control.

## Basic Usage

### 1. Directive + Store (Recommended)

```html
<div x-data>
  <!-- Tab List -->
  <div class="flex" role="tablist">
    <button
      id="tab-profile"
      role="tab"
      x-uno-tab="{ name: 'profile', group: 'main' }"
      :aria-selected="$store.uno.tab.isActive('profile', 'main')"
    >
      Profile
    </button>
    <button
      id="tab-settings"
      role="tab"
      x-uno-tab="{ name: 'settings', group: 'main' }"
      :aria-selected="$store.uno.tab.isActive('settings', 'main')"
    >
      Settings
    </button>
  </div>

  <!-- Panels: plain x-show against the store -->
  <div x-show="$store.uno.tab.isActive('profile', 'main')" role="tabpanel">Profile Content</div>
  <div x-show="$store.uno.tab.isActive('settings', 'main')" role="tabpanel">Settings Content</div>
</div>
```

The directive registers the tab on initialization and activates it on click (calling `e.preventDefault()`, so it is safe to use on links). The first registered tab of a group becomes active by default.

### 2. Shorthand String Form

```html
<button x-uno-tab="'profile'">Go to Profile</button>
```

Equivalent to `{ name: 'profile', group: 'default' }`.

## API Parameters

### `x-uno-tab` Directive Value

| Name    | Type      | Default     | Required | Description                                     |
| :------ | :-------- | :---------- | :------- | :---------------------------------------------- |
| `name`  | `String`  | -           | Yes      | Unique identifier for the tab.                  |
| `group` | `String`  | `'default'` | No       | Group namespace; each group has one active tab. |
| `hash`  | `Boolean` | `true`      | No       | Whether activation updates the URL hash.        |

Global defaults can be set via `Uno.setConfig({ tab: { group: '...', hash: false } })` before Alpine starts.

## Methods

Accessed via `Alpine.store('uno').tab`.

### `activate(options)`

Activates a specific tab.

- **Parameters**: `options` (Object) — **must be an object**; positional string arguments are not supported.
  - `name` (String, Required): Tab name.
  - `group` (String, Optional): Group name (default `'default'`).
  - `hash` (Boolean, Optional): Whether hash sync is enabled (default `true`).
  - `updateHash` (Boolean, Optional): Whether to write the hash to the URL (default `true`). Internally set to `false` when activation originates from a hash change, in which case a simulated HTMX click is fired for lazy content.
- **Returns**: `void`
- **Example**:
  ```javascript
  Alpine.store('uno').tab.activate({ name: 'settings', group: 'main' })
  ```

### `isActive(name, group)`

Checks if a tab is currently active.

- **Parameters**:
  - `name` (String): Tab name.
  - `group` (String): Group name (default `'default'`).
- **Returns**: `Boolean`
- **Example**:
  ```javascript
  if (Alpine.store('uno').tab.isActive('profile', 'main')) { ... }
  ```

### `register(name, group)`

Registers a tab in the store. Called automatically by the `x-uno-tab` directive on initialization; manual calls are rarely needed.

- **Parameters**:
  - `name` (String): Tab name.
  - `group` (String): Group name (default `'default'`).
- **Behavior**: If the group has no active tab yet, the newly registered tab becomes active. If `name` matches the current URL hash, the tab activates with `updateHash: false` (deep-link path).
- **Returns**: `void`

## Events

The component listens to native browser events; it does not emit custom DOM events.

| Event        | Source   | Behavior                                                      |
| :----------- | :------- | :------------------------------------------------------------ |
| `hashchange` | `window` | Activates the tab whose `name` matches the new hash.          |
| `load`       | `window` | Re-checks the hash (covers tabs registered after store init). |

The hash listener is registered lazily — only after the first tab registration or activation.

## Examples

### HTMX Deep Linking

If a user visits `page#users`, the "Users" tab activates AND the HTMX request fires automatically to load the content.

```html
<div x-data>
  <div class="flex" role="tablist">
    <!-- Standard Tab -->
    <button
      id="tab-dashboard"
      x-uno-tab="{ name: 'dashboard', group: 'app' }"
      :aria-selected="$store.uno.tab.isActive('dashboard', 'app')"
    >
      Dashboard
    </button>

    <!-- HTMX Tab: hx-get fires on click, and on deep-link activation -->
    <button
      id="tab-users"
      x-uno-tab="{ name: 'users', group: 'app' }"
      :aria-selected="$store.uno.tab.isActive('users', 'app')"
      hx-get="/api/users"
      hx-target="#panel-users"
    >
      Users
    </button>
  </div>

  <div x-show="$store.uno.tab.isActive('dashboard', 'app')" role="tabpanel">...</div>
  <div id="panel-users" x-show="$store.uno.tab.isActive('users', 'app')" role="tabpanel"></div>
</div>
```

> **Note**: The simulated click targets the element with id `tab-{name}` — keep your trigger `id` in that format if you rely on deep-link HTMX loading.

## Known Limitations

- **Registry is keyed by tab name globally**: if two groups contain tabs with the same `name`, hash-based deep linking resolves to whichever registered first. Use unique names across groups when hash sync matters.
- **ARIA attributes are user-managed**: the directive does not inject `role`/`aria-selected`; bind them yourself as shown above.

## Troubleshooting

### Common Issues

| Issue                                        | Possible Cause                                                            | Solution                                                                                    |
| :------------------------------------------- | :------------------------------------------------------------------------ | :------------------------------------------------------------------------------------------ |
| **Tab not switching**                        | `name` mismatch between directive and `isActive()` call.                  | Ensure both use the same `name` AND `group`.                                                |
| **Nothing is active initially**              | No tab registered for the group yet.                                      | The first registered tab activates automatically; check the directive expression evaluates. |
| **Hash not updating**                        | `hash: false` in config or global tab config.                             | Remove `hash: false`, or check `Uno.setConfig({ tab: ... })`.                               |
| **Content not loading (HTMX)**               | Trigger `id` not in `tab-{name}` format, or missing `hx-get`/`hx-target`. | Set `id="tab-yourName"` on the trigger and verify HTMX attributes.                          |
| **`activate('name', 'group')` does nothing** | Positional string arguments are not supported.                            | Pass an object: `activate({ name: 'name', group: 'group' })`.                               |

## Version History

| Version | Date       | Changes                                                                                                |
| :------ | :--------- | :----------------------------------------------------------------------------------------------------- |
| v0.8.0  | 2026-02-28 | Documentation realigned with the Store-First implementation (`x-uno-tab` + `Alpine.store('uno').tab`). |
| v0.7.0  | 2026-01-20 | **Lazy Listener**: Hash change listener now starts only after first tab usage.                         |
| v0.6.0  | 2025-12-18 | Initial stable release.                                                                                |
