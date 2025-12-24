/**
 * Uno - Notify Component (Store-Driven Refactor)
 *
 * @module uno.components.notify
 * @version 0.2.1
 * @license MIT
 */

import { POSITIONS, getFixedPositionClass, getStackingConfig } from '../utils/position'
import { once } from '../utils/dom'
import { Uno } from '../core/uno'

export default (Alpine, options = {}) => {
  const initConfig = options.plugins?.notify || options.notify || {}
  const getConfig = () => {
    const opts = typeof window !== 'undefined' ? window.Uno?.options : Uno.options
    return opts?.plugins?.notify || opts?.notify || initConfig
  }
  const getNonce = () =>
    Uno.options?.csp?.nonce || (typeof window !== 'undefined' && window.Uno?.options?.csp?.nonce)

  const getNotifyTypeClass = (type) => {
    return Uno.getStyle('notify', `types.${type}`, '')
  }

  const getNotifyIconClass = (type) => {
    return Uno.getStyle('notify', `icons.${type}`, '')
  }

  // 1. Store
  let containerEl = null
  const getProtocolEl = () => {
    if (typeof document === 'undefined') return null
    if (containerEl && document.contains(containerEl)) return containerEl
    const el =
      document.getElementById('uno-notify-container') || document.querySelector('[data-uno-notify]')
    if (el) containerEl = el
    return el
  }
  const getProtocolDefaults = () => {
    const el = getProtocolEl()
    if (!el) return {}

    const position = el.getAttribute('data-uno-position') || undefined
    const rawTimeout = el.getAttribute('data-uno-timeout')
    const timeout = rawTimeout !== null && rawTimeout !== '' ? Number(rawTimeout) : undefined
    const type = el.getAttribute('data-uno-type') || undefined

    return {
      position,
      timeout: Number.isFinite(timeout) ? timeout : undefined,
      type
    }
  }

  Alpine.store('notify', {
    items: [],
    nextId: 1,
    activeRequests: 0,
    activeProgress: 0,

    // Notification Methods
    show(text, callOptions = {}) {
      const config = getConfig()
      const protocolDefaults = getProtocolDefaults()
      const {
        type = callOptions.type || config.type || protocolDefaults.type || 'info',
        timeout = callOptions.timeout ??
          (config.duration !== undefined ? config.duration : (protocolDefaults.timeout ?? 3000)),
        position = callOptions.position ||
          config.position ||
          protocolDefaults.position ||
          'top-right',
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
        visible: true,
        indicatorWidth: 100
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
      return getNotifyTypeClass(type)
    },

    getIconClass(type) {
      return getNotifyIconClass(type)
    },

    ensureContainer() {
      if (containerEl && document.contains(containerEl)) return
      containerEl = document.getElementById('uno-notify-container')
      if (containerEl) return
      containerEl = document.querySelector('[data-uno-notify]')
      if (containerEl) return

      let zBase = Number(Uno.options?.ui?.layers?.notifyBase ?? 30000)
      if (!Number.isFinite(zBase)) zBase = 30000

      const container = document.createElement('div')
      container.id = 'uno-notify-container'
      container.setAttribute('data-uno-notify', '')

      const nonce = getNonce()
      if (nonce) container.setAttribute('nonce', nonce)

      container.setAttribute('x-data', '{}')

      const itemClass = Uno.getStyle('notify', 'item')
      const contentWrapClass = Uno.getStyle('notify', 'content')
      const rowClass = Uno.getStyle('notify', 'row')
      const iconWrapClass = Uno.getStyle('notify', 'iconWrap')
      const textWrapClass = Uno.getStyle('notify', 'textWrap')
      const textClass = Uno.getStyle('notify', 'text')

      const closeWrapClass = Uno.getStyle('notify', 'closeWrap')
      const closeBtnClass = Uno.getStyle('notify', 'closeBtn')
      const closeIconClass = Uno.getStyle('notify', 'closeIcon')

      const indicatorTrackClass = Uno.getStyle('notify', 'indicatorTrack')
      const indicatorBarClass = Uno.getStyle('notify', 'indicatorBar')
      const progressWrapClass = Uno.getStyle('notify', 'progressWrap')
      const progressBarClass = Uno.getStyle('notify', 'progressBar')
      const loadingWrapClass = Uno.getStyle('notify', 'loadingWrap')
      const spinnerClass = Uno.getStyle('notify', 'spinner')

      const buildRegion = (position) => {
        const positionClass = getFixedPositionClass(position)
        const stack = getStackingConfig(position)
        const widthClass =
          position === 'top-center' || position === 'bottom-center' || position === 'center'
            ? 'w-full'
            : ''
        const containerBase = Uno.getStyle('notify', 'container')
        const regionClass = `${containerBase} ${stack.flexDir} ${stack.alignItems} ${widthClass} ${positionClass}`
        return `
          <div class="${regionClass}" style="z-index: ${zBase};" x-show="$store.notify.itemsFor('${position}').length">
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
                <div class="${contentWrapClass}">
                  <div class="${rowClass}">
                    <div class="${iconWrapClass}">
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
                    <div class="${textWrapClass}">
                      <p class="${textClass}" x-text="item.text"></p>
                    </div>
                    <div class="${closeWrapClass}" x-show="item.showCloseButton">
                      <button @click="$store.notify.close(item.id)" class="${closeBtnClass}">
                        <span class="sr-only">Close</span>
                        <svg class="${closeIconClass}" viewBox="0 0 20 20" fill="currentColor">
                          <path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd" />
                        </svg>
                      </button>
                    </div>
                  </div>
                </div>
                <div x-show="item.closeIndicator" class="${indicatorTrackClass}">
                   <div
                        class="${indicatorBarClass}"
                        x-init="if(item.timeout > 0) { item.indicatorWidth = 100; $nextTick(() => (item.indicatorWidth = 0)) }"
                        :style="\`width: \${item.indicatorWidth}%; \${item.timeout > 0 ? 'transition: width ' + item.timeout + 'ms linear;' : ''}\`"
                    ></div>
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
          class="${progressWrapClass}"
          style="z-index: ${zBase + 1};"
        >
          <div class="${progressBarClass}"></div>
        </div>
        <div
          id="uno-loading-indicator"
          x-show="$store.notify.activeRequests > 0"
          class="${loadingWrapClass}"
          style="z-index: ${zBase + 2};"
        >
            <svg class="${spinnerClass}" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
        </div>
        ${POSITIONS.map(buildRegion).join('')}
      `

      document.body.appendChild(container)
      containerEl = container
      if (typeof Alpine?.initTree === 'function') {
        Alpine.initTree(container)
      }
    },
    itemsFor(position) {
      const pos = position || 'top-right'
      return this.items.filter((i) => (i.position || 'top-right') === pos)
    }
  })

  if (typeof window !== 'undefined' && typeof document !== 'undefined') {
    if (typeof Uno.htmx?.onRequest === 'function' && Uno.htmx.available()) {
      Uno.htmx.onRequest({
        start: () => {
          const store = Alpine.store('notify')
          store.loading(true)
          store.progress(true)
        },
        end: () => {
          const store = Alpine.store('notify')
          store.loading(false)
          store.progress(false)
        }
      })
    }
  }

  if (typeof window !== 'undefined' && typeof document !== 'undefined') {
    if (once(document, 'uno:notify:event')) {
      document.addEventListener('uno:notify', (event) => {
        const detail = event?.detail || {}
        const text = detail.text || detail.message || detail.title || ''
        if (!text) return
        Alpine.store('notify').show(text, detail)
      })
    }
  }

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
