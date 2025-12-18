import { Notify } from '../utils/uno-shim.js';

export function initAuth() {
  // Configure HTMX to allow swapping on error status codes (4xx, 5xx)
  document.body.addEventListener('htmx:beforeSwap', function (evt) {
    const xhr = evt.detail.xhr;
    const contentType = xhr.getResponseHeader('Content-Type');

    // Handle JSON responses (e.g. validation errors from generic middleware)
    if (contentType && contentType.includes('application/json')) {
      evt.detail.shouldSwap = false;
      evt.detail.isError = true;

      try {
        const result = JSON.parse(xhr.responseText);
        if (!result.success) {
          let msg = result.message || 'Operation failed';
          // Append details if it's a string
          if (result.details && typeof result.details === 'string') {
            msg += ': ' + result.details;
          }
          Notify.error(msg);
        } else {
          // Success JSON
          Notify.success(result.message || 'Success');
        }
      } catch (e) {
        console.error('Failed to parse JSON response', e);
      }
      return;
    }

    if ([400, 401, 403, 404, 409, 422, 500].includes(evt.detail.xhr.status)) {
      evt.detail.shouldSwap = true;
      evt.detail.isError = false; // Tell HTMX to treat it as a successful swap
    }
  });

  // Handle HTMX redirect header if needed (though HX-Redirect handles it automatically)
  // But we might want to do something custom if needed.
}

/**
 * Validation callbacks module
 * Registers custom validation callbacks for the Uno valid plugin
 */
export function initValidationCallbacks() {
  if (!window.Alpine) {
    console.warn(
      '[Validation] Alpine not found, cannot register validation callbacks',
    );
    return;
  }

  const register = (retryCount = 0) => {
    // Check if store exists
    // Alpine.store(name) returns the store object if it exists
    // However, if it was just initialized, it might be available.
    // Note: Alpine.store('valid') returns the reactive proxy.

    let validStore;
    try {
      validStore = Alpine.store('valid');
    } catch (e) {
      // ignore
    }

    if (!validStore || !validStore.registerCallback) {
      if (retryCount < 10) {
        // console.log(`[Validation] Store not ready, retrying... (${retryCount + 1})`);
        setTimeout(() => register(retryCount + 1), 100);
      } else {
        console.error(
          '[Validation] Failed to register callbacks: Uno valid store not found',
        );
      }
      return;
    }

    // Register checkUsername callback
    validStore.registerCallback('checkUsername', async (id, value) => {
      if (!value) return true;

      try {
        const formData = new FormData();
        formData.append('username', value);

        const response = await fetch('/check-username', {
          method: 'POST',
          body: formData,
        });

        const result = await response.json();

        // If success is true, it means available
        if (result.success) {
          return true;
        }

        // Otherwise return error message
        return result.message || 'Username already taken';
      } catch (error) {
        console.error('Validation error:', error);
        return 'Validation failed';
      }
    });

    // Register checkEmail callback
    validStore.registerCallback('checkEmail', async (id, value) => {
      if (!value) return true;

      // Check for original value (for edit forms)
      const element = document.getElementById(id);
      if (element) {
        const originalValue = element.getAttribute('data-original-value');
        if (originalValue && value === originalValue) {
          return true;
        }
      }

      try {
        const formData = new FormData();
        formData.append('email', value);

        const response = await fetch('/check-email', {
          method: 'POST',
          body: formData,
        });

        const result = await response.json();

        // If success is true, it means available
        if (result.success) {
          return true;
        }

        // Otherwise return error message
        return result.message || 'Email already taken';
      } catch (error) {
        console.error('Validation error:', error);
        return 'Validation failed';
      }
    });

    console.log('[Validation] Callbacks registered successfully');
  };

  register();
}
