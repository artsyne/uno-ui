(function ($) {
  const validators = {
    validname: {
      regex: /^[a-zA-Z_\u4e00-\u9fa5]+$/,
      msg: '请输入有效字符串，应为：英文字符串，下划线，中文',
    },
    en: {
      regex: /^[a-zA-Z]+$/,
      msg: '仅能输入英文字符',
    },
    domain: {
      regex: /^[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/,
      msg: '域名无效',
    },
    mobile: {
      regex: /^1[3-9]\d{9}$/,
      msg: '移动电话无效',
    },
    int: {
      regex: /^\d+$/,
      msg: '必须输入整数',
    },
    email: {
      regex: /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/,
      msg: '电子邮件地址无效',
    },
    date: {
      regex: /^\d{4}-\d{2}-\d{2}$/,
      msg: '日期不合法',
    },
    url: {
      regex: /^(https?:\/\/)?([\da-z.-]+)\.([a-z.]{2,6})([/\w .-]*)*\/?$/,
      msg: 'Url地址无效,如：http://www.domain.com',
    },
  };

  function getStringLength(str) {
    return [...str].length;
  }

  function updateCharCount(obj, json) {
    const remainingChars = json.max - getStringLength(obj.val());
    obj.nextAll('span:first').html(`可输入${remainingChars}个字符`);
  }

  function validateElement(obj) {
    const validJson = JSON.parse(obj.attr('valid'));
    const labelError = obj.nextAll('div:first');
    let hasError = false;
    let errorMsg = '';

    for (const json of validJson) {
      const type = json.type || json;
      const val = obj.val();

      switch (type) {
        case 'emptyquit':
          if (getStringLength(val) === 0) return false;
          break;
        case 'required':
          if (getStringLength(val) === 0) {
            errorMsg = json.msg || '必须填写';
          }
          break;
        case 'limit':
          const len = getStringLength(val);
          if (json.min !== undefined && len < json.min) {
            errorMsg = json.msg || `长度必须大于等于${json.min}个字符`;
          }
          if (json.max !== undefined && len > json.max) {
            errorMsg = json.msg || `长度必须小于等于${json.max}个字符`;
          }
          break;
        case 'validname':
        case 'en':
        case 'domain':
        case 'mobile':
        case 'int':
        case 'email':
        case 'date':
        case 'url':
          if (
            validators[type] &&
            !validators[type].regex.test(val) &&
            val !== ''
          ) {
            errorMsg = json.msg || validators[type].msg;
          }
          break;
        case 'regexp':
          const reg = new RegExp(`^${json.expression}$`, 'ig');
          if (!reg.test(val)) {
            errorMsg = json.msg || '输入数据不合法';
          }
          break;
        case 'range':
          if (/^\d+$/.test(val)) {
            const numVal = parseInt(val, 10);
            if (
              (json.min !== undefined && numVal < json.min) ||
              (json.max !== undefined && numVal > json.max)
            ) {
              errorMsg =
                json.msg || `输入数字应该在${json.min}到${json.max}之间`;
            }
          } else {
            errorMsg = '请输入数字';
          }
          break;
        case 'compareto':
          if ($(`#${json.destid}`).val() !== val || val === '') {
            errorMsg = json.msg || '请输入正确的内容';
          }
          break;
        case 'callback':
          if (eval(json.fn)(obj.attr('id'), val) === false) {
            errorMsg = json.msg || '输入内容无效';
          }
          break;
        case 'tag':
          if (val.split(json.mark || ' ').length > json.max && val !== '') {
            errorMsg = json.msg || `标签个数不能超过${json.max}个`;
          }
          break;
      }

      if (errorMsg) {
        labelError.text(errorMsg).attr('class', 'invalid-feedback');
        obj.removeClass('is-invalid is-valid').addClass('is-invalid');
        hasError = true;
        break;
      }
    }

    if (!hasError) {
      labelError.text('');
      obj.removeClass('is-invalid is-valid').addClass('is-valid');
    }

    return hasError;
  }

  const methods = {
    init() {
      this.find('[valid]').each(function () {
        const obj = $(this);
        if (obj.attr('loadedvalid') !== undefined) return;

        obj.after('<div />').attr('loadedvalid', 'true');
        obj.blur(() => validateElement(obj));

        const validJson = JSON.parse(obj.attr('valid'));
        validJson.forEach((json) => {
          if (json.show) {
            obj.after('<span />');
            obj.on('keydown keyup', () => {
              const value = obj.val();
              if (getStringLength(value) > json.max) {
                obj.val(value.slice(0, json.max));
              }
              updateCharCount(obj, json);
            });
            updateCharCount(obj, json);
          }
        });
      });
    },

    checkError() {
      const errors = this.find('[valid]')
        .map(function () {
          return validateElement($(this));
        })
        .get();
      return errors.some((error) => error);
    },
  };

  $.fn.valid = function (method) {
    if (methods[method]) {
      return methods[method].apply(
        this,
        Array.prototype.slice.call(arguments, 1),
      );
    } else if (typeof method === 'object' || !method) {
      return methods.init.apply(this, arguments);
    } else {
      $.error(`Method ${method} does not exist on jQuery.valid`);
    }
  };

  $(document).ready(() => $('form').valid());
})(jQuery);
