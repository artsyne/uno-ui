/**
 * Valid Component (Form Validation).
 *
 * A Store-First, Directive-Driven form validation plugin for Alpine.js and HTMX.
 * It manages validation state, rules, and error messages in a central store,
 * while using declarative directives for configuration.
 *
 * @module uno.components.valid
 * @requires alpinejs
 * @requires htmx.org
 * @since 2.0.0
 */

import { VALIDATORS } from '../utils/validation'
import { Uno } from '../core/uno'
import { registerPlugin } from '../core/plugin-registry'
import { log } from '../utils/log'

/**
 * Initializes the Valid component.
 *
 * @param {Object} Alpine - The Alpine.js instance
 * @param {Object} [options={}] - Global configuration options
 */
export default (Alpine, options = {}) => {
  // Register plugin defaults
  const pluginDefaults = {
    realtime: true,
    showError: true,
    scrollOnError: true
  }
  Uno.registerPluginDefaults('valid', {
    ...pluginDefaults,
    ...options.plugins?.valid,
    ...options.valid
  })

  /**
   * Valid Store.
   *
   * Manages validation state and logic.
   * Registered as `Alpine.store('uno').valid`.
   *
   * @class ValidStore
   */
  const validStore = {
    /**
     * Validation error messages keyed by field name.
     * @type {Object.<string, string>}
     */
    errors: {},

    /**
     * Validation rules keyed by field name.
     * @type {Object.<string, Object>}
     */
    rules: {},

    /**
     * Custom validators registry.
     * @type {Object.<string, Function>}
     */
    customValidators: {},

    /**
     * Adds or updates validation rules for a field.
     *
     * Merges configuration from 4 levels:
     * 1. Plugin Defaults
     * 2. Global User Config
     * 3. Global Plugin Config
     * 4. Instance Config (passed argument)
     *
     * @param {string} fieldName - The unique identifier for the field (name or id)
     * @param {Object} [instanceOpts={}] - Instance-specific configuration
     */
    addRules(fieldName, instanceOpts = {}) {
      // 4-level merge via Uno.getConfig
      this.rules[fieldName] = Uno.getConfig('valid', instanceOpts)
    },

    /**
     * Registers a custom validator function.
     *
     * @param {string} name - The name of the validator
     * @param {Function} validatorFn - The validator function. Receives (value, fieldId). Returns true, false, string (error msg), or Promise.
     */
    registerValidator(name, validatorFn) {
      if (typeof validatorFn === 'function') {
        this.customValidators[name] = validatorFn
      } else {
        log.error(`[Valid] Invalid validator function for ${name}`)
      }
    },

    /**
     * Validates a single field against its rules.
     *
     * @param {string} fieldName - The field identifier
     * @param {string|number} value - The value to validate
     * @returns {boolean} True if valid, false otherwise
     */
    validate(fieldName, value) {
      const config = this.rules[fieldName]
      if (!config) return true

      let isValid = true
      let errorMsg = ''
      const strValue = value === null || value === undefined ? '' : String(value)
      const isEmpty = strValue.trim() === ''

      // 1. Required Check
      if (config.required) {
        if (isEmpty) {
          isValid = false
          errorMsg =
            config.msgRequired || this.getI18nMsg('valid.required') || 'This field is required'
        }
      }

      // If empty and not required, skip other checks (EmptyQuit behavior)
      if (isValid && isEmpty && !config.required) {
        // Clear error if it was previously set
        delete this.errors[fieldName]
        return true
      }

      // 2. Type/Regex Validators (email, url, etc.)
      if (isValid) {
        // Check explicit type or inferred type
        const type = config.type || (config.email ? 'email' : null)

        if (type && VALIDATORS[type]) {
          if (!VALIDATORS[type].test(strValue)) {
            isValid = false
            errorMsg =
              config[`msg${type.charAt(0).toUpperCase() + type.slice(1)}`] ||
              config.msgType ||
              this.getI18nMsg(`valid.${type}`) ||
              `Invalid ${type} format`
          }
        }
      }

      // 3. Length Limits (MinLength/MaxLength)
      if (isValid && config.minLength !== null && config.minLength !== undefined) {
        if (strValue.length < config.minLength) {
          isValid = false
          errorMsg =
            config.msgMinLength ||
            this.getI18nMsg('valid.min', { n: config.minLength }) ||
            `Minimum length is ${config.minLength}`
        }
      }

      if (isValid && config.maxLength !== null && config.maxLength !== undefined) {
        if (strValue.length > config.maxLength) {
          isValid = false
          errorMsg =
            config.msgMaxLength ||
            this.getI18nMsg('valid.max', { n: config.maxLength }) ||
            `Maximum length is ${config.maxLength}`
        }
      }

      // 4. Numeric Range (Min/Max)
      if (
        isValid &&
        ((config.min !== null && config.min !== undefined) ||
          (config.max !== null && config.max !== undefined))
      ) {
        // Only validate if it looks like a number
        if (VALIDATORS.number.test(strValue)) {
          const numVal = parseFloat(strValue)

          if (config.min !== null && config.min !== undefined && numVal < config.min) {
            isValid = false
            errorMsg =
              config.msgMin ||
              this.getI18nMsg('valid.rangeMin', { n: config.min }) ||
              `Value must be at least ${config.min}`
          }

          if (isValid && config.max !== null && config.max !== undefined && numVal > config.max) {
            isValid = false
            errorMsg =
              config.msgMax ||
              this.getI18nMsg('valid.rangeMax', { n: config.max }) ||
              `Value must be at most ${config.max}`
          }
        } else {
          // If min/max are set but value isn't a number
          isValid = false
          errorMsg = config.msgNumber || this.getI18nMsg('valid.number') || 'Must be a number'
        }
      }

      // 5. Pattern (Custom Regex)
      if (isValid && config.pattern) {
        try {
          // Handle string pattern from HTML attribute
          const pattern =
            config.pattern instanceof RegExp ? config.pattern : new RegExp(config.pattern)

          if (!pattern.test(strValue)) {
            isValid = false
            errorMsg = config.msgPattern || this.getI18nMsg('valid.pattern') || 'Invalid format'
          }
        } catch (e) {
          log.error(`[Valid] Invalid regex pattern for ${fieldName}`, e)
        }
      }

      // 6. Match (CompareTo)
      if (isValid && config.match) {
        const target =
          document.querySelector(
            config.match.startsWith('#') ? config.match : `[name="${config.match}"]`
          ) || document.getElementById(config.match)

        if (target && strValue !== target.value) {
          isValid = false
          errorMsg = config.msgMatch || this.getI18nMsg('valid.match') || 'Fields do not match'
        }
      }

      // 7. Custom Validator (Function or Registered Name)
      if (isValid && config.validator) {
        let validatorFn = config.validator

        // Lookup if string
        if (typeof config.validator === 'string') {
          validatorFn = this.customValidators[config.validator]
          if (!validatorFn) {
            log.warn(`[Valid] Validator "${config.validator}" not found`)
          }
        }

        if (typeof validatorFn === 'function') {
          // Note: Async validators not fully supported in sync flow yet,
          // but we support returning Promise which will be truthy (valid) initially.
          // For full async, we'd need async/await here and in callers.
          // Currently assuming sync return or simple object.
          const result = validatorFn(value, fieldName)

          if (result === false) {
            isValid = false
            errorMsg = config.msgValidator || 'Validation failed'
          } else if (typeof result === 'string') {
            isValid = false
            errorMsg = result
          } else if (typeof result === 'object' && result.valid === false) {
            isValid = false
            errorMsg = result.message || config.msgValidator || 'Validation failed'
          }
        }
      }

      // Update State
      if (!isValid) {
        this.errors[fieldName] = errorMsg
      } else {
        delete this.errors[fieldName]
      }

      return isValid
    },

    /**
     * Validates all fields in a form.
     *
     * @param {HTMLFormElement} form - The form element to validate
     * @returns {boolean} True if the form is valid
     */
    validateForm(form) {
      const inputs = form.querySelectorAll('input, select, textarea')
      let isFormValid = true
      let firstErrorField = null

      inputs.forEach((input) => {
        // Skip hidden or disabled inputs if desired, but HTML5 validation usually checks them if they have rules
        if (input.disabled) return

        const fieldName = input.name || input.id
        if (this.rules[fieldName]) {
          const isValid = this.validate(fieldName, input.value)
          if (!isValid) {
            isFormValid = false
            if (!firstErrorField) firstErrorField = input
          }
        }
      })

      // Handle focus and scroll
      if (!isFormValid && firstErrorField) {
        const globalConfig = Uno.getConfig('valid') || {}
        if (globalConfig.scrollOnError !== false) {
          firstErrorField.scrollIntoView({ behavior: 'smooth', block: 'center' })
          firstErrorField.focus({ preventScroll: true })
        }
      }

      return isFormValid
    },

    /**
     * Manually sets an error for a field.
     *
     * @param {string} fieldName - The field identifier
     * @param {string} message - The error message
     */
    showError(fieldName, message) {
      this.errors[fieldName] = message
    },

    /**
     * Clears error for a specific field.
     *
     * @param {string} fieldName
     */
    hideError(fieldName) {
      delete this.errors[fieldName]
    },

    /**
     * Clears all validation errors.
     */
    clearErrors() {
      this.errors = {}
    },

    /**
     * Bulk sets errors (e.g., from server response).
     *
     * @param {Object.<string, string>} errorsObj - Object mapping field names to error messages
     */
    showErrors(errorsObj) {
      Object.assign(this.errors, errorsObj)
    },

    /**
     * Helper to get I18n message safely.
     *
     * @private
     * @param {string} key - Translation key
     * @param {Object} [params] - Translation parameters
     * @returns {string|null} The translated string or null
     */
    getI18nMsg(key, params) {
      const i18n = Alpine.store('i18n')
      if (i18n && typeof i18n.t === 'function') {
        return i18n.t(key, params)
      }
      return null
    }
  }

  // Register the plugin
  registerPlugin(Alpine, 'valid', validStore)

  /**
   * x-uno-valid Directive.
   *
   * Activates validation on an element.
   * Parses HTML5 attributes to build validation rules.
   *
   * Usage:
   * <input x-uno-valid name="email" required type="email">
   */
  Alpine.directive('uno-valid', (el, { expression }, { evaluateLater, cleanup }) => {
    // If used on a form, it acts as a marker for the global HTMX listener
    if (el.tagName === 'FORM') return

    const fieldName = el.name || el.id
    if (!fieldName) {
      console.warn('[Uno Valid] Element has no name or id:', el)
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
