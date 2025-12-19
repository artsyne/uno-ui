import { test, expect } from '@playwright/test'
import { mockCdnAssets } from './helpers/routes'

test('language switching updates texts', async ({ page }) => {
  await mockCdnAssets(page)

  await page.goto('/examples/i18n.html')
  const title = page.getByTestId('title')
  const desc = page.getByTestId('desc')
  await expect(title).toBeVisible()
  await expect(desc).toBeVisible()

  await page.getByTestId('btn-zh').click()
  await expect(title).toHaveText('Uno i18n 演示')

  await page.getByTestId('btn-zh-tw').click()
  await expect(title).toHaveText('Uno i18n 演示')

  await page.getByTestId('btn-en').click()
  await expect(title).toHaveText('Uno i18n Demo')
})
