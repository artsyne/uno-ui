import { jest } from '@jest/globals';
import { handleApiError } from '../../../website/static/src/js/utils/error-handler.js';

describe('handleApiError', () => {
  beforeEach(() => {
    // Reset window.Alpine before each test
    window.Alpine = undefined;
  });

  test('should return default error for empty response', () => {
    const result = handleApiError(null);
    expect(result.message).toBe('common.op_fail');
  });

  test('should use backend message as key', () => {
    const response = {
      success: false,
      code: 1204,
      message: 'auth.login_failed',
      data: {},
      timestamp: 123456789,
      request_id: 'req-1',
    };
    const result = handleApiError(response);
    expect(result.localizedMessage).toBe('auth.login_failed'); // No i18n store yet
  });

  test('should interpolate data without i18n store', () => {
    const response = {
      message: 'User {user} not found',
      data: { user: 'Alice' },
    };
    const result = handleApiError(response);
    expect(result.localizedMessage).toBe('User Alice not found');
  });

  test('should use i18n store if available', () => {
    // Mock Alpine and i18n store
    const tMock = jest.fn((key, params) => {
      if (key === 'auth.user_error') return `Error for ${params.user}`;
      return key;
    });

    window.Alpine = {
      store: jest.fn((name) => {
        if (name === 'i18n') return { t: tMock };
        return null;
      }),
    };

    const response = {
      message: 'auth.user_error',
      data: { user: 'Bob' },
    };

    const result = handleApiError(response);
    expect(window.Alpine.store).toHaveBeenCalledWith('i18n');
    expect(tMock).toHaveBeenCalledWith('auth.user_error', { user: 'Bob' });
    expect(result.localizedMessage).toBe('Error for Bob');
  });
});
