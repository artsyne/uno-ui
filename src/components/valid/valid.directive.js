import { getResponseInfo } from '../../utils/htmx-compat'
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

        // --- Length Rule Features (Counter & StopOnMax) ---
        const validStore = Alpine.store('uno').valid
        if (validStore && validStore.normalizeRules) {
          const rules = validStore.normalizeRules(value)
          // Find any length-related rule
          const lengthRule = rules.find((r) => r.type === 'length' || r.maxLength !== undefined || r.max !== undefined)

          if (lengthRule) {
            const max = lengthRule.max !== undefined ? lengthRule.max : lengthRule.maxLength

            // 1. StopOnMax Feature
            if (lengthRule.stopOnMax === true && max !== undefined) {
              el.setAttribute('maxlength', max)
            }

            // 2. Counter Feature
            if (lengthRule.showCounter === true && max !== undefined) {
              const counterId = `${fieldName}-counter`
              let counterEl = document.getElementById(counterId)

              // Ensure parent is relative for absolute positioning
              const parent = el.parentElement
              if (parent && window.getComputedStyle(parent).position === 'static' && parent.tagName !== 'FORM' && parent.tagName !== 'BODY') {
                parent.style.position = 'relative'
              }

              if (!counterEl) {
                counterEl = document.createElement('span')
                counterEl.id = counterId
                // Default styling: bottom-right inside the input container
                // Adjust margins to not overlap with border/padding
                counterEl.className = 'absolute bottom-1 right-2 text-xs text-gray-400 pointer-events-none transition-colors duration-200'
                el.after(counterEl)
              }

              const updateCounter = () => {
                const current = el.value.length
                counterEl.textContent = `${current}/${max}`

                if (current > max) {
                  counterEl.classList.add('text-red-500')
                  counterEl.classList.remove('text-gray-400')
                } else {
                  counterEl.classList.remove('text-red-500')
                  counterEl.classList.add('text-gray-400')
                }
              }

              // Initial update
              updateCounter()

              // Bind to input
              el.addEventListener('input', updateCounter)

              // Cleanup listener when directive is destroyed
              // Note: We don't remove the element itself automatically because it might be reused or handled by parent,
              // but cleaner is to remove it if we created it.
              // However, Alpine cleanup is called on element removal.
              cleanup(() => {
                el.removeEventListener('input', updateCounter)
                if (counterEl && counterEl.parentNode) {
                  counterEl.parentNode.removeChild(counterEl)
                }
              })
            }
          }
        }
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
      const info = getResponseInfo(evt)
      if (info && info.status === 422) {
        try {
          const response = JSON.parse(info.text)
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
