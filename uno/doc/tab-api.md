# Tab Component

A tab component built on Alpine.js, supporting multiple tab groups, AJAX content loading, HTMX integration, and Hash routing synchronization.

## Features

- **Declarative API**: Use `x-data` and `x-bind` to set up tabs easily.
- **Hash Navigation**: Deep linking support via URL Hash.
- **AJAX Loading**: Asynchronous content loading.
- **HTMX Support**: Native support for HTMX content loading (`hxUrl`).
- **Programmatic Control**: Extensible factory pattern.

## Quick Start

### HTML Structure

The component uses `x-data="window.unoTabs(defaultTab)"` to initialize, and `x-bind="trigger(tabName)"` and `x-bind="panel(tabName)"` to bind elements.

```html
<div x-data="window.unoTabs('tab1')">
  <!-- Tab List -->
  <div x-bind="tablist" class="flex border-b mb-4">
    <button x-bind="trigger('tab1')" class="tab-base px-4 py-2">Tab 1</button>
    <button x-bind="trigger('tab2')" class="tab-base px-4 py-2">Tab 2</button>
  </div>

  <!-- Tab Panels -->
  <div>
    <div x-bind="panel('tab1')">Content 1</div>
    <div x-bind="panel('tab2')">Content 2</div>
  </div>
</div>
```

## API Reference

### Global API (Uno.tab)

A global helper is provided for external control.

#### `Uno.tab.open(id)`
Opens the tab with the specified target ID.

```javascript
Uno.tab.open('tab-profile');
```

### Initialization

```javascript
// x-data="window.unoTabs(defaultTab, options)"
window.unoTabs('tab1', { 
  hash: true,      // Enable hash navigation (default: true)
  manual: false    // Manual activation via arrow keys (default: false)
})
```

### Bindings

| Binding | Description | Arguments |
|---------|-------------|-----------|
| `tablist` | Binds accessibility roles and keyboard navigation (Left/Right/Home/End) to the container. | - |
| `trigger(tab, ajaxUrl, hash, hxUrl)` | Binds click handlers, classes, and accessibility attributes to a tab button. | `tab`: Tab ID (string)<br>`ajaxUrl`: Optional URL for fetch<br>`hash`: Optional URL hash<br>`hxUrl`: Optional URL for HTMX |
| `panel(tab)` | Binds visibility (`x-show`) and accessibility roles to a content panel. | `tab`: Tab ID (string) |

## Advanced Usage

### AJAX Content Loading

Pass a URL as the second argument to `trigger()`. The content will be fetched via `fetch()` and injected into the panel.

```html
<button x-bind="trigger('tab-ajax', '/api/content.html')">Load Ajax</button>
<div x-bind="panel('tab-ajax')">Loading...</div>
```

### HTMX Content Loading

Pass a URL as the fourth argument to `trigger()` (or use named arguments if supported, but here strictly positional). The content will be fetched via `htmx.ajax()`.

```html
<!-- trigger(tabId, ajaxUrl, hash, hxUrl) -->
<button x-bind="trigger('tab-htmx', null, null, '/api/content.html')">Load HTMX</button>
<div x-bind="panel('tab-htmx')">Loading...</div>
```

### Hash Navigation (Deep Linking)

Pass a hash string as the third argument to `trigger()`.

```html
<button x-bind="trigger('tab-profile', null, 'profile')">Profile</button>
```

- Clicking updates URL to `...#profile`.
- Loading URL with `#profile` activates this tab.

## Security & Optimization

### Content Sanitization
AJAX and HTMX loaded content is sanitized via `Uno.dom.sanitizeHTML` before injection.
- Dangerous tags are removed.
- Use secure endpoints for `ajaxUrl` and `hxUrl`.

### Memory Management
- **Event Cleanup**: Hash change listeners are automatically managed and cleaned up when tabs are destroyed or re-initialized.
- **Idempotency**: Initialization is safe to call multiple times (via `Uno.utils.once`).

## Styling

The component manages state classes but leaves visual styling to you.
It applies the `current` class to the active tab trigger.

```css
.tab-base {
  /* Base styles */
}
.tab-base.current {
  /* Active styles */
  color: blue;
  border-bottom: 2px solid blue;
}
```
