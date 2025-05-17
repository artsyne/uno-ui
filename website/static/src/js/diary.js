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
    console.log('[Diary] 初始化日记本模块');
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
    id: 'show-diary-box',
    type: 'function',
    selector: '[jq-bind-click="addDaybook"]',
    /**
     * 处理显示日记本对话框请求
     */
    handler: function () {
      console.log('[Diary] 显示新增日记本对话框');

      // 显示对话框
      kiko.ui.box.dialog({
        title: appI18n.getMessage('operation.success'), // 可根据需要自定义key
        content:
          '<div class="modal-body" id="add-diary-form-container">' +
          appI18n.getMessage('loading') +
          '</div>',
        width: 700,
        height: 400,
        cache: true,
        boxId: 'addDaybook',
      });

      // 加载表单内容
      $.get(`${Diary.baseUrl}/form/daybook`, function (response) {
        $('#add-diary-form-container').html(response);

        // 初始化表单验证
        if ($.fn.valid) {
          $('#daybook_editor_form').valid();
        }

        // 绑定主题选择事件
        $('.d-theme>div')
          .mouseover(function () {
            $(this).addClass('dtfhover');
          })
          .mouseout(function () {
            $(this).removeClass('dtfhover');
          })
          .click(function () {
            $('.d-theme>div').removeClass('dtfselected');
            $(this).addClass('dtfselected');
            $('#theme_id').val($(this).attr('theme'));
          });
      });
    },
  },

  /**
   * 创建日记本行为
   * @type {Object}
   */
  createDiaryAction: {
    id: 'create-diary',
    type: 'function',
    selector: '#diary_submit_btn',
    /**
     * 处理创建日记本请求
     *
     * @param {string} id - 行为ID
     * @param {HTMLElement} target - 目标元素
     * @param {Object} eventData - 事件数据
     */
    handler: function (id, target, eventData) {
      const form = $('#daybook_editor_form');

      // 验证表单，如果有错误则停止提交
      if (form.valid && form.valid('checkError')) {
        return;
      }

      // 发送AJAX请求 - 全局ajaxSetup会处理表单状态
      $.ajax({
        url: `${Diary.baseUrl}/daybook`,
        type: 'POST',
        data: form.serialize(),
        dataType: 'json',
        success: function (response) {
          if (response.Status === 0) {
            // 创建成功，添加到列表
            const html = `
              <div class="col-sm-6 col-md-4 col-lg-3">
                <div class="card card-diary-theme dt${response.Data.theme_id}">
                  <div class="card-body">
                    <h5>
                      ${response.Data.view_setting === 2 ? "<i class='mdi mdi-lock-outline'></i>" : ''}
                      <a href='/d/${response.Data.diary_id}'>${response.Data.title}</a>
                    </h5>
                    <p>0篇</p>
                  </div>
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
