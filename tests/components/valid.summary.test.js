import { describe, it, expect, beforeEach, vi } from 'vitest'
import { TextEncoder, TextDecoder } from 'util'
// JSDOM 在模块加载阶段引入，避免冷启动时占用 beforeEach 的 10s 钩子预算。
import { JSDOM } from 'jsdom'

global.TextEncoder = TextEncoder
global.TextDecoder = TextDecoder

// E4（docs/a11y-design.md / ui-style-guide.md §10.4，WCAG 3.3.1/3.3.3）：
// 多字段表单提交失败时，validateForm 在表单顶部渲染可聚焦错误摘要
// （role="alert" tabindex="-1" aria-labelledby），每条链接到对应字段，
// 焦点移至摘要；单错维持聚焦首错字段现状；内联错误保留（摘要不替代）。
// 摘要视觉类 .ui-error-summary（容器/标题/列表/链接）由 website tailwind.css 提供。
describe('Valid error summary (E4, WCAG 3.3.1/3.3.3)', () => {
  let ValidPlugin
  let UnoInstance
  let Uno
  let Alpine
  let store
  let dom

  beforeEach(async () => {
    dom = new JSDOM('<!DOCTYPE html><body></body>')
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
    // 插件读取的是模块单例（export const Uno = new UnoInstance()），非 global.Uno；
    // 每个用例前重置插件配置，防跨用例泄漏
    Uno = unoModule.Uno
    Uno.options.valid = {}

    // 焦点/滚动 mock 在元素上留痕：共享 vi.fn() 无法区分哪个元素被处理
    dom.window.HTMLElement.prototype.focus = vi.fn(function () {
      this.dataset.testFocus = '1'
    })
    dom.window.HTMLElement.prototype.scrollIntoView = vi.fn(function () {
      this.dataset.testScroll = '1'
    })

    store = {}
    Alpine = {
      store: (name, value) => {
        if (value) store[name] = value
        return store[name]
      },
      directive: vi.fn(),
      magic: vi.fn(),
      effect: vi.fn(),
      initTree: vi.fn()
    }

    global.Uno = new UnoInstance()
    ValidPlugin(Alpine)
  })

  /** user 必填 / email 格式 / pwd 必填（无 id，测锚点回退） */
  const setupForm = () => {
    document.body.innerHTML = `
      <form id="testForm">
        <input name="user" id="user" />
        <input name="email" id="email" />
        <input name="pwd" />
      </form>
    `
    return document.getElementById('testForm')
  }

  const addRules = () => {
    store.uno.valid.addRules('user', ['required'])
    store.uno.valid.addRules('email', ['email'])
    store.uno.valid.addRules('pwd', ['required'])
  }

  const summaryOf = (form) => form.querySelector('.ui-error-summary')
  const titleOf = (summary) => document.getElementById(summary.getAttribute('aria-labelledby'))

  it('多错失败 → 表单顶部渲染摘要：role=alert / tabindex=-1 / aria-labelledby 指向标题', async () => {
    Alpine.store('i18n', {
      t: (key) => (key === 'valid.summary_title' ? '请修正以下问题' : key)
    })
    const form = setupForm()
    document.getElementById('email').value = 'not-an-email'
    addRules()

    const ok = await store.uno.valid.validateForm(form)

    expect(ok).toBe(false)
    const summary = summaryOf(form)
    expect(summary).not.toBeNull()
    expect(form.firstElementChild).toBe(summary) // 表单顶部
    expect(summary.getAttribute('role')).toBe('alert')
    expect(summary.getAttribute('tabindex')).toBe('-1')
    expect(titleOf(summary).textContent).toBe('请修正以下问题')
    // 内联错误保留（摘要不替代）
    expect(store.uno.valid.errors['user']).toBeTruthy()
    expect(store.uno.valid.errors['pwd']).toBeTruthy()
  })

  it('每条错误一条链接：href 锚点优先 input.id，无 id 回退 {field}-error', async () => {
    const form = setupForm()
    document.getElementById('email').value = 'not-an-email'
    addRules()

    await store.uno.valid.validateForm(form)

    const links = summaryOf(form).querySelectorAll('a')
    expect(links).toHaveLength(3)
    expect(links[0].getAttribute('href')).toBe('#user')
    expect(links[1].getAttribute('href')).toBe('#email')
    expect(links[2].getAttribute('href')).toBe('#pwd-error') // 无 id 回退错误元素锚点
    expect(links[0].textContent).toBe(store.uno.valid.errors['user']) // 链接文字即错误消息
  })

  it('焦点移至摘要（而非首错字段）', async () => {
    const form = setupForm()
    document.getElementById('email').value = 'not-an-email'
    addRules()

    await store.uno.valid.validateForm(form)

    const summary = summaryOf(form)
    expect(summary.dataset.testFocus).toBe('1')
    expect(summary.dataset.testScroll).toBe('1')
    expect(document.getElementById('user').dataset.testFocus).toBeUndefined()
  })

  it('单字段失败 → 不渲染摘要，焦点仍到首错字段（现状保持）', async () => {
    const form = setupForm()
    store.uno.valid.addRules('user', ['required'])

    await store.uno.valid.validateForm(form)

    expect(summaryOf(form)).toBeNull()
    expect(document.getElementById('user').dataset.testFocus).toBe('1')
    expect(document.getElementById('user').dataset.testScroll).toBe('1')
  })

  it('点击摘要链接 → 聚焦对应字段并阻止原生跳转（GOV.UK 行为）', async () => {
    const form = setupForm()
    document.getElementById('email').value = 'not-an-email'
    addRules()
    await store.uno.valid.validateForm(form)

    const link = summaryOf(form).querySelectorAll('a')[0] // → #user
    const evt = new dom.window.MouseEvent('click', { bubbles: true, cancelable: true })
    link.dispatchEvent(evt)

    expect(evt.defaultPrevented).toBe(true)
    expect(document.getElementById('user').dataset.testFocus).toBe('1')
  })

  it('重提交成功 → 移除旧摘要；连续失败 → 至多一个摘要（幂等）', async () => {
    const form = setupForm()
    document.getElementById('email').value = 'not-an-email'
    addRules()

    await store.uno.valid.validateForm(form)
    await store.uno.valid.validateForm(form) // 连续失败重提交
    expect(form.querySelectorAll('.ui-error-summary')).toHaveLength(1)

    document.getElementById('user').value = 'tom'
    document.getElementById('email').value = 'tom@example.com'
    document.querySelector('[name="pwd"]').value = 'secret'
    const ok = await store.uno.valid.validateForm(form)
    expect(ok).toBe(true)
    expect(summaryOf(form)).toBeNull()
  })

  it('scrollOnError=false → 摘要仍渲染但不抢焦点（配置语义保持）', async () => {
    const form = setupForm()
    document.getElementById('email').value = 'not-an-email'
    addRules()
    Uno.options.valid = { scrollOnError: false }

    await store.uno.valid.validateForm(form)

    const summary = summaryOf(form)
    expect(summary).not.toBeNull()
    expect(summary.dataset.testFocus).toBeUndefined()
    expect(summary.dataset.testScroll).toBeUndefined()
  })

  it('无 i18n store → 标题回退英文兜底', async () => {
    const form = setupForm()
    document.getElementById('email').value = 'not-an-email'
    addRules()

    await store.uno.valid.validateForm(form)

    expect(titleOf(summaryOf(form)).textContent).toBe('Please fix the following problems')
  })
})
