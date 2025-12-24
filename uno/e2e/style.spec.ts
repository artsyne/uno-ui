import { test, expect } from '@playwright/test'
import { mockCdnAssets } from './helpers/routes'

test.describe('Uno UI Style Customization', () => {
  test.beforeEach(async ({ page }) => {
    // Force serve latest build
    await mockCdnAssets(page)

    // Use a basic page
    await page.goto('/examples/components/box.html')
    await page.waitForFunction(() => (window as any).Alpine && (window as any).Uno)
  })

  test('Should apply global style overrides (rounded-none)', async ({ page }) => {
    // 1. Set global config BEFORE using components
    await page.evaluate(() => {
      ;(window as any).Uno.setConfig({
        ui: {
          base: 'rounded-none shadow-md border border-gray-200 dark:border-gray-700'
        }
      })
    })

    const base = await page.evaluate(() => {
      return (window as any).Uno.options.ui.base
    })
    expect(typeof base).toBe('string')
    expect(base).toContain('rounded-none')

    // 2. Box
    await page.evaluate(() => {
      ;(window as any).Alpine.store('box').alert('Test Box')
    })
    // Select the actual box container, not the wrapper or outer backdrop
    const box = page.locator('[id^="uno-box-outer-container"] [id^="uno-box-container"]').first()
    await expect(box).toBeVisible()
    await expect(box).toHaveClass(/rounded-none/)
    await expect(box).not.toHaveClass(/rounded-lg/)

    // Close box
    await page.evaluate(() => (window as any).Alpine.store('box').closeAll())

    // 3. Notify
    await page.evaluate(() => {
      ;(window as any).Alpine.store('notify').show('Test Notify')
    })
    const notify = page.locator('[data-testid="notify-item"]').first()
    await expect(notify).toBeVisible()
    // Notify container class is constructed using rounded
    await expect(notify).toHaveClass(/rounded-none/)
    await expect(notify).not.toHaveClass(/rounded-lg/) // Default is rounded (which usually implies rounded, or rounded-md/lg?)

    // 4. Tooltip
    // Trigger tooltip
    await page.evaluate(() => {
      const btn = document.createElement('button')
      btn.id = 'tooltip-btn'
      btn.textContent = 'Hover me'
      document.body.appendChild(btn)
      ;(window as any).Alpine.store('tooltip').show(btn, 'Tooltip Content')
    })
    const tooltip = page.locator('#uno-tooltip-content')
    await expect(tooltip).toBeVisible()
    // 5. getStyle append mode
    await page.evaluate(() => {
      if (!(window as any).Uno.options.ui.box) (window as any).Uno.options.ui.box = {}
      ;(window as any).Uno.options.ui.box.okBtn = '+test-append-class'
    })
    const result = await page.evaluate(() => {
      return (window as any).Uno.getStyle('box', 'okBtn', 'default-class')
    })
    expect(result).toBe('default-class test-append-class')
  })

  test('Tab should apply custom active class', async ({ page }) => {
    await page.evaluate(() => {
      ;(window as any).Uno.setConfig({
        ui: {
          tab: {
            triggerActive: 'text-red-600 font-bold',
            triggerInactive: 'text-gray-500'
          }
        }
      })

      const div = document.createElement('div')
      div.innerHTML = `
            <div x-data="unoTabs('tab1')">
                <div role="tablist">
                    <button id="t1" x-bind="trigger('tab1')">Tab 1</button>
                    <button id="t2" x-bind="trigger('tab2')">Tab 2</button>
                </div>
                <div x-bind="panel('tab1')">1</div>
                <div x-bind="panel('tab2')">2</div>
            </div>
        `
      document.body.appendChild(div)
      ;(window as any).Alpine.initTree(div)
    })

    const t1 = page.locator('#t1')
    const t2 = page.locator('#t2')

    await expect(t1).toHaveClass(/text-red-600/)
    await expect(t1).toHaveClass(/font-bold/)
    await expect(t2).toHaveClass(/text-gray-500/)

    // Switch
    await t2.click()

    await expect(t2).toHaveClass(/text-red-600/)
    await expect(t1).toHaveClass(/text-gray-500/)
  })
})
