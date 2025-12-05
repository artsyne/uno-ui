
/**
 * Uno - Tab Component
 *
 * @module uno.components.tab
 * @version 0.1.0
 * @license MIT
 */

export default (Alpine) => {
  const settings = {
    activeClass: 'text-blue-600 dark:text-blue-400 border-b-2 border-blue-500 dark:border-blue-400',
    inactiveClass: 'text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-200',
    baseTabClass: 'px-4 py-2 transition-colors duration-150 cursor-pointer'
  }

  /**
   * Apply active/inactive state to a tab element
   * @param {HTMLElement} el - Tab element
   * @param {boolean} isActive - Whether the tab is active
   */
  const applyTabState = (el, isActive) => {
    const activeCls = settings.activeClass.split(' ')
    const inactiveCls = settings.inactiveClass.split(' ')

    if (isActive) {
      el.classList.add('current')
      el.classList.remove(...inactiveCls)
      el.classList.add(...activeCls)
    } else {
      el.classList.remove('current')
      el.classList.remove(...activeCls)
      el.classList.add(...inactiveCls)
    }
  }

  /**
   * Switch to a specific tab
   * @param {string} group - Tab group name
   * @param {string} targetId - Target content ID
   */
  const switchTab = (group, targetId) => {
    // 1. Update Tab Buttons
    const groupTabs = document.querySelectorAll(`[data-tab-group="${group}"], [tabgroup="${group}"]`)
    groupTabs.forEach(tab => {
      const tid = tab.getAttribute('data-tab-target') || tab.getAttribute('tabtarget')
      const isMatch = tid === targetId
      applyTabState(tab, isMatch)
    })

    // 2. Update Content
    const targetsToManage = new Set()
    groupTabs.forEach(tab => {
      const tid = tab.getAttribute('data-tab-target') || tab.getAttribute('tabtarget')
      if (tid) targetsToManage.add(tid)
    })

    targetsToManage.forEach(tid => {
      const el = document.getElementById(tid)
      if (el) {
        if (tid === targetId) {
          el.classList.remove('hidden')
          el.style.display = ''
        } else {
          el.classList.add('hidden')
        }
      }
    })
  }

  /**
   * Handle hash change for deep linking
   */
  const handleHash = () => {
    const hash = window.location.hash.substring(1)
    if (!hash) return

    const tab = document.querySelector(`[data-tab-hash="${hash}"], [rel="${hash}"]`)
    if (tab) {
      const group = tab.getAttribute('data-tab-group') || tab.getAttribute('tabgroup')
      const targetId = tab.getAttribute('data-tab-target') || tab.getAttribute('tabtarget')
      if (group && targetId) {
        switchTab(group, targetId)
      }
    }
  }

  /**
   * Initialize tab styles
   */
  const initStyles = () => {
    const tabs = document.querySelectorAll('[data-tab-group], [tabgroup]')
    tabs.forEach(tab => {
      if (!tab.className.includes('transition-colors')) {
        const base = settings.baseTabClass.split(' ')
        tab.classList.add(...base)
      }
      const isActive = tab.classList.contains('current') || tab.getAttribute('data-tab-active') === 'true'
      applyTabState(tab, isActive)
    })
   }

  // Initialize
  document.addEventListener('click', (e) => {
    const target = e.target.closest('[data-tab-target], [tabtarget]')
    if (!target) return

    const group = target.getAttribute('data-tab-group') || target.getAttribute('tabgroup')
    const targetId = target.getAttribute('data-tab-target') || target.getAttribute('tabtarget')
    const hash = target.getAttribute('data-tab-hash') || target.getAttribute('rel')
    const ajaxUrl = target.getAttribute('data-tab-url')

    if (group && targetId) {
      e.preventDefault()
      
      if (ajaxUrl) {
        const contentEl = document.getElementById(targetId)
        if (contentEl && !contentEl.getAttribute('data-tab-loaded')) {
           contentEl.innerHTML = `<div class="animate-pulse p-4 space-y-2">
             <div class="h-4 bg-gray-200 rounded w-3/4"></div>
             <div class="h-4 bg-gray-200 rounded w-1/2"></div>
           </div>`
           fetch(ajaxUrl)
             .then(res => {
               if (!res.ok) throw new Error('Network response was not ok')
               return res.text()
             })
             .then(html => {
               contentEl.innerHTML = html
               contentEl.setAttribute('data-tab-loaded', 'true')
             })
             .catch(err => {
               contentEl.innerHTML = `<div class="text-red-500 p-4">Error loading content: ${err.message}</div>`
             })
        }
      }

      switchTab(group, targetId)
      if (hash) {
        history.replaceState(null, null, '#' + hash)
      }
    }
  })

  /**
   * Initialize component
   */
  const init = () => {
    handleHash()
    initStyles()
    window.addEventListener('hashchange', handleHash)
  }

  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', init)
  } else {
    init()
  }

  // Expose internal logic via Alpine store if needed for programmatic access
  /**
   * Tab store for programmatic control
   */
  Alpine.store('tabs', {
    switch(group, targetId) {
      switchTab(group, targetId)
    }
  })
}
