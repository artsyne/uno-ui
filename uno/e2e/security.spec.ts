import { test, expect } from '@playwright/test'
import { mockCdnAssets } from './helpers/routes'

test.describe('Uno UI Security & Features', () => {
  test.beforeEach(async ({ page }) => {
    // Force serve latest build
    await mockCdnAssets(page)

    // Log console
    page.on('console', (msg) => console.log(`[Browser Console] ${msg.text()}`))

    // Use box.html as a base since we know it works and loads Uno
    await page.goto('/examples/components/box.html')

    // Wait for Alpine and Uno
    await page.waitForFunction(() => (window as any).Alpine && (window as any).Uno)
  })

  test('Unit: sanitizeHTML should remove img tags', async ({ page }) => {
    const result = await page.evaluate(() => {
      const html = '<img src=x onerror=alert(1)><b>Safe</b>'
      // @ts-ignore
      return window.Uno.dom.sanitize(html)
    })
    expect(result).toBe('<b>Safe</b>')
  })

  test('Tab content should be sanitized (XSS prevention)', async ({ page }) => {
    // Inject a tab component dynamically to test
    await page.evaluate(() => {
      const div = document.createElement('div')
      div.innerHTML = `
            <div x-data="unoTabs('tab1')">
                <div role="tablist">
                    <button x-bind="trigger('tab1')">Tab 1</button>
                    <button id="malicious-tab" x-bind="trigger('tab2', { url: '/mock/malicious' })">Tab 2</button>
                </div>
                <div x-bind="panel('tab1')">Content 1</div>
                <div id="panel-tab2" x-bind="panel('tab2')"></div>
            </div>
        `
      document.body.appendChild(div)
      // Re-init Alpine tree
      ;(window as any).Alpine.initTree(div)
      // Process HTMX for the new content
      if ((window as any).htmx) (window as any).htmx.process(div)
    })

    // Intercept network request to return malicious content
    await page.route(/\/mock\/malicious/, (route) => {
      route.fulfill({
        status: 200,
        contentType: 'text/html',
        body: `
              <img src=x onerror=window.xssImg=true>
              <b>Safe</b>
            `
      })
    })

    // Click the tab
    await page.click('#malicious-tab')

    // Wait for content to load
    await expect(page.locator('#panel-tab2')).toContainText('Safe')

    // Check if img tag with onerror is sanitized (removed or stripped)
    const xssImg = await page.evaluate(() => (window as any).xssImg)
    expect(xssImg).toBeUndefined()

    // Check if tags are removed
    const imgCount = await page.locator('#panel-tab2 img').count()
    expect(imgCount).toBe(0)
  })

  test('Box content should be sanitized', async ({ page }) => {
    await page.evaluate(() => {
      ;(window as any).Alpine.store('box').open('dialog', {
        boxId: 'test-box',
        title: 'Test',
        content: '<img src=x onerror=window.boxXss=true><b>SafeBox</b>',
        trustedHTML: false // Explicitly false or default
      })
    })

    await expect(page.locator('#uno-box-content-test-box')).toContainText('SafeBox')

    const xss = await page.evaluate(() => (window as any).boxXss)
    expect(xss).toBeUndefined()

    const imgCount = await page.locator('#uno-box-content-test-box img').count()
    expect(imgCount).toBe(0)
  })

  test('Requests should have HTMX headers', async ({ page }) => {
    let headers: Record<string, string> = {}
    let captured = false

    // Use regex to match any API request to test-headers
    await page.route(/\/api\/test-headers/, (route) => {
      headers = route.request().headers()
      captured = true
      route.fulfill({ status: 200, body: 'ok' })
    })

    await page.evaluate(async () => {
      // Use relative path which resolves to page origin
      await (window as any).Alpine.store('box').open('dialog', {
        hxUrl: '/api/test-headers'
      })
    })

    // Wait for capture
    await expect.poll(() => captured).toBe(true)

    console.log('Captured Headers:', headers)

    expect(headers['hx-request']).toBe('true')
    expect(headers['hx-current-url']).toBeTruthy()
  })

  test('CSP Nonce should be injected', async ({ page }) => {
    await page.evaluate(() => {
      ;(window as any).Uno.options.csp = { nonce: 'test-nonce-123' }
      // Remove existing container to force recreation
      const existing = document.getElementById('uno-notify-container')
      if (existing)
        existing.remove()
        // Trigger a container injection
      ;(window as any).Alpine.store('notify').show('Test')
    })

    const notifyContainer = page.locator('#uno-notify-container')
    await expect(notifyContainer).toHaveAttribute('nonce', 'test-nonce-123')
  })
})
