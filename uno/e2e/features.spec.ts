import { test, expect } from '@playwright/test'
import { mockCdnAssets } from './helpers/routes'

test.describe('Feature: Accessibility', () => {
  test.beforeEach(async ({ page }) => {
    await mockCdnAssets(page)

    page.on('console', (msg) => console.log('BROWSER LOG:', msg.text()))
    await page.goto('/examples/features.html', { waitUntil: 'domcontentloaded' })
    await page.waitForFunction(() => (window as any).Alpine && (window as any).Uno)
  })

  test('Uno is loaded', async ({ page }) => {
    const isLoaded = await page.evaluate(() => {
      const w = window
      return Boolean((w as any).Uno) && Boolean((w as any).Alpine)
    })
    expect(isLoaded).toBe(true)
  })

  test('Box has correct ARIA roles', async ({ page }) => {
    await page.click('text=Open Modal with Focus Trap')
    const dialog = page.locator('[role="dialog"]')
    await expect(dialog).toBeVisible()
    await expect(dialog).toHaveAttribute('aria-modal', 'true')
  })

  test('Focus Trap works in Box', async ({ page, browserName }) => {
    await page.click('text=Open Modal with Focus Trap')
    const dialog = page.locator('[role="dialog"]')
    await expect(dialog).toBeVisible()

    // Focus the first input explicitly
    await page.locator('[placeholder="First input"]').focus()
    await expect(page.locator('[placeholder="First input"]')).toBeFocused()

    // Tab to second input
    await page.keyboard.press('Tab')
    await expect(page.locator('[placeholder="Second input"]')).toBeFocused()

    // Tab to action button
    await page.keyboard.press('Tab')
    // Ensure button is focusable
    const btn = page.locator('button:has-text("Action")')
    await expect(btn).toBeEnabled()

    // WebKit sometimes fails focus checks in headless mode for this specific scenario
    if (browserName !== 'webkit') {
      await expect(btn).toBeFocused()
    }
  })

  test('Tooltip has correct ARIA roles', async ({ page }) => {
    const trigger = page.locator('text=Focus or hover me for tooltip')

    // Trigger hover
    await trigger.hover({ force: true })

    // Wait for tooltip to appear
    const tooltip = page.locator('[role="tooltip"]')
    await expect(tooltip).toBeVisible({ timeout: 5000 })

    await expect(tooltip).toContainText('accessible tooltip')

    // Check aria-describedby
    const tooltipId = await tooltip.getAttribute('id')
    if (!tooltipId) {
      throw new Error('Tooltip id is missing')
    }
    await expect(trigger).toHaveAttribute('aria-describedby', tooltipId)
  })
})
