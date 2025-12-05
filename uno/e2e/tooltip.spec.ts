import { test, expect } from '@playwright/test'

test.describe('Tooltip Component', () => {
  test.beforeEach(async ({ page }) => {
    page.on('console', (msg) => console.log(`BROWSER LOG: ${msg.text()}`))
    page.on('pageerror', (err) => console.log(`BROWSER ERROR: ${err.message}`))
    page.on('requestfailed', (request) =>
      console.log(`REQUEST FAILED: ${request.url()} - ${request.failure()?.errorText}`)
    )

    await page.goto('/examples/components/tooltip.html')

    // 1. Check if Alpine loads
    try {
      await page.waitForFunction(() => (window as any).Alpine !== undefined, undefined, {
        timeout: 5000
      })
      console.log('Alpine loaded successfully')
    } catch (e) {
      console.log('Alpine failed to load within 5s')
      throw e
    }

    // 2. Check if Store exists
    try {
      await page.waitForFunction(
        () => {
          const win = window as any
          return win.Alpine && win.Alpine.store && win.Alpine.store('tooltip')
        },
        undefined,
        { timeout: 5000 }
      )
      console.log('Tooltip store registered successfully')
    } catch (e) {
      console.log('Tooltip store failed to register within 5s')
      // Log current stores
      const stores = await page.evaluate(() => {
        const win = window as any
        return win.Alpine ? Object.keys(win.Alpine.store || {}) : 'Alpine missing'
      })
      console.log('Available stores:', stores)
      throw e
    }
  })

  test('shows on hover (top)', async ({ page }) => {
    const trigger = page.getByTestId('trigger-top')
    const tooltipContent = page.locator('#uno-tooltip-content')

    await expect(tooltipContent).not.toBeVisible()
    await trigger.hover()

    await expect(tooltipContent).toBeVisible()
    await expect(tooltipContent).toHaveText('Hello Top!')

    const triggerBox = await trigger.boundingBox()
    const tooltipBox = await tooltipContent.boundingBox()

    if (triggerBox && tooltipBox) {
      expect(tooltipBox.y + tooltipBox.height).toBeLessThanOrEqual(triggerBox.y + 5)
    }
  })

  test('shows on hover (bottom)', async ({ page }) => {
    const trigger = page.getByTestId('trigger-bottom')
    const tooltipContent = page.locator('#uno-tooltip-content')

    await trigger.hover()
    await expect(tooltipContent).toBeVisible()
    await expect(tooltipContent).toHaveText('Hello Bottom!')

    const triggerBox = await trigger.boundingBox()
    const tooltipBox = await tooltipContent.boundingBox()

    if (triggerBox && tooltipBox) {
      // Bottom: tooltip top > trigger bottom
      expect(tooltipBox.y).toBeGreaterThanOrEqual(triggerBox.y + triggerBox.height)
    }
  })

  test('shows on click (right) with light theme', async ({ page }) => {
    const trigger = page.getByTestId('trigger-click')
    const tooltipContent = page.locator('#uno-tooltip-content')

    // Initial state: hidden
    await expect(tooltipContent).not.toBeVisible()

    // Click to show
    await trigger.click()
    await expect(tooltipContent).toBeVisible()
    await expect(tooltipContent).toHaveText('Clicked (Light)!')

    // Check theme (light theme has white background)
    await expect(tooltipContent).toHaveClass(/bg-white/)
    await expect(tooltipContent).toHaveClass(/text-gray-900/)

    // Check position (right)
    const triggerBox = await trigger.boundingBox()
    const tooltipBox = await tooltipContent.boundingBox()
    if (triggerBox && tooltipBox) {
      expect(tooltipBox.x).toBeGreaterThanOrEqual(triggerBox.x + triggerBox.width)
    }

    // Click again to hide (toggle)
    await trigger.click()
    await expect(tooltipContent).not.toBeVisible()

    // Click to show again, then click outside
    await trigger.click()
    await expect(tooltipContent).toBeVisible()
    await page.mouse.click(0, 0) // Click top-left corner of page
    await expect(tooltipContent).not.toBeVisible()
  })

  test('loads HTML content from DOM target', async ({ page }) => {
    const trigger = page.getByTestId('trigger-target')
    const tooltipContent = page.locator('#uno-tooltip-content')

    await trigger.hover()
    await expect(tooltipContent).toBeVisible()

    // Check for HTML elements inside tooltip
    const boldTitle = tooltipContent.locator('strong')
    await expect(boldTitle).toHaveText('Bold Title')

    const list = tooltipContent.locator('ul')
    await expect(list).toBeVisible()
  })
})
