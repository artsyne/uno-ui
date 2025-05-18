/**
 * Kiko UI库全局初始化
 *
 * 此文件必须在所有其他应用文件之前打包
 * 它不定义kiko命名空间或方法，而是在核心框架加载后提供应用级别的初始化
 *
 * @version 1.2.0
 * @license MIT
 */

// 在DOM就绪后初始化Kiko UI
$(function () {
  // 直接调用kiko.init初始化UI框架
  kiko.init({
    debug: true,
    animation: {
      enable: true,
      duration: 300,
    },
    errorHandler: function (type, error, data) {
      log('Kiko错误', type, error, 'error');
    },
  });

  log('Kiko UI框架初始化完成');
});
