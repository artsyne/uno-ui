import { Uno } from '../../core/uno'
import { getFlexAlignmentClass } from '../../utils/position'
import { getZIndex } from '../../core/layer'
import { getBoxTemplate } from './box.template'

export const createBoxStore = (Alpine) => ({
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
      const wrapperClass = Uno.getStyle('box', 'promptWrapper')
      const textClass = Uno.getStyle('box', 'promptText')
      const inputClass = Uno.getStyle('box', 'promptInput')

      // Escape HTML in message and defaultValue for security
      const escapedMessage = Uno.dom.sanitize(message)
      const escapedDefaultValue = Uno.dom.sanitize(defaultValue)

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

    // Note: Alpine(v3) reactivity is based on Proxy.
    // We must access the "reactive proxy object" via this.items to modify fields,
    // otherwise assigning to the original item object won't trigger view updates
    // (causing box fail to open or animation not working).
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

    const zBase = Uno.options?.ui?.layers?.box || getZIndex('BOX')

    const attributes = {
      'x-data': '{}',
      class: 'relative',
      style: `z-index: ${zBase};`,
      'x-show': '$store.uno.box.items.length > 0',
      'x-cloak': '' // Hide until Alpine is ready
    }

    const template = getBoxTemplate()

    const container = Uno.dom.ensureContainer(id, template, attributes)
    if (container && Alpine && typeof Alpine.initTree === 'function') {
      Alpine.initTree(container)
    }
    return container
  }
})
