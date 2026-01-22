export const User = {
  currentBoxId: null,

  init() {
    if (this._initialized) return;
    this._initialized = true;

    document.body.addEventListener('closeDiaryModal', () => {
      if (this.currentBoxId) {
        window.getUnoBox()?.close(this.currentBoxId);
        this.currentBoxId = null;
      }
    });
  },

  openCreateDiary() {
    const uno = window.getUno();
    if (!uno) return;

    const title = this.$store.i18n.t('diary.create');
    this.currentBoxId = uno.box.dialog({
      title: title,
      content: `<div class="modal-body" id="diary-form-container">${this.$store.i18n.t('common.loading')}</div>`,
      widthClass: '!w-[880px]',
      heightClass: '!h-[570px]',
      showOkButton: true,
      okText: this.$store.i18n.t('common.save'),
      containerClass: '!border',
      okButtonClass: 'btn',
      buttonAlign: 'right',
      onOk: () => this._handleCreateDiarySubmit(this.currentBoxId),
    });

    fetch('/form/diary?action=diary.create')
      .then((res) => res.text())
      .then((html) => {
        const container = document.getElementById('diary-form-container');
        if (container) {
          container.innerHTML = html;
          // Manually populate action field
          const actionInput = container.querySelector('#action');
          if (actionInput) actionInput.value = 'diary.create';

          if (window.Alpine) {
            window.Alpine.initTree(container);
          }
        }
      })
      .catch(() => {
        const container = document.getElementById('diary-form-container');
        if (container)
          container.innerHTML = `<p class="text-red-500">${this.$store.i18n.t('common.load_fail')}</p>`;
      });
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

  _renderDiaryCard(diary) {
    const t = (key) => this.$store.i18n.t(key);
    const id = diary.diary_id;
    const themeId = diary.theme_id || 0;
    const title = this._escapeHtml(diary.title);
    let lockIcon = '';
    if (diary.view_setting > 1) {
      lockIcon = '<i class="mdi mdi-lock-outline"></i>';
    }

    let infoHtml = '';
    if (diary.info_num > 0) {
      const date = new Date(diary.last_action_utc_time);
      const dateStr = date.toLocaleString();
      infoHtml = `
      <p class="mb-2">
        <span class="mr-1">${diary.info_num}</span>${t('user.profile.entries_count')}
      </p>
      <p class="mb-2">
        ${dateStr}
      </p>`;
    }

    return `
    <div class="border h-72 text-center dt${themeId}" data-diary-id="${id}" id="diary-card-${id}">
      <div class="h-full flex flex-col p-3">
        <div class="flex-grow flex items-center justify-center flex-col">
          <h1 class="text-xl mb-2 break-words w-full">
            <a href="/diary/${id}" class="block w-full">${title}
              ${lockIcon}
            </a>
          </h1>
          ${infoHtml}
        </div>
      </div>
    </div>`;
  },

  async _handleCreateDiarySubmit(boxId) {
    const uno = window.getUno();
    if (!uno) return false;

    const form = document.querySelector('#diary_editor_form');
    if (!form) return false;

    if (uno.valid) {
      const isValid = await uno.valid.validateForm(form);
      if (!isValid) {
        return false;
      }
    }

    const formData = new FormData(form);
    const searchParams = new URLSearchParams(formData);
    searchParams.delete('action');

    const url = '/diary';
    const method = 'POST';

    // Import helper dynamically or assume it's available via module system if bundled
    // Since we are in ES modules, we should import it at the top.
    // However, to minimize changes, I will use dynamic import or expect it to be passed.
    // Better: Add import at the top of file.
    
    // For now, let's assume we added the import at the top.
    // But since I can't edit the top easily with SearchReplace without context, I'll do it in two steps.
    
    // Step 1: Add loading state
    const { setFormLoading } = await import('./form.js');
    setFormLoading(form, true);

    try {
      const response = await fetch(url, {
        method: method,
        headers: {
          'Content-Type': 'application/x-www-form-urlencoded',
        },
        body: searchParams,
      });

      const result = await response.json();

      if (result.success) {
        uno.notify.success('common.save_success');
        uno.box.close(boxId);
        this.currentBoxId = null;

        const diary = result.data;
        const diaryHtml = this._renderDiaryCard(diary);

        const container = document.getElementById('diary_list');

        // Helper to create element from HTML string
        const createElement = (html) => {
          const div = document.createElement('div');
          div.innerHTML = html.trim();
          return div.firstElementChild;
        };

        if (container) {
          const newEl = createElement(diaryHtml);
          container.prepend(newEl);
          if (window.Alpine) {
            window.Alpine.initTree(newEl);
          }
        }
        return true;
      } else {
        uno.notify.error(result.message || 'common.save_fail');
        return false;
      }
    } catch (e) {
      console.error(e);
      const t = (key) => Alpine.store('i18n').t(key);
      const msg = t('common.net_error');
      uno.notify.error(msg);
      return false;
    } finally {
        setFormLoading(form, false);
    }
  },
};
