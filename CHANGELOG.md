# Changelog

All notable changes to this project will be documented in this file.

## [0.9.0] - 2026-08-26

### ✨ Features

- **Box**: Added `extraButtons` support — custom buttons with `position`, `buttonClass`, and async `onClick` (`return false` keeps the dialog open). New `handleExtraButton()` store method and `extraBtn` default style.
- **Notify**: New public `beginProgress()` / `endProgress()` API for external progress tracking with concurrent-request counting.

### 🐛 Fixes

- **Tooltip**: The `offset` option is now honored by `calculatePosition()` (previously parsed but ignored); horizontal placements now use the same default gap (8px) as vertical ones.
- **Tooltip**: The arrow element is now actually rendered — `calculatePosition()` returns an arrow anchor relative to the tooltip box and the template positions an 8px rotated square from it.
- **Tooltip**: Fixed `updatePosition()` reading a non-existent `placement` field from `calculatePosition()` results; the returned field is `position`, so auto-flip updates now reach the store.
- **Tooltip**: Direct string content is now sanitized like target-sourced content (XSS hardening).
- **Notify**: Fixed `remove()` writing `visible` while the template binds `x-show` to `isVisible` — leave transitions now actually play.
- **Box**: `confirm()` / `prompt()` now compose user-supplied `onOk` / `onCancel` / `onClose` with the internal Promise resolver instead of letting them overwrite it (which left the promise pending forever). An `onOk` returning `false` (or resolving to `false`) keeps the dialog open, matching `handleOk` semantics.

### 📚 Documentation & Types

- `TooltipPlacement` type narrowed to the four values actually resolved by `calculatePosition()` (`'top' | 'bottom' | 'left' | 'right'`); other strings fall back to `'bottom'` at runtime.
- Removed the long-deprecated `window.UnoConfig` global type declaration; use `Uno.setConfig()`.
- Rewrote `doc/tab-api.md` to document the real Store-First API (`x-uno-tab` + `Alpine.store('uno').tab`).
- Fixed README customization examples (real `ui` structure, `Uno.setConfig()` instead of `window.UnoConfig`).

## [0.8.0] - 2026-01-22

### 🚀 Major Improvements

- **Valid Plugin Overhaul**:
  - **Array-based Configuration**: Validation rules are now defined as an ordered array (e.g., `['required', 'email']`), ensuring predictable execution order.
  - **Debounce Support**: Added `debounce` option (e.g., `debounce: 500`) for async or heavy validators to reduce server load.
  - **Enhanced Custom Validators**: Improved API for functional, regex, async, and composite validators.

### ⚠️ Breaking Changes

- **Valid Plugin**:
  - **HTML5 Attributes Removed**: Support for defining rules via standard HTML5 attributes (e.g., `required`, `pattern=""`) has been **removed**. All rules must be explicitly defined in the `x-uno-valid` directive.
  - **Strict Configuration**: `x-uno-valid` now expects an Array (or a single string/object which is normalized to an Array).

### ✨ Features

- **Valid Plugin**:
  - `match` rule now supports `target` parameter (alias for `selector`) for easier field comparison (e.g., Password Confirmation).
  - `emptyQuit` option (enabled by default) allows skipping validation for empty optional fields.
  - `realtime: false` option to disable validation on input events for specific rules.

### 🐛 Fixes

- Fixed `valid` plugin not respecting `realtime: false` for input events.
- Fixed `valid` plugin `match` rule issues with password confirmation.
- Improved JSDoc comments for better IDE support.

## [0.7.0] - 2026-01-20

### 🚀 Major Improvements

- **Lazy Initialization**: The core plugins (`box`, `notify`, `tab`, `tooltip`) have been refactored to support lazy initialization. This means:
  - Global DOM elements (containers, overlays) are no longer injected immediately upon application startup.
  - Resources are allocated only when a component is actually used (e.g., calling `box.open()`, `notify.show()`, or using directives like `x-uno-tab`).
  - Improved initial page load performance and reduced DOM footprint.

### ⚠️ Breaking Changes

- **Internal API**: The `init()` methods on store instances (e.g., `Alpine.store('uno').box.init()`) no longer trigger immediate DOM injection. They are now effectively no-ops or handle internal state setup only. If you were manually calling these methods to force DOM injection, this will no longer work (nor is it necessary).
- **Tab Component**: The Tab component now only starts listening for `hashchange` events after the first tab is registered (via directive or JS). This prevents Uno from interfering with hash routing on pages that don't use tabs.

### 🐛 Fixes

- Fixed an issue where `ensureContainer` was being called eagerly, causing unnecessary DOM mutations.

## [0.6.0] - 2025-12-18

### ✨ Features

- **TypeScript Support**: Full TypeScript definitions included.
- **Valid Component**: New validation engine.
- **I18n**: Enhanced internationalization support.

### 📦 Chore

- Updated dependencies.
