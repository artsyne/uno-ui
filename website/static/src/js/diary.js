/**
 * 日记本模块 - 处理日记本相关功能
 *
 * 基于Kiko UI的页面行为管理组件实现
 *
 * @version 2.1.0
 * @license MIT
 */

/**
 * 日记本行为配置
 * @namespace Diary
 */
const Diary = {
  /**
   * 初始化日记本模块
   *
   * @param {string} baseUrl - 基础URL
   * @returns {Array} 日记本行为配置数组
   */
  init: function (baseUrl = '') {
    log('日记本模块初始化');
    this.baseUrl = baseUrl;
    return this.getActions();
  },

  /**
   * 获取所有日记本相关的行为配置
   *
   * @returns {Array} 日记本行为配置数组
   */
  getActions: function () {
    return [
      this.showDiaryBoxAction,
      this.postInfoAction,
      this.createDiaryAction,
    ];
  },

  /**
   * 显示新增日记本对话框行为
   * @type {Object}
   */
  showDiaryBoxAction: {
    id: 'showDiaryBox',
    type: 'function',
    /**
     * 处理显示日记本对话框请求
     */
    handler: function (id, target, eventData, data) {
      log('显示新增日记本对话框', id, target, eventData);

      //todo：data 中data-form的值，后面会判断是增加还是编辑逻辑
      log(data);

      // 显示对话框
      kiko.ui.box.dialog({
        title: appI18n.getMessage('diary.create.title') || '新建日记本',
        content:
          '<div class="modal-body" id="add-diary-form-container">' +
          appI18n.getMessage('loading') +
          '</div>',
        width: 680,
        height: 400,
        cache: true,
        boxId: 'addDaybook',
      });

      // 使用Promise加载表单内容
      $.get(`${Diary.baseUrl}/form/daybook`)
        .then(function (response) {
          // 将响应插入容器
          $('#add-diary-form-container').html(response);

          // 直接使用上下文选择器，确保选择的是新添加的元素
          const $form = $('#add-diary-form-container form#daybook_editor_form');

          // 绑定主题选择事件
          $form.find('.d-theme div').on('click', function () {
            $form
              .find('.d-theme div')
              .removeClass('ring-2 ring-black')
              .addClass('hover:ring-2 hover:ring-gray-400');

            $(this)
              .removeClass('hover:ring-2 hover:ring-gray-400')
              .addClass('ring-2 ring-black');

            $form.find('#theme_id').val($(this).attr('theme'));
          });

          // 初始化验证
          $form.valid();
          log('表单验证初始化完成');
        })
        .catch(function (error) {
          log('加载表单失败:', error, 'error');
          $('#add-diary-form-container').html('加载表单失败，请重试');
        });
    },
  },

  /**
   * 创建日记本行为
   * @type {Object}
   */
  createDiaryAction: {
    id: 'createDiary',
    type: 'function',
    /**
     * 处理创建日记本请求
     *
     * @param {string} id - 行为ID
     * @param {HTMLElement} target - 目标元素
     * @param {Object} eventData - 事件数据
     */
    handler: function (id, target, eventData) {
      log('提交日记本表单');

      // 获取表单 - 确保我们获取到对话框中的表单元素
      const $form = $('#add-diary-form-container form#daybook_editor_form');

      // 验证表单 - valid('checkError')返回true表示表单有错误
      if ($form.valid('checkError')) {
        log('表单验证失败');
        return;
      }

      // 发送AJAX请求 - 全局ajaxSetup会处理表单状态
      $.ajax({
        url: `${Diary.baseUrl}/daybook`,
        type: 'POST',
        data: $form.serialize(),
        dataType: 'json',
        success: function (response) {
          if (response.Status === 0) {
            // 创建成功，添加到列表
            const html = `
              <div class="border h-64 text-center dt${response.Data.theme_id}">
                <div class="h-full content-center p-3">
                  <h1 class="text-xl mb-2 break-words">
                    <a href="/d/${response.Data.diary_id}">${response.Data.title}</a>
                    ${response.Data.view_setting === 2 ? "<i class='mdi mdi-lock-outline'></i>" : ''}
                  </h1>
                  <p>0篇</p>
                </div>
              </div>`;

            $('#daybook_list').prepend(html);

            // 关闭对话框
            kiko.ui.box.close();
            kiko.ui.notify.success(appI18n.getMessage('operation.success'));
          }
        },
      });
    },
  },

  /**
   * 提交信息行为
   * @type {Object}
   */
  postInfoAction: {
    id: 'post-info',
    type: 'function',
    selector: '#post_btn',
    /**
     * 处理信息提交请求
     *
     * @param {string} id - 行为ID
     * @param {HTMLElement} target - 目标元素
     * @param {Object} eventData - 事件数据
     */
    handler: function (id, target, eventData) {
      const form = $(eventData.element).closest('form');

      // 验证表单，如果有错误则停止提交
      if (form.valid && form.valid('checkError')) {
        return;
      }

      // 发送AJAX请求 - 全局ajaxSetup会处理表单状态
      $.ajax({
        url: `${Diary.baseUrl}/info`,
        type: 'POST',
        data: form.serialize(),
        dataType: 'json',
        success: function (response) {
          // 处理成功响应
          if (response.Status === 0) {
            // 成功后可以刷新页面或显示成功消息
            if (response.Data && response.Data.redirect) {
              window.location = response.Data.redirect;
            } else {
              kiko.ui.notify.success(appI18n.getMessage('operation.success'));
            }
          }
        },
      });
    },
  },
};
