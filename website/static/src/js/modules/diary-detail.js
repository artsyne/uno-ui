import { SocialMixin } from './social.js';
import { EntryMixin } from './entry.js';
import { AuthorList } from './authors.js';

export default function DiaryDetail() {
  return {
    ...SocialMixin(),
    ...EntryMixin,
    ...AuthorList(),

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

      // Listen for diary updates globally to handle real-time changes
      // This avoids race conditions where the modal closes (removing the listener)
      // before the update event is processed.
      this._bindDiaryUpdateListener();
    },

    _bindDiaryUpdateListener() {
      const updateHandler = (evt) => {
        // Handle both simple key (if sent as closeDiaryModal payload) and data payload
        // When using HX-Trigger: {"diaryUpdated": data}, HTMX sets evt.detail to data directly
        const diary = evt.detail.diaryUpdated || evt.detail.value || evt.detail;
        if (diary && typeof diary === 'object' && diary.diary_id) {
          this._updateDiaryUI(diary);
        }
      };
      document.body.addEventListener('diaryUpdated', updateHandler);
    },

    openDiarySettings(diaryId) {
      const uno = window.getUno();
      if (!uno) return;
      const t = (key) => this.$store.i18n.t(key);

      let closeHandler;

      // Generate boxId upfront to ensure availability in closures and correct type
      const boxId =
        'box-diary-settings-' + Math.random().toString(36).slice(2, 10);
      uno.box.dialog({
        boxId: boxId,
        title: t('diary.settings'),
        hxUrl: `/form/diary?action=diary.edit&diary_id=${diaryId}`,
        widthClass: '!w-[880px]',
        heightClass: '!h-[570px]',
        showOkButton: true,
        okText: t('common.save'),
        containerClass: '!border',
        okButtonClass: 'btn',
        buttonAlign: 'right',
        onContentLoad: (el) => {
          if (!el) el = document.getElementById(boxId);
          if (window.Alpine) window.Alpine.initTree(el);
          if (window.htmx) window.htmx.process(el);

          // Attach boxId to the form for reliable event-driven closing
          const form = el.querySelector('form');
          if (form) {
            form.dataset.boxId = boxId;
          }
        },
        onOk: async () => {
          const form = document.getElementById('diary_editor_form');
          if (form && (await uno.valid.validateForm(form))) {
            form.dispatchEvent(new CustomEvent('validated', { bubbles: true }));
            return false;
          }
          return false;
        },
        onClose: () => {
          if (closeHandler)
            document.body.removeEventListener(
              'closeDiaryModal',
              closeHandler,
              true,
            );
        },
      });

      closeHandler = (evt) => {
        const eventBoxId = evt?.target?.dataset?.boxId;

        // If the event carries a specific boxId, only close if it matches this box
        if (eventBoxId && eventBoxId !== boxId) {
          return;
        }

        uno.box.close(boxId);
      };

      document.body.addEventListener('closeDiaryModal', closeHandler, true);
    },

    _updateDiaryUI(diary) {
      const titleEl = document.getElementById('diary_title');
      const remarkEl = document.getElementById('diary_remark');

      if (titleEl) {
        let iconHtml = '';
        if (diary.view_setting === 2) {
          iconHtml = '<i class="mdi mdi-lock-outline mx-3"></i>';
        }
        titleEl.innerHTML = diary.title + iconHtml;
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

      // 检测社交功能开关，处理边界情况
      if (this.socialSetting !== diary.social_setting) {
        this.socialSetting = diary.social_setting;

        // Dispatch global event for entryStats components
        window.dispatchEvent(
          new CustomEvent('social-setting-changed', {
            detail: { setting: diary.social_setting },
          }),
        );

        if (diary.social_setting === 0) {
          // 社交功能关闭：清理批量加载状态
          const batchStatsManager = window.batchStatsManager;
          if (batchStatsManager) {
            batchStatsManager.clear();
          }
        } else if (diary.social_setting === 1) {
          // 社交功能开启：加载日记本点赞状态
          this.loadDiaryLike();
        }
      }

      this.socialSetting = diary.social_setting;

      // Update Social Button
      // The wrapper is now always in DOM, controlled by x-show="socialSetting === 1"
      // So we just need to ensure data is loaded if setting is 1
      if (diary.social_setting === 1) {
        this.loadDiaryLike();
      }
    },

    // 初始化社交功能
    initSocial() {
      this.socialSetting = parseInt(
        document.querySelector('main')?.dataset.socialSetting || '2',
      );

      if (this.socialSetting === 1) {
        this.loadDiaryLike();
      }

      // 监听日记更新事件
      document.body.addEventListener('diaryUpdated', (evt) => {
        const diary = evt.detail.diaryUpdated || evt.detail.value || evt.detail;
        if (diary && typeof diary === 'object' && diary.diary_id) {
          this._updateDiaryUI(diary);
        }
      });
    },
  };
}
