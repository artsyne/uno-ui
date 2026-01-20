import { SocialMixin } from './social.js';
import { EntryMixin } from './entry.js';

export function DiaryDetail() {
  return {
    ...SocialMixin(),
    ...EntryMixin,

    init() {
      // Initialize social features
      this.initSocial();

      // Read isOwner from DOM
      const mainEl = document.querySelector('main');
      this.isOwner = mainEl?.dataset.isOwner === 'true';

      // Expose to global window
      window.openCreateEntry = (diaryId) => this.openCreateEntry(diaryId);
      window.openEditEntry = (diaryId, entryId) =>
        this.openEditEntry(diaryId, entryId);
      window.deleteEntry = (diaryId, entryId) =>
        this.deleteEntry(diaryId, entryId);
      window.openEntryModal = (url) => this.openEntryModal(url);

      window.openDiarySettings = (diaryId) => this.openDiarySettings(diaryId);
      window.toggleDiaryLike = (diaryId) => this.toggleDiaryLike(diaryId);
    },

    openDiarySettings(diaryId) {
      const uno = window.getUno();
      if (!uno) return;
      const t = (key) => this.$store.i18n.t(key);

      const boxId = uno.box.dialog({
        title: t('diary.settings'),
        content: `<div class="modal-body" id="diary-form-container">${t('common.loading')}</div>`,
        widthClass: '!w-[880px]',
        heightClass: '!h-[570px]',
        showOkButton: true,
        okText: t('common.save'),
        containerClass: '!border',
        okButtonClass: 'btn',
        buttonAlign: 'right',
        onOk: () => this._handleDiarySettingsSubmit(boxId),
      });

      fetch(`/form/diary?action=diary.edit&diary_id=${diaryId}`)
        .then((res) => res.text())
        .then((html) => {
          const container = document.getElementById('diary-form-container');
          if (container) {
            container.innerHTML = html;
            const actionInput = container.querySelector('#action');
            if (actionInput) actionInput.value = 'diary.edit';

            if (window.Alpine) {
              window.Alpine.initTree(container);
            }
          }
        })
        .catch((err) => {
          const container = document.getElementById('diary-form-container');
          if (container)
            container.innerHTML = `<p class="text-red-500">${t('common.load_fail')}</p>`;
        });
    },

    async _handleDiarySettingsSubmit(boxId) {
      const uno = window.getUno();
      if (!uno) return false;
      const t = (key) => this.$store.i18n.t(key);

      const form = document.querySelector('#diary-form-container form');
      if (!form) return false;

      if (!(await uno.valid.validateForm(form))) {
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
          uno.notify.success('common.save_success');
          uno.box.close(boxId);

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

          // Update Theme Class and Social Setting
          const mainEl = document.querySelector('main');
          if (mainEl) {
            if (diary.theme_id) {
              mainEl.classList.forEach((cls) => {
                if (cls.startsWith('dt') && !isNaN(cls.substring(2))) {
                  mainEl.classList.remove(cls);
                }
              });
              mainEl.classList.add(`dt${diary.theme_id}`);
            }
            mainEl.dataset.socialSetting = diary.social_setting;
          }

          this.socialSetting = diary.social_setting;

          // Update Social Button
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
                        :class="diaryLike.liked ? 'mdi-thumb-up' : 'mdi-thumb-up-outline'"></i>
                    <span x-text="diaryLike.count"></span>
                </button>`;
                btnContainer.appendChild(span);
                if (window.Alpine) {
                  window.Alpine.initTree(span);
                }
              }
            }
            this.loadStats();
            this.loadDiaryLike();
          } else {
            if (socialWrapper) socialWrapper.remove();
            this.stats = {};
          }

          return;
        } else {
          uno.notify.error(result.message || 'common.save_fail');
          return false;
        }
      } catch (e) {
        uno.notify.error('common.net_error');
        return false;
      }
    },
  };
}
