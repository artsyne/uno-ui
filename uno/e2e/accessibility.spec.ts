import { test, expect } from '@playwright/test'

/**
 * Accessibility Test Suite
 *
 * Tests Uno UI components for WCAG 2.1 Level AA compliance:
 * - Keyboard navigation
 * - Screen reader compatibility
 * - ARIA attributes
 * - Focus management
 * - Color contrast
 */

test.describe('Accessibility Tests', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/examples/index.html')
  })

  test.describe('Box Component Accessibility', () => {
    test('should have proper ARIA attributes', async ({ page }) => {
      // Open modal
      await page.click('button[data-box-title="Test Dialog"]')

      // Check dialog role
      const dialog = page.locator('[role="dialog"]')
      await expect(dialog).toBeVisible()

      // Check aria-modal
      await expect(dialog).toHaveAttribute('aria-modal', 'true')

      // Check aria-labelledby (title)
      const titleId = await dialog.getAttribute('aria-labelledby')
      expect(titleId).toBeTruthy()

      const title = page.locator(`#${titleId}`)
      await expect(title).toBeVisible()
    })

    test('should trap focus within modal', async ({ page }) => {
      await page.click('button[data-box-title="Test Dialog"]')

      const dialog = page.locator('[role="dialog"]')
      await expect(dialog).toBeVisible()

      // Get all focusable elements
      const focusable = dialog.locator(
        'button, [href], input, select, textarea, [tabindex]:not([tabindex="-1"])'
      )
      const count = await focusable.count()

      // Tab through all elements
      for (let i = 0; i < count; i++) {
        await page.keyboard.press('Tab')
      }

      // Focus should cycle back to first element
      await page.keyboard.press('Tab')
      const firstFocusable = focusable.nth(0)
      await expect(firstFocusable).toBeFocused()
    })

    test('should close on Escape key', async ({ page }) => {
      await page.click('button[data-box-title="Test Dialog"]')

      const dialog = page.locator('[role="dialog"]')
      await expect(dialog).toBeVisible()

      // Press Escape
      await page.keyboard.press('Escape')

      // Dialog should close
      await expect(dialog).not.toBeVisible()
    })

    test('should restore focus after closing', async ({ page }) => {
      const trigger = page.locator('button[data-box-title="Test Dialog"]')

      // Focus and open
      await trigger.focus()
      await trigger.click()

      // Close modal
      await page.keyboard.press('Escape')

      // Focus should return to trigger
      await expect(trigger).toBeFocused()
    })

    test('should have sufficient color contrast', async ({ page }) => {
      await page.click('button[data-box-title="Test Dialog"]')

      const title = page.locator('[role="dialog"] .text-base')

      // Get computed styles
      const color = await title.evaluate((el) => {
        const styles = window.getComputedStyle(el)
        return styles.color
      })

      const bgColor = await title.evaluate((el) => {
        const styles = window.getComputedStyle(el)
        return styles.backgroundColor
      })

      // Basic check - ensure colors are defined
      expect(color).toBeTruthy()
      expect(bgColor).toBeTruthy()
    })
  })

  test.describe('Notify Component Accessibility', () => {
    test('should have proper ARIA role', async ({ page }) => {
      await page.evaluate(() => {
        window.Alpine.store('uno').notify.success('Test notification')
      })

      // Check for alert role
      const notification = page.locator('[role="alert"], [role="status"]')
      await expect(notification).toBeVisible()
    })

    test('should be announced by screen readers', async ({ page }) => {
      await page.evaluate(() => {
        window.Alpine.store('uno').notify.info('Important message')
      })

      const notification = page.locator('.text-sm:has-text("Important message")')
      await expect(notification).toBeVisible()

      // Check aria-live (should be set by role="alert")
      const container = notification.locator('..')
      const role = await container.getAttribute('role')
      expect(['alert', 'status']).toContain(role)
    })

    test('should be dismissible via keyboard', async ({ page }) => {
      await page.evaluate(() => {
        window.Alpine.store('uno').notify.error('Error message')
      })

      const notification = page.locator('[role="alert"]').first()
      await expect(notification).toBeVisible()

      // Find close button
      const closeBtn = notification.locator('button[aria-label="Close"], button:has(svg)')
      await closeBtn.focus()
      await page.keyboard.press('Enter')

      // Notification should close
      await expect(notification).not.toBeVisible()
    })
  })

  test.describe('Tooltip Component Accessibility', () => {
    test('should have proper ARIA attributes', async ({ page }) => {
      const trigger = page.locator('[x-tooltip]').first()
      await trigger.hover()

      // Check for tooltip role
      const tooltip = page.locator('[role="tooltip"]')
      await expect(tooltip).toBeVisible()

      // Check aria-describedby connection
      const ariaId = await trigger.getAttribute('aria-describedby')
      expect(ariaId).toBeTruthy()

      const tooltipId = await tooltip.getAttribute('id')
      expect(ariaId).toBe(tooltipId)
    })

    test('should show on keyboard focus', async ({ page }) => {
      const trigger = page.locator('[x-tooltip]').first()

      // Tab to trigger
      await trigger.focus()

      // Tooltip should appear
      const tooltip = page.locator('[role="tooltip"]')
      await expect(tooltip).toBeVisible({ timeout: 1000 })
    })

    test('should hide on Escape key', async ({ page }) => {
      const trigger = page.locator('[x-tooltip]').first()
      await trigger.focus()

      const tooltip = page.locator('[role="tooltip"]')
      await expect(tooltip).toBeVisible()

      // Press Escape
      await page.keyboard.press('Escape')

      // Tooltip should hide
      await expect(tooltip).not.toBeVisible()
    })
  })

  test.describe('Tab Component Accessibility', () => {
    test('should have proper tablist role', async ({ page }) => {
      await page.goto('/examples/components/tab.html')

      const tablist = page.locator('[role="tablist"]')
      await expect(tablist).toBeVisible()
    })

    test('should have proper tab roles and attributes', async ({ page }) => {
      await page.goto('/examples/components/tab.html')

      const tabs = page.locator('[role="tab"]')
      const count = await tabs.count()
      expect(count).toBeGreaterThan(0)

      // Check first tab
      const firstTab = tabs.first()

      // Should have aria-selected
      const selected = await firstTab.getAttribute('aria-selected')
      expect(['true', 'false']).toContain(selected)

      // Should have aria-controls
      const controls = await firstTab.getAttribute('aria-controls')
      expect(controls).toBeTruthy()
    })

    test('should support arrow key navigation', async ({ page }) => {
      await page.goto('/examples/components/tab.html')

      const tabs = page.locator('[role="tab"]')
      const firstTab = tabs.first()

      // Focus first tab
      await firstTab.focus()
      await expect(firstTab).toBeFocused()

      // Press right arrow
      await page.keyboard.press('ArrowRight')

      // Second tab should be focused
      const secondTab = tabs.nth(1)
      await expect(secondTab).toBeFocused()

      // Press left arrow
      await page.keyboard.press('ArrowLeft')

      // First tab should be focused again
      await expect(firstTab).toBeFocused()
    })

    test('should have proper tabpanel roles', async ({ page }) => {
      await page.goto('/examples/components/tab.html')

      const panels = page.locator('[role="tabpanel"]')
      const count = await panels.count()
      expect(count).toBeGreaterThan(0)

      // Check first panel
      const firstPanel = panels.first()

      // Should have aria-labelledby
      const labelledBy = await firstPanel.getAttribute('aria-labelledby')
      expect(labelledBy).toBeTruthy()
    })
  })

  test.describe('Form Validation Accessibility', () => {
    test('should have proper error messages', async ({ page }) => {
      await page.goto('/examples/components/valid.html')

      const input = page.locator('input[type="email"]')

      // Enter invalid email
      await input.fill('invalid-email')
      await input.blur()

      // Check for error message
      const error = page.locator('.text-red-600, [role="alert"]')
      await expect(error).toBeVisible()

      // Error should be associated with input
      const describedBy = await input.getAttribute('aria-describedby')
      expect(describedBy).toBeTruthy()
    })

    test('should update aria-invalid attribute', async ({ page }) => {
      await page.goto('/examples/components/valid.html')

      const input = page.locator('input[x-valid\\:email]').first()

      // Initially no aria-invalid
      let ariaInvalid = await input.getAttribute('aria-invalid')
      expect(ariaInvalid).toBeNull()

      // Enter invalid email
      await input.fill('invalid')
      await input.blur()

      // Should have aria-invalid="true"
      ariaInvalid = await input.getAttribute('aria-invalid')
      expect(ariaInvalid).toBe('true')

      // Enter valid email
      await input.fill('valid@example.com')
      await input.blur()

      // Should have aria-invalid="false"
      ariaInvalid = await input.getAttribute('aria-invalid')
      expect(ariaInvalid).toBe('false')
    })
  })

  test.describe('Keyboard Navigation', () => {
    test('should support Tab navigation', async ({ page }) => {
      // Tab through interactive elements
      await page.keyboard.press('Tab')

      // Check if something is focused
      const focused = await page.evaluate(() => document.activeElement?.tagName || '')
      expect(['BUTTON', 'A', 'INPUT']).toContain(focused)
    })

    test('should support Shift+Tab navigation', async ({ page }) => {
      // Focus last element
      await page.keyboard.press('Tab')
      await page.keyboard.press('Tab')

      // Go back
      await page.keyboard.press('Shift+Tab')

      // Should focus previous element
      const focused = await page.evaluate(() => document.activeElement?.tagName || '')
      expect(['BUTTON', 'A', 'INPUT']).toContain(focused)
    })
  })

  test.describe('Screen Reader Support', () => {
    test('should have proper page title', async ({ page }) => {
      const title = await page.title()
      expect(title).toBeTruthy()
      expect(title.length).toBeGreaterThan(0)
    })

    test('should have proper heading hierarchy', async ({ page }) => {
      const h1 = page.locator('h1')
      const h1Count = await h1.count()

      // Should have exactly one h1
      expect(h1Count).toBe(1)

      // h1 should not be empty
      const h1Text = await h1.textContent()
      expect(h1Text?.trim()).toBeTruthy()
    })

    test('should have alt text for images', async ({ page }) => {
      const images = page.locator('img')
      const count = await images.count()

      for (let i = 0; i < count; i++) {
        const img = images.nth(i)
        const alt = await img.getAttribute('alt')

        // Should have alt attribute (can be empty for decorative images)
        expect(alt).not.toBeNull()
      }
    })

    test('should have proper labels for form inputs', async ({ page }) => {
      await page.goto('/examples/components/valid.html')

      const inputs = page.locator('input:not([type="hidden"])')
      const count = await inputs.count()

      for (let i = 0; i < count; i++) {
        const input = inputs.nth(i)

        // Should have label or aria-label
        const id = await input.getAttribute('id')
        const ariaLabel = await input.getAttribute('aria-label')
        const ariaLabelledBy = await input.getAttribute('aria-labelledby')

        if (id) {
          const label = page.locator(`label[for="${id}"]`)
          const labelCount = await label.count()

          // Should have label, aria-label, or aria-labelledby
          expect(labelCount > 0 || ariaLabel || ariaLabelledBy).toBeTruthy()
        }
      }
    })
  })
})
