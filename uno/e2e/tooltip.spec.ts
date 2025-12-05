import { test, expect } from '@playwright/test';

test.describe('Tooltip Component', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/examples/components/tooltip.html');
  });

  test('shows on hover (top)', async ({ page }) => {
    const trigger = page.getByTestId('trigger-top');
    const tooltip = page.getByTestId('tooltip-top');

    await expect(tooltip).not.toBeVisible();
    await trigger.hover();
    await expect(tooltip).toBeVisible();
    await expect(tooltip).toHaveText(/Hello Top!/);
    
    const triggerBox = await trigger.boundingBox();
    const tooltipBox = await tooltip.boundingBox();
    
    console.log('Trigger Box:', triggerBox);
    console.log('Tooltip Box:', tooltipBox);

    if (triggerBox && tooltipBox) {
        expect(tooltipBox.y + tooltipBox.height).toBeLessThanOrEqual(triggerBox.y);
    }
  });

  test('shows on hover (bottom)', async ({ page }) => {
    const trigger = page.getByTestId('trigger-bottom');
    const tooltip = page.getByTestId('tooltip-bottom');

    await trigger.hover();
    await expect(tooltip).toBeVisible();
    
    const triggerBox = await trigger.boundingBox();
    const tooltipBox = await tooltip.boundingBox();
    
    if (triggerBox && tooltipBox) {
        expect(tooltipBox.y).toBeGreaterThanOrEqual(triggerBox.y + triggerBox.height);
    }
  });

  test('toggle on click and light theme', async ({ page }) => {
    const trigger = page.getByTestId('trigger-click');
    const tooltip = page.getByTestId('tooltip-click');

    await trigger.click();
    await expect(tooltip).toBeVisible();
    await expect(tooltip).toHaveText(/Clicked \(Light\)!/);
    
    await expect(tooltip).toHaveClass(/bg-white/);

    await trigger.click();
    await expect(tooltip).not.toBeVisible();
  });

  test('target mode with HTML content', async ({ page }) => {
    const trigger = page.getByTestId('trigger-target');
    const tooltip = page.getByTestId('tooltip-target');

    await trigger.hover();
    await expect(tooltip).toBeVisible();
    
    const boldText = tooltip.locator('strong');
    await expect(boldText).toHaveText('Bold Title');
    const italicText = tooltip.locator('em');
    await expect(italicText).toHaveText('Italic description');
  });

  test('smart positioning (edge detection)', async ({ page }) => {
    const trigger = page.getByTestId('trigger-edge');
    const tooltip = page.getByTestId('tooltip-edge');

    // 这个 trigger 在页面顶部 (top=0)
    // 即使配置 position='top'，它应该翻转到 'bottom'
    await trigger.hover();
    await expect(tooltip).toBeVisible();
    
    const triggerBox = await trigger.boundingBox();
    const tooltipBox = await tooltip.boundingBox();
    
    if (triggerBox && tooltipBox) {
        // Tooltip 应该在 Trigger 下方
        expect(tooltipBox.y).toBeGreaterThanOrEqual(triggerBox.y + triggerBox.height);
    }
  });
});
