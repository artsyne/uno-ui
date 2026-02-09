export default {
  // Common
  'common.save': 'Save',
  'common.edit': 'Edit',
  'common.cancel': 'Cancel',
  'common.confirm': 'Confirm',
  'common.delete': 'Delete',
  'common.op_fail': 'Operation failed',
  'common.op_success': 'Operation successful',
  'common.loading': 'Loading...',
  'common.load_fail': 'Load failed',
  'common.save_success': 'Saved successfully',
  'common.save_fail': 'Save failed',
  'common.net_error': 'Network error',
  'common.views': 'Views',
  'common.comments': 'Comments',

  // Diary
  'diary.create': 'New Diary',
  'diary.settings': 'Diary Settings',
  'diary.view': 'View Diary',

  // Diary Author
  'diary.author.invite': 'Invite Author',
  'diary.author.invite_desc': 'Invite friends to join this diary',
  'diary.author.identity_label': 'Username or Email',
  'diary.author.identity_placeholder': 'Enter username or email',
  'diary.author.role': 'Role',
  'diary.author.admin': 'Admin',
  'diary.author.member': 'Member',
  'diary.author.owner': 'Owner',
  'diary.author.status_active': 'Active',
  'diary.author.status_pending': 'Pending',
  'diary.author.invite_success': 'Invitation sent successfully',
  'diary.author.remove_confirm':
    'Are you sure you want to remove author {name}?',
  'diary.author.remove_title': 'Remove Author',
  'diary.author.remove_author_title': 'Remove Author',
  'diary.author.confirm_remove_author':
    'Are you sure you want to remove this author?',
  'diary.author.quit_confirm': 'Are you sure you want to quit this diary?',
  'diary.author.quit_title': 'Quit Diary',
  'diary.author.remove_success': 'Removed successfully',
  'diary.author.quit_success': 'Quit successfully',
  'diary.author.update_success': 'Permission updated successfully',
  'diary.author.list_title': 'Authors',
  'diary.author.manage': 'Manage Authors',
  'diary.author.confirm_grant_admin':
    'Are you sure you want to grant admin access to this member?',
  'diary.author.confirm_revoke_admin':
    'Are you sure you want to revoke admin access from this member?',
  'diary.author.admin_toggle_title': 'Change Permission',

  // Entry
  'entry.create': 'New Entry',
  'entry.edit': 'Edit Entry',
  'entry.delete.confirm': 'Are you sure you want to delete this entry?',
  'entry.drafts.confirm_hard_delete':
    'Are you sure you want to permanently delete this draft? This action cannot be undone.',
  'entry.delete.success': 'Deleted successfully',
  'entry.delete.fail': 'Delete failed',
  'entry.delete.net_error': 'Network error, delete failed',
  'entry.like.success': 'Liked successfully',
  'entry.like.cancel': 'Unliked successfully',

  // Social
  'social.like.success': 'Liked successfully',
  'social.like.cancel': 'Unliked successfully',

  // Auth (Frontend Specific)
  'auth.username_taken': 'Username already taken',
  'auth.email_taken': 'Email already taken',
  'auth.success': 'Success',

  // Comment
  'comment.submit_success': 'Comment submitted successfully',
  'comment.delete_confirm': 'Are you sure you want to delete this comment?',
  'comment.delete_success': 'Comment deleted successfully',
  'comment.disabled': 'Comments disabled',
  'comment.load_fail': 'Failed to load comments',

  // --- Backend Error Keys (Synced with internal/utils/errors/errors.go) ---

  // Validation (1000-1099)
  validation_error: 'Validation Error',
  invalid_params: 'Invalid Parameters',
  missing_required_field: 'Missing Required Field',
  invalid_format: 'Invalid Format',
  value_out_of_range: 'Value Out of Range',
  invalid_password: 'Invalid Password',
  invalid_email: 'Invalid Email',
  invalid_username: 'Invalid Username',
  password_too_short: 'Password Too Short',
  password_too_long: 'Password Too Long',
  field_too_short: 'Field Too Short',
  field_too_long: 'Field Too Long',
  field_mismatch: 'Field Mismatch',

  // Auth (1100-1199)
  unauthorized: 'Unauthorized',
  invalid_token: 'Invalid Token',
  token_expired: 'Token Expired',
  login_required: 'Login Required',
  invalid_credentials: 'Invalid Credentials',

  // Permission (1200-1299)
  permission_denied: 'Permission Denied',
  access_forbidden: 'Access Forbidden',
  insufficient_privileges: 'Insufficient Privileges',
  operation_not_allowed: 'Operation Not Allowed',

  // Resource (1300-1399)
  user_not_found: 'User Not Found',
  resource_not_found: 'Resource Not Found',
  diary_not_found: 'Diary Not Found',
  entry_not_found: 'Entry Not Found',
  comment_not_found: 'Comment Not Found',
  page_not_found: 'Page Not Found',

  // User Business (2000-2999)
  user_exists: 'User Already Exists',
  email_exists: 'Email Already Exists',
  username_exists: 'Username Already Exists',
  duplicate_resource: 'Duplicate Resource',

  // User Status (2100-2199)
  user_inactive: 'User Inactive',
  user_banned: 'User Banned',
  account_locked: 'Account Locked',
  email_not_verified: 'Email Not Verified',

  // Password (2200-2299)
  password_incorrect: 'Incorrect Password',
  password_too_weak: 'Password Too Weak',
  password_expired: 'Password Expired',
  old_password_incorrect: 'Old Password Incorrect',

  // Limits (2300-2399)
  diary_limit_exceeded: 'Diary Limit Exceeded',
  content_too_long: 'Content Too Long',
  file_too_large: 'File Too Large',
  quota_exceeded: 'Quota Exceeded',

  // Communication (2400-2499)
  email_send_failed: 'Email Send Failed',
  sms_send_failed: 'SMS Send Failed',
  notification_failed: 'Notification Failed',

  // Files (2500-2599)
  upload_failed: 'Upload Failed',
  image_process_failed: 'Image Process Failed',
  file_format_unsupported: 'File Format Unsupported',

  // Database (6000-6099)
  database_error: 'Database Error',
  connection_failed: 'Connection Failed',
  query_timeout: 'Query Timeout',
  transaction_failed: 'Transaction Failed',

  // Cache (6100-6199)
  redis_error: 'Redis Error',
  cache_error: 'Cache Error',
  session_error: 'Session Error',

  // Service (6200-6299)
  config_error: 'Configuration Error',
  service_unavailable: 'Service Unavailable',
  external_service_error: 'External Service Error',

  // DI (6250-6269)
  di_service_already_registered: 'DI Service Already Registered',
  di_service_not_registered: 'DI Service Not Registered',

  // System (6300-6399)
  rate_limit_exceeded: 'Rate Limit Exceeded',
  server_overloaded: 'Server Overloaded',
  maintenance_mode: 'Maintenance Mode',

  // Internal (9999)
  internal_error: 'Internal Server Error',
};
