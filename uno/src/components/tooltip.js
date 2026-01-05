/**
 * Uno - Tooltip Component (Store-First, Directive-Driven Architecture)
 *
 * A lightweight, declarative, performance-optimized tooltip component.
 * Follows the "Store-First, Directive-Driven" pattern.
 *
 * Features:
 * - RAF-based positioning with auto-flip
 * - Global singleton container (performance)
 * - Declarative directive interface (x-uno-tooltip)
 * - HTMX native integration
 * - 4-level configuration merge
 * - Styles decoupled via Uno.getStyle
 *
 * @module uno.components.tooltip
 * @version 0.3.0
 * @license MIT
 */

import { calculatePosition } from '../utils/position'
import { DomManager as domManager } from '../utils/dom'
import { Uno } from '../core/uno'
import { mergeDeep } from '../utils/object'
import { registerPlugin } from '../core/plugin-registry'
import { log } from '../utils/log'
import { sanitizeHTML } from '../utils/dom'
import { getZIndex } from '../core/layer'

const tooltipPlugin = (Alpine, options = {}) => {
  // Register plugin defaults
  const pluginDefaults = {
    placement: 'top',
    trigger: 'hover',
    delay: 300,
    offset: 8,
    arrow: true,
    clickOutside: true
  }
  Uno.registerPluginDefaults('tooltip', {
    ...pluginDefaults,
    ...options.plugins?.tooltip,
    ...options.tooltip
  })

  // ============================================================
  // Shared Logic
  // ============================================================
  let rafId = 0
  let globalListenersCleanup = null

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

  // ============================================================
  // Store (Alpine.store('uno').tooltip)
  // ============================================================
  const tooltipStore = {
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

      const template = `
        <div id="uno-tooltip-content"
             x-show="$store.uno.tooltip.visible"
             role="tooltip"
             x-on:mouseenter="$store.uno.tooltip.cancelHide()"
             x-on:mouseleave="$store.uno.tooltip.scheduleHide()"
             x-transition:enter="transition ease-out duration-200"
             x-transition:enter-start="opacity-0 scale-95"
             x-transition:enter-end="opacity-100 scale-100"
             x-transition:leave="transition ease-in duration-150"
             x-transition:leave-start="opacity-100 scale-100"
             x-transition:leave-end="opacity-0 scale-95"
             :class="Uno.getStyle('tooltip', 'content')"
             :style="{ top: $store.uno.tooltip.top + 'px', left: $store.uno.tooltip.left + 'px', zIndex: $store.uno.tooltip.zIndex }">
             
             <!-- Arrow -->
             <div x-show="$store.uno.tooltip.arrow"
                  :style="{ top: $store.uno.tooltip.arrowPosition.top + 'px', left: $store.uno.tooltip.arrowPosition.left + 'px' }">
             </div>

             <!-- Content -->
             <div id="uno-tooltip-body" :class="Uno.getStyle('tooltip', 'body')">
                <template x-if="!$store.uno.tooltip.hxUrl">
                  <div x-html="$store.uno.tooltip.content"></div>
                </template>
                <template x-if="$store.uno.tooltip.hxUrl">
                   <div x-init="htmx.process($el)"
                        :hx-get="$store.uno.tooltip.hxUrl"
                        :hx-trigger="$store.uno.tooltip.hxTrigger"
                        hx-swap="innerHTML"
                        :hx-headers="$store.uno.tooltip.hxHeaders ? JSON.stringify($store.uno.tooltip.hxHeaders) : null">
                        <!-- Loading State -->
                        <div class="uno-tooltip-loading">Loading...</div>
                   </div>
                </template>
             </div>
        </div>
      `

      domManager.ensureContainer(id, template, {
        'x-data': '{}',
        style: `position: absolute; top: 0; left: 0; width: 0; height: 0; z-index: ${zBase};`
      })

      // Init tree if needed (ensureContainer usually handles it if domManager supports it,
      // but let's be safe as per box.js pattern)
      const el = document.getElementById(id)
      if (el && Alpine && !el._x_dataStack) {
        Alpine.initTree(el)
      }
    }
  }

  // Register Store using the registry
  registerPlugin(Alpine, 'tooltip', tooltipStore)

  // ============================================================
  // Directive Handler
  // ============================================================
  const directiveHandler = (el, { expression, modifiers }, { evaluateLater, effect, cleanup }) => {
    const getExpression = expression ? evaluateLater(expression) : null

    // Helper to resolve config from attributes and expression
    const resolveConfig = (val) => {
      // Strict Mode: Value must be an object or string (content)
      // Modifiers are NOT supported for configuration to enforce explicit config objects.
      if (modifiers.length > 0) {
        console.warn(
          '[Uno] x-uno-tooltip modifiers are deprecated in v2.0. Use configuration object instead.'
        )
      }

      const exprConfig = typeof val === 'object' && val !== null ? val : val ? { content: val } : {}
      return { ...exprConfig, element: el }
    }

    const show = (val) => {
      const store = Alpine.store('uno')?.tooltip
      if (!store) return
      store.show(resolveConfig(val))
    }

    const toggle = (val) => {
      const store = Alpine.store('uno')?.tooltip
      if (!store) return
      if (store.visible && store.triggerEl === el) {
        store.hide()
      } else {
        show(val)
      }
    }

    let currentVal = null

    // Reactively update config
    if (getExpression) {
      effect(() => {
        getExpression((val) => {
          currentVal = val
        })
      })
    }

    // Event Handlers
    const onMouseEnter = () => {
      // Safety check: if expression exists but hasn't evaluated, skip
      // This prevents default 'hover' trigger from firing before config is loaded
      if (expression && currentVal === null) return

      const config = resolveConfig(currentVal)
      const finalConfig = Uno.getConfig('tooltip', config)
      if (finalConfig.trigger === 'hover') {
        show(currentVal)
      }
    }

    const onMouseLeave = () => {
      const store = Alpine.store('uno')?.tooltip
      if (!store) return
      if (store.triggerMode === 'hover' && store.triggerEl === el) {
        store.scheduleHide()
      }
    }

    const onClick = (e) => {
      const config = resolveConfig(currentVal)
      const finalConfig = Uno.getConfig('tooltip', config)
      if (finalConfig.trigger === 'click') {
        e.preventDefault()
        toggle(currentVal)
      }
    }

    const onFocus = () => {
      const config = resolveConfig(currentVal)
      const finalConfig = Uno.getConfig('tooltip', config)
      if (finalConfig.trigger === 'focus') {
        show(currentVal)
      }
    }

    const onBlur = () => {
      const store = Alpine.store('uno')?.tooltip
      if (store && store.triggerMode === 'focus' && store.triggerEl === el) {
        store.hide()
      }
    }

    el.addEventListener('mouseenter', onMouseEnter)
    el.addEventListener('mouseleave', onMouseLeave)
    el.addEventListener('click', onClick)
    el.addEventListener('focus', onFocus)
    el.addEventListener('blur', onBlur)

    cleanup(() => {
      el.removeEventListener('mouseenter', onMouseEnter)
      el.removeEventListener('mouseleave', onMouseLeave)
      el.removeEventListener('click', onClick)
      el.removeEventListener('focus', onFocus)
      el.removeEventListener('blur', onBlur)
      const store = Alpine.store('uno')?.tooltip
      if (store && store.triggerEl === el) {
        store.hide()
      }
    })
  }

  Alpine.directive('uno-tooltip', directiveHandler)
}

export default tooltipPlugin
