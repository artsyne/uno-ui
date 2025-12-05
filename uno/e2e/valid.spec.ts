
import { test, expect } from '@playwright/test';

test.describe('Valid Plugin (Data-Driven)', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/examples/components/valid.html');
    await page.waitForLoadState('domcontentloaded');
    // Set language to English
    await page.evaluate(() => (window as any).Alpine?.store('i18n')?.setLang('en'));
    // Remove overlay if present
    await page.evaluate(() => {
      const ov = document.querySelector('vite-error-overlay');
      if (ov) ov.remove();
    });
  });

  test('validates required fields and format', async ({ page }) => {
    const form = page.locator('form').first();
    const username = form.locator('input[name="username"]');
    const email = form.locator('input[name="email"]');
    const submit = form.locator('button[type="submit"]');

    // 1. Submit empty -> Required errors
    await submit.click();
    await expect(username).toHaveClass(/is-invalid/);
    await expect(form.locator('[x-valid-error="username"]')).toContainText('Required');
    await expect(form.locator('[x-valid-error="email"]')).toContainText('Required');

    // 2. Invalid format
    await username.fill('Invalid Name!'); // validname allows a-z0-9_-
    await username.blur();
    await expect(form.locator('[x-valid-error="username"]')).toContainText('Letters, numbers, - and _ only');

    // 3. Valid format
    await username.fill('valid_user');
    await username.blur();
    await expect(username).toHaveClass(/is-valid/);
    await expect(form.locator('[x-valid-error="username"]')).toBeHidden();

    // 4. Email format
    await email.fill('not-an-email');
    await email.blur();
    await expect(form.locator('[x-valid-error="email"]')).toContainText('Invalid email address');
    
    await email.fill('test@example.com');
    await email.blur();
    await expect(email).toHaveClass(/is-valid/);
  });

  test('logic rules (password compare & range)', async ({ page }) => {
    // Second form
    const form = page.locator('form').nth(1);
    const pwd = form.locator('input[name="password"]');
    const confirm = form.locator('input[name="confirm"]');
    const age = form.locator('input[name="age"]');

    // Password compare
    await pwd.fill('Secret123');
    await confirm.fill('Different');
    await confirm.blur();
    await expect(form.locator('[x-valid-error="confirm"]')).toContainText('Must match password');

    await confirm.fill('Secret123');
    await confirm.blur();
    await expect(confirm).toHaveClass(/is-valid/);

    // Range (18-60)
    await age.fill('10');
    await age.blur();
    await expect(form.locator('[x-valid-error="age"]')).toContainText('Only 18-60 allowed');

    await age.fill('25');
    await age.blur();
    await expect(age).toHaveClass(/is-valid/);
  });

  test('character count limits', async ({ page }) => {
    // Third form
    const form = page.locator('form').nth(2);
    const bio = form.locator('textarea[name="bio"]');
    const counter = form.locator('[x-valid-count="bio"]');

    // Initial state (empty) -> 100 remaining
    await bio.click(); // focus to init if needed? no, x-validate inits on load
    // Wait for alpine to init directives
    await page.waitForTimeout(100); 
    
    // Type 10 chars
    await bio.fill('1234567890');
    await expect(counter).toContainText('Remaining: 90');

    // Exceed limit (should truncate)
    const longText = 'a'.repeat(120);
    await bio.fill(longText);
    
    // Value should be truncated to 100
    const val = await bio.inputValue();
    expect(val.length).toBe(100);
    await expect(counter).toContainText('Remaining: 0');
  });

  test('async validation and htmx submission', async ({ page }) => {
    // Fourth form
    const form = page.locator('form').last();
    const username = form.locator('input[name="username"]');
    const submit = form.locator('button[type="submit"]');
    const result = form.locator('#htmx-result');

    // 1. Async Error (admin)
    await username.fill('admin');
    await username.blur(); // Trigger validation
    // Spinner should appear
    const spinner = form.locator('.animate-spin');
    await expect(spinner).toBeVisible();
    // Wait for validation
    await expect(spinner).toBeHidden();
    await expect(form.locator('[x-valid-error="username"]')).toContainText('Username "admin" is taken');
    await expect(username).toHaveClass(/is-invalid/);

    // 2. Async Success (validuser)
    await username.fill('validuser');
    await username.blur();
    await expect(spinner).toBeVisible();
    await expect(spinner).toBeHidden();
    await expect(username).toHaveClass(/is-valid/);

    // 3. Submit (HTMX)
    // Ensure form is valid
    await submit.click();
    
    // HTMX request should happen
    // We check if the result container gets populated
    await expect(result).toContainText('Sample Data');
  });
});
