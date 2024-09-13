(($) => {
  $.fn.passport = function (options) {
    const settings = $.extend(
      {
        ajaxUrl: '',
      },
      options,
    );

    const methods = {
      login() {
        const form = $(this).closest('form');
        if (form.valid('checkError')) {
          return;
        }

        // 先序列化表单数据
        const formData = form.serialize();

        // 禁用表单中的所有输入元素
        $.utils.disableFormInputs(form, true);

        $.post(
          `${settings.ajaxUrl}/login`,
          formData,
          (data) => {
            window.location = data.AccessUrl;
          },
          'json',
        ).fail((res) => {
          console.log(res.responseJSON);
          // 如果 AJAX 调用失败，重新启用表单中的所有输入元素
          $.utils.disableFormInputs(form, false);
        });
      },
      reg() {
        const form = $(this).closest('form');
        if (form.valid('checkError')) {
          return;
        }

        // 先序列化表单数据
        let formData = form.serialize();

        const timeZone = Intl.DateTimeFormat().resolvedOptions().timeZone;
        const lang = navigator.language;

        formData += `&time_zone=${encodeURIComponent(timeZone)}`;
        formData += `&lang=${encodeURIComponent(lang)}`;

        // 禁用表单中的所有输入元素
        $.utils.disableFormInputs(form, true);

        $.post(
          `${settings.ajaxUrl}/register`,
          formData,
          (data) => {
            console.log(data);
          },
          'json',
        ).fail((res) => {
          console.log(res.responseJSON);
          // 如果 AJAX 调用失败，重新启用表单中的所有输入元素
          $.utils.disableFormInputs(form, false);
        });
      },
    };

    // 绑定事件列表
    const eventBindings = [
      { selector: '#register_btn', event: 'click', handler: methods.reg },
      { selector: '#login_btn', event: 'click', handler: methods.login },
    ];

    const bindEvents = () => {
      eventBindings.forEach(({ selector, event, handler }) => {
        $(document).on(event, selector, handler);
      });
    };

    return this.each(bindEvents);
  };
})(jQuery);