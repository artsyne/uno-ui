import { Uno } from '../../core/uno'
import { calculatePosition } from '../../utils/position'
import { DomManager } from '../../utils/dom'
import { log } from '../../utils/log'
import { sanitizeHTML } from '../../utils/dom'
import { getZIndex } from '../../core/layer'
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

    // HTMX support
    hxUrl: null,
    hxHeaders: null,
    hxTrigger: 'load',

    init() {
      this.ensureContainer()
    },

    show(instanceOpts = {}) {
      log.info('Tooltip show called', instanceOpts)
      // 4-level merge via Uno.getConfig
      const config = Uno.getConfig('tooltip', instanceOpts)
      log.info('Tooltip resolved config', config)

      this.triggerEl = config.element || this.triggerEl

      // Handle target (selector) for content
      if (config.target) {
        const targetEl = document.querySelector(config.target)
        if (targetEl) {
          this.content = sanitizeHTML(targetEl.innerHTML)
        } else {
          log.warn(`Tooltip target element not found: ${config.target}`)
          this.content = config.content || ''
        }
      } else {
        this.content = config.content || ''
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

      // Calculate position next tick to allow rendering
      Alpine.nextTick(() => {
        this.updatePosition()
      })
    },

    hide() {
      this.visible = false
      this.triggerEl = null

      if (globalListenersCleanup) {
        globalListenersCleanup()
        globalListenersCleanup = null
      }

      if (rafId) {
        cancelAnimationFrame(rafId)
        rafId = 0
      }
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
        offset: this.offset,
        arrow: this.arrow
      })

      this.top = pos.top
      this.left = pos.left
      this.placement = pos.placement // Update placement if flipped

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
