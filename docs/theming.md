# Theming and Tailwind

Uno is a class-string library, not a CSS package. It ships no stylesheet; every component
renders Tailwind utility classes that come from JavaScript template strings and the default
`ui` configuration. Your Tailwind build must see those strings, and you theme Uno by replacing
or appending classes.

## Tailwind 4 setup

### Application build

Add an `@source` rule for the package's `src` directory. The published package includes
`src/` in its `files` list.

```css
/* app.css */
@import 'tailwindcss';

@source '../node_modules/@artsyne/uno-ui/src';
@source '../src';
```

If you vendor the built bundle instead of installing the package, point at the directory that
contains `uno.es.js` / `uno.umd.js`:

```css
@source '../vendor/uno-ui/dist';
```

Both directories contain the class strings: `src/` has the readable templates, `dist/` has the
compiled copies. Use one, not both, unless your application also ships both.

### Browser build

```html
<script src="https://cdn.jsdelivr.net/npm/@tailwindcss/browser@4"></script>
```

`@tailwindcss/browser@4` scans the DOM at runtime. It sees Uno's templates after they are
inserted, which makes the demo pages work without a build step. This mode is intended for
prototypes; use the build integration for production.

### Classes outside the package

Class names in your own `ui` overrides are ordinary strings in your JavaScript or HTML. Keep
them as complete literals so Tailwind can detect them. Avoid constructing class names at
runtime (`'bg-' + color`) unless you also ensure the full names appear somewhere Tailwind
scans.

## Dark mode

Uno's defaults contain both light and `dark:` classes. Tailwind 4's default `dark:` variant
follows `prefers-color-scheme`. To use a class-based theme instead, declare the variant:

```css
@import 'tailwindcss';
@custom-variant dark (&:where(.dark, .dark *));
```

Then toggle `dark` on `<html>` from your own code:

```javascript
document.documentElement.classList.toggle('dark', isDark)
```

Uno does not manage the page theme and does not set `<html class>` or `<html lang>`.

## Changing the base style

`ui.base` is substituted wherever a component default contains `{base}`. The default is:

```text
rounded-lg shadow-md border border-gray-200 dark:border-gray-700
```

### Square, flat theme

```javascript
Uno.setConfig({
  ui: {
    base: 'rounded-none shadow-none border border-gray-200 dark:border-gray-700'
  }
})
```

### Brand colour

```javascript
Uno.setConfig({
  ui: {
    box: {
      okBtn:
        'inline-flex w-full justify-center rounded-md bg-indigo-600 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 sm:ms-3 sm:w-auto'
    },
    tab: {
      active: 'border-b-2 border-indigo-600 text-indigo-600',
      inactive: 'text-slate-500 hover:text-indigo-600'
    }
  }
})
```

### Append instead of replace

Prefix a class string with `+` to keep the default and add to it:

```javascript
Uno.setConfig({
  ui: {
    box: {
      okBtn: '+min-w-24',
      container: '+max-w-3xl'
    },
    notify: {
      item: '+backdrop-blur'
    }
  }
})
```

The merge produces a single space-separated string. `+` is handled both while merging
configuration and while resolving styles, so it works for keys that have a default and for keys
that do not.

## Per-component styling

Most behavior options accept a class override for a specific instance:

```javascript
Alpine.store('uno').box.open({
  title: 'Delete item',
  okText: 'Delete',
  okButtonClass: '!bg-red-600 hover:!bg-red-700',
  cancelButtonClass: '!text-slate-600',
  containerClass: 'sm:max-w-xl'
})
```

Use the instance-level class options for one-off changes and `ui` for global changes.
Instance-level classes are appended after the global class string.

## The `ui` map

The complete default map is documented in
[Configuration § UI token reference](configuration.md#ui-token-reference). The main groups are:

| Group | Keys |
| :--- | :--- |
| Shared | `base`, `layers.box`, `layers.notify`, `layers.tooltip` |
| `box` | overlay, dialog wrapper, container, content, footer, buttons, title, close button, prompt, fill mode |
| `notify` | container, item, content, text, close button, progress bar, type backgrounds, icon colours |
| `tooltip` | content, body, arrow, arrow sides |
| `tab` | active, inactive, loading, error |
| `valid` | helper text and reserved validator-state tokens |

Changing a key to a plain string replaces the default. Changing it to `+string` appends.

## Z-index layers

```javascript
Uno.setConfig({
  ui: {
    layers: { box: 2000, notify: 2100, tooltip: 2200 }
  }
})
```

Box dialogs stack by adding `index * 10` to the box layer. The global progress bar uses
`9999`. See [Layers](concepts.md#layers) for stacking-context caveats.

## Valid error text

The generated error element uses `ui.valid.helperText` by default:

```text
text-red-600 dark:text-red-400 text-xs mt-1
```

To use a different class globally, set `valid.errorClass`:

```javascript
Uno.setConfig({
  valid: {
    errorClass: 'text-danger small mt-2'
  }
})
```

`valid.autoError` (default `true`) controls whether the directive creates the error element.
See [Valid API](valid-api.md#auto-error-generation).

## What not to target

The injected templates use internal classes and IDs (`uno-box-container{id}`,
`uno-tooltip-content`, etc.). They are implementation details. Prefer `ui` tokens for styling
and the documented callbacks for behavior. Targeting internal nodes can break between minor
versions.

## Related

- [Configuration](configuration.md)
- [Getting Started](getting-started.md)
- [Box API](box-api.md)
- [Notify API](notify-api.md)
