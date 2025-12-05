/**
 * Uno - Body Scroll Lock
 *
 * @module uno.core.bodyScroll
 * @version 0.1.0
 * @license MIT
 */

let count = 0
let original = ''

/**
 * Lock body scroll
 */
export function lock() {
  if (count === 0) {
    original = document.body.style.overflow || ''
    document.body.style.overflow = 'hidden'
  }
  count++
}

/**
 * Unlock body scroll
 */
export function unlock() {
  count = Math.max(0, count - 1)
  if (count === 0) {
    document.body.style.overflow = original || ''
    original = ''
  }
}

/**
 * Force reset lock count and restore scroll
 */
export function reset() {
  count = 0
  document.body.style.overflow = original || ''
  original = ''
}

/**
 * Get current lock count
 * @returns {number} Lock count
 */
export function getCount() {
  return count
}
