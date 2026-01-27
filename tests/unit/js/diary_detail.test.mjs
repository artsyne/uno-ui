import { jest } from '@jest/globals';
import DiaryDetail from '../../../website/static/src/js/modules/diary-detail.js';

describe('DiaryDetail', () => {
  let mockUno;
  let mockStore;
  let diaryDetail;

  beforeEach(() => {
    // Mock Uno UI
    mockUno = {
      box: {
        dialog: jest.fn(() => ({ id: 'box-123' })),
        close: jest.fn(),
      },
      valid: {
        validateForm: jest.fn(),
      },
    };

    // Mock Alpine Store
    mockStore = {
      i18n: {
        t: jest.fn((key) => key),
      },
      uno: mockUno,
    };

    // Mock Global Window
    global.window = global;
    global.getUno = () => mockUno;
    global.window.getUno = () => mockUno;

    global.Alpine = {
      initTree: jest.fn(),
      store: (name) => {
        if (name === 'uno') return mockUno;
        if (name === 'i18n') return mockStore.i18n;
      },
    };
    global.window.Alpine = global.Alpine;

    global.htmx = {
      process: jest.fn(),
    };
    global.window.htmx = global.htmx;

    // Setup Document Body for Event Listeners
    document.body.addEventListener = jest.fn();
    document.body.removeEventListener = jest.fn();
    document.getElementById = jest.fn();
    document.querySelector = jest.fn();

    // Initialize Component
    diaryDetail = DiaryDetail();
    diaryDetail.$store = mockStore;
    // Mock internal methods
    diaryDetail._escapeHtml = jest.fn((str) => str);
    diaryDetail.loadStats = jest.fn();
    diaryDetail.loadDiaryLike = jest.fn();
    diaryDetail.initSocial = jest.fn();

    // Initialize
    diaryDetail.init();
  });

  test('openDiarySettings configures box correctly', () => {
    diaryDetail.openDiarySettings(1);

    expect(mockUno.box.dialog).toHaveBeenCalledWith(
      expect.objectContaining({
        title: 'diary.settings',
        hxUrl: '/form/diary?action=diary.edit&diary_id=1',
        showOkButton: true,
        okText: 'common.save',
      }),
    );
  });

  test('openDiarySettings registers close listener', () => {
    diaryDetail.openDiarySettings(1);

    const closeCall = document.body.addEventListener.mock.calls.find(
      (call) => call[0] === 'closeDiaryModal',
    );

    expect(closeCall[2]).toBe(true);
  });

  test('onOk validates form and dispatches event', async () => {
    // Setup dialog to get the config
    diaryDetail.openDiarySettings(1);
    const config = mockUno.box.dialog.mock.calls[0][0];

    // Mock Form
    const mockForm = {
      dispatchEvent: jest.fn(),
    };
    document.getElementById.mockReturnValue(mockForm);
    mockUno.valid.validateForm.mockResolvedValue(true);

    // Call onOk
    const result = await config.onOk();

    expect(mockUno.valid.validateForm).toHaveBeenCalledWith(mockForm);
    expect(mockForm.dispatchEvent).toHaveBeenCalledWith(
      expect.any(CustomEvent),
    );
    expect(result).toBe(false); // Should return false to keep box open (HTMX handles close)
  });

  test('onContentLoad initializes Alpine and HTMX', () => {
    diaryDetail.openDiarySettings(1);
    const config = mockUno.box.dialog.mock.calls[0][0];
    const mockForm = document.createElement('form');
    const mockEl = {
      querySelector: jest.fn((selector) => {
        if (selector === 'form') return mockForm;
        return null;
      }),
    };

    // Mock getElementById to return our mockEl because onContentLoad calls it
    document.getElementById.mockReturnValue(mockEl);

    config.onContentLoad(mockEl);

    expect(global.Alpine.initTree).toHaveBeenCalledWith(mockEl);
    expect(global.htmx.process).toHaveBeenCalledWith(mockEl);
    // 0.123456789.toString(36) is '0.4fzzzxj...' -> slice(2,10) -> '4fzzzxj'
    // id prefix: box-diary-settings-
    expect(mockForm.dataset.boxId).toMatch(/^box-diary-settings-/);
  });

  test('diaryUpdated event updates UI', () => {
    // We need to simulate the event listener being added and triggered
    // But since we mocked addEventListener, we can't trigger it via document.body.dispatchEvent naturally
    // So we will inspect the arguments passed to addEventListener

    // Find the 'diaryUpdated' listener (added in init)
    const calls = document.body.addEventListener.mock.calls;
    const updateCall = calls.find((call) => call[0] === 'diaryUpdated');
    expect(updateCall).toBeDefined();

    const handler = updateCall[1];

    // Mock UI elements
    const mockTitle = { innerHTML: '' };
    const mockRemark = { innerText: '' };
    document.getElementById.mockImplementation((id) => {
      if (id === 'diary_title') return mockTitle;
      if (id === 'diary_remark') return mockRemark;
      return null;
    });

    // Trigger handler
    const diaryData = {
      diary_id: 1,
      title: 'New Title',
      remark: 'New Remark',
      view_setting: 1,
      theme_id: 1,
      social_setting: 0,
    };
    handler({ detail: { diaryUpdated: diaryData } });

    expect(mockTitle.innerHTML).toContain('New Title');
    expect(mockRemark.innerText).toBe('New Remark');
  });

  test('onClose removes event listeners', () => {
    diaryDetail.openDiarySettings(1);
    const config = mockUno.box.dialog.mock.calls[0][0];

    // Simulate close
    config.onClose();

    expect(document.body.removeEventListener).toHaveBeenCalledWith(
      'closeDiaryModal',
      expect.any(Function),
      true,
    );
  });
});
