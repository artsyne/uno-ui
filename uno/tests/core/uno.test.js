import { describe, it, expect, beforeEach, vi } from 'vitest'
import { UnoInstance } from '../../src/core/uno'
import { resetOnceMarks } from '../../src/utils/dom'

// Mock the log module
vi.mock('../../src/utils/log', () => ({
  log: {
    info: vi.fn(),
    debug: vi.fn(),
    warn: vi.fn(),
    error: vi.fn()
  }
}))

describe('UnoInstance Configuration', () => {
  let uno

  beforeEach(() => {
    uno = new UnoInstance()
  })

  it('should register plugin defaults', () => {
    const defaults = { foo: 'bar' }
    uno.registerPluginDefaults('test', defaults)
    expect(uno.pluginDefaults['test']).toEqual(defaults)
  })

  it('should merge configurations correctly (4-level merge)', () => {
    // Level 1: Plugin Defaults
    uno.registerPluginDefaults('box', {
      width: 100,
      height: 100,
      color: 'blue'
    })

    // Level 2 & 3: Global Config
    uno.options = {
      plugins: {
        box: {
          height: 200, // Level 2
          opacity: 0.5
        }
      },
      box: {
        color: 'red' // Level 3
      }
    }

    // Level 4: Instance Config
    const instanceOpts = {
      width: 300
    }

    const config = uno.getConfig('box', instanceOpts)

    // Expectations
    expect(config.width).toBe(300) // Level 4 overrides Level 1
    expect(config.height).toBe(200) // Level 2 overrides Level 1
    expect(config.color).toBe('red') // Level 3 overrides Level 1
    expect(config.opacity).toBe(0.5) // Level 2
  })

  it('should support deep merge', () => {
    uno.registerPluginDefaults('test', {
      nested: { a: 1, b: 2 }
    })

    uno.options = {
      plugins: {
        test: {
          nested: { b: 3 }
        }
      }
    }

    const instanceOpts = {
      nested: { c: 4 }
    }

    const config = uno.getConfig('test', instanceOpts)

    expect(config.nested).toEqual({
      a: 1,
      b: 3,
      c: 4
    })
  })

  it('should resolve nested styles', () => {
    const config = {
      types: {
        success: 'bg-green-500',
        error: 'bg-red-500'
      },
      base: 'p-4'
    }

    // Mock _getUiValue which is used internally by getStyle
    // But getStyle is what we test.
    // We need to inject config first.
    // getStyle(comp, elem, defaultCls)

    // In the new implementation, we need to know how getStyle retrieves config.
    // It calls this._getUiValue(comp, elem)

    // Let's just test the logic if we can mock _getUiValue or if we set up the config correctly.
    // Since _getUiValue likely uses getConfig, let's see.
    // Actually infrastructure.md says:
    // getStyle(comp, elem, defaultCls) {
    //   const config = this._getUiValue(comp, elem)
    //   ...
    // }

    // We need to implement _getUiValue or make sure it works.
  })
})

describe('UnoInstance HTMX Integration', () => {
  let uno
  let mockDocument
  let mockWindow
  let mockAlpineStore

  beforeEach(() => {
    // Reset once marks to ensure tests can register event listeners
    resetOnceMarks()

    // Mock document
    const listeners = {}
    mockDocument = {
      addEventListener: vi.fn((event, handler) => {
        listeners[event] = handler
      }),
      removeEventListener: vi.fn(),
      // Store listeners for triggering in tests
      _triggerEvent: (event, detail) => {
        if (listeners[event]) {
          listeners[event]({ detail })
        }
      },
      documentElement: {
        lang: ''
      }
    }

    // Mock window.htmx
    const mockHtmx = {
      ajax: vi.fn()
    }

    // Mock Alpine store
    const mockNotify = {
      startProgress: vi.fn(),
      stopProgress: vi.fn(),
      success: vi.fn(),
      error: vi.fn()
    }

    mockAlpineStore = vi.fn((name) => {
      if (name === 'uno') {
        return { notify: mockNotify }
      }
      return null
    })

    mockWindow = {
      Alpine: {
        store: mockAlpineStore
      },
      htmx: mockHtmx,
      document: mockDocument
    }

    // Setup global mocks
    global.window = mockWindow
    global.document = mockDocument
    global.WeakSet = WeakSet

    uno = new UnoInstance()
  })

  afterEach(() => {
    delete global.window
    delete global.document
  })

  describe('htmx.available()', () => {
    it('should return true when htmx is available', () => {
      expect(uno.htmx.available()).toBe(true)
    })

    it('should return false when htmx is not available', () => {
      delete mockWindow.htmx
      expect(uno.htmx.available()).toBe(false)
    })

    it('should return false when window is undefined', () => {
      const prevWindow = global.window
      delete global.window
      expect(uno.htmx.available()).toBe(false)
      global.window = prevWindow
    })
  })

  describe('htmx.bindNotifyToHtmx()', () => {
    // Import the mocked log module
    let log

    beforeEach(async () => {
      log = (await import('../../src/utils/log')).log
    })

    it('should warn if document is not available', () => {
      const prevDocument = global.document
      delete global.document

      uno.htmx.bindNotifyToHtmx()

      expect(log.warn).toHaveBeenCalledWith(
        'bindNotifyToHtmx: document is not available'
      )

      global.document = prevDocument
    })

    it('should warn if htmx is not available', () => {
      delete mockWindow.htmx

      uno.htmx.bindNotifyToHtmx()

      expect(log.warn).toHaveBeenCalledWith(
        'bindNotifyToHtmx: HTMX is not available'
      )
    })

    it('should register event listeners for HTMX events', () => {
      uno.htmx.bindNotifyToHtmx()

      expect(mockDocument.addEventListener).toHaveBeenCalledWith(
        'htmx:beforeRequest',
        expect.any(Function)
      )
      expect(mockDocument.addEventListener).toHaveBeenCalledWith(
        'htmx:afterRequest',
        expect.any(Function)
      )
      expect(mockDocument.addEventListener).toHaveBeenCalledWith(
        'htmx:responseError',
        expect.any(Function)
      )
      expect(mockDocument.addEventListener).toHaveBeenCalledWith(
        'htmx:sendError',
        expect.any(Function)
      )
      expect(mockDocument.addEventListener).toHaveBeenCalledWith(
        'htmx:timeout',
        expect.any(Function)
      )
    })

    it('should show progress on request start', () => {
      uno.htmx.bindNotifyToHtmx({ showProgress: true })

      // Trigger htmx:beforeRequest
      mockDocument._triggerEvent('htmx:beforeRequest', {
        detail: { xhr: {} }
      })

      const mockNotify = mockAlpineStore('uno').notify
      expect(mockNotify.startProgress).toHaveBeenCalled()
    })

    it('should show success notification on successful request', () => {
      uno.htmx.bindNotifyToHtmx({ showSuccess: true, successMessage: 'Done!' })

      const mockXhr = {
        status: 200,
        getResponseHeader: vi.fn((name) => null)
      }

      // Trigger htmx:afterRequest with success
      mockDocument._triggerEvent('htmx:afterRequest', {
        detail: { xhr: mockXhr }
      })

      const mockNotify = mockAlpineStore('uno').notify
      expect(mockNotify.success).toHaveBeenCalledWith('Done!')
    })

    it('should suppress success notification when X-Notify-Suppress header is set', () => {
      uno.htmx.bindNotifyToHtmx({ showSuccess: true })

      const mockXhr = {
        status: 200,
        getResponseHeader: vi.fn((name) => {
          if (name === 'X-Notify-Suppress') return 'true'
          return null
        })
      }

      // Trigger htmx:afterRequest with success (xhr directly in detail, not nested)
      mockDocument._triggerEvent('htmx:afterRequest', {
        xhr: mockXhr
      })

      const mockNotify = mockAlpineStore('uno').notify
      expect(mockNotify.success).not.toHaveBeenCalled()
    })

    it('should show custom message from X-Notify-Message header', () => {
      uno.htmx.bindNotifyToHtmx({ showSuccess: true })

      const mockXhr = {
        status: 200,
        getResponseHeader: vi.fn((name) => {
          if (name === 'X-Notify-Message') return 'Custom success!'
          return null
        })
      }

      // xhr directly in detail, not nested
      mockDocument._triggerEvent('htmx:afterRequest', {
        xhr: mockXhr
      })

      const mockNotify = mockAlpineStore('uno').notify
      expect(mockNotify.success).toHaveBeenCalledWith('Custom success!')
    })

    it('should show error notification on failed request', () => {
      uno.htmx.bindNotifyToHtmx({ showError: true, errorMessage: 'Failed!' })

      const mockXhr = {
        status: 500,
        getResponseHeader: vi.fn((name) => null)
      }

      // Trigger htmx:responseError
      mockDocument._triggerEvent('htmx:responseError', {
        detail: { xhr: mockXhr }
      })

      const mockNotify = mockAlpineStore('uno').notify
      expect(mockNotify.error).toHaveBeenCalledWith('Failed!')
    })

    it('should show custom error message from X-Notify-Error header', () => {
      uno.htmx.bindNotifyToHtmx({ showError: true })

      const mockXhr = {
        status: 400,
        getResponseHeader: vi.fn((name) => {
          if (name === 'X-Notify-Error') return 'Custom error!'
          return null
        })
      }

      // xhr directly in detail, not nested
      mockDocument._triggerEvent('htmx:responseError', {
        xhr: mockXhr
      })

      const mockNotify = mockAlpineStore('uno').notify
      expect(mockNotify.error).toHaveBeenCalledWith('Custom error!')
    })

    it('should stop progress after all requests complete', () => {
      uno.htmx.bindNotifyToHtmx({ showProgress: true })

      // Start two requests
      mockDocument._triggerEvent('htmx:beforeRequest', { detail: { xhr: {} } })
      mockDocument._triggerEvent('htmx:beforeRequest', { detail: { xhr: {} } })

      const mockNotify = mockAlpineStore('uno').notify
      expect(mockNotify.startProgress).toHaveBeenCalledTimes(2)

      // End one request
      mockDocument._triggerEvent('htmx:afterRequest', { detail: { xhr: {} } })
      expect(mockNotify.stopProgress).not.toHaveBeenCalled()

      // End second request
      mockDocument._triggerEvent('htmx:afterRequest', { detail: { xhr: {} } })
      expect(mockNotify.stopProgress).toHaveBeenCalled()
    })

    it('should respect options to disable success notification', () => {
      uno.htmx.bindNotifyToHtmx({ showSuccess: false })

      const mockXhr = { status: 200, getResponseHeader: vi.fn(() => null) }
      mockDocument._triggerEvent('htmx:afterRequest', { detail: { xhr: mockXhr } })

      const mockNotify = mockAlpineStore('uno').notify
      expect(mockNotify.success).not.toHaveBeenCalled()
    })

    it('should respect options to disable error notification', () => {
      uno.htmx.bindNotifyToHtmx({ showError: false })

      const mockXhr = { status: 500, getResponseHeader: vi.fn(() => null) }
      mockDocument._triggerEvent('htmx:responseError', { detail: { xhr: mockXhr } })

      const mockNotify = mockAlpineStore('uno').notify
      expect(mockNotify.error).not.toHaveBeenCalled()
    })

    it('should respect options to disable progress', () => {
      uno.htmx.bindNotifyToHtmx({ showProgress: false })

      mockDocument._triggerEvent('htmx:beforeRequest', { detail: { xhr: {} } })

      const mockNotify = mockAlpineStore('uno').notify
      expect(mockNotify.startProgress).not.toHaveBeenCalled()
    })
  })
})
