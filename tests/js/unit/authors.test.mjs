
import { jest } from '@jest/globals';
import { AuthorList } from '../../../website/static/src/js/modules/authors.js';

// Mock dependencies
jest.mock('../../../website/static/src/js/stores/i18n.js', () => ({
  t: jest.fn((key, params) => key),
}));

describe('AuthorList Module', () => {
  let authorList;
  let mockUno;
  let mockFetch;
  let mockReload;

  beforeEach(() => {
    // Reset mocks
    jest.clearAllMocks();

    // Mock Window methods
    mockReload = jest.fn();
    Object.defineProperty(window, 'location', {
      value: { reload: mockReload },
      writable: true,
    });

    window.getCsrfToken = jest.fn(() => 'mock-csrf-token');

    // Mock Uno
    mockUno = {
      box: {
        confirm: jest.fn().mockResolvedValue(true),
      },
      notify: {
        error: jest.fn(),
      },
    };
    window.getUno = jest.fn(() => mockUno);

    // Mock Fetch
    mockFetch = jest.fn().mockResolvedValue({
      ok: true,
      json: jest.fn().mockResolvedValue({}),
    });
    global.fetch = mockFetch;

    // Initialize component
    authorList = AuthorList();
  });

  describe('toggleAdmin', () => {
    test('should call confirm and fetch with correct params when granting admin', async () => {
      await authorList.toggleAdmin(1, 101, false);

      expect(mockUno.box.confirm).toHaveBeenCalled();
      expect(mockFetch).toHaveBeenCalledWith('/diary/1/authors/101', expect.objectContaining({
        method: 'PUT',
        headers: expect.objectContaining({
          'X-CSRF-Token': 'mock-csrf-token',
        }),
        body: JSON.stringify({ is_admin: 1 }),
      }));
      expect(mockReload).toHaveBeenCalled();
    });

    test('should call confirm and fetch with correct params when revoking admin', async () => {
      await authorList.toggleAdmin(1, 101, true);

      expect(mockUno.box.confirm).toHaveBeenCalled();
      expect(mockFetch).toHaveBeenCalledWith('/diary/1/authors/101', expect.objectContaining({
        method: 'PUT',
        body: JSON.stringify({ is_admin: 0 }),
      }));
      expect(mockReload).toHaveBeenCalled();
    });

    test('should not fetch if confirm is cancelled', async () => {
      mockUno.box.confirm.mockResolvedValue(false);
      await authorList.toggleAdmin(1, 101, false);
      expect(mockFetch).not.toHaveBeenCalled();
    });

    test('should notify error on fetch failure', async () => {
      mockFetch.mockResolvedValue({
        ok: false,
        json: jest.fn().mockResolvedValue({ message: 'Error occurred' }),
      });

      await authorList.toggleAdmin(1, 101, false);

      // Wait for promise chain
      await new Promise(process.nextTick);

      expect(mockUno.notify.error).toHaveBeenCalledWith('Error occurred');
      expect(mockReload).not.toHaveBeenCalled();
    });
  });

  describe('removeAuthor', () => {
    test('should call confirm and DELETE request', async () => {
      await authorList.removeAuthor(1, 102, 'User B');

      expect(mockUno.box.confirm).toHaveBeenCalled();
      expect(mockFetch).toHaveBeenCalledWith('/diary/1/authors/102', expect.objectContaining({
        method: 'DELETE',
        headers: expect.objectContaining({
          'X-CSRF-Token': 'mock-csrf-token',
        }),
      }));
      expect(mockReload).toHaveBeenCalled();
    });
  });

  describe('quitDiary', () => {
    test('should call confirm and DELETE request to self', async () => {
      await authorList.quitDiary(1);

      expect(mockUno.box.confirm).toHaveBeenCalled();
      expect(mockFetch).toHaveBeenCalledWith('/diary/1/authors/self', expect.objectContaining({
        method: 'DELETE',
      }));
      expect(mockReload).toHaveBeenCalled();
    });
  });
});
