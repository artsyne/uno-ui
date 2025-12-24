import { VALIDATORS } from '../utils/validation'
import { Uno } from '../core/uno'

/**
 * Uno - Valid Component (Refactored for v0.5.0)
 *
 * @module uno.components.valid
 * @description HTMX-first validation interceptor with Alpine.js integration.
 */

export default (Alpine) => {
  // Global Store for Validation Errors (Reactive)
  Alpine.store('unoValid', {
    errors: {},
    setError(name, msg) {
      this.errors[name] = msg
    },
    clearError(name) {
      delete this.errors[name]
    },
    getMessage(name) {
      return this.errors[name] || ''
    }
  })

  // Global Helper for x-text
  // Usage: x-text="$uno.valid.getMessage('username')"
  if (typeof window !== 'undefined') {
    if (!window.$uno) window.$uno = {}
    window.$uno.valid = {
      getMessage: (name) => Alpine.store('unoValid').getMessage(name)
    }
  }

  // Initialize interceptor once
  if (typeof document !== 'undefined') {
    // Intercept HTMX requests
    document.body.addEventListener('htmx:confirm', (evt) => {
      const el = evt.detail.elt
      const form = el.tagName === 'FORM' ? el : el.closest('form')

      // Only intercept if form opts-in via data-uno-valid
      if (form && form.hasAttribute('data-uno-valid')) {
        if (!form.checkValidity()) {
          evt.preventDefault()
          evt.stopPropagation()
          // Trigger validation to show errors
          validateForm(form)
          form.dispatchEvent(new CustomEvent('uno:validation-failed', { bubbles: true }))
        }
      }
    })

    // Input listener for realtime validation (aggressive feedback)
    document.body.addEventListener('input', (evt) => {
      const input = evt.target
      if (input.form && input.form.hasAttribute('data-uno-valid')) {
        // If input is currently invalid (aria-invalid=true), re-validate to clear error
        if (input.getAttribute('aria-invalid') === 'true') {
          validateInput(input)
        }

        // Matchers logic (e.g. password confirm)
        if (input.name) {
          const matchers = input.form.querySelectorAll(`[data-uno-match="${input.name}"]`)
          matchers.forEach((matcher) => {
            if (matcher.getAttribute('aria-invalid') === 'true' || matcher.value) {
              validateInput(matcher)
            }
          })
        }
      }
    })

    // Blur listener
    document.body.addEventListener(
      'blur',
      (evt) => {
        const input = evt.target
        if (
          (input.tagName === 'INPUT' ||
            input.tagName === 'TEXTAREA' ||
            input.tagName === 'SELECT') &&
          input.form &&
          input.form.hasAttribute('data-uno-valid')
        ) {
          validateInput(input)
        }
      },
      true
    )

    // Character Count Listener
    document.body.addEventListener('input', (evt) => {
      const input = evt.target
      if (input.matches && input.matches('input, textarea')) {
        updateCount(input)
      }
    })
  }

  function updateCount(input) {
    // Check if input has count target
    // Convention: data-uno-count="#targetId" OR sibling .valid-count
    let target = null
    if (input.hasAttribute('data-uno-count')) {
      const selector = input.getAttribute('data-uno-count')
      if (selector) target = document.querySelector(selector)
    }

    // Fallback: look for .valid-count in parent/sibling
    if (!target && input.form) {
      // Look for a wrapper
      const wrapper = input.closest('div')
      if (wrapper) {
        target = wrapper.querySelector('.valid-count span') || wrapper.querySelector('.valid-count')
      }
    }

    if (target) {
      const len = input.value.length
      target.textContent = len
    }
  }

  function validateForm(form) {
    const inputs = form.querySelectorAll('input, select, textarea')
    let firstInvalid = null
    inputs.forEach((input) => {
      if (!validateInput(input)) {
        if (!firstInvalid) firstInvalid = input
      }
    })
    if (firstInvalid) firstInvalid.focus()
    return !firstInvalid
  }

  function validateInput(input) {
    if (input.disabled || input.type === 'hidden') return true

    // Reset validity
    input.setCustomValidity('')

    // Custom Rules Logic
    // We need to support 'valid' JSON attribute, 'data-valid-type', and custom validators
    // This logic mimics the previous implementation but updates state instead of DOM

    let rules = []

    // 1. JSON rules
    if (input.hasAttribute('valid')) {
      try {
        const jsonRules = JSON.parse(input.getAttribute('valid'))
        if (Array.isArray(jsonRules)) rules = jsonRules
      } catch (e) {
        console.warn('Uno Valid: Invalid JSON rules', e)
      }
    }

    // 2. data-valid-type
    const type = input.getAttribute('data-valid-type')
    if (type) rules.push({ type })

    // 3. Native required
    if (input.hasAttribute('required') && !rules.some((r) => r.type === 'required')) {
      rules.unshift({ type: 'required' })
    }

    // Execute rules
    // We use setCustomValidity to leverage native validation state
    for (const rule of rules) {
      // emptyquit
      if (rule.type === 'emptyquit') {
        if (!input.value) {
          input.setCustomValidity('')
          break
        }
        continue
      }

      // required
      if (rule.type === 'required') {
        if (input.validity.valueMissing) {
          input.setCustomValidity(
            rule.msg || getI18nMsg('valid.required') || 'This field is required'
          )
          break
        }
        continue
      }

      // Regex validators
      if (VALIDATORS[rule.type]) {
        if (input.value && !VALIDATORS[rule.type].test(input.value)) {
          input.setCustomValidity(rule.msg || getI18nMsg(`valid.${rule.type}`) || 'Invalid format')
          break
        }
        continue
      }

      // Global custom validators
      if (
        typeof window !== 'undefined' &&
        window.UnoValidators &&
        window.UnoValidators[rule.type]
      ) {
        const result = window.UnoValidators[rule.type](input.value, input)
        if (result !== true && result !== undefined && result !== '') {
          const errorMsg =
            typeof result === 'string'
              ? result
              : getI18nMsg(`valid.${rule.type}`) || 'Validation failed'
          input.setCustomValidity(rule.msg || errorMsg)
          break
        }
        continue
      }

      // Global function validators
      if (
        typeof window !== 'undefined' &&
        window[rule.type] &&
        typeof window[rule.type] === 'function'
      ) {
        const result = window[rule.type](input.value, input)
        if (result !== true && result !== undefined && result !== '') {
          const errorMsg = typeof result === 'string' ? result : 'Validation failed'
          input.setCustomValidity(rule.msg || errorMsg)
          break
        }
        continue
      }
    }

    // Handle data-uno-match
    if (input.hasAttribute('data-uno-match')) {
      const targetName = input.getAttribute('data-uno-match')
      const target = input.form.querySelector(`[name="${targetName}"]`)
      if (target && input.value !== target.value) {
        input.setCustomValidity(input.getAttribute('data-uno-error-match') || 'Fields do not match')
      }
    }

    // UPDATE STATE based on validity
    const isValid = input.checkValidity()

    // Toggle class on input (Box Style)
    // Use Uno.getStyle for industrial-grade configuration
    const errorClass = Uno.getStyle('valid', 'errorClass', 'is-invalid')

    if (!isValid) {
      input.classList.add(errorClass)
    } else {
      input.classList.remove(errorClass)
    }

    // Toggle class on parent (Box Container Style)
    const parentErrorClass = Uno.getStyle('valid', 'parentErrorClass')
    if (parentErrorClass) {
      // For parent selector, we still look at options or default
      const opts = (typeof window !== 'undefined' && window.Uno?.options?.valid) || {}
      const parentSelector = opts.parentSelector || '.field, .form-group, .input-group, div'
      const parent = input.closest(parentSelector)
      if (parent) {
        if (!isValid) parent.classList.add(parentErrorClass)
        else parent.classList.remove(parentErrorClass)
      }
    }

    // Update ARIA
    input.setAttribute('aria-invalid', !isValid)

    // Update Store
    if (!isValid) {
      // Prioritize custom message attributes if standard validity fails
      let msg = input.validationMessage

      // Fallback for standard HTML5 errors if no custom message set yet
      // (Though setCustomValidity above sets validationMessage, so this might be redundant but safe)
      if (input.validity.valueMissing) {
        msg = input.getAttribute('data-uno-error-required') || getI18nMsg('valid.required') || msg
      } else if (input.validity.typeMismatch) {
        let key = 'valid.pattern'
        if (input.type === 'email') key = 'valid.email'
        if (input.type === 'url') key = 'valid.url'
        msg = input.getAttribute('data-uno-error-type') || getI18nMsg(key) || msg
      } else if (input.validity.patternMismatch) {
        msg = input.getAttribute('data-uno-error-pattern') || getI18nMsg('valid.pattern') || msg
      } else if (input.validity.tooShort) {
        msg = input.getAttribute('data-uno-error-minlength') || getI18nMsg('valid.min') || msg
      } else if (input.validity.tooLong) {
        msg = input.getAttribute('data-uno-error-maxlength') || getI18nMsg('valid.max') || msg
      } else if (input.validity.rangeUnderflow) {
        msg = input.getAttribute('data-uno-error-min') || getI18nMsg('valid.rangeMin') || msg
      } else if (input.validity.rangeOverflow) {
        msg = input.getAttribute('data-uno-error-max') || getI18nMsg('valid.rangeMax') || msg
      }

      if (input.name) {
        Alpine.store('unoValid').setError(input.name, msg)
      }
    } else {
      if (input.name) {
        Alpine.store('unoValid').clearError(input.name)
      }
    }

    return isValid
  }

  function getI18nMsg(key) {
    if (typeof window !== 'undefined' && window.Alpine && window.Alpine.store('i18n')) {
      return window.Alpine.store('i18n').t(key)
    }
    return null
  }
}
