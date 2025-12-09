import { test, expect } from '@playwright/test'

test.describe('Tab Component HTMX', () => {
  test.beforeEach(async ({ page }) => {
    // Enable console logging
    page.on('console', (msg) => console.log(`[Browser Console] ${msg.text()}`))
    page.on('pageerror', (err) => console.log(`[Page Error] ${err.message}`))

    await page.goto('/examples/components/tab.html')
    await page.waitForLoadState('domcontentloaded')

    // Wait for Alpine to initialize
    try {
      await page.waitForFunction(() => (window as any).Alpine, null, { timeout: 5000 })
    } catch (e) {
      console.log('Alpine not found on window!')
    }
  })

  test('htmx content loading', async ({ page }) => {
    const htmxTab = page.locator('button[data-tab-hx-url="/examples/data/tab-content-htmx.html"]')
    const htmxContent = page.locator('#side-htmx')

    // Initially hidden
    await expect(htmxContent).toBeHidden()

    // Click to load
    await htmxTab.click()

    // Should be visible
    await expect(htmxContent).toBeVisible()

    // Should load content
    await expect(htmxContent).toContainText('HTMX Loaded Content', { timeout: 5000 })

    // Verify loaded attribute
    await expect(htmxContent).toHaveAttribute('data-tab-loaded', 'true')
  })

  test('htmx error handling', async ({ page }) => {
    // Mock htmx.ajax to simulate failure
    await page.evaluate(() => {
      ;(window as any).htmx.ajax = () => Promise.reject(new Error('Simulated HTMX Error'))
    })

    const htmxTab = page.locator('button[data-tab-hx-url="/examples/data/tab-content-htmx.html"]')
    const htmxContent = page.locator('#side-htmx')

    await htmxTab.click()

    // Should show error message
    await expect(htmxContent).toBeVisible()
    await expect(htmxContent).toContainText('Error loading content via HTMX', { timeout: 5000 })
  })
})
