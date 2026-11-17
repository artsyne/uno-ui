import { describe, it, expect, vi, beforeEach } from 'vitest'
import { createBoxStore } from '../../src/components/box/box.store'

vi.mock('../../src/core/uno', () => ({
  Uno: {
    getConfig: vi.fn(() => ({})),
    getStyle: vi.fn((comp, key) => `${comp}-${key}`),
    dom: { ensureContainer: vi.fn(() => document.createElement('div')) },
    options: { ui: { layers: { box: 5000 } } }
  }
}))

const Alpine = {
  store: vi.fn(),
  nextTick: vi.fn((fn) => setTimeout(fn, 0)),
  initTree: vi.fn()
}

// 模态顶层以原生 inert 隔离页面背景（弹窗管理器之外的 body 子节点），任何引擎的
// Tab/程序聚焦都无法逃逸；关闭后仅恢复本组件 inert 过的节点，应用原本 inert 的节点
// 保持原状。覆盖空态、非模态、叠层、closeAll 与 opener/背景销毁。
describe('Box background inert isolation', () => {
  let boxStore
  let mockStore

  beforeEach(() => {
    vi.clearAllMocks()
    mockStore = {}
    Alpine.store.mockImplementation((name, value) => {
      if (value) mockStore[name] = value
      return mockStore[name]
    })
    Alpine.store('uno', { box: createBoxStore(Alpine) })
    boxStore = Alpine.store('uno').box
    document.body.innerHTML = `
      <header id="hdr"><button id="opener">open</button></header>
      <main id="app"><input id="bg-input"></main>
      <div id="uno-box-manager"></div>
    `
  })

  const manager = () => document.getElementById('uno-box-manager')
  const addDialog = (id, { modal = true, isOpen = true } = {}) => {
    const el = document.createElement('div')
    el.id = 'uno-box-outer-container' + id
    el.innerHTML = '<button class="dlg-btn">ok</button>'
    manager().appendChild(el)
    boxStore.items.push({ id, isOpen, modal })
    return el
  }

  it('模态顶层 → 背景 body 子节点 inert，管理器除外', () => {
    addDialog('A')
    boxStore.syncUnderlayInert()
    expect(document.getElementById('hdr').hasAttribute('inert')).toBe(true)
    expect(document.getElementById('app').hasAttribute('inert')).toBe(true)
    expect(manager().hasAttribute('inert')).toBe(false)
  })

  it('关闭后恢复背景，且保留应用原本 inert 的节点', () => {
    const pre = document.createElement('aside')
    pre.id = 'pre'
    pre.setAttribute('inert', '')
    document.body.appendChild(pre)

    addDialog('A')
    boxStore.syncUnderlayInert()
    expect(document.getElementById('hdr').hasAttribute('inert')).toBe(true)

    boxStore.items.length = 0
    boxStore.syncUnderlayInert()
    expect(document.getElementById('hdr').hasAttribute('inert')).toBe(false)
    expect(document.getElementById('app').hasAttribute('inert')).toBe(false)
    // 应用自设 inert 未被本组件跟踪，恢复时保持原状
    expect(pre.hasAttribute('inert')).toBe(true)
  })

  it('非模态顶层 → 背景不隔离', () => {
    addDialog('A', { modal: false })
    boxStore.syncUnderlayInert()
    expect(document.getElementById('hdr').hasAttribute('inert')).toBe(false)
    expect(document.getElementById('app').hasAttribute('inert')).toBe(false)
  })

  it('叠层 → 下层弹窗 inert+aria-hidden，背景持续隔离直至最后一层关闭', () => {
    const a = addDialog('A')
    addDialog('B')
    boxStore.syncUnderlayInert()
    expect(a.hasAttribute('inert')).toBe(true)
    expect(a.getAttribute('aria-hidden')).toBe('true')
    expect(document.getElementById('uno-box-outer-containerB').hasAttribute('inert')).toBe(false)
    expect(document.getElementById('app').hasAttribute('inert')).toBe(true)

    // 关闭顶层 B：A 成为顶层，解除 inert，背景仍隔离
    boxStore.items.find((i) => i.id === 'B').isOpen = false
    boxStore.syncUnderlayInert()
    expect(a.hasAttribute('inert')).toBe(false)
    expect(a.hasAttribute('aria-hidden')).toBe(false)
    expect(document.getElementById('app').hasAttribute('inert')).toBe(true)

    // 关闭 A：背景恢复
    boxStore.items.find((i) => i.id === 'A').isOpen = false
    boxStore.syncUnderlayInert()
    expect(document.getElementById('app').hasAttribute('inert')).toBe(false)
  })

  it('closeAll → 恢复背景 inert 原状态', () => {
    addDialog('A')
    boxStore.syncUnderlayInert()
    expect(document.getElementById('app').hasAttribute('inert')).toBe(true)
    boxStore.closeAll()
    expect(boxStore.items).toEqual([])
    expect(document.getElementById('app').hasAttribute('inert')).toBe(false)
  })

  it('空态同步无副作用；背景节点销毁后恢复不抛错', () => {
    // 空态：未跟踪任何节点，同步为无操作
    boxStore.syncUnderlayInert()
    expect(document.getElementById('hdr').hasAttribute('inert')).toBe(false)

    addDialog('A')
    boxStore.syncUnderlayInert()
    // 模拟 opener 所在背景节点被外部销毁（htmx 替换/导航）
    document.getElementById('hdr').remove()
    boxStore.items.length = 0
    expect(() => boxStore.syncUnderlayInert()).not.toThrow()
    expect(document.getElementById('app').hasAttribute('inert')).toBe(false)
  })
})