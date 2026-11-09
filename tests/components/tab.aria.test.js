import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { createTabStore } from '../../src/components/tab/tab.store'
import { Uno } from '../../src/core/uno'

// Mock dependencies
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

const Alpine = {
  store: vi.fn(),
  nextTick: vi.fn((fn) => setTimeout(fn, 0))
}

// B4（docs/a11y-design.md）：APG tabs 模式的 store 契约——组内注册序（DOM 序代理）、
// 方向键/Home/End 旋转目标计算、共享面板（htmx swap 模式）aria-labelledby 随激活 tab 更新。
describe('Tab ARIA store contract (B4)', () => {
  let store

  beforeEach(() => {
    vi.clearAllMocks()
    vi.useFakeTimers()
    document.body.innerHTML = ''

    Uno.getConfig.mockImplementation((key, opts) => ({
      group: 'default',
      hash: true,
      ...opts
    }))

    store = createTabStore(Alpine)
    Alpine.store = vi.fn(() => ({ tab: store }))
  })

  afterEach(() => {
    vi.useRealTimers()
  })

  it('register preserves group order for arrow-key rotation', () => {
    store.register('a', 'g')
    store.register('b', 'g')
    store.register('c', 'g')
    store.register('b', 'g') // dedupe

    expect(store.order.g).toEqual(['a', 'b', 'c'])
  })

  it('rotate computes APG keyboard targets', () => {
    store.register('a', 'g')
    store.register('b', 'g')
    store.register('c', 'g')

    expect(store.rotate('g', 'a', 'ArrowRight')).toBe('b')
    expect(store.rotate('g', 'c', 'ArrowRight')).toBe('a') // wrap
    expect(store.rotate('g', 'a', 'ArrowLeft')).toBe('c') // wrap
    expect(store.rotate('g', 'c', 'ArrowLeft')).toBe('b')
    expect(store.rotate('g', 'b', 'Home')).toBe('a')
    expect(store.rotate('g', 'b', 'End')).toBe('c')
    expect(store.rotate('g', 'x', 'ArrowRight')).toBe('a') // unknown current
    expect(store.rotate('empty', 'a', 'ArrowRight')).toBeNull()
  })

  it('activate relabels shared tabpanel (htmx swap mode)', () => {
    const panel = document.createElement('div')
    panel.setAttribute('role', 'tabpanel')
    panel.setAttribute('data-tab-group', 'g')
    panel.setAttribute('aria-labelledby', 'tab-a')
    document.body.appendChild(panel)

    store.register('a', 'g')
    store.register('b', 'g')
    store.activate({ name: 'b', group: 'g', hash: false })

    expect(panel.getAttribute('aria-labelledby')).toBe('tab-b')
  })

  it('activate leaves per-tab panels (static labelledby) untouched', () => {
    const panelA = document.createElement('div')
    panelA.id = 'panel-a'
    panelA.setAttribute('role', 'tabpanel')
    panelA.setAttribute('aria-labelledby', 'tab-a')
    document.body.appendChild(panelA)

    store.register('a', 'g')
    store.register('b', 'g')
    store.activate({ name: 'b', group: 'g', hash: false })

    expect(panelA.getAttribute('aria-labelledby')).toBe('tab-a')
  })
})
