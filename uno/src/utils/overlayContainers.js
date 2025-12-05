/**
 * Uno - Overlay Positioning Utils
 *
 * @module uno.core.overlayContainers
 * @version 0.1.0
 * @license MIT
 */

/**
 * Get CSS style for a fixed position
 * @param {string} position - Position name (e.g. 'top-right')
 * @returns {Object} CSS style object
 */
export function getPositionStyle(position) {
  const s = { position: 'fixed' }
  switch (position) {
    case 'top':
      s.top = '0px'; s.left = '0px'; s.right = '0px'
      break
    case 'bottom':
      s.bottom = '0px'; s.left = '0px'; s.right = '0px'
      break
    case 'top-right':
      s.top = '20px'; s.right = '20px'
      break
    case 'bottom-right':
      s.bottom = '12px'; s.right = '20px'
      break
    case 'top-left':
      s.top = '20px'; s.left = '20px'
      break
    case 'bottom-left':
      s.bottom = '12px'; s.left = '20px'
      break
    case 'center':
      s.top = '50%'; s.left = '50%'; s.transform = 'translate(-50%, -50%)'
      break
    case 'top-center':
      s.top = '0px'; s.left = '50%'; s.transform = 'translateX(-50%)'
      break
    case 'bottom-center':
      s.bottom = '0px'; s.left = '50%'; s.transform = 'translateX(-50%)'
      break
    default:
      s.top = '20px'; s.right = '20px'
  }
  return s
}
