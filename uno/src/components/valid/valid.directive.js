import { log } from '../../utils/log'

/**
 * Registers the `x-uno-valid` directive.
 *
 * @param {Object} Alpine - The Alpine.js instance.
 */
export const registerValidDirectives = (Alpine) => {
  /**
   * x-uno-valid directive
   *
   * Usage:
   *  <input x-uno-valid="['required', 'email']">
   *  <input x-uno-valid="{ type: 'length', min: 5 }">
   */
  Alpine.directive('uno-valid', (el, { expression }, { evaluateLater, cleanup }) => {
    // If used on a form, it acts as a marker for the global HTMX listener
    if (el.tagName === 'FORM') return

    const fieldName = el.name || el.id
    if (!fieldName) {
      log.warn('[Uno Valid] Element has no name or id:', el)
      return
    }

    const getExpression = expression ? evaluateLater(expression) : () => {}

    // Register rules in store
    const register = (config) => {
      Alpine.store('uno').valid.addRules(fieldName, config)
    }

    if (expression) {
      getExpression((value) => {
        // We expect an array or object here, passed directly to store normalization
        register(value)
      })
    } else {
      // Empty directive x-uno-valid means no rules initially?
      // Or maybe we shouldn't allow empty.
      // For now, register empty config.
      register([])
    }

    // Event Listeners
    const onInput = () => {
      Alpine.store('uno').valid.validate(fieldName, el.value, 'input')
    }

    const onBlur = () => {
      Alpine.store('uno').valid.validate(fieldName, el.value, 'blur')
    }

    el.addEventListener('input', onInput)
    el.addEventListener('blur', onBlur)

    cleanup(() => {
      el.removeEventListener('input', onInput)
      el.removeEventListener('blur', onBlur)
    })
  })

  // ============================================================
  // HTMX Integration & Global Listeners
  // ============================================================
  if (typeof document !== 'undefined') {
    /**
     * Handles htmx:confirm event.
     * Validates the form before HTMX issues a request.
     */
    document.body.addEventListener('htmx:confirm', async (evt) => {
      const el = evt.detail.elt
      const form = el.tagName === 'FORM' ? el : el.closest('form')

      if (form && form.hasAttribute('x-uno-valid')) {
        evt.preventDefault()

        const isValid = await Alpine.store('uno').valid.validateForm(form)

        if (isValid) {
          evt.detail.issueRequest(true)
        }
      }
    })

    /**
     * Handles htmx:afterRequest event.
     * Processes server-side validation errors (422 status).
     */
    document.body.addEventListener('htmx:afterRequest', (evt) => {
      const xhr = evt.detail.xhr
      if (xhr.status === 422) {
        try {
          const response = JSON.parse(xhr.responseText)
          if (response.errors) {
            Alpine.store('uno').valid.showErrors(response.errors)
            log.info('[Valid] Server validation errors applied', response.errors)
          }
        } catch {
          // Ignore non-JSON responses
        }
      }
    })
  }
}
