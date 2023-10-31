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

Biuren.uploadImg = {
    upload: function (parametric) {
        var thisObj = {
            maxWidth: parametric.maxNum || 680,
            rate: parametric.rate || 0.8,
            uploadFileId: parametric.fileId,
            callback: parametric.callback
        }

        $("#" + thisObj.uploadFileId).on("change", function (e) {
            if ($(this).val() == "") {
                return;
            }
            if (this.files.length === 0) {
                return;
            }
            var file = this.files[0];
            var reader0 = new FileReader();
            reader0.readAsDataURL(file);
            reader0.onload = function (e) {
                var dataUrl = this.result;
                var image = new Image();
                image.src = dataUrl;
                image.onload = function (e) {
                    var width = image.width,
                        height = image.height;
                    var scale = width / height;
                    var canvas = document.createElement("canvas");
                    var ctx = canvas.getContext('2d');
                    ctx.clearRect(0, 0, canvas.width, canvas.height);
                    if (width - height >= 0) {
                        width1 = thisObj.maxWidth;
                        height1 = parseInt(width1 / scale);
                        if (width >= width1) {
                            canvas.width = width1;
                            canvas.height = height1;
                            ctx.drawImage(image, 0, 0, width1, height1);
                        } else {
                            canvas.width = width;
                            canvas.height = height;
                            ctx.drawImage(image, 0, 0, width, height);
                        }
                    } else {
                        scale = height / width;
                        height1 = thisObj.maxWidth;
                        width1 = parseInt(height1 / scale);
                        if (height >= height1) {
                            canvas.width = width1;
                            canvas.height = height1;
                            ctx.drawImage(image, 0, 0, width1, height1);
                        } else {
                            canvas.width = width;
                            canvas.height = height;
                            ctx.drawImage(image, 0, 0, width, height);
                        }
                    }
                    var cropStr = canvas.toDataURL("image/jpeg", thisObj.rate);
                    thisObj.callback(cropStr);
                }
            }
        })
    }
}

Biuren.Valid = {
    imagePath: "",
    init: function () {
        $("form").find("[valid]").each(function (i, obj) {
            var obj = $(obj);
            if (obj.attr("loadedvalid") != undefined) {
                return;
            }
            obj.after("<div />").attr("loadedvalid", "true");
            obj.blur(function () {
                var obj = $(this);
                Biuren.Valid.checkElement(obj);
            });
            $($.evalJSON(obj.attr("valid"))).each(function (i, json) {
                if (json.show == true) {
                    obj.after("<span />");
                    obj.keydown(Biuren.Valid.pressKey).keyup(Biuren.Valid.pressKey);
                    obj.nextAll("span:first").html("可输入" + (json.max - $.len(obj.val())) + "个字符");
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
                obj.nextAll("span:first").html("可输入" + (json.max - $.len(obj.val())) + "个字符");
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
        var labelError = obj.nextAll("div:first");
        var hasError = false;
        var errorMsg = "";
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

var AJAX_URL = '';

Biuren.Passport = {
    login: function () {
        if (Biuren.Valid.checkError("biurenform")) {
            return;
        }
        $.post(AJAX_URL + '/login', $("#biurenform").serialize(), function (data) {
            window.location = data.AccessUrl;
        }, 'json').fail(function (res) {
            console.log(res.responseJSON);
        });
    },
    reg: function () {
        if (Biuren.Valid.checkError("biurenform")) {
            return;
        }
        $('#time_zone').val(Intl.DateTimeFormat().resolvedOptions().timeZone);
        $('#lang').val(navigator.language);
        $.post(AJAX_URL + '/register', $("#biurenform").serialize(), function (data) {
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
                    "<div class=\"card-body\">" +
                    "<h5>" + $.HtmlEncode(data.Data.title) + "</h5>" +
                    "<p>0篇</p></div></div></div>";
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
    Biuren.Valid.init();
    Biuren.Pager.init();

    Biuren.Diary.init();

    //header fixed
    $(document).scroll(function () {
        if ($(document).scrollTop() > 1)
            $('.top-nav').html($('#page-nav').html());
        else if ($(document).scrollTop() == 0)
            $('.top-nav').html('');
    });
});