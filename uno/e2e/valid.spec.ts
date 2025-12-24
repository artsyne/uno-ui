import { test, expect } from '@playwright/test'
import { mockCdnAssets } from './helpers/routes'

test.describe('Valid Plugin (Standard HTML5)', () => {
  test.beforeEach(async ({ page }) => {
    await mockCdnAssets(page)

    await page.goto('/examples/components/valid.html')
    await page.waitForLoadState('domcontentloaded')
    // Wait for HTMX to load
    await page.waitForFunction(() => (window as any).htmx)
  })

  test('validates required fields and format', async ({ page }) => {
    const form = page.locator('form').first()
    const username = form.locator('input[name="username"]')
    const email = form.locator('input[name="email"]')
    const submit = form.locator('button[type="submit"]')

    // 1. Submit empty -> Required errors (Triggered by HTMX interception)
    await submit.click()
    
    // Check for is-invalid class
    await expect(username).toHaveClass(/is-invalid/)
    await expect(email).toHaveClass(/is-invalid/)
    
    // Check for error messages
    await expect(form.locator('[data-uno-valid-error-for="username"]')).toBeVisible() // Browser default message or custom?
    // valid.js uses validationMessage. If empty and required, it's "Please fill out this field."
    // But wait, valid.js uses getAttribute('data-uno-error-required') if available?
    // Let's check valid.js... yes it does check data-uno-error-*
    // But my example HTML didn't set data-uno-error-required.
    // It set data-uno-error-pattern.
    // So for required, it uses browser default.

    // 2. Invalid format
    await username.fill('Invalid Name!') 
    await username.blur()
    // Blur doesn't trigger validation in my logic unless aria-invalid is true?
    // valid.js listens to 'input'.
    // If I type, it validates if aria-invalid is true.
    // If I submit once, aria-invalid becomes true.
    // So if I fix it, it should clear.
    
    // But wait, 'Invalid Name!' is invalid pattern.
    // So it should stay invalid.
    await expect(username).toHaveClass(/is-invalid/)
    await expect(form.locator('[data-uno-valid-error-for="username"]')).toContainText('Only letters, numbers, _, - allowed')

    // 3. Valid format
    await username.fill('valid_user')
    // Input event fires, valid.js checks.
    await expect(username).not.toHaveClass(/is-invalid/)
    await expect(form.locator('[data-uno-valid-error-for="username"]')).toBeHidden()

    // 4. Email format
    await email.fill('not-an-email')
    // Previously invalid (from submit), so input triggers validation
    await expect(email).toHaveClass(/is-invalid/)
    await expect(form.locator('[data-uno-valid-error-for="email"]')).toContainText('Please enter a valid email address')

    await email.fill('test@example.com')
    await expect(email).not.toHaveClass(/is-invalid/)
  })

  test('logic rules (password compare & range)', async ({ page }) => {
    // Second form
    const form = page.locator('form').nth(1)
    const pwd = form.locator('input[name="password"]')
    const confirm = form.locator('input[name="confirm"]')
    const age = form.locator('input[name="age"]')
    const submit = form.locator('button[type="submit"]')

    // Trigger validation to enable aggressive checking
    await submit.click()

    // Password compare
    await pwd.fill('Secret123')
    await confirm.fill('Different')
    await confirm.blur() // Trigger validation
    await expect(confirm).toHaveClass(/is-invalid/)
    // Native validation doesn't have a standard message for custom validity unless set.
    // valid.js sets setCustomValidity('Mismatch') if match fails.
    // And error container textContent = input.validationMessage.
    // So it should be 'Mismatch' (browser specific? No, I set it).
    
    await confirm.fill('Secret123')
    await confirm.blur()
    await expect(confirm).not.toHaveClass(/is-invalid/)

    // Range (18-60)
    await age.fill('10')
    await age.blur() // Trigger validation
    await expect(age).toHaveClass(/is-invalid/)
    // Browser message for range underflow

    await age.fill('25')
    await age.blur()
    await expect(age).not.toHaveClass(/is-invalid/)
  })

  test('character count limits', async ({ page }) => {
    // Third form
    const form = page.locator('form').nth(2)
    const bio = form.locator('textarea[name="bio"]')
    const counter = form.locator('.valid-count span')

    // Initial state
    await expect(counter).toHaveText('0')

    // Type 10 chars
    await bio.fill('1234567890')
    await expect(counter).toHaveText('10')

    // Exceed limit (should truncate via maxlength)
    const longText = 'a'.repeat(120)
    await bio.fill(longText)

    // Value should be truncated to 100
    const val = await bio.inputValue()
    expect(val.length).toBe(100)
    await expect(counter).toHaveText('100')
  })
})
