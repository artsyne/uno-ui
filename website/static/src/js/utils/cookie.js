/**
 * Cookie Utilities
 * 简单的 Cookie 操作工具函数，替代 Kiko/jQuery 依赖
 */

/**
 * 获取 Cookie 值
 * @param {string} name Cookie 名称
 * @returns {string|null} Cookie 值或 null
 */
export function getCookie(name) {
  const value = `; ${document.cookie}`;
  const parts = value.split(`; ${name}=`);
  if (parts.length === 2) return parts.pop().split(';').shift();
  return null;
}

/**
 * 设置 Cookie 值
 * @param {string} name Cookie 名称
 * @param {string} value Cookie 值
 * @param {number} days 有效期（天）
 * @param {string} path 路径，默认为 '/'
 */
export function setCookie(name, value, days = 365, path = '/') {
  let expires = '';
  if (days) {
    const date = new Date();
    date.setTime(date.getTime() + days * 24 * 60 * 60 * 1000);
    expires = `; expires=${date.toUTCString()}`;
  }
  document.cookie = `${name}=${value || ''}${expires}; path=${path}; SameSite=Lax`;
}
