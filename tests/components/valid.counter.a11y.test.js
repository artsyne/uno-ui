import { describe, it, expect, beforeEach, vi } from 'vitest'
import { TextEncoder, TextDecoder } from 'util'

global.TextEncoder = TextEncoder
global.TextDecoder = TextDecoder

// P1-B6（docs/a11y-design.md）：showCounter 字符计数器必须通过 aria-describedby
// 关联到输入框，屏幕阅读器聚焦输入时才能播报剩余字数上下文（WCAG 1.3.1/4.1.2）。
describe('Valid counter a11y (P1-B6)', () => {
  let ValidPlugin
  let UnoInstance
  let Alpine
  let directives = {}

  beforeEach(async () => {
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

    const validModule = await import('../../src/components/valid')
    ValidPlugin = validModule.default
    const unoModule = await import('../../src/core/uno')
    UnoInstance = unoModule.UnoInstance

    dom.window.HTMLElement.prototype.scrollIntoView = vi.fn()
    dom.window.HTMLElement.prototype.focus = vi.fn()
    dom.window.HTMLElement.prototype.matches = function (selector) {
      return (
        (this.constructor.name === 'HTMLFormElement' && selector === 'form') ||
        (this.constructor.name === 'HTMLInputElement' && selector === 'input') ||
        dom.window.document.querySelector(selector) === this
      )
    }
    dom.window.Element.prototype.scrollIntoView = vi.fn()
    dom.window.Element.prototype.focus = vi.fn()

    const store = {}
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

    global.Uno = new UnoInstance()
    ValidPlugin(Alpine)
  })

  const mountCounter = (name, describedBy) => {
    document.body.innerHTML = `
      <form>
        <div class="relative">
          <input name="${name}" ${describedBy ? `aria-describedby="${describedBy}"` : ''} />
        </div>
      </form>
    `
    const input = document.querySelector('input')
    global.Uno.options.valid = {}

    const directiveFn = directives['uno-valid']
    const evaluateLater = (expr) => (cb) => cb([{ type: 'length', max: 20, showCounter: true }])
    const cleanup = () => {}
    directiveFn(input, { expression: '...' }, { evaluateLater, cleanup })
    return input
  }

  it('links counter to input via aria-describedby', () => {
    const input = mountCounter('bio')
    const counter = document.getElementById('bio-counter')

    expect(counter).not.toBeNull()
    const tokens = (input.getAttribute('aria-describedby') || '').split(/\s+/)
    expect(tokens).toContain('bio-counter')
  })

  it('preserves existing aria-describedby tokens (e.g. error element)', () => {
    const input = mountCounter('bio', 'bio-error')
    const tokens = (input.getAttribute('aria-describedby') || '').split(/\s+/)

    expect(tokens).toContain('bio-error')
    expect(tokens).toContain('bio-counter')
  })
})
