/**
 * Uno - Error Boundary
 *
 * Handles global error capturing and reporting.
 *
 * @module uno.core.error
 */

import { log } from '../utils/log'

/**
 * ErrorBoundary class for handling global errors and promise rejections.
 */
export class ErrorBoundary {
  /**
   * Create an ErrorBoundary instance.
   *
   * @param {Object} [options={}] - Configuration options
   * @param {string} [options.reportUrl] - URL to send error reports to
   * @param {Object} [options.context] - Additional context to send with reports
   */
  constructor(options = {}) {
    this.options = options
    this.handlers = []
    this.initialized = false
  }

  /**
   * Initialize the error boundary.
   * Sets up global event listeners for error and unhandledrejection.
   * Connects to the Uno logger.
   */
  init() {
    if (this.initialized) return
    this.initialized = true

    if (typeof window !== 'undefined') {
      // Global Error
      window.addEventListener('error', (event) => {
        this.handle(event.error || new Error(event.message), { type: 'uncaught', event })
      })

      // Promise Rejection
      window.addEventListener('unhandledrejection', (event) => {
        this.handle(event.reason || new Error('Unhandled Rejection'), {
          type: 'unhandledrejection',
          event
        })
      })
    }

    // Connect to Uno logger
    log.setReporter((...args) => {
      // Find the first Error object, or use the first argument
      const error = args.find((a) => a instanceof Error) || args[0]
      this.handle(error, { type: 'logger', args })
    })
  }

  /**
   * Handle an error.
   * Normalizes the error, notifies handlers, and reports if configured.
   *
   * @param {Error|any} error - The error to handle
   * @param {Object} [context={}] - Contextual information about the error
   */
  handle(error, context = {}) {
    // Normalize error
    const err = error instanceof Error ? error : new Error(String(error))

    // Notify handlers
    this.handlers.forEach((h) => {
      try {
        h(err, context)
      } catch (e) {
        console.error('[Uno ErrorBoundary] Handler failed', e)
      }
    })

    // Auto-report if configured
    if (this.options.reportUrl) {
      this.report(err, context)
    }
  }

  /**
   * Report an error to the configured URL.
   *
   * @param {Error} error - The error to report
   * @param {Object} context - Contextual information
   */
  report(error, context) {
    if (!this.options.reportUrl) return

    const payload = {
      message: error.message,
      stack: error.stack,
      context,
      timestamp: new Date().toISOString(),
      url: typeof window !== 'undefined' ? window.location.href : '',
      userAgent: typeof navigator !== 'undefined' ? navigator.userAgent : ''
    }

    try {
      if (typeof navigator !== 'undefined' && navigator.sendBeacon) {
        const blob = new Blob([JSON.stringify(payload)], { type: 'application/json' })
        navigator.sendBeacon(this.options.reportUrl, blob)
      } else if (typeof fetch !== 'undefined') {
        fetch(this.options.reportUrl, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(payload),
          keepalive: true
        }).catch(() => {
          // Ignore reporting failures to avoid loops
        })
      }
    } catch {
      // Ignore reporting errors
    }
  }

  /**
   * Register an error handler.
   *
   * @param {Function} fn - The handler function
   */
  onError(fn) {
    if (typeof fn === 'function') {
      this.handlers.push(fn)
    }
  }
}
