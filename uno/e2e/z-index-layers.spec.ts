import { test, expect } from '@playwright/test'

test.describe('Z-Index Layers Hierarchy', () => {
  test.beforeEach(async ({ page }) => {
    // Mock UMD if needed, or rely on file
    await page.route('**/dist/uno.umd.js', async (route) => {
      await route.fulfill({ path: 'dist/uno.umd.js' })
    })

    page.on('console', (msg) => console.log(`[Browser Console] ${msg.text()}`))
    await page.goto('/examples/debug-box-tooltip.html')
  })

  test('Box < Tooltip < Notify', async ({ page }) => {
    // 1. Open Box
    await page.getByText('Open Box').click()
    const boxContainer = page.locator('[id^="uno-box-outer-container"]').first()
    await expect(boxContainer).toBeVisible()

    // 2. Trigger Notify (inside Box)
    await page.getByText('Trigger Notify').click()
    const notifyItem = page.locator('[data-testid="notify-item"]')
    await expect(notifyItem).toBeVisible()
    
    // 3. Trigger Tooltip (inside Box)
    const tooltipBtn = page.getByText('Hover Me')
    await tooltipBtn.hover()
    const tooltip = page.locator('#uno-tooltip-content')
    await expect(tooltip).toBeVisible()
    await expect(page.locator('#uno-tooltip-body')).toHaveText('I should be ON TOP of the box')

    // 4. Verify Z-Indices
    
    // Box
    const boxZ = await boxContainer.evaluate((el) => {
      return parseInt(window.getComputedStyle(el).zIndex)
    })
    console.log('Box Z-Index:', boxZ)

    // Tooltip
    const tooltipZ = await tooltip.evaluate((el) => {
      return parseInt(window.getComputedStyle(el).zIndex)
    })
    console.log('Tooltip Z-Index:', tooltipZ)

    // Notify
    // The notify item is inside a fixed container which has the z-index
    const notifyContainer = notifyItem.locator('..') // Parent
    const notifyZ = await notifyContainer.evaluate((el) => {
      return parseInt(window.getComputedStyle(el).zIndex)
    })
    console.log('Notify Z-Index:', notifyZ)

    // Assertions
    expect(boxZ).toBeGreaterThanOrEqual(10000)
    expect(tooltipZ).toBeGreaterThanOrEqual(20000)
    expect(notifyZ).toBeGreaterThanOrEqual(30000)

    expect(boxZ).toBeLessThan(tooltipZ)
    expect(tooltipZ).toBeLessThan(notifyZ)
  })
})
