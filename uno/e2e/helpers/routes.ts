import type { Page } from '@playwright/test'

export async function mockCdnAssets(page: Page) {
  await page.route('**/dist/uno.umd.js*', async (route) => {
    await route.fulfill({ path: 'dist/uno.umd.js' })
  })
  await page.route('**/alpinejs@*/dist/cdn.min.js*', async (route) => {
    await route.fulfill({ path: 'node_modules/alpinejs/dist/cdn.min.js' })
  })
  await page.route('**/alpinejs@3.x.x*', async (route) => {
    await route.fulfill({ path: 'node_modules/alpinejs/dist/cdn.min.js' })
  })
  await page.route('**/htmx.org@*/dist/htmx.min.js*', async (route) => {
    await route.fulfill({ path: 'node_modules/htmx.org/dist/htmx.min.js' })
  })
  await page.route('**/htmx.org@*', async (route) => {
    await route.fulfill({ path: 'node_modules/htmx.org/dist/htmx.min.js' })
  })
}
