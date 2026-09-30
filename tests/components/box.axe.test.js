// @vitest-environment jsdom
// P2 组件级 axe 门禁（docs/a11y-design.md §3.5）：Box 对话框水合后 DOM 零违规。
// role="dialog"/aria-modal/aria-labelledby 由 Alpine :binding 运行时解析，
// 关闭按钮可访问名来自 x-t.aria-label——静态模板扫描会漏检/误报，故走真实运行时。
import { describe, it, expect } from 'vitest'
import { bootUno, scanA11y, tick } from '../helpers/axe'

describe('Box axe scan (P2 a11y gate)', () => {
  it('renders an open dialog with zero axe violations', async () => {
    const { Alpine } = await bootUno()
    const box = Alpine.store('uno').box

    box.open({ title: 'Dialog title', content: '<p>Dialog body</p>' })
    await tick(150)

    // 前置断言：对话框确已水合渲染（防扫描空 DOM 假绿）
    expect(document.querySelector('[role="dialog"]')).not.toBeNull()

    const results = await scanA11y()
    expect(results).toHaveNoViolations()
  })
})
