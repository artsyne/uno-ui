/**
 * Minimal Alpine.js type definitions for Uno UI
 *
 * @version 3.x
 */

declare module 'alpinejs' {
  export interface Alpine {
    /**
     * Alpine version
     */
    version: string

    /**
     * Register a global store
     * @param name - Store name
     * @param value - Store value or factory function
     */
    store<T = any>(name: string, value: T | (() => T)): void

    /**
     * Get a registered store
     * @param name - Store name
     */
    store<T = any>(name: string): T

    /**
     * Register an Alpine directive
     * @param name - Directive name
     * @param callback - Directive handler
     */
    directive(name: string, callback: DirectiveCallback): void

    /**
     * Register an Alpine magic property
     * @param name - Magic property name
     * @param callback - Magic property handler
     */
    magic(name: string, callback: MagicCallback): void

    /**
     * Register a plugin
     * @param callback - Plugin function
     */
    plugin(callback: PluginCallback): void

    /**
     * Start Alpine
     */
    start(): void

    /**
     * Clone an element's Alpine scope
     * @param element - Element to clone
     */
    clone(element: HTMLElement): any

    /**
     * Evaluate an Alpine expression
     * @param element - Element context
     * @param expression - Expression string
     * @param extras - Additional scope data
     */
    evaluate(element: HTMLElement, expression: string, extras?: any): any

    /**
     * Initialize Alpine on an element
     * @param element - Element to initialize
     */
    initTree(element: HTMLElement): void

    /**
     * Destroy Alpine on an element
     * @param element - Element to destroy
     */
    destroyTree(element: HTMLElement): void

    /**
     * Get Alpine data for an element
     * @param element - Element
     */
    $data(element: HTMLElement): any
  }

  /**
   * Directive callback function
   */
  export type DirectiveCallback = (
    el: HTMLElement,
    directive: {
      value: string
      modifiers: string[]
      expression: string
    },
    utilities: {
      Alpine: Alpine
      effect: (callback: () => void) => void
      cleanup: (callback: () => void) => void
      evaluate: (expression: string) => any
      evaluateLater: (expression: string) => (callback?: (result: any) => void) => void
    }
  ) => void

  /**
   * Magic property callback function
   */
  export type MagicCallback = (el: HTMLElement, utilities: { Alpine: Alpine }) => any

  /**
   * Plugin callback function
   */
  export type PluginCallback = (Alpine: Alpine) => void

  /**
   * Alpine.js store registry interface
   */
  export interface Stores {
    [key: string]: any
  }

  const Alpine: Alpine
  export default Alpine
}
