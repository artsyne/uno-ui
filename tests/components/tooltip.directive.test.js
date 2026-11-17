import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { registerTooltipDirectives } from '../../src/components/tooltip/tooltip.directive'
import { Uno } from '../../src/core/uno'

vi.mock('../../src/core/uno', () => ({
  Uno: {
    getConfig: vi.fn()
  }
}))

vi.mock('../../src/utils/log', () => ({
  log: { info: vi.fn(), warn: vi.fn(), error: vi.fn() }
}))

// 默认 hover 提示也响应键盘聚焦（APG tooltip 对 hover/focus
// 两种输入均显示），失焦隐藏；click 为披露模式不随聚焦自动打开。
describe('Tooltip directive focus response', () => {
  let handler
  let store
  let el

  const bind = (expression = '') => {
    handler(
      el,
      { expression, modifiers: [] },
      {
        evaluateLater: () => (cb) => cb(null),
        effect: (fn) => fn(),
        cleanup: vi.fn()
      }
    )
  }

  beforeEach(() => {
    vi.clearAllMocks()
    document.body.innerHTML = ''
    el = document.createElement('button')
    document.body.appendChild(el)

    store = {
      visible: false,
      triggerEl: null,
      triggerMode: 'hover',
      show: vi.fn(),
      hide: vi.fn(),
      scheduleHide: vi.fn()
    }
    const Alpine = {
      store: vi.fn(() => ({ tooltip: store })),
      directive: vi.fn((name, fn) => {
        handler = fn
      })
    }
    registerTooltipDirectives(Alpine)
  })

  afterEach(() => {
    document.body.innerHTML = ''
  })

  it('shows on focus in default hover mode', () => {
    Uno.getConfig.mockReturnValue({ trigger: 'hover' })
    bind()
    el.dispatchEvent(new FocusEvent('focus'))
    expect(store.show).toHaveBeenCalledTimes(1)
  })

  it('shows on focus in focus mode', () => {
    Uno.getConfig.mockReturnValue({ trigger: 'focus' })
    bind()
    el.dispatchEvent(new FocusEvent('focus'))
    expect(store.show).toHaveBeenCalledTimes(1)
  })

  it('does not auto-open on focus in click (disclosure) mode', () => {
    Uno.getConfig.mockReturnValue({ trigger: 'click' })
    bind()
    el.dispatchEvent(new FocusEvent('focus'))
    expect(store.show).not.toHaveBeenCalled()
  })

  it('hides on blur in focus mode', () => {
    Uno.getConfig.mockReturnValue({ trigger: 'focus' })
    bind()
    store.triggerEl = el
    store.triggerMode = 'focus'
    el.dispatchEvent(new FocusEvent('blur'))
    expect(store.hide).toHaveBeenCalledTimes(1)
  })

  it('schedules hide on blur in hover mode when the pointer is not over the trigger', () => {
    Uno.getConfig.mockReturnValue({ trigger: 'hover' })
    bind()
    store.triggerEl = el
    store.triggerMode = 'hover'
    el.dispatchEvent(new FocusEvent('blur'))
    expect(store.scheduleHide).toHaveBeenCalledTimes(1)
    expect(store.hide).not.toHaveBeenCalled()
  })
})