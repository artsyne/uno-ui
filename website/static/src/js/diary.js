/**
 * 日记本模块 - 处理日记本相关功能
 * 基于Kiko UI的页面行为管理组件实现
 * @version 2.1.1
 * @license MIT
 */

/**
 * 日记本行为配置
 * @namespace Diary
 */
const Diary = {
  /**
   * 初始化日记本模块
   * @param {string} baseUrl - 基础URL
   * @returns {Array} 日记本行为配置数组
   */
  init: function (baseUrl = '') {
    log('日记本模块初始化');
    this.baseUrl = baseUrl;
    return [this.showDiaryBoxAction, this.postInfoAction];
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
      log('显示新增日记本对话框');
      log('传入Action数据:', data);
      // 显示对话框并加载表单
      const boxId = kiko.ui.box.dialog({
        title: appI18n.getMessage('diary.create.title') || '新建日记本',
        content: `<div class="modal-body" id="add-diary-form-container">${appI18n.getMessage('loading')}</div>`,
        width: 780,
        height: 460,
        showOkButton: true,
        okButtonText: appI18n.getMessage('diary.save'),
        okButtonClass:'btn-ok-box',
        buttonAlign: 'right',

        // 表单提交处理
        onOk: () => {
          const $form = $('#add-diary-form-container form#daybook_editor_form');

          // 表单验证
          if (!$form.length || $form.valid('checkError')) {
            return false; // 阻止对话框关闭
          }

          // 提交表单
          $.ajax({
            url: `${Diary.baseUrl}/daybook`,
            type: 'POST',
            data: $form.serialize(),
            dataType: 'json',
            success: (response) => {
              if (response.Status === 0) {
                // 创建成功，添加到列表
                $('#daybook_list').prepend(`
                  <div class="border h-64 text-center dt${response.Data.theme_id}">
                    <div class="h-full content-center p-3">
                      <h1 class="text-xl mb-2 break-words">
                        <a href="/d/${response.Data.diary_id}">${response.Data.title}</a>
                        ${response.Data.view_setting === 2 ? "<i class='mdi mdi-lock-outline'></i>" : ''}
                      </h1>
                      <p>0篇</p>
                    </div>
                  </div>`);

                kiko.ui.box.close(boxId);
                kiko.ui.notify.success(appI18n.getMessage('operation.success'));
              }
            },
          });

          return true;
        },
      });

      // 加载表单内容
      $.get(`${Diary.baseUrl}/form/diary`)
        .then((response) => {
          // 初始化表单
          $('#add-diary-form-container').html(response);
          const $form = $('#add-diary-form-container form#daybook_editor_form');

          // 绑定主题选择事件
          $form.find('.d-theme div').on('click', function () {
            const $this = $(this);
            $form
              .find('.d-theme div')
              .removeClass('ring-2 ring-black')
              .addClass('hover:ring-2 hover:ring-gray-400');

            $this
              .removeClass('hover:ring-2 hover:ring-gray-400')
              .addClass('ring-2 ring-black');

            $form.find('#theme_id').val($this.attr('theme'));
          });

          //根据初始化data值formAction，设置表单行为
          $form.find('input[name="action"]').val(data.formAction);
          // 初始化验证
          $form.valid();
        })
        .catch((error) => {
          log('加载表单失败:', error, 'error');
          $('#add-diary-form-container').html('加载表单失败，请重试');
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
     */
    handler: function (id, target, eventData) {
      const form = $(eventData.element).closest('form');

      // 验证表单，如果有错误则停止提交
      if (form.valid && form.valid('checkError')) {
        return;
      }

      // 发送AJAX请求
      $.ajax({
        url: `${Diary.baseUrl}/info`,
        type: 'POST',
        data: form.serialize(),
        dataType: 'json',
        success: (response) => {
          if (response.Status === 0) {
            // 处理成功响应
            if (response.Data?.redirect) {
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
