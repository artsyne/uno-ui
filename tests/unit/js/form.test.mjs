
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
          <div class="form-mask" style="display: none;"></div>
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

    // Initialize handlers
    initFormHandlers();
  });

  afterEach(() => {
    jest.clearAllMocks();
    document.body.innerHTML = '';
  });

  describe('setFormLoading', () => {
    test('should show mask but NOT disable inputs when loading is true', () => {
      setFormLoading(form, true);

      expect(mask.style.display).toBe('flex');
      // Inputs should remain enabled because mask handles interaction blocking via z-index
      expect(input.disabled).toBe(false);
      expect(submitBtn.disabled).toBe(false);
      expect(externalBtn.disabled).toBe(false);
    });

    test('should hide mask when loading is false', () => {
      // First set to loading
      setFormLoading(form, true);
      
      // Then set to not loading
      setFormLoading(form, false);

      expect(mask.style.display).toBe('none');
    });

    test('should look for mask in parent if not found in form', () => {
      // Reconstruct DOM with mask outside form
      document.body.innerHTML = `
        <div id="container" class="relative">
          <div class="form-mask" style="display: none;"></div>
          <form id="test-form">
            <input type="text" name="username" />
            <button type="submit">Submit</button>
          </form>
        </div>
      `;
      
      form = document.getElementById('test-form');
      mask = document.querySelector('.form-mask'); // It's outside form now
      
      setFormLoading(form, true);
      expect(mask.style.display).toBe('flex');
    });
  });

  describe('HTMX Event Integration', () => {
    test('should trigger loading on htmx:beforeRequest', () => {
      const event = new CustomEvent('htmx:beforeRequest', {
        detail: { elt: form }
      });
      document.body.dispatchEvent(event);

      expect(mask.style.display).toBe('flex');
    });

    test('should stop loading on htmx:afterRequest', () => {
      // Set loading first
      setFormLoading(form, true);

      const event = new CustomEvent('htmx:afterRequest', {
        detail: { elt: form }
      });
      document.body.dispatchEvent(event);

      expect(mask.style.display).toBe('none');
    });
  });
});
