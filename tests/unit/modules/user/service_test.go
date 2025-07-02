package user_test

import (
	"testing"

	"github.com/biuren/hebe/internal/module/user"
	testhelpers "github.com/biuren/hebe/tests/helpers"
)

func TestUserService_RegisterUser(t *testing.T) {
	// 初始化测试环境
	dbHelper := testhelpers.SetupTestDatabase(t)
	defer dbHelper.Cleanup(t)

	// 开始事务
	dbHelper.BeginTransaction(t)
	defer dbHelper.RollbackTransaction(t)

	service := user.NewService()

	tests := []struct {
		name    string
		user    *user.User
		wantErr bool
		errType string
	}{
		{
			name: "正常注册",
			user: &user.User{
				UserName: "testuser",
				Email:    "test@example.com",
				Pwd:      "password123",
				TimeZone: "Asia/Shanghai",
				Lang:     "zh-cn",
				LastIp:   "127.0.0.1",
			},
			wantErr: false,
		},
		{
			name: "用户名太短",
			user: &user.User{
				UserName: "ab", // 太短
				Email:    "test2@example.com",
				Pwd:      "password123",
				TimeZone: "Asia/Shanghai",
				Lang:     "zh-cn",
				LastIp:   "127.0.0.1",
			},
			wantErr: true,
			errType: "输入参数错误",
		},
		{
			name: "密码太短",
			user: &user.User{
				UserName: "testuser2",
				Email:    "test3@example.com",
				Pwd:      "123", // 太短
				TimeZone: "Asia/Shanghai",
				Lang:     "zh-cn",
				LastIp:   "127.0.0.1",
			},
			wantErr: true,
			errType: "输入参数错误",
		},
		{
			name: "无效时区",
			user: &user.User{
				UserName: "testuser3",
				Email:    "test4@example.com",
				Pwd:      "password123",
				TimeZone: "Invalid/TimeZone", // 无效时区
				Lang:     "zh-cn",
				LastIp:   "127.0.0.1",
			},
			wantErr: true,
			errType: "输入参数错误",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.RegisterUser(tt.user)

			if tt.wantErr {
				if err == nil {
					t.Errorf("期望有错误，但没有错误")
				}
				return
			}

			if err != nil {
				t.Errorf("不期望有错误，但有错误: %v", err)
				return
			}

			// 验证用户ID已生成
			if tt.user.UserId <= 0 {
				t.Error("用户ID应该大于0")
			}

			// 验证密码已加密（不应该是明文）
			if tt.user.Pwd == "password123" {
				t.Error("密码应该已加密，不应该是明文")
			}

			// 验证盐值已生成
			if tt.user.PwdSalt == "" {
				t.Error("密码盐值应该已生成")
			}
		})
	}
}

func TestUserService_CheckNameAndPassword(t *testing.T) {
	// 初始化测试环境
	dbHelper := testhelpers.SetupTestDatabase(t)
	defer dbHelper.Cleanup(t)

	// 开始事务
	dbHelper.BeginTransaction(t)
	defer dbHelper.RollbackTransaction(t)

	service := user.NewService()

	// 先注册一个用户
	testUser := &user.User{
		UserName: "testuser",
		Email:    "test@example.com",
		Pwd:      "password123",
		TimeZone: "Asia/Shanghai",
		Lang:     "zh-cn",
		LastIp:   "127.0.0.1",
	}

	err := service.RegisterUser(testUser)
	if err != nil {
		t.Fatalf("注册用户失败: %v", err)
	}

	tests := []struct {
		name     string
		username string
		password string
		wantErr  bool
	}{
		{
			name:     "用户名登录成功",
			username: "testuser",
			password: "password123",
			wantErr:  false,
		},
		{
			name:     "邮箱登录成功",
			username: "test@example.com",
			password: "password123",
			wantErr:  false,
		},
		{
			name:     "密码错误",
			username: "testuser",
			password: "wrongpassword",
			wantErr:  true,
		},
		{
			name:     "用户不存在",
			username: "nonexistent",
			password: "password123",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.CheckNameAndPassword(tt.username, tt.password)

			if tt.wantErr {
				if err == nil {
					t.Errorf("期望有错误，但没有错误")
				}
				return
			}

			if err != nil {
				t.Errorf("不期望有错误，但有错误: %v", err)
				return
			}

			if result == nil {
				t.Error("结果不应为空")
				return
			}

			// 验证返回的用户信息
			if result.UserId != testUser.UserId {
				t.Errorf("用户ID不匹配: 期望 %d, 实际 %d", testUser.UserId, result.UserId)
			}

			if result.UserName != testUser.UserName {
				t.Errorf("用户名不匹配: 期望 '%s', 实际 '%s'", testUser.UserName, result.UserName)
			}
		})
	}
}

func TestUserService_UpdateUserProfile(t *testing.T) {
	// 初始化测试环境
	dbHelper := testhelpers.SetupTestDatabase(t)
	defer dbHelper.Cleanup(t)

	// 开始事务
	dbHelper.BeginTransaction(t)
	defer dbHelper.RollbackTransaction(t)

	service := user.NewService()

	// 先注册一个用户
	testUser := &user.User{
		UserName: "testuser",
		Email:    "test@example.com",
		Pwd:      "password123",
		TimeZone: "Asia/Shanghai",
		Lang:     "zh-cn",
		LastIp:   "127.0.0.1",
	}

	err := service.RegisterUser(testUser)
	if err != nil {
		t.Fatalf("注册用户失败: %v", err)
	}

	tests := []struct {
		name     string
		uid      int64
		fullName string
		remark   string
		wantErr  bool
	}{
		{
			name:     "更新全名和备注",
			uid:      testUser.UserId,
			fullName: "测试用户",
			remark:   "这是测试用户",
			wantErr:  false,
		},
		{
			name:     "只更新全名",
			uid:      testUser.UserId,
			fullName: "新的全名",
			remark:   "",
			wantErr:  false,
		},
		{
			name:     "用户不存在",
			uid:      999999, // 不存在的用户ID
			fullName: "测试",
			remark:   "测试",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.UpdateUserProfile(tt.uid, tt.uid, tt.fullName, tt.remark)

			if tt.wantErr {
				if err == nil {
					t.Errorf("期望有错误，但没有错误")
				}
				return
			}

			if err != nil {
				t.Errorf("不期望有错误，但有错误: %v", err)
				return
			}

			// 验证更新是否成功
			updatedUser, err := service.GetUserByUserId(tt.uid)
			if err != nil {
				t.Errorf("获取更新后的用户失败: %v", err)
				return
			}

			if tt.fullName != "" && updatedUser.FullName != tt.fullName {
				t.Errorf("全名未更新: 期望 '%s', 实际 '%s'", tt.fullName, updatedUser.FullName)
			}

			if tt.remark != "" && updatedUser.Remark != tt.remark {
				t.Errorf("备注未更新: 期望 '%s', 实际 '%s'", tt.remark, updatedUser.Remark)
			}
		})
	}
}

func TestUserService_ChangeUserPassword(t *testing.T) {
	// 初始化测试环境
	dbHelper := testhelpers.SetupTestDatabase(t)
	defer dbHelper.Cleanup(t)

	// 开始事务
	dbHelper.BeginTransaction(t)
	defer dbHelper.RollbackTransaction(t)

	service := user.NewService()

	// 先注册一个用户
	testUser := &user.User{
		UserName: "testuser",
		Email:    "test@example.com",
		Pwd:      "oldpassword123",
		TimeZone: "Asia/Shanghai",
		Lang:     "zh-cn",
		LastIp:   "127.0.0.1",
	}

	err := service.RegisterUser(testUser)
	if err != nil {
		t.Fatalf("注册用户失败: %v", err)
	}

	tests := []struct {
		name       string
		uid        int64
		currentPwd string
		newPwd     string
		wantErr    bool
	}{
		{
			name:       "正常修改密码",
			uid:        testUser.UserId,
			currentPwd: "oldpassword123",
			newPwd:     "newpassword123",
			wantErr:    false,
		},
		{
			name:       "当前密码错误",
			uid:        testUser.UserId,
			currentPwd: "wrongpassword",
			newPwd:     "newpassword123",
			wantErr:    true,
		},
		{
			name:       "新密码太短",
			uid:        testUser.UserId,
			currentPwd: "oldpassword123",
			newPwd:     "123", // 太短
			wantErr:    true,
		},
		{
			name:       "用户不存在",
			uid:        999999, // 不存在的用户ID
			currentPwd: "oldpassword123",
			newPwd:     "newpassword123",
			wantErr:    true,
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.ChangeUserPassword(tt.uid, tt.uid, tt.currentPwd, tt.newPwd)

			if tt.wantErr {
				if err == nil {
					t.Errorf("期望有错误，但没有错误")
				}
				return
			}

			if err != nil {
				t.Errorf("不期望有错误，但有错误: %v", err)
				return
			}

			// 只有第一个测试用例（正常修改密码）需要验证密码是否已更改
			if i == 0 {
				// 验证新密码是否可以登录
				_, loginErr := service.CheckNameAndPassword(testUser.UserName, tt.newPwd)
				if loginErr != nil {
					t.Errorf("使用新密码登录失败: %v", loginErr)
				}

				// 验证旧密码不能登录
				_, oldLoginErr := service.CheckNameAndPassword(testUser.UserName, tt.currentPwd)
				if oldLoginErr == nil {
					t.Error("旧密码仍然可以登录，密码修改失败")
				}
			}
		})
	}
}

func TestUserService_IsExistUserName(t *testing.T) {
	// 初始化测试环境
	dbHelper := testhelpers.SetupTestDatabase(t)
	defer dbHelper.Cleanup(t)

	// 开始事务
	dbHelper.BeginTransaction(t)
	defer dbHelper.RollbackTransaction(t)

	service := user.NewService()

	// 先注册一个用户
	testUser := &user.User{
		UserName: "existinguser",
		Email:    "existing@example.com",
		Pwd:      "password123",
		TimeZone: "Asia/Shanghai",
		Lang:     "zh-cn",
		LastIp:   "127.0.0.1",
	}

	err := service.RegisterUser(testUser)
	if err != nil {
		t.Fatalf("注册用户失败: %v", err)
	}

	tests := []struct {
		name     string
		username string
		wantErr  bool
	}{
		{
			name:     "用户名已存在",
			username: "existinguser",
			wantErr:  true,
		},
		{
			name:     "用户名不存在",
			username: "newuser",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.IsExistUserName(tt.username)

			if tt.wantErr {
				if err == nil {
					t.Errorf("期望有错误（用户名已存在），但没有错误")
				}
				return
			}

			if err != nil {
				t.Errorf("不期望有错误，但有错误: %v", err)
			}
		})
	}
}

func TestUserService_IsExistEmail(t *testing.T) {
	// 初始化测试环境
	dbHelper := testhelpers.SetupTestDatabase(t)
	defer dbHelper.Cleanup(t)

	// 开始事务
	dbHelper.BeginTransaction(t)
	defer dbHelper.RollbackTransaction(t)

	service := user.NewService()

	// 先注册一个用户
	testUser := &user.User{
		UserName: "testuser",
		Email:    "existing@example.com",
		Pwd:      "password123",
		TimeZone: "Asia/Shanghai",
		Lang:     "zh-cn",
		LastIp:   "127.0.0.1",
	}

	err := service.RegisterUser(testUser)
	if err != nil {
		t.Fatalf("注册用户失败: %v", err)
	}

	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{
			name:    "邮箱已存在",
			email:   "existing@example.com",
			wantErr: true,
		},
		{
			name:    "邮箱不存在",
			email:   "new@example.com",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.IsExistEmail(tt.email)

			if tt.wantErr {
				if err == nil {
					t.Errorf("期望有错误（邮箱已存在），但没有错误")
				}
				return
			}

			if err != nil {
				t.Errorf("不期望有错误，但有错误: %v", err)
			}
		})
	}
}
