# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

### ✨ Features

- **Valid**: error elements get `role="alert"` and a stable `${field}-error` id; the control is wired with `aria-describedby` and `aria-invalid`.
- **Valid**: `showCounter` counters are linked to the input via `aria-describedby` and carry a full-sentence `aria-label`.
- **Valid**: multiple field failures render a focusable error summary with jump links, keeping the inline errors.
- **Notify**: announcement politeness is graded by type — `error` / `warning` keep `alert` + `assertive`, `success` / `info` use `status` + `polite`.
- **Notify** / **Box**: close-button accessible names come from the Uno i18n keys.
- **Tab**: WAI-ARIA APG tabs pattern — `role="tab"`, `aria-selected`, roving `tabindex`, Arrow/Home/End navigation, per-tab and shared panels.
- **Box**: closing restores focus to the opener, and underlay dialogs are marked `inert` + `aria-hidden`.
- **Tooltip**: `trigger: 'click'` popovers implement the full keyboard contract (`aria-haspopup`, `aria-expanded`, `aria-describedby`, focus into the popup, Arrow/Home/End navigation, Esc/outside-click close with focus restore).
- **Tooltip**: the popup size is observed after opening, so content injected later (HTMX, late fonts, language switches) re-runs the collision-aware positioning instead of keeping the first measurement.
- **Tooltip**: `offset` accepts a function `({ triggerEl, tooltipEl, placement }) => number`, re-evaluated on every positioning pass (open, scroll, resize, content resize). A gap derived from live geometry — for example "panel top = container bottom + 8px" — stays correct across breakpoint switches and layout changes instead of being frozen when the trigger was configured.

### 🐛 Fixes

- **Tooltip**: viewport clamping never returns a negative horizontal coordinate, so a popup wider than the viewport stays anchored to the 10px safe margin instead of overflowing off screen.
- **Box**: the focus trap is DOM-order aware, so Tab no longer escapes a modal whose only tabbable element precedes the neutral focus target.
- **Box**: the floating close button reaches ≥3:1 contrast in light mode (WCAG 1.4.11).
- **Valid**: counters, focus rings and spinners use AA-compliant colour steps in both themes.
- **Valid**: `x-uno-valid` binding-tag validation.
- Icon-only controls: hit areas widened to at least 24×24 CSS px (WCAG 2.2 SC 2.5.8).

### 🧪 Tests

- Component-level axe scanning over all five components with real Alpine hydration, wired into CI.

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
