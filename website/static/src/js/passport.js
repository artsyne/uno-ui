/**
 * jQuery Passport Plugin
 * 
 * 这个插件提供了登录和注册功能，支持表单验证和AJAX提交。
 * 它还允许配置是否启用回车键提交表单。
 * 
 * @author nobody
 * @version 1.0.4
 * @license MIT
 * 
 * 使用示例：
 * 
 * $('#passport').passport({
 *   ajaxUrl: '/api',
 *   eventBindings: [
 *     { selector: '#login_btn', event: 'click', allowEnterKey: true, handler: 'login' },
 *     { selector: '#register_btn', event: 'click', allowEnterKey: false, handler: 'reg' },
 *     // 可以添加更多自定义绑定
 *   ]
 * });
 */

(($) => {
  $.fn.passport = function (options) {
    // 默认设置
    const settings = $.extend(
      {
        ajaxUrl: '', // AJAX请求的基础URL
        debug: false, // 新增的debug选项，默认为false
        eventBindings: [ // 默认的事件绑定配置
          { selector: '#login_btn', event: 'click', allowEnterKey: true, handler: 'login' },
          { selector: '#register_btn', event: 'click', allowEnterKey: true, handler: 'reg' },
        ]
      },
      options
    );

    // 添加一个调试日志函数
    const debug = (message, ...args) => {
      if (settings.debug) {
        console.log(`[Passport Debug] ${message}`, ...args);
      }
    };

    const methods = {
      /**
       * 处理登录逻辑
       * @param {Event} event - 触发事件的对象
       */
      login(event) {
        debug('Login method called', event);
        event.preventDefault();
        const form = $(this).closest('form');
        if (form.valid('checkError')) {
          return;
        }

        const formData = form.serialize();
        $.utils.disableFormInputs(form, true);

        $.post(
          `${settings.ajaxUrl}/login`,
          formData,
          (data) => {
            window.location = data.AccessUrl;
          },
          'json'
        ).fail((res) => {
          console.log(res.responseJSON);
          $.utils.disableFormInputs(form, false);
        });
      },

      /**
       * 处理注册逻辑
       * @param {Event} event - 触发事件的对象
       */
      reg(event) {
        debug('Register method called', event);
        event.preventDefault();
        const form = $(this).closest('form');
        if (form.valid('checkError')) {
          return;
        }

        let formData = form.serialize();
        const timeZone = Intl.DateTimeFormat().resolvedOptions().timeZone;
        const lang = navigator.language;

        formData += `&time_zone=${encodeURIComponent(timeZone)}`;
        formData += `&lang=${encodeURIComponent(lang)}`;

        $.utils.disableFormInputs(form, true);

        $.post(
          `${settings.ajaxUrl}/register`,
          formData,
          (data) => {
            console.log(data);
          },
          'json'
        ).fail((res) => {
          console.log(res.responseJSON);
          $.utils.disableFormInputs(form, false);
        });
      },

      /**
       * 处理回车键事件
       * @param {Event} event - 键盘事件对象
       * @param {Object} binding - 事件绑定配置对象
       */
      handleEnterKey(event, binding) {
        debug('Enter key pressed', event, binding);
        if (event.which === 13) { // 13 是回车键的键码
          event.preventDefault();
          $(binding.selector).click();
        }
      },
    };

    /**
     * 绑定所有配置的事件
     */
    const bindEvents = () => {
      debug('Binding events', settings.eventBindings);
      settings.eventBindings.forEach((binding) => {
        // 获取处理函数，可以是字符串（方法名）或函数
        const handler = typeof binding.handler === 'string' ? methods[binding.handler] : binding.handler;
        
        debug(`Binding ${binding.event} event to ${binding.selector}`);
        $(document).on(binding.event, binding.selector, handler);
        
        // 如果允许回车键提交，为相关表单的所有输入框绑定keypress事件
        if (binding.allowEnterKey) {
          debug(`Enabling Enter key for ${binding.selector}`);
          const $form = $(binding.selector).closest('form');
          $form.on('keypress', 'input', function(event) {
            methods.handleEnterKey.call(this, event, binding);
          });
        }
      });
    };

    debug('Passport plugin initialized with settings', settings);
    return this.each(bindEvents);
  };
})(jQuery);