/**
 * Uno - Validation Utils
 *
 * Common validation patterns and utilities.
 *
 * @module uno.utils.validation
 */

/**
 * Uno - Validation Utils
 *
 * Common validation patterns and utilities for form input validation.
 * Provides regex-based validators for common data types and formats.
 *
 * @module uno.utils.validation
 * @version 0.6.0
 */

/**
 * Pre-defined validation patterns for common input types.
 *
 * Pattern Categories:
 *
 * **Standard Types:**
 * - `email`: RFC 5322 compliant email address validation
 *   - Allows letters, numbers, and special chars before @
 *   - Requires domain with at least one dot
 *   - Example: user.name+tag@example.co.uk
 *
 * - `url`: HTTP/HTTPS URL validation
 *   - Optional protocol (http:// or https://)
 *   - Domain with TLD (2-6 characters)
 *   - Optional path and query parameters
 *   - Example: https://example.com/path?query=value
 *
 * - `date`: Multiple date format support
 *   - ISO format: YYYY-MM-DD (2023-12-31)
 *   - US format: MM/DD/YYYY (12/31/2023)
 *   - EU format: DD.MM.YYYY (31.12.2023)
 *
 * **Numeric:**
 * - `int`: Integer numbers (positive or negative)
 *   - Examples: 42, -10, 0
 *   - No decimal points allowed
 *
 * - `number`: Decimal numbers (floating point)
 *   - Supports optional decimal places
 *   - Examples: 3.14, -2.5, 42, 0.001
 *   - Optional leading minus sign
 *
 * **String Content:**
 * - `validname`: Alphanumeric usernames/identifiers
 *   - Letters, numbers, hyphens, underscores only
 *   - Examples: user_name, user-123, Username
 *
 * - `en`: English letters only (no numbers or special chars)
 *   - Both uppercase and lowercase
 *   - Examples: Hello, ABC, test
 *
 * - `domain`: Valid domain names (DNS format)
 *   - Multiple labels separated by dots
 *   - No leading/trailing hyphens in labels
 *   - TLD: 2-6 characters
 *   - Examples: example.com, sub.example.co.uk
 *
 * - `mobile`: International phone numbers
 *   - Optional leading + for country code
 *   - Digits, spaces, and hyphens allowed
 *   - Minimum 7 characters
 *   - Examples: +1-555-1234, 1234567890, +44 20 1234 5678
 *
 * **Security:**
 * - `pwd`: Basic password (minimum security)
 *   - At least one letter (uppercase or lowercase)
 *   - At least one number
 *   - 8-72 characters; upper bound aligns with the backend pkg.IsPwd rule
 *     (bcrypt silently truncates input beyond 72 bytes)
 *   - Allowed special characters: _!@#$%^&*()+-=
 *   - Example: Password1
 *
 * - `strongpwd`: Strong password (high security)
 *   - At least one letter
 *   - At least one number
 *   - At least one special character (@$!%*#?&)
 *   - Minimum 8 characters total
 *   - Example: P@ssw0rd!
 *
 * Usage Notes:
 * - All patterns are case-insensitive where applicable
 * - Patterns validate the entire string (implicit ^ and $)
 * - Empty strings will fail validation
 * - Whitespace is not trimmed automatically
 *
 * @type {Object.<string, RegExp>}
 * @readonly
 *
 * @example
 * // Validate email
 * VALIDATORS.email.test('user@example.com') // true
 * VALIDATORS.email.test('invalid.email') // false
 *
 * @example
 * // Validate password
 * VALIDATORS.pwd.test('Abc12345') // true
 * VALIDATORS.strongpwd.test('Abc@1234') // true
 * VALIDATORS.strongpwd.test('weakpass') // false
 */
export const VALIDATORS = {
  // --- Standard Types ---
  email: /^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9-]+(?:\.[a-zA-Z0-9-]+)+$/,
  url: /^(https?:\/\/)?([\da-z.-]+)\.([a-z.]{2,6})([/\w .-]*)*\/?$/i,
  date: /^(\d{4}-\d{2}-\d{2}|\d{2}\/\d{2}\/\d{4}|\d{2}\.\d{2}\.\d{4})$/,

  // --- Numeric ---
  int: /^-?\d+$/,
  number: /^-?\d*(\.\d+)?$/, // Added floating point support for general 'number'

  // --- String Content ---
  validname: /^[a-zA-Z0-9_-]{3,16}$/,
  en: /^[a-zA-Z]+$/,
  domain: /^((?!-)[A-Za-z0-9-]{1,63}(?<!-)\.)+[A-Za-z]{2,6}$/,
  mobile: /^\+?[\d\s-]{7,}$/,

  // --- Network & Colors (New) ---
  ipv4: /^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/,
  hex: /^#?([a-f0-9]{6}|[a-f0-9]{3})$/i,

  // --- Security ---
  pwd: /^(?=.*[A-Za-z])(?=.*\d)[A-Za-z\d!_@#$%^&*()+\-=]{8,72}$/, // Password must be 8-72 chars, contain letters and numbers, and can include _!@#$%^&*()+-= special characters
  strongpwd: /^(?=.*[A-Za-z])(?=.*\d)(?=.*[@$!%*#?&])[A-Za-z\d@$!%*#?&]{8,80}$/ // Password must be 8-80 chars, contain letters, numbers, and at least 1 special char (@$!%*#?&)
}

/**
 * Validate a value against a preset validation type.
 *
 * This function provides a convenient wrapper around the VALIDATORS regex patterns,
 * allowing validation by type name rather than direct regex testing.
 *
 * Validation Process:
 * 1. Checks if the type exists in VALIDATORS
 * 2. If type not found or empty, returns true (no validation)
 * 3. If type found, tests the value against the pattern
 * 4. Returns boolean result
 *
 * Type Safety:
 * - Non-existent types are treated as "no validation" (returns true)
 * - This is intentional to prevent validation errors breaking forms
 * - Alternative: Throw error for unknown types (breaking change)
 *
 * Edge Cases:
 * - Empty type: Returns true (no validation)
 * - Unknown type: Returns true (graceful degradation)
 * - Empty value: Depends on pattern (most patterns reject empty strings)
 * - null/undefined value: Coerced to string, likely fails validation
 *
 * Performance:
 * - O(1) lookup in VALIDATORS object
 * - Regex test performance varies by pattern complexity
 * - Suitable for real-time validation (keystroke events)
 *
 * @param {string} value - The value to validate
 * @param {string} type - The validation type (key from VALIDATORS object)
 * @returns {boolean} True if valid or type not found, false if invalid
 *
 * @example
 * // Validate email
 * validate('user@example.com', 'email') // true
 * validate('invalid', 'email') // false
 *
 * @example
 * // Validate number
 * validate('3.14', 'number') // true
 * validate('abc', 'number') // false
 *
 * @example
 * // Unknown type (no validation)
 * validate('anything', 'unknown-type') // true
 *
 * @example
 * // Empty type (no validation)
 * validate('anything', '') // true
 *
 * @see {@link VALIDATORS} for available validation types
 */
export function validate(value, type) {
  if (!type || !VALIDATORS[type]) return true
  return VALIDATORS[type].test(value)
}
