import { Uno } from '../../core/uno'
import { log } from '../../utils/log'
import { getZIndex } from '../../core/layer'
import { DomManager } from '../../utils/dom'
import { getFixedPositionClass, getStackingConfig } from '../../utils/position'
import { getNotifyTemplate, getProgressTemplate } from './notify.template'

export const createNotifyStore = (Alpine) => ({
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

  getContainerStyle() {
    // Use Layer Manager for Z-Index (NOTIFY layer)
    // If user provided a custom base in options, use that, otherwise use core layer
    const zBase = Uno.options?.ui?.layers?.notify || getZIndex('NOTIFY')
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
    const template = getNotifyTemplate()

    const container = DomManager.ensureContainer(id, template, {
      'x-data': '{}' // Independent Alpine component
    })

    // Manually init Alpine on the new element if it was just created
    if (container && Alpine.initTree) {
      Alpine.initTree(container)
    }

    // 3. Global Progress Bar
    const progressId = 'uno-global-progress'
    const progressTemplate = getProgressTemplate()

    const progressEl = DomManager.ensureContainer(progressId, progressTemplate, {
      'x-data': '{}'
    })

    if (progressEl && Alpine.initTree) {
      Alpine.initTree(progressEl)
    }
  },

  show(message, opts = {}) {
    this.ensureContainer()
    const config = Uno.getConfig('notify', opts)
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
})
