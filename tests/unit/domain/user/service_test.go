package user_test

import (
	"testing"

	"github.com/biuren/hebe/internal/domain/user"
	framework "github.com/biuren/hebe/tests/framework/core"
	
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	framework.RunTestsWithFlags(m)
}

// TestUserService_RegisterUser 测试用户注册功能
func TestUserService_RegisterUser(t *testing.T) {
	ctx := framework.NewTestContext("TestUserService_RegisterUser", framework.WithDatabase())
	defer ctx.Cleanup()

	// 创建真实的用户服务
	userService := user.NewService()

	// 准备测试数据 - 使用随机数据避免冲突
	newUser := &user.User{
		UserName: ctx.GetDataGenerator().GenerateUniqueString("user"),
		Email:    ctx.GetDataGenerator().GenerateUniqueEmail("user"),
		Pwd:      "password123",
		FullName: ctx.GetDataGenerator().GenerateUniqueString("fullname"),
		TimeZone: "Asia/Shanghai",
		Lang:     "zh-cn",
	}

	t.Run("成功注册用户", func(t *testing.T) {
		err := userService.RegisterUser(newUser)
		require.NoError(t, err)
		assert.True(t, newUser.UserId > 0)
		assert.NotEqual(t, "password123", newUser.Pwd) // 密码应该被加密
	})

	t.Run("用户名太短", func(t *testing.T) {
		invalidUser := &user.User{
			UserName: "ab",
			Email:    ctx.GetDataGenerator().GenerateUniqueEmail("user"),
			Pwd:      "password123",
			TimeZone: "Asia/Shanghai",
		}
		err := userService.RegisterUser(invalidUser)
		assert.Error(t, err)
	})

	t.Run("密码太短", func(t *testing.T) {
		invalidUser := &user.User{
			UserName: ctx.GetDataGenerator().GenerateUniqueString("user"),
			Email:    ctx.GetDataGenerator().GenerateUniqueEmail("user"),
			Pwd:      "123",
			TimeZone: "Asia/Shanghai",
		}
		err := userService.RegisterUser(invalidUser)
		assert.Error(t, err)
	})

	t.Run("无效时区", func(t *testing.T) {
		invalidUser := &user.User{
			UserName: ctx.GetDataGenerator().GenerateUniqueString("user"),
			Email:    ctx.GetDataGenerator().GenerateUniqueEmail("user"),
			Pwd:      "password123",
			TimeZone: "Invalid/TimeZone",
		}
		err := userService.RegisterUser(invalidUser)
		assert.Error(t, err)
	})

	t.Run("重复用户名", func(t *testing.T) {
		duplicateUser := &user.User{
			UserName: newUser.UserName, // 与第一个测试相同的用户名
			Email:    ctx.GetDataGenerator().GenerateUniqueEmail("user"),
			Pwd:      "password123",
			TimeZone: "Asia/Shanghai",
		}
		err := userService.RegisterUser(duplicateUser)
		assert.Error(t, err)
	})

	t.Run("重复邮箱", func(t *testing.T) {
		duplicateUser := &user.User{
			UserName: ctx.GetDataGenerator().GenerateUniqueString("user"),
			Email:    newUser.Email, // 与第一个测试相同的邮箱
			Pwd:      "password123",
			TimeZone: "Asia/Shanghai",
		}
		err := userService.RegisterUser(duplicateUser)
		assert.Error(t, err)
	})
}

// TestUserService_LoginWithUsernameOrEmail 测试用户登录验证
func TestUserService_LoginWithUsernameOrEmail(t *testing.T) {
	ctx := framework.NewTestContext("TestUserService_LoginWithUsernameOrEmail", framework.WithDatabase())
	defer ctx.Cleanup()
	
	// 创建真实的用户服务
	userService := user.NewService()

	// 先注册一个用户 - 使用随机数据
	newUser := &user.User{
		UserName: ctx.GetDataGenerator().GenerateUniqueString("user"),
		Email:    ctx.GetDataGenerator().GenerateUniqueEmail("user"),
		Pwd:      "loginpass123",
		FullName: ctx.GetDataGenerator().GenerateUniqueString("fullname"),
		TimeZone: "UTC",
		Lang:     "en",
	}

	err := userService.RegisterUser(newUser)
	require.NoError(t, err)
	require.True(t, newUser.UserId > 0)

	// 确保测试结束后清理数据
	defer func() {
		if newUser.UserId > 0 {
			// 清理测试数据
			userService.UpdateUserStatus(newUser.UserId, 0)
		}
	}()

	t.Run("用户名密码正确", func(t *testing.T) {
		user, err := userService.LoginWithUsernameOrEmail(newUser.UserName, "loginpass123")
		require.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, newUser.UserName, user.UserName)
		assert.Equal(t, newUser.Email, user.Email)
	})

	t.Run("邮箱密码正确", func(t *testing.T) {
		// 调试信息：打印注册的邮箱
		t.Logf("尝试使用邮箱登录: %s", newUser.Email)
		
		// 先尝试直接通过邮箱获取用户
		testUser, testErr := userService.GetUserByEmail(newUser.Email)
		if testErr != nil {
			t.Logf("直接通过邮箱获取用户失败: %v", testErr)
		} else {
			t.Logf("直接通过邮箱获取用户成功: %s", testUser.UserName)
		}
		
		user, err := userService.LoginWithUsernameOrEmail(newUser.Email, "loginpass123")
		require.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, newUser.UserName, user.UserName)
	})

	t.Run("密码错误", func(t *testing.T) {
		user, err := userService.LoginWithUsernameOrEmail(newUser.UserName, "wrongpassword")
		assert.Error(t, err)
		assert.Nil(t, user)
	})

	t.Run("用户不存在", func(t *testing.T) {
		user, err := userService.LoginWithUsernameOrEmail("nonexistent", "loginpass123")
		assert.Error(t, err)
		assert.Nil(t, user)
	})
}

// TestUserService_UpdateUserProfile 测试用户资料更新
func TestUserService_UpdateUserProfile(t *testing.T) {
	ctx := framework.NewTestContext("TestUserService_UpdateUserProfile", framework.WithDatabase())
	defer ctx.Cleanup()

	// 创建真实的用户服务
	userService := user.NewService()

	// 先注册一个用户 - 使用随机数据
	newUser := &user.User{
		UserName: ctx.GetDataGenerator().GenerateUniqueString("user"),
		Email:    ctx.GetDataGenerator().GenerateUniqueEmail("user"),
		Pwd:      "password123",
		FullName: "Original Name",
		Remark:   "Original Remark",
		TimeZone: "Asia/Shanghai",
	}

	err := userService.RegisterUser(newUser)
	require.NoError(t, err)
	require.True(t, newUser.UserId > 0)

	t.Run("成功更新资料", func(t *testing.T) {
		err := userService.UpdateUserProfile(newUser.UserId, "Updated Name", "Updated Remark")
		assert.NoError(t, err)

		// 验证更新是否成功
		updatedUser, err := userService.GetUserByUserId(newUser.UserId)
		require.NoError(t, err)
		assert.Equal(t, "Updated Name", updatedUser.FullName)
		assert.Equal(t, "Updated Remark", updatedUser.Remark)
	})

	t.Run("用户不存在", func(t *testing.T) {
		err := userService.UpdateUserProfile(99999, "Name", "Remark")
		assert.Error(t, err)
	})
}

// TestUserService_ChangeUserPassword 测试密码修改
func TestUserService_ChangeUserPassword(t *testing.T) {
	ctx := framework.NewTestContext("TestUserService_ChangeUserPassword", framework.WithDatabase())
	defer ctx.Cleanup()

	// 创建真实的用户服务
	userService := user.NewService()

	// 先注册一个用户 - 使用随机数据
	newUser := &user.User{
		UserName: ctx.GetDataGenerator().GenerateUniqueString("user"),
		Email:    ctx.GetDataGenerator().GenerateUniqueEmail("user"),
		Pwd:      "oldpassword123",
		TimeZone: "UTC",
	}

	err := userService.RegisterUser(newUser)
	require.NoError(t, err)
	require.True(t, newUser.UserId > 0)

	t.Run("成功修改密码", func(t *testing.T) {
		err := userService.ChangeUserPassword(newUser.UserId, "oldpassword123", "newpassword123")
		assert.NoError(t, err)

		// 验证新密码可以登录
		user, err := userService.LoginWithUsernameOrEmail(newUser.UserName, "newpassword123")
		require.NoError(t, err)
		assert.NotNil(t, user)

		// 验证旧密码不能登录
		user, err = userService.LoginWithUsernameOrEmail(newUser.UserName, "oldpassword123")
		assert.Error(t, err)
		assert.Nil(t, user)
	})

	t.Run("当前密码错误", func(t *testing.T) {
		err := userService.ChangeUserPassword(newUser.UserId, "wrongpassword", "newpassword456")
		assert.Error(t, err)
	})

	t.Run("新密码太短", func(t *testing.T) {
		err := userService.ChangeUserPassword(newUser.UserId, "newpassword123", "123")
		assert.Error(t, err)
	})

	t.Run("用户不存在", func(t *testing.T) {
		err := userService.ChangeUserPassword(99999, "oldpassword123", "newpassword123")
		assert.Error(t, err)
	})
}

// TestUserService_IsExistUserName 测试用户名存在检查
func TestUserService_IsExistUserName(t *testing.T) {
	ctx := framework.NewTestContext("TestUserService_IsExistUserName", framework.WithDatabase())
	defer ctx.Cleanup()

	// 创建真实的用户服务
	userService := user.NewService()

	// 先注册一个用户 - 使用随机数据
	newUser := &user.User{
		UserName: ctx.GetDataGenerator().GenerateUniqueString("user"),
		Email:    ctx.GetDataGenerator().GenerateUniqueEmail("user"),
		Pwd:      "password123",
		TimeZone: "Asia/Shanghai",
	}

	err := userService.RegisterUser(newUser)
	require.NoError(t, err)

	t.Run("已存在的用户名", func(t *testing.T) {
		isExist, err := userService.IsExistUserName(newUser.UserName)
		assert.NoError(t, err)
		assert.True(t, isExist)
	})

	t.Run("不存在的用户名", func(t *testing.T) {
		isExist, err := userService.IsExistUserName("nonexistentuser")
		assert.NoError(t, err)
		assert.False(t, isExist)
	})
}

// TestUserService_IsExistEmail 测试邮箱存在检查
func TestUserService_IsExistEmail(t *testing.T) {
	ctx := framework.NewTestContext("TestUserService_IsExistEmail", framework.WithDatabase())
	defer ctx.Cleanup()

	// 创建真实的用户服务
	userService := user.NewService()

	// 先注册一个用户 - 使用随机数据
	newUser := &user.User{
		UserName: ctx.GetDataGenerator().GenerateUniqueString("user"),
		Email:    ctx.GetDataGenerator().GenerateUniqueEmail("user"),
		Pwd:      "password123",
		TimeZone: "Asia/Shanghai",
	}

	err := userService.RegisterUser(newUser)
	require.NoError(t, err)

	t.Run("已存在的邮箱", func(t *testing.T) {
		isExist, err := userService.IsExistEmail(newUser.Email)
		assert.NoError(t, err)
		assert.True(t, isExist)
	})

	t.Run("不存在的邮箱", func(t *testing.T) {
		isExist, err := userService.IsExistEmail("nonexistent@example.com")
		assert.NoError(t, err)
		assert.False(t, isExist)
	})
}

// TestUserService_GetUserByUserId 测试根据ID获取用户
func TestUserService_GetUserByUserId(t *testing.T) {
	ctx := framework.NewTestContext("TestUserService_GetUserByUserId", framework.WithDatabase())
	defer ctx.Cleanup()

	// 创建真实的用户服务
	userService := user.NewService()

	// 先注册一个用户 - 使用随机数据
	newUser := &user.User{
		UserName: ctx.GetDataGenerator().GenerateUniqueString("user"),
		Email:    ctx.GetDataGenerator().GenerateUniqueEmail("user"),
		Pwd:      "loginpass123",
		FullName: ctx.GetDataGenerator().GenerateUniqueString("fullname"),
		TimeZone: "UTC",
	}

	err := userService.RegisterUser(newUser)
	require.NoError(t, err)
	require.True(t, newUser.UserId > 0)

	t.Run("成功获取用户", func(t *testing.T) {
		user, err := userService.GetUserByUserId(newUser.UserId)
		require.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, newUser.UserName, user.UserName)
		assert.Equal(t, newUser.Email, user.Email)
		assert.Equal(t, newUser.FullName, user.FullName)
	})

	t.Run("用户不存在", func(t *testing.T) {
		user, err := userService.GetUserByUserId(99999)
		assert.Error(t, err)
		assert.Nil(t, user)
	})
}

// TestUserService_GetUserByEmail 测试根据邮箱获取用户
func TestUserService_GetUserByEmail(t *testing.T) {
	ctx := framework.NewTestContext("TestUserService_GetUserByEmail", framework.WithDatabase())
	defer ctx.Cleanup()

	// 创建真实的用户服务
	userService := user.NewService()

	// 先注册一个用户 - 使用随机数据
	newUser := &user.User{
		UserName: ctx.GetDataGenerator().GenerateUniqueString("user"),
		Email:    ctx.GetDataGenerator().GenerateUniqueEmail("user"),
		Pwd:      "password123",
		FullName: ctx.GetDataGenerator().GenerateUniqueString("fullname"),
		TimeZone: "Asia/Shanghai",
	}

	err := userService.RegisterUser(newUser)
	require.NoError(t, err)

	t.Run("成功获取用户", func(t *testing.T) {
		user, err := userService.GetUserByEmail(newUser.Email)
		require.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, newUser.UserName, user.UserName)
		assert.Equal(t, newUser.Email, user.Email)
	})

	t.Run("邮箱不存在", func(t *testing.T) {
		user, err := userService.GetUserByEmail("nonexistent@example.com")
		assert.Error(t, err)
		assert.Nil(t, user)
	})
}

// TestUserService_UpdateUserStatus 测试用户状态更新
func TestUserService_UpdateUserStatus(t *testing.T) {
	ctx := framework.NewTestContext("TestUserService_UpdateUserStatus", framework.WithDatabase())
	defer ctx.Cleanup()

	// 创建真实的用户服务
	userService := user.NewService()

	// 先注册一个用户 - 使用随机数据
	newUser := &user.User{
		UserName: ctx.GetDataGenerator().GenerateUniqueString("user"),
		Email:    ctx.GetDataGenerator().GenerateUniqueEmail("user"),
		Pwd:      "password123",
		TimeZone: "UTC",
	}

	err := userService.RegisterUser(newUser)
	require.NoError(t, err)

	t.Run("成功更新邮箱", func(t *testing.T) {
		newEmail := ctx.GetDataGenerator().GenerateUniqueEmail("user")
		err := userService.UpdateUserEmail(newUser.UserId, newEmail)
		assert.NoError(t, err)

		// 验证更新是否成功
		user, err := userService.GetUserByUserId(newUser.UserId)
		require.NoError(t, err)
		assert.Equal(t, newEmail, user.Email)
	})

	t.Run("用户不存在", func(t *testing.T) {
		err := userService.UpdateUserEmail(99999, ctx.GetDataGenerator().GenerateUniqueEmail("user"))
		assert.Error(t, err)
	})
}

// TestUserService_GetSupportedLang 测试语言支持检查
func TestUserService_GetSupportedLang(t *testing.T) {
	ctx := framework.NewTestContext("TestUserService_GetSupportedLang", framework.WithDatabase())
	defer ctx.Cleanup()
	
	userService := user.NewService()

	tests := []struct {
			name     string
			input    string
			expected string
		}{
			{"支持的中文简体", "zh-cn", "zh-cn"},
			{"支持的中文繁体", "zh-tw", "zh-tw"},
			{"支持的英文", "en", "en"},
			{"不支持的语言", "fr", "en"}, // 默认返回中文
			{"空字符串", "", "en"},     // 默认返回中文
		}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := userService.GetSupportedLang(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
