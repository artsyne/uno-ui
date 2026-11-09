# Uno UI 0.10.0 — Open Source Release

Uno UI is now open source under the MIT License. This is the first public npm release of
`@artsyne/uno-ui`.

| Item | Value |
| :--- | :--- |
| Package | `@artsyne/uno-ui` |
| Version | `0.10.0` |
| Release date | 2026-11-09 |
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
<script src="https://cdn.jsdelivr.net/npm/@artsyne/uno-ui@0.10.0/dist/uno.umd.js" defer></script>
```

Load Uno before Alpine. See [Getting Started](docs/getting-started.md) for the full setup,
including Tailwind CSS 4 `@source` configuration.

## Release highlights

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
