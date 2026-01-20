# Migration Guide: v0.6.0 to v0.7.0

Uno UI v0.7.0 introduces a major architectural change: **Lazy Initialization**. This guide helps you upgrade your application to the new version.

## Key Changes

In previous versions (v0.6.x), Uno UI would "eagerly" initialize all components upon startup. This meant that:
- HTML containers for Modals, Notifications, and Tooltips were injected into the `<body>` immediately.
- Global event listeners (like `hashchange` for Tabs) were attached immediately.

In v0.7.0, we have switched to **Lazy Initialization**. Components now only initialize their DOM elements and event listeners when they are **first used**.

## Breaking Changes & Migration Steps

### 1. Manual Initialization Calls

**Old Behavior (v0.6.x):**
You might have manually called `init()` methods to ensure components were ready.

```javascript
// Old code - No longer needed/effective
Alpine.store('uno').box.init()
Alpine.store('uno').notify.init()
```

**New Behavior (v0.7.0):**
These `init()` methods have been modified to be lazy or removed from the public API contract. Calling them is no longer necessary and may not have the expected effect of forcing DOM injection.

**Migration Action:**
Remove any manual calls to `init()` methods for Uno stores. Rely on the automatic lazy initialization that happens when you call `open()`, `show()`, or use directives.

### 2. CSS/DOM Selectors

**Old Behavior:**
Since containers were injected immediately, you might have relied on CSS selectors targeting `#uno-box-manager` or `#uno-notify-container` existing in the DOM right after page load.

**New Behavior:**
These elements will **not** exist in the DOM until the first modal is opened or notification is shown.

**Migration Action:**
If you have custom CSS or JS that targets these containers on page load, update it to handle their absence or wait until they are created.

### 3. Tab Hash Sync

**Old Behavior:**
The Tab component would listen to `hashchange` events globally as soon as Uno loaded, even if no tabs were present on the page.

**New Behavior:**
Hash synchronization is now only enabled if:
1. You use the `x-uno-tab` directive.
2. You manually call `Alpine.store('uno').tab.register()`.
3. You manually call `Alpine.store('uno').tab.activate()`.

**Migration Action:**
No action required for standard usage. This is an improvement that prevents side effects on pages without tabs.

## Upgrade Steps

1. **Update Package:**
   ```bash
   npm install uno-ui@0.7.0
   ```
   Or update your CDN links:
   ```html
   <script src="https://unpkg.com/uno-ui@0.7.0/dist/uno.umd.js"></script>
   ```

2. **Verify Usage:**
   Check your codebase for any manual `init()` calls and remove them.

3. **Test:**
   Verify that Modals, Notifications, and Tabs still work as expected. The user experience should be identical, but with better initial performance.
