/**
 * Uno - DOM Utils
 *
 * @module uno.utils.dom
 */

import { log } from './log'

export function escapeHTML(input) {
  return String(input ?? '')
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&#39;')
}

const onceMarks = new WeakMap()

export function once(target, key) {
  if (!target) return false
  const k = String(key)
  let set = onceMarks.get(target)
  if (!set) {
    set = new Set()
    onceMarks.set(target, set)
  }
  if (set.has(k)) return false
  set.add(k)
  return true
}

export function sanitizeHTML(html) {
  if (!html) return ''
  // Basic XSS protection: strip script tags and on* attributes
  // This is a minimal implementation. For production with untrusted input, use DOMPurify.
  return html
    .replace(/<script\b[^>]*>([\s\S]*?)<\/script>/gim, '')
    .replace(/<img\b[^>]*>/gim, '')
    .replace(/\s+on[a-z]+="[^"]*"/gim, '')
    .replace(/\s+on[a-z]+='[^']*'/gim, '')
    .replace(/\s+href="javascript:[^"]*"/gim, '')
    .replace(/\s+href='javascript:[^']*'/gim, '')
    .replace(/javascript:/gim, '')
}

export function setHTML(el, html, { trusted = false, alpine } = {}) {
  if (!el) return
  // Use sanitizeHTML if not trusted
  el.innerHTML = trusted ? String(html ?? '') : sanitizeHTML(String(html ?? ''))
  alpine?.initTree?.(el)
}

// DomManager Implementation
export const DomManager = {
  /**
   * Ensure a container element exists in the DOM.
   * @param {string} id - Element ID
   * @param {string} html - Inner HTML content
   * @param {Object} attributes - Attributes to set
   * @param {HTMLElement} parent - Parent element (default: document.body)
   */
  ensureContainer(id, html = '', attributes = {}, parent = document.body) {
    // Return existing container if found
    const existing = document.getElementById(id)
    if (existing) return existing

    try {
      const container = document.createElement('div')
      container.id = id

      // Automatically inject CSP nonce
      if (window.Uno?.options?.csp?.nonce && !attributes.nonce) {
        attributes.nonce = window.Uno.options.csp.nonce
      }

      // Set attributes
      Object.entries(attributes).forEach(([key, value]) => {
        container.setAttribute(key, value)
      })

      // Set content
      if (html) {
        container.innerHTML = html
      }

      // Append
      if (parent) {
        parent.appendChild(container)
      } else {
        // Fallback if body not ready (should be handled by caller usually)
        document.documentElement.appendChild(container)
      }
      return container
    } catch (e) {
      log.error(`DomManager: Failed to inject container #${id}`, e)
      return null
    }
  },

  /**
   * Inject CSS styles with CSP nonce support.
   * @param {string} id - Style element ID
   * @param {string} css - CSS content
   */
  injectStyles(id, css) {
    if (document.getElementById(id)) return

    try {
      const style = document.createElement('style')
      style.id = id

      // CSP Nonce
      if (window.Uno?.options?.csp?.nonce) {
        style.setAttribute('nonce', window.Uno.options.csp.nonce)
      }

      style.textContent = css
      document.head.appendChild(style)
    } catch (e) {
      log.error(`DomManager: Failed to inject styles #${id}`, e)
    }
  },

  // Exports
  setHTML: setHTML, // Secure HTML setting
  escape: escapeHTML, // HTML escaping
  sanitize: sanitizeHTML // HTML sanitization
}
