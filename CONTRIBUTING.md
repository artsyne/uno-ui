# Contributing to Uno UI

Thanks for contributing. This guide covers the local workflow, tests, documentation, the
tracked `dist/` policy and the release process.

By participating, you agree to the [Code of Conduct](CODE_OF_CONDUCT.md).

## Prerequisites

- Node.js >= 20.19. CI uses the version in `.nvmrc` (24).
- npm (the lockfile is version 3).
- A public npm registry; `.npmrc` pins `https://registry.npmjs.org/`.

```bash
git clone https://github.com/artsyne/uno-ui.git
cd uno-ui
npm ci
```

Use `npm ci` for a clean, lockfile-exact install. `npm install` is fine for local experiments but
may update the lockfile.

## Repository layout

```text
src/
├── components/   # box, notify, tooltip, tab, valid (store / directive / template)
├── core/         # uno instance, i18n, layer, plugin registry
├── locales/      # 11 locale modules
├── utils/        # dom, position, validation, a11y, log, object
└── index.js      # public entry, plugin registration, Alpine auto-registration
types/            # public TypeScript declarations
tests/            # Vitest suites and the axe helper
examples/         # browser demos using the built dist/
docs/             # documentation (English)
dist/             # tracked build output
```

## Development commands

| Command | Purpose |
| :--- | :--- |
| `npm run dev` | Vite dev server on fixed port 5176; opens `examples/index.html`. |
| `npm run build` | Build `dist/uno.es.js` and `dist/uno.umd.js`. |
| `npm run docs:dev` | VitePress dev server for the documentation site. |
| `npm run docs:build` | Build the static site into `docs/.vitepress/dist/`. |
| `npm run test` | Vitest watch mode. |
| `npm run test:run` | Single test run; CI entry point. |
| `npm run lint` | ESLint. |
| `npm run format` | Prettier write. |

## The `dist/` policy

`dist/` is committed and is part of the published package.

- Do not edit `dist/` by hand.
- After any `src/` change, run `npm run build` and commit the resulting `dist/` files with the
  source change.
- CI runs `npm run build` and then `git diff --exit-code dist`. A source change with a stale
  bundle fails.
- Documentation-only changes must not modify `dist/`.

```bash
npm run build
git diff --exit-code dist
```

## Tests

- Vitest, with `happy-dom` as the default environment. Axe suites opt into jsdom per file.
- Component tests live in `tests/components/`, core in `tests/core/`, utilities in
  `tests/utils/`.
- `tests/types/types.test.ts` is not collected by `npm run test:run` (the include pattern is
  `tests/**/*.test.js`). Type changes should be checked with a separate config.
- There is no coverage provider configured and no coverage threshold. Do not claim a coverage
  percentage in documentation or PRs.
- Add tests for behavior changes and edge cases. Accessibility contracts belong in tests, not
  only in docs.

```bash
npm run test:run
npm run test:run -- tests/components/box.test.js
npm run lint
```

## Coding style

- Prettier config: no semicolons, single quotes, 2 spaces, no trailing commas, width 100.
- ESLint uses `@eslint/js` recommended plus `eslint-config-prettier`; `dist/` is ignored.
- Modern ESM syntax; browser globals plus `Alpine` and `htmx` are declared read-only.
- Public functions and stores use JSDoc.
- camelCase for variables/functions, PascalCase for classes/components, UPPER_CASE for
  constants, kebab-case for file names.

## Architecture

Uno follows **Store-First, Directive-Driven**:

1. **Store** — state and behavior live in `Alpine.store('uno').<component>`. The plugin
   registry re-registers the whole `uno` store so nested properties stay reactive.
2. **Directive** — `x-uno-*` reads its expression, attaches listeners, calls the store and cleans
   up with Alpine's `cleanup`.
3. **Template** — markup lives in `*.template.js` and is injected into a lazily created
   container. Templates use Alpine expressions and the global `Uno` object.
4. **Configuration** — `Uno.getConfig(name, instanceOpts)` merges plugin defaults,
   `plugins.<name>`, top-level `<name>` and instance options, in that order.
5. **Layer** — z-index constants live in `src/core/layer.js`; `ui.layers` overrides the three
   manager values.

Read [Core Concepts](docs/concepts.md) and [Configuration](docs/configuration.md) before adding a
component.

## Adding or changing a component

1. Create or update `src/components/<name>/`:
   - `index.js` — plugin factory, default registration, store/directive registration.
   - `<name>.store.js` — state and public actions.
   - `<name>.directive.js` — directives and DOM listeners.
   - `<name>.template.js` — template strings, if the component renders DOM.
2. Register plugin defaults with `Uno.registerPluginDefaults(name, defaults)`.
3. Register the store with `registerPlugin(Alpine, name, store)`.
4. Add tests under `tests/components/`. Include accessibility behavior in the tests.
5. Update the public types in `types/index.d.ts`.
6. Add or update the API page under `docs/` using the standard API order:
   purpose → minimal example → directive → store API → options → events/callbacks →
   accessibility → notes → related links.
7. Add a demo page under `examples/components/` and link it from `examples/index.html`.
8. Update `CHANGELOG.md` under `Unreleased` or the target version.
9. Run `npm run build` and commit `dist/`.

## Documentation

- `README.md` is the English landing page; `README.zh-CN.md` is its Chinese translation. All
  `docs/` pages are English.
- `docs/` is the documentation hub; keep `docs/README.md` links current.
- API option tables use **Name / Type / Default / Description** and defaults must match
  `src/core/uno.js` and the plugin factories.
- Do not document APIs, defaults, browser support or performance numbers that the source and
  tests do not support.
- Keep npm/CDN examples aligned with the published version; mark future versions explicitly.
- Do not reference private project docs or internal repositories.
- Prefer one canonical page per topic; link instead of duplicating a section.

## Commit messages

Use [Conventional Commits](https://www.conventionalcommits.org/):

```text
feat: add component X
fix: correct focus restore in nested dialogs
docs: document Valid error summary
test: cover tab RTL arrow handling
refactor: extract shared position helper
perf: avoid duplicate HTMX listeners
chore: release 0.10.0
```

Keep one logical change per commit. Do not mix formatting-only changes with behavior changes.

## Pull requests

1. Branch from `main`.
2. Keep the PR focused; update tests and docs.
3. Run `npm run lint` and `npm run test:run`.
4. If `src/` changed, rebuild and commit `dist/`.
5. Fill in the PR template: summary, type, tests, docs, `dist` confirmation.
6. Link related issues.
7. Wait for CI to pass.

Reviewers check behavior, tests, accessibility impact, documentation accuracy and the `dist/`
diff.

## Release process

The package is `@artsyne/uno-ui` under the `@artsyne` npm scope. `0.10.0` is published; subsequent
releases follow the process below.

Before tagging:

1. Confirm `main` is green and the working tree is clean.
2. Update the version in `package.json` (and the lockfile) and add the matching heading to
   `CHANGELOG.md`.
3. Run the full check:

```bash
npm ci
npm run lint
npm run test:run
npm run build
git diff --exit-code dist
npm pack --dry-run
```

4. Commit `chore(release): x.y.z`.
5. Create the tag `vX.Y.Z` and push both the commit and the tag.
6. Publish from a provenance-capable environment:

```bash
npm login            # account with access to the @artsyne scope
npm publish --access public --provenance
```

- `--access public` is required for a scoped package that should be public.
- `--provenance` requires a supported CI/OIDC environment; do not assume a local shell produces
  a valid attestation.
- `prepublishOnly` reruns `test:run` and `build`, so a stale `dist/` cannot be published.
- Never publish from a dirty tree or a branch other than the release commit.

7. Verify `npm view @artsyne/uno-ui version` and create the GitHub release.

2FA is required on the publishing account. Do not share tokens in issues or pull requests.

## Security

- Do not open a public issue for a vulnerability. Follow [SECURITY.md](SECURITY.md).
- Keep dependencies minimal; the only runtime dependency is `alpinejs`.
- Sanitize untrusted content on the server or with a dedicated sanitizer; Uno's regex sanitizer
  is not a general-purpose protection.
- CSP nonce support does not make Alpine's expression evaluator CSP-safe. See
  [CSP](docs/csp.md).

## Getting help

- Documentation: start at [docs/README.md](docs/README.md).
- Bugs and features: use the issue templates.
- Security: [SECURITY.md](SECURITY.md).
- Code of conduct reports: see [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md).

## License

By contributing, you agree that your contributions are licensed under the MIT License.
