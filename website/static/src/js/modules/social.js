import { Notify } from '../utils/uno-shim.js';

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
          `/api/diary/like/status?diaryid=${diaryId}`,
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
          Notify.success(
            result.data.is_liked
              ? this.$store.i18n.t('social.like.success')
              : this.$store.i18n.t('social.like.cancel'),
          );
        } else {
          Notify.error(result.message || this.$store.i18n.t('common.op_fail'));
        }
      } catch (e) {
        console.error('Diary like error:', e);
        Notify.error(this.$store.i18n.t('common.op_fail'));
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
          Notify.success(
            result.data.is_liked
              ? this.$store.i18n.t('social.like.success')
              : this.$store.i18n.t('social.like.cancel'),
          );
        } else {
          Notify.error(result.message || this.$store.i18n.t('common.op_fail'));
        }
      } catch (e) {
        console.error('Like error:', e);
        Notify.error(this.$store.i18n.t('common.op_fail'));
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
          Notify.success(
            this.liked
              ? this.$store.i18n.t('social.like.success')
              : this.$store.i18n.t('social.like.cancel'),
          );
        } else {
          Notify.error(result.message || this.$store.i18n.t('common.op_fail'));
        }
      } catch (e) {
        console.error('Like error:', e);
        Notify.error(this.$store.i18n.t('common.op_fail'));
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

    async submit() {
      if (!this.content.trim()) {
        Notify.warning(this.$store.i18n.t('comments.placeholder'));
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
          Notify.success(this.$store.i18n.t('comments.submit_success'));
          this.content = '';
          // Trigger HTMX reload of the comments list
          const listWrapper = document.getElementById('comments-list-wrapper');
          if (listWrapper) {
            htmx.trigger(listWrapper, 'reload-comments');
          }
        } else {
          Notify.error(result.message || this.$store.i18n.t('common.op_fail'));
        }
      } catch (e) {
        console.error('Submit comment error:', e);
        Notify.error(this.$store.i18n.t('common.op_fail'));
      } finally {
        this.submitting = false;
      }
    },

    async deleteComment(commentId) {
      if (!confirm(this.$store.i18n.t('comments.delete_confirm'))) {
        return;
      }

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
          Notify.success(this.$store.i18n.t('comments.delete_success'));
          // Trigger HTMX reload of the comments list
          const listWrapper = document.getElementById('comments-list-wrapper');
          if (listWrapper) {
            htmx.trigger(listWrapper, 'reload-comments');
          }
        } else {
          Notify.error(result.message || this.$store.i18n.t('common.op_fail'));
        }
      } catch (e) {
        console.error('Delete comment error:', e);
        Notify.error(this.$store.i18n.t('common.op_fail'));
      }
    },
  };
}
