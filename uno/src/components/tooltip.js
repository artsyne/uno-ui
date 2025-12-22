/**
 * Uno - Tooltip Directive
 *
 * @module uno.components.tooltip
 * @version 0.2.1
 * @license MIT
 */

import { calculatePosition } from '../utils/position'
import { log } from '../utils/log'
import { setHTML, DomManager as domManager } from '../utils/dom'
import { Uno } from '../core/uno'

export default (Alpine, options = {}) => {
  const initConfig = options.plugins?.tooltip || options.tooltip || {}
  const getConfig = () => {
    const opts = typeof window !== 'undefined' ? window.Uno?.options : Uno.options
    return opts?.plugins?.tooltip || opts?.tooltip || initConfig
  }

  if (!Alpine) {
    log.error('Tooltip: Alpine is undefined. Plugin initialization failed.')
    return
  }

  let containerEl = null
  let tooltipEl = null
  let bodyEl = null
  let rafId = 0
  let globalListenersCleanup = null

  const getTooltipZIndex = () => {
    let zBase = Number(Uno.options?.ui?.layers?.tooltipBase ?? 20000)
    if (!Number.isFinite(zBase)) zBase = 20000
    return zBase
  }

  const ensureEls = () => {
    if (
      containerEl &&
      tooltipEl &&
      bodyEl &&
      typeof document !== 'undefined' &&
      document.contains(containerEl)
    ) {
      return true
    }
    containerEl =
      typeof document !== 'undefined' ? document.getElementById('uno-tooltip-container') : null
    tooltipEl =
      typeof document !== 'undefined' ? document.getElementById('uno-tooltip-content') : null
    bodyEl = typeof document !== 'undefined' ? document.getElementById('uno-tooltip-body') : null
    return !!(containerEl && tooltipEl && bodyEl)
  }

  const scheduleUpdatePosition = () => {
    if (rafId) return
    if (typeof requestAnimationFrame !== 'function') return
    rafId = requestAnimationFrame(() => {
      rafId = 0
      Alpine.store('tooltip')?.updatePosition?.()
    })
  }

  const ensureGlobalListeners = () => {
    if (globalListenersCleanup) return
    if (typeof window === 'undefined' || typeof document === 'undefined') return

    const onScroll = () => scheduleUpdatePosition()
    const onResize = () => scheduleUpdatePosition()
    const onKeydown = (e) => {
      if (e.key !== 'Escape') return
      const store = Alpine.store('tooltip')
      if (!store?.visible) return
      store.hide()
    }
    const onWindowClick = (e) => {
      const store = Alpine.store('tooltip')
      if (!store?.visible) return
      if (store.triggerMode !== 'click' || store.clickOutside !== true) return

      if (!ensureEls()) return

      const t = e.target
      const insideTrigger =
        store.eventEl && typeof store.eventEl.contains === 'function' && store.eventEl.contains(t)
      const insideTooltip = tooltipEl && tooltipEl.contains(t)
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

  // 1. Global Store for Singleton Tooltip
  Alpine.store('tooltip', {
    visible: false,
    content: '',
    placement: 'top',
    top: -9999, // Start off-screen
    left: -9999,
    zIndex: getTooltipZIndex(),
    triggerEl: null,
    eventEl: null,
    triggerMode: 'hover',
    clickOutside: true,
    onShow: null,

    async show(el, content, config = {}) {
      this.triggerEl = el
      this.eventEl = config.eventEl || el
      this.content = content
      this.placement = config.placement || getConfig().placement || 'top'
      this.onShow = config.onShow
      this.triggerMode = config.trigger || 'hover'
      this.clickOutside =
        config.clickOutside !== undefined ? config.clickOutside : this.triggerMode === 'click'
      this.zIndex = getTooltipZIndex()

      this.visible = true

      // Ensure container exists (defensive)
      this.ensureContainer()
      ensureGlobalListeners()

      // Inject content and init Alpine for interactivity
      if (ensureEls()) {
        setHTML(bodyEl, this.content, { trusted: config.trustedHTML === true, alpine: Alpine })
      }

      // Wait for DOM update to measure tooltip size
      await Alpine.nextTick()

      this.updatePosition()

      if (this.onShow) {
        this.onShow(bodyEl)
      }
    },

    hide() {
      this.visible = false
      this.triggerEl = null
      this.eventEl = null
      if (globalListenersCleanup) globalListenersCleanup()
      if (rafId) {
        cancelAnimationFrame(rafId)
        rafId = 0
      }
    },

    updatePosition() {
      if (!this.triggerEl) return
      if (!ensureEls()) return

      // Calculate
      const pos = calculatePosition(this.triggerEl, tooltipEl, this.placement)
      this.top = pos.top
      this.left = pos.left
    },

    ensureContainer() {
      injectContainer()
      ensureEls()
    },

    destroy() {
      this.hide()
    }
  })

  // 2. Inject Global Tooltip Container
  const injectContainer = () => {
    if (typeof document === 'undefined') return
    const rounded = Uno.getStyle('base', 'rounded', 'rounded-lg')
    // Default content style
    const contentDefault = `fixed px-3 py-2 text-sm font-medium ${rounded} shadow-sm whitespace-nowrap pointer-events-none border bg-white text-gray-900 border-gray-200 dark:bg-slate-800 dark:text-white dark:border-slate-700`
    const contentClass = Uno.getStyle('tooltip', 'content', contentDefault)
    const bodyClass = Uno.getStyle('tooltip', 'body', 'pointer-events-auto')

    const template = `
          <div id="uno-tooltip-content"
               x-show="$store.tooltip.visible"
               role="tooltip"
               x-transition:enter="transition ease-out duration-200"
               x-transition:enter-start="opacity-0 scale-95"
               x-transition:enter-end="opacity-100 scale-100"
               x-transition:leave="transition ease-in duration-150"
               x-transition:leave-start="opacity-100 scale-100"
               x-transition:leave-end="opacity-0 scale-95"
               class="${contentClass}"
              :style="\`top: \${$store.tooltip.top}px; left: \${$store.tooltip.left}px; z-index: \${$store.tooltip.zIndex};\`"
              style="display: none;">
              <div id="uno-tooltip-body" class="${bodyClass}"></div>
          </div>
        `

    domManager.ensureContainer('uno-tooltip-container', template, {
      'x-data': '{}'
    })
  }

  // Eager injection removed for P1-10 (On-demand)
  // Container is injected when show() is called.

  // 3. Directive: x-tooltip
  Alpine.directive(
    'tooltip',
    (el, { expression, modifiers }, { evaluateLater, effect, cleanup }) => {
      const getExpression = evaluateLater(expression)

      let config = {
        content: '',
        placement: null,
        trigger: null,
        clickOutside: undefined
      }

      const resolvePlacement = () => {
        const configData = getConfig()
        let placement = configData.placement || 'top'
        if (config.placement) placement = config.placement
        if (modifiers.includes('top')) placement = 'top'
        if (modifiers.includes('bottom')) placement = 'bottom'
        if (modifiers.includes('left')) placement = 'left'
        if (modifiers.includes('right')) placement = 'right'
        return placement
      }

      const resolveTrigger = () => {
        const configData = getConfig()
        if (modifiers.includes('click')) return 'click'
        if (modifiers.includes('focus')) return 'focus'
        if (modifiers.includes('hover')) return 'hover'
        if (config.trigger) return config.trigger
        return configData.trigger || 'hover'
      }

      const resolveEventTarget = () => {
        const isDisabled =
          (typeof el.disabled === 'boolean' && el.disabled === true) ||
          el.getAttribute('aria-disabled') === 'true' ||
          el.getAttribute('disabled') !== null
        if (!isDisabled) return el
        return el.parentElement || el
      }

      // Helper to process value
      const processValue = (val) => {
        if (typeof val === 'string') {
          config.content = val
        } else if (typeof val === 'object') {
          config = { ...config, ...val }
        }

        // Check for ID selector or target
        if (config.content && config.content.startsWith('#') && !config.content.includes(' ')) {
          const target = document.querySelector(config.content)
          if (target) config.content = target.innerHTML
        } else if (config.target) {
          const target = document.querySelector(config.target)
          if (target) config.content = target.innerHTML
        }
      }

      // Watch for changes
      effect(() => {
        getExpression((val) => {
          processValue(val)
        })
      })

      // Event handlers
      const show = () => {
        const store = Alpine.store('tooltip')
        const mode = resolveTrigger()
        el.setAttribute('aria-describedby', 'uno-tooltip-content')
        store.show(el, config.content, {
          ...config,
          placement: resolvePlacement(),
          trigger: mode,
          eventEl: resolveEventTarget()
        })
      }
      const hide = () => {
        el.removeAttribute('aria-describedby')
        Alpine.store('tooltip').hide()
      }
      const toggle = () => {
        const store = Alpine.store('tooltip')
        // If visible and triggered by this element, hide it
        if (store.visible && store.triggerEl === el) {
          store.hide()
        } else {
          show()
        }
      }

      const onClick = () => {
        if (resolveTrigger() !== 'click') return
        toggle()
      }

      const onMouseEnter = () => {
        if (resolveTrigger() !== 'hover') return
        show()
      }

      const onMouseLeave = () => {
        if (resolveTrigger() !== 'hover') return
        hide()
      }

      const onFocus = () => {
        if (resolveTrigger() !== 'focus') return
        show()
      }

      const onBlur = () => {
        if (resolveTrigger() !== 'focus') return
        hide()
      }

      // Bind events
      const eventTarget = resolveEventTarget()
      eventTarget.addEventListener('click', onClick)
      eventTarget.addEventListener('mouseenter', onMouseEnter)
      eventTarget.addEventListener('mouseleave', onMouseLeave)
      eventTarget.addEventListener('focus', onFocus)
      eventTarget.addEventListener('blur', onBlur)

      // Cleanup
      cleanup(() => {
        eventTarget.removeEventListener('click', onClick)
        eventTarget.removeEventListener('mouseenter', onMouseEnter)
        eventTarget.removeEventListener('mouseleave', onMouseLeave)
        eventTarget.removeEventListener('focus', onFocus)
        eventTarget.removeEventListener('blur', onBlur)
        hide()
      })
    }
  )

  // Expose Public API
  const api = {
    show: (el, content, config) => Alpine.store('tooltip').show(el, content, config),
    hide: () => Alpine.store('tooltip').hide(),
    updatePosition: () => Alpine.store('tooltip').updatePosition(),
    destroy: () => Alpine.store('tooltip').destroy()
  }
  Uno.tooltip = api
  if (typeof window !== 'undefined' && window.Uno) {
    window.Uno.tooltip = api
  }
}
