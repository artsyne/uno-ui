package user_test

import (
	"testing"

	"github.com/biuren/hebe/internal/domain/task"
	"github.com/biuren/hebe/internal/domain/user"
	"github.com/biuren/hebe/internal/infra/repository"
	"github.com/biuren/hebe/internal/utils/i18n"
	framework "github.com/biuren/hebe/tests/framework"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	framework.RunTestsWithFlags(m)
}

// TestUserService_RegisterUser 测试用户注册功能
func TestUserService_RegisterUser(t *testing.T) {
	ctx := framework.NewTestContext(framework.WithDatabase())
	defer ctx.Cleanup()

	// 创建真实的用户服务
	userRepo := repository.NewUserRepository()
	taskRepo := repository.NewTaskRepository()
	config := ctx.GetConfig().GetAppConfig() // 通过测试上下文获取配置
	taskService := task.NewService(config, taskRepo)
	i18nManager := i18n.NewLocalizationManager(config)
	userService := user.NewService(config, userRepo, taskService, i18nManager, nil)

	// 准备测试数据 - 使用随机数据避免冲突
	newUser := &user.User{
		Handle: ctx.GetDataGenerator().GenerateUniqueString("user"),
		Email:    ctx.GetDataGenerator().GenerateUniqueEmail("user"),
		Pwd:      "password123",
		DisplayName: ctx.GetDataGenerator().GenerateUniqueString("fullname"),
		TimeZone: "Asia/Shanghai",
		Lang:     "zh-cn",
	}

	t.Run("成功注册用户", func(t *testing.T) {
		err := userService.RegisterUser(newUser)
		require.NoError(t, err)
		assert.True(t, newUser.UserId > 0)
		assert.NotEqual(t, "password123", newUser.Pwd) // 密码应该被加密
		// Check ModifiedUTCTime initialized
		assert.Equal(t, newUser.CreatedUTCTime, newUser.ModifiedUTCTime)
		assert.True(t, newUser.ModifiedUTCTime > 0)
	})

	t.Run("用户名太短", func(t *testing.T) {
		invalidUser := &user.User{
			Handle: "ab",
			Email:    ctx.GetDataGenerator().GenerateUniqueEmail("user"),
			Pwd:      "password123",
			TimeZone: "Asia/Shanghai",
		}
		err := userService.RegisterUser(invalidUser)
		assert.Error(t, err)
	})

	t.Run("密码太短", func(t *testing.T) {
		invalidUser := &user.User{
			Handle: ctx.GetDataGenerator().GenerateUniqueString("user"),
			Email:    ctx.GetDataGenerator().GenerateUniqueEmail("user"),
			Pwd:      "123",
			TimeZone: "Asia/Shanghai",
		}
		err := userService.RegisterUser(invalidUser)
		assert.Error(t, err)
	})

	t.Run("无效时区", func(t *testing.T) {
		invalidUser := &user.User{
			Handle: ctx.GetDataGenerator().GenerateUniqueString("user"),
			Email:    ctx.GetDataGenerator().GenerateUniqueEmail("user"),
			Pwd:      "password123",
			TimeZone: "Invalid/TimeZone",
		}
		err := userService.RegisterUser(invalidUser)
		assert.Error(t, err)
	})

	t.Run("重复用户名", func(t *testing.T) {
		duplicateUser := &user.User{
			Handle: newUser.Handle, // 与第一个测试相同的用户名
			Email:    ctx.GetDataGenerator().GenerateUniqueEmail("user"),
			Pwd:      "password123",
			TimeZone: "Asia/Shanghai",
		}
		err := userService.RegisterUser(duplicateUser)
		assert.Error(t, err)
	})

	t.Run("重复邮箱", func(t *testing.T) {
		duplicateUser := &user.User{
			Handle: ctx.GetDataGenerator().GenerateUniqueString("user"),
			Email:    newUser.Email, // 与第一个测试相同的邮箱
			Pwd:      "password123",
			TimeZone: "Asia/Shanghai",
		}
		err := userService.RegisterUser(duplicateUser)
		assert.Error(t, err)
	})
}

// TestUserService_UpdateUserHandle 测试用户Handle更新功能
func TestUserService_UpdateUserHandle(t *testing.T) {
	ctx := framework.NewTestContext(framework.WithDatabase())
	defer ctx.Cleanup()

	// Setup service
	userRepo := repository.NewUserRepository()
	taskRepo := repository.NewTaskRepository()
	config := ctx.GetConfig().GetAppConfig()
	
	taskService := task.NewService(config, taskRepo)
	i18nManager := i18n.NewLocalizationManager(config)
	userService := user.NewService(config, userRepo, taskService, i18nManager, nil)

	// Create user
	newUser := &user.User{
		Handle:      ctx.GetDataGenerator().GenerateUniqueString("user"),
		Email:       ctx.GetDataGenerator().GenerateUniqueEmail("user"),
		Pwd:         "password123",
		DisplayName: "Test User",
		TimeZone:    "Asia/Shanghai",
		Lang:        "zh-cn",
	}
	err := userService.RegisterUser(newUser)
	require.NoError(t, err)

	// Test 1: Cooldown active (should fail)
	// Use short handle to satisfy length limit (prefix + _ + h2 <= 16)
	newHandle := ctx.GetDataGenerator().GenerateUniqueString("h2")
	err = userService.UpdateUserHandle(newUser.UserId, newHandle)
	assert.Error(t, err)

	// Test 2: Cooldown passed (should success)
	// Manually update ModifiedUTCTime in DB to 31 days ago
	pastTime := newUser.CreatedUTCTime - 31*24*3600*1000
	_, err = ctx.GetDatabase().Execute("UPDATE br_user SET modified_utc_time = ? WHERE user_id = ?", pastTime, newUser.UserId)
	require.NoError(t, err)

	err = userService.UpdateUserHandle(newUser.UserId, newHandle)
	require.NoError(t, err)

	// Verify update
	updatedUser, err := userService.GetUserByUserId(newUser.UserId)
	require.NoError(t, err)
	assert.Equal(t, newHandle, updatedUser.Handle)
	assert.True(t, updatedUser.ModifiedUTCTime > pastTime)
}
