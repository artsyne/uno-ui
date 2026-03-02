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

// TestDiaryService_EmailI18n 测试邮件多语言功能
func TestDiaryService_EmailI18n(t *testing.T) {
	ctx := framework.NewTestContext(framework.WithDatabase())
	defer func() {
		if err := ctx.Cleanup(); err != nil {
			t.Log(err)
		}
	}()

	// 准备服务
	diaryRepo := repository.NewDiaryRepository()
	config := ctx.GetConfig().GetAppConfig()
	
	// 确保 Redis 清理干净，避免测试干扰
	// redisClient := ctx.GetRedisClient()
	// redisClient.FlushDB(context.Background())

	userRepo := repository.NewUserRepository()
	taskRepo := repository.NewTaskRepository()
	taskService := task.NewService(config, taskRepo)
	i18nManager := i18n.NewLocalizationManager(config)
	userService := user.NewService(config, userRepo, taskService, i18nManager, nil)
	diaryService := diary.NewService(config, diaryRepo, userService, taskService, i18nManager)

	// 1. 准备数据：创建拥有者
	owner := ctx.GetDataGenerator().CreateUser(
		core.WithHandle(ctx.GetIsolator().IsolateHandle("own")),
		core.WithEmail(ctx.GetDataGenerator().GenerateUniqueEmail("own")),
		core.WithLang("en"),
	)
	err := userService.RegisterUser(owner)
	require.NoError(t, err)

	// 创建日记本
	myDiary := &diary.Diary{
		Title:         "Global Diary",
		AddUserID:     owner.UserId,
		ViewSetting:   diary.ViewSettingPublic,
		SocialSetting: 1,
		ThemeID:       1,
		AddUTCTime:    time.Now().Unix(),
	}
	diaryID, err := diaryService.AddDiary(myDiary)
	require.NoError(t, err)

	// 2. 测试英文邀请
	t.Run("English Invitation", func(t *testing.T) {
		inviteeEn := ctx.GetDataGenerator().CreateUser(
			core.WithHandle(ctx.GetIsolator().IsolateHandle("en")),
			core.WithEmail(ctx.GetDataGenerator().GenerateUniqueEmail("en")),
			core.WithLang("en"),
		)
		err := userService.RegisterUser(inviteeEn)
		require.NoError(t, err)

		// 邀请，指定语言为 en
		err = diaryService.InviteAuthor(diaryID, owner.UserId, inviteeEn.Email, "en")
		require.NoError(t, err)

		// 等待异步任务
		time.Sleep(200 * time.Millisecond)

		// 验证任务
		tasks, err := taskRepo.GetJobTasks(100, "email")
		require.NoError(t, err)

		var foundTask *task.JobQueue
		for _, t := range tasks {
			// Payload 是 JSON，通常包含 email, subject, body
			// 这里我们简化检查，Task 结构体具体字段视实现而定，通常有 Payload 或 Args
			// 假设 TaskService.CreateEmailTask 存储了 email 在某个字段，或者我们通过 Payload 匹配
			// 假设 GetJobTasks 返回的是 task.Task 结构，查看 task.go 定义
			// CreateEmailTask 实现：
			// payload := map[string]string{"email": email, "subject": subject, "body": body}
			// taskRepo.Push(queueName, payload)
			// 但是 GetJobTasks 是怎么实现的？
			// 让我们检查 taskRepo.GetJobTasks 的实现或 task.Task 结构
			// 这里假设 tasks 是 []TaskEntity，并且 Payload 是 string
			if strings.Contains(t.Body, inviteeEn.Email) {
				foundTask = &t
				break
			}
		}
		require.NotNil(t, foundTask, "Email task for EN user not found")

		// 验证英文内容
		// Subject: Invitation to join diary
		assert.Contains(t, foundTask.Body, "Invitation to join diary")
		// Body: You have been invited
		assert.Contains(t, foundTask.Body, "You have been invited")
		// Button: View Invitation
		assert.Contains(t, foundTask.Body, "View Invitation")
	})

	// 3. 测试中文邀请
	t.Run("Chinese Invitation", func(t *testing.T) {
		inviteeZh := ctx.GetDataGenerator().CreateUser(
			core.WithHandle(ctx.GetIsolator().IsolateHandle("zh")),
			core.WithEmail(ctx.GetDataGenerator().GenerateUniqueEmail("zh")),
			core.WithLang("zh-cn"),
		)
		err := userService.RegisterUser(inviteeZh)
		require.NoError(t, err)

		// 邀请，指定语言为 zh-cn
		err = diaryService.InviteAuthor(diaryID, owner.UserId, inviteeZh.Email, "zh-cn")
		require.NoError(t, err)

		// 等待异步任务
		time.Sleep(200 * time.Millisecond)

		// 验证任务
		tasks, err := taskRepo.GetJobTasks(100, "email")
		require.NoError(t, err)

		var foundTask *task.JobQueue
		for _, t := range tasks {
			if strings.Contains(t.Body, inviteeZh.Email) {
				foundTask = &t
				break
			}
		}
		require.NotNil(t, foundTask, "Email task for ZH user not found")

		// 验证中文内容
		// Subject: 邀请加入日记本
		// JSON payload 中的 unicode 可能会被转义，也可能不会，取决于存储方式
		// 如果是 JSON 字符串，中文可能会是 \uXXXX
		// 但 strings.Contains 直接匹配字符串应该没问题，除非 payload 被编码过
		// 假设 payload 是存的原始 JSON string
		
		// 我们可以先打印出来看看，或者尝试匹配关键英文部分（如果有）
		// 但中文邮件全是中文。
		// 我们可以尝试匹配 "邀请加入日记本"
		// 如果失败，可能是编码问题。
		
		// 简单检查 Subject
		assert.Contains(t, foundTask.Body, "邀请加入日记本")
		// Body: 你已被邀请成为日记本
		assert.Contains(t, foundTask.Body, "你已被邀请成为日记本")
		// Button: 查看邀请
		assert.Contains(t, foundTask.Body, "查看邀请")
	
	})
}
