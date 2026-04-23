/**
 * TypeScript Type Definitions Tests
 *
 * This file tests the TypeScript type definitions to ensure they:
 * 1. Can be imported correctly
 * 2. Have valid type signatures
 * 3. Support basic type checking
 * 4. Maintain type safety
 *
 * @version 0.6.0
 */

import { describe, it, expect } from 'vitest'

// Type imports for testing
import type {
  UnoInstance,
  UnoOptions,
  BoxOptions,
  NotifyOptions,
  TooltipOptions,
  TabOptions,
  Plugin
} from '../../types/index'

describe('TypeScript Type Definitions', () => {
  describe('Core Types', () => {
    it('should allow valid UnoOptions configurations', () => {
      const config: UnoOptions = {
        debug: true,
        csp: { nonce: 'random-nonce' },
        onError: (error: Error) => console.error(error),
        plugins: {
          box: {
            defaultWidth: 500,
            defaultPosition: 'center'
          },
          notify: {
            duration: 5000,
            position: 'top-right'
          }
        }
      }

      expect(config.debug).toBe(true)
      expect(config.csp?.nonce).toBe('random-nonce')
      expect(config.plugins?.box?.defaultWidth).toBe(500)
    })

    it('should allow minimal UnoOptions', () => {
      const config: UnoOptions = {}
      expect(config).toEqual({})
    })

    it('should support Plugin type', () => {
      const plugin: Plugin = (Alpine, options) => {
        expect(Alpine).toBeDefined()
        expect(options).toBeDefined()
      }

      expect(typeof plugin).toBe('function')
    })
  })

  describe('Box Component Types', () => {
    it('should allow valid BoxOptions configurations', () => {
      const config1: BoxOptions = {
        type: 'dialog',
        title: 'Test Dialog',
        content: 'This is a test',
        position: 'center',
        width: 500
      }

      const config2: BoxOptions = {
        type: 'confirm',
        title: 'Confirm',
        showOkButton: true,
        showCancelButton: true,
        onOk: () => console.log('OK'),
        onCancel: () => console.log('Cancel')
      }

      const config3: BoxOptions = {
        title: 'Enter Name',
        content: 'Please enter your name:'
      }

      expect(config1.type).toBe('dialog')
      expect(config2.type).toBe('confirm')
      expect(config3.title).toBe('Enter Name')
    })

    it('should allow minimal BoxOptions', () => {
      const config: BoxOptions = {}
      expect(config).toEqual({})
    })

    it('should support callback return types', () => {
      const config: BoxOptions = {
        onOk: () => false, // Can return false to prevent closing
        onCancel: () => {}, // Or return void
        onClose: () => console.log('closed')
      }

      expect(typeof config.onOk).toBe('function')
    })
  })

  describe('Notify Component Types', () => {
    it('should allow valid NotifyOptions configurations', () => {
      const config1: NotifyOptions = {
        type: 'success',
        duration: 3000,
        position: 'top-right',
        closable: true
      }

      const config2: NotifyOptions = {
        type: 'error',
        duration: 0, // No auto-close
        position: 'bottom-left'
      }

      expect(config1.type).toBe('success')
      expect(config2.duration).toBe(0)
    })

    it('should support all notification types', () => {
      const types: Array<NotifyOptions['type']> = ['success', 'error', 'warning', 'info']
      expect(types).toHaveLength(4)
    })
  })

  describe('Tooltip Component Types', () => {
    it('should allow valid TooltipOptions configurations', () => {
      const config1: TooltipOptions = {
        content: 'Tooltip text',
        placement: 'top',
        trigger: 'hover'
      }

      const config2: TooltipOptions = {
        content: 'Click me',
        placement: 'bottom',
        trigger: 'click',
        arrow: true
      }

      expect(config1.placement).toBe('top')
      expect(config2.trigger).toBe('click')
    })

    it('should support all placement options', () => {
      const placements: Array<TooltipOptions['placement']> = ['top', 'bottom', 'left', 'right']
      expect(placements).toHaveLength(4)
    })

    it('should support all trigger modes', () => {
      const triggers: Array<TooltipOptions['trigger']> = ['hover', 'click', 'focus']
      expect(triggers).toHaveLength(3)
    })
  })

  describe('Tab Component Types', () => {
    it('should allow valid TabOptions configurations', () => {
      const config1: TabOptions = {
        name: 'tab1',
        group: 'main',
        hash: true
      }

      const config2: TabOptions = {
        name: 'tab2',
        group: 'sidebar',
        hash: false
      }

      expect(config1.name).toBe('tab1')
      expect(config2.hash).toBe(false)
    })
  })

  describe('Type Safety', () => {
    it('should enforce type constraints', () => {
      // These should compile without errors
      const boxType: BoxOptions['type'] = 'dialog'
      const notifyType: NotifyOptions['type'] = 'success'
      const tooltipPlacement: TooltipOptions['placement'] = 'top'

      expect(boxType).toBeDefined()
      expect(notifyType).toBeDefined()
      expect(tooltipPlacement).toBeDefined()
    })

    it('should allow optional properties', () => {
      // All configuration properties should be optional
      const box: BoxOptions = {}
      const notify: NotifyOptions = {}
      const tooltip: TooltipOptions = {}
      const tab: TabOptions = {}

      expect(box).toEqual({})
      expect(notify).toEqual({})
      expect(tooltip).toEqual({})
      expect(tab).toEqual({})
    })

    it('should handle null values correctly', () => {
      const config: UnoOptions = {
        csp: { nonce: null },
        onError: null
      }

      expect(config.csp?.nonce).toBeNull()
      expect(config.onError).toBeNull()
    })
  })

  describe('Plugin Configuration', () => {
    it('should support plugin-specific configurations', () => {
      const config: UnoOptions = {
        plugins: {
          box: {
            defaultWidth: 600,
            defaultPosition: 'center',
            okText: 'Confirm',
            cancelText: 'Close'
          },
          notify: {
            duration: 4000,
            position: 'bottom-right'
          },
          tooltip: {
            placement: 'bottom',
            trigger: 'hover',
            arrow: false
          },
          tab: {
            hash: false
          }
        }
      }

      expect(config.plugins?.box?.defaultWidth).toBe(600)
      expect(config.plugins?.notify?.duration).toBe(4000)
    })
  })

  describe('Backward Compatibility', () => {
    it('should support empty configurations', () => {
      const uno: UnoOptions = {}
      const box: BoxOptions = {}
      const notify: NotifyOptions = {}

      expect(uno).toEqual({})
      expect(box).toEqual({})
      expect(notify).toEqual({})
    })
  })

  describe('Complex Configurations', () => {
    it('should support nested configurations', () => {
      const config: UnoOptions = {
        debug: true,
        csp: {
          nonce: 'abc123'
        },
        plugins: {
          box: {
            defaultWidth: 500,
            okText: 'OK',
            cancelText: 'Cancel'
          }
        },
        ui: {
          base: 'rounded-lg shadow-md',
          layers: {
            box: 1000,
            tooltip: 2000,
            notify: 3000
          }
        }
      }

      expect(config.debug).toBe(true)
      expect(config.csp?.nonce).toBe('abc123')
      expect(config.ui?.layers?.box).toBe(1000)
    })
  })
})
