import { log } from '../../utils/log'

export const registerTabDirectives = (Alpine) => {
  // ============================================================
  // Tab Directive (x-uno-tab)  // Usage: <button x-uno-tab="{ name: 'tab1' }">
  // B4 APG tabs：role=tab、aria-selected/roving tabindex 响应式同步、
  // aria-controls 关联面板、←→/Home/End 键盘旋转（激活后焦点跟随）。
  // ============================================================
  Alpine.directive('uno-tab', (el, { expression }, { evaluateLater, effect, cleanup }) => {
    if (!expression) {
      log.warn('[Uno] x-uno-tab directive requires a configuration object.')
      return
    }

    const getExpression = evaluateLater(expression)
    let lastConfig = null

    const parseConfig = (value) => {
      if (typeof value === 'object' && value !== null) return { ...value }
      if (typeof value === 'string') return { name: value }
      return null
    }

    // Initial Registration
    getExpression((value) => {
      const config = parseConfig(value)
      if (!config) {
        log.warn('[Uno] x-uno-tab directive value must be a string or object', value)
        return
      }

      // Default group if not specified
      if (!config.group) config.group = 'default'
      lastConfig = config

      // Register on init
      if (config.name) {
        Alpine.store('uno').tab.register(config.name, config.group)
        // B4：tab 角色与面板关联；id 约定 tab-{name}（store 深链 HTMX 触发同口径）
        el.setAttribute('role', 'tab')
        if (!el.id) el.id = 'tab-' + config.name
        el.setAttribute('aria-controls', config.panel || 'panel-' + config.name)
      }
    })

    // B4：aria-selected + roving tabindex 随激活态响应式同步
    effect(() => {
      if (!lastConfig || !lastConfig.name) return
      const active = Alpine.store('uno').tab.isActive(lastConfig.name, lastConfig.group)
      el.setAttribute('aria-selected', active ? 'true' : 'false')
      el.setAttribute('tabindex', active ? '0' : '-1')
    })

    const onClick = (e) => {
      e.preventDefault()

      getExpression((value) => {
        const config = parseConfig(value)
        if (!config) return
        if (!config.group) config.group = 'default'
        lastConfig = config

        if (config.name) {
          Alpine.store('uno').tab.activate(config)
        }
      })
    }

    // B4 APG 键盘：←→ 循环旋转、Home/End 跳首尾；焦点跟随至目标 tab 并复用
    // click 路径激活（htmx swap 型 tab 的内容装载由锚点 hx-get 承担，裸 activate 会脱节）
    const onKeyDown = (e) => {
      if (!['ArrowRight', 'ArrowLeft', 'Home', 'End'].includes(e.key)) return
      if (!lastConfig || !lastConfig.name) return
      e.preventDefault()
      const store = Alpine.store('uno').tab
      // RTL 下左右方向键的视觉含义互换：先按最近的方向容器判定，未标注时回退文档方向。
      const scope = el.closest('[dir]') || document.documentElement
      const isRTL = typeof getComputedStyle === 'function' && getComputedStyle(scope).direction === 'rtl'
      const key = isRTL
        ? e.key === 'ArrowRight'
          ? 'ArrowLeft'
          : e.key === 'ArrowLeft'
            ? 'ArrowRight'
            : e.key
        : e.key
      const target = store.rotate(lastConfig.group, lastConfig.name, key)
      if (!target) return
      const targetEl = document.getElementById('tab-' + target)
      if (!targetEl) return
      targetEl.focus()
      targetEl.click()
    }

    el.addEventListener('click', onClick)
    el.addEventListener('keydown', onKeyDown)
    cleanup(() => {
      el.removeEventListener('click', onClick)
      el.removeEventListener('keydown', onKeyDown)
    })
  })
}
