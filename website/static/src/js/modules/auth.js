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
        const notify = window.getUnoNotify();

        if (!result.success) {
          let msg = result.message || 'common.op_fail';
          // Append details if it's a string
          if (result.details && typeof result.details === 'string') {
            msg += ': ' + result.details;
          }
          notify.error(msg);
        } else {
          // Success JSON
          notify.success(result.message || 'auth.success');
        }
      } catch (e) {
        // ignore
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

// Removed ValidForm: Obsolete. Uno now handles form interception via x-uno-valid and htmx:confirm.

/**
 * Validation callbacks module
 * Registers custom validation callbacks for the Uno valid plugin
 */
export function initValidationCallbacks() {
  // Get Uno valid store
  // Access via main 'uno' store is the canonical way
  const unoStore = Alpine.store('uno');
  const validStore = unoStore?.valid || Alpine.store('valid');

  // Register checkUsername callback
  // Pattern: Sync Return + Async Side Effect
  // Uno's validator expects a synchronous boolean return.
  // For async checks (server-side uniqueness), we return 'true' immediately to pass the sync check,
  // then perform the async fetch. If it fails, we manually set the error using validStore.showError().
  validStore.registerValidator('checkUsername', async (value, name) => {
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

      // Otherwise show error manually because validation is sync
      const t = (key) => Alpine.store('i18n').t(key);
      const msg = result.message || t('auth.username_taken');
      validStore.showError(name, msg);
      return true; // Return true to prevent sync validation failure, error is shown manually
    } catch (error) {
      const t = (key) => Alpine.store('i18n').t(key);
      const msg = t('common.op_fail');
      validStore.showError(name, msg);
      return true;
    }
  });

  // Register checkEmail callback
  validStore.registerValidator('checkEmail', async (value, name) => {
    if (!value) return true;

    // Check for original value (for edit forms)
    // Try to find element by name (as name is passed) or fallback to id if name matches id
    const element =
      document.querySelector(`[name="${name}"]`) ||
      document.getElementById(name);
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

      // Otherwise show error manually
      const t = (key) => Alpine.store('i18n').t(key);
      const msg = result.message || t('auth.email_taken');
      validStore.showError(name, msg);
      return true;
    } catch (error) {
      return true;
    }
  });

  // Register matchPwd callback
  validStore.registerValidator('matchPwd', (value, name) => {
    if (!value) return true;
    const newPwd = document.getElementById('new_pwd');
    if (!newPwd) return true;
    return value === newPwd.value;
  });
}
