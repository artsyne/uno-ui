import { describe, it, expect, beforeEach, vi } from 'vitest'
import { installAlpine } from '../../src/core/i18n'

describe('i18n Store', () => {
  let store
  let Alpine

  beforeEach(() => {
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
})
