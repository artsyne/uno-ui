/**
 * Uno - Logger
 *
 * @module uno.core.log
 * @version 0.1.0
 * @license MIT
 */

const DEBUG = true;

/**
 * Internal logger
 */
export const log = {
  /**
   * Log info message
   * @param {...any} args - Arguments to log
   */
  info(...args) {
    if (DEBUG) console.log('[Uno Info]', ...args);
  },
  /**
   * Log warning message
   * @param {...any} args - Arguments to log
   */
  warn(...args) {
    console.warn('[Uno Warn]', ...args);
  },
  /**
   * Log error message
   * @param {...any} args - Arguments to log
   */
  error(...args) {
    console.error('[Uno Error]', ...args);
  }
};
