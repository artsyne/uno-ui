import { describe, it, expect } from 'vitest'

// box closeBtnFloating 图标色非文本对比度（E7，WCAG 1.4.11，≥3:1）：亮色 slate-500
// （hover 加深 slate-600）；暗底（gray-900）配对 slate-400 / hover slate-300。
// 断言用 token 精确匹配，避开 'dark:text-slate-400' 对裸 'text-slate-400' 的子串误命中。
describe('box closeBtnFloating contrast (WCAG 1.4.11)', () => {
  it('默认配置 token：亮暗两态图标色均 ≥3:1，禁止回退 slate-400 亮色', async () => {
    const { UnoInstance } = await import('../../src/core/uno')
    const uno = new UnoInstance()
    const tokens = uno.options.ui.box.closeBtnFloating.split(/\s+/)

    // 亮色：白底常态灰阶须 ≥3:1，hover 加深
    expect(tokens).toContain('text-slate-500')
    expect(tokens).toContain('hover:text-slate-600')
    // 暗色配对：深底提亮
    expect(tokens).toContain('dark:text-slate-400')
    expect(tokens).toContain('dark:hover:text-slate-300')
    // 禁止回退：亮色 slate-400（<3:1）/ hover slate-500（无加深）
    expect(tokens).not.toContain('text-slate-400')
    expect(tokens).not.toContain('hover:text-slate-500')
  })
})
