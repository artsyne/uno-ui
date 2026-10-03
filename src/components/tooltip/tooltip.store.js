import { Uno } from '../../core/uno'
import { calculatePosition } from '../../utils/position'
import { DomManager } from '../../utils/dom'
import { log } from '../../utils/log'
import { sanitizeHTML } from '../../utils/dom'
import { getZIndex } from '../../core/layer'
import { trapFocus, FOCUSABLE_SELECTOR } from '../../utils/a11y'
import { getTooltipTemplate } from './tooltip.template'

// Shared Logic
let rafId = 0
let globalListenersCleanup = null

export const createTooltipStore = (Alpine) => {
  const scheduleUpdatePosition = () => {
    if (rafId) return
    if (typeof requestAnimationFrame !== 'function') return
    rafId = requestAnimationFrame(() => {
      rafId = 0
      const store = Alpine.store('uno')?.tooltip
      if (store) store.updatePosition()
    })
  }

  const ensureGlobalListeners = () => {
    if (globalListenersCleanup) return
    if (typeof window === 'undefined' || typeof document === 'undefined') return

    const onScroll = () => scheduleUpdatePosition()
    const onResize = () => scheduleUpdatePosition()
    const onKeydown = (e) => {
      if (e.key !== 'Escape') return
      const store = Alpine.store('uno')?.tooltip
      if (!store?.visible) return
      store.hide()
    }
    const onWindowClick = (e) => {
      const store = Alpine.store('uno')?.tooltip
      if (!store?.visible) return
      if (store.triggerMode !== 'click' || store.clickOutside !== true) return

      const t = e.target
      const container = document.getElementById('uno-tooltip-content')

      const insideTrigger = store.triggerEl && store.triggerEl.contains(t)
      const insideTooltip = container && container.contains(t)

      if (!insideTrigger && !insideTooltip) {
        store.hide()
      }
    }

    window.addEventListener('scroll', onScroll, true)
    window.addEventListener('resize', onResize)
    window.addEventListener('keydown', onKeydown)
    window.addEventListener('click', onWindowClick)

    globalListenersCleanup = () => {
      window.removeEventListener('scroll', onScroll, true)
      window.removeEventListener('resize', onResize)
      window.removeEventListener('keydown', onKeydown)
      window.removeEventListener('click', onWindowClick)
      globalListenersCleanup = null
    }
  }

  return {
    visible: false,
    content: '',
    placement: 'top',
    top: -9999,
    left: -9999,
    zIndex: getZIndex('TOOLTIP'),
    triggerEl: null,
    triggerMode: 'hover',
    clickOutside: true,
    arrow: true,
    arrowPosition: { top: 0, left: 0 },
    // A3 keyboard bindings (click trigger popup): trap cleanup + keydown handlers
    cleanupTrap: null,
    onPopupKeydown: null,
    onTriggerKeydown: null,

    // HTMX support
    hxUrl: null,
    hxHeaders: null,
    hxTrigger: 'load',

    /**
     * Active tooltip ID
     * @type {string|null}
     */
    active: null,

    show(instanceOpts = {}) {
      log.info('Tooltip show called', instanceOpts)
      // 4-level merge via Uno.getConfig
      const config = Uno.getConfig('tooltip', instanceOpts)
      log.info('Tooltip resolved config', config)

      // Smart Content Detection: If content is a selector and target is missing, use it as target
      if (!config.target && typeof config.content === 'string' && config.content.startsWith('#')) {
        // Check if element exists to avoid treating "#1" as a selector if it's just text
        // We only support simple ID selectors for this auto-detection
        if (!config.content.includes(' ')) {
          const targetEl = document.querySelector(config.content)
          if (targetEl) {
            config.target = config.content
          }
        }
      }

      this.triggerEl = config.element || this.triggerEl

      // Handle target (selector) for content
      if (config.target) {
        const targetEl = document.querySelector(config.target)
        if (targetEl) {
          this.content = sanitizeHTML(targetEl.innerHTML)
        } else {
          log.warn(`Tooltip target element not found: ${config.target}`)
          this.content = sanitizeHTML(config.content || '')
        }
      } else {
        // Direct content is rendered via x-html — sanitize it exactly like
        // target-sourced content to keep XSS handling consistent with Box.
        this.content = sanitizeHTML(config.content || '')
      }

      this.placement = config.placement
      this.triggerMode = config.trigger
      this.clickOutside = config.clickOutside
      this.offset = config.offset
      this.arrow = config.arrow !== false
      this.zIndex = Uno.options?.ui?.layers?.tooltip || getZIndex('TOOLTIP')

      // HTMX
      this.hxUrl = config.hxUrl
      this.hxHeaders = config.hxHeaders
      this.hxTrigger = config.hxTrigger || 'load'

      this.visible = true
      this.ensureContainer()
      ensureGlobalListeners()

      // A3：click 触发 = 弹层/披露模式。暴露触发器状态并与弹层关联、绑定键盘支持
      // （Space 合成/方向键导航/Tab 循环），渲染后焦点移入弹层，使菜单链接紧随触发器
      // 加入 Tab 序（WCAG 2.1.1 / 2.4.3 / 4.1.2）。hover/focus 触发为纯提示，不抢焦点不改状态。
      if (this.triggerMode === 'click' && this.triggerEl) {
        this.triggerEl.setAttribute('aria-haspopup', 'true')
        this.triggerEl.setAttribute('aria-expanded', 'true')
        this.triggerEl.setAttribute('aria-describedby', 'uno-tooltip-content')
        this.bindPopupKeyboard()
      }

      // Calculate position next tick to allow rendering
      Alpine.nextTick(() => {
        this.updatePosition()
        this.focusPopupIfNeeded()
      })
    },

    hide() {
      const trigger = this.triggerEl
      const wasClickMode = this.triggerMode === 'click'
      const container = document.getElementById('uno-tooltip-content')
      const active = document.activeElement
      const focusInPopup = !!container && !!active && container.contains(active)

      this.visible = false
      this.triggerEl = null

      // A3：解除键盘绑定并收起触发器状态
      if (this.cleanupTrap) {
        this.cleanupTrap()
        this.cleanupTrap = null
      }
      if (container && this.onPopupKeydown) {
        container.removeEventListener('keydown', this.onPopupKeydown)
        this.onPopupKeydown = null
      }
      if (trigger && this.onTriggerKeydown) {
        trigger.removeEventListener('keydown', this.onTriggerKeydown)
        this.onTriggerKeydown = null
      }
      if (wasClickMode && trigger && trigger.setAttribute) {
        trigger.setAttribute('aria-expanded', 'false')
        trigger.removeAttribute('aria-describedby')
      }

      if (globalListenersCleanup) {
        globalListenersCleanup()
        globalListenersCleanup = null
      }

      if (rafId) {
        cancelAnimationFrame(rafId)
        rafId = 0
      }

      // A3：焦点曾位于关闭中的弹层（Esc/动作/Space）或无处停留（outside click 落在
      // 不可聚焦区域）时还原至触发器，键盘用户不会坠回 body
      if (
        wasClickMode &&
        trigger &&
        typeof trigger.focus === 'function' &&
        (focusInPopup || active === document.body)
      ) {
        trigger.focus()
      }
    },

    /**
     * A3 click 弹层键盘绑定：Tab 循环（trapFocus）、role=menu 内容方向键导航、
     * 锚点触发器 Space 合成（Enter 在锚点/按钮上由原生 click 覆盖，按钮 Space 亦原生，
     * 仅锚点需合成以免双重 toggle）。
     */
    bindPopupKeyboard() {
      const container = document.getElementById('uno-tooltip-content')
      if (!container) return

      this.cleanupTrap = trapFocus(container, { initialFocus: false })

      this.onPopupKeydown = (e) => {
        if (!['ArrowDown', 'ArrowUp', 'Home', 'End'].includes(e.key)) return
        if (!container.querySelector('[role="menu"]')) return
        const items = Array.from(container.querySelectorAll(FOCUSABLE_SELECTOR))
        if (items.length === 0) return
        e.preventDefault()
        const idx = items.indexOf(document.activeElement)
        let next = 0
        if (e.key === 'ArrowDown') next = (idx + 1) % items.length
        else if (e.key === 'ArrowUp') next = (idx - 1 + items.length) % items.length
        else if (e.key === 'End') next = items.length - 1
        if (items[next]) items[next].focus()
      }
      container.addEventListener('keydown', this.onPopupKeydown)

      if (this.triggerEl && this.triggerEl.tagName === 'A') {
        this.onTriggerKeydown = (e) => {
          if (e.key !== ' ' && e.key !== 'Spacebar') return
          e.preventDefault()
          if (this.visible && this.triggerEl === e.currentTarget) this.hide()
        }
        this.triggerEl.addEventListener('keydown', this.onTriggerKeydown)
      }
    },

    /**
     * A3 焦点移入弹层首个可聚焦元素，键盘用户落在刚打开的披露内容内（2.4.3）；
     * 弹层无可聚焦内容时为 no-op（焦点留在触发器）。
     */
    focusPopupIfNeeded() {
      if (this.triggerMode !== 'click') return
      const container = document.getElementById('uno-tooltip-content')
      if (!container || container.contains(document.activeElement)) return
      const first = container.querySelector(FOCUSABLE_SELECTOR)
      if (first) first.focus()
    },

    cancelHide() {
      if (this.hideTimeout) {
        clearTimeout(this.hideTimeout)
        this.hideTimeout = null
      }
    },

    scheduleHide(delay = 300) {
      this.cancelHide()
      this.hideTimeout = setTimeout(() => {
        this.hide()
      }, delay)
    },

    updatePosition() {
      if (!this.visible || !this.triggerEl) return

      const tooltipEl = document.getElementById('uno-tooltip-content')
      if (!tooltipEl) return

      const pos = calculatePosition(this.triggerEl, tooltipEl, this.placement, {
        offset: this.offset
      })

      this.top = pos.top
      this.left = pos.left
      // calculatePosition() returns `{ top, left, position }` — keep the
      // previous placement as fallback so a missing field never wipes state.
      this.placement = pos.position || this.placement // Update placement if flipped

      // calculatePosition() always returns an arrow anchor relative to the
      // tooltip box; visibility is controlled by `x-show="...tooltip.arrow"`.
      if (pos.arrow) {
        this.arrowPosition = pos.arrow
      }
    },

    ensureContainer() {
      const id = 'uno-tooltip-manager'
      if (document.getElementById(id)) return
      const zBase = Uno.options?.ui?.layers?.tooltip || getZIndex('TOOLTIP')

      const template = getTooltipTemplate()

      DomManager.ensureContainer(id, template, {
        'x-data': '{}',
        style: `position: absolute; top: 0; left: 0; width: 0; height: 0; z-index: ${zBase};`
      })

      // Init tree if needed
      const el = document.getElementById(id)
      if (el && Alpine && !el._x_dataStack) {
        Alpine.initTree(el)
      }
    }
  }
}
