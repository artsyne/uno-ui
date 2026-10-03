// @vitest-environment jsdom
// 组件级 axe 门禁（docs/a11y-design.md §3.5）：Notify 通知项水合后 DOM 零违规。
// role/aria-live 由 :binding 按 type 运行时解析（success→status/polite），
// 关闭按钮可访问名来自 sr-only + x-t——静态模板扫描不可见，故走真实运行时。
import { describe, it, expect } from 'vitest'
import { bootUno, scanA11y, tick } from '../helpers/axe'

describe('Notify axe scan (a11y gate)', () => {
  it('renders a visible notification with zero axe violations', async () => {
    const { Alpine } = await bootUno()
    const notify = Alpine.store('uno').notify

    notify.success('Saved successfully')
    await tick(150)

    // 前置断言：通知项确已水合渲染（防扫描空 DOM 假绿）
    expect(document.querySelector('[data-testid="notify-item"]')).not.toBeNull()

    const results = await scanA11y()
    expect(results).toHaveNoViolations()
  })
})
