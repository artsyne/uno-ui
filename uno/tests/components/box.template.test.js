import { describe, it, expect, vi, beforeEach } from 'vitest'
import { getBoxTemplate } from '../../src/components/box/box.template'
import { Uno } from '../../src/core/uno'
import { getZIndex } from '../../src/core/layer'

// Mock dependencies
vi.mock('../../src/core/uno', () => ({
  Uno: {
    options: {
      ui: {
        layers: {
          box: 1400
        }
      }
    },
    getStyle: vi.fn(() => 'mock-class')
  }
}))

vi.mock('../../src/core/layer', () => ({
  getZIndex: vi.fn(() => 1400)
}))

describe('Box Template', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('should include z-index logic with index multiplier', () => {
    const template = getBoxTemplate()

    // Check for x-for with index
    expect(template).toContain('x-for="(item, index) in $store.uno.box.items"')

    // Check for z-index binding
    // We expect the literal string from the template, where ${zBase} is interpolated
    // Since we mocked zBase to be 1400 (either from options or getZIndex)
    const expectedZIndexLogic = ':style="\'z-index:\' + (1400 + index * 10)"'
    expect(template).toContain(expectedZIndexLogic)
  })

  it('should include grid layout classes for stacking', () => {
    const template = getBoxTemplate()

    // Check for grid classes on dialogWrapper
    // Note: removed grid-rows-1 to allow auto height for scrolling
    expect(template).toContain("'grid grid-cols-1'")

    // Check for grid item classes on itemWrapper
    expect(template).toContain("'col-start-1 row-start-1'")
  })

  it('should use z-index from Uno options if available', () => {
    Uno.options.ui.layers.box = 2000
    const template = getBoxTemplate()
    const expectedZIndexLogic = ':style="\'z-index:\' + (2000 + index * 10)"'
    expect(template).toContain(expectedZIndexLogic)
  })

  it('should fallback to getZIndex if options not set', () => {
    Uno.options.ui.layers.box = undefined
    const template = getBoxTemplate()
    const expectedZIndexLogic = ':style="\'z-index:\' + (1400 + index * 10)"'
    expect(template).toContain(expectedZIndexLogic)
  })
})
