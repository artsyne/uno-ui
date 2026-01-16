/**
 * Uno - A lightweight, modern UI component library.
 * Built on Alpine.js and Tailwind CSS.
 *
 * @module uno
 * @version 0.6.0
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
// Order matters: i18n should be first
Uno.use(installAlpine).use(tooltip).use(notify).use(box).use(valid).use(tab)

// 3. Auto-Registration with Alpine
// Uno automatically registers with Alpine when available.
// Configuration: Use Uno.setConfig() BEFORE Alpine loads.

if (typeof window !== 'undefined') {
  // Auto-register with Alpine when available
  const autoRegister = () => {
    if (window.Alpine && !Uno.initialized) {
      Uno.init(window.Alpine)
    }
  }

  // Try immediate registration (if Alpine already loaded)
  if (window.Alpine) {
    autoRegister()
  } else {
    // Wait for alpine:init event (defer scripts)
    document.addEventListener(
      'alpine:init',
      () => {
        autoRegister()
      },
      { once: true }
    )
  }
}

// Usage:
// <script type="module">
//   import Uno from './uno.es.js'
//   Uno.setConfig({ debug: true, ui: { /* ... */ } })
// </script>
// <script defer src="alpine.js"></script>
// <script defer src="uno.js"></script>

export { Uno }
export { DomManager }
export { default as tooltip } from './components/tooltip'
export { default as notify } from './components/notify'
export { default as box } from './components/box'
export { default as valid } from './components/valid'
export { default as tab } from './components/tab'

export default Uno
