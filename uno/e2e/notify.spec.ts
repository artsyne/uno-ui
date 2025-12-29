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
    await page.waitForFunction(() => (window as any).Alpine?.store?.('uno')?.notify !== undefined)
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
    const btn = page.getByText('Bottom Left', { exact: true })
    await btn.click()

    const item = page.locator('[data-testid="notify-item"]').filter({ hasText: 'Bottom Left' })
    await expect(item).toBeVisible()

    const box = await item.boundingBox()
    const vp = page.viewportSize()
    if (box && vp) {
      // Bottom Left: y should be large (bottom), x should be small (left)
      expect(box.y + box.height).toBeGreaterThanOrEqual(vp.height - 20)
      expect(box.x).toBeLessThanOrEqual(100)
    }
  })

  test('stacks multiple notifications without overlap (top-right)', async ({ page }) => {
    await page.getByTestId('btn-info').click()
    await page.getByTestId('btn-success').click()
    await page.getByTestId('btn-warning').click()

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
    const btn = page.getByRole('button', { name: 'Top Left', exact: true })
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
    const btn = page.getByText('Top Left', { exact: true })
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
    const btn = page.getByText('Bottom Right', { exact: true })
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

  test('top-center position', async ({ page }) => {
    const btn = page.getByText('Top Center', { exact: true })
    await btn.click()

    const item = page.locator('[data-testid="notify-item"]').filter({ hasText: 'Top Center' })
    await expect(item).toBeVisible()

    const box = await item.boundingBox()
    const vp = page.viewportSize()
    if (box && vp) {
      // Top Center: y small, x centered
      expect(box.y).toBeLessThanOrEqual(100)
      expect(Math.abs(box.x + box.width / 2 - vp.width / 2)).toBeLessThan(20)
    }
  })

  test('center-left position', async ({ page }) => {
    const btn = page.getByText('Center Left', { exact: true })
    await btn.click()

    const item = page.locator('[data-testid="notify-item"]').filter({ hasText: 'Center Left' })
    await expect(item).toBeVisible()

    const box = await item.boundingBox()
    const vp = page.viewportSize()
    if (box && vp) {
      // Center Left: y centered, x small
      expect(Math.abs(box.y + box.height / 2 - vp.height / 2)).toBeLessThan(20)
      expect(box.x).toBeLessThanOrEqual(100)
    }
  })

  test('center position', async ({ page }) => {
    const btn = page.getByText('Center', { exact: true })
    await btn.click()

    const item = page.locator('[data-testid="notify-item"]').filter({ hasText: 'Center' })
    // The notification contains "Center" text and also "Close" (sr-only).
    // So exact match on the container text will fail if we use /^Center$/.
    // We check that it contains "Center" and verify its position.
    // We can filter by the p tag text to be sure.
    const exactItem = item.filter({ has: page.locator('p', { hasText: /^Center$/ }) })
    await expect(exactItem).toBeVisible()

    const box = await exactItem.boundingBox()
    const vp = page.viewportSize()
    if (box && vp) {
      // Center: y centered, x centered
      expect(Math.abs(box.y + box.height / 2 - vp.height / 2)).toBeLessThan(20)
      expect(Math.abs(box.x + box.width / 2 - vp.width / 2)).toBeLessThan(20)
    }
  })

  test('center-right position', async ({ page }) => {
    const btn = page.getByText('Center Right', { exact: true })
    await btn.click()

    const item = page.locator('[data-testid="notify-item"]').filter({ hasText: 'Center Right' })
    await expect(item).toBeVisible()

    const box = await item.boundingBox()
    const vp = page.viewportSize()
    if (box && vp) {
      // Center Right: y centered, x large
      expect(Math.abs(box.y + box.height / 2 - vp.height / 2)).toBeLessThan(20)
      expect(box.x + box.width).toBeGreaterThanOrEqual(vp.width - 40)
    }
  })

  test('bottom-center position', async ({ page }) => {
    const btn = page.getByText('Bottom Center', { exact: true })
    await btn.click()

    const item = page.locator('[data-testid="notify-item"]').filter({ hasText: 'Bottom Center' })
    await expect(item).toBeVisible()

    const box = await item.boundingBox()
    const vp = page.viewportSize()
    if (box && vp) {
      // Bottom Center: y large, x centered
      expect(box.y + box.height).toBeGreaterThanOrEqual(vp.height - 20)
      expect(Math.abs(box.x + box.width / 2 - vp.width / 2)).toBeLessThan(20)
    }
  })

  test('global progress bar shows and hides', async ({ page }) => {
    const btnStart = page.getByTestId('btn-start-progress')
    await btnStart.click()

    // The container wrapper has id="uno-global-progress"
    // The actual progress bar element (with progressWrap class) is inside it
    // and controlled by x-show="$store.uno.notify.progressVisible"
    const progressWrapper = page.locator('#uno-global-progress').first()
    const progressBar = progressWrapper.locator('div').first()

    await expect(progressBar).toBeVisible()

    // Check if inner bar exists
    const innerBar = progressBar.locator('div').first()
    await expect(innerBar).toBeVisible()

    // Click Stop Progress
    const btnStop = page.getByTestId('btn-stop-progress')
    await btnStop.click()

    // Should hide after fade out (stopProgress has 500ms delay + 200ms fade)
    await expect(progressBar).not.toBeVisible({ timeout: 2000 })
  })

  test('directive usage (standard)', async ({ page }) => {
    const btn = page.getByRole('button', { name: 'Directive Success' })
    await btn.click()

    const item = page
      .locator('[data-testid="notify-item"]')
      .filter({ hasText: 'Hello via Directive!' })
    await expect(item).toBeVisible()

    // Check type class (success -> green)
    const itemClass = await item.getAttribute('class')
    expect(itemClass).toContain('bg-green-100')
  })

  test('directive usage (array syntax)', async ({ page }) => {
    const btn = page.getByRole('button', { name: "Directive Array Syntax (['Msg', {conf}])" })
    await btn.click()

    const item = page.locator('[data-testid="notify-item"]').filter({ hasText: 'Mixed Config' })
    await expect(item).toBeVisible()

    // Check type class (warning -> yellow)
    const itemClass = await item.getAttribute('class')
    expect(itemClass).toContain('bg-yellow-100')
  })

  test('htmx integration (hx-get)', async ({ page }) => {
    // We rely on the real file /examples/partials/test-content.html being served
    const btn = page.getByRole('button', { name: 'Load External HTML' })
    await btn.click()

    // Wait for HTMX to process.
    const item = page.locator('[data-testid="notify-item"]').last()
    await expect(item).toBeVisible()

    // It should load the content from the file
    await expect(item).toContainText('Content loaded via HTMX successfully!', { timeout: 10000 })
  })
})
