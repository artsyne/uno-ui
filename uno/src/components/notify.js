/**
 * Uno - Notify Component (Store-Driven Refactor)
 *
 * @module uno.components.notify
 * @version 0.2.1
 * @license MIT
 */

import { POSITIONS, getFixedPositionClass, getStackingConfig } from '../utils/position'
import { Uno } from '../core/uno'

export default (Alpine, options = {}) => {
  const initConfig = options.notify || {}
  const getConfig = () =>
    (typeof window !== 'undefined' && window.Uno?.options?.notify) || initConfig
  const getNonce = () =>
    Uno.options?.csp?.nonce || (typeof window !== 'undefined' && window.Uno?.options?.csp?.nonce)

  const ensureStyles = () => {
    if (document.getElementById('uno-notify-styles')) return

    const style = document.createElement('style')
    style.id = 'uno-notify-styles'

    const nonce = getNonce()
    if (nonce) style.setAttribute('nonce', nonce)

    style.textContent = `
      @keyframes notify-progress {
        from { width: 100%; }
        to { width: 0%; }
      }

      @keyframes uno-progress-indeterminate {
        0% { transform: translateX(-100%); }
        100% { transform: translateX(250%); }
      }

      #uno-progress-bar .uno-progress-anim {
        animation: uno-progress-indeterminate 1.2s ease-in-out infinite;
        width: 40%;
      }
    `

    document.head.appendChild(style)
  }

  // 1. Store
  Alpine.store('notify', {
    items: [],
    nextId: 1,
    activeRequests: 0,
    activeProgress: 0,

    // Notification Methods
    show(text, callOptions = {}) {
      const config = getConfig()
      const {
        type = 'info',
        timeout = config.duration !== undefined ? config.duration : 3000,
        position = callOptions.position || config.position || 'top-right',
        closeIndicator = false,
        showCloseButton = true
      } = callOptions

      const id = this.nextId++

      const item = {
        id,
        text,
        type,
        position,
        timeout,
        closeIndicator,
        showCloseButton,
        visible: true // For transition
      }

      // Add to start of array (newest first)
      this.items.unshift(item)

      // Auto dismiss
      if (timeout > 0) {
        setTimeout(() => {
          this.close(id)
        }, timeout)
      }

      this.ensureContainer()
    },

    close(id) {
      const item = this.items.find((i) => i.id === id)
      if (item) {
        // Trigger leave transition if handled by x-transition (removed from DOM after)
        // But since we use x-for, removing from array triggers leave transition.
        this.items = this.items.filter((i) => i.id !== id)
      }
    },

    success(message, options = {}) {
      const opts = typeof options === 'number' ? { timeout: options } : options
      this.show(message, { type: 'success', ...opts })
    },

    warning(message, options = {}) {
      const opts = typeof options === 'number' ? { timeout: options } : options
      this.show(message, { type: 'warning', ...opts })
    },

    error(message, options = {}) {
      const opts = typeof options === 'number' ? { timeout: options } : options
      this.show(message, { type: 'error', ...opts })
    },

    info(message, options = {}) {
      const opts = typeof options === 'number' ? { timeout: options } : options
      this.show(message, { type: 'info', ...opts })
    },

    loading(show) {
      if (show) {
        this.activeRequests++
      } else {
        this.activeRequests = Math.max(0, this.activeRequests - 1)
      }
      this.ensureContainer()
    },

    progress(show) {
      if (show) {
        this.activeProgress++
      } else {
        this.activeProgress = Math.max(0, this.activeProgress - 1)
      }
      this.ensureContainer()
    },

    // Style Helpers
    getTypeClass(type) {
      const defaults = {
        success: 'bg-green-50',
        error: 'bg-red-50',
        warning: 'bg-yellow-50',
        info: 'bg-blue-50'
      }
      return Uno.getStyle('notify', type, defaults[type] || 'bg-white')
    },

    getIconClass(type) {
      const defaults = {
        success: 'text-green-400',
        error: 'text-red-400',
        warning: 'text-yellow-400',
        info: 'text-blue-400'
      }
      return Uno.getStyle(
        'notify',
        `icon${type.charAt(0).toUpperCase() + type.slice(1)}`,
        defaults[type] || 'text-gray-400'
      )
    },

    ensureContainer() {
      if (document.getElementById('uno-notify-container')) return

      ensureStyles()

      const container = document.createElement('div')
      container.id = 'uno-notify-container'

      const nonce = getNonce()
      if (nonce) container.setAttribute('nonce', nonce)

      container.setAttribute('x-data', '{}')

      const rounded = Uno.getStyle('base', 'rounded', 'rounded-lg')
      const itemClass = Uno.getStyle(
        'notify',
        'item',
        `pointer-events-auto w-full max-w-sm overflow-hidden ${rounded} shadow-lg ring-1 ring-black ring-opacity-5`
      )

      const buildRegion = (position) => {
        const positionClass = getFixedPositionClass(position)
        const stack = getStackingConfig(position)
        const regionClass = Uno.getStyle(
          'notify',
          'container',
          `fixed z-[30000] flex ${stack.flexDir} ${stack.alignItems} gap-2 p-4 pointer-events-none ${positionClass}`
        )
        return `
          <div class="${regionClass}" x-show="$store.notify.itemsFor('${position}').length">
            <template x-for="item in $store.notify.itemsFor('${position}')" :key="item.id">
              <div
                data-testid="notify-item"
                x-show="item.visible"
                x-transition:enter="transition ease-out duration-300"
                x-transition:enter-start="opacity-0 transform translate-y-2 scale-95"
                x-transition:enter-end="opacity-100 transform translate-y-0 scale-100"
                x-transition:leave="transition ease-in duration-200"
                x-transition:leave-start="opacity-100 transform translate-y-0 scale-100"
                x-transition:leave-end="opacity-0 transform translate-y-2 scale-95"
                class="${itemClass}"
                :class="$store.notify.getTypeClass(item.type)"
              >
                <div class="p-4">
                  <div class="flex items-start">
                    <div class="flex-shrink-0">
                      <template x-if="item.type === 'success'">
                        <svg class="h-6 w-6" :class="$store.notify.getIconClass('success')" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                        </svg>
                      </template>
                      <template x-if="item.type === 'error'">
                        <svg class="h-6 w-6" :class="$store.notify.getIconClass('error')" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z" />
                        </svg>
                      </template>
                      <template x-if="item.type === 'warning'">
                        <svg class="h-6 w-6" :class="$store.notify.getIconClass('warning')" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
                        </svg>
                      </template>
                      <template x-if="item.type === 'info'">
                        <svg class="h-6 w-6" :class="$store.notify.getIconClass('info')" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                        </svg>
                      </template>
                    </div>
                    <div class="ml-3 w-0 flex-1 pt-0.5">
                      <p class="text-sm font-medium text-gray-900" x-text="item.text"></p>
                    </div>
                    <div class="ml-4 flex flex-shrink-0" x-show="item.showCloseButton">
                      <button @click="$store.notify.close(item.id)" class="inline-flex rounded-md bg-white text-gray-400 hover:text-gray-500 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2">
                        <span class="sr-only">Close</span>
                        <svg class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                          <path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd" />
                        </svg>
                      </button>
                    </div>
                  </div>
                </div>
                <div x-show="item.closeIndicator" class="h-1 bg-gray-200 overflow-hidden">
                   <div class="h-full bg-blue-500 origin-left"
                        :style="\`animation: notify-progress \${item.timeout}ms linear forwards\`"></div>
                </div>
              </div>
            </template>
          </div>
        `
      }

      container.innerHTML = `
        <div
          id="uno-progress-bar"
          x-show="$store.notify.activeProgress > 0"
          class="fixed top-0 left-0 right-0 h-1 bg-transparent z-[30001]"
        >
          <div class="uno-progress-anim h-full bg-blue-500"></div>
        </div>
        ${POSITIONS.map(buildRegion).join('')}
      `

      document.body.appendChild(container)
      if (typeof Alpine?.initTree === 'function') {
        Alpine.initTree(container)
      }
    },
    itemsFor(position) {
      const pos = position || getConfig().position || 'top-right'
      return this.items.filter((i) => (i.position || 'top-right') === pos)
    }
  })

  // Expose Public API
  const api = {
    show: (text, options) => Alpine.store('notify').show(text, options),
    success: (msg, options) => Alpine.store('notify').success(msg, options),
    error: (msg, options) => Alpine.store('notify').error(msg, options),
    warning: (msg, options) => Alpine.store('notify').warning(msg, options),
    info: (msg, options) => Alpine.store('notify').info(msg, options),
    loading: (show) => Alpine.store('notify').loading(show),
    progress: (show) => Alpine.store('notify').progress(show),
    close: (id) => Alpine.store('notify').close(id),
    closeAll: () => {
      Alpine.store('notify').items = []
    }
  }

  Uno.notify = api
  // Fallback for UMD/Global usage
  if (typeof window !== 'undefined' && window.Uno) {
    window.Uno.notify = api
  }
}
