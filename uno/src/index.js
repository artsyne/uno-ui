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
import tab from './components/tab/tab'

window.Alpine = Alpine
window.htmx = htmx

// Debug log
log.info('Uno: Script loaded. Import Alpine:', !!Alpine, 'Window Alpine:', !!window.Alpine)

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
  let alpine = Alpine || window.Alpine

  if (alpine) {
    log.info('Uno: Alpine found, registering')
    register(alpine)
  } else {
    log.info('Uno: Alpine not found, waiting for alpine:init')
    document.addEventListener('alpine:init', () => {
      log.info('Uno: alpine:init fired')
      register(window.Alpine)
    })

    // Fallback: Check on DOMContentLoaded in case alpine:init missed
    document.addEventListener('DOMContentLoaded', () => {
      if (!window.Alpine) return
      // If not yet registered (checked inside register)
      log.info('Uno: DOMContentLoaded fallback')
      register(window.Alpine)
    })
  }
}

init()

// Alpine.start(); // Let the user or Alpine CDN start it
