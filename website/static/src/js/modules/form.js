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

  // 如果表单有 no-loading-mask 类，跳过
  if (form.classList.contains('no-loading-mask')) return;

  // 查找或创建 mask
  let mask = form.querySelector('.form-mask');

  if (!mask && isLoading) {
    // 动态创建 mask（仅首次加载时创建）
    mask = document.createElement('div');
    mask.className = 'form-mask';
    mask.innerHTML = '<div class="loading-spinner"></div>';
    form.insertBefore(mask, form.firstChild);
  }

  if (mask) {
    // 使用 opacity + visibility 实现平滑过渡（CSS 已设置默认 opacity:0, visibility:hidden）
    if (isLoading) {
      mask.style.visibility = 'visible';
      // 下一帧设置 opacity，触发过渡动画
      requestAnimationFrame(() => {
        mask.style.opacity = '1';
      });
    } else {
      mask.style.opacity = '0';
      // 等待过渡完成后隐藏
      setTimeout(() => {
        if (mask.style.opacity === '0') {
          mask.style.visibility = 'hidden';
        }
      }, 300);
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
