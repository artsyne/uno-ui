/**
 * 主应用入口文件
 *
 * 使用Kiko UI的页面行为管理组件统一管理所有页面行为
 *
 * @version 2.3.0
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
  let errorMessage = appI18n.getMessage(fallbackMessage);
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

$(() => {
  const AJAX_URL = '';

  // 保留对全局handleError函数的引用
  const handleError = window.handleError;

  log('初始化应用...');

  // 设置全局AJAX拦截器 - 简化版
  $.ajaxSetup({
    beforeSend: (xhr, settings) => {
      // 显示加载指示器
      kiko.ui.notify.loading(true);

      // 只处理POST请求，寻找当前活动表单并禁用
      if (settings.type === 'POST') {
        const $focusedForm = $('form:has(:focus)');
        if ($focusedForm.length) {
          $focusedForm.addClass('form-processing').disableFormInputs(true);
        }
      }
    },

    error: (xhr, textStatus, errorThrown) => {
      // 当没有自定义错误处理时，显示默认错误消息
      if (!xhr.isHandled) {
        // 使用完整的错误处理函数，默认显示通知
        const result = handleError(xhr, 'operation.fail');
        xhr.isHandled = result.handled;
      }
    },

    complete: () => {
      // 隐藏加载指示器
      kiko.ui.notify.loading(false);
      // 恢复所有处于处理中的表单
      $('form.form-processing')
        .removeClass('form-processing')
        .disableFormInputs(false);
    },
  });

  // 初始化认证模块
  const authActions = Auth.init(AJAX_URL);

  // 初始化日记本模块
  const diaryActions = Diary.init(AJAX_URL);

  // 初始化设置模块
  const settingsActions = Settings.init(AJAX_URL);

  // 创建统一的页面行为管理器
  kiko.ui.pager.init(
    {
      debug: true,
      cache: false,
      enableAnchorCapture: true,
      loadingText: appI18n.getMessage('loading'),
      errorHandler: (error) => handleError(error),
      init: (action) => {
        log('行为初始化', action.id);
        return true; // 继续执行
      },
      complete: (success, action) => {
        log('行为执行完成', { success, actionId: action.id });

        // 成功完成的特定类型行为可以显示成功消息
        if (success && action.showSuccessMessage) {
          kiko.ui.notify.success(
            action.successMessage || appI18n.getMessage('operation.success'),
          );
        }
      },
    },
    [
      // 合并所有行为配置
      ...authActions,
      ...diaryActions,
      ...settingsActions,
    ],
  );

  // 初始化模块方法
  const initActionModules = (actions, moduleName) => {
    actions.forEach((action) => {
      if (action.init && typeof action.init === 'function') {
        try {
          action.init();
        } catch (e) {
          log(`错误: 执行${moduleName}初始化方法失败`, e);
          handleError(e, 'execution.failed');
        }
      }
    });
  };

  // 调用各模块的初始化方法
  log('调用初始化方法...');
  initActionModules(authActions, '认证模块');
  initActionModules(diaryActions, '日记本模块');
  initActionModules(settingsActions, '设置模块');

  // 初始化滚动监听
  initScrollHeader();

  log('应用初始化完成');
});

/**
 * 滚动监听初始化 - 处理顶部导航栏
 * 当页面滚动时，将页面导航内容复制到顶部导航栏
 */
function initScrollHeader() {
  // 检查所需元素是否存在
  const $topNav = $('#top-nav');
  const $pageNav = $('#page-nav');

  if (!$topNav.length || !$pageNav.length) return;

  // 缓存页面导航内容
  const pageNavHtml = $pageNav.html();

  // 绑定滚动事件，使用节流函数优化性能
  let lastScrollY = window.scrollY;
  $(window).on('scroll', () => {
    // 只有滚动位置发生变化时才更新
    if (window.scrollY !== lastScrollY) {
      lastScrollY = window.scrollY;
      $topNav.html(lastScrollY > 1 ? pageNavHtml : '');
    }
  });
}
