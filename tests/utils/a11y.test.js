import { describe, it, expect, beforeEach, vi, afterEach } from 'vitest'
import { trapFocus } from '../../src/utils/a11y'

describe('A11y Utils - trapFocus', () => {
  let container
  let btn1, btn2, btn3

  beforeEach(() => {
    document.body.innerHTML = `
      <div id="container">
        <button id="btn1">Button 1</button>
        <button id="btn2">Button 2</button>
        <div style="display:none">
            <button id="btn3">Hidden</button>
        </div>
      </div>
    `
    container = document.getElementById('container')
    btn1 = document.getElementById('btn1')
    btn2 = document.getElementById('btn2')

    // Mock getClientRects to simulate visibility
    // In happy-dom, default might be empty.
    HTMLElement.prototype.getClientRects = function () {
      // Check ancestry for display:none
      let el = this
      while (el && el instanceof HTMLElement) {
        const style = window.getComputedStyle(el)
        if (style.display === 'none' || style.visibility === 'hidden') {
          return []
        }
        el = el.parentElement
      }
      return [{ width: 10, height: 10 }]
    }

    vi.useFakeTimers()
    // Mock requestAnimationFrame
    global.requestAnimationFrame = (cb) => setTimeout(cb, 0)
  })

  afterEach(() => {
    vi.useRealTimers()
  })

  it('should focus first element initially', () => {
    const cleanup = trapFocus(container)

    vi.runAllTimers()

    expect(document.activeElement).toBe(btn1)
    cleanup()
  })

  it('should not take initial focus when initialFocus=false', () => {
    const cleanup = trapFocus(container, { initialFocus: false })

    vi.runAllTimers()

    // Trap only cycles Tab; initial focus is owned by the caller (box store)
    expect(document.activeElement).toBe(document.body)
    cleanup()
  })

  it('should cycle focus to first element when tabbing from last', () => {
    const cleanup = trapFocus(container)
    vi.runAllTimers()

    // Move focus to last
    btn2.focus()
    expect(document.activeElement).toBe(btn2)

    // Simulate Tab key on last element
    const event = new KeyboardEvent('keydown', {
      key: 'Tab',
      bubbles: true,
      cancelable: true
    })
    container.dispatchEvent(event)

    expect(document.activeElement).toBe(btn1)
    cleanup()
  })

  it('should cycle focus to last element when shift-tabbing from first', () => {
    const cleanup = trapFocus(container)
    vi.runAllTimers()

    // Focus first
    btn1.focus()

    // Simulate Shift+Tab
    const event = new KeyboardEvent('keydown', {
      key: 'Tab',
      shiftKey: true,
      bubbles: true,
      cancelable: true
    })
    container.dispatchEvent(event)

    expect(document.activeElement).toBe(btn2)
    cleanup()
  })

  it('should wrap Tab when the neutral target has no focusable after it in DOM order', () => {
    // Mirrors box.template.js layout: close button first, then the title
    // (tabindex="-1") neutral target — the browser would escape the container
    // from the title, so the trap must wrap to the first focusable instead.
    document.body.innerHTML = `
      <div id="dlg">
        <button id="close">x</button>
        <h3 id="title" tabindex="-1">Title</h3>
      </div>
    `
    const dlg = document.getElementById('dlg')
    const close = document.getElementById('close')
    const title = document.getElementById('title')
    const cleanup = trapFocus(dlg, { initialFocus: false })
    vi.runAllTimers()

    title.focus()
    expect(document.activeElement).toBe(title)

    const event = new KeyboardEvent('keydown', { key: 'Tab', bubbles: true, cancelable: true })
    title.dispatchEvent(event)

    expect(event.defaultPrevented).toBe(true)
    expect(document.activeElement).toBe(close)
    cleanup()
  })

  it('should not intercept Tab when a focusable exists after the neutral target', () => {
    document.body.innerHTML = `
      <div id="dlg">
        <button id="close">x</button>
        <h3 id="title" tabindex="-1">Title</h3>
        <button id="ok">OK</button>
      </div>
    `
    const dlg = document.getElementById('dlg')
    const title = document.getElementById('title')
    const cleanup = trapFocus(dlg, { initialFocus: false })
    vi.runAllTimers()

    title.focus()
    const event = new KeyboardEvent('keydown', { key: 'Tab', bubbles: true, cancelable: true })
    title.dispatchEvent(event)

    // Browser proceeds naturally (title -> ok in real engines); trap stays out
    expect(event.defaultPrevented).toBe(false)
    cleanup()
  })

  it('should wrap shift-tab to last focusable when tabbing back from the container target', () => {
    container.setAttribute('tabindex', '-1')
    const cleanup = trapFocus(container)
    vi.runAllTimers()

    container.focus()
    expect(document.activeElement).toBe(container)

    const event = new KeyboardEvent('keydown', {
      key: 'Tab',
      shiftKey: true,
      bubbles: true,
      cancelable: true
    })
    container.dispatchEvent(event)

    expect(event.defaultPrevented).toBe(true)
    expect(document.activeElement).toBe(btn2)
    container.removeAttribute('tabindex')
    cleanup()
  })

  it('should not intercept Tab from the container target when focusables follow', () => {
    container.setAttribute('tabindex', '-1')
    const cleanup = trapFocus(container)
    vi.runAllTimers()

    container.focus()
    const event = new KeyboardEvent('keydown', { key: 'Tab', bubbles: true, cancelable: true })
    container.dispatchEvent(event)

    expect(event.defaultPrevented).toBe(false)
    container.removeAttribute('tabindex')
    cleanup()
  })
})
