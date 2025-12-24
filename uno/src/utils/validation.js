/**
 * Uno - Validation Utils
 *
 * @module uno.utils.validation
 */

export const VALIDATORS = {
  // Common Regex Patterns
  email: /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/,
  url: /^(https?:\/\/)?([\da-z.-]+)\.([a-z.]{2,6})([/\w .-]*)*\/?$/,
  username: /^[a-zA-Z0-9_-]{3,16}$/,
  number: /^\d+$/,
  integer: /^-?\d+$/,
  date: /^\d{4}-\d{2}-\d{2}$/,
  time: /^\d{2}:\d{2}$/,
  phone: /^\+?[\d\s-]{10,}$/,
  alphanumeric: /^[a-zA-Z0-9]+$/,
}

export function validate(value, type) {
  if (!type || !VALIDATORS[type]) return true
  return VALIDATORS[type].test(value)
}
