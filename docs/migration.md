# Migration Guide

This page summarizes behaviour changes that affect upgrades. It is derived from
[CHANGELOG.md](../CHANGELOG.md) and the current source; dates are the ones already recorded in
the changelog.

## Version map

| Version | Date | Main theme |
| :--- | :--- | :--- |
| 0.10.0 | 2026-11-09 | Accessibility contracts, Box fill mode, Tooltip offset function and sizing |
| 0.9.0 | 2026-08-26 | Box extra buttons, Notify progress tracking, Tooltip fixes |
| 0.8.0 | 2026-01-22 | Valid rule arrays (breaking), debounce, match target |
| 0.7.0 | 2026-01-20 | Lazy component initialization |
| 0.6.0 | 2025-12-18 | TypeScript definitions, Valid component, i18n improvements |

TypeScript definitions have been part of the package since 0.6.0 according to the changelog and
the `types/index.d.ts` header. Older README text that said "0.7.0+" was inaccurate.

## Upgrading to 0.10.0

0.10.0 is primarily an accessibility and layout release. The public store and directive names
are unchanged.

### Accessibility

- Box now restores focus to the opener when a dialog closes, and marks underlay dialogs
  `inert` + `aria-hidden`. If you were manually moving focus after closing, remove that work.
- The Box focus trap wraps in DOM order, so a neutral target (the title) no longer lets Tab
  escape.
- Tab implements the APG tabs pattern: `role="tab"`, `aria-selected`, roving `tabindex`,
  Arrow/Home/End navigation, and per-tab/shared panels. If you were adding `aria-selected`
  yourself, the directive now owns it.
- Tooltip click popups set `aria-haspopup`, `aria-expanded`, `aria-describedby`, move focus into
  the popup, and restore focus on close.
- Notify announcement politeness is graded by type: `error`/`warning` use `alert` +
  `assertive`; `success`/`info` use `status` + `polite`.
- Valid generated errors get `role="alert"`, a `${field}-error` id, `aria-describedby` on the
  control, and `aria-invalid` sync. `showCounter` counters get `aria-describedby` and an
  `aria-label`. Multi-field failures render a focusable error summary.
- Box close button, Valid counter and icon hit areas were adjusted for contrast and target size.

### Box

- New `fill: true` mode. The content area, body and HTMX container become a flex column chain;
  without `height`/`heightClass`, `containerFill` is used. Loaded content must continue the
  chain (`flex-1 min-h-0`). See [Fill mode](box-api.md#fill-mode).
- `onContentLoad` still runs after an HTMX swap; focus is recalculated after it.

### Tooltip

- `offset` accepts a function `({ triggerEl, tooltipEl, placement }) => number`, re-evaluated on
  every positioning pass. Throwing or non-numeric results fall back to the default gap.
- The popup size is observed with `ResizeObserver`; late content/fonts re-run positioning.
- Viewport clamping never returns a negative horizontal coordinate; a popup wider than the
  viewport stays at the 10px safe margin.

### Valid

- No rule syntax change in 0.10.0; the accessibility wiring is additive.

## Upgrading to 0.9.0

### Box

- `extraButtons` support was added. Each item takes `text`, optional `position`
  (`'before'`/`'after'`), `buttonClass`, and an `onClick` that may return `false` (or a Promise
  resolving to `false`) to keep the dialog open.
- `confirm()` and `prompt()` now compose user `onOk`/`onCancel`/`onClose` callbacks with the
  internal Promise resolver instead of letting them overwrite it.
- `handleExtraButton()` was added to the store.

### Notify

- `beginProgress()` / `endProgress()` were added for concurrent progress tracking. `startProgress()`
  / `stopProgress()` remain the direct, uncounted pair.
- `remove()` now hides via `isVisible`, matching the template binding, so leave transitions play.

### Tooltip

- `offset` is honored by `calculatePosition()`; the default gap is 8px on all sides.
- The arrow is rendered from the position result.
- `updatePosition()` reads the returned `position` field (not a removed `placement` field), so
  auto-flip updates reach the store.
- Direct string content is sanitized like target-sourced content.
- `TooltipPlacement` is narrowed to `'top' | 'bottom' | 'left' | 'right'`. Other strings fall
  back to `'bottom'` at runtime.

### Removed

- `window.UnoConfig` is no longer a supported configuration path. Use `Uno.setConfig()`.

## Upgrading to 0.8.0 (Valid breaking change)

### Rules must be in `x-uno-valid`

HTML5 attributes such as `required` and `pattern` are **not** read. Move every rule into the
directive:

```html
<!-- Old: do not rely on HTML5 attributes -->
<input name="email" required type="email" />

<!-- 0.8+: explicit array -->
<input name="email" type="email" x-uno-valid="['required', 'email']" />
```

The directive normalizes a single string or object to an array, but the array form is the
documented contract and preserves order.

### New rule options

- `debounce: 500` delays input validation for expensive/async validators.
- `emptyQuit` (default `true`) skips non-required rules for empty values.
- `realtime: false` disables input-time validation for a rule.
- `match` supports `target` as an alias for `selector`; the implementation also accepts
  `rule.match`.
- `registerValidator()` plus `{ type: 'validator', name: '...' }` is the recommended custom
  validator path.

See [Valid API](valid-api.md).

## Upgrading to 0.7.0

Components are lazily initialized:

- Dialog, notification and tooltip containers are no longer injected at startup. They are
  created on first use.
- Calling a store's `init()` method manually no longer forces DOM injection. Remove manual
  `init()` calls; they are unnecessary.
- The Tab `hashchange` listener starts only after the first tab is registered or activated, so
  pages without tabs are not affected by hash routing.

## Upgrading to 0.6.0

- TypeScript definitions were introduced.
- The Valid component was introduced.
- i18n support was expanded.

## Upgrade checklist

1. Search for `window.UnoConfig`, manual `Uno.init(...)`, and HTML5 validation attributes.
2. Move configuration to `Uno.setConfig()` before Alpine starts.
3. Replace positional `tab.activate('name')` with `activate({ name, group })`.
4. If you override ARIA in Tab, remove the manual `aria-selected`/`tabindex` bindings; the
   directive owns them.
5. Rebuild your Tailwind CSS; new classes were introduced in 0.10.0.
6. Run the full suite:

```bash
npm ci
npm run lint
npm run test:run
npm run build
git diff --exit-code dist
```

7. Manually check focus return, Esc/overlay dismissal, error summaries and the HTMX flows you
   use.

## Related

- [CHANGELOG.md](../CHANGELOG.md)
- [Valid API](valid-api.md)
- [Box API](box-api.md)
- [Tooltip API](tooltip-api.md)
- [Tab API](tab-api.md)
