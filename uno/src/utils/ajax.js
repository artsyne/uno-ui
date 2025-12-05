/**
 * Uno - Ajax Utils
 *
 * @module uno.core.ajax
 * @version 0.1.0
 * @license MIT
 */

/**
 * Fetch content from URL (returns JSON or text)
 * @param {string} url - Request URL
 * @returns {Promise<any>} Response data
 */
export async function get(url) {
  const res = await fetch(url)
  const ct = res.headers.get('content-type') || ''
  if (ct.includes('application/json')) return res.json()
  return res.text()
}
