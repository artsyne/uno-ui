# Uno Box API Documentation

The Box component is a lightweight, declarative, HTMX-native modal system based on Alpine.js and TailwindCSS. It implements a "Store-First, Directive-Driven" architecture, offering interfaces for alerts, confirmations, prompts, and generic dialogs.

## Overview

- **Core Function**: Provides modal dialog management with support for stacking, positioning, and content injection.
- **Architecture**: Store-based state management (`Alpine.store('uno').box`) with directive triggers (`x-uno-box`).
- **Integration**: Native support for HTMX content loading (`hxUrl`).
- **A11y**: Built-in focus trapping and keyboard navigation (ESC to close). The floating close button's accessible name is localized via the `box.close` i18n key. On close, focus restores to the opener element (WCAG 2.4.3); stacked dialogs mark underlay layers `inert` + `aria-hidden` so AT never penetrates below the top dialog.

## Basic Usage

### 1. Directive Usage

Trigger dialogs directly from HTML without writing JavaScript.

```html
<!-- Simple Alert -->
<button x-uno-box='{ "title": "Hello", "content": "Welcome to Uno!" }'>Open Alert</button>

<!-- Confirmation Dialog -->
<button
  x-uno-box='{ 
  "type": "confirm",
  "title": "Delete Item?",
  "content": "This action cannot be undone.",
  "onOk": "deleteItem()",
  "okText": "Yes",
  "cancelText": "No",
  "okButtonClass": "bg-red-600 hover:bg-red-700",
  "cancelButtonClass": "ml-2"
}'
>
  Delete
</button>
```

### 2. JavaScript Usage

Use the global store for programmatic control.

```javascript
// Alert
Alpine.store('uno').box.alert('Operation successful!')

// Confirm
Alpine.store('uno')
  .box.confirm('Are you sure?', {
    title: 'Confirm Action',
    okText: 'Yes, do it',
    okButtonClass: 'bg-green-600',
    cancelButtonClass: 'text-gray-500'
  })
  .then((confirmed) => {
    if (confirmed) console.log('Confirmed')
  })

// Prompt
Alpine.store('uno')
  .box.prompt('Enter your name:', 'Guest')
  .then((name) => {
    if (name) console.log('Hello, ' + name)
  })
```

## API Parameters

The following options can be passed to `open()`, `dialog()`, or used in the `x-uno-box` configuration object.

| Name                  | Type       | Default    | Required | Description                                                                                                  |
| :-------------------- | :--------- | :--------- | :------- | :----------------------------------------------------------------------------------------------------------- |
| `type`                | `String`   | `'dialog'` | No       | Dialog type. Values: `'dialog'`, `'confirm'`.                                                                |
| `title`               | `String`   | `''`       | No       | The title text displayed in the header.                                                                      |
| `content`             | `String`   | `''`       | No       | HTML content for the body. Auto-sanitized.                                                                   |
| `width`               | `Number`   | `400`      | No       | Width in pixels. Ignored if `widthClass` is set.                                                             |
| `widthClass`          | `String`   | -          | No       | Tailwind CSS classes for width (e.g., `'w-full md:w-lg'`).                                                   |
| `height`              | `Number`   | -          | No       | Height in pixels. Ignored if `heightClass` is set.                                                           |
| `heightClass`         | `String`   | -          | No       | Tailwind CSS classes for height (e.g., `'h-full'`).                                                          |
| `position`            | `String`   | `'center'` | No       | Screen position. Values: `'center'`, `'top-center'`, `'bottom-right'`, etc.                                  |
| `modal`               | `Boolean`  | `true`     | No       | If `true`, shows a backdrop and blocks interaction with the page.                                            |
| `showCloseButton`     | `Boolean`  | `true`     | No       | Whether to show the "X" close button in the header.                                                          |
| `showOkButton`        | `Boolean`  | `false`    | No       | Whether to show the primary action button. (Default `true` for `confirm` type).                              |
| `showCancelButton`    | `Boolean`  | `false`    | No       | Whether to show the secondary action button. (Default `true` for `confirm` type).                            |
| `okText`              | `String`   | `'OK'`     | No       | Text for the OK button. Supports i18n keys.                                                                  |
| `cancelText`          | `String`   | `'Cancel'` | No       | Text for the Cancel button. Supports i18n keys.                                                              |
| `okButtonClass`       | `String`   | -          | No       | Tailwind CSS classes for the OK button.                                                                      |
| `cancelButtonClass`   | `String`   | -          | No       | Tailwind CSS classes for the Cancel button.                                                                  |
| `extraButtons`        | `Array`    | `[]`       | No       | Custom buttons rendered alongside OK/Cancel. See [ExtraButton](#extrabutton-options) for each item's schema. |
| `closeOnEsc`          | `Boolean`  | `false`    | No       | If `true`, pressing ESC closes the dialog.                                                                   |
| `closeOnOverlayClick` | `Boolean`  | `false`    | No       | If `true`, clicking the backdrop closes the dialog.                                                          |
| `autoFocus`           | `Boolean`  | `true`     | No       | If `true`, automatically focuses the first input field after opening.                                        |
| `hxUrl`               | `String`   | -          | No       | URL to fetch content from via HTMX.                                                                          |
| `hxMethod`            | `String`   | `'GET'`    | No       | HTTP method for HTMX request.                                                                                |
| `onOk`                | `Function` | -          | No       | Callback when OK button is clicked. Return `false` to prevent closing.                                       |
| `onCancel`            | `Function` | -          | No       | Callback when Cancel button is clicked.                                                                      |
| `onClose`             | `Function` | -          | No       | Callback when dialog is closed (after animation).                                                            |

## ExtraButton Options

Each item in the `extraButtons` array supports the following fields.

| Name          | Type       | Default    | Required | Description                                                                                                                                                        |
| :------------ | :--------- | :--------- | :------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `text`        | `String`   | -          | **Yes**  | Button label text.                                                                                                                                                 |
| `position`    | `String`   | `'before'` | No       | Placement relative to OK/Cancel. `'before'` = left of OK button; `'after'` = right of Cancel.                                                                      |
| `buttonClass` | `String`   | -          | No       | Additional CSS classes to override the default `extraBtn` style.                                                                                                   |
| `onClick`     | `Function` | -          | No       | Click handler. Receives the `BoxItem` instance. **Return `false` to keep the dialog open**; return `void`/`true` or a `Promise` resolving to non-`false` to close. |

### Behavior Parity with `onOk`

`extraButtons[n].onClick` follows the same close-control contract as `onOk`:

| Return value         | Dialog closes? |
| :------------------- | :------------- |
| `void` / `undefined` | ✅ Yes         |
| `true`               | ✅ Yes         |
| `false`              | ❌ No          |
| `Promise<void>`      | ✅ Yes         |
| `Promise<false>`     | ❌ No          |

## Auto Focus

When a box opens (and again after HTMX content loads), initial focus follows the WAI-ARIA APG dialog pattern:

1. **Inputable element** (`autoFocus: true` only): the first visible `input` (text, password, etc.), `textarea`, or `select`.
   - **Exclusions**: `hidden`, `disabled`, `readonly` elements, and buttons (submit/reset/button).
   - **Visibility Check**: the element and every ancestor must be visible (not `display: none` or `visibility: hidden`).
2. **Neutral target**: the title heading when the dialog has a title, otherwise the dialog container itself (both carry `tabindex="-1"` and suppress the focus outline). This step also runs when `autoFocus: false`.

Initial focus never lands on chrome controls such as the floating close button: script focus on them makes Safari draw a `:focus-visible` ring and lets Enter/Space dismiss the dialog by accident. The focus trap (`x-uno-trap`) cycles Tab/Shift+Tab within the dialog without taking initial focus. Wrapping is DOM-order aware: when focus sits on a neutral target (`tabindex="-1"`, e.g. the title) with no focusable element after/before it, the trap wraps to the first/last focusable instead of letting the browser's native tab order escape the dialog.

Set `autoFocus: false` to skip step 1 (e.g., read-only content where focusing an input would pop up the keyboard); the neutral target is still focused so focus stays inside the modal.

```javascript
Alpine.store('uno').box.open({
  title: 'Read Only',
  content: '<input value="Copy this" readonly>',
  autoFocus: false
})
```

## Events

The Box component triggers the following callbacks.

| Event / Callback | Trigger Timing                                              | Callback Parameters               |
| :--------------- | :---------------------------------------------------------- | :-------------------------------- |
| `onOk`           | Click on "OK" button.                                       | None                              |
| `onCancel`       | Click on "Cancel" button.                                   | None                              |
| `onBeforeClose`  | Immediately before the close animation starts.              | None                              |
| `onClose`        | After the close animation completes and element is removed. | None                              |
| `onContentLoad`  | After HTMX content has been loaded and injected.            | `el` (The loaded content element) |

## Methods

Accessed via `Alpine.store('uno').box`.

### `open(options)`

Opens a generic dialog.

- **Parameters**: `options` (Object) - Configuration object.
- **Returns**: `Object` - The dialog instance (reactive).
- **Example**:
  ```javascript
  Alpine.store('uno').box.open({ title: 'Info', content: 'Text' })
  ```

### `alert(message, options)`

Shows an alert dialog with an OK button.

- **Parameters**:
  - `message` (String) - The message text.
  - `options` (Object, optional) - Additional config.
- **Returns**: `Object` - The dialog instance.
- **Example**:
  ```javascript
  Alpine.store('uno').box.alert('Saved!')
  ```

### `confirm(message, options)`

Shows a confirmation dialog with OK/Cancel buttons.

- **Parameters**:
  - `message` (String) - The question text.
  - `options` (Object, optional) - Additional config.
- **Returns**: `Promise<Boolean>` - Resolves to `true` (OK) or `false` (Cancel/Close).
- **Example**:
  ```javascript
  const ok = await Alpine.store('uno').box.confirm('Delete?')
  ```

### `prompt(message, defaultValue, options)`

Shows a dialog with a text input.

- **Parameters**:
  - `message` (String) - The label text.
  - `defaultValue` (String, optional) - Initial input value.
  - `options` (Object, optional) - Additional config.
- **Returns**: `Promise<String|null>` - Resolves to input value or `null` if cancelled.
- **Example**:
  ```javascript
  const email = await Alpine.store('uno').box.prompt('Email:')
  ```

### `close(id)`

Closes a specific dialog instance.

- **Parameters**: `id` (String) - The unique ID of the dialog.
- **Returns**: `void`
- **Example**:
  ```javascript
  Alpine.store('uno').box.close('dialog-123')
  ```

### `closeAll()`

Closes all currently open dialogs.

- **Parameters**: None
- **Returns**: `void`
- **Example**:
  ```javascript
  Alpine.store('uno').box.closeAll()
  ```

## Examples

### Dialog with Extra Buttons

Add a "Save Draft" button next to the primary "Publish" action.

```javascript
const boxInstance = Alpine.store('uno').box.dialog({
  title: 'New Diary Entry',
  hxUrl: '/entry/new',
  showOkButton: true,
  okText: 'Publish',
  okButtonClass: 'btn-primary',
  extraButtons: [
    {
      text: 'Save Draft',
      position: 'before', // renders left of "Publish"
      buttonClass: '!font-normal !ring-0 !shadow-none !bg-transparent',
      onClick: async (item) => {
        await saveDraft() // your save logic
        showNotification('Draft saved!')
        return false // keep the dialog open
      }
    }
  ],
  onOk: async () => {
    await publishEntry()
  }
})

// Update extraButtons at runtime (e.g. after content loads)
boxInstance.extraButtons = isDraft
  ? [
      {
        text: 'Save Draft',
        position: 'before',
        onClick: () => {
          /* ... */ return false
        }
      }
    ]
  : []
```

### Custom Styled Dialog

```javascript
Alpine.store('uno').box.open({
  title: 'Profile Settings',
  widthClass: 'w-full max-w-2xl',
  content: `
    <div class="grid gap-4">
      <div class="p-4 bg-gray-50 rounded">User info...</div>
    </div>
  `,
  okText: 'Save Changes',
  showCancelButton: true,
  onOk: () => {
    // Perform save logic
    console.log('Saving...')
  }
})
```

### HTMX Dynamic Content

Load a form from the server and handle the response.

```javascript
Alpine.store('uno').box.open({
  title: 'Edit User',
  hxUrl: '/users/1/edit',
  onContentLoad: (el) => {
    console.log('Form loaded into dialog')
  }
})
```

## Troubleshooting

### Common Issues

| Issue                       | Possible Cause                                | Solution                                                                                 |
| :-------------------------- | :-------------------------------------------- | :--------------------------------------------------------------------------------------- |
| **Dialog not opening**      | Alpine.js not initialized or store not ready. | Ensure `Alpine.start()` is called and `uno` store is accessed after `alpine:init`.       |
| **HTMX content missing**    | `hxUrl` returned 404 or empty response.       | Check Network tab. Ensure server returns partial HTML, not full page.                    |
| **Styling broken**          | Tailwind CSS not loaded or version mismatch.  | Ensure Tailwind 4.x is loaded. Check overrides applied via `Uno.setConfig({ ui: ... })`. |
| **"Alpine is not defined"** | Script loading order incorrect.               | Load Alpine.js _after_ your app scripts or use `defer`.                                  |

## Version History

| Version    | Date       | Changes                                                                                                                                                                                                |
| :--------- | :--------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Unreleased | —          | Focus returns to the opener on close; underlay dialogs are marked `inert` + `aria-hidden`; the focus trap is DOM-order aware so Tab cannot escape.                                                     |
| v0.9.0     | 2026-08-26 | Added `extraButtons` support: custom buttons with `position`, `buttonClass`, and async `onClick` (`return false` to keep open). Added `handleExtraButton()` store method and `extraBtn` default style. |
| v0.8.0     | 2026-01-22 | Added `bindNotifyToHtmx()` method to HTMX integration for automatic notifications.                                                                                                                     |
| v0.7.0     | 2026-01-20 | Refactored for **Lazy Initialization** (DOM elements created on demand).                                                                                                                               |
| v0.6.0     | 2025-12-18 | Initial stable release.                                                                                                                                                                                |
