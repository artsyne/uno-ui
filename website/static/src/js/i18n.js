/**
 * Website 多语言初始化脚本
 * 基于 Kiko UI 的多语言机制，语言包已打包到 JS 文件中
 */

(function() {
  'use strict';

  // 存储键名（前后端约定）
  const STORAGE_KEY = 'idiary-lang';
  
  // 默认语言
  const DEFAULT_LANGUAGE = 'en';

  // 支持的语言代码（从navbar.html中的hardcode列表获取）
  const SUPPORTED_LANGUAGE_CODES = ['ar', 'de', 'en', 'es', 'fr', 'ja', 'ko', 'ru', 'zh-cn', 'zh-tw'];

  /**
   * 检测浏览器语言
   */
  function detectBrowserLanguage() {
    const browserLang = (navigator.language || navigator.userLanguage || DEFAULT_LANGUAGE).toLowerCase();
    
    // 处理中文语言变体
    if (browserLang.startsWith('zh')) {
      return browserLang.includes('tw') || browserLang.includes('hk') ? 'zh-tw' : 'zh-cn';
    }
    
    // 处理其他语言
    const langCode = browserLang.split('-')[0];
    return SUPPORTED_LANGUAGE_CODES.includes(langCode) ? langCode : DEFAULT_LANGUAGE;
  }
  
  /**
   * 获取存储的语言
   */
  function getStoredLanguage() {
    try {
      const stored = window.getCookie(STORAGE_KEY);
      return stored && SUPPORTED_LANGUAGE_CODES.includes(stored) ? stored : null;
    } catch (e) {
      return null;
    }
  }

  /**
   * 获取当前应该使用的语言
   * 简化优先级：cookie > 浏览器默认 > 默认en
   */
  function getCurrentLanguage() {
    // 1. cookie 优先级最高
    const storedLang = getStoredLanguage();
    if (storedLang) {
      return storedLang;
    }
    
    // 2. 浏览器默认语言
    const browserLang = detectBrowserLanguage();
    if (browserLang && browserLang !== DEFAULT_LANGUAGE) {
      return browserLang;
    }
    
    // 3. 系统默认
    return DEFAULT_LANGUAGE;
  }

  /**
   * 设置语言的便捷方法
   * @param {string} lang - 语言代码
   * @returns {boolean} 是否设置成功
   */
  function setLanguage(lang) {
    if (!SUPPORTED_LANGUAGE_CODES.includes(lang)) {
      console.warn(`[Website I18n] 不支持的语言: ${lang}`);
      return false;
    }
    
    try {
      window.setCookie(STORAGE_KEY, lang);
      
      // 使用 Kiko UI 设置语言
      if (typeof kiko !== 'undefined' && kiko.i18n) {
        const success = kiko.i18n.setLang(lang);
        
        if (success) {
          // 同步到全局变量
          window.kikoLang = lang;
          
          // 更新HTML lang属性
          document.documentElement.lang = lang;
          
          console.log(`[Website I18n] 语言已切换到: ${lang}`);
          
          // 触发语言变化事件
          document.dispatchEvent(new CustomEvent('languagechange', {
            detail: { language: lang }
          }));
        }
        
        return success;
      }
      
      return true;
    } catch (e) {
      console.error('[Website I18n] 设置语言失败:', e);
      return false;
    }
  }

  /**
   * 暴露全局 API
   */
  function exposeGlobalAPI() {
    // 暴露语言设置方法到全局
    window.setWebsiteLanguage = setLanguage;
    
    // 暴露支持的语言代码列表
    window.getSupportedLanguageCodes = function() {
      return SUPPORTED_LANGUAGE_CODES.slice();
    };
    
    // 暴露当前语言获取方法
    window.getCurrentWebsiteLanguage = function() {
      return getCurrentLanguage();
    };
    
    // 兼容旧API
    window.getCurrentLanguage = function() {
      return kiko.i18n ? kiko.i18n.getCurrentLang() : getCurrentLanguage();
    };
  }

  /**
   * 初始化 Website 多语言系统
   * 页面初始化时设置 idiary-lang cookie 和 html 的 lang 属性
   */
  function initWebsiteI18n() {
    try {
      // 暴露全局 API
      exposeGlobalAPI();
      
      // 获取当前语言（语言包已通过构建脚本打包）
      const currentLang = getCurrentLanguage();
      
      // 确保 cookie 和 HTML lang 属性同步
      window.setCookie(STORAGE_KEY, currentLang);
      document.documentElement.lang = currentLang;
      
      // 立即设置 kiko 语言
      kiko.i18n.setLang(currentLang);

      // 触发语言包加载完成事件
      document.dispatchEvent(new CustomEvent('websiteI18nReady', {
        detail: { language: currentLang }
      }));
      
    } catch (error) {
      console.error('[Website I18n] 初始化失败:', error);
      
      // 即使失败也要暴露 API
      exposeGlobalAPI();
    }
  }

  /**
   * 绑定语言切换事件
   */
  function bindLanguageEvents() {
    // 语言切换tooltip中的链接点击事件
    document.addEventListener('click', (e) => {
      const langLink = e.target.closest('[data-lang]');
      if (langLink) {
        e.preventDefault();
        const langCode = langLink.getAttribute('data-lang');
        
        // 设置语言（现在使用cookie存储）
        if (setLanguage(langCode)) {
          // 隐藏tooltip
          const tooltip = document.getElementById('lang_menu');
          if (tooltip) {
            tooltip.style.display = 'none';
          }
          
          // 简化刷新：直接重新加载当前页面
          window.location.reload();
        }
      }
    });
  }

  // 根据文档加载状态决定初始化时机
  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', initWebsiteI18n);
  } else {
    // 文档已加载完成，立即初始化
    initWebsiteI18n();
  }
  
  // 绑定事件（DOM加载完成后）
  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', bindLanguageEvents);
  } else {
    bindLanguageEvents();
  }

})();