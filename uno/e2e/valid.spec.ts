import { test, expect } from '@playwright/test';

test.describe('Valid Plugin', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/examples/components/valid.html');
    await page.waitForLoadState('domcontentloaded');
    await page.evaluate(() => (window as any).Alpine?.store('i18n')?.setLang('en'));
    await page.evaluate(() => {
      const ov = document.querySelector('vite-error-overlay');
      if (ov) ov.remove();
    });
  });

  test('email valid with type=email', async ({ page }) => {
    // Use first form (Built-in Validators)
    const form = page.locator('form').first();
    await form.locator('input[name="username"]').fill('user_123');
    const emailInput = form.locator('input[name="email"]');
    await emailInput.fill('biuren@gmail.com');
    await emailInput.focus();
    await page.keyboard.press('Tab');
    
    // Check hint specifically for email input
    const hint = emailInput.locator('xpath=following-sibling::div[contains(@class, "valid-hint")]');
    await expect(hint).not.toContainText('Invalid format');
  });

  test('textarea char count updates and truncates', async ({ page }) => {
    const ta = page.locator('textarea').first(); // There's only one textarea now
    await ta.focus();
    await ta.type('a'.repeat(10));
    
    // Count span is now inserted dynamically
    await page.waitForSelector('.valid-count');
    const count = page.locator('.valid-count');
    // max is 100 in new html
    await expect(count).toContainText(/Remaining:\s*90/);
    
    await ta.clear();
    await ta.type('a'.repeat(150));
    const value = await ta.inputValue();
    expect(value.length).toBe(100);
  });

  test('invalid prevents submit and shows hints (HTMX form)', async ({ page }) => {
    // Use the last form (Async & HTMX) which has hx-post
    const form = page.locator('form').last();
    
    // Trigger submit without filling required fields
    await form.evaluate((f: HTMLFormElement) => {
      f.dispatchEvent(new Event('submit', { bubbles: true, cancelable: true }));
    });

    const hint = form.locator('.valid-hint').first();
    await expect(hint).toContainText('Required');
    
    const result = page.locator('#htmx-result');
    await expect(result).toBeEmpty();
  });

  test('valid submits and swaps result (HTMX form)', async ({ page }) => {
    const form = page.locator('form').last();
    const input = form.locator('input[name="username"]');
    
    // 'checkUserUnique' allows anything except 'admin'/'guest'
    await input.fill('validuser');
    // Wait for async check to pass
    await page.waitForTimeout(1000);
    await expect(input).toHaveClass(/is-valid/);

    await form.evaluate((f: HTMLFormElement) => {
      f.dispatchEvent(new Event('submit', { bubbles: true, cancelable: true }));
    });
    
    // Wait for HTMX swap
    await page.waitForTimeout(500);
    const result = page.locator('#htmx-result');
    await expect(result).toContainText('Sample Data');
  });

  test('async callback validation', async ({ page }) => {
    // Async input is in the last form
    const input = page.locator('#async_user');
    
    await input.fill('admin');
    await input.blur();
    // Wait for async check (simulated 800ms)
    await page.waitForTimeout(1200); 
    
    const hint = input.locator('xpath=following-sibling::div[contains(@class, "valid-hint")]');
    await expect(hint).toContainText('Username "admin" is taken');

    await input.fill('newuser');
    await input.blur();
    await page.waitForTimeout(1200);
    await expect(input).toHaveClass(/is-valid/);
    await expect(input).not.toHaveClass(/is-invalid/);
  });

  test('range validation', async ({ page }) => {
    // Range input is in the second form (Logic Rules)
    const input = page.locator('input[name="age"]');
    await input.fill('10'); // min is 18
    await input.blur();
    const hint = input.locator('xpath=following-sibling::div[contains(@class, "valid-hint")]');
    await expect(hint).toContainText('Only 18-60 allowed');

    await input.fill('20');
    await input.blur();
    await expect(input).toHaveClass(/is-valid/);
  });

  test('compareto validation', async ({ page }) => {
    // Logic Rules form
    const pwd = page.locator('#logic_pwd');
    const confirm = page.locator('input[name="confirm"]');
    
    await pwd.fill('Secret123');
    await confirm.fill('Secret456');
    await confirm.blur();
    
    const hint = confirm.locator('xpath=following-sibling::div[contains(@class, "valid-hint")]');
    await expect(hint).toContainText('Must match password');

    await confirm.fill('Secret123');
    await confirm.blur();
    await expect(confirm).toHaveClass(/is-valid/);
  });
});
