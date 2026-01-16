import { log } from '../../utils/log'

export const registerValidDirectives = (Alpine) => {
  Alpine.directive('uno-valid', (el, { expression }, { evaluateLater, cleanup }) => {
    // If used on a form, it acts as a marker for the global HTMX listener
    if (el.tagName === 'FORM') return

    const fieldName = el.name || el.id
    if (!fieldName) {
      log.warn('[Uno Valid] Element has no name or id:', el)
      return
    }

    const getExpression = expression ? evaluateLater(expression) : () => {}

    // Parse configuration from attributes
    const buildFromAttrs = () => ({
      // Standard HTML5
      required: el.hasAttribute('required'),
      minLength: el.getAttribute('minlength') || null,
      maxLength: el.getAttribute('maxlength') || null,
      min: el.getAttribute('min') || null,
      max: el.getAttribute('max') || null,
      pattern: el.getAttribute('pattern') || null,
      type: el.getAttribute('type') || null
    })

    const initialAttrs = buildFromAttrs()

    // Register rules in store
    const register = (config) => {
      Alpine.store('uno').valid.addRules(fieldName, { ...initialAttrs, ...config })
    }

    if (expression) {
      getExpression((value) => {
        const config = typeof value === 'object' ? value : {}
        register(config)
      })
    } else {
      register({})
    }

    // Event Listeners
    const onInput = () => {
      const rules = Alpine.store('uno').valid.rules[fieldName]
      // Default to realtime unless explicitly disabled
      if (rules && rules.realtime !== false) {
        Alpine.store('uno').valid.validate(fieldName, el.value)
      }
    }

    const onBlur = () => {
      // Always validate on blur, regardless of realtime setting
      Alpine.store('uno').valid.validate(fieldName, el.value)
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
    document.body.addEventListener('htmx:confirm', (evt) => {
      const el = evt.detail.elt
      const form = el.tagName === 'FORM' ? el : el.closest('form')

      if (form && form.hasAttribute('x-uno-valid')) {
        const isValid = Alpine.store('uno').valid.validateForm(form)

        if (!isValid) {
          evt.preventDefault()
          evt.stopPropagation() // Stop HTMX
          log.info('[Valid] Form validation failed. Request blocked.')

          form.dispatchEvent(
            new CustomEvent('uno:validation-failed', {
              bubbles: true,
              detail: { errors: Alpine.store('uno').valid.errors }
            })
          )
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
