package post_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/biuren/hebe/internal/modules/diary"
	"github.com/biuren/hebe/internal/modules/post"
	"github.com/biuren/hebe/pkg"
	testhelpers "github.com/biuren/hebe/tests/helpers"
)

func TestPostService_AddPost(t *testing.T) {
	// 初始化测试环境
	dbHelper := testhelpers.SetupTestDatabase(t)
	defer dbHelper.Cleanup(t)

	// 开始事务
	dbHelper.BeginTransaction(t)
	defer dbHelper.RollbackTransaction(t)

	// 创建必要的服务
	diaryService := diary.NewService()
	postService := post.NewService()

	// 先创建一个日记本
	testDiary := &diary.Diary{
		Title:         "测试日记本",
		Remark:        "用于测试日记内容的日记本",
		ViewSetting:   1,
		SocialSetting: 1,
		ThemeID:       1,
		AddUserID:     1,
		AddUserName:   "testuser",
		AddUTCTime:    pkg.GetNowTimeUnixMilli(),
		Status:        1,
	}

	diaryID, err := diaryService.AddDiary(testDiary)
	if err != nil {
		t.Fatalf("创建日记本失败: %v", err)
	}

	tests := []struct {
		name    string
		post    *post.Post
		wantErr bool
	}{
		{
			name: "创建新文本日记",
			post: &post.Post{
				Title:       "我的第一篇日记",
				Body:        "今天是个好天气，心情很不错。这是详细的日记内容。",
				DiaryID:     diaryID,
				AddUserID:   1,
				AddUserName: "testuser",
				AddUTCTime:  pkg.GetNowTimeUnixMilli(),
				Status:      3, // 正常状态
				InfoType:    1, // 文本类型
				IsTop:       0, // 不置顶
			},
			wantErr: false,
		},
		{
			name: "创建图片日记",
			post: &post.Post{
				Title:       "图片日记",
				Body:        "分享一些美丽的风景照片，记录生活的美好瞬间。",
				DiaryID:     diaryID,
				AddUserID:   1,
				AddUserName: "testuser",
				AddUTCTime:  pkg.GetNowTimeUnixMilli(),
				Status:      3,
				InfoType:    2, // 图片类型
				IsTop:       0,
			},
			wantErr: false,
		},
		{
			name: "创建置顶日记",
			post: &post.Post{
				Title:       "重要日记",
				Body:        "这是一篇重要的日记，需要置顶显示。",
				DiaryID:     diaryID,
				AddUserID:   1,
				AddUserName: "testuser",
				AddUTCTime:  pkg.GetNowTimeUnixMilli(),
				Status:      3,
				InfoType:    1,
				IsTop:       1, // 置顶
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			postID, err := postService.AddPost(tt.post)

			if tt.wantErr {
				if err == nil {
					t.Errorf("期望有错误，但没有错误")
				}
				return
			}

			if err != nil {
				t.Errorf("不期望有错误，但有错误: %v", err)
				return
			}

			// 验证日记内容ID已生成
			if postID <= 0 {
				t.Error("日记内容ID应该大于0")
			}

			// 验证日记内容对象已更新
			if tt.post.PostID != postID {
				t.Errorf("日记内容对象的ID未更新: 期望 %d, 实际 %d", postID, tt.post.PostID)
			}

			// 验证摘要已生成（从Body中自动生成）
			if tt.post.Summary == "" {
				t.Error("摘要应该已生成")
			}

			// 验证摘要是Body的截取
			if !strings.Contains(tt.post.Body, strings.TrimSuffix(tt.post.Summary, "...")) {
				t.Error("摘要应该是正文的截取")
			}
		})
	}
}

func TestPostService_GetPostByID(t *testing.T) {
	// 初始化测试环境
	dbHelper := testhelpers.SetupTestDatabase(t)
	defer dbHelper.Cleanup(t)

	// 开始事务
	dbHelper.BeginTransaction(t)
	defer dbHelper.RollbackTransaction(t)

	// 创建必要的服务和数据
	diaryService := diary.NewService()
	postService := post.NewService()

	// 创建日记本
	testDiary := &diary.Diary{
		Title:         "测试日记本",
		Remark:        "测试用",
		ViewSetting:   1,
		SocialSetting: 1,
		ThemeID:       1,
		AddUserID:     1,
		AddUserName:   "testuser",
		AddUTCTime:    pkg.GetNowTimeUnixMilli(),
		Status:        1,
	}

	diaryID, err := diaryService.AddDiary(testDiary)
	if err != nil {
		t.Fatalf("创建日记本失败: %v", err)
	}

	// 创建日记内容
	testPost := &post.Post{
		Title:       "测试日记内容",
		Body:        "这是用于测试的日记内容，包含详细的文字描述。",
		DiaryID:     diaryID,
		AddUserID:   1,
		AddUserName: "testuser",
		AddUTCTime:  pkg.GetNowTimeUnixMilli(),
		Status:      3,
		InfoType:    1,
		IsTop:       0,
	}

	postID, err := postService.AddPost(testPost)
	if err != nil {
		t.Fatalf("创建日记内容失败: %v", err)
	}

	tests := []struct {
		name    string
		postID  int64
		wantErr bool
	}{
		{
			name:    "获取存在的日记内容",
			postID:  postID,
			wantErr: false,
		},
		{
			name:    "获取不存在的日记内容",
			postID:  999999, // 不存在的ID
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := postService.GetPostByID(tt.postID)

			if tt.wantErr {
				if err == nil {
					t.Errorf("期望有错误，但没有错误")
				}
				return
			}

			if err != nil {
				t.Errorf("不期望有错误，但有错误: %v", err)
				return
			}

			if result == nil {
				t.Error("结果不应为空")
				return
			}

			// 验证日记内容数据
			if result.PostID != tt.postID {
				t.Errorf("日记内容ID不匹配: 期望 %d, 实际 %d", tt.postID, result.PostID)
			}

			if result.Title != testPost.Title {
				t.Errorf("日记内容标题不匹配: 期望 '%s', 实际 '%s'", testPost.Title, result.Title)
			}

			if result.Body != testPost.Body {
				t.Errorf("日记内容不匹配: 期望 '%s', 实际 '%s'", testPost.Body, result.Body)
			}

			if result.DiaryID != diaryID {
				t.Errorf("日记本ID不匹配: 期望 %d, 实际 %d", diaryID, result.DiaryID)
			}
		})
	}
}

func TestPostService_UpdatePost(t *testing.T) {
	// 初始化测试环境
	dbHelper := testhelpers.SetupTestDatabase(t)
	defer dbHelper.Cleanup(t)

	// 开始事务
	dbHelper.BeginTransaction(t)
	defer dbHelper.RollbackTransaction(t)

	// 创建必要的服务和数据
	diaryService := diary.NewService()
	postService := post.NewService()

	// 创建日记本
	testDiary := &diary.Diary{
		Title:         "测试日记本",
		Remark:        "测试用",
		ViewSetting:   1,
		SocialSetting: 1,
		ThemeID:       1,
		AddUserID:     1,
		AddUserName:   "testuser",
		AddUTCTime:    pkg.GetNowTimeUnixMilli(),
		Status:        1,
	}

	diaryID, err := diaryService.AddDiary(testDiary)
	if err != nil {
		t.Fatalf("创建日记本失败: %v", err)
	}

	// 创建原始日记内容
	originalPost := &post.Post{
		Title:       "原始标题",
		Body:        "原始内容，这是最初的日记正文。",
		DiaryID:     diaryID,
		AddUserID:   1,
		AddUserName: "testuser",
		AddUTCTime:  pkg.GetNowTimeUnixMilli(),
		Status:      3,
		InfoType:    1,
		IsTop:       0,
	}

	postID, err := postService.AddPost(originalPost)
	if err != nil {
		t.Fatalf("创建日记内容失败: %v", err)
	}

	tests := []struct {
		name    string
		post    *post.Post
		wantErr bool
	}{
		{
			name: "更新标题和内容",
			post: &post.Post{
				PostID:        postID,
				Title:         "更新后的标题",
				Body:          "更新后的内容，更加详细和丰富，包含更多有用信息。",
				DiaryID:       diaryID,
				AddUserID:     1,
				AddUserName:   "testuser",
				AddUTCTime:    pkg.GetNowTimeUnixMilli(),
				UpdateUTCTime: pkg.GetNowTimeUnixMilli(),
				Status:        3,
				InfoType:      1,
				IsTop:         0,
			},
			wantErr: false,
		},
		{
			name: "设置为置顶",
			post: &post.Post{
				PostID:        postID,
				Title:         "重要日记",
				Body:          "这是一篇重要的日记，需要置顶显示给所有人看到。",
				DiaryID:       diaryID,
				AddUserID:     1,
				AddUserName:   "testuser",
				AddUTCTime:    pkg.GetNowTimeUnixMilli(),
				UpdateUTCTime: pkg.GetNowTimeUnixMilli(),
				Status:        3,
				InfoType:      1,
				IsTop:         1, // 置顶
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := postService.UpdatePost(tt.post)

			if tt.wantErr {
				if err == nil {
					t.Errorf("期望有错误，但没有错误")
				}
				return
			}

			if err != nil {
				t.Errorf("不期望有错误，但有错误: %v", err)
				return
			}

			// 验证更新是否成功
			updated, err := postService.GetPostByID(tt.post.PostID)
			if err != nil {
				t.Errorf("获取更新后的日记内容失败: %v", err)
				return
			}

			if updated.Title != tt.post.Title {
				t.Errorf("日记内容标题未更新: 期望 '%s', 实际 '%s'", tt.post.Title, updated.Title)
			}

			if updated.Body != tt.post.Body {
				t.Errorf("日记内容未更新: 期望 '%s', 实际 '%s'", tt.post.Body, updated.Body)
			}

			if updated.IsTop != tt.post.IsTop {
				t.Errorf("置顶设置未更新: 期望 %d, 实际 %d", tt.post.IsTop, updated.IsTop)
			}
		})
	}
}

func TestPostService_ListPostsByUserID(t *testing.T) {
	// 初始化测试环境
	dbHelper := testhelpers.SetupTestDatabase(t)
	defer dbHelper.Cleanup(t)

	// 开始事务
	dbHelper.BeginTransaction(t)
	defer dbHelper.RollbackTransaction(t)

	// 创建必要的服务和数据
	diaryService := diary.NewService()
	postService := post.NewService()

	// 创建两个日记本
	diary1 := &diary.Diary{
		Title:         "日记本1",
		Remark:        "第一个日记本",
		ViewSetting:   1,
		SocialSetting: 1,
		ThemeID:       1,
		AddUserID:     1,
		AddUserName:   "testuser",
		AddUTCTime:    pkg.GetNowTimeUnixMilli(),
		Status:        1,
	}

	diaryID1, err := diaryService.AddDiary(diary1)
	if err != nil {
		t.Fatalf("创建日记本1失败: %v", err)
	}

	diary2 := &diary.Diary{
		Title:         "日记本2",
		Remark:        "第二个日记本",
		ViewSetting:   1,
		SocialSetting: 1,
		ThemeID:       1,
		AddUserID:     2,
		AddUserName:   "otheruser",
		AddUTCTime:    pkg.GetNowTimeUnixMilli(),
		Status:        1,
	}

	diaryID2, err := diaryService.AddDiary(diary2)
	if err != nil {
		t.Fatalf("创建日记本2失败: %v", err)
	}

	// 在第一个日记本中创建多篇日记
	userID := int64(1)
	totalPostsInDiary1 := 5
	for i := 0; i < totalPostsInDiary1; i++ {
		post := &post.Post{
			Title:       fmt.Sprintf("日记本1的第%d篇日记", i+1),
			Body:        fmt.Sprintf("这是第%d篇日记的详细内容，记录了生活中的点点滴滴。", i+1),
			DiaryID:     diaryID1,
			AddUserID:   userID,
			AddUserName: "testuser",
			AddUTCTime:  pkg.GetNowTimeUnixMilli(),
			Status:      3,
			InfoType:    1,
			IsTop:       0,
		}

		_, err := postService.AddPost(post)
		if err != nil {
			t.Fatalf("创建日记内容%d失败: %v", i+1, err)
		}
	}

	// 在第二个日记本中创建一篇日记（其他用户）
	postInDiary2 := &post.Post{
		Title:       "日记本2的日记",
		Body:        "这是日记本2中的内容，属于其他用户。",
		DiaryID:     diaryID2,
		AddUserID:   2,
		AddUserName: "otheruser",
		AddUTCTime:  pkg.GetNowTimeUnixMilli(),
		Status:      3,
		InfoType:    1,
		IsTop:       0,
	}

	_, err = postService.AddPost(postInDiary2)
	if err != nil {
		t.Fatalf("创建日记本2的日记内容失败: %v", err)
	}

	tests := []struct {
		name          string
		userID        int64
		diaryID       int64
		pageIndex     int
		pageSize      int
		expectedCount int
		expectedTotal int
	}{
		{
			name:          "获取用户1在日记本1的第一页",
			userID:        userID,
			diaryID:       diaryID1,
			pageIndex:     0,
			pageSize:      3,
			expectedCount: 3,
			expectedTotal: totalPostsInDiary1,
		},
		{
			name:          "获取用户1在日记本1的第二页",
			userID:        userID,
			diaryID:       diaryID1,
			pageIndex:     1,
			pageSize:      3,
			expectedCount: 2, // 剩余2篇
			expectedTotal: totalPostsInDiary1,
		},
		{
			name:          "获取用户1在不存在的日记本",
			userID:        userID,
			diaryID:       999999,
			pageIndex:     0,
			pageSize:      10,
			expectedCount: 0,
			expectedTotal: 0,
		},
		{
			name:          "获取不存在的用户",
			userID:        999,
			diaryID:       diaryID1,
			pageIndex:     0,
			pageSize:      10,
			expectedCount: 0,
			expectedTotal: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, total, err := postService.ListPostsByUserID(tt.userID, tt.diaryID, tt.pageIndex, tt.pageSize)
			if err != nil {
				t.Fatalf("获取用户日记内容失败: %v", err)
			}

			if len(result) != tt.expectedCount {
				t.Errorf("期望获取到%d篇日记，实际获取到%d篇", tt.expectedCount, len(result))
			}

			if total != tt.expectedTotal {
				t.Errorf("期望总数为%d，实际为%d", tt.expectedTotal, total)
			}

			// 验证所有日记都属于指定用户和日记本
			for i, post := range result {
				if post.DiaryID != tt.diaryID {
					t.Errorf("日记%d不属于指定日记本: 期望日记本ID %d, 实际 %d", i, tt.diaryID, post.DiaryID)
				}

				if post.AddUserID != tt.userID {
					t.Errorf("日记%d不属于指定用户: 期望用户ID %d, 实际 %d", i, tt.userID, post.AddUserID)
				}

				if post.PostID <= 0 {
					t.Errorf("日记%d的ID应该大于0", i)
				}
			}
		})
	}
}

func TestPostService_GetPostCountByUserID(t *testing.T) {
	// 初始化测试环境
	dbHelper := testhelpers.SetupTestDatabase(t)
	defer dbHelper.Cleanup(t)

	// 开始事务
	dbHelper.BeginTransaction(t)
	defer dbHelper.RollbackTransaction(t)

	// 创建必要的服务和数据
	diaryService := diary.NewService()
	postService := post.NewService()

	// 创建日记本
	testDiary := &diary.Diary{
		Title:         "测试日记本",
		Remark:        "用于统计测试",
		ViewSetting:   1,
		SocialSetting: 1,
		ThemeID:       1,
		AddUserID:     1,
		AddUserName:   "testuser",
		AddUTCTime:    pkg.GetNowTimeUnixMilli(),
		Status:        1,
	}

	diaryID, err := diaryService.AddDiary(testDiary)
	if err != nil {
		t.Fatalf("创建日记本失败: %v", err)
	}

	// 创建多篇日记
	userID := int64(1)
	totalPosts := 3

	for i := 0; i < totalPosts; i++ {
		post := &post.Post{
			Title:       fmt.Sprintf("测试日记%d", i+1),
			Body:        fmt.Sprintf("这是第%d篇测试日记的内容。", i+1),
			DiaryID:     diaryID,
			AddUserID:   userID,
			AddUserName: "testuser",
			AddUTCTime:  pkg.GetNowTimeUnixMilli(),
			Status:      3,
			InfoType:    1,
			IsTop:       0,
		}

		_, err := postService.AddPost(post)
		if err != nil {
			t.Fatalf("创建日记内容%d失败: %v", i+1, err)
		}
	}

	tests := []struct {
		name          string
		userID        int64
		diaryID       int64
		expectedCount int
	}{
		{
			name:          "统计存在的用户和日记本",
			userID:        userID,
			diaryID:       diaryID,
			expectedCount: totalPosts,
		},
		{
			name:          "统计不存在的用户",
			userID:        999,
			diaryID:       diaryID,
			expectedCount: 0,
		},
		{
			name:          "统计不存在的日记本",
			userID:        userID,
			diaryID:       999999,
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			count, err := postService.GetPostCountByUserID(tt.userID, tt.diaryID)
			if err != nil {
				t.Fatalf("统计日记数量失败: %v", err)
			}

			if count != tt.expectedCount {
				t.Errorf("期望统计到%d篇日记，实际统计到%d篇", tt.expectedCount, count)
			}
		})
	}
}

func TestPostService_UpdatePostByID(t *testing.T) {
	// 初始化测试环境
	dbHelper := testhelpers.SetupTestDatabase(t)
	defer dbHelper.Cleanup(t)

	// 开始事务
	dbHelper.BeginTransaction(t)
	defer dbHelper.RollbackTransaction(t)

	// 创建必要的服务和数据
	diaryService := diary.NewService()
	postService := post.NewService()

	// 创建日记本
	testDiary := &diary.Diary{
		Title:         "测试日记本",
		Remark:        "测试用",
		ViewSetting:   1,
		SocialSetting: 1,
		ThemeID:       1,
		AddUserID:     1,
		AddUserName:   "testuser",
		AddUTCTime:    pkg.GetNowTimeUnixMilli(),
		Status:        1,
	}

	diaryID, err := diaryService.AddDiary(testDiary)
	if err != nil {
		t.Fatalf("创建日记本失败: %v", err)
	}

	// 创建日记内容
	originalPost := &post.Post{
		Title:       "原始标题",
		Body:        "原始内容",
		DiaryID:     diaryID,
		AddUserID:   1,
		AddUserName: "testuser",
		AddUTCTime:  pkg.GetNowTimeUnixMilli(),
		Status:      3,
		InfoType:    1,
		IsTop:       0,
		ViewNum:     0,
		LikeNum:     0,
		CommentNum:  0,
	}

	postID, err := postService.AddPost(originalPost)
	if err != nil {
		t.Fatalf("创建日记内容失败: %v", err)
	}

	// 测试按ID更新
	t.Run("按ID更新日记", func(t *testing.T) {
		updatePost := &post.Post{
			PostID:        postID,
			Title:         "按ID更新的标题",
			Body:          "按ID更新的内容",
			DiaryID:       diaryID,
			AddUserID:     1,
			AddUserName:   "testuser",
			UpdateUTCTime: pkg.GetNowTimeUnixMilli(),
			Status:        3,
			InfoType:      1,
			IsTop:         1, // 设置为置顶
			ViewNum:       10,
			LikeNum:       5,
			CommentNum:    2,
		}

		err := postService.UpdatePostByID(updatePost)
		if err != nil {
			t.Fatalf("按ID更新日记失败: %v", err)
		}

		// 验证更新是否成功
		updated, err := postService.GetPostByID(postID)
		if err != nil {
			t.Fatalf("获取更新后的日记失败: %v", err)
		}

		if updated.Title != updatePost.Title {
			t.Errorf("标题未更新: 期望 '%s', 实际 '%s'", updatePost.Title, updated.Title)
		}

		if updated.IsTop != updatePost.IsTop {
			t.Errorf("置顶状态未更新: 期望 %d, 实际 %d", updatePost.IsTop, updated.IsTop)
		}

		if updated.ViewNum != updatePost.ViewNum {
			t.Errorf("查看数量未更新: 期望 %d, 实际 %d", updatePost.ViewNum, updated.ViewNum)
		}
	})
}

func TestPostService_DeletePostByID(t *testing.T) {
	// 初始化测试环境
	dbHelper := testhelpers.SetupTestDatabase(t)
	defer dbHelper.Cleanup(t)

	// 开始事务
	dbHelper.BeginTransaction(t)
	defer dbHelper.RollbackTransaction(t)

	// 创建必要的服务和数据
	diaryService := diary.NewService()
	postService := post.NewService()

	// 创建日记本
	testDiary := &diary.Diary{
		Title:         "测试日记本",
		Remark:        "测试用",
		ViewSetting:   1,
		SocialSetting: 1,
		ThemeID:       1,
		AddUserID:     1,
		AddUserName:   "testuser",
		AddUTCTime:    pkg.GetNowTimeUnixMilli(),
		Status:        1,
	}

	diaryID, err := diaryService.AddDiary(testDiary)
	if err != nil {
		t.Fatalf("创建日记本失败: %v", err)
	}

	// 创建待删除的日记内容
	testPost := &post.Post{
		Title:       "待删除的日记",
		Body:        "这篇日记将被删除",
		DiaryID:     diaryID,
		AddUserID:   1,
		AddUserName: "testuser",
		AddUTCTime:  pkg.GetNowTimeUnixMilli(),
		Status:      3,
		InfoType:    1,
		IsTop:       0,
	}

	postID, err := postService.AddPost(testPost)
	if err != nil {
		t.Fatalf("创建日记内容失败: %v", err)
	}

	tests := []struct {
		name    string
		postID  int64
		wantErr bool
	}{
		{
			name:    "正常删除日记内容",
			postID:  postID,
			wantErr: false,
		},
		{
			name:    "删除不存在的日记内容",
			postID:  999999, // 不存在的ID
			wantErr: true,
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := postService.DeletePostByID(tt.postID)

			if tt.wantErr {
				if err == nil {
					t.Errorf("期望有错误，但没有错误")
				}
				return
			}

			if err != nil {
				t.Errorf("不期望有错误，但有错误: %v", err)
				return
			}

			// 只有第一个测试用例（正常删除）需要验证删除是否成功
			if i == 0 {
				// 尝试获取已删除的日记内容，应该失败
				_, getErr := postService.GetPostByID(tt.postID)
				if getErr == nil {
					t.Error("删除后仍能获取到日记内容，删除失败")
				}
			}
		})
	}
}
