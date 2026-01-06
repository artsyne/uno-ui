/**
 * Uno - Box Component (Store-First, Directive-Driven Architecture)
 *
 * A lightweight, declarative, HTMX-native modal component.
 * Follows the "Store-First, Directive-Driven" pattern.
 *
 * Features:
 * - 9-point grid positioning system
 * - HTMX native content loading
 * - Automatic XSS sanitization (trustedHTML removed for security)
 * - Focus trap and keyboard navigation
 * - i18n support for button text
 * - Callback mechanism (onOk, onCancel, onBeforeClose, onClose, onContentLoad)
 *
 * @module uno.components.box
 * @version 0.6.0
 * @license MIT
 */

import { Uno } from '../../core/uno'
import { registerPlugin } from '../../core/plugin-registry'
import { createBoxStore } from './box.store'
import { registerBoxDirectives } from './box.directive'

/**
 * Box Plugin Factory
 * @param {object} Alpine - Alpine.js instance
 * @param {object} options - Plugin options
 */
export default (Alpine, options = {}) => {
  // Register plugin defaults
  const pluginDefaults = {
    defaultWidth: 400,
    defaultPosition: 'center',
    okText: 'OK',
    cancelText: 'Cancel',
    closeOnEsc: false,
    closeOnOverlayClick: false
  }
  Uno.registerPluginDefaults('box', { ...pluginDefaults, ...options.plugins?.box, ...options.box })

  // Register Store
  const boxStore = createBoxStore(Alpine)
  registerPlugin(Alpine, 'box', boxStore)

  // Register Directives
  registerBoxDirectives(Alpine)
}
