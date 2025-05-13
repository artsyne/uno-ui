/**
 * jQuery选项卡插件
 *
 * 这个插件提供了一种简单而强大的方式创建基于标签的内容切换组件。
 * 支持通过URL哈希导航，实现页面刷新后保持当前标签状态。
 * 使用Tailwind CSS 4进行样式设计，无需额外的CSS文件。
 *
 * @author nobody
 * @version 1.0.0
 * @license MIT
 *
 * 使用说明：
 *
 * 1. 基础HTML结构：
 *    <div class="mb-4">
 *      <!-- 标签头部 -->
 *      <div class="border-b">
 *        <a href="#tab1" tabgroup="demo" tabtarget="content1">标签1</a>
 *        <a href="#tab2" tabgroup="demo" tabtarget="content2">标签2</a>
 *      </div>
 *      <!-- 标签内容 -->
 *      <div id="content1">内容1</div>
 *      <div id="content2">内容2</div>
 *    </div>
 *
 * 2. 初始化选项卡：
 *    $.tab.init();
 *
 * 3. 手动触发标签切换：
 *    $.tab.switchTo('demo', 'content2');
 *
 * 4. 支持的HTML属性：
 *    - tabgroup: 标签组名称，同组标签互相切换
 *    - tabtarget: 关联的内容区域ID
 *    - rel: 用于URL哈希导航
 *
 * 5. Tailwind CSS 4 样式：
 *    本插件默认使用Tailwind CSS 4样式，主要包括：
 *    - 标签头部: flex, border-b, space-x-2, mb-4
 *    - 非激活标签: px-4, py-2, text-gray-500, hover:text-gray-700
 *    - 激活标签: px-4, py-2, text-blue-600, border-b-2, border-blue-500
 *    - 标签内容: p-4
 *
 *    如需自定义样式，可以：
 *    1. 在初始化时通过选项设置自定义类名
 *    2. 使用CSS选择器覆盖默认样式
 *    3. 通过修改Tailwind配置扩展主题颜色和其他属性
 */

(function ($) {
  // 私有变量
  let idPrefix = 'jqtab-';
  let settings = {
    activeClass: 'text-blue-600 border-b-2 border-blue-500',
    inactiveClass: 'text-gray-500 hover:text-gray-700',
    baseTabClass: 'px-4 py-2 transition-colors duration-150',
  };

  // 获取URL锚点值
  const getAnchor = function () {
    return window.location.hash ? window.location.hash.substring(1) : '';
  };

  // 生成一个随机ID后缀，确保多个实例不会冲突
  const generateUniqueId = function () {
    return Math.random().toString(36).substring(2, 10);
  };

  // 私有方法集合
  const internal = {
    /**
     * 处理标签点击事件
     * @param {Event} e - 事件对象
     */
    handleTabClick: function (e) {
      e.preventDefault();

      const $target = $(this);
      const tabGroup = $target.attr('tabgroup');
      const tabTarget = $target.attr('tabtarget');

      if (!tabGroup || !tabTarget) return;

      // 切换到目标标签
      methods.switchTo(tabGroup, tabTarget);

      // 更新URL哈希以支持导航
      if ($target.attr('rel')) {
        window.location.hash = $target.attr('rel');
      }
    },

    /**
     * 初始化标签样式
     * @param {jQuery} $tab - 标签元素
     */
    setupTabStyles: function ($tab) {
      // 如果标签还没有基础样式类，添加它
      if (!$tab.hasClass(settings.baseTabClass)) {
        $tab.addClass(settings.baseTabClass);
      }

      // 默认添加非激活样式
      if (!$tab.hasClass('current')) {
        $tab.addClass(settings.inactiveClass);
      } else {
        $tab.addClass(settings.activeClass);
      }
    },
  };

  // 公共方法集合
  const methods = {
    /**
     * 初始化选项卡系统
     * @param {Object} options - 配置选项
     */
    init: function (options) {
      // 合并设置
      if (options) {
        settings = $.extend(settings, options);
      }

      // 获取URL锚点
      const anchor = getAnchor();

      // 查找并初始化所有未加载的选项卡
      $('a[tabgroup][loadedtab!=true]').each(function (i, el) {
        const $tab = $(el);

        // 绑定点击事件
        $tab.on('click', internal.handleTabClick);

        // 设置基础样式
        internal.setupTabStyles($tab);

        // 如果有锚点匹配，激活该标签
        if ($tab.attr('rel') === anchor && anchor.length > 0) {
          $(`a[tabgroup="${$tab.attr('tabgroup')}"]`)
            .removeClass('current ' + settings.activeClass)
            .addClass(settings.inactiveClass);
          $tab
            .addClass('current ' + settings.activeClass)
            .removeClass(settings.inactiveClass);

          // 也显示对应的内容
          if ($tab.attr('tabtarget')) {
            methods.switchTo(
              $tab.attr('tabgroup'),
              $tab.attr('tabtarget'),
              false,
            );
          }
        }

        // 标记为已加载
        $tab.attr('loadedtab', 'true');
      });

      // 美化标签容器
      $('div:has(a[tabgroup])')
        .not('.tab-styled')
        .each(function () {
          $(this).addClass('flex border-b space-x-2 mb-4 tab-styled');
        });

      return this;
    },

    /**
     * 切换到指定标签
     * @param {string} tabGroup - 标签组名称
     * @param {string} targetId - 目标内容区域ID
     * @param {boolean} updateTabs - 是否更新标签状态，默认为true
     */
    switchTo: function (tabGroup, targetId, updateTabs = true) {
      // 如果需要更新标签状态
      if (updateTabs) {
        // 移除所有同组标签的激活状态
        $(`a[tabgroup="${tabGroup}"]`)
          .removeClass('current ' + settings.activeClass)
          .addClass(settings.inactiveClass);

        // 激活目标标签
        $(`a[tabgroup="${tabGroup}"][tabtarget="${targetId}"]`)
          .addClass('current ' + settings.activeClass)
          .removeClass(settings.inactiveClass);
      }

      // 处理内容区域的显示/隐藏
      $(`a[tabgroup="${tabGroup}"][tabtarget]`).each(function () {
        const $tab = $(this);
        const $content = $(`#${$tab.attr('tabtarget')}`);

        if ($tab.attr('tabtarget') === targetId) {
          $content.removeClass('hidden').show();
        } else {
          $content.addClass('hidden').hide();
        }
      });

      return this;
    },

    /**
     * 设置选项卡系统选项
     * @param {Object} options - 配置选项
     */
    setOptions: function (options) {
      settings = $.extend(settings, options);
      return this;
    },

    /**
     * 刷新所有标签样式
     */
    refreshStyles: function () {
      $('a[tabgroup][loadedtab=true]').each(function () {
        const $tab = $(this);
        internal.setupTabStyles($tab);
      });
      return this;
    },
  };

  // 导出到 jQuery 命名空间
  $.tab = {
    init: methods.init,
    switchTo: methods.switchTo,
    setOptions: methods.setOptions,
    refreshStyles: methods.refreshStyles,
  };
})(jQuery);
