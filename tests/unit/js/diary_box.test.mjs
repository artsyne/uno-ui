import { jest } from '@jest/globals';
import DiaryDetail from '../../../website/static/src/js/modules/diary-detail.js';

// Mock dependencies
jest.mock('../../../website/static/src/js/modules/social.js', () => ({
  SocialMixin: () => ({ initSocial: jest.fn() }),
}));
jest.mock('../../../website/static/src/js/modules/entry.js', () => ({
  EntryMixin: {},
}));

describe('DiaryDetail Box Logic', () => {
  let diaryDetail;
  let mockUno;
  let mockBoxId = 'box-diary-settings';
  let mockForm;

  beforeEach(() => {
    // Setup DOM
    document.body.innerHTML = '';

    // Setup Mocks
    mockUno = {
      box: {
        dialog: jest.fn().mockReturnValue({ id: mockBoxId }),
        close: jest.fn(),
      },
      valid: {
        validateForm: jest.fn().mockResolvedValue(true),
      },
    };

    // Attach to JSDOM window directly
    window.getUno = jest.fn(() => mockUno);
    window.Alpine = { initTree: jest.fn() };
    window.htmx = { process: jest.fn() };
    window.openDiarySettings = jest.fn(); // Mock global export

    // Initialize module
    diaryDetail = DiaryDetail();
    diaryDetail.$store = {
      i18n: { t: (key) => key },
    };
  });

  afterEach(() => {
    jest.restoreAllMocks();
    delete window.getUno;
    delete window.Alpine;
    delete window.htmx;
    delete window.openDiarySettings;
  });

  test('openDiarySettings attaches boxId to form in onContentLoad', () => {
    // Trigger open
    diaryDetail.openDiarySettings(1);

    // Get the config passed to dialog
    const config = mockUno.box.dialog.mock.calls[0][0];
    const generatedBoxId = config.boxId;

    // Create a mock element with a form
    const container = document.createElement('div');
    const form = document.createElement('form');
    container.appendChild(form);

    // Simulate content load
    config.onContentLoad(container);

    // Verify boxId attached
    expect(form.dataset.boxId).toBe(generatedBoxId);
  });

  test('openDiarySettings registers event listener that closes specific boxId from target', () => {
    // Spy on addEventListener
    const addListenerSpy = jest.spyOn(document.body, 'addEventListener');

    diaryDetail.openDiarySettings(1);

    // Get generated ID
    const config = mockUno.box.dialog.mock.calls[0][0];
    const generatedBoxId = config.boxId;

    // Find the closeHandler
    const call = addListenerSpy.mock.calls.find(
      (c) => c[0] === 'closeDiaryModal',
    );
    expect(call).toBeDefined();
    const handler = call[1];
    expect(call[2]).toBe(true);

    // Simulate event with matching boxId on target
    const evt = { target: { dataset: { boxId: generatedBoxId } } };
    handler(evt);

    expect(mockUno.box.close).toHaveBeenCalledWith(generatedBoxId);
  });

  test('openDiarySettings fallback to closure boxId if target has no boxId', () => {
    const addListenerSpy = jest.spyOn(document.body, 'addEventListener');
    diaryDetail.openDiarySettings(1);

    // Get generated ID
    const config = mockUno.box.dialog.mock.calls[0][0];
    const generatedBoxId = config.boxId;

    const handler = addListenerSpy.mock.calls.find(
      (c) => c[0] === 'closeDiaryModal',
    )[1];

    // Simulate event without target boxId
    handler({});

    expect(mockUno.box.close).toHaveBeenCalledWith(generatedBoxId);
  });

  test('openDiarySettings listener ignores event with different boxId', () => {
    const addListenerSpy = jest.spyOn(document.body, 'addEventListener');
    diaryDetail.openDiarySettings(1);

    // Get generated ID
    const config = mockUno.box.dialog.mock.calls[0][0];
    const generatedBoxId = config.boxId;

    const handler = addListenerSpy.mock.calls.find(
      (c) => c[0] === 'closeDiaryModal',
    )[1];

    // Simulate event with DIFFERENT boxId
    const evt = { target: { dataset: { boxId: 'other-box' } } };
    handler(evt);

    expect(mockUno.box.close).not.toHaveBeenCalled();
  });
});
