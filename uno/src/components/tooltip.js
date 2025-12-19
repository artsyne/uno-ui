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
  const initConfig = options.tooltip || {}
  const getConfig = () =>
    (typeof window !== 'undefined' && window.Uno?.options?.tooltip) || initConfig

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
      this.placement = config.placement || getConfig().placement || 'top'
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
      injectContainer()
    }
  })

  // 2. Inject Global Tooltip Container
  const injectContainer = () => {
    const rounded = Uno.getStyle('base', 'rounded', 'rounded-lg')
    // Default content style
    const contentDefault = `fixed px-3 py-2 text-sm font-medium ${rounded} shadow-sm whitespace-nowrap pointer-events-none border bg-white text-gray-900 border-gray-200 dark:bg-slate-800 dark:text-white dark:border-slate-700 z-[20000]`
    const contentClass = Uno.getStyle('tooltip', 'content', contentDefault)

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
              :style="\`top: \${$store.tooltip.top}px; left: \${$store.tooltip.left}px;\`"
              style="display: none;">
              <div id="uno-tooltip-body" class="pointer-events-auto"></div>
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
        trigger: null
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
        if (config.trigger) return config.trigger
        return configData.trigger || 'hover'
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
        el.setAttribute('aria-describedby', 'uno-tooltip-content')
        store.show(el, config.content, { ...config, placement: resolvePlacement() })
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

      const clickOutside = (e) => {
        if (resolveTrigger() !== 'click') return
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

      // Bind events
      el.addEventListener('click', onClick)
      el.addEventListener('mouseenter', onMouseEnter)
      el.addEventListener('mouseleave', onMouseLeave)
      el.addEventListener('focus', onMouseEnter)
      el.addEventListener('blur', onMouseLeave)
      window.addEventListener('click', clickOutside)

      // Cleanup
      cleanup(() => {
        el.removeEventListener('click', onClick)
        el.removeEventListener('mouseenter', onMouseEnter)
        el.removeEventListener('mouseleave', onMouseLeave)
        el.removeEventListener('focus', onMouseEnter)
        el.removeEventListener('blur', onMouseLeave)
        window.removeEventListener('click', clickOutside)
        hide()
      })
    }
  )

  // Expose Public API
  const api = {
    show: (el, content, config) => Alpine.store('tooltip').show(el, content, config),
    hide: () => Alpine.store('tooltip').hide(),
    updatePosition: () => Alpine.store('tooltip').updatePosition()
  }
  Uno.tooltip = api
  if (typeof window !== 'undefined' && window.Uno) {
    window.Uno.tooltip = api
  }
}
