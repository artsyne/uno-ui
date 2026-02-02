package comment_test

import (
	"testing"

	"github.com/biuren/hebe/internal/domain/comment"
	"github.com/biuren/hebe/internal/domain/diary"
	"github.com/biuren/hebe/internal/domain/entry"
	"github.com/biuren/hebe/internal/domain/task"
	"github.com/biuren/hebe/internal/domain/user"
	"github.com/biuren/hebe/internal/infra/repository"
	"github.com/biuren/hebe/internal/utils/i18n"
	framework "github.com/biuren/hebe/tests/framework"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestData 创建测试用的用户、日记本和文章
func setupTestData(t *testing.T, ctx *framework.TestContext) (*user.User, *diary.Diary, *entry.Entry, *comment.Service) {
	t.Helper()

	// 创建真实的服务
	config := ctx.GetConfig().GetAppConfig()
	userRepo := repository.NewUserRepository()
	taskRepo := repository.NewTaskRepository()
	diaryRepo := repository.NewDiaryRepository()
	entryRepo := repository.NewEntryRepository()
	commentRepo := repository.NewCommentRepository()

	taskService := task.NewService(config, taskRepo)
	i18nManager := i18n.NewLocalizationManager(config)
	userService := user.NewService(config, userRepo, taskService, i18nManager, nil)
	diaryService := diary.NewService(config, diaryRepo)
	entryService := entry.NewService(config, userService, taskService, diaryRepo, entryRepo, i18nManager)
	commentService := comment.NewService(config, commentRepo, userService, entryService, diaryService)

	// 创建测试用户
	testUser := &user.User{
		Handle: ctx.GetDataGenerator().GenerateUniqueString("user"),
		Email:    ctx.GetDataGenerator().GenerateUniqueEmail("user"),
		Pwd:      "password123",
		DisplayName: "Test User",
		TimeZone: "UTC",
		Lang:     "zh-cn",
	}
	err := userService.RegisterUser(testUser)
	require.NoError(t, err)

	// 创建测试日记本
	testDiary := &diary.Diary{
		Title:         "测试日记本",
		Remark:        "用于评论测试",
		ViewSetting:   1, // 公开
		SocialSetting: 1, // 允许社交
		ThemeID:       1,
		AddUserID:     testUser.UserId,
		AddUserHandle: testUser.Handle,
		Status:        1,
	}
	diaryID, err := diaryService.AddDiary(testDiary)
	require.NoError(t, err)
	testDiary.DiaryID = diaryID

	// 创建测试文章
	testEntry := &entry.Entry{
		Title:       "测试文章",
		Body:        "这是一篇用于评论测试的文章",
		DiaryID:     testDiary.DiaryID,
		AddUserID:   testUser.UserId,
		AddUserHandle: testUser.Handle,
		AddUserDisplayName: testUser.DisplayName,
		Status:      1,
	}
	entryID, err := entryService.AddEntry(testEntry)
	require.NoError(t, err)
	testEntry.EntryID = entryID

	return testUser, testDiary, testEntry, commentService
}

func TestCommentService_AddComment(t *testing.T) {
	ctx := framework.NewTestContext(framework.WithDatabase())
	defer func() {
		if cleanupErr := ctx.Cleanup(); cleanupErr != nil {
			t.Logf("Cleanup failed: %v", cleanupErr)
		}
	}()

	testUser, _, testEntry, commentService := setupTestData(t, ctx.(*framework.TestContext))

	tests := []struct {
		name    string
		comment *comment.Comment
		wantErr bool
		errMsg  string
	}{
		{
			name: "成功添加评论",
			comment: &comment.Comment{
				EntryID:     testEntry.EntryID,
				Content:     "这是一条测试评论",
				AddUserID:   testUser.UserId,
				AddUserHandle: testUser.Handle,
				AddUserDisplayName: testUser.DisplayName,
				AddUserIP:   "127.0.0.1",
			},
			wantErr: false,
		},
		{
			name: "无效的用户ID",
			comment: &comment.Comment{
				EntryID:     testEntry.EntryID,
				Content:     "这是一条测试评论",
				AddUserID:   0,
				AddUserHandle: "invalid",
				AddUserIP:   "127.0.0.1",
			},
			wantErr: true,
			errMsg:  "invalid_user_id",
		},
		{
			name: "无效的文章ID",
			comment: &comment.Comment{
				EntryID:     99999,
				Content:     "这是一条测试评论",
				AddUserID:   testUser.UserId,
				AddUserHandle: testUser.Handle,
				AddUserDisplayName: testUser.DisplayName,
				AddUserIP:   "127.0.0.1",
			},
			wantErr: true,
			errMsg:  "entry_not_found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			commentID, err := commentService.AddComment(tt.comment)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				assert.Equal(t, int64(0), commentID)
			} else {
				assert.NoError(t, err)
				assert.Greater(t, commentID, int64(0))

				// 验证评论是否正确保存
				savedComment, err := commentService.GetCommentByID(commentID)
				assert.NoError(t, err)
				assert.Equal(t, tt.comment.EntryID, savedComment.EntryID)
				assert.Equal(t, tt.comment.Content, savedComment.Content)
				assert.Equal(t, tt.comment.AddUserID, savedComment.AddUserID)
				assert.NotEmpty(t, savedComment.AddUtcTime)
			}
		})
	}
}

func TestCommentService_GetCommentByID(t *testing.T) {
	ctx := framework.NewTestContext(framework.WithDatabase())
	defer func() {
		if cleanupErr := ctx.Cleanup(); cleanupErr != nil {
			t.Logf("Cleanup failed: %v", cleanupErr)
		}
	}()

	testUser, _, testEntry, commentService := setupTestData(t, ctx.(*framework.TestContext))

	// 先创建一个评论
	testComment := &comment.Comment{
		EntryID:     testEntry.EntryID,
		Content:     "测试评论内容",
		AddUserID:   testUser.UserId,
		AddUserHandle: testUser.Handle,
		AddUserDisplayName: testUser.DisplayName,
		AddUserIP:   "127.0.0.1",
	}
	commentID, err := commentService.AddComment(testComment)
	require.NoError(t, err)

	tests := []struct {
		name      string
		commentID int64
		wantErr   bool
	}{
		{
			name:      "成功获取评论",
			commentID: commentID,
			wantErr:   false,
		},
		{
			name:      "评论不存在",
			commentID: 99999,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			comment, err := commentService.GetCommentByID(tt.commentID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, comment)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, comment)
				assert.Equal(t, testComment.EntryID, comment.EntryID)
				assert.Equal(t, testComment.Content, comment.Content)
				assert.Equal(t, testComment.AddUserID, comment.AddUserID)
			}
		})
	}
}

func TestCommentService_ListCommentsByEntryID(t *testing.T) {
	ctx := framework.NewTestContext(framework.WithDatabase())
	defer func() {
		if cleanupErr := ctx.Cleanup(); cleanupErr != nil {
			t.Logf("Cleanup failed: %v", cleanupErr)
		}
	}()

	testUser, _, testEntry, commentService := setupTestData(t, ctx.(*framework.TestContext))

	// 创建多个评论
	comments := []*comment.Comment{
		{
			EntryID:     testEntry.EntryID,
			Content:     "第一条评论",
			AddUserID:   testUser.UserId,
			AddUserHandle: testUser.Handle,
			AddUserDisplayName: testUser.DisplayName,
			AddUserIP:   "127.0.0.1",
		},
		{
			EntryID:     testEntry.EntryID,
			Content:     "第二条评论",
			AddUserID:   testUser.UserId,
			AddUserHandle: testUser.Handle,
			AddUserDisplayName: testUser.DisplayName,
			AddUserIP:   "127.0.0.1",
		},
		{
			EntryID:     testEntry.EntryID,
			Content:     "第三条评论",
			AddUserID:   testUser.UserId,
			AddUserHandle: testUser.Handle,
			AddUserDisplayName: testUser.DisplayName,
			AddUserIP:   "127.0.0.1",
		},
	}

	for _, c := range comments {
		_, err := commentService.AddComment(c)
		require.NoError(t, err)
	}

	// 测试分页获取评论
	commentList, total, err := commentService.ListCommentsByEntryID(testEntry.EntryID, 1, 2)
	assert.NoError(t, err)
	assert.Equal(t, 3, total)
	assert.Len(t, commentList, 2)

	// 测试获取第二页
	commentList, total, err = commentService.ListCommentsByEntryID(testEntry.EntryID, 2, 2)
	assert.NoError(t, err)
	assert.Equal(t, 3, total)
	assert.Len(t, commentList, 1)
}

func TestCommentService_DeleteComment(t *testing.T) {
	ctx := framework.NewTestContext(framework.WithDatabase())
	defer func() {
		if cleanupErr := ctx.Cleanup(); cleanupErr != nil {
			t.Logf("Cleanup failed: %v", cleanupErr)
		}
	}()

	testUser, _, testEntry, commentService := setupTestData(t, ctx.(*framework.TestContext))

	// 先创建一个评论
	testComment := &comment.Comment{
		EntryID:     testEntry.EntryID,
		Content:     "待删除的评论",
		AddUserID:   testUser.UserId,
		AddUserHandle: testUser.Handle,
		AddUserDisplayName: testUser.DisplayName,
		AddUserIP:   "127.0.0.1",
	}
	commentID, err := commentService.AddComment(testComment)
	require.NoError(t, err)

	// 测试删除评论
	err = commentService.DeleteComment(commentID, testUser.UserId)
	assert.NoError(t, err)

	// 验证评论已被删除
	_, err = commentService.GetCommentByID(commentID)
	assert.Error(t, err)
}

func TestCommentService_GetCommentCountByEntryID(t *testing.T) {
	ctx := framework.NewTestContext(framework.WithDatabase())
	defer func() {
		if cleanupErr := ctx.Cleanup(); cleanupErr != nil {
			t.Logf("Cleanup failed: %v", cleanupErr)
		}
	}()

	testUser, _, testEntry, commentService := setupTestData(t, ctx.(*framework.TestContext))

	// 初始评论数应该为0
	count, err := commentService.GetCommentCountByEntryID(testEntry.EntryID)
	assert.NoError(t, err)
	assert.Equal(t, 0, count)

	// 添加几条评论
	for i := 0; i < 3; i++ {
		testComment := &comment.Comment{
			EntryID:     testEntry.EntryID,
			Content:     "测试评论",
			AddUserID:   testUser.UserId,
			AddUserHandle: testUser.Handle,
			AddUserDisplayName: testUser.DisplayName,
			AddUserIP:   "127.0.0.1",
		}
		_, err := commentService.AddComment(testComment)
		require.NoError(t, err)
	}

	// 验证评论数
	count, err = commentService.GetCommentCountByEntryID(testEntry.EntryID)
	assert.NoError(t, err)
	assert.Equal(t, 3, count)
}
