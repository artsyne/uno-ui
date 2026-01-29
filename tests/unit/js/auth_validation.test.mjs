import { jest } from '@jest/globals';
import { initValidationCallbacks } from '../../../website/static/src/js/modules/auth.js';

// Mock objects
const mockValidStore = {
  validators: {},
  registerValidator(name, fn) {
    this.validators[name] = fn;
  },
  showError() {
    throw new Error('Should not call showError in new implementation');
  },
};

// Setup global mocks
beforeAll(() => {
  global.Alpine = {
    store: (name) => {
      if (name === 'uno') return { valid: mockValidStore };
      if (name === 'valid') return mockValidStore;
      if (name === 'i18n') return { t: (key) => `TRANS:${key}` };
      return {};
    },
  };

  global.document = {
    querySelector: () => null,
    getElementById: (id) => {
      if (id === 'new_pwd') return { value: 'password123' };
      return null;
    },
    body: {
      addEventListener: jest.fn(),
      removeEventListener: jest.fn(),
    },
  };

  global.FormData = class {
    constructor() {
      this.data = {};
    }
    append(key, value) {
      this.data[key] = value;
    }
    get(key) {
      return this.data[key];
    }
  };

  global.fetch = jest.fn(async (url, options) => {
    if (url === '/check-handle') {
      const user = options.body.get('handle');
      if (user === 'taken_user') {
        return {
          json: async () => ({ success: false, message: 'user_exists' }),
        };
      }
      return {
        json: async () => ({ success: true }),
      };
    }
    if (url === '/check-email') {
      const email = options.body.get('email');
      if (email === 'taken@example.com') {
        return {
          json: async () => ({ success: false, message: 'email_exists' }),
        };
      }
      return {
        json: async () => ({ success: true }),
      };
    }
  });
});

describe('Auth Validation', () => {
  beforeEach(() => {
    mockValidStore.validators = {};
    initValidationCallbacks();
  });

  test('registers validators', () => {
    expect(mockValidStore.validators['checkHandle']).toBeDefined();
    expect(mockValidStore.validators['checkEmail']).toBeDefined();
    expect(mockValidStore.validators['matchPwd']).toBeDefined();
  });

  test('checkHandle returns true for valid user', async () => {
    const checkHandle = mockValidStore.validators['checkHandle'];
    const result = await checkHandle('testuser', 'handle');
    expect(result).toBe(true);
  });

  test('checkHandle returns error message for taken user', async () => {
    const checkHandle = mockValidStore.validators['checkHandle'];
    const result = await checkHandle('taken_user', 'handle');
    expect(result).toBe('TRANS:auth.username_taken');
  });

  test('checkEmail returns true for valid email', async () => {
    const checkEmail = mockValidStore.validators['checkEmail'];
    // Assuming valid email logic in fetch mock
    const result = await checkEmail('valid@example.com', 'email');
    expect(result).toBe(true);
  });

  test('checkEmail returns error message for taken email', async () => {
    const checkEmail = mockValidStore.validators['checkEmail'];
    const result = await checkEmail('taken@example.com', 'email');
    expect(result).toBe('TRANS:auth.email_taken');
  });

  test('matchPwd validates password matching', () => {
    const matchPwd = mockValidStore.validators['matchPwd'];

    // Mock document.querySelector for this specific test
    const originalQuerySelector = global.document.querySelector;
    global.document.querySelector = jest.fn((sel) => {
      if (sel === '#other_pwd') return { value: 'password123' };
      return null;
    });

    const pwdResult = matchPwd('password123', 'confirm', {
      target: '#other_pwd',
    });
    expect(pwdResult).toBe(true);

    const pwdFail = matchPwd('wrong', 'confirm', { target: '#other_pwd' });
    expect(pwdFail).toBe(false);

    global.document.querySelector = originalQuerySelector;
  });
});
