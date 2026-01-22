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
| `pwd`       | Basic password (8+ chars, letter+num).          | `^(?=.*[A-Za-z])(?=.*\d)[A-Za-z\d]{8,}$`                                    |
| `strongpwd` | Strong password (8+ chars, letter+num+special). | `^(?=.*[A-Za-z])(?=.*\d)(?=.*[@$!%*#?&])[A-Za-z\d@$!%*#?&]{8,}$`            |

> **Note**: In the regex column, pipes `|` are escaped as `\|` or `&#124;` to prevent table formatting errors. The actual regexes in source code use standard `|`.

### Logic Validators

| Type      | Description                    | Options                          |
| :-------- | :----------------------------- | :------------------------------- |
| `length`  | String length check.           | `min` (number), `max` (number)   |
| `range`   | Numeric value range check.     | `min` (number), `max` (number)   |
| `pattern` | Custom Regex match.            | `pattern` (RegExp string), `msg` |
| `match`   | Matches another field's value. | `match` (selector/name), `msg`   |

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

## Version History

| Version | Date       | Changes                                                                                                        |
| :------ | :--------- | :------------------------------------------------------------------------------------------------------------- |
| v0.8.0  | 2026-01-22 | **Breaking**: Removed HTML5 attribute support. Added `debounce`, `emptyQuit`, and `type: 'validator'` support. |
| v0.7.0  | 2026-01-20 | Initial release.                                                                                               |
