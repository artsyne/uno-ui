// @vitest-environment jsdom
// P2 组件级 axe 门禁（docs/a11y-design.md §3.5）：Tab 标签页水合后 DOM 零违规。
// fixture 采用项目规范用法（ul[role=tablist] > li[role=presentation] > a[x-uno-tab]），
// role=tab/aria-selected/aria-controls/roving tabindex 均由指令运行时注入——
// 静态模板扫描完全看不到，故走真实运行时。
import { describe, it, expect } from 'vitest'
import { bootUno, scanA11y, tick } from '../helpers/axe'

describe('Tab axe scan (P2 a11y gate)', () => {
  it('renders an APG tablist with zero axe violations', async () => {
    const { Alpine } = await bootUno()

    document.body.insertAdjacentHTML(
      'beforeend',
      `
      <ul role="tablist" aria-label="Sections">
        <li role="presentation"><a id="fixture-tab-a" href="#panel-a" x-uno-tab="{ name: 'fixture-a', group: 'axe' }">Tab A</a></li>
        <li role="presentation"><a id="fixture-tab-b" href="#panel-b" x-uno-tab="{ name: 'fixture-b', group: 'axe' }">Tab B</a></li>
      </ul>
      <div id="panel-fixture-a" role="tabpanel" aria-labelledby="fixture-tab-a"><p>Panel A</p></div>
      <div id="panel-fixture-b" role="tabpanel" aria-labelledby="fixture-tab-b" hidden><p>Panel B</p></div>
      `
    )
    Alpine.initTree(document.querySelector('ul[role="tablist"]'))
    await tick()

    // 激活首个 tab：验证指令运行时 ARIA 注入 + 面板关联（防扫描空 DOM 假绿）
    document.getElementById('fixture-tab-a').click()
    await tick()

    const tabA = document.getElementById('fixture-tab-a')
    expect(tabA.getAttribute('role')).toBe('tab')
    expect(tabA.getAttribute('aria-selected')).toBe('true')
    expect(tabA.getAttribute('aria-controls')).toBe('panel-fixture-a')

    const results = await scanA11y()
    expect(results).toHaveNoViolations()
  })
})
