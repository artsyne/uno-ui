import { getCookie, setCookie } from '../utils/cookie.js';

const STORAGE_KEY = 'idiary-lang';
const DEFAULT_LANGUAGE = 'en';
const SUPPORTED_LANGUAGES = [
  'ar',
  'de',
  'en',
  'es',
  'fr',
  'ja',
  'ko',
  'ru',
  'zh-cn',
  'zh-tw',
];

export default {
  current: DEFAULT_LANGUAGE,
  locale: DEFAULT_LANGUAGE, // Sync with Uno

  init() {
    this.current = this.detectLanguage();
    this.locale = this.current; // Sync with Uno
    this.apply(this.current);
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
      console.warn(`[I18n Store] Unsupported language: ${lang}`);
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

    // 兼容性支持：如果 Kiko 存在，同步语言设置
    // 这是 Strangler Fig 模式的一部分，确保旧组件也能感知语言变化
    if (typeof window.kiko !== 'undefined' && window.kiko.i18n) {
      window.kiko.i18n.setLang(lang);
    }
  },
};
