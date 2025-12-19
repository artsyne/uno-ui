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
import { log } from './utils/log'
import tooltip from './components/tooltip'
import notify from './components/notify'
import box from './components/box'
import valid from './components/valid'
import tab, { factory as unoTabsFactory } from './components/tab'

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
const getWindowConfig = () => {
  if (typeof window === 'undefined') return {}
  const config = window.UnoConfig
  return config && typeof config === 'object' ? config : {}
}

const autoInit = () => {
  if (Uno.initialized) return

  // Check for configuration on global object
  const config = getWindowConfig()
  log.setDebug(config.debug !== false)

  // Pass window.Alpine if available to ensure plugins use the global instance
  Uno.init(typeof window !== 'undefined' ? window.Alpine : undefined, config)
}

if (typeof window !== 'undefined') {
  // Scenario A: CDN / Script Tag with deferLoadingAlpine
  // Alpine waits for this callback before starting.
  const previous = window.deferLoadingAlpine
  window.deferLoadingAlpine = (start) => {
    // Initialize Uno before Alpine starts
    autoInit()

    if (typeof previous === 'function') {
      previous(start)
    } else {
      start()
    }
  }

  // Scenario B: Module / Bundle or Late Load
  // If Alpine is already initialized or we missed the hook
  document.addEventListener('alpine:init', () => {
    if (!Uno.initialized) {
      const config = getWindowConfig()
      log.setDebug(config.debug !== false)
      if (config.debug !== false) {
        log.info('Uno: alpine:init detected, initializing...')
      }
      autoInit()
    }
  })

  // Scenario C: DOMContentLoaded Fallback
  // If Alpine is not present or something else is managing it
  document.addEventListener('DOMContentLoaded', () => {
    if (!Uno.initialized && window.Alpine) {
      autoInit()
    }
  })

  // Scenario D: Immediate (Synchronous) Auto-Init
  // We prioritize backward compatibility for the bundle.
  if (!Uno.initialized && window.Alpine) {
    const config = getWindowConfig()
    log.setDebug(config.debug !== false)
    log.info('Uno: Auto-initializing synchronously')
    autoInit()
  }
}

export default Uno
