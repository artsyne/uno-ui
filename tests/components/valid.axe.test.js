// @vitest-environment jsdom
// 组件级 axe 门禁（docs/a11y-design.md §3.5）：Valid 校验错误态水合后 DOM 零违规。
// 错误元素（role=alert + 稳定 id）、输入框 aria-invalid/aria-describedby 联动均由
// A1 在运行时注入——静态模板扫描完全看不到，故走真实运行时。
import { describe, it, expect } from 'vitest'
import { bootUno, scanA11y, tick } from '../helpers/axe'

describe('Valid axe scan (a11y gate)', () => {
  it('renders a validation error state with zero axe violations', async () => {
    const { Alpine } = await bootUno()

    document.body.insertAdjacentHTML(
      'beforeend',
      `
      <form id="fixture-form">
        <label for="fixture-email">Email</label>
        <input id="fixture-email" name="fixture-email" type="email" x-uno-valid="['required', 'email']" />
      </form>
      `
    )
    Alpine.initTree(document.getElementById('fixture-form'))
    await tick()

    const input = document.getElementById('fixture-email')
    input.value = 'not-an-email'
    input.dispatchEvent(new Event('input', { bubbles: true }))
    await tick(150)

    // 前置断言：错误元素已生成且 ARIA 关联到位（A1 契约，防扫描空 DOM 假绿）
    const errorEl = document.getElementById('fixture-email-error')
    expect(errorEl).not.toBeNull()
    expect(errorEl.getAttribute('role')).toBe('alert')
    expect(input.getAttribute('aria-invalid')).toBe('true')
    expect((input.getAttribute('aria-describedby') || '').split(/\s+/)).toContain(
      'fixture-email-error'
    )

    const results = await scanA11y()
    expect(results).toHaveNoViolations()
  })
})
