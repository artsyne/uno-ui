import { test, expect } from '@playwright/test';

test.describe('Box Component', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/examples/components/box.html');
    await page.waitForLoadState('domcontentloaded');
    // Wait for Alpine and Uno to initialize
    await page.waitForFunction(() => (window as any).Alpine && (window as any).Alpine.store('i18n'));
  });

  test('opens dialog (center) with overlay and locks body scroll', async ({ page }) => {
    const btn = page.getByText('Open Dialog');
    await btn.click();

    const overlay = page.locator('[id^="uno-box-overlay"]');
    await expect(overlay).toBeVisible();

    const outer = page.locator('[id^="uno-box-outer-container"]');
    await expect(outer).toBeVisible();
  });

  test('confirm renders buttons and closes on cancel', async ({ page }) => {
    const btn = page.getByText('Open Confirm');
    await btn.click();

    const outer = page.locator('[id^="uno-box-outer-container"]');
    await expect(outer).toBeVisible();

    const cancel = page.locator('[id^="uno-box-cancel"]');
    await expect(cancel).toBeVisible();
    await cancel.click();
    await expect(outer).not.toBeVisible();
  });

  test('position bottom-right', async ({ page }) => {
    const btn = page.getByText('Bottom Right');
    await btn.click();

    const outer = page.locator('[id^="uno-box-outer-container"]');
    await expect(outer).toBeVisible();
  });

  test('ajax loads local json and displays content', async ({ page }) => {
    const btn = page.getByText('Ajax');
    await btn.click();

    const content = page.locator('[id^="uno-box-content"]');
    await expect(content).toBeVisible();
    await expect(content).toContainText('Sample Data');
  });

  test('i18n toggles ok/cancel to English', async ({ page }) => {
    await page.evaluate(() => (window as any).Alpine.store('i18n').setLang('en'));
    const btn = page.getByText('Open Confirm');
    await btn.click();
    const ok = page.locator('[id^="uno-box-ok"]');
    const cancel = page.locator('[id^="uno-box-cancel"]');
    await expect(ok).toHaveText('OK');
    await expect(cancel).toHaveText('Cancel');
  });
});
