/**
 * Uno - A lightweight, modern UI component library.
 * Built on Alpine.js and Tailwind CSS.
 *
 * @module uno
 * @version 0.1.0
 * @license MIT
 */

import Alpine from 'alpinejs'
import htmx from 'htmx.org'
import { installAlpine } from './core/i18n'
import { log } from './utils/log'
import tooltip from './components/tooltip'
import notify from './components/notify'
import box from './components/box'
import valid from './components/valid'
import tab, { factory as unoTabsFactory } from './components/tab'

if (typeof window !== 'undefined') {
  window.Alpine = window.Alpine || Alpine
  window.htmx = window.htmx || htmx
  window.unoTabs = window.unoTabs || unoTabsFactory
}

const register = (alpine) => {
  if (!alpine) {
    log.error('Uno: Alpine is undefined, cannot register plugins')
    return
  }
  // Check if already registered to avoid double registration
  if (alpine.store && alpine.store('tooltip')) {
    log.info('Uno: Plugins already registered')
    return
  }

  log.info('Uno: Registering plugins')
  installAlpine(alpine)
  tooltip(alpine)
  notify(alpine)
  box(alpine)
  valid(alpine)
  tab(alpine)
}

const init = () => {
  const alpine = Alpine || window.Alpine
  if (alpine) {
    log.info('Uno: Alpine found, registering')
    register(alpine)
  }

  // Ensure registration happens before Alpine starts when using CDN
  if (typeof window !== 'undefined') {
    const previous = window.deferLoadingAlpine
    window.deferLoadingAlpine = (start) => {
      let started = false
      const safeStart = () => {
        if (started) return
        started = true
        start()
      }

      try {
        register(window.Alpine)
      } catch {
        /* ignore */
      }

      if (typeof previous === 'function' && previous !== window.deferLoadingAlpine) {
        previous(safeStart)
      } else {
        safeStart()
      }
    }
  }

  // Safety: if Alpine is added later and the above hook missed
  document.addEventListener('alpine:init', () => {
    register(window.Alpine)
  })
}

if (typeof window !== 'undefined' && window.__UNO_INSTALLED__) {
  log.info('Uno: Already installed, skipping init')
} else {
  if (typeof window !== 'undefined') window.__UNO_INSTALLED__ = true
  init()
}

// Do not auto-start Alpine here; use CDN to start it
