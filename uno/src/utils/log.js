/**
 * Uno - Logger
 *
 * @module uno.core.log
 * @version 0.2.0
 * @license MIT
 */

let isDebug = true
let errorReporter = null

/**
 * Internal logger
 */
export const log = {
  /**
   * Set debug mode
   * @param {boolean} debug
   */
  setDebug(debug) {
    isDebug = debug
  },

  /**
   * Set custom error reporter
   * @param {Function} reporter - (error, ...args) => void
   */
  setReporter(reporter) {
    errorReporter = reporter
  },

  /**
   * Log info message
   * @param {...any} args - Arguments to log
   */
  info(...args) {
    if (isDebug) console.log('[Uno Info]', ...args)
  },
  /**
   * Log warning message
   * @param {...any} args - Arguments to log
   */
  warn(...args) {
    console.warn('[Uno Warn]', ...args)
  },
  /**
   * Log error message
   * @param {...any} args - Arguments to log
   */
  error(...args) {
    console.error('[Uno Error]', ...args)
    if (errorReporter) {
      try {
        errorReporter(...args)
      } catch (e) {
        console.error('[Uno Error] Reporter failed:', e)
      }
    }
  }
}
