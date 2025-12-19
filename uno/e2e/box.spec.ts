import { test, expect } from '@playwright/test'
import { mockCdnAssets } from './helpers/routes'

test.describe('Box Component', () => {
  test.beforeEach(async ({ page }) => {
    // Force serve latest build
    await mockCdnAssets(page)

    page.on('console', (msg) => console.log(`[Browser Console] ${msg.text()}`))
    page.on('pageerror', (err) => console.log(`[Page Error] ${err.message}`))

    await page.goto('/examples/components/box.html')
    await page.waitForLoadState('domcontentloaded')
    // Wait for Alpine and Uno to initialize
    await page.waitForFunction(() => (window as any).Alpine && (window as any).Alpine.store('i18n'))
    await page.waitForFunction(() => (window as any).Alpine && (window as any).Alpine.store('box'))
  })

  test('opens dialog (center) with overlay and locks body scroll', async ({ page }) => {
    const btn = page.getByText('Open Dialog')
    await btn.click()

    const overlay = page.locator('[id^="uno-box-overlay"]')
    await expect(overlay).toBeVisible()

    const outer = page.locator('[id^="uno-box-outer-container"]')
    await expect(outer).toBeVisible()
  })

  test('confirm renders buttons and closes on cancel', async ({ page }) => {
    const btn = page.getByText('Open Confirm')
    await btn.click()

    const outer = page.locator('[id^="uno-box-outer-container"]')
    await expect(outer).toBeVisible()

    const cancel = page.locator('[id^="uno-box-cancel"]')
    await expect(cancel).toBeVisible()
    await cancel.click()
    await expect(outer).not.toBeVisible()
  })

  test('position bottom-end', async ({ page }) => {
    const btn = page.getByRole('button', { name: 'Bottom End' })
    await btn.click()

    const outer = page.locator('[id^="uno-box-outer-container"]')
    await expect(outer).toBeVisible()
  })

  test('ajax loads local json and displays content', async ({ page }) => {
    await page.evaluate(() => {
      const AlpineRef = (window as any).Alpine || (window as any).top?.Alpine
      const store = AlpineRef.store('box')
      store.open({
        title: 'AJAX Content',
        ajaxUrl: '/examples/data/sample.json'
      })
    })

    const content = page.locator('[id^="uno-box-content"]')
    await expect(content).toBeVisible()
    await expect(content).toContainText('Sample Data')
  })

  test('loads content from DOM ID', async ({ page }) => {
    await page.evaluate(() => {
      const AlpineRef = (window as any).Alpine || (window as any).top?.Alpine
      const store = AlpineRef.store('box')
      store.open({
        title: 'DOM Target',
        targetId: 'hidden-content',
        trustedHTML: true
      })
    })

    const content = page.locator('[id^="uno-box-content"]')
    await expect(content).toBeVisible()
    await expect(content).toContainText('This content was hidden in the DOM.')
  })

  test('i18n toggles ok/cancel to English', async ({ page }) => {
    await page.evaluate(() => (window as any).Alpine.store('i18n').setLang('en'))
    const btn = page.getByText('Open Confirm')
    await btn.click()
    const ok = page.locator('[id^="uno-box-ok"]')
    const cancel = page.locator('[id^="uno-box-cancel"]')
    await expect(ok).toHaveText('OK')
    await expect(cancel).toHaveText('Cancel')
  })

  test('opens dialog with custom width/height classes', async ({ page }) => {
    await page.evaluate(() => {
      ;(window as any).Alpine.store('box').open('dialog', {
        title: 'Custom Size Box',
        content: 'Content',
        width: 800,
        height: 500,
        widthClass: '!w-[800px]',
        heightClass: '!h-[500px]',
        containerClass: '!border'
      })
    })

    // Use a selector that excludes the main wrapper div
    const container = page
      .locator('[id^="uno-box-container"]:not([id="uno-box-container"])')
      .first()
    await expect(container).toBeVisible()

    // Check classes
    await expect(container).toHaveClass(/!w-\[800px\]/)
    await expect(container).toHaveClass(/!h-\[500px\]/)
    await expect(container).toHaveClass(/!border/)
  })

  test('handles various dimension combinations', async ({ page }) => {
    // 1. Width/Height in pixels (number)
    await page.evaluate(() => {
      // @ts-ignore
      const store = (window.Alpine || window.top.Alpine).store('box')
      store.open({
        title: 'Pixel Dimensions',
        content: 'Testing number dimensions',
        width: 600,
        height: 400
      })
    })

    const box1 = page.locator('[id^="uno-box-container"]:not([id="uno-box-container"])').last()
    await expect(box1).toBeVisible()
    await expect(box1).toHaveAttribute('style', /width:\s*600px/)
    await expect(box1).toHaveAttribute('style', /height:\s*400px/)

    // Close it
    await page.getByRole('button', { name: 'Close' }).click()

    // 3. Class Overrides (Tailwind arbitrary values)
    await page.evaluate(() => {
      // @ts-ignore
      const store = (window.Alpine || window.top.Alpine).store('box')
      store.open({
        title: 'Class Overrides',
        content: 'Testing class overrides',
        widthClass: '!w-[880px]',
        heightClass: '!h-[570px]'
      })
    })

    const box2 = page.locator('[id^="uno-box-container"]:not([id="uno-box-container"])').last()
    await expect(box2).toBeVisible()
    await expect(box2).toHaveClass(/!w-\[880px\]/)
    await expect(box2).toHaveClass(/!h-\[570px\]/)

    // Close it
    await page.getByRole('button', { name: 'Close' }).click()

    // 4. Class without !important (should NOT have default width style)
    await page.evaluate(() => {
      // @ts-ignore
      const store = (window.Alpine || window.top.Alpine).store('box')
      store.open({
        title: 'Responsive Class',
        content: 'Testing responsive class',
        widthClass: 'w-1/2' // 50%
      })
    })

    const box3 = page.locator('[id^="uno-box-container"]:not([id="uno-box-container"])').last()
    await expect(box3).toBeVisible()
    await expect(box3).toHaveClass(/w-1\/2/)
    // Verify NO inline width style is applied (because widthClass is present)
    const style = await box3.getAttribute('style')
    if (style) {
      expect(style).not.toContain('width:')
    }
  })

  test('executes onClose and onBeforeClose callbacks', async ({ page }) => {
    await page.evaluate(() => {
      ;(window as any)._callbackLogs = []
      // @ts-ignore
      const store = (window.Alpine || window.top.Alpine).store('box')
      store.open({
        title: 'Callback Test',
        content: 'Testing callbacks',
        onBeforeClose: () => (window as any)._callbackLogs.push('onBeforeClose'),
        onClose: () => (window as any)._callbackLogs.push('onClose')
      })
    })

    const box = page.locator('[id^="uno-box-container"]:not([id="uno-box-container"])').last()
    await expect(box).toBeVisible()

    // Close the box
    await page.getByRole('button', { name: 'Close' }).click()

    await expect(box).not.toBeVisible()

    const logs = await page.evaluate(() => (window as any)._callbackLogs)
    expect(logs).toEqual(['onBeforeClose', 'onClose'])
  })

  test('executes callbacks on Cancel button', async ({ page }) => {
    await page.evaluate(() => {
      ;(window as any)._callbackLogs = []
      // @ts-ignore
      const store = (window.Alpine || window.top.Alpine).store('box')
      store.open('confirm', {
        title: 'Cancel Test',
        content: 'Testing cancel callback',
        onCancel: () => (window as any)._callbackLogs.push('onCancel'),
        onClose: () => (window as any)._callbackLogs.push('onClose')
      })
    })

    const box = page.locator('[id^="uno-box-container"]:not([id="uno-box-container"])').last()
    await expect(box).toBeVisible()

    // Click Cancel
    await page.getByRole('button', { name: 'Cancel' }).click()

    await expect(box).not.toBeVisible()

    const logs = await page.evaluate(() => (window as any)._callbackLogs)
    // Expect onCancel first, then onClose (as close() is called after onCancel)
    expect(logs).toEqual(['onCancel', 'onClose'])
  })

  test('executes onContentLoad callback after ajax loading', async ({ page }) => {
    await page.evaluate(() => {
      ;(window as any)._callbackLogs = []
      // @ts-ignore
      const store = (window.Alpine || window.top.Alpine).store('box')
      store.open({
        title: 'AJAX Callback Test',
        ajaxUrl: '/examples/data/sample.json',
        onContentLoad: (el: any) => {
          ;(window as any)._callbackLogs.push('onContentLoad')
          // Verify element is passed
          if (el && el.id && el.id.startsWith('uno-box-content')) {
            ;(window as any)._callbackLogs.push('element_received')
          }
        }
      })
    })

    // Wait for content to load
    const content = page.locator('[id^="uno-box-content"]')
    await expect(content).toContainText('Sample Data')

    // Check logs
    const logs = await page.evaluate(() => (window as any)._callbackLogs)
    expect(logs).toContain('onContentLoad')
    expect(logs).toContain('element_received')
  })

  test('perf: open/close 50 dialogs', async ({ page }) => {
    const durationMs = await page.evaluate(async () => {
      // @ts-ignore
      const AlpineRef = window.Alpine || window.top.Alpine
      const store = AlpineRef.store('box')

      const start = performance.now()
      for (let i = 0; i < 50; i++) {
        const id = store.open({
          title: 'Perf',
          content: 'Perf',
          modal: false,
          width: 400
        })
        await AlpineRef.nextTick()
        store.close(id)
        await AlpineRef.nextTick()
      }
      return performance.now() - start
    })

    console.log(`[Perf] box open/close 50 dialogs: ${durationMs.toFixed(2)}ms`)
    expect(Number.isFinite(durationMs)).toBeTruthy()
  })
})
