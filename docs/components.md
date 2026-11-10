# Components

Uno provides five components. They share the same Store-First API, lazy container
initialization and Tailwind-based styling.

| Component | Use it for | API |
| :--- | :--- | :--- |
| **Box** | Modals, alerts, confirms, prompts, HTMX-loaded dialogs and fill layouts. | [Box API](box-api.md) |
| **Notify** | Toast notifications in nine positions, a global progress bar and HTMX content. | [Notify API](notify-api.md) |
| **Tooltip** | Hover, click and focus popups with auto-flipping, viewport clamping and HTMX content. | [Tooltip API](tooltip-api.md) |
| **Tab** | Tab groups with hash deep links, HTMX panels and APG keyboard behaviour. | [Tab API](tab-api.md) |
| **Valid** | Array-based form validation, async validators, counters and error summaries. | [Valid API](valid-api.md) |

## Shared setup

All components require Alpine.js 3 and Uno to be registered before Alpine starts. Tailwind CSS 4
must scan Uno's templates; the package does not ship a stylesheet.

```bash
npm install @artsyne/uno-ui alpinejs
```

```javascript
import Alpine from 'alpinejs'
import Uno from '@artsyne/uno-ui'

Uno.setConfig({
  debug: false,
  ui: { base: 'rounded-lg shadow-md' }
})

window.Alpine = Alpine
Alpine.start()
```

```css
@import 'tailwindcss';
@source '../node_modules/@artsyne/uno-ui/src';
```

See [Getting Started](getting-started.md) for the CDN path, configuration timing and the
Alpine loading paths.

## Store and directives

Every component is available from the `uno` store and, where it makes sense, through a
directive:

| Component | Store | Directive |
| :--- | :--- | :--- |
| Box | `Alpine.store('uno').box` | `x-uno-box` |
| Notify | `Alpine.store('uno').notify` | `x-uno-notify` |
| Tooltip | `Alpine.store('uno').tooltip` | `x-uno-tooltip` |
| Tab | `Alpine.store('uno').tab` | `x-uno-tab` |
| Valid | `Alpine.store('uno').valid` | `x-uno-valid` |

See [Core Concepts](concepts.md) for the architecture and [Configuration](configuration.md) for
the four-level merge, `ui` tokens and the `{base}` / `+` conventions.

## Accessibility

Each component has documented keyboard behaviour and component-level axe scans in CI. The
project does not claim WCAG conformance; see [Accessibility](accessibility.md) for what the
scans cover and what still needs manual acceptance.
