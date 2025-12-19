import { test, expect } from '@playwright/test'
import { mockCdnAssets } from './helpers/routes'

test.describe('Uno Global Configuration', () => {
  test.beforeEach(async ({ page }) => {
    await mockCdnAssets(page)

    page.on('console', (msg) => console.log('PAGE LOG:', msg.text()))
    page.on('pageerror', (err) => console.log('PAGE ERROR:', err.message))
    await page.goto('/examples/config-demo.html')
    await page.waitForLoadState('domcontentloaded')
    await page.waitForFunction(
      () => (window as any).Alpine && (window as any).Uno?.initialized === true
    )
    await page.waitForFunction(() => (window as any).Alpine?.store?.('notify') !== undefined)
    await page.waitForFunction(() => (window as any).Alpine?.store?.('tooltip') !== undefined)

    // Debug: Check loaded config
    const unoOptions = await page.evaluate(() => {
      return (window as any).Uno ? (window as any).Uno.options : 'Uno not found'
    })
    console.log('DEBUG: Uno.options on page:', JSON.stringify(unoOptions, null, 2))
  })

  test('should respect global notify configuration (bottom-right)', async ({ page }) => {
    await page.evaluate(() => {
      ;(window as any).Alpine.store('notify').success('This respects global config!', {
        timeout: 0
      })
    })

    const item = page
      .locator('[data-testid="notify-item"]')
      .filter({ hasText: 'This respects global config!' })
    await expect(item).toBeVisible()

    const region = item.locator('xpath=ancestor::div[contains(@class,"fixed")][1]')
    const regionClass = await region.getAttribute('class')
    expect(regionClass || '').toContain('bottom-0')
    expect(regionClass || '').toContain('right-0')
    expect(regionClass || '').not.toContain('left-0')
  })

  test('should respect global box configuration (width: 600px)', async ({ page }) => {
    // Click button to show box
    await page.getByText('Test Box').click()

    // Expect modal to be visible
    const modal = page.locator('[role="dialog"]')
    await expect(modal).toBeVisible()

    const panel = page.locator('[id^="uno-box-container"]:not([id="uno-box-container"])').first()
    await expect(panel).toBeVisible()
    const style = await panel.getAttribute('style')
    expect(style || '').toMatch(/width:\s*600px/i)
  })

  test('should respect global tooltip configuration (placement: right, trigger: click)', async ({
    page
  }) => {
    // Inject a NEW button to test initialization with current config
    // Note: Existing elements won't update their event listeners dynamically
    await page.evaluate(() => {
      const btn = document.createElement('button')
      btn.id = 'dynamic-tooltip-btn'
      btn.setAttribute('x-tooltip', "'Dynamic Tooltip'")
      btn.textContent = 'Dynamic Tooltip'
      btn.className = 'mt-4 px-4 py-2 bg-blue-500 text-white rounded' // Styling to make it clickable
      document.body.appendChild(btn)
      if ((window as any).Alpine?.initTree) {
        ;(window as any).Alpine.initTree(btn)
      }
    })

    await page.evaluate(() => {
      try {
        ;(window as any).Alpine?.store('tooltip')?.hide?.()
      } catch {}
    })

    const trigger = page.locator('#dynamic-tooltip-btn')
    await trigger.scrollIntoViewIfNeeded()

    // Hover shouldn't trigger it (since trigger is click)
    await trigger.hover()
    await expect(page.locator('#uno-tooltip-content')).toBeHidden()

    // Click should trigger it
    await trigger.click()
    await expect(page.locator('#uno-tooltip-content')).toBeVisible()
    await expect(page.locator('#uno-tooltip-body')).toHaveText('Dynamic Tooltip')

    // Check placement (approximate by coordinates)
    const triggerBox = await trigger.boundingBox()
    const tooltipBox = await page.locator('#uno-tooltip-content').boundingBox()

    if (triggerBox && tooltipBox) {
      expect(tooltipBox.x).toBeGreaterThan(triggerBox.x + triggerBox.width - 10) // Allow small margin
    }
  })

  test('should support runtime configuration updates via Uno.setConfig()', async ({ page }) => {
    // Update config and test
    await page.getByText('Update & Test').click()

    // Wait for container
    const container = page.locator('#uno-notify-container')
    await container.waitFor({ state: 'attached' })

    const item = page.getByText('I am now top-center!')
    await expect(item).toBeVisible()

    const region = item.locator('xpath=ancestor::div[contains(@class,"fixed")][1]')
    const regionClass = await region.getAttribute('class')
    expect(regionClass || '').toContain('top-0')
    expect(regionClass || '').toContain('left-0')
    expect(regionClass || '').toContain('right-0')
  })
})
