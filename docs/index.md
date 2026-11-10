---
layout: home

hero:
  name: Uno UI
  text: Lightweight UI for Alpine.js and Tailwind CSS 4
  tagline: Store-First, Directive-Driven components with optional HTMX 2. Eleven locales, TypeScript types and a small runtime.
  actions:
    - theme: brand
      text: Get started
      link: /getting-started
    - theme: alt
      text: Components
      link: /components
    - theme: alt
      text: GitHub
      link: https://github.com/artsyne/uno-ui

features:
  - title: Five components
    details: Box dialogs, Notify toasts, Tooltip, Tab and Valid form validation.
  - title: Small runtime
    details: Alpine.js is the only runtime dependency. HTMX 2 is an optional peer.
  - title: Store-First
    details: Control everything through Alpine.store('uno'); directives provide the HTML API.
  - title: Lazy containers
    details: Dialog, notification and tooltip DOM is created on first use.
  - title: Eleven locales
    details: en, zh-CN, zh-TW, es, ja, de, fr, ko, pt, ru and ar are bundled.
  - title: Tested accessibility
    details: Component-level axe scans and keyboard contract tests run in CI.
---

## Install

```bash
npm install @artsyne/uno-ui alpinejs
```

Or load the UMD bundle from jsDelivr:

```html
<script src="https://cdn.jsdelivr.net/npm/@artsyne/uno-ui@0.10.0/dist/uno.umd.js" defer></script>
<script src="https://cdn.jsdelivr.net/npm/alpinejs@3.17.3/dist/cdn.min.js" defer></script>
```

Load Uno before Alpine, and let Tailwind CSS 4 scan the package source:

```css
@import 'tailwindcss';
@source '../node_modules/@artsyne/uno-ui/src';
```

See [Getting Started](getting-started.md) for the full setup, including the UMD path, CSP
nonce and i18n.

## Status

`@artsyne/uno-ui@0.10.0` is published on npm and released on GitHub. The project is pre-1.0;
the documented store and directive APIs are the public contract.

- [Open-source release notes](https://github.com/artsyne/uno-ui/blob/main/RELEASE.md)
- [Changelog](https://github.com/artsyne/uno-ui/blob/main/CHANGELOG.md)
- [Contributing](https://github.com/artsyne/uno-ui/blob/main/CONTRIBUTING.md)
- [Security policy](https://github.com/artsyne/uno-ui/blob/main/SECURITY.md)
