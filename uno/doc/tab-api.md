# Uno Tab API Documentation

The Tab component offers a flexible tab interface with support for deep linking, hash synchronization, and lazy loading of content via HTMX.

## Overview

- **Core Function**: Manages tab switching logic and visibility states.
- **Architecture**: Directive-based binding (`x-bind`) connected to a central store (`Alpine.store('uno').tab`).
- **Deep Linking**: Automatically syncs active tab with URL hash.
- **HTMX Support**: Handles "simulated clicks" to ensure remote content loads when landing directly on a tab URL.

## Basic Usage

### 1. Window Helper Usage (Recommended)

Use `window.unoTabs()` to initialize the component scope.

```html
<div x-data="window.unoTabs('profile')">
  <!-- Tab List -->
  <div x-bind="tablist" role="tablist" class="flex">
    <button x-bind="trigger('profile')" class="btn">Profile</button>
    <button x-bind="trigger('settings')" class="btn">Settings</button>
  </div>

  <!-- Panels -->
  <div class="content">
    <div x-bind="panel('profile')">Profile Content</div>
    <div x-bind="panel('settings')">Settings Content</div>
  </div>
</div>
```

### 2. Directive Usage

For simpler, standalone tab triggers.

```html
<button x-uno-tab="'profile'">Go to Profile</button>
```

## API Parameters

### `window.unoTabs(defaultTab, options)`

| Name            | Type      | Default     | Required | Description                                                     |
| :-------------- | :-------- | :---------- | :------- | :-------------------------------------------------------------- |
| `defaultTab`    | `String`  | -           | Yes      | The name of the tab to activate initially (if no hash matches). |
| `options`       | `Object`  | `{}`        | No       | Configuration object.                                           |
| `options.hash`  | `Boolean` | `true`      | No       | Whether to enable URL hash synchronization.                     |
| `options.group` | `String`  | `'default'` | No       | The group namespace for these tabs.                             |

### `trigger(name, ajaxUrl, hash, hxUrl)` (Binding)

| Name      | Type     | Default | Required | Description                                                        |
| :-------- | :------- | :------ | :------- | :----------------------------------------------------------------- |
| `name`    | `String` | -       | Yes      | Unique identifier for the tab.                                     |
| `ajaxUrl` | `String` | -       | No       | Simple URL to fetch content from (replaces panel content).         |
| `hash`    | `String` | -       | No       | Override the hash value (default uses `name`).                     |
| `hxUrl`   | `String` | -       | No       | **Legacy**. Use standard `hx-get` attribute on the button instead. |

## Events

| Event        | Trigger Timing                       | Callback Parameters                            |
| :----------- | :----------------------------------- | :--------------------------------------------- |
| `hashchange` | When URL hash changes (browser nav). | None (Store updates active tab automatically). |

## Methods

Accessed via `Alpine.store('uno').tab`.

### `activate(options)`

Activates a specific tab.

- **Parameters**: `options` (Object)
  - `name` (String, Required): Tab ID.
  - `group` (String, Optional): Tab group name.
  - `updateHash` (Boolean, Optional): Whether to update URL hash (default `true`).
- **Returns**: `void`
- **Example**:
  ```javascript
  Alpine.store('uno').tab.activate({ name: 'settings' })
  ```

### `isActive(name, group)`

Checks if a tab is currently active.

- **Parameters**:
  - `name` (String): Tab ID.
  - `group` (String): Group name (default `'default'`).
- **Returns**: `Boolean`
- **Example**:
  ```javascript
  if (Alpine.store('uno').tab.isActive('profile')) { ... }
  ```

### `register(name, group)`

Registers a tab in the store. Automatically called by `trigger` binding.

- **Parameters**:
  - `name` (String): Tab ID.
  - `group` (String): Group name.
- **Returns**: `void`

## Examples

### HTMX Deep Linking

This setup ensures that if a user visits `page#users`, the "Users" tab activates AND the HTMX request fires to load the content.

```html
<div x-data="window.unoTabs('dashboard')">
  <div x-bind="tablist">
    <!-- Standard Tab -->
    <button x-bind="trigger('dashboard')">Dashboard</button>

    <!-- HTMX Tab -->
    <!-- hx-get will be triggered automatically if landing on #users -->
    <button x-bind="trigger('users')" hx-get="/api/users" hx-target="#panel-users">Users</button>
  </div>

  <div>
    <div x-bind="panel('dashboard')">...</div>
    <div id="panel-users" x-bind="panel('users')"></div>
  </div>
</div>
```

## Troubleshooting

### Common Issues

| Issue | Possible Cause | Solution |
| :--- | :--- | :--- |
| **Tab not switching** | Name mismatch between trigger and panel. | Ensure `trigger('myTab')` matches `panel('myTab')`. |
| **Hash not updating** | `options.hash` set to `false`. | Enable hash sync in `window.unoTabs('default', { hash: true })`. |
| **Content not loading (HTMX)** | Trigger missing `hx-get` or invalid URL. | Check console for HTMX errors. Ensure `hx-target` is correct. |

## Version History

| Version | Date | Changes |
| :--- | :--- | :--- |
| v0.8.0 | 2026-02-28 | Added `bindNotifyToHtmx()` method to HTMX integration for automatic notifications. |
| v0.7.0 | 2026-01-20 | **Lazy Listener**: Hash change listener now starts only after first tab usage. |
| v0.6.0 | 2025-12-18 | Initial stable release. |
