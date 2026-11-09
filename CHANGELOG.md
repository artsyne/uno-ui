# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and this
project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Documentation

- Restructured the documentation under `docs/`, added getting-started, concepts, configuration,
  i18n, theming, HTMX, CSP, accessibility, troubleshooting and migration guides, rewrote the
  component API pages against the 0.10.0 source, and added community issue/PR templates and a
  security policy. The root README now defaults to English with a Chinese translation; all other
  documentation is English.

## [0.10.0] - 2026-11-09

### Added

- **Valid**: error elements get `role="alert"` and a stable `${field}-error` id; the control is
  wired with `aria-describedby` and `aria-invalid`.
- **Valid**: `showCounter` counters are linked to the input via `aria-describedby` and carry a
  full-sentence `aria-label`.
- **Valid**: multiple field failures render a focusable error summary with jump links, keeping
  the inline errors.
- **Notify**: announcement politeness is graded by type — `error` / `warning` keep `alert` +
  `assertive`, `success` / `info` use `status` + `polite`.
- **Notify** / **Box**: close-button accessible names come from the Uno i18n keys.
- **Tab**: WAI-ARIA APG tabs pattern — `role="tab"`, `aria-selected`, roving `tabindex`,
  Arrow/Home/End navigation, per-tab and shared panels.
- **Box**: closing restores focus to the opener, and underlay dialogs are marked `inert` +
  `aria-hidden`.
- **Box**: `fill` mode passes the remaining height down to loaded content.
- **Tooltip**: `trigger: 'click'` popovers implement the full keyboard contract
  (`aria-haspopup`, `aria-expanded`, `aria-describedby`, focus into the popup, Arrow/Home/End
  navigation, Esc/outside-click close with focus restore).
- **Tooltip**: the popup size is observed after opening, so content injected later (HTMX, late
  fonts, language switches) re-runs collision-aware positioning.
- **Tooltip**: `offset` accepts a function `({ triggerEl, tooltipEl, placement }) => number`,
  re-evaluated on every positioning pass.

### Fixed

- **Tooltip**: viewport clamping never returns a negative horizontal coordinate.
- **Box**: the focus trap is DOM-order aware, so Tab no longer escapes a modal whose only
  tabbable element precedes the neutral focus target.
- **Box**: the floating close button reaches ≥3:1 token contrast in light mode.
- **Valid**: counters, focus rings and spinners use AA-oriented colour tokens in both themes.
- **Valid**: `x-uno-valid` binding-tag validation.
- Icon-only controls: hit areas widened to at least 24×24 CSS px.

### Tests

- Component-level axe scanning over all five components with real Alpine hydration, wired into
  CI.

## [0.9.0] - 2026-08-26

### Added

- **Box**: `extraButtons` support — custom buttons with `position`, `buttonClass`, and async
  `onClick` (`return false` keeps the dialog open). New `handleExtraButton()` store method and
  `extraBtn` default style.
- **Notify**: public `beginProgress()` / `endProgress()` API for external progress tracking with
  concurrent-request counting.

### Fixed

- **Tooltip**: the `offset` option is now honored by `calculatePosition()`; horizontal
  placements use the same default gap (8px) as vertical ones.
- **Tooltip**: the arrow element is now rendered; `calculatePosition()` returns an arrow anchor
  relative to the tooltip box.
- **Tooltip**: `updatePosition()` reads the `position` field returned by `calculatePosition()`.
- **Tooltip**: direct string content is sanitized like target-sourced content.
- **Notify**: `remove()` writes `isVisible`, so leave transitions play.
- **Box**: `confirm()` / `prompt()` compose user `onOk` / `onCancel` / `onClose` callbacks with
  the internal Promise resolver instead of letting them overwrite it. An `onOk` returning
  `false` keeps the dialog open.

### Changed

- **Tooltip**: `TooltipPlacement` narrowed to `top | bottom | left | right`; other strings fall
  back to `bottom` at runtime.
- Removed the long-deprecated `window.UnoConfig` global type declaration; use
  `Uno.setConfig()`.

### Documentation

- Rewrote `docs/tab-api.md` to document the Store-First API (`x-uno-tab` +
  `Alpine.store('uno').tab`).
- Fixed README customization examples (real `ui` structure, `Uno.setConfig()`).

## [0.8.0] - 2026-01-22

### Added

- **Valid**: array-based rule configuration with ordered execution.
- **Valid**: `debounce` option for async/heavy validators.
- **Valid**: improved custom validators (function, regex, async, composite).
- **Valid**: `match` rule supports `target` as an alias for `selector`.
- **Valid**: `emptyQuit` option (default enabled) skips non-required rules for empty fields.
- **Valid**: `realtime: false` disables input-time validation for a rule.

### Changed

- **Valid (breaking)**: HTML5 attributes such as `required` and `pattern` are no longer read.
  All rules must be declared in `x-uno-valid`.
- **Valid**: `x-uno-valid` expects an array (a single string or object is normalized to one).

### Fixed

- **Valid**: `realtime: false` is respected for input events.
- **Valid**: `match` rule issues with password confirmation.
- Improved JSDoc comments for IDE support.

## [0.7.0] - 2026-01-20

### Changed

- **Lazy initialization**: `box`, `notify`, `tab` and `tooltip` no longer inject global DOM
  containers at startup. Containers are created when a component is first used.
- **Tab**: the `hashchange` listener starts only after the first tab is registered, so pages
  without tabs are not affected by hash routing.
- **Breaking (internal)**: store `init()` methods no longer trigger DOM injection. Manual calls
  are unnecessary.

### Fixed

- Avoided eager `ensureContainer` calls that caused unnecessary DOM mutations.

## [0.6.0] - 2025-12-18

### Added

- TypeScript definitions.
- Valid component.
- Expanded i18n support.

### Changed

- Updated dependencies.
