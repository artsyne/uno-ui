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
    this.baseUrl = baseUrl;
    return [
      this.showDiaryBoxAction,
      this.showPostBoxAction,
      this.deletePostAction,
      this.postInfoAction,
      this.gotoUrlAction,
      this.getAjaxContentAction,
    ];
  },
  /**
   * 获取ajax内容行为
   */
  getAjaxContentAction: {
    id: 'getAjaxContent',
    type: 'function',
    handler: function (id, target, eventData, data) {
      history.pushState(null, '', data.url);
      return $.ajax({
        url: data.url+ '?isajax=true',
        type: 'GET',
        dataType: 'html',
        success: function (data) {
          $('.post-container').first().html(data);
        },
      });
    },
  },

  /**
   * 删除日记行为
   */
  deletePostAction: {
    id: 'deletePost',
    type: 'function',
    handler: function (id, target, eventData, data) {
      // 确认删除
      const confirmMessage = kiko.i18n.getMessage('post.delete.confirm') || '确定要删除这篇日记吗？';
      
      kiko.ui.box.confirm({
        content: confirmMessage,
        width: 500,
        height: 80,
        okButtonClass:'btn',
        onOk: function() {
          // 发送删除请求
          $.ajax({
            url: `${Diary.baseUrl}/post/delete`,
            type: 'POST',
            data: {
              post_id: data.postId,
              diary_id: data.diaryId
            },
            success: (response) => {
              if (response.code === 0) {
                // 删除成功，从页面移除对应的日记卡片
                const $postCard = $(`.card[data-post-id="${data.postId}"]`);
                $postCard.fadeOut(300, function() {
                  $(this).remove();
                });
                kiko.ui.notify.success(kiko.i18n.getMessage('post.delete.success') || '删除成功');
              } else {
                kiko.ui.notify.error(response.message || '删除失败');
              }
            },
            error: (xhr, status, error) => {
              kiko.ui.notify.error('网络错误，删除失败');
            }
          });
        }
      });
    },
  },
  
  /**
   * 跳转url行为
   */
  gotoUrlAction: {
    id: 'gotoUrl',
    type: 'function',
    handler: function (id, target, eventData, data) {
      const originalUrl = window.location.href;
      kiko.ui.box.dialog({
        ajaxUrl: data.url+ '?isajax=true',
        containerClass:'!border',
        widthClass:'lg:w-[90%]',
        heightClass:'h-full',
        onClose: () => {
          history.pushState(null, '', originalUrl);
        },
        onOpen: () => {
          history.pushState(null, '', data.url);
        },
      });
    },
  },

  /*
   * 显示写日记对话框行为
   */
  showPostBoxAction: {
    id: 'showPostBox',
    type: 'function',

    /**
     * 处理显示写日记对话框请求
     * @param {string} id - 行为ID
     * @param {HTMLElement} target - 触发元素
     * @param {Object} eventData - 事件数据
     * @param {Object} data - 附加数据
     */
    handler: function (id, target, eventData, data) {
      const isEdit = data.formAction === 'post.edit';
      const dialogTitle = isEdit
        ? kiko.i18n.getMessage('post.edit.title')
        : kiko.i18n.getMessage('post.create.title');

      // 显示对话框并加载表单
      const boxId = kiko.ui.box.dialog({
        title: dialogTitle,
        content: `<div class="modal-body" id="post-form-container">${kiko.i18n.getMessage('loading')}</div>`,
        width: 880,
        height: 570,
        showOkButton: true,
        okButtonText: kiko.i18n.getMessage('post.save') || '保存',
        containerClass:'!border',
        okButtonClass: 'btn',
        buttonAlign: 'right',
        onOk: () => this._handlePostSubmit(data, boxId),
      });

      // 加载表单内容
      this._loadPostForm(data, boxId);
    },

    /**
     * 加载写日记表单内容
     * @private
     */
    _loadPostForm: function (data, boxId) {
      const params = {
        action: data.formAction,
        ...(data.postId && { post_id: data.postId }),
        ...(data.diaryId && { diary_id: data.diaryId }),
      };

      $.get(`${Diary.baseUrl}/form/post`, params)
        .then((response) => {
          this._initializePostForm(
            response,
            data.formAction,
            data.diaryId,
            data.postId,
          );
        })
        .catch((error) => {
          $('#post-form-container').html(`
            <div class="text-center text-red-500 py-8">
              <p>加载表单失败，请重试</p>
              <button class="btn btn-sm mt-2" onclick="location.reload()">刷新页面</button>
            </div>
          `);
        });
    },

    /**
     * 初始化写日记表单
     * @private
     */
    _initializePostForm: function (formHtml, formAction, diaryId, postId) {
      $('#post-form-container').html(formHtml);
      const $form = $('#post-form-container form');

      // 设置表单行为
      $form.find('input[name="action"]').val(formAction);
      $form.find('input[name="diary_id"]').val(diaryId || '');
      $form.find('input[name="post_id"]').val(postId || '');

      // 初始化验证
      if ($form.valid) {
        $form.valid();
      }

      // 光标focus到标题输入框
      const $titleInput = $form.find('input[name="info_title"]');
      if ($titleInput.length) {
        $titleInput.focus();
      }
    },

    /**
     * 处理写日记表单提交
     * @private
     */
    _handlePostSubmit: function (data, boxId) {
      const $form = $('#post-form-container form');

      // 表单验证
      if (!$form.length || ($form.valid && $form.valid('checkError'))) {
        return false; // 阻止对话框关闭
      }

      // 显示提交状态
      const $okButton = $('.kiko-dialog-footer .btn');
      const originalText = $okButton.text();
      $okButton.prop('disabled', true).text('保存中...');

      // 提交表单
      $.ajax({
        url: `${Diary.baseUrl}/post`,
        type: 'POST',
        data: $form.serialize(),
        dataType: 'json',
        success: (response) => {
          // 恢复按钮状态
          $okButton.prop('disabled', false).text(originalText);

          if (response.code === 0) {
            this._handlePostSubmitSuccess(data.formAction, response.data);
            kiko.ui.box.close(boxId);
            kiko.ui.notify.success(kiko.i18n.getMessage('post.save.success'));
          } else {
            this._handleSubmitError(response.message);
          }
        },
        error: (xhr, status, error) => {
          // 恢复按钮状态
          $okButton.prop('disabled', false).text(originalText);

          log('提交失败:', error, 'error');
          this._handleSubmitError('网络错误，请重试');
        },
      });

      return false; // 阻止对话框自动关闭
    },

    /**
     * 处理写日记提交成功
     * @private
     */
    _handlePostSubmitSuccess: function (formAction, responseData) {
      if (formAction === 'post.create') {
        // 新建日记成功，增加新的日记到列表
        this._addPostToList(responseData);
      } else if (formAction === 'post.edit') {
        // 编辑日记成功，更新列表中的日记信息
        this._updatePostInfo(responseData);
      }
    },
    /**
     * 更新日记信息
     * @private
     */
    _updatePostInfo: function (postData) {
      // 确保必要的数据存在
      if (!postData || !postData.post_id) {
        log('更新日记信息失败：缺少必要数据', 'error');
        return;
      }

      // 查找对应的日记卡片
      const $postCard = $(
        `#diary_post_container .card[data-post-id="${postData.post_id}"]`,
      );
      if (!$postCard.length) {
        log('未找到对应的日记卡片，无法更新', 'error');
        return;
      }

      // 更新标题和摘要
      $postCard.find('h1.post-title').text(escapeHtml(postData.title || ''));
      $postCard.find('p.text-relaxed a').text(escapeHtml(postData.summary || ''));
    },
    /**
     * 添加新日记到列表
     * @private
     */
    _addPostToList: function (postData) {
      // 确保必要的数据存在
      if (!postData || !postData.post_id) {
        log('添加日记到列表失败：缺少必要数据', 'error');
        return;
      }

      // 构建日记卡片HTML，基于最新的diary/index.html结构
      const postHtml = `
        <div class="card mb-6 break-inside-avoid" data-post-id="${postData.post_id}">
          <div class="post-edit">
            <a href="#" data-tooltip-target="post_edit${postData.post_id}" data-tooltip-position="bottom"
               data-tooltip-event="click,mouseover">
              <i class="mdi mdi-dots-horizontal"></i>
            </a>
          </div>
          <div id="post_edit${postData.post_id}" class="tooltip hidden">
            <ul>
              <li>
                <a href="" class="tooltip-link" data-action="showPostBox" data-form-action="post.edit"
                   data-diary-id="${postData.diary_id || ''}" data-post-id="${postData.post_id}">
                  <i class="mdi mdi-pencil-outline mr-2"></i> 编辑
                </a>
              </li>
              <li>
                <a href="" class="tooltip-link" data-action="deletePost"
                   data-diary-id="${postData.diary_id || ''}" data-post-id="${postData.post_id}">
                  <i class="mdi mdi-delete-outline mr-2"></i> 删除
                </a>
              </li>
            </ul>
          </div>
          <div class="mb-4">
            <h1 class="post-title">
              ${escapeHtml(postData.title)}
            </h1>
            <p class="text-relaxed cursor-pointer">
              <a href="/p/${postData.post_id}" data-action="gotoUrl" data-url="/p/${postData.post_id}" data-ajax="true"
                 style="all:unset !important">
                ${escapeHtml(postData.summary || '')}
              </a>
            </p>
          </div>
          <div class="flex justify-between">
            <div>
              <span class="post-time">
                ${formatTime(postData.add_utc_time || '')}
              </span>
            </div>
            <div>
              <a href="#" class="mr-2">
                <i class="mdi mdi-comment-text-outline mr-1"></i> <span>${postData.comment_num || 0}</span>
              </a>
              <a href="#">
                <i class="mdi mdi-thumb-up-outline mr-1"></i><span>${postData.like_num || 0}</span>
              </a>
            </div>
          </div>
        </div>
      `;

      // 将新日记添加到列表的第二个位置（第一个是写日记按钮）
      const $container = $('#diary_post_container');
      const $firstCard = $container.find('.card').first();

      if ($firstCard.length) {
        $firstCard.after(postHtml);
      } else {
        $container.prepend(postHtml);
      }

      // 添加动画效果
      const $newPost = $container.find(`[data-post-id="${postData.post_id}"]`);
      $newPost.hide().fadeIn(600);

      // 注意：使用事件委托后，tooltip会自动工作，无需手动初始化

      log('成功添加新日记到列表', postData);
    },

    /**
     * 处理提交错误
     * @private
     */
    _handleSubmitError: function (errorMessage) {
      kiko.ui.notify.error(errorMessage);
    },
  },

  /**
   * 显示日记对话框行为
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
        ? kiko.i18n.getMessage('diary.edit.title')
        : kiko.i18n.getMessage('diary.create.title');

      // 显示对话框并加载表单
      const boxId = kiko.ui.box.dialog({
        title: dialogTitle,
        content: `<div class="modal-body" id="add-diary-form-container">${kiko.i18n.getMessage('loading')}</div>`,
        width: 800,
        height: 500,
        showOkButton: true,
        containerClass:'!border',
        okButtonText: kiko.i18n.getMessage('diary.save') || '保存',
        okButtonClass: 'btn',
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
      const $okButton = $('.kiko-dialog-footer .btn');
      const originalText = $okButton.text();
      $okButton.prop('disabled', true).text('保存中...');

      // 提交表单
      $.ajax({
        url: `${Diary.baseUrl}/diary`,
        type: 'POST',
        data: $form.serialize(),
        dataType: 'json',
        success: (response) => {
          // 恢复按钮状态
          $okButton.prop('disabled', false).text(originalText);

          if (response.code === 0) {
            this._handleSubmitSuccess(data.formAction, response.data);
            kiko.ui.box.close(boxId);
            kiko.ui.notify.success(
              kiko.i18n.getMessage('operation.success') || '操作成功',
            );
          } else {
            this._handleSubmitError(response.message || '操作失败');
          }
        },
        error: (xhr, status, error) => {
          // 恢复按钮状态
          $okButton.prop('disabled', false).text(originalText);

          log('提交失败:', error, 'error');
          this._handleSubmitError('网络错误，请重试');
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
        <div class="border h-72 text-center dt${diaryData.theme_id}">
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
          .removeClass('ring-2 ring-black dark:ring-white ring-offset-2 dark:ring-offset-gray-800')
          .addClass('hover:ring-2 hover:ring-gray-400 dark:hover:ring-gray-300');

        // 设置选中样式 - 支持dark模式
        $this
          .removeClass('hover:ring-2 hover:ring-gray-400 dark:hover:ring-gray-300')
          .addClass('ring-2 ring-black dark:ring-white ring-offset-2 dark:ring-offset-gray-800');

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
    id: 'postInfo',
    type: 'function',

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
          if (response.code === 0) {
            // 处理成功响应
            if (response.data?.redirect) {
              window.location = response.data.redirect;
            } else {
              kiko.ui.notify.success(appI18n.getMessage('operation.success'));
            }
          }
        },
      });
    },
  },
};
