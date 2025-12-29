import { describe, it, expect, vi, beforeEach } from 'vitest'
import { ErrorBoundary } from '../../src/core/error'
import { log } from '../../src/utils/log'

describe('Error Boundary', () => {
  let boundary

  beforeEach(() => {
    boundary = new ErrorBoundary()
    // Reset log reporter
    log.setReporter(null)

    // Mock PromiseRejectionEvent if not available
    if (typeof PromiseRejectionEvent === 'undefined') {
      global.PromiseRejectionEvent = class PromiseRejectionEvent extends Event {
        constructor(type, options) {
          super(type, options)
          this.promise = options.promise
          this.reason = options.reason
        }
      }
    }
  })

  it('should capture global errors', () => {
    boundary.init()
    const handler = vi.fn()
    boundary.onError(handler)

    // Simulate error
    const error = new Error('Test Error')
    // ErrorEvent constructor might not be fully supported in happy-dom or requires specific args
    // happy-dom supports ErrorEvent
    const event = new ErrorEvent('error', { error, message: 'Test Error' })
    window.dispatchEvent(event)

    expect(handler).toHaveBeenCalledWith(error, expect.objectContaining({ type: 'uncaught' }))
  })

  it('should capture unhandled rejections', () => {
    boundary.init()
    const handler = vi.fn()
    boundary.onError(handler)

    const error = new Error('Async Error')
    // PromiseRejectionEvent
    const event = new PromiseRejectionEvent('unhandledrejection', {
      promise: Promise.resolve(),
      reason: error
    })
    window.dispatchEvent(event)

    expect(handler).toHaveBeenCalledWith(
      error,
      expect.objectContaining({ type: 'unhandledrejection' })
    )
  })

  it('should capture logger errors', () => {
    boundary.init()
    const handler = vi.fn()
    boundary.onError(handler)

    const error = new Error('Log Error')
    log.error(error)

    expect(handler).toHaveBeenCalledWith(error, expect.objectContaining({ type: 'logger' }))
  })
})
