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

const messages = {
  en,
  'zh-CN': zhCN,
  'zh-TW': zhTW
}

const supported = Object.keys(messages)

/**
 * Detect user language
 * @returns {string} Language code
 */
function detect() {
  const saved = typeof localStorage !== 'undefined' ? localStorage.getItem('uno.lang') : null
  if (saved && supported.includes(saved)) return saved
  const nav = typeof navigator !== 'undefined' ? navigator.language : 'en'
  const lang = nav.startsWith('zh') ? (nav.includes('TW') || nav.includes('Hant') ? 'zh-TW' : 'zh-CN') : 'en'
  return supported.includes(lang) ? lang : 'en'
}

let current = detect()

/**
 * Set current language
 * @param {string} lang - Language code
 */
function setLang(lang) {
  if (!supported.includes(lang)) return
  current = lang
  if (typeof localStorage !== 'undefined') localStorage.setItem('uno.lang', lang)
  if (typeof document !== 'undefined') document.documentElement.lang = lang
}

/**
 * Get current language
 * @returns {string} Language code
 */
function getLang() {
  return current
}

/**
 * Get value from object by path
 * @param {Object} obj - Source object
 * @param {string} path - Dot notation path
 * @returns {any} Value
 */
function get(obj, path) {
  const parts = path.split('.')
  let cur = obj
  for (const p of parts) {
    if (!cur || typeof cur !== 'object' || !(p in cur)) return undefined
    cur = cur[p]
  }
  return cur
}

/**
 * Translate key with variables
 * @param {string} key - Translation key
 * @param {Object} [vars={}] - Variables to replace
 * @returns {string} Translated string
 */
function t(key, vars = {}) {
  const msg = get(messages[current], key)
  const str = typeof msg === 'string' ? msg : key
  return str.replace(/\{(\w+)\}/g, (_, k) => (k in vars ? String(vars[k]) : `{${k}}`))
}

/**
 * Extend messages
 * @param {Object} locales - New locales to merge
 */
function extend(locales) {
  for (const lang in locales) {
    if (!messages[lang]) {
      messages[lang] = locales[lang]
      supported.push(lang)
    } else {
      const merge = (target, source) => {
        for (const key in source) {
          if (source[key] && typeof source[key] === 'object') {
            if (!target[key]) Object.assign(target, { [key]: {} })
            merge(target[key], source[key])
          } else {
            Object.assign(target, { [key]: source[key] })
          }
        }
      }
      merge(messages[lang], locales[lang])
    }
  }
}

/**
 * Install Alpine store
 * @param {Object} Alpine - Alpine instance
 */
function installAlpine(Alpine) {
  const store = {
    lang: current,
    setLang(lang) {
      if (!supported.includes(lang)) {
        // Try to find fuzzy match or extend supported?
        // For now, strict check but we might want to allow setting if it was extended.
        // If extend() was called, supported array is updated.
        // So this check is fine.
         if (!supported.includes(lang)) return 
      }
      setLang(lang)
      this.lang = lang
    },
    extend,
    t
  }
  Alpine.store('i18n', store)
}

export { messages, supported, setLang, getLang, t, extend, installAlpine }
