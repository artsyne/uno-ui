/**
 * Type definitions for Uno UI
 *
 * @version 0.6.0
 * @license MIT
 */

/// <reference path="./alpinejs.d.ts" />

import type { Alpine as AlpineType } from 'alpinejs'

declare module 'uno-ui' {
  // ============================================================
  // Core Types
  // ============================================================

  /**
   * Uno global configuration options
   */
  export interface UnoOptions {
    /** Enable debug logging */
    debug?: boolean
    /** Content Security Policy configuration */
    csp?: {
      /** CSP nonce for inline scripts and styles */
      nonce?: string | null
    }
    /** Custom error handler */
    onError?: ((error: Error) => void) | null
    /** Plugin-specific configuration */
    plugins?: {
      box?: Partial<BoxPluginDefaults>
      notify?: Partial<NotifyPluginDefaults>
      tooltip?: Partial<TooltipPluginDefaults>
      tab?: Partial<TabPluginDefaults>
      valid?: Partial<ValidPluginDefaults>
    }
    /** UI styling configuration */
    ui?: UiConfig
  }

  /**
   * UI styling configuration
   */
  export interface UiConfig {
    /** Base styles used across components */
    base?: string
    /** Z-index layer configuration */
    layers?: {
      box?: number
      tooltip?: number
      notify?: number
    }
    /** Box component styles */
    box?: BoxUiConfig
    /** Notify component styles */
    notify?: NotifyUiConfig
    /** Tooltip component styles */
    tooltip?: TooltipUiConfig
    /** Tab component styles */
    tab?: TabUiConfig
    /** Valid component styles */
    valid?: ValidUiConfig
  }

  /**
   * Box component UI configuration
   */
  export interface BoxUiConfig {
    overlay?: string
    dialogWrapper?: string
    itemWrapper?: string
    container?: string
    cancelBtn?: string
    okBtn?: string
    title?: string
    closeBtnFloating?: string
    closeIconSmall?: string
    content?: string
    loadingWrap?: string
    spinner?: string
    footer?: string
    promptWrapper?: string
    promptText?: string
    promptInput?: string
  }

  /**
   * Notify component UI configuration
   */
  export interface NotifyUiConfig {
    container?: string
    item?: string
    content?: string
    row?: string
    iconWrap?: string
    textWrap?: string
    text?: string
    closeWrap?: string
    closeBtn?: string
    closeIcon?: string
    progressWrap?: string
    progressBar?: string
    types?: {
      success?: string
      error?: string
      warning?: string
      info?: string
    }
    icons?: {
      success?: string
      error?: string
      warning?: string
      info?: string
    }
  }

  /**
   * Tooltip component UI configuration
   */
  export interface TooltipUiConfig {
    content?: string
    body?: string
  }

  /**
   * Tab component UI configuration
   */
  export interface TabUiConfig {
    active?: string
    inactive?: string
    loading?: string
    error?: string
  }

  /**
   * Valid component UI configuration
   */
  export interface ValidUiConfig {
    wrapper?: string
    iconBase?: string
    inputError?: string
    inputValid?: string
    iconError?: string
    iconValid?: string
    stateInvalid?: string
    stateValid?: string
    stateValidating?: string
    inputPaddingIcon?: string
    helperText?: string
  }

  /**
   * Valid plugin default configuration
   */
  export interface ValidPluginDefaults {
    emptyQuit?: boolean
    debounce?: number
    scrollOnError?: boolean
  }

  /**
   * Plugin function type
   */
  export type Plugin = (Alpine: AlpineType, options?: any) => void

  /**
   * Uno instance interface
   */
  export interface UnoInstance {
    /** Uno version */
    version: string
    /** Global options */
    options: UnoOptions
    /** Plugin default configurations */
    pluginDefaults: Record<string, any>
    /** Registered plugins */
    plugins: Array<{ plugin: Plugin; options: any }>
    /** Initialization status */
    initialized: boolean
    /** DOM manager utilities */
    dom: DomManager
    /** HTMX utilities */
    htmx: HtmxUtilities

    /**
     * Initialize Uno
     * @param alpineOrOptions - Alpine instance or configuration options
     * @param options - Configuration options (if first arg is Alpine)
     * @returns Alpine instance
     */
    init(alpineOrOptions?: AlpineType | UnoOptions, options?: UnoOptions): AlpineType

    /**
     * Register a plugin
     * @param plugin - Plugin function
     * @param options - Plugin options
     * @returns UnoInstance for chaining
     */
    use(plugin: Plugin, options?: any): UnoInstance

    /**
     * Set global configuration
     * @param config - Configuration object to merge
     * @returns UnoInstance for chaining
     */
    setConfig(config: Partial<UnoOptions>): UnoInstance

    /**
     * Get merged configuration for a component
     * @param componentName - Component name
     * @param instanceOpts - Instance-specific options
     * @returns Merged configuration
     */
    getConfig(componentName: string, instanceOpts?: any): any

    /**
     * Register plugin default configuration
     * @param name - Component name
     * @param defaults - Default configuration
     */
    registerPluginDefaults(name: string, defaults: any): void

    /**
     * Get style class from configuration
     * @param comp - Component name
     * @param elem - Element name
     * @param defaultCls - Default class string
     * @returns Final class string
     */
    getStyle(comp: string, elem: string, defaultCls?: string): string
  }

  /**
   * DOM Manager utilities
   */
  export interface DomManager {
    /**
     * Ensure a container element exists in the DOM
     * @param id - Element ID
     * @param html - Inner HTML content
     * @param attributes - Attributes to set
     * @param parent - Parent element
     */
    ensureContainer(
      id: string,
      html?: string,
      attributes?: Record<string, string>,
      parent?: HTMLElement
    ): HTMLElement | null

    /**
     * Inject CSS styles
     * @param id - Style element ID
     * @param css - CSS content
     */
    injectStyles(id: string, css: string): void

    /**
     * Set HTML content safely
     * @param el - Target element
     * @param html - HTML content
     * @param options - Options
     */
    setHTML(
      el: HTMLElement,
      html: string,
      options?: { trusted?: boolean; alpine?: AlpineType }
    ): void

    /**
     * Escape HTML special characters
     * @param input - Input string
     * @returns Escaped string
     */
    escape(input: string): string

    /**
     * Sanitize HTML content
     * @param html - HTML string
     * @returns Sanitized HTML
     */
    sanitize(html: string): string
  }

  /**
   * HTMX utilities
   */
  export interface HtmxUtilities {
    /**
     * Check if HTMX is available
     * @returns True if HTMX is loaded
     */
    available(): boolean

    /**
     * HTMX ajax wrapper
     * @param method - HTTP method
     * @param url - Request URL
     * @param ajaxOptions - HTMX ajax options
     * @returns Promise resolving to result
     */
    ajax(method: string, url: string, ajaxOptions?: any): Promise<any>

    /**
     * Fetch HTML fragment without inserting into DOM
     * @param method - HTTP method
     * @param url - Request URL
     * @param ajaxOptions - HTMX ajax options
     * @returns Promise resolving to HTML string
     */
    fragment(method: string, url: string, ajaxOptions?: any): Promise<string>

    /**
     * Register HTMX request lifecycle handlers
     * @param handlers - Lifecycle handlers
     */
    onRequest(handlers: {
      start?: (event: any) => void
      end?: (event: any) => void
      progress?: (event: any) => void
    }): void

    /**
     * Bind notify component to HTMX events
     */
    bindNotifyToHtmx(): void
  }

  // ============================================================
  // Box Component Types
  // ============================================================

  /**
   * Box dialog position
   */
  export type BoxPosition =
    | 'center'
    | 'top-left'
    | 'top-center'
    | 'top-right'
    | 'center-left'
    | 'center-right'
    | 'bottom-left'
    | 'bottom-center'
    | 'bottom-right'

  /**
   * Box dialog type
   */
  export type BoxType = 'dialog' | 'confirm'

  /**
   * Box plugin default configuration
   */
  export interface BoxPluginDefaults {
    defaultWidth?: number
    defaultPosition?: BoxPosition
    okText?: string
    cancelText?: string
    closeOnEsc?: boolean
    closeOnOverlayClick?: boolean
  }

  /**
   * Box dialog options
   */
  export interface BoxOptions {
    /** Dialog type */
    type?: BoxType
    /** Dialog title */
    title?: string
    /** Dialog content (auto-sanitized) */
    content?: string
    /** Position */
    position?: BoxPosition
    /** Width in pixels */
    width?: number
    /** Height in pixels */
    height?: number
    /** Width CSS class */
    widthClass?: string
    /** Height CSS class */
    heightClass?: string
    /** Additional container class */
    containerClass?: string
    /** Show modal overlay */
    modal?: boolean
    /** Close on Escape key */
    closeOnEsc?: boolean
    /** Close on overlay click */
    closeOnOverlayClick?: boolean
    /** Show close button */
    showCloseButton?: boolean
    /** Show OK button */
    showOkButton?: boolean
    /** Show Cancel button */
    showCancelButton?: boolean
    /** OK button text */
    okText?: string
    /** Cancel button text */
    cancelText?: string
    /** Custom box ID */
    boxId?: string
    /** HTMX content URL */
    hxUrl?: string
    /** HTMX request method */
    hxMethod?: 'GET' | 'POST'
    /** HTMX request headers */
    hxHeaders?: Record<string, string>
    /** HTMX trigger condition */
    hxTrigger?: string
    /** OK button callback (return false to prevent close) */
    onOk?: () => void | false
    /** Cancel button callback */
    onCancel?: () => void
    /** Before close callback */
    onBeforeClose?: () => void
    /** After close callback */
    onClose?: () => void
    /** Content loaded callback */
    onContentLoad?: (el: HTMLElement) => void
  }

  /**
   * Box item instance
   */
  export interface BoxItem extends Required<BoxOptions> {
    /** Box instance ID */
    id: string
    /** Open state */
    isOpen: boolean
    /** Close this box */
    close(): void
  }

  /**
   * Box store interface
   */
  export interface BoxStore {
    /** Active box items */
    items: BoxItem[]

    /**
     * Initialize box component
     */
    init(): void

    /**
     * Open a box dialog
     * @param typeOrOpts - Type string or options object
     * @param opts - Options if first arg is type
     * @returns Box item instance
     */
    open(typeOrOpts: BoxType | BoxOptions, opts?: BoxOptions): BoxItem

    /**
     * Show an alert dialog
     * @param message - Alert message
     * @param options - Additional options
     * @returns Box item instance
     */
    alert(message: string, options?: BoxOptions | string): BoxItem

    /**
     * Show a confirm dialog
     * @param message - Confirmation message
     * @param options - Additional options
     * @returns Promise resolving to true if confirmed, false otherwise
     */
    confirm(message: string, options?: BoxOptions): Promise<boolean>

    /**
     * Show a prompt dialog
     * @param message - Prompt message
     * @param defaultValue - Default input value
     * @param options - Additional options
     * @returns Promise resolving to input value or null
     */
    prompt(message: string, defaultValue?: string, options?: BoxOptions): Promise<string | null>

    /**
     * Show a dialog (alias for open)
     * @param options - Dialog options
     * @returns Box item instance
     */
    dialog(options: BoxOptions): BoxItem

    /**
     * Close a specific box by ID
     * @param id - Box ID
     */
    close(id: string): void

    /**
     * Close all open boxes
     */
    closeAll(): void

    /**
     * Get position CSS classes
     * @param position - Position name
     * @returns CSS classes
     */
    getPositionClasses(position: BoxPosition): string

    /**
     * Handle HTMX content load
     * @param el - Content element
     * @param item - Box item
     */
    onHtmxContentLoaded(el: HTMLElement, item: BoxItem): void

    /**
     * Ensure container exists in DOM
     */
    ensureContainer(): HTMLElement | null
  }

  // ============================================================
  // Notify Component Types
  // ============================================================

  /**
   * Notify position
   */
  export type NotifyPosition =
    | 'top-left'
    | 'top-center'
    | 'top-right'
    | 'center-left'
    | 'center'
    | 'center-right'
    | 'bottom-left'
    | 'bottom-center'
    | 'bottom-right'

  /**
   * Notify type
   */
  export type NotifyType = 'success' | 'error' | 'warning' | 'info'

  /**
   * Notify plugin default configuration
   */
  export interface NotifyPluginDefaults {
    duration?: number
    position?: NotifyPosition
  }

  /**
   * Notify options
   */
  export interface NotifyOptions {
    /** Notification type */
    type?: NotifyType
    /** Position */
    position?: NotifyPosition
    /** Duration in ms (0 = no auto-dismiss) */
    duration?: number
    /** Show close button */
    closable?: boolean
    /** Icon class */
    icon?: string
    /** HTMX content URL */
    hxUrl?: string
    /** HTMX request headers */
    hxHeaders?: Record<string, string>
    /** HTMX trigger condition */
    hxTrigger?: string
    /** HTMX swap strategy */
    hxSwap?: string
    /** Content loaded callback */
    onContentLoad?: (el: HTMLElement) => void
    /** Error callback */
    onError?: (error: any) => void
  }

  /**
   * Notify item
   */
  export interface NotifyItem extends NotifyOptions {
    /** Notification ID */
    id: number
    /** Message text */
    message: string
    /** Visible state */
    isVisible: boolean
    /** Content for HTMX-loaded notifications */
    content?: string
  }

  /**
   * Notify store interface
   */
  export interface NotifyStore {
    /** Notification groups by position */
    groups: Record<NotifyPosition, NotifyItem[]>
    /** Next notification ID */
    nextId: number
    /** Active progress requests */
    progressRequests: number
    /** Progress bar visible state */
    progressVisible: boolean
    /** Progress bar percentage */
    progressPercent: number

    /**
     * Initialize notify component
     */
    init(): void

    /**
     * Show a notification
     * @param message - Notification message
     * @param opts - Notification options
     * @returns Notification ID
     */
    show(message: string, opts?: NotifyOptions): number

    /**
     * Show success notification
     * @param message - Message text
     * @param opts - Additional options
     * @returns Notification ID
     */
    success(message: string, opts?: NotifyOptions): number

    /**
     * Show error notification
     * @param message - Message text
     * @param opts - Additional options
     * @returns Notification ID
     */
    error(message: string, opts?: NotifyOptions): number

    /**
     * Show warning notification
     * @param message - Message text
     * @param opts - Additional options
     * @returns Notification ID
     */
    warning(message: string, opts?: NotifyOptions): number

    /**
     * Show info notification
     * @param message - Message text
     * @param opts - Additional options
     * @returns Notification ID
     */
    info(message: string, opts?: NotifyOptions): number

    /**
     * Remove a notification by ID
     * @param id - Notification ID
     */
    remove(id: number): void

    /**
     * Start progress bar
     */
    startProgress(): void

    /**
     * Stop progress bar
     */
    stopProgress(): void

    /**
     * Get container style
     */
    getContainerStyle(): { zIndex: number }

    /**
     * Get container CSS classes
     * @param position - Position name
     * @returns CSS classes
     */
    getContainerClass(position: NotifyPosition): string

    /**
     * Ensure container exists in DOM
     */
    ensureContainer(): void

    /**
     * Handle HTMX content loaded
     * @param el - Content element
     * @param item - Notify item
     */
    onContentLoaded(el: HTMLElement, item: NotifyItem): void

    /**
     * Handle HTMX error
     * @param error - Error details
     * @param item - Notify item
     */
    onHtmxError(error: any, item: NotifyItem): void
  }

  // ============================================================
  // Tooltip Component Types
  // ============================================================

  /**
   * Tooltip placement
   */
  export type TooltipPlacement =
    | 'top'
    | 'top-start'
    | 'top-end'
    | 'bottom'
    | 'bottom-start'
    | 'bottom-end'
    | 'left'
    | 'left-start'
    | 'left-end'
    | 'right'
    | 'right-start'
    | 'right-end'

  /**
   * Tooltip trigger mode
   */
  export type TooltipTrigger = 'hover' | 'click' | 'focus'

  /**
   * Tooltip plugin default configuration
   */
  export interface TooltipPluginDefaults {
    placement?: TooltipPlacement
    trigger?: TooltipTrigger
    offset?: number
    arrow?: boolean
    clickOutside?: boolean
  }

  /**
   * Tooltip options
   */
  export interface TooltipOptions {
    /** Tooltip content */
    content?: string
    /** Content target selector */
    target?: string
    /** Placement */
    placement?: TooltipPlacement
    /** Trigger mode */
    trigger?: TooltipTrigger
    /** Offset in pixels */
    offset?: number
    /** Show arrow */
    arrow?: boolean
    /** Close on click outside */
    clickOutside?: boolean
    /** Trigger element */
    element?: HTMLElement
    /** HTMX content URL */
    hxUrl?: string
    /** HTMX request headers */
    hxHeaders?: Record<string, string>
    /** HTMX trigger condition */
    hxTrigger?: string
  }

  /**
   * Tooltip store interface
   */
  export interface TooltipStore {
    /** Visible state */
    visible: boolean
    /** Tooltip content */
    content: string
    /** Placement */
    placement: TooltipPlacement
    /** Top position */
    top: number
    /** Left position */
    left: number
    /** Z-index */
    zIndex: number
    /** Trigger element */
    triggerEl: HTMLElement | null
    /** Trigger mode */
    triggerMode: TooltipTrigger
    /** Click outside to close */
    clickOutside: boolean
    /** Show arrow */
    arrow: boolean
    /** Arrow position */
    arrowPosition: { top: number; left: number }
    /** HTMX URL */
    hxUrl: string | null
    /** HTMX headers */
    hxHeaders: Record<string, string> | null
    /** HTMX trigger */
    hxTrigger: string

    /**
     * Initialize tooltip component
     */
    init(): void

    /**
     * Show tooltip
     * @param instanceOpts - Tooltip options
     */
    show(instanceOpts?: TooltipOptions): void

    /**
     * Hide tooltip
     */
    hide(): void

    /**
     * Cancel scheduled hide
     */
    cancelHide(): void

    /**
     * Schedule tooltip hide
     * @param delay - Delay in ms
     */
    scheduleHide(delay?: number): void

    /**
     * Update tooltip position
     */
    updatePosition(): void

    /**
     * Ensure container exists in DOM
     */
    ensureContainer(): void
  }

  // ============================================================
  // Tab Component Types
  // ============================================================

  /**
   * Tab plugin default configuration
   */
  export interface TabPluginDefaults {
    group?: string
    hash?: boolean
  }

  /**
   * Tab options
   */
  export interface TabOptions {
    /** Tab name */
    name?: string
    /** Tab group */
    group?: string
    /** Enable hash synchronization */
    hash?: boolean
    /** Update hash on activation */
    updateHash?: boolean
  }

  /**
   * Tab store interface
   */
  export interface TabStore {
    /** Active tabs by group */
    groups: Record<string, string>
    /** Tab registry */
    registry: Record<string, { group: string; name: string }>

    /**
     * Initialize tab component
     */
    init(): void

    /**
     * Activate a tab
     * @param instanceOpts - Tab options
     */
    activate(instanceOpts?: TabOptions): void

    /**
     * Check if tab is active
     * @param name - Tab name
     * @param group - Tab group
     * @returns True if active
     */
    isActive(name: string, group?: string): boolean

    /**
     * Register a tab
     * @param name - Tab name
     * @param group - Tab group
     */
    register(name: string, group?: string): void
  }

  // ============================================================
  // Valid Component Types
  // ============================================================

  /**
   * Valid plugin default configuration
   */
  export interface ValidPluginDefaults {
    realtime?: boolean
    showError?: boolean
    scrollOnError?: boolean
  }

  /**
   * Validation rule
   */
  export type ValidationRule =
    | 'required'
    | 'email'
    | 'url'
    | 'number'
    | 'min'
    | 'max'
    | 'minlength'
    | 'maxlength'
    | 'pattern'
    | string

  /**
   * Validation result
   */
  export interface ValidationResult {
    /** Validation passed */
    valid: boolean
    /** Error message */
    message?: string
  }

  /**
   * Valid store interface
   */
  export interface ValidStore {
    /**
     * Validate a form element
     * @param form - Form element
     * @returns True if valid
     */
    validate(form: HTMLFormElement): boolean

    /**
     * Validate a single field
     * @param field - Input element
     * @returns Validation result
     */
    validateField(
      field: HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement
    ): ValidationResult

    /**
     * Clear validation errors
     * @param form - Form element
     */
    clearErrors(form: HTMLFormElement): void

    /**
     * Show validation error
     * @param field - Input element
     * @param message - Error message
     */
    showError(field: HTMLElement, message: string): void

    /**
     * Register custom validator
     * @param name - Validator name
     * @param fn - Validator function
     */
    registerValidator(
      name: string,
      fn: (value: string, rule: string, field: HTMLElement) => boolean | string
    ): void
  }

  // ============================================================
  // I18n Store Types
  // ============================================================

  /**
   * I18n messages structure
   */
  export interface I18nMessages {
    [key: string]: any
  }

  /**
   * I18n store interface
   */
  export interface I18nStore {
    /** Current locale */
    locale: string
    /** Messages by locale */
    messages: Record<string, I18nMessages>
    /** Loading state */
    loading: boolean

    /**
     * Initialize i18n
     */
    init(): void

    /**
     * Translate a key
     * @param key - Translation key
     * @param vars - Variables to interpolate
     * @returns Translated text
     */
    t(key: string, vars?: Record<string, any>): string

    /**
     * Set current locale
     * @param lang - Locale code
     */
    setLocale(lang: string): void

    /**
     * Set current locale (alias)
     * @param lang - Locale code
     */
    setLang(lang: string): void

    /**
     * Load locale from URL
     * @param lang - Locale code
     * @param url - JSON file URL
     * @returns Promise resolving to true if successful
     */
    loadLocale(lang: string, url: string): Promise<boolean>

    /**
     * Extend messages
     * @param locales - Messages by locale
     */
    extend(locales: Record<string, I18nMessages>): void
  }

  // ============================================================
  // Alpine Store Declaration
  // ============================================================

  /**
   * Uno store under Alpine
   */
  export interface UnoStores {
    box: BoxStore
    notify: NotifyStore
    tooltip: TooltipStore
    tab: TabStore
    valid: ValidStore
  }
}

// ============================================================
// Alpine.js Store Extension
// ============================================================

declare module 'alpinejs' {
  interface Stores {
    /** Uno component stores */
    uno: import('uno-ui').UnoStores
    /** I18n store */
    i18n: import('uno-ui').I18nStore
  }
}

// ============================================================
// Window Global Declaration
// ============================================================

declare global {
  interface Window {
    /** Global Uno instance */
    Uno: import('uno-ui').UnoInstance
    /** Global Uno configuration */
    UnoConfig?: import('uno-ui').UnoOptions
    /** Alpine.js instance */
    Alpine: AlpineType
    /** HTMX instance */
    htmx?: any
  }
}

// ============================================================
// Re-export for convenience
// ============================================================

import type { UnoInstance, DomManager as DomManagerType, Plugin } from 'uno-ui'

declare const Uno: UnoInstance
declare const DomManager: DomManagerType
declare const box: Plugin
declare const notify: Plugin
declare const tooltip: Plugin
declare const tab: Plugin
declare const valid: Plugin

export { Uno as default, DomManager, box, notify, tooltip, tab, valid }
