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
import tooltip from './components/tooltip/tooltip'
import notify from './components/notify/notify'
import box from './components/box/box'
import valid from './components/valid/valid'
import tab, { factory as unoTabsFactory } from './components/tab/tab'

window.Alpine = Alpine
window.htmx = htmx
if (typeof window !== 'undefined') {
  window.unoTabs = unoTabsFactory
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
    window.deferLoadingAlpine = (start) => {
      try {
        register(window.Alpine)
      } catch {
        /* ignore */
      }
      start()
    }
  }

  // Safety: if Alpine is added later and the above hook missed
  document.addEventListener('alpine:init', () => {
    register(window.Alpine)
  })
}

init()

// Do not auto-start Alpine here; use CDN to start it
