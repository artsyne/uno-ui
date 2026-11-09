// @vitest-environment jsdom
// 组件级 axe 门禁（docs/a11y-design.md §3.5）：Tooltip 弹层水合后 DOM 零违规。
// 弹层由全局容器运行时渲染，ARIA 关联（aria-describedby/expanded）由 A3 指令注入——
// 静态模板扫描不可见，故走真实运行时。
import { describe, it, expect } from 'vitest'
import { bootUno, scanA11y, tick } from '../helpers/axe'

describe('Tooltip axe scan (a11y gate)', () => {
  it('renders a click-triggered tooltip with zero axe violations', async () => {
    const { Alpine } = await bootUno()

    document.body.insertAdjacentHTML(
      'beforeend',
      `<button type="button" id="fixture-tt" x-uno-tooltip="{ content: 'Helpful hint', trigger: 'click' }">Help</button>`
    )
    Alpine.initTree(document.getElementById('fixture-tt'))
    await tick()

    document.getElementById('fixture-tt').click()
    await tick(150)

    // 前置断言：弹层确已水合渲染（防扫描空 DOM 假绿）
    expect(document.querySelector('[role="tooltip"]')).not.toBeNull()

    const results = await scanA11y()
    expect(results).toHaveNoViolations()
  })
})
