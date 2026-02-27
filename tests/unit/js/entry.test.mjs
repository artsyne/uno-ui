import { jest } from '@jest/globals';
import { EntryMixin } from '../../../website/static/src/js/modules/entry.js';

describe('EntryMixin', () => {
  let mockUno;
  let mockHtmx;
  let mockAlpine;
  let mockContext;
  let originalPushState;

  beforeEach(() => {
    // Reset mocks
    mockUno = {
      box: {
        dialog: jest.fn(() => ({ id: 'box-123', close: jest.fn() })),
        confirm: jest.fn(() => Promise.resolve(true)),
        close: jest.fn(),
        get: jest.fn(() => ({ close: jest.fn() })),
      },
      notify: {
        success: jest.fn(),
      },
    };
    mockHtmx = {
      process: jest.fn(),
      ajax: jest.fn().mockResolvedValue(),
    };
    mockAlpine = {
      initTree: jest.fn(),
    };

    // Attach to JSDOM window
    window.getUno = () => mockUno;
    window.getUnoBox = () => mockUno.box;
    window.htmx = mockHtmx;
    window.Alpine = mockAlpine;

    // Spy on history API
    originalPushState = window.history.pushState;
    window.history.pushState = jest.fn();

    // Mock document methods
    jest.spyOn(document, 'getElementById').mockReturnValue(null);
    jest.spyOn(document.body, 'addEventListener').mockImplementation(() => {});
    jest
      .spyOn(document.body, 'removeEventListener')
      .mockImplementation(() => {});

    // Mock Context for Mixin
    mockContext = {
      $store: {
        i18n: {
          t: (key) => key,
        },
      },
      ...EntryMixin,
      // Reset handler if attached to instance (though implementation attaches to this)
      _closeEntryModalHandler: null,
    };
  });

  afterEach(() => {
    jest.restoreAllMocks();
    window.history.pushState = originalPushState;
    delete window.getUno;
    delete window.getUnoBox;
    delete window.htmx;
    delete window.Alpine;
  });

  describe('openEntryModal', () => {
    test('updates history and opens dialog', () => {
      const url = '/entry/123';
      mockContext.openEntryModal(url);

      expect(window.history.pushState).toHaveBeenCalledWith(null, '', url);
      expect(mockUno.box.dialog).toHaveBeenCalledWith(
        expect.objectContaining({
          hxUrl: url + '?isajax=true',
          onContentLoad: expect.any(Function),
          onClose: expect.any(Function),
        }),
      );
    });

    test('onContentLoad initializes HTMX and Alpine', () => {
      mockContext.openEntryModal('/entry/1');
      const config = mockUno.box.dialog.mock.calls[0][0];
      const mockEl = {};

      config.onContentLoad(mockEl);

      expect(mockHtmx.process).toHaveBeenCalledWith(mockEl);
      expect(mockAlpine.initTree).toHaveBeenCalledWith(mockEl);
    });

    test('onClose restores history', () => {
      mockContext.openEntryModal('/entry/1');
      const config = mockUno.box.dialog.mock.calls[0][0];
      const originalUrl = 'http://localhost/'; // JSDOM default

      config.onClose();

      expect(window.history.pushState).toHaveBeenCalledWith(
        null,
        '',
        originalUrl,
      );
    });
  });

  describe('_openEntryForm (Create/Edit)', () => {
    test('opens create dialog with correct config', () => {
      mockContext.openCreateEntry(1);

      expect(mockUno.box.dialog).toHaveBeenCalledWith(
        expect.objectContaining({
          title: 'entry.create',
          hxUrl: expect.stringContaining(
            '/form/entry?action=entry.create&diary_id=1',
          ),
          onOk: expect.any(Function),
        }),
      );
    });

    test('opens edit dialog with correct config', () => {
      mockContext.openEditEntry(1, 100);

      expect(mockUno.box.dialog).toHaveBeenCalledWith(
        expect.objectContaining({
          title: 'entry.edit',
          hxUrl: expect.stringContaining(
            '/form/entry?action=entry.edit&diary_id=1&entry_id=100',
          ),
        }),
      );
    });

    test('onOk triggers requestSubmit', async () => {
      const mockForm = {
        requestSubmit: jest.fn(),
        querySelector: jest.fn().mockReturnValue({ value: 'entry.create' }),
      };
      document.getElementById.mockReturnValue(mockForm);

      mockContext.openCreateEntry(1);
      const config = mockUno.box.dialog.mock.calls[0][0];

      const result = await config.onOk();

      expect(document.getElementById).toHaveBeenCalledWith('entry_editor_form');
      expect(mockForm.requestSubmit).toHaveBeenCalled();
      expect(result).toBe(false);
    });

    test('registers closeEntryModal handler', () => {
      mockContext.openCreateEntry(1);

      expect(document.body.addEventListener).toHaveBeenCalledWith(
        'closeEntryModal',
        expect.any(Function),
        true,
      );
    });

    test('closeEntryModal handler closes the box', () => {
      mockContext.openCreateEntry(1);

      const config = mockUno.box.dialog.mock.calls[0][0];
      const generatedBoxId = config.boxId;

      const handler = document.body.addEventListener.mock.calls[0][1];
      handler({ detail: { boxId: generatedBoxId } });

      // Verify boxInstance.close() was called (primary path)
      const boxInstance = mockUno.box.dialog.mock.results[0].value;
      expect(boxInstance.close).toHaveBeenCalled();
    });
  });

  describe('deleteEntry', () => {
    test('confirms and sends DELETE request', async () => {
      await mockContext.deleteEntry(1, 100);

      expect(mockUno.box.confirm).toHaveBeenCalledWith('entry.delete.confirm');
      expect(mockHtmx.ajax).toHaveBeenCalledWith(
        'DELETE',
        '/entry/100?diary_id=1',
        expect.objectContaining({
          target: '#entry-card-100',
          swap: 'delete',
        }),
      );
    });

    test('does not send request if confirm cancelled', async () => {
      mockUno.box.confirm.mockResolvedValue(false);

      await mockContext.deleteEntry(1, 100);

      expect(mockHtmx.ajax).not.toHaveBeenCalled();
    });
  });
});
