/**
 * 用户状态管理模块
 * 提供统一的用户登录状态检查和用户信息获取接口
 */
const UserState = {
  // 私有变量存储用户信息
  _userInfo: null,
  _initialized: false,

  /**
   * 初始化用户状态
   * 从服务器端模板注入的数据中获取用户信息
   */
  init: function(userInfo = null) {
    if (userInfo) {
      this._userInfo = userInfo;
    } else {
      // 尝试从 meta 标签获取
      const userId = document.querySelector('meta[name="user-id"]')?.content;
      const userName = document.querySelector('meta[name="user-name"]')?.content;
      
      if (userId) {
        this._userInfo = {
          id: userId,
          name: userName || ''
        };
      }
    }
    this._initialized = true;
  },

  /**
   * 检查用户是否已登录
   * @returns {boolean} 用户是否已登录
   */
  isLoggedIn: function() {
    if (!this._initialized) {
      this.init();
    }
    return this._userInfo !== null && this._userInfo.id && this._userInfo.id !== '';
  },

  /**
   * 获取用户ID
   * @returns {string|null} 用户ID，未登录时返回null
   */
  getUserId: function() {
    if (!this.isLoggedIn()) {
      return null;
    }
    return this._userInfo.id;
  },

  /**
   * 获取用户名
   * @returns {string|null} 用户名，未登录时返回null
   */
  getUserName: function() {
    if (!this.isLoggedIn()) {
      return null;
    }
    return this._userInfo.name || '';
  },

  /**
   * 获取完整的用户信息
   * @returns {object|null} 用户信息对象，未登录时返回null
   */
  getUserInfo: function() {
    if (!this.isLoggedIn()) {
      return null;
    }
    return { ...this._userInfo };
  },

  /**
   * 设置用户信息（用于登录后更新状态）
   * @param {object} userInfo 用户信息对象
   */
  setUserInfo: function(userInfo) {
    this._userInfo = userInfo;
    this._initialized = true;
  },

  /**
   * 清除用户状态（用于登出）
   */
  clearUserInfo: function() {
    this._userInfo = null;
    this._initialized = true;
  }
};