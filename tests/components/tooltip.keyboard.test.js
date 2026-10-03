import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { createTooltipStore } from '../../src/components/tooltip/tooltip.store'
import { Uno } from '../../src/core/uno'
import { DomManager } from '../../src/utils/dom'

// Mock dependencies
vi.mock('../../src/core/uno', () => ({
  Uno: {
    getConfig: vi.fn(),
    options: {
      ui: {
        layers: { tooltip: 6000 }
      }
    }
  }
}))

vi.mock('../../src/utils/dom', () => ({
  DomManager: {
    ensureContainer: vi.fn()
  },
  sanitizeHTML: vi.fn((html) => html)
}))

vi.mock('../../src/utils/log', () => ({
  log: {
    info: vi.fn(),
    warn: vi.fn(),
    error: vi.fn()
  }
}))

vi.mock('../../src/utils/position', () => ({
  calculatePosition: vi.fn(() => ({
    top: 100,
    left: 200,
    position: 'bottom'
  }))
}))

// Mock Alpine（initTree 由 ensureContainer mock 的 _x_dataStack 标记跳过）
const Alpine = {
  store: vi.fn(),
  nextTick: vi.fn((fn) => setTimeout(fn, 0)),
  initTree: vi.fn()
}

// A3（docs/a11y-design.md）：click 触发 = 弹层/披露模式键盘契约（WCAG 2.1.1 / 2.4.3 / 4.1.2）。
// 打开时焦点移入弹层首个可聚焦元素、Tab 循环于弹层内、关闭时焦点还原触发器；
// 触发器暴露 aria-haspopup/aria-expanded 并与弹层 aria-describedby 关联；
// 弹层含 role="menu" 时支持方向键导航；锚点触发器补 Space 键合成（Enter 由原生 click 覆盖）。
// hover/focus 触发的提示不得抢占或还原焦点。
describe('Tooltip keyboard support (A3, click trigger popup)', () => {
  let store
  let trigger

  beforeEach(() => {
    vi.clearAllMocks()
    vi.useFakeTimers()
    document.body.innerHTML = ''

    Uno.getConfig.mockImplementation((key, opts) => ({
      placement: 'bottom',
      trigger: 'hover',
      offset: 8,
      arrow: true,
      clickOutside: true,
      hxTrigger: 'load',
      ...opts
    }))

    // 真实容器入 DOM（焦点断言前提）；_x_dataStack 标记跳过 store 的 Alpine.initTree 调用
    vi.mocked(DomManager.ensureContainer).mockImplementation((id, html) => {
      const container = document.createElement('div')
      container.id = id
      container.innerHTML = html
      container._x_dataStack = []
      document.body.appendChild(container)
      return container
    })

    store = createTooltipStore(Alpine)
    trigger = document.createElement('a')
    trigger.href = '#'
    document.body.appendChild(trigger)
  })

  afterEach(() => {
    vi.useRealTimers()
  })

  // 模拟 Alpine x-html 将 sanitize 后的内容注入弹层 body
  const injectPopupContent = (html) => {
    document.getElementById('uno-tooltip-body').innerHTML = html
  }

  const openClickPopup = (content) => {
    store.show({ content, trigger: 'click', element: trigger })
    injectPopupContent(content)
    vi.advanceTimersByTime(10)
  }

  it('opens click popup: focus moves into popup and trigger exposes expanded state', () => {
    openClickPopup('<ul role="menu"><li><a href="#" id="menu-edit">Edit</a></li></ul>')

    expect(document.activeElement).toBe(document.getElementById('menu-edit'))
    expect(trigger.getAttribute('aria-haspopup')).toBe('true')
    expect(trigger.getAttribute('aria-expanded')).toBe('true')
    expect(trigger.getAttribute('aria-describedby')).toBe('uno-tooltip-content')
  })

  it('hides click popup: focus restores to trigger and expanded state collapses', () => {
    openClickPopup('<ul role="menu"><li><a href="#" id="menu-edit">Edit</a></li></ul>')
    store.hide()

    expect(document.activeElement).toBe(trigger)
    expect(trigger.getAttribute('aria-expanded')).toBe('false')
    expect(trigger.getAttribute('aria-describedby')).toBeNull()
  })

  it('hover trigger never moves focus nor mutates aria state', () => {
    store.show({ content: '<p>hint</p>', trigger: 'hover', element: trigger })
    injectPopupContent('<p>hint</p>')
    vi.advanceTimersByTime(10)

    expect(document.activeElement).toBe(document.body)
    expect(trigger.getAttribute('aria-expanded')).toBeNull()

    store.hide()
    expect(document.activeElement).toBe(document.body)
    expect(trigger.getAttribute('aria-expanded')).toBeNull()
  })

  it('arrow keys navigate within role=menu popup', () => {
    openClickPopup(
      '<ul role="menu"><li><a href="#" id="menu-edit">Edit</a></li><li><a href="#" id="menu-del">Delete</a></li></ul>'
    )
    const edit = document.getElementById('menu-edit')
    const del = document.getElementById('menu-del')

    edit.dispatchEvent(new KeyboardEvent('keydown', { key: 'ArrowDown', bubbles: true }))
    expect(document.activeElement).toBe(del)

    del.dispatchEvent(new KeyboardEvent('keydown', { key: 'ArrowUp', bubbles: true }))
    expect(document.activeElement).toBe(edit)

    edit.dispatchEvent(new KeyboardEvent('keydown', { key: 'End', bubbles: true }))
    expect(document.activeElement).toBe(del)
  })

  it('anchor trigger synthesizes Space key to close popup (Enter covered by native click)', () => {
    openClickPopup('<ul role="menu"><li><a href="#" id="menu-edit">Edit</a></li></ul>')

    const event = new KeyboardEvent('keydown', { key: ' ', bubbles: true, cancelable: true })
    trigger.dispatchEvent(event)

    expect(event.defaultPrevented).toBe(true)
    expect(store.visible).toBe(false)
    expect(document.activeElement).toBe(trigger)
  })
})
