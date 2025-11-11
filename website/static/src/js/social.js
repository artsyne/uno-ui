/**
 * 社交功能模块 - 处理统计信息和点赞功能
 * 基于Kiko UI的页面行为管理组件实现
 * @version 1.0.0
 * @license MIT
 */

/**
 * 社交功能行为配置
 * @namespace Social
 */
const Social = {
  /**
   * 初始化社交功能模块
   * @param {string} baseUrl - 基础URL
   * @returns {Array} 社交功能行为配置数组
   */
  init: function (baseUrl = '') {
    this.baseUrl = baseUrl;
    // 自动加载统计信息
    this.loadPostsStats();
    // 自动加载日记本点赞状态
    this.loadDiaryLikeStatus();
    // 初始化评论模块
    this.initCommentsModule();

    return [
      this.likePostAction,
      this.likeDiaryAction,
      this.submitCommentAction,
      this.deleteCommentAction,
      this.loadMoreCommentsAction,
    ];
  },

  /**
   * 获取社交设置值
   * @returns {number} 社交设置值，1表示开启，2表示关闭
   */
  getSocialSetting: function () {
    // 首先尝试从页面的data属性获取
    const pageData = $('main').data('social-setting');
    if (pageData !== undefined) {
      return parseInt(pageData, 10);
    }
    return 2;
  },

  /**
   * 加载日记统计信息
   * 自动检测页面中的stats容器并加载对应的统计数据
   * 只有当SocialSetting为1时才会加载统计信息
   */
  loadPostsStats: function () {
    // 检查SocialSetting，只有为1时才加载统计信息
    const socialSetting = this.getSocialSetting();
    if (socialSetting !== 1) {
      return;
    }

    // 查找所有stats容器（包括日记列表页和详情页）
    const $statsContainers = $('[id^="stats_"], [id^="detail_stats_"]');

    if ($statsContainers.length === 0) {
      return;
    }

    // 提取所有PostID
    const postIds = [];
    $statsContainers.each(function () {
      const id = $(this).attr('id');
      let postId = '';
      if (id.startsWith('stats_')) {
        postId = id.replace('stats_', '');
      } else if (id.startsWith('detail_stats_')) {
        postId = id.replace('detail_stats_', '');
      }
      if (postId) {
        postIds.push(postId);
      }
    });

    if (postIds.length === 0) {
      return;
    }

    // 发送ajax请求获取统计信息
    $.ajax({
      url: `${this.baseUrl}/api/post/stats`,
      type: 'POST',
      contentType: 'application/json',
      data: JSON.stringify({
        postids: postIds.join(','),
      }),
      global: false,
      dataType: 'json',
    })
    .done((response) => {
      if (response.code === 0 && response.data && response.data.stats) {
        this._updatePostsStats(response.data.stats);
        // 加载点赞状态（仅在用户登录时）
        this._loadPostsLikeStatus(postIds);
      }
    })
    .fail((xhr, status, error) => {
      // 静默处理错误，不显示通知
      console.log('加载统计信息失败:', error);
    });
  },

  /**
   * 更新日记统计信息到页面
   * @private
   * @param {Array} statsData - 统计数据数组
   */
  _updatePostsStats: function (statsData) {
    // 获取当前日记本ID（从URL路径中提取）
    const pathParts = window.location.pathname.split('/');
    const diaryId = pathParts[2]; // /d/{diaryId} 格式

    // 遍历统计数据数组，更新对应的DOM元素
    statsData.forEach((stats) => {
      // 查找日记列表页的容器
      const $listContainer = $(`#stats_${stats.post_id}`);
      // 查找日记详情页的容器
      const $detailContainer = $(`#detail_stats_${stats.post_id}`);

      // 构建统计信息HTML
      const statsHtml = `
        <button class="btn-social mr-2">
          <i class="mdi mdi-eye mr-1"></i> <span>${stats.view_num || 0}</span>
        </button>
        <button class="btn-social mr-2">
          <i class="mdi mdi-comment-text-outline mr-1"></i> <span>${stats.comment_num || 0}</span>
        </button>
        <button class="btn-social" data-action="likePost" data-post-id="${stats.post_id}" data-diary-id="${diaryId}">
          <i class="mdi mdi-thumb-up-outline mr-1"></i><span>${stats.like_num || 0}</span>
        </button>`;

      // 更新日记列表页容器内容
      if ($listContainer.length && stats) {
        $listContainer.html(statsHtml);
      }

      // 更新日记详情页容器内容
      if ($detailContainer.length && stats) {
        $detailContainer.html(statsHtml);
      }
    });
  },

  /**
   * 加载日记点赞状态
   * 仅在用户登录时调用，更新点赞按钮的状态
   */
  _loadPostsLikeStatus: function (postIds) {
    // 检查用户是否登录
    if (!UserState.isLoggedIn()) {
      return; // 用户未登录，不加载点赞状态
    }

    // 发送ajax请求获取点赞状态
    $.ajax({
      url: `${this.baseUrl}/api/post/like/status`,
      type: 'POST',
      contentType: 'application/json',
      data: JSON.stringify({
        postids: postIds.join(','),
      }),
      global: false,
      dataType: 'json',
    })
    .done((response) => {
      if (response.code === 0 && response.data && response.data.status) {
        this._updatePostsLikeStatus(response.data.status);
      }
    })
    .fail((xhr, status, error) => {
      // 静默处理错误，不显示通知
      console.log('加载点赞状态失败:', error);
    });
  },

  /**
   * 更新日记点赞状态
   * 根据返回的点赞状态更新按钮样式
   */
  _updatePostsLikeStatus: function (likeStatusData) {
    likeStatusData.forEach((status) => {
      const $likeBtn = $(
        `button[data-action="likePost"][data-post-id="${status.postid}"]`,
      );

      if ($likeBtn.length) {
        const $icon = $likeBtn.find('i');

        if (status.liked) {
          // 已点赞状态
          $icon
            .removeClass('mdi-thumb-up-outline')
            .addClass('mdi-thumb-up like-icon-liked');
        } else {
          // 未点赞状态
          $icon.removeClass('mdi-thumb-up').addClass('mdi-thumb-up-outline');
        }
      }
    });
  },

  /**
   * 加载日记本点赞状态
   * 根据data-social-setting决定是否显示点赞按钮
   */
  loadDiaryLikeStatus: function () {
    // 检查社交设置
    const socialSetting = this.getSocialSetting();
    if (socialSetting !== 1) {
      // 如果社交设置不为1，清空容器
      $('#diary-like-container').empty();
      return;
    }

    // 获取日记ID
    const $container = $('#diary-like-container');
    const diaryId = $container.data('diary-id');
    if (!diaryId) {
      return;
    }

    // 发送AJAX请求获取日记点赞状态
    $.ajax({
      url: `${this.baseUrl}/api/diary/like/status?diaryid=${diaryId}`,
      method: 'GET',
      global: false,
    })
    .done((response) => {
      if (response.success && response.data) {
        this._renderDiaryLikeButton(response.data);
      }
    })
    .fail((xhr, status, error) => {
      // 静默处理错误，不显示通知
      console.log('加载日记点赞状态失败:', error);
    });
  },

  /**
   * 渲染日记本点赞按钮
   * @private
   * @param {Object} likeData - 点赞数据
   */
  _renderDiaryLikeButton: function (likeData) {
    const $container = $('#diary-like-container');
    const diaryId = $container.data('diary-id');

    // 根据点赞状态选择图标和样式类
    const iconClass = likeData.is_liked
      ? 'mdi mdi-thumb-up like-icon-liked'
      : 'mdi mdi-thumb-up-outline';

    const buttonHtml = `
      <button class="btn" data-action="likeDiary" data-diary-id="${diaryId}">
        <i class="${iconClass} mr-1"></i>
        <span>${likeData.like_count || 0}</span>
      </button>
    `;

    $container.html(buttonHtml);
  },

  /**
   * 日记点赞行为
   * @type {Object}
   */
  likePostAction: {
    id: 'likePost',
    type: 'function',

    /**
     * 处理日记点赞请求
     */
    handler: function (id, target, eventData, data) {
      const $button = $(eventData.element);
      const postId = data.postId;

      if (!postId) {
        kiko.ui.notify.error('参数错误');
        return;
      }

      // 防止重复点击
      if ($button.hasClass('liking')) {
        return;
      }

      $button.addClass('liking');
      const $icon = $button.find('i');
      const $count = $button.find('span');

      // 发送点赞请求
      $.ajax({
        url: `${Social.baseUrl}/api/post/like`,
        type: 'POST',
        contentType: 'application/json',
        data: JSON.stringify({
          post_id: parseInt(postId),
        }),
        dataType: 'json',
      })
      .done((response) => {
        if (response.code === 0 && response.data) {
          // 更新点赞状态
          const isLiked = response.data.is_liked;
          const likeCount = response.data.like_count || 0;

          // 更新图标和样式类
          if (isLiked) {
            $icon
              .removeClass('mdi-thumb-up-outline')
              .addClass('mdi-thumb-up like-icon-liked');
          } else {
            $icon
              .removeClass('mdi-thumb-up like-icon-liked')
              .addClass('mdi-thumb-up-outline');
          }

          // 更新计数
          $count.text(likeCount);

          // 显示成功提示
          const message = isLiked ? '点赞成功' : '取消点赞';
          kiko.ui.notify.success(message);
        } else {
          kiko.ui.notify.error(response.message || '操作失败');
        }
      })
      .fail((xhr, status, error) => {
        kiko.ui.notify.error('网络错误，请重试');
      })
      .always(() => {
        $button.removeClass('liking');
      });
    },
  },

  /**
   * 日记本点赞行为
   * @type {Object}
   */
  likeDiaryAction: {
    id: 'likeDiary',
    type: 'function',

    /**
     * 处理日记本点赞请求
     */
    handler: function (id, target, eventData, data) {
      const $button = $(eventData.element);
      const diaryId = data.diaryId;

      if (!diaryId) {
        kiko.ui.notify.error('参数错误');
        return;
      }

      // 防止重复点击
      if ($button.hasClass('liking')) {
        return;
      }

      $button.addClass('liking');
      const $icon = $button.find('i');
      const $count = $button.find('span');

      // 发送点赞请求
      $.ajax({
        url: `${Social.baseUrl}/api/diary/like`,
        type: 'POST',
        contentType: 'application/json',
        data: JSON.stringify({
          diary_id: parseInt(diaryId),
        }),
        dataType: 'json',
      })
      .done((response) => {
        if (response.code === 0 && response.data) {
          // 更新点赞状态
          const isLiked = response.data.is_liked;
          const likeCount = response.data.like_count || 0;

          // 更新图标和样式类
          if (isLiked) {
            $icon
              .removeClass('mdi-thumb-up-outline')
              .addClass('mdi-thumb-up like-icon-liked');
          } else {
            $icon
              .removeClass('mdi-thumb-up like-icon-liked')
              .addClass('mdi-thumb-up-outline');
          }

          // 更新计数
          $count.text(likeCount);

          // 显示成功提示
          const message = isLiked ? '点赞成功' : '取消点赞';
          kiko.ui.notify.success(message);
        } else {
          kiko.ui.notify.error(response.message || '操作失败');
        }
      })
      .fail((xhr, status, error) => {
        kiko.ui.notify.error('网络错误，请重试');
      })
      .always(() => {
        $button.removeClass('liking');
      });
    },
  },

  // 评论模块
  initCommentsModule: function () {
    const commentsSection = this._getCommentsContainer();
    if (!commentsSection.length) {
      return;
    }

    const socialSetting = commentsSection.data('social-setting');
    if (socialSetting !== 1) {
      return;
    }

    this._renderCommentsContainer();
    this.loadComments();
  },

  _getCommentsContainer: function () {
    return $('#comments-section');
  },

  _renderCommentsContainer: function () {
    const container = this._getCommentsContainer();
    const loggedIn = UserState.isLoggedIn();
    const placeholder = kiko.i18n.getMessage('comments.placeholder');
    const submitText = kiko.i18n.getMessage('comments.submit');
    const loginToCommentText = kiko.i18n.getMessage('comments.login_to_comment');

    const formHtml = loggedIn
      ? `
      <div class="mb-4">
        <textarea id="comment-textarea" class="form-control" rows="3" placeholder="${placeholder}"></textarea>
        <div class="text-right mt-2">
          <button class="btn" data-action="submitComment">${submitText}</button>
        </div>
      </div>
    `
      : `
      <div class="text-center p-4 border rounded">
        <p>${loginToCommentText}</p>
      </div>
    `;

    const listHtml = '<div id="comments-list"></div>';
    const paginationHtml = '<div id="comments-pagination" class="text-center mt-4"></div>';

    container.html(`
      <h3 class="text-lg">${kiko.i18n.getMessage('comments.title')}</h3>
      ${listHtml}
      ${paginationHtml}
      ${formHtml}
    `);
  },

  _formatDateTime: function (dateString) {
    const date = new Date(dateString);
    if (isNaN(date)) {
      return dateString; // Return original string if date is invalid
    }
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const day = String(date.getDate()).padStart(2, '0');
    const hours = String(date.getHours()).padStart(2, '0');
    const minutes = String(date.getMinutes()).padStart(2, '0');
    return `${year}-${month}-${day} ${hours}:${minutes}`;
  },

  loadComments: function (page = 1) {
    const container = this._getCommentsContainer();
    const postId = container.data('post-id');

    $.ajax({
      url: `/api/comment/list`,
      type: 'GET',
      data: { post_id: postId, page: page, limit: 10 },
      dataType: 'json',
    })
      .done((response) => {
        if (response.code === 0 && response.data) {
          this._renderCommentsList(response.data.comments);
          this._updatePagination(response.data.pagination);
        } else {
          $('#comments-list').html(`<p class="text-center text-gray-500">${kiko.i18n.getMessage('comments.no_comments')}</p>`);
        }
      })
      .fail(() => {
        $('#comments-list').html(`<p class="text-center text-red-500">Error loading comments.</p>`);
      });
  },

  _renderCommentsList: function (comments) {
    const listContainer = $('#comments-list');
    if (comments.length === 0) {
      listContainer.html(`<p class="text-center text-gray-500">${kiko.i18n.getMessage('comments.no_comments')}</p>`);
      return;
    }

    let commentsHtml = '';
    const currentUserId = UserState.getUserId();

    comments.forEach(comment => {
      const canDelete = currentUserId === comment.add_user_id;
      const deleteButton = canDelete
        ? `<button class="btn btn-sm btn-danger" data-action="deleteComment" data-comment-id="${comment.comment_id}">${kiko.i18n.getMessage('comments.delete')}</button>`
        : '';

      commentsHtml += `
        <div class="comment-item p-3" id="comment-${comment.comment_id}">
          <div class="flex items-center mb-2">
            <a href="/u/${comment.add_user_id}" class="font-semibold text-gray-800 dark:text-gray-200 no-underline hover:underline">${comment.add_user_name}</a>
            <p class="text-sm text-gray-500 dark:text-gray-400 ml-4">${this._formatDateTime(comment.gmt_create)}</p>
          </div>
          <p class="text-gray-800 dark:text-gray-300">${comment.content}</p>
          <div class="text-right mt-2">
            ${deleteButton}
          </div>
        </div>
      `;
    });

    listContainer.html(commentsHtml);
  },

  _updatePagination: function (pagination) {
    const paginationContainer = $('#comments-pagination');
    if (!pagination || pagination.total_pages <= 1) {
      paginationContainer.empty();
      return;
    }

    const loadMoreButton = `<button class="btn" data-action="loadMoreComments" data-next-page="${pagination.next_page}">${kiko.i18n.getMessage('comments.load_more')}</button>`;
    paginationContainer.html(loadMoreButton);
  },

  submitCommentAction: {
    id: 'submitComment',
    type: 'function',
    handler: function () {
      const postId = Social._getCommentsContainer().data('post-id');
      const content = $('#comment-textarea').val().trim();

      if (!content) {
        kiko.ui.notify.warning(kiko.i18n.getMessage('comments.placeholder'));
        return;
      }

      const sanitizedContent = $('<div>').text(content).html();

      $.ajax({
        url: '/api/comment/create',
        type: 'POST',
        contentType: 'application/json',
        data: JSON.stringify({
          post_id: postId,
          content: sanitizedContent,
          parent_id: 0,
        }),
        dataType: 'json',
      })
        .done((response) => {
          if (response.code === 0 && response.data) {
            kiko.ui.notify.success(kiko.i18n.getMessage('comments.submit_success'));
            
            const listContainer = $('#comments-list');
            // Remove 'no comments' placeholder if it exists
            if (listContainer.find('p').length > 0 && listContainer.children().length === 1) {
                listContainer.empty();
            }

            const comment = response.data;
            const currentUserId = UserState.getUserId();
            const canDelete = currentUserId === comment.add_user_id;
            const deleteButton = canDelete
              ? `<button class="btn btn-sm btn-danger" data-action="deleteComment" data-comment-id="${comment.comment_id}">${kiko.i18n.getMessage('comments.delete')}</button>`
              : '';

            const newCommentHtml = `
              <div class="comment-item p-3" id="comment-${comment.comment_id}">
                <div class="flex items-center mb-2">
                  <a href="/u/${comment.add_user_id}" class="font-semibold text-gray-800 dark:text-gray-200 no-underline hover:underline">${comment.add_user_name}</a>
                  <p class="text-sm text-gray-500 dark:text-gray-400 ml-4">${Social._formatDateTime(comment.gmt_create)}</p>
                </div>
                <p class="text-gray-800 dark:text-gray-300">${comment.content}</p>
                <div class="text-right mt-2">
                  ${deleteButton}
                </div>
              </div>
            `;

            listContainer.append(newCommentHtml);
            $('#comment-textarea').val('');
          } else {
            kiko.ui.notify.error(response.message);
          }
        })
        .fail(() => {
          kiko.ui.notify.error('Network error');
        });
    },
  },

  deleteCommentAction: {
    id: 'deleteComment',
    type: 'function',
    handler: function (id, target, eventData, data) {
      const commentId = data.commentId;
      if (!confirm(kiko.i18n.getMessage('comments.delete_confirm'))) {
        return;
      }

      $.ajax({
        url: '/api/comment/delete',
        type: 'POST',
        contentType: 'application/json',
        data: JSON.stringify({ comment_id: commentId }),
        dataType: 'json',
      })
        .done((response) => {
          if (response.code === 0) {
            kiko.ui.notify.success(kiko.i18n.getMessage('comments.delete_success'));
            $(`#comment-${commentId}`).remove();
            
            const listContainer = $('#comments-list');
            if (listContainer.children().length === 0) {
                listContainer.html(`<p class="text-center text-gray-500">${kiko.i18n.getMessage('comments.no_comments')}</p>`);
            }
          } else {
            kiko.ui.notify.error(response.message);
          }
        })
        .fail(() => {
          kiko.ui.notify.error('Network error');
        });
    },
  },

  loadMoreCommentsAction: {
    id: 'loadMoreComments',
    type: 'function',
    handler: function (id, target, eventData, data) {
      const nextPage = data.nextPage;
      Social.loadComments(nextPage);
    },
  },
};