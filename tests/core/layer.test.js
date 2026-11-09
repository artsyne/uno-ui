import { describe, it, expect } from 'vitest'
import { Z_INDEX, getZIndex } from '../../src/core/layer'

describe('Layer Manager', () => {
  it('should define base layers with correct hierarchy', () => {
    // Verify relative order
    expect(Z_INDEX.DROPDOWN).toBeGreaterThan(Z_INDEX.BASE)
    expect(Z_INDEX.STICKY).toBeGreaterThan(Z_INDEX.DROPDOWN)
    expect(Z_INDEX.FIXED).toBeGreaterThan(Z_INDEX.STICKY)
    expect(Z_INDEX.BOX_BACKDROP).toBeGreaterThan(Z_INDEX.FIXED)
    expect(Z_INDEX.BOX).toBeGreaterThan(Z_INDEX.BOX_BACKDROP)
    expect(Z_INDEX.POPOVER).toBeGreaterThan(Z_INDEX.BOX)
    expect(Z_INDEX.NOTIFY).toBeGreaterThan(Z_INDEX.POPOVER)
    expect(Z_INDEX.TOOLTIP).toBeGreaterThan(Z_INDEX.NOTIFY)
  })

  it('should retrieve z-index by name (case insensitive)', () => {
    expect(getZIndex('box')).toBe(Z_INDEX.BOX)
    expect(getZIndex('BOX')).toBe(Z_INDEX.BOX)
    expect(getZIndex('Notify')).toBe(Z_INDEX.NOTIFY)
  })

  it('should return BASE for unknown layers', () => {
    expect(getZIndex('unknown')).toBe(Z_INDEX.BASE)
    expect(getZIndex(null)).toBe(Z_INDEX.BASE)
  })
})
