export const EntryMixin = {
  openEntryModal(url, updateHistory = true) {
    const uno = window.getUno();
    if (!uno) return;

    const originalUrl = window.location.href;
    const dialogUrl = (() => {
      const urlObj = new URL(url, window.location.origin);
      urlObj.searchParams.set('isajax', 'true');
      return urlObj.pathname + urlObj.search + urlObj.hash;
    })();

    if (updateHistory) {
      history.pushState(null, '', url);
    }

    uno.box.dialog({
      hxUrl: dialogUrl,
      containerClass: '!border',
      widthClass: 'w-full lg:w-[90%]',
      heightClass: 'h-[90vh]',
      autoFocus: false,
      onContentLoad: (el) => {
        htmx.process(el);
        window.Alpine.initTree(el);
      },
      onClose: () => {
        if (updateHistory) {
          history.pushState(null, '', originalUrl);
        }
      },
    });
  },

  openCreateEntry(diaryId) {
    this._openEntryForm('entry.create', diaryId);
  },

  openEditEntry(diaryId, entryId, source = null) {
    this._openEntryForm('entry.edit', diaryId, entryId, source);
  },

  _openEntryForm(action, diaryId, entryId = null, source = null) {
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
      ...(source && { source: source }),
    });

    const url = `/form/entry?${params.toString()}`;

    let closeHandler;
    let autoSaveStopper;

    // Generate boxId upfront to ensure availability in closures and correct type
    const boxId = 'box-entry-' + Math.random().toString(36).slice(2, 10);
    const boxInstance = uno.box.dialog({
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
          autoSaveStopper = this._setupAutoSaveDraft(form);
        }
      },
      onOk: async () => {
        const form = document.getElementById('entry_editor_form');
        if (form) {
          if (autoSaveStopper) autoSaveStopper();

          const action = form.querySelector('#action')?.value;
          // Check if it's an edit action (and not from drafts page where we might want different behavior)
          // Actually, standard edit should update DOM.
          // We check if target card exists to decide if we do manual update.
          const entryId = form.querySelector('#entry_id')?.value;
          const card = entryId ? document.getElementById(`entry-card-${entryId}`) : null;

          if (action === 'entry.edit' && card) {
            // Manual Fetch Update to avoid card reload
            if (!(await uno.valid.validateForm(form))) {
              return false;
            }

            try {
              const formData = new FormData(form);
              const searchParams = new URLSearchParams(formData);
              // Ensure we don't trigger HTMX by not sending HX-Request
              // Backend UpdateEntry returns JSON if no HX-Request header
              
              const response = await fetch(`/entry/${entryId}`, {
                method: 'PUT',
                headers: {
                  'Content-Type': 'application/x-www-form-urlencoded',
                },
                body: searchParams,
              });

              const result = await response.json();
              
              if (result.success && result.data) {
                const entry = result.data;
                
                // Update Card Title
                const titleEl = card.querySelector('.entry-title');
                if (titleEl) titleEl.textContent = entry.title;
                
                // Update Card Summary
                // Template: <p class="text-relaxed cursor-pointer"><a ...>{{.Entry.Summary}}</a></p>
                const summaryEl = card.querySelector('.text-relaxed a') || card.querySelector('.text-relaxed');
                if (summaryEl) summaryEl.textContent = entry.summary;
                
                // Update Time? Maybe not needed for simple edit.
                
                // Dispatch event for other components
                window.dispatchEvent(new CustomEvent('entryUpdated', {
                  detail: { entryUpdated: entry }
                }));
                
                uno.notify.success(this.$store.i18n.t('common.save_success'));
                return true; // Close dialog
              } else {
                 // Try to handle error
                 const msg = result.message || 'Error';
                 uno.notify.error(msg);
                 return false;
              }
            } catch (e) {
               console.error(e);
               uno.notify.error('Network Error');
               return false;
            }
          }

          form.requestSubmit();
        }
        return false;
      },
      onClose: () => {
        if (closeHandler)
          document.body.removeEventListener(
            'closeEntryModal',
            closeHandler,
            true,
          );
        if (autoSaveStopper) autoSaveStopper();
      },
    });

    // Event listener for closing the modal from HTMX response
    closeHandler = (event) => {
      // Guard: Only close if the event targets this specific box
      // If event.detail.boxId is missing, we assume it's a legacy global event (risky but backward compat)
      // or we can strictly enforce it.
      // Given the new "Box Window Closure Logic Enhancement" memory, strict check is preferred.
      // However, the backend might not be sending boxId in the trigger yet unless we pass it.
      // For now, let's assume if boxId is present in event, it must match.
      if (event.detail && event.detail.boxId && event.detail.boxId !== boxId) {
        return;
      }

      if (boxInstance && typeof boxInstance.close === 'function') {
        boxInstance.close();
        return;
      }

      const box = uno.box.get(boxId);
      if (box) {
        box.close();
      }
    };
    document.body.addEventListener('closeEntryModal', closeHandler, true);

    // Listen for entry update events to refresh title and summary in box
    if (source === 'box' && entryId) {
      const updateHandler = (evt) => {
        const detail = evt.detail;
        if (detail && detail.entryUpdated) {
          const entry = detail.entryUpdated;
          // Update title in box
          const titleEl = document.querySelector('#entry-detail-container h1');
          if (titleEl && entry.Title) {
            titleEl.textContent = entry.Title;
          }
          // Update summary in box (body content)
          const bodyEl = document.querySelector('#entry-detail-container .prose');
          if (bodyEl && entry.Body) {
            // Check if markdown renderer is available
            const markdown = window.markdown?.render || ((text) => text);
            bodyEl.innerHTML = markdown(entry.Body);
          }
        }
      };
      document.body.addEventListener('entryUpdated', updateHandler);
    }
  },

  confirmDelete(entryId, el) {
    this._confirmEntryAction(
      entryId,
      el,
      false,
      'entry.delete.confirm',
      `/entry/${entryId}`,
    );
  },

  confirmHardDelete(entryId, el) {
    this._confirmEntryAction(
      entryId,
      el,
      true,
      'entry.drafts.confirm_hard_delete',
      `/entry/${entryId}/hard`,
    );
  },

  _confirmEntryAction(entryId, el, isHard, messageKey, url) {
    const uno = window.getUno();
    if (!uno) return;

    const message = this.$store.i18n.t(messageKey);
    uno.box.confirm(message).then((confirmed) => {
      if (confirmed) {
        if (window.htmx) {
          window.htmx.ajax('DELETE', url, {
            target: el.closest('tr'),
            swap: 'outerHTML',
          });
        }
      }
    });
  },

  _setupAutoSaveDraft(form) {
    const actionInput = form.querySelector('#action');
    const mode = actionInput?.value;
    if (mode !== 'entry.create' && mode !== 'entry.edit') {
      return () => {};
    }

    const enabled =
      (form.dataset.autoSaveEnabled || 'true').toLowerCase() === 'true';
    const interval =
      Number.parseInt(form.dataset.autoSaveInterval || '30000', 10) || 30000;
    if (!enabled || interval <= 0) {
      return () => {};
    }

    const csrfToken = document
      .querySelector('meta[name="csrf-token"]')
      ?.getAttribute('content');

    const getSnapshot = () => {
      const title = form.querySelector('#info_title')?.value || '';
      const body = form.querySelector('#info_body')?.value || '';
      return `${title}\n${body}`;
    };

    let lastSnapshot = getSnapshot();
    let isSaving = false;
    let timerId = null;

    const updateStatus = (msg) => {
      const el = form.querySelector('#auto-save-status');
      if (el) el.textContent = msg;
    };

    const tick = async () => {
      if (isSaving) return;

      const snapshot = getSnapshot();
      if (lastSnapshot !== null && snapshot === lastSnapshot) return;

      const body = form.querySelector('#info_body')?.value || '';
      if (body.trim().length === 0) return;

      isSaving = true;
      updateStatus('Saving...');
      try {
        const formData = new FormData(form);
        const resp = await fetch('/entry/auto-save', {
          method: 'POST',
          credentials: 'same-origin',
          headers: csrfToken ? { 'X-CSRF-Token': csrfToken } : {},
          body: new URLSearchParams(formData),
        });

        const json = await resp.json().catch(() => null);
        if (!resp.ok || !json?.success) {
          updateStatus('Auto-save failed');
          return;
        }

        const entryId = json?.data?.entry_id;
        if (entryId && !form.querySelector('#entry_id')?.value) {
          const entryIdInput = form.querySelector('#entry_id');
          if (entryIdInput) entryIdInput.value = String(entryId);
        }

        lastSnapshot = snapshot;
        updateStatus(`Last saved: ${new Date().toLocaleTimeString()}`);
      } catch {
        updateStatus('Auto-save error');
      } finally {
        isSaving = false;
      }
    };

    timerId = window.setInterval(tick, interval);

    const stop = () => {
      if (timerId) window.clearInterval(timerId);
      timerId = null;
    };

    form.addEventListener('submit', stop, true);
    return stop;
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
