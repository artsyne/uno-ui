package benchmark

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/biuren/hebe/internal/domain/diary"
	"github.com/biuren/hebe/internal/domain/post"
	"github.com/biuren/hebe/internal/domain/task"
	"github.com/biuren/hebe/internal/domain/user"
	"github.com/biuren/hebe/internal/infra/repository"
	framework "github.com/biuren/hebe/tests/framework/core"
)

// BenchmarkUserRegistration 用户注册性能基准测试
func BenchmarkUserRegistration(b *testing.B) {
	ctx := framework.NewTestContext("BenchmarkUserRegistration", framework.WithDatabase())
	defer ctx.Cleanup()

	userRepo := repository.NewUserRepository()
	taskRepo := repository.NewTaskRepository()
	taskService := task.NewService(taskRepo)
	userService := user.NewService(userRepo, taskService)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		testUser := &user.User{
			UserName: ctx.GetDataGenerator().GenerateUniqueString("benchuser"),
			Email:    ctx.GetDataGenerator().GenerateUniqueEmail("bench"),
			Pwd:      "password123",
			TimeZone: "UTC",
			Lang:     "en",
		}
		err := userService.RegisterUser(testUser)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkDiaryCreation 日记本创建性能基准测试
func BenchmarkDiaryCreation(b *testing.B) {
	ctx := framework.NewTestContext("BenchmarkDiaryCreation", framework.WithDatabase())
	defer ctx.Cleanup()

	userRepo := repository.NewUserRepository()
	taskRepo := repository.NewTaskRepository()
	taskService := task.NewService(taskRepo)
	userService := user.NewService(userRepo, taskService)
	diaryRepo := repository.NewDiaryRepository()
	diaryService := diary.NewService(diaryRepo)

	// 创建测试用户
	testUser := &user.User{
		UserName: ctx.GetDataGenerator().GenerateUniqueString("benchuser"),
		Email:    ctx.GetDataGenerator().GenerateUniqueEmail("bench"),
		Pwd:      "password123",
		TimeZone: "UTC",
		Lang:     "en",
	}
	err := userService.RegisterUser(testUser)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		testDiary := &diary.Diary{
			Title:         ctx.GetDataGenerator().GenerateUniqueString("benchdiary"),
			Remark:        "性能测试日记本",
			AddUserID:     testUser.UserId,
			AddUserName:   testUser.UserName,
			AddUTCTime:    time.Now().Unix(),
			Status:        1,
			ViewSetting:   1,
			SocialSetting: 1,
			ThemeID:       1,
		}
		_, err := diaryService.AddDiary(testDiary)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkPostCreation 文章创建性能基准测试
func BenchmarkPostCreation(b *testing.B) {
	ctx := framework.NewTestContext("BenchmarkPostCreation", framework.WithDatabase())
	defer ctx.Cleanup()

	userRepo := repository.NewUserRepository()
	taskRepo := repository.NewTaskRepository()
	taskService := task.NewService(taskRepo)
	userService := user.NewService(userRepo, taskService)
	diaryRepo := repository.NewDiaryRepository()
	diaryService := diary.NewService(diaryRepo)
	postRepo := repository.NewPostRepository(ctx.GetDatabase().GetDB())
	postService := post.NewService(userRepo, taskService, diaryRepo, postRepo)

	// 创建测试用户
	testUser := &user.User{
		UserName: ctx.GetDataGenerator().GenerateUniqueString("benchuser"),
		Email:    ctx.GetDataGenerator().GenerateUniqueEmail("bench"),
		Pwd:      "password123",
		TimeZone: "UTC",
		Lang:     "en",
	}
	err := userService.RegisterUser(testUser)
	if err != nil {
		b.Fatal(err)
	}

	// 创建测试日记本
	testDiary := &diary.Diary{
		Title:         ctx.GetDataGenerator().GenerateUniqueString("benchdiary"),
		Remark:        "性能测试日记本",
		AddUserID:     testUser.UserId,
		AddUserName:   testUser.UserName,
		AddUTCTime:    time.Now().Unix(),
		Status:        1,
		ViewSetting:   1,
		SocialSetting: 1,
		ThemeID:       1,
	}
	diaryID, err := diaryService.AddDiary(testDiary)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		testPost := &post.Post{
			Title:       ctx.GetDataGenerator().GenerateUniqueString("benchpost"),
			Body:        "这是一篇性能测试文章的内容",
			Summary:     "性能测试摘要",
			DiaryID:     diaryID,
			AddUserID:   testUser.UserId,
			AddUserName: testUser.UserName,
			AddUTCTime:  time.Now().Unix(),
			Status:      1,
			InfoType:    1,
		}
		_, err := postService.AddPost(testPost)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkDatabaseOperations 数据库操作性能基准测试
func BenchmarkDatabaseOperations(b *testing.B) {
	ctx := framework.NewTestContext("BenchmarkDatabaseOperations", framework.WithDatabase())
	defer ctx.Cleanup()

	db := ctx.GetDatabase().GetRawDB()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// 测试简单查询性能
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM br_user").Scan(&count)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkConcurrentUserRegistration 并发用户注册性能基准测试
func BenchmarkConcurrentUserRegistration(b *testing.B) {
	ctx := framework.NewTestContext("BenchmarkConcurrentUserRegistration", framework.WithDatabase())
	defer ctx.Cleanup()

	userRepo := repository.NewUserRepository()
	taskRepo := repository.NewTaskRepository()
	taskService := task.NewService(taskRepo)
	userService := user.NewService(userRepo, taskService)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// 并发用户注册 - 使用时间戳和随机数确保唯一性
			timestamp := time.Now().UnixNano()
			randomNum := rand.Int63()
			uniqueID := fmt.Sprintf("%d_%d", timestamp, randomNum)
			
			testUser := &user.User{
				UserName: fmt.Sprintf("%s_%s", ctx.GetDataGenerator().GenerateUniqueString("concurrent"), uniqueID),
				Email:    fmt.Sprintf("%s_%s@example.com", ctx.GetDataGenerator().GenerateUniqueString("concurrent"), uniqueID),
				Pwd:      "password123",
				TimeZone: "UTC",
				Lang:     "en",
			}
			err := userService.RegisterUser(testUser)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}