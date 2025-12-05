/**
 * Uno - Notify Component
 *
 * @module uno.components.notify
 * @version 0.1.0
 * @license MIT
 */

import { getPositionStyle } from '../../core/overlayContainers'
import { startProgressLoop, stopProgressLoop } from '../../core/progress'
import { log } from '../../core/log'

export default (Alpine) => {
  Alpine.data('notify', ({ position = 'top-right' } = {}) => ({
    items: [],
    position,
    nextId: 1,
    showCloseButton: false,
    closeIndicatorEnabled: false,
    activeRequests: 0,
    progressBarRunning: false,

    get containerClasses() {
      switch (this.position) {
        case 'top-left':
          return 'flex flex-col space-y-2 items-start';
        case 'bottom-right':
          return 'flex flex-col-reverse space-y-2 space-y-reverse items-end';
        case 'bottom-left':
          return 'flex flex-col-reverse space-y-2 space-y-reverse items-start';
        case 'top-right':
        default:
          return 'flex flex-col space-y-2 items-end';
      }
    },

    get containerStyle() {
      return getPositionStyle(this.position)
    },

    itemClasses(type) {
      const base = 'relative border-l-4 pl-3 pr-3 py-4 mb-2 rounded shadow-sm dark:shadow-md transition-colors duration-300 overflow-hidden animate-fadeIn';
      const typeClassMap = {
        success: 'bg-green-100 dark:bg-green-800 border-green-500 dark:border-green-400 text-green-700 dark:text-green-100',
        warning: 'bg-yellow-100 dark:bg-yellow-800 border-yellow-500 dark:border-yellow-400 text-yellow-700 dark:text-yellow-100',
        error: 'bg-red-100 dark:bg-red-800 border-red-500 dark:border-red-400 text-red-700 dark:text-red-100',
        info: 'bg-blue-100 dark:bg-blue-800 border-blue-500 dark:border-blue-400 text-blue-700 dark:text-blue-100'
      };
      const typeClass = typeClassMap[type] || typeClassMap.error;
      return `${base} ${typeClass}`;
    },

    /**
     * Show a notification
     * @param {string} text - Message text
     * @param {Object} options - Options { type, timeout, position, closeIndicator, showCloseButton }
     */
    show(text, { type = 'error', timeout, position, closeIndicator, showCloseButton } = {}) {
      const id = this.nextId++;
      const targetPosition = position || this.position;
      log.info('notify.show', { text, type, timeout, position: targetPosition })
      const item = { id, text, type, targetPosition, showCloseButton: !!showCloseButton, closeIndicator: !!closeIndicator };
      this.items = [item, ...this.items];

      let autoTimeout = timeout;
      if (autoTimeout === undefined) {
        switch (type) {
          case 'success': autoTimeout = 1000; break;
          case 'warning': autoTimeout = 2000; break;
          case 'error': default: autoTimeout = 3000; break;
        }
      }

      if (autoTimeout && autoTimeout > 0) {
        const idx = id;
        setTimeout(() => this.remove(idx), autoTimeout);
      }
    },

    /**
     * Remove a notification by ID
     * @param {number} id 
     */
    remove(id) {
      this.items = this.items.filter(i => i.id !== id);
    },

    success(message, optionsOrTimeout) {
      const msg = message || (this.$store?.i18n?.t('notify.operation.success') || 'Success');
      if (typeof optionsOrTimeout === 'number') {
        this.show(msg, { type: 'success', timeout: optionsOrTimeout });
      } else {
        this.show(msg, { type: 'success', ...(optionsOrTimeout || {}) });
      }
    },

    warning(message, optionsOrTimeout) {
      if (typeof optionsOrTimeout === 'number') {
        this.show(message, { type: 'warning', timeout: optionsOrTimeout });
      } else {
        this.show(message, { type: 'warning', ...(optionsOrTimeout || {}) });
      }
    },

    error(message, optionsOrTimeout) {
      if (typeof optionsOrTimeout === 'number') {
        this.show(message, { type: 'error', timeout: optionsOrTimeout });
      } else {
        this.show(message, { type: 'error', ...(optionsOrTimeout || {}) });
      }
    },

    /**
     * Show or hide loading indicator
     * @param {boolean} show 
     * @param {string|Object} styleOrOptions - 'indicator' | 'progressbar' | { style, position }
     */
    loading(show, styleOrOptions) {
      let loadingStyle = 'indicator';
      let position = null;
      if (typeof styleOrOptions === 'string') {
        loadingStyle = styleOrOptions;
      } else if (typeof styleOrOptions === 'object' && styleOrOptions) {
        loadingStyle = styleOrOptions.style || 'indicator';
        position = styleOrOptions.position;
      }

      if (show) {
        this.activeRequests++;
        if (loadingStyle === 'progressbar') {
          this._ensureProgressBar();
          this._showProgressBar();
        } else {
          this._ensureLoadingIndicator();
          this._showLoadingIndicator(position);
        }
      } else {
        this.activeRequests = Math.max(0, this.activeRequests - 1);
        if (this.activeRequests === 0) {
          this._hideLoadingIndicator();
          this._hideProgressBar();
        }
      }
    },

    _ensureLoadingIndicator() {
      if (!document.getElementById('uno-loading-indicator')) {
        const div = document.createElement('div');
        div.id = 'uno-loading-indicator';
        div.className = 'fixed top-4 right-4 z-[150] bg-blue-100 border-blue-500 text-blue-700 p-2 rounded-full shadow-md flex items-center hidden dark:bg-blue-900 dark:text-blue-100';
        div.innerHTML = '<svg class="animate-spin -ml-1 mr-2 h-5 w-5 text-blue-500 dark:text-blue-400" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"><circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle><path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path></svg><span>Loading</span>';
        document.body.appendChild(div);
      }
    },

    _showLoadingIndicator(position) {
      const el = document.getElementById('uno-loading-indicator');
      if (!el) return;
      el.style.display = 'flex';
    },

    _hideLoadingIndicator() {
      const el = document.getElementById('uno-loading-indicator');
      if (!el) return;
      el.style.display = 'none';
    },

    _progressFill() {
      const bar = document.getElementById('uno-progress-bar');
      return bar ? bar.querySelector('.uno-progress-fill') : null;
    },

    _ensureProgressBar() {
      if (!document.getElementById('uno-progress-bar')) {
        const bar = document.createElement('div');
        bar.id = 'uno-progress-bar';
        bar.className = 'fixed top-0 left-0 right-0 z-[9999] bg-blue-100 dark:bg-blue-900 hidden';
        const fill = document.createElement('div');
        fill.className = 'uno-progress-fill bg-blue-500 dark:bg-blue-400 transition-all duration-1000 ease-in-out';
        fill.style.height = '4px';
        fill.style.width = '0%';
        bar.appendChild(fill);
        document.body.appendChild(bar);
      }
    },

    _showProgressBar() {
      const bar = document.getElementById('uno-progress-bar');
      const fill = this._progressFill();
      if (!bar || !fill) return;
      bar.style.display = 'block';
      this.progressBarRunning = true;
      startProgressLoop(fill, 1000);
    },

    _hideProgressBar() {
      const bar = document.getElementById('uno-progress-bar');
      const fill = this._progressFill();
      if (!bar || !fill) return;
      bar.style.display = 'none';
      this.progressBarRunning = false;
      stopProgressLoop(fill);
    },
  }));
};
