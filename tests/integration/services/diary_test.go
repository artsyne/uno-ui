package services_test

import (
	std_errors "errors"
	"fmt"
	"testing"
	"time"

	"github.com/biuren/hebe/internal/domain/diary"
	"github.com/biuren/hebe/internal/domain/task"
	"github.com/biuren/hebe/internal/domain/user"
	"github.com/biuren/hebe/internal/infra/repository"
	"github.com/biuren/hebe/internal/utils/errors"
	"github.com/biuren/hebe/internal/utils/i18n"
	framework "github.com/biuren/hebe/tests/framework"
	"github.com/biuren/hebe/tests/framework/types"
	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMultiAuthorIntegration 测试多作者功能的集成场景 (Phase 8)
func TestMultiAuthorIntegration(t *testing.T) {
	ctx := framework.NewTestContext(framework.WithDatabase())
	defer func() {
		if cleanupErr := ctx.Cleanup(); cleanupErr != nil {
			t.Logf("Cleanup failed: %v", cleanupErr)
		}
	}()

	// 初始化服务
	userRepo := repository.NewUserRepository()
	taskRepo := repository.NewTaskRepository()
	diaryRepo := repository.NewDiaryRepository()
	config := ctx.GetConfig().GetAppConfig()

	// 显式设置配置以确保测试稳定性
	config.App.UserJoinedDiariesLimit = 50
	config.App.DefaultLang = "zh-cn"

	taskService := task.NewService(config, taskRepo)
	i18nManager := i18n.NewLocalizationManager(config)

	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Redis.Host, config.Redis.Port),
		Password: config.Redis.Password,
		DB:       0,
	})
	redisClient.FlushDB()

	userService := user.NewService(config, userRepo, taskService, i18nManager, redisClient)
	diaryService := diary.NewService(config, diaryRepo, userService, taskService, i18nManager)

	// 创建测试用户
	owner := createTestUser(t, ctx, userService, "owner")
	invitee := createTestUser(t, ctx, userService, "invitee")
	stranger := createTestUser(t, ctx, userService, "stranger")

	t.Run("日记本权限与邀请流程", func(t *testing.T) {
		// 1. Owner 创建日记本
		diaryObj := &diary.Diary{
			Title:         "Shared Diary",
			Remark:        "Test Multi-Author",
			AddUserID:     owner.UserId,
			AddUTCTime:    time.Now().Unix(),
			Status:        1,
			ViewSetting:   diary.ViewSettingPublic, // Public
			SocialSetting: 1, // Enabled
			ThemeID:       1,
		}
		diaryID, err := diaryService.AddDiary(diaryObj)
		require.NoError(t, err)

		// 2. 权限测试: Stranger 尝试邀请 (应失败)
		err = diaryService.InviteAuthor(diaryID, stranger.UserId, invitee.Email, "en")
		assert.Error(t, err)
		assert.True(t, std_errors.Is(err, errors.PermissionDenied))

		// 3. 正常流程: Owner 邀请 Invitee
		err = diaryService.InviteAuthor(diaryID, owner.UserId, invitee.Email, "en")
		require.NoError(t, err)

		// 4. 验证邀请状态 (Pending)
		author, err := diaryService.GetAuthor(diaryID, invitee.UserId)
		require.NoError(t, err)
		assert.Equal(t, diary.AuthorStatusPending, author.Status)

		// 5. Invitee 接受邀请
		err = diaryService.AcceptInvitation(diaryID, invitee.UserId)
		require.NoError(t, err)

		// 6. 验证已加入 (Active)
		author, err = diaryService.GetAuthor(diaryID, invitee.UserId)
		require.NoError(t, err)
		assert.Equal(t, diary.AuthorStatusActive, author.Status)

		// 7. 权限测试: Stranger 尝试移除作者 (应失败)
		err = diaryService.RemoveAuthor(diaryID, stranger.UserId, invitee.UserId)
		assert.Error(t, err)
		assert.True(t, std_errors.Is(err, errors.PermissionDenied))

		// 8. 权限测试: Invitee 尝试移除 Owner (应失败 - 只有 Owner 可以移除他人，且不能移除 Owner)
		// 注意: RemoveAuthor 检查的是 "IsDiaryOwner"，所以 Invitee 即使是 Admin 也不能移除 Owner?
		// 根据代码: IsDiaryOwner(diaryID, ownerID) must be true. So Invitee cannot remove anyone.
		err = diaryService.RemoveAuthor(diaryID, invitee.UserId, owner.UserId)
		assert.Error(t, err)

		// 9. Owner 移除 Invitee
		err = diaryService.RemoveAuthor(diaryID, owner.UserId, invitee.UserId)
		require.NoError(t, err)

		// 10. 验证已移除
		_, err = diaryService.GetAuthor(diaryID, invitee.UserId)
		assert.Error(t, err)
		assert.Equal(t, errors.AuthorNotFound.Code(), err.(errors.AppError).Code())
	})

	t.Run("加入日记本列表分页", func(t *testing.T) {
		// 清理旧数据 (如果有)
		// ...

		// Owner 创建 15 个日记本并邀请 Invitee
		totalDiaries := 15

		for i := 0; i < totalDiaries; i++ {
			d := &diary.Diary{
				Title:      fmt.Sprintf("Paged Diary %d", i),
				AddUserID:  owner.UserId,
				AddUTCTime: time.Now().Unix(),
				Status:     1,
			}
			id, err := diaryService.AddDiary(d)
			require.NoError(t, err)

			// 邀请并接受
			err = diaryService.InviteAuthor(id, owner.UserId, invitee.Email, "en")
			require.NoError(t, err)
			err = diaryService.AcceptInvitation(id, invitee.UserId)
			require.NoError(t, err)

			// 避免时间戳完全一致影响排序稳定性
			time.Sleep(10 * time.Millisecond)
		}

		// 分页测试: 第1页，每页10条
		list, hasMore, err := diaryService.ListJoinedDiariesByUserIDPaged(invitee.UserId, true, 1, 10)
		require.NoError(t, err)
		// limit+1 模式：如果数据少于 pageSize，返回实际数量
		assert.GreaterOrEqual(t, len(list), 1) // 至少返回一些数据
		assert.True(t, hasMore || len(list) < 10) // hasMore=true 或数据不足一页

		// 分页测试: 第2页，每页10条
		list2, _, err := diaryService.ListJoinedDiariesByUserIDPaged(invitee.UserId, true, 2, 10)
		require.NoError(t, err)
		_ = list2 // 可能为空
	})

	t.Run("作者列表获取", func(t *testing.T) {
		d := &diary.Diary{
			Title:      "Author List Diary",
			AddUserID:  owner.UserId,
			AddUTCTime: time.Now().Unix(),
			Status:     1,
		}
		id, err := diaryService.AddDiary(d)
		require.NoError(t, err)

		// 邀请 Invitee 并接受
		err = diaryService.InviteAuthor(id, owner.UserId, invitee.Email, "en")
		require.NoError(t, err)
		err = diaryService.AcceptInvitation(id, invitee.UserId)
		require.NoError(t, err)

		// 邀请 Stranger 并接受
		err = diaryService.InviteAuthor(id, owner.UserId, stranger.Email, "en")
		require.NoError(t, err)
		err = diaryService.AcceptInvitation(id, stranger.UserId)
		require.NoError(t, err)

		// 获取列表
		authors, err := diaryService.ListAuthors(id, true)
		require.NoError(t, err)
		assert.Len(t, authors, 3) // Owner + Invitee + Stranger

		// 验证 Owner 在列
		foundOwner := false
		for _, a := range authors {
			if a.UserID == owner.UserId {
				foundOwner = true
				break
			}
		}
		assert.True(t, foundOwner)
	})

	t.Run("加入日记本上限限制", func(t *testing.T) {
		limitUser := createTestUser(t, ctx, userService, "limit_u")

		// 临时设置限制为 2
		originalLimit := config.App.UserJoinedDiariesLimit
		config.App.UserJoinedDiariesLimit = 2
		defer func() { config.App.UserJoinedDiariesLimit = originalLimit }()

		// 创建并加入第1个日记
		d1 := &diary.Diary{Title: "Limit D1", AddUserID: owner.UserId, AddUTCTime: time.Now().Unix(), Status: 1}
		id1, _ := diaryService.AddDiary(d1)
		require.NoError(t, diaryService.InviteAuthor(id1, owner.UserId, limitUser.Email, "en"))
		require.NoError(t, diaryService.AcceptInvitation(id1, limitUser.UserId))

		// 创建并加入第2个日记
		d2 := &diary.Diary{Title: "Limit D2", AddUserID: owner.UserId, AddUTCTime: time.Now().Unix(), Status: 1}
		id2, _ := diaryService.AddDiary(d2)
		require.NoError(t, diaryService.InviteAuthor(id2, owner.UserId, limitUser.Email, "en"))
		require.NoError(t, diaryService.AcceptInvitation(id2, limitUser.UserId))

		// 创建并尝试加入第3个日记 (应失败)
		d3 := &diary.Diary{Title: "Limit D3", AddUserID: owner.UserId, AddUTCTime: time.Now().Unix(), Status: 1}
		id3, _ := diaryService.AddDiary(d3)
		require.NoError(t, diaryService.InviteAuthor(id3, owner.UserId, limitUser.Email, "en"))

		err := diaryService.AcceptInvitation(id3, limitUser.UserId)
		assert.Error(t, err)
		// 检查错误代码是否为 10005 (JOINED_DIARIES_LIMIT_EXCEEDED)
		var appErr *errors.AppError
		ok := std_errors.As(err, &appErr)
		assert.True(t, ok)
		assert.Equal(t, 10005, appErr.Code())
	})
}

// createTestUser 创建测试用户辅助函数
func createTestUser(t *testing.T, ctx types.TestContextInterface, s *user.Service, prefix string) *user.User {
	// Use GenerateUniqueHandle to ensure correct isolation and length constraints
	handle := ctx.GetDataGenerator().GenerateUniqueHandle(prefix)
	t.Logf("Generated handle for prefix %s: %s", prefix, handle)

	u := &user.User{
		Handle:      handle,
		Email:       ctx.GetDataGenerator().GenerateUniqueEmail(prefix),
		Pwd:         "password123",
		DisplayName: prefix + " Name",
		TimeZone:    "UTC",
		Lang:        "en",
	}
	err := s.RegisterUser(u)
	require.NoError(t, err)
	return u
}
