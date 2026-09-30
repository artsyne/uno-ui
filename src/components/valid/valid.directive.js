import { getResponseInfo } from '../../utils/htmx-compat'
import { log } from '../../utils/log'

/**
 * Marker attribute identifying plugin-managed counter containing blocks.
 */
const COUNTER_WRAPPER_ATTR = 'data-uno-counter-wrapper'

/**
 * Displays whose children take part in flex/grid placement. Reparenting such a
 * child into a wrapper div would change its outer layout role, so those
 * parents keep the legacy parent-relative anchoring instead.
 */
const isFlexOrGridDisplay = (display) =>
  display === 'flex' || display === 'inline-flex' || display === 'grid' || display === 'inline-grid'

/**
 * Ensures the input owns a tight positioned containing block for the counter.
 *
 * The input is wrapped in a plugin-managed `position: relative` div: the
 * wrapper is a plain block box (inline controls get an inline-block wrapper to
 * keep flowing with surrounding text) and the input's margins collapse through
 * it, so the outer layout stays unchanged while `absolute bottom-* right-*`
 * descendants anchor to the input box itself instead of an arbitrary ancestor.
 *
 * @param {HTMLElement} el - The validated input/textarea.
 * @returns {HTMLElement|null} The containing block to host the counter.
 */
const ensureCounterWrapper = (el) => {
  const parent = el.parentElement
  if (!parent) return null
  if (parent.hasAttribute(COUNTER_WRAPPER_ATTR)) return parent

  const parentStyle = window.getComputedStyle(parent)
  if (isFlexOrGridDisplay(parentStyle.display)) {
    // Legacy fallback: promote the parent itself rather than changing the
    // input's role as a flex/grid item.
    if (parentStyle.position === 'static' && parent.tagName !== 'FORM' && parent.tagName !== 'BODY') {
      parent.style.position = 'relative'
    }
    return parent
  }

  const wrapper = el.ownerDocument.createElement('div')
  wrapper.setAttribute(COUNTER_WRAPPER_ATTR, '')
  wrapper.style.position = 'relative'
  const display = window.getComputedStyle(el).display
  if (display === 'inline' || display === 'inline-block') {
    wrapper.style.display = 'inline-block'
  }
  parent.insertBefore(wrapper, el)
  wrapper.appendChild(el)
  return wrapper
}

/**
 * Restores the original DOM shape on directive teardown: moves the input back
 * before the managed wrapper and drops the wrapper.
 */
const releaseCounterWrapper = (el, wrapper) => {
  if (!wrapper || !wrapper.hasAttribute(COUNTER_WRAPPER_ATTR)) return
  const outer = wrapper.parentNode
  if (!outer) return
  outer.insertBefore(el, wrapper)
  outer.removeChild(wrapper)
}

/**
 * Builds the stable error element id used for `aria-describedby` wiring.
 *
 * @param {string} fieldName - The validated field's name/id.
 * @returns {string} The deterministic error element id.
 */
const errorIdFor = (fieldName) => `${fieldName}-error`

/**
 * Exposes the error element to assistive technology (WCAG 3.3.1 / 4.1.3):
 * gives it a stable id so the input can reference it and an `alert` role so
 * screen readers announce the message when it appears. Author-provided id or
 * role on a hand-written error element are preserved.
 *
 * @param {HTMLElement} errorEl - The error element (auto-generated or existing).
 * @param {string} fieldName - The validated field's name/id.
 */
const ensureErrorA11y = (errorEl, fieldName) => {
  if (!errorEl.id) errorEl.id = errorIdFor(fieldName)
  if (!errorEl.hasAttribute('role')) errorEl.setAttribute('role', 'alert')
}

/**
 * Points the input's `aria-describedby` at its error element, preserving any
 * existing describedby targets (e.g. help text) rather than clobbering them.
 *
 * @param {HTMLElement} input - The validated input/textarea.
 * @param {HTMLElement} errorEl - The error element to reference.
 */
const linkDescribedBy = (input, errorEl) => {
  if (!errorEl.id) return
  const tokens = (input.getAttribute('aria-describedby') || '')
    .split(/\s+/)
    .filter(Boolean)
  if (!tokens.includes(errorEl.id)) tokens.push(errorEl.id)
  input.setAttribute('aria-describedby', tokens.join(' '))
}

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
          ensureErrorA11y(errorEl, fieldName)
          el.after(errorEl)
          // Initialize Alpine on the new element
          Alpine.initTree(errorEl)
          linkDescribedBy(el, errorEl)
        }
      } else {
        // A hand-written error element still needs the a11y wiring so screen
        // readers announce it and the input exposes the association.
        ensureErrorA11y(existingError, fieldName)
        linkDescribedBy(el, existingError)
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
          const lengthRule = rules.find(
            (r) => r.type === 'length' || r.maxLength !== undefined || r.max !== undefined
          )

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

              // The counter is absolutely positioned, so its containing block is the
              // nearest positioned ancestor. Anchoring it to the input's parent is
              // unsafe: field containers usually carry sibling flow content (labels,
              // help text, status rows) whose height moves the anchor away from the
              // input box. Wrap the input in a plugin-managed relative container so
              // the counter stays pinned to the input's bottom-right corner
              // regardless of the surrounding form markup.
              const wrapper = ensureCounterWrapper(el)

              if (!counterEl) {
                counterEl = document.createElement('span')
                counterEl.id = counterId
                // Default styling: bottom-right inside the input container
                // Adjust margins to not overlap with border/padding
                counterEl.className =
                  'absolute bottom-1 right-2 text-xs text-gray-400 pointer-events-none transition-colors duration-200'
              }
              // (Re)attach into the containing block: a counter restored from a
              // previous init must not stay outside the wrapper
              if (counterEl.parentElement !== wrapper) {
                el.after(counterEl)
              }

              // P1-B6: expose the counter to screen readers via aria-describedby
              // (same association pattern as the error element)
              linkDescribedBy(el, counterEl)

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
                releaseCounterWrapper(el, wrapper)
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
     *
     * A form is intercepted when either the form itself carries the
     * `x-uno-valid` marker or any descendant declares validation rules via
     * `x-uno-valid`. The descendant check is required because a form whose
     * fields declare rules but whose <form> element is missing the marker
     * would otherwise bypass validation entirely (request sent unvalidated,
     * no error hints rendered).
     */
    document.body.addEventListener('htmx:confirm', async (evt) => {
      const el = evt.detail.elt
      const form = el.tagName === 'FORM' ? el : el.closest('form')

      if (form && (form.hasAttribute('x-uno-valid') || form.querySelector('[x-uno-valid]'))) {
        evt.preventDefault()

        try {
          const isValid = await Alpine.store('uno').valid.validateForm(form)

          if (isValid) {
            evt.detail.issueRequest(true)
          }
        } catch (err) {
          // Fail closed: an unexpected validation failure must block the
          // request rather than let an unvalidated submission through.
          log.error('[Valid] Form validation failed unexpectedly, request blocked', err)
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
