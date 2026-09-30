import { Uno } from '../../core/uno'
import { log } from '../../utils/log'

export const createTabStore = (Alpine) => ({
  // State: { [group]: activeTabName }
  groups: {},
  // Registry: { [tabName]: { group, name, ... } }
  registry: {},
  // P1-B4: { [group]: [tabName...] } 组内注册序（DOM 序代理），供方向键旋转
  order: {},
  // Initialization state
  initialized: false,

  /**
   * Ensure the store is initialized (lazy loading)
   * Sets up hash change listeners only when needed
   */
  ensureInit() {
    if (this.initialized) return
    this.initialized = true

    // Hash synchronization
    const checkHash = () => {
      let hash = window.location.hash.slice(1)

      // Fallback for some environments where location.hash might be empty initially
      if (!hash && window.location.href.includes('#')) {
        hash = window.location.href.split('#')[1]
      }

      log.info(`[Tab] checkHash: ${hash}`)

      if (hash && this.registry[hash]) {
        const { group, name } = this.registry[hash]
        // If hash is found in registry, activate it
        log.info(`[Tab] checkHash match: ${name}`)
        this.activate({ group, name, updateHash: false })
      }
    }

    // Initial check
    checkHash()

    // Listen for hash changes
    window.addEventListener('hashchange', checkHash)
    // Safety check on load (in case hash was not ready during init/register)
    window.addEventListener('load', checkHash)
  },

  /**
   * Active tab ID
   * @type {string|null}
   */
  active: null,

  /**
   * Activate a tab
   * @param {object} instanceOpts - Options { name, group, ... }
   */
  activate(instanceOpts = {}) {
    this.ensureInit()

    // 4-level merge via Uno.getConfig
    const config = Uno.getConfig('tab', {
      group: 'default',
      hash: true,
      ...instanceOpts
    })

    if (!config.name) return

    log.info(
      `[Tab] Activating ${config.name} in ${config.group}. Hash update: ${config.updateHash !== false}`
    )

    // Register if not known (lazy registration)
    if (!this.registry[config.name]) {
      this.registry[config.name] = { group: config.group, name: config.name }
    }

    // Update state
    this.groups[config.group] = config.name

    // P1-B4 共享面板模式（htmx swap）：面板容器声明 role=tabpanel + data-tab-group，
    // aria-labelledby 随激活 tab 更新；每 tab 独立面板用静态 labelledby，不受影响
    if (typeof document !== 'undefined') {
      document
        .querySelectorAll('[role="tabpanel"][data-tab-group="' + config.group + '"]')
        .forEach((p) => p.setAttribute('aria-labelledby', 'tab-' + config.name))
    }

    // Handle Hash
    if (config.hash && config.updateHash !== false) {
      if (window.location.hash.slice(1) !== config.name) {
        window.history.replaceState(null, null, `#${config.name}`)
      }
    }

    // Trigger HTMX if needed (Simulate click if hx-get exists)
    // This is crucial for deep linking where activation happens via JS, not user click.
    if (config.updateHash === false) {
      // Typically means activated from hash or init
      // Use setTimeout to allow DOM to be ready if called during init
      setTimeout(() => {
        const tabEl = document.getElementById(`tab-${config.name}`)
        if (tabEl && tabEl.hasAttribute('hx-get')) {
          log.info(`[Tab] Triggering HTMX for ${config.name}`)
          if (window.htmx) {
            window.htmx.trigger(tabEl, 'click')
          }
        }
      }, 0)
    }
  },

  /**
   * Check if a tab is active
   * @param {string} name - Tab name
   * @param {string} group - Group name
   * @returns {boolean}
   */
  isActive(name, group = 'default') {
    const active = this.groups[group] === name
    if (active) {
      // log.info(`[Tab] isActive true for ${name} in ${group}`)
    }
    return active
  },

  /**
   * Register a tab (helper for directives or manual setup)
   * @param {string} name
   * @param {string} group
   */
  register(name, group = 'default') {
    this.ensureInit()

    log.info(`[Tab] Registering ${name} in ${group}. Hash: ${window.location.hash}`)
    this.registry[name] = { name, group }
    // P1-B4：记录组内注册序（去重），供方向键旋转
    if (!this.order[group]) this.order[group] = []
    if (!this.order[group].includes(name)) this.order[group].push(name)
    // If no active tab for this group, set this one
    if (!this.groups[group]) {
      this.groups[group] = name
    }

    // Check if this new tab matches the current hash
    // This handles deep linking when tabs are registered after store init
    let hash = window.location.hash.slice(1)
    if (!hash && window.location.href.includes('#')) {
      hash = window.location.href.split('#')[1]
    }

    if (hash === name) {
      log.info(`[Tab] Hash match for ${name}. Activating.`)
      this.activate({ group, name, updateHash: false })
    }
  },

  /**
   * P1-B4 APG 键盘旋转：计算目标 tab 名（不激活），调用方负责 activate + focus
   * （roving tabindex）。←→ 循环旋转，Home/End 跳首尾。
   * @param {string} group
   * @param {string} current
   * @param {string} key - ArrowRight/ArrowLeft/Home/End
   * @returns {string|null}
   */
  rotate(group, current, key) {
    const list = this.order[group]
    if (!list || list.length === 0) return null
    const idx = list.indexOf(current)
    if (key === 'Home') return list[0]
    if (key === 'End') return list[list.length - 1]
    if (idx === -1) return list[0]
    if (key === 'ArrowRight') return list[(idx + 1) % list.length]
    if (key === 'ArrowLeft') return list[(idx - 1 + list.length) % list.length]
    return null
  }
})
