import { jest } from '@jest/globals';
import DiaryDetail from '../../../website/static/src/js/modules/diary-detail.js';

describe('Diary Modules', () => {
  let mockBoxOpen;
  let mockBoxClose;
  let mockBoxConfirm;
  let mockBoxDialog;
  let mockNotifySuccess;
  let mockNotifyError;

  beforeAll(() => {
    // Mock Alpine store which uno-shim.js uses
    mockBoxOpen = jest.fn();
    mockBoxClose = jest.fn();
    mockBoxConfirm = jest.fn();
    mockBoxDialog = jest.fn();
    mockNotifySuccess = jest.fn();
    mockNotifyError = jest.fn();

    const mockAlpine = {
      store: jest.fn((name) => {
        if (name === 'uno') {
          return {
            box: {
              open: mockBoxOpen,
              close: mockBoxClose,
              confirm: mockBoxConfirm,
              dialog: mockBoxDialog,
            },
            notify: { success: mockNotifySuccess, error: mockNotifyError },
            valid: {
              validateForm: jest.fn().mockResolvedValue(true),
            },
          };
        }
        if (name === 'i18n') {
          return {
            t: (key) => key,
          };
        }
        return null;
      }),
      initTree: jest.fn(),
    };

    // Attach Alpine to the existing window object
    if (typeof window !== 'undefined') {
      window.Alpine = mockAlpine;
      window.getUno = () => window.Alpine.store('uno');
      window.getUnoNotify = () => window.Alpine.store('uno').notify;
    } else {
      global.window = {
        Alpine: mockAlpine,
        getUno: () => mockAlpine.store('uno'),
        getUnoNotify: () => mockAlpine.store('uno').notify,
      };
    }

    // Setup global mocks if not present (JSDOM usually provides them)
    if (!global.fetch) global.fetch = jest.fn();
  });

  beforeEach(() => {
    jest.clearAllMocks();
    mockBoxOpen.mockClear();
    mockBoxClose.mockClear();
    mockNotifySuccess.mockClear();
    mockNotifyError.mockClear();
  });

  const createDiaryDetail = () => {
    const component = DiaryDetail();
    component.$store = {
      i18n: {
        t: (key) => {
          const map = {
            'entry.delete.confirm': '确认删除?',
            'common.op_fail': '操作失败',
            'entry.edit': '编辑日记',
            'entry.create': '写日记',
            'common.save': '保存',
            'common.loading': '加载中...',
            'common.save_success': '保存成功',
            'social.like.success': '点赞成功',
            'social.like.cancel': '取消点赞',
          };
          return map[key] || key;
        },
      },
      tooltip: {
        hide: jest.fn(),
      },
    };
    return component;
  };

  describe('DiaryDetail', () => {
    test('should initialize and load stats', () => {
      const diaryDetail = createDiaryDetail();

      // Mock document.querySelector for main
      const mockMain = { dataset: { socialSetting: '1' } };
      const mockQuerySelector = jest.fn().mockReturnValue(mockMain);
      Object.defineProperty(document, 'querySelector', {
        value: mockQuerySelector,
        writable: true,
      });

      // Mock document.querySelectorAll
      const mockQuerySelectorAll = jest.fn();
      Object.defineProperty(document, 'querySelectorAll', {
        value: mockQuerySelectorAll,
        writable: true,
      });

      const mockElements = [
        { dataset: { entryId: '1' } },
        { dataset: { entryId: '2' } },
      ];
      mockQuerySelectorAll.mockReturnValue(mockElements);

      // Mock getElementById for loadDiaryLike which is also called
      jest
        .spyOn(document, 'getElementById')
        .mockReturnValue({ dataset: { diaryId: '1' } });

      // Mock fetch responses
      global.fetch.mockResolvedValue({
        ok: true,
        headers: {
          get: (name) => (name === 'content-type' ? 'application/json' : ''),
        },
        json: async () => ({
          code: 0,
          data: {
            stats: [
              { entry_id: '1', view_num: 10, like_num: 5 },
              { entry_id: '2', view_num: 20, like_num: 8 },
            ],
          },
        }),
      });

      diaryDetail.init();
      expect(mockQuerySelectorAll).toHaveBeenCalledWith('[data-entry-id]');
    });

    test('deleteEntry should confirm and call fetch', async () => {
      const diaryDetail = createDiaryDetail();
      mockBoxConfirm.mockResolvedValueOnce(true);

      // Mock window.htmx
      window.htmx = {
        ajax: jest.fn().mockResolvedValue({}),
      };

      await diaryDetail.deleteEntry(1, 100);

      expect(mockBoxConfirm).toHaveBeenCalledWith('确认删除?');
      expect(window.htmx.ajax).toHaveBeenCalledWith(
        'DELETE',
        '/entry/100?diary_id=1',
        expect.objectContaining({
          target: '#entry-card-100',
          swap: 'delete',
        }),
      );
    });

    test('loadDiaryLike should update state on success', async () => {
      const diaryDetail = createDiaryDetail();

      // Mock getElementById for container
      const mockContainer = { dataset: { diaryId: '123' } };
      jest.spyOn(document, 'getElementById').mockReturnValue(mockContainer);

      global.fetch.mockResolvedValueOnce({
        headers: {
          get: () => 'application/json',
        },
        json: async () => ({
          success: true,
          data: { is_liked: true, like_count: 5 },
        }),
      });

      await diaryDetail.loadDiaryLike();

      expect(diaryDetail.diaryLike.liked).toBe(true);
      expect(diaryDetail.diaryLike.count).toBe(5);
      expect(global.fetch).toHaveBeenCalledWith(
        expect.stringContaining('/api/diary/like/status'),
        expect.any(Object),
      );
    });

    test('toggleDiaryLike should update state and notify', async () => {
      const diaryDetail = createDiaryDetail();

      global.fetch.mockResolvedValueOnce({
        headers: {
          get: () => 'application/json',
        },
        json: async () => ({
          success: true,
          data: { is_liked: true, like_count: 6 },
        }),
      });

      await diaryDetail.toggleDiaryLike(123);

      expect(diaryDetail.diaryLike.liked).toBe(true);
      expect(diaryDetail.diaryLike.count).toBe(6);
      expect(mockNotifySuccess).toHaveBeenCalledWith('social.like.success');
    });

    test('openCreateEntry should call Box.dialog with hxUrl', () => {
      const diaryDetail = createDiaryDetail();

      // Mock Box.dialog return value
      const mockBoxId = 'box-123';
      mockBoxDialog.mockReturnValue(mockBoxId);

      diaryDetail.openCreateEntry(1);

      expect(mockBoxDialog).toHaveBeenCalledWith(
        expect.objectContaining({
          title: '写日记',
          hxUrl: expect.stringContaining('/form/entry?action=entry.create&diary_id=1'),
        }),
      );
    });

    test('openEntryModal should call Box.dialog with ajaxUrl', () => {
      const diaryDetail = createDiaryDetail();

      // Mock window.location
      const originalLocation = window.location;
      delete window.location;
      window.location = { href: 'http://localhost/' };

      // Mock history.pushState
      const pushStateSpy = jest.spyOn(window.history, 'pushState');

      diaryDetail.openEntryModal('/entry/1');

      expect(mockBoxDialog).toHaveBeenCalledWith(
        expect.objectContaining({
          hxUrl: '/entry/1?isajax=true',
        }),
      );

      // Restore window.location
      window.location = originalLocation;
    });
  });
});
