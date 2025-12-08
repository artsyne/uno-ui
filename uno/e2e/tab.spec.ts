import { test, expect } from '@playwright/test'

test.describe('Tab Component', () => {
  test.beforeEach(async ({ page }) => {
    // Enable console logging
    page.on('console', (msg) => console.log(`[Browser Console] ${msg.text()}`))
    page.on('pageerror', (err) => console.log(`[Page Error] ${err.message}`))

    await page.goto('/examples/components/tab.html')
    await page.waitForLoadState('domcontentloaded')

    // Wait for Alpine to initialize
    try {
      await page.waitForFunction(() => (window as any).Alpine, null, { timeout: 5000 });
    } catch (e) {
      console.log('Alpine not found on window!');
    }
  })

  test('modern tabs switch content and hash', async ({ page }) => {
    const tab2 = page.locator('button[data-tab-target="tab2"]');
    const content1 = page.locator('#tab1')
    const content2 = page.locator('#tab2')

    // Initial state
    await expect(content1).toBeVisible()
    await expect(content2).toBeHidden()
    // await expect(tab2).not.toHaveClass(/current/); 

    // Click tab 2
    await tab2.click()

    // Verify switch
    await expect(content1).toBeHidden()
    await expect(content2).toBeVisible()
    // await expect(tab2).toHaveClass(/current/); 

    // Verify Hash
    expect(page.url()).toContain('#settings')
  })

  // Legacy test removed as legacy attributes are no longer supported in the new Alpine component
  // test('legacy attributes work correctly', ...);

  test('separated dom elements work', async ({ page }) => {
    const sTab2 = page.locator('button[data-tab-target="s-tab2"]');
    const sContent1 = page.locator('#s-tab1');
    const sContent2 = page.locator('#s-tab2');

    await sTab2.click();
    await expect(sContent1).toBeHidden();
    await expect(sContent2).toBeVisible();
  });

  test('deep linking loads correct tab', async ({ page }) => {
    // Navigate directly with hash
    await page.goto('/examples/components/tab.html#billing');
    await page.waitForLoadState('domcontentloaded');
    await page.waitForFunction(() => (window as any).Alpine);
    
    const content3 = page.locator('#tab3');
    const tab3 = page.locator('button[data-tab-target="tab3"]')

    await expect(content3).toBeVisible();
    await expect(tab3).toHaveClass(/current/);
  })

  test('ajax content loading', async ({ page }) => {
    const ajaxTab = page.locator('button[data-tab-target="s-tab3"]')
    const ajaxContent = page.locator('#s-tab3')

    // Initially hidden
    await expect(ajaxContent).toBeHidden()
    await expect(ajaxContent).not.toContainText('Ajax Loaded Content');

    // Click to load
    await ajaxTab.click()

    // Should be visible
    await expect(ajaxContent).toBeVisible()

    // Should load content
    await expect(ajaxContent).toContainText('Ajax Loaded Content')

    // Verify loaded attribute
    await expect(ajaxContent).toHaveAttribute('data-tab-loaded', 'true')
  })
})
