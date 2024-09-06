var Biuren = {}

var AJAX_URL = '';

Biuren.Passport = {
    login: function () {
        var form = $(this).closest('form');
        if (form.valid("checkError")) {
            return;
        }
        $.post(AJAX_URL + '/login', form.serialize(), function (data) {
            window.location = data.AccessUrl;
        }, 'json').fail(function (res) {
            console.log(res.responseJSON);
        });
    },
    reg: function () {
        var form = $(this).closest('form');
        if (form.valid("checkError")) {
            return;
        }
        var timeZone = Intl.DateTimeFormat().resolvedOptions().timeZone;
        var lang = navigator.language;
        
        // 序列化表单数据
        var formData = form.serialize();
        
        // 附加时区和语言信息
        formData += '&time_zone=' + encodeURIComponent(timeZone);
        formData += '&lang=' + encodeURIComponent(lang);
        
        $.post(AJAX_URL + '/register', formData, function (data) {
            console.log(data);
        }, 'json').fail(function (res) {
            console.log(res.responseJSON);
        });
    }
}

Biuren.Diary = {
    init: function () {
        $('#daybook_add_modal').on('shown.bs.modal', function (e) {
            $.get(AJAX_URL + '/form/daybook', function (data) {
                $('#daybook_add_modal .form_html:first').html(data);
            }).done(function () {
                Biuren.Valid.init();
                $('.d-theme>div').mouseover(
                    function () {
                        $(this).addClass('dtfhover');
                    }
                ).mouseout(function () {
                    $(this).removeClass('dtfhover');
                }).click(function () {
                    $('.d-theme>div').removeClass('dtfselected');
                    $(this).addClass('dtfselected');
                    $('#theme_id').val($(this).attr("theme"));
                });
            });
        });

        $('#daybook_add_modal').on('hidden.bs.modal', function (e) {
            $('#daybook_add_modal .form_html:first').html(" <div class=\"modal-body\">正在加载中......</div>");
        });
    },
    doDaybook: function () {
        if (Biuren.Valid.checkError("daybook_editor_form")) {
            return;
        }
        $.post(AJAX_URL + '/daybook', $("#daybook_editor_form").serialize(), function (data) {
            if (data.Status == 0) {
                var html = "<div class=\"col-sm-6 col-md-4 col-lg-3\">" +
                    "<div class=\"card card-diary-theme dt" + data.Data.theme_id + "\">" +
                    "<div class=\"card-body\"><h5>";
                if (data.Data.view_setting == 2) 
                    html+= "<i class='mdi mdi-lock-outline'></i>"
                html+= "<a href='/d/" + data.Data.diary_id + "'>" + $.HtmlEncode(data.Data.title) + "</a>" 
                html+= "</h5><p>0篇</p></div></div></div>";
                $('#daybook_list').prepend(html);
                $('#daybook_add_modal').modal('hide');
            }
        }, 'json');
    },
    postInfo:function(){
        var form = $(this).closest('form');
        if (form.valid("checkError")) {
            return;
        }
        $.post(AJAX_URL + '/info', form.serialize(), function (data) {
            if (data.Status == 0) {
                
            }
        }, 'json');
    }
}

$(function () {
    $('#register_btn').on('click', Biuren.Passport.reg);
    $('#login_btn').on('click', Biuren.Passport.login);
    $('#post_btn').on('click', Biuren.Diary.postInfo);

    Biuren.Diary.init();

    //header fixed
    $(document).scroll(function () {
        if ($(document).scrollTop() > 1)
            $('#top-nav').html($('#page-nav').html());
        else if ($(document).scrollTop() == 0)
            $('#top-nav').html('');
    });
});