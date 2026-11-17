import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { createTabStore } from '../../src/components/tab/tab.store'
import { registerTabDirectives } from '../../src/components/tab/tab.directive'
import { Uno } from '../../src/core/uno'

vi.mock('../../src/core/uno', () => ({
  Uno: {
    getConfig: vi.fn(),
    options: {}
  }
}))

vi.mock('../../src/utils/log', () => ({
  log: {
    info: vi.fn(),
    warn: vi.fn(),
    error: vi.fn()
  }
}))

// B4（docs/a11y-design.md）：APG tabs 指令键盘契约——Space 激活焦点 tab（锚点无原生
// Space 激活，须合成 click 并阻止滚动；按钮原生 Space 被接管不双重激活；长按不连发），
// ←→/Home/End 旋转焦点跟随并复用 click 路径激活。
describe('Tab directive keyboard (B4, APG tabs)', () => {
  let store
  let handler
  let Alpine
  let effects

  const bind = (el, config) => {
    const value = typeof config === 'string' ? config : config
    handler(
      el,
      { expression: 'cfg' },
      {
        evaluateLater: () => (cb) => cb(value),
        effect: (fn) => {
          fn()
          effects.push(fn)
        },
        cleanup: vi.fn()
      }
    )
  }

  // 模拟 Alpine effect 的响应式重跑：激活态变化后同步 aria-selected/roving tabindex
  const sync = () => effects.forEach((fn) => fn())

  const press = (el, key, opts = {}) => {
    const event = new KeyboardEvent('keydown', { key, bubbles: true, cancelable: true, ...opts })
    el.dispatchEvent(event)
    return event
  }

  beforeEach(() => {
    vi.clearAllMocks()
    document.body.innerHTML = ''
    effects = []

    Uno.getConfig.mockImplementation((key, opts) => ({
      group: 'default',
      hash: true,
      ...opts
    }))

    Alpine = {
      store: vi.fn(),
      directive: vi.fn((name, fn) => {
        handler = fn
      })
    }
    store = createTabStore(Alpine)
    Alpine.store = vi.fn(() => ({ tab: store }))
    registerTabDirectives(Alpine)
  })

  afterEach(() => {
    document.body.innerHTML = ''
  })

  it('Space on anchor tab prevents scroll and activates via click path', () => {
    const a = document.createElement('a')
    a.href = '#'
    document.body.appendChild(a)
    const b = document.createElement('a')
    b.href = '#'
    document.body.appendChild(b)
    bind(a, { name: 'a', group: 'g' })
    bind(b, { name: 'b', group: 'g' })
    store.activate({ name: 'a', group: 'g' })

    const onClick = vi.fn()
    b.addEventListener('click', onClick)
    const event = press(b, ' ')

    expect(event.defaultPrevented).toBe(true)
    expect(onClick).toHaveBeenCalledTimes(1)
    expect(store.isActive('b', 'g')).toBe(true)
    sync()
    expect(b.getAttribute('aria-selected')).toBe('true')
    expect(b.getAttribute('tabindex')).toBe('0')
    expect(a.getAttribute('tabindex')).toBe('-1')
  })

  it('Space on button tab activates exactly once (native keyup click preempted)', () => {
    const btn = document.createElement('button')
    document.body.appendChild(btn)
    bind(btn, { name: 'a', group: 'g' })

    const onClick = vi.fn()
    btn.addEventListener('click', onClick)
    const event = press(btn, ' ')

    expect(event.defaultPrevented).toBe(true)
    expect(onClick).toHaveBeenCalledTimes(1)
    expect(store.isActive('a', 'g')).toBe(true)
  })

  it('held Space (repeat) does not re-activate', () => {
    const a = document.createElement('a')
    a.href = '#'
    document.body.appendChild(a)
    bind(a, { name: 'a', group: 'g' })

    const onClick = vi.fn()
    a.addEventListener('click', onClick)
    press(a, ' ')
    press(a, ' ', { repeat: true })
    press(a, ' ', { repeat: true })

    expect(onClick).toHaveBeenCalledTimes(1)
  })

  it('Spacebar legacy key name also activates', () => {
    const a = document.createElement('a')
    a.href = '#'
    document.body.appendChild(a)
    bind(a, { name: 'a', group: 'g' })

    const event = press(a, 'Spacebar')
    expect(event.defaultPrevented).toBe(true)
    expect(store.isActive('a', 'g')).toBe(true)
  })

  it('arrow keys rotate focus and activation across the group', () => {
    const els = ['a', 'b', 'c'].map((name) => {
      const el = document.createElement('a')
      el.href = '#'
      el.id = 'tab-' + name
      document.body.appendChild(el)
      bind(el, { name, group: 'g' })
      return el
    })
    store.activate({ name: 'a', group: 'g' })
    els[0].focus()

    press(els[0], 'ArrowRight')
    expect(document.activeElement).toBe(els[1])
    expect(store.isActive('b', 'g')).toBe(true)

    press(els[1], 'End')
    expect(document.activeElement).toBe(els[2])
    expect(store.isActive('c', 'g')).toBe(true)

    press(els[2], 'ArrowRight')
    expect(document.activeElement).toBe(els[0])
    expect(store.isActive('a', 'g')).toBe(true)

    press(els[0], 'ArrowLeft')
    expect(document.activeElement).toBe(els[2])
    expect(store.isActive('c', 'g')).toBe(true)
  })
})