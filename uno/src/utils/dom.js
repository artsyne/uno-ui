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

export function sanitizeHTML(html, options = {}) {
  const {
    allowedTags = [
      'b',
      'i',
      'strong',
      'em',
      'u',
      'br',
      'p',
      'ul',
      'ol',
      'li',
      'code',
      'pre',
      'a',
      'span',
      'div'
    ],
    allowedAttrs = { a: ['href', 'target', 'rel'], '*': ['class'] }
  } = options

  const raw = String(html ?? '')
  const preStripped = raw
    .replace(/<(script|style|iframe|object|embed)\b[\s\S]*?<\/\1\s*>/gi, '')
    .replace(/<(video|audio|svg|math)\b[\s\S]*?<\/\1\s*>/gi, '')
    .replace(/<(img|link|meta|base|source|track)\b[^>]*\/?>/gi, '')

  const template = document.createElement('template')
  template.innerHTML = preStripped

  const walker = document.createTreeWalker(template.content, NodeFilter.SHOW_ELEMENT)
  const toRemove = []
  const toStrip = [] // Tags to remove completely including content

  while (walker.nextNode()) {
    const el = walker.currentNode
    const tag = el.tagName.toLowerCase()

    // Tags to strip completely (remove element AND content)
    if (['script', 'style', 'iframe', 'object', 'embed'].includes(tag)) {
      toStrip.push(el)
      continue
    }

    // Default whitelist strategy. Non-whitelisted tags are downgraded to plain text.
    if (!allowedTags.includes(tag)) {
      toRemove.push(el)
      continue
    }

    for (const attr of [...el.attributes]) {
      const name = attr.name.toLowerCase()
      const allow = (allowedAttrs[tag] || []).concat(allowedAttrs['*'] || [])
      if (!allow.includes(name)) el.removeAttribute(attr.name)
    }

    // Only allow safe protocols/internal links for href to prevent javascript: injection.
    if (tag === 'a') {
      const href = (el.getAttribute('href') || '').trim()
      const isSafe =
        href.startsWith('/') ||
        href.startsWith('https://') ||
        href.startsWith('http://') ||
        href.startsWith('#')
      if (!isSafe) el.removeAttribute('href')
      if (!el.getAttribute('rel')) el.setAttribute('rel', 'noopener noreferrer')
    }
  }

  // Remove unsafe tags completely (replace with empty text)
  for (const node of toStrip) {
    node.remove()
  }

  // Downgrade non-whitelisted tags to text (keep content)
  for (const node of toRemove) {
    const text = node.textContent || ''
    if (text) {
      node.replaceWith(document.createTextNode(text))
    } else {
      node.remove()
    }
  }

  return template.innerHTML
}

export function setHTML(el, html, { trusted = false, alpine } = {}) {
  if (!el) return
  el.innerHTML = trusted ? String(html ?? '') : sanitizeHTML(html)
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
    if (document.getElementById(id)) return

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
  sanitize: sanitizeHTML, // HTML sanitization
  escape: escapeHTML // HTML escaping
}
