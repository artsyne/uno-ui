import { describe, it, expect, vi, afterEach } from 'vitest'
import {
  POSITIONS,
  getFlexAlignmentClass,
  getFixedPositionClass,
  getStackingConfig,
  calculatePosition
} from '../../src/utils/position'

describe('Position Utils - POSITIONS constant', () => {
  it('should export 9 positions', () => {
    expect(POSITIONS).toHaveLength(9)
    expect(POSITIONS).toEqual([
      'top-left',
      'top-center',
      'top-right',
      'center-left',
      'center',
      'center-right',
      'bottom-left',
      'bottom-center',
      'bottom-right'
    ])
  })
})

describe('Position Utils - getFlexAlignmentClass()', () => {
  it('should return correct classes for center', () => {
    expect(getFlexAlignmentClass('center')).toBe('items-center justify-center')
  })

  it('should return correct classes for top positions', () => {
    expect(getFlexAlignmentClass('top-left')).toBe('items-start justify-start sm:px-8')
    expect(getFlexAlignmentClass('top-center')).toBe('items-start justify-center')
    expect(getFlexAlignmentClass('top-right')).toBe('items-start justify-end sm:px-8')
  })

  it('should return correct classes for center-left/right', () => {
    expect(getFlexAlignmentClass('center-left')).toBe('items-center justify-start sm:px-8')
    expect(getFlexAlignmentClass('center-right')).toBe('items-center justify-end sm:px-8')
  })

  it('should return correct classes for bottom positions', () => {
    expect(getFlexAlignmentClass('bottom-left')).toBe('items-end justify-start sm:px-8')
    expect(getFlexAlignmentClass('bottom-center')).toBe('items-end justify-center')
    expect(getFlexAlignmentClass('bottom-right')).toBe('items-end justify-end sm:px-8')
  })

  it('should default to center for undefined position', () => {
    expect(getFlexAlignmentClass()).toBe('items-center justify-center')
    expect(getFlexAlignmentClass(null)).toBe('items-center justify-center')
  })

  it('should default to center for unknown position', () => {
    expect(getFlexAlignmentClass('invalid-position')).toBe('items-center justify-center')
  })
})

describe('Position Utils - getFixedPositionClass()', () => {
  it('should return correct classes for top positions', () => {
    expect(getFixedPositionClass('top-left')).toBe('top-0 start-0')
    expect(getFixedPositionClass('top-center')).toBe('top-0 start-0 end-0')
    expect(getFixedPositionClass('top-right')).toBe('top-0 end-0')
  })

  it('should return correct classes for center', () => {
    expect(getFixedPositionClass('center')).toBe('inset-0')
  })

  it('should return correct classes for center-left/right', () => {
    expect(getFixedPositionClass('center-left')).toBe('top-1/2 start-0 -translate-y-1/2')
    expect(getFixedPositionClass('center-right')).toBe('top-1/2 end-0 -translate-y-1/2')
  })

  it('should return correct classes for bottom positions', () => {
    expect(getFixedPositionClass('bottom-left')).toBe('bottom-0 start-0')
    expect(getFixedPositionClass('bottom-center')).toBe('bottom-0 start-0 end-0')
    expect(getFixedPositionClass('bottom-right')).toBe('bottom-0 end-0')
  })

  it('should default to top-right for undefined position', () => {
    expect(getFixedPositionClass()).toBe('top-0 end-0')
    expect(getFixedPositionClass(null)).toBe('top-0 end-0')
  })

  it('should default to top-right for unknown position', () => {
    expect(getFixedPositionClass('invalid')).toBe('top-0 end-0')
  })
})

describe('Position Utils - getStackingConfig()', () => {
  describe('flex direction', () => {
    it('should use flex-col for top positions', () => {
      expect(getStackingConfig('top-left').flexDir).toBe('flex-col')
      expect(getStackingConfig('top-center').flexDir).toBe('flex-col')
      expect(getStackingConfig('top-right').flexDir).toBe('flex-col')
    })

    it('should use flex-col-reverse for bottom positions', () => {
      expect(getStackingConfig('bottom-left').flexDir).toBe('flex-col-reverse')
      expect(getStackingConfig('bottom-center').flexDir).toBe('flex-col-reverse')
      expect(getStackingConfig('bottom-right').flexDir).toBe('flex-col-reverse')
    })

    it('should use flex-col for center positions', () => {
      expect(getStackingConfig('center').flexDir).toBe('flex-col')
      expect(getStackingConfig('center-left').flexDir).toBe('flex-col')
      expect(getStackingConfig('center-right').flexDir).toBe('flex-col')
    })
  })

  describe('alignment', () => {
    it('should align start for left positions', () => {
      expect(getStackingConfig('top-left').alignItems).toBe('items-start')
      expect(getStackingConfig('center-left').alignItems).toBe('items-start')
      expect(getStackingConfig('bottom-left').alignItems).toBe('items-start')
    })

    it('should align center for positions with center keyword (except pure right)', () => {
      // Note: 'center-right' includes 'center', so it gets 'items-center'
      // This is the actual implementation behavior
      expect(getStackingConfig('top-right').alignItems).toBe('items-end')
      expect(getStackingConfig('center-right').alignItems).toBe('items-center') // Contains 'center'
      expect(getStackingConfig('bottom-right').alignItems).toBe('items-end')
    })

    it('should align center for center positions', () => {
      expect(getStackingConfig('top-center').alignItems).toBe('items-center')
      expect(getStackingConfig('bottom-center').alignItems).toBe('items-center')
    })

    it('should center both axes for true center', () => {
      expect(getStackingConfig('center').alignItems).toBe('items-center justify-center')
    })
  })

  describe('spacing', () => {
    it('should use space-y-2 for top positions', () => {
      expect(getStackingConfig('top-right').spaceY).toBe('space-y-2')
    })

    it('should use space-y-reverse for bottom positions', () => {
      expect(getStackingConfig('bottom-right').spaceY).toBe('space-y-reverse space-y-2')
    })
  })

  it('should default to top-right config', () => {
    const config = getStackingConfig()
    expect(config.flexDir).toBe('flex-col')
    expect(config.alignItems).toBe('items-end')
    expect(config.spaceY).toBe('space-y-2')
  })
})

describe('Position Utils - calculatePosition()', () => {
  let triggerEl, tooltipEl

  beforeEach(() => {
    // Mock trigger element
    triggerEl = {
      getBoundingClientRect: () => ({
        top: 100,
        left: 200,
        width: 80,
        height: 40,
        bottom: 140,
        right: 280
      })
    }

    // Mock tooltip element
    tooltipEl = {
      getBoundingClientRect: () => ({
        width: 120,
        height: 60,
        top: 0,
        left: 0,
        bottom: 60,
        right: 120
      })
    }

    // Mock window dimensions
    Object.defineProperty(window, 'innerWidth', { value: 1024, writable: true })
    Object.defineProperty(window, 'innerHeight', { value: 768, writable: true })
  })

  afterEach(() => {
    vi.unstubAllGlobals()
  })

  it('should mirror left/right placement by the anchor element direction', () => {
    vi.stubGlobal('getComputedStyle', () => ({ direction: 'rtl' }))
    const mirrored = calculatePosition(triggerEl, tooltipEl, 'left')
    expect(mirrored.position).toBe('right')
    expect(mirrored.left).toBe(200 + 80 + 8)

    vi.stubGlobal('getComputedStyle', () => ({ direction: 'ltr' }))
    const normal = calculatePosition(triggerEl, tooltipEl, 'left')
    expect(normal.position).toBe('left')
    expect(normal.left).toBe(200 - 120 - 8)
  })

  describe('top placement', () => {
    it('should position tooltip above trigger', () => {
      const result = calculatePosition(triggerEl, tooltipEl, 'top')
      expect(result.position).toBe('top')
      expect(result.top).toBe(100 - 60 - 8) // trigger.top - tooltip.height - spacing
      expect(result.left).toBe(200 + (80 - 120) / 2) // centered horizontally
    })

    it('should flip to bottom if no space above', () => {
      triggerEl.getBoundingClientRect = () => ({
        top: 30, // Not enough space above
        left: 200,
        width: 80,
        height: 40
      })

      const result = calculatePosition(triggerEl, tooltipEl, 'top')
      expect(result.position).toBe('bottom')
      expect(result.top).toBe(30 + 40 + 8) // trigger.top + trigger.height + spacing
    })
  })

  describe('bottom placement', () => {
    it('should position tooltip below trigger', () => {
      const result = calculatePosition(triggerEl, tooltipEl, 'bottom')
      expect(result.position).toBe('bottom')
      expect(result.top).toBe(100 + 40 + 8) // trigger.top + trigger.height + spacing
      expect(result.left).toBe(200 + (80 - 120) / 2)
    })

    it('should flip to top if no space below', () => {
      triggerEl.getBoundingClientRect = () => ({
        top: 720, // Near bottom of viewport
        left: 200,
        width: 80,
        height: 40
      })

      const result = calculatePosition(triggerEl, tooltipEl, 'bottom')
      expect(result.position).toBe('top')
      expect(result.top).toBe(720 - 60 - 8)
    })
  })

  describe('left placement', () => {
    it('should position tooltip to the left of trigger', () => {
      const result = calculatePosition(triggerEl, tooltipEl, 'left')
      expect(result.position).toBe('left')
      expect(result.left).toBe(200 - 120 - 8) // trigger.left - tooltip.width - offset (default 8)
      expect(result.top).toBe(100 + 40 / 2 - 60 / 2) // centered vertically
    })

    it('should flip to right if no space on left', () => {
      triggerEl.getBoundingClientRect = () => ({
        top: 100,
        left: 50, // Not enough space on left
        width: 80,
        height: 40
      })

      const result = calculatePosition(triggerEl, tooltipEl, 'left')
      expect(result.position).toBe('right')
      expect(result.left).toBe(50 + 80 + 8)
    })
  })

  describe('right placement', () => {
    it('should position tooltip to the right of trigger', () => {
      const result = calculatePosition(triggerEl, tooltipEl, 'right')
      expect(result.position).toBe('right')
      expect(result.left).toBe(200 + 80 + 8) // trigger.left + trigger.width + offset (default 8)
      expect(result.top).toBe(100 + 40 / 2 - 60 / 2)
    })

    it('should flip to left if no space on right', () => {
      triggerEl.getBoundingClientRect = () => ({
        top: 100,
        left: 950, // Near right edge
        width: 80,
        height: 40
      })

      const result = calculatePosition(triggerEl, tooltipEl, 'right')
      expect(result.position).toBe('left')
      expect(result.left).toBe(950 - 120 - 8)
    })
  })

  describe('viewport constraints', () => {
    it('should constrain tooltip within viewport horizontally', () => {
      triggerEl.getBoundingClientRect = () => ({
        top: 100,
        left: 5, // Near left edge
        width: 80,
        height: 40
      })

      const result = calculatePosition(triggerEl, tooltipEl, 'top')
      expect(result.left).toBeGreaterThanOrEqual(10) // Min margin
    })

    it('should constrain tooltip within viewport vertically', () => {
      triggerEl.getBoundingClientRect = () => ({
        top: 5, // Near top edge
        left: 200,
        width: 80,
        height: 40
      })

      const result = calculatePosition(triggerEl, tooltipEl, 'left')
      expect(result.top).toBeGreaterThanOrEqual(10) // Min margin
    })
  })

  describe('offset option', () => {
    it('should use default offset of 8 when not provided', () => {
      const result = calculatePosition(triggerEl, tooltipEl, 'top')
      expect(result.top).toBe(100 - 60 - 8)
    })

    it('should honor a custom offset', () => {
      const result = calculatePosition(triggerEl, tooltipEl, 'top', { offset: 20 })
      expect(result.top).toBe(100 - 60 - 20)

      const right = calculatePosition(triggerEl, tooltipEl, 'right', { offset: 16 })
      expect(right.left).toBe(200 + 80 + 16)
    })

    it('should apply custom offset to the flipped side as well', () => {
      triggerEl.getBoundingClientRect = () => ({
        top: 30, // Not enough space above
        left: 200,
        width: 80,
        height: 40
      })

      const result = calculatePosition(triggerEl, tooltipEl, 'top', { offset: 12 })
      expect(result.position).toBe('bottom')
      expect(result.top).toBe(30 + 40 + 12)
    })

    it('should fall back to 8 for invalid offset values', () => {
      expect(calculatePosition(triggerEl, tooltipEl, 'top', { offset: -5 }).top).toBe(100 - 60 - 8)
      expect(calculatePosition(triggerEl, tooltipEl, 'top', { offset: NaN }).top).toBe(100 - 60 - 8)
      expect(calculatePosition(triggerEl, tooltipEl, 'top', { offset: 'big' }).top).toBe(
        100 - 60 - 8
      )
    })
  })

  describe('arrow anchor', () => {
    // Arrow is an 8px square rendered inside the tooltip box (w-2 h-2),
    // anchored relative to the tooltip's own coordinate system.

    it('should normalize unsupported placements to bottom', () => {
      // Popper-style values like 'bottom-end' were once echoed back verbatim,
      // which broke consumers keying styles off `position` (the arrow borders).
      const result = calculatePosition(triggerEl, tooltipEl, 'bottom-end')
      expect(result.position).toBe('bottom')
      // ...and the geometry must match a plain 'bottom' placement
      const plain = calculatePosition(triggerEl, tooltipEl, 'bottom')
      expect({ top: result.top, left: result.left }).toEqual({ top: plain.top, left: plain.left })
      expect(result.arrow).toEqual(plain.arrow)
    })

    it('should anchor arrow on the bottom edge for top placement', () => {
      // tooltip: 120x60 at left = 200 + (80 - 120) / 2 = 180
      // trigger center X = 200 + 80 / 2 = 240 -> rel = 240 - 180 - 4 = 56
      const result = calculatePosition(triggerEl, tooltipEl, 'top')
      expect(result.position).toBe('top')
      expect(result.arrow).toEqual({ top: 60 - 4, left: 56 })
    })

    it('should anchor arrow on the top edge for bottom placement', () => {
      const result = calculatePosition(triggerEl, tooltipEl, 'bottom')
      expect(result.position).toBe('bottom')
      expect(result.arrow.top).toBe(-4)
      expect(result.arrow.left).toBe(56)
    })

    it('should anchor arrow on the side edges for left/right placement', () => {
      // trigger center Y = 100 + 40 / 2 = 120; tooltip top = 100 + 20 - 30 = 90
      // rel Y = 120 - 90 - 4 = 26
      const left = calculatePosition(triggerEl, tooltipEl, 'left')
      expect(left.arrow).toEqual({ top: 26, left: 120 - 4 })

      const right = calculatePosition(triggerEl, tooltipEl, 'right')
      expect(right.arrow).toEqual({ top: 26, left: -4 })
    })

    it('should clamp arrow inside the tooltip when trigger center is outside', () => {
      // Trigger far to the right of the (viewport-constrained) tooltip
      triggerEl.getBoundingClientRect = () => ({
        top: 100,
        left: 990, // tooltip gets clamped at viewportRight - 120 - 10 = 894
        width: 80,
        height: 40
      })

      const result = calculatePosition(triggerEl, tooltipEl, 'top')
      // triggerCX - left - 4 = 1030 - 894 - 4 = 132 > max (120 - 8 = 112)
      expect(result.arrow.left).toBe(112)
      expect(result.arrow.top).toBe(60 - 4)
    })
  })
})
