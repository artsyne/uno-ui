/**
 * Uno - Core Instance
 *
 * Manages plugins, initialization, and global configuration.
 *
 * @module uno.core.instance
 * @version 0.3.0
 */

import Alpine from 'alpinejs'
import { log } from '../utils/log'

function isPlainObject(value) {
  return Object.prototype.toString.call(value) === '[object Object]'
}

function mergeDeep(base, patch) {
  if (!isPlainObject(base)) return isPlainObject(patch) ? { ...patch } : patch
  if (!isPlainObject(patch)) return patch === undefined ? { ...base } : patch

  const out = { ...base }
  for (const [key, value] of Object.entries(patch)) {
    if (isPlainObject(value) && isPlainObject(out[key])) {
      out[key] = mergeDeep(out[key], value)
      continue
    }
    if (isPlainObject(value)) {
      out[key] = mergeDeep({}, value)
      continue
    }
    out[key] = value
  }
  return out
}

class UnoInstance {
  constructor() {
    this.version = '0.3.0'
    this.plugins = []
    this.initialized = false
    this.options = {
      debug: false,
      csp: {
        nonce: null // For Content Security Policy
      },
      onError: null, // Custom error handler
      plugins: {},
      ui: {
        base: { rounded: 'rounded-lg', primary: 'bg-blue-600', shadow: 'shadow-md' },
        layers: { boxBase: 10000, tooltipBase: 20000, notifyBase: 30000 },
        box: {},
        notify: {},
        tooltip: {},
        tab: {},
        valid: {}
      }
    }
  }

  _getUiValue(comp, elem) {
    const root = this.options.ui?.[comp]
    if (!root) return undefined
    if (typeof elem !== 'string' || elem.indexOf('.') === -1) return root?.[elem]

    const parts = elem.split('.').filter(Boolean)
    let cur = root
    for (const part of parts) {
      if (cur === null || cur === undefined) return undefined
      cur = cur[part]
    }
    return cur
  }

  /**
   * Get style class from configuration or fallback
   * Support "+" prefix to append class to default
   * @param {string} comp - Component name
   * @param {string} elem - Element name
   * @param {string} defaultCls - Default class string
   * @returns {string} Final class string
   */
  getStyle(comp, elem, defaultCls) {
    const config = this._getUiValue(comp, elem)

    if (config === undefined) return defaultCls

    // Check for append mode ("+class")
    if (typeof config === 'string' && config.trim().indexOf('+') === 0) {
      return `${defaultCls} ${config.trim().slice(1)}`.trim()
    }

    return config
  }

  /**
   * Set global configuration (before or after init)
   * @param {Object} config - Global configuration object
   */
  setConfig(config) {
    this.options = mergeDeep(this.options, config)

    if (this.initialized) {
      log.info('Uno configuration updated', this.options)
    }
    return this
  }

  /**
   * Register a plugin
   * @param {Function} plugin - Plugin function (Alpine, options) => void
   * @param {Object} options - Plugin options
   * @returns {UnoInstance}
   */
  use(plugin, options = {}) {
    if (this.initialized) {
      log.warn('Uno.use() called after initialization. Plugin might not work correctly.')
    }
    this.plugins.push({ plugin, options })
    return this
  }

  /**
   * Initialize Uno
   * @param {Object|Function} alpineOrOptions - Alpine instance or Configuration options
   * @param {Object} options - Configuration options (if first arg is Alpine)
   * @returns {Object} Alpine instance
   */
  init(alpineOrOptions = {}, options = {}) {
    if (this.initialized) {
      log.warn('Uno has already been initialized.')
      return Alpine
    }

    let alpine = Alpine
    let opts = options

    // Robust fallback: if Alpine import is undefined (e.g. UMD external issue), try window.Alpine
    if (!alpine && typeof window !== 'undefined' && window.Alpine) {
      alpine = window.Alpine
    }

    // Detect if first arg is Alpine instance (duck typing)
    if (
      alpineOrOptions &&
      (alpineOrOptions.version ||
        typeof alpineOrOptions.start === 'function' ||
        typeof alpineOrOptions.store === 'function')
    ) {
      alpine = alpineOrOptions
    } else {
      opts = alpineOrOptions || options
    }

    this.options = mergeDeep(this.options, opts)

    // Configure Logging
    log.setDebug(this.options.debug !== false)
    if (this.options.onError) {
      log.setReporter(this.options.onError)
    }

    log.info(`Uno v${this.version} initializing...`)

    // Register plugins
    this.plugins.forEach(({ plugin, options: pluginOptions }) => {
      try {
        // Merge global options with plugin-specific options from .use()
        const mergedOptions = mergeDeep(this.options, pluginOptions)
        plugin(alpine, mergedOptions)
      } catch (e) {
        log.error('Failed to load plugin', e)
      }
    })

    this.initialized = true
    log.info('Uno initialized')
    return alpine
  }
}

export const Uno = new UnoInstance()
