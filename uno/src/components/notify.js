/**
 * Uno - Notify Component (Store-Driven Refactor)
 *
 * @module uno.components.notify
 * @version 0.2.1
 * @license MIT
 */

import { getFixedPositionClass, getStackingConfig } from '../utils/position'
import { domManager } from '../core/dom'

export default (Alpine, options = {}) => {
  const initConfig = options.notify || {}
  const getConfig = () =>
    (typeof window !== 'undefined' && window.Uno?.options?.notify) || initConfig

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
        position = config.position || 'top-right',
        closeIndicator = false,
        showCloseButton = true
      } = callOptions

      const id = this.nextId++

      const item = {
        id,
        text,
        type,
        position,
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

    // Loading Indicator Methods
    loading(show) {
      if (show) {
        this.activeRequests++
      } else {
        this.activeRequests = Math.max(0, this.activeRequests - 1)
      }
      this.ensureContainer()
    },

    // Progress Bar Methods
    progress(show) {
      if (show) {
        this.activeProgress++
      } else {
        this.activeProgress = Math.max(0, this.activeProgress - 1)
      }
      this.ensureContainer()
    },

    // Container Injection
    ensureContainer() {
      injectContainer()
    }
  })

  // 2. Helper to inject the global container
  const injectContainer = () => {
    // 1. Inject Styles
    domManager.injectStyles(
      'uno-notify-styles',
      `
      @keyframes uno-progress-indeterminate {
        0% { left: -35%; width: 35%; }
        100% { left: 100%; width: 35%; }
      }
      .uno-progress-anim {
        animation: uno-progress-indeterminate 1.5s infinite linear;
      }
    `
    )

    // 2. Inject Container
    // We define regions for positions.
    const regions = [
      'top-right',
      'top-left',
      'top-center',
      'center',
      'bottom-right',
      'bottom-left',
      'bottom-center'
    ]

    const regionHtml = regions
      .map((pos) => {
        const fixedClass = getFixedPositionClass(pos)
        const { flexDir, alignItems, spaceY } = getStackingConfig(pos)
        // Backward compatibility for 'top' and 'bottom'
        let filter = `i.position === '${pos}'`
        if (pos === 'top-center') filter += " || i.position === 'top'"
        if (pos === 'bottom-center') filter += " || i.position === 'bottom'"

        return `
      <!-- ${pos} -->
      <div class="fixed ${fixedClass} z-[30000] p-4 flex ${flexDir} ${spaceY} ${alignItems} pointer-events-none">
          <template x-for="item in $store.notify.items.filter(i => ${filter})" :key="item.id">
              ${renderItem()}
          </template>
      </div>`
      })
      .join('\n')

    const template = `
      <!-- Progress Bar (Top Fixed) -->
      <div id="uno-progress-bar"
           x-show="$store.notify.activeProgress > 0"
           class="fixed top-0 left-0 w-full h-1 bg-blue-100 z-[30200]"
           style="display: none;">
           <div class="absolute top-0 h-full bg-blue-600 uno-progress-anim"></div>
      </div>

      <!-- Loading Indicator (Top Right Fixed) -->
      <div id="uno-loading-indicator" 
           x-show="$store.notify.activeRequests > 0"
           class="fixed top-4 right-4 z-[30150] bg-blue-100 border-blue-500 text-blue-700 p-2 rounded-full shadow-md flex items-center dark:bg-blue-900 dark:text-blue-100"
           style="display: none;">
           <svg class="animate-spin -ml-1 mr-2 h-5 w-5 text-blue-500 dark:text-blue-400" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"><circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle><path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path></svg>
           <span>Loading</span>
      </div>

      ${regionHtml}
    `

    domManager.ensureContainer('uno-notify-container', template, {
      'x-data': '{}'
    })
  }

  const renderItem = () => {
    // Note: pointer-events-auto needed because container is pointer-events-none
    return `
    <div x-show="true"
         data-testid="notify-item"
         x-transition:enter="transition ease-out duration-300"
         x-transition:enter-start="opacity-0 transform scale-90 translate-y-2"
         x-transition:enter-end="opacity-100 transform scale-100 translate-y-0"
         x-transition:leave="transition ease-in duration-300"
         x-transition:leave-start="opacity-100 transform scale-100 translate-y-0"
         x-transition:leave-end="opacity-0 transform scale-90 translate-y-2"
         class="pointer-events-auto relative w-80 max-w-full border-l-4 pl-3 pr-3 py-4 rounded shadow-sm dark:shadow-md overflow-hidden"
         :class="{
            'bg-green-100 dark:bg-green-800 border-green-500 text-green-700 dark:text-green-100': item.type === 'success',
            'bg-yellow-100 dark:bg-yellow-800 border-yellow-500 text-yellow-700 dark:text-yellow-100': item.type === 'warning',
            'bg-red-100 dark:bg-red-800 border-red-500 text-red-700 dark:text-red-100': item.type === 'error',
            'bg-blue-100 dark:bg-blue-800 border-blue-500 text-blue-700 dark:text-blue-100': item.type === 'info'
         }"
         role="alert">
         
         <div class="flex items-start justify-between">
            <div x-text="item.text" class="text-sm font-medium"></div>
            <button x-show="item.showCloseButton" 
                    @click="$store.notify.close(item.id)" 
                    class="ml-4 opacity-50 hover:opacity-100 transition-opacity focus:outline-none">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path></svg>
            </button>
         </div>
    </div>
    `
  }

  // Eager injection removed for P1-10 (On-demand)
  // Container is injected when show() is called.
}
