package post_test

import (
	"fmt"
	"testing"

	"github.com/biuren/hebe/internal/utils/config"
	"github.com/biuren/hebe/internal/domain/diary"
	"github.com/biuren/hebe/internal/domain/post"
	"github.com/biuren/hebe/internal/domain/task"
	"github.com/biuren/hebe/internal/domain/user"
	"github.com/biuren/hebe/internal/infra/repository"
	framework "github.com/biuren/hebe/tests/framework/core"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestUserAndDiary 创建测试用的用户和日记本
func setupTestUserAndDiary(t *testing.T, ctx *framework.TestContext, userService user.ServiceInterface, diaryService diary.ServiceInterface) (*user.User, *diary.Diary) {
	t.Helper()

	// 创建用户 - 使用随机数据
	testUser := &user.User{
		UserName: ctx.GetDataGenerator().GenerateUniqueString("user"),
		Email:    ctx.GetDataGenerator().GenerateUniqueEmail("user"),
		Pwd:      "password123",
		TimeZone: "UTC",
		Lang:     "zh-cn",
	}
	err := userService.RegisterUser(testUser)
	require.NoError(t, err)

	// 创建日记本 - 使用随机数据
	testDiary := &diary.Diary{
		Title:         "测试日记本",
		Remark:        "用于文章测试",
		ViewSetting:   1, // 必需字段：1-公开，2-私有
		SocialSetting: 1, // 必需字段：1-允许社交，2-不允许社交
		ThemeID:       1, // 必需字段：主题ID，范围1-16
		AddUserID:     testUser.UserId,
		AddUserName:   testUser.UserName,
		Status:        1,
	}
	diaryID, err := diaryService.AddDiary(testDiary)
	require.NoError(t, err)
	testDiary.DiaryID = diaryID

	return testUser, testDiary
}

// TestPostService_AddPost 测试添加文章功能
func TestPostService_AddPost(t *testing.T) {
	// 使用新测试框架
	ctx := framework.NewTestContext("TestPostService_AddPost", framework.WithDatabase())
	defer ctx.Cleanup()

	// 创建真实的服务
	config := config.GetConfig()
	userRepo := repository.NewUserRepository()
	taskRepo := repository.NewTaskRepository()
	taskService := task.NewService(taskRepo, config)
	userService := user.NewService(userRepo, taskService, config)
	diaryRepo := repository.NewDiaryRepository()
	diaryService := diary.NewService(diaryRepo)
	postRepo := repository.NewPostRepository(ctx.GetDatabase().GetDB())
	postService := post.NewService(userRepo, taskService, diaryRepo, postRepo)

	// 设置测试数据
	testUser, testDiary := setupTestUserAndDiary(t, ctx, userService, diaryService)

	t.Run("成功添加文章", func(t *testing.T) {
		newPost := &post.Post{
			Title:       "测试文章标题",
			Body:        "这是一篇测试文章的内容",
			DiaryID:     testDiary.DiaryID,
			AddUserID:   testUser.UserId,
			AddUserName: testUser.UserName,
			Status:      1,
		}

		postID, err := postService.AddPost(newPost)
		require.NoError(t, err)
		assert.True(t, postID > 0)
		assert.Equal(t, postID, newPost.PostID)
	})

	t.Run("标题为空", func(t *testing.T) {
		invalidPost := &post.Post{
			Title:       "", // 空标题
			Body:        "内容",
			DiaryID:     testDiary.DiaryID,
			AddUserID:   testUser.UserId,
			AddUserName: testUser.UserName,
			Status:      1,
		}

		_, err := postService.AddPost(invalidPost)
		assert.Error(t, err)
	})

	t.Run("内容为空", func(t *testing.T) {
		invalidPost := &post.Post{
			Title:       "标题",
			Body:        "", // 空内容
			DiaryID:     testDiary.DiaryID,
			AddUserID:   testUser.UserId,
			AddUserName: testUser.UserName,
			Status:      1,
		}

		_, err := postService.AddPost(invalidPost)
		assert.Error(t, err)
	})

	t.Run("无效的日记本ID", func(t *testing.T) {
		invalidPost := &post.Post{
			Title:       "标题",
			Body:        "内容",
			DiaryID:     99999, // 不存在的日记本ID
			AddUserID:   testUser.UserId,
			AddUserName: testUser.UserName,
			Status:      1,
		}

		_, err := postService.AddPost(invalidPost)
		assert.Error(t, err)
	})

	t.Run("无效的用户ID", func(t *testing.T) {
		invalidPost := &post.Post{
			Title:       "标题",
			Body:        "内容",
			DiaryID:     testDiary.DiaryID,
			AddUserID:   0, // 无效用户ID
			AddUserName: "invalid",
			Status:      1,
		}

		_, err := postService.AddPost(invalidPost)
		assert.Error(t, err)
	})
}

// TestPostService_GetPostByID 测试根据ID获取文章
func TestPostService_GetPostByID(t *testing.T) {
	// 使用新测试框架
	ctx := framework.NewTestContext("TestPostService_GetPostByID", framework.WithDatabase())
	defer ctx.Cleanup()

	// 创建真实的服务
	config := config.GetConfig()
	userRepo := repository.NewUserRepository()
	taskRepo := repository.NewTaskRepository()
	taskService := task.NewService(taskRepo, config)
	userService := user.NewService(userRepo, taskService, config)
	diaryRepo := repository.NewDiaryRepository()
		diaryService := diary.NewService(diaryRepo)
		postRepo := repository.NewPostRepository(ctx.GetDatabase().GetDB())
	postService := post.NewService(userRepo, taskService, diaryRepo, postRepo)

	// 设置测试数据
	testUser, testDiary := setupTestUserAndDiary(t, ctx, userService, diaryService)

	// 添加一篇文章
	newPost := &post.Post{
		Title:       "可获取的文章",
		Body:        "测试获取功能的文章内容",
		Summary:     "测试摘要",
		DiaryID:     testDiary.DiaryID,
		AddUserID:   testUser.UserId,
		AddUserName: testUser.UserName,
		Status:      1,
	}
	postID, err := postService.AddPost(newPost)
	require.NoError(t, err)

	t.Run("成功获取文章", func(t *testing.T) {
		retrievedPost, err := postService.GetPostByID(postID)
		require.NoError(t, err)
		assert.NotNil(t, retrievedPost)
		assert.Equal(t, "可获取的文章", retrievedPost.Title)
		assert.Equal(t, "测试获取功能的文章内容", retrievedPost.Body)
		assert.Equal(t, testUser.UserId, retrievedPost.AddUserID)
		assert.Equal(t, testDiary.DiaryID, retrievedPost.DiaryID)
	})

	t.Run("文章不存在", func(t *testing.T) {
		retrievedPost, err := postService.GetPostByID(99999)
		assert.Error(t, err)
		assert.Nil(t, retrievedPost)
	})
}

// TestPostService_UpdatePost 测试更新文章
func TestPostService_UpdatePost(t *testing.T) {
	// 使用新测试框架
	ctx := framework.NewTestContext("TestPostService_UpdatePost", framework.WithDatabase())
	defer ctx.Cleanup()

	// 创建真实的服务
	config := config.GetConfig()
	userRepo := repository.NewUserRepository()
	taskRepo := repository.NewTaskRepository()
	taskService := task.NewService(taskRepo, config)
	userService := user.NewService(userRepo, taskService, config)
	diaryRepo := repository.NewDiaryRepository()
		diaryService := diary.NewService(diaryRepo)
		postRepo := repository.NewPostRepository(ctx.GetDatabase().GetDB())
		postService := post.NewService(userRepo, taskService, diaryRepo, postRepo)

	// 设置测试数据
	testUser, testDiary := setupTestUserAndDiary(t, ctx, userService, diaryService)

	// 添加一篇文章
	originalPost := &post.Post{
		Title:       "原始标题",
		Body:        "原始内容",
		Summary:     "原始摘要",
		DiaryID:     testDiary.DiaryID,
		AddUserID:   testUser.UserId,
		AddUserName: testUser.UserName,
		Status:      1,
	}
	postID, err := postService.AddPost(originalPost)
	require.NoError(t, err)

	t.Run("成功更新文章", func(t *testing.T) {
		updatedPost := &post.Post{
			PostID:        postID,
			Title:         "更新后的标题",
			Body:          "更新后的内容",
			Summary:       "更新后的摘要",
			DiaryID:       testDiary.DiaryID,
			AddUserID:     testUser.UserId,
			AddUserName:   testUser.UserName,
			UpdateUTCTime: 1642680200,
			Status:        1,
		}

		err := postService.UpdatePost(updatedPost)
		assert.NoError(t, err)

		// 验证更新是否成功
		retrievedPost, err := postService.GetPostByID(postID)
		require.NoError(t, err)
		assert.Equal(t, "更新后的标题", retrievedPost.Title)
		assert.Equal(t, "更新后的内容", retrievedPost.Body)
		assert.Equal(t, "更新后的摘要", retrievedPost.Summary)
	})

	t.Run("更新不存在的文章", func(t *testing.T) {
		nonExistentPost := &post.Post{
			PostID:      99999,
			Title:       "不存在的文章",
			Body:        "测试",
			DiaryID:     testDiary.DiaryID,
			AddUserID:   testUser.UserId,
			AddUserName: testUser.UserName,
			Status:      1,
		}

		err := postService.UpdatePost(nonExistentPost)
		assert.Error(t, err)
	})

	t.Run("标题为空时更新失败", func(t *testing.T) {
		invalidPost := &post.Post{
			PostID:      postID,
			Title:       "", // 空标题
			Body:        "内容",
			DiaryID:     testDiary.DiaryID,
			AddUserID:   testUser.UserId,
			AddUserName: testUser.UserName,
			Status:      1,
		}

		err := postService.UpdatePost(invalidPost)
		assert.Error(t, err)
	})

	t.Run("内容为空时更新失败", func(t *testing.T) {
		invalidPost := &post.Post{
			PostID:      postID,
			Title:       "标题",
			Body:        "", // 空内容
			DiaryID:     testDiary.DiaryID,
			AddUserID:   testUser.UserId,
			AddUserName: testUser.UserName,
			Status:      1,
		}

		err := postService.UpdatePost(invalidPost)
		assert.Error(t, err)
	})
}

// TestPostService_ListPostsByUserID 测试分页获取用户文章列表
func TestPostService_ListPostsByUserID(t *testing.T) {
	// 使用新测试框架
	ctx := framework.NewTestContext("TestPostService_ListPostsByUserID", framework.WithDatabase())
	defer ctx.Cleanup()

	// 创建真实的服务
	config := config.GetConfig()
	userRepo := repository.NewUserRepository()
	taskRepo := repository.NewTaskRepository()
	taskService := task.NewService(taskRepo, config)
	userService := user.NewService(userRepo, taskService, config)
	diaryRepo := repository.NewDiaryRepository()
	diaryService := diary.NewService(diaryRepo)
	postRepo := repository.NewPostRepository(ctx.GetDatabase().GetDB())
	postService := post.NewService(userRepo, taskService, diaryRepo, postRepo)

	// 设置测试数据
	testUser, testDiary := setupTestUserAndDiary(t, ctx, userService, diaryService)

	// 创建另一个用户和日记本
	otherUser := &user.User{
		UserName: ctx.GetDataGenerator().GenerateUniqueString("other_user"),
		Email:    ctx.GetDataGenerator().GenerateUniqueEmail("other"),
		Pwd:      "password123",
		TimeZone: "UTC",
		Lang:     "zh-cn",
	}
	err := userService.RegisterUser(otherUser)
	require.NoError(t, err)

	otherDiary := &diary.Diary{
		Title:         "其他用户的日记本",
		Remark:        "测试",
		AddUserID:     otherUser.UserId,
		AddUserName:   otherUser.UserName,
		Status:        1,
		ViewSetting:   1,
		SocialSetting: 1,
		ThemeID:       1,
	}
	otherDiaryID, err := diaryService.AddDiary(otherDiary)
	require.NoError(t, err)

	// 为testUser添加5篇文章
	for i := 1; i <= 5; i++ {
		post := &post.Post{
			Title:       fmt.Sprintf("文章%d", i),
			Body:        fmt.Sprintf("第%d篇文章内容", i),
			Summary:     fmt.Sprintf("第%d篇文章摘要", i),
			DiaryID:     testDiary.DiaryID,
			AddUserID:   testUser.UserId,
			AddUserName: testUser.UserName,
			Status:      1,
			InfoType:    1,
		}
		_, err := postService.AddPost(post)
		require.NoError(t, err)
	}

	// 为otherUser添加1篇文章
	otherPost := &post.Post{
		Title:       "其他用户的文章",
		Body:        "其他用户的文章内容",
		Summary:     "其他用户的文章摘要",
		DiaryID:     otherDiaryID,
		AddUserID:   otherUser.UserId,
		AddUserName: otherUser.UserName,
		Status:      1,
		InfoType:    1,
	}
	_, err = postService.AddPost(otherPost)
	require.NoError(t, err)

	t.Run("第一页，每页3篇", func(t *testing.T) {
		posts, total, err := postService.ListPostsByUserID(testUser.UserId, testDiary.DiaryID, 1, 3)
		require.NoError(t, err)
		assert.Len(t, posts, 3)
		assert.Equal(t, 5, total)

		// 验证返回的都是指定用户和日记本的文章
		for _, p := range posts {
			assert.Equal(t, testUser.UserId, p.AddUserID)
			assert.Equal(t, testDiary.DiaryID, p.DiaryID)
		}
	})

	t.Run("第二页，每页3篇", func(t *testing.T) {
		posts, total, err := postService.ListPostsByUserID(testUser.UserId, testDiary.DiaryID, 2, 3)
		require.NoError(t, err)
		assert.Len(t, posts, 2) // 剩余2篇
		assert.Equal(t, 5, total)
	})

	t.Run("获取所有文章", func(t *testing.T) {
		posts, total, err := postService.ListPostsByUserID(testUser.UserId, testDiary.DiaryID, 1, 10)
		require.NoError(t, err)
		assert.Len(t, posts, 5)
		assert.Equal(t, 5, total)
	})

	t.Run("获取其他用户的文章", func(t *testing.T) {
		posts, total, err := postService.ListPostsByUserID(otherUser.UserId, otherDiaryID, 1, 10)
		require.NoError(t, err)
		assert.Len(t, posts, 1)
		assert.Equal(t, 1, total)
		assert.Equal(t, "其他用户的文章", posts[0].Title)
	})

	t.Run("不存在的用户", func(t *testing.T) {
		posts, total, err := postService.ListPostsByUserID(99999, testDiary.DiaryID, 1, 10)
		assert.NoError(t, err) // 应该返回空列表而不是错误
		assert.Len(t, posts, 0)
		assert.Equal(t, 0, total)
	})

	t.Run("不存在的日记本", func(t *testing.T) {
		posts, total, err := postService.ListPostsByUserID(testUser.UserId, 99999, 1, 10)
		assert.NoError(t, err) // 应该返回空列表而不是错误
		assert.Len(t, posts, 0)
		assert.Equal(t, 0, total)
	})
}

// TestPostService_GetPostCountByUserID 测试获取用户文章数量
func TestPostService_GetPostCountByUserID(t *testing.T) {
	// 使用新测试框架
	ctx := framework.NewTestContext("TestPostService_GetPostCountByUserID", framework.WithDatabase())
	defer ctx.Cleanup()

	// 创建真实的服务
	config := config.GetConfig()
	userRepo := repository.NewUserRepository()
	taskRepo := repository.NewTaskRepository()
	taskService := task.NewService(taskRepo, config)
	userService := user.NewService(userRepo, taskService, config)
	diaryRepo := repository.NewDiaryRepository()
		diaryService := diary.NewService(diaryRepo)
		postRepo := repository.NewPostRepository(ctx.GetDatabase().GetDB())
		postService := post.NewService(userRepo, taskService, diaryRepo, postRepo)

	// 设置测试数据
		testUser, testDiary := setupTestUserAndDiary(t, ctx, userService, diaryService)

	t.Run("初始文章数量为0", func(t *testing.T) {
		count, err := postService.GetPostCountByUserID(testUser.UserId, testDiary.DiaryID)
		require.NoError(t, err)
		assert.Equal(t, 0, count)
	})

	// 添加3篇文章
	for i := 1; i <= 3; i++ {
		post := &post.Post{
			Title:       fmt.Sprintf("文章%d", i),
			Body:        fmt.Sprintf("第%d篇文章内容", i),
			DiaryID:     testDiary.DiaryID,
			AddUserID:   testUser.UserId,
			AddUserName: testUser.UserName,
			Status:      1,
		}
		_, err := postService.AddPost(post)
		require.NoError(t, err)
	}

	t.Run("添加文章后的数量", func(t *testing.T) {
		count, err := postService.GetPostCountByUserID(testUser.UserId, testDiary.DiaryID)
		require.NoError(t, err)
		assert.Equal(t, 3, count)
	})

	t.Run("不存在的用户", func(t *testing.T) {
		count, err := postService.GetPostCountByUserID(99999, testDiary.DiaryID)
		require.NoError(t, err)
		assert.Equal(t, 0, count)
	})

	t.Run("不存在的日记本", func(t *testing.T) {
		count, err := postService.GetPostCountByUserID(testUser.UserId, 99999)
		require.NoError(t, err)
		assert.Equal(t, 0, count)
	})
}

// TestPostService_UpdatePostByID 测试根据ID更新文章
func TestPostService_UpdatePostByID(t *testing.T) {
	// 使用新测试框架
	ctx := framework.NewTestContext("TestPostService_UpdatePostByID", framework.WithDatabase())
	defer ctx.Cleanup()

	// 创建真实的服务
	config := config.GetConfig()
	userRepo := repository.NewUserRepository()
	taskRepo := repository.NewTaskRepository()
	taskService := task.NewService(taskRepo, config)
	userService := user.NewService(userRepo, taskService, config)
	diaryRepo := repository.NewDiaryRepository()
		diaryService := diary.NewService(diaryRepo)
		postRepo := repository.NewPostRepository(ctx.GetDatabase().GetDB())
	postService := post.NewService(userRepo, taskService, diaryRepo, postRepo)

	// 设置测试数据
	testUser, testDiary := setupTestUserAndDiary(t, ctx, userService, diaryService)

	// 添加一篇文章
	originalPost := &post.Post{
		Title:       "待更新的文章",
		Body:        "原始内容",
		Summary:     "原始摘要",
		DiaryID:     testDiary.DiaryID,
		AddUserID:   testUser.UserId,
		AddUserName: testUser.UserName,
		Status:      1,
	}
	postID, err := postService.AddPost(originalPost)
	require.NoError(t, err)

	t.Run("成功更新文章", func(t *testing.T) {
		updatedPost := &post.Post{
			PostID:        postID,
			Title:         "已更新的文章",
			Body:          "更新后的内容",
			Summary:       "更新后的摘要",
			DiaryID:       testDiary.DiaryID,
			AddUserID:     testUser.UserId,
			AddUserName:   testUser.UserName,
			UpdateUTCTime: 1642680300,
			Status:        1,
		}

		err := postService.UpdatePostByID(updatedPost)
		assert.NoError(t, err)

		// 验证更新是否成功
		retrievedPost, err := postService.GetPostByID(postID)
		require.NoError(t, err)
		assert.Equal(t, "已更新的文章", retrievedPost.Title)
		assert.Equal(t, "更新后的内容", retrievedPost.Body)
		assert.Equal(t, "更新后的摘要", retrievedPost.Summary)
	})

	t.Run("更新不存在的文章", func(t *testing.T) {
		nonExistentPost := &post.Post{
			PostID:      99999,
			Title:       "不存在",
			Body:        "测试",
			DiaryID:     testDiary.DiaryID,
			AddUserID:   testUser.UserId,
			AddUserName: testUser.UserName,
			Status:      1,
		}

		err := postService.UpdatePostByID(nonExistentPost)
		assert.Error(t, err)
	})
}

// TestPostService_DeletePostByID 测试删除文章
func TestPostService_DeletePostByID(t *testing.T) {
	// 使用新测试框架
	ctx := framework.NewTestContext("TestPostService_DeletePostByID", framework.WithDatabase())
	defer ctx.Cleanup()

	// 创建真实的服务
	config := config.GetConfig()
	userRepo := repository.NewUserRepository()
	taskRepo := repository.NewTaskRepository()
	taskService := task.NewService(taskRepo, config)
	userService := user.NewService(userRepo, taskService, config)
	diaryRepo := repository.NewDiaryRepository()
		diaryService := diary.NewService(diaryRepo)
		postRepo := repository.NewPostRepository(ctx.GetDatabase().GetDB())
	postService := post.NewService(userRepo, taskService, diaryRepo, postRepo)

	// 设置测试数据
	testUser, testDiary := setupTestUserAndDiary(t, ctx, userService, diaryService)

	// 添加一篇文章
	newPost := &post.Post{
		Title:       "待删除的文章",
		Body:        "这篇文章将被删除",
		Summary:     "删除测试",
		DiaryID:     testDiary.DiaryID,
		AddUserID:   testUser.UserId,
		AddUserName: testUser.UserName,
		Status:      1,
	}
	postID, err := postService.AddPost(newPost)
	require.NoError(t, err)

	// 验证文章存在
	retrievedPost, err := postService.GetPostByID(postID)
	require.NoError(t, err)
	assert.NotNil(t, retrievedPost)

	t.Run("成功删除文章", func(t *testing.T) {
		err := postService.DeletePostByID(postID)
		assert.NoError(t, err)

		// 验证文章已被删除
		deletedPost, err := postService.GetPostByID(postID)
		assert.Error(t, err)
		assert.Nil(t, deletedPost)
	})

	t.Run("删除不存在的文章", func(t *testing.T) {
		err := postService.DeletePostByID(99999)
		assert.Error(t, err)
	})

	t.Run("删除已删除的文章", func(t *testing.T) {
		err := postService.DeletePostByID(postID)
		assert.Error(t, err)
	})
}
