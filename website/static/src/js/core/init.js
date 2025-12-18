import themeStore from '../stores/theme.js';
import i18nStore from '../stores/i18n.js';
import { Navbar, Footer } from '../modules/layout.js';
import { DiaryList } from '../modules/diary.js';
import { LikeButton, CommentsMixin } from '../modules/social.js';
import { initAuth, initValidationCallbacks, ValidForm } from '../modules/auth.js';
import { Member } from '../modules/member.js';
import { Settings } from '../modules/settings.js';
import { Notify } from '../utils/uno-shim.js';

// Initialize Auth module (HTMX config)
initAuth();

// Listen for HTMX triggered toasts
document.body.addEventListener('showToast', (evt) => {
  const { type, message } = evt.detail;
  if (type === 'success') Notify.success(message);
  else if (type === 'error') Notify.error(message);
  else Notify.success(message); // Default to success
});

// 使用 CDN 模式，Alpine 和 HTMX 通过 script 标签引入，存在于全局作用域
// 我们需要在 alpine:init 事件中注册 Stores

document.addEventListener('alpine:init', () => {
  // Initialize validation callbacks
  initValidationCallbacks();

  // 注册全局组件
  Alpine.data('navbar', Navbar);
  Alpine.data('footer', Footer);
  Alpine.data('diaryList', DiaryList);
  Alpine.data('likeButton', LikeButton);
  Alpine.data('CommentsMixin', CommentsMixin);
  Alpine.data('member', () => Member);
  Alpine.data('settings', Settings);
  Alpine.data('validForm', ValidForm);

  // 注册 Alpine Stores
  Alpine.store('theme', themeStore);

  // Safe registration for i18n store (merge with Uno if exists)
  const existingI18n = Alpine.store('i18n');
  if (existingI18n) {
    console.log('[Init] Merging custom i18n store into Uno store');
    if (existingI18n.messages) {
      console.log(
        '[Init] Found existing Uno messages:',
        Object.keys(existingI18n.messages),
      );
    } else {
      console.warn('[Init] Uno store exists but has no messages!');
    }

    // Merge properties and methods
    Object.assign(existingI18n, {
      ...i18nStore,
      // Preserve Uno's messages if we don't have them (we don't)
      messages: existingI18n.messages,
      // Merge init
      init() {
        if (i18nStore.init) i18nStore.init.call(this);
        // Uno's init is usually empty or just sets lang, which we do in our init too
      },
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
