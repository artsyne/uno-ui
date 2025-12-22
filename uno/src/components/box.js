/**
 * Uno - Box Component (Store-Driven Refactor)
 *
 * @module uno.components.box
 * @version 0.2.1
 * @license MIT
 */

import { escapeHTML, once, sanitizeHTML, DomManager as domManager } from '../utils/dom'
import { log } from '../utils/log'
import { getFlexAlignmentClass } from '../utils/position'
import { requestHTML, requestData } from '../utils/request'
import { trapFocus } from '../utils/a11y'
import { Uno } from '../core/uno'

export default (Alpine, options = {}) => {
  const initConfig = options.plugins?.box || options.box || {}
  const getConfig = () => {
    const opts = typeof window !== 'undefined' ? window.Uno?.options : Uno.options
    return opts?.plugins?.box || opts?.box || initConfig
  }
  const getUi = (key, defaultCls) => Uno.getStyle('box', key, defaultCls)
  const getErrorHTML = (message) => {
    const cls = getUi('error', 'p-4 text-red-500')
    return `<div class="${cls}">${escapeHTML(message)}</div>`
  }

  let rootEl = null
  const elementCache = new Map()
  const getCachedEl = (id, prefix) => {
    const key = `${prefix}:${id}`
    const cached = elementCache.get(key)
    if (cached && typeof document !== 'undefined' && document.contains(cached)) {
      return cached
    }
    const el =
      typeof document !== 'undefined'
        ? document.getElementById(prefix + id) || rootEl?.querySelector?.(`#${prefix + id}`)
        : null
    if (el) elementCache.set(key, el)
    return el
  }

  // Register Focus Trap Directive
  Alpine.directive('uno-trap', (el, { expression: _expression }, { cleanup }) => {
    void _expression
    const cleanupTrap = trapFocus(el)
    cleanup(cleanupTrap)
  })

  // 1. Store
  Alpine.store('box', {
    items: [],

    // Public API
    alert(message, options = {}) {
      const opts = typeof options === 'string' ? { title: options } : options
      return this.open('dialog', {
        content: message,
        showOkButton: true,
        ...opts
      })
    },

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

    prompt(message, defaultValue = '', options = {}) {
      return new Promise((resolve) => {
        const id = 'prompt-' + Math.random().toString(36).slice(2, 11)
        const wrapperClass = getUi('promptWrapper', 'space-y-3')
        const textClass = getUi('promptText', 'text-slate-700 dark:text-slate-200')
        const inputClass = getUi(
          'promptInput',
          'w-full px-3 py-2 border rounded focus:ring-2 focus:ring-blue-500 outline-none'
        )
        const content = `
          <div class="${wrapperClass}">
            <p class="${textClass}">${escapeHTML(message)}</p>
            <input id="${id}" type="text" class="${inputClass}" value="${escapeHTML(defaultValue)}">
          </div>
        `
        this.open('confirm', {
          title: options.title || 'Prompt',
          content,
          trustedHTML: true,
          onOk: () => {
            const input = document.getElementById(id)
            resolve(input ? input.value : null)
          },
          onCancel: () => resolve(null),
          ...options
        })
      })
    },

    async open(typeOrOpts, opts = {}) {
      let type = 'dialog'
      let options = opts
      const config = getConfig()

      if (typeof typeOrOpts === 'object') {
        options = typeOrOpts
        type = options.type || 'dialog'
      } else {
        type = typeOrOpts
      }

      const id = options.boxId || Math.random().toString(36).slice(2, 10)

      const item = {
        id,
        type,
        title: options.title || '',
        content: options.content || '',
        trustedHTML: options.trustedHTML === true,
        width:
          options.width !== undefined
            ? options.width
            : options.widthClass
              ? undefined
              : config.defaultWidth || 400,
        height: options.height,
        widthClass: options.widthClass,
        heightClass: options.heightClass,
        containerClass: options.containerClass,
        modal: options.modal !== false, // default true
        overlayOpacity: options.overlayOpacity ?? 0.5,
        position: options.position || 'center',
        showCloseButton: options.showCloseButton !== undefined ? options.showCloseButton : true,
        // Buttons
        okText: options.okButtonText,
        cancelText: options.cancelButtonText,
        onOk: options.onOk,
        onCancel: options.onCancel,
        onClose: options.onClose,
        onBeforeClose: options.onBeforeClose,
        onContentLoad: options.onContentLoad,
        showOk:
          type === 'confirm' || type === 'message' || (type === 'dialog' && options.showOkButton),
        showCancel: type === 'confirm' || (type === 'dialog' && options.showCancelButton),

        // Cache Option
        cache: options.cache === true,

        // State
        loading: false,
        visible: true,
        abortController: null
      }

      // Sanitize Initial Content if needed
      if (!item.trustedHTML) {
        item.content = sanitizeHTML(item.content)
        item.trustedHTML = true // Now it's safe
      }

      // I18n fallbacks
      const i18n = Alpine.store('i18n')
      if (!item.okText) item.okText = i18n ? i18n.t('box.ok') : 'OK'
      if (!item.cancelText) item.cancelText = i18n ? i18n.t('box.cancel') : 'Cancel'

      // Content Loading
      if (options.targetId) {
        const target = document.getElementById(options.targetId)
        if (target) {
          item.content = target.innerHTML
          // targetId content assumes trusted or needs sanitization?
          // If we took it from DOM, it might contain scripts if not careful.
          // But usually targetId implies developer control.
          // Let's check trustedHTML option again.
          if (!options.trustedHTML) {
            item.content = sanitizeHTML(item.content)
          }
          item.trustedHTML = true
        } else {
          log.warn('Box: Target element not found:', options.targetId)
          item.content = getErrorHTML('Error: Content target not found')
          item.trustedHTML = true
        }
      } else if (options.hxUrl) {
        item.loading = true
        this.add(item)

        try {
          const html = await requestHTML(options.hxUrl, { method: 'GET', cache: item.cache })
          const current = this.items.find((i) => i.id === id)
          if (current) {
            current.content = options.trustedHTML ? html : sanitizeHTML(html)
            current.trustedHTML = true
            current.loading = false
            if (current.onContentLoad) {
              Alpine.nextTick(() => {
                current.onContentLoad(getCachedEl(id, 'uno-box-content'))
              })
            }
          }
        } catch (err) {
          log.error('Box: HTMX Load Failed', err)
          const current = this.items.find((i) => i.id === id)
          if (current) {
            current.loading = false
            current.trustedHTML = true
            current.content = getErrorHTML('Error loading content via HTMX')
          }
        }
        return id
      } else if (options.ajaxUrl) {
        item.loading = true
        item.abortController = new AbortController()
        this.add(item)

        try {
          const data = await requestData(options.ajaxUrl, {
            cache: item.cache,
            signal: item.abortController.signal
          })
          const preClass = getUi(
            'ajaxPre',
            'bg-gray-100 dark:bg-slate-800 p-3 rounded overflow-x-auto whitespace-pre-wrap'
          )
          let content =
            typeof data === 'object'
              ? `<pre class="${preClass}">${escapeHTML(JSON.stringify(data, null, 2))}</pre>`
              : String(data)

          if (!options.trustedHTML) {
            content = sanitizeHTML(content)
          }

          const current = this.items.find((i) => i.id === id)
          if (current) {
            current.content = content
            current.trustedHTML = true
            current.loading = false

            Alpine.nextTick(() => {
              const el = getCachedEl(id, 'uno-box-content')
              if (el) {
                if (current.onContentLoad) {
                  current.onContentLoad(el)
                }
              }
            })
          }
        } catch (e) {
          if (e.name === 'AbortError') {
            log.info('Box: Request cancelled')
            return id
          }
          log.error('Box: Ajax failed', e)
          const current = this.items.find((i) => i.id === id)
          if (current) {
            current.content = getErrorHTML('Error loading content')
            current.trustedHTML = true
            current.loading = false
          }
        }
        return id
      } else if (options.iframeUrl) {
        const iframeClass = getUi('iframe', 'w-full h-full')
        item.content = `<iframe frameborder="0" class="${iframeClass}" src="${escapeHTML(options.iframeUrl)}"></iframe>`
        item.trustedHTML = true
      }

      this.add(item)
      return id
    },

    add(item) {
      this.items.push(item)
      this.ensureContainer()
    },

    close(id) {
      const item = this.items.find((i) => i.id === id)
      if (item) {
        // Cancel pending request
        if (item.abortController) {
          item.abortController.abort()
        }

        if (item.onBeforeClose) item.onBeforeClose()
        this.items = this.items.filter((i) => i.id !== id)
        elementCache.delete(`uno-box-content:${id}`)
        if (item.onClose) item.onClose()
      }
    },

    closeAll() {
      // Cancel all
      this.items.forEach((item) => {
        if (item.abortController) item.abortController.abort()
      })
      this.items = []
      elementCache.clear()
    },

    getPositionClass(position) {
      const base = getUi(
        'positionBase',
        'fixed inset-0 flex transition-opacity duration-300 p-4 pointer-events-none'
      )
      const alignment = getFlexAlignmentClass(position)
      return `${base} ${alignment}`
    },

    ensureContainer() {
      injectContainer()
    }
  })

  // 2. Container Injection
  const injectContainer = () => {
    if (typeof document !== 'undefined' && document.getElementById('uno-box-container')) return

    const rounded = Uno.getStyle('base', 'rounded', 'rounded-lg')
    const shadow = Uno.getStyle('base', 'shadow', 'shadow-md')
    const primary = Uno.getStyle('base', 'primary', 'bg-blue-600')

    let zBase = Number(Uno.options.ui?.layers?.boxBase ?? 10000)
    if (!Number.isFinite(zBase)) zBase = 10000
    let zStep = Number(Uno.options.ui?.layers?.boxStep ?? 1)
    if (!Number.isFinite(zStep) || zStep <= 0) zStep = 1

    // Overlay
    const overlayDefault =
      'absolute inset-0 bg-black transition-opacity duration-300 pointer-events-auto'
    const overlayClass = Uno.getStyle('box', 'overlay', overlayDefault)

    // Container
    const containerDefault = `relative bg-white dark:bg-gray-900 ${rounded} ${shadow} ring-1 ring-black/5 overflow-hidden flex flex-col max-h-[90vh] transition-all duration-300 transform pointer-events-auto outline-none`
    const containerClass = Uno.getStyle('box', 'container', containerDefault)

    // Buttons
    const cancelDefault = `px-4 py-2 bg-slate-100 dark:bg-slate-700 border border-slate-200 dark:border-slate-600 rounded-md text-sm font-medium text-slate-800 dark:text-slate-200 hover:bg-slate-200 dark:hover:bg-slate-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 transition-colors`
    const cancelClass = Uno.getStyle('box', 'cancelBtn', cancelDefault)

    const okDefault = `px-4 py-2 ${primary} border border-transparent rounded-md text-sm font-medium text-white hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500`
    const okClass = Uno.getStyle('box', 'okBtn', okDefault)

    const headerDefault =
      'flex justify-between items-center p-4 border-b border-slate-200 dark:border-slate-700'
    const headerClass = getUi('header', headerDefault)

    const titleDefault = 'text-lg font-semibold text-slate-800 dark:text-slate-200'
    const titleClass = getUi('title', titleDefault)

    const closeBtnDefault = 'text-slate-400 hover:text-slate-500 focus:outline-none'
    const closeBtnClass = getUi('closeBtn', closeBtnDefault)
    const closeBtnFloatingDefault =
      'absolute top-2 right-2 z-10 text-slate-400 hover:text-slate-500 focus:outline-none'
    const closeBtnFloatingClass = getUi('closeBtnFloating', closeBtnFloatingDefault)

    const closeIconDefault = 'h-6 w-6'
    const closeIconClass = getUi('closeIcon', closeIconDefault)
    const closeIconSmallDefault = 'h-5 w-5'
    const closeIconSmallClass = getUi('closeIconSmall', closeIconSmallDefault)

    const contentWrapperDefault = 'p-4 overflow-auto flex-1 text-slate-600 dark:text-slate-300'
    const contentWrapperClass = getUi('content', contentWrapperDefault)

    const loadingWrapDefault = 'flex justify-center items-center h-full'
    const loadingWrapClass = getUi('loadingWrap', loadingWrapDefault)
    const spinnerDefault = 'animate-spin h-8 w-8 text-blue-500'
    const spinnerClass = getUi('spinner', spinnerDefault)

    const footerDefault =
      'p-4 border-t border-slate-200 dark:border-slate-700 flex justify-end space-x-3 bg-slate-50 dark:bg-slate-800'
    const footerClass = getUi('footer', footerDefault)

    const template = `
      <template x-for="(item, index) in $store.box.items" :key="item.id">
          <div :id="'uno-box-outer-container' + item.id"
               :class="$store.box.getPositionClass(item.position)"
               :style="\`z-index: \${${zBase} + (index * ${zStep})}\`"
               role="dialog"
               aria-modal="true"
               @keydown.escape.window="if(item.modal) $store.box.close(item.id)">
               
               <!-- Overlay -->
               <div :id="'uno-box-overlay' + item.id"
                    class="${overlayClass}"
                    :style="\`opacity: \${item.overlayOpacity}\`"
                    x-show="item.modal"
                    x-transition:enter="ease-out duration-300"
                    x-transition:enter-start="opacity-0"
                    x-transition:enter-end="opacity-100"
                    x-transition:leave="ease-in duration-200"
                    x-transition:leave-start="opacity-100"
                    x-transition:leave-end="opacity-0"></div>

               <!-- Box Container -->
               <div :id="'uno-box-container' + item.id"
                    class="${containerClass}"
                    :class="[item.widthClass, item.heightClass, item.containerClass]"
                    :style="\`\${item.width ? 'width: ' + item.width + 'px;' : ''} \${item.height ? 'height:' + item.height + 'px' : ''}\`"
                    x-show="item.visible"
                    x-transition:enter="ease-out duration-300"
                    x-transition:enter-start="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95"
                    x-transition:enter-end="opacity-100 translate-y-0 sm:scale-100"
                    x-transition:leave="ease-in duration-200"
                    x-transition:leave-start="opacity-100 translate-y-0 sm:scale-100"
                    x-transition:leave-end="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95"
                    x-uno-trap="item.visible"
                    tabindex="-1">
                    
                    <!-- Header -->
                    <div :id="'uno-box-header' + item.id" x-show="item.title" class="${headerClass}">
                        <h3 :id="'uno-box-title' + item.id" class="${titleClass}" x-text="item.title"></h3>
                        <button :id="'uno-box-close' + item.id"
                                x-show="item.showCloseButton"
                                @click="$store.box.close(item.id)"
                                class="${closeBtnClass}">
                            <span class="sr-only">Close</span>
                            <svg class="${closeIconClass}" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                            </svg>
                        </button>
                    </div>
                    
                    <!-- Close button if no title -->
                    <button x-show="!item.title && item.showCloseButton" 
                            :id="'uno-box-close' + item.id"
                            @click="$store.box.close(item.id)" 
                            class="${closeBtnFloatingClass}">
                        <span class="sr-only">Close</span>
                        <svg class="${closeIconSmallClass}" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                        </svg>
                    </button>

                    <!-- Content -->
                    <div :id="'uno-box-wrapper' + item.id" class="${contentWrapperClass}">
                        <div x-show="item.loading" class="${loadingWrapClass}">
                            <svg class="${spinnerClass}" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                            </svg>
                        </div>
                        <div :id="'uno-box-content' + item.id" x-show="!item.loading" x-html="item.content"></div>
                    </div>

                    <!-- Footer -->
                    <div :id="'uno-box-footer' + item.id" x-show="item.showOk || item.showCancel" class="${footerClass}">
                        <button :id="'uno-box-cancel' + item.id"
                                x-show="item.showCancel"
                                @click="item.onCancel ? (item.onCancel(), $store.box.close(item.id)) : $store.box.close(item.id)"
                                x-text="item.cancelText"
                                class="${cancelClass}">
                        </button>
                        <button :id="'uno-box-ok' + item.id"
                                x-show="item.showOk"
                                @click="async () => {
                                    let shouldClose = true;
                                    if (item.onOk) {
                                        const result = item.onOk();
                                        shouldClose = result instanceof Promise ? await result : result;
                                    }
                                    if (shouldClose !== false) $store.box.close(item.id);
                                }"
                                x-text="item.okText"
                                class="${okClass}">
                        </button>
                    </div>
               </div>
          </div>
      </template>
    `
    rootEl = domManager.ensureContainer('uno-box-container', template, {
      'x-data': '{}'
    })
  }

  // Eager injection removed for P1-10 (On-demand)
  // Container is injected when open() is called.

  // HTMX Integration: Re-initialize Alpine on swapped content
  // This ensures that if HTMX is used within a box (or anywhere), Alpine components are initialized.
  if (typeof window !== 'undefined') {
    if (once(document, 'uno:box:htmx-afterSwap')) {
      document.addEventListener('htmx:afterSwap', (event) => {
        if (event.detail.elt && Alpine) {
          Alpine.initTree(event.detail.elt)
        }
      })
    }
  }

  // Expose Public API to Uno Global
  Uno.box = {
    dialog: (opts) => Alpine.store('box').open('dialog', opts),
    confirm: (msg, opts) => Alpine.store('box').confirm(msg, opts),
    alert: (msg, opts) => Alpine.store('box').alert(msg, opts),
    prompt: (msg, def, opts) => Alpine.store('box').prompt(msg, def, opts),
    close: (id) => Alpine.store('box').close(id),
    closeAll: () => Alpine.store('box').closeAll()
  }
}
