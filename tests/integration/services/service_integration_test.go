package services

import (
	"testing"

	"github.com/biuren/hebe/internal/domain/diary"
	"github.com/biuren/hebe/internal/domain/post"
	"github.com/biuren/hebe/internal/domain/user"
	framework "github.com/biuren/hebe/tests/framework/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestUserDiaryIntegration 测试用户和日记服务的集成
func TestUserDiaryIntegration(t *testing.T) {
	ctx := framework.NewTestContext("TestUserDiaryIntegration", framework.WithDatabase())
	defer ctx.Cleanup()

	// 创建真实的服务
	userService := user.NewService()
	diaryService := diary.NewService()

	t.Run("用户创建和日记本操作流程", func(t *testing.T) {
			// 1. 注册用户 - 使用随机数据
			testUser := &user.User{
				UserName: ctx.GetDataGenerator().GenerateUniqueString("user"),
				Email:    ctx.GetDataGenerator().GenerateUniqueEmail("user"),
				Pwd:      "password123",
				FullName: ctx.GetDataGenerator().GenerateUniqueString("fullname"),
				TimeZone: "Asia/Shanghai",
				Lang:     "zh-cn",
			}
			err := userService.RegisterUser(testUser)
			require.NoError(t, err)
			require.True(t, testUser.UserId > 0)

			// 2. 验证用户登录
			authenticatedUser, err := userService.LoginWithUsernameOrEmail(testUser.UserName, "password123")
			require.NoError(t, err)
			assert.Equal(t, testUser.UserId, authenticatedUser.UserId)

			// 3. 为用户创建日记本 - 使用随机数据
			testDiary := &diary.Diary{
				Title:         ctx.GetDataGenerator().GenerateUniqueString("diary"),
				Remark:        "用于集成测试的日记本",
				AddUserID:     testUser.UserId,
				AddUserName:   testUser.UserName,
				AddUTCTime:    1642680000,
				Status:        1,
				ViewSetting:   1,
				SocialSetting: 1,
				ThemeID:       1,
			}
			diaryID, err := diaryService.AddDiary(testDiary)
			require.NoError(t, err)
			assert.True(t, diaryID > 0)

			// 4. 验证用户可以访问自己的日记本
			userDiaries, err := diaryService.ListDiariesForWrite(testUser.UserId)
			require.NoError(t, err)
			assert.Len(t, userDiaries, 1)
			assert.Equal(t, testDiary.Title, userDiaries[0].Title)

			// 5. 检查日记本数量限制
			err = diaryService.CheckDiaryLimit(testUser.UserId)
			assert.NoError(t, err)

			// 6. 更新用户资料
			err = userService.UpdateUserProfile(testUser.UserId, testUser.UserId, "Updated Name", "Updated Remark")
			assert.NoError(t, err)

			// 7. 验证用户信息更新
			updatedUser, err := userService.GetUserByUserId(testUser.UserId)
			require.NoError(t, err)
			assert.Equal(t, "Updated Name", updatedUser.FullName)
			assert.Equal(t, "Updated Remark", updatedUser.Remark)
		})
}

// TestUserPostDiaryIntegration 测试用户、日记和文章服务的完整集成
func TestUserPostDiaryIntegration(t *testing.T) {
	ctx := framework.NewTestContext("TestUserPostDiaryIntegration", framework.WithDatabase())
	defer ctx.Cleanup()

	// 创建真实的服务
	userService := user.NewService()
	diaryService := diary.NewService()
	postService := post.NewService()

	t.Run("完整的内容创建流程", func(t *testing.T) {
			// 1. 创建用户 - 使用随机数据
			testUser := &user.User{
				UserName: ctx.GetDataGenerator().GenerateUniqueString("user"),
				Email:    ctx.GetDataGenerator().GenerateUniqueEmail("user"),
				Pwd:      "password123",
				FullName: ctx.GetDataGenerator().GenerateUniqueString("fullname"),
				TimeZone: "UTC",
				Lang:     "en",
			}
			err := userService.RegisterUser(testUser)
			require.NoError(t, err)

			// 2. 创建日记本 - 使用随机数据
			testDiary := &diary.Diary{
				Title:        ctx.GetDataGenerator().GenerateUniqueString("diary"),
				Remark:       "用于内容测试",
				AddUserID:    testUser.UserId,
				AddUserName:  testUser.UserName,
				AddUTCTime:   1642680000,
				Status:       1,
				ViewSetting:  1,
				SocialSetting: 1,
				ThemeID:      1,
			}
			diaryID, err := diaryService.AddDiary(testDiary)
			require.NoError(t, err)

			// 3. 在日记本中创建文章 - 使用随机数据
			testPost := &post.Post{
				Title:       ctx.GetDataGenerator().GenerateUniqueString("post"),
				Body:        ctx.GetDataGenerator().GenerateUniqueString("content"),
				DiaryID:     diaryID,
				AddUserID:   testUser.UserId,
				AddUserName: testUser.UserName,
				AddUTCTime:  1642680100,
				Status:      1,
			}
			postID, err := postService.AddPost(testPost)
			require.NoError(t, err)
			assert.True(t, postID > 0)

			// 更新日记本的文章计数
			err = diaryService.IncrementPostCount(diaryID, testUser.UserId, testUser.UserName, testPost.AddUTCTime)
			require.NoError(t, err)

			// 4. 验证文章计数
			postCount, err := postService.GetPostCountByUserID(testUser.UserId, diaryID)
			require.NoError(t, err)
			assert.Equal(t, 1, postCount)

			// 5. 验证日记本的文章计数已更新
			updatedDiary, err := diaryService.GetDiaryByID(diaryID)
			require.NoError(t, err)
			assert.Equal(t, int32(1), updatedDiary.InfoNum)

			// 6. 创建更多文章
			for i := 2; i <= 3; i++ {
				additionalPost := &post.Post{
					Title:       "第" + string(rune('0'+i)) + "篇文章",
					Body:        "这是第" + string(rune('0'+i)) + "篇文章的内容",
					Summary:     "第" + string(rune('0'+i)) + "篇文章摘要",
					DiaryID:     diaryID,
					AddUserID:   testUser.UserId,
					AddUserName: testUser.UserName,
					AddUTCTime:  int64(1642680100 + i),
					Status:      1,
				}
				_, err := postService.AddPost(additionalPost)
				require.NoError(t, err)

				// 更新日记本的文章计数
				err = diaryService.IncrementPostCount(diaryID, testUser.UserId, testUser.UserName, additionalPost.AddUTCTime)
				require.NoError(t, err)
			}

			// 7. 验证分页获取文章
			posts, total, err := postService.ListPostsByUserID(testUser.UserId, diaryID, 1, 2)
			require.NoError(t, err)
			assert.Len(t, posts, 2)   // 第一页2篇
			assert.Equal(t, 3, total) // 总共3篇

			// 8. 获取第二页
			posts2, total2, err := postService.ListPostsByUserID(testUser.UserId, diaryID, 2, 2)
			require.NoError(t, err)
			assert.Len(t, posts2, 1)   // 第二页1篇
			assert.Equal(t, 3, total2) // 总共还是3篇

			// 9. 更新文章
			firstPost := &posts[0]
			firstPost.Title = "更新后的标题"
			firstPost.Body = "更新后的内容"
			firstPost.UpdateUTCTime = 1642680300
			err = postService.UpdatePost(firstPost)
			assert.NoError(t, err)

			// 10. 验证文章更新
			updatedPost, err := postService.GetPostByID(firstPost.PostID)
			require.NoError(t, err)
			assert.Equal(t, "更新后的标题", updatedPost.Title)
			assert.Equal(t, "更新后的内容", updatedPost.Body)

			// 11. 删除一篇文章
			err = postService.DeletePostByID(firstPost.PostID)
			assert.NoError(t, err)

			// 12. 验证文章已删除
			deletedPost, err := postService.GetPostByID(firstPost.PostID)
			assert.Error(t, err)
			assert.Nil(t, deletedPost)

			// 13. 验证文章计数已更新
			finalCount, err := postService.GetPostCountByUserID(testUser.UserId, diaryID)
			require.NoError(t, err)
			assert.Equal(t, 2, finalCount) // 删除一篇后剩2篇
		})
}

// TestMultiUserScenario 测试多用户场景
func TestMultiUserScenario(t *testing.T) {
	ctx := framework.NewTestContext("TestMultiUserScenario", framework.WithDatabase())
	defer ctx.Cleanup()

	// 创建真实的服务
	userService := user.NewService()
	diaryService := diary.NewService()
	postService := post.NewService()

	t.Run("多用户数据隔离", func(t *testing.T) {
			// 创建两个用户
			user1 := &user.User{
				UserName: ctx.GetDataGenerator().GenerateUniqueString("testuser1"),
				Email:    ctx.GetDataGenerator().GenerateUniqueEmail("user1"),
				Pwd:      "password123",
				TimeZone: "Asia/Shanghai",
			}
			err := userService.RegisterUser(user1)
			require.NoError(t, err)

			user2 := &user.User{
				UserName: ctx.GetDataGenerator().GenerateUniqueString("testuser2"),
				Email:    ctx.GetDataGenerator().GenerateUniqueEmail("user2"),
				Pwd:      "password123",
				TimeZone: "UTC",
			}
			err = userService.RegisterUser(user2)
			require.NoError(t, err)

			// 为每个用户创建日记本
			diary1 := &diary.Diary{
				Title:         "用户1的日记本",
				Remark:        "用户1专用",
				AddUserID:     user1.UserId,
				AddUserName:   user1.UserName,
				AddUTCTime:    1642680000,
				Status:        1,
				ViewSetting:   1,
				SocialSetting: 1,
				ThemeID:       1,
			}
			diary1ID, err := diaryService.AddDiary(diary1)
			require.NoError(t, err)

			diary2 := &diary.Diary{
				Title:         "用户2的日记本",
				Remark:        "用户2专用",
				AddUserID:     user2.UserId,
				AddUserName:   user2.UserName,
				AddUTCTime:    1642680000,
				Status:        1,
				ViewSetting:   1,
				SocialSetting: 1,
				ThemeID:       1,
			}
			diary2ID, err := diaryService.AddDiary(diary2)
			require.NoError(t, err)

			// 为每个用户在各自的日记本中创建文章
			post1 := &post.Post{
				Title:       "用户1的文章",
				Body:        "用户1的内容",
				DiaryID:     diary1ID,
				AddUserID:   user1.UserId,
				AddUserName: user1.UserName,
				AddUTCTime:  1642680100,
				Status:      1,
			}
			_, err = postService.AddPost(post1)
			require.NoError(t, err)

			post2 := &post.Post{
				Title:       "用户2的文章",
				Body:        "用户2的内容",
				DiaryID:     diary2ID,
				AddUserID:   user2.UserId,
				AddUserName: user2.UserName,
				AddUTCTime:  1642680100,
				Status:      1,
			}
			_, err = postService.AddPost(post2)
			require.NoError(t, err)

			// 验证数据隔离：用户1只能看到自己的日记本
			user1Diaries, err := diaryService.ListDiariesForWrite(user1.UserId)
			require.NoError(t, err)
			assert.Len(t, user1Diaries, 1)
			assert.Equal(t, "用户1的日记本", user1Diaries[0].Title)

			// 验证数据隔离：用户2只能看到自己的日记本
			user2Diaries, err := diaryService.ListDiariesForWrite(user2.UserId)
			require.NoError(t, err)
			assert.Len(t, user2Diaries, 1)
			assert.Equal(t, "用户2的日记本", user2Diaries[0].Title)

			// 验证文章数据隔离
			user1Posts, _, err := postService.ListPostsByUserID(user1.UserId, diary1ID, 1, 10)
			require.NoError(t, err)
			assert.Len(t, user1Posts, 1)
			assert.Equal(t, "用户1的文章", user1Posts[0].Title)

			user2Posts, _, err := postService.ListPostsByUserID(user2.UserId, diary2ID, 1, 10)
			require.NoError(t, err)
			assert.Len(t, user2Posts, 1)
			assert.Equal(t, "用户2的文章", user2Posts[0].Title)

			// 验证用户不能访问其他用户的数据
			crossUserPosts, _, err := postService.ListPostsByUserID(user1.UserId, diary2ID, 1, 10)
			require.NoError(t, err)
			assert.Len(t, crossUserPosts, 0) // 用户1不能访问用户2的日记本中的文章
		})
}

// TestServiceErrorHandling 测试服务错误处理
func TestServiceErrorHandling(t *testing.T) {
	ctx := framework.NewTestContext("TestServiceErrorHandling", framework.WithDatabase())
	defer ctx.Cleanup()

	// 创建真实的服务
	userService := user.NewService()
	diaryService := diary.NewService()
	postService := post.NewService()

	t.Run("服务间错误传播", func(t *testing.T) {
			// 1. 尝试为不存在的用户创建日记本
			invalidDiary := &diary.Diary{
				Title:        "无效日记本",
				Remark:       "测试",
				AddUserID:    99999, // 不存在的用户ID
				AddUserName:  "nonexistent",
				AddUTCTime:   1642680000,
				Status:       1,
				ViewSetting:  1,
				SocialSetting: 1,
				ThemeID:      1,
			}
			_, err := diaryService.AddDiary(invalidDiary)
			// 注意：系统设计上允许为不存在的用户创建日记本，所以这里不会出错
			assert.NoError(t, err)

			// 2. 尝试在不存在的日记本中创建文章
			invalidPost := &post.Post{
				Title:       "无效文章",
				Body:        "测试内容",
				DiaryID:     99999, // 不存在的日记本ID
				AddUserID:   1,
				AddUserName: "testuser",
				AddUTCTime:  1642680100,
				Status:      1,
			}
			_, err = postService.AddPost(invalidPost)
			assert.Error(t, err) // 应该出错

			// 3. 尝试重复注册相同的用户名
			duplicateUserName := ctx.GetDataGenerator().GenerateUniqueString("duplicateuser")
			user1 := &user.User{
				UserName: duplicateUserName,
				Email:    ctx.GetDataGenerator().GenerateUniqueEmail("dup1"),
				Pwd:      "password123",
				TimeZone: "UTC",
			}
			err = userService.RegisterUser(user1)
			require.NoError(t, err)

			user2 := &user.User{
				UserName: duplicateUserName, // 重复的用户名
				Email:    ctx.GetDataGenerator().GenerateUniqueEmail("dup2"),
				Pwd:      "password123",
				TimeZone: "UTC",
			}
			err = userService.RegisterUser(user2)
			assert.Error(t, err) // 应该出错

			// 4. 尝试重复注册相同的邮箱
			user3 := &user.User{
				UserName: "anotheruser",
				Email:    "dup1@example.com", // 重复的邮箱
				Pwd:      "password123",
				TimeZone: "UTC",
			}
			err = userService.RegisterUser(user3)
			assert.Error(t, err) // 应该出错
		})
}

// TestServicePerformance 测试服务性能
func TestServicePerformance(t *testing.T) {
	ctx := framework.NewTestContext("TestServicePerformance", framework.WithDatabase())
	defer ctx.Cleanup()

	// 创建真实的服务
	userService := user.NewService()
	diaryService := diary.NewService()
	postService := post.NewService()

	t.Run("批量操作性能", func(t *testing.T) {
			// 创建一个用户
			testUser := &user.User{
				UserName: ctx.GetDataGenerator().GenerateUniqueString("perfuser"),
				Email:    ctx.GetDataGenerator().GenerateUniqueEmail("perf"),
				Pwd:      "password123",
				TimeZone: "UTC",
			}
			err := userService.RegisterUser(testUser)
			require.NoError(t, err)

			// 创建一个日记本
			testDiary := &diary.Diary{
				Title:         "性能测试日记本",
				Remark:        "用于性能测试",
				AddUserID:     testUser.UserId,
				AddUserName:   testUser.UserName,
				AddUTCTime:    1642680000,
				Status:        1,
				ViewSetting:   1,
				SocialSetting: 1,
				ThemeID:       1,
			}
			diaryID, err := diaryService.AddDiary(testDiary)
			require.NoError(t, err)

			// 批量创建文章（小批量，避免测试时间过长）
			const articleCount = 10
			for i := 1; i <= articleCount; i++ {
				article := &post.Post{
					Title:       "性能测试文章" + string(rune('0'+i)),
					Body:        "这是第" + string(rune('0'+i)) + "篇性能测试文章",
					Summary:     "摘要" + string(rune('0'+i)),
					DiaryID:     diaryID,
					AddUserID:   testUser.UserId,
					AddUserName: testUser.UserName,
					AddUTCTime:  int64(1642680100 + i),
					Status:      1,
				}
				_, err := postService.AddPost(article)
				require.NoError(t, err)
			}

			// 验证分页查询性能
			posts, total, err := postService.ListPostsByUserID(testUser.UserId, diaryID, 1, 5)
			require.NoError(t, err)
			assert.Len(t, posts, 5)
			assert.Equal(t, articleCount, total)

			// 验证文章计数
			count, err := postService.GetPostCountByUserID(testUser.UserId, diaryID)
			require.NoError(t, err)
			assert.Equal(t, articleCount, count)
		})
}
