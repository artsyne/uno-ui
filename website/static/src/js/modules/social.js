/**
 * Social Module - 社交功能组件
 * 
 * 包含：批量统计加载器、点赞组件、评论组件
 * 设计原则：尽可能利用 Alpine.js 响应式和 HTMX 原生能力
 */

import { handleApiError } from '../utils/error-handler.js';

// ============ 常量定义 ============
const EVENTS = {
  STATS_LOADED: 'batch-stats-loaded',
  LIKE_STATUS_LOADED: 'batch-like-status-loaded',
  ENTRY_LIKED: 'entry-liked',
  COMMENT_UPDATED: 'entry-comment-updated',
  SOCIAL_CHANGED: 'social-setting-changed',
};

// ============ 工具函数 ============
const isJsonResponse = (res) => (res.headers.get('content-type') || '').includes('application/json');

const safeJson = async (res) => {
  if (!isJsonResponse(res)) return null;
  try { return await res.json(); } catch { return null; }
};

const getNotify = () => window.getUnoNotify?.();
const getUno = () => window.getUno?.();

// ============ 批量统计管理器 ============
const batchStatsManager = {
  pendingIds: new Set(),
  loadedIds: new Set(),
  cache: {},
  timer: null,
  loading: false,
  observer: null,
  observed: new Map(),

  // 获取缓存
  get(id) {
    return this.cache[String(id)];
  },

  // 更新缓存
  set(id, data) {
    const key = String(id);
    this.cache[key] = { ...this.cache[key], ...data };
  },

  // 清除缓存
  clear(id) {
    const key = String(id);
    delete this.cache[key];
    this.loadedIds.delete(key);
  },

  // 添加到加载队列
  add(id) {
    const key = String(id);
    if (this.loadedIds.has(key) || this.pendingIds.has(key)) return;
    this.pendingIds.add(key);
    if (this.timer) clearTimeout(this.timer);
    this.timer = setTimeout(() => this.load(), 100);
  },

  // 观察元素可见性
  observe(el, id) {
    if (!el || !id) return;
    if (!this.observer) {
      this.observer = new IntersectionObserver(
        (entries) => {
          entries.forEach((entry) => {
            if (entry.isIntersecting) {
              const eid = this.observed.get(entry.target);
              if (eid) {
                this.add(eid);
                this.observer.unobserve(entry.target);
                this.observed.delete(entry.target);
              }
            }
          });
        },
        { rootMargin: '200px', threshold: [0, 0.1, 0.5, 1] }
      );
    }
    this.observed.set(el, String(id));
    this.observer.observe(el);
  },

  // 取消观察
  unobserve(el) {
    if (this.observer && el) {
      this.observer.unobserve(el);
      this.observed.delete(el);
    }
  },

  // 批量加载
  async load() {
    if (this.loading || this.pendingIds.size === 0) return;
    this.loading = true;
    const ids = Array.from(this.pendingIds);
    this.pendingIds.clear();

    try {
      const [statsRes, likeRes] = await Promise.all([
        fetch('/api/entry/stats', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ entryids: ids.join(',') }),
        }),
        fetch('/api/entry/like/status', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ entryids: ids.join(',') }),
        }),
      ]);

      const statsMap = {};
      const statsResult = await safeJson(statsRes);
      if (statsResult?.success && statsResult.data?.stats) {
        statsResult.data.stats.forEach((s) => {
          const key = String(s.entry_id);
          statsMap[key] = { view_num: s.view_num || 0, comment_num: s.comment_num || 0, like_num: s.like_num || 0 };
          this.set(key, statsMap[key]);
        });
      }

      const statusMap = {};
      const likeResult = await safeJson(likeRes);
      if (likeResult?.success && likeResult.data?.status) {
        likeResult.data.status.forEach((s) => {
          const key = String(s.entry_id);
          statusMap[key] = s.liked;
          this.set(key, { liked: s.liked });
        });
      }

      // 派发事件
      window.dispatchEvent(new CustomEvent(EVENTS.STATS_LOADED, { detail: { stats: statsMap } }));
      window.dispatchEvent(new CustomEvent(EVENTS.LIKE_STATUS_LOADED, { detail: { statuses: statusMap } }));

      ids.forEach((id) => this.loadedIds.add(String(id)));
    } catch (e) {
      console.error('Batch stats load failed:', e);
    } finally {
      this.loading = false;
      if (this.pendingIds.size > 0) {
        this.timer = setTimeout(() => this.load(), 50);
      }
    }
  },

  // 重置
  reset() {
    this.pendingIds.clear();
    this.loadedIds.clear();
    if (this.timer) clearTimeout(this.timer);
    this.timer = null;
  },
};

// 导出到全局
if (typeof window !== 'undefined') {
  Object.defineProperty(window, 'batchStatsManager', { value: batchStatsManager, writable: false, configurable: false });
}

export function getBatchStatsManager() {
  return batchStatsManager;
}

// ============ 全局事件监听器（单例）============
let globalListenersReady = false;
// 每个 entryId 对应一个 Set，支持多个组件实例（列表卡片 + box 详情）
const componentRegistry = new Map();

// 注册组件到 registry
function registerComponent(entryId, comp) {
  const key = String(entryId);
  if (!componentRegistry.has(key)) {
    componentRegistry.set(key, new Set());
  }
  componentRegistry.get(key).add(comp);
}

// 从 registry 注销组件
function unregisterComponent(entryId, comp) {
  const key = String(entryId);
  const set = componentRegistry.get(key);
  if (set) {
    set.delete(comp);
    if (set.size === 0) {
      componentRegistry.delete(key);
    }
  }
}

// 获取 entryId 对应的所有组件
function getComponents(entryId) {
  return componentRegistry.get(String(entryId)) || new Set();
}

function initGlobalListeners() {
  if (globalListenersReady) return;
  globalListenersReady = true;

  window.addEventListener(EVENTS.STATS_LOADED, (evt) => {
    Object.keys(evt.detail.stats).forEach((id) => {
      getComponents(id).forEach((comp) => {
        if (!comp.stats.loaded) {
          const data = evt.detail.stats[id];
          if (data) {
            comp.stats.view_num = data.view_num;
            comp.stats.comment_num = data.comment_num;
            comp.stats.like_num = data.like_num;
            comp.stats.loaded = true;
          }
        }
      });
    });
  });

  window.addEventListener(EVENTS.LIKE_STATUS_LOADED, (evt) => {
    Object.keys(evt.detail.statuses).forEach((id) => {
      getComponents(id).forEach((comp) => {
        comp.stats.liked = evt.detail.statuses[id];
      });
    });
  });

  window.addEventListener(EVENTS.ENTRY_LIKED, (evt) => {
    getComponents(evt.detail.entryId).forEach((comp) => {
      comp.stats.liked = evt.detail.liked;
      comp.stats.like_num = evt.detail.count;
    });
  });

  window.addEventListener(EVENTS.COMMENT_UPDATED, (evt) => {
    getComponents(evt.detail.entryId).forEach((comp) => {
      comp.stats.comment_num = evt.detail.count;
    });
  });

  window.addEventListener(EVENTS.SOCIAL_CHANGED, (evt) => {
    componentRegistry.forEach((set) => {
      set.forEach((comp) => {
        comp.socialSetting = parseInt(evt.detail.setting);
        if (comp.socialSetting === 1) {
          batchStatsManager.clear(comp.entryId);
          comp.initStats();
        }
      });
    });
  });
}

// ============ SocialMixin - 日记本级别社交功能 ============
export function SocialMixin() {
  return {
    stats: {},
    diaryLike: { liked: false, count: 0, loading: false },
    socialSetting: 2,

    initSocial() {
      this.socialSetting = parseInt(document.querySelector('main')?.dataset.socialSetting || '2');
      if (this.socialSetting === 1) this.loadDiaryLike();

      window.addEventListener(EVENTS.ENTRY_LIKED, (evt) => {
        if (this.stats[evt.detail.entryId]) {
          this.stats[evt.detail.entryId].liked = evt.detail.liked;
          this.stats[evt.detail.entryId].like_num = evt.detail.count;
        }
      });

      window.addEventListener(EVENTS.COMMENT_UPDATED, (evt) => {
        if (this.stats[evt.detail.entryId]) {
          this.stats[evt.detail.entryId].comment_num = evt.detail.count;
        }
      });
    },

    async loadDiaryLike() {
      const container = document.getElementById('diary-like-wrapper');
      const diaryId = container?.dataset.diaryId;
      if (!diaryId) return;
      try {
        const res = await fetch(`/api/diary/like/status?diary_id=${diaryId}`, {
          method: 'GET',
          headers: { 'Content-Type': 'application/json' },
        });
        const result = await res.json();
        if (result.success && result.data) {
          this.diaryLike.liked = result.data.is_liked;
          this.diaryLike.count = result.data.like_count;
        }
      } catch { /* ignore */ }
    },

    async toggleDiaryLike(diaryId) {
      if (this.diaryLike.loading) return;
      this.diaryLike.loading = true;
      const notify = getNotify();
      try {
        const res = await fetch('/api/diary/like', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ diary_id: parseInt(diaryId) }),
        });
        const result = await res.json();
        if (result.success && result.data) {
          this.diaryLike.liked = result.data.is_liked;
          this.diaryLike.count = result.data.like_count;
          notify?.success(result.data.is_liked ? 'social.like.success' : 'social.like.cancel');
        } else {
          notify?.error(handleApiError(result).localizedMessage);
        }
      } catch {
        notify?.error(handleApiError(null).localizedMessage);
      } finally {
        this.diaryLike.loading = false;
      }
    },
  };
}

// ============ LikeButton - 独立点赞按钮组件 ============
export function LikeButton(initialCount, initialLiked, entryId) {
  return {
    count: initialCount,
    liked: initialLiked,
    entryId,
    loading: false,

    async toggle() {
      if (this.loading) return;
      this.loading = true;
      const notify = getNotify();
      try {
        const res = await fetch('/api/entry/like', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ entry_id: parseInt(this.entryId) }),
        });
        const result = await res.json();
        if (result.success && result.data) {
          this.liked = result.data.is_liked;
          this.count = result.data.like_count;
          notify?.success(this.liked ? 'social.like.success' : 'social.like.cancel');
          window.dispatchEvent(new CustomEvent(EVENTS.ENTRY_LIKED, {
            detail: { entryId: parseInt(this.entryId), liked: this.liked, count: this.count },
          }));
        } else {
          notify?.error(handleApiError(result).localizedMessage);
        }
      } catch {
        notify?.error(handleApiError(null).localizedMessage);
      } finally {
        this.loading = false;
      }
    },
  };
}

// ============ CommentsMixin - 评论组件 ============
export function CommentsMixin(entryId) {
  return {
    entryId,
    content: '',
    loadingComments: false,
    currentPage: 1,

    init() {
      this.loadComments(1);
    },

    async loadComments(page = 1) {
      const wrapper = this.$root?.querySelector('#comments-list-wrapper');
      if (!wrapper) return;

      this.currentPage = page;
      this.loadingComments = true;
      const t = (key) => this.$store?.i18n?.t(key) || key;

      try {
        const res = await fetch(
          `/api/comment/list?entry_id=${encodeURIComponent(this.entryId)}&page_index=${page}&_t=${Date.now()}`,
          { headers: { 'HX-Request': 'true', 'X-Requested-With': 'XMLHttpRequest' } }
        );

        const text = await res.text();

        // JSON 错误响应处理
        if (res.headers.get('content-type')?.includes('application/json') || (text.trim().startsWith('{') && text.trim().endsWith('}'))) {
          try {
            const json = JSON.parse(text);
            wrapper.innerHTML = `<div class="text-center text-red-500 py-8">${handleApiError(json).localizedMessage}</div>`;
          } catch {
            wrapper.innerHTML = `<div class="text-center text-red-500 py-8">${t('comment.load_fail')}</div>`;
          }
          return;
        }

        if (!res.ok) {
          wrapper.innerHTML = `<div class="text-center text-red-500 py-8">${t('comment.load_fail')} (${res.status})</div>`;
          return;
        }

        wrapper.innerHTML = text;

        // 更新评论数
        const totalCount = res.headers.get('X-Total-Count');
        const countEl = document.getElementById('post-comment-count');
        if (countEl && totalCount) {
          const icon = countEl.querySelector('i');
          countEl.innerHTML = '';
          if (icon) countEl.appendChild(icon);
          countEl.append(` ${totalCount}`);
        }

        // 同步事件
        if (totalCount && this.entryId) {
          window.dispatchEvent(new CustomEvent(EVENTS.COMMENT_UPDATED, {
            detail: { entryId: parseInt(this.entryId), count: parseInt(totalCount) },
          }));
        }

        // 初始化 Alpine 组件
        window.Alpine?.initTree?.(wrapper);
      } catch {
        wrapper.innerHTML = `<div class="text-center text-red-500 py-8">${t('comment.load_fail')}</div>`;
      } finally {
        this.loadingComments = false;
      }
    },

    async submit() {
      const notify = getNotify();
      const valid = window.getUnoValid?.();
      if (!notify || !valid) return;

      const form = this.$root?.querySelector('#comment_editor_form');
      if (!form || !valid.validateForm(form)) return;

      const { setFormLoading } = await import('./form.js');
      setFormLoading(form, true);

      try {
        const res = await fetch('/api/comment/create', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ entry_id: parseInt(this.entryId), content: this.content, parent_id: 0 }),
        });
        const result = await res.json();

        if (result.success) {
          notify.success('comment.submit_success');
          this.content = '';
          await this.loadComments(1);
        } else {
          notify.error(handleApiError(result).localizedMessage);
        }
      } catch {
        notify.error(handleApiError(null).localizedMessage);
      } finally {
        setFormLoading(form, false);
      }
    },

    async deleteComment(commentId) {
      const uno = getUno();
      if (!uno || !(await uno.box.confirm(this.$store.i18n.t('comment.delete_confirm')))) return;

      const notify = getNotify();
      try {
        const res = await fetch('/api/comment/delete', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ comment_id: parseInt(commentId) }),
        });
        const result = await res.json();
        if (result.success) {
          notify.success('comment.delete_success');
          await this.loadComments(this.currentPage);
        } else {
          notify.error(handleApiError(result).localizedMessage);
        }
      } catch {
        notify.error(handleApiError(null).localizedMessage);
      }
    },
  };
}

// ============ entryStats - 日记卡片统计组件 ============
export function entryStats(entryId, socialSetting = null) {
  return {
    entryId,
    socialSetting,
    stats: { loaded: false, loading: false, view_num: 0, comment_num: 0, like_num: 0, liked: false },

    async initStats() {
      // 检查社交设置
      if (this.socialSetting !== 1) {
        const globalSetting = parseInt(document.querySelector('main')?.dataset.socialSetting || '2');
        if (globalSetting !== 1) return;
        this.socialSetting = 1;
      }

      // 初始化全局监听器并注册组件
      initGlobalListeners();
      registerComponent(this.entryId, this);

      // 检查缓存
      const cached = batchStatsManager.get(this.entryId);
      if (cached) {
        Object.assign(this.stats, cached, { loaded: true });
        if (cached.liked !== undefined) this.stats.liked = cached.liked;
      }

      // Box 内直接加载，列表中懒加载
      const isInBox = this.$el?.closest('.uno-box');
      if (isInBox) {
        await this.loadSingle();
      } else {
        await this.loadWhenVisible();
      }
    },

    async loadWhenVisible() {
      if (this.stats.loaded) return;
      const cached = batchStatsManager.get(this.entryId);
      if (cached) {
        Object.assign(this.stats, cached, { loaded: true });
        return;
      }
      batchStatsManager.observe(this.$el, this.entryId);
    },

    async loadSingle() {
      if (this.stats.loading || this.stats.loaded) return;
      this.stats.loading = true;

      try {
        const [statsRes, likeRes] = await Promise.all([
          fetch('/api/entry/stats', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ entryids: String(this.entryId) }),
          }),
          fetch('/api/entry/like/status', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ entryids: String(this.entryId) }),
          }),
        ]);

        const statsResult = await safeJson(statsRes);
        if (statsResult?.success && statsResult.data?.stats?.[0]) {
          const s = statsResult.data.stats[0];
          this.stats.view_num = s.view_num || 0;
          this.stats.comment_num = s.comment_num || 0;
          this.stats.like_num = s.like_num || 0;
          batchStatsManager.set(this.entryId, { view_num: this.stats.view_num, comment_num: this.stats.comment_num, like_num: this.stats.like_num });
        }

        const likeResult = await safeJson(likeRes);
        if (likeResult?.success && likeResult.data?.status?.[0]) {
          this.stats.liked = likeResult.data.status[0].liked;
          batchStatsManager.set(this.entryId, { liked: this.stats.liked });
        }

        this.stats.loaded = true;
      } catch (e) {
        console.error('Load stats failed:', e);
        this.stats.loaded = true;
      } finally {
        this.stats.loading = false;
      }
    },

    async toggleLike() {
      if (this.stats.loading) return;

      // 乐观更新
      const oldLiked = this.stats.liked;
      const oldCount = this.stats.like_num;
      this.stats.liked = !oldLiked;
      this.stats.like_num = oldLiked ? oldCount - 1 : oldCount + 1;
      batchStatsManager.set(this.entryId, { liked: this.stats.liked, like_num: this.stats.like_num });

      try {
        const res = await fetch('/api/entry/like', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ entry_id: parseInt(this.entryId) }),
        });
        const result = await res.json();
        if (result.success && result.data) {
          this.stats.liked = result.data.is_liked;
          this.stats.like_num = result.data.like_count;
          batchStatsManager.set(this.entryId, { liked: this.stats.liked, like_num: this.stats.like_num });
          window.dispatchEvent(new CustomEvent(EVENTS.ENTRY_LIKED, {
            detail: { entryId: parseInt(this.entryId), liked: this.stats.liked, count: this.stats.like_num },
          }));
        } else {
          throw new Error('API error');
        }
      } catch (e) {
        // 回滚
        this.stats.liked = oldLiked;
        this.stats.like_num = oldCount;
        batchStatsManager.set(this.entryId, { liked: oldLiked, like_num: oldCount });
        console.error('Toggle like failed:', e);
      }
    },

    destroy() {
      batchStatsManager.unobserve(this.$el);
      unregisterComponent(this.entryId, this);
    },
  };
}