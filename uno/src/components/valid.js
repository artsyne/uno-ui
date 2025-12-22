import { VALIDATORS, buildSafeRegExp, safeParseJSON } from '../utils/validation'
import { Uno } from '../core/uno'
import { log } from '../utils/log'
import en from '../locales/en'

const regexCache = new WeakMap()
const DEBOUNCE_CANCELLED = Symbol('uno.valid.debounce.cancelled')

const splitClasses = (cls) =>
  String(cls || '')
    .split(/\s+/)
    .filter(Boolean)

const applyClasses = (el, cls, enabled) => {
  if (!el) return
  for (const c of splitClasses(cls)) {
    if (enabled) el.classList.add(c)
    else el.classList.remove(c)
  }
}

const formatTemplate = (str, vars = {}) =>
  String(str ?? '').replace(/\{(\w+)\}/g, (_, k) => (k in vars ? String(vars[k]) : `{${k}}`))

const debounceAsyncByKey = (field, key, delayMs, fn) => {
  if (!field) return fn()
  if (!field._debounceMap) field._debounceMap = new Map()

  const map = field._debounceMap
  const prev = map.get(key)
  if (prev) {
    clearTimeout(prev.timer)
    prev.resolve(DEBOUNCE_CANCELLED)
    map.delete(key)
  }

  const ms = Math.max(0, Number(delayMs) || 0)
  if (ms === 0) return fn()

  return new Promise((resolve, reject) => {
    const timer = setTimeout(async () => {
      map.delete(key)
      try {
        const out = await fn()
        resolve(out)
      } catch (e) {
        reject(e)
      }
    }, ms)
    map.set(key, { timer, resolve })
  })
}

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
    const fullKey = `valid.${key}`
    if (s && typeof s.t === 'function') {
      const out = s.t(fullKey, vars)
      if (out !== fullKey) return out
    }
    const fallback = en?.valid?.[key]
    if (typeof fallback === 'string') return formatTemplate(fallback, vars)
    return fullKey
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
      return false
    },
    clearForm(form) {
      if (!form) return
      if (form._x_dataStack) {
        const validComponent = form._x_dataStack.find((d) => d.reset)
        if (validComponent) {
          validComponent.reset()
        }
      }
    }
  })

  // Expose Public API
  const api = {
    // Validate all forms or specific form?
    // Mostly helper to set state manually
    set: (el, { valid, message }) => {
      if (el) {
        el.dispatchEvent(
          new CustomEvent('valid-set', {
            detail: { valid, message, name: el.name },
            bubbles: true
          })
        )
      }
    },
    validate: async (target) => {
      const form = target?.tagName === 'FORM' ? target : target?.closest?.('form')
      return Alpine.store('valid').validateForm(form)
    },
    clear: (target) => {
      const form = target?.tagName === 'FORM' ? target : target?.closest?.('form')
      Alpine.store('valid').clearForm(form)
    }
  }

  Uno.valid = api
  if (typeof window !== 'undefined' && window.Uno) {
    window.Uno.valid = api
  }

  // 4. Main Data Component
  Alpine.data('validForm', () => ({
    fields: {}, // Reactive state: { [name]: { value, error, valid, ... } }

    init() {
      this.$el.noValidate = true
      // Handle HTMX trigger manually if needed, but we prefer hx-trigger="valid-submit"

      // Listen for manual set events from Uno.valid.set
      this.$el.addEventListener('valid-set', (e) => {
        const { name, valid, message } = e.detail
        if (name && this.fields[name]) {
          this.fields[name].valid = valid
          this.fields[name].error = message || null
          // If setting invalid manually, mark as touched so error shows
          if (!valid) {
            this.fields[name].touched = true
          }
        } else if (e.target && e.target.name && this.fields[e.target.name]) {
          // Target based dispatch
          const fieldName = e.target.name
          this.fields[fieldName].valid = valid
          this.fields[fieldName].error = message || null
          if (!valid) {
            this.fields[fieldName].touched = true
          }
        }
      })
    },

    // Reset validation state (API: Uno.valid.clear)
    reset() {
      const stateInvalid = Uno.getStyle('valid', 'stateInvalid', 'is-invalid')
      const stateValid = Uno.getStyle('valid', 'stateValid', 'is-valid')
      const stateValidating = Uno.getStyle('valid', 'stateValidating', 'validating')
      const inputError = Uno.getStyle('valid', 'inputError', '!border-red-500')
      const inputValid = Uno.getStyle('valid', 'inputValid', '!border-green-500')
      const inputPaddingIcon = Uno.getStyle('valid', 'inputPaddingIcon', 'pr-10')

      Object.keys(this.fields).forEach((name) => {
        const f = this.fields[name]
        f._validationSeq += 1
        f.error = null
        f.valid = true
        f.touched = false
        f.dirty = false
        f.validating = false

        if (f._debounceMap) {
          for (const v of f._debounceMap.values()) {
            clearTimeout(v.timer)
            v.resolve(DEBOUNCE_CANCELLED)
          }
          f._debounceMap.clear()
        }

        applyClasses(f.el, stateInvalid, false)
        applyClasses(f.el, stateValid, false)
        applyClasses(f.el, stateValidating, false)
        applyClasses(f.el, inputError, false)
        applyClasses(f.el, inputValid, false)
        applyClasses(f.el, inputPaddingIcon, false)

        if (f.icon) {
          f.icon.innerHTML = ''
          f.icon.style.opacity = '0'
        }
        if (f.spinner) {
          f.spinner.style.display = ''
        }

        const wrapper = f.wrapper || f.el?.closest?.('.uno-input-wrapper')
        if (wrapper) {
          wrapper.querySelectorAll('[data-valid-error-for]').forEach((n) => {
            if (n.getAttribute('data-valid-error-for') === name) n.remove()
          })
        }
      })
    },

    // Register field (called by x-validate directive)
    register(name, el, rules = [], meta = null) {
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

      for (const r of fieldRules) {
        const ruleObj = typeof r === 'object' && r ? r : null
        if (
          ruleObj?.type === 'compareto' &&
          typeof ruleObj.destid === 'string' &&
          !ruleObj._destEl
        ) {
          const dest = document.getElementById(ruleObj.destid)
          ruleObj._destEl = dest || null
          ruleObj._destResolved = true
        }
      }

      this.fields[name] = {
        el,
        wrapper: meta?.wrapper || null,
        icon: meta?.icon || null,
        spinner: meta?.spinner || null,
        rules: fieldRules,
        value: el.value,
        error: null,
        valid: true,
        dirty: false,
        touched: false,
        validating: false,
        _validationSeq: 0
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

      const seq = (field._validationSeq += 1)
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
          // P2-6: WeakMap Cache for RegExp
          let re = regexCache.get(ruleObj)
          if (!re) {
            const cacheKey = String(ruleObj.expression)
            re = buildSafeRegExp(cacheKey, { defaultFlags: 'i' })
            if (re) {
              regexCache.set(ruleObj, re)
            }
          }

          if (!re) error = ruleObj.msg || t('pattern')
          else if (!re.test(val)) error = ruleObj.msg || t('pattern')
        } else if (type === 'limit') {
          if (ruleObj.min && val.length < ruleObj.min)
            error = ruleObj.msg || t('min', { n: ruleObj.min })
          if (ruleObj.max && val.length > ruleObj.max)
            error = ruleObj.msg || t('max', { n: ruleObj.max })
        } else if (type === 'range') {
          const n = parseFloat(val)
          if (
            ruleObj.min !== undefined &&
            ruleObj.max !== undefined &&
            (n < ruleObj.min || n > ruleObj.max)
          )
            error = ruleObj.msg || t('range', { min: ruleObj.min, max: ruleObj.max })
          else if (ruleObj.min !== undefined && n < ruleObj.min)
            error = ruleObj.msg || t('rangeMin', { n: ruleObj.min })
          else if (ruleObj.max !== undefined && n > ruleObj.max)
            error = ruleObj.msg || t('rangeMax', { n: ruleObj.max })
        } else if (type === 'compareto') {
          if (!ruleObj._destResolved && typeof ruleObj.destid === 'string') {
            ruleObj._destEl = document.getElementById(ruleObj.destid) || null
            ruleObj._destResolved = true
          }
          const dest = ruleObj._destEl
          if (dest && val !== dest.value) error = ruleObj.msg || t('compareto')
        } else if (type === 'tag') {
          const mark = ruleObj.mark || ','
          const parts = val.split(mark).filter((s) => s.trim())
          if (ruleObj.max && parts.length > ruleObj.max)
            error = ruleObj.msg || t('tagMax', { n: ruleObj.max })
        }

        if (error) break

        // 3. Async Callbacks
        if (type === 'callback' && ruleObj.fn) {
          const cb = Alpine.store('valid').callbacks[ruleObj.fn]
          if (cb) {
            try {
              const delayMs = Number.isFinite(ruleObj.debounce)
                ? ruleObj.debounce
                : Number.isFinite(Uno.options?.plugins?.valid?.callbackDebounce)
                  ? Uno.options.plugins.valid.callbackDebounce
                  : Number.isFinite(Uno.options?.valid?.callbackDebounce)
                    ? Uno.options.valid.callbackDebounce
                    : 300
              const res = await debounceAsyncByKey(field, `callback:${ruleObj.fn}`, delayMs, () =>
                cb(field.el.id, val)
              )
              if (res === DEBOUNCE_CANCELLED) return field.valid
              if (field._validationSeq !== seq) return field.valid
              if (res === false || (typeof res === 'object' && res.valid === false)) {
                error = (typeof res === 'object' ? res.message : '') || ruleObj.msg || t('callback')
              } else if (typeof res === 'string') {
                error = res
              }
            } catch (err) {
              log.warn('[Uno.valid] Async validation failed', err)
              error = ruleObj.msg || t('callback')
            }
          }
        }

        if (error) break
      }

      if (field._validationSeq !== seq) return field.valid
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

    const wrapperCls = Uno.getStyle('valid', 'wrapper', 'relative w-full')
    const iconBaseCls = Uno.getStyle(
      'valid',
      'iconBase',
      'absolute right-0 top-0 bottom-0 flex items-center justify-center w-10 pointer-events-none transition-opacity duration-200'
    )
    const inputError = Uno.getStyle('valid', 'inputError', '!border-red-500')
    const inputValid = Uno.getStyle('valid', 'inputValid', '!border-green-500')
    const iconError = Uno.getStyle('valid', 'iconError', 'text-red-500 dark:text-red-400')
    const iconValid = Uno.getStyle('valid', 'iconValid', 'text-green-500 dark:text-green-400')
    const stateInvalid = Uno.getStyle('valid', 'stateInvalid', 'is-invalid')
    const stateValid = Uno.getStyle('valid', 'stateValid', 'is-valid')
    const stateValidating = Uno.getStyle('valid', 'stateValidating', 'validating')
    const inputPaddingIcon = Uno.getStyle('valid', 'inputPaddingIcon', 'pr-10')

    // Ensure DOM Structure (Wrapper + Icon)
    const setupDom = () => {
      // Check if already wrapped
      let wrapper = el.parentNode
      if (!wrapper.classList.contains('uno-input-wrapper')) {
        wrapper = document.createElement('div')
        wrapper.className = 'uno-input-wrapper'
        applyClasses(wrapper, wrapperCls, true)

        // Insert wrapper before input
        el.parentNode.insertBefore(wrapper, el)
        // Move input inside wrapper
        wrapper.appendChild(el)
      } else {
        applyClasses(wrapper, wrapperCls, true)
      }

      // Check if icon exists
      let icon = wrapper.querySelector('.uno-valid-icon')
      if (!icon) {
        icon = document.createElement('div')
        icon.className = `uno-valid-icon ${iconBaseCls}`
        icon.style.opacity = '0'
        wrapper.appendChild(icon)
      } else {
        icon.className = `uno-valid-icon ${iconBaseCls}`
      }

      return { wrapper, icon }
    }

    // Initialize DOM structure immediately to prevent layout shifts/focus loss later
    const { wrapper, icon } = setupDom()

    let spinner = null
    let sibling = wrapper.nextElementSibling
    let checks = 0
    while (sibling && checks < 3) {
      if (sibling.nodeType === 1) {
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

    // Get Data Scope
    // We assume x-validate is used within x-data="validForm"
    // We can't easily inject, so we rely on the user setting up the form correctly.

    // Parse Rules from attribute
    let rules = []
    const validAttr = el.getAttribute('valid')
    if (validAttr) {
      const parsed = safeParseJSON(validAttr)
      rules = Array.isArray(parsed) ? parsed : []
    }

    const limitRule = rules.find(
      (r) => (r.type === 'limit' || r === 'limit') && typeof r === 'object'
    )
    const limitMax = limitRule?.max
    const limitShowClip = limitRule?.show === true && Number.isFinite(limitMax)

    // Register Loop
    // We use effect to ensure we register when data is ready (sometimes init takes a tick)
    let dataRef = null
    effect(() => {
      const data = Alpine.$data(el)
      if (data && typeof data.register === 'function') {
        dataRef = data
        data.register(name, el, rules, { wrapper, icon, spinner })
      }
    })

    // Events
    const onInput = () => {
      const data = dataRef || Alpine.$data(el)
      if (limitShowClip && limitMax && el.value.length > limitMax) {
        el.value = el.value.substring(0, limitMax)
      }
      if (data && data.update) data.update(name, el.value)
    }

    const onBlur = () => {
      const data = dataRef || Alpine.$data(el)
      if (data && data.touch) data.touch(name)
    }

    el.addEventListener('input', onInput)
    el.addEventListener('blur', onBlur)

    // UI Effect (Classes & Styles)
    // Keeping this for "Plugin" ease-of-use, so user doesn't need verbose :class bindings
    effect(() => {
      const data = dataRef || Alpine.$data(el)
      if (data && data.fields && data.fields[name]) {
        const field = data.fields[name]
        // Usually we show error after touch or dirty+invalid
        // Let's show error if field.error is set (validate sets it)

        const showError = !!field.error
        const showValid = field.valid && field.dirty && !field.error

        // Classes on Input
        applyClasses(el, stateInvalid, showError)
        applyClasses(el, stateValid, showValid)
        applyClasses(el, stateValidating, field.validating)

        // Reset icon state
        icon.innerHTML = ''
        icon.style.opacity = '0'

        // Tailwind borders
        if (showError) {
          applyClasses(el, inputError, true)
          applyClasses(el, inputValid, false)
          applyClasses(el, inputPaddingIcon, true)

          // Set Icon
          icon.innerHTML = SVG_ICONS.invalid
          icon.className = `uno-valid-icon ${iconBaseCls} ${iconError}`
          icon.style.opacity = '1'
        } else if (showValid) {
          applyClasses(el, inputValid, true)
          applyClasses(el, inputError, false)
          applyClasses(el, inputPaddingIcon, true)

          // Set Icon
          icon.innerHTML = SVG_ICONS.valid
          icon.className = `uno-valid-icon ${iconBaseCls} ${iconValid}`
          icon.style.opacity = '1'
        } else {
          applyClasses(el, inputError, false)
          applyClasses(el, inputValid, false)
          applyClasses(el, inputPaddingIcon, false)
        }

        // Spinner
        // Since we wrapped the input, we need to look outside the wrapper for the spinner.
        // It might be the next sibling, or the one after (e.g. if error container is in between).
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
          const errorText = Uno.getStyle(
            'valid',
            'helperText',
            'text-red-600 dark:text-red-400 text-xs'
          )
          el.className = errorText
        } else {
          el.textContent = ''
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
