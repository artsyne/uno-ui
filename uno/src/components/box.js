/**
 * Uno - Box Component (Store-First, Directive-Driven Architecture)
 *
 * A lightweight, declarative, HTMX-native modal component.
 * Follows the "Store-First, Directive-Driven" pattern.
 *
 * Features:
 * - 9-point grid positioning system
 * - HTMX native content loading
 * - Automatic XSS sanitization (trustedHTML removed for security)
 * - Focus trap and keyboard navigation
 * - i18n support for button text
 * - Callback mechanism (onOk, onCancel, onBeforeClose, onClose, onContentLoad)
 *
 * @module uno.components.box
 * @version 0.6.0
 * @license MIT
 */

import { trapFocus } from '../utils/a11y'
import { Uno } from '../core/uno'
import { registerPlugin } from '../core/plugin-registry'
import { getFlexAlignmentClass } from '../utils/position'

/**
 * Box Plugin Factory
 * @param {object} Alpine - Alpine.js instance
 * @param {object} options - Plugin options
 */
export default (Alpine, options = {}) => {
  // Register plugin defaults
  const pluginDefaults = {
    defaultWidth: 400,
    defaultPosition: 'center',
    okText: 'OK',
    cancelText: 'Cancel',
    closeOnEsc: false,
    closeOnOverlayClick: false
  }
  Uno.registerPluginDefaults('box', { ...pluginDefaults, ...options.plugins?.box, ...options.box })

  /**
   * Helper: Get UI style class
   * @param {string} key - Style key
   * @returns {string} Style class string
   */
  const getUi = (key) => Uno.getStyle('box', key)

  // ============================================================
  // Focus Trap Directive
  // ============================================================
  Alpine.directive('uno-trap', (el, { expression: _expression }, { cleanup }) => {
    void _expression
    const cleanupTrap = trapFocus(el)
    cleanup(cleanupTrap)
  })

  // ============================================================
  // Box Directive (x-uno-box)
  // Declarative interface for triggering box dialogs
  // ============================================================
  Alpine.directive('uno-box', (el, { expression }, { evaluateLater, cleanup }) => {
    const hasExpression = expression && expression.trim().length > 0
    const getExpression = hasExpression ? evaluateLater(expression) : null

    /**
     * Parse boolean from attribute value
     * @param {string|null} value - Attribute value
     * @returns {boolean|undefined} Parsed boolean
     */
    const parseBool = (value) => {
      if (value === null || value === undefined) return undefined
      const v = String(value).trim().toLowerCase()
      if (v === 'true' || v === '1') return true
      if (v === 'false' || v === '0') return false
      return undefined
    }

    /**
     * Parse number from attribute value
     * @param {string|null} value - Attribute value
     * @returns {number|undefined} Parsed number
     */
    const parseNumber = (value) => {
      if (value === null || value === undefined || value === '') return undefined
      const n = Number(value)
      return Number.isFinite(n) ? n : undefined
    }

    /**
     * Build configuration from data attributes
     * @returns {object} Configuration object
     */
    const buildFromAttrs = () => {
      const type =
        el.getAttribute('data-uno-box') || el.getAttribute('data-uno-box-type') || 'dialog'
      const title = el.getAttribute('data-uno-title') || ''
      const content = el.getAttribute('data-uno-content') || ''
      const hxUrl = el.getAttribute('data-uno-hx-url') || undefined
      const position = el.getAttribute('data-uno-position') || undefined
      const modal = parseBool(el.getAttribute('data-uno-modal'))
      const closeOnEsc = parseBool(el.getAttribute('data-uno-close-on-esc'))
      const closeOnOverlayClick = parseBool(el.getAttribute('data-uno-close-on-overlay-click'))
      const width = parseNumber(el.getAttribute('data-uno-width'))
      const height = parseNumber(el.getAttribute('data-uno-height'))
      const widthClass = el.getAttribute('data-uno-width-class') || undefined
      const heightClass = el.getAttribute('data-uno-height-class') || undefined

      return {
        type,
        title,
        content,
        hxUrl,
        position,
        modal,
        closeOnEsc,
        closeOnOverlayClick,
        width,
        height,
        widthClass,
        heightClass
      }
    }

    /**
     * Click handler for directive
     * @param {Event} e - Click event
     */
    const onClick = (e) => {
      e.preventDefault()
      // Access store via uno namespace
      const store = Alpine.store('uno')?.box
      if (!store) {
        console.warn('[Uno] Box store not found. Is Uno initialized?')
        return
      }

      if (getExpression) {
        getExpression((evaluated) => {
          if (evaluated && typeof evaluated === 'object') {
            store.open(evaluated)
          } else {
            // Fallback: If expression evaluates to falsy (or we want to support conditional opening),
            // we might want to do nothing.
            // But if it's just empty/undefined, we check attributes.
            // Since we checked hasExpression before creating getExpression,
            // we are here because expression string was not empty.
            // If the user passed a variable that is null, we probably shouldn't open.
            // BUT, strictly following the logic "attribute fallback":
            if (!evaluated) {
              const attrs = buildFromAttrs()
              // Only open if we have at least a type or content from attributes to justify it
              if (attrs.title || attrs.content || attrs.hxUrl) {
                store.open(attrs)
              }
            }
          }
        })
        return
      }

      store.open(buildFromAttrs())
    }

    el.addEventListener('click', onClick)
    cleanup(() => el.removeEventListener('click', onClick))
  })

  // ============================================================
  // Box Store (Store-First Architecture)
  // Registered under uno namespace: Alpine.store('uno').box
  // ============================================================
  const boxStore = {
    items: [],

    /**
     * Initialize the box component
     * Ensures container exists when store is initialized
     */
    init() {
      this.ensureContainer()
    },

    // ============================================================
    // Public API - Convenience Methods
    // ============================================================

    /**
     * Show an alert dialog
     * @param {string} message - Alert message
     * @param {object} options - Additional options
     * @returns {object} Box item instance
     */
    alert(message, options = {}) {
      const opts = typeof options === 'string' ? { title: options } : options
      return this.open('dialog', {
        content: message,
        showOkButton: true,
        ...opts
      })
    },

    /**
     * Show a confirm dialog
     * @param {string} message - Confirmation message
     * @param {object} options - Additional options
     * @returns {Promise<boolean>} Resolves true if OK, false otherwise
     */
    confirm(message, options = {}) {
      return new Promise((resolve) => {
        this.open('confirm', {
          content: message,
          onOk: () => resolve(true),
          onCancel: () => resolve(false),
          onClose: () => resolve(false), // Closing without choice = cancel
          ...options
        })
      })
    },

    /**
     * Show a prompt dialog
     * @param {string} message - Prompt message
     * @param {string} defaultValue - Default input value
     * @param {object} options - Additional options
     * @returns {Promise<string|null>} Resolves with input value or null
     */
    prompt(message, defaultValue = '', options = {}) {
      return new Promise((resolve) => {
        const id = 'prompt-' + Math.random().toString(36).slice(2, 11)
        const wrapperClass = getUi('promptWrapper')
        const textClass = getUi('promptText')
        const inputClass = getUi('promptInput')

        // Escape HTML in message and defaultValue for security
        const escapedMessage = Uno.dom.escape(message)
        const escapedDefaultValue = Uno.dom.escape(defaultValue)

        const content = `
          <div class="${wrapperClass}">
            <p class="${textClass}">${escapedMessage}</p>
            <input id="${id}" type="text" class="${inputClass}" value="${escapedDefaultValue}">
          </div>
        `
        this.open('confirm', {
          title: options.title || 'Prompt',
          content,
          _isTrustedPrompt: true, // Internal flag for prompt content
          onOk: () => {
            const input = document.getElementById(id)
            resolve(input ? input.value : null)
          },
          onCancel: () => resolve(null),
          ...options
        })
      })
    },

    /**
     * Show a dialog (alias for open)
     * @param {object} options - Dialog options
     * @returns {object} Box item instance
     */
    dialog(options = {}) {
      return this.open('dialog', options)
    },

    // ============================================================
    // Core Methods
    // ============================================================

    /**
     * Open a box dialog
     * Uses 4-level configuration merge: Plugin Defaults -> Global -> Plugin Config -> Instance
     * @param {string|object} typeOrOpts - Type string or options object
     * @param {object} opts - Options if first arg is type
     * @returns {object} Box item instance
     */
    open(typeOrOpts, opts = {}) {
      let type = 'dialog'
      let instanceOpts = opts

      if (typeof typeOrOpts === 'object') {
        instanceOpts = typeOrOpts
        type = instanceOpts.type || 'dialog'
      } else {
        type = typeOrOpts
      }

      // 4-level configuration merge via Uno.getConfig
      const config = Uno.getConfig('box', instanceOpts)

      const id = config.boxId || Math.random().toString(36).slice(2, 10)

      // Determine button visibility based on type
      let showOkButton = config.showOkButton
      let showCancelButton = config.showCancelButton

      if (type === 'confirm') {
        showOkButton = showOkButton !== false
        showCancelButton = showCancelButton !== false
      } else if (type === 'dialog') {
        showOkButton = showOkButton || false
        showCancelButton = showCancelButton || false
      }

      const item = {
        id,
        type,
        title: config.title || '',
        content: config._isTrustedPrompt
          ? config.content || ''
          : Uno.dom.sanitize(config.content || ''),
        // Security: Always sanitize content (trustedHTML removed)
        // Exception: Internal prompt content is pre-escaped
        _isTrustedPrompt: config._isTrustedPrompt || false,
        width:
          config.width !== undefined
            ? config.width
            : config.widthClass
              ? undefined
              : config.defaultWidth || 400,
        height: config.height,
        widthClass: config.widthClass,
        heightClass: config.heightClass,
        containerClass: config.containerClass,
        position: config.position || config.defaultPosition || 'center',
        modal: config.modal !== undefined ? config.modal : true,
        closeOnEsc: config.closeOnEsc !== undefined ? config.closeOnEsc : false,
        closeOnOverlayClick:
          config.closeOnOverlayClick !== undefined ? config.closeOnOverlayClick : false,
        showCloseButton: config.showCloseButton !== undefined ? config.showCloseButton : true,
        showOkButton,
        showCancelButton,
        okText: config.okText || 'OK',
        cancelText: config.cancelText || 'Cancel',
        onOk: config.onOk,
        onCancel: config.onCancel,
        onBeforeClose: config.onBeforeClose,
        onClose: config.onClose,
        onContentLoad: config.onContentLoad,
        hxUrl: config.hxUrl,
        hxMethod: config.hxMethod || 'GET',
        hxHeaders: config.hxHeaders,
        hxTrigger: config.hxTrigger || 'load',

        // Internal state - start closed to trigger animation
        isOpen: false,

        /**
         * Close this box dialog
         */
        close() {
          // Execute onBeforeClose callback
          if (this.onBeforeClose) {
            this.onBeforeClose()
          }

          this.isOpen = false

          // Use transitionend for precise timing instead of hardcoded timeout
          const el =
            typeof document !== 'undefined' ? document.getElementById('uno-box-' + this.id) : null

          const cleanup = () => {
            const store = Alpine.store('uno')?.box
            if (store) {
              store.remove(this.id)
            }
            if (this.onClose) this.onClose()
          }

          if (!el) {
            cleanup()
            return
          }

          let resolved = false
          const onEnd = (e) => {
            // Ignore bubbling events from children
            if (e && e.target !== el) return
            if (resolved) return
            resolved = true
            el.removeEventListener('transitionend', onEnd)
            cleanup()
          }

          el.addEventListener('transitionend', onEnd)
          // Safety fallback in case transition doesn't fire (e.g. hidden/detached)
          setTimeout(onEnd, 350)
        }
      }

      this.items.push(item)
      this.ensureContainer()

      // 注意：Alpine(v3) 的响应式基于 Proxy。
      // 这里必须通过 this.items 取到“响应式代理对象”再改字段，
      // 否则对原始对象 item 的赋值不会触发视图更新（会导致 box 无法打开/动画不生效）。
      const reactiveItem = this.items[this.items.length - 1]

      const openNow = () => {
        reactiveItem.isOpen = true
      }

      if (Alpine && typeof Alpine.nextTick === 'function') {
        Alpine.nextTick(openNow)
      } else {
        setTimeout(openNow, 0)
      }

      return reactiveItem
    },

    /**
     * Close a specific box by ID
     * @param {string} id - Box ID
     */
    close(id) {
      const item = this.items.find((i) => i.id === id)
      if (item) {
        item.close()
      }
    },

    /**
     * Remove a box from the items array
     * @param {string} id - Box ID
     */
    remove(id) {
      const index = this.items.findIndex((i) => i.id === id)
      if (index !== -1) {
        this.items.splice(index, 1)
      }
    },

    /**
     * Close all open boxes
     */
    closeAll() {
      this.items.forEach((item) => {
        if (item.onBeforeClose) item.onBeforeClose()
        if (item.onClose) item.onClose()
      })
      this.items = []
    },

    /**
     * Get position CSS classes for 9-point grid
     * @param {string} position - Position name
     * @returns {string} CSS classes
     */
    getPositionClasses(position) {
      return getFlexAlignmentClass(position)
    },

    /**
     * Handle HTMX content load completion
     * @param {HTMLElement} el - Content element
     * @param {object} item - Box item
     */
    onHtmxContentLoaded(el, item) {
      if (item.onContentLoad && typeof item.onContentLoad === 'function') {
        item.onContentLoad(el)
      }
      // Re-initialize Alpine components in loaded content
      if (window.Alpine && window.Alpine.initTree) {
        window.Alpine.initTree(el)
      }
    },

    /**
     * Ensure the box container exists in the DOM
     */
    ensureContainer() {
      const id = 'uno-box-manager'
      const existing = typeof document !== 'undefined' && document.getElementById(id)

      if (existing) {
        if (Alpine && typeof Alpine.initTree === 'function' && !existing._x_dataStack) {
          Alpine.initTree(existing)
        }
        return existing
      }

      const zBase = Uno.options?.ui?.layers?.boxBase || 10000

      const attributes = {
        'x-data': '{}',
        class: 'relative',
        style: `z-index: ${zBase};`,
        'x-show': '$store.uno.box.items.length > 0',
        'x-cloak': '' // Hide until Alpine is ready
      }

      // Template using $store.uno.box namespace
      const template = `
        <!-- Backdrop / Overlay -->
        <div id="uno-box-overlay"
             :class="Uno.getStyle('box', 'overlay')"
             x-show="$store.uno.box.items.length > 0"
             x-transition:enter="ease-out duration-300"
             x-transition:enter-start="opacity-0"
             x-transition:enter-end="opacity-100"
             x-transition:leave="ease-in duration-200"
             x-transition:leave-start="opacity-100"
             x-transition:leave-end="opacity-0"
             @click="$store.uno.box.items.length && $store.uno.box.items[$store.uno.box.items.length - 1].closeOnOverlayClick && $store.uno.box.items[$store.uno.box.items.length - 1].close()"></div>

        <div :class="Uno.getStyle('box', 'dialogWrapper')"
             @click.self="$store.uno.box.items.length && $store.uno.box.items[$store.uno.box.items.length - 1].closeOnOverlayClick && $store.uno.box.items[$store.uno.box.items.length - 1].close()">
          <template x-for="item in $store.uno.box.items" :key="item.id">
            <div :id="'uno-box-outer-container' + item.id"
                 :class="[Uno.getStyle('box', 'itemWrapper'), $store.uno.box.getPositionClasses(item.position)]"
                 :style="'z-index:' + (Uno.options?.ui?.layers?.boxBase || 10000)"
                 x-show="item.isOpen"
                 x-transition:enter="ease-out duration-300"
                 x-transition:enter-start="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95"
                 x-transition:enter-end="opacity-100 translate-y-0 sm:scale-100"
                 x-transition:leave="ease-in duration-200"
                 x-transition:leave-start="opacity-100 translate-y-0 sm:scale-100"
                 x-transition:leave-end="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95"
                 role="dialog"
                 :aria-modal="item.modal !== false ? 'true' : 'false'"
                 :aria-labelledby="item.title ? 'uno-box-title-' + item.id : undefined"
                 :aria-describedby="'uno-box-content-' + item.id"
                 @keydown.escape.window="item.isOpen && item.closeOnEsc && item.close()">

              <div :id="'uno-box-container' + item.id"
                   :class="[Uno.getStyle('box', 'container'), item.widthClass, item.heightClass, item.containerClass]"
                   :style="(item.width ? 'width:' + item.width + 'px; max-width:' + item.width + 'px;' : '') + (item.height ? ' height:' + item.height + 'px;' : '')"
                   x-uno-trap>

                <!-- Close Button -->
                <button type="button"
                        x-show="item.showCloseButton"
                        :class="Uno.getStyle('box', 'closeBtnFloating')"
                        @click="item.close()"
                        aria-label="Close">
                  <svg :class="Uno.getStyle('box', 'closeIconSmall')" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                  </svg>
                </button>

                <!-- Content Area -->
                <div :class="Uno.getStyle('box', 'content')">
                  <h3 :id="'uno-box-title-' + item.id"
                      :class="Uno.getStyle('box', 'title')"
                      x-text="item.title"
                      x-show="item.title"></h3>
                  <div class="mt-2" :id="'uno-box-content-' + item.id">
                    <template x-if="!item.hxUrl">
                      <div x-html="item.content"></div>
                    </template>
                    <template x-if="item.hxUrl">
                      <div x-init="htmx.process($el); $el.addEventListener('htmx:afterSwap', () => $store.uno.box.onHtmxContentLoaded($el, item))"
                           :hx-get="item.hxMethod === 'GET' ? item.hxUrl : undefined"
                           :hx-post="item.hxMethod === 'POST' ? item.hxUrl : undefined"
                           :hx-trigger="item.hxTrigger"
                           hx-swap="innerHTML"
                           :hx-headers="item.hxHeaders ? JSON.stringify(item.hxHeaders) : null">
                        <div :class="Uno.getStyle('box', 'loadingWrap')">
                          <svg :class="Uno.getStyle('box', 'spinner')" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                          </svg>
                        </div>
                      </div>
                    </template>
                  </div>
                </div>

                <!-- Footer with Buttons -->
                <div :class="Uno.getStyle('box', 'footer')"
                     x-show="item.showOkButton || item.showCancelButton">
                  <button type="button"
                          :id="'uno-box-ok' + item.id"
                          :class="Uno.getStyle('box', 'okBtn')"
                          x-show="item.showOkButton"
                          x-text="item.okText"
                          @click="item.onOk ? (item.onOk() !== false && item.close()) : item.close()"></button>
                  <button type="button"
                          :id="'uno-box-cancel' + item.id"
                          :class="Uno.getStyle('box', 'cancelBtn')"
                          x-show="item.showCancelButton"
                          x-text="item.cancelText"
                          @click="item.onCancel ? (item.onCancel() !== false && item.close()) : item.close()"></button>
                </div>
              </div>
            </div>
          </template>
        </div>
      `

      const container = Uno.dom.ensureContainer(id, template, attributes)
      if (container && Alpine && typeof Alpine.initTree === 'function') {
        Alpine.initTree(container)
      }
      return container
    }
  }

  // ============================================================
  // Register Store under uno namespace using unified registry
  // ============================================================
  registerPlugin(Alpine, 'box', boxStore)
}
