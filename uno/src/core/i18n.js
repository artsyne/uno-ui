/**
 * Uno - Internationalization Utils
 *
 * @module uno.core.i18n
 * @version 0.1.0
 * @license MIT
 */

import en from '../locales/en'
import zhCN from '../locales/zh-CN'
import zhTW from '../locales/zh-TW'
import { log } from '../utils/log'

const initialMessages = {
  en,
  'zh-CN': zhCN,
  'zh-TW': zhTW
}

const supported = Object.keys(initialMessages)

function detect() {
  if (typeof localStorage !== 'undefined') {
    const saved = localStorage.getItem('uno.lang')
    if (saved && supported.includes(saved)) return saved
  }
  if (typeof navigator !== 'undefined') {
    const nav = navigator.language
    const lang = nav.startsWith('zh')
      ? nav.includes('TW') || nav.includes('Hant')
        ? 'zh-TW'
        : 'zh-CN'
      : 'en'
    return supported.includes(lang) ? lang : 'en'
  }
  return 'en'
}

function get(obj, path) {
  const parts = path.split('.')
  let cur = obj
  for (const p of parts) {
    if (!cur || typeof cur !== 'object' || !(p in cur)) return undefined
    cur = cur[p]
  }
  return cur
}

function deepMerge(target, source) {
  for (const key in source) {
    if (source[key] && typeof source[key] === 'object') {
      if (!target[key]) Object.assign(target, { [key]: {} })
      deepMerge(target[key], source[key])
    } else {
      Object.assign(target, { [key]: source[key] })
    }
  }
}

export function installAlpine(Alpine) {
  const existingStore = Alpine.store('i18n')
  const newStore = {
    locale: detect(),
    messages: initialMessages,
    loading: false,

    init() {
      if (typeof document !== 'undefined') {
        document.documentElement.lang = this.locale
      }
    },

    t(key, vars = {}) {
      // Access locale to ensure reactivity
      const currentLang = this.locale
      const dict = this.messages[currentLang] || this.messages['en']
      const msg = get(dict, key)
      const str = typeof msg === 'string' ? msg : key
      return str.replace(/\{(\w+)\}/g, (_, k) => (k in vars ? String(vars[k]) : `{${k}}`))
    },

    setLocale(lang) {
      if (!this.messages[lang]) {
        log.warn(`[Uno] Language ${lang} not found. Please load it first using loadLocale().`)
        return
      }
      this.locale = lang
      if (typeof localStorage !== 'undefined') localStorage.setItem('uno.lang', lang)
      if (typeof document !== 'undefined') document.documentElement.lang = lang
    },

    // Alias for backward compatibility
    setLang(lang) {
      this.setLocale(lang)
    },

    extend(newMessages) {
      deepMerge(this.messages, newMessages)
    },

    async loadLocale(lang, url) {
      if (this.messages[lang]) return // Already loaded

      this.loading = true
      try {
        const response = await fetch(url)
        if (!response.ok) throw new Error(`Failed to load ${lang}`)
        const data = await response.json()
        this.messages[lang] = data
        this.setLocale(lang)
      } catch (e) {
        log.error(`[Uno] Failed to load locale ${lang}:`, e)
      } finally {
        this.loading = false
      }
    }
  }

  if (existingStore) {
    log.info('[Uno] Merging into existing i18n store')
    // Merge messages and methods into existing store
    // We don't overwrite locale if it's already set
    if (!existingStore.messages) existingStore.messages = {}
    deepMerge(existingStore.messages, initialMessages)

    // Add methods if missing
    if (!existingStore.t) existingStore.t = newStore.t
    if (!existingStore.setLocale) existingStore.setLocale = newStore.setLocale
    if (!existingStore.setLang) existingStore.setLang = newStore.setLang
    if (!existingStore.loadLocale) existingStore.loadLocale = newStore.loadLocale

    // Ensure init is called or merged
    // Alpine stores don't support multiple init functions easily if we just attach methods
    // But if we are in the setup phase, we can wrap it?
    // Usually existingStore is the reactive proxy.
    // We just assume the user handles init if they created the store.
  } else {
    Alpine.store('i18n', newStore)
  }

  // Directive: x-t="key" or x-t.placeholder="key"
  Alpine.directive('t', (el, { expression, modifiers }, { effect }) => {
    const getI18n = () => Alpine.store('i18n')

    effect(() => {
      const store = getI18n()
      // Access locale to trigger reactivity
      // eslint-disable-next-line no-unused-vars
      const locale = store.locale
      const key = expression
      const text = store.t(key)

      if (modifiers.length > 0) {
        modifiers.forEach((attr) => {
          el.setAttribute(attr, text)
        })
      } else {
        el.textContent = text
      }
    })
  })
}
