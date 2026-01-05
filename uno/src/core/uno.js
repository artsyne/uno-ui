/**
 * Uno - Core Instance
 *
 * Manages plugins, initialization, and global configuration.
 *
 * @module uno.core.instance
 * @version 0.5.0
 */

import Alpine from 'alpinejs'
import { log } from '../utils/log'
import { once } from '../utils/dom'
import { mergeDeep } from '../utils/object'
import { Z_INDEX } from './layer'
import box from '../components/box'

/**
 * UnoInstance class.
 * Central hub for Uno UI, managing plugins, configuration, and HTMX integration.
 */
export class UnoInstance {
  /**
   * Create an UnoInstance.
   *
   * @param {Object} [options={}] - Global configuration options
   */
  constructor(options = {}) {
    this.version = '0.6.0'
    this.options = options
    this.pluginDefaults = {} // New: Plugin default configuration storage
    this.plugins = []
    this.initialized = false
    this._htmxRequestHandlers = []
    this._htmxActiveXhrs = typeof WeakSet === 'function' ? new WeakSet() : null
    this.htmx = {
      /**
       * Check if HTMX is available.
       * @returns {boolean} True if HTMX is loaded and available
       */
      available: () => typeof window !== 'undefined' && typeof window.htmx?.ajax === 'function',

      /**
       * Wrapper for htmx.ajax.
       *
       * @param {string} method - HTTP method (GET, POST, etc.)
       * @param {string} url - The URL to request
       * @param {Object} [ajaxOptions={}] - HTMX ajax options
       * @returns {Promise} Promise resolving to the result
       */
      ajax: (method, url, ajaxOptions = {}) => {
        if (typeof window === 'undefined' || typeof window.htmx?.ajax !== 'function') {
          return Promise.reject(new Error('htmx.ajax is not available'))
        }
        return window.htmx.ajax(method, url, ajaxOptions)
      },

      /**
       * Fetch an HTML fragment without inserting it into the DOM.
       * Useful for getting content to display in a modal or tooltip.
       *
       * @param {string} method - HTTP method
       * @param {string} url - The URL to request
       * @param {Object} [ajaxOptions={}] - HTMX ajax options
       * @returns {Promise<string>} The HTML fragment string
       */
      fragment: async (method, url, ajaxOptions = {}) => {
        if (typeof document === 'undefined') throw new Error('document is not available')
        const target = document.createElement('div')
        // Detached element to prevent immediate script execution
        // We intercept beforeSwap to get the content without inserting it
        let responseText = ''
        const onBeforeSwap = (evt) => {
          evt.preventDefault() // Prevent HTMX from inserting content
          responseText = evt.detail.serverResponse
          evt.detail.shouldSwap = false
        }
        target.addEventListener('htmx:beforeSwap', onBeforeSwap)

        try {
          // No need to append to document.body for htmx.ajax to work with explicit target
          const swap = ajaxOptions.swap || 'innerHTML'
          await this.htmx.ajax(method, url, { ...ajaxOptions, target, swap })
          return responseText
        } finally {
          target.removeEventListener('htmx:beforeSwap', onBeforeSwap)
          target.remove()
        }
      },

      /**
       * Register handlers for HTMX request lifecycle events.
       *
       * @param {Object} handlers - The handlers object
       * @param {Function} [handlers.start] - Called when a request starts
       * @param {Function} [handlers.end] - Called when a request ends (success, error, or timeout)
       */
      onRequest: (handlers = {}) => {
        if (!handlers || typeof handlers !== 'object') return
        this._htmxRequestHandlers.push(handlers)

        if (typeof window === 'undefined' || typeof document === 'undefined') return
        if (!window.htmx) return

        if (!once(document, 'uno:htmx:request-hooks')) return

        const getXhr = (event) => event?.detail?.xhr
        const callEach = (name, event) => {
          for (const h of this._htmxRequestHandlers) {
            const fn = h?.[name]
            if (typeof fn === 'function') {
              try {
                fn(event)
              } catch (e) {
                log.error('Uno.htmx.onRequest handler failed', e)
              }
            }
          }
        }

        const begin = (event) => {
          const xhr = getXhr(event)
          if (xhr && this._htmxActiveXhrs) {
            if (this._htmxActiveXhrs.has(xhr)) return
            this._htmxActiveXhrs.add(xhr)
          }
          callEach('start', event)
        }

        const end = (event) => {
          const xhr = getXhr(event)
          if (xhr && this._htmxActiveXhrs) {
            if (!this._htmxActiveXhrs.has(xhr)) return
            this._htmxActiveXhrs.delete(xhr)
          }
          callEach('end', event)
        }

        document.addEventListener('htmx:beforeRequest', begin)
        document.addEventListener('htmx:afterRequest', end)
        document.addEventListener('htmx:responseError', end)
        document.addEventListener('htmx:sendError', end)
        document.addEventListener('htmx:timeout', end)

        document.addEventListener('htmx:xhr:progress', (event) => callEach('progress', event))
      },
      // Global HTMX event binding for notify functionality - progress bar only
      bindNotifyToHtmx: () => {
        if (typeof window === 'undefined' || typeof document === 'undefined') return
        if (!window.htmx || !window.Alpine) return

        if (!once(document, 'uno:htmx:bind-notify')) return

        document.addEventListener('htmx:beforeRequest', () => {
          window.Alpine.store('uno').notify.startProgress()
        })

        document.addEventListener('htmx:afterRequest', () => {
          window.Alpine.store('uno').notify.stopProgress()
        })

        document.addEventListener('htmx:responseError', () => {
          window.Alpine.store('uno').notify.stopProgress()
          window.Alpine.store('uno').notify.error('Request failed')
        })

        document.addEventListener('htmx:sendError', () => {
          window.Alpine.store('uno').notify.stopProgress()
          window.Alpine.store('uno').notify.error('Request failed')
        })

        document.addEventListener('htmx:timeout', () => {
          window.Alpine.store('uno').notify.stopProgress()
          window.Alpine.store('uno').notify.warning('Request timed out')
        })
      }
    }
    const base = 'rounded-lg shadow-md border border-gray-200 dark:border-gray-700'

    const defaultOptions = {
      debug: false,
      csp: {
        nonce: null // For Content Security Policy
      },
      onError: null, // Custom error handler
      plugins: {},
      ui: {
        base,
        layers: { box: Z_INDEX.BOX, tooltip: Z_INDEX.TOOLTIP, notify: Z_INDEX.NOTIFY },
        box: {
          overlay: 'fixed inset-0 bg-gray-500/75 transition-opacity',
          dialogWrapper: 'fixed inset-0 z-10 w-screen overflow-y-auto',
          itemWrapper: 'flex min-h-full p-4 text-center sm:p-0 pointer-events-none',
          container:
            'relative transform overflow-hidden {base} bg-white dark:bg-gray-900 text-left shadow-xl transition-all sm:my-8 sm:w-full sm:max-w-lg flex flex-col max-h-[90vh] pointer-events-auto outline-none',
          cancelBtn:
            'mt-3 inline-flex w-full justify-center rounded-md bg-white px-3 py-2 text-sm font-semibold text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 hover:bg-gray-50 sm:mt-0 sm:w-auto',
          okBtn:
            'inline-flex w-full justify-center rounded-md bg-blue-600 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-blue-500 sm:ml-3 sm:w-auto',
          title: 'text-base font-semibold leading-6 text-gray-900 dark:text-slate-200',
          closeBtnFloating:
            'absolute top-2 right-2 z-10 text-slate-400 hover:text-slate-500 focus:outline-none',
          closeIconSmall: 'h-5 w-5',
          content:
            'px-4 pb-4 pt-5 sm:p-6 sm:pb-4 overflow-auto flex-1 text-slate-600 dark:text-slate-300',
          loadingWrap: 'flex justify-center items-center h-full',
          spinner: 'animate-spin h-8 w-8 text-blue-500',
          footer:
            'bg-gray-50 px-4 py-3 sm:flex sm:flex-row-reverse sm:px-6 border-t border-slate-200 dark:border-slate-700 dark:bg-slate-800',
          promptWrapper: 'space-y-3',
          promptText: 'text-slate-700 dark:text-slate-200',
          promptInput:
            'w-full px-3 py-2 border rounded focus:ring-2 focus:ring-blue-500 outline-none'
        },
        notify: {
          container: 'fixed flex gap-2 p-4 pointer-events-none',
          item: 'pointer-events-auto max-w-[calc(100vw-2rem)] overflow-hidden {base}',
          content: 'p-4',
          row: 'flex items-start',
          iconWrap: 'flex-shrink-0',
          textWrap: 'ml-3 min-w-0 flex-1 pt-0.5',
          text: 'text-sm font-medium text-gray-900 dark:text-slate-100 break-words',
          closeWrap: 'ml-4 flex flex-shrink-0',
          closeBtn:
            'inline-flex rounded-md bg-white dark:bg-slate-900 text-gray-400 hover:text-gray-500 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2',
          closeIcon: 'h-5 w-5',
          progressWrap:
            'fixed top-0 left-0 right-0 h-1 bg-gray-200 dark:bg-gray-700 pointer-events-none',
          progressBar:
            'h-full bg-gradient-to-r from-blue-500 to-indigo-600 shadow-sm shadow-blue-500/50',
          types: {
            success: 'bg-green-100 dark:bg-green-950',
            error: 'bg-red-100 dark:bg-red-950',
            warning: 'bg-yellow-100 dark:bg-yellow-950',
            info: 'bg-blue-100 dark:bg-blue-950'
          },
          icons: {
            success: 'text-green-500 dark:text-green-300',
            error: 'text-red-500 dark:text-red-300',
            warning: 'text-yellow-500 dark:text-yellow-300',
            info: 'text-blue-500 dark:text-blue-300'
          }
        },
        tooltip: {
          content:
            'fixed px-3 py-2 text-sm font-medium {base} whitespace-nowrap bg-white text-gray-900 dark:bg-slate-800 dark:text-white',
          body: 'pointer-events-auto'
        },
        tab: {
          active: 'border-b-2 border-blue-500',
          inactive: 'text-slate-500 hover:text-slate-700',
          loading: 'animate-pulse h-32 w-full rounded-lg bg-gray-200 dark:bg-slate-700',
          error: 'p-4 text-red-600 dark:text-red-400'
        },
        valid: {
          wrapper: 'relative w-full',
          iconBase:
            'absolute right-0 top-0 bottom-0 flex items-center justify-center w-10 pointer-events-none transition-opacity duration-200',
          inputError: '!border-red-500',
          inputValid: '!border-green-500',
          iconError: 'text-red-500 dark:text-red-400',
          iconValid: 'text-green-500 dark:text-green-400',
          stateInvalid: 'is-invalid',
          stateValid: 'is-valid',
          stateValidating: 'validating',
          inputPaddingIcon: 'pr-10',
          helperText: 'text-red-600 dark:text-red-400 text-xs'
        }
      }
    }

    this.options = mergeDeep(defaultOptions, options)
  }

  /**
   * ✅ Unified configuration retrieval method
   * Supports 4-level merge: Plugin Defaults -> Global User Config -> Global Plugin Config -> Direct Global Config -> Instance Config
   * Note: We use mergeDeep to ensure deep merging of options.
   *
   * Level 1: Plugin Defaults (this.pluginDefaults[componentName])
   * Level 2: Global User Config (this.options.plugins[componentName])
   * Level 3: Direct Global Config (this.options[componentName]) - Legacy/Shorthand
   * Level 4: Instance Config (instanceOpts)
   *
   * @param {string} componentName - Name of the component
   * @param {Object} [instanceOpts={}] - Instance-specific options
   * @returns {Object} The merged configuration object
   */
  getConfig(componentName, instanceOpts = {}) {
    const globalOpts = this.options || {}

    // Level 1: Plugin Defaults
    let config = this.pluginDefaults[componentName] ? { ...this.pluginDefaults[componentName] } : {}

    // Level 2: Global User Config (plugins)
    if (globalOpts.plugins && globalOpts.plugins[componentName]) {
      config = mergeDeep(config, globalOpts.plugins[componentName])
    }

    // Level 3: Direct Global Config
    if (globalOpts[componentName]) {
      config = mergeDeep(config, globalOpts[componentName])
    }

    // Level 4: Instance Config
    if (instanceOpts) {
      config = mergeDeep(config, instanceOpts)
    }

    return config
  }

  /**
   * ✅ Register plugin defaults
   *
   * @param {string} name - Component name
   * @param {Object} defaults - Default configuration object
   */
  registerPluginDefaults(name, defaults) {
    this.pluginDefaults[name] = defaults
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

  _resolveUiTokens(value) {
    if (typeof value !== 'string') return value
    const base = this.options.ui?.base
    if (typeof base === 'string' && value.indexOf('{base}') !== -1) {
      return value.replaceAll('{base}', base)
    }
    return value
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

    if (config === undefined) return this._resolveUiTokens(defaultCls)

    // Check for append mode ("+class")
    if (typeof config === 'string' && config.trim().indexOf('+') === 0) {
      return this._resolveUiTokens(`${defaultCls} ${config.trim().slice(1)}`.trim())
    }

    return this._resolveUiTokens(config)
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

    if (typeof window !== 'undefined' && typeof document !== 'undefined') {
      if (window.htmx && typeof alpine?.initTree === 'function') {
        if (once(document, 'uno:htmx-afterSwap:initTree')) {
          document.addEventListener('htmx:afterSwap', (event) => {
            const elt = event?.detail?.elt
            if (elt) alpine.initTree(elt)
          })
        }
      }
    }

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
