import { trapFocus } from '../../utils/a11y'

export const registerBoxDirectives = (Alpine) => {
  // ============================================================
  // Focus Trap Directive
  // ============================================================
  Alpine.directive('uno-trap', (el, { expression: _expression }, { cleanup }) => {
    void _expression
    const cleanupTrap = trapFocus(el)
    cleanup(cleanupTrap)
  })

  // ============================================================
  // Box Directive (x-uno-box)
  // Declarative interface for triggering box dialogs
  // ============================================================
  Alpine.directive('uno-box', (el, { expression }, { evaluateLater, cleanup }) => {
    if (!expression) {
      console.warn('[Uno] x-uno-box directive requires a configuration object.')
      return
    }

    const getExpression = evaluateLater(expression)

    /**
     * Click handler for directive
     * @param {Event} e - Click event
     */
    const onClick = (e) => {
      e.preventDefault()
      // Access store via uno namespace
      const store = Alpine.store('uno')?.box
      if (!store) {
        console.warn('[Uno] Box store not found. Is Uno initialized?')
        return
      }

      getExpression((evaluated) => {
        if (evaluated && typeof evaluated === 'object') {
          store.open(evaluated)
        } else {
          console.warn('[Uno] x-uno-box directive value must be an object', evaluated)
        }
      })
    }

    el.addEventListener('click', onClick)
    cleanup(() => el.removeEventListener('click', onClick))
  })
}
