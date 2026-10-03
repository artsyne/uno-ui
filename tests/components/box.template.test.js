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

  // ============================================================
  // extraButtons - Template Rendering
  // ============================================================

  it('should include handleExtraButton click binding', () => {
    const template = getBoxTemplate()
    expect(template).toContain('handleExtraButton(item, btn)')
  })

  it('should render extraButtons with x-for loop', () => {
    const template = getBoxTemplate()
    // position='before' loop
    expect(template).toContain(
      "(item.extraButtons || []).filter(b => (b.position || 'before') === 'before')"
    )
    // position='after' loop
    expect(template).toContain("(item.extraButtons || []).filter(b => b.position === 'after')")
  })

  it('should include extraButtons.length in footer x-show condition', () => {
    const template = getBoxTemplate()
    expect(template).toContain('item.extraButtons && item.extraButtons.length > 0')
  })

  it('should use buttonClass as replacement (not additive) for extraBtn default style', () => {
    const template = getBoxTemplate()
    // btn.buttonClass || Uno.getStyle('box', 'extraBtn') means:
    // if buttonClass is set, it REPLACES the default; otherwise fallback to extraBtn
    expect(template).toContain("btn.buttonClass || Uno.getStyle('box', 'extraBtn')")
  })

  it('should apply extraBtn style as fallback when buttonClass is not set', () => {
    const template = getBoxTemplate()
    // extraBtn is still referenced as the fallback
    expect(template).toContain("Uno.getStyle('box', 'extraBtn')")
  })

  it("should apply sm:order-last to position='before' buttons for correct visual ordering", () => {
    const template = getBoxTemplate()
    // position='before' buttons need sm:order-last to appear visually LEFT of OK
    // because the footer uses sm:flex-row-reverse (DOM-first = visual-right)
    const beforeLoopStart = template.indexOf(
      "(item.extraButtons || []).filter(b => (b.position || 'before') === 'before')"
    )
    const afterLoopStart = template.indexOf(
      "(item.extraButtons || []).filter(b => b.position === 'after')"
    )
    const beforeSection = template.slice(beforeLoopStart, afterLoopStart)
    expect(beforeSection).toContain("'sm:order-last'")
    // position='after' buttons should NOT have sm:order-last
    const afterSection = template.slice(afterLoopStart)
    const afterButtonClass = afterSection.match(/:class="\[([^"]+)\]"/)?.[1] || ''
    expect(afterButtonClass).not.toContain('sm:order-last')
  })

  it('should reference btn.buttonClass in class binding', () => {
    const template = getBoxTemplate()
    expect(template).toContain('btn.buttonClass')
  })

  it('should localize close button aria-label via i18n instead of hardcoding (B7)', () => {
    const template = getBoxTemplate()
    // 关闭按钮 aria-label 走 x-t 指令本地化（box.close），不硬编码 "Close"
    expect(template).toContain('x-t.aria-label="box.close"')
    expect(template).not.toContain('aria-label="Close"')
  })
})
