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
})
