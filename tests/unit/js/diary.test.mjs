import { jest } from '@jest/globals';
import { DiaryList } from '../../../website/static/src/js/modules/diary.js';

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

  const createDiaryList = () => {
    const component = DiaryList();
    component.$store = {
      i18n: {
        t: (key) => {
          const map = {
            'diary.delete.confirm': '确认删除?',
            'diary.delete.success': '删除成功',
            'diary.delete.fail': '删除失败',
            'diary.delete.net_error': '网络错误',
            'common.op_fail': '操作失败',
            'diary.edit.title': '编辑日记',
            'diary.create.title': '写日记',
            'common.save': '保存',
            'common.loading': '加载中...',
            'operation.success': '操作成功',
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

  describe('DiaryList', () => {
    test('should initialize and load stats', () => {
      const diaryList = createDiaryList();

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
        { dataset: { postId: '1' } },
        { dataset: { postId: '2' } },
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
              { post_id: '1', view_num: 10, like_num: 5 },
              { post_id: '2', view_num: 20, like_num: 8 },
            ],
          },
        }),
      });

      diaryList.init();
      expect(mockQuerySelectorAll).toHaveBeenCalledWith('[data-post-id]');
    });

    test('confirmDelete should call Uno.box.confirm', () => {
      const diaryList = createDiaryList();
      mockBoxConfirm.mockResolvedValueOnce(true);
      diaryList.confirmDelete(1, 100);
      expect(mockBoxConfirm).toHaveBeenCalledWith(
        '确认删除?',
        expect.objectContaining({
          title: expect.any(String),
          okText: expect.any(String),
          cancelText: expect.any(String),
        }),
      );
    });

    test('deletePost should call fetch and remove element on success', async () => {
      const diaryList = createDiaryList();

      // Mock document.querySelector
      const mockElement = { remove: jest.fn() };
      const mockQuerySelector = jest.fn().mockReturnValue(mockElement);
      Object.defineProperty(document, 'querySelector', {
        value: mockQuerySelector,
        writable: true,
      });

      global.fetch.mockResolvedValueOnce({
        json: async () => ({ success: true }),
      });

      await diaryList.deletePost(1, 100);

      expect(global.fetch).toHaveBeenCalledWith(
        '/post/100?diary_id=1',
        expect.any(Object),
      );
      expect(mockElement.remove).toHaveBeenCalled();
      expect(mockNotifySuccess).toHaveBeenCalled();
    });

    test('loadDiaryLike should update state on success', async () => {
      const diaryList = createDiaryList();

      // Mock getElementById for container
      const mockContainer = { dataset: { diaryId: '123' } };
      jest.spyOn(document, 'getElementById').mockReturnValue(mockContainer);

      global.fetch.mockResolvedValueOnce({
        json: async () => ({
          success: true,
          data: { is_liked: true, like_count: 5 },
        }),
      });

      await diaryList.loadDiaryLike();

      expect(diaryList.diaryLike.liked).toBe(true);
      expect(diaryList.diaryLike.count).toBe(5);
      expect(global.fetch).toHaveBeenCalledWith(
        expect.stringContaining('/api/diary/like/status'),
        expect.any(Object),
      );
    });

    test('toggleDiaryLike should update state and notify', async () => {
      const diaryList = createDiaryList();

      global.fetch.mockResolvedValueOnce({
        json: async () => ({
          success: true,
          data: { is_liked: true, like_count: 6 },
        }),
      });

      await diaryList.toggleDiaryLike(123);

      expect(diaryList.diaryLike.liked).toBe(true);
      expect(diaryList.diaryLike.count).toBe(6);
      expect(mockNotifySuccess).toHaveBeenCalledWith('social.like.success');
    });
    test('openCreatePost should call Box.dialog and fetch form', async () => {
      const diaryList = createDiaryList();

      // Mock Box.dialog return value
      const mockBoxId = 'box-123';
      mockBoxDialog.mockReturnValue(mockBoxId);

      // Mock fetch for form
      global.fetch.mockResolvedValueOnce({
        text: async () => '<form>Create Post</form>',
      });

      // Mock getElementById for form container
      const mockContainer = { innerHTML: '' };
      jest.spyOn(document, 'getElementById').mockReturnValue(mockContainer);

      diaryList.openCreatePost(1);

      expect(mockBoxDialog).toHaveBeenCalledWith(
        expect.objectContaining({
          title: '写日记',
        }),
      );
      expect(global.fetch).toHaveBeenCalledWith(
        expect.stringContaining('/form/post?action=post.create&diary_id=1'),
      );
    });

    test('openPostModal should call Box.dialog with ajaxUrl', () => {
      const diaryList = createDiaryList();

      diaryList.openPostModal('/p/1');

      expect(mockBoxDialog).toHaveBeenCalledWith(
        expect.objectContaining({
          hxUrl: '/p/1?isajax=true',
        }),
      );
    });
  });
});
