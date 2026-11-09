# Configuration

Uno resolves configuration at four levels and merges them deeply. This page documents the
levels, the merge rules, the global options, the registered plugin defaults and the complete
`ui` style map.

## Four levels

From lowest to highest priority:

| Level | Source | Example |
| ---: | :--- | :--- |
| 1 | Plugin defaults registered by the component | `Uno.registerPluginDefaults('notify', { duration: 3000 })` |
| 2 | `plugins.<component>` in `Uno.setConfig()` | `Uno.setConfig({ plugins: { notify: { duration: 5000 } } })` |
| 3 | Top-level `<component>` in `Uno.setConfig()` | `Uno.setConfig({ notify: { duration: 5000 } })` |
| 4 | Instance options | `notify.show('Saved', { duration: 1000 })` |

Levels 2 and 3 are both global. Use `plugins` as the primary namespace; the top-level key is
supported by `Uno.getConfig()` and by the plugin factories, but it is not represented in the
`UnoOptions` TypeScript interface yet.

Resolution happens when an Action runs, not when the script loads:

```javascript
const config = Uno.getConfig('notify', { duration: 1000 })
```

## Setting configuration

`Uno.setConfig(config)` accepts a plain object, merges it deeply, and returns the Uno instance
so calls can be chained.

### Before initialization (recommended)

```javascript
// Runs before Alpine starts.
Uno.setConfig({
  debug: true,
  csp: { nonce: window.__CSP_NONCE__ },
  ui: { base: 'rounded-lg shadow-md' },
  plugins: {
    box: { defaultWidth: 480, closeOnEsc: true },
    notify: { duration: 4000 }
  }
})
```

Configuration set before initialization is stored separately and merged during `Uno.init()`.

### After initialization

```javascript
// Alpine is running.
Uno.setConfig({
  ui: { box: { okBtn: '+min-w-24' } }
})
```

Styles and component options are read at call time, so most UI changes take effect on the next
Action. Two options are applied only during initialization:

- `debug` calls `log.setDebug()` once. Changing it later does not change the logger.
- `onError` calls `log.setReporter()` once. Changing it later does not replace the reporter.

Call `setConfig()` before Alpine starts when you need either one.

### Reading configuration

```javascript
Uno.getGlobalConfig() // pre-init config before init, resolved options after
Uno.options // resolved options object
Uno.getConfig('box', { title: 'Example' }) // one component, merged
Uno.getStyle('box', 'okBtn') // resolved style string
```

`getStyle(component, path, fallback)` accepts dot paths (`types.success`,
`arrowSides.top`), substitutes `{base}`, and handles a leading `+`.

## Merge semantics

- Plain objects merge recursively.
- Arrays and primitives replace the previous value, including an explicit `undefined`.
- A string beginning with `+` appends to an existing string value:

```javascript
Uno.setConfig({
  ui: {
    box: {
      // Appends to the default okBtn class string.
      okBtn: '+min-w-24'
    }
  }
})
```

`+` appending is applied both by `mergeDeep()` and by `getStyle()`, depending on whether the
key already has a value. The result is a single space-separated class string.

- `{base}` is substituted only when a style is resolved through `getStyle()`. The default
  `ui.base` is:

```text
rounded-lg shadow-md border border-gray-200 dark:border-gray-700
```

If you override `{base}` in a custom style string, keep the token in the component default or
write the full class list yourself.

## Global options

| Name | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `debug` | `boolean` | `false` | Enables `[Uno Info]`, `[Uno Debug]`, `[Uno Warn]`, `[Uno Error]` console output. Applied at initialization. |
| `csp.nonce` | `string \| null` | `null` | Nonce written to Uno-created containers and to `<style>` elements created by `Uno.dom.injectStyles()`. |
| `onError` | `function \| null` | `null` | Receives the arguments passed to `log.error(...)`. It is not always an `Error` instance. Applied at initialization. |
| `plugins` | `object` | `{}` | Behavior options per component. |
| `ui` | `object` | see below | Style configuration consumed by `Uno.getStyle()`. |

The top-level component keys (`box`, `notify`, `tooltip`, `tab`, `valid`) are also accepted by
`getConfig()`. They are not declared in `UnoOptions` yet; see
[TypeScript usage](typescript-usage.md#known-type-gaps).

## Plugin defaults

These are the built-in values registered by each plugin and then merged with user configuration.

### Box

| Name | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `defaultWidth` | `number` | `400` | Width in pixels when neither `width` nor `widthClass` is set. |
| `defaultPosition` | `BoxPosition` | `'center'` | Position when the instance does not set `position`. |
| `okText` | `string` | `'OK'` | Default OK label. Per-instance `okText` wins. |
| `cancelText` | `string` | `'Cancel'` | Default Cancel label. Per-instance `cancelText` wins. |
| `closeOnEsc` | `boolean` | `false` | Default for `closeOnEsc`. |
| `closeOnOverlayClick` | `boolean` | `false` | Default for `closeOnOverlayClick`. |

### Notify

| Name | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `duration` | `number` | `3000` | Auto-dismiss delay in milliseconds. `0` disables it. |
| `position` | `NotifyPosition` | `'top-right'` | Default container position. |

`type` has no registered default. If you call `notify.show('Saved')` without a `type`, the
notification renders without a type-specific background or icon. Pass `type: 'info'` (or one
of the other three values) when you want the default visual treatment.

### Tooltip

| Name | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `placement` | `'top' \| 'bottom' \| 'left' \| 'right'` | `'top'` | Preferred side; the position engine may flip. |
| `trigger` | `'hover' \| 'click' \| 'focus'` | `'hover'` | Interaction mode. |
| `offset` | `number` | `8` | Gap between trigger and popup, in pixels. |
| `arrow` | `boolean` | `true` | Render the arrow square. |
| `clickOutside` | `boolean` | `true` | Close a click popup on an outside click. |
| `delay` | `number` | `300` | Registered default; the directive currently calls `scheduleHide()` with its own 300 ms default. Treat this key as not applied. |

### Tab

| Name | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `hash` | `boolean` | `true` | Enables hash updates for `activate()`. An instance-level `hash` wins. |
| `defaultTab` | `string` | `''` | Registered default; not read by the current store. |
| `activeClass` | `string` | `'is-active'` | Registered default; not read by the current store. |
| `inactiveClass` | `string` | `''` | Registered default; not read by the current store. |

### Valid

| Name | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `realtime` | `boolean` | `true` | Registered default; the global value is not read. Per-rule `realtime: false` still skips the input pass. |
| `showError` | `boolean` | `true` | Registered default; not read by the current store. |
| `scrollOnError` | `boolean` | `true` | When `false`, `validateForm()` does not scroll or focus the first error/summary. |

Additional Valid configuration is read through `Uno.getConfig('valid')`: `emptyQuit`
(default `true`), `autoError` (default `true`), `errorClass` and `createErrorElement`.
Per-rule options such as `emptyQuit`, `debounce` and `trigger` are documented in the
[Valid API](valid-api.md).

## UI token reference

### Shared

| Name | Default | Notes |
| :--- | :--- | :--- |
| `base` | `rounded-lg shadow-md border border-gray-200 dark:border-gray-700` | Substituted wherever a component class string contains `{base}`. |
| `layers.box` | `1400` | Box manager z-index. |
| `layers.notify` | `1600` | Notify container z-index. |
| `layers.tooltip` | `1700` | Tooltip container z-index. |

### Box

| Token | Default |
| :--- | :--- |
| `overlay` | `fixed inset-0 bg-gray-500/75 transition-opacity` |
| `dialogWrapper` | `fixed inset-0 z-10 w-screen overflow-y-auto` |
| `itemWrapper` | `flex min-h-full p-4 text-center sm:p-0 pointer-events-none` |
| `container` | `relative transform overflow-hidden {base} bg-white dark:bg-gray-900 text-start shadow-xl transition-all sm:my-8 sm:w-full sm:max-w-lg flex flex-col max-h-[calc(100dvh-2rem)] sm:max-h-[calc(100dvh-4rem)] pointer-events-auto outline-none` |
| `containerFill` | `h-[calc(100dvh-2rem)] sm:h-[min(85dvh,56rem)]` |
| `cancelBtn` | `mt-3 inline-flex w-full justify-center rounded-md bg-white px-3 py-2 text-sm font-semibold text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 hover:bg-gray-50 sm:mt-0 sm:w-auto` |
| `extraBtn` | `mt-3 inline-flex w-full justify-center rounded-md bg-white px-3 py-2 text-sm font-normal text-gray-700 hover:bg-gray-50 sm:mt-0 sm:w-auto` |
| `okBtn` | `inline-flex w-full justify-center rounded-md bg-blue-600 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-blue-500 sm:ms-3 sm:w-auto` |
| `title` | `text-base font-semibold leading-6 text-gray-900 dark:text-slate-200` |
| `closeBtnFloating` | `absolute top-2 end-2 z-10 p-1 -m-1 text-slate-500 hover:text-slate-600 dark:text-slate-400 dark:hover:text-slate-300 focus:outline-none` |
| `closeIconSmall` | `h-5 w-5` |
| `content` | `px-4 pb-4 pt-5 sm:p-6 sm:pb-4 overflow-auto flex-1 text-slate-600 dark:text-slate-300` |
| `contentFill` | `flex flex-col min-h-0` |
| `bodyFill` | `flex-1 min-h-0 flex flex-col` |
| `loadingWrap` | `flex justify-center items-center h-full` |
| `spinner` | `animate-spin h-8 w-8 text-blue-500` |
| `footer` | `bg-gray-50 px-4 py-3 sm:flex sm:flex-row-reverse sm:px-6 border-t border-slate-200 dark:border-slate-700 dark:bg-slate-800` |
| `promptWrapper` | `space-y-3` |
| `promptText` | `text-slate-700 dark:text-slate-200` |
| `promptInput` | `w-full px-3 py-2 border rounded focus:ring-2 focus:ring-blue-500 outline-none` |

### Notify

| Token | Default |
| :--- | :--- |
| `container` | `fixed flex gap-2 p-4 pointer-events-none` |
| `item` | `pointer-events-auto max-w-[calc(100vw-2rem)] overflow-hidden {base}` |
| `content` | `p-4` |
| `row` | `flex items-start` |
| `iconWrap` | `flex-shrink-0` |
| `textWrap` | `ms-3 min-w-0 flex-1 pt-0.5` |
| `text` | `text-sm font-medium text-gray-900 dark:text-slate-100 break-words` |
| `closeWrap` | `ms-4 flex flex-shrink-0` |
| `closeBtn` | `inline-flex rounded-md p-1 -m-1 text-gray-500 hover:text-gray-700 hover:bg-black/5 dark:text-slate-400 dark:hover:text-slate-200 dark:hover:bg-white/10 focus:outline-none focus:ring-2 focus:ring-indigo-500` |
| `closeIcon` | `h-5 w-5` |
| `progressWrap` | `fixed top-0 start-0 end-0 h-1 bg-gray-200 dark:bg-gray-700 pointer-events-none` |
| `progressBar` | `h-full bg-gradient-to-r from-blue-500 to-indigo-600 shadow-sm shadow-blue-500/50 rtl:bg-gradient-to-l` |
| `types.success` | `bg-green-100 dark:bg-green-950` |
| `types.error` | `bg-red-100 dark:bg-red-950` |
| `types.warning` | `bg-yellow-100 dark:bg-yellow-950` |
| `types.info` | `bg-blue-100 dark:bg-blue-950` |
| `icons.success` | `text-green-500 dark:text-green-300` |
| `icons.error` | `text-red-500 dark:text-red-300` |
| `icons.warning` | `text-yellow-500 dark:text-yellow-300` |
| `icons.info` | `text-blue-500 dark:text-blue-300` |

### Tooltip

| Token | Default |
| :--- | :--- |
| `content` | `fixed px-3 py-2 text-sm font-medium {base} whitespace-nowrap bg-white text-gray-900 dark:bg-slate-800 dark:text-white` |
| `body` | `pointer-events-auto` |
| `arrow` | `absolute w-2 h-2 rotate-45 bg-inherit border-gray-200 dark:border-gray-700` |
| `arrowSides.top` | `border-b border-e` |
| `arrowSides.bottom` | `border-t border-s` |
| `arrowSides.left` | `border-t border-e` |
| `arrowSides.right` | `border-b border-s` |

### Tab

| Token | Default |
| :--- | :--- |
| `active` | `border-b-2 border-blue-500` |
| `inactive` | `text-slate-500 hover:text-slate-700` |
| `loading` | `animate-pulse h-32 w-full rounded-lg bg-gray-200 dark:bg-slate-700` |
| `error` | `p-4 text-red-600 dark:text-red-400` |

### Valid

| Token | Default |
| :--- | :--- |
| `wrapper` | `relative w-full` |
| `iconBase` | `absolute end-0 top-0 bottom-0 flex items-center justify-center w-10 pointer-events-none transition-opacity duration-200` |
| `inputError` | `!border-red-500` |
| `inputValid` | `!border-green-500` |
| `iconError` | `text-red-500 dark:text-red-400` |
| `iconValid` | `text-green-500 dark:text-green-400` |
| `stateInvalid` | `is-invalid` |
| `stateValid` | `is-valid` |
| `stateValidating` | `validating` |
| `inputPaddingIcon` | `pe-10` |
| `helperText` | `text-red-600 dark:text-red-400 text-xs mt-1` |

The `wrapper`, `iconBase`, `iconError`, `iconValid`, `state*` and `inputPaddingIcon` tokens are
reserved for custom templates/plugins; the current Valid directive only reads `helperText` and
uses its own counter classes. They remain part of the public `ui.valid` object.

## Related

- [Theming and Tailwind](theming.md)
- [i18n](i18n.md)
- [Box API](box-api.md), [Notify API](notify-api.md), [Tooltip API](tooltip-api.md),
  [Tab API](tab-api.md), [Valid API](valid-api.md)
