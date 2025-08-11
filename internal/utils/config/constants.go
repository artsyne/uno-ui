package config

// 系统操作状态码
const (
	OPS_SUCCESS int = 0 // 操作成功
)

// 用户相关常量
const (
	// 密码规则
	PWD_MIN_LEN int = 6  // 最小密码长度
	PWD_MAX_LEN int = 50 // 最大密码长度

	// 用户状态
	USER_STATUS_OK          int8 = 1 // 用户正常状态
	USER_STATUS_WAIT_ACTIVE int8 = 2 // 用户待激活
	USER_STATUS_DEL         int8 = 0 // 用户已删除状态

	// VIP状态
	USER_NOT_VIP int8 = 0 // 用户不是vip
	USER_IS_VIP  int8 = 1 // 用户是vip

	// 用户分组
	DEFAULT_USER_GROUP int8 = 1 // 用户注册默认分组
)

// 验证码相关常量
const (
	VALID_CODE_LEN int = 64 // 确认码长度

	// 验证码类型
	VALID_CODE_TYPE_REG       int8 = 1 // 注册验证码类型
	VALID_CODE_TYPE_RESET_PWD int8 = 2 // 重置密码验证码类型

	// 验证码过期时间
	VALID_CODE_EXPIRE_DAYS_REG       int = 1 // 注册验证码过期天数
	VALID_CODE_EXPIRE_DAYS_RESET_PWD int = 1 // 重置密码验证码过期天数

	// 旧常量（兼容性）
	VALID_CODE_ACTIVE_USER_EMAIL        int8 = 1 // 激活用户
	VALID_CODE_RESET_PWD_EMAIL          int8 = 2 // 重设密码
	VALID_CODE_INVITE_USER_JOIN_DAYBOOK int8 = 3 // 邀请用户加入日记本

	// 验证码过期时间
	REG_EMAIL_VAILID_CODE_EXPIRE_DAY      int = 1 // 注册email过期天数
	FIND_PWD_EMAIL_VAILID_CODE_EXPIRE_DAY int = 1 // 重设密码邮件过期天数
)

// 任务类型常量
const (
	JOB_TYPE_EMAIL int8 = 1 // email任务
	JOB_TYPE_DB    int8 = 2 // 数据库操作任务
	JOB_TYPE_OAUTH int8 = 3 // OAuth清理任务

	// 具体任务类型
	JOB_TYPE_SEND_ACTIVATE_EMAIL  int8 = 10 // 发送激活邮件
	JOB_TYPE_SEND_RESET_PWD_EMAIL int8 = 11 // 发送重置密码邮件
)

// 任务状态常量
const (
	JOB_RESULT_PENDING int8 = 0 // 任务待处理
	JOB_RESULT_RUNNING int8 = 1 // 任务处理中
	JOB_RESULT_SUCCESS int8 = 2 // 任务处理成功
	JOB_RESULT_FAILED  int8 = 3 // 任务处理失败
)

// 任务重试常量
const (
	MAX_RETRY_COUNT int8 = 3 // 最大重试次数
)

// 日记本权限常量
const (
	DAYBOOK_VIEW_PUBLIC  int8 = 1 // 公开
	DAYBOOK_VIEW_PRIVATE int8 = 2 // 私有
)

// 环境常量
const (
	ENV_DEVELOPMENT = "development"
	ENV_PRODUCTION  = "production"
	ENV_TEST        = "test"
)
