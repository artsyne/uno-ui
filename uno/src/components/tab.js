/**
 * Uno - Tab Component (Store-First, Directive-Driven Architecture)
 *
 * A lightweight, declarative, HTMX-native tab component.
 * Follows the "Store-First, Directive-Driven" pattern.
 *
 * Features:
 * - Store-based state management
 * - Declarative x-uno-tab directive
 * - HTMX integration
 * - Hash synchronization
 * - Accessibility (ARIA)
 *
 * @module uno.components.tab
 * @version 0.6.0
 * @license MIT
 */

import { Uno } from '../core/uno'
import { registerPlugin } from '../core/plugin-registry'
import { log } from '../utils/log'

export default (Alpine, options = {}) => {
  // Register plugin defaults
  const pluginDefaults = {
    defaultTab: '',
    hash: true,
    activeClass: 'is-active',
    inactiveClass: ''
  }
  Uno.registerPluginDefaults('tab', { ...pluginDefaults, ...options.plugins?.tab, ...options.tab })

  // ============================================================
  // Tab Store
  // Registered under uno namespace: Alpine.store('uno').tab
  // ============================================================
  const tabStore = {
    // State: { [group]: activeTabName }
    groups: {},
    // Registry: { [tabName]: { group, name, ... } }
    registry: {},

    init() {
      // Hash synchronization
      const checkHash = () => {
        let hash = window.location.hash.slice(1)

        // Fallback for some environments where location.hash might be empty initially
        if (!hash && window.location.href.includes('#')) {
          hash = window.location.href.split('#')[1]
        }

        log.info(`[Tab] checkHash: ${hash}`)

        if (hash && this.registry[hash]) {
          const { group, name } = this.registry[hash]
          // If hash is found in registry, activate it
          log.info(`[Tab] checkHash match: ${name}`)
          this.activate({ group, name, updateHash: false })
        }
      }

      // Initial check
      checkHash()

      // Listen for hash changes
      window.addEventListener('hashchange', checkHash)
      // Safety check on load (in case hash was not ready during init/register)
      window.addEventListener('load', checkHash)
    },

    /**
     * Activate a tab
     * @param {object} instanceOpts - Options { name, group, ... }
     */
    activate(instanceOpts = {}) {
      // 4-level merge via Uno.getConfig
      const config = Uno.getConfig('tab', {
        group: 'default',
        hash: true,
        ...instanceOpts
      })

      if (!config.name) return

      log.info(
        `[Tab] Activating ${config.name} in ${config.group}. Hash update: ${config.updateHash !== false}`
      )

      // Register if not known (lazy registration)
      if (!this.registry[config.name]) {
        this.registry[config.name] = { group: config.group, name: config.name }
      }

      // Update state
      this.groups[config.group] = config.name

      // Handle Hash
      if (config.hash && config.updateHash !== false) {
        if (window.location.hash.slice(1) !== config.name) {
          window.history.replaceState(null, null, `#${config.name}`)
        }
      }

      // Trigger HTMX if needed (Simulate click if hx-get exists)
      // This is crucial for deep linking where activation happens via JS, not user click.
      if (config.updateHash === false) {
        // Typically means activated from hash or init
        // Use setTimeout to allow DOM to be ready if called during init
        setTimeout(() => {
          const tabEl = document.getElementById(`tab-${config.name}`)
          if (tabEl && tabEl.hasAttribute('hx-get')) {
            log.info(`[Tab] Triggering HTMX for ${config.name}`)
            if (window.htmx) {
              window.htmx.trigger(tabEl, 'click')
            }
          }
        }, 0)
      }
    },

    /**
     * Check if a tab is active
     * @param {string} name - Tab name
     * @param {string} group - Group name
     * @returns {boolean}
     */
    isActive(name, group = 'default') {
      const active = this.groups[group] === name
      if (active) {
        // log.info(`[Tab] isActive true for ${name} in ${group}`)
      }
      return active
    },

    /**
     * Register a tab (helper for directives or manual setup)
     * @param {string} name
     * @param {string} group
     */
    register(name, group = 'default') {
      log.info(`[Tab] Registering ${name} in ${group}. Hash: ${window.location.hash}`)
      this.registry[name] = { name, group }
      // If no active tab for this group, set this one
      if (!this.groups[group]) {
        this.groups[group] = name
      }

      // Check if this new tab matches the current hash
      // This handles deep linking when tabs are registered after store init
      let hash = window.location.hash.slice(1)
      if (!hash && window.location.href.includes('#')) {
        hash = window.location.href.split('#')[1]
      }

      if (hash === name) {
        log.info(`[Tab] Hash match for ${name}. Activating.`)
        this.activate({ group, name, updateHash: false })
      }
    }
  } // Register store
  registerPlugin(Alpine, 'tab', tabStore)

  // Explicitly initialize the store since Alpine doesn't auto-init nested stores
  // created via registerPlugin (which uses spread syntax or property assignment)
  Alpine.store('uno').tab.init()

  // ============================================================
  // Tab Directive (x-uno-tab)  // Usage: <button x-uno-tab="{ name: 'tab1' }">
  // ============================================================
  Alpine.directive('uno-tab', (el, { expression }, { evaluateLater, cleanup }) => {
    if (!expression) {
      console.warn('[Uno] x-uno-tab directive requires a configuration object.')
      return
    }

    const getExpression = evaluateLater(expression)

    // Initial Registration
    getExpression((value) => {
      let config = {}
      if (typeof value === 'object' && value !== null) {
        config = value
      } else if (typeof value === 'string') {
        config = { name: value }
      } else {
        console.warn('[Uno] x-uno-tab directive value must be a string or object', value)
        return
      }

      // Default group if not specified
      if (!config.group) config.group = 'default'

      // Register on init
      if (config.name) {
        Alpine.store('uno').tab.register(config.name, config.group)
      }
    })

    const onClick = (e) => {
      e.preventDefault()

      getExpression((value) => {
        let config = {}
        if (typeof value === 'object' && value !== null) {
          config = value
        } else if (typeof value === 'string') {
          config = { name: value }
        }

        if (config.name) {
          Alpine.store('uno').tab.activate(config)
        }
      })
    }

    el.addEventListener('click', onClick)
    cleanup(() => el.removeEventListener('click', onClick))
  })
}
