import themeStore from '../stores/theme.js';
import i18nStore from '../stores/i18n.js';
import { Navbar, Footer } from '../modules/layout.js';

// 使用 CDN 模式，Alpine 和 HTMX 通过 script 标签引入，存在于全局作用域
// 我们需要在 alpine:init 事件中注册 Stores

document.addEventListener('alpine:init', () => {
  // 注册全局组件
  Alpine.data('navbar', Navbar);
  Alpine.data('footer', Footer);

  // 注册 Alpine Stores
  Alpine.store('theme', themeStore);
  
  // Safe registration for i18n store (merge with Uno if exists)
  const existingI18n = Alpine.store('i18n');
  if (existingI18n) {
      console.log('[Init] Merging custom i18n store into Uno store');
      // Merge properties and methods
      Object.assign(existingI18n, {
          ...i18nStore,
          // Preserve Uno's messages if we don't have them (we don't)
          messages: existingI18n.messages,
          // Merge init
          init() {
             if (i18nStore.init) i18nStore.init.call(this);
             // Uno's init is usually empty or just sets lang, which we do in our init too
          }
      });
      // Run init manually as Alpine won't re-run it for existing store
      existingI18n.init();
  } else {
      Alpine.store('i18n', i18nStore);
  }

  console.log('[Init] Alpine stores registered');
});

// HTMX 自动初始化，无需手动 start
console.log('[Init] Modular JS bundle loaded');
