import { describe, it, expect } from 'vitest'
import { Z_INDEX, getZIndex } from '../../src/core/layer'

describe('Layer Manager', () => {
  it('should define base layers with correct hierarchy', () => {
    // Verify relative order
    expect(Z_INDEX.DROPDOWN).toBeGreaterThan(Z_INDEX.BASE)
    expect(Z_INDEX.STICKY).toBeGreaterThan(Z_INDEX.DROPDOWN)
    expect(Z_INDEX.FIXED).toBeGreaterThan(Z_INDEX.STICKY)
    expect(Z_INDEX.MODAL_BACKDROP).toBeGreaterThan(Z_INDEX.FIXED)
    expect(Z_INDEX.MODAL).toBeGreaterThan(Z_INDEX.MODAL_BACKDROP)
    expect(Z_INDEX.POPOVER).toBeGreaterThan(Z_INDEX.MODAL)
    expect(Z_INDEX.TOAST).toBeGreaterThan(Z_INDEX.POPOVER)
    expect(Z_INDEX.TOOLTIP).toBeGreaterThan(Z_INDEX.TOAST)
  })

  it('should retrieve z-index by name (case insensitive)', () => {
    expect(getZIndex('modal')).toBe(Z_INDEX.MODAL)
    expect(getZIndex('MODAL')).toBe(Z_INDEX.MODAL)
    expect(getZIndex('Toast')).toBe(Z_INDEX.TOAST)
  })

  it('should return BASE for unknown layers', () => {
    expect(getZIndex('unknown')).toBe(Z_INDEX.BASE)
    expect(getZIndex(null)).toBe(Z_INDEX.BASE)
  })
})
