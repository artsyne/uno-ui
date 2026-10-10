# Test Reports

This document describes how the Uno test suite is organised and what its results do and do not
prove. It is an entry point, not a results snapshot: counts and timings change with every release
and are produced by the run itself.

## How to Run

| Command                      | Scope                                                                           |
| :--------------------------- | :------------------------------------------------------------------------------ |
| `npm run test:run`           | Single run of the whole suite — the entry point CI uses.                        |
| `npm run test`               | Vitest in watch mode.                                                           |
| `npm run test:run -- <path>` | One file or directory, e.g. `npm run test:run -- tests/components/box.test.js`. |
| `npm run lint`               | ESLint over the repository.                                                     |

`uno/package.json` defines no coverage script and `vitest.config.js` configures no coverage
provider, so no coverage number is maintained here.

## Layout

```text
tests/
├── components/   # box, notify, tooltip, tab, valid — behaviour, templates, ARIA, axe scans
├── core/         # uno, i18n, layer, plugin-registry
├── utils/        # dom, position, validation, a11y, log, object
├── types/        # types.test.ts — TypeScript definitions
└── helpers/      # axe.js — shared axe wrapper, not a suite
```

`vitest.config.js` runs with `environment: 'happy-dom'`, `globals: true` and
`include: ['tests/**/*.test.js']`. Two consequences are worth knowing:

- The pattern matches `.js` only, so `tests/types/types.test.ts` is not part of `npm run test:run`.
- `tests/helpers/axe.js` is imported by the component suites and is not collected as a suite.

## Component Accessibility Scans

The `.axe.test.js` suites scan Alpine-rendered components through `tests/helpers/axe.js`. That
wrapper deliberately disables the rules that depend on a full page or on a real layout engine:

| Disabled rule                                                                                      | Reason                                                      |
| :------------------------------------------------------------------------------------------------- | :---------------------------------------------------------- |
| `color-contrast`                                                                                   | Needs real layout; covered by manual review instead.        |
| `region`, `landmark-one-main`, `page-has-heading-one`, `bypass`, `html-has-lang`, `document-title` | Host-page responsibilities, not component responsibilities. |

A passing scan therefore proves component-level semantics for the configuration under test. It is
**not** a WCAG conformance statement and does not replace keyboard or screen-reader acceptance: the
component contracts live in [a11y-design.md](../../docs/a11y-design.md) §3.5, and the open review
findings in [core-review-summary.md](../../docs/core-review-summary.md) §5.

## Security-Related Checks

- Content paths that render HTML (Box body, Tooltip content) go through the shared `sanitizeHTML`
  helper in `src/utils/dom.js`.
- Dynamically created elements and `<style>` blocks receive the CSP nonce configured through
  `Uno.setConfig({ csp: { nonce } })`; see the [CSP Guide](./csp-guide.md).

## Manual Verification

```bash
npm run build
npm run dev
# The dev server opens http://localhost:5176/examples/ (fixed port, strictPort)
```

Interactive demos for every component live under `examples/components/`.

## Quality Gates

`npm run test:run` runs in the CI `js` job alongside the root project's i18n, icon and Jest checks;
see [.github/workflows/ci.yml](../../.github/workflows/ci.yml). Release history is recorded in the
[CHANGELOG](../CHANGELOG.md).
