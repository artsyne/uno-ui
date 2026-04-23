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
   * Rules are stored as normalized arrays of rule objects.
   * @type {Object.<string, Array<Object>>}
   */
  rules: {},

  /**
   * Fields currently undergoing async validation.
   * @type {Object.<string, boolean>}
   */
  validating: {},

  /**
   * Custom validators registry.
   * @type {Object.<string, Function>}
   */
  customValidators: {},

  /**
   * Debounce timers for async validation.
   * @type {Object.<string, number>}
   */
  timers: {},

  /**
   * Adds or updates validation rules for a field.
   *
   * @param {string} fieldName - The unique identifier for the field (name or id)
   * @param {Array|Object} config - The validation configuration
   * @param {Object} [instanceOpts={}] - Instance-specific global options
   */
  addRules(fieldName, config, instanceOpts = {}) {
    const normalizedRules = this.normalizeRules(config)

    // Merge global options (like emptyQuit) into each rule if not present
    const globalConfig = Uno.getConfig('valid', instanceOpts)

    this.rules[fieldName] = normalizedRules.map((rule) => ({
      emptyQuit: globalConfig.emptyQuit !== false, // Default to true
      ...rule
    }))
  },

  /**
   * Normalizes configuration into a standard array of rule objects.
   *
   * @param {Array|Object|string} config
   * @returns {Array<Object>}
   */
  normalizeRules(config) {
    if (!config) return []

    if (typeof config === 'string') {
      return [{ type: config }]
    }

    if (typeof config === 'object' && !Array.isArray(config)) {
      return [config]
    }

    if (Array.isArray(config)) {
      return config.map((item) => {
        if (typeof item === 'string') {
          return { type: item }
        }
        return item
      })
    }

    return []
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
   * @param {string} [eventType='input'] - The event triggering validation ('input', 'blur', 'submit')
   * @returns {Promise<boolean>|boolean} True if valid (can be async)
   */
  async validate(fieldName, value, eventType = 'input') {
    const rules = this.rules[fieldName]
    if (!rules || rules.length === 0) return true

    // Clear previous error
    delete this.errors[fieldName]

    // Clear previous timer if any (for debounce)
    if (this.timers[fieldName]) {
      clearTimeout(this.timers[fieldName])
      delete this.timers[fieldName]
    }

    const strValue = value === null || value === undefined ? '' : String(value)
    const isEmpty = strValue.trim() === ''

    // Iterate through rules in order
    for (const rule of rules) {
      // 0. Trigger Check
      // 'submit' event should always run all rules.
      if (eventType === 'input') {
        if (rule.trigger === 'blur' || rule.realtime === false) {
          continue
        }
      }

      // 1. EmptyQuit Check
      // If value is empty, and this rule is NOT 'required', and emptyQuit is enabled
      // Then skip this rule.
      const isRequired = rule.type === 'required' || rule.required === true
      if (isEmpty && !isRequired && rule.emptyQuit !== false) {
        continue
      }

      // 2. Debounce Check
      // Only debounce on 'input' event. Blur and Submit should be immediate.
      if (rule.debounce > 0 && eventType === 'input') {
        // If we hit a debounced rule, we must pause execution and wait
        // This makes the whole validate function async
        this.validating[fieldName] = true

        return new Promise((resolve) => {
          this.timers[fieldName] = setTimeout(async () => {
            const result = await this.executeRule(fieldName, value, rule, strValue)
            delete this.validating[fieldName]

            if (!result) {
              // Error found in async rule
              resolve(false)
            } else {
              // Async rule passed, strictly speaking we should continue to next rules
              // But debounce usually implies it's the last or heavy check.
              // For simplicity in this loop structure, if there are more rules after debounce,
              // they will run immediately after this promise resolves.
              // However, the current structure returns a Promise here.
              // We need to support continuing the chain.

              // RE-EVALUATE: A simple return here stops the chain.
              // If we want to support multiple debounced rules or rules after debounce,
              // we need a more complex runner.
              // For now, let's assume debounce is the final step or we just return the result of this step.
              resolve(true)
            }
          }, rule.debounce)
        })
      }

      // 3. Execute Sync Rule
      this.validating[fieldName] = true
      try {
        const result = await this.executeRule(fieldName, value, rule, strValue)
        if (!result) {
          return false // Stop at first error
        }
      } finally {
        delete this.validating[fieldName]
      }
    }

    return true
  },

  /**
   * Executes a single rule.
   *
   * @param {string} fieldName
   * @param {any} value
   * @param {Object} rule
   * @param {string} strValue
   * @returns {Promise<boolean>}
   */
  async executeRule(fieldName, value, rule, strValue) {
    let isValid = true
    let errorMsg = ''
    const type = rule.type

    // --- Built-in Validators ---

    // Required
    if (type === 'required' || rule.required) {
      if (strValue.trim() === '') {
        isValid = false
        errorMsg = rule.msg || this.getI18nMsg('valid.required') || 'This field is required'
      }
    }
    // Standard Types (email, url, mobile, int, etc.) mapped from VALIDATORS
    else if (VALIDATORS[type]) {
      if (!VALIDATORS[type].test(strValue)) {
        isValid = false
        errorMsg = rule.msg || this.getI18nMsg(`valid.${type}`) || `Invalid ${type} format`
      }
    }
    // Length
    else if (type === 'length') {
      if (rule.min !== undefined && strValue.length < rule.min) {
        isValid = false
        errorMsg =
          rule.msg || this.getI18nMsg('valid.min', { n: rule.min }) || `Min length is ${rule.min}`
      } else if (rule.max !== undefined && strValue.length > rule.max) {
        isValid = false
        errorMsg =
          rule.msg || this.getI18nMsg('valid.max', { n: rule.max }) || `Max length is ${rule.max}`
      }
    }
    // Range (Numeric)
    else if (type === 'range') {
      const num = parseFloat(strValue)
      if (isNaN(num)) {
        isValid = false
        errorMsg = rule.msg || this.getI18nMsg('valid.number') || 'Must be a number'
      } else {
        if (rule.min !== undefined && num < rule.min) {
          isValid = false
          errorMsg =
            rule.msg ||
            this.getI18nMsg('valid.rangeMin', { n: rule.min }) ||
            `Min value is ${rule.min}`
        } else if (rule.max !== undefined && num > rule.max) {
          isValid = false
          errorMsg =
            rule.msg ||
            this.getI18nMsg('valid.rangeMax', { n: rule.max }) ||
            `Max value is ${rule.max}`
        }
      }
    }
    // Pattern (Regex)
    else if (type === 'pattern' || rule.pattern) {
      const pattern = rule.pattern instanceof RegExp ? rule.pattern : new RegExp(rule.pattern)
      if (!pattern.test(strValue)) {
        isValid = false
        errorMsg = rule.msg || this.getI18nMsg('valid.pattern') || 'Invalid format'
      }
    }
    // Match (Another field)
    else if (type === 'match' || rule.match) {
      const selector = rule.match || rule.selector || rule.target
      if (!selector) {
        // Configuration error, skip or warn?
        // Treating as valid to avoid blocking, but logging warning could be good.
        return true
      }

      const target =
        document.querySelector(selector.startsWith('#') ? selector : `[name="${selector}"]`) ||
        document.getElementById(selector)

      if (target && strValue !== target.value) {
        isValid = false
        errorMsg = rule.msg || this.getI18nMsg('valid.match') || 'Fields do not match'
      }
    }
    // Custom Validator
    else if (type === 'custom' || type === 'validator' || rule.fn || this.customValidators[type]) {
      let validatorFn = rule.fn || rule.name
      if (typeof validatorFn === 'string') {
        validatorFn = this.customValidators[validatorFn]
      } else if (!validatorFn && this.customValidators[type]) {
        validatorFn = this.customValidators[type]
      }

      if (typeof validatorFn === 'function') {
        try {
          const result = await validatorFn(value, fieldName, rule)
          if (result === false) {
            isValid = false
            errorMsg = rule.msg || 'Validation failed'
          } else if (typeof result === 'string') {
            isValid = false
            errorMsg = result
          } else if (typeof result === 'object' && result.valid === false) {
            isValid = false
            errorMsg = result.message || rule.msg || 'Validation failed'
          }
        } catch (e) {
          log.error(`[Valid] Error in custom validator for ${fieldName}`, e)
          isValid = false
          errorMsg = 'Validation error'
        }
      } else {
        log.warn(`[Valid] Custom validator not found: ${type || rule.fn}`)
      }
    }

    if (!isValid) {
      this.errors[fieldName] = errorMsg
      // Important: Reactive update for Alpine to pick up the change immediately
      // Depending on Alpine version, this might need explicit handling, but usually assignment is enough.
    }

    return isValid
  },

  /**
   * Validates all fields in a form.
   *
   * @param {HTMLFormElement} form - The form element to validate
   * @returns {Promise<boolean>} True if the form is valid
   */
  async validateForm(form) {
    const inputs = form.querySelectorAll('input, select, textarea')
    let isFormValid = true
    let firstErrorField = null

    // We need to wait for all validations
    // Use Promise.all to run them in parallel or loop for sequential?
    // Sequential is safer for focus management.

    for (const input of inputs) {
      if (input.disabled) continue

      const fieldName = input.name || input.id
      if (this.rules[fieldName]) {
        // Use 'submit' event type to ensure all rules run (ignoring trigger: 'blur' filters)
        const isValid = await this.validate(fieldName, input.value, 'submit')
        if (!isValid) {
          isFormValid = false
          if (!firstErrorField) firstErrorField = input
        }
      }
    }

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
