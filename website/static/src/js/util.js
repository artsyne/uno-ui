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

(($) => {
    const methods = {
      disableFormInputs(form, disabled) {
        form.find(':input').prop('disabled', disabled);
      },
    };
  
    $.utils = methods;
  })(jQuery);