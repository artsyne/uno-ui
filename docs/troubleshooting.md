# Troubleshooting

Start with the quick checks, then jump to the component or integration section. Most reports
come down to one of four causes: script order, Tailwind scanning, a missing HTMX 2 runtime, or
a configuration that was applied after initialization.

## Quick checks

```html
<script src="/vendor/uno-ui/dist/uno.umd.js" defer></script>
<script src="/js/uno-config.js" defer></script>
<script src="https://cdn.jsdelivr.net/npm/alpinejs@3.17.3/dist/cdn.min.js" defer></script>
```

1. Is Uno loaded before Alpine? See [Alpine loading paths](getting-started.md#alpine-loading-paths).
2. Does `window.Uno` exist?
3. Does `window.Alpine` exist, and has `Alpine.start()` run?
4. Does `Alpine.store('uno')` expose the components?
5. Is Tailwind generating the classes? Open DevTools and check a dialog's computed styles.
6. Enable `debug: true` and read the console messages.

```javascript
window.Uno.setConfig({ debug: true })
```

Debug output is applied when Uno initializes. Setting it after initialization does not turn the
logger on.

## Installation and tooling

| Symptom | Cause | Fix |
| :--- | :--- | :--- |
| `npm install @artsyne/uno-ui` fails with 404 | The requested version does not exist, or a registry mirror is stale. | Check `npm view @artsyne/uno-ui version` against the public registry, then retry. |
| CDN URL returns 404 | The version/tag has not propagated, or the URL is wrong. | Use the exact published version, or the local `dist/uno.umd.js` while developing. |
| `npm ci` fails | Node version or registry. | Use Node >= 20.19 (`.nvmrc` pins 24) and the public npm registry (`.npmrc` sets it). |
| `npm run dev` refuses to start | Port `5176` is already in use (`strictPort`). | Stop the other process or change the port locally. |

## Alpine does not initialize Uno

| Symptom | Cause | Fix |
| :--- | :--- | :--- |
| Button does nothing, no error | Alpine ran before Uno registered its directives. | Load Uno before Alpine, or import Uno before calling `Alpine.start()`. |
| `Alpine.store('uno')` is undefined | Uno never initialized. | Check script order and the console; `[Uno Error]` / `[Uno Warn]` messages are suppressed unless `debug: true`. |
| `Uno has already been initialized` warning | You called `Uno.init()` manually after auto-init. | Remove the manual call; use `Uno.setConfig()` and let registration happen. |
| Config has no effect | `setConfig()` ran after the first component used the value. | Configure before Alpine starts, or before the affected Action. `debug` and `onError` only apply at init. |

## Styles are missing

| Symptom | Cause | Fix |
| :--- | :--- | :--- |
| Dialog renders but is unstyled | Tailwind did not scan Uno's templates. | Add `@source '../node_modules/@artsyne/uno-ui/src';` to your CSS. |
| Some classes work, component classes do not | Only application files were scanned. | Add the package's `src` or `dist` directory as an `@source`. |
| Browser prototype has no styles | Tailwind browser build not loaded or loaded after the template was inserted. | Load `@tailwindcss/browser@4` and let it rescan; use a build for production. |
| Dark mode does nothing | Tailwind's `dark:` variant is media-based by default. | Use `@custom-variant dark (&:where(.dark, .dark *));` and toggle `.dark`. |
| Custom `ui` classes are missing | Class names were built dynamically. | Keep full class names as literals, or configure Tailwind to scan the file. |

Uno does not ship a CSS file. There is no "Uno stylesheet" to import.

## HTMX

| Symptom | Cause | Fix |
| :--- | :--- | :--- |
| `Uno.htmx.available()` is false | HTMX is not loaded, or not HTMX 2. | Load `htmx.org@2.x` before Uno. |
| `Uno.htmx` helpers are missing | Uno initialized before HTMX. | Load HTMX before the Uno script / before `Alpine.start()`. |
| Loaded content has no Alpine behaviour | No `initTree` after the swap. | Load HTMX before Uno so the `afterSwap` hook is installed, or call `Alpine.initTree(event.detail.elt)` yourself. |
| `hxUrl` opens a dialog but stays empty | The response is a full page, not a fragment, or the request failed. | Return an HTML fragment; check the Network tab and `htmx:responseError`. |
| Server errors do not show inline | The 422 body is not JSON or does not have `errors`. | Return `422` with `{ "errors": { "field": "message" } }`. |
| Form submits without validation | The form has no `x-uno-valid` marker and no descendant with it. | Add `x-uno-valid` to the form or to at least one field. |

See [HTMX integration](htmx.md) for the full event and helper reference.

## CSP

| Symptom | Cause | Fix |
| :--- | :--- | :--- |
| `Refused to evaluate a string as JavaScript` | Alpine's standard build needs `'unsafe-eval'`. | Adjust the policy or evaluate Alpine's CSP build (not tested with Uno). |
| `Refused to apply inline style` | A runtime style build is blocked. | Use a build-time stylesheet, or intentionally allow the source. |
| Components render but do nothing | Expressions blocked or Alpine never started. | Check the console, `window.Alpine`, `window.Uno`. |
| Nonce mismatch | Header and `Uno.setConfig()` use different values. | Generate one nonce per response and render it into both places. |

A nonce alone does not make the application CSP-safe. See [CSP](csp.md).

## i18n

| Symptom | Cause | Fix |
| :--- | :--- | :--- |
| `setLang('fr')` logs a warning | The locale is not loaded. | All eleven bundled locales are loaded; custom codes need `extend()` or `loadLocale()`. |
| A regional code does not load | `normalizeLocale()` maps `fr-CA` to `fr`, which is already loaded. | Use `extend()` for variants, or a distinct base code. |
| Text does not change | You changed `locale` directly instead of calling `setLocale()`. | Use `setLocale()` / `setLang()`. |
| `{n}` remains in the message | Missing interpolation variable. | Pass `{ n: value }` as the second argument to `t()`. |
| `<html lang>` is unchanged | Uno does not manage document attributes. | Set `lang` in your template/application. |
| A Box button is not translated | `okText`/`cancelText` are literal strings; only the close button uses i18n. | Translate the strings before passing them. |

See [i18n](i18n.md).

## Box

| Symptom | Cause | Fix |
| :--- | :--- | :--- |
| Dialog does not open | Alpine/Uno not initialized, or the directive value is not an object. | Check `Alpine.store('uno').box` and the console warning from `x-uno-box`. |
| ESC does not close | `closeOnEsc` defaults to `false`. | Pass `closeOnEsc: true`. |
| Overlay click does not close | `closeOnOverlayClick` defaults to `false`. | Pass `closeOnOverlayClick: true`. |
| `confirm()` never resolves | An `onOk` callback returned `false` (dialog stays open) or the dialog was closed through a path that does not compose the resolver. | Resolve or close via OK/Cancel; see the note in [Box API](box-api.md#store-api). |
| Content is stripped | `content` is sanitized; scripts, images and event handlers are removed. | Pass safe markup, or render complex content in a custom component. |
| Fill layout does not scroll | The loaded content chain lacks `flex-1 min-h-0`. | Continue the flex chain down to the growing element; see [Fill mode](box-api.md#fill-mode). |

## Notify

| Symptom | Cause | Fix |
| :--- | :--- | :--- |
| Toast has no background | `type` has no default. | Pass `type: 'info'` (or success/error/warning). |
| Toast never disappears | `duration: 0` or a long duration. | Set a finite `duration`. |
| Progress bar stays visible | A request never completed. | Pair `beginProgress()` with `endProgress()`, or call `stopProgress()` in `finally`. |
| Progress bar hides early with concurrent requests | You used `startProgress()` / `stopProgress()` manually. | Use `beginProgress()` / `endProgress()` for counting. |
| HTMX success toast is missing | `bindNotifyToHtmx()` not called, HTMX missing, or `X-Notify-Suppress` set. | Call it after HTMX loads; check the response headers. |

## Tooltip

| Symptom | Cause | Fix |
| :--- | :--- | :--- |
| Tooltip stays behind a dialog | Stacking context or layer. | Raise `ui.layers.tooltip`, and check ancestor `transform`/`filter`. |
| Wrong position | Anchor moved after opening. | Call `tooltip.updatePosition()` or rely on scroll/resize/size observation. |
| Flicker on hover | The pointer leaves the trigger before entering the popup. | Increase `offset` or adjust layout; the popup itself cancels hide on mouseenter. |
| Click popup does not close | `clickOutside: false`. | Keep the default `true`, or close from your own action. |
| Menu arrows do nothing | Content has no `role="menu"`. | Add APG menu semantics. |

## Tab

| Symptom | Cause | Fix |
| :--- | :--- | :--- |
| Tab does not switch | `name`/`group` mismatch between the directive and `isActive()`. | Use the same pair in both places. |
| Nothing is active | The directive expression did not evaluate. | Check the console and Alpine scope. |
| Hash does not update | `hash: false` on the instance. | Remove it; a global `tab.hash` does not override the instance default. |
| Deep-link HTMX does not load | Trigger id is not `tab-{name}` or lacks `hx-get`. | Set `id="tab-{name}"` and the HTMX attributes. |
| Panel is not associated | Panel lacks `id="panel-{name}"` and `aria-labelledby`. | Add both, or use shared-panel mode with `data-tab-group`. |
| `activate('name')` does nothing | Positional arguments are not supported. | Pass an object: `activate({ name: 'name', group: 'group' })`. |

## Valid

| Symptom | Cause | Fix |
| :--- | :--- | :--- |
| Validator does not run | Field empty and `emptyQuit` is true. | Add `required`, or set `emptyQuit: false` on the rule. |
| Validator name not found | `registerValidator()` name mismatch. | Register before the field initializes; use the same name. |
| Async validator fires on every keystroke | No `debounce`. | Add `debounce: 500`. |
| Validation runs while typing | Default trigger includes `input`. | Add `trigger: 'blur'` to the rule. |
| `required` attribute is ignored | Rules must be declared in `x-uno-valid`; HTML5 attributes are not read. | Use `x-uno-valid="['required']"`. |
| Custom error text does not show | Validator returned `false` without `msg`, or the manual error element is not bound. | Return a string, or set `msg`; bind a manual error element with `x-text`/`x-show`. |
| Generated error element missing | `autoError: false`, or no name/id on the field. | Enable `autoError`; give the field a `name` or `id`. |
| Error text class is not applied | `valid.errorClass` vs `ui.valid.helperText`. | Set `valid.errorClass`; it takes precedence. |

## Debugging

```javascript
window.Uno.setConfig({ debug: true, onError: (value) => console.error('[app]', value) })
```

```javascript
console.log(window.Uno?.version)
console.log(window.Uno?.options)
console.log(window.Alpine?.store('uno'))
console.log(window.Alpine?.store('i18n'))
```

- `debug` only affects Uno's logger; your own console output is separate.
- The global handlers report through `log.error()`. Without `debug` or `onError`, they are
  effectively silent.
- Re-run `npm run build` after changing `src/`; the examples load `dist/`.

## Getting help

- Check the component API page and [Configuration](configuration.md).
- Reproduce with the [examples](https://github.com/artsyne/uno-ui/blob/main/examples/index.html) and the committed `dist/`.
- Open an issue with version, Alpine/HTMX/Tailwind versions, browser, minimal HTML and console
  output. Use the issue template.
- Security issues go through [SECURITY.md](https://github.com/artsyne/uno-ui/blob/main/SECURITY.md), not a public issue.

## Related

- [Getting Started](getting-started.md)
- [Error handling](error-handling.md)
- [HTMX integration](htmx.md)
- [CSP](csp.md)
- [i18n](i18n.md)
