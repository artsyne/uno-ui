# Valid Component (Form Validation)

A robust, Store-First form validation plugin for Alpine.js and HTMX. It leverages standard HTML5 attributes for configuration and integrates deeply with HTMX events for a seamless "SPA-like" form experience.

## Features

- **Store-First Architecture**: All validation state lives in `Alpine.store('uno').valid`.
- **Directive-Driven**: Zero config boilerplate. Just add `x-uno-valid` and standard HTML attributes.
- **HTML5 Native**: Uses standard `required`, `pattern`, `minlength`, `type="email"` attributes.
- **HTMX Integration**: Automatically intercepts `htmx:confirm` to prevent invalid requests and handles 422 server errors.
- **Declarative Errors**: Reactive error display via Alpine store.
- **I18n Support**: Integrated with `Uno` I18n system.

## Quick Start

### 1. Basic Setup

Enable validation on your form and inputs:

```html
<form x-uno-valid hx-post="/api/submit" novalidate>
  <!-- Input Group -->
  <div>
    <label>Username</label>
    <input
      x-uno-valid
      name="username"
      required
      pattern="^[a-zA-Z0-9]+$"
      class="border p-2"
      :class="$store.uno.valid.errors['username'] && 'border-red-500'"
    />
    <!-- Error Message -->
    <div x-text="$store.uno.valid.errors['username']" class="text-red-500 text-sm"></div>
  </div>

  <button type="submit">Submit</button>
</form>
```

### 2. Supported Attributes

The plugin reads configuration directly from standard HTML5 attributes:

| Attribute                 | Type      | Description                                                                                                        |
| ------------------------- | --------- | ------------------------------------------------------------------------------------------------------------------ |
| `x-uno-valid`             | Directive | **Required**. Activates validation. Can pass JSON config override: `x-uno-valid="{ realtime: false }"`             |
| `name` or `id`            | String    | **Required**. Unique identifier for the field in the store.                                                        |
| `required`                | Boolean   | Native required check.                                                                                             |
| `pattern`                 | Regex     | Native regex pattern check.                                                                                        |
| `type="email"`            | String    | Validates email format.                                                                                            |
| `type="url"`              | String    | Validates URL format.                                                                                              |
| `minlength` / `maxlength` | Number    | String length limits.                                                                                              |
| `min` / `max`             | Number    | Numeric range limits (requires `type="number"` or numeric content).                                                |

### 3. Custom Error Messages & Advanced Config

Use the `x-uno-valid` configuration object to provide custom error messages or advanced rules like `match`.

```html
<input
  x-uno-valid="{
    msgRequired: 'Please fill this out!',
    msgPattern: 'Numbers only please',
    match: 'password_confirm',
    msgMatch: 'Passwords do not match'
  }"
  name="password"
  required
  pattern="\d+"
/>
```

**Available Config Properties:**

- `realtime`: (Boolean) Default `true`. Set `false` to validate only on submit/blur.
- `match`: (Selector) Validates value matches another field (e.g., `match: 'password'`).
- `msgRequired`: Custom message for required error.
- `msgEmail`, `msgUrl`: Custom message for type errors.
- `msgPattern`: Custom message for pattern errors.
- `msgMin`, `msgMax`: Custom message for range errors.
- `msgMinLength`, `msgMaxLength`: Custom message for length errors.
- `msgMatch`: Custom message for match errors.

## HTMX Integration

### Client-Side Interception (`htmx:confirm`)

When a form with `x-uno-valid` triggers an HTMX request:

1. The plugin intercepts `htmx:confirm`.
2. Runs validation on all inputs.
3. If invalid: `preventDefault()` stops the request, focuses the first error field, and dispatches `uno:validation-failed`.
4. If valid: Allows the request to proceed.

### Server-Side Errors (`htmx:afterRequest`)

If the server returns a **422 Unprocessable Entity** status with a JSON body:

```json
{
  "errors": {
    "username": "Username already taken",
    "email": "Invalid domain"
  }
}
```

The plugin automatically populates `$store.uno.valid.errors` with these messages.

## Programmatic API

Access the internal store via `Alpine.store('uno').valid`:

```javascript
// Add rules programmatically
Alpine.store('uno').valid.addRules('username', { required: true, minLength: 5 })

// Validate a single field
const isValid = Alpine.store('uno').valid.validate('username', 'value')

// Validate entire form
const form = document.getElementById('myForm')
const isFormValid = Alpine.store('uno').valid.validateForm(form)

// Clear errors
Alpine.store('uno').valid.clearErrors()

// Show errors manually
Alpine.store('uno').valid.showError('username', 'Invalid value')
```

## Configuration

Global configuration can be set via `Uno.registerPluginDefaults` or global options:

```javascript
Uno.registerPluginDefaults('valid', {
  realtime: true,
  scrollOnError: true,
  showError: true // Auto-toggle error classes (if implemented)
})
```
