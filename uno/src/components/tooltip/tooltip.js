/**
 * Uno - Tooltip Component
 *
 * @module uno.components.tooltip
 * @version 0.1.0
 * @license MIT
 */

import { calculatePosition } from '../../core/positioning';
import { log } from '../../core/log';

export default (Alpine) => {
  /**
   * Tooltip component data
   * @param {string} initialContent - Tooltip text content
   * @param {Object} options - Configuration options
   */
  Alpine.data('tooltip', (initialContent = '', { position = 'top', trigger = 'hover', timeout = 100, target = null } = {}) => ({
    open: false,
    content: initialContent,
    position: position,
    currentPosition: position,
    trigger: trigger,
    timeout: timeout,
    targetId: target,
    timer: null,
    top: 0,
    left: 0,
    positioned: false,
    
    get classes() {
      return 'fixed z-50 px-3 py-2 text-sm font-medium rounded-lg shadow-sm whitespace-nowrap transition-opacity duration-200 border border-gray-200 bg-white text-gray-900 dark:bg-gray-900 dark:text-white dark:border-gray-700';
    },

    get arrowClasses() {
      const base = 'absolute w-2 h-2 rotate-45 bg-white dark:bg-gray-900';
      let posClass = '';
      switch (this.currentPosition) {
        case 'top': posClass = 'top-full left-1/2 -translate-x-1/2 -mt-1'; break;
        case 'bottom': posClass = 'bottom-full left-1/2 -translate-x-1/2 -mb-1'; break;
        case 'left': posClass = 'left-full top-1/2 -translate-y-1/2 -ml-1'; break;
        case 'right': posClass = 'right-full top-1/2 -translate-y-1/2 -mr-1'; break;
      }
      return `${base} ${posClass}`;
    },

    /**
     * Initialize component
     */
    init() {
      if (this.targetId) {
        const targetEl = document.getElementById(this.targetId);
        if (targetEl) {
          this.content = targetEl.innerHTML;
        }
      }
      
      const update = () => { if(this.open) this.updatePosition() };
      window.addEventListener('scroll', update, true);
      window.addEventListener('resize', update);
    },
    
    /**
     * Show tooltip
     */
    show() {
        clearTimeout(this.timer);
        this.positioned = false;
        this.open = true;
        
        const tryUpdate = () => {
            if (!this.open) return;
            const tooltipEl = this.$refs.tooltipRef;
            // Check if element is rendered and has dimensions
            if (tooltipEl && tooltipEl.getBoundingClientRect().width > 0) {
                this.updatePosition();
                this.positioned = true;
            } else {
                requestAnimationFrame(tryUpdate);
            }
        };
        this.$nextTick(tryUpdate);
    },
    
    /**
     * Hide tooltip with delay
     */
    hide() {
        this.timer = setTimeout(() => {
            this.open = false;
        }, this.timeout);
    },

    onMouseEnter() {
      if (this.trigger === 'hover') this.show();
    },
    
    onMouseLeave() {
      if (this.trigger === 'hover') this.hide();
    },
    
    onClick() {
      if (this.trigger === 'click') {
        if (this.open) {
            this.open = false; // Immediate hide
        } else {
            this.show();
        }
      }
    },
    
    onClickOutside() {
      if (this.trigger === 'click' && this.open) {
        this.open = false;
      }
    },

    /**
     * Calculate and update tooltip position
     */
    updatePosition() {
      // Handle case where Alpine binds $el to the trigger element (e.g. in event handlers)
      let triggerEl = this.$el.querySelector('[data-trigger]');
      if (!triggerEl) {
          if (this.$el.hasAttribute('data-trigger')) {
              triggerEl = this.$el;
          } else {
              triggerEl = this.$el.firstElementChild;
          }
      }

      const tooltipEl = this.$refs.tooltipRef;
      
      if (!triggerEl || !tooltipEl) return;

      try {
          const result = calculatePosition(triggerEl, tooltipEl, this.position);
          log.info('Position calculated:', { pos: this.position, result });
          
          this.top = result.top;
          this.left = result.left;
          this.positioned = true;
          this.currentPosition = result.position;
      } catch (e) {
          log.error('Error in calculatePosition', e);
      }
    }
  }))
}
