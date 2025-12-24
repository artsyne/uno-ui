/**
 * Uno - Tab Component (Refactored for v0.5.0)
 *
 * @module uno.components.tab
 * @description Lightweight tab switching with HTMX integration and Hash navigation.
 */

import { sanitizeHTML } from '../utils/dom'

export default (Alpine) => {
  // Global HTMX interceptor for sanitization
  if (typeof document !== 'undefined') {
    document.body.addEventListener('htmx:beforeSwap', (evt) => {
      // Check if this request is from a tab
      const trigger = evt.detail.requestConfig.elt
      if (trigger && trigger.getAttribute('role') === 'tab') {
        // Sanitize content to prevent XSS
        evt.detail.serverResponse = sanitizeHTML(evt.detail.serverResponse)
      }
    })
  }

  Alpine.data('unoTabs', (defaultTab = '', options = {}) => ({
    activeTab: defaultTab,
    // Registry to store configuration for each tab (url, target, etc.)
    tabs: {},

    init() {
      // Hash synchronization
      this.$nextTick(() => {
        if (options.hash !== false) {
          const checkHash = () => {
            const hash = window.location.hash.slice(1)
            if (hash) {
              this.activeTab = hash
              this.ensureLoad(hash)
              return
            }

            if (this.activeTab) {
              this.ensureLoad(this.activeTab)
            }
          }

          checkHash()
          window.addEventListener('hashchange', checkHash)
        } else {
          // If hash is disabled, load the default tab if needed
          if (this.activeTab) {
            this.ensureLoad(this.activeTab)
          }
        }
      })
    },

    isActive(tab) {
      return this.activeTab === tab
    },

    /**
     * Load tab content if it's an AJAX tab.
     * Uses htmx.ajax() directly to avoid DOM attribute binding timing issues.
     */
    load(tab) {
      const config = this.tabs[tab]
      if (!config || !config.url) return

      // Find the tab element to use as source (for hx-indicator etc.)
      const sourceEl = this.$root.querySelector(`[role="tab"][data-tab="${tab}"]`)

      // Determine target: use config.target or default to #panel-{tab}
      const targetSelector = config.target || `#panel-${tab}`
      const targetEl = document.querySelector(targetSelector)

      if (typeof htmx !== 'undefined' && targetEl) {
        // Use htmx.ajax API directly
        // https://htmx.org/api/#ajax
        htmx.ajax('GET', config.url, {
          target: targetSelector,
          swap: config.swap || 'innerHTML',
          source: sourceEl // Allows hx-indicator to work if placed on the tab button
        })
      }
    },

    ensureLoad(tab) {
      const attemptLoad = (attempt = 0) => {
        const config = this.tabs[tab]

        if (config) {
          if (!config.url) return
          this.load(tab)
          return
        }

        if (attempt >= 20) return

        if (typeof requestAnimationFrame !== 'undefined') {
          requestAnimationFrame(() => attemptLoad(attempt + 1))
          return
        }

        setTimeout(() => attemptLoad(attempt + 1), 16)
      }

      attemptLoad(0)
      // Start attempt
      attemptLoad()
    },

    open(tab, openOptions = {}) {
      const { updateHash = true } = openOptions
      this.activeTab = tab
      this.ensureLoad(tab)
      // Trigger load for AJAX tabs
      this.load(tab)

      if (options.hash !== false && updateHash && tab) {
        // Update URL hash without jumping
        if (window.location.hash.slice(1) !== tab) {
          window.history.replaceState(null, null, `#${tab}`)
        }
      }
    },
      // Register tab configuration
      // This runs during x-bind evaluation
      this.tabs[tab] = options


    trigger(tab, options = {}) {
      // Get config for classes
      const opts = (typeof window !== 'undefined' && window.Uno?.options?.ui?.tab) || {}

      const activeClass = opts.triggerActive || 'is-active'
      const inactiveClass = opts.triggerInactive || ''

      const bindings = {
        role: 'tab',
        'data-tab': tab, // Critical for locating the element later
        ':aria-selected': `isActive('${tab}').toString()`,
        ':class': `isActive('${tab}') ? '${activeClass}' : '${inactiveClass}'`,
        '@click.prevent': `open('${tab}')`
      // We still add data attributes for debugging or other scripts,
      // but the primary logic is now driven by JS state.

      // Add data attributes for AJAX
      if (options.url) {
        bindings['data-hx-url'] = options.url
        bindings['data-hx-target'] = options.target || `#panel-${tab}`
      }

      return bindings
    },

    tablist: {
      role: 'tablist'
    },

    panel(tab) {
        // We do NOT set ID here to avoid overwriting user-defined IDs (crucial for hx-target)
      return {
        'x-show': `isActive('${tab}')`,
        role: 'tabpanel',
        'aria-labelledby': `tab-${tab}`
      }
    }
  }))
}
