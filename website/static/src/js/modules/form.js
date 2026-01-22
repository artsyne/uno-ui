/**
 * Global Form Handler Module
 * Handles anti-double-click and loading states for all forms
 */

/**
 * Manually toggle form loading state
 * @param {HTMLFormElement} form
 * @param {boolean} isLoading
 */
export function setFormLoading(form, isLoading) {
  if (!form || form.tagName !== 'FORM') return;

  const submitBtn = form.querySelector('button[type="submit"]');
  if (!submitBtn) return;

  if (isLoading) {
    if (!submitBtn.disabled) {
      submitBtn.dataset.originalContent = submitBtn.innerHTML;
      submitBtn.disabled = true;
      // Add loading spinner
      const loadingIcon = '<span class="mdi mdi-loading mdi-spin mr-2"></span>';
      if (!submitBtn.innerHTML.includes('mdi-loading')) {
        submitBtn.innerHTML = loadingIcon + submitBtn.innerHTML;
      }
    }
  } else {
    if (submitBtn.disabled) {
      submitBtn.disabled = false;
      // Restore original content
      if (submitBtn.dataset.originalContent) {
        submitBtn.innerHTML = submitBtn.dataset.originalContent;
        delete submitBtn.dataset.originalContent;
      } else {
        // Fallback cleanup if dataset is missing
        submitBtn.innerHTML = submitBtn.innerHTML.replace(
          '<span class="mdi mdi-loading mdi-spin mr-2"></span>',
          '',
        );
      }
    }
  }
}

export function initFormHandlers() {
  // Anti-double-click: Disable submit button on request start
  document.body.addEventListener('htmx:beforeRequest', function (evt) {
    const form = evt.detail.elt;
    setFormLoading(form, true);
  });

  // Restore button state on request end (success, error, timeout)
  document.body.addEventListener('htmx:afterRequest', function (evt) {
    const form = evt.detail.elt;
    setFormLoading(form, false);
  });
}
