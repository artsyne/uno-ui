# Uno UI 0.10.1 — Open Source Release

Uno UI is open source under the MIT License. `0.10.1` is a follow-up to the first public npm
release, refining the accessibility contracts across all five components.

| Item | Value |
| :--- | :--- |
| Package | `@artsyne/uno-ui` |
| Version | `0.10.1` |
| Release date | 2026-11-17 |
| License | MIT |
| Repository | https://github.com/artsyne/uno-ui |
| npm | https://www.npmjs.com/package/@artsyne/uno-ui |
| Website | https://artsyne.github.io/uno-ui/ |

## About the project

Uno is a lightweight UI component library for Alpine.js and Tailwind CSS 4, built around a
Store-First, Directive-Driven architecture with optional HTMX 2 integration.

The library includes Box (dialogs), Notify (toasts), Tooltip, Tab and Valid (form validation).
It ships eleven locales, TypeScript declarations and a small runtime whose only dependency is
Alpine.js.

## Install

```bash
npm install @artsyne/uno-ui alpinejs
```

```html
<script src="https://cdn.jsdelivr.net/npm/@artsyne/uno-ui@0.10.1/dist/uno.umd.js" defer></script>
```

Load Uno before Alpine. See [Getting Started](docs/getting-started.md) for the full setup,
including Tailwind CSS 4 `@source` configuration.

## 0.10.1 highlights

- **Box**: accessible name resolution with `labelledBy`, a localized fallback, native `inert`
  isolation of the page background for the topmost modal, and no backdrop/trap for non-modal
  dialogs.
- **Tooltip**: hover hints also respond to keyboard focus; all triggers associate
  `aria-describedby`; delayed hides are cancelled when switching triggers.
- **Notify**: `error()` is persistent by default; auto-dismiss pauses on hover/focus and resumes
  with the remaining time; live regions exclude action buttons; focus hands off on close.
- **Tab**: anchor tabs activate on Space with scroll prevention and key-repeat guarding.
- **Valid**: error-summary links carry the field name and target the real control, with
  reduced-motion-aware scrolling.

## Project highlights

- Five accessible, dependency-light components.
- Store-First API with lazy container initialization.
- Optional HTMX 2 integration using XHR events only.
- Eleven bundled locales and `x-t` translations.
- TypeScript definitions included in the package.
- Component-level axe scans and keyboard contract tests in CI.

## Project status

Uno is pre-1.0. The public contract is the documented store and directive API; internal
templates and containers may change between minor versions. Security reports should follow
[SECURITY.md](SECURITY.md).

## Documentation

- [English README](README.md)
- [中文 README](README.zh-CN.md)
- [Documentation index](docs/README.md)
- [API reference](docs/README.md#using-components)
- [Migration guide](docs/migration.md)
- [Changelog](CHANGELOG.md)

## License

MIT. See [LICENSE](LICENSE).
