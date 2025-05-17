/**
 * i18n.js - 多语言支持工具
 *
 * 使用外部语言文件，从 window.kikoI18nMessages 加载语言包。
 * 优先使用 window.kikoLang 作为当前语言。
 * 用于 static/js/ 应用层，和 kiko ui 解耦。
 *
 * @module appI18n
 * @example
 * // 获取本地化消息
 * const msg = appI18n.getMessage('operation.success');
 * // 设置当前语言
 * appI18n.setLang('zh-cn');
 */
(function (global) {
  /**
   * 支持的语言列表
   * @type {Array<string>}
   */
  const supportedLanguages = ['zh-cn', 'en', 'zh-tw'];

  /**
   * 默认语言
   * @type {string}
   */
  const defaultLanguage = 'en';

  /**
   * 检测当前语言，优先使用 window.kikoLang
   * 若未设置，则使用浏览器语言，找不到则回退到默认语言
   * @returns {string} 当前语言代码
   */
  function detectLang() {
    // 优先使用全局设置的语言
    if (window.kikoLang && typeof window.kikoLang === 'string') {
      const lang = window.kikoLang.toLowerCase();
      if (isLanguageSupported(lang)) {
        return lang;
      }
    }

    // 使用浏览器语言
    const browserLang = (
      navigator.language ||
      navigator.userLanguage ||
      defaultLanguage
    ).toLowerCase();

    // 处理中文特殊情况
    if (browserLang.startsWith('zh')) {
      if (browserLang.includes('tw') || browserLang.includes('hk')) {
        return 'zh-tw';
      }
      return 'zh-cn';
    }

    // 取主语言代码
    const mainLang = browserLang.split('-')[0];
    return isLanguageSupported(mainLang) ? mainLang : defaultLanguage;
  }

  /**
   * 检查语言是否被支持
   * @param {string} lang - 语言代码
   * @returns {boolean} 是否支持该语言
   */
  function isLanguageSupported(lang) {
    return supportedLanguages.includes(lang);
  }

  /**
   * 获取语言包
   * @param {string} lang - 语言代码
   * @returns {Object} 语言包对象
   */
  function getLanguagePack(lang) {
    return (window.kikoI18nMessages && window.kikoI18nMessages[lang]) || {};
  }

  /**
   * 当前语言代码
   * @type {string}
   */
  let currentLang = detectLang();

  /**
   * 获取多语言消息
   * @param {string} key 消息键
   * @param {Object} [params] 可选参数，支持 {name} 占位符替换
   * @returns {string} 本地化消息
   */
  function getMessage(key, params) {
    const messages = getLanguagePack(currentLang);
    let msg = messages[key];

    // 回退到英文
    if (!msg && currentLang !== defaultLanguage) {
      const defaultMessages = getLanguagePack(defaultLanguage);
      msg = defaultMessages[key];
    }

    // 如果仍然没有找到，返回键名
    if (!msg) {
      msg = key;
    }

    // 参数替换
    if (params && typeof msg === 'string') {
      return msg.replace(/\{([^}]+)\}/g, function (m, k) {
        return params[k] !== undefined ? params[k] : m;
      });
    }

    return msg;
  }

  /**
   * 设置当前语言
   * @param {string} lang 语言代码
   * @returns {boolean} 是否设置成功
   */
  function setLang(lang) {
    if (isLanguageSupported(lang)) {
      currentLang = lang;
      // 同时设置全局语言变量，便于其他模块使用
      window.kikoLang = lang;
      return true;
    }
    return false;
  }

  /**
   * 获取当前语言
   * @returns {string} 当前语言代码
   */
  function getLang() {
    return currentLang;
  }

  /**
   * 获取支持的所有语言
   * @returns {Array<string>} 支持的语言代码列表
   */
  function getSupportedLanguages() {
    return supportedLanguages.slice();
  }

  /**
   * 全局导出 appI18n 对象
   * @namespace appI18n
   * @property {function} getMessage 获取本地化消息
   * @property {function} setLang 设置当前语言
   * @property {function} getLang 获取当前语言
   * @property {function} getSupportedLanguages 获取支持的所有语言
   */
  global.appI18n = {
    getMessage,
    setLang,
    getLang,
    getSupportedLanguages,
  };
})(window);
