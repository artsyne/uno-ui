/**
 * Uno - A lightweight, modern UI component library.
 * Built on Alpine.js and Tailwind CSS.
 *
 * @module uno
 * @version 0.5.0
 * @license MIT
 */

import { Uno } from './core/uno'
import { installAlpine } from './core/i18n'
import { DomManager } from './utils/dom'
import tooltip from './components/tooltip'
import notify from './components/notify'
import box from './components/box'
import valid from './components/valid'
import tab from './components/tab'

// Attach DomManager to Uno instance
Uno.dom = DomManager

// 1. Expose Dependencies & Uno Instance
if (typeof window !== 'undefined') {
  window.Uno = Uno
}

// 2. Register Built-in Plugins
Uno.use(installAlpine).use(tooltip).use(notify).use(box).use(valid).use(tab)

// 3. Initialization Logic
if (typeof window !== 'undefined' && typeof document !== 'undefined') {
  const initUno = () => {
    if (Uno.initialized) return
    const config = {}
    Uno.init(window.Alpine, config)
  }

  if (window.Alpine) {
    initUno()
  } else {
    document.addEventListener('alpine:init', initUno)
  }
}

export default Uno
