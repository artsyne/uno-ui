export const registerTabDirectives = (Alpine) => {
  // ============================================================
  // Tab Directive (x-uno-tab)  // Usage: <button x-uno-tab="{ name: 'tab1' }">
  // ============================================================
  Alpine.directive('uno-tab', (el, { expression }, { evaluateLater, cleanup }) => {
    if (!expression) {
      console.warn('[Uno] x-uno-tab directive requires a configuration object.')
      return
    }

    const getExpression = evaluateLater(expression)

    // Initial Registration
    getExpression((value) => {
      let config = {}
      if (typeof value === 'object' && value !== null) {
        config = value
      } else if (typeof value === 'string') {
        config = { name: value }
      } else {
        console.warn('[Uno] x-uno-tab directive value must be a string or object', value)
        return
      }

      // Default group if not specified
      if (!config.group) config.group = 'default'

      // Register on init
      if (config.name) {
        Alpine.store('uno').tab.register(config.name, config.group)
      }
    })

    const onClick = (e) => {
      e.preventDefault()

      getExpression((value) => {
        let config = {}
        if (typeof value === 'object' && value !== null) {
          config = value
        } else if (typeof value === 'string') {
          config = { name: value }
        }

        if (config.name) {
          Alpine.store('uno').tab.activate(config)
        }
      })
    }

    el.addEventListener('click', onClick)
    cleanup(() => el.removeEventListener('click', onClick))
  })
}
