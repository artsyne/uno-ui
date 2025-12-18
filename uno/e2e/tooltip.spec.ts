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

  test('sanitizes dangerous HTML from DOM target', async ({ page }) => {
    const trigger = page.getByTestId('trigger-target')
    const tooltipContent = page.locator('#uno-tooltip-content')

    // Reset flags (img onerror might have fired on page load)
    await page.evaluate(() => {
      const win = window as any
      win.__UNO_XSS_LINK__ = undefined
      win.__UNO_XSS_IMG__ = undefined
    })

    await trigger.hover()
    await expect(tooltipContent).toBeVisible()

    const unsafeLink = tooltipContent.getByText('Unsafe Link')
    await expect(unsafeLink).toBeVisible()
    await expect(unsafeLink).not.toHaveAttribute('href', /.+/)

    await expect(tooltipContent.locator('img')).toHaveCount(0)

    const flags = await page.evaluate(() => {
      const win = window as any
      return {
        link: win.__UNO_XSS_LINK__,
        img: win.__UNO_XSS_IMG__
      }
    })
    expect(flags.link).toBeUndefined()
    expect(flags.img).toBeUndefined()
  })

  test('perf: show/hide 100 tooltips', async ({ page }) => {
    const durationMs = await page.evaluate(async () => {
      // @ts-ignore
      const AlpineRef = window.Alpine || window.top.Alpine
      const store = AlpineRef.store('tooltip')
      const trigger = document.querySelector('[data-testid="trigger-top"]')
      if (!trigger) return -1

      const start = performance.now()
      for (let i = 0; i < 100; i++) {
        await store.show(trigger, 'Perf', { placement: 'top' })
        store.hide()
        await AlpineRef.nextTick()
      }
      return performance.now() - start
    })

    console.log(`[Perf] tooltip show/hide 100 times: ${durationMs.toFixed(2)}ms`)
    expect(durationMs).toBeGreaterThan(0)
  })
})
