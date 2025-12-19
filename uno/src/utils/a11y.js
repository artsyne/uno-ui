/**
 * Uno - Accessibility Utils
 *
 * @module uno.utils.a11y
 */

const FOCUSABLE_SELECTOR = [
  'a[href]',
  'area[href]',
  'input:not([disabled])',
  'select:not([disabled])',
  'textarea:not([disabled])',
  'button:not([disabled])',
  'iframe',
  'object',
  'embed',
  '[contenteditable]',
  '[tabindex]:not([tabindex^="-"])'
].join(',')

/**
 * Trap focus within an element
 * @param {HTMLElement} el - The element to trap focus in
 * @param {Object} options
 * @returns {Function} cleanup function
 */
export function trapFocus(el, options = {}) {
  void options
  const focusables = () => Array.from(el.querySelectorAll(FOCUSABLE_SELECTOR))
    .filter(node => node.offsetParent !== null) // Visible only

  const firstFocusable = () => focusables()[0]
  const lastFocusable = () => {
    const list = focusables()
    return list[list.length - 1]
  }

  const handleKeyDown = (e) => {
    if (e.key !== 'Tab') return

    if (e.shiftKey) {
      if (document.activeElement === firstFocusable()) {
        e.preventDefault()
        const last = lastFocusable()
        if (last) last.focus()
      }
    } else {
      if (document.activeElement === lastFocusable()) {
        e.preventDefault()
        const first = firstFocusable()
        if (first) first.focus()
      }
    }
  }

  // Initial focus
  setTimeout(() => {
    const first = firstFocusable()
    if (first) first.focus()
  }, 50)

  el.addEventListener('keydown', handleKeyDown)

  return () => {
    el.removeEventListener('keydown', handleKeyDown)
    // Restore focus? usually handled by caller or browser history
  }
}
