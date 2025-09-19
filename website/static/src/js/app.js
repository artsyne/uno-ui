/**
 * 主应用入口文件
 *
 * 使用Kiko UI的页面行为管理组件统一管理所有页面行为
 * 通用工具函数已移至 utils.js
 *
 * @version 2.4.0
 * @license MIT
 */

$(() => {
  const AJAX_URL = '';

  // 保留对全局handleError函数的引用
  const handleError = window.handleError;

  log('初始化应用...');

  // 设置全局AJAX拦截器 - 简化版
  $.ajaxSetup({
    beforeSend: (xhr, settings) => {
      // 显示加载指示器
      kiko.ui.notify.loading(true, 'progressbar');

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

  // 初始化社交功能模块
  const socialActions = Social.init(AJAX_URL);

  // 初始化设置模块
  const settingsActions = Settings.init(AJAX_URL);

  // 创建统一的页面行为管理器
  kiko.ui.pager.init(
    {
      debug: true,
      cache: false,
      enableAnchorCapture: true,
      loadingText: kiko.i18n.getMessage('loading'),
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
            action.successMessage || kiko.i18n.getMessage('operation.success'),
          );
        }
      },
    },
    [
      // 合并所有行为配置
      ...authActions,
      ...diaryActions,
      ...socialActions,
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
  initActionModules(socialActions, '社交功能模块');
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

      // 重新初始化新插入的tooltip元素
      if (lastScrollY > 1 && kiko.ui && kiko.ui.tooltip) {
        kiko.ui.tooltip.initElement(
          $topNav.find(
            '[data-tooltip], [data-tooltip-content], [data-tooltip-target]',
          ),
        );
      }
    }
  });
}
