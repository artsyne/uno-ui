/**
 * Uno - Notify Component (Store-First, Directive-Driven Architecture)
 *
 * A lightweight, flexible toast notification system.
 * Follows the "Store-First, Directive-Driven" pattern.
 *
 * Features:
 * - Global notification management via Store
 * - Multiple positioning support (top-right, bottom-right, etc.)
 * - Auto-dismiss with visual progress bar
 * - HTMX integration (server-triggered notifications)
 * - i18n support for messages
 *
 * @module uno.components.notify
 * @version 0.6.0
 * @license MIT
 */

import { Uno } from '../../core/uno'
import { registerPlugin } from '../../core/plugin-registry'
import { createNotifyStore } from './notify.store'
import { registerNotifyDirectives } from './notify.directive'

export default (Alpine, options = {}) => {
  // Register plugin defaults
  const pluginDefaults = {
    duration: 3000,
    position: 'top-right'
  }

  Uno.registerPluginDefaults('notify', {
    ...pluginDefaults,
    ...options.plugins?.notify,
    ...options.notify
  })

  // Register Store
  const notifyStore = createNotifyStore(Alpine)
  registerPlugin(Alpine, 'notify', notifyStore)

  // Lazy initialization: Container is injected on first show() or progress
  // notifyStore.init() is removed

  // Register Directive
  registerNotifyDirectives(Alpine)
}
