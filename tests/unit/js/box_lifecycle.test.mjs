import { jest } from '@jest/globals';
import DiaryDetail from '../../../website/static/src/js/modules/diary-detail.js';

describe('Box Lifecycle Complex Scenarios', () => {
  let mockUno;
  let diaryDetail;
  const mockBoxId1 = 'box-diary-settings-1';
  const mockBoxId2 = 'box-diary-settings-2';

  beforeEach(() => {
    document.body.innerHTML = '';
    let dialogCallCount = 0;
    mockUno = {
      box: {
        dialog: jest.fn((config) => {
          dialogCallCount++;
          if (dialogCallCount === 1) return { id: mockBoxId1 };
          if (dialogCallCount === 2) return { id: mockBoxId2 };
          return { id: config.id || `box-generated-${dialogCallCount}` };
        }),
        close: jest.fn(),
      },
    };
    window.getUno = jest.fn(() => mockUno);
    window.Alpine = { initTree: jest.fn() };
    window.htmx = { process: jest.fn() };

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
  });

  test('should close the correct box when two are opened concurrently', () => {
    // 1. Open first box
    diaryDetail.openDiarySettings(1);
    const config1 = mockUno.box.dialog.mock.calls[0][0];
    const generatedBoxId1 = config1.boxId;

    const form1 = document.createElement('form');
    const el1 = {
      querySelector: (selector) => {
        if (selector === 'form') return form1;
        return null;
      },
    };
    config1.onContentLoad(el1);

    // 2. Open second box
    diaryDetail.openDiarySettings(2);
    const config2 = mockUno.box.dialog.mock.calls[1][0];
    const generatedBoxId2 = config2.boxId;

    const form2 = document.createElement('form');
    const el2 = {
      querySelector: (selector) => {
        if (selector === 'form') return form2;
        return null;
      },
    };
    config2.onContentLoad(el2);

    // 3. Trigger close event for the second box
    const closeEvent = new CustomEvent('closeDiaryModal', {
      detail: { target: { dataset: { boxId: generatedBoxId2 } } },
    });
    Object.defineProperty(closeEvent, 'target', {
      value: { dataset: { boxId: generatedBoxId2 } },
    });
    document.body.dispatchEvent(closeEvent);

    // 4. Assert only the second box is closed
    expect(mockUno.box.close).toHaveBeenCalledTimes(1);
    expect(mockUno.box.close).toHaveBeenCalledWith(generatedBoxId2);
    expect(mockUno.box.close).not.toHaveBeenCalledWith(generatedBoxId1);

    // 5. Clean up the second box's listener
    config2.onClose();

    // 6. Trigger close for the first box
    const closeEvent1 = new CustomEvent('closeDiaryModal', {
      detail: { target: { dataset: { boxId: generatedBoxId1 } } },
    });
    Object.defineProperty(closeEvent1, 'target', {
      value: { dataset: { boxId: generatedBoxId1 } },
    });
    document.body.dispatchEvent(closeEvent1);

    // Assert the first box is now also closed
    expect(mockUno.box.close).toHaveBeenCalledTimes(2);
    expect(mockUno.box.close).toHaveBeenCalledWith(generatedBoxId1);
  });

  test('should still close box after HTMX swaps the form element', () => {
    // 1. Open the box
    diaryDetail.openDiarySettings(1);
    const config = mockUno.box.dialog.mock.calls[0][0];
    const generatedBoxId = config.boxId;

    const originalForm = document.createElement('form');
    const el = {
      querySelector: (selector) => {
        if (selector === 'form') return originalForm;
        return null;
      },
    };
    config.onContentLoad(el);

    // 2. Simulate HTMX swap: The new form element still has the correct data-box-id
    const newForm = { dataset: { boxId: generatedBoxId } };

    // 3. Trigger close event, targeting the *new* form-like object
    const closeEvent = new CustomEvent('closeDiaryModal', {
      detail: { target: newForm },
    });
    Object.defineProperty(closeEvent, 'target', { value: newForm });
    document.body.dispatchEvent(closeEvent);

    // 4. Assert the box is closed
    expect(mockUno.box.close).toHaveBeenCalledTimes(1);
    expect(mockUno.box.close).toHaveBeenCalledWith(generatedBoxId);
  });
});
