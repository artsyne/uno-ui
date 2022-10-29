(function ($) {
    $.setCookie = function (name, value, days) {
        var date = new Date();
        date.setTime(date.getTime() + days * 24 * 60 * 60 * 1000);
        document.cookie = escape(name) + "=" + escape(value) + "; expires=" + date.toGMTString();
    };
    $.getCookie = function (name) {
        var arr = document.cookie.match(new RegExp("(^| )" + name + "=([^;]*)(;|$)"));
        if (arr != null) {
            return unescape(arr[2]);
        }
        return null;
    };
    $.delCookie = function (name) {
        var exp = new Date();
        exp.setTime(exp.getTime() - 1);
        var cval = $.getCookie(name);
        if (cval != null) {
            document.cookie = name + "=" + cval + ";expires=" + exp.toGMTString();
        }
    };

    function toIntegersAtLease(n) {
        return n < 10 ? "0" + n : n;
    }
    Date.prototype.toJSON = function (date) {
        return this.getUTCFullYear() + "-" + toIntegersAtLease(this.getUTCMonth()) + "-" + toIntegersAtLease(this.getUTCDate());
    };
    var escapeable = /["\\\x00-\x1f\x7f-\x9f]/g;
    var meta = {
        "\b": "\\b",
        "\t": "\\t",
        "\n": "\\n",
        "\f": "\\f",
        "\r": "\\r",
        '"': '\\"',
        "\\": "\\\\"
    };
    $.quoteString = function (string) {
        if (escapeable.test(string)) {
            return '"' + string.replace(escapeable, function (a) {
                var c = meta[a];
                if (typeof c === "string") {
                    return c;
                }
                c = a.charCodeAt();
                return "\\u00" + Math.floor(c / 16).toString(16) + (c % 16).toString(16);
            }) + '"';
        }
        return '"' + string + '"';
    };
    $.toJSON = function (o, compact) {
        var type = typeof (o);
        if (type == "undefined") {
            return "undefined";
        } else {
            if (type == "number" || type == "boolean") {
                return o + "";
            } else {
                if (o === null) {
                    return "null";
                }
            }
        }
        if (type == "string") {
            return $.quoteString(o);
        }
        if (type == "object" && typeof o.toJSON == "function") {
            return o.toJSON(compact);
        }
        if (type != "function" && typeof (o.length) == "number") {
            var ret = [];
            for (var i = 0; i < o.length; i++) {
                ret.push($.toJSON(o[i], compact));
            }
            if (compact) {
                return "[" + ret.join(",") + "]";
            } else {
                return "[" + ret.join(", ") + "]";
            }
        }
        if (type == "function") {
            throw new TypeError("Unable to convert object of type 'function' to json.");
        }
        var ret = [];
        for (var k in o) {
            var name;
            type = typeof (k);
            if (type == "number") {
                name = '"' + k + '"';
            } else {
                if (type == "string") {
                    name = $.quoteString(k);
                } else {
                    continue;
                }
            }
            var val = $.toJSON(o[k], compact);
            if (typeof (val) != "string") {
                continue;
            }
            if (compact) {
                ret.push(name + ":" + val);
            } else {
                ret.push(name + ": " + val);
            }
        }
        return "{" + ret.join(", ") + "}";
    };
    $.compactJSON = function (o) {
        return $.toJSON(o, true);
    };
    $.evalJSON = function (src) {
        return eval("(" + src + ")");
    };
    $.secureEvalJSON = function (src) {
        var filtered = src;
        filtered = filtered.replace(/\\["\\\/bfnrtu]/g, "@");
        filtered = filtered.replace(/"[^"\\\n\r]*"|true|false|null|-?\d+(?:\.\d*)?(?:[eE][+\-]?\d+)?/g, "]");
        filtered = filtered.replace(/(?:^|:|,)(?:\s*\[)+/g, "");
        if (/^[\],:{}\s]*$/.test(filtered)) {
            return eval("(" + src + ")");
        } else {
            throw new SyntaxError("Error parsing JSON, source is not valid.");
        }
    };
    $.HtmlEncode = function (strHTML) {
        var div = document.createElement("div");
        var text = document.createTextNode(strHTML);
        div.appendChild(text);
        return div.innerHTML;
    };
    $.srcElement = function (event) {
        return $(event.srcElement || event.target);
    };
    $.hasCnChar = function (str) {
        str = $.trim(str);
        var reg = /[^\x00-\xff]/;
        return reg.test(str);
    };
    $.isValidName = function (str) {
        var reg = /^(\w|_|[\u4e00-\u9fa5]){1,40}$/;
        return reg.test(str);
    };
    $.isInt = function (str) {
        var reg = /^\d+$/;
        return reg.test(str);
    };
    $.isEnChar = function (str) {
        var reg = /^[a-zA-Z]+$/;
        return reg.test(str);
    };
    $.isDoubleChar = function (str) {
        var reg = /^[^\x00-\xff]+$/;
        return reg.test(str);
    };
    $.isEmail = function (str) {
        var reg = /^[_\.0-9a-z-]+@([0-9a-z][0-9a-z-]+\.){1,4}[a-z]{2,3}$/i;
        return reg.test(str);
    };
    $.isMobile = function (str) {
        var reg = /^1[358][0-9]{9}$/;
        return reg.test(str);
    };
    $.isUrl = function (str) {
        var reg = /(ftp|http|https):\/\/(\w+:{0,1}\w*@)?(\S+)(:[0-9]+)?(\/|\/([\w#!:.?+=&%@!\-\/]))?/;
        return reg.test(str);
    };
    $.isDate = function (sDate) {
        var iaMonthDays = [31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31];
        var iaDate = new Array(3);
        var year, month, day;
        if (arguments.length != 1) {
            return false;
        }
        iaDate = sDate.toString().split("-");
        if (iaDate.length != 3) {
            return false;
        }
        if (iaDate[1].length > 2 || iaDate[2].length > 2) {
            return false;
        }
        year = parseFloat(iaDate[0]);
        month = parseFloat(iaDate[1]);
        day = parseFloat(iaDate[2]);
        if (year < 1900 || year > 2100) {
            return false;
        }
        if (((year % 4 == 0) && (year % 100 != 0)) || (year % 400 == 0)) {
            iaMonthDays[1] = 29;
        }
        if (month < 1 || month > 12) {
            return false;
        }
        if (day < 1 || day > iaMonthDays[month - 1]) {
            return false;
        }
        return true;
    };
    $.showDate = function (date) {
        var year = date.getFullYear();
        var month = date.getMonth() + 1;
        var day = date.getDate();
        return year + "年" + month + "月" + day + "日";
    };
    $.showWeek = function (date) {
        if (date.getDay() == 0) {
            return ("星期日");
        }
        if (date.getDay() == 1) {
            return ("星期一");
        }
        if (date.getDay() == 2) {
            return ("星期二");
        }
        if (date.getDay() == 3) {
            return ("星期三");
        }
        if (date.getDay() == 4) {
            return ("星期四");
        }
        if (date.getDay() == 5) {
            return ("星期五");
        }
        if (date.getDay() == 6) {
            return ("星期六");
        }
    };
    $.isDomain = function (str) {
        var bCheck = false;
        var reg = /^([0-9a-z\u4E00-\u9FA5]([a-z0-9\u4E00-\u9FA5\-]*[a-z0-9\u4E00-\u9FA5])?\.){1,3}[a-z]{2,6}$/gi;
        if (reg.test(str)) {
            reg = /(--)/;
            if (!reg.test(str)) {
                bCheck = true;
            }
        }
        return bCheck;
    };
    $.showText = function (emid, str, strCss) {
        var e = $("#" + emid);
        if (e.length == 0) {
            return;
        }
        e.show().html(str).attr("class", strCss);
    };
    $.showTextInTag = function (tagid, text, className, needhide) {
        var id = "helpmsg" + tagid;
        if ($("#" + id).length == 0) {
            $("#" + tagid).prepend('<div id="' + id + '"></div>');
        }
        $("#" + id).html(text).attr("class", className).show();
        if (needhide == true) {
            $("#" + id).fadeOut(3000);
        }
    };
    $.clearText = function () {
        for (var i = 0; i < arguments.length; i++) {
            var element = $("#" + arguments[i]);
            if (element) {
                element.hide().html("").attr("class", "");
            }
        }
    };
    $.len = function (str) {
        var str = str.replace(/[^\x00-\xff]/g, "**");
        return str.length;
    };
    $.getByteOfNum = function (val, num) {
        var ch, bytenum = 0;
        var rs = "";
        var pt = /[^\x00-\xff]/;
        for (var i = 0; i < num; i++) {
            ch = val.substr(i, 1);
            if (ch.match(pt)) {
                bytenum += 2;
                if (bytenum > num) {
                    return rs;
                }
            } else {
                bytenum += 1;
            }
            rs += ch;
            if (bytenum == num) {
                return rs;
            }
        }
        return rs;
    };
    $.loadCSS = function (path) {
        var ls = document.getElementsByTagName("link");
        for (i = 0; i < ls.length; i++) {
            if (ls[i].href && ls[i].href.indexOf(path) != -1) {
                return;
            }
        }
        s = document.createElement("link");
        s.rel = "stylesheet";
        s.type = "text/css";
        s.href = path;
        s.disabled = false;
        var head = document.getElementsByTagName("head")[0];
        head.appendChild(s);
    };
    $.rnd = function () {
        var day = new Date();
        var i = 9;
        var time = day.getYear().toString() + day.getMonth().toString() + day.getDate().toString() + day.getHours().toString() + day.getMinutes().toString() + day.getSeconds().toString() + day.getMilliseconds().toString();
        return time;
    };
    $.anchor = function () {
        var hash = window.location.hash;
        hash = hash.substr(1);
        if (hash.indexOf("?") > -1) {
            var ary = hash.split("?", 1);
            hash = ary[0];
        }
        return hash;
    };
    $.QueryString = function (paramName) {
        var sValue = location.search.match(new RegExp("[?&]" + paramName + "=([^&]*)(&?)", "i"));
        return sValue ? sValue[1] : sValue;
    };
    $.setWindowStatus = function (text) {
        window.status = text;
        return true;
    };
})(jQuery);
var Biuren = {
    html: function (container, html) {
        if ($("#" + container).length > 0) {
            $("#" + container).get(0).innerHTML = html;
        }
    },
    getEditorHTMLContents: function (editorName) {
        var editor = CKEDITOR.instances[editorName];
        return (editor.getData());
    },
    insertEditorHtml: function (editorName, html) {
        var editor = CKEDITOR.instances[editorName];
        editor.insertHtml(html);
    },
    getErrorText: function (json, simple) {
        var code = json.Result;
        var msg = json.Desc;
        if (code.toString() == "1020") {
            Biuren.Comm.refreshValidateCode();
        }
        if (simple) {
            return msg;
        } else {
            return "错误码:" + code + "  描述:" + msg;
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
Biuren.Pager = new function () {
    this.updateID = "ajaxcontainer";
    this.defaultAction = "notfound";
    this.init = function () {
        this.blindClick();
        var action = this.getAction();
        if (action != "" && typeof Pager != "undefined" && typeof Pager[action] != "undefined") {
            if (typeof Pager[action].doAction != "undefined") {
                Pager[action].doAction(Pager[action].container, Pager[action].url, action);
            } else {
                Biuren.Pager.doRequest(Pager[action].container, Pager[action].url, action);
            }
        }
    };
    this.initForm = function () {
        var action = this.getAction();
        this.updateElement(Biuren.Pager.updateID, Pager[action].url);
        Biuren.html("biunav", Pager[action].nav);
    };
    this.updateElement = function (destID, url) {
        $("#" + destID).html(requestLoading).load(url);
    };
    this.getAction = function () {
        var action = $.anchor();
        return action == "" ? Biuren.Pager.defaultAction : action;
    };
    this.doClick = function (event) {
        var target = $(this).get(0);
        var url = window.location.href.toLowerCase();
        var clickurl = target.href.toLowerCase().split("#", 2);
        if (clickurl.length != 2) {
            return;
        }
        var href = clickurl[0];
        if (url.indexOf(href) > -1) {
            var action = target.rel;
            if (typeof Pager[action] == "undefined") {
                action = this.defaultAction;
            }
            if (typeof Pager[action].doAction != "undefined") {
                if (typeof (Pager[action].oncreate) != "undefined") {
                    Pager[action].oncreate();
                }
                Pager[action].doAction(Pager[action].container, Pager[action].url, action);
            } else {
                Biuren.Pager.doRequest(Pager[action].container, Pager[action].url, action);
            }
        }
    };
    this.blindClick = function () {
        $(".relaction").bind("click", {}, Biuren.Pager.doClick);
    };
    this.doRequest = function (container, url, action) {
        var options = {
            cache: Pager[action].cache,
            url: url,
            beforeSend: function () {
                if (typeof (Pager[action].oncreate) != "undefined") {
                    Pager[action].oncreate();
                }
                $("#" + container).html(requestLoading);
            },
            success: function (data) {
                if (Biuren.checkError(data, container)) {
                    return;
                }
                Biuren.html(container, data);
                if (typeof (Pager[action].oncomplete) != "undefined") {
                    Pager[action].oncomplete();
                }
            }
        };
        Biuren.get(options);
        if (action != "" && typeof action != "undefined" && $("#biunav")) {
            Biuren.html("biunav", Pager[action].nav);
        }
    };
    this.setAjaxStatus = function (text, className, speed) {
        if ($("#ajaxstatus").length == 0) {
            $('<div id="ajaxstatus" class="help">').appendTo("body");
        }
        $("#ajaxstatus").show().html(text).attr("class", className);
        if (typeof (speed) != "undefined") {
            $("#ajaxstatus").fadeOut(speed);
        }
    };
    this.clearAjaxStatus = function () {
        $("#ajaxstatus").fadeOut(2000);
    };
};
Biuren.Box = {
    options: {},
    ver: "1.2",
    isIE6: false,
    _initBox: function (options) {
        if ($("#biurenbox_overlay").length == 0) {
            var boxHtml = '<div id="biurenbox" style="display: none;"><div id="biurenbox_title"><div id="biurenbox_title_text"></div><div id="biurenbox_line"></div></div><div id="biurenbox_container"><div id="biurenbox_content"></div></div><div id="biuren_footer"></div><div id="biurenbox_title_close" title="关闭"></div></div><div id="biurenbox_overlay" class="biurenbox_overlay" ></div>';
            $(boxHtml).appendTo("body");
        }
        var options = jQuery.extend({
            title: "提示信息",
            content: "暂无内容",
            modal: true,
            position: "center"
        }, options || {});
        if ($.browser.msie && $.browser.version == "6.0") {
            Biuren.Box.isIE6 = true;
        }
        Biuren.Box.options = options;
        $("#biurenbox_title_close").click(Biuren.Box.close);
        $("#biurenbox_title_text").html(options.title);
        if (typeof (options.footerTxt) == "string") {
            $("#biuren_footer").html(options.footerTxt).show();
        } else {
            $("#biuren_footer").hide();
        }
        if (options.modal) {
            $(".biurenbox_overlay").addClass("biurenbox_overlay").css({
                width: $(window).width() + "px",
                height: $(document).height() + "px"
            }).show();
        } else {
            $(".biurenbox_overlay").removeClass("biurenbox_overlay");
        }
        if (typeof (options.content) == "string") {
            Biuren.html("biurenbox_content", options.content);
        }
        if (typeof (options.targetId) == "string") {
            Biuren.html("biurenbox_content", $("#" + options.targetId).html());
        }
        if (typeof (options.ajaxUrl) == "string") {
            $("#biurenbox_content").html(requestLoading);
            Biuren.get({
                url: options.ajaxUrl,
                success: function (data) {
                    Biuren.html("biurenbox_content", data);
                    if ($.isFunction(options.ajaxCallback)) {
                        options.ajaxCallback();
                    }
                }
            });
        }
        if (typeof (options.iframeUrl) == "string") {
            Biuren.html("biurenbox_content", '<iframe frameborder="0" name="biurenbox_content" id="biurenbox_framecontent" src="' + options.iframeUrl + '" width="100%" height="100%"> </iframe>');
        }
    },
    dialog: function (options) {
        if ($.isFunction(options.onBeforeOpen)) {
            options.onBeforeOpen();
        }
        Biuren.Box._initBox(options);
        Biuren.Box._showBox();
        if ($.isFunction(Biuren.Box.options.onOpen)) {
            Biuren.Box.options.onOpen();
        }
    },
    confirm: function (options) {
        if ($.isFunction(options.onBeforeOpen)) {
            options.onBeforeOpen();
        }
        Biuren.Box._initBox(options);
        if ($("#biurenbox_button").length == 0) {
            $("#biurenbox_container").after('<div id="biurenbox_button"><input id="biurenbox_okbtn" type="button" value="确定" >&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<input id="biurenbox_canclebtn" type="button" value="取消" ></div>');
            $("#biurenbox_okbtn").click(function () {
                Biuren.Box.close();
                if ($.isFunction(options.onOk)) {
                    options.onOk();
                }
            });
            $("#biurenbox_canclebtn").click(function () {
                Biuren.Box.close();
                if ($.isFunction(options.onCancle)) {
                    options.onCancle();
                }
            });
        }
        Biuren.Box._showBox();
        if ($.isFunction(Biuren.Box.options.onOpen)) {
            Biuren.Box.options.onOpen();
        }
    },
    message: function (options) {
        if ($.isFunction(options.onBeforeOpen)) {
            options.onBeforeOpen();
        }
        Biuren.Box._initBox(options);
        if ($("#biurenbox_button").length == 0) {
            $("#biurenbox_container").after('<div id="biurenbox_button"><input id="biurenbox_okbtn" type="button" value="确定" ></div>');
            $("#biurenbox_okbtn").click(function () {
                Biuren.Box.close();
                if ($.isFunction(options.onOk)) {
                    options.onOk();
                }
            });
        }
        Biuren.Box._showBox();
        if ($.isFunction(Biuren.Box.options.onOpen)) {
            Biuren.Box.options.onOpen();
        }
    },
    close: function () {
        var options = Biuren.Box.options;
        if ($.isFunction(options.onBeforeClose)) {
            options.onBeforeClose();
        }
        $("#biurenbox").hide();
        $(".biurenbox_overlay").hide();
        Biuren.Box._showElement();
        $("#biurenbox_content").empty();
        $("#biurenbox_button").remove();
        if ($.isFunction(options.onClose)) {
            options.onClose();
        }
    },
    _showBox: function () {
        var options = Biuren.Box.options;
        var width = options.width + 10 * 2 + 8 * 2;
        var height = options.height + 10 * 2 + (1 + 4) + 25 + 25 + 10 * 2 + 8 * 2;
        var csspos = {};
        if (!$.isArray(options.position)) {
            csspos = $.extend(Biuren.Box._getPositonCss(0, options.position, width, height), Biuren.Box._getPositonCss(1, "default", width, height));
        } else {
            csspos = $.extend(Biuren.Box._getPositonCss(0, options.position[0], width, height), Biuren.Box._getPositonCss(1, options.position[1], width, height));
        }
        var boxStyle = "width:" + options.width + "px;";
        $.each(csspos, function (i, n) {
            boxStyle += i + ":" + n + ";";
        });
        $("#biurenbox").attr("style", boxStyle).show();
        $("#biurenbox_content").css({
            height: options.height + "px"
        });
        Biuren.Box._hideElement();
    },
    _getPositonCss: function (index, pos, width, height) {
        var csspos = {};
        var windowWidth, windowHeight;
        if (self.innerHeight) {
            windowWidth = self.innerWidth;
            windowHeight = self.innerHeight;
        } else {
            if (document.documentElement && document.documentElement.clientHeight) {
                windowWidth = document.documentElement.clientWidth;
                windowHeight = document.documentElement.clientHeight;
            } else {
                if (document.body) {
                    windowWidth = document.body.clientWidth;
                    windowHeight = document.body.clientHeight;
                }
            }
        }
        switch (pos) {
            case "center":
                csspos = index == 0 ? {
                    left: Math.floor((windowWidth - width) / 2)
                } : {
                    top: Math.floor((windowHeight - height) / 2)
                };
                break;
            case "top":
                csspos = {
                    top: 0
                };
                break;
            case "bottom":
                csspos = {
                    bottom: 0
                };
                break;
            case "left":
                csspos = {
                    left: 0
                };
                break;
            case "right":
                csspos = {
                    right: 0
                };
                break;
            default:
                if (typeof (pos) == "number") {
                    csspos = index == 0 ? {
                        left: pos
                    } : {
                        top: pos
                    };
                } else {
                    csspos = index == 0 ? {
                        left: Math.floor((windowWidth - width) / 2)
                    } : {
                        top: Math.floor((windowHeight - height) / 2)
                    };
                }
                break;
        }
        $.each(csspos, function (i, n) {
            if (Biuren.Box.isIE6 && (i == "left" || i == "top")) {
                csspos[i] = "expression(documentElement.scroll" + i.substring(0, 1).toUpperCase() + i.substring(1) + "+" + (n < 0 ? 0 : n) + '+"px")';
                csspos[i.toUpperCase()] = (n < 0 ? 0 : n);
            } else {
                csspos[i] = (n < 0 ? 0 : n) + "px";
            }
        });
        return csspos;
    },
    _hideElement: function () {
        if ($.browser.msie) {
            $("select").hide();
        }
        $("object").hide();
        $("embed").hide();
    },
    _showElement: function () {
        if ($.browser.msie) {
            $("select").show();
        }
        $("embed").show();
        $("object").show();
    }
};
Biuren.Grid = {
    create: function (options) {
        var parentTag = $("#" + options.container);
        var aryHtml = new Array();
        aryHtml.push('<table id="' + options.id + '"  class="' + options.className + '">');
        if (typeof (options.caption) != "undefined") {
            aryHtml.push("<caption>" + options.caption + "</caption>");
        }
        if (typeof (options.template.head) != "undefined") {
            aryHtml.push("<thead><tr>");
            var tempHead = options.template.head;
            var temptitle;
            for (var i = 0; i < tempHead.length; i++) {
                if (typeof (tempHead[i].text) != "undefined") {
                    temptitle = tempHead[i].text;
                } else {
                    temptitle = tempHead[i];
                }
                aryHtml.push("<td>" + temptitle + "</td>");
            }
            aryHtml.push("</tr></thead>");
        }
        aryHtml.push("<tbody>");
        var tempRow = options.template.row;
        var ds = options.ds;
        if (ds.length == 0) {
            aryHtml.push("<tr><td>" + options.noData + "</td></tr>");
        }
        for (var i = 0; i < ds.length; i++) {
            if ((i % options.col) == 0) {
                if (i != 0) {
                    aryHtml.push("</tr>");
                }
                aryHtml.push("<tr ");
                if (typeof (tempRow.className) != "undefined") {
                    var temp = tempRow.className;
                    if (typeof (tempRow.switchClassName) != "undefined" && i % 2 != 0) {
                        temp = tempRow.switchClassName;
                    }
                    aryHtml.push('class="' + temp + '" ');
                }
                if (typeof (tempRow.mark) != "undefined") {
                    aryHtml.push(" " + tempRow.mark + " ");
                }
                if (typeof (tempRow.event) != "undefined") {
                    for (var l = 0; l < tempRow.event.length; l++) {
                        aryHtml.push(" " + tempRow.event[l].name + '= "' + tempRow.event[l].func + '" ');
                    }
                }
                aryHtml.push(">");
            }
            var tempCell = options.template.cell;
            var cellMark = "";
            if (typeof (options.template.cellAttr) != "undefined") {
                if (typeof (options.template.cellAttr.event) != "undefined") {
                    for (var l = 0; l < options.template.cellAttr.event.length; l++) {
                        cellMark += " " + options.template.cellAttr.event[l].name + '= "' + options.template.cellAttr.event[l].func + '" ';
                    }
                }
            }
            for (var j = 0; j < tempCell.length; j++) {
                aryHtml.push("<td " + cellMark);
                if (typeof (tempCell[j].className) != "undefined") {
                    aryHtml.push(' class="' + tempCell[j].className + '" ');
                }
                aryHtml.push(">" + this.replaceParam(tempCell[j].text, this.getValueFromAry(ds[i], tempCell[j].value)) + "</td>");
            }
        }
        aryHtml.push("</tr></tbody>");
        if (typeof (options.template.foot) != "undefined" && ds.length != 0) {
            aryHtml.push('<tfoot><tr><td colspan="' + options.template.cell.length * options.col + '"> ');
            var tempFoot = options.template.foot;
            temp = tempFoot.text;
            if (typeof (options.pageLink) != "undefined") {
                temp = temp.replace("{0}", '<div class="biurennav">' + options.pageLink + "</div>");
            }
            aryHtml.push('<div class="pages">' + temp + "</div>");
            aryHtml.push('<input type="hidden" name="' + options.id + 'ViewState" id="' + options.id + 'ViewState" />');
            aryHtml.push("</td></tr></tfoot>");
        }
        aryHtml.push("</table>");
        Biuren.html(options.container, aryHtml.join(""));
    },
    chkAll: function (obj) {
        var ary = $(".biurenchk");
        if (obj.checked == true) {
            $(".biurenchk").each(function (i, obj) {
                obj.checked = true;
            });
        } else {
            $(".biurenchk").each(function (i, obj) {
                obj.checked = false;
            });
        }
    },
    getChkValue: function () {
        var ary = $(".biurenchk");
        var result = [];
        ary.each(function (i, obj) {
            if (obj.checked == true) {
                result.push(obj.value);
            }
        });
        return result.toString();
    },
    truncate: function (val, num, suffix) {
        if (val.replace(/[^\x00-\xff]/g, "**").length <= num) {
            return val;
        }
        var ch, bytenum = 0;
        var rs = "";
        var pt = /[^\x00-\xff]/;
        num = num - suffix.length;
        var result;
        for (var i = 0; i < num; i++) {
            ch = val.substr(i, 1);
            if (ch.match(pt)) {
                bytenum += 2;
                if (bytenum > num) {
                    result = rs;
                    break;
                }
            } else {
                bytenum += 1;
            }
            rs += ch;
            if (bytenum == num) {
                result = rs;
                break;
            }
        }
        return result + suffix;
    },
    refresh: function (id) {
        if ($("#BiurenGridcurrent" + id).length == 0) {
            return;
        }
        var func = $("#BiurenGridcurrent" + id).attr("rel");
        eval(func);
    },
    getValueFromAry: function (ary, value) {
        var temp;
        var result = new Array(value.length);
        var pos;
        for (var i = 0; i < value.length; i++) {
            pos = typeof (value[i].pos) == "undefined" ? value[i] : value[i].pos;
            temp = ary[pos];
            if (temp == null || typeof (temp) == "undefined" || temp.length == 0) {
                if (typeof (value[i].def) != "undefined") {
                    temp = value[i].def;
                } else {
                    temp = "";
                }
            }
            temp = temp.toString();
            if (typeof (value[i].parse) != "undefined") {
                if (typeof (value[i].parse[temp]) != "undefined") {
                    temp = value[i].parse[temp];
                } else {
                    temp = "";
                }
            }
            if (typeof (value[i].func) != "undefined") {
                temp = eval(value[i].func)(temp);
            }
            if (typeof (value[i].maxlen) != "undefined") {
                temp = Biuren.Grid.truncate(temp, value[i].maxlen, "...");
            }
            if (value[i].isHtml != true) {
                temp = $.HtmlEncode(temp);
            }
            result[i] = temp;
        }
        return result;
    },
    replaceParam: function (text, value) {
        var result = text;
        var reg;
        for (var i = 0; i < value.length; i++) {
            reg = new RegExp("{(" + i + ")}", "g");
            result = result.replace(reg, value[i]);
        }
        return result;
    },
    _createTable: function (json, options) {
        $("#" + options.container).html("");
        if (typeof (json.Table) != "undefined") {
            jQuery.extend(options, {
                ds: json.Table.rows
            });
        } else {
            jQuery.extend(options, {
                ds: json.rows
            });
        }
        this.create(options);
    },
    _getInterval: function (pageIndex, pages) {
        var ne_half = Math.ceil(8 / 2);
        var upper_limit = pages - 8;
        var start = pageIndex > ne_half ? Math.max(Math.min(pageIndex - ne_half, upper_limit), 0) : 0;
        var end = pageIndex > ne_half ? Math.min(pageIndex + ne_half, pages) : Math.min(8, pages);
        return [start, end];
    },
    _appendItem: function (id, pageIndex, currentPage, pages, appendopts) {
        pageIndex = pageIndex < 0 ? 0 : (pageIndex < pages ? pageIndex : pages - 1);
        appendopts = jQuery.extend({
            text: pageIndex + 1,
            classes: ""
        }, appendopts || {});
        if (pageIndex == currentPage) {
            var lnk = '<span id="BiurenGridcurrent' + id + '" class="current" rel="Biuren.Grid.show(\'' + escape($.toJSON({
                id: id,
                pageIndex: pageIndex + 1
            })) + "')\">" + (appendopts.text) + "</span>";
        } else {
            var lnk = '<a href="javascript:void(0)" onclick="Biuren.Grid.show(\'' + escape($.toJSON({
                id: id,
                pageIndex: pageIndex + 1
            })) + "')\"> " + (appendopts.text) + "</a>";
        }
        return lnk;
    },
    _pagination: function (id, listOptions) {
        var pageSize = listOptions.row * listOptions.col;
        var pageIndex = listOptions.pageIndex - 1;
        var pages = Math.floor((listOptions.total - 1) / pageSize) + 1;
        var aryNav = new Array();
        var interval = Biuren.Grid._getInterval(pageIndex, pages);
        if (pageIndex > 0) {
            aryNav.push(Biuren.Grid._appendItem(id, pageIndex - 1, pageIndex, pages, {
                text: "上页",
                classes: "prev"
            }));
        }
        if (interval[0] > 0) {
            var end = Math.min(2, interval[0]);
            for (var i = 0; i < end; i++) {
                aryNav.push(Biuren.Grid._appendItem(id, i, pageIndex, pages));
            }
            if (interval[0] > 2) {
                aryNav.push("<span>...</span>");
            }
        }
        for (var i = interval[0]; i < interval[1]; i++) {
            aryNav.push(Biuren.Grid._appendItem(id, i, pageIndex, pages));
        }
        if (interval[1] < pages) {
            if (pages - 2 > interval[1]) {
                aryNav.push("<span>...</span>");
            }
            var begin = Math.max(pages - 2 + 1, interval[1]);
            for (var i = begin; i < pages; i++) {
                aryNav.push(Biuren.Grid._appendItem(id, i, pageIndex, pages));
            }
        }
        if (pageIndex < pages - 1) {
            aryNav.push(Biuren.Grid._appendItem(id, pageIndex + 1, pageIndex, pages, {
                text: "下页",
                classes: "next"
            }));
        }
        return aryNav.join("");
    },
    show: function (ajaxOptions, templateOptions) {
        if (arguments.length == 1) {
            var obj = $.evalJSON(unescape(arguments[0]));
            var json = $("#" + obj.id + "ViewState").val();
            json = $.evalJSON(unescape(json));
            ajaxOptions = jQuery.extend(json.ajax, {
                pageIndex: obj.pageIndex
            });
            templateOptions = json.template;
        }
        jQuery.extend(templateOptions, {
            container: ajaxOptions.container
        });
        $("#" + ajaxOptions.container).html(requestLoading);
        var options = {
            type: ajaxOptions.type || "get",
            cache: ajaxOptions.cache || false,
            dataType: "json",
            url: ajaxOptions.url,
            container: ajaxOptions.container,
            data: ajaxOptions.queryParams + "&pageSize=" + ajaxOptions.row * ajaxOptions.col + "&pageIndex=" + ajaxOptions.pageIndex,
            success: function (json) {
                if (Biuren.checkError(json, ajaxOptions.container)) {
                    return;
                }
                if (typeof (templateOptions.template.foot) != "undefined") {
                    if (ajaxOptions.pageIndex == 1) {
                        jQuery.extend(ajaxOptions, {
                            total: json.Table1.rows[0][0]
                        });
                    }
                    var pagelink = Biuren.Grid._pagination(templateOptions.id, ajaxOptions);
                    jQuery.extend(templateOptions, {
                        pageLink: pagelink
                    });
                }
                jQuery.extend(templateOptions, {
                    row: ajaxOptions.row,
                    col: ajaxOptions.col
                });
                Biuren.Grid._createTable(json, templateOptions);
                if (typeof (templateOptions.template.foot) != "undefined") {
                    templateOptions = jQuery.extend(templateOptions, {
                        ds: []
                    });
                    var param = {
                        ajax: ajaxOptions,
                        template: templateOptions
                    };
                    $("#" + templateOptions.id + "ViewState").val(escape($.toJSON(param)));
                }
                if (typeof (templateOptions.callback) != "undefined") {
                    eval(templateOptions.callback)(json, ajaxOptions);
                }
            }
        };
        Biuren.get(options);
    }
};
Biuren.Toggle = {
    init: function () {
        $(".biurentoggle").each(function (i, obj) {
            obj = $(obj);
            obj.bind("click", {}, Biuren.Toggle.doClick);
            var target = obj.attr("toggle");
            if (obj.attr("cookie") == "true") {
                if ($.getCookie(target) == "1") {
                    $("#" + target).hide();
                    obj.toggleClass("biurentoggleclick");
                    obj.attr("title", "展开");
                }
            } else {
                if (obj.hasClass("biurentoggleclick")) {
                    $("#" + target).hide();
                }
            }
        });
    },
    doClick: function () {
        var link = $(this);
        var target = link.attr("toggle");
        var cookie;
        $("#" + target).slideToggle("slow");
        link.toggleClass("biurentoggleclick");
        if (link.attr("title") == "收起") {
            link.attr("title", "展开");
            cookie = "1";
        } else {
            link.attr("title", "收起");
            cookie = "0";
        }
        if (link.attr("cookie") == "true") {
            $.setCookie(target, cookie, 36);
        }
    }
};
Biuren.Tab = {
    init: function () {
        var anchor = $.anchor();
        $("a[tabgroup][loadedtab!=true]").each(function (i, obj) {
            var obj = $(obj);
            obj.click(Biuren.Tab.doClick);
            if (obj.attr("rel") == anchor && anchor.length > 0) {
                $('a[tabgroup="' + obj.attr("tabgroup") + '"]').removeClass("current");
                obj.addClass("current");
            }
            obj.attr("loadedtab", "true");
        });
    },
    doClick: function () {
        var target = $(this);
        var links = $('a[tabgroup="' + target.attr("tabgroup") + '"]').removeClass("current");
        target.addClass("current");
        $('a[tabgroup="' + target.attr("tabgroup") + '"][tabtarget]').each(function (i, obj) {
            var obj = $(this);
            var targetdiv = $("#" + obj.attr("tabtarget"));
            if (obj.attr("tabtarget") == target.attr("tabtarget")) {
                targetdiv.show();
            } else {
                targetdiv.hide();
            }
        });
    }
};
Biuren.Valid = {
    imagePath: "",
    init: function () {
        $("form").find("[valid]").each(function (i, obj) {
            var obj = $(obj);
            if (obj.attr("loadedvalid") != undefined) {
                return;
            }
            obj.after("<label />").attr("loadedvalid", "true");
            obj.blur(function () {
                var obj = $(this);
                Biuren.Valid.checkElement(obj);
            });
            $($.evalJSON(obj.attr("valid"))).each(function (i, json) {
                if (json.show == true) {
                    obj.after("<span />");
                    obj.keydown(Biuren.Valid.pressKey).keyup(Biuren.Valid.pressKey);
                    obj.nextAll("span:first").html("<br/>可输入" + (json.max - $.len(obj.val())) + "个字符");
                }
            });
        });
    },
    pressKey: function (event) {
        var obj = $(this);
        var value = obj.val();
        var json = $.evalJSON(obj.attr("valid"));
        $($.evalJSON(obj.attr("valid"))).each(function (i, json) {
            if (json.show == true) {
                var maxlen = json.max;
                if ($.len(value) > maxlen) {
                    var str = $.getByteOfNum(value, maxlen);
                    obj.val(str);
                }
                obj.nextAll("span:first").html("<br/>可输入" + (json.max - $.len(value)) + "个字符");
                return false;
            }
        });
    },
    checkError: function (id) {
        var result = [];
        $("#" + id).find("[valid]").each(function (i, obj) {
            result.push(Biuren.Valid.checkElement($(obj)));
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
        var validJson = $.evalJSON(obj.attr("valid"));
        var labelError = obj.nextAll("label:first");
        var hasError = false;
        var errorMsg = "";
        var errorclass = obj.attr("validerrorclass");
        var errclass = typeof (errclass) == "undefined" ? "" : " " + errorclass;
        $(validJson).each(function (i, json) {
            var type = json;
            if (typeof (json) == "object") {
                type = json.type;
            }
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
                    if (!$.isValidName(val) && val != "") {
                        errorMsg = json.msg || "请输入有效字符串，应为：英文字符串，下划线，中文";
                    }
                    break;
                case "en":
                    if (!$.isEnChar(val) && val != "") {
                        errorMsg = json.msg || "仅能输入英文字符";
                    }
                    break;
                case "domain":
                    if (!$.isDomain(val) && val != "") {
                        errorMsg = json.msg || "域名无效";
                    }
                    break;
                case "mobile":
                    if (!$.isMobile(val) && val != "") {
                        errorMsg = json.msg || "移动电话无效";
                    }
                    break;
                case "int":
                    if (!$.isInt(val) && val != "") {
                        errorMsg = json.msg || "必须输入整数";
                    }
                    break;
                case "email":
                    if (!$.isEmail(val) && val != "") {
                        errorMsg = json.msg || "电子邮件地址无效";
                    }
                    break;
                case "date":
                    if (!$.isDate(val) && val != "") {
                        errorMsg = json.msg || "日期不合法";
                    }
                    break;
                case "regexp":
                    var reg = eval("/^" + json.expression + "$/ig");
                    if (reg.test(val) == false) {
                        errorMsg = json.msg || "输入数据不合法";
                    }
                    break;
                case "range":
                    if ($.isInt(val)) {
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
                case "url":
                    if (!$.isUrl(val) && val != "") {
                        errorMsg = json.msg || "Url地址无效,如：http://www.domain.com";
                    }
                    break;
                case "callback":
                    if (eval(json.fn)(obj.attr("id"), val) == false) {
                        errorMsg = json.msg || "输入内容无效";
                    }
                    break;
                case "tag":
                    if (val.split(json.mark || " ").length > json.max && val != "") {
                        errorMsg = json.msg || "标签个数不能操过" + json.max + "个";
                    }
                    break;
            }
            if (errorMsg != "") {
                labelError.html(errorMsg).attr("class", "errormsg" + errclass).show().stop();
                hasError = true;
                return false;
            }
        });
        if (hasError == false && errorMsg.length == 0) {
            labelError.html("&nbsp;").attr("class", "accept").show();
        }
        return hasError;
    }
};
Biuren.TextEditor = {
    init: function () {
        $(".texteditor").each(function (index, obj) {
            $(obj).hover(function () {
                $(this).addClass("texteditorhover");
            }, function () {
                $(this).removeClass("texteditorhover");
            }).attr("title", "点击修改").click(function () {
                var jobj = $(this);
                jobj.hide();
                var editorid = jobj.attr("editorid");
                if ($("#editorform" + editorid).length == 0) {
                    jobj.after('<form id="editorform' + editorid + '"><input id="inputtext' + editorid + '" name="inputtext" type="input" class="inputs"/><span style="display:none;color:red;">*</span><input id="inputid' + editorid + '" value="' + editorid + '" name="inputid" type="hidden" /><a href="javascript:void(0)" id="linkupdate' + editorid + '">更新</a>  <a href="javascript:void(0)" id="linkcancel' + editorid + '">取消</a>  </form><div></div>');
                    $("#linkcancel" + editorid).click(Biuren.TextEditor.cancelEditor);
                    $("#linkupdate" + editorid).click(Biuren.TextEditor.updateEditor);
                } else {
                    $("#editorform" + editorid).show();
                    $("#inputtext" + editorid).next().hide();
                }
                $("#inputtext" + editorid).val(jobj.text());
            });
        });
    },
    updateEditor: function () {
        var jobj = $(this).parent();
        var input = jobj.find("input[name=inputtext]");
        if (input.val().length == 0) {
            input.next().show();
            return;
        }
        input.next().hide();
        jobj.hide();
        jobj.next().show().html(requestLoading);
        eval(jobj.prev().attr("editorcallback"))(jobj.attr("id"));
    },
    cancelEditor: function () {
        $(this).parent().hide().prev().show();
    },
    resetEditor: function (formid) {
        var jobj = $("#" + formid);
        jobj.next().hide();
        jobj.hide().prev().show();
    },
    updateText: function (formid) {
        var jobj = $("#" + formid);
        var text = jobj.find("input[name=inputtext]").val();
        jobj.next().hide();
        jobj.prev().show().text(text);
    }
};
String.prototype.trim = function () {
    var eleContains = function (arg) {
        var elestr = "$(*+.[?^{|";
        if (arg.length == 1) {
            if (elestr.indexOf(arg) > -1) {
                arg = "\\" + arg;
            } else {
                if (arg == " ") {
                    arg = "\\s";
                }
            }
        } else {
            if (arg.length > 1) {
                var elearray = elestr.split("");
                for (var i = 0; i < elearray.length; i++) {
                    if (arg.indexOf(elearray[i]) > -1) {
                        arg = arg.replace(eval("/\\" + elearray[i] + "/g"), "\\" + elearray[i]);
                    } else {
                        if (arg.indexOf(" ") > -1) {
                            arg = arg.replace(" ", "\\s");
                        }
                    }
                }
            }
        }
        return arg;
    };
    if (arguments.length > 0) {
        var rstr = this;
        var arg = "";
        for (var i = 0; i < arguments.length; i++) {
            arg = arguments[i];
            arg = eleContains(arg);
            rstr = rstr.replace(eval("/(^" + arg + "*)|(" + arg + "*$)/g"), "");
        }
        return rstr;
    } else {
        return this.replace(/(^\s*)|(\s*$)/g, "");
    }
};

var scriptDomain = "http://res.sehui.com";
var siteDomain = "http://www.sehui.com";
var loadingPic = '<img src="' + scriptDomain + '/images/loading.gif">';
var ajaxForm = "/Ajax/";
var requestLoading = loadingPic + "正在处理请求，请稍候";
var waitTime = 800;
var listWaitTime = 500;
var autosavetime = 180000;
var waitCkeditor = 2000;
var waitSwfUpload = 2000;
var waitInitCity = 2000;
var passportUrl = "/AccountHandler.ashx";
var uploadUrl = siteDomain + "/UploadHandler.ashx";
var spaceUrl = "/SpaceHandler.ashx";
var msgUrl = "/MessageHandler.ashx";
var groupUrl = "/GroupHandler.ashx";
var friendUrl = "/FriendHandler.ashx";
var partyUrl = "/PartyHandler.ashx";
var goodsUrl = "/GoodsHandler.ashx";
var adminUrl = "/AdminHandler.ashx";
var sectionUrl = "/SectionHandler.ashx";
var validcodeUrl = "/ValidateCodeHandler.ashx";
var avatarUrl = "/Avatar.ashx";
var delConfirm = "请确认本操作\n同意请点击确定按钮，否则点击取消";
var postSuccessText = "<div class='help'>操作成功</div>";
var nophoto = scriptDomain + "/Images/nophoto.jpg";
var maxHeight = 700;
var maxWidth = 700;
var swfu;
var checkLeave = false;
var boxFooterTxt = "慎勿与色会，色会即祸生";
var tempDiscussOptions = {
    id: "listcommenttable",
    className: "listbiuren",
    noData: "暂时没有评论，快来抢沙发吧！",
    template: {
        head: [],
        row: {
            event: [{
                name: "onmouseover",
                func: "Biuren.setBgColor(this);Biuren.Comm.showedit(this);"
            }, {
                name: "onmouseout",
                func: "Biuren.moveBgColor(this);Biuren.Comm.hideedit(this);"
            }],
            className: "listdark",
            switchClassName: "listlight"
        },
        cell: [{
            text: '<img src="' + avatarUrl + '?uid={0}&size=s" alt="{1}" />',
            value: [6, 5],
            className: "c15 top center"
        }, {
            text: '<a href="/profiles/{5}" >{0}</a>&nbsp;&nbsp;{1}&nbsp;&nbsp;<a href="javascript:void(0)" onclick="Biuren.Comm.replyComment({3},\'{4}\')">回复</a><br><span id="commenttxt{3}">{2}</span><br /><span edit="listcomment" editid="{3}" class="hidden right"></span>',
            value: [5, 4, {
                pos: 3,
                isHtml: true
            }, 0, 5, 6],
            className: "c80 top"
        }],
        foot: {
            text: '{0}<div id="batchsetup" class="right"></div>'
        }
    },
    callback: "Biuren.Edit.init"
};
var tempFCKPhotoOptions = {
    id: "tbllistphoto",
    className: "listbiuren",
    noData: "还没有任何图片",
    template: {
        head: [],
        row: {
            mark: 'height="70px"'
        },
        cellAttr: {
            event: [{
                name: "onmouseover",
                func: "Biuren.setBgColor(this)"
            }, {
                name: "onmouseout",
                func: "Biuren.moveBgColor(this)"
            }]
        },
        cell: [{
            text: "<a href='javascript:void(0)' onclick='Biuren.Comm.insertPhoto(\"{1}\",undefined,undefined,{2},\"{0}\")' title='点击插入图片'><img src='{0}' width='60px' height='60px' /></a>",
            value: [0, 1, 2],
            className: "c15 center"
        }],
        foot: {
            text: "{0}"
        }
    }
};
var noticefrequency = 30000;
var noticecounter = 0;
var noticetimer = 0;
Biuren.Comm = {
    showedit: function (obj) {
        $(obj).find("[edit]").show();
    },
    hideedit: function (obj) {
        $(obj).find("[edit]").hide();
    },
    getErrorText: function (json, simple) {
        var code = json.Result;
        var msg = json.Desc;
        if (code.toString() == "1020") {
            Biuren.Comm.refreshValidateCode();
        }
        if (simple) {
            return msg;
        } else {
            return "错误码:" + code + "  描述:" + msg;
        }
    },
    joinGroup: function (id) {
        Biuren.Box.dialog({
            title: "申请加入群组",
            ajaxUrl: ajaxForm + "Group/Apply.aspx?groupid=" + id,
            width: 550,
            height: 270,
            footerTxt: boxFooterTxt
        });
    },
    viewPointRule: function (id) {
        Biuren.Box.dialog({
            title: "积分规则",
            ajaxUrl: ajaxForm + "Space/PointRule.aspx",
            width: 550,
            height: 270,
            footerTxt: boxFooterTxt
        });
    },
    getNotice: function () {
        if (typeof (loginuid) == "undefined" || loginuid == 0) {
            return;
        }
        window.setTimeout(function () {
            Biuren.Comm.doNotice();
        }, noticefrequency);
    },
    doNotice: function () {
        Biuren.get({
            url: passportUrl,
            data: "action=Profile.getNotice",
            global: false,
            cache: false,
            dataType: "json",
            success: function (json) {
                var msg = "";
                var notice = parseInt(json.Notice);
                if ((notice & 1) == 1) {
                    msg += ' <a href="/message/' + loginuid + '#Message.listInbox" title="有新消息"><img src="' + scriptDomain + '/Images/Notice/email.png" /></a> ';
                }
                if ((notice & 2) == 2) {
                    msg += ' <a href="/message/' + loginuid + '#Message.listFriendRequest" title="有新好友请求"><img src="' + scriptDomain + '/Images/Notice/friend.png" /></a> ';
                }
                if ((notice & 4) == 4) {
                    msg += ' <a href="/message/' + loginuid + '#Message.listNewComment" title="有新评论"><img src="' + scriptDomain + '/Images/Notice/comment.png" /></a> ';
                }
                if (msg.length > 0) {
                    $("#nav_notice").html(msg);
                    noticefrequency += 60000;
                } else {
                    noticefrequency += 10000;
                }
                Biuren.Comm.getNotice();
            }
        });
    },
    showHelp: function () {
        $.setCookie("checkcookie", "ok", 1);
        if ($.getCookie("checkcookie") != "ok") {
            $.showText("help", "需要将cookie功能打开，否则无法正常使用我们的服务", "error");
            return;
        }
        $.delCookie("checkcookie");
        if (!$.browser.msie && ($.browser.version == "6.0" || $.browser.version == "5.5")) {
            $.showText("help", "为了获取更好用户体验，推荐使用以下版本浏览器访问本站：IE 7+，Firefox 2+，Opera 9+，Chrome 2+", "help");
            return;
        }
    },
    resetHeight: function () {
        if (!$.browser.msie) {
            return;
        }
        if ($.browser.version != "6.0") {
            return;
        }
        $(".box").find(".alpha").each(function (index, obj) {
            var jobj = $(obj);
            var pobj = jobj.parent(".box");
            jobj.height(pobj.height());
        });
    },
    loadValidateCode: function (obj) {
        if ($(obj).next(".validcode").length == 0) {
            $(obj).after('&nbsp;&nbsp;&nbsp;&nbsp;<img src="/ValidateCodeHandler.ashx?' + $.rnd() + '" alt="点击刷新验证码" align="absmiddle" class="hand validcode" title="点击图片更换新的验证码，不区分大小写" onclick="Biuren.Comm.refreshValidateCode(this)" />');
        }
    },
    refreshValidateCode: function (obj) {
        if (typeof (obj) != "undefined") {
            $(obj).attr("src", "/ValidateCodeHandler.ashx?" + $.rnd());
        } else {
            $(".validcode").each(function () {
                $(this).attr("src", "/ValidateCodeHandler.ashx?" + $.rnd());
            });
        }
    },
    popClassWin: function (action, typeId, id, targetid) {
        var url = ajaxForm + "MyClass/Edit.aspx?typeid=" + typeId + "&action=" + action + "&target=" + targetid;
        if (typeof (id) != "undefined") {
            url += "&id=" + id + "&rnd=" + $.rnd();
        }
        Biuren.Box.dialog({
            title: action == "Space.insertMyClass" ? "创建分类" : "修改分类",
            ajaxUrl: url,
            width: 400,
            height: 110,
            footerTxt: boxFooterTxt
        });
    },
    getPageActionByType: function (typeid) {
        var pageaction = "";
        switch (typeid.toString()) {
            case "1":
                pageaction = "Space.listGoodsClass";
                break;
            case "2":
                pageaction = "Space.listBlogClass";
                break;
            case "3":
                pageaction = "Space.listPhotoClass";
                break;
            case "4":
                pageaction = "Space.listUrlClass";
                break;
        }
        return pageaction;
    },
    delClass: function (type, id) {
        Biuren.Comm.postAction(spaceUrl, "action=Space.delMyClass&classtype=" + type + "&classid=" + id, null, true, "help", function () {
            var pageaction = Biuren.Comm.getPageActionByType(type);
            Biuren.get({
                url: Pager[pageaction].url,
                cache: false,
                success: function (data) {
                    Biuren.showOK();
                    Biuren.html(Pager[pageaction].container, data);
                }
            });
            Biuren.Box.close();
        });
    },
    doClass: function () {
        if (Biuren.Valid.checkError("popbiurenform")) {
            return;
        }
        Biuren.Comm.postAction(spaceUrl, $("#popbiurenform").serialize(), null, false, "pophelp", function (json) {
            var target = $("#target").val();
            if (target != "undefined") {
                var selectObj = $("#" + target).get(0);
                selectObj.options[selectObj.options.length] = new Option($("#classname").val(), json.id);
                selectObj[selectObj.options.length - 1].selected = true;
                Biuren.Box.close();
            } else {
                var pageaction = Biuren.Comm.getPageActionByType($("#typeid").val());
                Biuren.get({
                    url: Pager[pageaction].url,
                    cache: false,
                    success: function (data) {
                        Biuren.html(Pager[pageaction].container, data);
                        Biuren.showOK();
                    }
                });
                Biuren.Box.close();
            }
        });
    },
    insertPhoto: function (url, target, type, photoId, surl) {
        var inserttarget = typeof (target) == "undefined" ? $("#cachedata").data("target") : target;
        var inserttype = typeof (type) == "undefined" ? $("#cachedata").data("type") : type;
        if (inserttype == "url") {
            $("#" + inserttarget).val(url);
        } else {
            if (inserttype == "party") {
                var imgnum = $("#" + inserttarget + " img").length;
                if (imgnum < 10) {
                    $("#" + inserttarget).append('<li style="width:20%" class="left"><a href="javascript:void(0)" title="点击取消选择图片" onclick="$(this).parent().remove();"><img src="' + surl + '" width="70px" height="70px" pid="' + photoId + '"></a><div class="clearfloat"></div></li>');
                } else {
                    $.showTextInTag("uploadstatus", "最多能够提交10张图片到活动", "error", true);
                }
            } else {
                Biuren.insertEditorHtml(inserttarget, '<p><img class="resizeimg" src="' + url + '" onload="Biuren.Comm.DrawImage(this)" /></p><br />');
            }
        }
    },
    listPhoto: function (id) {
        if (id == "0") {
            return;
        }
        var ajaxOptions = {
            container: "listphoto",
            url: spaceUrl,
            queryParams: "action=Space.listMyPhoto&id=" + id,
            row: 1,
            col: 6,
            pageIndex: 1,
            cache: false
        };
        Biuren.Grid.show(ajaxOptions, tempFCKPhotoOptions);
    },
    selectImage: function (target, type) {
        $("#divuploadarea").hide();
        $("#divinsertarea").show();
        if ($("#divinsertarea").find("select").length > 0) {
            return;
        }
        $("#classoption").html("<select id='myclassoption' name='myclassoption' onchange=\"Biuren.Comm.listPhoto($('#myclassoption').val())\"><option value='0'>请选择相册</option></select>").show();
        Biuren.get({
            data: "action=Space.listMyClass&classtype=3&uid=" + loginuid,
            url: spaceUrl,
            dataType: "json",
            cache: false,
            success: function (json) {
                var selectObj = $("#myclassoption").get(0);
                for (var i = 0; i < json.length; i++) {
                    selectObj.options[selectObj.options.length] = new Option(json[i].classname, json[i].classid);
                }
            }
        });
        $("#cachedata").data("target", target);
        $("#cachedata").data("type", type);
    },
    initCkeditor: function (callback) {
        $.getScript(scriptDomain + "/ckeditor/ckeditor.js", function () {
            $.getScript(scriptDomain + "/ckeditor/config.js", function () {
                window.setTimeout(function () {
                    $(".myckeditor").each(function (i, obj) {
                        var jobj = $(obj);
                        var editor = CKEDITOR.instances[jobj.attr("name")];
                        if (editor) {
                            CKEDITOR.remove(editor);
                        }
                    });
                    CKEDITOR.replaceAll("myckeditor");
                    if (callback) {
                        callback();
                    }
                }, waitCkeditor);
            });
        });
    },
    destorySwfupload: function () {
        if (typeof (swfu) != "undefined") {
            swfu.destroy();
        }
    },
    initSwfupload: function (postParams, uploadFiles, uploadLimit, queueLimit, completeHandler, successHandler) {
        $.loadCSS(scriptDomain + "/Script/swfupload/default.css", "css", "");
        $.getScript(scriptDomain + "/Script/swfupload/swfupload.js", function () {
            window.setTimeout(function () {
                swfu = new SWFUpload({
                    upload_url: uploadUrl,
                    file_post_name: "MyUploadFile",
                    post_params: postParams,
                    button_image_url: scriptDomain + "/Script/swfupload/XPButtonNoText_160x22.png",
                    button_placeholder_id: "spanButtonPlaceholder",
                    button_width: 160,
                    button_height: 22,
                    button_text: '<span class="button">选择图像文件 <span class="buttonSmall">(1 MB Max)</span></span>',
                    button_text_style: ".button { font-family: Helvetica, Arial, sans-serif; font-size: 13pt; } .buttonSmall { font-size: 10pt; }",
                    button_text_top_padding: 1,
                    button_text_left_padding: 5,
                    file_size_limit: "1000",
                    file_types: "*.jpg;*.gif;*.png;*.bmp;",
                    file_types_description: "Images",
                    file_upload_limit: uploadLimit,
                    file_queue_limit: queueLimit,
                    file_queue_error_handler: fileQueueError,
                    file_dialog_complete_handler: completeHandler,
                    upload_progress_handler: uploadProgress,
                    upload_error_handler: uploadError,
                    upload_success_handler: successHandler,
                    upload_complete_handler: uploadComplete,
                    file_queued_handler: fileQueued,
                    flash_url: scriptDomain + "/Script/swfupload/swfupload.swf",
                    custom_settings: {
                        upload_target: "divFileProgressContainer",
                        uploadfiles: uploadFiles
                    },
                    debug: false
                });
            }, waitSwfUpload);
        });
    },
    uploadImage: function (target, type, uploadlimit) {
        $("#divuploadarea").show();
        $("#divinsertarea").hide();
        if ($("#divuploadarea").html().length > 50) {
            return;
        }
        $("#cachedata").data("target", target);
        $("#cachedata").data("type", type);
        Biuren.get({
            cache: false,
            url: ajaxForm + "UploadImage.aspx",
            success: function (data) {
                if (Biuren.checkError(data, "divuploadarea")) {
                    return;
                }
                Biuren.html("divuploadarea", data);
                $.getScript(scriptDomain + "/Script/uploadhandlers.js?v=4", function () {
                    Biuren.Comm.initSwfupload({
                        action: "photo",
                        userdata: $("#userdata").val()
                    }, 1, 20, (type == "img" || type == "party") ? 8 : 1, fileFCKDialogComplete, uploadFCKSuccess);
                });
            }
        });
    },
    postAction: function (dourl, params, refreshid, needconfirm, errrorid, successcallbackfunc, errorcallbackfunc) {
        if (needconfirm) {
            Biuren.Box.confirm({
                title: "操作确认",
                content: delConfirm,
                width: 300,
                height: 50,
                onOk: function () {
                    Biuren.post({
                        data: params,
                        url: dourl,
                        success: function (json) {
                            if (json.Result == 0) {
                                Biuren.showOK();
                                if (typeof (refreshid) != "undefined") {
                                    window.setTimeout('Biuren.Grid.refresh("' + refreshid + '")', waitTime);
                                }
                                if (typeof (successcallbackfunc) != "undefined") {
                                    successcallbackfunc(json);
                                }
                            } else {
                                Biuren.Comm.showErrorMsg(errrorid || "help", Biuren.getErrorText(json), "error");
                                if (typeof (errorcallbackfunc) != "undefined") {
                                    errorcallbackfunc(json);
                                }
                            }
                        }
                    });
                }
            });
        } else {
            Biuren.post({
                data: params,
                url: dourl,
                success: function (json) {
                    if (json.Result == 0) {
                        Biuren.showOK();
                        if (typeof (refreshid) != "undefined") {
                            window.setTimeout('Biuren.Grid.refresh("' + refreshid + '")', waitTime);
                        }
                        if (typeof (successcallbackfunc) != "undefined") {
                            successcallbackfunc(json);
                        }
                    } else {
                        Biuren.Comm.showErrorMsg(errrorid || "help", Biuren.getErrorText(json), "error");
                        if (typeof (errorcallbackfunc) != "undefined") {
                            errorcallbackfunc(json);
                        }
                    }
                }
            });
        }
    },
    doShare: function () {
        if (Biuren.Valid.checkError("popbiurenform")) {
            return;
        }
        Biuren.post({
            data: $("#popbiurenform").serialize(),
            url: spaceUrl,
            success: function (json) {
                if (json.Result == 0) {
                    Biuren.Box.close();
                    Biuren.showOK();
                } else {
                    $.showText("pophelp", Biuren.getErrorText(json), "error");
                }
            }
        });
    },
    winShare: function (classid, id) {
        var title;
        switch (classid) {
            case 6:
                title = "分享空间";
                break;
            case 3:
                title = "分享日志";
                break;
            case 4:
                title = "分享相册";
                break;
            case 5:
                title = "分享图片";
                break;
            case 7:
                title = "分享群组";
                break;
            case 8:
                title = "分享话题";
                break;
            case 1:
                title = "分享网址";
                break;
        }
        Biuren.Box.dialog({
            title: title,
            ajaxUrl: ajaxForm + "Space/Share.aspx?classid=" + classid + "&id=" + id,
            width: 490,
            height: 230,
            footerTxt: boxFooterTxt
        });
    },
    writeMsg: function (title, receiver) {
        var param = new Array();
        if (typeof (title) != "undefined") {
            param.push("title=" + title);
        }
        if (typeof (receiver) != "undefined") {
            param.push("receiver=" + receiver);
        }
        Biuren.Box.dialog({
            title: "写消息",
            ajaxUrl: ajaxForm + "Message/Write.aspx?" + param.join("&"),
            width: 600,
            height: 350,
            footerTxt: boxFooterTxt
        });
    },
    saveMsg: function () {
        if (Biuren.Valid.checkError("popbiurenform")) {
            return;
        }
        Biuren.post({
            data: $("#popbiurenform").serialize(),
            url: msgUrl,
            success: function (json) {
                if (json.Result == 0) {
                    Biuren.Box.close();
                    Biuren.showOK();
                } else {
                    $.showText("pophelp", Biuren.getErrorText(json), "error");
                }
            }
        });
    },
    applyFriend: function (uid) {
        Biuren.Box.dialog({
            title: "申请加为好友",
            ajaxUrl: ajaxForm + "Friend/Add.aspx?uid=" + uid,
            width: 460,
            height: 170,
            footerTxt: boxFooterTxt
        });
    },
    inviteEmail: function () {
        Biuren.Box.dialog({
            title: "邮件邀请好友",
            ajaxUrl: ajaxForm + "Email.aspx",
            width: 460,
            height: 200,
            footerTxt: boxFooterTxt
        });
    },
    showDiy: function () {
        Biuren.Box.dialog({
            title: "自定义",
            ajaxUrl: ajaxForm + "Email.aspx",
            width: 460,
            height: 200,
            footerTxt: boxFooterTxt
        });
    },
    doInviteEmail: function () {
        if (Biuren.Valid.checkError("popbiurenform")) {
            return;
        }
        Biuren.post({
            data: $("#popbiurenform").serialize(),
            url: msgUrl,
            success: function (json) {
                if (json.Result == 0) {
                    Biuren.Box.close();
                    Biuren.showOK();
                } else {
                    $.showText("pophelp", Biuren.getErrorText(json), "error");
                }
            }
        });
    },
    doApplyFriend: function () {
        if (Biuren.Valid.checkError("popbiurenform")) {
            return;
        }
        Biuren.post({
            data: $("#popbiurenform").serialize(),
            url: friendUrl,
            success: function (json) {
                if (json.Result == 0) {
                    Biuren.Box.close();
                    Biuren.showOK();
                } else {
                    $.showText("pophelp", Biuren.getErrorText(json), "error");
                }
            }
        });
    },
    DrawImage: function (ImgD) {
        var image = new Image();
        image.src = ImgD.src;
        if (image.width > 0 && image.height > 0) {
            if (image.width / image.height >= maxHeight / maxWidth) {
                if (image.width > maxHeight) {
                    ImgD.width = maxHeight;
                    ImgD.height = (image.height * maxWidth) / image.width;
                } else {
                    ImgD.width = image.width;
                    ImgD.height = image.height;
                }
            } else {
                if (image.height > maxHeight) {
                    ImgD.height = maxHeight;
                    ImgD.width = (image.width * maxHeight) / image.height;
                } else {
                    ImgD.width = image.width;
                    ImgD.height = image.height;
                }
            }
        }
    },
    refresh: function (obj) {
        var jobj = $(obj);
        var path = window.location.pathname;
        if (jobj.attr("href").indexOf(path) != -1) {
            window.location.href = jobj.attr("href");
            window.location.reload();
        }
    },
    logout: function () {
        $.post(passportUrl, "action=logout", function () {
            window.location = "/index";
        });
    },
    delComment: function (id) {
        Biuren.Comm.postAction(spaceUrl, "action=Space.delBatchComment&chkvalue=" + id, "listcommenttable", true);
    },
    showErrorMsg: function (id, msg, className) {
        $.showText(id, msg, className);
        $("html").animate({
            scrollTop: $("#" + id).offset().top - 45
        }, 500);
    },
    initEditComment: function () {
        $("[edit=listcomment]").each(function (pos, obj) {
            var jobj = $(obj);
            var id = jobj.attr("editid");
            jobj.get(0).innerHTML = '<a href="javascript:void(0)" onclick="Biuren.Comm.delComment(' + id + ')">删除</a> ';
        });
    },
    initMenu: function () {
        var menus = $(".ml li");
        if (menus.length == 0) {
            return;
        }
        menus.mouseover(function () {
            $(this).addClass("nav_a_hover");
        }).mouseout(function () {
            $(this).removeClass("nav_a_hover");
        });
        var url = window.location.pathname;
        $(".ml a").each(function (pos, obj) {
            var jobj = $(obj);
            var href = jobj.attr("href");
            if (url.toLowerCase() == jobj.attr("href")) {
                jobj.parent().addClass("nav_a_focus");
                return false;
            }
        });
    }
};
Biuren.Edit = {
    menuInit: function () {
        if (typeof (loginuid) == "undefined" || loginuid == 0) {
            return;
        }
        if ((typeof (nowuid) != "undefined" && nowuid == loginuid) || typeof (groupadmin) != "undefined" && groupadmin == true) {
            $(".editmenulink").show();
            $(".editlink").show();
        }
    },
    init: function () {
        var anchor = $("#cachedata").data("anchor") != undefined ? $("#cachedata").data("anchor") : $.anchor();
        if (typeof (loginuid) == "undefined" || loginuid == 0) {
            return;
        }
        if (typeof (groupadmin) != "undefined" && groupadmin == false) {
            return;
        }
        if (typeof (nowuid) != "undefined" && nowuid != loginuid) {
            return;
        }
        $(".editlink").show();
        switch (anchor) {
            case "Space.listMyBlog":
            case "Space.listMyBlogDraft":
                Biuren.Space.initEditBlog();
                break;
            case "Space.listPhoto":
                Biuren.Space.initEditPhoto();
                Biuren.TextEditor.init();
                break;
            case "Space.listAlbum":
                Biuren.Space.initEditAlbum();
                Biuren.TextEditor.init();
                break;
            case "Group.listPlacard":
                Biuren.Group.initEditPlacard();
                break;
            case "Group.listGroupUrl":
                Biuren.Group.initEditGroupUrl();
                break;
            case "Group.listGroupUser":
                Biuren.Group.initEditGroupUser();
                break;
            case "Group.listGroupManager":
                Biuren.Group.initEditGroupManager();
                break;
            case "Friend.listMyFriend":
                Biuren.Friend.initEditMyFriend();
                break;
            case "Space.listPlaylist":
                Biuren.Space.initListPlaylist();
                break;
            case "Space.listShare":
                Biuren.Space.initListShare();
                break;
            case "Space.listDiyTemplate":
                Biuren.TextEditor.init();
                Biuren.Space.initListTemplate();
                break;
            case "Group.listGroupApply":
                Biuren.Group.initListGroupApply();
                break;
            case "Group.listJoinGroup":
            case "Group.listManageGroup":
                Biuren.Group.initEditGroup();
                break;
            case "Space.listMyGoods":
                Biuren.Goods.initEditGoods();
                break;
            case "Space.listMyUrl":
                Biuren.Space.initEditUrl();
                break;
            case "Party.listMyParty":
                Biuren.Party.initEditParty();
                break;
            case "Party.listFocusParty":
                Biuren.Party.initEditFocusParty();
                break;
            case "Party.listMemberByAdmin":
                Biuren.Party.initAdminMember();
                break;
            case "Party.listJoinedParty":
                Biuren.Party.initEditJoinedParty();
                break;
            case "Space.listComment":
                Biuren.Comm.initEditComment();
                break;
        }
    }
};
$(function () {
    $.ajaxSetup({
        cache: false,
        global: true,
        error: function (xhr, desc) {
            Biuren.Comm.showErrorMsg("help", Biuren.getErrorText({
                Result: 1000,
                Desc: "服务器处理错误，请待会再尝试"
            }), "error");
        }
    });
    $("body").ajaxSend(function (event, request, options) {
        $.clearText("help", "pophelp", "status");
        Biuren.Pager.setAjaxStatus(requestLoading, "help");
        if (options.type.toUpperCase() == "POST") {
            $(":button").attr("disabled", true);
        }
    });
    $("body").ajaxComplete(function (event, request, options) {
        Biuren.Pager.clearAjaxStatus();
        if (options.type.toUpperCase() == "GET") {
            Biuren.Valid.init();
            Biuren.Tooltips.init();
            Biuren.Tab.init();
        } else {
            Biuren.Pager.setAjaxStatus("操作成功", "help");
            $(":button").attr("disabled", false);
        }
    });
    Biuren.Pager.init();
    Biuren.Tab.init();
    Biuren.Toggle.init();
    Biuren.Tooltips.init();
    Biuren.Edit.menuInit();
    Biuren.Valid.init();
    Biuren.Slider.init();
    Biuren.Comm.getNotice();
    Biuren.Comm.initMenu();
    Biuren.Comm.showHelp();
    window.onbeforeunload = function () {
        if (checkLeave) {
            return ("您有设置未保存，请确定是否要放弃！");
        }
    };
});

Biuren.Passport = {
    login: function () {
        if (Biuren.Valid.checkError("biurenform")) {
            return;
        }
        Biuren.post({
            data: $("#biurenform").serialize() + "&action=login",
            url: passportUrl,
            success: function (json) {
                if (json.Result == 0) {
                    switch (json.State) {
                        case 1:
                            Biuren.Pager.setAjaxStatus("登录成功，正在跳转页面", "help");
                            var gotourl = $.QueryString("back");
                            window.location = gotourl ? decodeURIComponent(gotourl) : "/index";
                            break;
                        case 2:
                            $.showText("help", '注册没有确认，请点击注册邮箱中确认链接，或<a href="javascript:void(0)" onclick="Biuren.Passport.openGetCode()">重新接收确认邮件</a>', "error");
                            break;
                        case 3:
                            Biuren.Comm.showErrorMsg("help", '你需要<a href="/reg/upload">上传个人图片</a>完成注册，网页将自动跳转', "help");
                            window.setTimeout(function () {
                                window.location = "/reg/upload";
                            }, 2000);
                            break;
                        case 4:
                            Biuren.Comm.showErrorMsg("help", '你账户已经审核通过，但需要<a href="/reg/profiles">完善个人资料</a>才能使用网站，网页将自动跳转', "help");
                            window.setTimeout(function () {
                                window.location = "/reg/profiles";
                            }, 2000);
                            break;
                    }
                } else {
                    Biuren.Comm.showErrorMsg("help", Biuren.getErrorText(json), "error");
                }
            }
        });
    },
    reg: function () {
        if (Biuren.Valid.checkError("biurenform")) {
            return;
        }
        Biuren.post({
            url: passportUrl,
            data: $("#biurenform").serialize() + "&action=reg",
            success: function (json) {
                if (json.Result == 0) {
                    if (json.State == 1) {
                        Biuren.Pager.setAjaxStatus("登录成功，正在跳转页面", "help");
                        var gotourl = $.QueryString("back");
                        window.location = gotourl ? decodeURIComponent(gotourl) : "/index";
                    } else {
                        $.showText("biurenform", "注册成功！您注册邮箱将会收到激活邮件，确认后完成激活！", "help");
                    }
                } else {
                    Biuren.Comm.showErrorMsg("help", Biuren.getErrorText(json), "error");
                }
            }
        });
    },
    checkName: function () {
        var result = true;
        Biuren.post({
            url: passportUrl,
            async: false,
            data: "action=checkname&username=" + encodeURIComponent($("#username").val()),
            success: function (json) {
                if (json.Result != true) {
                    result = false;
                }
            }
        });
        return result;
    },
    checkEmail: function (id, val) {
        var result = true;
        if (val.length == 0) {
            return;
        }
        Biuren.post({
            url: passportUrl,
            async: false,
            data: "action=checkemail&email=" + encodeURIComponent(val),
            success: function (json) {
                if (json.Result != true) {
                    result = false;
                }
            }
        });
        return result;
    },
    sendValidCodeEmail: function () {
        Biuren.Comm.postAction(passportUrl, $("#popbiurenform").serialize() + "&action=Profile.sendValidCodeEmail", "", false, "pophelp", function () {
            $("#popbiurenform").html('<div class="notice">已给注册邮箱发送修改密码链接，请查收！</div>');
        });
    },
    updatePwdByAnswer: function () {
        if (Biuren.Valid.checkError("popbiurenform")) {
            return;
        }
        Biuren.Comm.postAction(passportUrl, $("#popbiurenform").serialize() + "&action=Profile.updatePwdByAnswer", "", false, "pophelp", function () {
            $("#popbiurenform").html('<div class="notice">密码已经成功修改！</div>');
        });
    },
    updatePwdByValidCode: function () {
        if (Biuren.Valid.checkError("biurenform")) {
            return;
        }
        Biuren.Comm.postAction(passportUrl, $("#biurenform").serialize() + "&action=Profile.updatePwdByValidCode", "", false, "help", function () {
            $("#biurenform").html('<div class="notice">密码已经成功修改！页面将自动跳转到首页</div>');
            window.setTimeout(function () {
                window.location = "/index";
            }, 2000);
        });
    },
    checkAnswer: function () {
        if (Biuren.Valid.checkError("popbiurenform")) {
            return;
        }
        Biuren.post({
            url: passportUrl,
            data: $("#popbiurenform").serialize() + "&action=Profile.checkAnswer",
            success: function (json) {
                if (json.Result == 0) {
                    Biuren.html("resetarea", '<br />新&nbsp;密&nbsp;码：<input type="password" valid="[{type:\'limit\',min:6,max:30}]" class="inputs" id="rpwd" name="rpwd" /><p>确认密码：<input type="password" valid="[{type:\'compareto\',destid:\'rpwd\'}]" class="inputs" id="rconfirmpwd" name="rconfirmpwd" /></p><input type="button" class="formbtn" onclick="Biuren.Passport.updatePwdByAnswer()" value="重设密码" />');
                    $("#findpwdbyemail").hide();
                    $("#answerarea").hide();
                    Biuren.Valid.init();
                } else {
                    $.showText("pophelp", Biuren.getErrorText(json), "error");
                }
            }
        });
    },
    getPwdInfo: function () {
        if (Biuren.Valid.checkError("popbiurenform")) {
            return;
        }
        Biuren.post({
            url: passportUrl,
            data: $("#popbiurenform").serialize() + "&action=Profile.getPwdInfo",
            success: function (json) {
                if (json.Result == 0) {
                    $("#findpwdsetp1").hide();
                    $("#findpwdsetp2").show();
                    var question = json.Question;
                    if (question.length > 0) {
                        Biuren.html("answerarea", '<span id="findpwdquestion">问题：' + question + '</span><p>答案：<input name="answer" id="answer" type="text" maxlength="250" size="30" valid="[{type:\'limit\',min:4,max:150}]" /></p><input type="button" class="formbtn" onclick="Biuren.Passport.checkAnswer()" value="提交" />');
                        $("#findpwdbyquestion").show();
                        Biuren.Valid.init();
                    } else {
                        $("#findpwdbyquestion").hide();
                    }
                    $("#findusername").val(json.UserName);
                } else {
                    $.showText("pophelp", Biuren.getErrorText(json), "error");
                }
            }
        });
    },
    resendCode: function () {
        if (Biuren.Valid.checkError("popbiurenform")) {
            return;
        }
        Biuren.Comm.postAction(passportUrl, $("#popbiurenform").serialize() + "&action=Profile.resendCode", "", false, "pophelp", function () {
            $("#popbiurenform").html('<div class="notice">已发送确认邮件，请查收！</div>');
        });
    },
    openFindPwd: function () {
        Biuren.Box.dialog({
            title: "找回密码",
            ajaxUrl: ajaxForm + "ForgetPwd.html",
            width: 490,
            height: 230,
            footerTxt: boxFooterTxt
        });
    },
    openGetCode: function () {
        Biuren.Box.dialog({
            title: "重新获取确认码",
            ajaxUrl: ajaxForm + "GetCode.html",
            width: 490,
            height: 260,
            footerTxt: boxFooterTxt
        });
    }
};
var Pager = {
    reg: {
        container: "ajaxcontainer",
        url: ajaxForm + "Reg.html",
        nav: "免费注册"
    },
    login: {
        container: "ajaxcontainer",
        url: ajaxForm + "Login.html",
        nav: "登录"
    }
};