/**
 * Kiko UI库全局初始化
 *
 * 此文件必须在所有其他应用文件之前打包
 *
 * @version 1.2.0
 * @license MIT
 */

// 创建kiko命名空间
window.kiko = window.kiko || {};
window.kiko.ui = window.kiko.ui || {};

/**
 * Kiko默认配置
 * 将配置与逻辑分离，便于维护和扩展
 */
kiko.defaultConfig = {
  debug: true, // 开启调试模式
  logLevel: 'info', // 日志级别
  animation: {
    enable: true, // 启用动画
    duration: 300, // 动画持续时间
  },
  ui: {
    notify: {
      position: 'top',
      messageContainer: '#ajax-message',
      loadingIndicator: '#global-loading-indicator',
    },
  },
  errorHandler: function (type, error, data) {
    // 记录错误
    console.error('[Kiko Error]', type, error);
  },
};

/**
 * 记录日志的工具函数
 * @param {string} message - 日志消息
 * @param {...any} data - 附加数据
 */
kiko.log = function (message, ...data) {
  if (this.config && this.config.debug) {
    console.log(`[Kiko] ${message}`, ...data);
  }
};

/**
 * 初始化Kiko框架
 * @param {Object} [options] - 自定义配置选项
 * @returns {Object} kiko实例
 */
kiko.init = function (options) {
  // 合并配置选项
  this.config = Object.assign({}, this.defaultConfig, options || {});

  kiko.log('配置已初始化', this.config);
  return this;
};

/**
 * 初始化UI组件
 * 负责初始化所有UI相关组件
 */
kiko.initUI = function () {
  kiko.log('初始化UI组件...');

  // 初始化notify组件
  this._initNotify();

  kiko.log('UI组件初始化完成');
  return this;
};

/**
 * 初始化通知组件
 * 私有方法，由initUI调用
 */
kiko._initNotify = function () {
  if (typeof kiko.ui.notify?.init === 'function') {
    kiko.ui.notify.init(this.config.ui.notify);
    kiko.log('Notify组件初始化完成');
  } else {
    console.warn('[Kiko] Notify组件未找到，无法初始化');
  }
};

// 确保在DOM就绪后初始化
$(function () {
  kiko.log('Kiko UI库初始化开始...');

  // 初始化框架并设置组件
  kiko.init().initUI();

  kiko.log('Kiko UI库初始化完成');
});
