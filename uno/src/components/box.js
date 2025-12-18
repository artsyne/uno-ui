/**
 * Uno - Box Component (Store-Driven Refactor)
 *
 * @module uno.components.box
 * @version 0.2.0
 * @license MIT
 */

import { lock, unlock } from '../utils/bodyScroll'
import { escapeHTML, once, setHTML } from '../utils/dom'
import { log } from '../utils/log'
import { getFlexAlignmentClass } from '../utils/position'
import { requestHTML, requestData } from '../utils/request'

export default (Alpine) => {
  if (once(Alpine, 'uno:directive:safe-html')) {
    Alpine.directive('safe-html', (el, { expression }, { evaluateLater, effect }) => {
      const get = evaluateLater(expression)
      let lastHTML = null
      let lastTrusted = null

      effect(() => {
        get((value) => {
          let html = ''
          let trusted = false

          if (typeof value === 'string') {
            html = value
          } else if (typeof value === 'object' && value !== null) {
            html = value.html || ''
            trusted = value.trusted === true
          }

          if (html === lastHTML && trusted === lastTrusted) return
          lastHTML = html
          lastTrusted = trusted

          setHTML(el, html, { trusted, alpine: Alpine })
        })
      })
    })
  }

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
      // TODO: Implement prompt UI with input
      return new Promise((resolve) => {
        // Simple fallback or custom content with input?
        // For now, let's use a custom content with input
        const id = 'prompt-' + Math.random().toString(36).substr(2, 9)
        const content = `
          <div class="space-y-3">
            <p>${escapeHTML(message)}</p>
            <input id="${id}" type="text" class="w-full px-3 py-2 border rounded focus:ring-2 focus:ring-blue-500 outline-none" value="${escapeHTML(
          defaultValue
        )}">
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

      if (typeof typeOrOpts === 'object') {
        options = typeOrOpts
        type = options.type || 'dialog'
      } else {
        type = typeOrOpts
      }

      const id = options.boxId || Math.random().toString(36).slice(2, 10)

      // Check if exists (for cache) - though data-driven usually means state is source of truth
      // If boxId provided and already open, do nothing or bring to front?
      // For simplicity, we assume new open.

      const item = {
        id,
        type,
        title: options.title || '',
        content: options.content || '',
        trustedHTML: options.trustedHTML === true,
        width: options.width !== undefined ? options.width : options.widthClass ? null : 400,
        height: options.height,
        widthClass: options.widthClass,
        heightClass: options.heightClass,
        containerClass: options.containerClass,
        modal: options.modal !== false, // default true
        overlayOpacity: options.overlayOpacity ?? 0.5,
        position: options.position || 'center',
        showCloseButton: true, // default true?
        // Buttons
        okText: options.okButtonText,
        cancelText: options.cancelButtonText,
        onOk: options.onOk,
        onCancel: options.onCancel,
        onClose: options.onClose,
        onBeforeClose: options.onBeforeClose,
        onContentLoad: options.onContentLoad, // Callback after content is loaded (ajax/htmx)
        showOk:
          type === 'confirm' || type === 'message' || (type === 'dialog' && options.showOkButton),
        showCancel: type === 'confirm' || (type === 'dialog' && options.showCancelButton),

        // State
        loading: false,
        visible: true
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
          item.trustedHTML = true
        } else {
          log.warn('Box: Target element not found:', options.targetId)
          item.content = '<div class="p-4 text-red-500">Error: Content target not found</div>'
          item.trustedHTML = true
        }
      } else if (options.hxUrl) {
        item.loading = true
        this.add(item)

        try {
          const html = await requestHTML(options.hxUrl, { method: 'GET' })
          const current = this.items.find((i) => i.id === id)
          if (current) {
            current.content = html
            current.trustedHTML = true
            current.loading = false
            if (current.onContentLoad) {
              Alpine.nextTick(() => {
                current.onContentLoad(document.getElementById('uno-box-content' + id))
              })
            }
          }
        } catch (err) {
          log.error('Box: HTMX Load Failed', err)
          const current = this.items.find((i) => i.id === id)
          if (current) {
            current.loading = false
            current.trustedHTML = true
            current.content = '<div class="p-4 text-red-500">Error loading content via HTMX</div>'
          }
        }
        return id
      } else if (options.ajaxUrl) {
        item.loading = true
        // Add immediately to show loading state
        this.add(item)

        try {
          const data = await requestData(options.ajaxUrl)
          const content =
            typeof data === 'object'
              ? `<pre class="bg-gray-100 p-3 rounded overflow-x-auto whitespace-pre-wrap">${JSON.stringify(data, null, 2)}</pre>`
              : String(data)

          // Update item in store
          const current = this.items.find((i) => i.id === id)
          if (current) {
            current.content = content
            current.trustedHTML = true
            current.loading = false

            // Re-initialize Alpine components in the new content
            Alpine.nextTick(() => {
              const el = document.getElementById('uno-box-content' + id)
              if (el) {
                if (current.onContentLoad) {
                  current.onContentLoad(el)
                }
              }
            })
          }
        } catch (e) {
          log.error('Box: Ajax failed', e)
          const current = this.items.find((i) => i.id === id)
          if (current) {
            current.content = `<div class="text-red-500">Error loading content</div>`
            current.trustedHTML = true
            current.loading = false
          }
        }
        return id
      } else if (options.iframeUrl) {
        item.content = `<iframe frameborder="0" class="w-full h-full" src="${escapeHTML(options.iframeUrl)}"></iframe>`
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
        if (item.onBeforeClose) item.onBeforeClose()
        this.items = this.items.filter((i) => i.id !== id)
        if (item.onClose) item.onClose()
      }
    },

    closeAll() {
      this.items = []
    },

    getPositionClass(position) {
      const base = 'fixed inset-0 flex transition-opacity duration-300 p-4 pointer-events-none'
      const alignment = getFlexAlignmentClass(position)
      return `${base} ${alignment}`
    },

    ensureContainer() {
      if (document.getElementById('uno-box-container')) return
      injectContainer()
    }
  })

  // 2. Container Injection
  const injectContainer = () => {
    try {
      const container = document.createElement('div')
      container.id = 'uno-box-container'
      container.setAttribute('x-data', '{}')

      // Template
      // We use x-for to render stacked modals.
      // Each modal is a fixed inset-0 container (to center itself).
      // z-index starts at 100 and increases?
      // Alpine x-for index is available. :style="{ zIndex: 100 + index }"

      const template = `
        <template x-for="(item, index) in $store.box.items" :key="item.id">
            <div :id="'uno-box-outer-container' + item.id"
                 :class="$store.box.getPositionClass(item.position)"
                 :style="\`z-index: \${10000 + index}\`"
                 role="dialog"
                 aria-modal="true">
                 
                 <!-- Overlay: Covers viewport, blocks interaction, does NOT close modal on click -->
                 <div :id="'uno-box-overlay' + item.id"
                      class="absolute inset-0 bg-black transition-opacity duration-300 pointer-events-auto"
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
                      class="relative bg-white dark:bg-gray-900 rounded-lg shadow-2xl ring-1 ring-black/5 overflow-hidden flex flex-col max-h-[90vh] transition-all duration-300 transform pointer-events-auto"
                      :class="[item.widthClass, item.heightClass, item.containerClass]"
                      :style="\`\${item.width ? 'width: ' + item.width + 'px;' : ''} \${item.height ? 'height:' + item.height + 'px' : ''}\`"
                      x-show="item.visible"
                      x-transition:enter="ease-out duration-300"
                      x-transition:enter-start="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95"
                      x-transition:enter-end="opacity-100 translate-y-0 sm:scale-100"
                      x-transition:leave="ease-in duration-200"
                      x-transition:leave-start="opacity-100 translate-y-0 sm:scale-100"
                      x-transition:leave-end="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95">
                      
                      <!-- Header -->
                      <div :id="'uno-box-header' + item.id" x-show="item.title" class="flex justify-between items-center p-4 border-b border-slate-200 dark:border-slate-700">
                          <h3 :id="'uno-box-title' + item.id" class="text-lg font-semibold text-slate-800 dark:text-slate-200" x-text="item.title"></h3>
                          <button :id="'uno-box-close' + item.id" @click="$store.box.close(item.id)" class="text-slate-400 hover:text-slate-500 focus:outline-none">
                              <span class="sr-only">Close</span>
                              <svg class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                              </svg>
                          </button>
                      </div>
                      
                      <!-- Close button if no title -->
                      <button x-show="!item.title && item.showCloseButton" 
                              :id="'uno-box-close' + item.id"
                              @click="$store.box.close(item.id)" 
                              class="absolute top-2 right-2 z-10 text-slate-400 hover:text-slate-500 focus:outline-none">
                          <svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                          </svg>
                      </button>

                      <!-- Content -->
                      <div :id="'uno-box-wrapper' + item.id" class="p-4 overflow-auto flex-1 text-slate-600 dark:text-slate-300">
                          <div x-show="item.loading" class="flex justify-center items-center h-full">
                              <svg class="animate-spin h-8 w-8 text-blue-500" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                              </svg>
                          </div>
                          <div :id="'uno-box-content' + item.id" x-show="!item.loading" x-safe-html="{ html: item.content, trusted: item.trustedHTML }"></div>
                      </div>

                      <!-- Footer -->
                      <div :id="'uno-box-footer' + item.id" x-show="item.showOk || item.showCancel" class="p-4 border-t border-slate-200 dark:border-slate-700 flex justify-end space-x-3 bg-slate-50 dark:bg-slate-800">
                          <button :id="'uno-box-cancel' + item.id"
                                  x-show="item.showCancel"
                                  @click="item.onCancel ? (item.onCancel(), $store.box.close(item.id)) : $store.box.close(item.id)"
                                  x-text="item.cancelText"
                                  class="px-4 py-2 bg-slate-100 dark:bg-slate-700 border border-slate-200 dark:border-slate-600 rounded-md text-sm font-medium text-slate-800 dark:text-slate-200 hover:bg-slate-200 dark:hover:bg-slate-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 transition-colors">
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
                                  class="px-4 py-2 bg-blue-600 border border-transparent rounded-md text-sm font-medium text-white hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500">
                          </button>
                      </div>
                 </div>
            </div>
        </template>
      `
      container.innerHTML = template
      document.body.appendChild(container)
    } catch (e) {
      log.error('Box: Failed to inject container', e)
    }
  }

  // Eager injection
  if (document.body) {
    injectContainer()
  } else {
    if (once(document, 'uno:box:domcontentloaded')) {
      document.addEventListener('DOMContentLoaded', injectContainer)
    }
  }

  // Global delegator for data-box-type (backward compatibility / ease of use)
  if (once(document, 'uno:box:delegator')) {
    document.addEventListener('click', (e) => {
      const t = e.target.closest('[data-box-type]')
      if (!t) return

      const type = t.getAttribute('data-box-type')
      if (!['dialog', 'confirm', 'message'].includes(type)) return

      e.preventDefault()

      // Parse data attributes
      const d = t.dataset
      const opts = {
        title: d.boxTitle,
        content: d.boxContent,
        width: d.boxWidth ? parseInt(d.boxWidth) : undefined,
        height: d.boxHeight ? parseInt(d.boxHeight) : undefined,
        ajaxUrl: d.boxAjaxUrl,
        hxUrl: d.boxHxUrl,
        iframeUrl: d.boxIframeUrl,
        targetId: d.boxTargetId,
        position: d.boxPosition,
        okButtonText: d.boxOkText,
        cancelButtonText: d.boxCancelText,
        overlayOpacity: d.boxOverlayOpacity ? parseFloat(d.boxOverlayOpacity) : undefined,
        showCloseButton: d.boxShowClose !== 'false',
        boxId: d.boxId,
        cache: d.boxCache === 'true',
        trustedHTML: d.boxTrusted === 'true'
      }

      if (d.boxModal === 'false') opts.modal = false

      Alpine.store('box').open(type, opts)
    })
  }

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
}
