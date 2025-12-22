import { Uno } from '../core/uno'
import { escapeHTML, once, setHTML } from '../utils/dom'
import { requestData } from '../utils/request'

let storedAlpine

const factory = (initialTab = null, initOptions = {}) => ({
  activeTab: initialTab,
  options: initOptions,
  _cleanup: [],
  _pending: new Map(),
  _seq: 0,

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
    for (const pending of this._pending.values()) {
      pending?.abortController?.abort?.()
    }
    this._pending.clear()
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

    const url = ajaxUrl || hxUrl
    if (!url) return

    const panel = document.getElementById(tab)
    if (!panel || panel.getAttribute('data-tab-loaded') === 'true') return

    const prev = this._pending.get(tab)
    if (prev?.abortController) prev.abortController.abort()

    const abortController = new AbortController()
    const seq = (this._seq += 1)
    this._pending.set(tab, { abortController, seq })

    const alpine = storedAlpine || window.Alpine

    const rounded = Uno.getStyle('base', 'rounded', 'rounded-lg')
    const loadingClass = Uno.getStyle(
      'tab',
      'loading',
      `animate-pulse h-32 w-full ${rounded} bg-gray-200 dark:bg-slate-700`
    )
    const errorClass = Uno.getStyle('tab', 'error', 'p-4 text-red-600 dark:text-red-400')

    setHTML(panel, `<div class="${loadingClass}"></div>`, { trusted: false, alpine })

    try {
      const data = await requestData(url, { cache: true, signal: abortController.signal })
      const current = this._pending.get(tab)
      if (!current || current.seq !== seq) return

      const content =
        typeof data === 'string'
          ? data
          : `<pre class="whitespace-pre-wrap text-sm">${escapeHTML(
              JSON.stringify(data, null, 2)
            )}</pre>`

      setHTML(panel, content, { trusted: false, alpine })
      panel.setAttribute('data-tab-loaded', 'true')
    } catch (err) {
      const current = this._pending.get(tab)
      if (!current || current.seq !== seq) return
      if (err?.name === 'AbortError') return

      const msg = err?.message || 'Unknown error'
      const prefix = hxUrl ? 'Error loading content via HTMX: ' : 'Error loading content: '
      setHTML(panel, `<div class="${errorClass}">${escapeHTML(prefix)}${escapeHTML(msg)}</div>`, {
        trusted: false,
        alpine
      })
    } finally {
      const current = this._pending.get(tab)
      if (current && current.seq === seq) this._pending.delete(tab)
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
