# Uno Valid API Documentation

The Valid component is a comprehensive form validation system for Alpine.js and HTMX. It supports HTML5 constraints, custom rules, real-time feedback, and integration with HTMX request lifecycles.

## Overview

- **Core Function**: Client-side form validation with server-side error handling integration.
- **Architecture**: Store (`Alpine.store('uno').valid`) manages rules and errors; Directive (`x-uno-valid`) binds fields.
- **HTMX Integration**: Intercepts `htmx:confirm` to prevent invalid requests; handles 422 responses from `htmx:afterRequest`.
- **Extensibility**: Supports custom validator functions and regex patterns.

## Basic Usage

### 1. Field Validation

Add `x-uno-valid` to inputs along with standard HTML5 attributes.

```html
<input type="email" name="email" required x-uno-valid />
<!-- Error Display -->
<span x-text="$store.uno.valid.errors['email']" class="text-red-500"></span>
```

### 2. Form Setup (HTMX)

Add `x-uno-valid` to the form to enable HTMX interception.

```html
<form hx-post="/save" x-uno-valid>
  <input name="username" required x-uno-valid />
  <button type="submit">Submit</button>
</form>
```

## API Parameters

### Configuration Object (`x-uno-valid="{...}"`)

| Name          | Type                 | Default | Required | Description                                                                                                    |
| :------------ | :------------------- | :------ | :------- | :------------------------------------------------------------------------------------------------------------- |
| `required`    | `Boolean`            | `false` | No       | Marks field as required. (Or use `required` attribute).                                                        |
| `type`        | `String`             | -       | No       | Validator type. Values: `'email'`, `'url'`, `'date'`, `'number'`, `'int'`, `'mobile'`, `'pwd'`, `'strongpwd'`. |
| `minLength`   | `Number`             | -       | No       | Minimum character length.                                                                                      |
| `maxLength`   | `Number`             | -       | No       | Maximum character length.                                                                                      |
| `min`         | `Number`             | -       | No       | Minimum numeric value.                                                                                         |
| `max`         | `Number`             | -       | No       | Maximum numeric value.                                                                                         |
| `pattern`     | `RegExp`             | -       | No       | Custom Regex pattern.                                                                                          |
| `match`       | `String`             | -       | No       | Name or ID of another field to match (e.g., password confirmation).                                            |
| `validator`   | `String`\|`Function` | -       | No       | Name of a registered custom validator or a function.                                                           |
| `realtime`    | `Boolean`            | `true`  | No       | If `false`, validation only runs on blur, not input.                                                           |
| `msgRequired` | `String`             | -       | No       | Custom error message for required check.                                                                       |
| `msgType`     | `String`             | -       | No       | Custom error message for type check.                                                                           |
| `msgMatch`    | `String`             | -       | No       | Custom error message for match check.                                                                          |

### Built-in Validator Types

- `email`: Standard email format.
- `url`: HTTP/HTTPS URLs.
- `date`: YYYY-MM-DD formats.
- `int`: Integers only.
- `number`: Numbers (including decimals).
- `validname`: Alphanumeric + underscore/dash.
- `en`: English letters only.
- `domain`: Domain names.
- `mobile`: Phone numbers (basic loose check).
- `pwd`: Minimum 8 chars, 1 letter, 1 number.
- `strongpwd`: Min 8 chars, letter, number, special char.

## Events

The Valid component dispatches and listens to specific events for flow control.

| Event                   | Trigger Timing                                                  | Callback Parameters                                     |
| :---------------------- | :-------------------------------------------------------------- | :------------------------------------------------------ |
| `uno:validation-failed` | Dispatched on `<form>` when validation fails during submission. | `detail: { errors: Object }`                            |
| `htmx:confirm`          | Intercepted on `document.body`.                                 | Prevents request if form is invalid.                    |
| `htmx:afterRequest`     | Intercepted on `document.body`.                                 | Parses JSON response. If status is 422, applies errors. |

## Methods

Accessed via `Alpine.store('uno').valid`.

### `addRules(fieldName, rules)`

Adds validation rules to a field programmatically.

- **Parameters**:
  - `fieldName` (String): Field name or ID.
  - `rules` (Object): Configuration object.
- **Returns**: `void`

### `validate(fieldName, value)`

Validates a single field.

- **Parameters**:
  - `fieldName` (String): Field identifier.
  - `value` (Any): Current value.
- **Returns**: `Boolean` - `true` if valid.

### `validateForm(formElement)`

Validates all registered fields within a form.

- **Parameters**:
  - `formElement` (HTMLFormElement): The form.
- **Returns**: `Boolean` - `true` if form is valid.

### `registerValidator(name, fn)`

Registers a global custom validator.

- **Parameters**:
  - `name` (String): Validator name.
  - `fn` (Function): `(value, fieldName) => true | String | {valid: false, message: ''}`
- **Returns**: `void`
- **Example**:
  ```javascript
  Alpine.store('uno').valid.registerValidator('checkUser', (val) => {
    return val !== 'admin' || 'Admin is reserved'
  })
  ```

### `showErrors(errors)`

Manually sets error messages (e.g., from server response).

- **Parameters**:
  - `errors` (Object): `{ fieldName: 'Error Message' }`
- **Returns**: `void`

### `clearErrors()`

Clears all current error messages.

- **Parameters**: None.
- **Returns**: `void`

## Examples

### Password Confirmation

```html
<form x-uno-valid>
  <div>
    <label>Password</label>
    <input type="password" name="password" x-uno-valid="{ type: 'pwd' }" />
  </div>

  <div>
    <label>Confirm Password</label>
    <!-- Use 'match' to reference the first password field -->
    <input type="password" name="confirm" x-uno-valid="{ match: 'password' }" />
    <span x-text="$store.uno.valid.errors['confirm']"></span>
  </div>
</form>
```

### Custom Async Validator (Server Check)

While client-side custom validators are synchronous, you can combine them with HTMX for server checks.

```html
<!-- HTMX checks availability on blur -->
<input
  name="username"
  x-uno-valid="{ required: true }"
  hx-post="/check-username"
  hx-trigger="blur"
/>
<div id="username-error"></div>
```

_Note: Pure client-side custom validators are registered via `registerValidator`._

## Version History

| Version | Date       | Changes                                                                       |
| :------ | :--------- | :---------------------------------------------------------------------------- |
| v1.0.0  | 2025-01-09 | Initial release with HTML5 support, Custom Validators, and HTMX 422 handling. |
