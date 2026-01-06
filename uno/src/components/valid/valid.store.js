import { VALIDATORS } from '../../utils/validation'
import { Uno } from '../../core/uno'
import { log } from '../../utils/log'

export const createValidStore = (Alpine) => ({
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
   * @param {string} fieldName - The unique identifier for the field (name or id)
   * @param {Object} [instanceOpts={}] - Instance-specific configuration
   */
  addRules(fieldName, instanceOpts = {}) {
    this.rules[fieldName] = Uno.getConfig('valid', instanceOpts)
  },

  /**
   * Registers a custom validator function.
   *
   * @param {string} name - The name of the validator
   * @param {Function} validatorFn - The validator function.
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
      delete this.errors[fieldName]
      return true
    }

    // 2. Type/Regex Validators (email, url, etc.)
    if (isValid) {
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

    // 3. Length Limits
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

    // 4. Numeric Range
    if (
      isValid &&
      ((config.min !== null && config.min !== undefined) ||
        (config.max !== null && config.max !== undefined))
    ) {
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
        isValid = false
        errorMsg = config.msgNumber || this.getI18nMsg('valid.number') || 'Must be a number'
      }
    }

    // 5. Pattern
    if (isValid && config.pattern) {
      try {
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

    // 6. Match
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

    // 7. Custom Validator
    if (isValid && config.validator) {
      let validatorFn = config.validator

      if (typeof config.validator === 'string') {
        validatorFn = this.customValidators[config.validator]
        if (!validatorFn) {
          log.warn(`[Valid] Validator "${config.validator}" not found`)
        }
      }

      if (typeof validatorFn === 'function') {
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

    if (!isFormValid && firstErrorField) {
      const globalConfig = Uno.getConfig('valid') || {}
      if (globalConfig.scrollOnError !== false) {
        firstErrorField.scrollIntoView({ behavior: 'smooth', block: 'center' })
        firstErrorField.focus({ preventScroll: true })
      }
    }

    return isFormValid
  },

  showError(fieldName, message) {
    this.errors[fieldName] = message
  },

  hideError(fieldName) {
    delete this.errors[fieldName]
  },

  clearErrors() {
    this.errors = {}
  },

  showErrors(errorsObj) {
    Object.assign(this.errors, errorsObj)
  },

  getI18nMsg(key, params) {
    const i18n = Alpine.store('i18n')
    if (i18n && typeof i18n.t === 'function') {
      return i18n.t(key, params)
    }
    return null
  }
})
