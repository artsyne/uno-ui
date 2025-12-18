import { VALIDATORS, sanitizeRegexString, safeParseJSON } from '../../utils/validation'

/**
 * Uno - Validation Component (Data-Driven Refactor)
 *
 * @module uno.components.valid
 * @version 0.3.0
 * @license MIT
 */

export default (Alpine) => {
  // 1. Validation Icons (SVG Strings)
  const SVG_ICONS = {
    invalid: `<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 12 12' fill='none' stroke='currentColor' class='w-4 h-4'><circle cx='6' cy='6' r='4.5'/><path stroke-linecap='round' d='M5.8 3.6h.4L6 6.5zM6 8.2a.6.6 0 110-1.2.6.6 0 010 1.2z'/></svg>`,
    valid: `<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 8 8' class='w-4 h-4'><path fill='currentColor' d='M2.3 6.73L.6 4.53c-.4-1.04.46-1.4 1.1-.8l1.1 1.4 3.4-3.8c.6-.63 1.6-.27 1.2.7l-4 4.6c-.43.5-.8.4-1.1.1z'/></svg>`
  }

  // 2. Helper: I18n
  const t = (key, vars) => {
    const s = Alpine.store('i18n')
    if (!s) return key
    // Fallback messages if store translation missing
    const defaults = {
      required: 'Required',
      pattern: 'Invalid format',
      min: `At least ${vars?.n} characters`,
      max: `Max ${vars?.n} characters`,
      int: 'Integer required',
      email: 'Invalid email address',
      url: 'Invalid URL'
    }
    return s.t(`valid.${key}`, vars) || defaults[key] || key
  }

  // 3. Global Store for Async Callbacks
  Alpine.store('valid', {
    callbacks: {},
    registerCallback(name, fn) {
      if (typeof fn === 'function') {
        this.callbacks[name] = fn
      }
    },
    async validateForm(form) {
      if (!form) return false
      if (form._x_dataStack) {
        const validComponent = form._x_dataStack.find((d) => d.validateAll)
        if (validComponent) {
          return await validComponent.validateAll()
        }
      }
      return true
    }
  })

  // 4. Main Data Component
  Alpine.data('validForm', () => ({
    fields: {}, // Reactive state: { [name]: { value, error, valid, ... } }

    init() {
      this.$el.noValidate = true
      // Handle HTMX trigger manually if needed, but we prefer hx-trigger="valid-submit"
    },

    // Register field (called by x-validate directive)
    register(name, el, rules = []) {
      if (this.fields[name]) return

      // Initial rules processing
      const fieldRules = [...rules]

      // Extract HTML5 attributes into rules
      if (el.hasAttribute('required') || el.dataset.validRequired === 'true') {
        fieldRules.push({ type: 'required' })
      }
      if (el.type === 'email' || el.dataset.validType === 'email') fieldRules.push('email')
      if (el.type === 'url' || el.dataset.validType === 'url') fieldRules.push('url')
      // Add other types based on dataset
      if (el.dataset.validType) {
        const type = el.dataset.validType
        if (type !== 'email' && type !== 'url') fieldRules.push(type)
      }
      if (el.dataset.validPattern) {
        fieldRules.push({ type: 'regexp', expression: el.dataset.validPattern })
      }

      this.fields[name] = {
        el,
        rules: fieldRules,
        value: el.value,
        error: null,
        valid: true,
        dirty: false,
        touched: false,
        validating: false
      }
    },

    update(name, value) {
      if (!this.fields[name]) return
      this.fields[name].value = value
      this.fields[name].dirty = true
      // Validate on input if it was already invalid or touched?
      // For better UX, we often validate on input only if already dirty/touched.
      // Here we validate immediately to update state, but UI might choose when to show error.
      this.validate(name)
    },

    touch(name) {
      if (!this.fields[name]) return
      this.fields[name].touched = true
      this.validate(name)
    },

    async validate(name) {
      const field = this.fields[name]
      if (!field) return true

      field.validating = true
      field.error = null

      const val = (field.value || '').trim()
      let error = null

      // Validation Pipeline
      for (const rule of field.rules) {
        const type = rule.type || rule
        const ruleObj = typeof rule === 'object' ? rule : {}

        // 1. Empty Check (Skip subsequent rules if empty and not required)
        if (!val) {
          if (type === 'required') {
            error = ruleObj.msg || t('required')
            break
          }
          if (type === 'emptyquit') break
          // If not required and empty, skip other validators except 'callback' which might handle empty?
          // Usually empty value is valid if not required.
          continue
        }

        // 2. Sync Validators
        if (VALIDATORS[type]) {
          if (!VALIDATORS[type].test(val)) {
            error = ruleObj.msg || t(type) || t('pattern')
          }
        } else if (type === 'regexp' && ruleObj.expression) {
          const { pattern, flags } = sanitizeRegexString(ruleObj.expression)
          const re = new RegExp(pattern, flags || 'i')
          if (!re.test(val)) error = ruleObj.msg || t('pattern')
        } else if (type === 'limit') {
          if (ruleObj.min && val.length < ruleObj.min)
            error = ruleObj.msg || t('min', { n: ruleObj.min })
          if (ruleObj.max && val.length > ruleObj.max)
            error = ruleObj.msg || t('max', { n: ruleObj.max })
        } else if (type === 'range') {
          const n = parseFloat(val)
          if (ruleObj.min !== undefined && n < ruleObj.min)
            error = ruleObj.msg || `Min ${ruleObj.min}`
          if (ruleObj.max !== undefined && n > ruleObj.max)
            error = ruleObj.msg || `Max ${ruleObj.max}`
        } else if (type === 'compareto') {
          const dest = document.getElementById(ruleObj.destid)
          if (dest && val !== dest.value) error = ruleObj.msg || 'Not match'
        } else if (type === 'tag') {
          const mark = ruleObj.mark || ','
          const parts = val.split(mark).filter((s) => s.trim())
          if (ruleObj.max && parts.length > ruleObj.max)
            error = ruleObj.msg || `Max ${ruleObj.max} items`
        }

        if (error) break

        // 3. Async Callbacks
        if (type === 'callback' && ruleObj.fn) {
          const cb = Alpine.store('valid').callbacks[ruleObj.fn]
          if (cb) {
            try {
              const res = await cb(field.el.id, val)
              if (res === false || (typeof res === 'object' && res.valid === false)) {
                error = (typeof res === 'object' ? res.message : '') || 'Invalid'
              } else if (typeof res === 'string') {
                error = res
              }
            } catch {
              // console.warn('Async validation failed')
            }
          }
        }

        if (error) break
      }

      field.error = error
      field.valid = !error
      field.validating = false
      return field.valid
    },

    async validateAll() {
      const names = Object.keys(this.fields)
      const results = await Promise.all(names.map((name) => this.validate(name)))
      return results.every((r) => r)
    },

    async submit(e) {
      // Use e.preventDefault() in the template @click.prevent or @submit.prevent
      // But if passed here, we ensure it.
      if (e) {
        e.preventDefault()
        e.stopPropagation()
      }

      // Touch all fields to show errors
      Object.keys(this.fields).forEach((name) => {
        this.fields[name].touched = true
      })

      const isValid = await this.validateAll()

      if (isValid) {
        // Emit event for HTMX or other listeners
        this.$el.dispatchEvent(new CustomEvent('valid-submit', { bubbles: true }))

        // If no hx-post/get/trigger, native submit
        if (!this.$el.hasAttribute('hx-post') && !this.$el.hasAttribute('hx-get')) {
          this.$el.submit()
        }
      } else {
        // Focus first invalid
        const first = Object.keys(this.fields).find((name) => !this.fields[name].valid)
        if (first && this.fields[first].el) this.fields[first].el.focus()
      }
    }
  }))

  // 5. Directive: x-validate
  // Usage: <input x-validate name="username">
  Alpine.directive('validate', (el, _, { effect, cleanup }) => {
    const name = el.getAttribute('name')
    if (!name) return

    // Ensure DOM Structure (Wrapper + Icon)
    const setupDom = () => {
      // Check if already wrapped
      let wrapper = el.parentNode
      if (!wrapper.classList.contains('uno-input-wrapper')) {
        wrapper = document.createElement('div')
        wrapper.className = 'relative w-full uno-input-wrapper' // w-full to match typical input width

        // Insert wrapper before input
        el.parentNode.insertBefore(wrapper, el)
        // Move input inside wrapper
        wrapper.appendChild(el)
      }

      // Check if icon exists
      let icon = wrapper.querySelector('.uno-valid-icon')
      if (!icon) {
        icon = document.createElement('div')
        icon.className =
          'uno-valid-icon absolute right-0 top-0 bottom-0 flex items-center justify-center w-10 pointer-events-none transition-opacity duration-200'
        icon.style.opacity = '0'
        wrapper.appendChild(icon)
      }

      return { wrapper, icon }
    }

    // Initialize DOM structure immediately to prevent layout shifts/focus loss later
    const { icon } = setupDom()

    // Get Data Scope
    // We assume x-validate is used within x-data="validForm"
    // We can't easily inject, so we rely on the user setting up the form correctly.

    // Parse Rules from attribute
    let rules = []
    const validAttr = el.getAttribute('valid')
    if (validAttr) rules = safeParseJSON(validAttr)

    // Register Loop
    // We use effect to ensure we register when data is ready (sometimes init takes a tick)
    effect(() => {
      const data = Alpine.$data(el)
      if (data && typeof data.register === 'function') {
        data.register(name, el, rules)
      }
    })

    // Events
    const onInput = () => {
      const data = Alpine.$data(el)
      if (data && data.update) data.update(name, el.value)

      // Handle Limit/Count Logic (Clipping)
      if (rules.find((r) => (r.type === 'limit' || r === 'limit') && r.show === true)) {
        const limitRule = rules.find((r) => r.type === 'limit' || r === 'limit')
        if (limitRule && limitRule.max && el.value.length > limitRule.max) {
          el.value = el.value.substring(0, limitRule.max)
          if (data && data.update) data.update(name, el.value)
        }
      }
    }

    const onBlur = () => {
      const data = Alpine.$data(el)
      if (data && data.touch) data.touch(name)
    }

    el.addEventListener('input', onInput)
    el.addEventListener('blur', onBlur)

    // UI Effect (Classes & Styles)
    // Keeping this for "Plugin" ease-of-use, so user doesn't need verbose :class bindings
    effect(() => {
      const data = Alpine.$data(el)
      if (data && data.fields && data.fields[name]) {
        const field = data.fields[name]
        // Usually we show error after touch or dirty+invalid
        // Let's show error if field.error is set (validate sets it)

        const showError = !!field.error
        const showValid = field.valid && field.dirty && !field.error

        // Classes on Input
        el.classList.toggle('is-invalid', showError)
        el.classList.toggle('is-valid', showValid)
        el.classList.toggle('validating', field.validating)

        // Reset icon state
        icon.innerHTML = ''
        icon.style.opacity = '0'

        // Tailwind borders
        if (showError) {
          el.classList.add('!border-red-500', 'pr-10')
          el.classList.remove('!border-green-500')

          // Set Icon
          icon.innerHTML = SVG_ICONS.invalid
          icon.className =
            'uno-valid-icon absolute right-0 top-0 bottom-0 flex items-center justify-center w-10 pointer-events-none transition-opacity duration-200 text-red-500 dark:text-red-400'
          icon.style.opacity = '1'
        } else if (showValid) {
          el.classList.add('!border-green-500', 'pr-10')
          el.classList.remove('!border-red-500')

          // Set Icon
          icon.innerHTML = SVG_ICONS.valid
          icon.className =
            'uno-valid-icon absolute right-0 top-0 bottom-0 flex items-center justify-center w-10 pointer-events-none transition-opacity duration-200 text-green-500 dark:text-green-400'
          icon.style.opacity = '1'
        } else {
          el.classList.remove('!border-red-500', '!border-green-500', 'pr-10')
        }

        // Spinner
        // Since we wrapped the input, we need to look outside the wrapper for the spinner.
        // It might be the next sibling, or the one after (e.g. if error container is in between).
        const wrapper = el.parentNode
        let sibling = wrapper.nextElementSibling
        let spinner = null
        let checks = 0

        while (sibling && checks < 3) {
          if (sibling.nodeType === 1) {
            // Element node
            if (sibling.classList.contains('animate-spin')) {
              spinner = sibling
            } else {
              spinner = sibling.querySelector('.animate-spin')
            }
            if (spinner) break
          }
          sibling = sibling.nextElementSibling
          checks++
        }

        if (spinner) spinner.style.display = field.validating ? 'block' : ''
      }
    })

    cleanup(() => {
      el.removeEventListener('input', onInput)
      el.removeEventListener('blur', onBlur)
    })
  })

  // 6. Directive: x-valid-error
  // Usage: <div x-valid-error="username"></div>
  Alpine.directive('valid-error', (el, { expression }, { effect }) => {
    const name = expression
    el.style.display = 'none' // Default hidden

    effect(() => {
      const data = Alpine.$data(el)
      if (data && data.fields && data.fields[name]) {
        const field = data.fields[name]
        const error = field.error
        if (error) {
          el.textContent = error
          el.style.display = 'block'
          el.classList.add('text-red-600', 'dark:text-red-400', 'text-xs') // Default styling
          el.classList.remove('text-gray-400')
        } else {
          el.style.display = 'none'
        }
      }
    })
  })

  // 7. Directive: x-valid-count
  // Usage: <div x-valid-count="bio"></div>
  Alpine.directive('valid-count', (el, { expression }, { effect }) => {
    const name = expression
    effect(() => {
      const data = Alpine.$data(el)
      if (data && data.fields && data.fields[name]) {
        const field = data.fields[name]
        // Find max rule
        const limitRule = field.rules.find((r) => r.type === 'limit' || r === 'limit')
        if (limitRule && limitRule.max) {
          const len = (field.value || '').length
          const remaining = limitRule.max - len
          const s = Alpine.store('i18n')
          el.textContent =
            (s && s.t('charCount', { count: remaining })) || `Remaining: ${remaining}`
          el.style.display = 'block'
        } else {
          el.style.display = 'none'
        }
      }
    })
  })
}
