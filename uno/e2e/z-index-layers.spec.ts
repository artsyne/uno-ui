import { test, expect } from '@playwright/test'
import { mockCdnAssets } from './helpers/routes'

test.describe('Z-Index Layers Hierarchy', () => {
  test.beforeEach(async ({ page }) => {
    // Mock UMD if needed, or rely on file
    await mockCdnAssets(page)

    page.on('console', (msg) => console.log(`[Browser Console] ${msg.text()}`))
    await page.goto('/examples/debug-box-tooltip.html')
  })

  test('Box < Tooltip < Notify', async ({ page }) => {
    // 1. Open Box
    await page.evaluate(() => {
      ;(window as any).Alpine.store('box').open('dialog', {
        title: 'Box with Tooltip',
        trustedHTML: true,
        content: `<div x-data class='p-4 space-y-4'>
            <p class='text-slate-700 dark:text-slate-300'>Hover over the button below to see the tooltip.</p>
            <button class='px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600 transition-colors'
                    x-tooltip='&quot;I should be ON TOP of the box&quot;'>
                Hover Me
            </button>
            <div class='pt-4 border-t border-slate-200 dark:border-slate-700'>
                <button class='px-4 py-2 bg-green-500 text-white rounded hover:bg-green-600 transition-colors'
                        @click='$store.notify.show(&quot;I am a notification&quot;, { type: &quot;success&quot;, timeout: 0 })'>
                    Trigger Notify
                </button>
            </div>
        </div>`
      })
    })
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
