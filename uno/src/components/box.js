/**
 * Uno - Box Component (Store-Driven Refactor)
 *
 * @module uno.components.box
 * @version 0.5.0
 * @license MIT
 */

import { trapFocus } from '../utils/a11y'
import { Uno } from '../core/uno'

export default (Alpine, options = {}) => {
  const initConfig = options.plugins?.box || options.box || {}
  const getConfig = () => {
    const opts = typeof window !== 'undefined' ? window.Uno?.options : Uno.options
    return opts?.plugins?.box || opts?.box || initConfig
  }
  const getUi = (key) => Uno.getStyle('box', key)

  // Register Focus Trap Directive
  Alpine.directive('uno-trap', (el, { expression: _expression }, { cleanup }) => {
    void _expression
    const cleanupTrap = trapFocus(el)
    cleanup(cleanupTrap)
  })

  Alpine.directive('uno-box', (el, { expression }, { evaluateLater, cleanup }) => {
    const getExpression = expression ? evaluateLater(expression) : null

    const parseBool = (value) => {
      if (value === null || value === undefined) return undefined
      const v = String(value).trim().toLowerCase()
      if (v === 'true' || v === '1') return true
      if (v === 'false' || v === '0') return false
      return undefined
    }

    const parseNumber = (value) => {
      if (value === null || value === undefined || value === '') return undefined
      const n = Number(value)
      return Number.isFinite(n) ? n : undefined
    }

    const buildFromAttrs = () => {
      const type =
        el.getAttribute('data-uno-box') || el.getAttribute('data-uno-box-type') || 'dialog'
      const title = el.getAttribute('data-uno-title') || ''
      const content = el.getAttribute('data-uno-content') || ''
      const hxUrl = el.getAttribute('data-uno-hx-url') || undefined
      const position = el.getAttribute('data-uno-position') || undefined
      const modal = parseBool(el.getAttribute('data-uno-modal'))
      const trustedHTML = parseBool(el.getAttribute('data-uno-trusted-html'))
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
        trustedHTML,
        width,
        height,
        widthClass,
        heightClass
      }
    }

    const onClick = (e) => {
      e.preventDefault()
      const store = Alpine.store('box')
      if (!store) return

      if (getExpression) {
        let evaluated
        getExpression((value) => {
          evaluated = value
        })
        if (evaluated && typeof evaluated === 'object') {
          store.open(evaluated)
          return
        }
      }

      store.open(buildFromAttrs())
    }

    el.addEventListener('click', onClick)
    cleanup(() => el.removeEventListener('click', onClick))
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
        const wrapperClass = getUi('promptWrapper')
        const textClass = getUi('promptText')
        const inputClass = getUi('promptInput')

        const content = `
          <div class="${wrapperClass}">
            <p class="${textClass}">${message}</p>
            <input id="${id}" type="text" class="${inputClass}" value="${defaultValue}">
          </div>
        `
        this.open('confirm', {
          title: options.title || 'Prompt',
          content,
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
        position: options.position || config.defaultPosition || 'center',
        modal: options.modal !== undefined ? options.modal : true,
        showCloseButton: options.showCloseButton !== undefined ? options.showCloseButton : true,
        showOkButton: options.showOkButton || false,
        showCancelButton: options.showCancelButton || false,
        okText: options.okText || config.okText || 'OK',
        cancelText: options.cancelText || config.cancelText || 'Cancel',
        onOk: options.onOk,
        onCancel: options.onCancel,
        onClose: options.onClose,
        hxUrl: options.hxUrl,
        hxHeaders: options.hxHeaders,

        // Internal state
        isOpen: false,

        // Methods
        close() {
          this.isOpen = false
          setTimeout(() => {
            Alpine.store('box').remove(this.id)
            if (this.onClose) this.onClose()
          }, 300) // Animation duration
        }
      }

      this.items.push(item)
      this.ensureContainer()

      // Wait for Alpine to render
      setTimeout(() => {
        // Find the reactive item in the store
        const reactiveItem = this.items.find((i) => i.id === item.id)
        if (reactiveItem) reactiveItem.isOpen = true

        // Focus management could go here
      }, 50)

      return item
    },

    remove(id) {
      this.items = this.items.filter((i) => i.id !== id)
    },

    closeAll() {
      this.items = []
    },

    ensureContainer() {
      const id = 'uno-box-manager'
      if (typeof document !== 'undefined' && document.getElementById(id)) return

      const zBase = Uno.options?.ui?.layers?.boxBase || 10000

      const attributes = {
        'x-data': '',
        class: 'relative',
        style: `z-index: ${zBase}; display: none;`,
        'x-show': '$store.box.items.length > 0',
        role: 'dialog',
        'aria-modal': 'true'
      }

      const template = `
             <!-- Backdrop -->
             <div :class="Uno.getStyle('box', 'overlay')"
                  x-show="$store.box.items.length > 0"
                  x-transition:enter="ease-out duration-300"
                  x-transition:enter-start="opacity-0"
                  x-transition:enter-end="opacity-100"
                  x-transition:leave="ease-in duration-200"
                  x-transition:leave-start="opacity-100"
                  x-transition:leave-end="opacity-0"></div>

             <div :class="Uno.getStyle('box', 'dialogWrapper')">
                <template x-for="item in $store.box.items" :key="item.id">
                   <div :id="'uno-box-outer-container' + item.id"
                        :class="Uno.getStyle('box', 'itemWrapper')"
                        x-show="item.isOpen">

                        <div :id="'uno-box-container' + item.id"
                             :class="[Uno.getStyle('box', 'container'), item.widthClass]"
                             :style="item.width ? 'width:' + item.width + 'px; max-width:' + item.width + 'px;' : ''"
                             x-show="item.isOpen">

                             <!-- Content -->
                             <div :class="Uno.getStyle('box', 'content')">
                                 <h3 :class="Uno.getStyle('box', 'title')"
                                     x-text="item.title"
                                     x-show="item.title"></h3>
                                 <div class="mt-2" :id="'uno-box-content' + item.id">
                                    <template x-if="!item.hxUrl">
                                      <div x-html="item.trustedHTML ? item.content : Uno.dom.sanitize(item.content)"></div>
                                    </template>
                                    <template x-if="item.hxUrl">
                                      <div x-init="htmx.process($el)"
                                           :hx-get="item.hxUrl"
                                           hx-trigger="load"
                                           hx-swap="innerHTML"
                                           :hx-headers="item.hxHeaders ? JSON.stringify(item.hxHeaders) : null">
                                        <!-- Loading state could go here -->
                                        <div class="animate-pulse">Loading...</div>
                                      </div>
                                    </template>
                                 </div>
                             </div>

                             <!-- Footer -->
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
    }
  })
}
