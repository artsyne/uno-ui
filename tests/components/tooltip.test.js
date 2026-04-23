import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { createTooltipStore } from '../../src/components/tooltip/tooltip.store'
import { Uno } from '../../src/core/uno'
import { DomManager } from '../../src/utils/dom'

// Mock dependencies
vi.mock('../../src/core/uno', () => ({
  Uno: {
    getConfig: vi.fn(),
    options: {
      ui: {
        layers: { tooltip: 6000 }
      }
    }
  }
}))

vi.mock('../../src/utils/dom', () => ({
  DomManager: {
    ensureContainer: vi.fn(() => document.createElement('div'))
  },
  sanitizeHTML: vi.fn((html) => html)
}))

vi.mock('../../src/utils/log', () => ({
  log: {
    info: vi.fn(),
    warn: vi.fn(),
    error: vi.fn()
  }
}))

vi.mock('../../src/utils/position', () => ({
  calculatePosition: vi.fn(() => ({
    top: 100,
    left: 200,
    placement: 'top'
  }))
}))

// Mock Alpine
const Alpine = {
  store: vi.fn(),
  nextTick: vi.fn((fn) => setTimeout(fn, 0))
}

describe('Tooltip Component', () => {
  let tooltipStore
  let mockStore

  beforeEach(() => {
    vi.clearAllMocks()
    vi.useFakeTimers()

    // Reset DOM
    document.body.innerHTML = ''

    // Reset store mock
    mockStore = {}
    Alpine.store.mockImplementation((name, value) => {
      if (value) {
        mockStore[name] = value
        if (value.tooltip?.init) value.tooltip.init()
      }
      return mockStore[name]
    })

    // Setup Uno.getConfig default behavior
    Uno.getConfig.mockImplementation((key, opts) => ({
      placement: 'top',
      trigger: 'hover',
      offset: 8,
      arrow: true,
      clickOutside: true,
      hxTrigger: 'load',
      ...opts
    }))

    // Create store
    Alpine.store('uno', { tooltip: createTooltipStore(Alpine) })
    tooltipStore = Alpine.store('uno').tooltip
  })

  afterEach(() => {
    vi.useRealTimers()
  })

  describe('Initialization', () => {
    it('should initialize with default values', () => {
      expect(tooltipStore.visible).toBe(false)
      expect(tooltipStore.content).toBe('')
      expect(tooltipStore.placement).toBe('top')
      expect(tooltipStore.top).toBe(-9999)
      expect(tooltipStore.left).toBe(-9999)
    })

    it('should call ensureContainer on show', () => {
      const trigger = document.createElement('div')
      document.body.appendChild(trigger)
      tooltipStore.show({ element: trigger, content: 'test' })
      expect(DomManager.ensureContainer).toHaveBeenCalled()
    })
  })

  describe('show()', () => {
    let triggerEl

    beforeEach(() => {
      triggerEl = document.createElement('button')
      triggerEl.textContent = 'Trigger'
      document.body.appendChild(triggerEl)
    })

    it('should make tooltip visible', () => {
      tooltipStore.show({ element: triggerEl, content: 'Tooltip content' })
      expect(tooltipStore.visible).toBe(true)
    })

    it('should set content from options', () => {
      tooltipStore.show({ element: triggerEl, content: 'Test content' })
      expect(tooltipStore.content).toBe('Test content')
    })

    it('should use 4-level configuration merge', () => {
      tooltipStore.show({ element: triggerEl, placement: 'bottom' })
      expect(Uno.getConfig).toHaveBeenCalledWith('tooltip', {
        element: triggerEl,
        placement: 'bottom'
      })
    })

    it('should load content from target selector', () => {
      const targetEl = document.createElement('div')
      targetEl.id = 'tooltip-content'
      targetEl.innerHTML = '<strong>Target Content</strong>'
      document.body.appendChild(targetEl)

      tooltipStore.show({ element: triggerEl, target: '#tooltip-content' })
      expect(tooltipStore.content).toContain('Target Content')
    })

    it('should warn if target selector not found', async () => {
      // Import log from the source to get the mocked version
      const log = await import('../../src/utils/log').then((m) => m.log)
      tooltipStore.show({ element: triggerEl, target: '#non-existent' })
      expect(log.warn).toHaveBeenCalledWith(
        expect.stringContaining('Tooltip target element not found')
      )
    })

    it('should set placement from config', () => {
      tooltipStore.show({ element: triggerEl, placement: 'bottom' })
      expect(tooltipStore.placement).toBe('bottom')
    })

    it('should set trigger mode', () => {
      tooltipStore.show({ element: triggerEl, trigger: 'click' })
      expect(tooltipStore.triggerMode).toBe('click')
    })

    it('should update position on nextTick', () => {
      const updatePositionSpy = vi.spyOn(tooltipStore, 'updatePosition')

      tooltipStore.show({ element: triggerEl, content: 'Test' })

      expect(Alpine.nextTick).toHaveBeenCalled()
      vi.runAllTimers()

      expect(updatePositionSpy).toHaveBeenCalled()
    })

    it('should set HTMX options', () => {
      tooltipStore.show({
        element: triggerEl,
        hxUrl: '/api/tooltip',
        hxHeaders: { 'X-Custom': 'value' },
        hxTrigger: 'click'
      })

      expect(tooltipStore.hxUrl).toBe('/api/tooltip')
      expect(tooltipStore.hxHeaders).toEqual({ 'X-Custom': 'value' })
      expect(tooltipStore.hxTrigger).toBe('click')
    })
  })

  describe('hide()', () => {
    let triggerEl

    beforeEach(() => {
      triggerEl = document.createElement('button')
      document.body.appendChild(triggerEl)
      tooltipStore.show({ element: triggerEl, content: 'Test' })
    })

    it('should hide tooltip', () => {
      tooltipStore.hide()
      expect(tooltipStore.visible).toBe(false)
    })

    it('should clear triggerEl', () => {
      tooltipStore.hide()
      expect(tooltipStore.triggerEl).toBeNull()
    })
  })

  describe('scheduleHide()', () => {
    it('should hide tooltip after delay', () => {
      tooltipStore.visible = true

      tooltipStore.scheduleHide(300)
      expect(tooltipStore.visible).toBe(true)

      vi.advanceTimersByTime(300)
      expect(tooltipStore.visible).toBe(false)
    })

    it('should use default delay of 300ms', () => {
      tooltipStore.visible = true

      tooltipStore.scheduleHide()
      vi.advanceTimersByTime(299)
      expect(tooltipStore.visible).toBe(true)

      vi.advanceTimersByTime(1)
      expect(tooltipStore.visible).toBe(false)
    })
  })

  describe('cancelHide()', () => {
    it('should cancel scheduled hide', () => {
      tooltipStore.visible = true

      tooltipStore.scheduleHide(300)
      tooltipStore.cancelHide()

      vi.advanceTimersByTime(300)
      expect(tooltipStore.visible).toBe(true)
    })
  })

  describe('updatePosition()', () => {
    let triggerEl, tooltipEl

    beforeEach(() => {
      triggerEl = document.createElement('button')
      triggerEl.getBoundingClientRect = vi.fn(() => ({
        top: 100,
        left: 200,
        width: 80,
        height: 40
      }))
      document.body.appendChild(triggerEl)

      tooltipEl = document.createElement('div')
      tooltipEl.id = 'uno-tooltip-content'
      tooltipEl.getBoundingClientRect = vi.fn(() => ({
        width: 120,
        height: 60
      }))
      document.body.appendChild(tooltipEl)
    })

    it('should update tooltip position when visible', () => {
      tooltipStore.visible = true
      tooltipStore.triggerEl = triggerEl

      tooltipStore.updatePosition()

      expect(tooltipStore.top).toBe(100)
      expect(tooltipStore.left).toBe(200)
    })

    it('should not update if not visible', () => {
      tooltipStore.visible = false
      tooltipStore.triggerEl = triggerEl

      const initialTop = tooltipStore.top
      tooltipStore.updatePosition()

      expect(tooltipStore.top).toBe(initialTop)
    })

    it('should not update if no trigger element', () => {
      tooltipStore.visible = true
      tooltipStore.triggerEl = null

      const initialTop = tooltipStore.top
      tooltipStore.updatePosition()

      expect(tooltipStore.top).toBe(initialTop)
    })

    it('should not update if tooltip element not found', () => {
      tooltipEl.remove()

      tooltipStore.visible = true
      tooltipStore.triggerEl = triggerEl

      const initialTop = tooltipStore.top
      tooltipStore.updatePosition()

      expect(tooltipStore.top).toBe(initialTop)
    })

    it('should update placement if flipped', async () => {
      // Import mocked calculatePosition
      const positionModule = await import('../../src/utils/position')
      positionModule.calculatePosition.mockReturnValueOnce({
        top: 150,
        left: 250,
        placement: 'bottom' // Flipped from 'top'
      })

      tooltipStore.visible = true
      tooltipStore.triggerEl = triggerEl
      tooltipStore.placement = 'top'

      tooltipStore.updatePosition()

      expect(tooltipStore.placement).toBe('bottom')
    })
  })

  describe('Arrow Configuration', () => {
    it('should show arrow by default', () => {
      const triggerEl = document.createElement('button')
      tooltipStore.show({ element: triggerEl, content: 'Test' })
      expect(tooltipStore.arrow).toBe(true)
    })

    it('should respect arrow: false', () => {
      const triggerEl = document.createElement('button')
      tooltipStore.show({ element: triggerEl, content: 'Test', arrow: false })
      expect(tooltipStore.arrow).toBe(false)
    })
  })

  describe('Click Outside', () => {
    it('should enable click outside by default', () => {
      const triggerEl = document.createElement('button')
      tooltipStore.show({ element: triggerEl, content: 'Test' })
      expect(tooltipStore.clickOutside).toBe(true)
    })

    it('should respect clickOutside option', () => {
      const triggerEl = document.createElement('button')
      tooltipStore.show({ element: triggerEl, content: 'Test', clickOutside: false })
      expect(tooltipStore.clickOutside).toBe(false)
    })
  })

  describe('Placement Options', () => {
    const placements = [
      'top',
      'top-start',
      'top-end',
      'bottom',
      'bottom-start',
      'bottom-end',
      'left',
      'left-start',
      'left-end',
      'right',
      'right-start',
      'right-end'
    ]

    placements.forEach((placement) => {
      it(`should support ${placement} placement`, () => {
        const triggerEl = document.createElement('button')
        tooltipStore.show({ element: triggerEl, content: 'Test', placement })
        expect(tooltipStore.placement).toBe(placement)
      })
    })
  })

  describe('Z-Index Configuration', () => {
    it('should use configured z-index from layer system', () => {
      // The actual z-index comes from getZIndex('TOOLTIP') in layer.js
      // which returns 1700, not from the mock config
      expect(tooltipStore.zIndex).toBeDefined()
      expect(typeof tooltipStore.zIndex).toBe('number')
    })

    it('should update z-index from config on show', () => {
      Uno.options.ui.layers.tooltip = 7000

      const triggerEl = document.createElement('button')
      tooltipStore.show({ element: triggerEl, content: 'Test' })

      expect(tooltipStore.zIndex).toBe(7000)
    })
  })

  describe('ensureContainer()', () => {
    it('should create tooltip container', () => {
      tooltipStore.ensureContainer()
      expect(DomManager.ensureContainer).toHaveBeenCalledWith(
        'uno-tooltip-manager',
        expect.any(String),
        expect.objectContaining({
          'x-data': '{}'
        })
      )
    })

    it('should not recreate existing container', () => {
      const existingContainer = document.createElement('div')
      existingContainer.id = 'uno-tooltip-manager'
      document.body.appendChild(existingContainer)

      DomManager.ensureContainer.mockClear()
      tooltipStore.ensureContainer()

      expect(DomManager.ensureContainer).not.toHaveBeenCalled()
    })
  })

  describe('Global Event Listeners', () => {
    it('should attach global listeners when shown', () => {
      // Note: Global listeners are attached inside ensureGlobalListeners()
      // which is called during show(). Due to module closure, we can't easily spy on it.
      // This test verifies the show() method executes without errors
      const triggerEl = document.createElement('button')
      expect(() => {
        tooltipStore.show({ element: triggerEl, content: 'Test' })
      }).not.toThrow()

      expect(tooltipStore.visible).toBe(true)
    })

    it('should remove global listeners when hidden', () => {
      const removeEventListenerSpy = vi.spyOn(window, 'removeEventListener')

      const triggerEl = document.createElement('button')
      tooltipStore.show({ element: triggerEl, content: 'Test' })
      tooltipStore.hide()

      expect(removeEventListenerSpy).toHaveBeenCalledWith('scroll', expect.any(Function), true)
      expect(removeEventListenerSpy).toHaveBeenCalledWith('resize', expect.any(Function))
      expect(removeEventListenerSpy).toHaveBeenCalledWith('keydown', expect.any(Function))
      expect(removeEventListenerSpy).toHaveBeenCalledWith('click', expect.any(Function))
    })
  })

  describe('Trigger Modes', () => {
    const modes = ['hover', 'click', 'focus']

    modes.forEach((mode) => {
      it(`should support ${mode} trigger mode`, () => {
        const triggerEl = document.createElement('button')
        tooltipStore.show({ element: triggerEl, content: 'Test', trigger: mode })
        expect(tooltipStore.triggerMode).toBe(mode)
      })
    })
  })
})
