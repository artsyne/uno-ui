/**
 * Uno - Position Utils
 * Handles screen-relative positioning logic (Box, Notify) and element anchoring (Tooltip).
 *
 * @module uno.utils.position
 */

// Available positions
export const POSITIONS = [
  'top-left',
  'top-center',
  'top-right',
  'center-left',
  'center',
  'center-right',
  'bottom-left',
  'bottom-center',
  'bottom-right'
]

/**
 * Get flexbox alignment classes for a given position.
 * Used by Box component (centering content in a full-screen overlay).
 *
 * @param {string} position - The position name (e.g., 'top-right')
 * @returns {string} Tailwind CSS classes for flex alignment
 */
export const getFlexAlignmentClass = (position) => {
  const pos = position || 'center'
  const map = {
    center: 'items-center justify-center',
    'top-center': 'items-start justify-center',
    'top-left': 'items-start justify-start',
    'top-right': 'items-start justify-end',
    'bottom-center': 'items-end justify-center',
    'bottom-left': 'items-end justify-start',
    'bottom-right': 'items-end justify-end',
    'center-left': 'items-center justify-start',
    'center-right': 'items-center justify-end'
  }
  return map[pos] || map['center']
}

/**
 * Get fixed positioning classes for a container.
 * Used by Notify component (pinning containers to screen edges).
 *
 * @param {string} position - The position name
 * @returns {string} Tailwind CSS classes for fixed positioning
 */
export const getFixedPositionClass = (position) => {
  const pos = position || 'top-right'
  const map = {
    'top-right': 'top-0 right-0',
    'top-left': 'top-0 left-0',
    'top-center': 'top-0 left-0 right-0', // Full width for center alignment
    center: 'inset-0', // Full screen for center
    'center-right': 'top-1/2 right-0 -translate-y-1/2',
    'center-left': 'top-1/2 left-0 -translate-y-1/2',
    'bottom-right': 'bottom-0 right-0',
    'bottom-left': 'bottom-0 left-0',
    'bottom-center': 'bottom-0 left-0 right-0'
  }
  return map[pos] || map['top-right']
}

/**
 * Get stacking configuration for lists of items.
 * Used by Notify component to determine flex direction and item alignment.
 *
 * @param {string} position - The position name
 * @returns {Object} { flexDir, alignItems, spaceY }
 */
export const getStackingConfig = (position) => {
  const pos = position || 'top-right'

  // Default: Top -> Bottom (flex-col), Align Start/End based on side
  let flexDir = 'flex-col'
  let alignItems = 'items-end' // default for right
  let spaceY = 'space-y-2'

  if (pos.includes('bottom')) {
    flexDir = 'flex-col-reverse'
    spaceY = 'space-y-reverse space-y-2'
  }

  if (pos.includes('left')) {
    alignItems = 'items-start'
  } else if (pos.includes('center')) {
    if (pos === 'center') {
      alignItems = 'items-center justify-center'
    } else {
      alignItems = 'items-center'
    }
  } else {
    // right
    alignItems = 'items-end'
  }

  return { flexDir, alignItems, spaceY }
}

/**
 * Calculate tooltip position relative to trigger
 * @param {HTMLElement} triggerEl - Trigger element
 * @param {HTMLElement} tooltipEl - Tooltip element
 * @param {string} preferredPosition - Preferred position (top/bottom/left/right)
 * @returns {Object} { top, left, position }
 */
export function calculatePosition(triggerEl, tooltipEl, preferredPosition) {
  const triggerRect = triggerEl.getBoundingClientRect()
  const tooltipRect = tooltipEl.getBoundingClientRect()

  const containerWidth = tooltipRect.width
  const containerHeight = tooltipRect.height
  const triggerWidth = triggerRect.width
  const triggerHeight = triggerRect.height

  // Use fixed positioning, relative to viewport
  const pos = { top: triggerRect.top, left: triggerRect.left }

  const viewportWidth = window.innerWidth
  const viewportHeight = window.innerHeight
  const viewportTop = 0
  const viewportLeft = 0
  const viewportBottom = viewportHeight
  const viewportRight = viewportWidth

  const constrainToViewport = (position, isHorizontal) => {
    if (isHorizontal) {
      if (position + containerWidth > viewportRight) {
        return viewportRight - containerWidth - 10
      }
      if (position < viewportLeft) {
        return viewportLeft + 10
      }
    } else {
      if (position + containerHeight > viewportBottom) {
        return viewportBottom - containerHeight - 10
      }
      if (position < viewportTop) {
        return viewportTop + 10
      }
    }
    return position
  }

  let top, left
  let finalPosition = preferredPosition

  switch (preferredPosition) {
    case 'left':
      left = pos.left - containerWidth - 5
      if (left < viewportLeft) {
        left = pos.left + triggerWidth + 5
        finalPosition = 'right'
      }
      top = constrainToViewport(pos.top + triggerHeight / 2 - containerHeight / 2, false)
      break
    case 'top':
      top = pos.top - containerHeight - 8 // Add some spacing
      if (top < viewportTop) {
        top = pos.top + triggerHeight + 8
        finalPosition = 'bottom'
      }
      left = pos.left + (triggerWidth - containerWidth) / 2
      left = constrainToViewport(left, true)
      break
    case 'right':
      left = pos.left + triggerWidth + 5
      if (left + containerWidth > viewportRight) {
        left = pos.left - containerWidth - 5
        finalPosition = 'left'
      }
      top = constrainToViewport(pos.top + triggerHeight / 2 - containerHeight / 2, false)
      break
    case 'bottom':
    default:
      top = pos.top + triggerHeight + 8
      if (top + containerHeight > viewportBottom) {
        top = pos.top - containerHeight - 8
        finalPosition = 'top'
      }
      left = pos.left + (triggerWidth - containerWidth) / 2
      left = constrainToViewport(left, true)
      break
  }

  return { top, left, position: finalPosition }
}
