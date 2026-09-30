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
    // sr-only 关闭文案走 x-t 本地化（复用 uno 既有 box.close），不再硬编码 "Close"
    expect(template).toContain('x-t="box.close"')
    expect(template).not.toContain('<span class="sr-only">Close</span>')
  })
})
