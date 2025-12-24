import { test, expect } from '@playwright/test'
import { mockCdnAssets } from './helpers/routes'

test.describe('Tooltip Component', () => {
  test.beforeEach(async ({ page }) => {
    page.on('console', (msg) => console.log(`BROWSER LOG: ${msg.text()}`))
    page.on('pageerror', (err) => console.log(`BROWSER ERROR: ${err.message}`))
    page.on('requestfailed', (request) =>
      console.log(`REQUEST FAILED: ${request.url()} - ${request.failure()?.errorText}`)
    )

    await mockCdnAssets(page)

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
    // Note: The example HTML uses x-tooltip.top, which is still supported as alias.
    // We should also test the new x-uno-tooltip directive in a separate test or modify this one.
    // Let's modify the test to use data attributes if we can, or add a new test.
    // Since we are running e2e against existing HTML files in examples/, I should probably not break those files.
    // Instead, I will inject new elements using the new directive in the test.
    
    await page.evaluate(() => {
        const btn = document.createElement('button')
        btn.setAttribute('data-testid', 'trigger-uno-top')
        btn.setAttribute('x-uno-tooltip', "'Hello Uno Top!'")
        btn.setAttribute('data-uno-placement', 'top')
        btn.textContent = 'Uno Top'
        btn.className = 'px-3 py-2 bg-blue-500 text-white rounded'
        document.body.appendChild(btn)
        
        // @ts-ignore
        const AlpineRef = window.Alpine || window.top.Alpine
        AlpineRef.initTree(btn)
    })

    const trigger = page.getByTestId('trigger-uno-top')
    const tooltipContent = page.locator('#uno-tooltip-content')

    await expect(tooltipContent).not.toBeVisible()
    await trigger.hover()

    await expect(tooltipContent).toBeVisible()
    await expect(tooltipContent).toHaveText('Hello Uno Top!')

    const triggerBox = await trigger.boundingBox()
    const tooltipBox = await tooltipContent.boundingBox()

    if (triggerBox && tooltipBox) {
      expect(tooltipBox.y + tooltipBox.height).toBeLessThanOrEqual(triggerBox.y + 5)
    }
  })

  test('shows on hover (bottom) via data attribute', async ({ page }) => {
    await page.evaluate(() => {
        const btn = document.createElement('button')
        btn.setAttribute('data-testid', 'trigger-uno-bottom')
        btn.setAttribute('x-uno-tooltip', "'Hello Uno Bottom!'")
        btn.setAttribute('data-uno-placement', 'bottom')
        btn.textContent = 'Uno Bottom'
        btn.className = 'px-3 py-2 bg-blue-500 text-white rounded'
        document.body.appendChild(btn)
        
        // @ts-ignore
        const AlpineRef = window.Alpine || window.top.Alpine
        AlpineRef.initTree(btn)
    })

    const trigger = page.getByTestId('trigger-uno-bottom')
    const tooltipContent = page.locator('#uno-tooltip-content')

    await trigger.hover()
    await expect(tooltipContent).toBeVisible()
    await expect(tooltipContent).toHaveText('Hello Uno Bottom!')

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

  test('updates position on scroll when visible (click trigger)', async ({ page }) => {
    const trigger = page.getByTestId('trigger-click')
    const tooltipContent = page.locator('#uno-tooltip-content')

    await trigger.scrollIntoViewIfNeeded()
    await trigger.click()
    await expect(tooltipContent).toBeVisible()

    const before = {
      trigger: await trigger.boundingBox(),
      tooltip: await tooltipContent.boundingBox()
    }
    if (!before.trigger || !before.tooltip) throw new Error('Missing bounding box before scroll')

    await page.evaluate(() => window.scrollBy(0, 240))
    await page.evaluate(
      () => new Promise((r) => requestAnimationFrame(() => requestAnimationFrame(r)))
    )

    const after = {
      trigger: await trigger.boundingBox(),
      tooltip: await tooltipContent.boundingBox()
    }
    if (!after.trigger || !after.tooltip) throw new Error('Missing bounding box after scroll')

    const beforeOffsetY = before.tooltip.y - before.trigger.y
    const afterOffsetY = after.tooltip.y - after.trigger.y
    expect(Math.abs(afterOffsetY - beforeOffsetY)).toBeLessThan(12)
    expect(after.tooltip.x).toBeGreaterThanOrEqual(after.trigger.x + after.trigger.width - 2)
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

  test('auto flips near top edge (top -> bottom)', async ({ page }) => {
    await page.evaluate(() => {
      const wrap = document.createElement('div')
      wrap.style.position = 'fixed'
      wrap.style.left = '16px'
      wrap.style.top = '16px'
      wrap.style.zIndex = '999999'

      const btn = document.createElement('button')
      btn.setAttribute('data-testid', 'trigger-top-edge')
      btn.setAttribute('x-tooltip.top', "'I flipped to bottom!'")
      btn.textContent = 'Top Edge'
      btn.className = 'px-3 py-2 bg-slate-800 text-white rounded'

      wrap.appendChild(btn)
      document.body.appendChild(wrap)

      // @ts-ignore
      const AlpineRef = window.Alpine || window.top.Alpine
      AlpineRef.initTree(wrap)
    })

    const trigger = page.getByTestId('trigger-top-edge')
    const tooltipContent = page.locator('#uno-tooltip-content')

    await trigger.hover()
    await expect(tooltipContent).toBeVisible()
    await expect(tooltipContent).toHaveText('I flipped to bottom!')

    const triggerBox = await trigger.boundingBox()
    const tooltipBox = await tooltipContent.boundingBox()
    if (!triggerBox || !tooltipBox) throw new Error('Missing bounding box')

    expect(tooltipBox.y).toBeGreaterThanOrEqual(triggerBox.y + triggerBox.height - 1)
  })

  test('auto flips near bottom edge (bottom -> top)', async ({ page }) => {
    await page.evaluate(() => {
      const wrap = document.createElement('div')
      wrap.style.position = 'fixed'
      wrap.style.left = '16px'
      wrap.style.bottom = '16px'
      wrap.style.zIndex = '1'

      const btn = document.createElement('button')
      btn.setAttribute('data-testid', 'trigger-bottom-edge')
      btn.setAttribute('x-tooltip.bottom', "'I flipped to top!'")
      btn.textContent = 'Bottom Edge'
      btn.className = 'px-3 py-2 bg-slate-800 text-white rounded'

      wrap.appendChild(btn)
      document.body.appendChild(wrap)

      // @ts-ignore
      const AlpineRef = window.Alpine || window.top.Alpine
      AlpineRef.initTree(wrap)
    })

    const trigger = page.getByTestId('trigger-bottom-edge')
    const tooltipContent = page.locator('#uno-tooltip-content')

    await trigger.hover()
    await expect(tooltipContent).toBeVisible()
    await expect(tooltipContent).toHaveText('I flipped to top!')

    const triggerBox = await trigger.boundingBox()
    const tooltipBox = await tooltipContent.boundingBox()
    if (!triggerBox || !tooltipBox) throw new Error('Missing bounding box')

    expect(tooltipBox.y + tooltipBox.height).toBeLessThanOrEqual(triggerBox.y + 1)
  })

  test('click trigger does not show on hover', async ({ page }) => {
    await page.evaluate(() => {
      const btn = document.createElement('button')
      btn.setAttribute('data-testid', 'trigger-click-only')
      btn.setAttribute('x-tooltip.click', "'Click Only'")
      btn.textContent = 'Click Only'
      btn.className = 'px-3 py-2 bg-slate-800 text-white rounded'
      document.body.appendChild(btn)

      // @ts-ignore
      const AlpineRef = window.Alpine || window.top.Alpine
      AlpineRef.initTree(btn)
    })

    const trigger = page.getByTestId('trigger-click-only')
    const tooltipContent = page.locator('#uno-tooltip-content')

    await trigger.hover()
    await expect(tooltipContent).not.toBeVisible()

    await trigger.click()
    await expect(tooltipContent).toBeVisible()
    await expect(tooltipContent).toHaveText('Click Only')
  })

  test('shows on focus and hides on blur', async ({ page }) => {
    await page.evaluate(() => {
      const input = document.createElement('input')
      input.setAttribute('data-testid', 'trigger-focus')
      input.setAttribute('x-tooltip.focus', "'Focus Tip'")
      input.className = 'border px-2 py-1'
      document.body.appendChild(input)

      // @ts-ignore
      const AlpineRef = window.Alpine || window.top.Alpine
      AlpineRef.initTree(input)
    })

    const trigger = page.getByTestId('trigger-focus')
    const tooltipContent = page.locator('#uno-tooltip-content')

    await trigger.focus()
    await expect(tooltipContent).toBeVisible()
    await expect(tooltipContent).toHaveText('Focus Tip')

    await page.mouse.click(0, 0)
    await expect(tooltipContent).not.toBeVisible()
  })

  test('disabled trigger degrades to parent wrapper for hover', async ({ page }) => {
    await page.evaluate(() => {
      const wrap = document.createElement('span')
      wrap.setAttribute('data-testid', 'disabled-wrap')
      wrap.style.display = 'inline-block'

      const btn = document.createElement('button')
      btn.disabled = true
      btn.setAttribute('x-tooltip', "'Disabled Tooltip'")
      btn.textContent = 'Disabled'
      btn.className = 'px-4 py-2 bg-gray-300 text-gray-700 rounded'

      wrap.appendChild(btn)
      document.body.appendChild(wrap)

      // @ts-ignore
      const AlpineRef = window.Alpine || window.top.Alpine
      AlpineRef.initTree(wrap)
    })

    const wrap = page.getByTestId('disabled-wrap')
    const tooltipContent = page.locator('#uno-tooltip-content')

    await wrap.hover()
    await expect(tooltipContent).toBeVisible()
    await expect(tooltipContent).toHaveText('Disabled Tooltip')
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
    await page.waitForSelector('[data-testid="trigger-top"]')
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
