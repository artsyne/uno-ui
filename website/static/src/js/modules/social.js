import { handleApiError } from '../utils/error-handler.js';

const isJsonResponse = (response) => {
  const contentType = response.headers.get('content-type') || '';
  return contentType.includes('application/json');
};

const safeJson = async (response) => {
  if (!isJsonResponse(response)) {
    return null;
  }
  try {
    return await response.json();
  } catch {
    return null;
  }
};

export function SocialMixin() {
  return {
    stats: {}, // Map of entryId -> { view_num, comment_num, like_num, liked, like_count }
    diaryLike: { liked: false, count: 0, loading: false },
    socialSetting: 2,

    getStats(entryId) {
      return (
        this.stats[entryId] || {
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

      // Listen for global entry-liked events to sync state
      window.addEventListener('entry-liked', (evt) => {
        const { entryId, liked, count } = evt.detail;
        if (this.stats[entryId]) {
          this.stats[entryId].liked = liked;
          this.stats[entryId].like_num = count;
          this.stats[entryId].like_count = count;
        }
      });

      // Listen for global entry-comment-updated events to sync state
      window.addEventListener('entry-comment-updated', (evt) => {
        const { entryId, count } = evt.detail;
        if (this.stats[entryId]) {
          this.stats[entryId].comment_num = count;
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
      } catch {
        // ignore
      }
    },

    async toggleDiaryLike(diaryId) {
      if (this.diaryLike.loading) return;
      this.diaryLike.loading = true;

      const notify = window.getUnoNotify();

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
          const { localizedMessage } = handleApiError(result);
          notify.error(localizedMessage);
        }
      } catch {
        const { localizedMessage } = handleApiError(null);
        notify.error(localizedMessage);
      } finally {
        this.diaryLike.loading = false;
      }
    },

    async loadStats() {
      // 1. Get all entry IDs from the DOM
      const entryElements = document.querySelectorAll('[data-entry-id]');
      const entryIds = Array.from(entryElements)
        .map((el) => el.dataset.entryId)
        .filter((id) => id);

      if (entryIds.length === 0) return;

      try {
        const response = await fetch('/api/entry/stats', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ entryids: entryIds.join(',') }),
        });

        if (!response.ok) {
          return;
        }

        const statsResult = await safeJson(response, '/api/entry/stats');
        if (!statsResult) return;

        if (statsResult.success && statsResult.data && statsResult.data.stats) {
          statsResult.data.stats.forEach((stat) => {
            if (!this.stats[stat.entry_id]) this.stats[stat.entry_id] = {};
            Object.assign(this.stats[stat.entry_id], {
              view_num: stat.view_num || 0,
              comment_num: stat.comment_num || 0,
              like_num: stat.like_num || 0,
              like_count: stat.like_num || 0, // normalized for LikeButton
            });
          });
        }

        // 3. Fetch like status
        const likeResponse = await fetch('/api/entry/like/status', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ entryids: entryIds.join(',') }),
        });
        const likeResult = await safeJson(
          likeResponse,
          '/api/entry/like/status',
        );
        if (!likeResult) return;

        if (likeResult.success && likeResult.data && likeResult.data.status) {
          likeResult.data.status.forEach((status) => {
            if (!this.stats[status.entry_id]) this.stats[status.entry_id] = {};
            this.stats[status.entry_id].liked = status.liked;
          });
        }
      } catch {
        // ignore
      }
    },

    // Like action wrapper for child components (if needed by parent)
    async toggleLike(entryId) {
      const id = parseInt(entryId);
      if (!this.stats[id]) return;

      const notify = window.getUnoNotify();

      try {
        const response = await fetch('/api/entry/like', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ entry_id: id }),
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
          const { localizedMessage } = handleApiError(result);
          notify.error(localizedMessage);
        }
      } catch {
        const { localizedMessage } = handleApiError(null);
        notify.error(localizedMessage);
      }
    },
  };
}

export function LikeButton(initialCount, initialLiked, entryId) {
  return {
    count: initialCount,
    liked: initialLiked,
    entryId: entryId,
    loading: false,

    async toggle() {
      if (this.loading) return;
      this.loading = true;

      const notify = window.getUnoNotify();

      try {
        const response = await fetch('/api/entry/like', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({ entry_id: parseInt(this.entryId) }),
        });

        const result = await response.json();
        if (result.success && result.data) {
          this.liked = result.data.is_liked;
          this.count = result.data.like_count;
          notify.success(
            this.liked ? 'social.like.success' : 'social.like.cancel',
          );
          // Dispatch global event to sync with other components (e.g., list view)
          window.dispatchEvent(
            new CustomEvent('entry-liked', {
              detail: {
                entryId: parseInt(this.entryId),
                liked: this.liked,
                count: this.count,
              },
            }),
          );
        } else {
          const { localizedMessage } = handleApiError(result);
          notify.error(localizedMessage);
        }
      } catch {
        const { localizedMessage } = handleApiError(null);
        notify.error(localizedMessage);
      } finally {
        this.loading = false;
      }
    },
  };
}
export function CommentsMixin(entryId) {
  return {
    entryId: entryId,
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
      const t = (key) => this.$store?.i18n?.t(key) || key;

      if (!this.entryId) {
        wrapper.innerHTML = `<div class="text-center text-gray-400 py-8">${t('comment.disabled')}</div>`;
        return;
      }

      this.loadingComments = true;
      this.currentPage = pageIndex;

      try {
        const response = await fetch(
          `/api/comment/list?entry_id=${encodeURIComponent(
            this.entryId,
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
            const { localizedMessage } = handleApiError(json);
            wrapper.innerHTML = `<div class="text-center text-red-500 py-8">${localizedMessage}</div>`;
          } catch {
            wrapper.innerHTML = `<div class="text-center text-red-500 py-8">${t('comment.load_fail')}</div>`;
          }
          return;
        }

        if (!response.ok) {
          wrapper.innerHTML = `<div class="text-center text-red-500 py-8">${t('comment.load_fail')} (Status ${response.status})</div>`;
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

        // Dispatch global event to sync with other components (e.g., list view)
        if (totalCountHeader && this.entryId) {
          window.dispatchEvent(
            new CustomEvent('entry-comment-updated', {
              detail: {
                entryId: parseInt(this.entryId),
                count: parseInt(totalCountHeader),
              },
            }),
          );
        }

        if (window.Alpine?.initTree) {
          window.Alpine.initTree(wrapper);
        }
      } catch {
        wrapper.innerHTML = `<div class="text-center text-red-500 py-8">${t('comment.load_fail')}</div>`;
      } finally {
        this.loadingComments = false;
      }
    },

    async submit() {
      const notify = window.getUnoNotify();
      const valid = window.getUnoValid();

      if (!notify || !valid) return;

      const root = this.$root || document;
      const form = root.querySelector('#comment_editor_form');
      if (!form) return;

      // Use Uno Validation engine
      // This handles showing errors on the field automatically
      if (!valid.validateForm(form)) {
        return;
      }

      // Logic Note:
      // We removed the `this.submitting` check here intentionally to rely purely on the UI mask
      // (handled by setFormLoading) for preventing double submission.
      // This ensures a cleaner separation of concerns where UI state (mask) dictates interaction.

      const { setFormLoading } = await import('./form.js');
      setFormLoading(form, true);

      try {
        const response = await fetch('/api/comment/create', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            entry_id: parseInt(this.entryId),
            content: this.content,
            parent_id: 0,
          }),
        });

        const result = await response.json();

        if (result.success) {
          notify.success('comment.submit_success');
          this.content = '';
          await this.loadComments(1);
        } else {
          const { localizedMessage } = handleApiError(result);
          notify.error(localizedMessage);
        }
      } catch (e) {
        if (process.env.NODE_ENV === 'development') {
          console.error('Submit comment error:', e);
        }
        const { localizedMessage } = handleApiError(null);
        notify.error(localizedMessage);
      } finally {
        // Removed this.submitting = false;
        const { setFormLoading } = await import('./form.js');
        setFormLoading(form, false);
      }
    },

    async deleteComment(commentId) {
      const uno = window.getUno();

      if (
        !(await uno.box.confirm(this.$store.i18n.t('comment.delete_confirm')))
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
          notify.success('comment.delete_success');
          await this.loadComments(this.currentPage);
        } else {
          const { localizedMessage } = handleApiError(result);
          notify.error(localizedMessage);
        }
      } catch {
        const { localizedMessage } = handleApiError(null);
        notify.error(localizedMessage);
      }
    },
  };
}
