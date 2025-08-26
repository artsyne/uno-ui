package diary

import (
	"github.com/biuren/hebe/internal/app/controller/validation"
	"github.com/biuren/hebe/internal/domain/diary"
	errors "github.com/biuren/hebe/internal/utils/error"
	"github.com/gin-gonic/gin"
)

// 注意：已移除重复的parseDiaryID方法，直接使用validation.ParseID

// validateDiaryOwnership 验证日记本所有权
// 检查指定用户是否拥有指定日记本的操作权限
//
// 参数:
//   - diaryID: 日记本ID
//   - userID: 用户ID
//
// 返回值:
//   - *diary.Diary: 日记本对象，如果验证失败返回nil
//   - error: 错误信息
func (ctrl *DiaryController) validateDiaryOwnership(diaryID, userID int64) (*diary.Diary, error) {
	diaryData, err := ctrl.diaryService.GetDiaryByID(diaryID)
	if err != nil {
		return nil, errors.DiaryNotFound.Err()
	}

	// 验证所有权
	if diaryData.AddUserID != userID {
		return nil, errors.PermissionDenied.Err()
	}

	return diaryData, nil
}

// getAndValidateDiary 获取并验证日记本权限（用于表单渲染）
// 从URL参数中获取diary_id，验证其有效性和用户权限
//
// 参数:
//   - c: Gin上下文
//   - userID: 当前用户ID
//
// 返回值:
//   - *diary.Diary: 日记本对象，如果验证失败返回nil
//   - error: 错误信息
func (ctrl *DiaryController) getAndValidateDiary(c *gin.Context, userID int64) (*diary.Diary, error) {
	diaryID, err := validation.ParseID(c.Query("diary_id"))
	if err != nil {
		return nil, err
	}

	return ctrl.validateDiaryOwnership(diaryID, userID)
}