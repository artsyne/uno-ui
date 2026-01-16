import { Uno } from '../../core/uno'
import { log } from '../../utils/log'

export const registerTooltipDirectives = (Alpine) => {
  const directiveHandler = (el, { expression, modifiers }, { evaluateLater, effect, cleanup }) => {
    const getExpression = expression ? evaluateLater(expression) : null

    // Helper to resolve config from attributes and expression
    const resolveConfig = (val) => {
      // Strict Mode: Value must be an object or string (content)
      // Modifiers are NOT supported for configuration to enforce explicit config objects.
      if (modifiers.length > 0) {
        log.warn(
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
