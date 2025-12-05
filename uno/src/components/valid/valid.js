/**
 * Uno - Validation Component
 *
 * @module uno.components.valid
 * @version 0.1.0
 * @license MIT
 */

import { log } from '../../core/log'

export default (Alpine) => {
  // Validation Icons (SVG Data URIs)
  const INVALID_ICON = `url("data:image/svg+xml,%3csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 12 12' width='12' height='12' fill='none' stroke='%23ef4444'%3e%3ccircle cx='6' cy='6' r='4.5'/%3e%3cpath stroke-linecap='round' d='M5.8 3.6h.4L6 6.5zM6 8.2af.6.6 0 110-1.2.6.6 0 010 1.2z'/%3e%3c/svg%3e")`
  const VALID_ICON = `url("data:image/svg+xml,%3csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 8 8'%3e%3cpath fill='%2322c55e' d='M2.3 6.73L.6 4.53c-.4-1.04.46-1.4 1.1-.8l1.1 1.4 3.4-3.8c.6-.63 1.6-.27 1.2.7l-4 4.6c-.43.5-.8.4-1.1.1z'/%3e%3c/svg%3e")`

  /**
   * Update validation visual state using Tailwind utility classes and inline styles
   * @param {HTMLElement} el - Input element
   * @param {string} state - 'valid', 'invalid', 'validating', 'clear'
   */
  const updateValidationState = (el, state) => {
    // Reset visual classes
    el.classList.remove('!border-red-500', '!border-green-500', 'pr-10', 'bg-no-repeat', 'bg-[right_0.75rem_center]', 'bg-[length:1rem_1rem]')
    el.style.backgroundImage = ''

    // Reset state markers
    el.classList.remove('is-valid', 'is-invalid')

    // Spinner logic
    const nextDiv = el.nextElementSibling
    if (nextDiv && nextDiv.tagName === 'DIV') {
      const spinner = nextDiv.querySelector('.animate-spin')
      if (spinner) spinner.style.display = (state === 'validating') ? 'block' : ''
    }

    if (state === 'validating') {
      el.classList.add('validating')
      return
    }
    el.classList.remove('validating')

    if (state === 'invalid') {
      el.classList.add('is-invalid')
      el.classList.add('!border-red-500', 'pr-10', 'bg-no-repeat', 'bg-[right_0.75rem_center]', 'bg-[length:1rem_1rem]')
      el.style.backgroundImage = INVALID_ICON
    } else if (state === 'valid') {
      el.classList.add('is-valid')
      el.classList.add('!border-green-500', 'pr-10', 'bg-no-repeat', 'bg-[right_0.75rem_center]', 'bg-[length:1rem_1rem]')
      el.style.backgroundImage = VALID_ICON
    }
  }

  Alpine.store('valid', {
    callbacks: {},
    events: {},
    
    /**
     * Register a custom validation callback
     * @param {string} name - Callback name
     * @param {Function} callback - Async validation function
     * @param {string} [eventType='blur'] - Trigger event (blur/input)
     */
    registerCallback(name, callback, eventType = 'blur') {
      if (typeof callback === 'function') {
        this.callbacks[name] = callback
        this.events[name] = eventType
      }
    }
  })

  Alpine.data('validForm', ({ rules = {} } = {}) => ({
    errors: {},
    rules,

    /**
     * Translate message key
     * @param {string} key - Message key
     * @param {Object} [vars] - Variables to replace
     * @returns {string} Translated message
     */
    t(key, vars) {
      try {
        const s = Alpine.store('i18n')
        let msg = (s && s.t(`valid.${key}`, vars)) || ''
        if (!msg && key === 'min' && vars?.n != null) msg = `At least ${vars.n} characters`
        if (!msg && key === 'required') msg = 'Required'
        if (!msg && key === 'pattern') msg = 'Invalid format'
        return msg
      } catch { return '' }
    },

    VALIDATORS: {
      validname: /^[a-zA-Z0-9_-]+$/,
      en: /^[a-zA-Z]+$/,
      domain: /^((?!-)[A-Za-z0-9-]{1,63}(?<!-)\.)+[A-Za-z]{2,6}$/,
      mobile: /^\+?[\d\s-]{7,}$/,
      int: /^-?\d+$/,
      email: /^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9-]+(?:\.[a-zA-Z0-9-]+)*$/,
      date: /^(\d{4}-\d{2}-\d{2}|\d{2}\/\d{2}\/\d{4}|\d{2}\.\d{2}\.\d{4})$/,
      url: /^(https?:\/\/)?([\da-z\.-]+)\.([a-z\.-]{2,6})([\/\w \.-]*)*\/?$/i,
      pwd: /^(?=.*[A-Za-z])(?=.*\d)[A-Za-z\d]{8,}$/,
      strongpwd: /^(?=.*[A-Za-z])(?=.*\d)(?=.*[@$!%*#?&])[A-Za-z\d@$!%*#?&]{8,}$/
    },

    resolvePattern(el, name) {
      const typeKey = el.dataset.validType
      if (typeKey && this.VALIDATORS[typeKey]) return this.VALIDATORS[typeKey]
      // Prefer built-in by input type
      if (el.type === 'email') return this.VALIDATORS.email
      if (el.type === 'url') return this.VALIDATORS.url
      // rules.pattern
      const rule = this.rules[name]
      if (rule && rule.pattern) {
        if (typeof rule.pattern === 'string') {
          if (this.VALIDATORS[rule.pattern]) return this.VALIDATORS[rule.pattern]
          const { pattern, flags } = sanitizeRegexString(rule.pattern)
          return new RegExp(pattern, flags || 'i')
        }
        return rule.pattern
      }
      // data-valid-pattern
      if (el.dataset.validPattern) {
        const { pattern, flags } = sanitizeRegexString(el.dataset.validPattern)
        return new RegExp(pattern, flags || 'i')
      }
      return null
    },

    init() {
      try {
      const form = this.$el.tagName === 'FORM' ? this.$el : this.$el.closest('form')
      const root = form || this.$el
      root.querySelectorAll('[name]').forEach((el) => {
        el.addEventListener('input', () => this.validateField(el))
        el.addEventListener('blur', () => this.validateField(el))
      })

      root.querySelectorAll('[valid]').forEach((el) => {
        try {
        if (!el.dataset.loadedvalid) {
          const next = el.nextElementSibling
          let hasFeedback = false
          if (next && (next.classList.contains('valid-feedback') || next.classList.contains('valid-hint'))) {
            hasFeedback = true
          }
          
          if (!hasFeedback) {
            const div = document.createElement('div')
            div.className = 'valid-feedback text-xs text-red-600'
            el.after(div)
          }
          const rulesArr = safeParseJSON(el.getAttribute('valid'))
          const hasCallback = rulesArr.some(r => (r.type || r) === 'callback')
          if (hasCallback) {
            el.addEventListener('input', (e) => this.checkElement(el, { event: e.type }))
            el.addEventListener('blur', (e) => this.checkElement(el, { event: e.type }))
            rulesArr.forEach(r => {
              const type = r.type || r
              if (type === 'callback') {
                const cb = r.fn
                const evt = Alpine.store('valid').events[cb] || 'blur'
                if (evt !== 'input' && evt !== 'blur') {
                  el.addEventListener(evt, () => this.checkElement(el, { event: evt, onlyCallback: true, callbackName: cb }))
                }
              }
            })
          } else {
            el.addEventListener('input', (e) => this.checkElement(el, { event: e.type }))
            el.addEventListener('blur', (e) => this.checkElement(el, { event: e.type }))
          }
          rulesArr.forEach(json => {
            if (json.show === true) {
              const span = ensureCountSpan(el)
              span.className = 'valid-count block w-full text-xs text-gray-500 dark:text-gray-400 text-right mt-2'
              el.addEventListener('input', () => this.updateCharCount(el))
              this.updateCharCount(el)
            }
          })
          el.dataset.loadedvalid = 'true'
        }
        } catch (err) {
          log.error('Error in valid loop for element:', el, err)
        }
      })

      if (form) {
        // ... existing submit listeners ...
        form.addEventListener('submit', async (e) => {
          const ok = this.validateForm()
          if (!ok) {
            e.preventDefault()
            e.stopPropagation()
            e.stopImmediatePropagation()
            log.warn('valid: prevent submit due to errors', this.errors)
            // form.querySelectorAll('[name]').forEach((el) => this.validateField(el))
            return
          }
          // Dev fallback: emulate hx-post/hx-get for static responses
          const targetSel = form.getAttribute('hx-target')
          const url = form.getAttribute('hx-post') || form.getAttribute('hx-get')
          if (url && targetSel) {
            e.preventDefault()
            const swapEl = document.querySelector(targetSel)
            try {
              const res = await fetch(url)
              const ct = res.headers.get('content-type') || ''
              let content = ''
              if (ct.includes('application/json')) {
                const json = await res.json()
                content = `<pre class="bg-gray-100 p-3 rounded overflow-x-auto whitespace-pre-wrap">${JSON.stringify(json, null, 2)}</pre>`
              } else {
                content = await res.text()
              }
              swapEl.innerHTML = content
            } catch (err) {
              log.error('valid: submit fetch error', err)
            }
          }
        }, { capture: true })

        const submitBtn = form.querySelector('[type="submit"]')
        if (submitBtn) {
          submitBtn.addEventListener('click', async (e) => {
            const ok = this.validateForm()
            if (!ok) {
              e.preventDefault()
              // form.querySelectorAll('[name]').forEach((el) => this.validateField(el))
              return
            }
            const targetSel = form.getAttribute('hx-target')
            const url = form.getAttribute('hx-post') || form.getAttribute('hx-get')
            if (url && targetSel) {
              e.preventDefault()
              const swapEl = document.querySelector(targetSel)
              try {
                const res = await fetch(url)
                const ct = res.headers.get('content-type') || ''
                let content = ''
                if (ct.includes('application/json')) {
                  const json = await res.json()
                  content = `<pre class="bg-gray-100 p-3 rounded overflow-x-auto whitespace-pre-wrap">${JSON.stringify(json, null, 2)}</pre>`
                } else {
                  content = await res.text()
                }
                swapEl.innerHTML = content
              } catch (err) {
                // ignore
              }
            }
          }, { capture: true })
        }
      }
      } catch(e) {
        log.error('Valid Init Error:', e)
      }
    },

    validateField(el) {
      const name = el.getAttribute('name')
      const raw = el.value || ''
      const val = raw.trim()
      const required = el.dataset.validRequired === 'true'
      const pattern = this.resolvePattern(el, name)
      const min = el.dataset.validMin ? parseInt(el.dataset.validMin, 10) : null
      const validRulesArr = safeParseJSON(el.getAttribute('valid'))
      const hasEmptyQuit = validRulesArr.some(r => (r.type || r) === 'emptyquit')
      let err = ''
      if (required && !val) err = this.t('required')
      else if (hasEmptyQuit && val.length === 0) err = ''
      else if (pattern && val.length > 0 && !pattern.test(val)) {
        let key = el.dataset.validType || (el.type === 'email' ? 'email' : (el.type === 'url' ? 'url' : ''))
        err = (key && this.t(key)) || this.t('pattern')
      }
      else if (min !== null && val.length > 0 && val.length < min) err = this.t('min', { n: min })
      if (!err && this.rules[name]) {
        const r = this.rules[name]
        if (r.required && !val) err = this.t('required')
        if (!err && r.pattern && val.length > 0) {
          const isBuiltin = typeof r.pattern === 'string' && this.VALIDATORS[r.pattern]
          const rp = isBuiltin ? this.VALIDATORS[r.pattern] : (typeof r.pattern === 'string' ? new RegExp(r.pattern) : r.pattern)
          if (!(rp).test(val)) err = (isBuiltin && this.t(r.pattern)) || this.t('pattern')
        }
        if (!err && r.min && val.length > 0 && val.length < r.min) err = this.t('min', { n: r.min })
      }
      if (err) this.errors[name] = err
      else delete this.errors[name]
      
      // Toggle validation classes
      if (err) {
        updateValidationState(el, 'invalid')
      } else if (!el.classList.contains('validating')) {
        // Only mark valid if we are not waiting for async validation
        updateValidationState(el, 'valid')
      }

      const hint = el.nextElementSibling
      if (hint && (hint.classList.contains('valid-hint') || hint.classList.contains('valid-feedback'))) {
        hint.textContent = err || ''
        const isFeedback = hint.classList.contains('valid-feedback')
        const base = 'valid-hint text-xs transition-colors duration-200'
        const color = err ? 'valid-hint-error text-red-600' : 'valid-hint-normal text-gray-400'
        hint.className = `${base} ${color} ${isFeedback ? 'valid-feedback' : ''}`
      }
      return !err
    },

    updateCharCount(el) {
      const rulesArr = safeParseJSON(el.getAttribute('valid'))
      if (!rulesArr) return
      const span = ensureCountSpan(el)
      if (!span) return
      for (const json of rulesArr) {
        if (json.show === true) {
          const maxlen = json.max
          const value = el.value || ''
          if (value.length > maxlen) el.value = value.substring(0, maxlen)
          const remaining = maxlen - (el.value || '').length
          const s = Alpine?.store('i18n')
          const msg = (s && s.t('charCount', { count: remaining })) || `Remaining: ${remaining}`
          span.innerHTML = msg
          break
        }
      }
    },

    checkError() {
      const form = this.$el.tagName === 'FORM' ? this.$el : this.$el.closest('form')
      const elements = Array.from(form.querySelectorAll('[valid]'))
      const results = elements.map(el => ({ el, hasError: this.checkElement(el, { event: 'submit' }) }))
      const bad = results.find(r => r.hasError)
      if (bad) bad.el.focus()
      return results.some(r => r.hasError)
    },

    checkElement(el, options = {}) {
      const rulesArr = safeParseJSON(el.getAttribute('valid')) || []
      const next = el.nextElementSibling
      const labelError = next && (next.classList.contains('valid-feedback') || next.classList.contains('valid-hint')) ? next : null
      let hasError = false
      let errorMsg = ''
      const eventType = options.event || 'blur'
      for (const json of rulesArr) {
        const type = json.type || json
        const val = el.value || ''
        if (options.onlyCallback && type === 'callback' && json.fn !== options.callbackName) continue
        if (options.onlyCallback && type !== 'callback') continue
        if (type === 'callback') {
          const cbName = json.fn
          const cbEvent = Alpine.store('valid').events[cbName] || 'blur'
          if (eventType !== cbEvent && eventType !== 'submit') continue
        }
        switch (type) {
          case 'emptyquit':
            if (val.length === 0) return false
            break
          case 'required':
            if (val.length === 0) errorMsg = json.msg || this.t('required')
            break
          case 'limit':
            if (json.min !== undefined && val.length < json.min) errorMsg = json.msg || `Min ${json.min}`
            if (!errorMsg && json.max !== undefined && val.length > json.max) errorMsg = json.msg || `Max ${json.max}`
            break
          case 'validname':
          case 'en':
          case 'domain':
          case 'mobile':
          case 'int':
          case 'email':
          case 'date':
          case 'url':
          case 'pwd':
          case 'strongpwd':
            if (val.length > 0) {
              const regex = this.VALIDATORS[type]
              if (!regex.test(val)) errorMsg = json.msg || this.t(type)
            }
            break
          case 'regexp':
            if (val.length > 0) {
              const { pattern, flags } = sanitizeRegexString(json.expression)
              const re = new RegExp(pattern, flags || 'i')
              if (!re.test(val)) errorMsg = json.msg || this.t('pattern')
            }
            break
          case 'range':
            if (val.length > 0) {
              const n = parseFloat(val)
              if (json.min !== undefined && n < json.min) errorMsg = json.msg || `Min ${json.min}`
              if (!errorMsg && json.max !== undefined && n > json.max) errorMsg = json.msg || `Max ${json.max}`
            }
            break
          case 'tag':
            const mark = json.mark || ','
            const parts = val.split(mark).filter(s => s.trim())
            if (json.max && parts.length > json.max) errorMsg = json.msg || `Max ${json.max} items`
            break
          case 'compareto':
            const dest = document.getElementById(json.destid)
            if (dest && val !== dest.value) errorMsg = json.msg || 'Not match'
            break
          case 'callback':
            // Async callback
            const cbName = json.fn
            if (Alpine.store('valid').callbacks[cbName]) {
              updateValidationState(el, 'validating')
              Alpine.store('valid').callbacks[cbName](el.id, val).then(res => {
                updateValidationState(el, 'clear') // Clear validating state
                let cbError = ''
                if (res === false) cbError = 'Invalid'
                else if (typeof res === 'string') cbError = res
                else if (typeof res === 'object' && res.valid === false) cbError = res.message
                
                if (cbError) {
                  updateValidationState(el, 'invalid')
                  if (labelError) labelError.textContent = cbError
                } else {
                   // If no previous sync error, mark valid
                   if (!el.classList.contains('is-invalid')) updateValidationState(el, 'valid')
                }
              })
            }
            break
        }
        if (errorMsg) {
          hasError = true
          break
        }
      }

      if (hasError) {
        updateValidationState(el, 'invalid')
        if (labelError) labelError.textContent = errorMsg
      } else if (!el.classList.contains('validating')) {
         // Check simple rules too
         const simpleValid = this.validateField(el)
         if (simpleValid) {
            updateValidationState(el, 'valid')
            if (labelError) labelError.textContent = ''
         }
      }
      return hasError
    },
    
    /**
     * Validate the entire form
     * @returns {boolean} True if valid
     */
    validateForm() {
      const form = this.$el.tagName === 'FORM' ? this.$el : this.$el.closest('form')
      let ok = true
      // 1. Check basic fields
      form.querySelectorAll('[name]').forEach((el) => {
        if (!this.validateField(el)) ok = false
      })
      // 2. Check [valid] complex rules
      if (this.checkError()) ok = false
      
      return ok
    }

  }))

  function safeParseJSON(str) {
    try {
      if (!str) return []
      return JSON.parse(str)
    } catch (e) {
      return []
    }
  }

  function sanitizeRegexString(str) {
    const match = str.match(/^\/(.*?)\/([gimsuy]*)$/)
    if (match) {
      return { pattern: match[1], flags: match[2] }
    }
    return { pattern: str, flags: '' }
  }
  
  function ensureCountSpan(el) {
    let span = el.parentNode.querySelector('.valid-count')
    if (!span) {
      span = document.createElement('span')
      span.className = 'valid-count block w-full text-xs text-gray-500 dark:text-gray-400 text-right mt-2'
      el.parentNode.appendChild(span)
    }
    return span
  }
}
