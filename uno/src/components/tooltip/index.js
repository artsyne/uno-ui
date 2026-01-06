/**
 * Uno - Tooltip Component (Store-First, Directive-Driven Architecture)
 *
 * A lightweight, declarative, performance-optimized tooltip component.
 * Follows the "Store-First, Directive-Driven" pattern.
 *
 * Features:
 * - RAF-based positioning with auto-flip
 * - Global singleton container (performance)
 * - Declarative directive interface (x-uno-tooltip)
 * - HTMX native integration
 * - 4-level configuration merge
 * - Styles decoupled via Uno.getStyle
 *
 * @module uno.components.tooltip
 * @version 0.3.0
 * @license MIT
 */

import { Uno } from '../../core/uno'
import { registerPlugin } from '../../core/plugin-registry'
import { createTooltipStore } from './tooltip.store'
import { registerTooltipDirectives } from './tooltip.directive'

const tooltipPlugin = (Alpine, options = {}) => {
  // Register plugin defaults
  const pluginDefaults = {
    placement: 'top',
    trigger: 'hover',
    delay: 300,
    offset: 8,
    arrow: true,
    clickOutside: true
  }
  Uno.registerPluginDefaults('tooltip', {
    ...pluginDefaults,
    ...options.plugins?.tooltip,
    ...options.tooltip
  })

  // Register Store
  const tooltipStore = createTooltipStore(Alpine)
  registerPlugin(Alpine, 'tooltip', tooltipStore)

  // Register Directive
  registerTooltipDirectives(Alpine)
}

export default tooltipPlugin
