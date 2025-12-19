/**
 * Uno - DOM Manager
 * Handles centralized DOM injection with CSP support.
 *
 * @module uno.core.dom
 */

import { Uno } from './uno'
import { log } from '../utils/log'

export const domManager = {
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
      if (Uno.options.csp?.nonce) {
        style.setAttribute('nonce', Uno.options.csp.nonce)
      }

      style.textContent = css
      document.head.appendChild(style)
    } catch (e) {
      log.error(`DomManager: Failed to inject styles #${id}`, e)
    }
  }
}
