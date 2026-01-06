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
 * Calculate tooltip position relative to trigger element with automatic flipping.
 *
 * This function implements an intelligent positioning algorithm that:
 * 1. Calculates the ideal position based on the preferred placement
 * 2. Checks if the tooltip fits within the viewport boundaries
 * 3. Automatically flips to the opposite side if there's not enough space
 * 4. Constrains the position to prevent overflow from the viewport
 *
 * Algorithm Details:
 * - For 'top' placement: Position above the trigger, flip to 'bottom' if no space
 * - For 'bottom' placement: Position below the trigger, flip to 'top' if no space
 * - For 'left' placement: Position to the left, flip to 'right' if no space
 * - For 'right' placement: Position to the right, flip to 'left' if no space
 * - Always centers the tooltip relative to the trigger when possible
 * - Adds padding (5px horizontal, 8px vertical) to prevent edge touching
 *
 * Edge Cases Handled:
 * - Tooltip wider than trigger: Centers tooltip over trigger
 * - Near viewport edge: Shifts tooltip away from edge (10px minimum padding)
 * - Insufficient space on preferred side: Flips to opposite side
 * - Insufficient space on both sides: Shows on preferred side with viewport constraint
 *
 * @param {HTMLElement} triggerEl - The trigger element that activates the tooltip
 * @param {HTMLElement} tooltipEl - The tooltip element to be positioned
 * @param {string} preferredPosition - Desired placement ('top'|'bottom'|'left'|'right')
 * @returns {Object} Position object with properties:
 *   - {number} top - Top coordinate in pixels (viewport-relative)
 *   - {number} left - Left coordinate in pixels (viewport-relative)
 *   - {string} position - Actual position used (may differ from preferred if flipped)
 *
 * @example
 * const result = calculatePosition(button, tooltip, 'top')
 * tooltip.style.top = `${result.top}px`
 * tooltip.style.left = `${result.left}px`
 * console.log(result.position) // 'top' or 'bottom' if flipped
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

  /**
   * Constrain a position coordinate to stay within viewport boundaries.
   *
   * This helper function prevents tooltips from overflowing the visible area by:
   * - Checking if the tooltip would extend beyond the viewport edge
   * - Pulling it back within bounds if needed
   * - Adding a 10px safety margin from viewport edges
   *
   * @param {number} position - The proposed coordinate value
   * @param {boolean} isHorizontal - True for left/right (horizontal), false for top/bottom (vertical)
   * @returns {number} Constrained coordinate value
   *
   * @private
   */
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
