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
    return [
      this.likePostAction,
      this.likeDiaryAction,
    ];
  },

  /**
   * 获取社交设置值
   * @returns {number} 社交设置值，1表示开启，2表示关闭
   */
  _getSocialSetting: function () {
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
    const socialSetting = this._getSocialSetting();
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
    const socialSetting = this._getSocialSetting();
    if (socialSetting !== 1) {
      // 如果社交设置不为1，清空容器
      $('#diary-like-container').empty();
      return;
    }

    // 获取日记ID
    const $container = $('#diary-like-container');
    const diaryId = $container.data('diary-id');

    if (!diaryId) {
      console.log('未找到日记ID，无法加载点赞状态');
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
};