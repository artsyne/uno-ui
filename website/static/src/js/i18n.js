/**
 * Website 多语言初始化脚本
 * 基于 Kiko UI 的多语言机制，语言包已打包到 JS 文件中
 */

(function() {
  'use strict';

  // 支持的语言列表
  const SUPPORTED_LANGUAGES = ['en', 'zh-cn', 'zh-tw'];
  
  // 默认语言
  const DEFAULT_LANGUAGE = 'en';

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
    return SUPPORTED_LANGUAGES.includes(langCode) ? langCode : DEFAULT_LANGUAGE;
  }

  /**
   * 获取当前应该使用的语言
   */
  function getCurrentLanguage() {
    // 优先级：window.kikoLang > 浏览器语言 > 默认语言
    if (window.kikoLang && typeof window.kikoLang === 'string') {
      const lang = window.kikoLang.toLowerCase();
      if (SUPPORTED_LANGUAGES.includes(lang)) {
        return lang;
      }
    }
    
    return detectBrowserLanguage();
  }

  /**
   * 设置语言的便捷方法
   * @param {string} lang - 语言代码
   * @returns {boolean} 是否设置成功
   */
  function setLanguage(lang) {
    if (!SUPPORTED_LANGUAGES.includes(lang)) {
      console.warn(`[Website I18n] 不支持的语言: ${lang}`);
      return false;
    }
    
    // 使用 Kiko UI 设置语言
    const success = kiko.i18n.setLang(lang);
    
    if (success) {
      // 同步到全局变量
      window.kikoLang = lang;
      
      // 保存到本地存储
      localStorage.setItem('website-lang', lang);
      
      console.log(`[Website I18n] 语言已切换到: ${lang}`);
    }
    
    return success;
  }

  /**
   * 暴露全局 API
   */
  function exposeGlobalAPI() {
    // 暴露语言设置方法到全局
    window.setWebsiteLanguage = setLanguage;
    
    // 暴露支持的语言列表
    window.getSupportedLanguages = function() {
      return SUPPORTED_LANGUAGES.slice();
    };
    
    // 暴露当前语言获取方法
    window.getCurrentLanguage = function() {
      return kiko.i18n.getCurrentLang();
    };
  }

  /**
   * 初始化 Website 多语言系统
   */
  function initWebsiteI18n() {
    try {
      // 等待 Kiko UI 初始化完成
      if (typeof kiko === 'undefined' || !kiko.i18n) {
        console.warn('[Website I18n] 等待 Kiko UI 初始化...');
        
        // 等待 Kiko UI 可用
        const checkKiko = () => {
          if (typeof kiko !== 'undefined' && kiko.i18n) {
            // 暴露全局 API
            exposeGlobalAPI();
            
            // 设置当前语言（语言包已通过构建脚本打包）
            const currentLang = getCurrentLanguage();
            kiko.i18n.setLang(currentLang);
            
            console.log(`[Website I18n] 初始化完成，当前语言: ${currentLang}`);
            
            // 触发语言包加载完成事件
            document.dispatchEvent(new CustomEvent('websiteI18nReady', {
              detail: { language: currentLang }
            }));
          } else {
            setTimeout(checkKiko, 50);
          }
        };
        checkKiko();
        return;
      }

      // 暴露全局 API
      exposeGlobalAPI();
      
      // 设置当前语言（语言包已通过构建脚本打包）
      const currentLang = getCurrentLanguage();
      kiko.i18n.setLang(currentLang);
      
      console.log(`[Website I18n] 初始化完成，当前语言: ${currentLang}`);
      
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

  // 根据文档加载状态决定初始化时机
  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', initWebsiteI18n);
  } else {
    // 文档已加载完成，立即初始化
    initWebsiteI18n();
  }

})();