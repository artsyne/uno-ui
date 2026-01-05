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

    // Check for invalid class (HTML uses !border-red-500)
    await expect(username).toHaveClass(/!border-red-500/)
    await expect(email).toHaveClass(/!border-red-500/)

    // Check for error messages
    await expect(form.locator('div[x-text="$store.uno.valid.errors[\'username\']"]')).toBeVisible()

    // 2. Invalid format
    await username.fill('Invalid Name!')
    await username.blur()

    await expect(username).toHaveClass(/!border-red-500/)
    await expect(form.locator('div[x-text="$store.uno.valid.errors[\'username\']"]')).toContainText(
      'Only letters, numbers, _, - allowed'
    )

    // 3. Valid format
    await username.fill('valid_user')
    // Input event fires, valid.js checks.
    await expect(username).not.toHaveClass(/!border-red-500/)
    await expect(form.locator('div[x-text="$store.uno.valid.errors[\'username\']"]')).toBeHidden()

    // 4. Email format
    await email.fill('not-an-email')
    // Previously invalid (from submit), so input triggers validation
    await expect(email).toHaveClass(/!border-red-500/)
    await expect(form.locator('div[x-text="$store.uno.valid.errors[\'email\']"]')).toContainText(
      'Please enter a valid email address'
    )

    await email.fill('test@example.com')
    await expect(email).not.toHaveClass(/!border-red-500/)
  })

  test('logic rules (password compare & range)', async ({ page }) => {
    // Second form
    const form = page.locator('form').nth(1)
    const pwd = form.locator('input[name="logic_password"]')
    const confirm = form.locator('input[name="confirm"]')
    const age = form.locator('input[name="logic_age"]')
    const submit = form.locator('button[type="submit"]')

    // Trigger validation to enable aggressive checking
    await submit.click()

    // Password compare
    await pwd.fill('Secret123')
    await confirm.fill('Different')
    await confirm.blur() // Trigger validation
    await expect(confirm).toHaveClass(/!border-red-500/)
    // Native validation doesn't have a standard message for custom validity unless set.
    // valid.js sets setCustomValidity('Mismatch') if match fails.
    // And error container textContent = input.validationMessage.
    // So it should be 'Mismatch' (browser specific? No, I set it).

    await confirm.fill('Secret123')
    await confirm.blur()
    await expect(confirm).not.toHaveClass(/!border-red-500/)

    // Range (18-60)
    await age.fill('10')
    await age.blur() // Trigger validation
    await expect(age).toHaveClass(/!border-red-500/)
    // Browser message for range underflow

    await age.fill('25')
    await age.blur()
    await expect(age).not.toHaveClass(/!border-red-500/)
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
    const longText = 'a'.repeat(250)
    await bio.fill(longText)

    // Value should be truncated to 200
    const val = await bio.inputValue()
    expect(val.length).toBe(200)
    await expect(counter).toHaveText('200')
  })
})
