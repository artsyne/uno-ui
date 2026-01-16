export const Member = {
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

  openAddDiary() {
    const uno = window.getUno();
    if (!uno) return;

    const title = this.$store.i18n.t('diary.notebook.add');
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
      onOk: async () => {
        const form = document.querySelector('#daybook_editor_form');
        if (!form) return;
        if (await uno.valid.validateForm(form)) {
          // Trigger HTMX submission
          if (window.htmx) window.htmx.trigger(form, 'submit');
        }
      },
    });

    fetch('/form/diary?action=diary.add')
      .then((res) => res.text())
      .then((html) => {
        const container = document.getElementById('diary-form-container');
        if (container) {
          container.innerHTML = html;
          // Manually populate action field
          const actionInput = container.querySelector('#action');
          if (actionInput) actionInput.value = 'diary.add';

          if (window.Alpine) {
            window.Alpine.initTree(container);
          }
          if (window.htmx) {
            window.htmx.process(container);
          }
        }
      })
      .catch((err) => {
        const container = document.getElementById('diary-form-container');
        if (container)
          container.innerHTML = `<p class="text-red-500">${this.$store.i18n.t('common.load_fail')}</p>`;
      });
  },
};
