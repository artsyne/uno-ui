package benchmark

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/biuren/hebe/internal/domain/diary"
	"github.com/biuren/hebe/internal/domain/post"
	"github.com/biuren/hebe/internal/domain/task"
	"github.com/biuren/hebe/internal/domain/user"
	"github.com/biuren/hebe/internal/infra/repository"
	framework "github.com/biuren/hebe/tests/framework"
)

// 全局测试上下文和服务，避免重复初始化
var (
	benchCtx      *framework.TestContext
	benchServices *benchmarkServices
	benchOnce     sync.Once
)

type benchmarkServices struct {
	userService  *user.Service
	diaryService *diary.Service
	postService  *post.Service
	taskService  *task.Service
}

// 初始化基准测试环境
func initBenchmarkEnv() {
	benchOnce.Do(func() {
		benchCtx = framework.NewTestContext(framework.WithDatabase()).(*framework.TestContext)

		userRepo := repository.NewUserRepository()
		taskRepo := repository.NewTaskRepository()
		diaryRepo := repository.NewDiaryRepository()
		postRepo := repository.NewPostRepository(benchCtx.GetDatabase().GetDB())
		config := benchCtx.GetConfig().GetAppConfig()

		taskService := task.NewService(config, taskRepo)
		userService := user.NewService(config, userRepo, taskService)
		diaryService := diary.NewService(config, diaryRepo)
		postService := post.NewService(config, userRepo, taskService, diaryRepo, postRepo)

		benchServices = &benchmarkServices{
			userService:  userService,
			diaryService: diaryService,
			postService:  postService,
			taskService:  taskService,
		}
	})
}

// BenchmarkUserRegistration 用户注册性能基准测试
func BenchmarkUserRegistration(b *testing.B) {
	initBenchmarkEnv()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		timestamp := time.Now().UnixNano()
		randomNum := rand.Int63()
		uniqueID := fmt.Sprintf("%d_%d", timestamp, randomNum)
		testUser := &user.User{
			UserName: fmt.Sprintf("%s_%s", benchCtx.GetDataGenerator().GenerateUniqueString("benchuser"), uniqueID),
			Email:    fmt.Sprintf("%s_%s@example.com", benchCtx.GetDataGenerator().GenerateUniqueString("bench"), uniqueID),
			Pwd:      "password123",
			TimeZone: "UTC",
			Lang:     "en",
		}
		err := benchServices.userService.RegisterUser(testUser)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkDiaryCreation 日记本创建性能基准测试
func BenchmarkDiaryCreation(b *testing.B) {
	initBenchmarkEnv()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// 为每个测试创建独立的用户，避免数量限制
		timestamp := time.Now().UnixNano()
		randomNum := rand.Int63()
		uniqueID := fmt.Sprintf("%d_%d", timestamp, randomNum)

		testUser := &user.User{
			UserName: fmt.Sprintf("diary_bench_%s_%s", uniqueID, benchCtx.GetDataGenerator().GenerateUniqueString("user")),
			Email:    fmt.Sprintf("diary_bench_%s_%s@example.com", uniqueID, benchCtx.GetDataGenerator().GenerateUniqueString("user")),
			Pwd:      "password123",
			TimeZone: "UTC",
			Lang:     "en",
		}
		err := benchServices.userService.RegisterUser(testUser)
		if err != nil {
			b.Fatal(err)
		}

		testDiary := &diary.Diary{
			Title:         fmt.Sprintf("%s_%s", benchCtx.GetDataGenerator().GenerateUniqueString("benchdiary"), uniqueID),
			Remark:        "性能测试日记本",
			AddUserID:     testUser.UserId,
			AddUserName:   testUser.UserName,
			AddUTCTime:    time.Now().Unix(),
			Status:        1,
			ViewSetting:   1,
			SocialSetting: 1,
			ThemeID:       1,
		}
		_, err = benchServices.diaryService.AddDiary(testDiary)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkPostCreation 文章创建性能基准测试
func BenchmarkPostCreation(b *testing.B) {
	initBenchmarkEnv()

	// 创建测试用户 - 使用唯一前缀
	timestamp := time.Now().UnixNano()
	uniqueID := fmt.Sprintf("%d_%d", timestamp, rand.Int63())
	testUser := &user.User{
		UserName: fmt.Sprintf("post_bench_%s_%s", uniqueID, benchCtx.GetDataGenerator().GenerateUniqueString("user")),
		Email:    fmt.Sprintf("post_bench_%s_%s@example.com", uniqueID, benchCtx.GetDataGenerator().GenerateUniqueString("user")),
		Pwd:      "password123",
		TimeZone: "UTC",
		Lang:     "en",
	}
	err := benchServices.userService.RegisterUser(testUser)
	if err != nil {
		b.Fatal(err)
	}

	// 创建测试日记本
	testDiary := &diary.Diary{
		Title:         fmt.Sprintf("%s_%s", benchCtx.GetDataGenerator().GenerateUniqueString("benchdiary"), uniqueID),
		Remark:        "性能测试日记本",
		AddUserID:     testUser.UserId,
		AddUserName:   testUser.UserName,
		AddUTCTime:    time.Now().Unix(),
		Status:        1,
		ViewSetting:   1,
		SocialSetting: 1,
		ThemeID:       1,
	}
	diaryID, err := benchServices.diaryService.AddDiary(testDiary)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		testPost := &post.Post{
			Title:       benchCtx.GetDataGenerator().GenerateUniqueString("benchpost"),
			Body:        "这是一篇性能测试文章的内容",
			Summary:     "性能测试摘要",
			DiaryID:     diaryID,
			AddUserID:   testUser.UserId,
			AddUserName: testUser.UserName,
			AddUTCTime:  time.Now().Unix(),
			Status:      1,
			InfoType:    1,
		}
		_, err := benchServices.postService.AddPost(testPost)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkUserQuery 用户查询性能基准测试
func BenchmarkUserQuery(b *testing.B) {
	initBenchmarkEnv()

	// 创建测试用户 - 使用唯一前缀
	timestamp := time.Now().UnixNano()
	uniqueID := fmt.Sprintf("%d_%d", timestamp, rand.Int63())
	testUser := &user.User{
		UserName: fmt.Sprintf("query_bench_%s_%s", uniqueID, benchCtx.GetDataGenerator().GenerateUniqueString("user")),
		Email:    fmt.Sprintf("query_bench_%s_%s@example.com", uniqueID, benchCtx.GetDataGenerator().GenerateUniqueString("user")),
		Pwd:      "password123",
		TimeZone: "UTC",
		Lang:     "en",
	}
	err := benchServices.userService.RegisterUser(testUser)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// 查询用户
		_, err := benchServices.userService.GetUserByUserId(testUser.UserId)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkConcurrentUserRegistration 并发用户注册性能基准测试
func BenchmarkConcurrentUserRegistration(b *testing.B) {
	initBenchmarkEnv()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			timestamp := time.Now().UnixNano()
			randomNum := rand.Int63()
			uniqueID := fmt.Sprintf("%d_%d", timestamp, randomNum)
			testUser := &user.User{
				UserName: fmt.Sprintf("%s_%s", benchCtx.GetDataGenerator().GenerateUniqueString("concurrentuser"), uniqueID),
				Email:    fmt.Sprintf("%s_%s@example.com", benchCtx.GetDataGenerator().GenerateUniqueString("concurrent"), uniqueID),
				Pwd:      "password123",
				TimeZone: "UTC",
				Lang:     "en",
			}
			err := benchServices.userService.RegisterUser(testUser)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
