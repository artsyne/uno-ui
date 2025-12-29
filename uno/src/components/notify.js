/**
 * Uno - Notify Component
 *
 * Store-First, Directive-Driven notification system.
 *
 * @module uno.components.notify
 */

import { Uno } from '../core/uno'
import { log } from '../utils/log'
import { getZIndex } from '../core/layer'
import { DomManager } from '../utils/dom'
import { getFixedPositionClass, getStackingConfig } from '../utils/position'

export default (Alpine, options = {}) => {
  // 1. Config Helpers
  const getConfig = (instanceOpts = {}) => {
    // Level 1: Plugin Defaults (handled by Uno.registerPluginDefaults if used, or local default)
    // Level 2: Global User Config (plugins.notify)
    // Level 3: Direct Global Config (notify)
    // Level 4: Instance Config (instanceOpts)

    // We use Uno.getConfig helper if available, or manual merge
    return Uno.getConfig('notify', instanceOpts)
  }

  // 2. Alpine Store
  const notifyStore = {
    groups: {
      'top-left': [],
      'top-center': [],
      'top-right': [],
      'center-left': [],
      center: [],
      'center-right': [],
      'bottom-left': [],
      'bottom-center': [],
      'bottom-right': []
    },
    nextId: 1,
    progressRequests: 0,
    progressVisible: false,
    progressPercent: 0,
    _progressInterval: null,

    init() {
      // Ensure container exists when store is initialized
      this.ensureContainer()
    },

    // HTMX Event Handlers
    onContentLoaded(el, item) {
      this._onHtmxEnd(item)

      // Update the content of the notification item with the loaded HTML
      const loadedContent = el.innerHTML || el.textContent || ''
      const notification = this._findNotificationById(item.id)
      if (notification) {
        // Update content property which is used in x-html binding
        notification.content = loadedContent
      }

      if (item.onContentLoad) {
        item.onContentLoad(el)
      }
      if (window.Alpine && window.Alpine.initTree) {
        window.Alpine.initTree(el)
      }
    },

    // Helper method to find notification by ID
    _findNotificationById(id) {
      for (const position of Object.keys(this.groups)) {
        const group = this.groups[position]
        const notification = group.find((item) => item.id === id)
        if (notification) return notification
      }
      return null
    },

    onHtmxError(error, item) {
      this._onHtmxEnd(item)
      console.error('HTMX Error:', error)
      this.error('Content load failed') // Could be localized if needed
      if (item.onError) {
        item.onError(error)
      }
    },

    _onHtmxBegin(item) {
      if (!item || item.__unoHtmxActive) return

      item.__unoHtmxActive = true
      this.progressRequests++
      if (this.progressRequests === 1) this.startProgress()
    },

    _onHtmxEnd(item) {
      if (!item || !item.__unoHtmxActive) return

      item.__unoHtmxActive = false
      if (this.progressRequests > 0) this.progressRequests--
      if (this.progressRequests === 0) this.stopProgress()
    },

    startProgress() {
      if (this._progressInterval) clearInterval(this._progressInterval)
      this.progressVisible = true
      this.progressPercent = 0

      // Fast loop from 0% to 100% then reset to 0% for continuous effect
      this._progressInterval = setInterval(() => {
        if (this.progressPercent >= 100) {
          // When it reaches 100%, instantly reset to 0% for next cycle
          this.progressPercent = 0
        } else {
          this.progressPercent += 10 // Increase by 10% for faster animation
        }
      }, 100) // Update every 100ms for smooth animation
    },

    stopProgress() {
      if (this._progressInterval) clearInterval(this._progressInterval)
      // Complete the progress smoothly
      this.progressPercent = 100
      setTimeout(() => {
        this.progressVisible = false
        setTimeout(() => {
          this.progressPercent = 0
        }, 200) // Reset after fade out
      }, 300)
    },

    getContainerStyle(position) {
      // Use Layer Manager for Z-Index (TOAST layer)
      // If user provided a custom base in options, use that, otherwise use core layer
      const zBase = Uno.options?.ui?.layers?.notifyBase || getZIndex('TOAST')
      return {
        zIndex: zBase
      }
    },

    getContainerClass(position) {
      const base = Uno.getStyle('notify', 'container')
      const posClass = getFixedPositionClass(position)
      const stack = getStackingConfig(position)

      // We combine:
      // 1. Base Uno styles (fixed flex gap-2 p-4 pointer-events-none)
      // 2. Fixed positioning (top-0 right-0 etc)
      // 3. Flex direction (flex-col etc)
      // 4. Alignment (items-end etc)

      return [base, posClass, stack.flexDir, stack.alignItems].join(' ')
    },

    ensureContainer() {
      const id = 'uno-notify-container'

      // 1. Notify Container
      // We use a template string for the Alpine component
      const template = `
          <template x-for="(items, position) in $store.uno.notify.groups" :key="position">
            <div :class="$store.uno.notify.getContainerClass(position)"
                 :style="$store.uno.notify.getContainerStyle(position)">
              
              <template x-for="item in items" :key="item.id">
                <div :class="[
                    Uno.getStyle('notify', 'item'),
                    Uno.getStyle('notify', 'types.' + item.type)
                  ]"
                  role="alert"
                  aria-live="assertive"
                  x-show="item.isVisible"
                  x-transition:enter="transition ease-out duration-300"
                  x-transition:enter-start="opacity-0 transform scale-90"
                  x-transition:enter-end="opacity-100 transform scale-100"
                  x-transition:leave="transition ease-in duration-200"
                  x-transition:leave-start="opacity-100 transform scale-100"
                  x-transition:leave-end="opacity-0 transform scale-90"
                  :data-id="item.id"
                  data-testid="notify-item"
                 >
                  
                  <div :class="Uno.getStyle('notify', 'content')">
                    <!-- Standard Content -->
                    <template x-if="!item.hxUrl">
                      <div :class="Uno.getStyle('notify', 'row')">
                         <!-- Icon -->
                         <div x-show="item.icon" :class="Uno.getStyle('notify', 'iconWrap')">
                            <i :class="[item.icon, Uno.getStyle('notify', 'icons.' + item.type)]"></i>
                         </div>
                         
                         <!-- Content -->
                         <div :class="Uno.getStyle('notify', 'textWrap')">
                            <p :class="Uno.getStyle('notify', 'text')" x-text="item.message"></p>
                         </div>
                         
                         <!-- Close Button -->
                         <div x-show="item.closable" :class="Uno.getStyle('notify', 'closeWrap')">
                            <button type="button" 
                                    @click="$store.uno.notify.remove(item.id)"
                                    :class="Uno.getStyle('notify', 'closeBtn')">
                              <span class="sr-only">Close</span>
                              <svg :class="Uno.getStyle('notify', 'closeIcon')" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                              </svg>
                            </button>
                         </div>
                      </div>
                    </template>

                    <!-- HTMX Content -->
                    <template x-if="item.hxUrl">
                       <div class="w-full"
                            x-init="
                              $el.setAttribute('hx-get', item.hxUrl);
                              $el.setAttribute('hx-trigger', item.hxTrigger || 'load');
                              $el.setAttribute('hx-swap', item.hxSwap || 'innerHTML');
                              if (item.hxHeaders) $el.setAttribute('hx-headers', JSON.stringify(item.hxHeaders));
                              htmx.process($el);
                            "
                            @htmx:beforeRequest.stop="$store.uno.notify._onHtmxBegin(item)"
                            @htmx:afterRequest.stop="$store.uno.notify._onHtmxEnd(item)"
                            @htmx:afterSwap.stop="$store.uno.notify.onContentLoaded($el, item)"
                            @htmx:responseError.stop="$store.uno.notify.onHtmxError($event.detail, item)"
                            @htmx:sendError.stop="$store.uno.notify.onHtmxError($event.detail, item)"
                            @htmx:timeout.stop="$store.uno.notify.onHtmxError($event.detail, item)"
                       >
                          <div x-html="item.content || ''"></div>
                       </div>
                    </template>
                  </div>

                </div>
              </template>
            </div>
          </template>
        `

      const container = DomManager.ensureContainer(id, template, {
        'x-data': '{}' // Independent Alpine component
      })

      // Manually init Alpine on the new element if it was just created
      if (container && Alpine.initTree) {
        Alpine.initTree(container)
      }

      // 3. Global Progress Bar
      const progressId = 'uno-global-progress'
      // Use a high z-index (e.g. same as loading or slightly below/above)
      const progressZ = getZIndex('MODAL') + 2100

      const progressTemplate = `
          <div :class="Uno.getStyle('notify', 'progressWrap')"
               style="z-index: ${progressZ}"
               x-show="$store.uno.notify.progressVisible"
               x-transition:enter="transition-opacity duration-200"
               x-transition:enter-start="opacity-0"
               x-transition:enter-end="opacity-100"
               x-transition:leave="transition-opacity duration-200"
               x-transition:leave-start="opacity-100"
               x-transition:leave-end="opacity-0">
              <div :class="Uno.getStyle('notify', 'progressBar')"
                   :style="{ width: $store.uno.notify.progressPercent + '%' }">
              </div>
          </div>
      `

      const progressEl = DomManager.ensureContainer(progressId, progressTemplate, {
        'x-data': '{}'
      })

      if (progressEl && Alpine.initTree) {
        Alpine.initTree(progressEl)
      }
    },

    show(message, opts = {}) {
      const config = getConfig(opts)
      const id = this.nextId++

      // i18n support: if message is a key, try to translate it
      let text = message
      if (Alpine.store('i18n') && typeof Alpine.store('i18n').t === 'function') {
        text = Alpine.store('i18n').t(message)
      }

      // Default duration 3000ms, unless explicitly set to 0 or another value
      const duration = config.duration !== undefined ? config.duration : 3000

      log.info('[Notify] Show:', { id, message: text, duration, config })

      const notification = {
        id,
        message: text,
        isVisible: true, // For x-show
        closable: true,
        duration,
        content: config.hxUrl ? '' : undefined, // Initialize content for HTMX-loaded content
        ...config
      }

      const position = config.position || 'top-right'
      if (!this.groups[position]) {
        // Fallback or create
        if (!this.groups['top-right']) this.groups['top-right'] = []
        this.groups['top-right'].push(notification)
      } else {
        this.groups[position].push(notification)
      }

      // Auto dismiss
      if (duration > 0) {
        setTimeout(() => {
          this.remove(id)
        }, duration)
      }

      return id
    },

    remove(id) {
      Object.keys(this.groups).forEach((position) => {
        // We find the item and set visible to false to trigger leave transition
        // Then remove it after transition
        const idx = this.groups[position].findIndex((item) => item.id === id)
        if (idx !== -1) {
          this.groups[position][idx].visible = false
          // Wait for transition (default 200ms in template)
          setTimeout(() => {
            // Re-fetch index in case it shifted
            const currentIdx = this.groups[position].findIndex((item) => item.id === id)
            if (currentIdx !== -1) {
              this.groups[position].splice(currentIdx, 1)
            }
          }, 300)
        }
      })
    },

    // 3. Proxy Methods
    success(message, opts = {}) {
      return this.show(message, { type: 'success', ...opts })
    },
    error(message, opts = {}) {
      return this.show(message, { type: 'error', ...opts })
    },
    warning(message, opts = {}) {
      return this.show(message, { type: 'warning', ...opts })
    },
    info(message, opts = {}) {
      return this.show(message, { type: 'info', ...opts })
    }

    // 4. Loading Helpers - Removed to simplify the component
    // Loading indicator functionality has been removed to reduce complexity
  }

  // Register to Alpine.store('uno') namespace
  if (!Alpine.store('uno')) {
    Alpine.store('uno', { notify: notifyStore })
  } else {
    // If uno store already exists, attach notify
    Alpine.store('uno').notify = notifyStore
  }

  // Ensure initialization (container creation)
  notifyStore.init()

  // 3. Directive Registration
  Alpine.directive('uno-notify', (el, { expression }, { evaluate }) => {
    // Usage: <button x-uno-notify="'Message'" data-uno-type="success">
    // Usage: <button x-uno-notify="['Message', { type: 'success' }]">

    const handler = () => {
      const store = Alpine.store('uno')?.notify
      if (!store) return

      // Parse attributes
      const attrConfig = {
        type: el.getAttribute('data-uno-type') || 'info',
        position: el.getAttribute('data-uno-position'),
        duration: el.getAttribute('data-uno-duration')
          ? Number(el.getAttribute('data-uno-duration'))
          : undefined,
        closable: el.getAttribute('data-uno-closable') !== 'false'
      }

      // Clean undefined
      Object.keys(attrConfig).forEach(
        (key) => attrConfig[key] === undefined && delete attrConfig[key]
      )

      // Message can come from expression or attribute
      let message = ''
      let exprConfig = {}

      if (expression) {
        const result = evaluate(expression)
        if (Array.isArray(result)) {
          message = result[0]
          exprConfig = result[1] || {}
        } else if (typeof result === 'object' && result !== null) {
          // Fallback if user passed object (maybe just config?)
          // Or if they used comma operator (unlikely to get message)
          exprConfig = result
        } else {
          message = String(result)
        }
      } else {
        message = el.getAttribute('data-uno-message') || ''
      }

      store.show(message, { ...attrConfig, ...exprConfig })
    }

    el.addEventListener('click', handler)
  })
}
