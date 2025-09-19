/**
 * 通用工具函数库
 * 
 * 提供全站通用的工具函数，包括cookie操作、HTML转义、时间格式化、错误处理等
 * 
 * @version 1.0.0
 * @license MIT
 */

/**
 * 全局日志函数
 *
 * @param {string} message - 日志消息
 * @param {...any} data - 可选的额外数据
 */
window.log = (message, ...data) => kiko.core.log(`[App] ${message}`, ...data);

/**
 * Cookie 工具函数
 * 统一的cookie操作方法，用于替代localStorage
 */

/**
 * 获取cookie值
 * @param {string} name - cookie名称
 * @returns {string|null} cookie值，不存在时返回null
 */
window.getCookie = (name) => {
  const value = `; ${document.cookie}`;
  const parts = value.split(`; ${name}=`);
  if (parts.length === 2) {
    return decodeURIComponent(parts.pop().split(';').shift());
  }
  return null;
};

/**
 * 设置cookie
 * @param {string} name - cookie名称
 * @param {string} value - cookie值
 * @param {Object} options - 选项
 * @param {number} [options.days=365] - 过期天数，默认1年
 * @param {string} [options.path='/'] - 路径，默认根路径
 * @param {string} [options.domain] - 域名
 * @param {boolean} [options.secure] - 是否仅HTTPS
 * @param {string} [options.sameSite='Lax'] - SameSite属性
 */
window.setCookie = (name, value, options = {}) => {
  const {
    days = 365,
    path = '/',
    domain,
    secure,
    sameSite = 'Lax'
  } = options;
  
  let cookieString = `${name}=${encodeURIComponent(value)}`;
  
  if (days) {
    const date = new Date();
    date.setTime(date.getTime() + (days * 24 * 60 * 60 * 1000));
    cookieString += `; expires=${date.toUTCString()}`;
  }
  
  if (path) cookieString += `; path=${path}`;
  if (domain) cookieString += `; domain=${domain}`;
  if (secure) cookieString += `; secure`;
  if (sameSite) cookieString += `; samesite=${sameSite}`;
  
  document.cookie = cookieString;
};

/**
 * 删除cookie
 * @param {string} name - cookie名称
 * @param {string} [path='/'] - 路径，默认根路径
 * @param {string} [domain] - 域名
 */
window.removeCookie = (name, path = '/', domain) => {
  let cookieString = `${name}=; expires=Thu, 01 Jan 1970 00:00:00 UTC`;
  if (path) cookieString += `; path=${path}`;
  if (domain) cookieString += `; domain=${domain}`;
  document.cookie = cookieString;
};

/**
 * HTML转义工具函数
 * 将文本中的HTML特殊字符转义，防止XSS攻击
 *
 * @param {string} text - 需要转义的文本
 * @returns {string} 转义后的安全文本
 */
window.escapeHtml = (text) => {
  if (!text) return '';

  const div = document.createElement('div');
  div.textContent = text;
  return div.innerHTML;
};

/**
 * 时间格式化工具函数
 * 将时间戳格式化为 YYYY-MM-DD HH:mm 格式
 *
 * @param {number} timestamp - Unix时间戳（秒或毫秒）
 * @param {string} [timezone] - 时区标识符（如 'Asia/Shanghai', 'America/New_York'）
 * @returns {string} 格式化后的时间字符串，格式为 YYYY-MM-DD HH:mm
 */
window.formatTime = (timestamp, timezone) => {
  if (!timestamp) return '';

  // 处理时间戳格式（支持秒和毫秒）
  const date = new Date(timestamp < 10000000000 ? timestamp * 1000 : timestamp);

  // 如果传入了时区，使用指定时区；否则使用浏览器本地时区
  const options = {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    hour12: false,
  };

  if (timezone) {
    options.timeZone = timezone;
  }

  try {
    // 使用 Intl.DateTimeFormat 进行格式化
    const formatter = new Intl.DateTimeFormat('sv-SE', options);
    return formatter.format(date);
  } catch (error) {
    // 如果时区无效，回退到本地时区
    const fallbackFormatter = new Intl.DateTimeFormat('sv-SE', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
      hour12: false,
    });
    return fallbackFormatter.format(date);
  }
};

/**
 * 统一的错误处理函数
 * 从错误对象或响应中提取消息并显示，或者根据错误码返回自定义消息
 * 已更新以适配新的ResponseService格式：{code, message, data}
 *
 * @param {Object|string} error - 错误对象、消息或XHR
 * @param {string} [fallbackMessage] - 默认错误消息键
 * @param {Object} [options] - 额外选项
 * @param {boolean} [options.showNotify=true] - 是否显示通知
 * @param {Object} [options.errorCodeMessages] - 错误代码到消息的映射
 * @param {Object} [options.messageData] - 用于替换消息中占位符的数据
 * @returns {Object} 包含处理结果的对象 {handled, message, isValid}
 */
window.handleError = (
  error,
  fallbackMessage = 'operation.fail',
  options = {},
) => {
  const {
    showNotify = true,
    errorCodeMessages = {},
    messageData = {},
  } = options;

  // 错误码（来自服务端 AppError.Code）
  let errorCode = null;

  // 确定错误消息文本
  let errorMessage = kiko.i18n.getMessage(fallbackMessage);
  let isUserError = false; // 是否为用户输入导致的预期错误

  // 解析错误对象，提取错误代码和消息
  // 适配新的ResponseService格式：{code, message, data}
  if (error?.responseJSON) {
    errorCode = error.responseJSON.code;

    // 优先使用映射的自定义错误消息
    if (errorCode && errorCodeMessages[errorCode]) {
      errorMessage = errorCodeMessages[errorCode];
      isUserError = true;
    }
    // 否则使用服务端提供的错误描述
    else if (error.responseJSON.message) {
      errorMessage = error.responseJSON.message;
    }
  } else if (error?.message) {
    errorMessage = error.message;
  } else if (typeof error === 'string' && error) {
    errorMessage = error;
  }

  // 替换消息中的占位符
  Object.keys(messageData).forEach((key) => {
    errorMessage = errorMessage.replace(`{${key}}`, messageData[key]);
  });

  // 显示错误消息（如果需要）
  if (showNotify) {
    kiko.ui.notify.error(errorMessage);
  }

  return {
    handled: true,
    message: errorMessage,
    code: errorCode,
    // 对于表单验证，返回是否有效（没有错误为有效）
    isValid: !isUserError,
  };
};