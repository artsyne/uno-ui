import { describe, it, expect, beforeEach, vi } from 'vitest'
import { TextEncoder, TextDecoder } from 'util'

global.TextEncoder = TextEncoder
global.TextDecoder = TextDecoder

describe('Valid Component', () => {
  let ValidPlugin
  let UnoInstance
  let Alpine
  let store
  let directives = {}

  beforeEach(async () => {
    // 1. Fresh JSDOM for each test
    const { JSDOM } = await import('jsdom')
    const dom = new JSDOM('<!DOCTYPE html><body></body>')
    global.document = dom.window.document
    global.window = dom.window
    global.HTMLElement = dom.window.HTMLElement
    global.NodeList = dom.window.NodeList

    global.CustomEvent = dom.window.CustomEvent
    global.Event = dom.window.Event
    global.MutationObserver =
      dom.window.MutationObserver ||
      class {
        observe() {}
        disconnect() {}
        takeRecords() {
          return []
        }
      }

    // Dynamic imports
    const validModule = await import('../../src/components/valid')
    ValidPlugin = validModule.default
    const unoModule = await import('../../src/core/uno')
    UnoInstance = unoModule.UnoInstance

    // Mock Element methods
    dom.window.HTMLElement.prototype.scrollIntoView = vi.fn()
    dom.window.HTMLElement.prototype.focus = vi.fn()
    dom.window.HTMLElement.prototype.matches = function (selector) {
      return (
        (this.constructor.name === 'HTMLFormElement' && selector === 'form') ||
        (this.constructor.name === 'HTMLInputElement' && selector === 'input') ||
        dom.window.document.querySelector(selector) === this
      )
    }
    // Also mock on Element just in case
    dom.window.Element.prototype.scrollIntoView = vi.fn()
    dom.window.Element.prototype.focus = vi.fn()

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
      effect: vi.fn(),
      initTree: vi.fn()
    }

    // 3. Mock Uno Global
    global.Uno = new UnoInstance()

    // 4. Initialize Plugin
    ValidPlugin(Alpine)
  })

  describe('Store: valid (Normalization)', () => {
    it('should register store', () => {
      expect(store.uno).toBeDefined()
      expect(store.uno.valid).toBeDefined()
    })

    it('should normalize string config to array', () => {
      store.uno.valid.addRules('field1', 'required')
      const rules = store.uno.valid.rules['field1']
      expect(Array.isArray(rules)).toBe(true)
      expect(rules[0]).toMatchObject({ type: 'required', emptyQuit: true })
    })

    it('should normalize array of strings', () => {
      store.uno.valid.addRules('field2', ['required', 'email'])
      const rules = store.uno.valid.rules['field2']
      expect(rules).toHaveLength(2)
      expect(rules[0]).toMatchObject({ type: 'required' })
      expect(rules[1]).toMatchObject({ type: 'email' })
    })

    it('should normalize mixed array', () => {
      store.uno.valid.addRules('field3', ['required', { type: 'length', min: 5 }])
      const rules = store.uno.valid.rules['field3']
      expect(rules[0]).toMatchObject({ type: 'required' })
      expect(rules[1]).toMatchObject({ type: 'length', min: 5 })
    })

    it('should handle single object config (legacy-ish but supported)', () => {
      store.uno.valid.addRules('field4', { type: 'custom' })
      const rules = store.uno.valid.rules['field4']
      expect(Array.isArray(rules)).toBe(true)
      expect(rules[0]).toMatchObject({ type: 'custom' })
    })
  })

  describe('Validation Logic', () => {
    it('should validate required field', async () => {
      store.uno.valid.addRules('username', ['required'])

      const resultFail = await store.uno.valid.validate('username', '')
      expect(resultFail).toBe(false)
      expect(store.uno.valid.errors['username']).toBeTruthy()

      const resultPass = await store.uno.valid.validate('username', 'valid')
      expect(resultPass).toBe(true)
      expect(store.uno.valid.errors['username']).toBeUndefined()
    })

    it('should validate email field (using built-in VALIDATORS)', async () => {
      store.uno.valid.addRules('email', ['email'])

      expect(await store.uno.valid.validate('email', 'invalid')).toBe(false)
      expect(await store.uno.valid.validate('email', 'test@example.com')).toBe(true)
    })

    it('should respect EmptyQuit (skip validation if empty and not required)', async () => {
      // Rule: email type, but NOT required (implicit)
      store.uno.valid.addRules('optionalEmail', ['email'])

      // Empty value -> should be valid (skipped due to EmptyQuit default)
      expect(await store.uno.valid.validate('optionalEmail', '')).toBe(true)
      expect(store.uno.valid.errors['optionalEmail']).toBeUndefined()

      // Non-empty invalid value -> should fail
      expect(await store.uno.valid.validate('optionalEmail', 'not-an-email')).toBe(false)
    })

    it('should validate numeric limits (min/max) via range type', async () => {
      // Note: validation.js doesn't have 'range' type in VALIDATORS regex,
      // but valid.store.js handles 'range' type logic manually?
      // Let's check valid.store.js implementation.
      // It has `else if (type === 'range')` block.

      store.uno.valid.addRules('age', [{ type: 'range', min: 18, max: 100 }])

      // Valid range
      expect(await store.uno.valid.validate('age', '20')).toBe(true)

      // Too low
      expect(await store.uno.valid.validate('age', '17')).toBe(false)
      expect(store.uno.valid.errors['age']).toMatch(/Min value/)

      // Too high
      expect(await store.uno.valid.validate('age', '101')).toBe(false)
      expect(store.uno.valid.errors['age']).toMatch(/Max value/)
    })

    it('should validate match rule (password confirm)', async () => {
      document.body.innerHTML = `
        <form>
          <input name="password" value="secret123" />
          <input name="confirm" value="secret123" />
          <input name="confirm_bad" value="wrong" />
        </form>
      `

      store.uno.valid.addRules('confirm', [{ type: 'match', target: 'password' }])
      store.uno.valid.addRules('confirm_bad', [{ type: 'match', target: 'password' }])

      // Good match
      const isValid = await store.uno.valid.validate('confirm', 'secret123')
      expect(isValid).toBe(true)
      expect(store.uno.valid.errors['confirm']).toBeUndefined()

      // Bad match
      const isValidBad = await store.uno.valid.validate('confirm_bad', 'wrong')
      expect(isValidBad).toBe(false)
      expect(store.uno.valid.errors['confirm_bad']).toBeTruthy()
    })

    it('should support custom validators', async () => {
      const customFn = vi.fn((val) => {
        return val === 'secret' ? true : 'Wrong secret'
      })

      store.uno.valid.registerValidator('checkSecret', customFn)
      store.uno.valid.addRules('secretField', [{ type: 'custom', fn: 'checkSecret' }])

      // Invalid
      expect(await store.uno.valid.validate('secretField', 'guess')).toBe(false)
      expect(store.uno.valid.errors['secretField']).toBe('Wrong secret')
      expect(customFn).toHaveBeenCalled()

      // Valid
      expect(await store.uno.valid.validate('secretField', 'secret')).toBe(true)
    })

    it('should support custom validators with type="validator" and name property', async () => {
      const customFn = vi.fn((val) => val === 'secret')
      store.uno.valid.registerValidator('checkSecret2', customFn)
      store.uno.valid.addRules('secretField2', [{ type: 'validator', name: 'checkSecret2' }])

      // Invalid
      await store.uno.valid.validate('secretField2', 'guess')
      expect(customFn).toHaveBeenCalled()
      expect(store.uno.valid.errors['secretField2']).toBeTruthy()

      // Valid
      await store.uno.valid.validate('secretField2', 'secret')
      expect(store.uno.valid.errors['secretField2']).toBeUndefined()
    })
  })

  describe('Validation Sequence', () => {
    it('should stop at first error', async () => {
      // Order: required -> length -> email
      // We use a custom validator to spy on execution
      const spyValidator = vi.fn(() => true)
      store.uno.valid.registerValidator('spy', spyValidator)

      store.uno.valid.addRules('seq', [
        'required',
        { type: 'length', min: 5 },
        { type: 'custom', fn: 'spy' }
      ])

      // 1. Fails required (empty)
      // spy should NOT be called
      await store.uno.valid.validate('seq', '')
      expect(store.uno.valid.errors['seq']).toMatch(/required/i)
      expect(spyValidator).not.toHaveBeenCalled()

      // 2. Fails length ('abc')
      // spy should NOT be called
      await store.uno.valid.validate('seq', 'abc')
      expect(store.uno.valid.errors['seq']).toMatch(/Min length/i)
      expect(spyValidator).not.toHaveBeenCalled()

      // 3. Passes both
      // spy SHOULD be called
      await store.uno.valid.validate('seq', 'abcde')
      expect(store.uno.valid.errors['seq']).toBeUndefined()
      expect(spyValidator).toHaveBeenCalled()
    })
  })

  describe('Trigger Logic', () => {
    it('should respect trigger: "blur"', async () => {
      const spyValidator = vi.fn(() => true)
      store.uno.valid.registerValidator('spy', spyValidator)
      store.uno.valid.addRules('blurField', [{ type: 'custom', fn: 'spy', trigger: 'blur' }])

      // 1. Triggered by INPUT -> Should SKIP
      await store.uno.valid.validate('blurField', 'val', 'input')
      expect(spyValidator).not.toHaveBeenCalled()

      // 2. Triggered by BLUR -> Should RUN
      await store.uno.valid.validate('blurField', 'val', 'blur')
      expect(spyValidator).toHaveBeenCalledTimes(1)
    })

    it('should respect realtime: false (legacy)', async () => {
      const spyValidator = vi.fn(() => true)
      store.uno.valid.registerValidator('spyLegacy', spyValidator)
      store.uno.valid.addRules('legacyField', [
        { type: 'custom', fn: 'spyLegacy', realtime: false }
      ])

      // 1. Triggered by INPUT -> Should SKIP
      await store.uno.valid.validate('legacyField', 'val', 'input')
      expect(spyValidator).not.toHaveBeenCalled()

      // 2. Triggered by BLUR -> Should RUN
      await store.uno.valid.validate('legacyField', 'val', 'blur')
      expect(spyValidator).toHaveBeenCalledTimes(1)
    })

    it('should always run on submit trigger', async () => {
      const spyValidator = vi.fn(() => true)
      store.uno.valid.registerValidator('spySubmit', spyValidator)
      store.uno.valid.addRules('submitField', [
        { type: 'custom', fn: 'spySubmit', trigger: 'blur' }
      ])

      // 1. Triggered by SUBMIT -> Should RUN regardless of trigger: 'blur'
      await store.uno.valid.validate('submitField', 'val', 'submit')
      expect(spyValidator).toHaveBeenCalledTimes(1)
    })
  })

  describe('Async & Debounce', () => {
    it('should debounce async validation', async () => {
      vi.useFakeTimers()

      const asyncFn = vi.fn(async () => true)
      store.uno.valid.registerValidator('asyncCheck', asyncFn)

      store.uno.valid.addRules('asyncField', [{ type: 'custom', fn: 'asyncCheck', debounce: 500 }])

      // Start validation
      const promise = store.uno.valid.validate('asyncField', 'test')

      // Should be in validating state
      expect(store.uno.valid.validating['asyncField']).toBe(true)
      expect(asyncFn).not.toHaveBeenCalled()

      // Advance time
      vi.advanceTimersByTime(500)

      // Wait for promise
      const result = await promise

      expect(result).toBe(true)
      expect(asyncFn).toHaveBeenCalled()
      expect(store.uno.valid.validating['asyncField']).toBeUndefined()

      vi.useRealTimers()
    })

    it('should NOT debounce on blur event', async () => {
      vi.useFakeTimers()
      const asyncFn = vi.fn(async () => true)
      store.uno.valid.registerValidator('blurAsync', asyncFn)
      store.uno.valid.addRules('blurAsyncField', [
        { type: 'custom', fn: 'blurAsync', debounce: 500 }
      ])

      // Trigger by BLUR
      const promise = store.uno.valid.validate('blurAsyncField', 'test', 'blur')

      // Should call immediately (no waiting for timer)
      // But since it's async function, we await it.
      await promise

      expect(asyncFn).toHaveBeenCalled()

      vi.useRealTimers()
    })

    it('should debounce custom validator with type="validator" and name property', async () => {
      vi.useFakeTimers()
      const asyncFn = vi.fn(async () => true)
      store.uno.valid.registerValidator('checkDebounce', asyncFn)
      store.uno.valid.addRules('debounceField', [
        { type: 'validator', name: 'checkDebounce', debounce: 300 }
      ])

      // Start validation
      const promise = store.uno.valid.validate('debounceField', 'test')

      // Should be in validating state
      expect(store.uno.valid.validating['debounceField']).toBe(true)
      expect(asyncFn).not.toHaveBeenCalled()

      // Advance time partially
      vi.advanceTimersByTime(200)
      expect(asyncFn).not.toHaveBeenCalled()

      // Advance time fully
      vi.advanceTimersByTime(100)
      const result = await promise

      expect(result).toBe(true)
      expect(asyncFn).toHaveBeenCalled()
      expect(store.uno.valid.validating['debounceField']).toBeUndefined()

      vi.useRealTimers()
    })
  })

  describe('Form Validation', () => {
    it('should validate all fields in form', async () => {
      document.body.innerHTML = `
        <form id="testForm">
          <input name="user" value="" />
          <input name="email" value="bad-email" />
        </form>
      `
      const form = document.getElementById('testForm')
      // Explicitly mock scrollIntoView/focus on elements to avoid prototype issues
      form.querySelectorAll('input').forEach((el) => {
        el.scrollIntoView = vi.fn()
        el.focus = vi.fn()
      })

      store.uno.valid.addRules('user', ['required'])
      store.uno.valid.addRules('email', ['email'])

      const isValid = await store.uno.valid.validateForm(form)

      expect(isValid).toBe(false)
      expect(store.uno.valid.errors['user']).toBeTruthy()
      expect(store.uno.valid.errors['email']).toBeTruthy()
    })
  })

  describe('HTMX Integration', () => {
    it('should intercept htmx:confirm and validate', async () => {
      document.body.innerHTML = `
        <form id="htmxForm" x-uno-valid>
          <input name="field" value="" />
        </form>
      `
      const form = document.getElementById('htmxForm')
      store.uno.valid.addRules('field', ['required'])

      const preventDefault = vi.fn()
      const issueRequest = vi.fn()

      // Use legacy createEvent to ensure maximum compatibility with JSDOM
      const evt = form.ownerDocument.createEvent('CustomEvent')
      evt.initCustomEvent('htmx:confirm', true, true, { issueRequest, elt: form })

      // We can't easily mock preventDefault on the instance created this way if it's strictly checked
      // But we can check if defaultPrevented becomes true
      // Or we can try to spy on it?
      // Let's just mock it for now as before, hoping createEvent makes a valid event
      Object.defineProperty(evt, 'preventDefault', { value: preventDefault })

      // Mock scrollIntoView/focus on inputs
      form.querySelectorAll('input').forEach((el) => {
        el.scrollIntoView = vi.fn()
        el.focus = vi.fn()
      })

      form.dispatchEvent(evt) // We need to await the event handler?
      // The event handler is async but dispatchEvent is sync.
      // We can't easily await the handler.
      // But since we use await inside the handler, the checks after dispatchEvent might run before handler finishes?
      // Actually, since we are in JSDOM/Node, the microtask queue handling matters.

      form.dispatchEvent(evt)

      // Wait for promises to resolve
      await new Promise(process.nextTick)

      expect(preventDefault).toHaveBeenCalled()
      expect(issueRequest).not.toHaveBeenCalled()
      expect(store.uno.valid.errors['field']).toBeTruthy()
    })

    it('should intercept htmx:confirm when only fields carry x-uno-valid', async () => {
      document.body.innerHTML = `
        <form id="htmxForm">
          <input name="field" value="" x-uno-valid="['required']" />
        </form>
      `
      const form = document.getElementById('htmxForm')
      store.uno.valid.addRules('field', ['required'])

      const preventDefault = vi.fn()
      const issueRequest = vi.fn()

      const evt = form.ownerDocument.createEvent('CustomEvent')
      evt.initCustomEvent('htmx:confirm', true, true, { issueRequest, elt: form })
      Object.defineProperty(evt, 'preventDefault', { value: preventDefault })

      form.querySelectorAll('input').forEach((el) => {
        el.scrollIntoView = vi.fn()
        el.focus = vi.fn()
      })

      form.dispatchEvent(evt)
      await new Promise(process.nextTick)

      expect(preventDefault).toHaveBeenCalled()
      expect(issueRequest).not.toHaveBeenCalled()
      expect(store.uno.valid.errors['field']).toBeTruthy()
    })

    it('should allow request if validation passes', async () => {
      document.body.innerHTML = `
        <form id="htmxForm" x-uno-valid>
          <input name="field" value="valid" />
        </form>
      `
      const form = document.getElementById('htmxForm')
      store.uno.valid.addRules('field', ['required'])

      const preventDefault = vi.fn()
      const issueRequest = vi.fn()

      // Use legacy createEvent
      const evt = form.ownerDocument.createEvent('CustomEvent')
      evt.initCustomEvent('htmx:confirm', true, true, { issueRequest, elt: form })
      Object.defineProperty(evt, 'preventDefault', { value: preventDefault })

      // Mock scrollIntoView/focus just in case validation fails and tries to focus
      form.querySelectorAll('input').forEach((el) => {
        el.scrollIntoView = vi.fn()
        el.focus = vi.fn()
      })

      form.dispatchEvent(evt)

      // Wait for promises
      await new Promise(process.nextTick)

      // Expectation: preventDefault called (to pause), then issueRequest called
      expect(preventDefault).toHaveBeenCalled()
      expect(issueRequest).toHaveBeenCalled()
    })

    it('should apply server-side errors from 422 responses', () => {
      const showErrorsSpy = vi.spyOn(store.uno.valid, 'showErrors')

      const evt = document.createEvent('CustomEvent')
      evt.initCustomEvent('htmx:afterRequest', true, true, {
        xhr: {
          status: 422,
          responseText: JSON.stringify({ errors: { email: 'Email already taken' } }),
          getResponseHeader: () => null
        }
      })
      document.body.dispatchEvent(evt)

      expect(showErrorsSpy).toHaveBeenCalledWith({ email: 'Email already taken' })
      expect(store.uno.valid.errors['email']).toBe('Email already taken')
    })

    it('should ignore 422 responses with non-JSON bodies', () => {
      const showErrorsSpy = vi.spyOn(store.uno.valid, 'showErrors')

      const evt = document.createEvent('CustomEvent')
      evt.initCustomEvent('htmx:afterRequest', true, true, {
        xhr: {
          status: 422,
          responseText: '<html>Unprocessable</html>',
          getResponseHeader: () => null
        }
      })
      document.body.dispatchEvent(evt)

      expect(showErrorsSpy).not.toHaveBeenCalled()
    })

    it('should ignore non-422 responses', () => {
      const showErrorsSpy = vi.spyOn(store.uno.valid, 'showErrors')

      const evt = document.createEvent('CustomEvent')
      evt.initCustomEvent('htmx:afterRequest', true, true, {
        xhr: {
          status: 200,
          responseText: JSON.stringify({ errors: { email: 'should not apply' } }),
          getResponseHeader: () => null
        }
      })
      document.body.dispatchEvent(evt)

      expect(showErrorsSpy).not.toHaveBeenCalled()
    })
  })

  describe('Auto Error Generation', () => {
    it('should automatically create error element if not exists', async () => {
      document.body.innerHTML = `
        <form>
          <input name="auto" />
        </form>
      `
      const input = document.querySelector('input')

      // Mock Uno.options.valid default
      global.Uno.options.valid = {}

      const directiveFn = directives['uno-valid']
      expect(directiveFn).toBeDefined()

      // Call the directive function with mock params
      const el = input
      const expression = "['required']"
      const evaluateLater = (expr) => (cb) => cb(['required']) // mock evaluateLater returning array via callback

      const cleanup = () => {}
      directiveFn(el, { expression }, { evaluateLater, cleanup })

      // Check if error element was created
      const errorEl = input.nextElementSibling
      expect(errorEl).not.toBeNull()
      expect(errorEl.tagName).toBe('DIV')
      expect(errorEl.getAttribute('x-text')).toBe("$store.uno.valid.errors['auto']")
      expect(errorEl.getAttribute('x-show')).toBe("$store.uno.valid.errors['auto']")
      // Default class includes mt-1 (and text-red-500 etc)
      // Note: we can't easily check default class without full Uno init, but we can check if it has a class
      expect(errorEl.className).toBeTruthy()
      expect(Alpine.initTree).toHaveBeenCalledWith(errorEl)
    })

    it('should NOT create error element if one already exists', async () => {
      document.body.innerHTML = `
        <form>
          <input name="exists" />
          <div x-text="$store.uno.valid.errors['exists']"></div>
        </form>
      `
      const input = document.querySelector('input')
      const originalError = document.querySelector('div')

      // Mock Uno.options.valid default
      global.Uno.options.valid = {}

      const directiveFn = directives['uno-valid']
      const evaluateLater = (expr) => (cb) => cb(['required'])
      const cleanup = () => {}
      directiveFn(input, { expression: "['required']" }, { evaluateLater, cleanup })

      // Should still have only one error element (the original one)
      expect(input.nextElementSibling).toBe(originalError)
      // No new sibling inserted
      expect(input.parentElement.children.length).toBe(2)
      // initTree should NOT be called for existing element
      expect(Alpine.initTree).not.toHaveBeenCalled()
    })

    it('should use custom createErrorElement if configured', async () => {
      document.body.innerHTML = `
        <form>
          <input name="custom" />
        </form>
      `
      const input = document.querySelector('input')

      // Mock Uno config
      global.Uno.options.valid = {
        createErrorElement: (fieldName, errorKey) => {
          const span = document.createElement('span')
          span.id = 'custom-error-' + fieldName
          span.textContent = 'Error Place'
          return span
        }
      }

      const directiveFn = directives['uno-valid']
      const evaluateLater = (expr) => (cb) => cb(['required'])
      const cleanup = () => {}
      directiveFn(input, { expression: "['required']" }, { evaluateLater, cleanup })

      const errorEl = input.nextElementSibling
      expect(errorEl.tagName).toBe('SPAN')
      expect(errorEl.id).toBe('custom-error-custom')
      expect(Alpine.initTree).toHaveBeenCalledWith(errorEl)
    })

    it('should respect autoError: false config', async () => {
      document.body.innerHTML = `
        <form>
          <input name="noauto" />
        </form>
      `
      const input = document.querySelector('input')

      global.Uno.options.valid = {
        autoError: false
      }

      const directiveFn = directives['uno-valid']
      const evaluateLater = (expr) => (cb) => cb(['required'])
      const cleanup = () => {}
      directiveFn(input, { expression: "['required']" }, { evaluateLater, cleanup })

      const errorEl = input.nextElementSibling
      expect(errorEl).toBeNull()
    })
  })

  describe('Length Rule Features', () => {
    it('should set maxlength attribute when stopOnMax is true', () => {
      document.body.innerHTML = `
        <form>
          <input name="limited" />
        </form>
      `
      const input = document.querySelector('input')
      global.Uno.options.valid = {}

      const directiveFn = directives['uno-valid']
      const evaluateLater = (expr) => (cb) => cb([{ type: 'length', max: 10, stopOnMax: true }])
      const cleanup = () => {}

      directiveFn(input, { expression: '...' }, { evaluateLater, cleanup })

      expect(input.getAttribute('maxlength')).toBe('10')
    })

    it('should NOT set maxlength attribute when stopOnMax is false (default)', () => {
      document.body.innerHTML = `
        <form>
          <input name="unlimited" />
        </form>
      `
      const input = document.querySelector('input')
      global.Uno.options.valid = {}

      const directiveFn = directives['uno-valid']
      const evaluateLater = (expr) => (cb) => cb([{ type: 'length', max: 10 }])
      const cleanup = () => {}

      directiveFn(input, { expression: '...' }, { evaluateLater, cleanup })

      expect(input.hasAttribute('maxlength')).toBe(false)
    })

    it('should create counter element when showCounter is true', () => {
      document.body.innerHTML = `
        <form>
          <div class="relative">
            <input name="counter" />
          </div>
        </form>
      `
      const input = document.querySelector('input')
      global.Uno.options.valid = {}

      const directiveFn = directives['uno-valid']
      const evaluateLater = (expr) => (cb) => cb([{ type: 'length', max: 20, showCounter: true }])
      const cleanup = () => {}

      directiveFn(input, { expression: '...' }, { evaluateLater, cleanup })

      const counter = document.getElementById('counter-counter')
      expect(counter).not.toBeNull()
      expect(counter.textContent).toBe('0/20')
      expect(counter.tagName).toBe('SPAN')
    })

    it('should update counter on input', () => {
      document.body.innerHTML = `
        <form>
          <div class="relative">
            <input name="updating" />
          </div>
        </form>
      `
      const input = document.querySelector('input')
      global.Uno.options.valid = {}

      const directiveFn = directives['uno-valid']
      const evaluateLater = (expr) => (cb) => cb([{ type: 'length', max: 5, showCounter: true }])
      const cleanup = () => {}

      directiveFn(input, { expression: '...' }, { evaluateLater, cleanup })

      const counter = document.getElementById('updating-counter')

      // Simulate input
      input.value = 'abc'
      input.dispatchEvent(new Event('input'))

      expect(counter.textContent).toBe('3/5')
      expect(counter.classList.contains('text-gray-600')).toBe(true)

      // Exceed max
      input.value = 'abcdef'
      input.dispatchEvent(new Event('input'))

      expect(counter.textContent).toBe('6/5')
      expect(counter.classList.contains('text-red-600')).toBe(true)
    })

    it('should anchor the counter in a plugin-managed wrapper independent of sibling flow content', () => {
      document.body.innerHTML = `
        <form>
          <div class="mb-4">
            <label for="bio">Bio</label>
            <input name="bio" id="bio" />
            <p class="help">help text that grows the field container</p>
          </div>
        </form>
      `
      const input = document.querySelector('input')
      global.Uno.options.valid = {}

      const directiveFn = directives['uno-valid']
      const evaluateLater = (expr) => (cb) => cb([{ type: 'length', max: 500, showCounter: true }])
      const cleanup = () => {}

      directiveFn(input, { expression: '...' }, { evaluateLater, cleanup })

      // Input owns a tight positioned containing block
      const wrapper = input.parentElement
      expect(wrapper.hasAttribute('data-uno-counter-wrapper')).toBe(true)
      expect(wrapper.style.position).toBe('relative')

      // Counter shares that containing block instead of the field container
      const counter = document.getElementById('bio-counter')
      expect(counter.parentElement).toBe(wrapper)

      // Sibling flow content stays outside so it can never shift the anchor
      const field = wrapper.parentElement
      expect(field.querySelector('p.help').parentElement).toBe(field)
    })

    it('should keep the auto error element outside the counter wrapper', () => {
      document.body.innerHTML = `
        <form>
          <div>
            <input name="witherror" />
          </div>
        </form>
      `
      const input = document.querySelector('input')
      global.Uno.options.valid = {}

      const directiveFn = directives['uno-valid']
      const evaluateLater = (expr) => (cb) => cb(['required', { type: 'length', max: 10, showCounter: true }])
      const cleanup = () => {}

      directiveFn(input, { expression: '...' }, { evaluateLater, cleanup })

      const wrapper = input.parentElement
      expect(wrapper.hasAttribute('data-uno-counter-wrapper')).toBe(true)
      // Error element stays in normal flow after the wrapper: showing an error
      // can no longer stretch the counter's containing block
      const errorEl = wrapper.nextElementSibling
      expect(errorEl.getAttribute('x-text')).toBe("$store.uno.valid.errors['witherror']")
    })

    it('should restore the original DOM shape on cleanup', () => {
      document.body.innerHTML = `
        <form>
          <div id="field">
            <input name="restored" />
          </div>
        </form>
      `
      const input = document.querySelector('input')
      const field = document.getElementById('field')
      global.Uno.options.valid = {}

      const directiveFn = directives['uno-valid']
      const evaluateLater = (expr) => (cb) => cb([{ type: 'length', max: 10, showCounter: true }])
      const cleanups = []
      const cleanup = (cb) => {
        cleanups.push(cb)
      }

      directiveFn(input, { expression: '...' }, { evaluateLater, cleanup })
      expect(input.parentElement.hasAttribute('data-uno-counter-wrapper')).toBe(true)

      cleanups.forEach((cb) => cb())
      expect(input.parentElement).toBe(field)
      expect(field.querySelector('[data-uno-counter-wrapper]')).toBeNull()
      expect(document.getElementById('restored-counter')).toBeNull()
    })

    it('should fall back to parent anchoring inside flex containers', () => {
      document.body.innerHTML = `
        <form>
          <div id="row" style="display: flex; position: static;">
            <input name="flexed" />
          </div>
        </form>
      `
      const input = document.querySelector('input')
      const row = document.getElementById('row')
      global.Uno.options.valid = {}

      const directiveFn = directives['uno-valid']
      const evaluateLater = (expr) => (cb) => cb([{ type: 'length', max: 10, showCounter: true }])
      const cleanup = () => {}

      directiveFn(input, { expression: '...' }, { evaluateLater, cleanup })

      // No wrapper: the input keeps its role as flex item, parent is promoted
      expect(input.parentElement).toBe(row)
      expect(row.style.position).toBe('relative')
      expect(document.getElementById('flexed-counter').parentElement).toBe(row)
    })

    it('should reuse the managed wrapper and counter on re-init', () => {
      document.body.innerHTML = `
        <form>
          <div>
            <input name="reinit" />
          </div>
        </form>
      `
      const input = document.querySelector('input')
      global.Uno.options.valid = {}

      const directiveFn = directives['uno-valid']
      const evaluateLater = (expr) => (cb) => cb([{ type: 'length', max: 10, showCounter: true }])
      const cleanup = () => {}

      directiveFn(input, { expression: '...' }, { evaluateLater, cleanup })
      const wrapper = input.parentElement
      // Simulate Alpine re-running the directive (e.g. after HTMX swap-in)
      directiveFn(input, { expression: '...' }, { evaluateLater, cleanup })

      expect(input.parentElement).toBe(wrapper)
      expect(document.querySelectorAll('[data-uno-counter-wrapper]').length).toBe(1)
      expect(document.querySelectorAll('#reinit-counter').length).toBe(1)
    })
  })

  describe('A11y: Error Element Wiring (A1)', () => {
    const runDirective = (input, config = {}) => {
      global.Uno.options.valid = config
      const directiveFn = directives['uno-valid']
      const evaluateLater = () => (cb) => cb(['required'])
      directiveFn(input, { expression: "['required']" }, { evaluateLater, cleanup: () => {} })
    }

    it('auto-generated error element gets role="alert", a stable id, and wires aria-describedby', () => {
      document.body.innerHTML = `<form><input name="title" /></form>`
      const input = document.querySelector('input')

      runDirective(input)

      const errorEl = input.nextElementSibling
      expect(errorEl.getAttribute('role')).toBe('alert')
      expect(errorEl.id).toBe('title-error')
      expect(input.getAttribute('aria-describedby')).toBe('title-error')
    })

    it('preserves an author-provided error element id/role and still wires describedby', () => {
      document.body.innerHTML = `
        <form>
          <input name="handle" />
          <div id="my-error" role="status" x-text="$store.uno.valid.errors['handle']"></div>
        </form>`
      const input = document.querySelector('input')

      runDirective(input)

      const errorEl = document.getElementById('my-error')
      expect(errorEl.getAttribute('role')).toBe('status')
      expect(errorEl.id).toBe('my-error')
      expect(input.getAttribute('aria-describedby')).toBe('my-error')
      expect(Alpine.initTree).not.toHaveBeenCalled()
    })

    it('custom createErrorElement keeps its id and gains role="alert"', () => {
      document.body.innerHTML = `<form><input name="custom" /></form>`
      const input = document.querySelector('input')

      runDirective(input, {
        createErrorElement: (fieldName) => {
          const span = document.createElement('span')
          span.id = 'custom-error-' + fieldName
          return span
        }
      })

      const errorEl = input.nextElementSibling
      expect(errorEl.id).toBe('custom-error-custom')
      expect(errorEl.getAttribute('role')).toBe('alert')
      expect(input.getAttribute('aria-describedby')).toBe('custom-error-custom')
    })

    it('appends to an existing aria-describedby without clobbering prior targets', () => {
      document.body.innerHTML = `<form><input name="bio" aria-describedby="bio-help" /></form>`
      const input = document.querySelector('input')

      runDirective(input)

      expect(input.getAttribute('aria-describedby')).toBe('bio-help bio-error')
    })

    it('does not wire describedby when autoError is disabled and no error element exists', () => {
      document.body.innerHTML = `<form><input name="noauto" /></form>`
      const input = document.querySelector('input')

      runDirective(input, { autoError: false })

      expect(input.nextElementSibling).toBeNull()
      expect(input.hasAttribute('aria-describedby')).toBe(false)
    })
  })

  describe('A11y: aria-invalid Sync (A1)', () => {
    it('sets aria-invalid="true" on failure and removes it on success', async () => {
      document.body.innerHTML = `<form><input name="req" /></form>`
      const input = document.querySelector('input')
      store.uno.valid.addRules('req', ['required'])

      await store.uno.valid.validate('req', '')
      expect(input.getAttribute('aria-invalid')).toBe('true')

      await store.uno.valid.validate('req', 'ok')
      expect(input.hasAttribute('aria-invalid')).toBe(false)
    })

    it('sets aria-invalid on fields delivered via 422 showErrors', () => {
      document.body.innerHTML = `<form><input name="email" /></form>`
      const input = document.querySelector('input')

      store.uno.valid.showErrors({ email: 'Email already taken' })

      expect(input.getAttribute('aria-invalid')).toBe('true')
    })

    it('clears aria-invalid through hideError and clearErrors', () => {
      document.body.innerHTML = `<form><input name="a" /><input name="b" /></form>`
      const a = document.querySelector('[name="a"]')
      const b = document.querySelector('[name="b"]')

      store.uno.valid.showErrors({ a: 'err a', b: 'err b' })
      expect(a.getAttribute('aria-invalid')).toBe('true')

      store.uno.valid.hideError('a')
      expect(a.hasAttribute('aria-invalid')).toBe(false)
      expect(b.getAttribute('aria-invalid')).toBe('true')

      store.uno.valid.clearErrors()
      expect(b.hasAttribute('aria-invalid')).toBe(false)
    })
  })
})
