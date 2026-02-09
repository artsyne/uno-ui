import themeStore from '../stores/theme.js';
import i18nStore from '../stores/i18n.js';
import { Navbar, Footer } from '../modules/layout.js';
import DiaryDetail from '../modules/diary-detail.js';
import { LikeButton, CommentsMixin } from '../modules/social.js';
import { initAuth, initValidationCallbacks } from '../modules/auth.js';
import { initFormHandlers, validateAndTrigger } from '../modules/form.js';
import { User } from '../modules/user.js';
import { Settings, HandleSettings } from '../modules/settings.js';
import { DiaryForm } from '../modules/diary-form.js';
import { HeroAnimation } from '../modules/home.js';
import Drafts from '../modules/drafts.js';
import {
  InviteAuthorModal,
  AuthorList,
  initAuthorListeners,
} from '../modules/authors.js';

/**
 * Application Initialization
 * Relies on script loading order in footer.html:
 * 1. htmx
 * 2. uno.umd.js (Window.Uno)
 * 3. bundle.js (This file)
 * 4. alpine.js (Fires alpine:init)
 */

// 1. Uno UI Configuration
window.Uno.setConfig({
  debug: false,
  notify: {
    position: 'bottom-right',
  },
  ui: {
    base: 'shadow-md border border-gray-200 dark:border-gray-700',
    box: {
      container:
        'relative transform overflow-hidden shadow-md border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-900 text-left shadow-xl transition-all sm:my-8 sm:w-full sm:max-w-lg flex flex-col max-h-[90vh] pointer-events-auto outline-none',
      cancelBtn:
        'mt-3 inline-flex w-full justify-center bg-white px-3 py-2 text-sm font-semibold text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 hover:bg-gray-50 sm:mt-0 sm:w-auto',
      okBtn: 'btn ml-4',
      promptInput:
        'w-full px-3 py-2 border focus:ring-2 focus:ring-blue-500 outline-none',
    },
    notify: {
      item: 'pointer-events-auto max-w-[calc(100vw-2rem)] overflow-hidden shadow-md border border-gray-200 dark:border-gray-700',
      progressWrap:
        'fixed top-0 left-0 right-0 h-0.5 bg-transparent pointer-events-none z-[9999]',
      progressBar:
        'h-full bg-slate-400 dark:bg-slate-500 transition-all duration-200 ease-out',
    },
    tooltip: {
      content:
        'fixed px-3 py-2 text-sm font-medium shadow-md border border-gray-200 dark:border-gray-700 whitespace-nowrap bg-white text-gray-900 dark:bg-slate-800 dark:text-white',
    },
  },
});

// 2. Global Helper Functions
window.getUno = () =>
  window.Alpine && typeof window.Alpine.store === 'function'
    ? window.Alpine.store('uno')
    : null;
window.getUnoNotify = () => window.getUno()?.notify;
window.getUnoBox = () => window.getUno()?.box;
window.getUnoValid = () => window.getUno()?.valid;
window.getUnoTooltip = () => window.getUno()?.tooltip;

// 2.0.1 Form Validation Helper
window.validateAndTrigger = validateAndTrigger;

// 2.0.2 Helper to get CSRF Token
window.getCsrfToken = () => {
  return document
    .querySelector('meta[name="csrf-token"]')
    ?.getAttribute('content');
};

// 2.1 Fetch Interceptor for Global Progress Bar
const originalFetch = window.fetch;
window.fetch = function (input, init) {
  const notify = window.getUnoNotify();
  const shouldNotify = notify && (!init || !init.skipProgress);

  // Inject CSRF Token
  const csrfToken = window.getCsrfToken();
  if (csrfToken) {
    init = init || {};
    init.headers = init.headers || {};
    if (init.headers instanceof Headers) {
      init.headers.append('X-CSRF-Token', csrfToken);
    } else {
      init.headers['X-CSRF-Token'] = csrfToken;
    }
  }

  if (shouldNotify) {
    notify.progressRequests++;
    if (notify.progressRequests === 1) notify.startProgress();
  }

  return originalFetch(input, init).finally(() => {
    if (shouldNotify) {
      if (notify.progressRequests > 0) notify.progressRequests--;
      if (notify.progressRequests === 0) notify.stopProgress();
    }
  });
};

// 3. Initialize Auth (HTMX Config)
initAuth();
initFormHandlers();
initAuthorListeners();

// 4. Global Event Listeners
document.body.addEventListener('htmx:configRequest', (evt) => {
  const csrfToken = window.getCsrfToken();
  if (csrfToken) {
    evt.detail.headers['X-CSRF-Token'] = csrfToken;
  }
});

document.body.addEventListener('htmx:beforeRequest', () => {
  const notify = window.getUnoNotify();
  if (notify) {
    notify.progressRequests++;
    if (notify.progressRequests === 1) notify.startProgress();
  }
});

document.body.addEventListener('htmx:afterRequest', () => {
  const notify = window.getUnoNotify();
  if (notify) {
    if (notify.progressRequests > 0) notify.progressRequests--;
    if (notify.progressRequests === 0) notify.stopProgress();
  }
});

document.body.addEventListener('showToast', (evt) => {
  const { type, message } = evt.detail;
  const notify = window.getUnoNotify();
  if (notify) {
    if (type === 'error') notify.error(message);
    else notify.success(message);
  }
});

document.body.addEventListener('htmx:afterSwap', (evt) => {
  const target = evt?.detail?.target;
  if (target) {
    window.Alpine.initTree(target);
  }
});

// 5. Alpine Initialization
document.addEventListener('alpine:init', () => {
  // Register Components
  Alpine.data('navbar', Navbar);
  Alpine.data('footer', Footer);
  Alpine.data('diaryDetail', DiaryDetail);
  Alpine.data('likeButton', LikeButton);
  Alpine.data('CommentsMixin', CommentsMixin);
  Alpine.data('user', () => User);
  Alpine.data('settings', Settings);
  Alpine.data('handleSettings', HandleSettings);
  Alpine.data('diaryForm', DiaryForm);
  Alpine.data('heroAnimation', HeroAnimation);
  Alpine.data('drafts', Drafts);
  Alpine.data('authorList', AuthorList);

  // Register Stores
  Alpine.store('theme', themeStore);

  // Merge i18n Store (Preserve Uno messages if present)
  const existingI18n = Alpine.store('i18n');
  if (existingI18n) {
    Object.assign(existingI18n, {
      ...i18nStore,
      messages: existingI18n.messages,
      init() {
        if (i18nStore.init) i18nStore.init.call(this);
      },
    });
    // Manually call init() because Alpine won't do it again for an existing store
    if (existingI18n.init) {
      existingI18n.init();
    }
  } else {
    Alpine.store('i18n', i18nStore);
  }
});

// 6. Post-Initialization Setup
document.addEventListener('alpine:initialized', () => {
  initValidationCallbacks();
});
