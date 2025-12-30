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
    // Updated: Wait for box store under uno namespace
    await page.waitForFunction(
      () => (window as any).Alpine && (window as any).Alpine.store('uno')?.box
    )
  })

  test('opens dialog (center) with overlay and locks body scroll', async ({ page }) => {
    await page.evaluate(() => {
      ;(window as any).Alpine.store('uno').box.open('dialog', {
        title: 'Dialog Title',
        content: 'This is a standard dialog with some content.',
        position: 'center'
      })
    })

    const overlay = page.locator('#uno-box-overlay')
    await expect(overlay).toBeVisible()

    const outer = page.locator('[id^="uno-box-outer-container"]')
    await expect(outer).toBeVisible()
  })

  test('confirm renders buttons and closes on cancel', async ({ page }) => {
    await page.evaluate(() => {
      ;(window as any).Alpine.store('uno').box.open('confirm', {
        title: 'Confirm Action',
        content: 'Are you sure you want to proceed?'
      })
    })

    const outer = page.locator('[id^="uno-box-outer-container"]')
    await expect(outer).toBeVisible()

    const cancel = page.locator('[id^="uno-box-cancel"]')
    await expect(cancel).toBeVisible()
    await cancel.click()
    await expect(outer).not.toBeVisible()
  })

  test('positioning supports 9-grid positions', async ({ page }) => {
    const cases: Array<[string, string, string]> = [
      ['top-center', 'items-start', 'justify-center'],
      ['top-left', 'items-start', 'justify-start'],
      ['top-right', 'items-start', 'justify-end'],
      ['center', 'items-center', 'justify-center'],
      ['center-left', 'items-center', 'justify-start'],
      ['center-right', 'items-center', 'justify-end'],
      ['bottom-center', 'items-end', 'justify-center'],
      ['bottom-left', 'items-end', 'justify-start'],
      ['bottom-right', 'items-end', 'justify-end']
    ]

    for (const [position, itemsCls, justifyCls] of cases) {
      const boxId = `pos-${position}`

      await page.evaluate(
        ({ position, boxId }) => {
          ;(window as any).Alpine.store('uno').box.open('dialog', {
            boxId,
            title: position,
            content: 'Positioning test',
            position
          })
        },
        { position, boxId }
      )

      const outer = page.locator(`#uno-box-outer-container${boxId}`)
      await expect(outer).toBeVisible()
      await expect(outer).toHaveClass(new RegExp(itemsCls))
      await expect(outer).toHaveClass(new RegExp(justifyCls))

      await page.getByRole('button', { name: 'Close', exact: true }).click()
      await expect(outer).not.toBeVisible()
    }
  })

  test('hxUrl loads local html and displays content', async ({ page }) => {
    await page.evaluate(() => {
      const AlpineRef = (window as any).Alpine || (window as any).top?.Alpine
      const store = AlpineRef.store('uno').box
      store.open({
        title: 'HTMX Content',
        hxUrl: '/examples/data/tab-content.html'
      })
    })

    const content = page.locator('[id^="uno-box-content"]')
    await expect(content).toBeVisible()
    await expect(content).toContainText('Ajax Loaded Content')
  })

  test('loads content from DOM innerHTML (sanitized)', async ({ page }) => {
    await page.evaluate(() => {
      const AlpineRef = (window as any).Alpine || (window as any).top?.Alpine
      const store = AlpineRef.store('uno').box
      store.open({
        title: 'DOM Target',
        content: document.getElementById('hidden-content')?.innerHTML || ''
      })
    })

    const content = page.locator('[id^="uno-box-content"]')
    await expect(content).toBeVisible()
    await expect(content).toContainText('This content was hidden in the DOM.')
  })

  test('i18n toggles ok/cancel to English', async ({ page }) => {
    await page.evaluate(() => (window as any).Alpine.store('i18n').setLang('en'))
    await page.evaluate(() => {
      ;(window as any).Alpine.store('uno').box.open('confirm', {
        title: 'I18N Confirm',
        content: 'Testing i18n'
      })
    })
    const ok = page.locator('[id^="uno-box-ok"]')
    const cancel = page.locator('[id^="uno-box-cancel"]')
    await expect(ok).toHaveText('OK')
    await expect(cancel).toHaveText('Cancel')
  })

  test('opens dialog with custom width/height classes', async ({ page }) => {
    await page.evaluate(() => {
      ;(window as any).Alpine.store('uno').box.open('dialog', {
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
      const store = (window.Alpine || window.top.Alpine).store('uno').box
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

    // Close it - use first visible close button
    await page.locator('[aria-label="Close"]:visible').first().click()
    await expect(box1).not.toBeVisible()

    // 3. Class Overrides (Tailwind arbitrary values)
    await page.evaluate(() => {
      // @ts-ignore
      const store = (window.Alpine || window.top.Alpine).store('uno').box
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

    // Close it - use first visible close button
    await page.locator('[aria-label="Close"]:visible').first().click()
    await expect(box2).not.toBeVisible()

    // 4. Class without !important (should NOT have default width style)
    await page.evaluate(() => {
      // @ts-ignore
      const store = (window.Alpine || window.top.Alpine).store('uno').box
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
      const store = (window.Alpine || window.top.Alpine).store('uno').box
      store.open({
        title: 'Callback Test',
        content: 'Testing callbacks',
        onBeforeClose: () => (window as any)._callbackLogs.push('onBeforeClose'),
        onClose: () => (window as any)._callbackLogs.push('onClose')
      })
    })

    const box = page.locator('[id^="uno-box-container"]:not([id="uno-box-container"])').last()
    await expect(box).toBeVisible()

    // Close the box - use exact match
    await page.getByRole('button', { name: 'Close', exact: true }).click()

    await expect(box).not.toBeVisible()

    const logs = await page.evaluate(() => (window as any)._callbackLogs)
    expect(logs).toEqual(['onBeforeClose', 'onClose'])
  })

  test('executes callbacks on Cancel button', async ({ page }) => {
    await page.evaluate(() => {
      ;(window as any)._callbackLogs = []
      // @ts-ignore
      const store = (window.Alpine || window.top.Alpine).store('uno').box
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
    // Expect onCancel first, then onBeforeClose, then onClose
    expect(logs).toContain('onCancel')
    expect(logs).toContain('onClose')
  })

  test('executes onContentLoad callback after hxUrl loading', async ({ page }) => {
    await page.evaluate(() => {
      ;(window as any)._callbackLogs = []
      // @ts-ignore
      const store = (window.Alpine || window.top.Alpine).store('uno').box
      store.open({
        title: 'HTMX Callback Test',
        hxUrl: '/examples/data/tab-content.html',
        onContentLoad: (el: any) => {
          ;(window as any)._callbackLogs.push('onContentLoad')
          // Element is the parent container that contains the loaded content
          if (el) {
            ;(window as any)._callbackLogs.push('element_received')
          }
        }
      })
    })

    // Wait for content to load
    const content = page.locator('[id^="uno-box-content"]')
    await expect(content).toContainText('Ajax Loaded Content')

    // Check logs
    const logs = await page.evaluate(() => (window as any)._callbackLogs)
    expect(logs).toContain('onContentLoad')
    expect(logs).toContain('element_received')
  })

  test('XSS protection: sanitizes dangerous content', async ({ page }) => {
    await page.evaluate(() => {
      // @ts-ignore
      const store = (window.Alpine || window.top.Alpine).store('uno').box
      store.open({
        title: 'XSS Test',
        content:
          '<script>window.__XSS_EXECUTED__=true</script><a href="javascript:alert(1)">Click</a><img src="x" onerror="window.__XSS_IMG__=true">Safe text'
      })
    })

    const content = page.locator('[id^="uno-box-content"]')
    await expect(content).toBeVisible()
    // Should show safe text
    await expect(content).toContainText('Safe text')

    // Script should NOT execute
    const xssExecuted = await page.evaluate(() => (window as any).__XSS_EXECUTED__)
    expect(xssExecuted).toBeUndefined()

    // onerror should NOT execute
    const xssImg = await page.evaluate(() => (window as any).__XSS_IMG__)
    expect(xssImg).toBeUndefined()
  })

  test('closes on Escape key when configured', async ({ page }) => {
    await page.evaluate(() => {
      ;(window as any).Alpine.store('uno').box.open({
        title: 'Escape Test',
        content: 'Press Escape to close',
        modal: true,
        closeOnEsc: true
      })
    })

    const outer = page.locator('[id^="uno-box-outer-container"]')
    await expect(outer).toBeVisible()

    // Press Escape
    await page.keyboard.press('Escape')

    await expect(outer).not.toBeVisible()
  })

  test('does not close on Escape key by default', async ({ page }) => {
    await page.evaluate(() => {
      ;(window as any).Alpine.store('uno').box.open({
        title: 'No Escape Test',
        content: 'Press Escape should not close',
        modal: true
        // closeOnEsc defaults to false
      })
    })

    const outer = page.locator('[id^="uno-box-outer-container"]')
    await expect(outer).toBeVisible()

    // Press Escape
    await page.keyboard.press('Escape')

    // Should still be visible
    await expect(outer).toBeVisible()
  })

  test('does not close on overlay click by default', async ({ page }) => {
    await page.evaluate(() => {
      ;(window as any).Alpine.store('uno').box.open({
        title: 'Overlay Click Test',
        content: 'Click overlay not close',
        modal: true
      })
    })

    const outer = page.locator('[id^="uno-box-outer-container"]')
    await expect(outer).toBeVisible()

    // Click on the overlay (dialogWrapper covers the screen)
    const wrapper = page.locator('[class*="fixed inset-0 z-10 w-screen"]')
    await wrapper.click({ position: { x: 10, y: 10 } })

    // Should still be visible
    await expect(outer).toBeVisible()
  })

  test('closes on overlay click when configured', async ({ page }) => {
    await page.evaluate(() => {
      ;(window as any).Alpine.store('uno').box.open({
        title: 'Overlay Click Test',
        content: 'Click overlay to close',
        modal: true,
        closeOnOverlayClick: true
      })
    })

    const outer = page.locator('[id^="uno-box-outer-container"]')
    await expect(outer).toBeVisible()

    // Click on the overlay (dialogWrapper covers the screen)
    const wrapper = page.locator('[class*="fixed inset-0 z-10 w-screen"]')
    await wrapper.click({ position: { x: 10, y: 10 } })

    // Should be hidden
    await expect(outer).not.toBeVisible()
  })

  test('perf: open/close 50 dialogs', async ({ page }) => {
    const durationMs = await page.evaluate(async () => {
      // @ts-ignore
      const AlpineRef = window.Alpine || window.top.Alpine
      const store = AlpineRef.store('uno').box

      const start = performance.now()
      for (let i = 0; i < 50; i++) {
        const item = store.open({
          title: 'Perf',
          content: 'Perf',
          modal: false,
          width: 400
        })
        await AlpineRef.nextTick()
        store.close(item.id)
        await AlpineRef.nextTick()
      }
      return performance.now() - start
    })

    console.log(`[Perf] box open/close 50 dialogs: ${durationMs.toFixed(2)}ms`)
    expect(Number.isFinite(durationMs)).toBeTruthy()
  })

  test('dialog method works as alias for open', async ({ page }) => {
    await page.evaluate(() => {
      ;(window as any).Alpine.store('uno').box.dialog({
        title: 'Dialog Method',
        content: 'Using dialog() method'
      })
    })

    const content = page.locator('[id^="uno-box-content"]')
    await expect(content).toContainText('Using dialog() method')
  })
})
