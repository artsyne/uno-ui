import { getCookie, setCookie } from '../utils/cookie.js';
import localeData from '../locales/data.js';

const STORAGE_KEY = 'idiary-lang';
const DEFAULT_LANGUAGE = 'en';
const SUPPORTED_LANGUAGES = ['en', 'zh-cn', 'zh-tw', 'fr', 'ja', 'ko'];

export default {
  current: DEFAULT_LANGUAGE,
  locale: DEFAULT_LANGUAGE, // Sync with Uno

  init() {
    this.current = this.detectLanguage();
    this.locale = this.current; // Sync with Uno
    this.apply(this.current);
  },

  /**
   * Get translation for a key
   * @param {string} key - Translation key
   * @param {Object} params - Parameters for replacement (optional)
   * @returns {string} Translated text
   */
  t(key, params = {}) {
    const lang = this.current.toLowerCase();
    const fallbackLang = DEFAULT_LANGUAGE;

    // Try current language from localeData
    let text = localeData[lang]?.[key];

    // Helper to get from Uno messages
    const getFromUno = (langCode) => {
      if (!this.messages) {
        return undefined;
      }

      // Map 'zh-cn' to 'zh-CN' for Uno
      const mapToUnoLang = (code) => {
        if (code === 'zh-cn') return 'zh-CN';
        if (code === 'zh-tw') return 'zh-TW';
        return code;
      };

      const unoLang = mapToUnoLang(langCode);
      const dict = this.messages[unoLang];
      if (!dict) {
        return undefined;
      }

      const parts = key.split('.');
      let cur = dict;
      for (const p of parts) {
        if (!cur || typeof cur !== 'object' || !(p in cur)) return undefined;
        cur = cur[p];
      }
      return typeof cur === 'string' ? cur : undefined;
    };

    // If not in localeData, try Uno messages
    if (!text) {
      text = getFromUno(lang);
    }

    // Fallback to default language
    if (!text && lang !== fallbackLang) {
      text = localeData[fallbackLang]?.[key];
      // If still not found, try Uno messages in default language
      if (!text) {
        text = getFromUno(fallbackLang);
      }
    }

    // Return key if not found
    if (!text) {
      return key;
    }

    // Replace params if any
    Object.keys(params).forEach((param) => {
      text = text.replace(new RegExp(`{${param}}`, 'g'), params[param]);
    });

    return text;
  },

  detectLanguage() {
    // 1. 优先检查 Cookie
    const stored = getCookie(STORAGE_KEY);
    if (stored && SUPPORTED_LANGUAGES.includes(stored)) {
      return stored;
    }

    // 2. 检查浏览器语言
    const browserLang = (
      navigator.language ||
      navigator.userLanguage ||
      DEFAULT_LANGUAGE
    ).toLowerCase();

    // 处理中文变体
    if (browserLang.startsWith('zh')) {
      return browserLang.includes('tw') || browserLang.includes('hk')
        ? 'zh-tw'
        : 'zh-cn';
    }

    // 处理其他语言
    const langCode = browserLang.split('-')[0];
    return SUPPORTED_LANGUAGES.includes(langCode) ? langCode : DEFAULT_LANGUAGE;
  },

  set(lang) {
    if (!SUPPORTED_LANGUAGES.includes(lang)) {
      return;
    }

    // Prevent unnecessary updates/reloads if language is unchanged
    if (this.current === lang) return;

    this.current = lang;
    this.locale = lang; // Sync with Uno
    setCookie(STORAGE_KEY, lang);
    this.apply(lang);

    // 触发自定义事件，以便其他组件响应
    document.dispatchEvent(
      new CustomEvent('languagechange', {
        detail: { language: lang },
      }),
    );

    // Reload the page to apply server-side translations (Go templates)
    window.location.reload();
  },

  apply(lang) {
    document.documentElement.lang = lang;
  },
};
