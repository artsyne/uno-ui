# Uno Documentation

Uno is a lightweight UI component library for Alpine.js and Tailwind CSS 4, with optional HTMX
2 integration. The package is at `0.10.0` and is not published to npm yet; the repository's
committed `dist/` bundles are the reference artifacts.

The repository landing page ([README.md](../README.md)) is in Chinese. All other documentation
is in English.

## Choose a path

### New to Uno

1. [Getting Started](getting-started.md) — run the demo, load the UMD bundle, set up Tailwind.
2. [Core Concepts](concepts.md) — Store, Directive, Template, Layer, Action, lifecycle.
3. [Box API](box-api.md) — the quickest component to see working.

### Using components

- [Box API](box-api.md) — dialogs, alerts, confirms, prompts, fill mode.
- [Notify API](notify-api.md) — toasts, positions, progress bar.
- [Tooltip API](tooltip-api.md) — hover, click and focus popups.
- [Tab API](tab-api.md) — tab groups, hash sync, HTMX panels.
- [Valid API](valid-api.md) — rules, async validators, counters, error summaries.

### Integrating

- [Configuration](configuration.md) — four-level merge, defaults, `ui` tokens.
- [Theming and Tailwind](theming.md) — `@source`, dark mode, class overrides.
- [i18n](i18n.md) — eleven locales, detection, `extend`, `loadLocale`.
- [HTMX integration](htmx.md) — XHR events, helpers, `afterSwap`, 422 errors.
- [CSP](csp.md) — nonce scope and Alpine's evaluation requirement.
- [TypeScript Usage](typescript-usage.md) — types, tsconfig, known gaps.

### Operating and upgrading

- [Troubleshooting](troubleshooting.md) — symptom-based fixes.
- [Error Handling](error-handling.md) — global handlers, `onError`, HTMX failures.
- [Migration Guide](migration.md) — 0.7/0.8/0.9 → 0.10.
- [Accessibility](accessibility.md) — implemented semantics, keyboard, test scope.
- [Test Reports](test-reports.md) — commands, axe scope, CI.

### Contributing and maintaining

- [CONTRIBUTING.md](../CONTRIBUTING.md) — local setup, tests, `dist` policy, release.
- [SECURITY.md](../SECURITY.md) — private reporting, scope, limitations.
- [CHANGELOG.md](../CHANGELOG.md) — release history.
- Issue and PR templates are in [.github](../.github/).

## Documentation conventions

- API pages use the same order: purpose, minimal example, directive, Store API, options,
  events/callbacks, accessibility, notes, related links.
- Option tables use **Name / Type / Default / Description**, and defaults are taken from
  `src/core/uno.js` and the plugin factories.
- Code blocks are language-tagged.
- "After first publish" marks commands or URLs that need a published npm package.
- The docs do not claim WCAG conformance, browser matrices or security guarantees that the
  repository does not test.

## Status and versions

| Item | Value |
| :--- | :--- |
| Package | `@artsyne/uno-ui` |
| Version | `0.10.0` |
| npm | Not published yet |
| License | MIT |
| Node (development) | >= 20.19; CI uses 24 |
| Runtime dependencies | `alpinejs` |
| Optional peer | `htmx.org` 2.x |

## Related

- [Repository README](../README.md)
- [Contributing](../CONTRIBUTING.md)
- [Security](../SECURITY.md)
- [License](../LICENSE)
