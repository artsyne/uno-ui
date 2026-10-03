import { describe, it, expect, beforeEach, vi } from 'vitest'
import { TextEncoder, TextDecoder } from 'util'

global.TextEncoder = TextEncoder
global.TextDecoder = TextDecoder

// B6（docs/a11y-design.md）：showCounter 字符计数器必须通过 aria-describedby
// 关联到输入框，屏幕阅读器聚焦输入时才能播报剩余字数上下文（WCAG 1.3.1/4.1.2）。
describe('Valid counter a11y (B6)', () => {
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

  it('counter exposes full semantic label via i18n', () => {
    // 注册 i18n store（真实环境由 core/i18n 提供 t(key, vars) 插值）
    Alpine.store('i18n', {
      t: (key, vars) =>
        key === 'valid.counter' ? `已输入 ${vars.current} 字，共 ${vars.max} 字` : key
    })
    const input = mountCounter('bio')
    const counter = document.getElementById('bio-counter')

    expect(counter.textContent).toBe('0/20') // 视觉保持紧凑
    expect(counter.getAttribute('aria-label')).toBe('已输入 0 字，共 20 字')

    input.value = 'hello'
    input.dispatchEvent(new global.Event('input'))
    expect(counter.getAttribute('aria-label')).toBe('已输入 5 字，共 20 字')
  })

  it('counter falls back to "current/max" label without i18n store', () => {
    mountCounter('bio')
    const counter = document.getElementById('bio-counter')
    expect(counter.getAttribute('aria-label')).toBe('0/20')
  })

  it('over-limit input gets aria-invalid, cleared when back within limit', () => {
    const input = mountCounter('bio')

    input.value = 'x'.repeat(21)
    input.dispatchEvent(new global.Event('input'))
    expect(input.getAttribute('aria-invalid')).toBe('true')

    input.value = 'ok'
    input.dispatchEvent(new global.Event('input'))
    expect(input.hasAttribute('aria-invalid')).toBe(false)
  })

  // 对比度专项（docs/a11y-design.md Decision 4）：计数器色阶须达 WCAG AA——
  // 默认态 text-gray-600、超限态 text-red-600（与错误文本 helperText 同源）；
  // 暗色配对 gray-400 red-400。禁止回退 gray-400/red-500。
  // 断言用 classList.contains 精确匹配 token：避开 'dark:text-gray-400' 子串误命中裸类。
  it('counter color tokens meet AA contrast (gray-600 default / red-600 over-limit)', () => {
    const input = mountCounter('bio')
    const counter = document.getElementById('bio-counter')

    // 默认态：亮色 gray-600，暗色 gray-400
    expect(counter.classList.contains('text-gray-600')).toBe(true)
    expect(counter.classList.contains('dark:text-gray-400')).toBe(true)
    expect(counter.classList.contains('text-gray-400')).toBe(false)

    // 超限态：亮色 red-600，暗色 red-400
    input.value = 'x'.repeat(21)
    input.dispatchEvent(new global.Event('input'))
    expect(counter.classList.contains('text-red-600')).toBe(true)
    expect(counter.classList.contains('dark:text-red-400')).toBe(true)
    expect(counter.classList.contains('text-red-500')).toBe(false)

    // 回到限内：恢复默认态（red 类清除、gray 类复原）
    input.value = 'ok'
    input.dispatchEvent(new global.Event('input'))
    expect(counter.classList.contains('text-gray-600')).toBe(true)
    expect(counter.classList.contains('dark:text-gray-400')).toBe(true)
    expect(counter.classList.contains('text-red-600')).toBe(false)
    expect(counter.classList.contains('dark:text-red-400')).toBe(false)
  })
})
