import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import createNotify from '../../src/components/notify'
import { Uno } from '../../src/core/uno'
import { log } from '../../src/utils/log'
import { DomManager } from '../../src/utils/dom'

// Mock dependencies
vi.mock('../../src/core/uno', () => ({
  Uno: {
    getConfig: vi.fn(),
    registerPluginDefaults: vi.fn(),
    getStyle: vi.fn((comp, key) => `${comp}-${key}`), // Mock style class
    options: {
      ui: {
        layers: { notifyBase: 5000 }
      }
    }
  }
}))

vi.mock('../../src/utils/log', () => ({
  log: {
    info: vi.fn(),
    error: vi.fn()
  }
}))

vi.mock('../../src/utils/dom', () => ({
  DomManager: {
    ensureContainer: vi.fn(() => document.createElement('div'))
  }
}))

// Mock Alpine
const Alpine = {
  store: vi.fn(),
  directive: vi.fn(),
  initTree: vi.fn()
}

describe('Notify Component', () => {
  let notifyStore
  let mockStore // Helper to simulate Alpine store registry

  beforeEach(() => {
    vi.clearAllMocks()
    vi.useFakeTimers()

    // Reset store mock
    mockStore = {}
    Alpine.store.mockImplementation((name, value) => {
      if (value) {
        mockStore[name] = value
        if (value.init) value.init()
      }
      return mockStore[name]
    })

    // Setup Uno.getConfig mock default behavior
    Uno.getConfig.mockImplementation((key, opts) => ({
      ...opts
    }))

    // Initialize component
    createNotify(Alpine)
    notifyStore = mockStore.uno.notify
  })

  afterEach(() => {
    vi.useRealTimers()
  })

  it('should register store under uno namespace', () => {
    expect(Alpine.store).toHaveBeenCalledWith(
      'uno',
      expect.objectContaining({ notify: expect.any(Object) })
    )
    expect(notifyStore).toBeDefined()
  })

  it('should initialize container on init', () => {
    expect(DomManager.ensureContainer).toHaveBeenCalled()

    // Verify template contains correct HTMX init logic
    const templateCall = DomManager.ensureContainer.mock.calls.find(
      (call) => call[0] === 'uno-notify-container'
    )
    expect(templateCall).toBeDefined()
    const templateHtml = templateCall[1]
    expect(templateHtml).toContain('htmx.process($el)')
    expect(templateHtml).toContain("$el.setAttribute('hx-get', item.hxUrl)")
    expect(templateHtml).toContain('@htmx:beforeRequest.stop')
    expect(templateHtml).toContain('@htmx:afterRequest.stop')
    expect(templateHtml).toContain('@htmx:afterSwap.stop')
    expect(templateHtml).toContain('@htmx:responseError.stop')
  })

  describe('show()', () => {
    it('should add notification to correct group', () => {
      const id = notifyStore.show('Test Message', { position: 'top-left' })

      expect(notifyStore.groups['top-left']).toHaveLength(1)
      expect(notifyStore.groups['top-left'][0]).toMatchObject({
        id,
        message: 'Test Message',
        position: 'top-left'
      })
    })

    it('should default to top-right if position not specified', () => {
      notifyStore.show('Default Pos')
      expect(notifyStore.groups['top-right']).toHaveLength(1)
    })

    it('should handle auto-dismiss', () => {
      const id = notifyStore.show('Auto Dismiss', { duration: 1000 })

      expect(notifyStore.groups['top-right']).toHaveLength(1)

      vi.advanceTimersByTime(1000)

      // Should set visible to false first
      expect(notifyStore.groups['top-right'][0].visible).toBe(false)

      // Then remove after transition (300ms)
      vi.advanceTimersByTime(300)
      expect(notifyStore.groups['top-right']).toHaveLength(0)
    })

    it('should support i18n translation', () => {
      // Mock i18n store
      mockStore.i18n = {
        t: vi.fn((key) => `Translated ${key}`)
      }

      notifyStore.show('hello')
      expect(notifyStore.groups['top-right'][0].message).toBe('Translated hello')
    })
  })

  describe('remove()', () => {
    it('should remove notification by id', () => {
      const id = notifyStore.show('To Remove')
      expect(notifyStore.groups['top-right']).toHaveLength(1)

      notifyStore.remove(id)

      // Should set visible false immediately
      expect(notifyStore.groups['top-right'][0].visible).toBe(false)

      // And remove after delay
      vi.advanceTimersByTime(300)
      expect(notifyStore.groups['top-right']).toHaveLength(0)
    })
  })

  describe('Convenience Methods', () => {
    it('should support success/error/warning/info shortcuts', () => {
      notifyStore.success('Success')
      expect(notifyStore.groups['top-right'][0].type).toBe('success')

      notifyStore.error('Error')
      expect(notifyStore.groups['top-right'][1].type).toBe('error')
    })
  })

  describe('HTMX Integration', () => {
    it('should auto-manage progress for hxUrl requests', () => {
      const item1 = {}
      const item2 = {}

      const spyStartProgress = vi.spyOn(notifyStore, 'startProgress')
      const spyStopProgress = vi.spyOn(notifyStore, 'stopProgress')

      notifyStore._onHtmxBegin(item1)
      expect(spyStartProgress).toHaveBeenCalledTimes(1)
      expect(notifyStore.progressRequests).toBe(1)

      notifyStore._onHtmxBegin(item2)
      expect(spyStartProgress).toHaveBeenCalledTimes(1) // Should only call once for first request
      expect(notifyStore.progressRequests).toBe(2)

      notifyStore._onHtmxEnd(item1)
      expect(spyStopProgress).toHaveBeenCalledTimes(0) // Should not stop until all requests finish
      expect(notifyStore.progressRequests).toBe(1)

      notifyStore._onHtmxEnd(item2)
      expect(spyStopProgress).toHaveBeenCalledTimes(1) // Should stop when all requests finish
      expect(notifyStore.progressRequests).toBe(0)
    })

    it('should keep progress visible until all hxUrl requests finish', () => {
      const item1 = {}
      const item2 = {}

      const spyStartProgress = vi.spyOn(notifyStore, 'startProgress')
      const spyStopProgress = vi.spyOn(notifyStore, 'stopProgress')

      notifyStore._onHtmxBegin(item1)
      notifyStore._onHtmxBegin(item2)
      expect(spyStartProgress).toHaveBeenCalledTimes(1)
      expect(notifyStore.progressRequests).toBe(2)

      notifyStore._onHtmxEnd(item1)
      expect(spyStopProgress).toHaveBeenCalledTimes(0)
      expect(notifyStore.progressRequests).toBe(1)

      notifyStore._onHtmxEnd(item2)
      expect(spyStopProgress).toHaveBeenCalledTimes(1)
      expect(notifyStore.progressRequests).toBe(0)
    })

    it('should handle onContentLoaded', () => {
      const el = document.createElement('div')
      const item = { onContentLoad: vi.fn() }

      // Mock window.Alpine
      window.Alpine = { initTree: vi.fn() }

      notifyStore.onContentLoaded(el, item)

      expect(item.onContentLoad).toHaveBeenCalledWith(el)
      expect(window.Alpine.initTree).toHaveBeenCalledWith(el)
    })

    it('should handle onHtmxError', () => {
      const item = { onError: vi.fn() }
      const error = { status: 500 }

      // Spy on store.error to ensure it calls show error
      const spyError = vi.spyOn(notifyStore, 'error')

      notifyStore.onHtmxError(error, item)

      expect(spyError).toHaveBeenCalled()
      expect(item.onError).toHaveBeenCalledWith(error)
    })
  })

  describe('Directive', () => {
    it('should register uno-notify directive', () => {
      expect(Alpine.directive).toHaveBeenCalledWith('uno-notify', expect.any(Function))
    })

    // Testing the directive handler logic would require simulating the Alpine directive callback
    // which is complex to mock perfectly, but we can verify registration.
  })
})
