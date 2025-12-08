import { test, expect } from '@playwright/test'

test.describe('Notify Component', () => {
  test.beforeEach(async ({ page }) => {
    // Force serve latest build
    await page.route('**/dist/uno.umd.js', async (route) => {
      await route.fulfill({ path: 'dist/uno.umd.js' })
    })

    page.on('console', (msg) => console.log(`[Browser Console] ${msg.text()}`))
    page.on('pageerror', (err) => console.log(`[Page Error] ${err.message}`))
    await page.goto('/examples/components/notify.html')

    // Wait for Alpine to be ready
    await page.waitForFunction(() => (window as any).Alpine !== undefined)
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
      expect(box.x + box.width).toBeGreaterThanOrEqual(vp.width - 20)
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
      expect(box.x + box.width).toBeGreaterThanOrEqual(vp.width - 20)
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
      expect(box.x + box.width).toBeGreaterThanOrEqual(vp.width - 20)
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
    // We wait 3.5s to be safe
    await expect(progressBar).not.toBeVisible({ timeout: 4000 })
  })
})
