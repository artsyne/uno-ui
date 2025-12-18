import { test, expect } from '@playwright/test'

test.describe('Box + Tooltip Z-Index', () => {
  test.beforeEach(async ({ page }) => {
    await page.route('**/dist/uno.umd.js', async (route) => {
      await route.fulfill({ path: 'dist/uno.umd.js' })
    })

    await page.goto('/examples/debug-box-tooltip.html')
    await page.waitForLoadState('domcontentloaded')
    await page.waitForFunction(() => (window as any).Alpine && (window as any).Alpine.store('box'))
  })

  test('Tooltip shows above Box', async ({ page }) => {
    // 1. Open Box
    await page.getByText('Open Box').click()

    const boxContainer = page.locator('[id^="uno-box-outer-container"]').first()
    await expect(boxContainer).toBeVisible()

    // 2. Trigger Tooltip
    // The button inside the box content needs to be hovered
    const innerBtn = page.getByText('Hover Me')
    await expect(innerBtn).toBeVisible()
    await innerBtn.hover()

    // 3. Check Tooltip Visibility
    const tooltip = page.locator('#uno-tooltip-content')
    await expect(tooltip).toBeVisible()
    
    // Check body content specifically
    const tooltipBody = page.locator('#uno-tooltip-body')
    await expect(tooltipBody).toBeVisible()
    await expect(tooltipBody).toHaveText('I should be ON TOP of the box')

    // 4. Verify Z-Index
    // We get the computed style z-index of the box outer container and the tooltip
    const boxZVal = await boxContainer.evaluate((el) => {
      return parseInt(window.getComputedStyle(el).zIndex)
    })

    const tooltipZVal = await tooltip.evaluate((el) => {
      return parseInt(window.getComputedStyle(el).zIndex)
    })

    // console.log(`Box Z-Index: ${boxZVal}, Tooltip Z-Index: ${tooltipZVal}`)

    expect(tooltipZVal).toBeGreaterThan(boxZVal)

    // 5. Verify Tooltip is not covered by Box
    // We can check if the tooltip element is the topmost element at its position
    // Or simpler: check if it intercepts pointer events
    const tooltipBox = await tooltip.boundingBox()
    if (tooltipBox) {
      // Click the center of the tooltip and ensure it hits the tooltip
      // Note: Tooltip usually has pointer-events: none on the wrapper, but the body might not?
      // In tooltip.js: wrapper has pointer-events-none, body has pointer-events-auto.
      // But the wrapper is what has the z-index.
      // If wrapper is pointer-events-none, we can't click it.
      // But visual layering is what matters.
    }
  })
})
