/**
 * 主应用入口文件
 *
 * 使用Kiko UI的页面行为管理组件统一管理所有页面行为
 *
 * @version 2.3.0
 * @license MIT
 */

$(() => {
  const AJAX_URL = '';

  /**
   * 日志函数
   *
   * @param {string} message - 日志消息
   * @param {...any} data - 可选的额外数据
   */
  const log = (message, ...data) => {
    console.log(`[App] ${message}`, ...data);
  };

  /**
   * 统一的错误处理函数
   * 从错误对象或响应中提取消息并显示
   *
   * @param {Object|string} error - 错误对象、消息或XHR
   * @param {string} [fallbackMessage] - 默认错误消息键
   * @returns {Object} 包含处理结果的对象
   */
  const handleError = (error, fallbackMessage = 'operation.fail') => {
    // 确定错误消息文本
    let errorMessage = appI18n.getMessage(fallbackMessage);

    // 处理XHR响应
    if (error && error.responseJSON && error.responseJSON.Message) {
      errorMessage = error.responseJSON.Message;
    }
    // 处理错误对象
    else if (error && error.message) {
      errorMessage = error.message;
    }
    // 处理字符串错误
    else if (typeof error === 'string' && error) {
      errorMessage = error;
    }

    // 显示错误消息
    kiko.ui.notify.error(errorMessage);

    return { handled: true, message: errorMessage };
  };

  log('初始化应用...');

  // 设置全局AJAX拦截器 - 简化版
  $.ajaxSetup({
    /**
     * 全局Ajax发送前回调
     * @param {Object} xhr - XMLHttpRequest对象
     * @param {Object} settings - Ajax设置对象
     */
    beforeSend: function (xhr, settings) {
      // 显示加载指示器
      kiko.ui.notify.loading(true);

      // 只处理POST请求，寻找当前活动表单
      if (settings.type === 'POST') {
        // 获取当前焦点所在表单
        const $focusedForm = $('form:has(:focus)');
        if ($focusedForm.length) {
          // 标记表单为处理中并禁用输入
          $focusedForm.addClass('form-processing').disableFormInputs(true);
        }
      }
    },

    /**
     * 全局Ajax错误处理
     * @param {Object} xhr - XMLHttpRequest对象
     * @param {string} textStatus - 错误状态描述
     * @param {string} errorThrown - 抛出的错误信息
     */
    error: function (xhr, textStatus, errorThrown) {
      // 当没有自定义错误处理时，显示默认错误消息
      if (!xhr.isHandled) {
        // 使用共享错误处理函数
        const result = handleError(
          xhr.responseJSON?.Message || errorThrown || xhr,
        );
        xhr.isHandled = result.handled;
      }
    },

    /**
     * 全局Ajax完成处理
     * @param {Object} xhr - XMLHttpRequest对象
     * @param {string} textStatus - 状态描述
     */
    complete: function (xhr, textStatus) {
      // 隐藏加载指示器
      kiko.ui.notify.loading(false);

      // 恢复所有处于处理中的表单（简化的兜底措施）
      $('form.form-processing')
        .removeClass('form-processing')
        .disableFormInputs(false);
    },
  });

  // 初始化认证模块
  const authActions = Auth.init(AJAX_URL);

  // 初始化日记本模块
  const diaryActions = Diary.init(AJAX_URL);

  // 创建统一的页面行为管理器
  kiko.ui.pager.init(
    {
      debug: true,
      cache: false,
      enableAnchorCapture: true,
      // 移除复杂的loadingText，使用简单文本，因为全局ajax加载状态已由$.ajaxSetup处理
      loadingText: appI18n.getMessage('loading'),
      errorHandler: function (error, type, data) {
        // 使用共享错误处理函数
        handleError(error);
      },
      init: function (action, eventData) {
        log('行为初始化', action.id);
        return true; // 继续执行
      },
      complete: function (success, action, eventData) {
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
    ],
  );

  log('调用初始化方法...');

  // 调用各模块的初始化方法（如果有）
  authActions.forEach((action) => {
    if (action.init && typeof action.init === 'function') {
      try {
        action.init();
      } catch (e) {
        log(`错误: 执行认证模块初始化方法失败`, e);
        handleError(e, 'execution.failed');
      }
    }
  });

  diaryActions.forEach((action) => {
    if (action.init && typeof action.init === 'function') {
      try {
        action.init();
      } catch (e) {
        log(`错误: 执行日记本模块初始化方法失败`, e);
        handleError(e, 'execution.failed');
      }
    }
  });

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
  if ($('#top-nav').length === 0 || $('#page-nav').length === 0) {
    return;
  }

  // 绑定滚动事件
  $(window).on('scroll', function () {
    if ($(window).scrollTop() > 1) {
      $('#top-nav').html($('#page-nav').html());
    } else if ($(window).scrollTop() === 0) {
      $('#top-nav').html('');
    }
  });
}
