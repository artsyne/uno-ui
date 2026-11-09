# Getting Started

This guide takes you from an empty page to a working Uno dialog. It covers the two supported
loading paths, the Tailwind CSS setup, and the optional integrations.

**Status:** Uno UI is at version `0.10.0` and has not been published to npm yet. The
repository's committed `dist/` bundles are the reference artifacts. Commands and URLs that
require a published package are marked **after first publish**.

## Prerequisites

| Requirement | Version | Notes |
| :--- | :--- | :--- |
| Alpine.js | 3.x (tested with 3.17.3) | The only runtime dependency. |
| Tailwind CSS | 4.x | Uno ships no stylesheet; Tailwind must generate the classes. |
| HTMX | 2.x (optional) | Only needed for server-driven content and HTMX helpers. |
| Node.js | >= 20.19, `.nvmrc` pins 24 | Only for building, testing and the local demo. |

## 1. Run the repository example

The fastest way to see the components is the bundled demo. It uses the committed `dist/`
bundle, so it does not require a published npm package.

```bash
git clone https://github.com/artsyne/uno-ui.git
cd uno-ui
npm ci
npm run build
npm run dev
```

`npm run dev` opens `examples/index.html` on `http://localhost:5176` (the port is fixed and
`strictPort` is enabled). The component pages are under `examples/components/`.

> `npm run build` is only necessary after changing `src/`. `dist/` is tracked in the repository
> and CI verifies that it matches the source.

## 2. Use the UMD bundle from a local copy

Uno must be served over HTTP. Opening an HTML file with `file://` prevents `fetch` and HTMX
requests from working, and may block module scripts.

```html
<!doctype html>
<html lang="en">
  <head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <title>Uno quick start</title>

    <!-- 1. Tailwind 4 (browser build, prototyping only) -->
    <script src="https://cdn.jsdelivr.net/npm/@tailwindcss/browser@4"></script>

    <!-- 2. Uno runtime: load BEFORE Alpine -->
    <script src="/vendor/uno-ui/dist/uno.umd.js" defer></script>

    <!-- 3. Your Uno configuration, still before Alpine starts -->
    <script src="/js/uno-config.js" defer></script>

    <!-- 4. Alpine.js: starts the framework and triggers Uno registration -->
    <script src="https://cdn.jsdelivr.net/npm/alpinejs@3.17.3/dist/cdn.min.js" defer></script>
  </head>
  <body>
    <button x-data @click="$store.uno.box.alert('Hello from Uno')">Open dialog</button>
  </body>
</html>
```

`/js/uno-config.js` is optional. It runs after `uno.umd.js` and before Alpine, so both
`window.Uno` and the pre-initialization configuration are ready:

```javascript
// /js/uno-config.js
window.Uno.setConfig({
  debug: false,
  ui: {
    base: 'rounded-lg shadow-md border border-gray-200 dark:border-gray-700'
  }
})
```

**Why this order matters:** the deferred UMD bundle evaluates before Alpine only if the
`<script>` for Uno comes first. At that moment `window.Alpine` does not exist, so Uno waits
for `alpine:init`. When Alpine later starts, Uno registers its stores and directives, then
Alpine initializes the page. Loading Alpine first is different: Alpine starts and walks the DOM
before Uno registers `x-uno-*` directives, so the components on that page are ignored. See
[Alpine loading paths](#alpine-loading-paths) below for the supported alternative.

## 3. Use ESM and a bundler

**After first publish:**

```bash
npm install @artsyne/uno-ui alpinejs
```

Until then, point the import at the local build output or install the repository as a file
dependency:

```bash
npm install /path/to/uno-ui
```

```javascript
// main.js
import Alpine from 'alpinejs'
import Uno from '@artsyne/uno-ui'

// Configure before Alpine starts.
Uno.setConfig({
  debug: false,
  ui: { base: 'rounded-lg shadow-md' }
})

// Your Alpine components.
Alpine.data('myComponent', () => ({
  init() {
    // ...
  }
}))

// Expose Alpine and start it. Uno registers during Alpine's alpine:init.
window.Alpine = Alpine
Alpine.start()
```

The package's `sideEffects` field keeps `src/index.js` and `dist/*` from being dropped by the
bundler. Import Uno for its side effects even if you only use the directives.

## Alpine loading paths

Uno registers in one of two ways. Both are supported; they are not interchangeable at runtime.

### Path A: Uno is evaluated before Alpine (recommended for UMD)

Uno sets `window.Uno`, registers a one-time `alpine:init` listener, and waits. Alpine starts,
fires `alpine:init`, and Uno calls `Uno.init(window.Alpine)`.

```html
<script src="/dist/uno.umd.js" defer></script>
<script src="/js/uno-config.js" defer></script>
<script src="https://cdn.jsdelivr.net/npm/alpinejs@3.17.3/dist/cdn.min.js" defer></script>
```

### Path B: `window.Alpine` already exists when Uno is evaluated

This path applies when another script or a dynamic import sequence has already set
`window.Alpine` before Uno is evaluated. Uno registers immediately; `Alpine.start()` must still
run after that registration.

```javascript
// This path is useful when another script has already set window.Alpine.
const { Uno } = await import('@artsyne/uno-ui')

// Uno has already initialized here; UI config still applies to later actions.
Uno.setConfig({ ui: { base: 'rounded-lg shadow-md' } })

window.Alpine.start()
```

Because initialization happens at import time in this path, `debug` and `onError` have already
been applied from the defaults. Use Path A when you need to set those before initialization.

Do not call `Alpine.start()` before Uno has been imported or before the Uno script has run.
Alpine would initialize the DOM without the `x-uno-*` directives, and those elements would not
be revisited.

## Tailwind CSS 4

Uno does not publish a compiled CSS file. Its component classes live inside JavaScript
templates and the default `ui` configuration, so Tailwind has to scan those files.

### Production build

Add an `@source` entry that points at the installed package's `src` directory. The package
ships `src/` in its `files` list, so this works from `node_modules`.

```css
/* app.css */
@import 'tailwindcss';

/* Uno templates and default `ui` class strings */
@source '../node_modules/@artsyne/uno-ui/src';

/* Your own templates and scripts */
@source '../src';
```

If you vendor a built copy instead, point at the directory that contains `uno.es.js` or
`uno.umd.js`:

```css
@source '../vendor/uno-ui/dist';
```

### Browser build

The `@tailwindcss/browser@4` script scans the DOM at runtime, so classes introduced by Uno's
injected templates are generated when the template is inserted. This is convenient for demos
and prototypes; use the build integration for production.

### Dark mode

Tailwind 4's default `dark:` variant follows `prefers-color-scheme`. Uno's defaults ship both
light and `dark:` classes. If your application uses a class-based theme, declare the variant:

```css
@import 'tailwindcss';
@custom-variant dark (&:where(.dark, .dark *));
```

Toggle the `dark` class on `<html>` (or a wrapper) from your own theme code. Uno does not own
the page theme.

## HTMX (optional)

HTMX 2 is an optional peer dependency. Load it before Uno so that Uno can install its
`htmx:afterSwap` hook when it initializes:

```html
<script src="https://cdn.jsdelivr.net/npm/htmx.org@2.0.10/dist/htmx.min.js"></script>
<script src="/dist/uno.umd.js" defer></script>
<script src="https://cdn.jsdelivr.net/npm/alpinejs@3.17.3/dist/cdn.min.js" defer></script>
```

Without HTMX, the directives and store APIs still work; only `hxUrl`, `hx-get`-driven tabs and
the `Uno.htmx` helpers are unavailable. See [HTMX integration](htmx.md).

## CSP nonce

Configure the nonce before the first component creates a container:

```javascript
window.Uno.setConfig({
  csp: { nonce: window.__CSP_NONCE__ }
})
```

Uno applies the nonce to the containers it creates and to `<style>` elements injected through
`Uno.dom.injectStyles()`. It does not make Alpine's expression evaluator CSP-safe: the standard
Alpine build needs `'unsafe-eval'` for `x-data`/`x-on` expressions. See the
[CSP guide](csp.md) before promising a strict CSP deployment.

## Internationalization

All eleven bundled locales are registered at startup. Switch languages after Alpine is
running:

```javascript
Alpine.store('i18n').setLang('zh-CN')
```

The choice is stored in `localStorage` under `uno.lang`. See [i18n](i18n.md) for the locale
list, detection order, placeholders and custom locales.

## Theming

Override component classes with `Uno.setConfig({ ui: ... })`. The configuration is deep-merged
and supports the `{base}` token and the `+` append prefix:

```javascript
Uno.setConfig({
  ui: {
    base: 'rounded-lg shadow-md border border-gray-200 dark:border-gray-700',
    box: {
      okBtn: '+min-w-24'
    }
  }
})
```

See [Configuration](configuration.md) and [Theming](theming.md) for the full token map.

## Verify the setup

1. Open the page in a browser.
2. Check the console for `[Uno Error]` or `[Uno Warn]` output (enable `debug: true` first).
3. Click the button. `$store.uno.box.alert()` should open a dialog with a close button.
4. If the dialog is unstyled, Tailwind is not generating the classes; see
   [Troubleshooting](troubleshooting.md).
5. If the button does nothing, inspect the script order; see
   [Alpine loading paths](#alpine-loading-paths).

## Related

- [Core concepts](concepts.md)
- [Configuration](configuration.md)
- [Theming and Tailwind](theming.md)
- [HTMX integration](htmx.md)
- [Troubleshooting](troubleshooting.md)
