/**
 * Uno Plugin Registry - Unified plugin registration system
 *
 * Solves the Alpine.js reactivity issue when adding properties to existing stores.
 * All plugins should register their stores through this registry.
 *
 * @module uno.core.plugin-registry
 * @version 0.6.0
 * @license MIT
 */

/**
 * Register a plugin store under the 'uno' namespace
 *
 * This ensures proper Alpine.js reactivity by re-registering the entire
 * store object when adding new plugins, rather than mutating existing stores.
 *
 * @param {object} Alpine - Alpine.js instance
 * @param {string} pluginName - Plugin name (e.g., 'box', 'notify', 'tooltip')
 * @param {object} pluginStore - Plugin store object
 * @example
 * registerPlugin(Alpine, 'box', boxStore)
 */
export function registerPlugin(Alpine, pluginName, pluginStore) {
  if (!Alpine || typeof Alpine.store !== 'function') {
    console.error('[Uno] Alpine.store is not available')
    return
  }

  const existingUnoStore = Alpine.store('uno')

  if (!existingUnoStore) {
    // First plugin: create new 'uno' store
    Alpine.store('uno', { [pluginName]: pluginStore })
  } else {
    // Subsequent plugins: re-register entire store with new plugin
    // This ensures Alpine.js reactivity for all nested properties
    Alpine.store('uno', {
      ...existingUnoStore,
      [pluginName]: pluginStore
    })
  }

  // 兼容层：保留 Alpine.store('<pluginName>') 的旧用法。
  // 注意：这里保存的是已注册到 Alpine 的响应式对象（Proxy），避免后续失去响应性。
  if (!Alpine.store(pluginName)) {
    Alpine.store(pluginName, Alpine.store('uno')?.[pluginName])
  }
}

/**
 * Get the 'uno' store or a specific plugin store
 *
 * @param {object} Alpine - Alpine.js instance
 * @param {string} [pluginName] - Optional plugin name to get specific store
 * @returns {object|null} Store object or null if not found
 * @example
 * const unoStore = getStore(Alpine)
 * const boxStore = getStore(Alpine, 'box')
 */
export function getStore(Alpine, pluginName = null) {
  if (!Alpine || typeof Alpine.store !== 'function') {
    return null
  }

  const unoStore = Alpine.store('uno')
  if (!unoStore) {
    return null
  }

  return pluginName ? unoStore[pluginName] : unoStore
}
