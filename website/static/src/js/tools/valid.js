/**
 * jQuery表单验证插件
 * 
 * 这个插件提供了一种简单而强大的方式来验证表单输入。
 * 它支持多种验证类型，包括内置验证和自定义回调验证。
 * 
 * @author nobody
 * @version 1.0.3
 * @license MIT
 * 
 * 使用说明：
 * 
 * 1. 在HTML中为需要验证的输入元素添加valid属性：
 *    <input type="text" name="username" valid='["required","validname"]'>
 * 
 * 2. 支持的验证类型：
 *    - required: 必填字段
 *    - validname: 有效的用户名（字母、数字、下划线、连字符）
 *    - en: 仅英文字母
 *    - domain: 有效的域名
 *    - mobile: 手机号码
 *    - int: 整数
 *    - email: 电子邮件地址
 *    - date: 日期
 *    - url: URL地址
 *    - pwd: 密码（至少8个字符，包含字母和数字）
 *    - strongpwd: 强密码（至少8个字符，包含字母、数字和特殊字符）
 *    - limit: 字符长度限制
 *    - range: 数值范围
 *    - compareto: 与其他字段比较
 *    - callback: 自定义回调函数
 *    - tag: 标签数量限制
 * 
 * 3. 自定义回调函数示例：
 *    $.fn.valid.registerCallback('checkUsername', function(id, value) {
 *      return value.length >= 3 && value.length <= 20;
 *    });
 * 
 *    然后在HTML中使用：
 *    <input type="text" name="username" valid='[{"type":"callback", "fn":"checkUsername", "msg":"用户名长度必须在3到20个字符之间"}]'>
 * 
 * 4. 手动触发验证：
 *    $('#myForm').valid('checkError');
 * 
 * 5. 显示字符计数：
 *    <input type="text" name="description" valid='[{"type":"limit", "max": 100, "show": true}]'>
 * 
 * 注意：确保在使用此插件之前已经加载了jQuery库。
 */

(function ($) {
  // 私有函数：计算字符串长度（使用标准的字符串长度属性）
  const strLen = (str) => str.length;

  // 私有函数：按字符数截取字符串（使用标准的substring方法）
  const substrByChar = (str, charLength) => str.substring(0, charLength);

  // 回调函数注册对象
  const callbackRegistry = {};

  $.fn.valid = function () {
    const validators = {
      validname: {
        regex: /^[a-zA-Z0-9_-]+$/,
        msg: '请输入有效的用户名（包含字母、数字、下划线和连字符）',
      },
      en: {
        regex: /^[a-zA-Z]+$/,
        msg: '请只输入英文字母（A-Z, a-z）',
      },
      domain: {
        regex: /^((?!-)[A-Za-z0-9-]{1,63}(?<!-)\.)+[A-Za-z]{2,6}$/,
        msg: '请输入有效的域名（如：example.com）',
      },
      mobile: {
        regex: /^\+?[\d\s-]{7,}$/,
        msg: '请输入有效的手机号码（如：+86 123 4567 8901 或 12345678901）',
      },
      int: {
        regex: /^-?\d+$/,
        msg: '请输入整数（如：42 或 -17）',
      },
      email: {
        regex:
          /^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9-]+(?:\.[a-zA-Z0-9-]+)*$/,
        msg: '请输入有效的电子邮件地址（如：user@example.com）',
      },
      date: {
        regex: /^(\d{4}-\d{2}-\d{2}|\d{2}\/\d{2}\/\d{4}|\d{2}\.\d{2}\.\d{4})$/,
        msg: '请输入有效的日期（格式：YYYY-MM-DD、DD/MM/YYYY 或 DD.MM.YYYY）',
      },
      url: {
        regex:
          /^(https?:\/\/)?([\da-z\.-]+)\.([a-z\.]{2,6})([\/\w \.-]*)*\/?$/i,
        msg: '请输入有效的URL地址（如：http://www.example.com 或 www.example.com）',
      },
      pwd: {
        regex: /^(?=.*[A-Za-z])(?=.*\d)[A-Za-z\d]{8,}$/,
        msg: '请输入有效的密码（至少8个字符，包含至少一个字母和一个数字）',
      },
      strongpwd: {
        regex: /^(?=.*[A-Za-z])(?=.*\d)(?=.*[@$!%*#?&])[A-Za-z\d@$!%*#?&]{8,}$/,
        msg: '请输入有效的强密码（至少8个字符，包含至少一个字母、一个数字和一个特殊字符 @$!%*#?&）',
      },
    };

    const methods = {
      /**
       * 初始化验证
       */
      init() {
        this.find('[valid]').each((_, obj) => {
          const $obj = $(obj);
          if ($obj.data('loadedvalid')) return;

          $obj.after('<div />').data('loadedvalid', true);
          $obj.on('input blur', () => methods.checkElement($obj));

          const validRules = JSON.parse($obj.attr('valid'));
          validRules.forEach((json) => {
            if (json.show === true) {
              $obj.after('<span />');
              $obj.on('input', methods.updateCharCount);
              methods.updateCharCount.call($obj[0]);
            }
          });
        });
      },

      /**
       * 更新字符计数
       */
      updateCharCount() {
        const $obj = $(this);
        const value = $obj.val();
        const validRules = JSON.parse($obj.attr('valid'));

        for (const json of validRules) {
          if (json.show === true) {
            const maxlen = json.max;
            if (strLen(value) > maxlen) {
              $obj.val(substrByChar(value, maxlen));
            }
            $obj
              .nextAll('span:first')
              .html(`可输入${maxlen - strLen($obj.val())}个字符`);
            break;
          }
        }
      },

      /**
       * 检查所有元素的错误
       * @returns {boolean} 是否有错误
       */
      checkError() {
        const errors = Array.from($(this).find('[valid]')).map((obj) =>
          methods.checkElement($(obj)),
        );
        return errors.some(Boolean);
      },

      /**
       * 检查单个元素的验证
       * @param {jQuery} $obj - 要检查的jQuery对象
       * @returns {boolean} 是否有错误
       */
      checkElement($obj) {
        const validRules = JSON.parse($obj.attr('valid'));
        const $labelError = $obj.nextAll('div:first');
        let hasError = false;
        let errorMsg = '';

        for (const json of validRules) {
          const type = json.type || json;
          const val = $obj.val();

          switch (type) {
            case 'emptyquit':
              if (strLen(val) === 0) return false;
              break;
            case 'required':
              if (strLen(val) === 0) errorMsg = json.msg || '必须填写';
              break;
            case 'limit':
              const len = strLen(val);
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
              const reg = new RegExp(`^${json.expression}$`, 'iu');
              if (!reg.test(val)) errorMsg = json.msg || '输入数据不合法';
              break;
            case 'range':
              if (/^\d+$/.test(val)) {
                const numVal = parseInt(val, 10);
                if (numVal < json.min || numVal > json.max) {
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
              if (typeof callbackRegistry[json.fn] === 'function') {
                if (callbackRegistry[json.fn]($obj.attr('id'), val) === false) {
                  errorMsg = json.msg || '输入内容无效';
                }
              } else {
                console.error(`Callback function ${json.fn} is not registered`);
              }
              break;
            case 'tag':
              if (val.split(json.mark || ' ').length > json.max && val !== '') {
                errorMsg = json.msg || `标签个数不能超过${json.max}个`;
              }
              break;
          }

          if (errorMsg) {
            $labelError.text(errorMsg).attr('class', 'invalid-feedback');
            $obj.removeClass('is-invalid is-valid').addClass('is-invalid');
            hasError = true;
            break;
          }
        }

        if (!hasError) {
          $labelError.text('');
          $obj.removeClass('is-invalid is-valid').addClass('is-valid');
        }
        return hasError;
      },
    };

    if (typeof arguments[0] === 'string' && methods[arguments[0]]) {
      return methods[arguments[0]].apply(
        this,
        Array.prototype.slice.call(arguments, 1),
      );
    } else if (typeof arguments[0] === 'object' || !arguments[0]) {
      return methods.init.apply(this, arguments);
    } else {
      $.error('Method ' + arguments[0] + ' does not exist on jQuery.valid');
    }
  };

  /**
   * 注册自定义回调函数
   * @param {string} name - 回调函数的名称
   * @param {function} callback - 回调函数
   */
  $.fn.valid.registerCallback = function(name, callback) {
    if (typeof callback === 'function') {
      callbackRegistry[name] = callback;
    } else {
      console.error(`Invalid callback function for ${name}`);
    }
  };

  // 使用 jQuery 的简写语法和方法链
  $(() => $('form').valid());
})(jQuery);