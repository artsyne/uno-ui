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

  // 1. Handle Mask (if present)
  // Look for mask inside the form first, then in the container
  let mask = form.querySelector('.form-mask');
  if (!mask) {
    mask = form.parentElement?.querySelector('.form-mask');
  }

  if (mask) {
    mask.style.display = isLoading ? 'flex' : 'none';

    // Mask Visual Feedback Policy:
    // We retain the mask display logic above (to show loading spinner and block clicks via z-index).
    // However, we REMOVED the logic that explicitly disables inputs/buttons (setting .disabled = true).
    // This is to avoid visual state changes on the buttons (like graying out or opacity changes)
    // while still preventing interaction through the full-screen mask.
    // The mask itself (z-9999) acts as the interaction blocker.
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

/**
 * Validates form and triggers 'validated' event if successful
 * @param {HTMLFormElement} form
 */
export async function validateAndTrigger(form) {
  const valid = window.getUnoValid();
  if (valid && (await valid.validateForm(form))) {
    form.dispatchEvent(new CustomEvent('validated', { bubbles: true }));
  }
}
