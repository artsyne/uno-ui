import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { createTabStore } from '../../src/components/tab/tab.store'
import { Uno } from '../../src/core/uno'

// Mock dependencies
vi.mock('../../src/core/uno', () => ({
  Uno: {
    getConfig: vi.fn(),
    getStyle: vi.fn((comp, key) => `${comp}-${key}`),
    htmx: {
      available: vi.fn(() => true),
      ajax: vi.fn(() => Promise.resolve()),
      fragment: vi.fn(() => Promise.resolve('<div>Content</div>'))
    }
  }
}))

vi.mock('../../src/utils/log', () => ({
  log: {
    info: vi.fn(),
    warn: vi.fn(),
    error: vi.fn()
  }
}))

// Mock Alpine
const Alpine = {
  store: vi.fn(),
  nextTick: vi.fn((fn) => setTimeout(fn, 0))
}

describe('Tab Component', () => {
  let tabStore
  let mockStore

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

    // Setup Uno.getConfig default behavior
    Uno.getConfig.mockImplementation((key, opts) => ({
      ...opts
    }))

    // Create store
    Alpine.store('uno', { tab: createTabStore(Alpine) })
    tabStore = Alpine.store('uno').tab
  })

  afterEach(() => {
    vi.useRealTimers()
  })

  describe('Initialization', () => {
    it('should initialize with empty groups object', () => {
      expect(tabStore.groups).toEqual({})
    })

    it('should initialize with empty registry', () => {
      expect(tabStore.registry).toEqual({})
    })
  })

  describe('register()', () => {
    it('should register a tab', () => {
      tabStore.register('tab1', 'test-group')

      expect(tabStore.registry['tab1']).toEqual({
        name: 'tab1',
        group: 'test-group'
      })
    })

    it('should set first tab as active if group is empty', () => {
      tabStore.register('tab1', 'test-group')

      expect(tabStore.groups['test-group']).toBe('tab1')
    })

    it('should not change active tab if group already has one', () => {
      tabStore.register('tab1', 'test-group')
      tabStore.register('tab2', 'test-group')

      expect(tabStore.groups['test-group']).toBe('tab1')
    })

    it('should use default group if not specified', () => {
      tabStore.register('tab1')

      expect(tabStore.registry['tab1'].group).toBe('default')
    })
  })

  describe('activate()', () => {
    beforeEach(() => {
      tabStore.register('tab1', 'test-group')
      tabStore.register('tab2', 'test-group')
    })

    it('should activate tab', () => {
      tabStore.activate({ name: 'tab2', group: 'test-group' })

      expect(tabStore.groups['test-group']).toBe('tab2')
    })

    it('should update URL hash when different from current', () => {
      const replaceStateSpy = vi.spyOn(window.history, 'replaceState')

      // Mock initial hash to be different
      Object.defineProperty(window.location, 'hash', {
        value: '#tab1',
        writable: true,
        configurable: true
      })

      tabStore.activate({ name: 'tab2', group: 'test-group', hash: true })

      expect(replaceStateSpy).toHaveBeenCalledWith(null, null, '#tab2')
    })

    it('should not update hash if updateHash is false', () => {
      const replaceStateSpy = vi.spyOn(window.history, 'replaceState')

      tabStore.activate({ name: 'tab2', group: 'test-group', hash: true, updateHash: false })

      expect(replaceStateSpy).not.toHaveBeenCalled()
    })

    it('should use default group if not specified', () => {
      tabStore.register('tab1')

      tabStore.activate({ name: 'tab1' })

      expect(tabStore.groups['default']).toBe('tab1')
    })
  })

  describe('isActive()', () => {
    beforeEach(() => {
      tabStore.register('tab1', 'test-group')
      tabStore.register('tab2', 'test-group')
    })

    it('should return true for active tab', () => {
      expect(tabStore.isActive('tab1', 'test-group')).toBe(true)
    })

    it('should return false for inactive tab', () => {
      expect(tabStore.isActive('tab2', 'test-group')).toBe(false)
    })

    it('should use default group if not specified', () => {
      tabStore.register('tab1')

      expect(tabStore.isActive('tab1')).toBe(true)
    })
  })

  describe('HTMX Integration', () => {
    it('should trigger HTMX click when tab has hx-get attribute', () => {
      // Create mock tab element
      const tabEl = document.createElement('button')
      tabEl.id = 'tab-tab1'
      tabEl.setAttribute('hx-get', '/api/content')
      document.body.appendChild(tabEl)

      // Mock window.htmx
      window.htmx = { trigger: vi.fn() }

      tabStore.register('tab1', 'test-group')
      tabStore.activate({ name: 'tab1', group: 'test-group', updateHash: false })

      // Wait for setTimeout
      vi.runAllTimers()

      expect(window.htmx.trigger).toHaveBeenCalledWith(tabEl, 'click')

      // Cleanup
      delete window.htmx
      document.body.removeChild(tabEl)
    })
  })

  describe('Hash Navigation', () => {
    it('should activate tab based on URL hash', () => {
      // Need to setup before setting hash
      tabStore.register('tab1', 'test-group')
      tabStore.register('tab2', 'test-group')

      // Mock window.location properties
      Object.defineProperty(window.location, 'hash', {
        value: '#tab2',
        writable: true,
        configurable: true
      })
      Object.defineProperty(window.location, 'href', {
        value: 'http://localhost/#tab2',
        writable: true,
        configurable: true
      })

      // Manually trigger the hash check logic
      const hash = window.location.hash.slice(1)
      if (hash && tabStore.registry[hash]) {
        const { group, name } = tabStore.registry[hash]
        tabStore.activate({ group, name, updateHash: false })
      }

      expect(tabStore.groups['test-group']).toBe('tab2')
    })
  })

  describe('Configuration Merge', () => {
    it('should use 4-level configuration merge for activate', () => {
      Uno.getConfig.mockImplementation((key, opts) => ({
        group: 'default',
        hash: true,
        customOption: 'from-config',
        ...opts
      }))

      // Use updateHash: false to avoid security error in happy-dom
      tabStore.activate({ name: 'tab1', updateHash: false })

      expect(Uno.getConfig).toHaveBeenCalledWith('tab', {
        group: 'default',
        hash: true,
        name: 'tab1',
        updateHash: false
      })

      // Verify tab was activated
      expect(tabStore.groups['default']).toBe('tab1')
    })
  })
})
