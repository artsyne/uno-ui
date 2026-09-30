import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { createBoxStore } from '../../src/components/box/box.store'
import { Uno } from '../../src/core/uno'

// Mock dependencies
vi.mock('../../src/core/uno', () => ({
  Uno: {
    getConfig: vi.fn(),
    getStyle: vi.fn((comp, key) => `${comp}-${key}`),
    dom: {
      ensureContainer: vi.fn(() => document.createElement('div')),
      sanitize: vi.fn((html) => html.replace(/<script[^>]*>.*?<\/script>/gi, ''))
    },
    options: {
      ui: {
        layers: { box: 5000 }
      }
    }
  }
}))

// Mock Alpine
const Alpine = {
  store: vi.fn(),
  nextTick: vi.fn((fn) => setTimeout(fn, 0)),
  initTree: vi.fn()
}

describe('Box Component', () => {
  let boxStore
  let mockStore

  beforeEach(() => {
    vi.clearAllMocks()
    vi.useFakeTimers()

    // Reset store mock
    mockStore = {}
    Alpine.store.mockImplementation((name, value) => {
      if (value) {
        mockStore[name] = value
        if (value.box?.init) value.box.init()
      }
      return mockStore[name]
    })

    // Setup Uno.getConfig default behavior
    Uno.getConfig.mockImplementation((key, opts) => ({
      defaultWidth: 400,
      defaultPosition: 'center',
      okText: 'OK',
      cancelText: 'Cancel',
      closeOnEsc: false,
      closeOnOverlayClick: false,
      ...opts
    }))

    // Create store
    Alpine.store('uno', { box: createBoxStore(Alpine) })
    boxStore = Alpine.store('uno').box
  })

  afterEach(() => {
    vi.useRealTimers()
  })

  describe('Initialization', () => {
    it('should initialize with empty items array', () => {
      expect(boxStore.items).toEqual([])
    })

    it('should call ensureContainer on open', () => {
      boxStore.open('dialog')
      expect(Uno.dom.ensureContainer).toHaveBeenCalled()
    })
  })

  describe('Auto Focus', () => {
    it('should auto-focus first input for static content', () => {
      const originalRaf = global.requestAnimationFrame
      global.requestAnimationFrame = (cb) => setTimeout(cb, 0)
      document.body.innerHTML = `
        <div id="uno-box-containerSTATIC"></div>
      `
      const html = `<input type="text" id="firstInput"><input type="text" id="secondInput">`
      boxStore.open({ boxId: 'STATIC', content: html, autoFocus: true })
      const container = document.getElementById('uno-box-containerSTATIC')
      container.innerHTML = html
      const first = document.getElementById('firstInput')
      expect(first).toBeDefined()
      vi.runAllTimers()
      global.requestAnimationFrame = originalRaf
      expect(document.activeElement).toBe(first)
    })

    it('should focus the container (not inputs) when autoFocus=false', () => {
      const originalRaf = global.requestAnimationFrame
      global.requestAnimationFrame = (cb) => setTimeout(cb, 0)
      document.body.innerHTML = `<div id="uno-box-containerNOFOCUS"></div>`
      const html = `<input type="text" id="nf1"><input type="text" id="nf2">`
      boxStore.open({ boxId: 'NOFOCUS', content: html, autoFocus: false })
      const container = document.getElementById('uno-box-containerNOFOCUS')
      container.innerHTML = html
      const first = document.getElementById('nf1')
      vi.runAllTimers()
      global.requestAnimationFrame = originalRaf
      // autoFocus=false skips input focus but focus must stay in the dialog
      expect(document.activeElement).not.toBe(first)
      expect(document.activeElement).toBe(container)
    })

    it('should fall back to the container when content has no inputable element', () => {
      const originalRaf = global.requestAnimationFrame
      global.requestAnimationFrame = (cb) => setTimeout(cb, 0)
      document.body.innerHTML = `<div id="uno-box-containerNOINPUT"></div>`
      const html = `<p>Read only content</p>`
      boxStore.open({ boxId: 'NOINPUT', content: html, autoFocus: true })
      const container = document.getElementById('uno-box-containerNOINPUT')
      container.innerHTML = html
      vi.runAllTimers()
      global.requestAnimationFrame = originalRaf
      expect(document.activeElement).toBe(container)
    })

    it('should fall back to the title heading when titled dialog has no inputable element', () => {
      const originalRaf = global.requestAnimationFrame
      global.requestAnimationFrame = (cb) => setTimeout(cb, 0)
      document.body.innerHTML = `<div id="uno-box-containerTITLEFB"></div>`
      const html = `<h3 id="uno-box-title-TITLEFB">Hello</h3><p>Read only</p>`
      boxStore.open({ boxId: 'TITLEFB', title: 'Hello', content: html, autoFocus: true })
      const container = document.getElementById('uno-box-containerTITLEFB')
      container.innerHTML = html
      const title = document.getElementById('uno-box-title-TITLEFB')
      vi.runAllTimers()
      global.requestAnimationFrame = originalRaf
      expect(document.activeElement).toBe(title)
    })

    it('should skip inputs hidden by a display:none ancestor wrapper', () => {
      const originalRaf = global.requestAnimationFrame
      global.requestAnimationFrame = (cb) => setTimeout(cb, 0)
      document.body.innerHTML = `<div id="uno-box-containerWRAP"></div>`
      const html = `
        <div style="display:none"><input type="text" id="wrapped"></div>
        <input type="text" id="shown">
      `
      boxStore.open({ boxId: 'WRAP', content: html, autoFocus: true })
      const container = document.getElementById('uno-box-containerWRAP')
      container.innerHTML = html
      const shown = document.getElementById('shown')
      vi.runAllTimers()
      global.requestAnimationFrame = originalRaf
      expect(document.activeElement).toBe(shown)
    })

    it('should skip hidden and disabled inputs', () => {
      const originalRaf = global.requestAnimationFrame
      global.requestAnimationFrame = (cb) => setTimeout(cb, 0)
      document.body.innerHTML = `<div id="uno-box-containerSKIP"></div>`
      const html = `
        <input type="text" id="hidden" style="display:none">
        <input type="text" id="disabled" disabled>
        <input type="text" id="visible">
      `
      boxStore.open({ boxId: 'SKIP', content: html, autoFocus: true })
      const container = document.getElementById('uno-box-containerSKIP')
      container.innerHTML = html
      const visible = document.getElementById('visible')
      vi.runAllTimers()
      global.requestAnimationFrame = originalRaf
      expect(document.activeElement).toBe(visible)
    })
  })
  describe('open() - Core Method', () => {
    it('should create a box item with default config', () => {
      const item = boxStore.open('dialog')

      expect(item).toBeDefined()
      expect(item.type).toBe('dialog')
      expect(boxStore.items).toHaveLength(1)
      expect(boxStore.items[0]).toBe(item)
    })

    it('should accept type string as first parameter', () => {
      const item = boxStore.open('confirm')
      expect(item.type).toBe('confirm')
    })

    it('should accept options object as first parameter', () => {
      const item = boxStore.open({ type: 'dialog', title: 'Test' })
      expect(item.type).toBe('dialog')
      expect(item.title).toBe('Test')
    })

    it('should pass button classes to item', () => {
      const item = boxStore.open({
        okButtonClass: 'custom-ok',
        cancelButtonClass: 'custom-cancel'
      })
      expect(item.okButtonClass).toBe('custom-ok')
      expect(item.cancelButtonClass).toBe('custom-cancel')
    })

    it('should use 4-level configuration merge', () => {
      Uno.getConfig.mockImplementation((key, opts) => ({
        defaultWidth: 400,
        title: 'Merged Title',
        customOption: true,
        ...opts
      }))

      const item = boxStore.open({ title: 'Instance Title' })
      expect(Uno.getConfig).toHaveBeenCalledWith('box', { title: 'Instance Title' })
      expect(item.title).toBe('Instance Title')
    })

    it('should sanitize content by default', () => {
      const item = boxStore.open({
        content: '<p>Safe</p><script>alert(1)</script>'
      })

      expect(Uno.dom.sanitize).toHaveBeenCalledWith('<p>Safe</p><script>alert(1)</script>')
      expect(item.content).not.toContain('<script>')
    })

    it('should set isOpen to false initially', () => {
      const item = boxStore.open('dialog')
      expect(item.isOpen).toBe(false)
    })

    it('should set isOpen to true after nextTick', () => {
      const item = boxStore.open('dialog')
      expect(item.isOpen).toBe(false)

      vi.runAllTimers()
      expect(item.isOpen).toBe(true)
    })

    it('should generate unique ID if not provided', () => {
      const item1 = boxStore.open('dialog')
      const item2 = boxStore.open('dialog')
      expect(item1.id).not.toBe(item2.id)
    })

    it('should use custom boxId if provided', () => {
      const item = boxStore.open({ boxId: 'custom-id' })
      expect(item.id).toBe('custom-id')
    })
  })

  describe('Button Visibility Logic', () => {
    it('should show both buttons for confirm type by default', () => {
      const item = boxStore.open('confirm')
      expect(item.showOkButton).toBe(true)
      expect(item.showCancelButton).toBe(true)
    })

    it('should hide both buttons for dialog type by default', () => {
      const item = boxStore.open('dialog')
      expect(item.showOkButton).toBe(false)
      expect(item.showCancelButton).toBe(false)
    })

    it('should respect explicit button visibility for confirm', () => {
      const item = boxStore.open('confirm', { showCancelButton: false })
      expect(item.showOkButton).toBe(true)
      expect(item.showCancelButton).toBe(false)
    })

    it('should respect explicit button visibility for dialog', () => {
      const item = boxStore.open('dialog', { showOkButton: true })
      expect(item.showOkButton).toBe(true)
    })
  })

  describe('alert() - Convenience Method', () => {
    it('should create dialog with OK button', () => {
      const item = boxStore.alert('Alert Message')
      expect(item.type).toBe('dialog')
      expect(item.content).toContain('Alert Message')
      expect(item.showOkButton).toBe(true)
    })

    it('should accept title as string option', () => {
      const item = boxStore.alert('Message', 'Alert Title')
      expect(item.title).toBe('Alert Title')
    })

    it('should accept options object', () => {
      const item = boxStore.alert('Message', { title: 'Title', width: 500 })
      expect(item.title).toBe('Title')
      expect(item.width).toBe(500)
    })
  })

  describe('confirm() - Promise-based', () => {
    it('should return a Promise', () => {
      const result = boxStore.confirm('Confirm?')
      expect(result).toBeInstanceOf(Promise)
    })

    it('should resolve true when OK is clicked', async () => {
      const promise = boxStore.confirm('Confirm?')
      const item = boxStore.items[0]

      // Simulate OK click
      if (item.onOk) item.onOk()

      const result = await promise
      expect(result).toBe(true)
    })

    it('should resolve false when Cancel is clicked', async () => {
      const promise = boxStore.confirm('Confirm?')
      const item = boxStore.items[0]

      // Simulate Cancel click
      if (item.onCancel) item.onCancel()

      const result = await promise
      expect(result).toBe(false)
    })

    it('should resolve false when closed without choice', async () => {
      const promise = boxStore.confirm('Confirm?')
      const item = boxStore.items[0]

      // Simulate close
      if (item.onClose) item.onClose()

      const result = await promise
      expect(result).toBe(false)
    })

    it('should compose a user-supplied onOk with the Promise resolver', async () => {
      const onOk = vi.fn()
      const promise = boxStore.confirm('Confirm?', { onOk })
      const item = boxStore.items[0]

      const ret = item.onOk()
      expect(onOk).toHaveBeenCalledTimes(1)
      expect(ret).toBeUndefined() // allow the dialog to close (handleOk semantics)

      await expect(promise).resolves.toBe(true)
    })

    it('should keep the dialog open and the promise pending when user onOk returns false', async () => {
      const onOk = vi.fn(() => false)
      const promise = boxStore.confirm('Confirm?', { onOk })
      const item = boxStore.items[0]

      const ret = item.onOk()
      expect(onOk).toHaveBeenCalledTimes(1)
      expect(ret).toBe(false) // keep open, like handleOk

      // Promise must remain pending while the dialog stays open
      let settled = false
      promise.then(() => {
        settled = true
      })
      await new Promise(process.nextTick)
      expect(settled).toBe(false)
    })

    it('should resolve after an async user onOk completes', async () => {
      const onOk = vi.fn(async () => undefined)
      const promise = boxStore.confirm('Confirm?', { onOk })
      const item = boxStore.items[0]

      item.onOk()
      expect(onOk).toHaveBeenCalledTimes(1)
      await expect(promise).resolves.toBe(true)
    })

    it('should compose a user-supplied onCancel with the Promise resolver', async () => {
      const onCancel = vi.fn()
      const promise = boxStore.confirm('Confirm?', { onCancel })
      const item = boxStore.items[0]

      item.onCancel()
      expect(onCancel).toHaveBeenCalledTimes(1)
      await expect(promise).resolves.toBe(false)
    })
  })

  describe('prompt() - Input Dialog', () => {
    it('should return a Promise', () => {
      const result = boxStore.prompt('Enter name:')
      expect(result).toBeInstanceOf(Promise)
    })

    it('should create confirm dialog with input field', () => {
      boxStore.prompt('Enter name:', 'default')
      const item = boxStore.items[0]

      expect(item.type).toBe('confirm')
      expect(item.content).toContain('Enter name:')
      expect(item.content).toContain('type="text"')
      expect(item.content).toContain('default')
    })

    it('should sanitize message and defaultValue', () => {
      boxStore.prompt('<script>alert(1)</script>', '<img onerror="alert(2)">')
      expect(Uno.dom.sanitize).toHaveBeenCalledTimes(2)
    })

    it('should resolve with input value when OK is clicked', async () => {
      const promise = boxStore.prompt('Name:', 'John')
      const item = boxStore.items[0]

      // Extract input ID from content
      const match = item.content.match(/id="([^"]+)"/)
      const inputId = match ? match[1] : null

      // Create mock input element
      const input = document.createElement('input')
      input.id = inputId
      input.value = 'Test Value'
      document.body.appendChild(input)

      // Simulate OK click
      if (item.onOk) item.onOk()

      const result = await promise
      expect(result).toBe('Test Value')

      document.body.removeChild(input)
    })

    it('should resolve null when Cancel is clicked', async () => {
      const promise = boxStore.prompt('Name:')
      const item = boxStore.items[0]

      // Simulate Cancel click
      if (item.onCancel) item.onCancel()

      const result = await promise
      expect(result).toBeNull()
    })
  })

  describe('dialog() - Alias Method', () => {
    it('should be an alias for open()', () => {
      const item = boxStore.dialog({ title: 'Dialog Title' })
      expect(item.type).toBe('dialog')
      expect(item.title).toBe('Dialog Title')
    })
  })

  describe('close() - by ID', () => {
    it('should close specific box by ID', () => {
      const item = boxStore.open({ boxId: 'test-box' })
      vi.runAllTimers() // Open the box

      boxStore.close('test-box')
      expect(item.isOpen).toBe(false)
    })

    it('should do nothing if ID not found', () => {
      boxStore.open('dialog')
      expect(() => boxStore.close('non-existent')).not.toThrow()
    })
  })

  describe('closeAll()', () => {
    it('should close all boxes', () => {
      const onClose1 = vi.fn()
      const onClose2 = vi.fn()

      boxStore.open({ onClose: onClose1 })
      boxStore.open({ onClose: onClose2 })

      expect(boxStore.items).toHaveLength(2)

      boxStore.closeAll()

      expect(boxStore.items).toHaveLength(0)
      expect(onClose1).toHaveBeenCalled()
      expect(onClose2).toHaveBeenCalled()
    })

    it('should call onBeforeClose for all boxes', () => {
      const onBeforeClose1 = vi.fn()
      const onBeforeClose2 = vi.fn()

      boxStore.open({ onBeforeClose: onBeforeClose1 })
      boxStore.open({ onBeforeClose: onBeforeClose2 })

      boxStore.closeAll()

      expect(onBeforeClose1).toHaveBeenCalled()
      expect(onBeforeClose2).toHaveBeenCalled()
    })
  })

  describe('Item.close() - Instance Method', () => {
    it('should call onBeforeClose callback', () => {
      const onBeforeClose = vi.fn()
      const item = boxStore.open({ onBeforeClose })
      vi.runAllTimers()

      item.close()
      expect(onBeforeClose).toHaveBeenCalled()
    })

    it('should set isOpen to false', () => {
      const item = boxStore.open('dialog')
      vi.runAllTimers()
      expect(item.isOpen).toBe(true)

      item.close()
      expect(item.isOpen).toBe(false)
    })

    it('should remove item from store after transition', () => {
      const item = boxStore.open('dialog')
      vi.runAllTimers()

      expect(boxStore.items).toHaveLength(1)

      item.close()
      vi.advanceTimersByTime(350) // Safety timeout

      expect(boxStore.items).toHaveLength(0)
    })

    it('should call onClose callback after cleanup', () => {
      const onClose = vi.fn()
      const item = boxStore.open({ onClose })
      vi.runAllTimers()

      item.close()
      vi.advanceTimersByTime(350)

      expect(onClose).toHaveBeenCalled()
    })
  })

  describe('Callbacks', () => {
    it('should call onContentLoad if provided', () => {
      const onContentLoad = vi.fn()
      const item = boxStore.open({ onContentLoad })

      const mockEl = document.createElement('div')
      boxStore.onHtmxContentLoaded(mockEl, item)

      expect(onContentLoad).toHaveBeenCalledWith(mockEl)
    })

    it('should initialize Alpine tree in loaded content', () => {
      window.Alpine = { initTree: vi.fn() }
      const item = boxStore.open({})
      const mockEl = document.createElement('div')

      boxStore.onHtmxContentLoaded(mockEl, item)

      expect(window.Alpine.initTree).toHaveBeenCalledWith(mockEl)
      delete window.Alpine
    })
  })

  describe('HTMX Support', () => {
    it('should pass HTMX options to item', () => {
      const item = boxStore.open({
        hxUrl: '/api/content',
        hxMethod: 'POST',
        hxHeaders: { 'X-Custom': 'value' },
        hxTrigger: 'click'
      })

      expect(item.hxUrl).toBe('/api/content')
      expect(item.hxMethod).toBe('POST')
      expect(item.hxHeaders).toEqual({ 'X-Custom': 'value' })
      expect(item.hxTrigger).toBe('click')
    })

    it('should default hxMethod to GET', () => {
      const item = boxStore.open({ hxUrl: '/api/content' })
      expect(item.hxMethod).toBe('GET')
    })

    it('should default hxTrigger to load', () => {
      const item = boxStore.open({ hxUrl: '/api/content' })
      expect(item.hxTrigger).toBe('load')
    })
  })

  describe('Position Classes', () => {
    it('should get position classes', () => {
      const classes = boxStore.getPositionClasses('top-left')
      expect(classes).toBeTruthy()
      expect(typeof classes).toBe('string')
    })
  })

  describe('Width and Height Configuration', () => {
    it('should use numeric width if provided', () => {
      const item = boxStore.open({ width: 600 })
      expect(item.width).toBe(600)
    })

    it('should use default width if no width specified', () => {
      const item = boxStore.open({})
      expect(item.width).toBe(400)
    })

    it('should set width to undefined if widthClass is provided', () => {
      const item = boxStore.open({ widthClass: 'w-96' })
      expect(item.width).toBeUndefined()
      expect(item.widthClass).toBe('w-96')
    })

    it('should support height and heightClass', () => {
      const item = boxStore.open({ height: 500, heightClass: 'h-screen' })
      expect(item.height).toBe(500)
      expect(item.heightClass).toBe('h-screen')
    })
  })

  describe('Modal Behavior', () => {
    it('should be modal by default', () => {
      const item = boxStore.open('dialog')
      expect(item.modal).toBe(true)
    })

    it('should respect modal option', () => {
      const item = boxStore.open({ modal: false })
      expect(item.modal).toBe(false)
    })

    it('should disable closeOnEsc by default', () => {
      const item = boxStore.open('dialog')
      expect(item.closeOnEsc).toBe(false)
    })

    it('should disable closeOnOverlayClick by default', () => {
      const item = boxStore.open('dialog')
      expect(item.closeOnOverlayClick).toBe(false)
    })
  })

  // ============================================================
  // extraButtons - Custom Button Extension
  // ============================================================

  describe('extraButtons Configuration', () => {
    it('should default extraButtons to empty array', () => {
      const item = boxStore.open('dialog')
      expect(item.extraButtons).toEqual([])
    })

    it('should accept extraButtons array', () => {
      const btns = [
        { text: 'Save Draft', position: 'before', onClick: vi.fn() },
        { text: 'Preview', position: 'after', onClick: vi.fn() }
      ]
      const item = boxStore.open({ extraButtons: btns })
      expect(item.extraButtons).toHaveLength(2)
      expect(item.extraButtons[0].text).toBe('Save Draft')
      expect(item.extraButtons[1].text).toBe('Preview')
    })

    it('should normalize non-array extraButtons to empty array', () => {
      const item = boxStore.open({ extraButtons: 'invalid' })
      expect(item.extraButtons).toEqual([])
    })

    it('should normalize null extraButtons to empty array', () => {
      const item = boxStore.open({ extraButtons: null })
      expect(item.extraButtons).toEqual([])
    })

    it('should allow runtime mutation of extraButtons', () => {
      const item = boxStore.open('dialog')
      expect(item.extraButtons).toHaveLength(0)

      item.extraButtons = [{ text: 'New Btn', onClick: vi.fn() }]
      expect(item.extraButtons).toHaveLength(1)
    })
  })

  describe('handleExtraButton()', () => {
    it('should do nothing if btn is null', async () => {
      const item = boxStore.open('dialog')
      vi.runAllTimers()
      const closeSpy = vi.spyOn(item, 'close')
      await boxStore.handleExtraButton(item, null)
      expect(closeSpy).not.toHaveBeenCalled()
    })

    it('should do nothing if btn has no onClick', async () => {
      const item = boxStore.open('dialog')
      vi.runAllTimers()
      const closeSpy = vi.spyOn(item, 'close')
      await boxStore.handleExtraButton(item, { text: 'No Handler' })
      expect(closeSpy).not.toHaveBeenCalled()
    })

    it('should close dialog when sync onClick returns void', async () => {
      const item = boxStore.open('dialog')
      vi.runAllTimers()
      const closeSpy = vi.spyOn(item, 'close')
      await boxStore.handleExtraButton(item, { text: 'Btn', onClick: () => {} })
      expect(closeSpy).toHaveBeenCalledOnce()
    })

    it('should NOT close dialog when sync onClick returns false', async () => {
      const item = boxStore.open('dialog')
      vi.runAllTimers()
      const closeSpy = vi.spyOn(item, 'close')
      await boxStore.handleExtraButton(item, { text: 'Btn', onClick: () => false })
      expect(closeSpy).not.toHaveBeenCalled()
    })

    it('should close dialog when async onClick resolves undefined', async () => {
      const item = boxStore.open('dialog')
      vi.runAllTimers()
      const closeSpy = vi.spyOn(item, 'close')
      await boxStore.handleExtraButton(item, {
        text: 'Btn',
        onClick: async () => undefined
      })
      expect(closeSpy).toHaveBeenCalledOnce()
    })

    it('should NOT close dialog when async onClick resolves false', async () => {
      const item = boxStore.open('dialog')
      vi.runAllTimers()
      const closeSpy = vi.spyOn(item, 'close')
      await boxStore.handleExtraButton(item, {
        text: 'Btn',
        onClick: async () => false
      })
      expect(closeSpy).not.toHaveBeenCalled()
    })

    it('should invoke onClick with the box item as argument', async () => {
      const item = boxStore.open('dialog')
      vi.runAllTimers()
      const onClick = vi.fn().mockReturnValue(false)
      await boxStore.handleExtraButton(item, { text: 'Btn', onClick })
      expect(onClick).toHaveBeenCalledWith(item)
    })
  })
})
