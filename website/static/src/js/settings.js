/**
 * 设置模块 - 处理用户设置相关功能
 * 基于Kiko UI的页面行为管理组件实现
 * @version 1.0.0
 * @license MIT
 */

/**
 * 设置模块行为配置
 * @namespace Settings
 */
const Settings = {
  /**
   * 初始化设置模块
   * @param {string} baseUrl - 基础URL
   * @returns {Array} 设置模块行为配置数组
   */
  init: function (baseUrl = '') {
    log('设置模块初始化');
    this.baseUrl = baseUrl;
    return [
      this.saveProfileAction,
      this.saveEmailAction,
      this.changePasswordAction,
    ];
  },

  /**
   * 保存个人资料行为
   * @type {Object}
   */
  saveProfileAction: {
    id: 'save-profile',
    type: 'function',

    /**
     * 处理保存个人资料请求
     */
    handler: function (id, target, eventData) {
      const form = $('#profile-form');

      // 验证表单
      if (form.valid('checkError')) {
        return;
      }

      // 收集表单数据
      const data = {
        full_name: form.find('#full_name').val(),
        remark: form.find('#remark').val(),
      };

      // TODO: 实现服务端保存逻辑
      log('保存个人资料:', data);
      kiko.ui.notify.success(appI18n.getMessage('operation.success'));
    },
  },

  /**
   * 保存邮箱行为
   * @type {Object}
   */
  saveEmailAction: {
    id: 'save-email',
    type: 'function',

    /**
     * 处理保存邮箱请求
     */
    handler: function (id, target, eventData) {
      const form = $('#email-form');

      // 验证表单
      if (form.valid('checkError')) {
        return;
      }

      // 收集表单数据
      const data = {
        email: form.find('#email').val(),
      };

      // TODO: 实现服务端保存逻辑
      log('保存邮箱:', data);
      kiko.ui.notify.success(appI18n.getMessage('operation.success'));
    },
  },

  /**
   * 修改密码行为
   * @type {Object}
   */
  changePasswordAction: {
    id: 'change-password',
    type: 'function',

    /**
     * 处理修改密码请求
     */
    handler: function (id, target, eventData) {
      const form = $('#password-form');

      // 验证表单
      if (form.valid('checkError')) {
        return;
      }

      // 收集表单数据
      const data = {
        current_pwd: form.find('#current_pwd').val(),
        new_pwd: form.find('#new_pwd').val(),
        confirm_pwd: form.find('#confirm_pwd').val(),
      };

      // TODO: 实现服务端保存逻辑
      $.ajax({
        url: '/u/settings/password',
        method: 'POST',
        data: data,
        success: function (response) {
          log('修改密码:', response);
          kiko.ui.notify.success(appI18n.getMessage('operation.success'));
          // 清空密码字段
          form.find('#current_pwd, #new_pwd, #confirm_pwd').val('');
        },
        error: function (response) {
          kiko.ui.notify.error(appI18n.getMessage('operation.fail'));
        },
      });
    },
  },
};
