String.prototype.trim = function (arg) {
    // 特殊字符转义函数
    const escapeSpecialChars = (arg) => {
        const specialChars = "$(*+.[?^{|";
        if (arg.length === 1) {
            return specialChars.includes(arg) ? `\\${arg}` : arg === " " ? "\\s" : arg;
        } else {
            return arg.split("").map(char => specialChars.includes(char) ? `\\${char}` : char === " " ? "\\s" : char).join("");
        }
    };

    // 如果有参数，按参数进行 trim
    if (arg !== undefined) {
        const escapedArg = escapeSpecialChars(arg);
        const regex = new RegExp(`(^${escapedArg}*)|(${escapedArg}*$)`, "g");
        return this.replace(regex, "");
    } else {
        // 如果没有参数，默认去除空白字符
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
