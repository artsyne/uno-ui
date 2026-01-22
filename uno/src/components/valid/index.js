/**
 * Uno - Valid Component (Store-First, Directive-Driven Architecture)
 *
 * A powerful, declarative form validation engine.
 * Follows the "Store-First, Directive-Driven" pattern.
 *
 * Features:
 * - Declarative rules via x-uno-valid directive (Array Config)
 * - Real-time and submit-time validation
 * - Custom validator functions support
 * - HTMX integration (blocks invalid requests, handles 422 errors)
 * - Auto-scroll to error and focus management
 * - i18n support for error messages
 *
 * @module uno.components.valid
 * @version 0.6.0
 * @license MIT
 */

import { Uno } from '../../core/uno'
import { registerPlugin } from '../../core/plugin-registry'
import { createValidStore } from './valid.store'
import { registerValidDirectives } from './valid.directive'

export default (Alpine, options = {}) => {
  // Register plugin defaults
  const pluginDefaults = {
    realtime: true,
    showError: true,
    scrollOnError: true
  }
  Uno.registerPluginDefaults('valid', {
    ...pluginDefaults,
    ...options.plugins?.valid,
    ...options.valid
  })

  // Initialize Store
  const validStore = createValidStore(Alpine)
  registerPlugin(Alpine, 'valid', validStore)

  // Register Directives
  registerValidDirectives(Alpine)
}
