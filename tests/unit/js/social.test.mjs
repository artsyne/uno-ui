import { jest } from '@jest/globals';
import { getBatchStatsManager, entryStats } from '../../../website/static/src/js/modules/social.js';

describe('Social Module', () => {
  let batchStatsManager;
  let mockFetch;

  beforeEach(() => {
    // Reset mocks
    mockFetch = jest.fn(() =>
      Promise.resolve({
        json: () => Promise.resolve({ success: true, data: { stats: [], status: [] } }),
        ok: true,
        headers: {
          get: (name) => name === 'content-type' ? 'application/json' : null,
        },
      })
    );
    global.fetch = mockFetch;

    // Mock IntersectionObserver
    global.IntersectionObserver = class {
      constructor(cb) {
        this.cb = cb;
      }
      observe() {}
      unobserve() {}
      disconnect() {}
    };

    // Mock window events
    window.dispatchEvent = jest.fn();
    window.addEventListener = jest.fn();
    window.removeEventListener = jest.fn();

    batchStatsManager = getBatchStatsManager();
    // Reset manager state - 使用新的属性名
    batchStatsManager.pendingIds.clear();
    batchStatsManager.loadedIds.clear();
    batchStatsManager.cache = {};
    if (batchStatsManager.timer) clearTimeout(batchStatsManager.timer);
    batchStatsManager.timer = null;
    batchStatsManager.observer = null;
    batchStatsManager.observed = new Map();
  });

  describe('batchStatsManager', () => {
    test('should batch requests', async () => {
      jest.useFakeTimers();

      batchStatsManager.add(1);
      batchStatsManager.add(2);

      expect(batchStatsManager.pendingIds.size).toBe(2);

      // Fast-forward time
      jest.advanceTimersByTime(100);

      expect(mockFetch).toHaveBeenCalledTimes(2); // stats + like/status
      const body = JSON.parse(mockFetch.mock.calls[0][1].body);
      expect(body.entryids).toContain('1');
      expect(body.entryids).toContain('2');
      
      jest.useRealTimers();
    });

    test('should use shared IntersectionObserver', async () => {
      const el1 = document.createElement('div');
      const el2 = document.createElement('div');
      
      let observerCallback;
      global.IntersectionObserver = class {
        constructor(cb) {
          observerCallback = cb;
        }
        observe() {}
        unobserve() {}
      };

      batchStatsManager.observe(el1, 1);
      batchStatsManager.observe(el2, 2);

      expect(batchStatsManager.observed.size).toBe(2);

      // Simulate intersection
      jest.useFakeTimers();
      observerCallback([
        { isIntersecting: true, target: el1 },
        { isIntersecting: true, target: el2 }
      ]);

      // Fast-forward time for debounce
      jest.advanceTimersByTime(100);

      expect(mockFetch).toHaveBeenCalledTimes(2);
      const body = JSON.parse(mockFetch.mock.calls[0][1].body);
      expect(body.entryids.split(',')).toEqual(expect.arrayContaining(['1', '2']));
      
      jest.useRealTimers();
    });

    test('should cache stats data', () => {
      batchStatsManager.set(1, { view_num: 10, like_num: 5 });
      expect(batchStatsManager.get(1)).toEqual({ view_num: 10, like_num: 5 });
      
      batchStatsManager.set(1, { liked: true });
      expect(batchStatsManager.get(1)).toEqual({ view_num: 10, like_num: 5, liked: true });
    });

    test('should clear cache', () => {
      batchStatsManager.set(1, { view_num: 10 });
      batchStatsManager.loadedIds.add('1');
      
      batchStatsManager.clear(1);
      
      expect(batchStatsManager.get(1)).toBeUndefined();
      expect(batchStatsManager.loadedIds.has('1')).toBe(false);
    });
  });

  describe('entryStats', () => {
    test('toggleLike should update cache', async () => {
      const stats = entryStats(1, 1);
      stats.stats.loaded = true;
      stats.stats.loading = false;
      stats.stats.liked = false;
      stats.stats.like_num = 0;

      // Mock fetch response for toggleLike
      mockFetch.mockResolvedValueOnce({
        json: () => Promise.resolve({ 
          success: true, 
          data: { is_liked: true, like_count: 1 } 
        }),
      });

      await stats.toggleLike();

      expect(stats.stats.liked).toBe(true);
      expect(batchStatsManager.cache['1'].liked).toBe(true);
      expect(batchStatsManager.cache['1'].like_num).toBe(1);
    });

    test('toggleLike should rollback on error', async () => {
      const stats = entryStats(1, 1);
      stats.stats.loaded = true;
      stats.stats.loading = false;
      stats.stats.liked = false;
      stats.stats.like_num = 5;

      // Mock fetch failure
      mockFetch.mockRejectedValueOnce(new Error('Network error'));

      await stats.toggleLike();

      // Should rollback to original values
      expect(stats.stats.liked).toBe(false);
      expect(stats.stats.like_num).toBe(5);
    });
  });
});