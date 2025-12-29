/**
 * Uno - Error Boundary
 *
 * Handles global error capturing and reporting.
 *
 * @module uno.core.error
 */

import { log } from '../utils/log'

export class ErrorBoundary {
  constructor(options = {}) {
    this.options = options
    this.handlers = []
    this.initialized = false
  }

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
        this.handle(event.reason || new Error('Unhandled Rejection'), { type: 'unhandledrejection', event })
      })
    }
    
    // Connect to Uno logger
    log.setReporter((...args) => {
        // Find the first Error object, or use the first argument
        const error = args.find(a => a instanceof Error) || args[0]
        this.handle(error, { type: 'logger', args })
    })
  }

  handle(error, context = {}) {
    // Normalize error
    const err = error instanceof Error ? error : new Error(String(error))
    
    // Notify handlers
    this.handlers.forEach(h => {
        try {
            h(err, context)
        } catch (e) {
            console.error('[Uno ErrorBoundary] Handler failed', e)
        }
    })
    
    // Auto-report if configured
    if (this.options.reportUrl) {
        // TODO: Implement reporting logic (beacon/fetch)
    }
  }

  onError(fn) {
    if (typeof fn === 'function') {
        this.handlers.push(fn)
    }
  }
}
