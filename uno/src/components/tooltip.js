/**
 * Uno - Tooltip Directive
 *
 * @module uno.components.tooltip
 * @version 0.2.0
 * @license MIT
 */

import { calculatePosition } from '../utils/position'
import { log } from '../utils/log'
import { setHTML } from '../utils/dom'

export default (Alpine) => {
  if (!Alpine) {
    log.error('Tooltip: Alpine is undefined. Plugin initialization failed.')
    return
  }

  // 1. Global Store for Singleton Tooltip
  Alpine.store('tooltip', {
    visible: false,
    content: '',
    placement: 'top',
    top: -9999, // Start off-screen
    left: -9999,
    triggerEl: null,
    onShow: null,

    async show(el, content, config = {}) {
      this.triggerEl = el
      this.content = content
      this.placement = config.placement || 'top'
      this.onShow = config.onShow

      this.visible = true

      // Ensure container exists (defensive)
      this.ensureContainer()

      // Inject content and init Alpine for interactivity
      const body = document.getElementById('uno-tooltip-body')
      if (body) {
        setHTML(body, this.content, { trusted: config.trustedHTML === true, alpine: Alpine })
      }

      // Wait for DOM update to measure tooltip size
      await Alpine.nextTick()

      this.updatePosition()

      if (this.onShow) {
        this.onShow(document.getElementById('uno-tooltip-body'))
      }
    },

    hide() {
      this.visible = false
      this.triggerEl = null
    },

    updatePosition() {
      if (!this.triggerEl) return
      const tooltipEl = document.getElementById('uno-tooltip-content')
      if (!tooltipEl) return

      // Calculate
      const pos = calculatePosition(this.triggerEl, tooltipEl, this.placement)
      this.top = pos.top
      this.left = pos.left
    },

    ensureContainer() {
      if (!document.getElementById('uno-tooltip-container')) {
        injectContainer()
      }
    }
  })

  // 2. Inject Global Tooltip Container
  const injectContainer = () => {
    if (document.getElementById('uno-tooltip-container')) return

    try {
      const container = document.createElement('div')
      container.id = 'uno-tooltip-container'
      // Initialize as an Alpine component to ensure reactivity context if needed,
      // though we rely on $store.
      container.setAttribute('x-data', '{}')

      container.innerHTML = `
          <div id="uno-tooltip-content"
               x-show="$store.tooltip.visible"
               x-transition:enter="transition ease-out duration-200"
               x-transition:enter-start="opacity-0 scale-95"
               x-transition:enter-end="opacity-100 scale-100"
               x-transition:leave="transition ease-in duration-150"
               x-transition:leave-start="opacity-100 scale-100"
               x-transition:leave-end="opacity-0 scale-95"
               class="fixed px-3 py-2 text-sm font-medium rounded-lg shadow-sm whitespace-nowrap pointer-events-none border bg-white text-gray-900 border-gray-200 dark:bg-slate-800 dark:text-white dark:border-slate-700 z-[20000]"
              :style="\`top: \${$store.tooltip.top}px; left: \${$store.tooltip.left}px;\`"
              style="display: none;">
              <div id="uno-tooltip-body" class="pointer-events-auto"></div>
          </div>
        `

      // Prefer appending to body
      if (document.body) {
        document.body.appendChild(container)
      } else {
        document.documentElement.appendChild(container)
      }
    } catch (e) {
      log.error('Tooltip: Failed to inject container', e)
    }
  }

  // Eager injection to ensure container is ready
  if (document.body) {
    injectContainer()
  } else {
    document.addEventListener('DOMContentLoaded', injectContainer)
  }

  // 3. Directive: x-tooltip
  Alpine.directive(
    'tooltip',
    (el, { expression, modifiers }, { evaluateLater, effect, cleanup }) => {
      const getExpression = evaluateLater(expression)
      let config = {
        content: '',
        placement: modifiers[0] || 'top',
        trigger: modifiers.includes('click') ? 'click' : 'hover'
      }

      // Handle modifiers
      if (modifiers.includes('top')) config.placement = 'top'
      if (modifiers.includes('bottom')) config.placement = 'bottom'
      if (modifiers.includes('left')) config.placement = 'left'
      if (modifiers.includes('right')) config.placement = 'right'

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
      const show = () => Alpine.store('tooltip').show(el, config.content, config)
      const hide = () => Alpine.store('tooltip').hide()
      const toggle = () => {
        const store = Alpine.store('tooltip')
        // If visible and triggered by this element, hide it
        if (store.visible && store.triggerEl === el) {
          store.hide()
        } else {
          show()
        }
      }

      if (config.trigger === 'click') {
        el.addEventListener('click', toggle)

        // Click outside to close
        const clickOutside = (e) => {
          const store = Alpine.store('tooltip')
          const tooltipContent = document.getElementById('uno-tooltip-content')
          if (
            store.visible &&
            store.triggerEl === el &&
            !el.contains(e.target) &&
            (!tooltipContent || !tooltipContent.contains(e.target))
          ) {
            store.hide()
          }
        }
        window.addEventListener('click', clickOutside)
        cleanup(() => window.removeEventListener('click', clickOutside))
      } else {
        el.addEventListener('mouseenter', show)
        el.addEventListener('mouseleave', hide)
        el.addEventListener('focus', show)
        el.addEventListener('blur', hide)
      }

      cleanup(() => {
        if (config.trigger === 'click') {
          el.removeEventListener('click', toggle)
        } else {
          el.removeEventListener('mouseenter', show)
          el.removeEventListener('mouseleave', hide)
          el.removeEventListener('focus', show)
          el.removeEventListener('blur', hide)
        }
      })
    }
  )
}
