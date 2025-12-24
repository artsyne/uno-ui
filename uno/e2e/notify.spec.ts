import { test, expect } from '@playwright/test'
import { mockCdnAssets } from './helpers/routes'

test.describe('Notify Component', () => {
  test.beforeEach(async ({ page }) => {
    // Force serve latest build
    await mockCdnAssets(page)

    page.on('console', (msg) => console.log(`[Browser Console] ${msg.text()}`))
    page.on('pageerror', (err) => console.log(`[Page Error] ${err.message}`))
    await page.goto('/examples/components/notify.html')

    await page.waitForFunction(() => (window as any).Alpine !== undefined)
    await page.waitForFunction(() => (window as any).Uno?.initialized === true)
    await page.evaluate(() => {
      const w = window as any
      if (
        w.Alpine &&
        typeof w.Alpine.start === 'function' &&
        !w.__uno_e2e_alpine_started &&
        !(document.body as any)._x_dataStack
      ) {
        try {
          w.Alpine.start()
        } catch {}
        w.__uno_e2e_alpine_started = true
      }
    })
    await page.waitForFunction(() => !!(document.body as any)._x_dataStack)
    await page.waitForFunction(() => (window as any).Alpine?.store?.('notify') !== undefined)
  })

  test('shows info notification (top-right) and auto dismiss', async ({ page }) => {
    const btn = page.getByTestId('btn-info')
    await btn.click()

    // Use a locator that matches the text to ensure we get the right one
    const item = page.locator('[data-testid="notify-item"]').filter({ hasText: 'Info message' })
    await expect(item).toBeVisible()
    await expect(item).toHaveText(/Info message/)

    const box = await item.boundingBox()
    const vp = page.viewportSize()
    if (box && vp) {
      // Top Right: y should be small (top), x should be large (right)
      expect(box.y).toBeLessThanOrEqual(100)
      expect(box.x + box.width).toBeGreaterThanOrEqual(vp.width - 40)
      expect(box.width).toBeGreaterThanOrEqual(80)
      expect(box.width).toBeLessThanOrEqual(vp.width)
    }

    // Auto dismiss within ~3s (timeout is 3000ms)
    // We wait a bit longer to ensure animation completes
    await expect(item).not.toBeVisible({ timeout: 5000 })
  })

  test('sticky notification requires click to close', async ({ page }) => {
    const btn = page.getByTestId('btn-sticky')
    await btn.click()

    const item = page.locator('[data-testid="notify-item"]').filter({ hasText: 'Sticky message' })
    await expect(item).toBeVisible()

    // Wait 2s, should still be there (timeout 0)
    await page.waitForTimeout(2000)
    await expect(item).toBeVisible()

    // Click close button
    // The close button is inside the item
    const closeBtn = item.locator('button')
    await closeBtn.click()
    await expect(item).not.toBeVisible()
  })

  test('bottom-left position', async ({ page }) => {
    const btn = page.getByText('Bottom Left')
    await btn.click()

    const item = page.locator('[data-testid="notify-item"]').filter({ hasText: 'Bottom Left' })
    await expect(item).toBeVisible()

    const box = await item.boundingBox()
    const vp = page.viewportSize()
    if (box && vp) {
      // Bottom Left: y should be large (bottom), x should be large (bottom), x should be small (left)
      expect(box.y + box.height).toBeGreaterThanOrEqual(vp.height - 20)
      expect(box.x).toBeLessThanOrEqual(100)
    }
  })

  test('center position', async ({ page }) => {
    const btn = page.getByText('Center', { exact: true })
    await btn.click()

    const item = page.locator('[data-testid="notify-item"]').filter({ hasText: 'Center (Toast)' })
    await expect(item).toBeVisible()

    const box = await item.boundingBox()
    const vp = page.viewportSize()
    if (box && vp) {
      // Center: roughly centered
      const centerX = box.x + box.width / 2
      const centerY = box.y + box.height / 2
      const viewportCenterX = vp.width / 2
      const viewportCenterY = vp.height / 2

      expect(Math.abs(centerX - viewportCenterX)).toBeLessThan(20)
      expect(Math.abs(centerY - viewportCenterY)).toBeLessThan(20)
    }
  })

  test('stacks multiple notifications without overlap (top-right)', async ({ page }) => {
    await page.getByTestId('btn-info').click()
    await page.getByText('Success').click()
    await page.getByText('Warning').click()

    const info = page.locator('[data-testid="notify-item"]').filter({ hasText: 'Info message' })
    const success = page
      .locator('[data-testid="notify-item"]')
      .filter({ hasText: 'Success message' })
    const warning = page
      .locator('[data-testid="notify-item"]')
      .filter({ hasText: 'Warning message' })

    await expect(info).toBeVisible()
    await expect(success).toBeVisible()
    await expect(warning).toBeVisible()

    await page.waitForTimeout(400)

    const boxes = await Promise.all([
      info.boundingBox(),
      success.boundingBox(),
      warning.boundingBox()
    ])
    const filtered = boxes.filter(Boolean) as NonNullable<(typeof boxes)[number]>[]
    expect(filtered.length).toBe(3)

    const sorted = [...filtered].sort((a, b) => a.y - b.y)
    expect(sorted[1].y).toBeGreaterThanOrEqual(sorted[0].y + sorted[0].height - 1)
    expect(sorted[2].y).toBeGreaterThanOrEqual(sorted[1].y + sorted[1].height - 1)
  })

  test('stacks multiple notifications without overlap (top-left)', async ({ page }) => {
    // We need to trigger notifications in top-left.
    // Since the example buttons are hardcoded to specific positions or defaults,
    // we might need to use the console or modify the DOM to trigger top-left specifically for this test
    // OR, use the "Top Left" button if it exists and allows multiple clicks.
    // The "Top Left" button in example:
    // <button ... data-uno-notify ... data-uno-position="top-left">Top Left</button>

    const btn = page.getByRole('button', { name: 'Top Left' })
    await btn.click()
    await btn.click()
    await btn.click()

    const items = page.locator('[data-testid="notify-item"]').filter({ hasText: 'Top Left' })
    await expect(items).toHaveCount(3)

    const boxes = await Promise.all([
      items.nth(0).boundingBox(),
      items.nth(1).boundingBox(),
      items.nth(2).boundingBox()
    ])
    const filtered = boxes.filter(Boolean) as NonNullable<(typeof boxes)[number]>[]
    expect(filtered.length).toBe(3)

    const sorted = [...filtered].sort((a, b) => a.y - b.y)
    // Ensure they are stacked vertically (y increases)
    expect(sorted[1].y).toBeGreaterThanOrEqual(sorted[0].y + sorted[0].height - 1)
    expect(sorted[2].y).toBeGreaterThanOrEqual(sorted[1].y + sorted[1].height - 1)

    // Ensure they are left aligned (roughly)
    expect(Math.abs(sorted[0].x - sorted[1].x)).toBeLessThan(5)
    expect(Math.abs(sorted[0].x - sorted[2].x)).toBeLessThan(5)
  })

  test('avoids squeezing content (adaptive width)', async ({ page }) => {
    const btn = page.getByTestId('btn-sticky')
    await btn.click()

    const item = page.locator('[data-testid="notify-item"]').filter({ hasText: 'Sticky message' })
    await expect(item).toBeVisible()

    const box = await item.boundingBox()
    if (box) {
      // "Sticky message (click to close)" should be reasonably wide, not a vertical strip
      // Text is approx 30 chars. 30 chars * ~8px = ~240px.
      // Squeezed width would be much smaller (e.g. < 100px)
      expect(box.width).toBeGreaterThan(150)

      // Also ensure it's not full width
      const vp = page.viewportSize()
      if (vp) {
        expect(box.width).toBeLessThan(vp.width - 40)
      }
    }

    // Cleanup
    await item.locator('button').click()
  })

  test('dark mode styles applied', async ({ page }) => {
    // Enable dark mode
    await page.evaluate(() => document.documentElement.classList.add('dark'))

    const btn = page.getByTestId('btn-info')
    await btn.click()

    const item = page.locator('[data-testid="notify-item"]').filter({ hasText: 'Info message' })
    await expect(item).toBeVisible()

    // Check for dark mode specific classes on the item container (bg)
    const itemClass = await item.getAttribute('class')
    expect(itemClass).toContain('dark:bg-blue-950')

    // Check for dark mode text color on the text element
    const textEl = item.locator('p') // The text is in a <p> tag
    const textClass = await textEl.getAttribute('class')
    expect(textClass).toContain('dark:text-slate-100')

    // Verify background color is different in dark mode (via computed style if needed, but class check is faster for now)
    // The type class is applied dynamically.
    // Let's check if the computed background color is not white.

    const bgColor = await item.evaluate((el) => {
      return window.getComputedStyle(el).backgroundColor
    })

    // bg-blue-950/40 is roughly dark blue. White is rgb(255, 255, 255).
    expect(bgColor).not.toBe('rgb(255, 255, 255)')

    // Cleanup
    await page.evaluate(() => document.documentElement.classList.remove('dark'))
  })

  test('top-center position', async ({ page }) => {
    const btn = page.getByText('Top Center')
    await btn.click()

    const item = page.locator('[data-testid="notify-item"]').filter({ hasText: 'Top Center' })
    await expect(item).toBeVisible()

    const box = await item.boundingBox()
    const vp = page.viewportSize()
    if (box && vp) {
      // Top Center: y should be small (top), x should be centered
      expect(box.y).toBeLessThanOrEqual(100)

      const centerX = box.x + box.width / 2
      const viewportCenterX = vp.width / 2
      expect(Math.abs(centerX - viewportCenterX)).toBeLessThan(20)
    }
  })

  test('bottom-center position', async ({ page }) => {
    const btn = page.getByText('Bottom Center')
    await btn.click()

    const item = page.locator('[data-testid="notify-item"]').filter({ hasText: 'Bottom Center' })
    await expect(item).toBeVisible()

    const box = await item.boundingBox()
    const vp = page.viewportSize()
    if (box && vp) {
      // Bottom Center: y should be large (bottom), x should be centered
      expect(box.y + box.height).toBeGreaterThanOrEqual(vp.height - 20)

      const centerX = box.x + box.width / 2
      const viewportCenterX = vp.width / 2
      expect(Math.abs(centerX - viewportCenterX)).toBeLessThan(20)
    }
  })

  test('top-right position (explicit)', async ({ page }) => {
    const btn = page.getByText('Top Right', { exact: true })
    await btn.click()

    const item = page.locator('[data-testid="notify-item"]').filter({ hasText: 'Top Right' })
    await expect(item).toBeVisible()

    const box = await item.boundingBox()
    const vp = page.viewportSize()
    if (box && vp) {
      // Top Right: y small, x large
      expect(box.y).toBeLessThanOrEqual(100)
      expect(box.x + box.width).toBeGreaterThanOrEqual(vp.width - 40)
    }
  })

  test('top-left position', async ({ page }) => {
    const btn = page.getByText('Top Left')
    await btn.click()

    const item = page.locator('[data-testid="notify-item"]').filter({ hasText: 'Top Left' })
    await expect(item).toBeVisible()

    const box = await item.boundingBox()
    const vp = page.viewportSize()
    if (box && vp) {
      // Top Left: y small, x small
      expect(box.y).toBeLessThanOrEqual(100)
      expect(box.x).toBeLessThanOrEqual(100)
    }
  })

  test('bottom-right position', async ({ page }) => {
    const btn = page.getByText('Bottom Right')
    await btn.click()

    const item = page.locator('[data-testid="notify-item"]').filter({ hasText: 'Bottom Right' })
    await expect(item).toBeVisible()

    const box = await item.boundingBox()
    const vp = page.viewportSize()
    if (box && vp) {
      // Bottom Right: y large, x large
      expect(box.y + box.height).toBeGreaterThanOrEqual(vp.height - 20)
      expect(box.x + box.width).toBeGreaterThanOrEqual(vp.width - 40)
    }
  })

  test('progress bar shows and hides', async ({ page }) => {
    const btn = page.getByTestId('btn-progress')
    await btn.click()

    const progressBar = page.locator('#uno-progress-bar')
    await expect(progressBar).toBeVisible()

    // Check if it has the animation class
    const innerBar = progressBar.locator('.uno-progress-anim')
    await expect(innerBar).toBeVisible()

    // Should hide after 3s (timeout in the button handler)
    // We wait 3.5s to be safe, increasing to 5s for slower environments
    await expect(progressBar).not.toBeVisible({ timeout: 6000 })
  })

  test('loading indicator shows and hides', async ({ page }) => {
    const btn = page.getByRole('button', { name: 'Show Loading (2s)' })
    await btn.click()

    await page.waitForFunction(() => (window as any).Alpine?.store?.('notify')?.activeRequests > 0)

    const overlay = page.locator('#uno-loading-indicator')
    await expect(overlay).toBeVisible()

    const spinner = overlay.locator('svg')
    await expect(spinner).toBeVisible()

    await page.waitForFunction(
      () => (window as any).Alpine?.store?.('notify')?.activeRequests === 0
    )
    await expect(overlay).not.toBeVisible({ timeout: 4000 })
  })

  test('dark mode: avoids low-contrast class pairing (bg-white + text-gray-200)', async ({
    page
  }) => {
    const html = page.locator('html')
    if (!(await html.evaluate((el) => el.classList.contains('dark')))) {
      await page.locator('[data-theme-toggle]').click({ force: true })
    }
    await expect(html).toHaveClass(/dark/)

    await page.getByTestId('btn-info').click()

    const item = page.locator('[data-testid="notify-item"]').filter({ hasText: 'Info message' })
    await expect(item).toBeVisible()

    const hasBadPairing = await item.evaluate((root) => {
      const nodes = [root, ...Array.from(root.querySelectorAll('*'))]
      return nodes.some(
        (el) => el.classList.contains('bg-white') && el.classList.contains('text-gray-200')
      )
    })
    expect(hasBadPairing).toBe(false)
  })
})
