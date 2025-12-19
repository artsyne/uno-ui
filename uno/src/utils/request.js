/**
 * Uno - Request Utils
 *
 * @module uno.utils.request
 */

class LRUCache {
  constructor(maxSize = 50) {
    this.maxSize = maxSize
    this.cache = new Map()
  }

  get(key) {
    if (!this.cache.has(key)) return undefined
    const value = this.cache.get(key)
    // Update access order: delete and re-insert
    this.cache.delete(key)
    this.cache.set(key, value)
    return value
  }

  set(key, value) {
    if (this.cache.has(key)) {
      this.cache.delete(key)
    }
    this.cache.set(key, value)

    // Delete oldest if capacity exceeded
    if (this.cache.size > this.maxSize) {
      const firstKey = this.cache.keys().next().value
      this.cache.delete(firstKey)
    }
  }

  has(key) {
    return this.cache.has(key)
  }

  clear() {
    this.cache.clear()
  }
}

const requestCache = new LRUCache(50)

/**
 * Clear request cache
 */
export function clearCache() {
  requestCache.clear()
}

function parseHXTrigger(value) {
  const raw = String(value ?? '').trim()
  if (!raw) return []

  try {
    const json = JSON.parse(raw)
    if (typeof json === 'string') return [[json, {}]]
    if (json && typeof json === 'object') return Object.entries(json)
  } catch {
    // ignore
  }

  return raw
    .split(',')
    .map((s) => s.trim())
    .filter(Boolean)
    .map((name) => [name, {}])
}

function dispatchTrigger(target, name, detail) {
  const t = target || document.body || document.documentElement
  if (!t) return

  // Prefer htmx.trigger if available to ensure consistency with HTMX ecosystem.
  if (typeof window !== 'undefined' && window.htmx?.trigger) {
    window.htmx.trigger(t, name, detail)
    return
  }

  t.dispatchEvent(new CustomEvent(name, { bubbles: true, detail }))
}

function applyHXTriggerHeaders(res, triggerTarget) {
  const hxTrigger = res.headers.get('HX-Trigger')
  if (!hxTrigger) return

  for (const [name, detail] of parseHXTrigger(hxTrigger)) {
    dispatchTrigger(triggerTarget, name, detail)
  }
}

export async function requestHTML(url, options = {}) {
  const {
    method = 'GET',
    htmxTarget,
    swap = 'innerHTML',
    htmxConfig = {},
    cache = false
  } = options

  // Cache check for GET
  if (cache && method === 'GET' && requestCache.has(url)) {
    return requestCache.get(url)
  }

  if (typeof window !== 'undefined' && window.htmx?.ajax) {
    const temp = htmxTarget || document.createElement('div')
    if (!htmxTarget) temp.style.display = 'none'

    if (!htmxTarget) {
      if (document.body) document.body.appendChild(temp)
      else document.documentElement.appendChild(temp)
    }

    try {
      await window.htmx.ajax(method, url, { target: temp, swap, ...htmxConfig })
      const html = temp.innerHTML
      if (cache && method === 'GET') requestCache.set(url, html)
      return html
    } finally {
      if (!htmxTarget) temp.remove()
    }
  }

  return requestText(url, options)
}

export async function requestData(url, options = {}) {
  const {
    method = 'GET',
    headers = {},
    body,
    credentials = 'same-origin',
    triggerTarget = typeof document !== 'undefined' ? document.body : null,
    cache = false,
    signal // AbortSignal
  } = options

  // Cache check for GET
  if (cache && method === 'GET' && requestCache.has(url)) {
    // Return cached promise/result?
    // Store result to avoid stale promise issues
    return requestCache.get(url)
  }

  // Unify HTMX-style request headers
  const requestHeaders = {
    'HX-Request': 'true',
    'X-Requested-With': 'XMLHttpRequest',
    ...headers
  }

  const res = await fetch(url, { method, headers: requestHeaders, body, credentials, signal })
  applyHXTriggerHeaders(res, triggerTarget)

  if (!res.ok) {
    const error = new Error(`Request failed: ${res.status}`)
    error.status = res.status
    throw error
  }

  const ct = res.headers.get('content-type') || ''
  let data
  if (ct.includes('application/json')) {
    data = await res.json()
  } else {
    data = await res.text()
  }

  if (cache && method === 'GET') {
    requestCache.set(url, data)
  }

  return data
}

export async function requestText(url, options = {}) {
  const data = await requestData(url, options)
  return typeof data === 'string' ? data : JSON.stringify(data, null, 2)
}
