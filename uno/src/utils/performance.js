/**
 * Uno - Performance Utils
 *
 * @module uno.utils.performance
 */

/**
 * Debounce a function
 * @param {Function} func - Function to debounce
 * @param {number} wait - Wait time in ms
 * @param {boolean} immediate - Trigger immediately
 * @returns {Function}
 */
export function debounce(func, wait, immediate) {
  let timeout
  return function () {
    const context = this
    const args = arguments
    const later = function () {
      timeout = null
      if (!immediate) func.apply(context, args)
    }
    const callNow = immediate && !timeout
    clearTimeout(timeout)
    timeout = setTimeout(later, wait)
    if (callNow) func.apply(context, args)
  }
}

/**
 * Throttle a function
 * @param {Function} func - Function to throttle
 * @param {number} limit - Limit time in ms
 * @returns {Function}
 */
export function throttle(func, limit) {
  let inThrottle
  return function () {
    const args = arguments
    const context = this
    if (!inThrottle) {
      func.apply(context, args)
      inThrottle = true
      setTimeout(() => (inThrottle = false), limit)
    }
  }
}

/**
 * Simple performance measure wrapper
 * @param {string} name
 * @param {Function} fn
 */
export function measure(name, fn) {
  if (typeof performance === 'undefined') return fn()
  const start = performance.now()
  try {
    return fn()
  } finally {
    const end = performance.now()
    console.debug(`[Uno Perf] ${name}: ${(end - start).toFixed(2)}ms`)
  }
}
