/**
 * Uno - A lightweight, modern UI component library.
 * Built on Alpine.js and Tailwind CSS.
 *
 * @module uno
 * @version 0.3.0
 * @license MIT
 */

import Alpine from 'alpinejs'
import htmx from 'htmx.org'
import { Uno } from './core/uno'
import { installAlpine } from './core/i18n'
import { DomManager } from './utils/dom'
import tooltip from './components/tooltip'
import notify from './components/notify'
import box from './components/box'
import valid from './components/valid'
import tab, { factory as unoTabsFactory } from './components/tab'

// Attach DomManager to Uno instance
Uno.dom = DomManager

// 1. Expose Dependencies & Uno Instance
if (typeof window !== 'undefined') {
  window.Alpine = window.Alpine || Alpine
  window.htmx = window.htmx || htmx
  window.unoTabs = window.unoTabs || unoTabsFactory
  window.Uno = Uno
}

// 2. Register Built-in Plugins
Uno.use(installAlpine).use(tooltip).use(notify).use(box).use(valid).use(tab)

// 3. Initialization Logic
function initializeUno() {
  if (Uno.initialized) return

  const config = window.UnoConfig || {}
  Uno.init(window.Alpine, config)
}

if (typeof window !== 'undefined') {
  // Scenario A: CDN / Script Tag with deferLoadingAlpine
  const previous = window.deferLoadingAlpine
  window.deferLoadingAlpine = (start) => {
    initializeUno()
    if (typeof previous === 'function') {
      previous(start)
    } else {
      start()
    }
  }

  // Scenario B: Listen for Alpine init
  document.addEventListener('alpine:init', () => {
    if (!Uno.initialized) {
      initializeUno()
    }
  })

  // Scenario C: DOMContentLoaded
  document.addEventListener('DOMContentLoaded', () => {
    if (!Uno.initialized && window.Alpine) {
      initializeUno()
    }
  })

  // Scenario D: Sync
  if (!Uno.initialized && window.Alpine) {
    initializeUno()
  }
}

export default Uno
