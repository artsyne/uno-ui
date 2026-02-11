import { EntryMixin } from './entry.js';

export const User = {
  ...EntryMixin,

  openCreateDiary() {
    const uno = window.getUno();
    if (!uno) return;

    const title = this.$store.i18n.t('diary.create');

    let closeHandler;

    // Generate boxId upfront to ensure availability in closures and correct type
    const boxId = 'box-diary-' + Math.random().toString(36).slice(2, 10);

    // Create the dialog and store its ID locally for the event listener
    uno.box.dialog({
      boxId: boxId,
      title: title,
      hxUrl: '/form/diary?action=diary.create',
      widthClass: '!w-[880px]',
      heightClass: '!h-[570px]',
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
        // Trigger submission which will be intercepted by Alpine for validation
        const form = document.getElementById('diary_editor_form');
        if (form) {
          form.requestSubmit();
        }
        return false; // Prevent default closing, let HTMX event handle it
      },
      onClose: () => {
        // Cleanup event listener when box is closed
        if (closeHandler) {
          document.body.removeEventListener(
            'closeDiaryModal',
            closeHandler,
            true,
          );
        }
      },
    });

    // Specific event handler for this dialog instance
    closeHandler = (evt) => {
      const eventBoxId = evt?.target?.dataset?.boxId;

      // If the event carries a specific boxId, only close if it matches this box
      if (eventBoxId && eventBoxId !== boxId) {
        return;
      }

      uno.box.close(boxId);
    };

    // Register the specific listener
    document.body.addEventListener('closeDiaryModal', closeHandler, true);
  },
};
