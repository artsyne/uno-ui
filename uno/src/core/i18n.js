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
  Alpine.store('i18n', {
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

    async loadLocale(lang, url) {
      if (this.messages[lang]) return // Already loaded

      this.loading = true
      try {
        const res = await fetch(url)
        if (!res.ok) throw new Error(`Failed to load locale: ${res.statusText}`)
        const data = await res.json()
        this.extend({ [lang]: data })
        return true
      } catch (e) {
        log.error(`[Uno] Failed to load locale ${lang}`, e)
        throw e
      } finally {
        this.loading = false
      }
    },

    extend(locales) {
      for (const lang in locales) {
        if (!this.messages[lang]) {
          this.messages[lang] = locales[lang]
        } else {
          deepMerge(this.messages[lang], locales[lang])
        }
      }
    }
  })

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
