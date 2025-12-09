/**
 * Uno - Tab Component
 *
 * @module uno.components.tab
 * @version 0.2.1
 * @license MIT
 */

export const factory = (defaultTab = '', options = {}) => ({
  active: defaultTab,
  options: {
    hash: true,
    manual: false, // if true, arrow keys only move focus, don't select
    ...options
  },
  loaded: {}, // Track loaded AJAX content

  init() {
    // Handle Hash on load
    if (this.options && this.options.hash) {
      // Use nextTick to ensure DOM is ready if needed, though init usually runs after DOM is parsed for this component
      this.$nextTick(() => {
        const hash = window.location.hash.substring(1)
        if (hash) {
          // Find the trigger associated with this hash
          const trigger = this.$el.querySelector(`[data-tab-hash="${hash}"]`)
          if (trigger) {
            const target = trigger.getAttribute('data-tab-target')
            if (target) this.active = target
          }
        }
      })
    }

    // Listen for hash changes
    if (this.options.hash) {
      window.addEventListener('hashchange', () => {
        const hash = window.location.hash.substring(1)
        if (hash) {
          const trigger = this.$el.querySelector(`[data-tab-hash="${hash}"]`)
          if (trigger) {
            const target = trigger.getAttribute('data-tab-target')
            if (target) this.active = target
          }
        }
      })
    }
  },

  isActive(tab) {
    return this.active === tab
  },

  async open(tab, ajaxUrl = null, hash = null, hxUrl = null) {
    this.active = tab

    // Handle Hash
    if (this.options.hash && hash) {
      // Use replaceState to avoid cluttering history if desired, or pushState
      // The spec test expects URL to contain hash
      // If we just set window.location.hash, it triggers hashchange event
      // which might cause loop if not careful, but here we set active first so it's fine.
      // But typically tab switching updates URL.
      // The test says "modern tabs switch content and hash".
      // Let's use history API to be clean.
      if (window.location.hash.substring(1) !== hash) {
        history.pushState(null, null, '#' + hash)
      }
    }

    // Handle HTMX
    if (hxUrl) {
      if (typeof window !== 'undefined' && window.htmx) {
        const contentEl = document.getElementById(tab)
        if (contentEl && !this.loaded[tab]) {
          // Loading state
          contentEl.innerHTML = `<div class="animate-pulse p-4 space-y-2">
               <div class="h-4 bg-gray-200 dark:bg-gray-700 rounded w-3/4"></div>
               <div class="h-4 bg-gray-200 dark:bg-gray-700 rounded w-1/2"></div>
             </div>`

          try {
            await window.htmx.ajax('GET', hxUrl, {
              target: `#${tab}`,
              swap: 'innerHTML'
            })
            contentEl.setAttribute('data-tab-loaded', 'true')
            this.loaded[tab] = true
            contentEl.dispatchEvent(new CustomEvent('tab-loaded', { bubbles: true }))
          } catch {
            contentEl.innerHTML = `<div class="text-red-500 p-4">Error loading content via HTMX</div>`
          }
        }
      } else {
        console.warn('Uno Tab: HTMX requested but window.htmx is not defined.')
        const contentEl = document.getElementById(tab)
        if (contentEl) {
          contentEl.innerHTML = `<div class="text-amber-600 bg-amber-50 p-4 rounded border border-amber-200">
            <strong>Warning:</strong> HTMX library not found. Please include htmx.org script.
          </div>`
        }
      }
      return
    }

    // Handle AJAX
    if (ajaxUrl) {
      const contentEl = document.getElementById(tab)
      // Check if already loaded or currently loading?
      // We use this.loaded[tab] to track success.
      if (contentEl && !this.loaded[tab]) {
        // Loading state
        contentEl.innerHTML = `<div class="animate-pulse p-4 space-y-2">
             <div class="h-4 bg-gray-200 dark:bg-gray-700 rounded w-3/4"></div>
             <div class="h-4 bg-gray-200 dark:bg-gray-700 rounded w-1/2"></div>
           </div>`

        try {
          const res = await fetch(ajaxUrl)
          if (!res.ok) throw new Error('Network response was not ok')
          const html = await res.text()
          contentEl.innerHTML = html
          contentEl.setAttribute('data-tab-loaded', 'true')
          this.loaded[tab] = true

          // Dispatch event for listeners (e.g. HTMX)
          contentEl.dispatchEvent(new CustomEvent('tab-loaded', { bubbles: true }))
        } catch (err) {
          contentEl.innerHTML = `<div class="text-red-500 p-4">Error loading content: ${err.message}</div>`
        }
      }
    }
  },

  trigger(tab, ajaxUrl = null, hash = null, hxUrl = null) {
    return {
      ['@click.prevent']: () => this.open(tab, ajaxUrl, hash, hxUrl),
      [':class']: () => (this.isActive(tab) ? 'current' : ''),
      ['role']: 'tab',
      ['aria-selected']: () => this.isActive(tab).toString(),
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
      ['x-show']: () => this.isActive(tab),
      ['role']: 'tabpanel',
      ['id']: tab
    }
  }
})

if (typeof window !== 'undefined') {
  window.unoTabs = factory
}

export default (Alpine) => {
  Alpine.data('unoTabs', factory)
}
