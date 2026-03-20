
import { jest } from '@jest/globals';
import { setFormLoading, initFormHandlers } from '../../../website/static/src/js/modules/form.js';

describe('Form Module', () => {
  let form;
  let mask;
  let input;
  let submitBtn;
  let externalBtn;

  beforeEach(() => {
    // Reset DOM
    document.body.innerHTML = `
      <div id="container">
        <form id="test-form">
          <div class="form-mask" style="visibility: hidden; opacity: 0;"></div>
          <input type="text" name="username" />
          <button type="submit">Submit</button>
        </form>
        <button id="external-btn">External</button>
      </div>
    `;

    form = document.getElementById('test-form');
    mask = form.querySelector('.form-mask');
    input = form.querySelector('input');
    submitBtn = form.querySelector('button[type="submit"]');
    externalBtn = document.getElementById('external-btn');

    // Mock requestAnimationFrame for synchronous testing
    global.requestAnimationFrame = (callback) => {
      setTimeout(callback, 0);
    };

    // Initialize handlers
    initFormHandlers();
  });

  afterEach(() => {
    jest.clearAllMocks();
    document.body.innerHTML = '';
  });

  describe('setFormLoading', () => {
    test('should show mask but NOT disable inputs when loading is true', async () => {
      setFormLoading(form, true);

      expect(mask.style.visibility).toBe('visible');

      // Wait a bit for requestAnimationFrame
      await new Promise(resolve => setTimeout(resolve, 10));

      expect(mask.style.opacity).toBe('1');
      // Inputs should remain enabled because mask handles interaction blocking via z-index
      expect(input.disabled).toBe(false);
      expect(submitBtn.disabled).toBe(false);
      expect(externalBtn.disabled).toBe(false);
    });

    test('should hide mask when loading is false', async () => {
      // First set to loading
      setFormLoading(form, true);

      // Wait for loading to be applied
      await new Promise(resolve => setTimeout(resolve, 10));

      // Then set to not loading
      setFormLoading(form, false);

      expect(mask.style.opacity).toBe('0');

      // Wait for transition to complete
      await new Promise(resolve => setTimeout(resolve, 350));

      expect(mask.style.visibility).toBe('hidden');
    });

    test('should create mask inside form if not found', async () => {
      // Remove existing mask
      form.removeChild(mask);

      setFormLoading(form, true);

      // Should create a new mask inside the form
      const newMask = form.querySelector('.form-mask');
      expect(newMask).toBeTruthy();
      expect(newMask.style.visibility).toBe('visible');

      // Wait for requestAnimationFrame
      await new Promise(resolve => setTimeout(resolve, 10));

      expect(newMask.style.opacity).toBe('1');
    });
  });

  describe('HTMX Event Integration', () => {
    test('should trigger loading on htmx:beforeRequest', async () => {
      const event = new CustomEvent('htmx:beforeRequest', {
        detail: { elt: form }
      });
      document.body.dispatchEvent(event);

      expect(mask.style.visibility).toBe('visible');

      // Wait for requestAnimationFrame
      await new Promise(resolve => setTimeout(resolve, 10));

      expect(mask.style.opacity).toBe('1');
    });

    test('should stop loading on htmx:afterRequest', async () => {
      // Set loading first
      setFormLoading(form, true);

      // Wait for loading to be applied
      await new Promise(resolve => setTimeout(resolve, 10));

      const event = new CustomEvent('htmx:afterRequest', {
        detail: { elt: form }
      });
      document.body.dispatchEvent(event);

      expect(mask.style.opacity).toBe('0');

      // Wait for visibility transition
      await new Promise(resolve => setTimeout(resolve, 350));

      expect(mask.style.visibility).toBe('hidden');
    });
  });
});
