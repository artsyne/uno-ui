import { test, expect } from '@playwright/test'
import { mockCdnAssets } from './helpers/routes'

test.describe('Tab Component HTMX', () => {
  test.beforeEach(async ({ page }) => {
    await mockCdnAssets(page)

    // Enable console logging
    page.on('console', (msg) => console.log(`[Browser Console] ${msg.text()}`))
    page.on('pageerror', (err) => console.log(`[Page Error] ${err.message}`))

    await page.goto('/examples/components/tab.html')
    await page.waitForLoadState('domcontentloaded')

    await page.waitForFunction(() => (window as any).Alpine && (window as any).Uno)
  })

  test('htmx content loading', async ({ page }) => {
    const htmxTab = page.locator('#tab-side-htmx')
    const htmxContent = page.locator('#panel-side-htmx-target')

    // Initially hidden
    await expect(htmxContent).toBeHidden()

    // 关键点：确保点击后确实触发了对 tab-content.html 的请求
    const contentResponse = page.waitForResponse((resp) => {
      return resp.url().includes('/examples/data/tab-content.html') && resp.status() === 200
    })

    // Click to load
    await htmxTab.click()

    // Should be visible
    await expect(htmxContent).toBeVisible()

    await contentResponse

    // Should load content
    await expect(htmxContent).toContainText('Ajax Loaded Content', { timeout: 5000 })
  })

  test('AJAX Content: loads content into its own panel', async ({ page }) => {
    const ajaxTab = page.locator('#tab-side-ajax')
    const ajaxContent = page.locator('#panel-side-ajax')

    await expect(ajaxContent).toBeHidden()

    // 关键点：这个 case 之前经常出现“永远 Loading”，这里直接用网络响应 + 文本双重断言
    const contentResponse = page.waitForResponse((resp) => {
      return resp.url().includes('/examples/data/tab-content.html') && resp.status() === 200
    })

    await ajaxTab.click()

    await expect(ajaxContent).toBeVisible()
    await contentResponse
    await expect(ajaxContent).toContainText('Ajax Loaded Content', { timeout: 5000 })
    await expect(ajaxContent).not.toContainText('Loading...', { timeout: 5000 })
  })
})
