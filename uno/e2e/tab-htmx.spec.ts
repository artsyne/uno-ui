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
    // Force network error for the HTMX content URL
    await page.route('**/tab-content-htmx.html', (route) => route.abort('failed'))

    const htmxTab = page.locator('button[data-tab-hx-url="/examples/data/tab-content-htmx.html"]')
    const htmxContent = page.locator('#side-htmx')

    await htmxTab.click()

    // Should show error message
    await expect(htmxContent).toBeVisible()
    await expect(htmxContent).toContainText('Error loading content via HTMX', { timeout: 5000 })
  })
})
