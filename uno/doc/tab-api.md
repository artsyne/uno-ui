# Tab Component

A tab component built on Alpine.js, supporting multiple tab groups, AJAX content loading, Hash routing synchronization, and programmatic control.

## Features

- **Multiple Groups**: Supports multiple independent tab groups on a single page.
- **Hash Navigation**: Deep linking support via URL Hash.
- **AJAX Loading**: Asynchronous content loading with caching.
- **Programmatic Control**: Switch tabs via Alpine Store.

## Quick Start

### HTML Structure

The component uses `data-tab-group` to link buttons with content, and `data-tab-target` to specify the content ID to show.

```html
<!-- Tab Buttons -->
<div class="flex border-b mb-4">
  <button data-tab-group="main" data-tab-target="tab1" class="tab-base current px-4 py-2 border-b-2 border-blue-500 text-blue-500">Tab 1</button>
  <button data-tab-group="main" data-tab-target="tab2" class="tab-base px-4 py-2 text-gray-500">Tab 2</button>
</div>

<!-- Tab Content -->
<div id="tab1">
  <p>Content for Tab 1</p>
</div>
<div id="tab2" class="hidden">
  <p>Content for Tab 2</p>
</div>
```

## Configuration / Options

### Tab Triggers (Buttons/Links)

| Attribute | Description | Example |
|-----------|-------------|---------|
| `data-tab-group` | **Required**. Defines the group namespace. Only tabs in the same group affect each other. | `data-tab-group="settings"` |
| `data-tab-target` | **Required**. The ID of the content element to show. | `data-tab-target="profile"` |
| `data-tab-hash` | Optional. Updates URL Hash on click. | `data-tab-hash="user-profile"` |
| `data-tab-url` | Optional. URL to fetch content via AJAX. | `data-tab-url="/api/profile"` |
| `class` | `tab-base` provides base styles; `current` marks the active tab. | `class="tab-base current"` |

### Tab Content Containers

| Attribute | Description |
|-----------|-------------|
| `id` | **Required**. Must match the `data-tab-target` of the trigger. |
| `class` | Use `hidden` class for content that should be initially hidden. |
| `data-tab-loaded`| **Internal**. Added automatically after AJAX content is loaded to prevent re-fetching. |

## Advanced Usage

### AJAX Content Loading

Add `data-tab-url` attribute to enable dynamic loading. Content will be fetched via GET request and injected into the target element.

```html
<button 
  data-tab-group="main" 
  data-tab-target="ajax-content" 
  data-tab-url="/partials/content.html">
  Load Remote Content
</button>

<div id="ajax-content">
  <!-- Skeleton loader will appear here automatically -->
</div>
```

**Behavior Logic:**
1. First click shows a skeleton loader.
2. Fetches content from URL.
3. On success, content replaces the inner HTML of the target element.
4. Target element is marked with `data-tab-loaded="true"`.
5. Subsequent clicks show cached content immediately.

### Hash Navigation (Deep Linking)

Use `data-tab-hash` to sync tab state with URL Hash.

```html
<button 
  data-tab-group="main" 
  data-tab-target="tab1" 
  data-tab-hash="section1">
  Chapter 1
</button>
```

- Clicking this tab changes URL to `...#section1`.
- Visiting URL with `...#section1` automatically activates this tab.

## Programmatic Usage

The component exposes an Alpine store for programmatic control.

```javascript
// Switch Tab
// Alpine.store('tabs').switch(groupName, targetId)
Alpine.store('tabs').switch('main', 'tab2');
```

## Styling

Component uses standard Tailwind CSS classes:

- `.tab-base`: Base style for tab buttons.
- `.tab-active` / `.current`: Active state style (usually blue text, bottom border).
- `.tab-inactive`: Inactive state style (gray text).

Override these classes or add utility classes directly to elements for customization.
