/**
 * Uno - Positioning Utils
 *
 * @module uno.core.positioning
 * @version 0.1.0
 * @license MIT
 */

/**
 * Calculate tooltip position relative to trigger
 * @param {HTMLElement} triggerEl - Trigger element
 * @param {HTMLElement} tooltipEl - Tooltip element
 * @param {string} preferredPosition - Preferred position (top/bottom/left/right)
 * @returns {Object} { top, left, position }
 */
export function calculatePosition(triggerEl, tooltipEl, preferredPosition) {
  const triggerRect = triggerEl.getBoundingClientRect();
  const tooltipRect = tooltipEl.getBoundingClientRect();
  
  const containerWidth = tooltipRect.width;
  const containerHeight = tooltipRect.height;
  const triggerWidth = triggerRect.width;
  const triggerHeight = triggerRect.height;
  
  // Use fixed positioning, relative to viewport
  const pos = { top: triggerRect.top, left: triggerRect.left };
  
  const viewportWidth = window.innerWidth;
  const viewportHeight = window.innerHeight;
  const viewportTop = 0;
  const viewportLeft = 0;
  const viewportBottom = viewportHeight;
  const viewportRight = viewportWidth;
  
  const constrainToViewport = (position, isHorizontal) => {
      if (isHorizontal) {
        if (position + containerWidth > viewportRight) {
          return viewportRight - containerWidth - 10;
        }
        if (position < viewportLeft) {
          return viewportLeft + 10;
        }
      } else {
        if (position + containerHeight > viewportBottom) {
          return viewportBottom - containerHeight - 10;
        }
        if (position < viewportTop) {
          return viewportTop + 10;
        }
      }
      return position;
  };
  
  let top, left;
  let finalPosition = preferredPosition;
  
  switch (preferredPosition) {
      case 'left':
        left = pos.left - containerWidth - 5;
        if (left < viewportLeft) {
            left = pos.left + triggerWidth + 5;
            finalPosition = 'right';
        }
        top = constrainToViewport(pos.top + triggerHeight / 2 - containerHeight / 2, false);
        break;
      case 'top':
        top = pos.top - containerHeight - 8; // Add some spacing
        if (top < viewportTop) {
            top = pos.top + triggerHeight + 8;
            finalPosition = 'bottom';
        }
        left = pos.left + (triggerWidth - containerWidth) / 2;
        left = constrainToViewport(left, true);
        break;
      case 'right':
        left = pos.left + triggerWidth + 5;
        if (left + containerWidth > viewportRight) {
            left = pos.left - containerWidth - 5;
            finalPosition = 'left';
        }
        top = constrainToViewport(pos.top + triggerHeight / 2 - containerHeight / 2, false);
        break;
      case 'bottom':
      default:
        top = pos.top + triggerHeight + 8;
        if (top + containerHeight > viewportBottom) {
            top = pos.top - containerHeight - 8;
            finalPosition = 'top';
        }
        left = pos.left + (triggerWidth - containerWidth) / 2;
        left = constrainToViewport(left, true);
        break;
  }
  
  return { top, left, position: finalPosition };
}
