import { strict as assert } from 'assert';
import { initValidationCallbacks } from '../../website/static/src/js/modules/auth.js';

// Mock Alpine
global.Alpine = {
  store: (name) => {
    if (name === 'uno') return { valid: mockValidStore };
    if (name === 'valid') return mockValidStore;
    if (name === 'i18n') return { t: (key) => `TRANS:${key}` };
    return {};
  },
};

// Mock Document
global.document = {
  querySelector: () => null,
  getElementById: (id) => {
    if (id === 'new_pwd') return { value: 'password123' };
    return null;
  },
};

// Mock FormData
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

// Mock Fetch
global.fetch = async (url, options) => {
  if (url === '/check-username') {
    const user = options.body.get('user_name');
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
};

const mockValidStore = {
  validators: {},
  registerValidator(name, fn) {
    this.validators[name] = fn;
  },
  showError() {
    throw new Error('Should not call showError in new implementation');
  },
};

// Run Test
async function runTests() {
  console.log('Running Auth Validation Tests...');
  initValidationCallbacks();

  // Test checkUsername (Async)
  const checkUsername = mockValidStore.validators['checkUsername'];
  assert(checkUsername, 'checkUsername validator should be registered');

  // Success case
  const userResult = await checkUsername('testuser', 'user_name');
  assert.equal(
    userResult,
    true,
    'checkUsername should return true for valid user',
  );

  // Fail case (user_exists)
  const userTakenResult = await checkUsername('taken_user', 'user_name');
  assert.equal(
    userTakenResult,
    'TRANS:auth.username_taken',
    'checkUsername should return translated auth.username_taken when backend returns user_exists',
  );

  // Test checkEmail (Async)
  const checkEmail = mockValidStore.validators['checkEmail'];
  assert(checkEmail, 'checkEmail validator should be registered');

  // Fail case (email_exists)
  const emailResult = await checkEmail('taken@example.com', 'email');
  assert.equal(
    emailResult,
    'TRANS:auth.email_taken',
    'checkEmail should return translated auth.email_taken when backend returns email_exists',
  );

  // Test matchPwd (Sync)
  const matchPwd = mockValidStore.validators['matchPwd'];
  assert(matchPwd, 'matchPwd validator should be registered');

  // Rule with explicit target (simulating Uno passing rule object)
  // Mock document.querySelector for this specific test
  const originalQuerySelector = global.document.querySelector;
  global.document.querySelector = (sel) => {
    if (sel === '#other_pwd') return { value: 'password123' };
    return null;
  };

  const pwdResult = matchPwd('password123', 'confirm', {
    target: '#other_pwd',
  });
  assert.equal(pwdResult, true, 'matchPwd should match with explicit target');

  const pwdFail = matchPwd('wrong', 'confirm', { target: '#other_pwd' });
  assert.equal(pwdFail, false, 'matchPwd should fail mismatch');

  global.document.querySelector = originalQuerySelector;

  console.log('All tests passed!');
}

runTests().catch((e) => {
  console.error(e);
  process.exit(1);
});
