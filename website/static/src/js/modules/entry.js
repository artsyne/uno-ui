export const EntryMixin = {
  _escapeHtml(unsafe) {
    if (!unsafe) return '';
    return unsafe
      .replace(/&/g, '&amp;')
      .replace(/</g, '&lt;')
      .replace(/>/g, '&gt;')
      .replace(/"/g, '&quot;')
      .replace(/'/g, '&#039;');
  },

  _updateEntryStats(entry) {
    const id = entry.entry_id;
    if (!this.stats[id]) {
      this.stats[id] = {};
    }
    Object.assign(this.stats[id], {
      view_num: entry.view_num || 0,
      comment_num: entry.comment_num || 0,
      like_num: entry.like_num || 0,
      liked: this.stats[id].liked || false,
    });
  },

  _renderEntryCard(entry) {
    const t = (key) => this.$store.i18n.t(key);
    const date = new Date(entry.add_utc_time);
    const dateStr = date.toLocaleString();
    const id = entry.entry_id;

    // Owner menu
    let ownerMenu = '';
    if (this.isOwner) {
      ownerMenu = `
          <div class="entry-edit">
              <a href="javascript:void(0)" x-uno-tooltip="{ content: '#entry-menu-${id}', trigger: 'click', placement: 'bottom' }">
                  <i class="mdi mdi-dots-horizontal"></i>
              </a>
          </div>
          <div id="entry-menu-${id}" class="hidden">
              <ul class="text-left">
                  <li class="py-1"><a href="#" class="block hover:text-blue-500"
                          @click.prevent="openEditEntry(${entry.diary_id}, ${id}); $store.uno?.tooltip?.hide()">
                          <i class="mdi mdi-pencil-outline mr-2"></i> ${t('common.edit')}</a>
                  </li>
                  <li class="py-1"><button class="block w-full text-left hover:text-red-500"
                          @click="deleteEntry(${entry.diary_id}, ${id}); $store.uno?.tooltip?.hide()">
                          <i class="mdi mdi-delete-outline mr-2"></i> ${t('common.delete')}</button>
                  </li>
              </ul>
          </div>`;
    }

    // Title
    let titleHtml = '';
    if (entry.title) {
      titleHtml = `<h1 class="entry-title">${this._escapeHtml(entry.title)}</h1>`;
    }

    return `
        <div class="card mb-6 break-inside-avoid" id="entry-card-${id}" data-entry-id="${id}">
            ${ownerMenu}
            <div class="mb-4">
                ${titleHtml}
                <p class="text-relaxed cursor-pointer">
                    <a href="/entry/${id}" @click.prevent="openEntryModal('/entry/${id}')"
                        style="all:unset !important">
                        ${this._escapeHtml(entry.summary)}
                    </a>
                </p>
            </div>
            <div class="flex justify-between">
                <div>
                    <span class="entry-time">
                        ${dateStr}
                    </span>
                </div>
                <!-- Alpine-based Stats Container -->
                <div class="flex items-center" x-show="stats[${id}]" style="display: none;">
                    <button class="btn-social mr-2" title="${t('common.views')}">
                        <i class="mdi mdi-eye mr-1"></i> <span x-text="getStats(${id}).view_num"></span>
                    </button>
                    <button class="btn-social mr-2" title="${t('common.comments')}">
                        <i class="mdi mdi-comment-text-outline mr-1"></i> <span
                            x-text="getStats(${id}).comment_num"></span>
                    </button>
                    <button class="btn-social" @click="toggleLike(${id})">
                        <i class="mdi mr-1"
                            :class="getStats(${id}).liked ? 'mdi-thumb-up' : 'mdi-thumb-up-outline'"></i>
                        <span x-text="getStats(${id}).like_num"></span>
                    </button>
                </div>
            </div>
        </div>
      `;
  },

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

    const boxId = uno.box.dialog({
      title: title,
      content: `<div class="modal-body" id="entry-form-container">${this.$store.i18n.t('common.loading')}</div>`,
      widthClass: 'w-full lg:w-[70%]',
      heightClass: 'h-[70vh]',
      showOkButton: true,
      okText: this.$store.i18n.t('common.save'),
      containerClass: '!border',
      okButtonClass: 'btn',
      buttonAlign: 'right',
      onOk: () => this._handleEntrySubmit(boxId, action),
    });

    const params = new URLSearchParams({
      action: action,
      ...(diaryId && { diary_id: diaryId }),
      ...(entryId && { entry_id: entryId }),
    });

    fetch(`/form/entry?${params.toString()}`)
      .then((res) => res.text())
      .then((html) => {
        const container = document.getElementById('entry-form-container');
        if (container) {
          container.innerHTML = html;
          if (window.Alpine) {
            window.Alpine.initTree(container);
          }
        }
      })
      .catch((err) => {
        const container = document.getElementById('entry-form-container');
        if (container)
          container.innerHTML = `<p class="text-red-500">${this.$store.i18n.t('common.load_fail')}</p>`;
      });
  },

  async _handleEntrySubmit(boxId, action) {
    const uno = window.getUno();
    if (!uno) return false;

    const form = document.querySelector('#entry-form-container form');
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

    let url = '/entry';
    let method = 'POST';
    const entryId = formData.get('entry_id');

    if (action === 'entry.edit' && entryId) {
      url = `/entry/${entryId}`;
      method = 'PUT';
    }

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

        const entry = result.data;
        this._updateEntryStats(entry);
        const entryHtml = this._renderEntryCard(entry);

        // Handle pagination redirect if needed
        const pathParts = window.location.pathname.split('/');
        let isPagegt1 = false;
        if (pathParts.length >= 4) {
          const pageIdx = parseInt(pathParts[3]);
          if (!isNaN(pageIdx) && pageIdx > 1) {
            isPagegt1 = true;
          }
        }

        if (isPagegt1) {
          const targetUrl = `/diary/${entry.diary_id}`;
          if (window.htmx) {
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

        const container = document.getElementById('diary_entry_container');

        // Helper to create element from HTML string
        const createElement = (html) => {
          const div = document.createElement('div');
          div.innerHTML = html.trim();
          return div.firstElementChild;
        };

        const id = entry.entry_id;

        if (action === 'entry.create') {
          if (container) {
            const newEl = createElement(entryHtml);
            const writeEntryCard = document.getElementById('write-entry-card');
            if (writeEntryCard) {
              writeEntryCard.insertAdjacentElement('afterend', newEl);
            } else {
              container.prepend(newEl);
            }
            if (window.Alpine) {
              window.Alpine.initTree(newEl);
            }
          }
        } else if (action === 'entry.edit') {
          const existingCard = document.querySelector(
            `.card[data-entry-id="${id}"]`,
          );
          if (existingCard) {
            const newEl = createElement(entryHtml);
            existingCard.replaceWith(newEl);
            if (window.Alpine) {
              window.Alpine.initTree(newEl);
            }
          }
        }
        return;
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
