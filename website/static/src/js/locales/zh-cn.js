export default {
  // Common
  'common.save': '保存',
  'common.edit': '编辑',
  'common.cancel': '取消',
  'common.delete': '删除',
  'common.op_fail': '操作失败',
  'common.op_success': '操作成功',
  'common.loading': '加载中...',
  'common.load_fail': '加载失败',
  'common.save_success': '保存成功',
  'common.save_fail': '保存失败',
  'common.net_error': '网络错误',
  'common.views': '浏览',
  'common.comments': '评论',

  // Diary
  'diary.create': '新建日记本',
  'diary.settings': '日记本设置',

  // Entry
  'entry.create': '写日记',
  'entry.edit': '编辑日记',
  'entry.delete.confirm': '确定要删除这篇日记吗？',
  'entry.drafts.confirm_hard_delete': '确定要彻底删除这篇草稿吗？此操作无法撤销。',
  'entry.delete.success': '删除成功',
  'entry.delete.fail': '删除失败',
  'entry.delete.net_error': '网络错误，删除失败',
  'entry.like.success': '点赞成功',
  'entry.like.cancel': '取消点赞',

  // Social
  'social.like.success': '点赞成功',
  'social.like.cancel': '取消点赞',

  // Auth (Frontend Specific)
  'auth.username_taken': '用户名已被占用',
  'auth.email_taken': '邮箱已被占用',
  'auth.success': '操作成功',

  // Comment
  'comment.submit_success': '留言提交成功',
  'comment.delete_confirm': '确定要删除这条留言吗？',
  'comment.delete_success': '留言删除成功',
  'comment.disabled': '评论已关闭',
  'comment.load_fail': '加载评论失败',

  // --- Backend Error Keys (Synced with internal/utils/errors/errors.go) ---

  // Validation (1000-1099)
  validation_error: '验证错误',
  invalid_params: '参数无效',
  missing_required_field: '缺少必填字段',
  invalid_format: '格式无效',
  value_out_of_range: '数值超出范围',
  invalid_password: '密码无效',
  invalid_email: '邮箱无效',
  invalid_username: '用户名无效',
  password_too_short: '密码太短',
  password_too_long: '密码太长',
  field_too_short: '字段太短',
  field_too_long: '字段太长',
  field_mismatch: '字段不匹配',

  // Auth (1100-1199)
  unauthorized: '未授权',
  invalid_token: 'Token无效',
  token_expired: 'Token已过期',
  login_required: '请先登录',
  invalid_credentials: '凭证无效',

  // Permission (1200-1299)
  permission_denied: '权限被拒绝',
  access_forbidden: '禁止访问',
  insufficient_privileges: '权限不足',
  operation_not_allowed: '操作不允许',

  // Resource (1300-1399)
  user_not_found: '用户不存在',
  resource_not_found: '资源不存在',
  diary_not_found: '日记本不存在',
  entry_not_found: '日记不存在',
  comment_not_found: '评论不存在',
  page_not_found: '页面不存在',

  // User Business (2000-2999)
  user_exists: '用户已存在',
  email_exists: '邮箱已存在',
  username_exists: '用户名已存在',
  duplicate_resource: '资源重复',

  // User Status (2100-2199)
  user_inactive: '用户未激活',
  user_banned: '用户已被封禁',
  account_locked: '账号已锁定',
  email_not_verified: '邮箱未验证',

  // Password (2200-2299)
  password_incorrect: '密码错误',
  password_too_weak: '密码太弱',
  password_expired: '密码已过期',
  old_password_incorrect: '旧密码错误',

  // Limits (2300-2399)
  diary_limit_exceeded: '日记本数量超出限制',
  content_too_long: '内容太长',
  file_too_large: '文件太大',
  quota_exceeded: '配额超限',

  // Communication (2400-2499)
  email_send_failed: '邮件发送失败',
  sms_send_failed: '短信发送失败',
  notification_failed: '通知发送失败',

  // Files (2500-2599)
  upload_failed: '上传失败',
  image_process_failed: '图片处理失败',
  file_format_unsupported: '不支持的文件格式',

  // Database (6000-6099)
  database_error: '数据库错误',
  connection_failed: '连接失败',
  query_timeout: '查询超时',
  transaction_failed: '事务失败',

  // Cache (6100-6199)
  redis_error: 'Redis错误',
  cache_error: '缓存错误',
  session_error: '会话错误',

  // Service (6200-6299)
  config_error: '配置错误',
  service_unavailable: '服务不可用',
  external_service_error: '外部服务错误',

  // DI (6250-6269)
  di_service_already_registered: 'DI服务已注册',
  di_service_not_registered: 'DI服务未注册',

  // System (6300-6399)
  rate_limit_exceeded: '请求过于频繁',
  server_overloaded: '服务器过载',
  maintenance_mode: '维护模式',

  // Internal (9999)
  internal_error: '内部服务器错误',
};
