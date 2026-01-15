const isJsonResponse = (response) => {
  const contentType = response.headers.get('content-type') || '';
  return contentType.includes('application/json');
};

const safeJson = async (response, context) => {
  if (!isJsonResponse(response)) {
    console.warn(
      '[Social] Non-JSON response for',
      context,
      response.headers.get('content-type'),
    );
    return null;
  }
  try {
    return await response.json();
  } catch (e) {
    console.warn('[Social] Failed to parse JSON for', context, e);
    return null;
  }
};

export function SocialMixin() {
  return {
    stats: {}, // Map of postId -> { view_num, comment_num, like_num, liked, like_count }
    diaryLike: { liked: false, count: 0, loading: false },
    socialSetting: 2,

    getStats(postId) {
      return (
        this.stats[postId] || {
          view_num: 0,
          comment_num: 0,
          like_num: 0,
          liked: false,
        }
      );
    },

    initSocial() {
      this.socialSetting = parseInt(
        document.querySelector('main')?.dataset.socialSetting || '2',
      );

      if (this.socialSetting === 1) {
        this.loadStats();
        this.loadDiaryLike();
      }

      // Reload stats after HTMX content swap
      document.body.addEventListener('htmx:afterSwap', (evt) => {
        if (
          evt.detail.target.id === 'diary-content' &&
          this.socialSetting === 1
        ) {
          this.loadStats();
        }
      });
    },

    async loadDiaryLike() {
      const container = document.getElementById('diary-like-wrapper');
      if (!container) return;
      const diaryId = container.dataset.diaryId;
      if (!diaryId) return;

      try {
        const response = await fetch(
          `/api/diary/like/status?diary_id=${diaryId}`,
          {
            method: 'GET',
            headers: { 'Content-Type': 'application/json' },
          },
        );
        const result = await response.json();
        if (result.success && result.data) {
          this.diaryLike.liked = result.data.is_liked;
          this.diaryLike.count = result.data.like_count;
        }
      } catch (e) {
        console.error('[Social] Error loading diary like:', e);
      }
    },

    async toggleDiaryLike(diaryId) {
      if (this.diaryLike.loading) return;
      this.diaryLike.loading = true;

      const notify = window.getUnoNotify();
      if (!notify) {
        this.diaryLike.loading = false;
        return;
      }

      try {
        const response = await fetch('/api/diary/like', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ diary_id: parseInt(diaryId) }),
        });
        const result = await response.json();

        if (result.success && result.data) {
          this.diaryLike.liked = result.data.is_liked;
          this.diaryLike.count = result.data.like_count;
          notify.success(
            result.data.is_liked ? 'social.like.success' : 'social.like.cancel',
          );
        } else {
          notify.error(result.message || 'common.op_fail');
        }
      } catch (e) {
        console.error('Diary like error:', e);
        notify.error('common.op_fail');
      } finally {
        this.diaryLike.loading = false;
      }
    },

    async loadStats() {
      // 1. Get all post IDs from the DOM
      const postElements = document.querySelectorAll('[data-post-id]');
      const postIds = Array.from(postElements)
        .map((el) => el.dataset.postId)
        .filter((id) => id);

      if (postIds.length === 0) return;

      try {
        const response = await fetch('/api/post/stats', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ postids: postIds.join(',') }),
        });

        if (!response.ok) {
          console.warn('[Social] Failed to fetch stats:', response.status);
          return;
        }

        const statsResult = await safeJson(response, '/api/post/stats');
        if (!statsResult) return;

        if (statsResult.success && statsResult.data && statsResult.data.stats) {
          statsResult.data.stats.forEach((stat) => {
            if (!this.stats[stat.post_id]) this.stats[stat.post_id] = {};
            Object.assign(this.stats[stat.post_id], {
              view_num: stat.view_num || 0,
              comment_num: stat.comment_num || 0,
              like_num: stat.like_num || 0,
              like_count: stat.like_num || 0, // normalized for LikeButton
            });
          });
        }

        // 3. Fetch like status
        const likeResponse = await fetch('/api/post/like/status', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ postids: postIds.join(',') }),
        });
        const likeResult = await safeJson(
          likeResponse,
          '/api/post/like/status',
        );
        if (!likeResult) return;

        if (likeResult.success && likeResult.data && likeResult.data.status) {
          likeResult.data.status.forEach((status) => {
            if (!this.stats[status.postid]) this.stats[status.postid] = {};
            this.stats[status.postid].liked = status.liked;
          });
        }
      } catch (e) {
        console.error('[Social] Error loading stats:', e);
      }
    },

    // Like action wrapper for child components (if needed by parent)
    async toggleLike(postId) {
      const id = parseInt(postId);
      if (!this.stats[id]) return;

      const notify = window.getUnoNotify();
      if (!notify) return;

      const isLiked = this.stats[id].liked;

      try {
        const response = await fetch('/api/post/like', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ post_id: id }),
        });
        const result = await response.json();

        if (result.success && result.data) {
          const stats = this.stats[id];
          if (stats) {
            stats.liked = result.data.is_liked;
            stats.like_count = result.data.like_count;
            stats.like_num = result.data.like_count; // sync both
          }
          notify.success(
            result.data.is_liked ? 'social.like.success' : 'social.like.cancel',
          );
        } else {
          notify.error(result.message || 'common.op_fail');
        }
      } catch (e) {
        console.error('Like error:', e);
        notify.error('common.op_fail');
      }
    },
  };
}

export function LikeButton(initialCount, initialLiked, postId) {
  return {
    count: initialCount,
    liked: initialLiked,
    postId: postId,
    loading: false,

    async toggle() {
      if (this.loading) return;
      this.loading = true;

      const notify = window.getUnoNotify();
      if (!notify) {
        this.loading = false;
        return;
      }

      try {
        const response = await fetch('/api/post/like', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({ post_id: parseInt(this.postId) }),
        });

        const result = await response.json();
        if (result.success && result.data) {
          this.liked = result.data.is_liked;
          this.count = result.data.like_count;
          notify.success(
            this.liked ? 'social.like.success' : 'social.like.cancel',
          );
        } else {
          notify.error(result.message || 'common.op_fail');
        }
      } catch (e) {
        console.error('Like error:', e);
        notify.error('common.op_fail');
      } finally {
        this.loading = false;
      }
    },
  };
}
export function CommentsMixin(postId) {
  return {
    postId: postId,
    content: '',
    submitting: false,
    loadingComments: false,
    currentPage: 1,

    init() {
      this.loadComments(1);
    },

    async loadComments(pageIndex = 1) {
      const root = this.$root || document;
      const wrapper = root.querySelector('#comments-list-wrapper');
      if (!wrapper) return;

      const countEl = document.getElementById('post-comment-count');

      if (!this.postId) {
        console.warn('[CommentsMixin] No postId provided');
        wrapper.innerHTML =
          '<div class="text-center text-gray-400 py-8">Comments disabled</div>';
        return;
      }

      this.loadingComments = true;
      this.currentPage = pageIndex;

      try {
        const response = await fetch(
          `/api/comment/list?post_id=${encodeURIComponent(
            this.postId,
          )}&page_index=${encodeURIComponent(pageIndex)}&_t=${new Date().getTime()}`,
          {
            method: 'GET',
            headers: {
              'HX-Request': 'true',
              'X-Requested-With': 'XMLHttpRequest',
            },
          },
        );

        const result = await response.text();

        // Check if response is JSON (error or unexpected)
        if (
          response.headers.get('content-type')?.includes('application/json') ||
          (result.trim().startsWith('{') && result.trim().endsWith('}'))
        ) {
          try {
            const json = JSON.parse(result);
            wrapper.innerHTML = `<div class="text-center text-red-500 py-8">${json.message || 'Failed to load comments'}</div>`;
          } catch (e) {
            wrapper.innerHTML = `<div class="text-center text-red-500 py-8">Failed to load comments</div>`;
          }
          return;
        }

        if (!response.ok) {
          wrapper.innerHTML = `<div class="text-center text-red-500 py-8">Failed to load comments (Status ${response.status})</div>`;
          return;
        }

        wrapper.innerHTML = result;

        const totalCountHeader = response.headers.get('X-Total-Count');
        if (countEl && totalCountHeader) {
          const icon = countEl.querySelector('i');
          countEl.innerHTML = '';
          if (icon) countEl.appendChild(icon);
          countEl.append(` ${totalCountHeader}`);
        }

        if (window.Alpine?.initTree) {
          window.Alpine.initTree(wrapper);
        }
      } catch (e) {
        console.error('Load comments error:', e);
        wrapper.innerHTML =
          '<div class="text-center text-red-500 py-8">Failed to load comments</div>';
      } finally {
        this.loadingComments = false;
      }
    },

    async submit() {
      const notify = window.getUnoNotify();
      const valid = window.getUnoValid();

      if (!notify || !valid) return;

      const form = document.getElementById('comment_editor_form');
      if (!form) return;

      // Use Uno Validation engine
      // This handles showing errors on the field automatically
      if (!valid.validateForm(form)) {
        return;
      }

      if (this.submitting) return;
      this.submitting = true;

      try {
        const response = await fetch('/api/comment/create', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            post_id: parseInt(this.postId),
            content: this.content,
            parent_id: 0,
          }),
        });

        const result = await response.json();

        if (result.success) {
          notify.success('comments.submit_success');
          this.content = '';
          await this.loadComments(1);
        } else {
          notify.error(result.message || 'common.op_fail');
        }
      } catch (e) {
        console.error('Submit comment error:', e);
        notify.error('common.op_fail');
      } finally {
        this.submitting = false;
      }
    },

    async deleteComment(commentId) {
      const uno = window.getUno();
      if (!uno) return;

      if (
        !(await uno.box.confirm(this.$store.i18n.t('comments.delete_confirm')))
      ) {
        return;
      }

      const notify = window.getUnoNotify();
      if (!notify) return;

      try {
        const response = await fetch('/api/comment/delete', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({ comment_id: parseInt(commentId) }),
        });

        const result = await response.json();

        if (result.success) {
          notify.success('comments.delete_success');
          await this.loadComments(this.currentPage);
        } else {
          notify.error(result.message || 'common.op_fail');
        }
      } catch (e) {
        console.error('Delete comment error:', e);
        notify.error('common.op_fail');
      }
    },
  };
}
