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

    // Auto-generate error element
    const validConfig = Uno.getConfig('valid')
    const uiConfig = Uno.options.ui?.valid || {}

    // Check if autoError is explicitly disabled in config
    if (validConfig.autoError !== false) {
      const errorKey = `$store.uno.valid.errors['${fieldName}']`
      const errorSelector = `[x-text="${errorKey.replace(/"/g, '\\"')}"]`

      // Search within the form context if possible
      const form = el.form || el.closest('form')
      const existingError = form
        ? form.querySelector(errorSelector)
        : document.querySelector(errorSelector)

      if (!existingError) {
        let errorEl
        if (typeof validConfig.createErrorElement === 'function') {
          errorEl = validConfig.createErrorElement(fieldName, errorKey)
        } else {
          errorEl = document.createElement('div')
          errorEl.setAttribute('x-text', errorKey)
          errorEl.setAttribute('x-show', errorKey)

          // Use helperText class from UI config or fallback
          const cls = validConfig.errorClass || uiConfig.helperText || 'text-xs text-red-500 mt-1'
          errorEl.className = cls
        }

        if (errorEl) {
          el.after(errorEl)
          // Initialize Alpine on the new element
          Alpine.initTree(errorEl)
        }
      }
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
