import { jest } from '@jest/globals';
import { User } from '../../../website/static/src/js/modules/user.js';

describe('User Module Box Logic', () => {
  let mockUno;
  let mockStore;
  let userModule;

  beforeEach(() => {
    mockUno = {
      box: {
        dialog: jest.fn(() => ({ id: 'box-123' })),
        close: jest.fn(),
      },
    };

    mockStore = {
      i18n: {
        t: jest.fn((key) => key),
      },
    };

    global.window = global;
    global.getUno = () => mockUno;
    global.window.getUno = () => mockUno;
    global.htmx = { process: jest.fn() };
    global.Alpine = { initTree: jest.fn() };

    document.body.addEventListener = jest.fn();
    document.body.removeEventListener = jest.fn();
    document.getElementById = jest.fn();

    userModule = { ...User, $store: mockStore };
  });

  test('openCreateDiary registers event listener that closes specific boxId', () => {
    userModule.openCreateDiary();

    // 1. Verify dialog created
    expect(mockUno.box.dialog).toHaveBeenCalled();
    const config = mockUno.box.dialog.mock.calls[0][0];
    const generatedBoxId = config.boxId;
    expect(generatedBoxId).toMatch(/^box-diary-/);

    // 2. Verify listener registered
    expect(document.body.addEventListener).toHaveBeenCalledWith(
      'closeDiaryModal',
      expect.any(Function),
      true,
    );

    // 3. Get the handler
    const calls = document.body.addEventListener.mock.calls;
    const handler = calls.find((call) => call[0] === 'closeDiaryModal')[1];

    // 4. Trigger handler
    handler();

    // 5. Verify close called with correct ID
    expect(mockUno.box.close).toHaveBeenCalledWith(generatedBoxId);
  });

  test('onClose removes event listener', () => {
    userModule.openCreateDiary();
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
