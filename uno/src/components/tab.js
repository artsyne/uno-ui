import { Uno } from '../core/uno'
import { once, setHTML } from '../utils/dom'

let storedAlpine

const factory = () => ({
  activeTab: null,
  _cleanup: [],

  init() {
    if (!once(this.$el, 'tabs_init')) return

    this.$nextTick(() => {
      const hash = window.location.hash.slice(1)
      let found = false

      if (hash) {
        const tabWithHash = this.$el.querySelector(`[data-tab-hash="${hash}"]`)
        if (tabWithHash) {
          this.activeTab = tabWithHash.getAttribute('data-tab-target')
          found = true
        }
      }

      if (!found && !this.activeTab) {
        const firstTab = this.$el.querySelector('[role="tab"]')
        if (firstTab) {
          this.activeTab = firstTab.getAttribute('data-tab-target')
        }
      }
    })

    const onHashChange = () => {
      const hash = window.location.hash.slice(1)
      if (hash) {
        const tabWithHash = this.$el.querySelector(`[data-tab-hash="${hash}"]`)
        if (tabWithHash) {
          this.activeTab = tabWithHash.getAttribute('data-tab-target')
        }
      }
    }

    window.addEventListener('hashchange', onHashChange)
    this._cleanup.push(() => window.removeEventListener('hashchange', onHashChange))
  },

  destroy() {
    this._cleanup.forEach((fn) => fn())
    this._cleanup = []
  },

  isActive(tab) {
    return this.activeTab === tab
  },

  async open(tab, ajaxUrl = null, hash = null, hxUrl = null) {
    this.activeTab = tab

    if (hash) {
      if (window.location.hash.slice(1) !== hash) {
        window.history.replaceState(null, null, `#${hash}`)
      }
    }

    if (ajaxUrl) {
      const panel = document.getElementById(tab)
      if (panel && panel.getAttribute('data-tab-loaded') !== 'true') {
        try {
          panel.innerHTML =
            Uno.options.ui?.box?.loading ||
            '<div class="animate-pulse bg-gray-200 h-32 w-full rounded"></div>'

          const response = await fetch(ajaxUrl)
          if (!response.ok) throw new Error('Network response was not ok')
          const content = await response.text()

          setHTML(panel, content, { trusted: false, alpine: storedAlpine || window.Alpine })
          panel.setAttribute('data-tab-loaded', 'true')
        } catch (err) {
          panel.innerHTML = `<div class="text-red-500 p-4">Error loading content: ${err?.message || 'Unknown error'}</div>`
        }
      }
    }

    if (hxUrl) {
      const panel = document.getElementById(tab)
      if (panel && panel.getAttribute('data-tab-loaded') !== 'true') {
        try {
          panel.innerHTML =
            Uno.options.ui?.box?.loading ||
            '<div class="animate-pulse bg-gray-200 h-32 w-full rounded"></div>'

          const response = await fetch(hxUrl, { headers: { 'HX-Request': 'true' } })
          if (!response.ok) throw new Error('Network response was not ok')
          const content = await response.text()

          setHTML(panel, content, { trusted: false, alpine: storedAlpine || window.Alpine })
          panel.setAttribute('data-tab-loaded', 'true')
        } catch {
          panel.innerHTML = '<div class="text-red-500 p-4">Error loading content via HTMX</div>'
        }
      }
    }
  },

  trigger(tab, ajaxUrl = null, hash = null, hxUrl = null) {
    if (hash && typeof window !== 'undefined' && window.location.hash.slice(1) === hash) {
      if (!this.activeTab) this.activeTab = tab
    }

    const baseActiveClass = Uno.getStyle('tab', 'active', 'border-b-2 border-blue-500')
    const baseInactiveClass = Uno.getStyle('tab', 'inactive', 'text-slate-500 hover:text-slate-700')
    const activeClass = Uno.getStyle('tab', 'triggerActive', baseActiveClass)
    const inactiveClass = Uno.getStyle('tab', 'triggerInactive', baseInactiveClass)

    return {
      ['@click.prevent']: `open('${tab}', ${ajaxUrl ? `'${ajaxUrl}'` : 'null'}, ${hash ? `'${hash}'` : 'null'}, ${hxUrl ? `'${hxUrl}'` : 'null'})`,
      [':class']: `isActive('${tab}') ? '${activeClass}' : '${inactiveClass}'`,
      ['role']: 'tab',
      [':aria-selected']: `isActive('${tab}') ? 'true' : 'false'`,
      ['data-tab-target']: tab,
      ['data-tab-hash']: hash || null,
      ['data-tab-url']: ajaxUrl || null,
      ['data-tab-hx-url']: hxUrl || null
    }
  },

  tablist: {
    ['role']: 'tablist',
    ['@keydown.right.prevent']: 'next()',
    ['@keydown.left.prevent']: 'prev()',
    ['@keydown.home.prevent']: 'first()',
    ['@keydown.end.prevent']: 'last()'
  },

  panel(tab) {
    return {
      ['x-show']: `isActive('${tab}')`,
      ['role']: 'tabpanel',
      ['id']: tab,
      ['tabindex']: '0'
    }
  }
})

export { factory }

if (typeof window !== 'undefined') {
  window.unoTabs = factory
}

export default (Alpine) => {
  storedAlpine = Alpine
  Alpine.data('unoTabs', factory)

  Uno.tab = {
    open: (id) => {
      const trigger = document.querySelector(`[data-tab-target="${id}"]`)
      if (trigger) {
        trigger.click()
      } else {
        window.location.hash = id
      }
    }
  }
}
