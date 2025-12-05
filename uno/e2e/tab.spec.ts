import { test, expect } from '@playwright/test';

test.describe('Tab Component', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/examples/components/tab.html');
    await page.waitForLoadState('domcontentloaded');
  });

  test('modern tabs switch content and hash', async ({ page }) => {
    const tab2 = page.locator('button[data-tab-target="tab2"]');
    const content1 = page.locator('#tab1');
    const content2 = page.locator('#tab2');

    // Initial state
    await expect(content1).toBeVisible();
    await expect(content2).toBeHidden();
    await expect(tab2).not.toHaveClass(/current/); // active class check

    // Click tab 2
    await tab2.click();

    // Verify switch
    await expect(content1).toBeHidden();
    await expect(content2).toBeVisible();
    await expect(tab2).toHaveClass(/current/); // Active class applied
    
    // Verify Hash
    expect(page.url()).toContain('#settings');
  });

  test('legacy attributes work correctly', async ({ page }) => {
    const lTab2 = page.locator('a[tabtarget="l-tab2"]');
    const lContent1 = page.locator('#l-tab1');
    const lContent2 = page.locator('#l-tab2');

    await lTab2.click();

    await expect(lContent1).toBeHidden();
    await expect(lContent2).toBeVisible();
    expect(page.url()).toContain('#legacy2');
  });

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
    
    const content3 = page.locator('#tab3');
    const tab3 = page.locator('button[data-tab-target="tab3"]');

    await expect(content3).toBeVisible();
    await expect(tab3).toHaveClass(/current/);
  });

  test('ajax content loading', async ({ page }) => {
    const ajaxTab = page.locator('button[data-tab-target="s-tab3"]');
    const ajaxContent = page.locator('#s-tab3');

    // Initially hidden and empty (or just comment)
    await expect(ajaxContent).toBeHidden();
    await expect(ajaxContent).not.toContainText('Ajax Loaded Content');

    // Click to load
    await ajaxTab.click();
    
    // Should be visible
    await expect(ajaxContent).toBeVisible();
    
    // Should load content
    // It might take a moment for fetch to complete, expect automatically retries
    await expect(ajaxContent).toContainText('Ajax Loaded Content');
    
    // Verify loaded attribute
    await expect(ajaxContent).toHaveAttribute('data-tab-loaded', 'true');
  });
});
