import { Uno } from '../../core/uno'
import { getFlexAlignmentClass } from '../../utils/position'
import { getZIndex } from '../../core/layer'
import { getBoxTemplate } from './box.template'

export const createBoxStore = (Alpine) => ({
  items: [],

  /**
   * Dialog stack management
   * @type {Array}
   */
  stack: [],

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
   * @param {object} options - Additional options. User-supplied onOk/onCancel/onClose
   *   callbacks are composed with (not replacing) the Promise resolver: the promise
   *   resolves true on OK / false otherwise. Like handleOk, an onOk returning false
   *   (or resolving to false) keeps the dialog open and the promise pending.
   * @returns {Promise<boolean>} Resolves true if OK, false otherwise
   */
  confirm(message, options = {}) {
    return new Promise((resolve) => {
      const { onOk, onCancel, onClose, ...rest } = options
      this.open('confirm', {
        ...rest,
        content: message,
        onOk: () => {
          const result = typeof onOk === 'function' ? onOk() : undefined
          if (result === false) return false // Keep open, promise stays pending
          if (result instanceof Promise) {
            return result.then((value) => {
              if (value === false) return false
              resolve(true)
              return undefined
            })
          }
          resolve(true)
          return undefined
        },
        onCancel: () => {
          if (typeof onCancel === 'function') onCancel()
          resolve(false)
        },
        onClose: () => {
          if (typeof onClose === 'function') onClose()
          resolve(false) // Closing without choice = cancel (no-op if already resolved)
        }
      })
    })
  },

  /**
   * Show a prompt dialog
   * @param {string} message - Prompt message
   * @param {string} defaultValue - Default input value
   * @param {object} options - Additional options. User-supplied onOk/onCancel are
   *   composed with the Promise resolver (see confirm()); onOk returning false
   *   keeps the dialog open and the promise pending.
   * @returns {Promise<string|null>} Resolves with input value or null
   */
  prompt(message, defaultValue = '', options = {}) {
    return new Promise((resolve) => {
      const { onOk, onCancel, ...rest } = options
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
        title: rest.title || 'Prompt',
        ...rest,
        content,
        _isTrustedPrompt: true, // Internal flag for prompt content
        onOk: () => {
          const result = typeof onOk === 'function' ? onOk() : undefined
          if (result === false) return false // Keep open, promise stays pending
          if (result instanceof Promise) {
            return result.then((value) => {
              if (value === false) return false
              const input = document.getElementById(id)
              resolve(input ? input.value : null)
              return undefined
            })
          }
          const input = document.getElementById(id)
          resolve(input ? input.value : null)
          return undefined
        },
        onCancel: () => {
          if (typeof onCancel === 'function') onCancel()
          resolve(null)
        }
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
      height:
        config.height !== undefined ? config.height : config.heightClass ? undefined : undefined,
      widthClass: config.widthClass,
      heightClass: config.heightClass,
      containerClass: config.containerClass,
      autoFocus: config.autoFocus !== undefined ? config.autoFocus : true,
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
      okButtonClass: config.okButtonClass,
      cancelButtonClass: config.cancelButtonClass,
      extraButtons: Array.isArray(config.extraButtons) ? config.extraButtons : [],
      onOk: config.onOk,
      onCancel: config.onCancel,
      onBeforeClose: config.onBeforeClose,
      onClose: config.onClose,
      onContentLoad: config.onContentLoad,
      hxUrl: config.hxUrl,
      hxMethod: config.hxMethod || 'GET',
      hxHeaders: config.hxHeaders,
      hxTrigger: config.hxTrigger || 'load',

      // B5：记录 opener，关闭时焦点还原（2.4.3）
      opener: typeof document !== 'undefined' ? document.activeElement : null,

      // Internal state - start closed to trigger animation
      isOpen: false,

      /**
       * B5 焦点还原：焦点仍在本 box 内（OK/Cancel/关闭按钮）或无处停留（body）时
       * 还原至 opener；焦点已被外部接管（如新弹窗/页面导航）时不抢占。
       * 容器锚点必须与 box.template.js 的 :id 绑定一致（uno-box-outer-container{id}，
       * 即 role="dialog" 元素）——与模板漂移会导致 contains 判定恒为 false。
       */
      restoreFocus() {
        const opener = this.opener
        if (!opener || opener === document.body || !opener.isConnected) return
        const el = document.getElementById('uno-box-outer-container' + this.id)
        const active = document.activeElement
        const focusInBox = !!el && !!active && el.contains(active)
        if (focusInBox || active === document.body) opener.focus()
      },

      /**
       * Close this box dialog
       */
      close() {
        // Execute onBeforeClose callback
        if (this.onBeforeClose) {
          this.onBeforeClose()
        }

        this.isOpen = false
        // B5：先重算下层 inert（多层弹窗），再还原焦点——顺序不可颠倒：
        // 若 opener 位于下层弹窗内，下层仍为 inert 时 focus() 会静默失败
        const store = Alpine.store('uno')?.box
        if (store) store.syncUnderlayInert()
        this.restoreFocus()

        // Use transitionend for precise timing instead of hardcoded timeout
        const el =
          typeof document !== 'undefined'
            ? document.getElementById('uno-box-outer-container' + this.id)
            : null

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

      this.scheduleAutoFocus(reactiveItem, Alpine)
      // B5：新弹窗置顶后，下层弹窗 inert + aria-hidden
      this.syncUnderlayInert()
    }

    if (Alpine && typeof Alpine.nextTick === 'function') {
      Alpine.nextTick(openNow)
    }

    return reactiveItem
  },

  /**
   * Handle extra button click
   * Supports synchronous and asynchronous onClick callbacks, consistent with handleOk
   * @param {object} item - Box item
   * @param {object} btn - Extra button config
   */
  async handleExtraButton(item, btn) {
    if (!item || !btn?.onClick) return

    const result = btn.onClick(item)

    if (result instanceof Promise) {
      const value = await result
      if (value !== false) {
        item.close()
      }
    } else if (result !== false) {
      item.close()
    }
  },

  /**
   * Handle OK button click
   * Supports synchronous and asynchronous onOk callbacks
   * @param {object} item - Box item
   */
  async handleOk(item) {
    if (!item) return

    if (item.onOk) {
      // Handle both sync and async return values
      const result = item.onOk()

      // If promise, wait for it
      if (result instanceof Promise) {
        const value = await result
        if (value !== false) {
          item.close()
        }
      }
      // If sync value
      else if (result !== false) {
        item.close()
      }
    } else {
      // No onOk handler, just close
      item.close()
    }
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
      // B5：顶层移除后重算下层 inert 状态
      this.syncUnderlayInert()
    }
  },

  /**
   * B5 多层弹窗：仅顶层可交互/可感知，下层 inert + aria-hidden，
   * 屏幕阅读器与键盘不穿透底层（不改现有 scheduleAutoFocus 契约）。
   */
  syncUnderlayInert() {
    if (typeof document === 'undefined') return
    const openItems = this.items.filter((i) => i.isOpen)
    openItems.forEach((it, idx) => {
      const el = document.getElementById('uno-box-outer-container' + it.id)
      if (!el) return
      const underlay = idx < openItems.length - 1
      el.toggleAttribute('inert', underlay)
      if (underlay) el.setAttribute('aria-hidden', 'true')
      else el.removeAttribute('aria-hidden')
    })
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

    this.scheduleAutoFocus(item, window.Alpine)
  },

  scheduleAutoFocus(item, alpine) {
    const containerId = 'uno-box-container' + item.id
    const run = () => {
      const container = document.getElementById(containerId)
      if (!container) return
      // APG initial-focus policy: prefer the first inputable element when
      // autoFocus is on; otherwise fall through to the neutral target so
      // focus never lands on chrome controls (e.g. the floating close
      // button, which Safari rings on script focus).
      if (item.autoFocus && this.focusFirstInput(container)) return
      this.focusNeutralTarget(container, item)
    }
    const defer = (fn) => {
      if (typeof requestAnimationFrame === 'function') {
        requestAnimationFrame(() => requestAnimationFrame(fn))
      } else {
        setTimeout(fn, 0)
      }
    }
    if (alpine && typeof alpine.nextTick === 'function') {
      alpine.nextTick(() => defer(run))
    } else {
      defer(run)
    }
  },

  /**
   * Focus the first input element in the container
   * @param {HTMLElement} container - Container element
   * @returns {boolean} True if an element was focused
   */
  focusFirstInput(container) {
    if (!container) return false

    // Selectors for input fields (excluding buttons, hidden, disabled, readonly)
    const selector = [
      'input:not([type="hidden"]):not([type="button"]):not([type="submit"]):not([type="reset"]):not([type="image"]):not([disabled]):not([readonly])',
      'textarea:not([disabled]):not([readonly])',
      'select:not([disabled]):not([readonly])'
    ].join(',')

    const inputs = container.querySelectorAll(selector)

    // Visibility by ancestor walk: offsetParent is unreliable in WebKit
    // (null inside position:fixed dialog wrappers), so never consult it here.
    const isVisible = (node) => {
      for (let el = node; el && el.nodeType === 1; el = el.parentElement) {
        const style = window.getComputedStyle(el)
        if (style.display === 'none' || style.visibility === 'hidden') return false
      }
      return true
    }

    // Find first visible input
    for (const input of inputs) {
      if (isVisible(input)) {
        input.focus()
        return true
      }
    }
    return false
  },

  /**
   * Focus the dialog's neutral target (WAI-ARIA APG dialog pattern):
   * the title heading when present, otherwise the container itself.
   * Keeps initial focus off chrome controls such as the close button.
   * @param {HTMLElement} container - Container element
   * @param {object} item - Box item
   * @returns {boolean} True if a neutral target was focused
   */
  focusNeutralTarget(container, item) {
    if (!container) return false
    // Never steal focus the user has already placed inside the dialog
    if (container.contains(document.activeElement)) return false

    const title = item && item.title ? document.getElementById('uno-box-title-' + item.id) : null
    const target = title && container.contains(title) ? title : container
    // Programmatic focus target must be focusable (template sets tabindex)
    if (!target.hasAttribute('tabindex')) target.setAttribute('tabindex', '-1')
    target.focus({ preventScroll: true })
    return true
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
