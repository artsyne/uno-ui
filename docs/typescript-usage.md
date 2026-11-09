# TypeScript Usage

Type definitions ship with the package (`types/index.d.ts`) and are wired through `package.json`:

```json
{
  "types": "./types/index.d.ts",
  "exports": {
    ".": {
      "types": "./types/index.d.ts",
      "import": "./dist/uno.es.js",
      "require": "./dist/uno.umd.js"
    }
  }
}
```

No separate `@types/uno-ui` package exists.

## tsconfig

Use a bundler-oriented resolution so the `exports.types` condition is honored:

```json
{
  "compilerOptions": {
    "target": "ES2020",
    "module": "ESNext",
    "moduleResolution": "bundler",
    "strict": true,
    "skipLibCheck": true,
    "types": []
  }
}
```

The repository's [tsconfig.example.json](../tsconfig.example.json) is for building Uno itself,
not a consumer template. Copy the relevant compiler options, not the `rootDir`/`outDir` values.

## Imports

```typescript
import Uno, { DomManager, box, notify, tooltip, tab, valid } from '@artsyne/uno-ui'
import type {
  UnoOptions,
  BoxOptions,
  BoxItem,
  NotifyOptions,
  NotifyType,
  TooltipOptions,
  TooltipPlacement,
  TabOptions,
  UiConfig,
  BoxUiConfig
} from '@artsyne/uno-ui'
```

`Uno` is both the default and a named export. `DomManager` and the plugin factories are named
exports. Import the package for its side effects even if you only use directives.

## Configuration

```typescript
import Alpine from 'alpinejs'
import Uno from '@artsyne/uno-ui'

const config: UnoOptions = {
  debug: false,
  csp: { nonce: window.__CSP_NONCE__ },
  plugins: {
    box: { defaultWidth: 480, closeOnEsc: true },
    notify: { duration: 4000 }
  },
  ui: {
    base: 'rounded-lg shadow-md',
    layers: { box: 1400, notify: 1600, tooltip: 1700 }
  }
}

Uno.setConfig(config)

window.Alpine = Alpine
Alpine.start()
```

`UnoOptions` currently declares `plugins`, `ui`, `debug`, `csp` and `onError`. The runtime also
accepts top-level component keys (`box`, `notify`, ...); see
[Known type gaps](#known-type-gaps).

## Box

```typescript
import type { BoxItem, BoxOptions } from '@artsyne/uno-ui'

const options: BoxOptions = {
  type: 'confirm',
  title: 'Delete item?',
  content: 'This cannot be undone.',
  position: 'center',
  showOkButton: true,
  showCancelButton: true,
  closeOnEsc: true,
  onOk: () => {
    return true
  }
}

const item: BoxItem = Alpine.store('uno').box.open(options)
item.close()

const confirmed: boolean = await Alpine.store('uno').box.confirm('Continue?', {
  title: 'Confirm'
})
```

`BoxOptions.onOk` is typed as `() => void | false`; async callbacks returning a Promise are
accepted at runtime but the declaration is narrower. `BoxItem` is declared as a required form of
`BoxOptions` plus `id`, `isOpen`, `extraButtons` and `close()`; the runtime object contains the
resolved fields, so the two are not a perfect structural match today.

## Notify

```typescript
import type { NotifyOptions, NotifyType } from '@artsyne/uno-ui'

const options: NotifyOptions = {
  type: 'success',
  position: 'top-right',
  duration: 3000,
  closable: true
}

const id: number = Alpine.store('uno').notify.show('Saved', options)
Alpine.store('uno').notify.remove(id)

const types: NotifyType[] = ['success', 'error', 'warning', 'info']
```

`type` is optional in the type, but there is no runtime default. Pass a value when you expect the
typed visual treatment.

## Tooltip

```typescript
import type { TooltipOptions, TooltipPlacement } from '@artsyne/uno-ui'

const placement: TooltipPlacement = 'top'

const options: TooltipOptions = {
  content: 'Help text',
  placement,
  trigger: 'hover',
  arrow: true,
  offset: 8
}

Alpine.store('uno').tooltip.show({
  ...options,
  element: document.querySelector('#help') as HTMLElement
})
```

`offset` is typed as `number`; the runtime also accepts a function. The declaration has not been
widened yet.

## Tab

```typescript
import type { TabOptions } from '@artsyne/uno-ui'

const options: TabOptions = {
  name: 'profile',
  group: 'account',
  hash: true
}

Alpine.store('uno').tab.activate(options)
const active: boolean = Alpine.store('uno').tab.isActive('profile', 'account')
```

`activate()` ignores positional arguments; always pass an object.

## Store typing

The package augments Alpine's `Stores` interface and `Window`:

```typescript
Alpine.store('uno').notify.success('Saved')
Alpine.store('i18n').setLang('zh-CN')
window.Uno.setConfig({ debug: true })
window.Alpine = Alpine
```

If `Alpine.store('uno')` is `any` in your project, check that:

1. The package is imported (even for side effects).
2. `moduleResolution` is `bundler`, `node16` or `nodenext`.
3. `skipLibCheck` is not hiding a conflict between your Alpine declarations and Uno's.

Alpine 3.17 does not publish a `types` field in its package metadata. If your project
type-checks `import Alpine from 'alpinejs'`, you may need `@types/alpinejs` or a local module
declaration. Uno's `types/alpinejs.d.ts` is a minimal internal declaration, not a replacement
for Alpine's own types.

## Extending types

```typescript
import type { BoxOptions } from '@artsyne/uno-ui'

function showDialog<T extends BoxOptions>(options: T): Promise<boolean> {
  return Alpine.store('uno').box.confirm('Proceed?', options)
}
```

To add your own store, use Alpine's own augmentation pattern. Uno does not provide a generic
store registry type.

## Known type gaps

These are current limitations in `types/index.d.ts`; do not rely on them without a local cast or
declaration:

| Gap | Detail |
| :--- | :--- |
| Top-level component config | `UnoOptions` lacks `box`, `notify`, `tooltip`, `tab`, `valid`. Use `plugins.*` for typed config. |
| `BoxItem` | Declared as `Required<Omit<BoxOptions, 'extraButtons'>>`; the runtime item is not fully required and contains extra internals. |
| `DomManager.escape` | Declared in the interface, not implemented in `src/utils/dom.js`. |
| Store `init()` | Declared on Box/Notify/Tooltip/Tab stores, not implemented. Do not call it. |
| `ValidStore` methods | `showError`, `hideError`, `clearErrors` are implemented but not declared. |
| `loadLocale` | Typed as `Promise<boolean>`; the implementation resolves `undefined`. |
| `init()` return | Typed as `AlpineType`; it can resolve `null` when Alpine is unavailable. |
| `offset` | `TooltipOptions.offset` is `number`; the runtime accepts a function. |
| `EmptyOptions` | `UnoOptions` is not exported as a runtime value; it is a type only. |

The repository's `tests/types/types.test.ts` is not collected by `npm run test:run` because the
Vitest `include` pattern is `tests/**/*.test.js`. Type changes are therefore not covered by the
CI test step; see [Test reports](test-reports.md).

## Common errors

| Error | Fix |
| :--- | :--- |
| `Cannot find module '@artsyne/uno-ui'` | Install the package or use a file dependency until the first publish; check `moduleResolution`. |
| `Property 'box' does not exist on type ...` | Import the package for side effects and use `Alpine.store('uno')`; ensure the store augmentation is visible. |
| `Cannot find module 'alpinejs'` | Install `alpinejs`; add `@types/alpinejs` or a local declaration if your compiler needs it. |
| `Type 'string' is not assignable to type 'BoxPosition'` | Use one of the nine documented positions. |
| `onOk` async error | The runtime accepts Promises; the type is narrower. Cast or use a typed wrapper. |

## Related

- [Configuration](configuration.md)
- [TypeScript tests](test-reports.md)
- [Box API](box-api.md)
- [Notify API](notify-api.md)
- [Tooltip API](tooltip-api.md)
- [Tab API](tab-api.md)
