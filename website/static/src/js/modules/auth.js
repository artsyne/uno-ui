export function initAuth() {
  // Configure HTMX to allow swapping on error status codes (4xx, 5xx)
  document.body.addEventListener('htmx:beforeSwap', function (evt) {
    const xhr = evt.detail.xhr;
    const contentType = xhr.getResponseHeader('Content-Type');

    // Handle JSON responses (e.g. validation errors from generic middleware)
    if (contentType && contentType.includes('application/json')) {
      evt.detail.shouldSwap = false;

      try {
        const result = JSON.parse(xhr.responseText);
        const notify = window.getUnoNotify();

        if (!result.success) {
          evt.detail.isError = true;
          let msg = result.message || 'common.op_fail';
          // Append details if it's a string
          if (result.details && typeof result.details === 'string') {
            msg += ': ' + result.details;
          }

          // Locate error container:
          // Priority 1: Check if the target itself is the error container or contains it
          let errorContainer = null;
          const target = evt.detail.target;

          if (target) {
            if (target.classList.contains('error-message')) {
              errorContainer = target;
            } else if (target.querySelector) {
              errorContainer = target.querySelector('.error-message');
            }
          }

          // Priority 2: Look inside the form (triggering element) if not found yet
          if (
            !errorContainer &&
            evt.detail.elt &&
            evt.detail.elt.querySelector
          ) {
            errorContainer = evt.detail.elt.querySelector('.error-message');
          }

          if (errorContainer) {
            // Translate if possible
            const i18n = window.Alpine?.store('i18n');
            const translatedMsg = i18n ? i18n.t(msg) : msg;

            errorContainer.textContent = translatedMsg;
            errorContainer.classList.remove('hidden');

            errorContainer.style.display = ''; // Ensure visibility
          } else {
            // Fallback: Log warning but DO NOT use notify as requested
            console.warn(
              'Error container .error-message not found. Message:',
              msg,
            );
          }
        } else {
          // Success JSON
          evt.detail.isError = false;
          notify.success(result.message || 'auth.success');
        }
      } catch (e) {
        evt.detail.isError = true;
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

  // Handle Login Success Redirect
  document.body.addEventListener('htmx:afterRequest', function (evt) {
    const target = evt.detail.elt;
    // Check if it is the login or register form submission
    const action = target.getAttribute('hx-post');
    if (action === '/login' || action === '/register') {
      if (evt.detail.successful) {
        try {
          const res = JSON.parse(evt.detail.xhr.responseText);
          if (res.success && res.data && res.data.UserId) {
            window.location.href = '/u/' + res.data.UserId;
          }
        } catch (e) {
          // JSON parse error or invalid structure, ignore
          console.error('Auth response processing error:', e);
        }
      }
    }
  });
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
  // Updated for Uno Valid v0.8.0+: Native Async Support
  validStore.registerValidator('checkUsername', async (value, name) => {
    if (!value) return true;

    try {
      const formData = new FormData();
      formData.append('user_name', value);

      const response = await fetch('/check-username', {
        method: 'POST',
        body: formData,
      });

      const result = await response.json();

      // If success is true, it means available
      if (result.success) {
        return true;
      }

      // Return error message directly
      const t = (key) => Alpine.store('i18n').t(key);
      // Map backend error codes to frontend translation keys if necessary
      let msgKey = result.message;
      if (msgKey === 'user_exists') {
        msgKey = 'auth.username_taken';
      }
      return t(msgKey) || t('auth.username_taken');
    } catch (error) {
      const t = (key) => Alpine.store('i18n').t(key);
      return t('common.op_fail');
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

      // Return error message directly
      const t = (key) => Alpine.store('i18n').t(key);
      let msgKey = result.message;
      if (msgKey === 'email_exists') {
        msgKey = 'auth.email_taken';
      }
      return t(msgKey) || t('auth.email_taken');
    } catch (error) {
      // On network error, we might want to fail gracefully or show error
      // For now, let's return true to avoid blocking user if network is flaky,
      // or return error string if strict. Original code returned true.
      return true;
    }
  });

  // Register matchPwd callback
  validStore.registerValidator('matchPwd', (value, name, rule) => {
    if (!value) return true;
    // Support rule.target selector if provided, otherwise default to #new_pwd
    // This makes it compatible with both old implicit behavior and new explicit config
    const targetSelector = rule?.target || '#new_pwd';
    // If target is ID without #, add it (backward compat for hardcoded #new_pwd logic)
    const targetEl =
      document.querySelector(targetSelector) ||
      document.getElementById('new_pwd');

    if (!targetEl) return true;
    return value === targetEl.value;
  });
}
