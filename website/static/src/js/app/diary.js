(($) => {
  $.fn.diary = function (options) {
    const settings = $.extend(
      {
        ajaxUrl: '',
      },
      options,
    );

    const methods = {
      init() {
        $('#daybook_add_modal').on('shown.bs.modal', () => {
          $.get(`${settings.ajaxUrl}/form/daybook`, (data) => {
            $('#daybook_add_modal .form_html:first').html(data);
          }).done(() => {
            Biuren.Valid.init();
            $('.d-theme>div')
              .mouseover(function () {
                $(this).addClass('dtfhover');
              })
              .mouseout(function () {
                $(this).removeClass('dtfhover');
              })
              .click(function () {
                $('.d-theme>div').removeClass('dtfselected');
                $(this).addClass('dtfselected');
                $('#theme_id').val($(this).attr('theme'));
              });
          });
        });

        $('#daybook_add_modal').on('hidden.bs.modal', () => {
          $('#daybook_add_modal .form_html:first').html(
            '<div class="modal-body">正在加载中......</div>',
          );
        });
      },
      doDaybook() {
        if (Biuren.Valid.checkError('daybook_editor_form')) {
          return;
        }
        $.post(
          `${settings.ajaxUrl}/daybook`,
          $('#daybook_editor_form').serialize(),
          (data) => {
            if (data.Status === 0) {
              const html = `
                            <div class="col-sm-6 col-md-4 col-lg-3">
                                <div class="card card-diary-theme dt${data.Data.theme_id}">
                                    <div class="card-body">
                                        <h5>
                                            ${data.Data.view_setting === 2 ? "<i class='mdi mdi-lock-outline'></i>" : ''}
                                            <a href='/d/${data.Data.diary_id}'>${$.HtmlEncode(data.Data.title)}</a>
                                        </h5>
                                        <p>0篇</p>
                                    </div>
                                </div>
                            </div>`;
              $('#daybook_list').prepend(html);
              $('#daybook_add_modal').modal('hide');
            }
          },
          'json',
        );
      },
      postInfo() {
        const form = $(this).closest('form');
        if (form.valid('checkError')) {
          return;
        }
        $.post(
          `${settings.ajaxUrl}/info`,
          form.serialize(),
          (data) => {
            if (data.Status === 0) {
              // Handle success
            }
          },
          'json',
        );
      },
    };

    return this.each(function () {
      $('#post_btn').on('click', methods.postInfo);

      methods.init();
    });
  };
})(jQuery);
