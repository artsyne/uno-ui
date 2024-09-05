(function ($) {
    $.fn.valid = function () {
        var validators = {
            validname: {
                regex: /^[a-zA-Z_\u4e00-\u9fa5]+$/,
                msg: "请输入有效字符串，应为：英文字符串，下划线，中文"
            },
            en: {
                regex: /^[a-zA-Z]+$/,
                msg: "仅能输入英文字符"
            },
            domain: {
                regex: /^[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/,
                msg: "域名无效"
            },
            mobile: {
                regex: /^1[3-9]\d{9}$/,
                msg: "移动电话无效"
            },
            int: {
                regex: /^\d+$/,
                msg: "必须输入整数"
            },
            email: {
                regex: /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/,
                msg: "电子邮件地址无效"
            },
            date: {
                regex: /^\d{4}-\d{2}-\d{2}$/,
                msg: "日期不合法"
            },
            url: {
                regex: /^(https?:\/\/)?([\da-z.-]+)\.([a-z.]{2,6})([/\w .-]*)*\/?$/,
                msg: "Url地址无效,如：http://www.domain.com"
            }
        };

        var methods = {
            init: function () {
                this.find("[valid]").each(function (i, obj) {
                    var obj = $(obj);
                    if (obj.attr("loadedvalid") != undefined) {
                        return;
                    }
                    obj.after("<div />").attr("loadedvalid", "true");
                    obj.blur(function () {
                        var obj = $(this);
                        methods.checkElement(obj);
                    });
                    $(JSON.parse(obj.attr("valid"))).each(function (i, json) {
                        if (json.show == true) {
                            obj.after("<span />");
                            obj.keydown(methods.pressKey).keyup(methods.pressKey);
                            obj.nextAll("span:first").html("可输入" + (json.max - $.len(obj.val())) + "个字符");
                        }
                    });
                });
            },
            pressKey: function (event) {
                var obj = $(this);
                var value = obj.val();
                var json = JSON.parse(obj.attr("valid"));
                $(json).each(function (i, json) {
                    if (json.show == true) {
                        var maxlen = json.max;
                        if ($.len(value) > maxlen) {
                            var str = $.getByteOfNum(value, maxlen);
                            obj.val(str);
                        }
                        obj.nextAll("span:first").html("可输入" + (json.max - $.len(obj.val())) + "个字符");
                        return false;
                    }
                });
            },
            checkError: function () {
                var result = [];
                $(this).find("[valid]").each(function (i, obj) {
                    result.push(methods.checkElement($(obj)));
                });
                var hasError = false;
                $(result).each(function (i, value) {
                    if (value == true) {
                        hasError = true;
                        return false;
                    }
                });
                return hasError;
            },
            checkElement: function (obj) {
                var validJson = JSON.parse(obj.attr("valid"));
                var labelError = obj.nextAll("div:first");
                var hasError = false;
                var errorMsg = "";
                $(validJson).each(function (i, json) {
                    var type = json.type || json;
                    var val = obj.val();

                    switch (type) {
                        case "emptyquit":
                            if ($.len(val) == 0) {
                                hasError = false;
                                return false;
                            }
                        case "required":
                            if ($.len(val) == 0) {
                                errorMsg = json.msg || "必须填写";
                            }
                            break;
                        case "limit":
                            var len = $.len(val);
                            if (typeof (json.min) != "undefined" && len < json.min) {
                                errorMsg = json.msg || "长度必须大于等于" + json.min + "个字符";
                            }
                            if (typeof (json.max) != "undefined" && len > json.max) {
                                errorMsg = json.msg || "长度必须小于等于" + json.max + "个字符";
                            }
                            break;
                        case "validname":
                        case "en":
                        case "domain":
                        case "mobile":
                        case "int":
                        case "email":
                        case "date":
                        case "url":
                            if (validators[type] && !validators[type].regex.test(val) && val !== "") {
                                errorMsg = json.msg || validators[type].msg;
                            }
                            break;
                        case "regexp":
                            var reg = new RegExp("^" + json.expression + "$", "ig");
                            if (!reg.test(val)) {
                                errorMsg = json.msg || "输入数据不合法";
                            }
                            break;
                        case "range":
                            if (/^\d+$/.test(val)) {
                                val = parseInt(val);
                                if (val < json.min || val > json.max) {
                                    errorMsg = json.msg || "输入数字应该在" + json.min + "到" + json.max + "之间";
                                }
                            } else {
                                errorMsg = "请输入数字";
                            }
                            break;
                        case "compareto":
                            if ($("#" + json.destid).val() != val || val == "") {
                                errorMsg = json.msg || "请输入正确的内容";
                            }
                            break;
                        case "callback":
                            if (eval(json.fn)(obj.attr("id"), val) == false) {
                                errorMsg = json.msg || "输入内容无效";
                            }
                            break;
                        case "tag":
                            if (val.split(json.mark || " ").length > json.max && val != "") {
                                errorMsg = json.msg || "标签个数不能超过" + json.max + "个";
                            }
                            break;
                    }
                    if (errorMsg != "") {
                        labelError.text(errorMsg).attr("class", "invalid-feedback");
                        obj.removeClass("is-invalid is-valid").addClass("is-invalid");
                        hasError = true;
                        return false;
                    }
                });
                if (hasError == false && errorMsg.length == 0) {
                    labelError.text("");
                    obj.removeClass("is-invalid is-valid").addClass("is-valid");
                }
                return hasError;
            }
        };

        if (methods[arguments[0]]) {
            return methods[arguments[0]].apply(this, Array.prototype.slice.call(arguments, 1));
        } else if (typeof arguments[0] === 'object' || !arguments[0]) {
            return methods.init.apply(this, arguments);
        } else {
            $.error('Method ' + arguments[0] + ' does not exist on jQuery.valid');
        }
    };

    // 文档加载完成后自动初始化表单验证
    $(document).ready(function () {
        $("form").valid();
    });
})(jQuery);