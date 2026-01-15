export const Member = {
  init() {
    if (this._initialized) return;
    this._initialized = true;

    document.addEventListener('diary:added', (e) => {
      this.appendDiaryCard(e.detail);
    });
  },

  appendDiaryCard(diary) {
    // Check if we are not on the first page using path parsing (format: /u/{uid}/{page})
    // or query param as fallback
    const pathParts = window.location.pathname.split('/');
    let isPagegt1 = false;
    // Expected: "", "u", "uid", "page"
    if (pathParts.length >= 4) {
      const pageIdx = parseInt(pathParts[3]);
      if (!isNaN(pageIdx) && pageIdx > 1) {
        isPagegt1 = true;
      }
    } else {
      // Fallback to query param
      const urlParams = new URLSearchParams(window.location.search);
      const page = parseInt(urlParams.get('p') || '1');
      if (page > 1) isPagegt1 = true;
    }

    if (isPagegt1) {
      // SPA redirect to first page
      // Construct target URL (remove page segment or query params)
      // Assuming /u/{uid} is the base for page 1
      let targetUrl = window.location.pathname;
      if (pathParts.length >= 4) {
        // Remove the page part
        targetUrl = '/' + pathParts[1] + '/' + pathParts[2];
      } else {
        // Remove query params
        targetUrl = window.location.pathname;
      }

      if (window.htmx) {
        window.htmx
          .ajax('GET', targetUrl, {
            target: '#diary-notebook-container',
            select: '#diary-notebook-container',
            swap: 'outerHTML',
          })
          .then(() => {
            window.history.pushState(null, '', targetUrl);
          });
        return;
      } else {
        window.location.href = targetUrl;
        return;
      }
    }

    const list = document.getElementById('daybook_list');
    if (!list) return;

    const div = document.createElement('div');
    div.className = `border h-72 text-center dt${diary.theme_id}`;
    div.dataset.diaryId = diary.diary_id;

    const isLocked = diary.view_setting > 1;
    const lockIcon = isLocked ? '<i class="mdi mdi-lock-outline"></i>' : '';

    div.innerHTML = `
        <div class="h-full flex flex-col p-3">
            <div class="flex-grow flex items-center justify-center flex-col">
                <h1 class="text-xl mb-2 break-words w-full">
                    <a href="/d/${diary.diary_id}" class="block w-full">${diary.title}
                        ${lockIcon}
                    </a>
                </h1>
            </div>
        </div>
    `;

    // Prepend to show the new diary at the beginning of the list
    list.prepend(div);

    if (window.Alpine) {
      window.Alpine.initTree(div);
    }
  },

  openAddDiary() {
    const uno = window.getUno();
    if (!uno) return;

    const title = this.$store.i18n.t('diary.create.title');
    const boxId = uno.box.dialog({
      title: title,
      ajaxUrl: '/form/diary?action=diary.add',
      width: 800,
      showOkButton: true,
      okButtonText: this.$store.i18n.t('common.save'),
      onOk: () => this.submitDiaryForm(boxId),
      onContentLoad: (el) => {
        if (window.htmx) htmx.process(el);
      },
    });
  },

  async submitDiaryForm(boxId) {
    const uno = window.getUno();
    if (!uno) return;

    const form = document.querySelector('#daybook_editor_form');
    if (!form) return;

    if (!(await uno.valid.validateForm(form))) {
      return;
    }

    const formData = new FormData(form);
    const searchParams = new URLSearchParams(formData);
    searchParams.delete('action');

    try {
      const response = await fetch('/diary', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/x-www-form-urlencoded',
        },
        body: searchParams,
      });

      const result = await response.json();

      if (result.success) {
        uno.notify.success('operation.success');
        uno.box.close(boxId);

        // Dispatch event to update diary list if listener exists
        document.dispatchEvent(
          new CustomEvent('diary:added', { detail: result.data }),
        );
      } else {
        uno.notify.error(
          result.message || 'common.op_fail',
        );
      }
    } catch (e) {
      uno.notify.error('common.net_error');
    }
  },
};
