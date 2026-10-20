import { describe, it, expect, beforeEach, vi } from 'vitest'
import { installAlpine } from '../../src/core/i18n'
import en from '../../src/locales/en'
import zhCN from '../../src/locales/zh-CN'
import zhTW from '../../src/locales/zh-TW'

function flatten(obj, prefix = '', out = {}) {
  for (const [k, v] of Object.entries(obj)) {
    const key = prefix ? `${prefix}.${k}` : k
    if (v && typeof v === 'object') flatten(v, key, out)
    else out[key] = v
  }
  return out
}

function placeholders(str) {
  return [...String(str).matchAll(/\{(\w+)\}/g)].map((m) => m[1]).sort()
}

describe('i18n Store', () => {
  let store
  let Alpine

  beforeEach(() => {
    // Reset store so each case re-runs detection
    store = undefined
    // Mock Alpine
    Alpine = {
      store: vi.fn((name, definition) => {
        if (definition) {
          store = definition
          store.init()
        }
        return store
      }),
      directive: vi.fn()
    }

    // Mock localStorage
    global.localStorage = {
      getItem: vi.fn(),
      setItem: vi.fn()
    }

    // Mock document
    global.document = {
      documentElement: { lang: 'en' }
    }
  })

  it('should initialize with default locale', () => {
    installAlpine(Alpine)
    expect(store.locale).toBe('en')
    expect(document.documentElement.lang).toBe('en')
  })

  it('should translate keys', () => {
    installAlpine(Alpine)
    // Assuming 'en' has some keys. If not, we might need to mock initialMessages or extend it.
    // Let's extend it to be sure.
    store.extend({
      en: { hello: 'Hello' },
      'zh-CN': { hello: '你好' }
    })

    expect(store.t('hello')).toBe('Hello')

    store.setLang('zh-CN')
    expect(store.t('hello')).toBe('你好')
  })

  it('should interpolate variables', () => {
    installAlpine(Alpine)
    store.extend({
      en: { welcome: 'Welcome, {name}!' }
    })

    // Explicitly set locale to ensure we are using the updated messages
    store.setLang('en')

    expect(store.t('welcome', { name: 'Alice' })).toBe('Welcome, Alice!')
  })

  it('should handle missing keys', () => {
    installAlpine(Alpine)
    expect(store.t('missing.key')).toBe('missing.key')
  })

  it('should support async loading (loadLocale)', async () => {
    installAlpine(Alpine)

    // Mock fetch
    global.fetch = vi.fn().mockResolvedValue({
      ok: true,
      json: async () => ({ loaded: 'Loaded Message' })
    })

    // We expect a loadLocale method to exist (TDD - Red)
    // Currently it doesn't exist.
    if (store.loadLocale) {
      await store.loadLocale('fr', '/locales/fr.json')

      expect(global.fetch).toHaveBeenCalledWith('/locales/fr.json')
      expect(store.messages['fr']).toEqual({ loaded: 'Loaded Message' })

      store.setLang('fr')
      expect(store.t('loaded')).toBe('Loaded Message')
    } else {
      // Fail if not implemented (which is expected now)
      throw new Error('loadLocale not implemented')
    }
  })

  it('should normalize saved language codes on detection', () => {
    for (const lang of ['zh-HK', 'zh-MO', 'zh-Hant-TW', 'zh_tw', 'zh-TW']) {
      store = undefined
      global.localStorage.getItem.mockReturnValue(lang)
      installAlpine(Alpine)
      expect(store.locale, lang).toBe('zh-TW')
    }
  })

  it('should detect zh-HK/zh-MO navigator locales as zh-TW', () => {
    for (const lang of ['zh-HK', 'zh-MO', 'zh-Hant']) {
      store = undefined
      vi.stubGlobal('navigator', { language: lang, languages: [lang] })
      installAlpine(Alpine)
      expect(store.locale, lang).toBe('zh-TW')
    }
    vi.unstubAllGlobals()
  })

  it('should normalize language codes in setLocale', () => {
    installAlpine(Alpine)
    store.setLocale('zh-tw')
    expect(store.locale).toBe('zh-TW')
    expect(global.localStorage.setItem).toHaveBeenCalledWith('uno.lang', 'zh-TW')
    expect(document.documentElement.lang).toBe('zh-TW')
  })
})

describe('locale files', () => {
  const locales = {
    en: flatten(en),
    'zh-CN': flatten(zhCN),
    'zh-TW': flatten(zhTW)
  }
  const baseKeys = Object.keys(locales.en).sort()

  it('should keep the same key set across locales', () => {
    expect(Object.keys(locales['zh-CN']).sort()).toEqual(baseKeys)
    expect(Object.keys(locales['zh-TW']).sort()).toEqual(baseKeys)
  })

  it('should keep placeholders consistent across locales', () => {
    for (const key of baseKeys) {
      const expected = placeholders(locales.en[key])
      expect(placeholders(locales['zh-CN'][key]), `zh-CN: ${key}`).toEqual(expected)
      expect(placeholders(locales['zh-TW'][key]), `zh-TW: ${key}`).toEqual(expected)
    }
  })
})
