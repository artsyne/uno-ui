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
   * 显示日记本对话框行为
   * @type {Object}
   */
  showDiaryBoxAction: {
    id: 'showDiaryBox',
    type: 'function',

    /**
     * 处理显示日记本对话框请求
     */
    handler: function (id, target, eventData, data) {
      log('显示日记本对话框');
      log('传入Action数据:', data);

      const isEdit = data.formAction === 'diary.edit';
      const dialogTitle = isEdit
        ? appI18n.getMessage('diary.edit.title') || '编辑日记本'
        : appI18n.getMessage('diary.create.title') || '新建日记本';

      // 显示对话框并加载表单
      const boxId = kiko.ui.box.dialog({
        title: dialogTitle,
        content: `<div class="modal-body" id="add-diary-form-container">${appI18n.getMessage('loading') || '加载中...'}</div>`,
        width: 780,
        height: 460,
        showOkButton: true,
        okButtonText: appI18n.getMessage('diary.save') || '保存',
        okButtonClass: 'btn-red',
        buttonAlign: 'right',
        onOk: () => this._handleFormSubmit(data, boxId),
      });

      // 加载表单内容
      this._loadForm(data, boxId);
    },

    /**
     * 处理表单提交
     * @private
     */
    _handleFormSubmit: function (data, boxId) {
      const $form = $('#add-diary-form-container form#daybook_editor_form');

      // 表单验证
      if (!$form.length || $form.valid('checkError')) {
        return false; // 阻止对话框关闭
      }

      // 显示提交状态
      const $okButton = $('.kiko-dialog-footer .btn-red');
      const originalText = $okButton.text();
      $okButton.prop('disabled', true).text('保存中...');

      // 提交表单
      $.ajax({
        url: `${Diary.baseUrl}/diary`,
        type: 'POST',
        data: $form.serialize(),
        dataType: 'json',
        success: (response) => {
          if (response.Status === 0) {
            this._handleSubmitSuccess(data.formAction, response.Data);
            kiko.ui.box.close(boxId);
            kiko.ui.notify.success(
              appI18n.getMessage('operation.success') || '操作成功',
            );
          } else {
            this._handleSubmitError(response.Message || '操作失败');
          }
        },
        error: (xhr, status, error) => {
          log('提交失败:', error, 'error');
          this._handleSubmitError('网络错误，请重试');
        },
        complete: () => {
          // 恢复按钮状态
          $okButton.prop('disabled', false).text(originalText);
        },
      });

      return false; // 阻止对话框自动关闭
    },

    /**
     * 处理提交成功
     * @private
     */
    _handleSubmitSuccess: function (formAction, responseData) {
      if (formAction === 'diary.add') {
        this._addDiaryToList(responseData);
      } else if (formAction === 'diary.edit') {
        this._updateDiaryInfo(responseData);
      }
    },

    /**
     * 添加日记本到列表
     * @private
     */
    _addDiaryToList: function (diaryData) {
      const lockIcon =
        diaryData.view_setting === 2
          ? "<i class='mdi mdi-lock-outline'></i>"
          : '';
      const diaryHtml = `
        <div class="border h-64 text-center dt${diaryData.theme_id}">
          <div class="h-full content-center p-3">
            <h1 class="text-xl mb-2 break-words">
              <a href="/d/${diaryData.diary_id}">${diaryData.title}</a>
              ${lockIcon}
            </h1>
            <p>0篇</p>
          </div>
        </div>`;
      $('#daybook_list').prepend(diaryHtml);
    },

    /**
     * 更新日记本信息
     * @private
     */
    _updateDiaryInfo: function (diaryData) {
      // 更新标题
      const $title = $('#diary_title');
      $title.text(diaryData.title);

      // 处理锁图标
      $title.find('i.mdi-lock-outline').remove();
      if (diaryData.view_setting === 2) {
        $title.append("<i class='mdi mdi-lock-outline'></i>");
      }

      // 更新描述
      $('#diary_remark').text(diaryData.remark || '');

      // 更新主题样式
      const $main = $('.main');
      $main.removeClass((index, className) => {
        return (className.match(/(^|\s)dt\d+/g) || []).join(' ');
      });
      $main.addClass(`dt${diaryData.theme_id}`);
    },

    /**
     * 处理提交错误
     * @private
     */
    _handleSubmitError: function (errorMessage) {
      kiko.ui.notify.error(errorMessage);
    },

    /**
     * 加载表单内容
     * @private
     */
    _loadForm: function (data, boxId) {
      const params = {
        action: data.formAction,
        ...(data.diaryId && { diary_id: data.diaryId }),
      };

      $.get(`${Diary.baseUrl}/form/diary`, params)
        .then((response) => {
          this._initializeForm(response, data.formAction);
        })
        .catch((error) => {
          log('加载表单失败:', error, 'error');
          $('#add-diary-form-container').html(`
            <div class="text-center text-red-500 py-8">
              <p>加载表单失败，请重试</p>
              <button class="btn btn-sm mt-2" onclick="location.reload()">刷新页面</button>
            </div>
          `);
        });
    },

    /**
     * 初始化表单
     * @private
     */
    _initializeForm: function (formHtml, formAction) {
      $('#add-diary-form-container').html(formHtml);
      const $form = $('#add-diary-form-container form#daybook_editor_form');

      // 绑定主题选择事件
      this._bindThemeSelection($form);

      // 设置表单行为
      $form.find('input[name="action"]').val(formAction);

      // 初始化验证
      if ($form.valid) {
        $form.valid();
      }
    },

    /**
     * 绑定主题选择事件
     * @private
     */
    _bindThemeSelection: function ($form) {
      $form.find('.d-theme div').on('click', function () {
        const $this = $(this);
        const $themes = $form.find('.d-theme div');

        // 重置所有主题样式
        $themes
          .removeClass('ring-2 ring-black')
          .addClass('hover:ring-2 hover:ring-gray-400');

        // 设置选中样式
        $this
          .removeClass('hover:ring-2 hover:ring-gray-400')
          .addClass('ring-2 ring-black');

        // 更新隐藏字段值
        $form.find('#theme_id').val($this.attr('theme'));
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
