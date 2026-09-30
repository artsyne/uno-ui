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

// P1-B5（docs/a11y-design.md）：Box 关闭后焦点还原至 opener（2.4.3），多层弹窗下层
// inert + aria-hidden（屏幕阅读器不穿透底层）。不改现有 scheduleAutoFocus 契约。
describe('Box focus restore & multilayer inert (P1-B5)', () => {
  let boxStore
  let mockStore

  beforeEach(() => {
    vi.clearAllMocks()
    vi.useFakeTimers()
    document.body.innerHTML = ''

    mockStore = {}
    Alpine.store.mockImplementation((name, value) => {
      if (value) {
        mockStore[name] = value
        if (value.box?.init) value.box.init()
      }
      return mockStore[name]
    })

    Uno.getConfig.mockImplementation((key, opts) => ({
      defaultWidth: 400,
      defaultPosition: 'center',
      okText: 'OK',
      cancelText: 'Cancel',
      closeOnEsc: false,
      closeOnOverlayClick: false,
      ...opts
    }))

    Alpine.store('uno', { box: createBoxStore(Alpine) })
    boxStore = Alpine.store('uno').box
  })

  afterEach(() => {
    vi.useRealTimers()
  })

  // 模拟渲染层产出 box 容器（真实应用由 Alpine 模板渲染）
  const mountBoxContainer = (item, innerHtml = '') => {
    const el = document.createElement('div')
    el.id = 'uno-box-' + item.id
    el.innerHTML = innerHtml
    document.body.appendChild(el)
    return el
  }

  it('open records opener and close restores focus to it', () => {
    const opener = document.createElement('button')
    opener.textContent = 'open'
    document.body.appendChild(opener)
    opener.focus()

    const item = boxStore.open({ content: '<p>body</p>' })
    expect(item.opener).toBe(opener)

    const el = mountBoxContainer(item, '<button id="ok-in-box">OK</button>')
    vi.advanceTimersByTime(10)
    document.getElementById('ok-in-box').focus()

    item.close()
    expect(document.activeElement).toBe(opener)

    vi.advanceTimersByTime(400)
    expect(boxStore.items.length).toBe(0)
  })

  it('close leaves outside focus untouched', () => {
    const opener = document.createElement('button')
    const outside = document.createElement('button')
    document.body.appendChild(opener)
    document.body.appendChild(outside)
    opener.focus()

    const item = boxStore.open({ content: '<p>body</p>' })
    mountBoxContainer(item)
    vi.advanceTimersByTime(10)
    outside.focus()

    item.close()
    expect(document.activeElement).toBe(outside)
    vi.advanceTimersByTime(400)
  })

  it('multilayer: underlay box gets inert + aria-hidden, released when top closes', () => {
    const itemA = boxStore.open({ content: '<p>a</p>' })
    const elA = mountBoxContainer(itemA)
    vi.advanceTimersByTime(10)
    expect(elA.hasAttribute('inert')).toBe(false)

    const itemB = boxStore.open({ content: '<p>b</p>' })
    const elB = mountBoxContainer(itemB)
    vi.advanceTimersByTime(10)

    expect(elA.getAttribute('aria-hidden')).toBe('true')
    expect(elA.hasAttribute('inert')).toBe(true)
    expect(elB.hasAttribute('inert')).toBe(false)

    itemB.close()
    vi.advanceTimersByTime(400)

    expect(elA.getAttribute('aria-hidden')).toBeNull()
    expect(elA.hasAttribute('inert')).toBe(false)
  })
})
