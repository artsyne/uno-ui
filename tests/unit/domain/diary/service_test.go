package diary_test

import (
	"fmt"
	"testing"

	"github.com/biuren/hebe/internal/core/di"
	"github.com/biuren/hebe/internal/domain/diary"
	"github.com/biuren/hebe/pkg"
	testhelpers "github.com/biuren/hebe/tests/helpers"
)

func TestDiaryService_AddDiary(t *testing.T) {
	// 初始化测试环境
	suite := testhelpers.SetupTestSuite(t)
	defer suite.Cleanup(t)

	// 清理表数据
	if suite.Database != nil {
		suite.Database.CleanupAllTables(t, "diaries")
	}

	diaryService, err := di.TryGet[diary.ServiceInterface](suite.Container)
	if err != nil {
		t.Fatalf("获取日记服务失败: %v", err)
	}

	tests := []struct {
		name    string
		diary   *diary.Diary
		wantErr bool
	}{
		{
			name: "创建新日记本",
			diary: &diary.Diary{
				Title:         "我的日记本",
				Remark:        "这是我的第一个日记本",
				ViewSetting:   1, // 私有
				SocialSetting: 1, // 不允许评论
				ThemeID:       1, // 主题ID
				AddUserID:     1,
				AddUserName:   "testuser",
				AddUTCTime:    pkg.GetNowTimeUnixMilli(),
				Status:        1, // 正常状态
			},
			wantErr: false,
		},
		{
			name: "创建公开日记本",
			diary: &diary.Diary{
				Title:         "公开日记",
				Remark:        "分享生活点滴",
				ViewSetting:   2, // 公开
				SocialSetting: 2, // 允许评论
				ThemeID:       2,
				AddUserID:     1,
				AddUserName:   "testuser",
				AddUTCTime:    pkg.GetNowTimeUnixMilli(),
				Status:        1,
			},
			wantErr: false,
		},
		{
			name: "日记本标题为空",
			diary: &diary.Diary{
				Title:         "", // 空标题
				Remark:        "测试",
				ViewSetting:   1,
				SocialSetting: 1,
				ThemeID:       1,
				AddUserID:     1,
				AddUserName:   "testuser",
				AddUTCTime:    pkg.GetNowTimeUnixMilli(),
				Status:        1,
			},
			wantErr: true,
		},
		{
			name: "无效的浏览设置",
			diary: &diary.Diary{
				Title:         "测试日记本",
				Remark:        "测试",
				ViewSetting:   3, // 无效设置
				SocialSetting: 1,
				ThemeID:       1,
				AddUserID:     1,
				AddUserName:   "testuser",
				AddUTCTime:    pkg.GetNowTimeUnixMilli(),
				Status:        1,
			},
			wantErr: true,
		},
		{
			name: "无效的主题ID",
			diary: &diary.Diary{
				Title:         "测试日记本",
				Remark:        "测试",
				ViewSetting:   1,
				SocialSetting: 1,
				ThemeID:       20, // 超出范围的主题ID
				AddUserID:     1,
				AddUserName:   "testuser",
				AddUTCTime:    pkg.GetNowTimeUnixMilli(),
				Status:        1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diaryID, err := diaryService.AddDiary(tt.diary)

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

			// 验证日记本ID已生成
			if diaryID <= 0 {
				t.Error("日记本ID应该大于0")
			}

			// 验证日记本对象已更新
			if tt.diary.DiaryID != diaryID {
				t.Errorf("日记本对象的ID未更新: 期望 %d, 实际 %d", diaryID, tt.diary.DiaryID)
			}
		})
	}
}

func TestDiaryService_GetDiaryByID(t *testing.T) {
	// 初始化测试环境
	suite := testhelpers.SetupTestSuite(t)
	defer suite.Cleanup(t)

	// 清理表数据
	if suite.Database != nil {
		suite.Database.CleanupAllTables(t, "diaries")
	}

	diaryService, err := di.TryGet[diary.ServiceInterface](suite.Container)
	if err != nil {
		t.Fatalf("获取日记服务失败: %v", err)
	}

	// 先创建一个日记本
	testDiary := &diary.Diary{
		Title:         "测试日记本",
		Remark:        "用于测试的日记本",
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
		diaryID int64
		wantErr bool
	}{
		{
			name:    "获取存在的日记本",
			diaryID: diaryID,
			wantErr: false,
		},
		{
			name:    "获取不存在的日记本",
			diaryID: 999999, // 不存在的ID
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := diaryService.GetDiaryByID(tt.diaryID)

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

			// 验证日记本数据
			if result.DiaryID != tt.diaryID {
				t.Errorf("日记本ID不匹配: 期望 %d, 实际 %d", tt.diaryID, result.DiaryID)
			}

			if result.Title != testDiary.Title {
				t.Errorf("日记本标题不匹配: 期望 '%s', 实际 '%s'", testDiary.Title, result.Title)
			}

			if result.AddUserID != testDiary.AddUserID {
				t.Errorf("创建者ID不匹配: 期望 %d, 实际 %d", testDiary.AddUserID, result.AddUserID)
			}
		})
	}
}

func TestDiaryService_UpdateDiary(t *testing.T) {
	// 初始化测试环境
	suite := testhelpers.SetupTestSuite(t)
	defer suite.Cleanup(t)

	// 清理表数据
	if suite.Database != nil {
		suite.Database.CleanupAllTables(t, "diaries")
	}

	diaryService, err := di.TryGet[diary.ServiceInterface](suite.Container)
	if err != nil {
		t.Fatalf("获取日记服务失败: %v", err)
	}

	// 先创建一个日记本
	originalDiary := &diary.Diary{
		Title:         "原始日记本",
		Remark:        "原始描述",
		ViewSetting:   1,
		SocialSetting: 1,
		ThemeID:       1,
		AddUserID:     1,
		AddUserName:   "testuser",
		AddUTCTime:    pkg.GetNowTimeUnixMilli(),
		Status:        1,
	}

	diaryID, err := diaryService.AddDiary(originalDiary)
	if err != nil {
		t.Fatalf("创建日记本失败: %v", err)
	}

	tests := []struct {
		name    string
		diary   *diary.Diary
		wantErr bool
	}{
		{
			name: "更新日记本标题和描述",
			diary: &diary.Diary{
				DiaryID:       diaryID,
				Title:         "更新后的日记本",
				Remark:        "更新后的描述",
				ViewSetting:   1,
				SocialSetting: 1,
				ThemeID:       1,
				AddUserID:     1,
				AddUserName:   "testuser",
				AddUTCTime:    pkg.GetNowTimeUnixMilli(),
				Status:        1,
			},
			wantErr: false,
		},
		{
			name: "更新为公开状态",
			diary: &diary.Diary{
				DiaryID:       diaryID,
				Title:         "公开日记本",
				Remark:        "现在是公开的",
				ViewSetting:   2, // 公开
				SocialSetting: 2, // 允许评论
				ThemeID:       2,
				AddUserID:     1,
				AddUserName:   "testuser",
				AddUTCTime:    pkg.GetNowTimeUnixMilli(),
				Status:        1,
			},
			wantErr: false,
		},
		{
			name: "无效的标题",
			diary: &diary.Diary{
				DiaryID:       diaryID,
				Title:         "", // 空标题
				Remark:        "测试",
				ViewSetting:   1,
				SocialSetting: 1,
				ThemeID:       1,
				AddUserID:     1,
				AddUserName:   "testuser",
				AddUTCTime:    pkg.GetNowTimeUnixMilli(),
				Status:        1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := diaryService.UpdateDiary(tt.diary)

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
			updated, err := diaryService.GetDiaryByID(tt.diary.DiaryID)
			if err != nil {
				t.Errorf("获取更新后的日记本失败: %v", err)
				return
			}

			if updated.Title != tt.diary.Title {
				t.Errorf("日记本标题未更新: 期望 '%s', 实际 '%s'", tt.diary.Title, updated.Title)
			}

			if updated.Remark != tt.diary.Remark {
				t.Errorf("日记本描述未更新: 期望 '%s', 实际 '%s'", tt.diary.Remark, updated.Remark)
			}

			if updated.ViewSetting != tt.diary.ViewSetting {
				t.Errorf("浏览设置未更新: 期望 %d, 实际 %d", tt.diary.ViewSetting, updated.ViewSetting)
			}
		})
	}
}

func TestDiaryService_ListDiariesForWrite(t *testing.T) {
	// 初始化测试环境
	suite := testhelpers.SetupTestSuite(t)
	defer suite.Cleanup(t)

	// 清理表数据
	if suite.Database != nil {
		suite.Database.CleanupAllTables(t, "diaries")
	}

	diaryService, err := di.TryGet[diary.ServiceInterface](suite.Container)
	if err != nil {
		t.Fatalf("获取日记服务失败: %v", err)
	}

	// 创建多个用户的日记本
	userID := int64(1)
	diaries := []*diary.Diary{
		{
			Title:         "个人日记",
			Remark:        "私人日记本",
			ViewSetting:   1,
			SocialSetting: 1,
			ThemeID:       1,
			AddUserID:     userID,
			AddUserName:   "testuser",
			AddUTCTime:    pkg.GetNowTimeUnixMilli(),
			Status:        1,
		},
		{
			Title:         "旅行日记",
			Remark:        "记录旅行足迹",
			ViewSetting:   2,
			SocialSetting: 2,
			ThemeID:       2,
			AddUserID:     userID,
			AddUserName:   "testuser",
			AddUTCTime:    pkg.GetNowTimeUnixMilli(),
			Status:        1,
		},
		{
			Title:         "其他用户的日记",
			Remark:        "不应该在结果中",
			ViewSetting:   1,
			SocialSetting: 1,
			ThemeID:       1,
			AddUserID:     2, // 不同的用户ID
			AddUserName:   "otheruser",
			AddUTCTime:    pkg.GetNowTimeUnixMilli(),
			Status:        1,
		},
	}

	// 创建日记本
	for _, diaryObj := range diaries {
		_, err := diaryService.AddDiary(diaryObj)
		if err != nil {
			t.Fatalf("创建日记本失败: %v", err)
		}
	}

	// 获取用户的日记本列表（用于写日记）
	result, err := diaryService.ListDiariesForWrite(userID)
	if err != nil {
		t.Fatalf("获取用户日记本失败: %v", err)
	}

	// 验证结果 - 应该只有属于指定用户的日记本
	expectedCount := 2
	if len(result) != expectedCount {
		t.Errorf("期望获取到%d个日记本，实际获取到%d个", expectedCount, len(result))
	}

	// 验证所有日记本都属于指定用户
	for i, diary := range result {
		if diary.AddUserID != userID {
			t.Errorf("日记本%d不属于指定用户: 期望用户ID %d, 实际 %d", i, userID, diary.AddUserID)
		}

		if diary.DiaryID <= 0 {
			t.Errorf("日记本%d的ID应该大于0", i)
		}
	}
}

func TestDiaryService_ListDiariesByUser(t *testing.T) {
	// 初始化测试环境
	suite := testhelpers.SetupTestSuite(t)
	defer suite.Cleanup(t)

	// 清理表数据
	if suite.Database != nil {
		suite.Database.CleanupAllTables(t, "diaries")
	}

	diaryService, err := di.TryGet[diary.ServiceInterface](suite.Container)
	if err != nil {
		t.Fatalf("获取日记服务失败: %v", err)
	}

	// 创建测试日记本
	userID := int64(1)
	totalDiaries := 5

	for i := 0; i < totalDiaries; i++ {
		diary := &diary.Diary{
			Title:         fmt.Sprintf("日记本%d", i+1),
			Remark:        fmt.Sprintf("第%d个日记本", i+1),
			ViewSetting:   1,
			SocialSetting: 1,
			ThemeID:       1,
			AddUserID:     userID,
			AddUserName:   "testuser",
			AddUTCTime:    pkg.GetNowTimeUnixMilli(),
			Status:        1,
		}

		_, err := diaryService.AddDiary(diary)
		if err != nil {
			t.Fatalf("创建日记本%d失败: %v", i+1, err)
		}
	}

	tests := []struct {
		name          string
		userID        int64
		pageIndex     int
		pageSize      int
		expectedCount int
		expectedTotal int
	}{
		{
			name:          "获取第一页",
			userID:        userID,
			pageIndex:     0,
			pageSize:      3,
			expectedCount: 3,
			expectedTotal: totalDiaries,
		},
		{
			name:          "获取第二页",
			userID:        userID,
			pageIndex:     1,
			pageSize:      3,
			expectedCount: 2, // 剩余2个
			expectedTotal: totalDiaries,
		},
		{
			name:          "获取不存在用户的日记本",
			userID:        999,
			pageIndex:     0,
			pageSize:      10,
			expectedCount: 0,
			expectedTotal: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, total, err := diaryService.ListDiariesByUser(tt.userID, tt.pageIndex, tt.pageSize)
			if err != nil {
				t.Fatalf("获取用户日记本失败: %v", err)
			}

			if len(result) != tt.expectedCount {
				t.Errorf("期望获取到%d个日记本，实际获取到%d个", tt.expectedCount, len(result))
			}

			if total != tt.expectedTotal {
				t.Errorf("期望总数为%d，实际为%d", tt.expectedTotal, total)
			}

			// 验证所有日记本都属于指定用户
			for i, diary := range result {
				if diary.AddUserID != tt.userID {
					t.Errorf("日记本%d不属于指定用户: 期望用户ID %d, 实际 %d", i, tt.userID, diary.AddUserID)
				}
			}
		})
	}
}

func TestDiaryService_PostCountOperations(t *testing.T) {
	// 初始化测试环境
	suite := testhelpers.SetupTestSuite(t)
	defer suite.Cleanup(t)

	// 清理表数据
	if suite.Database != nil {
		suite.Database.CleanupAllTables(t, "diaries")
	}

	diaryService, err := di.TryGet[diary.ServiceInterface](suite.Container)
	if err != nil {
		t.Fatalf("获取日记服务失败: %v", err)
	}

	// 先创建一个日记本
	testDiary := &diary.Diary{
		Title:         "测试日记本",
		Remark:        "用于测试日记数量操作",
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

	actionTime := pkg.GetNowTimeUnixMilli()
	userID := int64(1)
	userName := "testuser"

	// 测试增加日记数量
	t.Run("增加日记数量", func(t *testing.T) {
		err := diaryService.IncrementPostCount(diaryID, userID, userName, actionTime)
		if err != nil {
			t.Errorf("增加日记数量失败: %v", err)
		}

		// 获取更新后的日记本验证
		updated, err := diaryService.GetDiaryByID(diaryID)
		if err != nil {
			t.Errorf("获取更新后的日记本失败: %v", err)
		} else {
			// 验证最后操作信息已更新
			if updated.LastActionUserID != userID {
				t.Errorf("最后操作用户ID未更新: 期望 %d, 实际 %d", userID, updated.LastActionUserID)
			}
		}
	})

	// 测试减少日记数量
	t.Run("减少日记数量", func(t *testing.T) {
		err := diaryService.DecrementPostCount(diaryID, userID, userName, actionTime)
		if err != nil {
			t.Errorf("减少日记数量失败: %v", err)
		}
	})

	// 测试操作不存在的日记本
	t.Run("操作不存在的日记本", func(t *testing.T) {
		err := diaryService.IncrementPostCount(999999, userID, userName, actionTime)
		if err == nil {
			t.Error("期望操作不存在的日记本时有错误")
		}
	})
}
