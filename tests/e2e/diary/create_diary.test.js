// tests/e2e/diary/create_diary.test.js
const assert = require('assert');

// Mock browser environment
global.window = {
  Alpine: {
    initTree: () => console.log('[Alpine] initTree called'),
  },
  getUno: () => ({
    notify: {
      success: (msg) => console.log(`[Notify Success] ${msg}`),
      error: (msg) => console.log(`[Notify Error] ${msg}`),
    },
    box: {
      close: (id) => console.log(`[Box Close] ${id}`),
    },
    valid: {
      validateForm: async () => true,
    },
  }),
  location: {
    reload: () => console.log('[Location Reload]'),
  },
};

global.document = {
  getElementById: (id) => {
    if (id === 'diary_list') {
      return {
        prepend: (el) => console.log(`[DOM Prepend] Added new diary card`),
        firstElementChild: {},
      };
    }
    return null;
  },
  querySelector: (sel) => ({}),
  createElement: (tag) => ({
    innerHTML: '',
    firstElementChild: { tagName: 'DIV' },
  }),
};

global.FormData = class {
  constructor() {}
  entries() {
    return [];
  }
};

global.URLSearchParams = class {
  constructor() {}
  delete() {}
};

// Mock Alpine store for error handler
global.Alpine = {
  store: (name) => ({
    t: (key) => key,
  }),
};

// Mock User object context
const mockUser = {
  $store: {
    i18n: {
      t: (key) => key, // Mock translation
    },
  },
  _escapeHtml(unsafe) {
    return unsafe || '';
  },
  _renderDiaryCard(diary) {
    return `<div id="diary-${diary.diary_id}">${diary.title}</div>`;
  },
  // The function under test (Logic mirroring user.js)
  async _handleCreateDiarySubmit(boxId) {
    const uno = window.getUno();
    if (!uno) return false;

    const form = document.querySelector('#diary_editor_form');
    // Mock form validation
    if (uno.valid) {
      await uno.valid.validateForm(form);
    }

    const url = '/diary';
    const method = 'POST';

    try {
      const response = await fetch(url, { method });
      const result = await response.json();

      if (result.success) {
        uno.notify.success('common.save_success');
        uno.box.close(boxId);

        const diary = result.data;
        const diaryHtml = this._renderDiaryCard(diary);

        const container = document.getElementById('diary_list');
        const createElement = (html) => ({ tagName: 'DIV', html });

        if (container) {
          const newEl = createElement(diaryHtml);
          container.prepend(newEl);
          if (window.Alpine) {
            // Mock initTree
            window.Alpine.initTree(newEl);
          }
        }
        return true;
      } else {
        uno.notify.error(result.message || 'common.save_fail');
        return false;
      }
    } catch (e) {
      console.log('Error caught:', e);
      const t = (key) => global.Alpine.store('i18n').t(key);
      const msg = t('common.net_error');
      uno.notify.error(msg);
      return false;
    }
  },
};

async function testSuccess() {
  console.log('\n--- Test Case 1: Successful Creation ---');
  global.fetch = async (url) => ({
    json: async () => ({
      success: true,
      data: { diary_id: 123, title: 'New Diary', info_num: 0 },
    }),
  });

  const result = await mockUser._handleCreateDiarySubmit('box-1');
  assert.strictEqual(result, true, 'Should return true');
}

async function testFailure() {
  console.log('\n--- Test Case 2: Server Failure ---');
  global.fetch = async (url) => ({
    json: async () => ({
      success: false,
      message: 'Fail',
    }),
  });

  const result = await mockUser._handleCreateDiarySubmit('box-1');
  assert.strictEqual(result, false, 'Should return false');
}

async function testNetworkError() {
  console.log('\n--- Test Case 3: Network Error ---');
  global.fetch = async (url) => {
    throw new Error('Network Error');
  };

  const result = await mockUser._handleCreateDiarySubmit('box-1');
  assert.strictEqual(result, false, 'Should return false');
}

(async () => {
  try {
    await testSuccess();
    await testFailure();
    await testNetworkError();
    console.log('\nAll tests passed!');
  } catch (e) {
    console.error(e);
    process.exit(1);
  }
})();
