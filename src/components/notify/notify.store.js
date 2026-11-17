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

  /**
   * Notification queue
   * @type {Array}
   */
  notifications: [],

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
    log.error('HTMX Error:', error)
    this.error('Content load failed') // Could be localized if needed
    if (item.onError) {
      item.onError(error)
    }
  },

  _onHtmxBegin(item) {
    if (!item || item.__unoHtmxActive) return

    item.__unoHtmxActive = true
    this.beginProgress()
  },

  _onHtmxEnd(item) {
    if (!item || !item.__unoHtmxActive) return

    item.__unoHtmxActive = false
    this.endProgress()
  },

  /**
   * Begin tracking an active request. Shows the progress bar on the first
   * concurrent request. Public API — always pair with endProgress().
   */
  beginProgress() {
    this.progressRequests++
    if (this.progressRequests === 1) this.startProgress()
  },

  /**
   * End tracking of an active request. Hides the progress bar when the last
   * concurrent request finishes. Safe to call more often than beginProgress().
   */
  endProgress() {
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
    // 2. Fixed positioning (top-0 end-0 etc)
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
      ...config,
      // 可暂停的自动关闭计时——保留剩余时长，hover/focus 期间不走表
      _remaining: duration,
      _timer: null,
      _startedAt: 0,
      _paused: { hover: false, focus: false }
    }

    const position = config.position || 'top-right'
    if (!this.groups[position]) {
      // Fallback or create
      if (!this.groups['top-right']) this.groups['top-right'] = []
      this.groups['top-right'].push(notification)
    } else {
      this.groups[position].push(notification)
    }

    // Auto dismiss（经响应式代理启动，pause/resume 与 remove 操作同一对象）
    this._startDismiss(this._findNotificationById(id))

    return id
  },

  /**
   * 启动自动关闭计时。仅对 duration>0、可见、未暂停且未在计时的通知生效；
   * 以 _remaining 为延时，暂停后恢复时续走剩余时长而非重置。
   */
  _startDismiss(item) {
    if (!item || !(item.duration > 0) || !item.isVisible || item._timer) return
    if (item._paused.hover || item._paused.focus) return
    item._startedAt = Date.now()
    item._timer = setTimeout(() => {
      item._timer = null
      this.remove(item.id)
    }, item._remaining)
  },

  /**
   * 暂停自动关闭（鼠标悬停或键盘聚焦）。结算已走时长并保留剩余，
   * 使读者/键盘用户有充足时间阅读与操作，不被超时打断。
   * @param {number} id
   * @param {'hover'|'focus'} reason
   */
  pause(id, reason) {
    const item = this._findNotificationById(id)
    if (!item) return
    if (reason) item._paused[reason] = true
    if (item._timer) {
      clearTimeout(item._timer)
      item._timer = null
      item._remaining = Math.max(0, item._remaining - (Date.now() - item._startedAt))
    }
  },

  /**
   * 恢复自动关闭。仅当 hover 与 focus 两种暂停都解除后续走剩余时长。
   * @param {number} id
   * @param {'hover'|'focus'} reason
   */
  resume(id, reason) {
    const item = this._findNotificationById(id)
    if (!item) return
    if (reason) item._paused[reason] = false
    if (!item._paused.hover && !item._paused.focus) this._startDismiss(item)
  },

  remove(id) {
    // 关闭前清理计时并把焦点合理交接——若焦点正落在被关通知内，
    // 移到另一条可见通知的操作按钮，避免 DOM 移除后焦点坠回 body
    const closing = this._findNotificationById(id)
    if (closing && closing._timer) {
      clearTimeout(closing._timer)
      closing._timer = null
    }
    if (typeof document !== 'undefined') {
      const el = document.querySelector(`[data-id="${id}"]`)
      if (el && el.contains(document.activeElement)) {
        const next = Array.from(document.querySelectorAll('[data-testid="notify-item"]'))
          .filter((n) => n !== el)
          .map((n) => n.querySelector('button'))
          .find(Boolean)
        if (next) next.focus()
      }
    }

    Object.keys(this.groups).forEach((position) => {
      // We find the item and set visible to false to trigger leave transition
      // Then remove it after transition
      const idx = this.groups[position].findIndex((item) => item.id === id)
      if (idx !== -1) {
        // The template binds x-show to `item.isVisible` — setting `visible`
        // here would never trigger the leave transition.
        this.groups[position][idx].isVisible = false
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
    // 关键错误默认持久可读（duration 0），不被超时删除；调用方显式传 duration 可覆盖
    return this.show(message, { type: 'error', duration: 0, ...opts })
  },
  warning(message, opts = {}) {
    return this.show(message, { type: 'warning', ...opts })
  },
  info(message, opts = {}) {
    return this.show(message, { type: 'info', ...opts })
  }
})
