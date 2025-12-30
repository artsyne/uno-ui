import { test, expect } from '@playwright/test'
import { mockCdnAssets } from './helpers/routes'

test.describe('Tab Component', () => {
  test.beforeEach(async ({ page }) => {
    await mockCdnAssets(page)

    // Enable console logging
    page.on('console', (msg) => console.log(`[Browser Console] ${msg.text()}`))
    page.on('pageerror', (err) => console.log(`[Page Error] ${err.message}`))

    await page.goto('/examples/components/tab.html')
    await page.waitForLoadState('domcontentloaded')

    await page.waitForFunction(() => (window as any).Alpine && (window as any).Uno)
  })

  test('modern tabs switch content and hash', async ({ page }) => {
    // Selectors based on x-bind: trigger('settings') -> id="tab-settings"
    const tab2 = page.locator('#tab-settings')
    const content1 = page.locator('#panel-tab1')
    const content2 = page.locator('#panel-settings')

    // Initial state
    await expect(content1).toBeVisible()
    await expect(content2).toBeHidden()

    // Click tab 2
    await tab2.click()

    // Verify switch
    await expect(content1).toBeHidden()
    await expect(content2).toBeVisible()

    // Verify Hash
    expect(page.url()).toContain('#settings')
  })

  test('separated dom elements work', async ({ page }) => {
    const sTab2 = page.locator('#tab-s-tab2')
    const sContent1 = page.locator('#panel-s-tab1')
    const sContent2 = page.locator('#panel-s-tab2')

    await sTab2.click()
    await expect(sContent1).toBeHidden()
    await expect(sContent2).toBeVisible()
  })

  test('deep linking loads correct tab', async ({ page }) => {
    // Navigate directly with hash
    await page.goto('/examples/components/tab.html#billing')
    await page.waitForLoadState('domcontentloaded')
    await page.waitForFunction(() => (window as any).Alpine)

    const content3 = page.locator('#panel-billing')
    const tab3 = page.locator('#tab-billing')

    await expect(content3).toBeVisible()
    await expect(tab3).toHaveAttribute('aria-selected', 'true')
  })

  test('htmx content loading (tab panel)', async ({ page }) => {
    const tab = page.locator('#tab-s-tab3')
    const panel = page.locator('#panel-s-tab3')

    // Initially hidden (it is a tab panel, activeTab is s-tab1 initially)
    await expect(panel).toBeHidden()

    // Click to load
    await tab.click()

    // Should be visible
    await expect(panel).toBeVisible()

    // Should load content via HTMX
    // Note: tab.html uses /examples/data/tab-content.html which contains "Ajax Loaded Content"
    await expect(panel).toContainText('Ajax Loaded Content')
  })

  test('S-Tab 3: refresh with hash loads content and clears loading', async ({ page }) => {
    // Enable debug logging
    await page.addInitScript(() => {
      // @ts-ignore
      window.UnoConfig = { debug: true }
    })

    // 关键点：模拟用户直接刷新到带 hash 的 URL（历史上这里容易出现一直 Loading）
    const contentResponse = page.waitForResponse((resp) => {
      return resp.url().includes('/examples/data/tab-content.html') && resp.status() === 200
    })

    await page.goto('/examples/components/tab.html#s-tab3')
    await page.waitForLoadState('domcontentloaded')
    
    // Log current hash for debugging
    const hash = await page.evaluate(() => window.location.hash)
    console.log('Test Page Hash:', hash)

    await page.waitForFunction(() => (window as any).Alpine && (window as any).Uno)

    const panel = page.locator('#panel-s-tab3')
    await expect(panel).toBeVisible()

    // 等待 HTMX 请求确实发起并成功返回，避免只靠文本断言造成的偶发误判
    await contentResponse

    await expect(panel).toContainText('Ajax Loaded Content', { timeout: 5000 })
    await expect(panel).not.toContainText('Loading...', { timeout: 5000 })
  })
})
