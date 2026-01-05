/**
 * Uno - Validation Utils
 *
 * Common validation patterns and utilities.
 *
 * @module uno.utils.validation
 */

export const VALIDATORS = {
  // --- Standard Types ---
  email: /^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9-]+(?:\.[a-zA-Z0-9-]+)*$/,
  url: /^(https?:\/\/)?([\da-z.-]+)\.([a-z.]{2,6})([/\w .-]*)*\/?$/i,
  date: /^(\d{4}-\d{2}-\d{2}|\d{2}\/\d{2}\/\d{4}|\d{2}\.\d{2}\.\d{4})$/,

  // --- Numeric ---
  int: /^-?\d+$/,
  number: /^-?\d*(\.\d+)?$/, // Added floating point support for general 'number'

  // --- String Content ---
  validname: /^[a-zA-Z0-9_-]+$/,
  en: /^[a-zA-Z]+$/,
  domain: /^((?!-)[A-Za-z0-9-]{1,63}(?<!-)\.)+[A-Za-z]{2,6}$/,
  mobile: /^\+?[\d\s-]{7,}$/,

  // --- Security ---
  pwd: /^(?=.*[A-Za-z])(?=.*\d)[A-Za-z\d]{8,}$/, // Letters + Numbers, 8+
  strongpwd: /^(?=.*[A-Za-z])(?=.*\d)(?=.*[@$!%*#?&])[A-Za-z\d@$!%*#?&]{8,}$/ // + Special Char
}

/**
 * Validate a value against a preset type
 *
 * @param {string} value
 * @param {string} type
 * @returns {boolean}
 */
export function validate(value, type) {
  if (!type || !VALIDATORS[type]) return true
  return VALIDATORS[type].test(value)
}
