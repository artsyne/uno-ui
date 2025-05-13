/**
 * jQuery对话框插件
 *
 * 这个插件提供了一种简单而强大的方式创建自定义对话框、确认框和消息框。
 * 它支持多种内容加载方式，包括HTML字符串、AJAX和iframe。
 * 使用Tailwind CSS 4进行样式设计，无需额外的CSS文件。
 *
 * @author nobody
 * @version 2.2.0
 * @license MIT
 *
 * 使用说明：
 *
 * 1. 基础对话框：
 *    $.box.dialog({
 *      title: '标题',
 *      content: '内容',
 *      width: 400,
 *      height: 300
 *    });
 *
 * 2. 确认对话框：
 *    $.box.confirm({
 *      title: '确认',
 *      content: '确定要执行此操作吗？',
 *      onOk: function() { ... },
 *      onCancel: function() { ... }
 *    });
 *
 * 3. 消息对话框：
 *    $.box.message({
 *      title: '提示',
 *      content: '操作成功！',
 *      onOk: function() { ... }
 *    });
 *
 * 4. 支持的配置项：
 *    - title: 对话框标题
 *    - content: 字符串内容
 *    - targetId: 加载指定ID元素的内容
 *    - ajaxUrl: 通过AJAX加载内容的URL
 *    - ajaxCallback: AJAX加载完成后的回调
 *    - iframeUrl: 通过iframe加载的URL
 *    - modal: 是否显示遮罩层
 *    - width: 对话框宽度
 *    - height: 对话框高度
 *    - position: 位置 ('center', 'top', 'left', 'right', 'bottom' 或数组)
 *    - footerTxt: 底部文本内容
 *    - idPrefix: 自定义ID前缀，默认为'jqbox-'
 *    - onBeforeOpen: 打开前的回调
 *    - onOpen: 打开后的回调
 *    - onBeforeClose: 关闭前的回调
 *    - onClose: 关闭后的回调
 *    - onOk: 确认框确定按钮的回调
 *    - onCancel: 确认框取消按钮的回调
 *
 * 5. Tailwind CSS 4 样式：
 *    本插件使用Tailwind CSS 4进行样式定义，主要使用了以下类：
 *    - 布局类: fixed, flex, justify-center, items-center, space-x-4, z-50
 *    - 间距类: p-4, px-4, py-2
 *    - 尺寸类: w-full, h-full, h-5, w-5
 *    - 边框类: border-t, border-b, rounded, rounded-lg
 *    - 颜色类: bg-white, bg-blue-500, bg-gray-300, text-white, text-gray-700
 *    - 效果类: shadow-lg, overflow-hidden, hover:bg-blue-600, hover:bg-gray-400
 *    - 状态类: hidden
 *
 *    如需自定义样式，可以:
 *    1. 在初始化时修改HTML模板中的类名
 *    2. 使用CSS选择器覆盖默认样式
 *    3. 通过修改Tailwind配置扩展主题颜色和其他属性
 */

(function ($) {
  // 私有变量
  let options = {};
  let idPrefix = 'jqbox-';
  let uniqueId = '';

  // 生成一个随机ID后缀，确保多个实例不会冲突
  const generateUniqueId = function () {
    return Math.random().toString(36).substring(2, 10);
  };

  // 获取带前缀的元素ID
  const getId = function (elementName) {
    return idPrefix + elementName + uniqueId;
  };

  // 获取元素选择器
  const getSelector = function (elementName) {
    return '#' + getId(elementName);
  };

  // 私有方法：初始化对话框
  const initBox = function (opts) {
    // 设置选项
    options = $.extend(
      {
        title: '提示信息',
        content: '暂无内容',
        modal: true,
        position: 'center',
        width: 400,
        height: 300,
        idPrefix: 'jqbox-',
      },
      opts || {},
    );

    // 如果提供了自定义前缀，则使用它
    if (options.idPrefix) {
      idPrefix = options.idPrefix;
    }

    // 为当前实例生成一个唯一ID
    if (!uniqueId) {
      uniqueId = generateUniqueId();
    }

    // 创建对话框DOM（如果不存在）
    if ($(getSelector('container')).length === 0) {
      const boxHtml = `
        <div id="${getId('container')}" class="hidden fixed z-50 bg-white rounded-lg shadow-lg overflow-hidden">
          <div id="${getId('header')}" class="flex justify-between items-center p-4 border-b">
            <div id="${getId('title')}" class="text-lg font-semibold"></div>
            <button id="${getId('close')}" class="text-gray-500 hover:text-gray-700" title="关闭">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>
          <div id="${getId('wrapper')}" class="p-4">
            <div id="${getId('content')}" class="overflow-auto"></div>
          </div>
          <div id="${getId('footer')}" class="p-4 border-t hidden"></div>
        </div>
        <div id="${getId('overlay')}" class="hidden fixed inset-0 bg-black bg-opacity-50 z-40"></div>
      `;
      $(boxHtml).appendTo('body');
    }

    // 设置事件和内容
    $(getSelector('close')).off('click').on('click', methods.close);
    $(getSelector('title')).text(options.title);

    // 设置底部文本
    if (typeof options.footerTxt === 'string') {
      $(getSelector('footer')).html(options.footerTxt).removeClass('hidden');
    } else {
      $(getSelector('footer')).addClass('hidden');
    }

    // 设置模态遮罩
    if (options.modal) {
      $(getSelector('overlay')).removeClass('hidden');
    } else {
      $(getSelector('overlay')).addClass('hidden');
    }

    // 加载内容
    if (typeof options.content === 'string') {
      $(getSelector('content')).html(options.content);
    } else if (typeof options.targetId === 'string') {
      $(getSelector('content')).html($('#' + options.targetId).html());
    } else if (typeof options.ajaxUrl === 'string') {
      $(getSelector('content')).html(
        '<div class="flex justify-center items-center h-full"><div class="animate-spin rounded-full h-8 w-8 border-b-2 border-gray-900"></div></div>',
      );
      $.ajax({
        url: options.ajaxUrl,
        success: function (data) {
          $(getSelector('content')).html(data);
          if ($.isFunction(options.ajaxCallback)) {
            options.ajaxCallback();
          }
        },
      });
    } else if (typeof options.iframeUrl === 'string') {
      $(getSelector('content')).html(
        `<iframe frameborder="0" class="w-full h-full" src="${options.iframeUrl}"></iframe>`,
      );
    }
  };

  // 私有方法：显示对话框
  const showBox = function () {
    const width = options.width;
    const height = options.height;
    let positionStyles = {};

    // 设置位置
    if (!$.isArray(options.position)) {
      positionStyles = getPositionStyles(options.position, width, height);
    } else {
      positionStyles = getPositionStyles(
        options.position[0],
        width,
        height,
        options.position[1],
      );
    }

    // 应用样式
    $(getSelector('container'))
      .css({
        width: width + 'px',
        ...positionStyles,
      })
      .removeClass('hidden');

    // 设置内容高度
    $(getSelector('content')).css({
      height: height + 'px',
    });
  };

  // 私有方法：获取位置样式
  const getPositionStyles = function (
    horizontalPos,
    width,
    height,
    verticalPos = 'center',
  ) {
    const styles = {};
    const windowWidth =
      window.innerWidth ||
      document.documentElement.clientWidth ||
      document.body.clientWidth;
    const windowHeight =
      window.innerHeight ||
      document.documentElement.clientHeight ||
      document.body.clientHeight;

    // 水平位置
    switch (horizontalPos) {
      case 'center':
        styles.left = Math.floor((windowWidth - width) / 2) + 'px';
        break;
      case 'left':
        styles.left = '0px';
        break;
      case 'right':
        styles.right = '0px';
        break;
      default:
        if (typeof horizontalPos === 'number') {
          styles.left = horizontalPos + 'px';
        } else {
          styles.left = Math.floor((windowWidth - width) / 2) + 'px';
        }
        break;
    }

    // 垂直位置
    switch (verticalPos) {
      case 'center':
        styles.top = Math.floor((windowHeight - height) / 2) + 'px';
        break;
      case 'top':
        styles.top = '0px';
        break;
      case 'bottom':
        styles.bottom = '0px';
        break;
      default:
        if (typeof verticalPos === 'number') {
          styles.top = verticalPos + 'px';
        } else {
          styles.top = Math.floor((windowHeight - height) / 2) + 'px';
        }
        break;
    }

    return styles;
  };

  // 公共方法集合
  const methods = {
    /**
     * 显示对话框
     * @param {Object} opts - 配置选项
     */
    dialog: function (opts) {
      // 重置唯一ID用于新的对话框实例
      uniqueId = generateUniqueId();

      if ($.isFunction(opts.onBeforeOpen)) {
        opts.onBeforeOpen();
      }

      initBox(opts);
      showBox();

      if ($.isFunction(options.onOpen)) {
        options.onOpen();
      }

      // 返回生成的唯一ID，以便外部引用
      return uniqueId;
    },

    /**
     * 显示确认对话框
     * @param {Object} opts - 配置选项
     */
    confirm: function (opts) {
      // 重置唯一ID用于新的对话框实例
      uniqueId = generateUniqueId();

      if ($.isFunction(opts.onBeforeOpen)) {
        opts.onBeforeOpen();
      }

      initBox(opts);

      // 添加确认和取消按钮
      if ($(getSelector('buttons')).length === 0) {
        $(getSelector('wrapper')).after(`
          <div id="${getId('buttons')}" class="flex justify-center space-x-4 p-4 border-t">
            <button id="${getId('ok')}" class="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600">确定</button>
            <button id="${getId('cancel')}" class="px-4 py-2 bg-gray-300 text-gray-700 rounded hover:bg-gray-400">取消</button>
          </div>
        `);

        $(getSelector('ok'))
          .off('click')
          .on('click', function () {
            methods.close();
            if ($.isFunction(opts.onOk)) {
              opts.onOk();
            }
          });

        $(getSelector('cancel'))
          .off('click')
          .on('click', function () {
            methods.close();
            if ($.isFunction(opts.onCancel)) {
              opts.onCancel();
            }
          });
      }

      showBox();

      if ($.isFunction(options.onOpen)) {
        options.onOpen();
      }

      // 返回生成的唯一ID，以便外部引用
      return uniqueId;
    },

    /**
     * 显示消息对话框
     * @param {Object} opts - 配置选项
     */
    message: function (opts) {
      // 重置唯一ID用于新的对话框实例
      uniqueId = generateUniqueId();

      if ($.isFunction(opts.onBeforeOpen)) {
        opts.onBeforeOpen();
      }

      initBox(opts);

      // 添加确定按钮
      if ($(getSelector('buttons')).length === 0) {
        $(getSelector('wrapper')).after(`
          <div id="${getId('buttons')}" class="flex justify-center p-4 border-t">
            <button id="${getId('ok')}" class="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600">确定</button>
          </div>
        `);

        $(getSelector('ok'))
          .off('click')
          .on('click', function () {
            methods.close();
            if ($.isFunction(opts.onOk)) {
              opts.onOk();
            }
          });
      }

      showBox();

      if ($.isFunction(options.onOpen)) {
        options.onOpen();
      }

      // 返回生成的唯一ID，以便外部引用
      return uniqueId;
    },

    /**
     * 关闭对话框
     * @param {string} id - 可选的指定实例ID，如果不提供则关闭当前实例
     */
    close: function (id) {
      // 如果提供了ID，则设置该ID为当前ID
      if (id) {
        uniqueId = id;
      }

      if ($.isFunction(options.onBeforeClose)) {
        options.onBeforeClose();
      }

      $(getSelector('container')).addClass('hidden');
      $(getSelector('overlay')).addClass('hidden');
      $(getSelector('content')).empty();
      $(getSelector('buttons')).remove();

      if ($.isFunction(options.onClose)) {
        options.onClose();
      }
    },

    /**
     * 获取当前对话框实例的唯一ID
     * @returns {string} 当前对话框的唯一ID
     */
    getId: function () {
      return uniqueId;
    },

    /**
     * 设置默认ID前缀
     * @param {string} prefix - 新的ID前缀
     */
    setIdPrefix: function (prefix) {
      if (typeof prefix === 'string' && prefix) {
        idPrefix = prefix;
      }
    },
  };

  // 导出到 jQuery 命名空间
  $.box = {
    dialog: methods.dialog,
    confirm: methods.confirm,
    message: methods.message,
    close: methods.close,
    getId: methods.getId,
    setIdPrefix: methods.setIdPrefix,
  };
})(jQuery);
