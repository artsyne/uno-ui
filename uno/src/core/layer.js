/**
 * Uno - Layer Manager (Z-Index)
 *
 * Defines z-index constants and utilities to manage stacking context.
 *
 * @module uno.core.layer
 */

export const Z_INDEX = {
  BASE: 0,
  DROPDOWN: 1000,
  STICKY: 1100,
  FIXED: 1200,
  BOX_BACKDROP: 1300,
  BOX: 1400,
  POPOVER: 1500,
  NOTIFY: 1600,
  TOOLTIP: 1700,
  MAX: 9999
}

/**
 * Get z-index value for a named layer.
 *
 * @param {string} layer - The layer name (e.g., 'modal', 'toast')
 * @returns {number} The z-index value
 */
export const getZIndex = (layer) => {
  if (!layer) return Z_INDEX.BASE
  const key = layer.toUpperCase()
  return Z_INDEX[key] !== undefined ? Z_INDEX[key] : Z_INDEX.BASE
}
