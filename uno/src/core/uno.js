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

class UnoInstance {
  constructor() {
    this.version = '0.3.0'
    this.plugins = []
    this.initialized = false
    this.options = {
      debug: true,
      csp: {
        nonce: null // For Content Security Policy
      },
      onError: null // Custom error handler
    }
  }

  /**
   * Set global configuration (before or after init)
   * @param {Object} config - Global configuration object
   */
  setConfig(config) {
    const newOptions = { ...this.options }
    
    // Deep merge (1-level) for component configs to preserve other properties (e.g. duration)
    Object.keys(config).forEach(key => {
        if (typeof config[key] === 'object' && config[key] !== null && !Array.isArray(config[key]) && newOptions[key]) {
            newOptions[key] = { ...newOptions[key], ...config[key] }
        } else {
            newOptions[key] = config[key]
        }
    })
    
    this.options = newOptions

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

    this.options = { ...this.options, ...opts }

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
        const mergedOptions = { ...this.options, ...pluginOptions }
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
