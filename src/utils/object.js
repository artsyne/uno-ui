/**
 * Uno - Object Utils
 *
 * Utility functions for object manipulation.
 *
 * @module uno.utils.object
 */

/**
 * Check if a value is a plain object.
 *
 * @param {*} value - The value to check
 * @returns {boolean} True if the value is a plain object
 */
function isPlainObject(value) {
  return Object.prototype.toString.call(value) === '[object Object]'
}

/**
 * Deep merge two objects.
 * Supports special syntax for string concatenation (e.g. "+class").
 *
 * @param {Object} base - The base object
 * @param {Object} patch - The patch object to merge into base
 * @returns {Object} The merged object
 */
export function mergeDeep(base, patch) {
  if (!isPlainObject(base)) return isPlainObject(patch) ? { ...patch } : patch
  if (!isPlainObject(patch)) return patch === undefined ? { ...base } : patch

  const out = { ...base }
  for (const [key, value] of Object.entries(patch)) {
    if (isPlainObject(value) && isPlainObject(out[key])) {
      out[key] = mergeDeep(out[key], value)
      continue
    }
    if (isPlainObject(value)) {
      out[key] = mergeDeep({}, value)
      continue
    }
    // Handle string append (e.g. "+class")
    if (
      typeof value === 'string' &&
      value.trim().indexOf('+') === 0 &&
      typeof out[key] === 'string'
    ) {
      out[key] = `${out[key]} ${value.trim().slice(1)}`.trim()
      continue
    }
    out[key] = value
  }
  return out
}
