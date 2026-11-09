# Valid API

Valid is a directive-driven form validation engine. Rules are declared as an ordered array in
`x-uno-valid`; the store exposes validation state and methods for programmatic use. It supports
sync and async validators, debounce, character counters, server-side `422` errors and HTMX form
interception.

## Prerequisites

- Alpine.js 3 with Uno registered. See [Getting Started](getting-started.md).
- HTMX 2 for the `htmx:confirm` interception and `422` handling.

## Minimal example

```html
<form x-uno-valid @submit.prevent>
  <label for="email">Email</label>
  <input id="email" name="email" type="email" x-uno-valid="['required', 'email']" />
  <button type="submit">Save</button>
</form>
```

The directive attaches `input` and `blur` listeners and creates the error element after the
input. The form marker `x-uno-valid` on `<form>` is what makes HTMX submission validate the
whole form.

## Rules

The value is an array. A single string or object is normalized to an array, but the array form
is the documented contract and preserves order.

```html
<input
  name="username"
  x-uno-valid="[
    'required',
    { type: 'length', min: 3, max: 16, showCounter: true },
    { type: 'pattern', pattern: '^[a-z0-9_]+$', msg: 'Lowercase letters, digits and _ only' }
  ]"
/>
```

Each item is either a string (shorthand for `{ type: 'name' }`) or an object.

### Rule options

| Name | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `type` | `string` | — | Validator name. |
| `msg` | `string` | — | Error message. Falls back to the i18n key for the validator. |
| `required` | `boolean` | `false` | Treat an empty value as an error even when `type` is not `required`. |
| `min` | `number` | — | Minimum length or numeric value, depending on the validator. |
| `max` | `number` | — | Maximum length or numeric value. |
| `pattern` | `string \| RegExp` | — | Regular expression for the `pattern` validator. |
| `match` | `string` | — | Target selector/name for `match`; `selector` and `target` are aliases. |
| `minLength` / `maxLength` | `number` | — | Read by the `length` counter/`stopOnMax` helper, not by the validator itself. Use `min`/`max` for validation. |
| `debounce` | `number` | `0` | Milliseconds to wait before validating on `input`. Ignored on `blur`/`submit`. |
| `emptyQuit` | `boolean` | `true` | Skip non-required rules when the value is empty. |
| `trigger` | `string` | — | `'blur'` skips this rule on the input pass. Other values have no effect. |
| `realtime` | `boolean` | — | `false` skips this rule on the input pass (same effect as `trigger: 'blur'`). |
| `showCounter` | `boolean` | `false` | For `length`: render a character counter when `max` is set. |
| `stopOnMax` | `boolean` | `false` | For `length`: set `maxlength` when `max` is set. |
| `name` | `string` | — | Registered validator name for `{ type: 'validator' }`. |
| `fn` | `function` | — | Inline validator function. |

### Empty values

By default an empty value skips every rule except `required`:

```html
<!-- Optional, but if filled must be an email -->
<input name="backup" x-uno-valid="['email']" />

<!-- Runs even when empty -->
<input name="code" x-uno-valid="[{ type: 'complexCheck', emptyQuit: false }]" />
```

`emptyQuit` is merged from `Uno.getConfig('valid')` into every rule; a rule-level value wins.
Set `valid.emptyQuit: false` globally if you need the opposite default.

### Debounce

```html
<input
  name="username"
  x-uno-valid="[
    'required',
    { type: 'validator', name: 'uniqueUsername', debounce: 500 }
  ]"
/>
```

- Debounce only applies to the `input` event.
- While waiting, `validating[fieldName]` is `true`.
- The current implementation pauses at the first debounced rule: rules after it in the same
  input pass are not evaluated. Put the debounced rule last.
- `blur` and `submit` run immediately.

## Built-in validators

| Type | Description | Rule |
| :--- | :--- | :--- |
| `required` | Non-empty after trim. | `String(value).trim() !== ''` |
| `email` | Email address. | ``^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9-]+(?:\.[a-zA-Z0-9-]+)+$`` |
| `url` | HTTP/HTTPS URL. | `^(https?:\/\/)?([\da-z.-]+)\.([a-z.]{2,6})([/\w .-]*)*\/?$` (case-insensitive) |
| `date` | `YYYY-MM-DD`, `DD/MM/YYYY` or `DD.MM.YYYY`. | `^(\d{4}-\d{2}-\d{2}\|\d{2}\/\d{2}\/\d{4}\|\d{2}\.\d{2}\.\d{4})$` |
| `int` | Integer. | `^-?\d+$` |
| `number` | Integer or decimal. | `^-?\d*(\.\d+)?$` |
| `validname` | Letters, digits, `-`, `_`; 3–16 chars. | `^[a-zA-Z0-9_-]{3,16}$` |
| `en` | English letters only. | `^[a-zA-Z]+$` |
| `domain` | DNS domain. | `^((?!-)[A-Za-z0-9-]{1,63}(?<!-)\.)+[A-Za-z]{2,6}$` |
| `mobile` | Optional `+`, digits/space/`-`, 7+ chars. | `^\+?[\d\s-]{7,}$` |
| `ipv4` | IPv4 address. | `^(?:(?:25[0-5]\|2[0-4][0-9]\|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]\|2[0-4][0-9]\|[01]?[0-9][0-9]?)$` |
| `hex` | Hex color, with or without `#`, 3 or 6 digits. | `^#?([a-f0-9]{6}\|[a-f0-9]{3})$` (case-insensitive) |
| `pwd` | 8–72 chars, at least one letter and one digit; `_!@#$%^&*()+-=` allowed. | `^(?=.*[A-Za-z])(?=.*\d)[A-Za-z\d!_@#$%^&*()+\-=]{8,72}$` |
| `strongpwd` | 8–80 chars, letter + digit + one of `@$!%*#?&`. | `^(?=.*[A-Za-z])(?=.*\d)(?=.*[@$!%*#?&])[A-Za-z\d@$!%*#?&]{8,80}$` |

The `|` characters inside the regex column are escaped for table rendering; the source uses
normal alternation. See [src/utils/validation.js](../src/utils/validation.js) for the exact
patterns.

### Length and range

| Type | Options | Behavior |
| :--- | :--- | :--- |
| `length` | `min`, `max`, `showCounter`, `stopOnMax` | Compares `String(value).length`. |
| `range` | `min`, `max` | Parses with `parseFloat`; a non-number fails with the number message. |

### Pattern and match

```html
<input x-uno-valid="[{ type: 'pattern', pattern: '^[A-Z]{2}-\d{4}$', msg: 'Use AB-1234' }]" />
<input name="confirm" x-uno-valid="[{ type: 'match', target: 'password', msg: 'Passwords do not match' }]" />
```

`match` resolves `rule.match || rule.selector || rule.target` as a selector. A value starting
with `#` is used as a CSS selector; otherwise it is used as `[name="..."]`, with an id fallback.
Comparison is strict `!==` against `target.value`.

## Auto error generation

When `autoError` is not `false`, the directive creates an error element after the input:

- Id: `${fieldName}-error`.
- `role="alert"`.
- `x-text` and `x-show` bound to `$store.uno.valid.errors['fieldName']`.
- Class from `valid.errorClass`, falling back to `ui.valid.helperText`:
  `text-red-600 dark:text-red-400 text-xs mt-1`.
- Inserted immediately after the input; `Alpine.initTree()` is called on it.
- The input's `aria-describedby` receives the error id, preserving existing tokens.

If a hand-written element already matches
`[x-text="$store.uno.valid.errors['fieldName']"]`, it is reused and receives the same id/role
wiring.

```javascript
Uno.setConfig({
  valid: {
    autoError: false,
    errorClass: 'text-danger small mt-2',
    createErrorElement: (fieldName, errorKey) => {
      const div = document.createElement('div')
      div.setAttribute('x-text', errorKey)
      div.setAttribute('x-show', errorKey)
      div.className = 'custom-error'
      return div
    }
  }
})
```

The custom factory must produce an element that displays the message (for example by setting
`x-text`/`x-show` as above).

## Character counter

```html
<textarea
  name="bio"
  x-uno-valid="[{ type: 'length', max: 200, showCounter: true, stopOnMax: true }]"
></textarea>
```

- Counter id: `${fieldName}-counter`.
- Text: `${current}/${max}`.
- `aria-describedby` is added to the input; the counter's `aria-label` uses
  `valid.counter` with `{current}` and `{max}`.
- The counter is gray at or below the limit and red above it.
- `stopOnMax` sets `maxlength` so the browser blocks extra input.
- The input is wrapped in a plugin-managed positioned container; the wrapper is removed on
  Alpine cleanup.

## Custom validators

Register a function before the field initializes, then reference it by `type` or via
`{ type: 'validator', name }`.

```javascript
document.addEventListener('alpine:init', () => {
  const valid = Alpine.store('uno').valid

  valid.registerValidator('blacklist', (value, fieldName, rule) => {
    const forbidden = rule.words || []
    return forbidden.includes(value.toLowerCase())
      ? rule.msg || `The word "${value}" is not allowed.`
      : true
  })

  valid.registerValidator('checkStock', async (value) => {
    const response = await fetch(`/api/stock?sku=${encodeURIComponent(value)}`)
    const data = await response.json()
    return data.available ? true : 'Out of stock'
  })
})
```

```html
<input
  name="sku"
  x-uno-valid="[
    'required',
    { type: 'validator', name: 'checkStock', debounce: 500 }
  ]"
/>
```

### Callback signature

```javascript
(value, fieldName, rule) => result
```

| Argument | Type | Description |
| :--- | :--- | :--- |
| `value` | `any` | Current field value. |
| `fieldName` | `string` | The field's `name` or `id`. |
| `rule` | `object` | The complete rule object. |

### Return values

| Return | Result |
| :--- | :--- |
| `true` / `undefined` | Valid. |
| `false` | Invalid; uses `rule.msg` or `Validation failed`. |
| `string` | Invalid; the string is the message. |
| `{ valid: false, message }` | Invalid; `message` (or `rule.msg`). |
| `Promise` of the above | Awaited. |

- A thrown exception is treated as invalid with the message `Validation error`.
- A validator that returns `{ valid: true }` is treated as valid.
- Unknown validator names log a warning and do not fail the field.

`{ type: 'validator', name: 'x' }` is the supported explicit form. A registered name can also be
used directly as `type` (`{ type: 'blacklist' }`), and `rule.fn` accepts an inline function.

## Store API

```javascript
const valid = Alpine.store('uno').valid
```

### `addRules(fieldName, config, instanceOpts)`

Adds or replaces the rules for a field and merges global options.

```javascript
valid.addRules('email', ['required', 'email'])
valid.addRules('password', { type: 'length', min: 8 })
```

### `validate(fieldName, value, eventType)`

Validates one field and returns a Promise of `boolean`. `eventType` defaults to `'input'`;
`'blur'` and `'submit'` skip the trigger filters. A debounced rule may resolve later.

```javascript
const ok = await valid.validate('email', 'user@example.com', 'blur')
```

### `validateForm(form)`

Validates every form control that has registered rules. Returns a Promise of `boolean`.

- Uses the `submit` event type so every rule runs.
- With two or more failures, renders a focusable `role="alert"` summary at the top of the form
  and focuses it.
- With one failure, focuses the first failing field.
- `scrollOnError: false` disables scrolling/focus.
- Removes a previous summary before re-evaluating.

### `showError(fieldName, message)` / `hideError(fieldName)`

Set or remove one field error, updating `aria-invalid`.

### `clearErrors()`

Clears all errors and updates `aria-invalid` for each affected field.

### `showErrors(errors)`

Merges a map of field errors and updates `aria-invalid` for each:

```javascript
valid.showErrors({
  email: 'Email is already registered',
  password: 'Password is too weak'
})
```

### `registerValidator(name, fn)`

Registers a custom validator function.

## HTMX

The plugin listens on `document.body` for `htmx:confirm`. A form is intercepted when it has
`x-uno-valid` **or** contains a field with `x-uno-valid`:

```html
<form x-uno-valid hx-post="/save" hx-swap="none">
  <input name="email" x-uno-valid="['required', 'email']" />
  <button type="submit">Save</button>
</form>
```

- Invalid form: the request is blocked.
- Valid form: `evt.detail.issueRequest(true)` continues.
- HTTP `422` with `{ "errors": { "field": "message" } }` calls `showErrors()`.
- Non-JSON 422 bodies are ignored.
- An unexpected exception during validation fails closed (request blocked, error logged).

See [HTMX integration](htmx.md#valid).

## Accessibility

- Generated and reused error elements get `role="alert"` and `${field}-error`.
- The control's `aria-describedby` receives the error id, appended to existing tokens.
- `aria-invalid` is set on error and removed when the field becomes valid, including `422`
  errors.
- The counter is linked by `aria-describedby` and has a full-sentence `aria-label`.
- Multi-field failures render a focusable summary with links that focus the field.
- If `autoError: false` and no matching manual error element exists, there is nothing to point
  `aria-describedby` at; `aria-invalid` still syncs.

## Notes

- Rules are keyed by `name || id`; the field must have one.
- `x-uno-valid` on a `<form>` returns early in the directive and acts as an HTMX marker; it does
  not create a field error.
- Validation always listens to `input` and `blur`; there is no separate `change` event.
- The error element is created once at directive initialization. Changing rules later does not
  re-create counters/errors.
- `maxLength` alone does not fail validation; use `max` with `type: 'length'`.
- The current `valid.range` i18n key is not used by `executeRule()`; range failures use
  `valid.rangeMin` / `valid.rangeMax`.

## Version history

| Version | Date | Changes |
| :--- | :--- | :--- |
| 0.10.0 | 2026-11-09 | `role="alert"` + `${field}-error`, `aria-describedby`, `aria-invalid`, counter association, error summary. |
| 0.8.0 | 2026-01-22 | Breaking: HTML5 attribute rules removed. Added debounce, emptyQuit, `type: 'validator'`; documented store API. |
| 0.7.0 | 2026-01-20 | Initial release. |

## Related

- [HTMX integration](htmx.md#valid)
- [Accessibility](accessibility.md)
- [i18n](i18n.md)
- [Configuration](configuration.md)
