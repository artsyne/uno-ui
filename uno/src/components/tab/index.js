/**
 * Uno - Tab Component (Store-First, Directive-Driven Architecture)
 *
 * A lightweight, declarative, HTMX-native tab component.
 * Follows the "Store-First, Directive-Driven" pattern.
 *
 * Features:
 * - Store-based state management
 * - Declarative x-uno-tab directive
 * - HTMX integration
 * - Hash synchronization
 * - Accessibility (ARIA)
 *
 * @module uno.components.tab
 * @version 0.6.0
 * @license MIT
 */

import { Uno } from '../../core/uno'
import { registerPlugin } from '../../core/plugin-registry'
import { createTabStore } from './tab.store'
import { registerTabDirectives } from './tab.directive'

export default (Alpine, options = {}) => {
  // Register plugin defaults
  const pluginDefaults = {
    defaultTab: '',
    hash: true,
    activeClass: 'is-active',
    inactiveClass: ''
  }
  Uno.registerPluginDefaults('tab', { ...pluginDefaults, ...options.plugins?.tab, ...options.tab })

  // Register Store
  const tabStore = createTabStore(Alpine)
  registerPlugin(Alpine, 'tab', tabStore)

  // Explicitly initialize the store
  Alpine.store('uno').tab.init()

  // Register Directive
  registerTabDirectives(Alpine)
}
