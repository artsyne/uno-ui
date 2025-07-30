package diary_test

import (
	"fmt"
	"testing"

	"github.com/biuren/hebe/internal/domain/diary"
	"github.com/biuren/hebe/internal/domain/user"
	testhelpers "github.com/biuren/hebe/tests/helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDiaryService_AddDiary 测试添加日记本功能
func TestDiaryService_AddDiary(t *testing.T) {
	testhelpers.WithDatabaseTest(t, func(db *testhelpers.DatabaseTestHelper) {
		// 开始事务，确保测试数据隔离
		db.BeginTransaction(t)
		defer db.RollbackTransaction(t)

		// 创建真实的服务
		userService := user.NewService()
		diaryService := diary.NewService()

		// 先创建一个用户 - 使用随机数据
		testUser := &user.User{
			UserName: testhelpers.GenerateRandomUsername(),
			Email:    testhelpers.GenerateRandomEmail(),
			Pwd:      "password123",
			TimeZone: "Asia/Shanghai",
		}
		err := userService.RegisterUser(testUser)
		require.NoError(t, err)

		t.Run("成功添加日记本", func(t *testing.T) {
			newDiary := &diary.Diary{
				Title:       testhelpers.GenerateRandomDiaryTitle(),
				Remark:      "这是一个测试日记本",
				AddUserID:   testUser.UserId,
				AddUserName: testUser.UserName,
				AddUTCTime:  1642680000,
				Status:      1,
			}

			diaryID, err := diaryService.AddDiary(newDiary)
			require.NoError(t, err)
			assert.True(t, diaryID > 0)
			assert.Equal(t, diaryID, newDiary.DiaryID)
		})

		t.Run("标题为空", func(t *testing.T) {
			invalidDiary := &diary.Diary{
				Title:       "", // 空标题
				Remark:      "测试备注",
				AddUserID:   testUser.UserId,
				AddUserName: testUser.UserName,
				AddUTCTime:  1642680000,
				Status:      1,
			}

			_, err := diaryService.AddDiary(invalidDiary)
			assert.Error(t, err)
		})

		t.Run("用户ID无效", func(t *testing.T) {
			invalidDiary := &diary.Diary{
				Title:       testhelpers.GenerateRandomDiaryTitle(),
				Remark:      "测试备注",
				AddUserID:   0, // 无效用户ID
				AddUserName: testhelpers.GenerateRandomUsername(),
				AddUTCTime:  1642680000,
				Status:      1,
			}

			_, err := diaryService.AddDiary(invalidDiary)
			assert.Error(t, err)
		})

		t.Run("检查日记本数量限制", func(t *testing.T) {
			// 首先验证限制检查功能
			err := diaryService.CheckDiaryLimit(testUser.UserId)
			assert.NoError(t, err) // 当前没有日记本，应该不会超限

			// 添加一个日记本
			testDiary := &diary.Diary{
				Title:       testhelpers.GenerateRandomDiaryTitle(),
				Remark:      "测试",
				AddUserID:   testUser.UserId,
				AddUserName: testUser.UserName,
				AddUTCTime:  1642680000,
				Status:      1,
			}
			_, err = diaryService.AddDiary(testDiary)
			require.NoError(t, err)

			// 再次检查限制
			err = diaryService.CheckDiaryLimit(testUser.UserId)
			assert.NoError(t, err) // 应该还没超限
		})
	})
}

// TestDiaryService_GetDiaryByID 测试根据ID获取日记本
func TestDiaryService_GetDiaryByID(t *testing.T) {
	testhelpers.WithDatabaseTest(t, func(db *testhelpers.DatabaseTestHelper) {
		// 开始事务，确保测试数据隔离
		db.BeginTransaction(t)
		defer db.RollbackTransaction(t)

		// 创建真实的服务
		userService := user.NewService()
		diaryService := diary.NewService()

		// 先创建一个用户
		testUser := &user.User{
			UserName: testhelpers.GenerateRandomUsername(),
			Email:    testhelpers.GenerateRandomEmail(),
			Pwd:      "password123",
			TimeZone: "UTC",
		}
		err := userService.RegisterUser(testUser)
		require.NoError(t, err)

		// 添加一个日记本
		newDiary := &diary.Diary{
			Title:       testhelpers.GenerateRandomDiaryTitle(),
			Remark:      "用于测试获取功能",
			AddUserID:   testUser.UserId,
			AddUserName: testUser.UserName,
			AddUTCTime:  1642680000,
			Status:      1,
		}
		diaryID, err := diaryService.AddDiary(newDiary)
		require.NoError(t, err)

		t.Run("成功获取日记本", func(t *testing.T) {
			retrievedDiary, err := diaryService.GetDiaryByID(diaryID)
			require.NoError(t, err)
			assert.NotNil(t, retrievedDiary)
			assert.Equal(t, newDiary.Title, retrievedDiary.Title)
			assert.Equal(t, "用于测试获取功能", retrievedDiary.Remark)
			assert.Equal(t, testUser.UserId, retrievedDiary.AddUserID)
		})

		t.Run("日记本不存在", func(t *testing.T) {
			retrievedDiary, err := diaryService.GetDiaryByID(99999)
			assert.Error(t, err)
			assert.Nil(t, retrievedDiary)
		})
	})
}

// TestDiaryService_UpdateDiary 测试更新日记本
func TestDiaryService_UpdateDiary(t *testing.T) {
	testhelpers.WithDatabaseTest(t, func(db *testhelpers.DatabaseTestHelper) {
		// 开始事务，确保测试数据隔离
		db.BeginTransaction(t)
		defer db.RollbackTransaction(t)

		// 创建真实的服务
		userService := user.NewService()
		diaryService := diary.NewService()

		// 先创建一个用户
		testUser := &user.User{
			UserName: testhelpers.GenerateRandomUsername(),
			Email:    testhelpers.GenerateRandomEmail(),
			Pwd:      "password123",
			TimeZone: "Asia/Shanghai",
		}
		err := userService.RegisterUser(testUser)
		require.NoError(t, err)

		// 添加一个日记本
		originalDiary := &diary.Diary{
			Title:       testhelpers.GenerateRandomDiaryTitle(),
			Remark:      "原始备注",
			AddUserID:   testUser.UserId,
			AddUserName: testUser.UserName,
			AddUTCTime:  1642680000,
			Status:      1,
		}
		diaryID, err := diaryService.AddDiary(originalDiary)
		require.NoError(t, err)

		t.Run("成功更新日记本", func(t *testing.T) {
			updatedTitle := testhelpers.GenerateRandomDiaryTitle()
			updatedDiary := &diary.Diary{
				DiaryID:     diaryID,
				Title:       updatedTitle,
				Remark:      "更新后的备注",
				AddUserID:   testUser.UserId,
				AddUserName: testUser.UserName,
				AddUTCTime:  1642680000,
				Status:      1,
			}

			err := diaryService.UpdateDiary(updatedDiary)
			assert.NoError(t, err)

			// 验证更新是否成功
			retrievedDiary, err := diaryService.GetDiaryByID(diaryID)
			require.NoError(t, err)
			assert.Equal(t, updatedTitle, retrievedDiary.Title)
			assert.Equal(t, "更新后的备注", retrievedDiary.Remark)
		})

		t.Run("更新不存在的日记本", func(t *testing.T) {
			nonExistentDiary := &diary.Diary{
				DiaryID:     99999,
				Title:       testhelpers.GenerateRandomDiaryTitle(),
				Remark:      "测试",
				AddUserID:   testUser.UserId,
				AddUserName: testUser.UserName,
				AddUTCTime:  1642680000,
				Status:      1,
			}

			err := diaryService.UpdateDiary(nonExistentDiary)
			assert.Error(t, err)
		})

		t.Run("标题为空时更新失败", func(t *testing.T) {
			invalidDiary := &diary.Diary{
				DiaryID:     diaryID,
				Title:       "", // 空标题
				Remark:      "备注",
				AddUserID:   testUser.UserId,
				AddUserName: testUser.UserName,
				AddUTCTime:  1642680000,
				Status:      1,
			}

			err := diaryService.UpdateDiary(invalidDiary)
			assert.Error(t, err)
		})
	})
}

// TestDiaryService_ListDiariesForWrite 测试获取可写入的日记本列表
func TestDiaryService_ListDiariesForWrite(t *testing.T) {
	testhelpers.WithDatabaseTest(t, func(db *testhelpers.DatabaseTestHelper) {
		// 开始事务，确保测试数据隔离
		db.BeginTransaction(t)
		defer db.RollbackTransaction(t)

		// 创建真实的服务
		userService := user.NewService()
		diaryService := diary.NewService()

		// 创建两个用户
		user1 := &user.User{
			UserName: testhelpers.GenerateRandomUsername(),
			Email:    testhelpers.GenerateRandomEmail(),
			Pwd:      "password123",
			TimeZone: "UTC",
		}
		err := userService.RegisterUser(user1)
		require.NoError(t, err)

		user2 := &user.User{
			UserName: testhelpers.GenerateRandomUsername(),
			Email:    testhelpers.GenerateRandomEmail(),
			Pwd:      "password123",
			TimeZone: "UTC",
		}
		err = userService.RegisterUser(user2)
		require.NoError(t, err)

		// 为用户1添加两个日记本
		diary1 := &diary.Diary{
			Title:       testhelpers.GenerateRandomDiaryTitle(),
			Remark:      "测试",
			AddUserID:   user1.UserId,
			AddUserName: user1.UserName,
			AddUTCTime:  1642680000,
			Status:      1,
		}
		_, err = diaryService.AddDiary(diary1)
		require.NoError(t, err)

		diary2 := &diary.Diary{
			Title:       testhelpers.GenerateRandomDiaryTitle(),
			Remark:      "测试",
			AddUserID:   user1.UserId,
			AddUserName: user1.UserName,
			AddUTCTime:  1642680001,
			Status:      1,
		}
		_, err = diaryService.AddDiary(diary2)
		require.NoError(t, err)

		// 为用户2添加一个日记本
		diary3 := &diary.Diary{
			Title:       testhelpers.GenerateRandomDiaryTitle(),
			Remark:      "测试",
			AddUserID:   user2.UserId,
			AddUserName: user2.UserName,
			AddUTCTime:  1642680002,
			Status:      1,
		}
		_, err = diaryService.AddDiary(diary3)
		require.NoError(t, err)

		t.Run("获取用户1的日记本列表", func(t *testing.T) {
			diaries, err := diaryService.ListDiariesForWrite(user1.UserId)
			require.NoError(t, err)
			assert.Len(t, diaries, 2)

			// 验证返回的都是用户1的日记本
			for _, d := range diaries {
				assert.Equal(t, user1.UserId, d.AddUserID)
			}
		})

		t.Run("获取用户2的日记本列表", func(t *testing.T) {
			diaries, err := diaryService.ListDiariesForWrite(user2.UserId)
			require.NoError(t, err)
			assert.Len(t, diaries, 1)
			assert.Equal(t, diary3.Title, diaries[0].Title)
		})

		t.Run("获取不存在用户的日记本列表", func(t *testing.T) {
			diaries, err := diaryService.ListDiariesForWrite(99999)
			assert.NoError(t, err) // 应该返回空列表而不是错误
			assert.Len(t, diaries, 0)
		})
	})
}

// TestDiaryService_ListDiariesByUser 测试分页获取用户日记本列表
func TestDiaryService_ListDiariesByUser(t *testing.T) {
	testhelpers.WithDatabaseTest(t, func(db *testhelpers.DatabaseTestHelper) {
		// 开始事务，确保测试数据隔离
		db.BeginTransaction(t)
		defer db.RollbackTransaction(t)

		// 创建真实的服务
		userService := user.NewService()
		diaryService := diary.NewService()

		// 创建一个用户
		testUser := &user.User{
			UserName: testhelpers.GenerateRandomUsername(),
			Email:    testhelpers.GenerateRandomEmail(),
			Pwd:      "password123",
			TimeZone: "Asia/Shanghai",
		}
		err := userService.RegisterUser(testUser)
		require.NoError(t, err)

		// 添加5个日记本
		for i := 1; i <= 5; i++ {
			testDiary := &diary.Diary{
				Title:       testhelpers.GenerateRandomDiaryTitle(),
				Remark:      fmt.Sprintf("第%d个日记本", i),
				AddUserID:   testUser.UserId,
				AddUserName: testUser.UserName,
				AddUTCTime:  int64(1642680000 + i),
				Status:      1,
			}
			_, err = diaryService.AddDiary(testDiary)
			require.NoError(t, err)
		}

		t.Run("第一页，每页3个", func(t *testing.T) {
			diaries, total, err := diaryService.ListDiariesByUser(testUser.UserId, 1, 3)
			require.NoError(t, err)
			assert.Len(t, diaries, 3)
			assert.Equal(t, 5, total)
		})

		t.Run("第二页，每页3个", func(t *testing.T) {
			diaries, total, err := diaryService.ListDiariesByUser(testUser.UserId, 2, 3)
			require.NoError(t, err)
			assert.Len(t, diaries, 2) // 剩余2个
			assert.Equal(t, 5, total)
		})

		t.Run("获取所有", func(t *testing.T) {
			diaries, total, err := diaryService.ListDiariesByUser(testUser.UserId, 1, 10)
			require.NoError(t, err)
			assert.Len(t, diaries, 5)
			assert.Equal(t, 5, total)
		})

		t.Run("用户不存在", func(t *testing.T) {
			diaries, total, err := diaryService.ListDiariesByUser(99999, 1, 10)
			assert.NoError(t, err) // 应该返回空列表而不是错误
			assert.Len(t, diaries, 0)
			assert.Equal(t, 0, total)
		})
	})
}

// TestDiaryService_PostCountOperations 测试文章数量操作
func TestDiaryService_PostCountOperations(t *testing.T) {
	testhelpers.WithDatabaseTest(t, func(db *testhelpers.DatabaseTestHelper) {
		// 开始事务，确保测试数据隔离
		db.BeginTransaction(t)
		defer db.RollbackTransaction(t)

		// 创建真实的服务
		userService := user.NewService()
		diaryService := diary.NewService()

		// 创建一个用户
		testUser := &user.User{
			UserName: testhelpers.GenerateRandomUsername(),
			Email:    testhelpers.GenerateRandomEmail(),
			Pwd:      "password123",
			TimeZone: "UTC",
		}
		err := userService.RegisterUser(testUser)
		require.NoError(t, err)

		// 添加一个日记本
		testDiary := &diary.Diary{
			Title:       testhelpers.GenerateRandomDiaryTitle(),
			Remark:      "测试文章数量操作",
			AddUserID:   testUser.UserId,
			AddUserName: testUser.UserName,
			AddUTCTime:  1642680000,
			Status:      1,
		}
		diaryID, err := diaryService.AddDiary(testDiary)
		require.NoError(t, err)

		// 获取初始日记本状态
		initialDiary, err := diaryService.GetDiaryByID(diaryID)
		require.NoError(t, err)
		initialCount := initialDiary.InfoNum

		t.Run("增加文章计数", func(t *testing.T) {
			actionTime := int64(1642680100)
			err := diaryService.IncrementPostCount(diaryID, testUser.UserId, testUser.UserName, actionTime)
			assert.NoError(t, err)

			// 验证计数是否增加
			updatedDiary, err := diaryService.GetDiaryByID(diaryID)
			require.NoError(t, err)
			assert.Equal(t, initialCount+1, updatedDiary.InfoNum)
			assert.Equal(t, testUser.UserId, updatedDiary.LastActionUserID)
			assert.Equal(t, testUser.UserName, updatedDiary.LastActionUserName)
			assert.Equal(t, actionTime, updatedDiary.LastActionUTCTime)
		})

		t.Run("减少文章计数", func(t *testing.T) {
			actionTime := int64(1642680200)

			// 获取当前计数
			currentDiary, err := diaryService.GetDiaryByID(diaryID)
			require.NoError(t, err)
			currentCount := currentDiary.InfoNum

			err = diaryService.DecrementPostCount(diaryID, testUser.UserId, testUser.UserName, actionTime)
			assert.NoError(t, err)

			// 验证计数是否减少
			updatedDiary, err := diaryService.GetDiaryByID(diaryID)
			require.NoError(t, err)
			assert.Equal(t, currentCount-1, updatedDiary.InfoNum)
			assert.Equal(t, actionTime, updatedDiary.LastActionUTCTime)
		})

		t.Run("对不存在的日记本操作", func(t *testing.T) {
			actionTime := int64(1642680300)

			err := diaryService.IncrementPostCount(99999, testUser.UserId, testUser.UserName, actionTime)
			assert.Error(t, err)

			err = diaryService.DecrementPostCount(99999, testUser.UserId, testUser.UserName, actionTime)
			assert.Error(t, err)
		})

		t.Run("计数不会变为负数", func(t *testing.T) {
			// 先重置日记本计数为0（通过减少足够多次）
			currentDiary, err := diaryService.GetDiaryByID(diaryID)
			require.NoError(t, err)

			// 尝试减少计数，即使当前计数已经很低
			actionTime := int64(1642680400)
			for i := 0; i < int(currentDiary.InfoNum)+5; i++ { // 减少超过当前计数的次数
				err = diaryService.DecrementPostCount(diaryID, testUser.UserId, testUser.UserName, actionTime)
				// 可能有些减少操作会成功，有些可能失败，这取决于具体实现
			}

			// 验证计数没有变为负数
			finalDiary, err := diaryService.GetDiaryByID(diaryID)
			require.NoError(t, err)
			assert.True(t, finalDiary.InfoNum >= 0, "文章计数不应该为负数")
		})
	})
}
