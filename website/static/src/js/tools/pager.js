/**
 * jQuery页面行为管理插件
 *
 * 这个插件提供了一种简单而强大的方式管理页面中的交互行为和动态内容加载。
 * 使用现代浏览器的History API实现无刷新页面切换体验。
 * 支持多种行为类型，包括AJAX内容加载、对话框显示、事件触发和回调函数执行。
 *
 * @author nobody
 * @version 1.0.0
 * @license MIT
 *
 * 使用说明：
 *
 * 1. 基础初始化：
 *    $.pageAction.init({
 *      // 可选的全局配置
 *      loadingText: '加载中...',
 *      defaultSelector: '[data-action]',
 *      cache: true
 *    }, [
 *      // 行为配置数组
 *      {
 *        id: 'user-profile',       // 行为ID，必须唯一
 *        type: 'ajax',             // 行为类型
 *        title: '用户资料',         // 页面标题
 *        url: '/api/user/profile', // AJAX请求URL
 *        target: '#content-area'   // 内容加载目标
 *      },
 *      {
 *        id: 'settings',
 *        type: 'dialog',
 *        title: '系统设置',
 *        content: '设置对话框内容',
 *        width: 500,
 *        height: 400
 *      }
 *    ]);
 *
 * 2. HTML元素绑定：
 *    <!-- 带多个data属性的元素，这些属性会自动传递到行为处理中 -->
 *    <button data-action="user-profile" data-user-id="123" data-view="detailed">查看用户资料</button>
 *    <a data-action="product-detail" data-id="456" data-category="electronics">产品详情</a>
 *
 * 3. 手动触发行为：
 *    $.pageAction.trigger('user-profile', { userId: 123 });
 *
 * 4. 动态添加行为：
 *    $.pageAction.add('new-action', {
 *      type: 'ajax',
 *      url: '/api/new-content',
 *      target: '#content'
 *    });
 *
 * 5. 支持的行为类型：
 *    - ajax: 通过AJAX加载内容，data-*属性会作为请求参数传递
 *    - dialog: 显示对话框，data-*属性会合并到对话框配置中
 *    - event: 触发自定义事件，data-*属性会作为事件数据传递
 *    - callback: 执行回调函数，data-*属性会作为参数传递给回调
 *    - redirect: 页面重定向，data-*属性会作为查询参数添加到URL
 *
 * 6. 支持的配置项：
 *    - defaultSelector: 默认选择器，用于查找绑定行为的元素
 *    - loadingText: 加载提示文本
 *    - cache: 是否缓存AJAX内容
 *    - errorHandler: 全局错误处理函数
 *    - lazyLoad: 延迟加载配置函数
 *    - baseUrl: 页面基础URL，用于构建完整URL
 *    - animation: 内容切换动画效果
 *    - animationDuration: 动画持续时间
 *    - scrollToTop: 加载新内容后是否滚动到顶部
 *    - additionalSelectors: 额外的选择器列表
 *    - autoInit: 是否自动初始化事件监听
 */

(function ($) {
  // 私有变量
  let actions = {}; // 行为配置存储
  let config = {
    // 默认配置
    defaultSelector: '[data-action]', // 默认查找带data-action属性的元素
    additionalSelectors: [], // 额外的选择器列表
    loadingText: '加载中...', // 加载提示文本
    cache: true, // 是否缓存AJAX内容
    errorHandler: null, // 全局错误处理函数
    lazyLoad: null, // 延迟加载配置函数
    baseUrl: '', // 基础URL
    animation: 'fade', // 内容切换动画效果: 'none', 'fade', 'slide'
    animationDuration: 300, // 动画持续时间(毫秒)
    scrollToTop: true, // 加载新内容后是否滚动到顶部
    autoInit: true, // 是否自动初始化事件监听
    updatePageTitle: true, // 是否自动更新页面标题
    debug: false, // 是否启用调试模式
  };

  // 缓存容器
  const contentCache = {};

  // 历史状态管理
  let currentState = null;

  // 检查浏览器是否支持history API
  const historySupported = window.history && window.history.pushState;
  if (!historySupported) {
    console.warn('PageAction: 当前浏览器不支持History API，部分功能将不可用');
  }

  // 私有方法：日志输出
  const log = function (message, data) {
    if (config.debug) {
      console.log(`[PageAction] ${message}`, data || '');
    }
  };

  // 私有方法：初始化历史状态监听
  const initHistoryListener = function () {
    if (!historySupported) return;

    log('初始化历史状态监听');

    // 监听popstate事件(浏览器前进后退) - jQuery 4语法
    $(window).on('popstate.pageAction', function (e) {
      const state = e.originalEvent.state;
      if (state && state._pageActionId) {
        log('检测到历史状态变化', state);
        // 从历史记录中恢复状态
        executeAction(state._pageActionId, state.data || {}, true);
      }
    });

    // 处理初始URL中的状态
    if (window.history.state && window.history.state._pageActionId) {
      const state = window.history.state;
      log('处理初始URL状态', state);
      setTimeout(function () {
        executeAction(state._pageActionId, state.data || {}, true);
      }, 100);
    }
  };

  // 私有方法：更新历史状态
  const updateHistory = function (actionId, data, title) {
    if (!historySupported) return;

    title = title || document.title;

    // 准备状态对象
    currentState = {
      _pageActionId: actionId,
      data: data || {},
    };

    // 构建URL
    let url = config.baseUrl ? config.baseUrl + '/' + actionId : '/' + actionId;

    // 如果有data.data且包含参数，将它们添加为查询字符串
    if (data && data.data && typeof data.data === 'object') {
      const params = new URLSearchParams();
      for (const key in data.data) {
        if (data.data.hasOwnProperty(key)) {
          params.append(key, data.data[key]);
        }
      }
      const paramString = params.toString();
      if (paramString) {
        url += '?' + paramString;
      }
    }

    log('更新历史状态', { actionId, title, url });

    try {
      // 更新文档标题
      if (config.updatePageTitle) {
        document.title = title;
      }

      // 更新URL并添加到历史记录
      window.history.pushState(currentState, title, url);
    } catch (e) {
      handleError('history_error', {
        message: '无法更新浏览器历史记录',
        error: e,
        actionId: actionId,
      });
    }
  };

  // 私有方法：从元素收集data属性
  const collectDataAttributes = function ($element) {
    const elementData = {};

    // 使用jQuery的data方法获取所有data-*属性
    const allData = $element.data();

    // 手动删除action属性，因为我们不想把它作为参数传递
    delete allData.action;

    // 合并获取的数据
    $.extend(elementData, allData);

    log('收集的数据属性', elementData);
    return elementData;
  };

  // 私有方法：处理元素点击
  const handleElementClick = function (e) {
    const $this = $(this);
    let actionId = $this.data('action');

    // 如果没有data-action属性，则不处理
    if (!actionId) return true;

    log('元素点击', { element: this, actionId });

    // 收集元素上的数据属性
    const elementData = collectDataAttributes($this);

    // 阻止默认行为
    e.preventDefault();

    // 执行行为
    executeAction(actionId, {
      element: $this,
      event: e,
      data: elementData,
    });

    return false;
  };

  // 私有方法：执行行为
  const executeAction = function (actionId, extraData, skipHistory) {
    log('执行行为', { actionId, extraData, skipHistory });

    // 检查行为是否存在
    if (!actions[actionId]) {
      // 尝试延迟加载
      if (config.lazyLoad && typeof config.lazyLoad === 'function') {
        log('尝试延迟加载行为配置', actionId);
        const promise = config.lazyLoad(actionId);
        if (promise && typeof promise.then === 'function') {
          promise
            .then(function (actionConfig) {
              if (actionConfig) {
                log('延迟加载成功', actionConfig);
                actions[actionId] = actionConfig;
                processAction(actionId, extraData, skipHistory);
              } else {
                handleError('action_not_found', {
                  message: `未找到ID为"${actionId}"的行为配置`,
                  actionId: actionId,
                });
              }
            })
            .catch(function (error) {
              handleError('lazy_load_error', {
                message: '延迟加载配置失败',
                error: error,
                actionId: actionId,
              });
            });
        } else {
          handleError('action_not_found', {
            message: `未找到ID为"${actionId}"的行为配置`,
            actionId: actionId,
          });
        }
        return;
      }

      handleError('action_not_found', {
        message: `未找到ID为"${actionId}"的行为配置`,
        actionId: actionId,
      });
      return;
    }

    processAction(actionId, extraData, skipHistory);
  };

  // 私有方法：处理行为
  const processAction = function (actionId, extraData, skipHistory) {
    // 获取行为配置
    const actionConfig = actions[actionId];

    log('处理行为', { actionId, actionConfig });

    // 合并额外数据
    const data = $.extend({}, actionConfig.data || {}, extraData || {});

    // 执行前回调
    if (typeof actionConfig.before === 'function') {
      try {
        log('执行before回调');
        const result = actionConfig.before(actionId, data);
        // 允许中断执行链
        if (result === false) {
          log('执行链被中断');
          return;
        }
      } catch (error) {
        handleError('before_callback_error', {
          message: '执行before回调时发生错误',
          error: error,
          actionId: actionId,
        });
        return;
      }
    }

    try {
      // 根据类型执行对应操作
      switch (actionConfig.type) {
        case 'ajax':
          if (!actionConfig.url || !actionConfig.target) {
            throw new Error('ajax类型行为必须指定url和target');
          }

          loadContent(
            actionConfig.url,
            actionConfig.target,
            data,
            function (success, response) {
              if (success && !skipHistory) {
                updateHistory(actionId, data, actionConfig.title);
              }
              if (typeof actionConfig.callback === 'function') {
                actionConfig.callback(success, response, data);
              }
            },
            actionConfig.animation || config.animation,
            actionConfig.animationDuration || config.animationDuration,
          );
          break;

        case 'event':
          if (!actionConfig.event) {
            throw new Error('event类型行为必须指定event');
          }

          triggerEvent(actionConfig.event, data);
          if (!skipHistory) {
            updateHistory(actionId, data, actionConfig.title);
          }
          if (typeof actionConfig.callback === 'function') {
            actionConfig.callback(true, null, data);
          }
          break;

        case 'callback':
          if (typeof actionConfig.callback !== 'function') {
            throw new Error('callback类型行为必须指定callback函数');
          }

          actionConfig.callback(actionId, data);
          if (!skipHistory) {
            updateHistory(actionId, data, actionConfig.title);
          }
          break;

        case 'dialog':
          showDialog(actionConfig, data);
          if (!skipHistory) {
            updateHistory(actionId, data, actionConfig.title);
          }
          break;

        case 'redirect':
          if (!actionConfig.url) {
            throw new Error('redirect类型行为必须指定url');
          }

          // 准备重定向URL
          let redirectUrl = actionConfig.url;

          // 如果有data.data且包含参数，将它们添加为查询字符串
          if (data && data.data && typeof data.data === 'object') {
            const params = new URLSearchParams();
            for (const key in data.data) {
              if (data.data.hasOwnProperty(key)) {
                params.append(key, data.data[key]);
              }
            }
            const paramString = params.toString();
            if (paramString) {
              redirectUrl +=
                (redirectUrl.includes('?') ? '&' : '?') + paramString;
            }
          }

          log('重定向到', redirectUrl);
          window.location.href = redirectUrl;
          break;

        default:
          throw new Error(`未知的行为类型: ${actionConfig.type}`);
      }

      // 执行后回调
      if (typeof actionConfig.after === 'function') {
        setTimeout(function () {
          try {
            log('执行after回调');
            actionConfig.after(actionId, data);
          } catch (error) {
            handleError('after_callback_error', {
              message: '执行after回调时发生错误',
              error: error,
              actionId: actionId,
            });
          }
        }, 0);
      }
    } catch (error) {
      handleError('action_execution_error', {
        message: '执行行为时发生错误',
        error: error,
        actionId: actionId,
        config: actionConfig,
      });
    }
  };

  // 私有方法：加载内容
  const loadContent = function (
    url,
    target,
    data,
    callback,
    animation,
    duration,
  ) {
    const $target = $(target);

    if (!$target.length) {
      handleError('target_not_found', {
        message: '目标元素不存在',
        target: target,
        url: url,
      });
      if (typeof callback === 'function') callback(false, null);
      return;
    }

    // 显示加载状态
    $target.html(`<div class="text-center p-4">${config.loadingText}</div>`);

    // 构建缓存键
    const cacheKey = url + (data ? JSON.stringify(data) : '');

    // 准备加载内容的函数
    const renderContent = function (content) {
      // 应用动画效果
      switch (animation) {
        case 'fade':
          // jQuery 4风格的动画链
          $target
            .fadeOut({ duration: duration / 2 })
            .promise()
            .then(function () {
              $target.html(content);
              return $target.fadeIn({ duration: duration / 2 }).promise();
            })
            .then(function () {
              if (config.scrollToTop) {
                $('html, body').animate(
                  { scrollTop: 0 },
                  { duration: duration / 2 },
                );
              }
            });
          break;
        case 'slide':
          // jQuery 4风格的动画链
          $target
            .slideUp({ duration: duration / 2 })
            .promise()
            .then(function () {
              $target.html(content);
              return $target.slideDown({ duration: duration / 2 }).promise();
            })
            .then(function () {
              if (config.scrollToTop) {
                $('html, body').animate(
                  { scrollTop: 0 },
                  { duration: duration / 2 },
                );
              }
            });
          break;
        default: // 'none'或其他值
          $target.html(content);
          if (config.scrollToTop) {
            $('html, body').animate(
              { scrollTop: 0 },
              { duration: duration / 2 },
            );
          }
          break;
      }
    };

    // 检查缓存
    if (config.cache && contentCache[cacheKey]) {
      log('使用缓存内容', url);
      renderContent(contentCache[cacheKey]);
      if (typeof callback === 'function')
        callback(true, contentCache[cacheKey]);
      return;
    }

    // 准备AJAX参数
    const ajaxParams = {
      url: url,
      method: 'GET',
      success: function (response) {
        log('AJAX加载成功');
        if (config.cache) {
          contentCache[cacheKey] = response;
        }
        renderContent(response);
        if (typeof callback === 'function') callback(true, response);
      },
      error: function (xhr, status, error) {
        log('AJAX加载失败', { status, error });
        const errorMsg = `<div class="text-red-500 p-4">加载失败: ${error}</div>`;
        $target.html(errorMsg);

        handleError('ajax_error', {
          message: '加载内容失败',
          error: error,
          status: status,
          xhr: xhr,
          url: url,
        });

        if (typeof callback === 'function') callback(false, errorMsg);
      },
    };

    // 如果有数据，添加到请求中
    if (data && data.data) {
      if (typeof data.data === 'object' && Object.keys(data.data).length > 0) {
        ajaxParams.data = data.data;
        // 如果数据比较多或包含复杂结构，使用POST方法
        if (
          Object.keys(data.data).length > 5 ||
          JSON.stringify(data.data).length > 100
        ) {
          ajaxParams.method = 'POST';
        }
      }
    }

    log('AJAX请求参数', ajaxParams);

    // 发送AJAX请求
    $.ajax(ajaxParams);
  };

  // 私有方法：触发事件
  const triggerEvent = function (eventName, data) {
    log('触发事件', { eventName, data });
    $(document).trigger(eventName, [data]); // jQuery 4风格的事件触发，参数作为数组传递
  };

  // 私有方法：显示对话框
  const showDialog = function (dialogConfig, extraData) {
    // 准备对话框配置
    const finalConfig = {};

    // 先复制原始配置
    $.extend(finalConfig, dialogConfig);

    // 合并来自元素的data-*属性，仅对特定属性进行合并
    if (extraData && extraData.data) {
      // 常用的对话框属性，可以从data-*中获取
      const dialogProps = ['width', 'height', 'title', 'content', 'footerTxt'];

      for (const key in extraData.data) {
        if (dialogProps.includes(key)) {
          finalConfig[key] = extraData.data[key];
        }
      }
    }

    log('显示对话框', finalConfig);

    try {
      $.box.dialog(finalConfig);
    } catch (error) {
      handleError('dialog_error', {
        message: '显示对话框失败',
        error: error,
        config: dialogConfig,
      });
    }
  };

  // 私有方法：错误处理
  const handleError = function (type, errorData) {
    console.error(`[PageAction Error] ${type}:`, errorData.message, errorData);

    // 调用全局错误处理函数
    if (typeof config.errorHandler === 'function') {
      try {
        config.errorHandler(type, errorData);
      } catch (error) {
        console.error('[PageAction Error] 执行错误处理时发生异常:', error);
      }
    }
  };

  // 私有方法：绑定选择器
  const bindSelectors = function () {
    log('绑定选择器');

    // 使用jQuery 4风格的事件委托
    // 移除旧的绑定
    $(document).off('click.pageAction');

    // 绑定默认选择器
    $(document).on(
      'click.pageAction',
      config.defaultSelector,
      handleElementClick,
    );

    // 绑定额外的选择器
    for (let i = 0; i < config.additionalSelectors.length; i++) {
      $(document).on(
        'click.pageAction',
        config.additionalSelectors[i],
        handleElementClick,
      );
    }

    // 绑定动态配置的选择器
    for (let id in actions) {
      if (actions[id].selector) {
        $(document).on('click.pageAction', actions[id].selector, function (e) {
          e.preventDefault();
          // 收集data属性
          const elementData = collectDataAttributes($(this));
          executeAction(id, {
            element: $(this),
            event: e,
            data: elementData,
          });
          return false;
        });
      }
    }
  };

  // 公共方法
  const methods = {
    /**
     * 初始化插件
     * @param {Object} options - 全局配置选项
     * @param {Array} actionsArray - 行为配置数组
     */
    init: function (options, actionsArray) {
      // 处理参数
      if (Array.isArray(options) && actionsArray === undefined) {
        // 兼容旧版本调用方式，只传入了行为数组
        actionsArray = options;
        options = {};
      }

      // 更新全局配置
      if (options && typeof options === 'object') {
        config = $.extend({}, config, options);
      }

      log('初始化PageAction', { config, actionsArray });

      if (!Array.isArray(actionsArray)) {
        console.error('[PageAction Error] 初始化需要一个行为配置数组');
        return this;
      }

      // 重置行为存储
      actions = {};

      // 处理每个行为配置
      for (let i = 0; i < actionsArray.length; i++) {
        const action = actionsArray[i];

        if (!action.id) {
          console.error('[PageAction Error] 行为配置必须包含id字段:', action);
          continue;
        }

        // 存储行为配置
        actions[action.id] = action;
      }

      // 如果启用自动初始化
      if (config.autoInit) {
        // 绑定选择器
        bindSelectors();

        // 初始化历史监听
        initHistoryListener();
      }

      return this;
    },

    /**
     * 添加单个行为配置
     * @param {String} id - 行为ID
     * @param {Object} actionConfig - 行为配置
     */
    add: function (id, actionConfig) {
      if (!id || typeof id !== 'string') {
        console.error('[PageAction Error] add方法的第一个参数必须是字符串ID');
        return this;
      }

      log('添加行为', { id, actionConfig });
      actions[id] = actionConfig;

      // 如果有选择器，重新绑定
      if (actionConfig.selector) {
        $(document).off('click.pageAction', actionConfig.selector);
        $(document).on('click.pageAction', actionConfig.selector, function (e) {
          e.preventDefault();
          // 收集data属性
          const elementData = collectDataAttributes($(this));
          executeAction(id, {
            element: $(this),
            event: e,
            data: elementData,
          });
          return false;
        });
      }

      return this;
    },

    /**
     * 移除行为配置
     * @param {String} id - 行为ID
     */
    remove: function (id) {
      if (actions[id]) {
        log('移除行为', id);

        // 如果有选择器，解除绑定
        if (actions[id].selector) {
          $(document).off('click.pageAction', actions[id].selector);
        }

        delete actions[id];
      }
      return this;
    },

    /**
     * 设置全局配置
     * @param {Object} options - 配置选项
     */
    configure: function (options) {
      log('更新配置', options);
      config = $.extend(config, options);

      // 重新绑定选择器
      if (options.defaultSelector || options.additionalSelectors) {
        bindSelectors();
      }

      return this;
    },

    /**
     * 清除内容缓存
     * @param {String} url - 可选的URL，不指定则清除所有缓存
     */
    clearCache: function (url) {
      if (url) {
        log('清除URL缓存', url);
        // 查找并清除匹配的缓存
        for (let key in contentCache) {
          if (key.startsWith(url)) {
            delete contentCache[key];
          }
        }
      } else {
        log('清除所有缓存');
        // 清除所有缓存
        for (let key in contentCache) {
          delete contentCache[key];
        }
      }
      return this;
    },

    /**
     * 手动触发行为
     * @param {String} id - 行为ID
     * @param {Object} data - 可选的额外数据
     */
    trigger: function (id, data) {
      log('手动触发行为', { id, data });
      // 如果data不是一个对象，则包装成一个对象
      const wrappedData = data
        ? typeof data === 'object'
          ? { data: data }
          : { data: { value: data } }
        : {};
      executeAction(id, wrappedData);
      return this;
    },

    /**
     * 获取当前所有行为配置
     * @return {Object} 行为配置对象
     */
    getActions: function () {
      return $.extend({}, actions);
    },

    /**
     * 注销插件，清理所有事件监听
     */
    destroy: function () {
      log('注销插件');
      // 移除所有事件监听
      $(document).off('click.pageAction');
      $(window).off('popstate.pageAction');
      // 清空缓存
      for (let key in contentCache) {
        delete contentCache[key];
      }
      // 清空行为配置
      actions = {};
      return this;
    },
  };

  // 导出到jQuery命名空间
  $.pageAction = {
    init: methods.init,
    add: methods.add,
    remove: methods.remove,
    configure: methods.configure,
    clearCache: methods.clearCache,
    trigger: methods.trigger,
    getActions: methods.getActions,
    destroy: methods.destroy,
  };
})(jQuery);
