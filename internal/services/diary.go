package services

import (
	"github.com/biuren/hebe/internal/configs"
	"github.com/biuren/hebe/internal/errors"
	"github.com/biuren/hebe/internal/models"
)

// AddDiary 创建一个新的日记本
func AddDiary(diary *models.Diary) (int64, error) {
	// 校验参数
	if err := validateDiary(diary); err != nil {
		return 0, err
	}

	// 判断用户增加日记本数量是否超过限制
	if err := CheckDiaryLimit(diary.AddUserID); err != nil {
		return 0, err
	}

	// 新增日记本
	diaryID, err := diary.Insert()
	if err != nil {
		return 0, err
	}
	diary.DiaryID = diaryID

	return diaryID, nil
}

// validateDiary 校验日记本参数
func validateDiary(diary *models.Diary) error {
	if len(diary.Title) < 1 {
		return errors.ParseInputParamErr
	}
	if diary.ViewSetting != 1 && diary.ViewSetting != 2 {
		return errors.ParseInputParamErr
	}
	if diary.SocialSetting != 1 && diary.SocialSetting != 2 {
		return errors.ParseInputParamErr
	}
	if diary.ThemeID < 1 || diary.ThemeID > 16 {
		return errors.ParseInputParamErr
	}
	return nil
}

// CheckDiaryLimit 检查用户日记本数量是否超过限制
func CheckDiaryLimit(userID int64) error {
	daybookNum, err := models.GetDiariesCountByUserID(userID)
	if err != nil {
		return errors.ParseInputParamErr
	}

	if daybookNum > configs.GetConfig().App.MaxDaybookLimit {
		return errors.MaxDaybookLimitErr
	}

	return nil
}

// UpdateDiary 更新日记本信息
func UpdateDiary(diary *models.Diary) error {
	// 校验参数
	if err := validateDiary(diary); err != nil {
		return err
	}

	// 更新日记本
	err := diary.Update()
	if err != nil {
		return err
	}

	return nil
}
