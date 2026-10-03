import { describe, it, expect, vi } from 'vitest'
import { getNotifyTemplate } from '../../src/components/notify/notify.template'

// getNotifyTemplate 返回纯字符串；getZIndex 仅被 getProgressTemplate 使用，
// 此处 mock core/layer 以隔离模块导入副作用。
vi.mock('../../src/core/layer', () => ({
  getZIndex: vi.fn(() => 5000)
}))

describe('Notify Template', () => {
  it('should localize close button label via i18n instead of hardcoding (B7)', () => {
    const template = getNotifyTemplate()
    // sr-only 关闭文案走 x-t 本地化（复用 uno 既有 box.close），不硬编码 "Close"
    expect(template).toContain('x-t="box.close"')
    expect(template).not.toContain('<span class="sr-only">Close</span>')
  })

  // B7：成功类通知一律 assertive 会打断当前朗读；播报策略按类型分级——
  // error/warning 保留 alert+assertive，success/info 降为 status+polite。
  it('should scale politeness by type: success/info polite status, error/warning assertive alert (B7)', () => {
    const template = getNotifyTemplate()
    expect(template).toContain(
      "item.type === 'error' || item.type === 'warning' ? 'assertive' : 'polite'"
    )
    expect(template).toContain(
      "item.type === 'error' || item.type === 'warning' ? 'alert' : 'status'"
    )
    expect(template).not.toContain('aria-live="assertive"')
  })
})
