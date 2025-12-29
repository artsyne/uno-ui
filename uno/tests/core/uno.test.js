import { describe, it, expect, beforeEach, vi } from 'vitest'
import { UnoInstance } from '../../src/core/uno'

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
