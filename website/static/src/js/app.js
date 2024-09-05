

var Biuren = {
    html: function (container, html) {
        if ($("#" + container).length > 0) {
            $("#" + container).get(0).innerHTML = html;
        }
    },
    inputErr: function (fid, emid, str, strCss, setTop) {
        var f = $("#" + fid);
        this.showText(emid, str, strCss);
        if (f.length == 0) {
            return false;
        }
        if (f.disabled) {
            return false;
        }
        f.focus();
        if (typeof setTop != "undefined") {
            this.scrollTop();
        }
        return false;
    },
    post: function (options) {
        var options = jQuery.extend({
            type: "post",
            dataType: "json",
            success: function (json) {
                if (Biuren.checkError(json, ajaxOptions.container)) {
                    return;
                }
                if (json.Result == 0) {
                    Biuren.Pager.setAjaxStatus("本次操作成功！", "help");
                } else {
                    Biuren.Comm.showErrorMsg("help", Biuren.getErrorText(json), "error");
                }
            }
        }, options);
        $.ajax(options);
    },
    get: function (options) {
        var options = jQuery.extend({
            cache: false,
            type: "get"
        }, options);
        $.ajax(options);
    },
    showOK: function () {
        $("#help").show().attr("class", "help").html("操作成功").fadeOut(4000);
    },
    setBgColor: function (obj) {
        $(obj).css({
            backgroundColor: "#EFEFEF"
        });
    },
    moveBgColor: function (obj) {
        $(obj).css({
            backgroundColor: ""
        });
    },
    checkError: function (data, container) {
        var result = false;
        if (data == null) {
            return true;
        }
        if (typeof (data.Result) != "undefined") {
            if (data.Result != 0) {
                Biuren.Comm.showErrorMsg("help", Biuren.getErrorText(data), "error");
                result = true;
            }
        }
        if (result) {
            $("#" + container).text("");
        }
        return result;
    }
};

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
        form.find('#time_zone').val(Intl.DateTimeFormat().resolvedOptions().timeZone);
        form.find('#lang').val(navigator.language);
        $.post(AJAX_URL + '/register', form.serialize(), function (data) {
            console.log(data);
        }, 'json').fail(function (res) {
            console.log(res.responseJSON);
        });
    }
}

$(document).ready(function () {
    $('#register_btn').on('click', Biuren.Passport.reg);
    $('#login_btn').on('click', Biuren.Passport.login);
});

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
        if (Biuren.Valid.checkError("diary_editor_form")) {
            return;
        }
        $.post(AJAX_URL + '/info', $("#diary_editor_form").serialize(), function (data) {
            if (data.Status == 0) {
                
            }
        }, 'json');
    }
}

$(function () {
   // Biuren.Valid.init();
  //  Biuren.Pager.init();
    Biuren.Diary.init();

    //header fixed
    $(document).scroll(function () {
        if ($(document).scrollTop() > 1)
            $('#top-nav').html($('#page-nav').html());
        else if ($(document).scrollTop() == 0)
            $('#top-nav').html('');
    });
});