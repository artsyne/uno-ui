// @vitest-environment jsdom
// fill 模式：类名由 Alpine :class 运行时合成，走真实运行时断言水合后 DOM。
import { describe, it, expect, beforeAll, afterEach } from 'vitest'
import { bootUno, tick } from '../helpers/axe'

describe('Box fill mode', () => {
  let box
  let Uno

  beforeAll(async () => {
    const booted = await bootUno()
    box = booted.Alpine.store('uno').box
    Uno = booted.Uno
  })

  afterEach(async () => {
    box.closeAll()
    await tick()
  })

  const classesOf = (el) => (el ? el.className.split(/\s+/).filter(Boolean) : [])
  const styleClasses = (key) => Uno.getStyle('box', key).split(/\s+/).filter(Boolean)
  const nodes = (id) => {
    const container = document.getElementById('uno-box-container' + id)
    const body = document.getElementById('uno-box-content-' + id)
    // x-if 渲染的节点插在 <template> 之后，取首个 div 子节点
    return { container, content: body?.parentElement, body, inner: body?.querySelector(':scope > div') }
  }

  it('defaults fill to false', () => {
    const item = box.open({ content: 'x' })
    expect(item.fill).toBe(false)
  })

  it('chains flex classes down to the content container and applies the fill height', async () => {
    const item = box.open({ content: '<p>body</p>', widthClass: 'w-full', fill: true })
    await tick(100)
    const { container, content, body, inner } = nodes(item.id)

    expect(classesOf(container)).toEqual(expect.arrayContaining(styleClasses('containerFill')))
    expect(classesOf(content)).toEqual(expect.arrayContaining(styleClasses('contentFill')))
    expect(classesOf(body)).toEqual(expect.arrayContaining(['mt-2', ...styleClasses('bodyFill')]))
    expect(classesOf(inner)).toEqual(expect.arrayContaining(styleClasses('bodyFill')))
  })

  it('lets heightClass replace the fill height', async () => {
    const item = box.open({ content: 'x', heightClass: 'h-[30rem]', fill: true })
    await tick(100)
    const cls = classesOf(nodes(item.id).container)

    expect(cls).toContain('h-[30rem]')
    styleClasses('containerFill').forEach((c) => expect(cls).not.toContain(c))
  })

  it('lets numeric height replace the fill height', async () => {
    const item = box.open({ content: 'x', height: 480, fill: true })
    await tick(100)
    const { container } = nodes(item.id)

    styleClasses('containerFill').forEach((c) => expect(classesOf(container)).not.toContain(c))
    expect(container.getAttribute('style')).toContain('height:480px')
  })

  it('keeps the block layout when fill is off', async () => {
    const item = box.open({ content: 'x', heightClass: 'h-[30rem]' })
    await tick(100)
    const { content, body } = nodes(item.id)

    expect(classesOf(content)).not.toContain('min-h-0')
    expect(classesOf(body)).toEqual(['mt-2'])
  })

  it('stretches the HTMX loading placeholder in fill mode', async () => {
    window.htmx = { process: () => {} }
    const item = box.open({ hxUrl: '/fragment', fill: true })
    await tick(100)
    const { inner } = nodes(item.id)

    expect(classesOf(inner)).toEqual(expect.arrayContaining(styleClasses('bodyFill')))
    expect(classesOf(inner.firstElementChild)).toContain('flex-1')
    delete window.htmx
  })
})
