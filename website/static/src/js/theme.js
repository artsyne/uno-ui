/**
 * 极简主题管理器
 * 基于 Tailwind CSS 的 dark 类机制
 * 
 * @version 2.0.1
 * @author iDiary.ai
 */

// 立即执行的主题预加载 - 防止页面闪烁
// 必须在任何其他代码执行前运行
(function() {
  try {
    const theme = localStorage.getItem('idiary-theme');
    const isDark = theme === 'dark' || (!theme && window.matchMedia('(prefers-color-scheme: dark)').matches);
    if (isDark) document.documentElement.classList.add('dark');
  } catch (e) {
    document.documentElement.classList.add('dark'); // 默认暗黑模式
  }
})();

(function() {
  'use strict';

  // 主题管理器
  const ThemeManager = {
    storageKey: 'idiary-theme',
    
    /**
     * 初始化主题系统
     */
    init() {
      this.applyStoredTheme();
      this.bindEvents();
    },
    
    /**
     * 应用存储的主题
     */
    applyStoredTheme() {
      let savedTheme = this.getStoredTheme();
      let effectiveTheme;
      if (savedTheme === null) {
        effectiveTheme = window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
      } else {
        effectiveTheme = savedTheme;
      }
      this.applyTheme(effectiveTheme);
    },
    
    /**
     * 获取存储的主题
     */
    getStoredTheme() {
      try {
        return localStorage.getItem(this.storageKey);
      } catch (e) {
        return 'system';
      }
    },
    
    /**
     * 解析主题
     */
    resolveTheme(theme) {
      return theme;
    },
    
    /**
     * 应用主题到 HTML
     */
    applyTheme(theme) {
      const html = document.documentElement;
      
      if (theme === 'dark') {
        html.classList.add('dark');
      } else {
        html.classList.remove('dark');
      }
      
      // 更新 meta theme-color
      this.updateMetaThemeColor(theme);
      
      // 更新主题切换按钮图标
      this.updateThemeIcon();
    },
    
    /**
     * 更新 meta theme-color
     */
    updateMetaThemeColor(theme) {
      let metaThemeColor = document.querySelector('meta[name="theme-color"]');
      
      if (!metaThemeColor) {
        metaThemeColor = document.createElement('meta');
        metaThemeColor.name = 'theme-color';
        document.head.appendChild(metaThemeColor);
      }
      
      metaThemeColor.content = theme === 'dark' ? '#1f2937' : '#ffffff';
    },
    
    /**
     * 更新主题切换按钮图标
     */
    updateThemeIcon() {
      const themeButton = document.querySelector('[data-theme-toggle]');
      if (!themeButton) return;
      
      const icon = themeButton.querySelector('i');
      if (!icon) return;
      
      const currentTheme = this.getCurrentTheme();
      
      // 移除所有可能的图标类
      icon.className = icon.className.replace(/mdi-[\w-]+/g, '');
      
      // 根据当前主题设置对应图标
      switch (currentTheme) {
        case 'light':
          icon.classList.add('mdi-weather-sunny');
          themeButton.title = '切换到暗色主题';
          break;
        case 'dark':
          icon.classList.add('mdi-weather-night');
          themeButton.title = '切换到亮主题';
          break;
        default:
          icon.classList.add('mdi-weather-sunny');
          themeButton.title = '切换到亮色主题';
          break;
      }
    },
    
    /**
     * 切换主题
     */
    toggle() {
      const currentStored = this.getStoredTheme();
      const nextTheme = this.getNextTheme(currentStored);
      
      this.setTheme(nextTheme);
    },
    
    /**
     * 获取下一个主题
     */
    getNextTheme(current) {
      const themes = ['light', 'dark'];
      const currentIndex = themes.indexOf(current);
      return themes[(currentIndex + 1) % themes.length];
    },
    
    /**
     * 设置主题
     */
    setTheme(theme) {
      try {
        localStorage.setItem(this.storageKey, theme);
      } catch (e) {
        console.warn('无法保存主题设置:', e);
      }
      
      const effectiveTheme = this.resolveTheme(theme);
      this.applyTheme(effectiveTheme);
      
      // 触发自定义事件
      this.dispatchThemeChangeEvent(theme, effectiveTheme);
    },
    
    /**
     * 获取当前主题
     */
    getCurrentTheme() {
      return this.getStoredTheme();
    },
    
    /**
     * 获取当前有效主题
     */
    getEffectiveTheme() {
      return this.resolveTheme(this.getCurrentTheme());
    },
    
    /**
     * 绑定事件
     */
    bindEvents() {
      // 主题切换按钮点击事件
      document.addEventListener('click', (e) => {
        if (e.target.closest('[data-theme-toggle]')) {
          e.preventDefault();
          this.toggle();
        }
      });
      
      // 键盘快捷键 (Ctrl/Cmd + Shift + T)
      document.addEventListener('keydown', (e) => {
        if ((e.ctrlKey || e.metaKey) && e.shiftKey && e.key === 'T') {
          e.preventDefault();
          this.toggle();
        }
      });
      

    },
    
    /**
     * 触发主题变化事件
     */
    dispatchThemeChangeEvent(theme, effectiveTheme) {
      const event = new CustomEvent('themechange', {
        detail: {
          theme: theme,
          effectiveTheme: effectiveTheme
        }
      });
      
      document.dispatchEvent(event);
    }
  };
  
  // 暴露到全局
  window.ThemeManager = ThemeManager;
  
  // 立即初始化（在 DOM 加载前就执行，防止闪烁）
  ThemeManager.init();
  
  // DOM 加载完成后再次确保主题正确
  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', () => {
      ThemeManager.applyStoredTheme();
    });
  }
  
})();