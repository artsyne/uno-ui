package user_test

import (
	"testing"

	"github.com/biuren/hebe/internal/domain/user"
	framework "github.com/biuren/hebe/tests/framework"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestUserService_Integration 使用真实数据库的集成测试示例
// 这种测试方式可以验证真实的数据库交互和业务逻辑
func TestUserService_Integration_RegisterAndLogin(t *testing.T) {
	ctx := framework.NewTestContext(framework.WithDatabase())
	defer func() {
		if cleanupErr := ctx.Cleanup(); cleanupErr != nil {
			t.Logf("Cleanup failed: %v", cleanupErr)
		}
	}()

	// 获取用户服务
	userService := ctx.GetServices().GetUserService()

	// 测试用户注册
	testUser := &user.User{
		Handle: ctx.GetIsolator().IsolateHandle("u1"), // 短用户名: prefix(10) + _ + u1(2) = 13 <= 16
		Email:    ctx.GetIsolator().IsolateEmail("integration@example.com"),
		Pwd:      "password123",
		TimeZone: "Asia/Shanghai",
		Lang:     "zh-cn",
		LastIp:   "127.0.0.1",
		DisplayName: "集成测试用户",
		Remark:   "这是一个集成测试用户",
	}

	// 1. 测试用户注册
	t.Logf("注册前用户名: %s", testUser.Handle)
	err := userService.RegisterUser(testUser)
	assert.NoError(t, err, "用户注册失败")
	t.Logf("注册后用户名: %s, ID: %d", testUser.Handle, testUser.UserId)
	assert.Greater(t, testUser.UserId, int64(0), "用户ID应该大于0")
	assert.NotEqual(t, "password123", testUser.Pwd, "密码应该已加密")

	// 2. 测试用户名重复检查
	isExistHandle, err := userService.IsExistHandle(testUser.Handle)
	t.Logf("检查用户名: %s, 结果: %v, 错误: %v", testUser.Handle, isExistHandle, err)
	assert.NoError(t, err, "检查用户名存在时出错")
	assert.True(t, isExistHandle, "期望用户名已存在")

	// 3. 测试邮箱重复检查
	isExistEmail, err := userService.IsExistEmail(testUser.Email)
	assert.NoError(t, err, "检查邮箱存在时出错")
	assert.True(t, isExistEmail, "期望邮箱已存在")

	// 4. 测试用户登录（用户名）
	loginUser, err := userService.LoginWithUsernameOrEmail(testUser.Handle, "password123")
	require.NoError(t, err, "用户名登录失败")
	assert.Equal(t, testUser.UserId, loginUser.UserId, "登录用户ID不匹配")

	// 5. 测试用户登录（邮箱）
	loginUser2, err := userService.LoginWithUsernameOrEmail(testUser.Email, "password123")
	require.NoError(t, err, "邮箱登录失败")
	assert.Equal(t, testUser.UserId, loginUser2.UserId, "登录用户ID不匹配")

	// 6. 测试错误密码
	_, err = userService.LoginWithUsernameOrEmail(testUser.Handle, "wrongpassword")
	assert.Error(t, err, "期望密码错误")

	// 7. 测试获取用户信息
	userInfo, err := userService.GetUserByUserId(testUser.UserId)
	require.NoError(t, err, "获取用户信息失败")
	assert.Equal(t, testUser.Handle, userInfo.Handle, "用户名不匹配")
	assert.Equal(t, testUser.Email, userInfo.Email, "邮箱不匹配")

	// 8. 测试更新用户资料
	err = userService.UpdateUserProfile(testUser.UserId, "新的全名", "新的备注", "UTC", "en")
	assert.NoError(t, err, "更新用户资料失败")

	// 验证更新是否成功
	updatedUser, err := userService.GetUserByUserId(testUser.UserId)
	require.NoError(t, err, "获取更新后的用户信息失败")
	assert.Equal(t, "新的全名", updatedUser.DisplayName, "全名未更新")
	assert.Equal(t, "新的备注", updatedUser.Remark, "备注未更新")

	// 9. 测试修改密码
	err = userService.ChangeUserPassword(testUser.UserId, "password123", "newpassword456")
	assert.NoError(t, err, "修改密码失败")

	// 验证新密码可以登录
	_, err = userService.LoginWithUsernameOrEmail(testUser.Handle, "newpassword456")
	require.NoError(t, err, "使用新密码登录失败")

	// 验证旧密码不能登录
	_, err = userService.LoginWithUsernameOrEmail(testUser.Handle, "password123")
	assert.Error(t, err, "旧密码仍然可以登录，密码修改失败")

	t.Log("集成测试完成：用户注册、登录、资料更新、密码修改等功能均正常")
}

// TestUserService_Integration_ErrorCases 测试各种错误情况
func TestUserService_Integration_ErrorCases(t *testing.T) {
	ctx := framework.NewTestContext(framework.WithDatabase())
	defer func() {
		if cleanupErr := ctx.Cleanup(); cleanupErr != nil {
			t.Logf("Cleanup failed: %v", cleanupErr)
		}
	}()

	// 获取用户服务
	userService := ctx.GetServices().GetUserService()

	// 测试获取不存在的用户
	_, err := userService.GetUserByUserId(999999)
	assert.Error(t, err, "期望获取不存在用户时返回错误")

	// 测试修改不存在用户的密码
	err = userService.ChangeUserPassword(999999, "oldpwd", "newpwd")
	assert.Error(t, err, "期望修改不存在用户密码时返回错误")

	// 测试更新不存在用户的资料
	err = userService.UpdateUserProfile(999999, "新名字", "新备注", "", "")
	assert.Error(t, err, "期望更新不存在用户资料时返回错误")

	t.Log("错误情况测试完成：所有错误情况都正确处理")
}
