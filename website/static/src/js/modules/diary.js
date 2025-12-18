import { Box, Notify, Valid } from '../utils/uno-shim.js';
import { SocialMixin } from './social.js';

export function DiaryList() {
  return {
    ...SocialMixin(),

    init() {
      // Initialize social features
      this.initSocial();

      // Read isOwner from DOM
      const mainEl = document.querySelector('main');
      this.isOwner = mainEl?.dataset.isOwner === 'true';

      // Expose to global window for tooltip access
      window.openEditPost = (diaryId, postId) =>
        this.openEditPost(diaryId, postId);
      window.confirmDelete = (diaryId, postId) =>
        this.confirmDelete(diaryId, postId);
    },

    _escapeHtml(unsafe) {
      if (!unsafe) return '';
      return unsafe
        .replace(/&/g, '&amp;')
        .replace(/</g, '&lt;')
        .replace(/>/g, '&gt;')
        .replace(/"/g, '&quot;')
        .replace(/'/g, '&#039;');
    },

    _updatePostStats(post) {
      if (!this.stats[post.post_id]) {
        this.stats[post.post_id] = {};
      }
      Object.assign(this.stats[post.post_id], {
        view_num: post.view_num || 0,
        comment_num: post.comment_num || 0,
        like_num: post.like_num || 0,
        // For edit, preserve existing 'liked' status if available
        liked: this.stats[post.post_id].liked || false,
      });
    },

    _renderPostCard(post) {
      const t = (key) => this.$store.i18n.t(key);
      const date = new Date(post.add_utc_time);
      const dateStr = date.toLocaleString();

      // Owner menu
      let ownerMenu = '';
      if (this.isOwner) {
        ownerMenu = `
            <div class="post-edit">
                <a href="javascript:void(0)" x-tooltip.click.bottom="{ content: '#post-menu-${post.post_id}' }">
                    <i class="mdi mdi-dots-horizontal"></i>
                </a>
            </div>
            <div id="post-menu-${post.post_id}" class="hidden">
                <ul class="text-left">
                    <li class="py-1"><a href="#" class="block hover:text-blue-500"
                            @click.prevent="openEditPost(${post.diary_id}, ${post.post_id}); $store.tooltip.hide()">
                            <i class="mdi mdi-pencil-outline mr-2"></i> ${t('common.edit')}</a>
                    </li>
                    <li class="py-1"><a href="#" class="block hover:text-red-500"
                            @click.prevent="confirmDelete(${post.diary_id}, ${post.post_id}); $store.tooltip.hide()">
                            <i class="mdi mdi-delete-outline mr-2"></i> ${t('common.delete')}</a>
                    </li>
                </ul>
            </div>`;
      }

      // Title
      let titleHtml = '';
      if (post.title) {
        titleHtml = `<h1 class="post-title">${this._escapeHtml(post.title)}</h1>`;
      }

      return `
          <div class="card mb-6 break-inside-avoid" data-post-id="${post.post_id}">
              ${ownerMenu}
              <div class="mb-4">
                  ${titleHtml}
                  <p class="text-relaxed cursor-pointer">
                      <a href="/p/${post.post_id}" @click.prevent="openPostModal('/p/${post.post_id}')"
                          style="all:unset !important">
                          ${this._escapeHtml(post.summary)}
                      </a>
                  </p>
              </div>
              <div class="flex justify-between">
                  <div>
                      <span class="post-time">
                          ${dateStr}
                      </span>
                  </div>
                  <!-- Alpine-based Stats Container -->
                  <div class="flex items-center" x-show="stats[${post.post_id}]" style="display: none;">
                      <button class="btn-social mr-2">
                          <i class="mdi mdi-eye mr-1"></i> <span x-text="getStats(${post.post_id}).view_num"></span>
                      </button>
                      <button class="btn-social mr-2">
                          <i class="mdi mdi-comment-text-outline mr-1"></i> <span
                              x-text="getStats(${post.post_id}).comment_num"></span>
                      </button>
                      <button class="btn-social" @click="toggleLike(${post.post_id})"
                          :class="getStats(${post.post_id}).liked ? 'text-red-500' : ''">
                          <i class="mdi mr-1"
                              :class="getStats(${post.post_id}).liked ? 'mdi-thumb-up' : 'mdi-thumb-up-outline'"></i>
                          <span x-text="getStats(${post.post_id}).like_num"></span>
                      </button>
                  </div>
              </div>
          </div>
        `;
    },

    confirmDelete(diaryId, postId) {
      Box.confirm({
        content: this.$store.i18n.t('diary.delete.confirm'),
        okButtonClass: 'btn',
        onOk: () => {
          this.deletePost(diaryId, postId);
        },
      });
    },
    async deletePost(diaryId, postId) {
      try {
        const response = await fetch(`/post/${postId}?diary_id=${diaryId}`, {
          method: 'DELETE',
        });
        const result = await response.json();
        if (result.success) {
          const el = document.querySelector(`.card[data-post-id="${postId}"]`);
          if (el) el.remove();
          Notify.success(this.$store.i18n.t('diary.delete.success'));
        } else {
          Notify.error(
            result.message || this.$store.i18n.t('diary.delete.fail'),
          );
        }
      } catch (error) {
        console.error('Delete error:', error);
        Notify.error(this.$store.i18n.t('diary.delete.net_error'));
      }
    },

    // Modal Actions
    openPostModal(url) {
      const originalUrl = window.location.href;

      // Update URL to the post URL immediately
      history.pushState(null, '', url);

      Box.dialog({
        ajaxUrl: url + '?isajax=true',
        containerClass: '!border',
        widthClass: 'w-full lg:w-[90%]',
        heightClass: 'h-full',
        onContentLoad: (el) => {
          if (window.htmx) {
            htmx.process(el);
          }
        },
        onClose: () => {
          // Restore original URL when closing the modal
          history.pushState(null, '', originalUrl);
        },
      });
    },

    openCreatePost(diaryId) {
      this._openPostForm('post.create', diaryId);
    },

    openEditPost(diaryId, postId) {
      this._openPostForm('post.edit', diaryId, postId);
    },

    _openPostForm(action, diaryId, postId = null) {
      const isEdit = action === 'post.edit';
      const title = isEdit
        ? this.$store.i18n.t('diary.edit.title')
        : this.$store.i18n.t('diary.create.title');

      const boxId = Box.dialog({
        title: title,
        content: `<div class="modal-body" id="post-form-container">${this.$store.i18n.t('common.loading')}</div>`,
        widthClass: '!w-[800px]',
        heightClass: '!h-[500px]',
        showOkButton: true,
        okButtonText: this.$store.i18n.t('common.save'),
        containerClass: '!border',
        okButtonClass: 'btn',
        buttonAlign: 'right',
        onOk: () => this._handlePostSubmit(boxId, action),
      });

      const params = new URLSearchParams({
        action: action,
        ...(diaryId && { diary_id: diaryId }),
        ...(postId && { post_id: postId }),
      });

      fetch(`/form/post?${params.toString()}`)
        .then((res) => res.text())
        .then((html) => {
          const container = document.getElementById('post-form-container');
          if (container) {
            container.innerHTML = html;
          }
        })
        .catch((err) => {
          console.error('Error loading form:', err);
          const container = document.getElementById('post-form-container');
          if (container)
            container.innerHTML = `<p class="text-red-500">${this.$store.i18n.t('common.load_fail')}</p>`;
        });
    },

    async _handlePostSubmit(boxId, action) {
      const form = document.querySelector('#post-form-container form');
      if (!form) return false;

      // Uno Validation
      if (!(await Valid.validateForm(form))) {
        return false;
      }

      const formData = new FormData(form);
      const data = Object.fromEntries(formData.entries());

      try {
        // Using URLSearchParams for x-www-form-urlencoded
        const searchParams = new URLSearchParams(formData);

        // Remove action parameter if present as we use RESTful routes
        searchParams.delete('action');

        let url = '/post';
        let method = 'POST';
        const postId = formData.get('post_id');

        if (action === 'post.edit' && postId) {
          url = `/post/${postId}`;
          method = 'PUT';
        }

        const response = await fetch(url, {
          method: method,
          headers: {
            'Content-Type': 'application/x-www-form-urlencoded',
          },
          body: searchParams,
        });

        const result = await response.json();

        if (result.success) {
          Notify.success(this.$store.i18n.t('common.save_success'));
          Box.close(boxId);

          const post = result.data;
          this._updatePostStats(post);
          const postHtml = this._renderPostCard(post);

          // Check if we are on page > 1 (format: /d/{diaryId}/{page})
          const pathParts = window.location.pathname.split('/');
          let isPagegt1 = false;
          if (pathParts.length >= 4) {
            const pageIdx = parseInt(pathParts[3]);
            if (!isNaN(pageIdx) && pageIdx > 1) {
              isPagegt1 = true;
            }
          }

          if (isPagegt1) {
            const targetUrl = `/d/${post.diary_id}`;
            if (window.htmx) {
              // SPA redirect to page 1
              // We return the promise to let the caller await it (though caller might not strictly await undefined)
              // But importantly, we do NOT proceed to manual insertion.
              return new Promise((resolve) => {
                window.htmx
                  .ajax('GET', targetUrl, {
                    target: '#diary-content',
                    select: '#diary-content',
                    swap: 'outerHTML',
                  })
                  .then(() => {
                    window.history.pushState(null, '', targetUrl);
                    resolve();
                  });
              });
            } else {
              window.location.href = targetUrl;
              return;
            }
          }

          const container = document.getElementById('diary_post_container');

          // Helper to create element from HTML string
          const createElement = (html) => {
            const div = document.createElement('div');
            div.innerHTML = html.trim();
            return div.firstElementChild;
          };

          if (action === 'post.create') {
            if (container) {
              const newEl = createElement(postHtml);
              // Try to find the "Write Entry" card to insert after
              const writeEntryCard =
                document.getElementById('write-entry-card');
              if (writeEntryCard) {
                writeEntryCard.insertAdjacentElement('afterend', newEl);
              } else {
                container.prepend(newEl);
              }

              if (window.Alpine) {
                window.Alpine.initTree(newEl);
              }
            }
          } else if (action === 'post.edit') {
            const existingCard = document.querySelector(
              `.card[data-post-id="${post.post_id}"]`,
            );
            if (existingCard) {
              const newEl = createElement(postHtml);
              existingCard.replaceWith(newEl);
              if (window.Alpine) {
                window.Alpine.initTree(newEl);
              }
            }
          }

          // Return nothing (undefined) to allow Uno to close the dialog if it was waiting
          return;
        } else {
          Notify.error(
            result.message || this.$store.i18n.t('common.save_fail'),
          );
          return false;
        }
      } catch (e) {
        console.error('Submit error:', e);
        Notify.error(this.$store.i18n.t('common.net_error'));
        return false;
      }
    },

    openEditDiary(diaryId) {
      const boxId = Box.dialog({
        title: this.$store.i18n.t('diary.settings.title'), // Should use i18n
        content: `<div class="modal-body" id="diary-form-container">${this.$store.i18n.t('common.loading')}</div>`,
        widthClass: '!w-[880px]',
        heightClass: '!h-[570px]',
        showOkButton: true,
        okButtonText: this.$store.i18n.t('common.save'),
        containerClass: '!border',
        okButtonClass: 'btn',
        buttonAlign: 'right',
        onOk: () => this._handleDiarySubmit(boxId),
      });

      fetch(`/form/diary?action=diary.edit&diary_id=${diaryId}`)
        .then((res) => res.text())
        .then((html) => {
          const container = document.getElementById('diary-form-container');
          if (container) {
            container.innerHTML = html;
            // Manually populate action field since backend doesn't populate it in form/diary template
            const actionInput = container.querySelector('#action');
            if (actionInput) actionInput.value = 'diary.edit';
          }
        })
        .catch((err) => {
          const container = document.getElementById('diary-form-container');
          if (container)
            container.innerHTML = `<p class="text-red-500">${this.$store.i18n.t('common.load_fail')}</p>`;
        });
    },

    async _handleDiarySubmit(boxId) {
      const form = document.querySelector('#diary-form-container form');
      if (!form) return false;

      // Uno Validation
      if (!(await Valid.validateForm(form))) {
        return false;
      }

      const formData = new FormData(form);
      const searchParams = new URLSearchParams(formData);
      searchParams.delete('action');

      const diaryId = formData.get('diary_id');

      try {
        const response = await fetch(`/diary/${diaryId}`, {
          method: 'PUT',
          headers: {
            'Content-Type': 'application/x-www-form-urlencoded',
          },
          body: searchParams,
        });

        const result = await response.json();

        if (result.success) {
          Notify.success(this.$store.i18n.t('common.save_success'));
          Box.close(boxId);

          const diary = result.data;
          const titleEl = document.getElementById('diary_title');
          const remarkEl = document.getElementById('diary_remark');

          if (titleEl) {
            let iconHtml = '';
            if (diary.view_setting === 2) {
              iconHtml = '<i class="mdi mdi-lock-outline mx-3"></i>';
            }
            titleEl.innerHTML = this._escapeHtml(diary.title) + iconHtml;
          }
          if (remarkEl) {
            remarkEl.innerText = diary.remark || '';
          }

          // Update Theme Class and Social Setting on Main Element
          const mainEl = document.querySelector('main');
          if (mainEl) {
            if (diary.theme_id) {
              // Remove existing dt* classes
              mainEl.classList.forEach((cls) => {
                if (cls.startsWith('dt') && !isNaN(cls.substring(2))) {
                  mainEl.classList.remove(cls);
                }
              });
              // Add new theme class
              mainEl.classList.add(`dt${diary.theme_id}`);
            }
            // Update social setting data attribute
            mainEl.dataset.socialSetting = diary.social_setting;
          }

          // Update Component State
          this.socialSetting = diary.social_setting;

          // Update Social Button on Diary Detail Page
          const socialWrapper = document.getElementById('diary-like-wrapper');
          if (diary.social_setting === 1) {
            if (!socialWrapper) {
              const pageNav = document.getElementById('page-nav');
              if (pageNav && pageNav.nextElementSibling) {
                const btnContainer = pageNav.nextElementSibling;
                const span = document.createElement('span');
                span.id = 'diary-like-wrapper';
                span.dataset.diaryId = diary.diary_id;
                span.innerHTML = `
                <button class="btn" @click="toggleDiaryLike(${diary.diary_id})" :disabled="diaryLike.loading">
                    <i class="mdi mr-1"
                        :class="diaryLike.liked ? 'mdi-thumb-up text-red-500' : 'mdi-thumb-up-outline'"></i>
                    <span x-text="diaryLike.count"></span>
                </button>`;
                btnContainer.appendChild(span);
                if (window.Alpine) {
                  window.Alpine.initTree(span);
                }
              }
            }
            // Reload stats and like status
            this.loadStats();
            this.loadDiaryLike();
          } else {
            if (socialWrapper) {
              socialWrapper.remove();
            }
            // Clear stats to hide them immediately
            this.stats = {};
          }

          return;
        } else {
          Notify.error(
            result.message || this.$store.i18n.t('common.save_fail'),
          );
          return false;
        }
      } catch (e) {
        Notify.error(this.$store.i18n.t('common.net_error'));
        return false;
      }
    },
  };
}
