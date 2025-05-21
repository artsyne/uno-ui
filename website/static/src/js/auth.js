/**
 * 认证模块 - 处理登录和注册功能
 *
 * 基于Kiko UI的页面行为管理组件实现
 *
 * @version 2.1.0
 * @license MIT
 */

/**
 * 认证行为配置对象
 * @namespace Auth
 */
const Auth = {
  /**
   * 初始化认证模块
   *
   * @param {string} baseUrl - API基础URL
   * @returns {Array} 认证行为配置数组
   */
  init: function (baseUrl = '') {
    log('认证模块初始化');
    this.baseUrl = baseUrl;

    // 注册自定义验证回调
    this.registerValidationCallbacks();

    // 返回认证行为配置数组
    return [this.loginAction, this.registerAction, this.forgotPasswordAction];
  },

  /**
   * 注册自定义验证回调
   */
  registerValidationCallbacks: function () {
    // 用户名查重验证
    $.fn.valid.registerCallback('checkUsername', function (id, value) {
      // 如果用户名为空，不进行查重
      if (!value) return true;
      // 返回Promise对象，以支持异步验证
      return new Promise((resolve) => {
        // 调用后端API检查用户名是否存在
        $.ajax({
          url: `${Auth.baseUrl}/check-username`,
          type: 'POST',
          data: { username: value },
          dataType: 'json',
          success: function (response) {
            // 返回false表示验证失败（用户名已存在）
            resolve(!response.exists);
          },
          error: function () {
            // 出错时默认通过验证
            resolve(true);
          },
        });
      });
    });

    // 邮箱查重验证
    $.fn.valid.registerCallback('checkEmail', function (id, value) {
      // 如果邮箱为空，不进行查重
      if (!value) return true;

      // 返回Promise对象，以支持异步验证
      return new Promise((resolve) => {
        // 调用后端API检查邮箱是否存在
        $.ajax({
          url: `${Auth.baseUrl}/check-email`,
          type: 'POST',
          data: { email: value },
          dataType: 'json',
          success: function (response) {
            // 返回false表示验证失败（邮箱已存在）
            resolve(!response.exists);
          },
          error: function () {
            // 出错时默认通过验证
            resolve(true);
          },
        });
      });
    });
  },

  /**
   * 登录行为定义
   * @type {Object}
   */
  loginAction: {
    id: 'login',
    type: 'function',
    selector: '#login_btn',
    /**
     * 处理登录请求
     * @param {string} id - 行为ID
     * @param {HTMLElement} target - 目标元素
     * @param {Object} eventData - 事件数据
     */
    handler: function (id, target, eventData) {
      const form = $(eventData.element).closest('form');

      // 验证表单，如果有错误则停止提交
      if (form.valid && form.valid('checkError')) {
        return;
      }

      // 发送AJAX请求 - 全局ajaxSetup会处理表单状态
      $.ajax({
        url: `${Auth.baseUrl}/login`,
        type: 'POST',
        data: form.serialize(),
        dataType: 'json',
        success: function (response) {
          window.location = `/u/${response.UserId}`;
        },
      });
    },

    /**
     * 初始化事件处理
     */
    init: function () {
      // 绑定回车键提交
      $(document).on('keypress', 'form input', function (e) {
        if (e.which === 13) {
          e.preventDefault();
          $('#login_btn').trigger('click');
        }
      });
    },
  },

  /**
   * 注册行为定义
   * @type {Object}
   */
  registerAction: {
    id: 'register',
    type: 'function',
    selector: '#register_btn',
    /**
     * 处理注册请求
     * @param {string} id - 行为ID
     * @param {HTMLElement} target - 目标元素
     * @param {Object} eventData - 事件数据
     */
    handler: function (id, target, eventData) {
      const form = $(eventData.element).closest('form');

      // 验证表单，如果有错误则停止提交
      if (form.valid && form.valid('checkError')) {
        return;
      }

      // 准备额外数据 - 时区和语言信息
      const timeZone = Intl.DateTimeFormat().resolvedOptions().timeZone;
      const lang = navigator.language;
      const extraData = `&time_zone=${encodeURIComponent(timeZone)}&lang=${encodeURIComponent(lang)}`;

      // 发送AJAX请求 - 全局ajaxSetup会处理表单状态
      $.ajax({
        url: `${Auth.baseUrl}/register`,
        type: 'POST',
        data: form.serialize() + extraData,
        dataType: 'json',
        success: function (response) {
          // 注册成功后跳转
          if (response.Data && response.Data.redirect) {
            window.location = response.Data.redirect;
          }
        },
      });
    },

    /**
     * 初始化事件处理
     */
    init: function () {
      // 绑定回车键提交
      $(document).on(
        'keypress',
        'form[name="register_form"] input',
        function (e) {
          if (e.which === 13) {
            e.preventDefault();
            $('#register_btn').trigger('click');
          }
        },
      );
    },
  },

  /**
   * 找回密码行为定义
   * @type {Object}
   */
  forgotPasswordAction: {
    id: 'forgot-password',
    type: 'function',
    selector: '#forgot_password_btn',
    /**
     * 处理找回密码请求
     * @param {string} id - 行为ID
     * @param {HTMLElement} target - 目标元素
     * @param {Object} eventData - 事件数据
     */
    handler: function (id, target, eventData) {
      const form = $(eventData.element).closest('form');

      // 验证表单，如果有错误则停止提交
      if (form.valid && form.valid('checkError')) {
        return;
      }

      // 发送AJAX请求 - 全局ajaxSetup会处理表单状态
      $.ajax({
        url: `${Auth.baseUrl}/forgot-password`,
        type: 'POST',
        data: form.serialize(),
        dataType: 'json',
        success: function (response) {
          kiko.ui.notify.success(appI18n.getMessage('operation.success'));
        },
      });
    },
  },
};
