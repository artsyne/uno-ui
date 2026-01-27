import { jest } from '@jest/globals';
import { User } from '../../../website/static/src/js/modules/user.js';

describe('User Module', () => {
  let mockUno;
  let mockHtmx;
  let mockAlpine;
  let mockContext;

  beforeEach(() => {
    // Reset mocks
    mockUno = {
      box: {
        dialog: jest.fn(() => ({ id: 'box-diary-123' })),
        close: jest.fn(),
      },
    };
    mockHtmx = {
      process: jest.fn(),
      trigger: jest.fn(),
    };
    mockAlpine = {
      initTree: jest.fn(),
    };

    // Attach to JSDOM window
    window.getUno = () => mockUno;
    window.getUnoBox = () => mockUno.box;
    window.htmx = mockHtmx;
    window.Alpine = mockAlpine;

    // Mock document methods
    jest.spyOn(document, 'getElementById').mockReturnValue(null);
    jest.spyOn(document.body, 'addEventListener').mockImplementation(() => {});
    jest
      .spyOn(document.body, 'removeEventListener')
      .mockImplementation(() => {});

    // Mock Math.random for predictable IDs
    jest.spyOn(Math, 'random').mockReturnValue(0.123456789);

    // Mock Context for User
    mockContext = {
      $store: {
        i18n: {
          t: (key) => key,
        },
      },
      ...User,
    };
  });

  afterEach(() => {
    jest.restoreAllMocks();
    delete window.getUno;
    delete window.getUnoBox;
    delete window.htmx;
    delete window.Alpine;
  });

  test('openCreateDiary uses hxUrl and configures callbacks correctly', () => {
    mockContext.openCreateDiary();

    expect(mockUno.box.dialog).toHaveBeenCalledWith(
      expect.objectContaining({
        title: 'diary.create',
        hxUrl: '/form/diary?action=diary.create',
        onContentLoad: expect.any(Function),
        onOk: expect.any(Function),
        onClose: expect.any(Function),
      }),
    );
  });

  test('onContentLoad initializes HTMX and Alpine', () => {
    mockContext.openCreateDiary();
    const config = mockUno.box.dialog.mock.calls[0][0];
    const mockForm = document.createElement('form');
    const mockEl = {
      querySelector: jest.fn((selector) => {
        if (selector === 'form') return mockForm;
        return null;
      }),
    };

    config.onContentLoad(mockEl);

    expect(mockHtmx.process).toHaveBeenCalledWith(mockEl);
    expect(mockAlpine.initTree).toHaveBeenCalledWith(mockEl);
    // 0.123456789.toString(36) is '0.4fzzzxj...'
    // slice(2, 10) is '4fzzzxj'
    // id is 'box-diary-4fzzzxj'
    expect(mockForm.dataset.boxId).toMatch(/^box-diary-/);
  });

  test('onOk triggers form submission and returns false', () => {
    const mockForm = { requestSubmit: jest.fn() };
    document.getElementById.mockReturnValue(mockForm);

    mockContext.openCreateDiary();
    const config = mockUno.box.dialog.mock.calls[0][0];

    const result = config.onOk();

    expect(document.getElementById).toHaveBeenCalledWith('diary_editor_form');
    expect(mockForm.requestSubmit).toHaveBeenCalled();
    expect(result).toBe(false);
  });

  test('onClose removes event listener', () => {
    mockContext.openCreateDiary();
    const config = mockUno.box.dialog.mock.calls[0][0];

    config.onClose();

    expect(document.body.removeEventListener).toHaveBeenCalledWith(
      'closeDiaryModal',
      expect.any(Function),
      true,
    );
  });

  test('closeHandler closes the box', () => {
    // To test the internal closeHandler, we need to capture it from addEventListener
    mockContext.openCreateDiary();

    // Check that event listener was added
    expect(document.body.addEventListener).toHaveBeenCalledWith(
      'closeDiaryModal',
      expect.any(Function),
      true,
    );

    // Get the handler
    const handler = document.body.addEventListener.mock.calls[0][1];

    // Get the generated boxId from the dialog config
    const config = mockUno.box.dialog.mock.calls[0][0];
    const generatedBoxId = config.boxId;

    // Call it
    handler();

    expect(mockUno.box.close).toHaveBeenCalledWith(generatedBoxId);
  });
});
