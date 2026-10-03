# Uno Valid Plugin API

The `Valid` plugin provides a lightweight, directive-driven form validation engine for Alpine.js. It supports synchronous and asynchronous validation, custom rules, and seamless HTMX integration.

## Key Features

- **Directive-Driven**: Validation rules are defined entirely within `x-uno-valid`. No dependency on HTML5 attributes (`required`, `pattern`, etc.).
- **User-Defined Order**: Rules are executed in the exact order defined in the configuration array.
- **EmptyQuit Strategy**: By default, validation skips non-required rules if the field is empty, reducing unnecessary errors.
- **Async Support**: Native support for asynchronous validators with `debounce` to reduce server load.
- **HTMX Integration**: Automatically intercepts `htmx:confirm` to prevent requests if validation fails.

## Installation

The Valid plugin is initialized automatically by Uno.

```javascript
import Uno from 'uno-ui'
Uno.init() // Includes Valid plugin
```

## Usage

### Basic Syntax

Apply the `x-uno-valid` directive to input fields. The value should be an **Array** of rules.

```html
<input name="username" x-uno-valid="['required', { type: 'min', min: 3 }]" />
```

### Rule Configuration

Rules can be specified as:

1.  **String**: Short for `{ type: 'ruleName' }`.
2.  **Object**: Full configuration `{ type: 'ruleName', msg: 'Custom error', ...args }`.

**Example: Mixed Configuration**

```html
<input
  name="email"
  x-uno-valid="[
    'required', 
    { type: 'email', msg: 'Invalid email format' },
    { type: 'blacklist', words: ['admin', 'root'] }
  ]"
/>
```

### EmptyQuit Behavior

By default, if a field is empty and **not** required, subsequent rules are skipped. This prevents errors like "Format invalid" for optional fields that are left blank.

To force validation on empty fields, set `emptyQuit: false` on the specific rule.

```html
<!-- Optional field, but if filled, must be an email -->
<input x-uno-valid="['email']" />

<!-- Optional field, but runs 'complexCheck' even if empty -->
<input x-uno-valid="[{ type: 'complexCheck', emptyQuit: false }]" />
```

### Async Validation & Debounce

For expensive async validations (e.g., checking username availability), use the `debounce` option (in milliseconds).

```html
<input
  name="username"
  x-uno-valid="[
    'required',
    { 
      type: 'uniqueUsername', 
      debounce: 500, 
      msg: 'Username already taken' 
    }
  ]"
/>
```

The validator will wait 500ms after the last input before executing `uniqueUsername`.

### Validation Triggers

By default, validation runs on both `input` and `blur` events. You can restrict a rule to run only on `blur` by setting `trigger: 'blur'`.

```html
<!-- Validates format on input, but checks availability only on blur -->
<input x-uno-valid="['required', 'validname', { type: 'checkRemote', trigger: 'blur' }]" />
```

### Auto Error Generation (New in v0.8.0)

The plugin automatically generates and manages error message elements. You do NOT need to manually add error containers.

- **Placement**: Error messages are inserted immediately after the input element.
- **Styling**: Uses the configured error class (default: `text-red-500 text-xs mt-1`).
- **Behavior**:
  - Shows error when validation fails.
  - Hides/removes error when validation passes.
  - Updates message automatically based on the failing rule.

#### Accessibility (WCAG 3.3.1 / 4.1.3)

Error exposure to assistive technology is wired automatically — no extra markup required:

- **Error element**: gets `role="alert"` (so the message is announced when it appears) and a stable id `${fieldName}-error`. An author-provided `id` or `role` on a hand-written error element is preserved.
- **Input**: its `aria-describedby` is pointed at the error element id (appended to any existing describedby targets, never clobbering them). This applies both to auto-generated and hand-written error elements.
- **`aria-invalid`**: the store mirrors each field's error state onto its control's `aria-invalid` attribute at every `errors` mutation site (`validate`, `showError`, `hideError`, `clearErrors`, `showErrors`), so client-side failures and `422` server errors both toggle it. The attribute is removed once the field is valid again.

> When `autoError: false` and no matching hand-written error element exists, no `aria-describedby` is wired (there is nothing to point at); `aria-invalid` sync still applies.

#### Configuration

You can configure the auto-error behavior via `Uno.options.valid`:

- `autoError` (boolean, default: `true`): Enable or disable auto-error generation globally.
- `errorClass` (string, default: `text-red-500 text-xs mt-1`): CSS classes to apply to the generated error element.
- `createErrorElement` (function): Custom function to create the error element.

```javascript
// Example: Custom configuration
Uno.options.valid = {
  autoError: true,
  errorClass: 'text-danger small mt-2',
  createErrorElement: (fieldName, errorKey) => {
    const div = document.createElement('div')
    div.className = 'custom-error'
    return div
  }
}
```

```html
<!-- No manual error div needed! -->
<input name="email" x-uno-valid="['required', 'email']" />
<!-- Auto-generated error will appear here -->
```

## Built-in Validators

The plugin includes a comprehensive set of validators. Use them by specifying the `type` property (or just the string name).

### Standard Types

| Type        | Description                                     | Regex / Logic                                                               |
| :---------- | :---------------------------------------------- | :-------------------------------------------------------------------------- |
| `required`  | Value must not be empty.                        | Checks `trim() !== ''`                                                      |
| `email`     | Email address format.                           | `^[a-zA-Z0-9.!#$%&'*+/=?^_`{&#124;}~-]+@[a-zA-Z0-9-]+(?:\.[a-zA-Z0-9-]+)+$` |
| `url`       | URL format (http/https).                        | `^(https?:\/\/)?([\da-z.-]+)\.([a-z.]{2,6})([/\w .-]*)*\/?$/i`              |
| `date`      | Date format (YYYY-MM-DD, DD/MM/YYYY, etc.).     | `^(\d{4}-\d{2}-\d{2}\|\d{2}\/\d{2}\/\d{4}\|\d{2}\.\d{2}\.\d{4})$`           |
| `int`       | Integer (positive or negative).                 | `^-?\d+$`                                                                   |
| `number`    | General number (integer or float).              | `^-?\d*(\.\d+)?$`                                                           |
| `validname` | Alphanumeric, underscores, hyphens.             | `^[a-zA-Z0-9_-]+$`                                                          |
| `en`        | English letters only.                           | `^[a-zA-Z]+$`                                                               |
| `domain`    | Domain name format.                             | `^((?!-)[A-Za-z0-9-]{1,63}(?<!-)\.)+[A-Za-z]{2,6}$`                         |
| `mobile`    | Mobile phone number (7+ digits).                | `^\+?[\d\s-]{7,}$`                                                          |
| `ipv4`      | IPv4 Address.                                   | `^(?:(?:25[0-5]\|2[0-4][0-9]\|[01]?[0-9][0-9]?)\.){3}(?:...)$`              |
| `hex`       | Hex color code.                                 | `^#?([a-f0-9]{6}\|[a-f0-9]{3})$/i`                                          |
| `pwd`       | Basic password (8-72 chars, letter+num).        | `^(?=.*[A-Za-z])(?=.*\d)[A-Za-z\d!_@#$%^&*()+\-=]{8,72}$`                   |
| `strongpwd` | Strong password (8+ chars, letter+num+special). | `^(?=.*[A-Za-z])(?=.*\d)(?=.*[@$!%*#?&])[A-Za-z\d@$!%*#?&]{8,}$`            |

> **Note**: In the regex column, pipes `|` are escaped as `\|` or `&#124;` to prevent table formatting errors. The actual regexes in source code use standard `|`.

### Logic Validators

| Type      | Description                    | Options                          |
| :-------- | :----------------------------- | :------------------------------- |
| `length`  | String length check.           | `min` (number), `max` (number)   |
| `range`   | Numeric value range check.     | `min` (number), `max` (number)   |
| `pattern` | Custom Regex match.            | `pattern` (RegExp string), `msg` |
| `match`   | Matches another field's value. | `match` (selector/name), `msg`   |

## Feature: Length Control & Counter

The `length` rule supports additional UI controls:

- **`showCounter`** (boolean, default: `false`): Automatically displays a character counter (e.g., "5/100") in the bottom-right corner of the input. The counter element is linked to the input via `aria-describedby` (B6), so screen readers announce the remaining-length context on focus; existing describedby tokens (e.g. the error element) are preserved. The counter also carries a full-sentence `aria-label` (i18n key `valid.counter`, fallback `"5/100"`) while the visual text stays compact, and over-limit input is exposed via `aria-invalid` by the standard length-rule validation.
- **`stopOnMax`** (boolean, default: `false`): If `true`, prevents the user from typing more characters once the `max` limit is reached (sets `maxlength` attribute).

```html
<textarea
  x-uno-valid="[{ 
    type: 'length', 
    max: 200, 
    showCounter: true, 
    stopOnMax: true 
  }]"
></textarea>
```

## Custom Validators

Uno Valid provides a flexible system for defining custom validation logic. You can register reusable validator functions and invoke them via the `type` property.

### 1. Registration

First, register your validator function using `Alpine.store('uno').valid.registerValidator`.

```javascript
// app.js
document.addEventListener('alpine:init', () => {
  const uno = Alpine.store('uno')

  // Synchronous Validator
  uno.valid.registerValidator('blacklist', (value, fieldName, rule) => {
    const forbidden = rule.words || []
    if (forbidden.includes(value.toLowerCase())) {
      return rule.msg || `The word "${value}" is not allowed.`
    }
    return true
  })

  // Asynchronous Validator
  uno.valid.registerValidator('checkStock', async (value, fieldName, rule) => {
    try {
      const response = await fetch(`/api/check-stock?sku=${value}`)
      const data = await response.json()
      return data.available ? true : 'Out of stock'
    } catch (e) {
      console.error('Validation error', e)
      return 'Validation failed (Network Error)'
    }
  })
})
```

### 2. Usage Strategy: The `validator` Type

While you can often just use the registered name as the type (e.g., `{ type: 'blacklist' }`), the dedicated `type: 'validator'` pattern is recommended for advanced scenarios where you need to pass specific parameters or want explicit configuration.

#### Basic Syntax

```javascript
{
  type: 'validator',
  name: 'registeredName',  // MUST match the name used in registerValidator
  // ... custom parameters passed to the rule object
  // ... standard parameters (msg, debounce, etc.)
}
```

#### Callback Signature

Your registered function receives three arguments:

```javascript
(value, fieldName, rule) => { ... }
```

1.  **`value`** (`Any`): The current value of the input field.
2.  **`fieldName`** (`String`): The `name` or `id` attribute of the input field.
3.  **`rule`** (`Object`): The complete configuration object defined in `x-uno-valid`. This gives you access to custom properties (e.g., `words`, `apiEndpoint`).

#### Return Value Specifications

The validator function must return one of the following:

- **`true`** (Boolean): Validation Passed.
- **`false`** (Boolean): Validation Failed. Uses `rule.msg` or a default generic error.
- **`String`**: Validation Failed. The string becomes the error message.
- **`Object`**: `{ valid: false, message: 'Custom error' }`.
- **`Promise`**: For async validators. Must resolve to one of the above.

#### Error Handling

- **Exceptions**: If a validator throws an exception (sync or async), the validation is treated as **Failed**.
- **Network Errors**: For async validators, it is best practice to `catch` fetch errors and return a user-friendly string (e.g., "Network error, please try again") rather than letting the promise reject.

### 3. Examples

#### Scenario A: Advanced Configuration (Params)

Passing custom parameters (`minAge`, `country`) to a reusable validator.

```javascript
// Registration
uno.valid.registerValidator('ageCheck', (value, field, rule) => {
  const min = rule.minAge || 18
  return parseInt(value) >= min ? true : `Must be at least ${min} years old.`
})
```

```html
<!-- Usage -->
<input
  name="age"
  x-uno-valid="[
    'required',
    { 
      type: 'validator', 
      name: 'ageCheck', 
      minAge: 21, 
      msg: 'Underage' 
    }
  ]"
/>
```

#### Scenario B: Async Integration with Debounce

Validating against a backend API with debouncing to prevent flooding.

```javascript
// Registration
uno.valid.registerValidator('uniqueUser', async (value) => {
  if (!value) return true
  const res = await fetch(`/api/users/check?name=${value}`)
  const isAvailable = await res.json()
  return isAvailable ? true : 'Username is taken'
})
```

```html
<!-- Usage -->
<input
  name="username"
  x-uno-valid="[
    'required',
    { 
      type: 'validator', 
      name: 'uniqueUser', 
      debounce: 500,
      msg: 'User exists' 
    }
  ]"
/>
```

#### Scenario C: Cross-Field Validation (Integration)

Validating one field based on another (e.g., Password Confirmation).

```javascript
// Registration
uno.valid.registerValidator('matchField', (value, field, rule) => {
  const targetSelector = rule.target
  const targetEl = document.querySelector(targetSelector)
  if (!targetEl) return true // Skip if target not found

  return value === targetEl.value ? true : rule.msg || 'Fields do not match'
})
```

```html
<!-- Usage -->
<input type="password" id="pwd" name="password" />
<input
  type="password"
  name="confirm"
  x-uno-valid="[
    { 
      type: 'validator', 
      name: 'matchField', 
      target: '#pwd',
      msg: 'Passwords must match' 
    }
  ]"
/>
```

### 4. Inline Regular Expressions

For simple pattern matching without registering a global function, use the `pattern` type.

```html
<input
  x-uno-valid="[{ 
    type: 'pattern', 
    pattern: '^[A-Z0-9]+$', 
    msg: 'Only uppercase alphanumeric characters allowed' 
  }]"
/>
```

## Troubleshooting

### Common Issues

| Issue                              | Possible Cause                                                        | Solution                                                                 |
| :--------------------------------- | :-------------------------------------------------------------------- | :----------------------------------------------------------------------- |
| **Validator not running**          | 1. Field is empty and `emptyQuit` is true.<br>2. Validator name typo. | 1. Set `emptyQuit: false` in rule.<br>2. Check `registerValidator` name. |
| **Async loop / Too many requests** | Missing `debounce`.                                                   | Add `debounce: 500` to your rule config.                                 |
| **"HTML5 validation ignored"**     | v0.8.0+ removed support for `required` attribute.                     | Move all rules into `x-uno-valid="['required', ...]"`.                   |
| **Validation runs on typing**      | Default trigger is `input`.                                           | Add `trigger: 'blur'` to the rule object.                                |
| **Custom Error not showing**       | Validator returns `false` but no `msg` provided.                      | Return a string from validator or set `msg` property in rule.            |

## Methods

除了声明式指令，你还可以通过 Store API 编程式地控制验证：

```js
const valid = Alpine.store('uno').valid
```

### `addRules(fieldName, config, instanceOpts = {})`

添加或更新字段的验证规则。

| Parameter      | Type                    | Description                 |
| :------------- | :---------------------- | :-------------------------- |
| `fieldName`    | `string`                | 字段唯一标识符 (name 或 id) |
| `config`       | `Array\|Object\|string` | 验证配置                    |
| `instanceOpts` | `Object`                | 实例级全局选项 (可选)       |

```js
valid.addRules('email', ['required', 'email'])
valid.addRules('password', { type: 'minLength', min: 8 })
```

### `validate(fieldName, value, eventType = 'input')`

验证单个字段。

| Parameter   | Type             | Description                                    |
| :---------- | :--------------- | :--------------------------------------------- |
| `fieldName` | `string`         | 字段标识符                                     |
| `value`     | `string\|number` | 要验证的值                                     |
| `eventType` | `string`         | 触发验证的事件类型 ('input', 'blur', 'submit') |

```js
const isValid = await valid.validate('email', 'test@example.com', 'blur')
```

### `showError(fieldName, message)`

手动显示字段错误。

| Parameter   | Type     | Description |
| :---------- | :------- | :---------- |
| `fieldName` | `string` | 字段标识符  |
| `message`   | `string` | 错误消息    |

```js
valid.showError('email', 'Invalid email format')
```

### `hideError(fieldName)`

手动隐藏字段错误。

```js
valid.hideError('email')
```

### `clearErrors()`

清除所有字段错误。

```js
valid.clearErrors()
```

### `showErrors(errorsObj)`

批量显示多个字段错误。

| Parameter   | Type     | Description                |
| :---------- | :------- | :------------------------- |
| `errorsObj` | `Object` | 字段名到错误消息的映射对象 |

```js
valid.showErrors({
  email: 'Invalid email',
  password: 'Password too weak'
})
```

### `registerValidator(name, validatorFn)`

注册自定义验证器。

| Parameter     | Type       | Description |
| :------------ | :--------- | :---------- |
| `name`        | `string`   | 验证器名称  |
| `validatorFn` | `Function` | 验证函数    |

```js
valid.registerValidator('customRule', ({ value, rule }) => {
  // 自定义验证逻辑
  return value.length > rule.min
})
```

## Version History

| Version | Date       | Changes                                                                                                                                                                                                                                                    |
| :------ | :--------- | :--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| v0.9.0  | 2026-09-30 | **A11y**: auto error elements get `role="alert"` + stable id; inputs are wired with `aria-describedby` and `aria-invalid` (synced on client validation and 422 server errors); `showCounter` counters are linked to the input via `aria-describedby` (B6). |
| v0.8.1  | 2026-02-28 | Added programmatic Store API documentation (addRules, validate, showError, hideError, clearErrors, showErrors).                                                                                                                                            |
| v0.8.0  | 2026-01-22 | **Breaking**: Removed HTML5 attribute support. Added `debounce`, `emptyQuit`, and `type: 'validator'` support.                                                                                                                                             |
| v0.7.0  | 2026-01-20 | Initial release.                                                                                                                                                                                                                                           |
