import { test, expect } from '@playwright/test';

test.describe('Notify Component', () => {
  test.beforeEach(async ({ page }) => {
    page.on('console', msg => console.log(`[Browser Console] ${msg.text()}`));
    page.on('pageerror', err => console.log(`[Page Error] ${err.message}`));
    await page.goto('/examples/components/notify.html');
  });

  test('shows info notification (top-right) and auto dismiss', async ({ page }) => {
    const btn = page.getByTestId('btn-info');
    const container = page.getByTestId('notify-container-default');
    await btn.click();

    const item = page.getByTestId('notify-item');
    await expect(item).toBeVisible();
    await expect(item).toHaveText(/Info message/);

    const box = await item.boundingBox();
    const vp = page.viewportSize();
    if (box && vp) {
      expect(box.y).toBeLessThanOrEqual(100);
      expect(box.x + box.width).toBeGreaterThanOrEqual(vp.width - 20);
    }

    // Auto dismiss within ~3s
    await page.waitForTimeout(3200);
    await expect(item).not.toBeVisible();
  });

  test('sticky notification requires click to close', async ({ page }) => {
    const btn = page.getByTestId('btn-sticky');
    await btn.click();
    const item = page.getByTestId('notify-item');
    await expect(item).toBeVisible();

    await page.waitForTimeout(2000);
    await expect(item).toBeVisible();

    await item.click();
    await expect(item).not.toBeVisible();
  });

  test('bottom-left position', async ({ page }) => {
    const btn = page.getByTestId('btn-bl');
    await btn.click();

    const item = page.getByTestId('notify-item-bl');
    await expect(item).toBeVisible();

    const box = await item.boundingBox();
    const vp = page.viewportSize();
    if (box && vp) {
      expect(box.y + box.height).toBeGreaterThanOrEqual(vp.height - 20);
      expect(box.x).toBeLessThanOrEqual(100);
    }
  });
});
