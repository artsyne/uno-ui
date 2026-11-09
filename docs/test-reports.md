# Test Reports

This page explains how the Uno test suite is organized, what it proves and how to reproduce the
CI checks. It is an entry point, not a stored results snapshot: counts and timings change with
every release.

## Commands

| Command | Scope |
| :--- | :--- |
| `npm run test:run` | Single run of the whole suite; the CI entry point. |
| `npm run test` | Vitest in watch mode. |
| `npm run test:run -- tests/components/box.test.js` | One file (or directory). |
| `npm run lint` | ESLint over the repository. |
| `npm run format` | Prettier write. |
| `npm run build` | Vite build into `dist/`. |
| `npm pack --dry-run` | Inspect the published package contents. |

`vitest.config.js` uses `environment: 'happy-dom'`, `globals: true`, and
`include: ['tests/**/*.test.js']`. Two consequences:

- The pattern is `.js` only, so `tests/types/types.test.ts` is not collected by
  `npm run test:run`.
- `tests/helpers/axe.js` is a helper, not a suite.

The axe suites opt into jsdom with a file-level `// @vitest-environment jsdom` comment.

## Layout

```text
tests/
├── components/   # box, notify, tooltip, tab, valid: behavior, templates, ARIA, axe
├── core/         # uno, i18n, layer, plugin-registry
├── utils/        # dom, position, validation, a11y, log, object
├── types/        # types.test.ts (not collected by test:run)
└── helpers/      # axe.js: shared axe wrapper, not a suite
```

## Component accessibility scans

The `.axe.test.js` suites boot real Alpine, render the component, and scan the hydrated DOM with
`axe-core`. The wrapper disables rules that need a full page or a layout engine:

| Disabled rule | Reason |
| :--- | :--- |
| `color-contrast` | Requires real layout and computed colors. |
| `region`, `landmark-one-main`, `page-has-heading-one`, `bypass`, `html-has-lang`, `document-title` | Host-page responsibilities, not component responsibilities. |

`scanA11y()` also fails when the engine analyzed zero nodes, so an empty DOM cannot produce a
false green result.

Two suites assert color tokens rather than computed ratios:

- `box.closebtn.contrast.test.js` — floating Box close button tokens.
- `valid.counter.a11y.test.js` — counter default/over-limit tokens.

A passing scan proves component-level semantics for the rendered fixture. It is **not** a WCAG
conformance statement and does not replace keyboard or screen-reader acceptance. Keyboard
contracts have their own tests (`box.focus`, `tooltip.keyboard`, `tab.aria`, `valid.summary`,
`valid.counter.a11y`).

## Security-related checks

- HTML-bearing content paths (Box body, Tooltip content) go through the shared sanitizer in
  `src/utils/dom.js`.
- Dynamically created containers and `<style>` elements receive the configured CSP nonce. See
  [CSP](csp.md).
- HTMX responses are not sanitized by Uno; the tests do not treat them as untrusted input.

## Manual verification

```bash
npm run build
npm run dev
```

The dev server opens `http://localhost:5176/examples/` (fixed port, `strictPort`). Component
demos live under `examples/components/`.

## CI

`.github/workflows/ci.yml` runs on `main` pushes and pull requests:

1. `npm ci` with the Node version from `.nvmrc` (24).
2. `npm run lint`.
3. `npm run test:run`.
4. `npm run build`, then `git diff --exit-code dist` — the tracked `dist/` must match the
   source build.
5. `npm pack --dry-run` — the publishable package must be valid.

See [.github/workflows/ci.yml](https://github.com/artsyne/uno-ui/blob/main/.github/workflows/ci.yml).

## Coverage

`package.json` defines no coverage script and `vitest.config.js` configures no coverage provider,
so no coverage number is maintained here. Add a coverage provider only as a deliberate change,
and record the threshold in this document when you do.

## Interpreting failures

- A component axe failure usually means missing ARIA state or an invalid role/aria relationship
  in the hydrated DOM; the error includes the rule and target.
- A `git diff --exit-code dist` failure means `src/` changed without rebuilding.
- `npm pack --dry-run` failures usually mean `files`/`exports` changed or a required artifact is
  missing.
- The TypeScript test file is not part of CI. To exercise it, run it with a separate Vitest
  config whose `include` matches `tests/types/**/*.test.ts`.

## Related

- [Accessibility](accessibility.md)
- [CONTRIBUTING.md](https://github.com/artsyne/uno-ui/blob/main/CONTRIBUTING.md)
- [CHANGELOG.md](https://github.com/artsyne/uno-ui/blob/main/CHANGELOG.md)
