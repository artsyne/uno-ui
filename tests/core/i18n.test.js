import { describe, it, expect, beforeEach, vi } from 'vitest'
import { installAlpine } from '../../src/core/i18n'

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
      await store.loadLocale('it', '/locales/it.json')

      expect(global.fetch).toHaveBeenCalledWith('/locales/it.json')
      expect(store.messages['it']).toEqual({ loaded: 'Loaded Message' })

      store.setLang('it')
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
    // <html lang> 只由服务端输出：Uno 不再回写文档属性
    expect(document.documentElement.lang).toBe('en')
  })

  it('should fall back to the base language for regional tags', () => {
    for (const [tag, expected] of [
      ['ja-JP', 'ja'],
      ['ko-KR', 'ko'],
      ['pt-BR', 'pt'],
      ['de-CH', 'de']
    ]) {
      store = undefined
      global.localStorage.getItem.mockReturnValue(null)
      vi.stubGlobal('navigator', { language: tag, languages: [tag] })
      installAlpine(Alpine)
      expect(store.locale, tag).toBe(expected)
    }
    vi.unstubAllGlobals()
  })
})

describe('locale files', () => {
  // 键集与占位符一致性遍历 uno/src/locales 目录，新增语言无需改测试
  const modules = import.meta.glob('../../src/locales/*.js', { eager: true })
  const locales = Object.fromEntries(
    Object.entries(modules).map(([file, mod]) => [
      file.split('/').pop().replace(/\.js$/, ''),
      flatten(mod.default)
    ])
  )
  const baseKeys = Object.keys(locales.en).sort()

  it('should keep the same key set across every locale', () => {
    for (const [code, dict] of Object.entries(locales)) {
      expect(Object.keys(dict).sort(), code).toEqual(baseKeys)
    }
  })

  it('should keep placeholders consistent across every locale', () => {
    for (const [code, dict] of Object.entries(locales)) {
      for (const key of baseKeys) {
        const expected = placeholders(locales.en[key])
        expect(placeholders(dict[key]), `${code}: ${key}`).toEqual(expected)
      }
    }
  })
})
