# Valid Component (Form Validation)

A lightweight form validation plugin based on Alpine.js + HTMX + TailwindCSS, providing declarative and programmatic configuration, real-time validation, async validation, and i18n support.

## Features

- **Diverse Validation Rules**: Supports required, regex, length limits, numeric ranges, custom functions, etc.
- **Built-in Validators**: Pre-built rules for email, URL, mobile phone, password strength, etc.
- **Feedback Rendering**: Automatically handles validation state styles (`is-valid`/`is-invalid`) and error messages.
- **Form-Level Validation**: Unified validation before submission to prevent illegal submissions.
- **HTMX Integration**: Allows HTMX requests only after validation passes, supporting `hx-post`/`hx-target`.
- **Internationalization**: Supports multi-language error messages.

## Quick Start

### Declarative Configuration (HTML Attributes)

The simplest way is to use `data-valid-*` attributes directly on the Input element.

```html
<form x-data="validForm()">
  <div class="mb-4">
    <label class="block mb-2">Username</label>
    <!-- Required, min length 3, use built-in validname rule -->
    <input 
      type="text" 
      name="username" 
      class="border p-2 w-full"
      data-valid-required="true" 
      data-valid-min="3" 
      data-valid-type="validname" 
    />
    <!-- Error container -->
    <div class="valid-feedback text-red-500 text-sm mt-1"></div>
  </div>

  <div class="mb-4">
    <label class="block mb-2">Email</label>
    <input 
      type="email" 
      name="email" 
      class="border p-2 w-full"
      data-valid-required="true" 
    />
    <div class="valid-feedback text-red-500 text-sm mt-1"></div>
  </div>

  <button type="submit" class="bg-blue-500 text-white px-4 py-2">Submit</button>
</form>
```

### JSON Rule Configuration (`valid` attribute)

For complex rules (such as combined rules, custom error messages), use the `valid` attribute to pass a JSON array.

```html
<input 
  type="password" 
  name="password" 
  valid='[
    {"type":"required", "msg":"Password cannot be empty"},
    {"type":"limit", "min":6, "max":20, "msg":"Length must be between 6-20 characters"},
    {"type":"strongpwd", "msg":"Password must contain letters, numbers, and special characters"}
  ]'
/>
```

## Built-in Validators

The plugin has built-in common validation rules, which can be used via `data-valid-type` or specifying `type` in JSON rules:

| Validator Name | Description | Rule/Regex |
|----------------|-------------|------------|
| `validname` | Valid Name | Allows only letters, numbers, underscores, hyphens `^[a-zA-Z0-9_-]+$` |
| `en` | English Only | Allows only letters `^[a-zA-Z]+$` |
| `domain` | Domain | Domain name format |
| `mobile` | Mobile | International format, supports `+` prefix |
| `int` | Integer | Allows negative numbers |
| `email` | Email | Standard email format |
| `date` | Date | `YYYY-MM-DD`, `MM/DD/YYYY`, `DD.MM.YYYY` |
| `url` | URL | HTTP/HTTPS link |
| `pwd` | Simple Password | At least letters and numbers, 8+ characters |
| `strongpwd`| Strong Password | Letters, numbers, special characters, 8+ characters |

## Security: ReDoS Protection

To prevent Regular Expression Denial of Service (ReDoS) attacks:
- **Length Limit**: Regex patterns are limited to 1000 characters.
- **Safe Execution**: Pattern execution is wrapped in try-catch blocks.
- **Complexity**: Avoid nested quantifiers in custom regex rules.

## Rule Details (JSON Config)

Rule objects defined in the `valid` attribute support the following fields:

### General Rules

- `required`: Required field.
  ```json
  {"type": "required", "msg": "This field is required"}
  ```
- `emptyquit`: Passes if empty (subsequent rules are not executed). Often used for "optional, but must match format if filled".
  ```json
  {"type": "emptyquit"}
  ```

### Format & Length

- `limit`: Length limit.
  ```json
  {"type": "limit", "min": 5, "max": 10}
  ```
- `range`: Numeric range.
  ```json
  {"type": "range", "min": 1, "max": 100}
  ```
- `regexp`: Custom regex.
  ```json
  {"type": "regexp", "expression": "^[A-Z]+$", "msg": "Must be uppercase letters"}
  ```
- `tag`: Tag count limit (calculated by separator).
  ```json
  {"type": "tag", "mark": ",", "max": 5}
  ```
- `compareto`: Value comparison (e.g., confirm password).
  ```json
  {"type": "compareto", "destid": "password_input_id", "msg": "Inputs do not match"}
  ```

### Async Validation (Callback)

Supports async validation logic (e.g., checking if username exists).

1. **Register Callback**:
   Register validation function in Alpine store.
   ```javascript
   Alpine.store('valid').registerCallback('checkUser', async (id, value) => {
     const res = await fetch(`/api/check-user?name=${value}`);
     const data = await res.json();
     // Return true/false, or return error message string, or return { valid: false, message: '...' }
     return data.exists ? 'Username already exists' : true;
   });
   ```

2. **Apply Rule**:
   ```html
   <input valid='[{"type":"callback", "fn":"checkUser"}]' ... />
   ```

## Internationalization (i18n)

The component relies on `Alpine.store('i18n')` for multi-language support. If not configured, default English messages are used.

### Configure Language Pack

```javascript
Alpine.store('i18n', {
  locale: 'zh-CN',
  messages: {
    'zh-CN': {
      valid: {
        required: 'This field is required',
        pattern: 'Invalid format',
        min: 'At least {n} characters',
        email: 'Please enter a valid email',
        mobile: 'Please enter a valid mobile number',
        // ...other built-in validator keys match validator names
      }
    }
  },
  t(key, vars) {
    // Simple translation implementation example
    let msg = this.messages[this.locale]?.valid?.[key.replace('valid.', '')] || key;
    if (vars) {
      Object.keys(vars).forEach(k => msg = msg.replace(`{${k}}`, vars[k]));
    }
    return msg;
  }
});
```

## Events & API

### Events
- **Field Input/Blur**: Validation triggered on input or blur.
- **Form Submit**: Full validation triggered on form submission.

### Class States
The component automatically toggles the following class names on the Input element:
- `is-valid`: Validation passed
- `is-invalid`: Validation failed
- `validating`: Async validation in progress

You can define styles for these class names in CSS, or rely on default Tailwind styles applied.
