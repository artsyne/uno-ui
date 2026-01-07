import { describe, it, expect, beforeEach, vi } from 'vitest'
import { registerPlugin, getStore } from '../../src/core/plugin-registry'

describe('Plugin Registry', () => {
  let Alpine
  let mockStoreData

  beforeEach(() => {
    mockStoreData = {}

    // Mock Alpine.store
    Alpine = {
      store: vi.fn((name, value) => {
        if (value !== undefined) {
          mockStoreData[name] = value
        }
        return mockStoreData[name]
      })
    }
  })

  describe('registerPlugin()', () => {
    it('should create uno store with first plugin', () => {
      const boxStore = { open: vi.fn(), close: vi.fn() }

      registerPlugin(Alpine, 'box', boxStore)

      expect(Alpine.store).toHaveBeenCalledWith('uno', { box: boxStore })
      expect(mockStoreData.uno).toEqual({ box: boxStore })
    })

    it('should add plugin to existing uno store', () => {
      const boxStore = { open: vi.fn() }
      const notifyStore = { show: vi.fn() }

      registerPlugin(Alpine, 'box', boxStore)
      registerPlugin(Alpine, 'notify', notifyStore)

      expect(mockStoreData.uno).toEqual({
        box: boxStore,
        notify: notifyStore
      })
    })

    it('should preserve existing plugins when adding new ones', () => {
      const boxStore = { open: vi.fn() }
      const notifyStore = { show: vi.fn() }
      const tooltipStore = { show: vi.fn() }

      registerPlugin(Alpine, 'box', boxStore)
      registerPlugin(Alpine, 'notify', notifyStore)
      registerPlugin(Alpine, 'tooltip', tooltipStore)

      expect(mockStoreData.uno).toEqual({
        box: boxStore,
        notify: notifyStore,
        tooltip: tooltipStore
      })
    })

    it('should maintain Alpine reactivity by re-registering entire store', () => {
      const boxStore = { open: vi.fn() }
      const notifyStore = { show: vi.fn() }

      registerPlugin(Alpine, 'box', boxStore)

      const storeCallCount1 = Alpine.store.mock.calls.length

      registerPlugin(Alpine, 'notify', notifyStore)

      const storeCallCount2 = Alpine.store.mock.calls.length

      // Should call Alpine.store multiple times (uno store + individual stores)
      expect(storeCallCount2).toBeGreaterThan(storeCallCount1)

      // Verify uno store contains both plugins
      expect(mockStoreData.uno).toEqual({
        box: boxStore,
        notify: notifyStore
      })
    })

    it('should create compatibility layer with individual plugin stores', () => {
      const boxStore = { open: vi.fn() }

      registerPlugin(Alpine, 'box', boxStore)

      // Should also register Alpine.store('box') for backward compatibility
      expect(Alpine.store).toHaveBeenCalledWith('box', boxStore)
    })

    it('should not create individual store if already exists', () => {
      const boxStore = { open: vi.fn() }

      // Pre-register box store
      mockStoreData.box = { existing: true }

      registerPlugin(Alpine, 'box', boxStore)

      // Should not overwrite existing individual store
      expect(mockStoreData.box).toEqual({ existing: true })
    })

    it('should handle null Alpine gracefully', () => {
      console.error = vi.fn()

      registerPlugin(null, 'box', {})

      expect(console.error).toHaveBeenCalledWith(
        expect.stringContaining('Alpine.store is not available')
      )
    })

    it('should handle Alpine without store function', () => {
      console.error = vi.fn()

      registerPlugin({}, 'box', {})

      expect(console.error).toHaveBeenCalledWith(
        expect.stringContaining('Alpine.store is not available')
      )
    })
  })

  describe('getStore()', () => {
    beforeEach(() => {
      const boxStore = { open: vi.fn() }
      const notifyStore = { show: vi.fn() }

      registerPlugin(Alpine, 'box', boxStore)
      registerPlugin(Alpine, 'notify', notifyStore)
    })

    it('should return entire uno store when no plugin specified', () => {
      const store = getStore(Alpine)

      expect(store).toEqual(mockStoreData.uno)
      expect(store.box).toBeDefined()
      expect(store.notify).toBeDefined()
    })

    it('should return specific plugin store', () => {
      const boxStore = getStore(Alpine, 'box')

      expect(boxStore).toBe(mockStoreData.uno.box)
    })

    it('should return null if Alpine is invalid', () => {
      expect(getStore(null)).toBeNull()
      expect(getStore({})).toBeNull()
    })

    it('should return null if uno store not registered', () => {
      Alpine.store.mockImplementation(() => null)

      expect(getStore(Alpine)).toBeNull()
    })

    it('should return null for non-existent plugin', () => {
      const store = getStore(Alpine, 'nonexistent')

      expect(store).toBeUndefined()
    })
  })

  describe('Reactivity Preservation', () => {
    it('should ensure nested properties remain reactive', () => {
      const boxStore = {
        items: [],
        open: vi.fn(function () {
          this.items.push({ id: 1 })
        })
      }

      registerPlugin(Alpine, 'box', boxStore)

      const unoStore = mockStoreData.uno

      // Call method
      unoStore.box.open()

      // Verify state change
      expect(unoStore.box.items).toHaveLength(1)
    })

    it('should allow multiple plugins to coexist reactively', () => {
      const boxStore = {
        count: 0,
        increment: vi.fn(function () {
          this.count++
        })
      }
      const notifyStore = {
        messages: [],
        add: vi.fn(function (msg) {
          this.messages.push(msg)
        })
      }

      registerPlugin(Alpine, 'box', boxStore)
      registerPlugin(Alpine, 'notify', notifyStore)

      const uno = mockStoreData.uno

      uno.box.increment()
      uno.notify.add('test')

      expect(uno.box.count).toBe(1)
      expect(uno.notify.messages).toEqual(['test'])
    })
  })

  describe('Edge Cases', () => {
    it('should handle plugin with undefined store gracefully', () => {
      expect(() => {
        registerPlugin(Alpine, 'undefined-plugin', undefined)
      }).not.toThrow()
    })

    it('should handle empty plugin name', () => {
      const store = { test: true }

      registerPlugin(Alpine, '', store)

      expect(mockStoreData.uno['']).toBe(store)
    })

    it('should handle plugin overwrite (same name registered twice)', () => {
      const store1 = { version: 1 }
      const store2 = { version: 2 }

      registerPlugin(Alpine, 'test', store1)
      registerPlugin(Alpine, 'test', store2)

      expect(mockStoreData.uno.test).toBe(store2)
    })
  })
})
