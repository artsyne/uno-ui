# Uno UI

<p align="left">
  <a href="README.md"><img alt="English" src="https://img.shields.io/badge/Language-English-blue?style=flat-square"></a>
  <a href="README.zh-CN.md"><img alt="简体中文" src="https://img.shields.io/badge/Language-%E7%AE%80%E4%BD%93%E4%B8%AD%E6%96%87-red?style=flat-square"></a>
</p>

[![CI](https://github.com/artsyne/uno-ui/actions/workflows/ci.yml/badge.svg)](https://github.com/artsyne/uno-ui/actions/workflows/ci.yml)
[![npm](https://img.shields.io/npm/v/@artsyne/uno-ui.svg?style=flat-square)](https://www.npmjs.com/package/@artsyne/uno-ui)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg?style=flat-square)](https://opensource.org/licenses/MIT)
[![Version](https://img.shields.io/badge/version-0.10.0-blue.svg?style=flat-square)](CHANGELOG.md)

Uno is a lightweight UI component library built on **Alpine.js** and **Tailwind CSS 4**. It
follows a **Store-First, Directive-Driven** architecture and supports optional **HTMX 2**
integration.

## Status

- Version: `0.10.0`
- npm: published as `@artsyne/uno-ui@0.10.0`.
- Website: https://artsyne.github.io/uno-ui/
- Open-source release notes: [RELEASE.md](RELEASE.md).
- The committed `dist/uno.umd.js` and `dist/uno.es.js` are the current artifacts; CI verifies
  that they match `src/`.
- Uno does not ship a standalone CSS file. Tailwind generates the component classes from the
  class strings in Uno's templates.

## Features

- **Five components**: Box (dialogs), Notify (toasts), Tooltip, Tab and Valid (form validation).
- **Store-First**: control components through `Alpine.store('uno')`; directives provide the
  declarative HTML surface.
- **Small runtime**: the only runtime dependency is Alpine.js; HTMX is an optional peer.
- **Lazy initialization**: dialog, notification and tooltip containers are created on first use.
- **Eleven bundled locales**: `en`, `zh-CN`, `zh-TW`, `es`, `ja`, `de`, `fr`, `ko`, `pt`, `ru`,
  `ar`.
- **TypeScript definitions**: shipped as `types/index.d.ts`.
- **HTMX 2**: XHR-only integration for content loading, progress, form interception and `422`
  error mapping.

## Requirements

| Dependency | Version | Notes |
| :--- | :--- | :--- |
| Alpine.js | 3.x (tested with 3.17.3) | Required; the only runtime dependency. |
| Tailwind CSS | 4.x | Required; Uno does not publish CSS. |
| HTMX | 2.x | Optional. |
| Node.js | >= 20.19 (CI uses 24 from `.nvmrc`) | Development, build and tests only. |

Uno does not maintain a formal browser test matrix. The practical support range is constrained
by Alpine.js 3 and Tailwind CSS 4; check their documentation and test against your own target
browsers.

## Quick start

### Run the repository examples

```bash
git clone https://github.com/artsyne/uno-ui.git
cd uno-ui
npm ci
npm run build
npm run dev
```

`npm run dev` opens `http://localhost:5176/examples/` (fixed port, `strictPort`). Component
pages live under `examples/components/`.

### Use the local UMD bundle

```html
<!doctype html>
<html lang="en">
  <head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <title>Uno example</title>

    <!-- Tailwind 4 browser build: prototypes only -->
    <script src="https://cdn.jsdelivr.net/npm/@tailwindcss/browser@4"></script>

    <!-- Load Uno before Alpine -->
    <script src="/vendor/uno-ui/dist/uno.umd.js" defer></script>
    <script src="/js/uno-config.js" defer></script>
    <script src="https://cdn.jsdelivr.net/npm/alpinejs@3.17.3/dist/cdn.min.js" defer></script>
  </head>
  <body>
    <button x-data @click="$store.uno.box.alert('Hello from Uno')">Open dialog</button>
  </body>
</html>
```

```javascript
// /js/uno-config.js
window.Uno.setConfig({
  debug: false,
  ui: {
    base: 'rounded-lg shadow-md border border-gray-200 dark:border-gray-700'
  }
})
```

**Load order matters**: Uno must register before Alpine starts. The order above (Uno first,
Alpine second) is the UMD path. For ESM and the second registration path, see
[Getting Started](docs/getting-started.md#alpine-loading-paths).

### Let Tailwind 4 scan Uno classes

Production builds must add the package's `src` directory to Tailwind's `@source`; otherwise the
component classes are not generated:

```css
@import 'tailwindcss';
@source '../node_modules/@artsyne/uno-ui/src';
```

See [Theming](docs/theming.md) for dark mode, the browser build and custom themes.

### npm and CDN

```bash
npm install @artsyne/uno-ui alpinejs
```

```html
<script src="https://cdn.jsdelivr.net/npm/@artsyne/uno-ui@0.10.0/dist/uno.umd.js" defer></script>
```

Both jsDelivr and unpkg serve the published tarball. The repository `dist/` remains available for
local development and unreleased changes.

## Components

| Component | Description | Documentation |
| :--- | :--- | :--- |
| **Box** | Modals, alerts, confirms and prompts; stacking, positioning, HTMX content, fill layout. | [Box API](docs/box-api.md) |
| **Notify** | Toasts in nine positions, a global progress bar, HTMX content. | [Notify API](docs/notify-api.md) |
| **Tooltip** | Hover, click and focus popups with flipping, viewport clamping and HTMX content. | [Tooltip API](docs/tooltip-api.md) |
| **Tab** | Tab groups, hash deep links, HTMX panels, APG keyboard behaviour. | [Tab API](docs/tab-api.md) |
| **Valid** | Array rules, async validators, character counters, error summaries and `422` mapping. | [Valid API](docs/valid-api.md) |

## Documentation

### Start here

- [Getting Started](docs/getting-started.md) — installation, Uno/Alpine load order, Tailwind,
  HTMX, CSP and i18n.
- [Core Concepts](docs/concepts.md) — Store, Directive, Template, Layer, Action and lifecycle.
- [Configuration](docs/configuration.md) — four-level merge, defaults, `ui` tokens, `{base}`
  and `+`.
- [Theming and Tailwind](docs/theming.md) — `@source`, dark mode and class overrides.
- [i18n](docs/i18n.md) — eleven locales, detection, `extend`, `loadLocale` and `x-t`.

### Integrations

- [HTMX Integration](docs/htmx.md) — XHR-only, events, `Uno.htmx.*`, `afterSwap` and `422`.
- [CSP](docs/csp.md) — the real nonce scope and Alpine's `'unsafe-eval'` requirement.
- [Accessibility](docs/accessibility.md) — per-component semantics, keyboard behaviour, axe
  scope and manual acceptance.
- [TypeScript Usage](docs/typescript-usage.md) — types, tsconfig and known gaps.

### Operations and upgrades

- [Troubleshooting](docs/troubleshooting.md) — symptom-based fixes.
- [Error Handling](docs/error-handling.md) — global handlers, `onError` and HTMX failures.
- [Migration Guide](docs/migration.md) — 0.7 / 0.8 / 0.9 → 0.10.
- [Test Reports](docs/test-reports.md) — test commands, axe scope and CI gates.

## Configuration and theming

```javascript
window.Uno.setConfig({
  debug: false,
  plugins: {
    box: { closeOnEsc: true },
    notify: { duration: 4000 }
  },
  ui: {
    base: 'rounded-lg shadow-md border border-gray-200 dark:border-gray-700',
    box: {
      okBtn: '+min-w-24'
    },
    layers: { box: 1400, notify: 1600, tooltip: 1700 }
  }
})
```

`{base}` is replaced with `ui.base` when a style is resolved. A leading `+` appends instead of
replacing. See [Configuration](docs/configuration.md#ui-token-reference) for the full token
table.

## Internationalization

```javascript
document.addEventListener('alpine:init', () => {
  Alpine.store('i18n').setLang('zh-CN')
})
```

The selection is stored in `localStorage` under `uno.lang`. All eleven locales are bundled and
loaded with the library. See [i18n](docs/i18n.md).

## Security

- Box `content` and Tooltip `content` pass through `Uno.dom.sanitize()`, which is a regex pass,
  not a complete HTML sanitizer. It can be bypassed.
- HTMX responses in Box, Notify and Tooltip are treated as trusted server HTML; Uno does not
  sanitize them again.
- The CSP nonce applies to Uno-created containers and `<style>` elements created by
  `injectStyles()`. It does not make Alpine's expression evaluation compatible with a strict
  CSP.
- Report vulnerabilities through the private channel in [SECURITY.md](SECURITY.md).

## Development

```bash
npm ci
npm run lint
npm run test:run
npm run build
```

- `npm run dev`: Vite dev server for the component examples.
- `npm run build`: build the npm package into `dist/`.
- `npm run docs:dev`: VitePress dev server for this documentation site.
- `npm run docs:build`: build the static site into `docs/.vitepress/dist/`.
- `npm run test:run`: single test run used by CI.
- `npm run test`: Vitest watch mode.

`dist/` is a tracked release artifact: after changing `src/`, rebuild and commit it. CI runs
`git diff --exit-code dist`.

## Project structure

```text
uno-ui/
├── dist/                 # build output (ESM + UMD)
├── docs/                 # documentation
├── examples/             # runnable component demos
├── src/
│   ├── components/       # box / notify / tooltip / tab / valid
│   ├── core/             # uno / i18n / layer / plugin-registry
│   ├── locales/          # 11 locale modules
│   ├── utils/            # dom / position / validation / a11y / log / object
│   └── index.js          # entry point
├── tests/                # unit tests, axe scans, type tests
├── types/                # TypeScript declarations
├── CONTRIBUTING.md
├── SECURITY.md
├── CHANGELOG.md
└── package.json
```

## Contributing and release

- Contribution workflow, style, tests and the `dist` policy: [CONTRIBUTING.md](CONTRIBUTING.md).
- Release history: [CHANGELOG.md](CHANGELOG.md).
- Security policy: [SECURITY.md](SECURITY.md).
- Chinese translation: [README.zh-CN.md](README.zh-CN.md).

## License

MIT License. See [LICENSE](LICENSE).
