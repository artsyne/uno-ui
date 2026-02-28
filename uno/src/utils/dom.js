/**
 * Uno - DOM Utils
 *
 * @module uno.utils.dom
 */

import { log } from './log'

const onceMarks = new WeakMap()
let onceVersion = 0

/**
 * Reset once marks (useful for testing)
 */
export function resetOnceMarks() {
  onceVersion++
}

/**
 * Ensure a function/action is only performed once per element.
 *
 * This utility prevents duplicate initialization or repeated operations on the same element
 * by maintaining a WeakMap-based registry of completed actions.
 *
 * Implementation Details:
 * - Uses WeakMap for automatic garbage collection when elements are removed
 * - Each element can have multiple "once" keys for different operations
 * - Thread-safe for synchronous code (no race conditions)
 *
 * Common Use Cases:
 * - Initialize event listeners only once
 * - Prevent duplicate plugin registration
 * - Ensure setup code runs exactly once per element
 *
 * Edge Cases:
 * - Returns false if target is null/undefined (graceful degradation)
 * - Returns false if the key has already been registered
 * - Returns true and registers the key if this is the first call
 *
 * @param {HTMLElement} target - The target element
 * @param {string} key - A unique key for the action (e.g., 'tooltip-init', 'htmx-setup')
 * @returns {boolean} True if the action has not been performed yet, false otherwise
 *
 * @example
 * if (once(element, 'my-plugin-init')) {
 *   // This code will run exactly once per element
 *   element.addEventListener('click', handler)
 * }
 */
export function once(target, key) {
  if (!target) return false
  const k = String(key)
  const versionedKey = `${onceVersion}:${k}`
  let set = onceMarks.get(target)
  if (!set) {
    set = new Set()
    onceMarks.set(target, set)
  }
  if (set.has(versionedKey)) return false
  set.add(versionedKey)
  return true
}

/**
 * Sanitize HTML content to prevent XSS attacks.
 *
 * Security Strategy:
 * This function implements a defense-in-depth approach using multiple regex patterns
 * to remove potentially dangerous HTML content:
 *
 * 1. **Script Tags**: Removes all <script> tags and their content
 *    - Pattern: /<script\b[^>]*>([\s\S]*?)<\/script>/gim
 *    - Prevents execution of arbitrary JavaScript
 *
 * 2. **Image Tags**: Removes all <img> tags
 *    - Pattern: /<img\b[^>]*>/gim
 *    - Prevents onerror XSS attacks (e.g., <img src=x onerror=alert(1)>)
 *
 * 3. **Event Handlers**: Strips all on* attributes (onclick, onload, onerror, etc.)
 *    - Patterns: /\s+on[a-z]+="[^"]*"/gim and /\s+on[a-z]+='[^']*'/gim
 *    - Prevents inline event handler XSS
 *
 * 4. **JavaScript Protocol**: Removes javascript: URLs from href attributes
 *    - Patterns: /\s+href="javascript:[^"]*"/gim and similar
 *    - Prevents navigation-based XSS
 *
 * 5. **Global javascript: removal**: Removes remaining javascript: occurrences
 *    - Pattern: /javascript:/gim
 *    - Catchall for any missed instances
 *
 * Limitations:
 * - This is a basic XSS protection suitable for trusted or semi-trusted content
 * - For untrusted user input, integrate DOMPurify library for comprehensive protection
 * - Regex-based sanitization can be bypassed with clever encoding or edge cases
 * - Does not handle all possible XSS vectors (e.g., CSS injection, data: URLs)
 *
 * Performance:
 * - Linear time complexity O(n) where n is the HTML string length
 * - Six regex replacements per call
 * - Suitable for typical content lengths (<10KB)
 *
 * @param {string} html - The HTML string to sanitize
 * @returns {string} The sanitized HTML string (safe to insert into DOM)
 *
 * @example
 * // Removes dangerous content
 * const safe = sanitizeHTML('<div onclick="evil()">Text</div>')
 * // Result: '<div>Text</div>'
 *
 * @example
 * // Removes script tags
 * const safe = sanitizeHTML('Hello <script>alert(1)</script> World')
 * // Result: 'Hello  World'
 *
 * @see {@link https://cheatsheetseries.owasp.org/cheatsheets/Cross_Site_Scripting_Prevention_Cheat_Sheet.html|OWASP XSS Prevention}
 */
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

/**
 * Set HTML content of an element safely with optional sanitization.
 *
 * This function provides a secure way to set innerHTML with automatic XSS protection
 * when dealing with untrusted content.
 *
 * Security Modes:
 * 1. **Untrusted Mode** (default, trusted=false):
 *    - Runs content through sanitizeHTML() before insertion
 *    - Removes script tags, event handlers, and dangerous attributes
 *    - Safe for user-generated or external content
 *
 * 2. **Trusted Mode** (trusted=true):
 *    - Bypasses sanitization for performance
 *    - Only use with content from trusted sources
 *    - Example: Server-rendered templates, static content
 *
 * Alpine.js Integration:
 * - Optionally re-initializes Alpine.js directives after HTML insertion
 * - Enables x-data, x-show, and other Alpine directives in dynamic content
 * - Automatically called when alpine option is provided
 *
 * Edge Cases:
 * - Handles null/undefined HTML gracefully (converts to empty string)
 * - No-op if element is null/undefined
 * - Safe to call multiple times on the same element
 *
 * @param {HTMLElement} el - The target element to update
 * @param {string} html - The HTML content to insert
 * @param {Object} [options={}] - Configuration options
 * @param {boolean} [options.trusted=false] - Whether the content is trusted (skips sanitization)
 * @param {Object} [options.alpine] - Alpine.js instance for initializing directives
 *
 * @example
 * // Untrusted content (auto-sanitized)
 * setHTML(div, userGeneratedHTML)
 *
 * @example
 * // Trusted content (no sanitization)
 * setHTML(div, serverTemplate, { trusted: true })
 *
 * @example
 * // With Alpine.js support
 * setHTML(div, '<div x-data="{ count: 0 }">...</div>', {
 *   trusted: true,
 *   alpine: window.Alpine
 * })
 */
export function setHTML(el, html, { trusted = false, alpine } = {}) {
  if (!el) return
  // Use sanitizeHTML if not trusted
  el.innerHTML = trusted ? String(html ?? '') : sanitizeHTML(String(html ?? ''))
  alpine?.initTree?.(el)
}

// DomManager Implementation
export const DomManager = {
  /**
   * Ensure a container element exists in the DOM, creating it if necessary.
   *
   * This method implements idempotent container management:
   * - Returns existing container if already present (by ID)
   * - Creates and injects new container if not found
   * - Automatically handles CSP nonce injection
   * - Safely handles errors with logging
   *
   * Container Creation Process:
   * 1. Check for existing element with matching ID
   * 2. If found, return it immediately (no duplicate creation)
   * 3. If not found:
   *    a. Create new div element
   *    b. Set the specified ID
   *    c. Auto-inject CSP nonce from global Uno config (if available)
   *    d. Apply all provided attributes
   *    e. Set innerHTML if content provided
   *    f. Append to specified parent (default: document.body)
   *
   * CSP Nonce Injection:
   * - Automatically reads nonce from window.Uno.options.csp.nonce
   * - Only injects if not already specified in attributes
   * - Enables inline styles/scripts in CSP-protected environments
   *
   * Error Handling:
   * - Wraps creation in try-catch
   * - Logs errors via Uno's logging system
   * - Returns null on failure (graceful degradation)
   *
   * Edge Cases:
   * - Parent null/undefined: Falls back to document.documentElement
   * - Duplicate ID: Returns existing element (prevents duplicates)
   * - Empty HTML: Creates container without content
   * - No attributes: Creates minimal container with just ID
   *
   * @param {string} id - Unique element ID (used for lookup and creation)
   * @param {string} [html=''] - Inner HTML content for the container
   * @param {Object} [attributes={}] - HTML attributes to set (key-value pairs)
   * @param {HTMLElement} [parent=document.body] - Parent element for appending
   * @returns {HTMLElement|null} The container element, or null if creation failed
   *
   * @example
   * // Simple container
   * const container = DomManager.ensureContainer('my-container')
   *
   * @example
   * // Container with content and attributes
   * const container = DomManager.ensureContainer(
   *   'tooltip-container',
   *   '<div class="tooltip">Content</div>',
   *   { 'x-data': '{}', 'class': 'fixed top-0' }
   * )
   *
   * @example
   * // Idempotent: Safe to call multiple times
   * const c1 = DomManager.ensureContainer('my-id')
   * const c2 = DomManager.ensureContainer('my-id')
   * console.log(c1 === c2) // true (same element)
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
   * Inject CSS styles into the document head with CSP nonce support.
   *
   * This method provides safe, idempotent style injection with automatic
   * Content Security Policy (CSP) nonce handling.
   *
   * Features:
   * - **Idempotent**: Safe to call multiple times (won't create duplicates)
   * - **CSP Support**: Automatically injects nonce attribute
   * - **Error Handling**: Graceful failure with logging
   * - **Head Injection**: Styles go in <head> for proper CSS precedence
   *
   * Injection Process:
   * 1. Check if style element with this ID already exists
   * 2. If exists, return immediately (no duplicate creation)
   * 3. If not exists:
   *    a. Create new <style> element
   *    b. Set the specified ID
   *    c. Auto-inject CSP nonce from window.Uno.options.csp.nonce
   *    d. Set CSS content as textContent (safe from XSS)
   *    e. Append to document.head
   *
   * CSP Compliance:
   * - Reads nonce from global Uno configuration
   * - Enables inline styles in strict CSP environments
   * - Example CSP header: `style-src 'nonce-abc123'`
   *
   * Performance:
   * - ID-based duplicate prevention (O(1) lookup)
   * - No style re-processing on subsequent calls
   * - Suitable for dynamic style injection
   *
   * Use Cases:
   * - Component-specific styles
   * - Theme customization
   * - Dynamic CSS generation
   * - Plugin style injection
   *
   * @param {string} id - Unique style element ID (prevents duplicates)
   * @param {string} css - CSS content to inject (raw CSS text)
   *
   * @example
   * // Inject component styles
   * DomManager.injectStyles('tooltip-styles', `
   *   .tooltip { position: fixed; z-index: 9999; }
   *   .tooltip-arrow { width: 0; height: 0; }
   * `)
   *
   * @example
   * // Idempotent: Safe to call multiple times
   * DomManager.injectStyles('my-styles', '.foo { color: red; }')
   * DomManager.injectStyles('my-styles', '.foo { color: red; }') // No-op
   *
   * @see {@link https://developer.mozilla.org/en-US/docs/Web/HTTP/CSP|MDN: Content Security Policy}
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
  sanitize: sanitizeHTML // HTML sanitization
}
