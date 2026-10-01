/**
 * Uno - Accessibility Utils
 *
 * @module uno.utils.a11y
 */

export const FOCUSABLE_SELECTOR = [
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
  const { initialFocus = true } = options
  const isVisible = (node) => {
    if (!(node instanceof HTMLElement)) return false
    const style = window.getComputedStyle(node)
    if (style.display === 'none' || style.visibility === 'hidden') return false
    return node.getClientRects().length > 0
  }

  const focusables = () =>
    Array.from(el.querySelectorAll(FOCUSABLE_SELECTOR)).filter((node) => isVisible(node))

  const firstFocusable = () => focusables()[0]

  const handleKeyDown = (e) => {
    if (e.key !== 'Tab') return

    const list = focusables()
    if (!list.length) return
    const first = list[0]
    const last = list[list.length - 1]
    const active = document.activeElement

    // Wrap when focus is at the boundary OR when the active element (e.g. a
    // tabindex="-1" neutral target such as the dialog title) has no focusable
    // element after/before it in DOM order — otherwise the browser's native
    // tab order would escape the container entirely.
    const hasNext = list.some(
      (node) => active.compareDocumentPosition(node) & Node.DOCUMENT_POSITION_FOLLOWING
    )
    const hasPrev = list.some(
      (node) => active.compareDocumentPosition(node) & Node.DOCUMENT_POSITION_PRECEDING
    )

    if (e.shiftKey) {
      if (active === first || !hasPrev) {
        e.preventDefault()
        last.focus()
      }
    } else {
      if (active === last || !hasNext) {
        e.preventDefault()
        first.focus()
      }
    }
  }

  // Initial focus
  const focusFirstIfNeeded = () => {
    if (el.contains(document.activeElement)) return
    const first = firstFocusable()
    if (first) first.focus()
  }

  if (initialFocus) {
    if (typeof requestAnimationFrame === 'function') {
      requestAnimationFrame(() => requestAnimationFrame(focusFirstIfNeeded))
    } else {
      setTimeout(focusFirstIfNeeded, 0)
    }
  }

  el.addEventListener('keydown', handleKeyDown)

  return () => {
    el.removeEventListener('keydown', handleKeyDown)
    // Restore focus? usually handled by caller or browser history
  }
}
