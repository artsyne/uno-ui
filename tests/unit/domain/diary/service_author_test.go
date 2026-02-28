package diary_test

import (
	"strings"
	"testing"
	"time"

	"github.com/biuren/hebe/internal/domain/diary"
	"github.com/biuren/hebe/internal/domain/task"
	"github.com/biuren/hebe/internal/domain/user"
	"github.com/biuren/hebe/internal/infra/repository"
	"github.com/biuren/hebe/internal/utils/i18n"
	framework "github.com/biuren/hebe/tests/framework"
	core "github.com/biuren/hebe/tests/framework/core"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDiaryService_MultiAuthor 测试多人日记本功能
func TestDiaryService_MultiAuthor(t *testing.T) {
	ctx := framework.NewTestContext(framework.WithDatabase())
	defer func() {
		if err := ctx.Cleanup(); err != nil {
			t.Log(err)
		}
	}()

	// 准备服务
	diaryRepo := repository.NewDiaryRepository()
	config := ctx.GetConfig().GetAppConfig()
	
	userRepo := repository.NewUserRepository()
	taskRepo := repository.NewTaskRepository()
	taskService := task.NewService(config, taskRepo)
	i18nManager := i18n.NewLocalizationManager(config)
	userService := user.NewService(config, userRepo, taskService, i18nManager, nil)
	diaryService := diary.NewService(config, diaryRepo, userService, taskService, i18nManager)

	// 1. 准备数据：创建拥有者和被邀请者
	owner := ctx.GetDataGenerator().CreateUser(
		core.WithHandle(ctx.GetIsolator().IsolateHandle("own")),
		core.WithEmail(ctx.GetDataGenerator().GenerateUniqueEmail("own")),
		core.WithLang("en"),
	)
	err := userService.RegisterUser(owner)
	require.NoError(t, err)

	invitee := ctx.GetDataGenerator().CreateUser(
		core.WithHandle(ctx.GetIsolator().IsolateHandle("inv")),
		core.WithEmail(ctx.GetDataGenerator().GenerateUniqueEmail("inv")),
		core.WithLang("zh-cn"),
	)
	err = userService.RegisterUser(invitee)
	require.NoError(t, err)

	// 创建日记本
	myDiary := &diary.Diary{
		Title:         "Shared Diary",
		AddUserID:     owner.UserId,
		ViewSetting:   1,
		SocialSetting: 1,
		ThemeID:       1,
		AddUTCTime:    time.Now().Unix(),
	}
	diaryID, err := diaryService.AddDiary(myDiary)
	require.NoError(t, err)

	t.Run("邀请作者_通过Email", func(t *testing.T) {
		err := diaryService.InviteAuthor(diaryID, owner.UserId, invitee.Email, "en")
		require.NoError(t, err)

		// 验证状态为Pending
		authors, err := diaryService.ListAuthors(diaryID, false)
		require.NoError(t, err)
		assert.Len(t, authors, 2) // Owner + Invitee

		var invitedAuthor diary.Author
		found := false
		for _, a := range authors {
			if a.UserID == invitee.UserId {
				invitedAuthor = a
				found = true
				break
			}
		}
		require.True(t, found, "Invited author not found")
		assert.Equal(t, diary.AuthorStatusPending, invitedAuthor.Status)

		// Wait for async email task
		time.Sleep(100 * time.Millisecond)

		// Verify email task
		tasks, err := taskRepo.GetJobTasks(10, "email")
		require.NoError(t, err)

		var foundEmail bool
		for _, task := range tasks {
			if strings.Contains(task.Body, invitee.Email) {
				foundEmail = true
				// Verify template content was rendered
				assert.Contains(t, task.Body, "Invitation to join diary")
				assert.Contains(t, task.Body, "Hello "+invitee.DisplayName)
				assert.Contains(t, task.Body, "Shared Diary")
				assert.Contains(t, task.Body, "View Invitation")
				break
			}
		}
		assert.True(t, foundEmail, "Invitation email task not found in queue")
	})

	t.Run("重复邀请_应该失败", func(t *testing.T) {
		err := diaryService.InviteAuthor(diaryID, owner.UserId, invitee.Email, "en")
		assert.Error(t, err)
	})

	t.Run("接受邀请", func(t *testing.T) {
		err := diaryService.AcceptInvitation(diaryID, invitee.UserId)
		require.NoError(t, err)

		// 验证状态为Active
		authors, err := diaryService.ListAuthors(diaryID, true)
		require.NoError(t, err)

		found := false
		for _, a := range authors {
			if a.UserID == invitee.UserId {
				assert.Equal(t, diary.AuthorStatusActive, a.Status)
				found = true
			}
		}
		require.True(t, found)
	})

	t.Run("删除作者", func(t *testing.T) {
		err := diaryService.RemoveAuthor(diaryID, owner.UserId, invitee.UserId)
		require.NoError(t, err)

		authors, err := diaryService.ListAuthors(diaryID, false)
		require.NoError(t, err)
		assert.Len(t, authors, 1) // Only owner
	})

	t.Run("非拥有者无法邀请", func(t *testing.T) {
		// invitee已经被移除了，我们重新邀请并接受，让他成为作者
		// 但这里我们测试的是非拥有者（比如invitee如果再次加入）是否有权限邀请
		// 或者我们创建一个新的第三方用户
		otherUser := ctx.GetDataGenerator().CreateUser(
			core.WithHandle(ctx.GetIsolator().IsolateHandle("other")),
			core.WithEmail(ctx.GetDataGenerator().GenerateUniqueEmail("other")),
		)
		err := userService.RegisterUser(otherUser)
		require.NoError(t, err)

		// otherUser 尝试邀请 invitee (otherUser 不是 diary owner 也不是 author)
		err = diaryService.InviteAuthor(diaryID, otherUser.UserId, invitee.Email, "en")
		assert.Error(t, err) 
	})

	t.Run("通过Handle邀请", func(t *testing.T) {
		// 先清理之前的状态，invitee已经被移除了
		
		err := diaryService.InviteAuthor(diaryID, owner.UserId, invitee.Handle, "en")
		require.NoError(t, err)

		// 验证邀请成功
		authors, err := diaryService.ListAuthors(diaryID, false)
		require.NoError(t, err)
		found := false
		for _, a := range authors {
			if a.UserID == invitee.UserId && a.Status == diary.AuthorStatusPending {
				found = true
				break
			}
		}
		assert.True(t, found)
	})

	t.Run("退出日记本", func(t *testing.T) {
		// 先接受邀请
		err := diaryService.AcceptInvitation(diaryID, invitee.UserId)
		require.NoError(t, err)

		// 退出
		err = diaryService.QuitDiary(diaryID, invitee.UserId)
		require.NoError(t, err)

		// 验证已退出
		authors, err := diaryService.ListAuthors(diaryID, false)
		require.NoError(t, err)
		assert.Len(t, authors, 1)
		assert.Equal(t, owner.UserId, authors[0].UserID)
	})

	t.Run("拥有者无法退出日记本", func(t *testing.T) {
		err := diaryService.QuitDiary(diaryID, owner.UserId)
		assert.Error(t, err)
	})
}
