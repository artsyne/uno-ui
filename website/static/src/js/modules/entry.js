export const EntryMixin = {
  openEntryModal(url) {
    const uno = window.getUno();
    if (!uno) return;

    const originalUrl = window.location.href;

    history.pushState(null, '', url);

    uno.box.dialog({
      hxUrl: url + '?isajax=true',
      containerClass: '!border',
      widthClass: 'w-full lg:w-[90%]',
      heightClass: 'h-[90vh]',
      onContentLoad: (el) => {
        htmx.process(el);
        window.Alpine.initTree(el);
      },
      onClose: () => {
        history.pushState(null, '', originalUrl);
      },
    });
  },

  openCreateEntry(diaryId) {
    this._openEntryForm('entry.create', diaryId);
  },

  openEditEntry(diaryId, entryId) {
    this._openEntryForm('entry.edit', diaryId, entryId);
  },

  _openEntryForm(action, diaryId, entryId = null) {
    const uno = window.getUno();
    if (!uno) return;

    const isEdit = action === 'entry.edit';
    const title = isEdit
      ? this.$store.i18n.t('entry.edit')
      : this.$store.i18n.t('entry.create');

    const params = new URLSearchParams({
      action: action,
      ...(diaryId && { diary_id: diaryId }),
      ...(entryId && { entry_id: entryId }),
    });

    const url = `/form/entry?${params.toString()}`;

    let closeHandler;

    // Generate boxId upfront to ensure availability in closures and correct type
    const boxId = 'box-entry-' + Math.random().toString(36).slice(2, 10);
    uno.box.dialog({
      boxId: boxId,
      title: title,
      hxUrl: url,
      widthClass: 'w-full lg:w-[70%]',
      heightClass: 'h-[70vh]',
      showOkButton: true,
      okText: this.$store.i18n.t('common.save'),
      containerClass: '!border',
      okButtonClass: 'btn',
      buttonAlign: 'right',
      onContentLoad: (el) => {
        if (window.htmx) {
          window.htmx.process(el);
        }
        if (window.Alpine) {
          window.Alpine.initTree(el);
        }
        // Attach boxId to the form for reliable event-driven closing
        const form = el.querySelector('form');
        if (form) {
          form.dataset.boxId = boxId;
        }
      },
      onOk: () => {
        const form = document.getElementById('entry_editor_form');
        if (form) form.requestSubmit();
        return false;
      },
      onClose: () => {
        if (closeHandler)
          document.body.removeEventListener(
            'closeEntryModal',
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

    document.body.addEventListener('closeEntryModal', closeHandler, true);
  },

  async deleteEntry(diaryId, entryId) {
    const uno = window.getUno();
    if (!uno) return;
    const t = (key) => this.$store.i18n.t(key);

    const confirmed = await uno.box.confirm(t('entry.delete.confirm'));
    if (confirmed) {
      if (window.htmx) {
        htmx
          .ajax('DELETE', `/entry/${entryId}?diary_id=${diaryId}`, {
            target: `#entry-card-${entryId}`,
            swap: 'delete',
          })
          .then(() => {
            uno.notify.success(t('common.delete_success'));
          });
      }
    }
  },
};
