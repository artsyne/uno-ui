import { describe, it, expect, beforeEach, vi } from 'vitest'
import { JSDOM } from 'jsdom'

describe('Valid Component', () => {
  let ValidPlugin
  let UnoInstance
  let Alpine
  let store
  let directives = {}

  beforeEach(async () => {
    // 1. Fresh JSDOM for each test to avoid listener pollution
    const dom = new JSDOM('<!DOCTYPE html><body></body>')
    global.document = dom.window.document
    global.window = dom.window
    global.HTMLElement = dom.window.HTMLElement
    global.NodeList = dom.window.NodeList

    // Fix: Make sure CustomEvent is available globally and matches JSDOM's
    global.CustomEvent = dom.window.CustomEvent
    global.MutationObserver =
      dom.window.MutationObserver ||
      class {
        observe() {}
        disconnect() {}
        takeRecords() {
          return []
        }
      }

    // Dynamic imports to ensure MutationObserver is ready before Alpine loads
    const validModule = await import('../../src/components/valid')
    ValidPlugin = validModule.default
    const unoModule = await import('../../src/core/uno')
    UnoInstance = unoModule.UnoInstance

    // Mock Element methods
    global.HTMLElement.prototype.scrollIntoView = vi.fn()
    global.HTMLElement.prototype.focus = vi.fn()
    global.HTMLElement.prototype.matches = function (selector) {
      return (
        (this.constructor.name === 'HTMLFormElement' && selector === 'form') ||
        (this.constructor.name === 'HTMLInputElement' && selector === 'input') ||
        dom.window.document.querySelector(selector) === this
      )
    }

    // 2. Mock Alpine
    store = {}
    Alpine = {
      store: (name, value) => {
        if (value) store[name] = value
        return store[name]
      },
      directive: (name, cb) => {
        directives[name] = cb
      },
      magic: vi.fn(),
      effect: vi.fn()
    }

    // 3. Mock Uno Global
    global.Uno = new UnoInstance()

    // 4. Initialize Plugin
    ValidPlugin(Alpine)
  })

  describe('Store: valid', () => {
    it('should register store', () => {
      expect(store.uno).toBeDefined()
      expect(store.uno.valid).toBeDefined()
    })

    it('should add rules correctly', () => {
      store.uno.valid.addRules('username', { required: true })
      // Use toMatchObject because default config is merged
      expect(store.uno.valid.rules['username']).toMatchObject({ required: true })
    })

    it('should validate required field', () => {
      store.uno.valid.addRules('username', { required: true })

      const resultFail = store.uno.valid.validate('username', '')
      expect(resultFail).toBe(false)
      expect(store.uno.valid.errors['username']).toBe('This field is required')

      const resultPass = store.uno.valid.validate('username', 'valid')
      expect(resultPass).toBe(true)
      // undefined because we delete the error key
      expect(store.uno.valid.errors['username']).toBeUndefined()
    })

    it('should validate email field', () => {
      store.uno.valid.addRules('email', { email: true })

      expect(store.uno.valid.validate('email', 'invalid')).toBe(false)
      expect(store.uno.valid.validate('email', 'test@example.com')).toBe(true)
    })

    it('should validate minLength', () => {
      store.uno.valid.addRules('pass', { minLength: 5 })

      expect(store.uno.valid.validate('pass', '1234')).toBe(false)
      expect(store.uno.valid.validate('pass', '12345')).toBe(true)
    })

    it('should support emptyquit behavior (skip validation if empty and not required)', () => {
      // Rule: email type, but NOT required
      store.uno.valid.addRules('optionalEmail', { email: true, required: false })

      // Empty value -> should be valid (skipped)
      expect(store.uno.valid.validate('optionalEmail', '')).toBe(true)
      expect(store.uno.valid.errors['optionalEmail']).toBeUndefined()

      // Non-empty invalid value -> should fail
      expect(store.uno.valid.validate('optionalEmail', 'not-an-email')).toBe(false)
      expect(store.uno.valid.errors['optionalEmail']).toBeTruthy()
    })

    it('should validate numeric limits (min/max)', () => {
      store.uno.valid.addRules('age', { min: 18, max: 100, type: 'number' })

      // Valid range
      expect(store.uno.valid.validate('age', '20')).toBe(true)

      // Too low
      expect(store.uno.valid.validate('age', '17')).toBe(false)
      expect(store.uno.valid.errors['age']).toMatch(/at least 18/)

      // Too high
      expect(store.uno.valid.validate('age', '101')).toBe(false)
      expect(store.uno.valid.errors['age']).toMatch(/at most 100/)

      // Not a number
      expect(store.uno.valid.validate('age', 'abc')).toBe(false)
      expect(store.uno.valid.errors['age']).toMatch(/Invalid number format/)
    })

    it('should support custom validators', () => {
      const customFn = vi.fn((val) => {
        return val === 'secret' ? true : 'Wrong secret'
      })

      store.uno.valid.registerValidator('checkSecret', customFn)
      store.uno.valid.addRules('secretField', { validator: 'checkSecret' })

      // Invalid
      expect(store.uno.valid.validate('secretField', 'guess')).toBe(false)
      expect(store.uno.valid.errors['secretField']).toBe('Wrong secret')
      expect(customFn).toHaveBeenCalledWith('guess', 'secretField')

      // Valid
      expect(store.uno.valid.validate('secretField', 'secret')).toBe(true)
      expect(store.uno.valid.errors['secretField']).toBeUndefined()
    })
  })

  describe('Form Validation', () => {
    it('should validate all fields in form', () => {
      document.body.innerHTML = `
        <form id="testForm">
          <input name="user" value="" />
          <input name="email" value="bad-email" />
        </form>
      `
      const form = document.getElementById('testForm')

      store.uno.valid.addRules('user', { required: true })
      store.uno.valid.addRules('email', { email: true })

      const isValid = store.uno.valid.validateForm(form)

      expect(isValid).toBe(false)
      expect(store.uno.valid.errors['user']).toBeTruthy()
      expect(store.uno.valid.errors['email']).toBeTruthy()
    })
  })

  describe('HTMX Integration', () => {
    it('should intercept htmx:confirm and validate', () => {
      document.body.innerHTML = `
        <form id="htmxForm" x-uno-valid>
          <input name="field" value="" />
        </form>
      `
      const form = document.getElementById('htmxForm')
      store.uno.valid.addRules('field', { required: true })

      const preventDefault = vi.fn()
      const issueRequest = vi.fn()

      // IMPORTANT: bubbles: true is required for event delegation
      const evt = new window.CustomEvent('htmx:confirm', {
        bubbles: true,
        detail: { issueRequest, elt: form }
      })
      evt.preventDefault = preventDefault

      form.dispatchEvent(evt)

      expect(preventDefault).toHaveBeenCalled()
      expect(issueRequest).not.toHaveBeenCalled()
      expect(store.uno.valid.errors['field']).toBeTruthy()
    })

    it('should allow request if validation passes', () => {
      document.body.innerHTML = `
        <form id="htmxForm" x-uno-valid>
          <input name="field" value="valid" />
        </form>
      `
      const form = document.getElementById('htmxForm')
      store.uno.valid.addRules('field', { required: true })

      const preventDefault = vi.fn()
      const issueRequest = vi.fn()

      const evt = new window.CustomEvent('htmx:confirm', {
        bubbles: true,
        detail: { issueRequest, elt: form }
      })
      evt.preventDefault = preventDefault

      form.dispatchEvent(evt)

      expect(preventDefault).not.toHaveBeenCalled()
      expect(issueRequest).not.toHaveBeenCalled()
    })
  })
})
