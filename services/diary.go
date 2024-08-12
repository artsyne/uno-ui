package services

import (
	"github.com/biuren/hebe/errors"
	"github.com/biuren/hebe/models"
	"github.com/biuren/hebe/settings"
)

// 创建一个新的日记本
func CreateNewDiary(diary models.Diary) (int64, error) {

	// 校验参数，标题必须填写，设置访问与社交类型必须为1或者2，主题id必须在1-16之间
	if len(diary.Title) < 1 || (diary.ViewSetting != 1 &&
		diary.ViewSetting != 2 || diary.SocialSetting != 1 &&
		diary.SocialSetting != 2 || diary.ThemeID < 1 || diary.ThemeID > 16) {
		return 0, errors.ParseInputParamErr
	}

	//判断用户增加日记本数量是否超过限制
	daybookNum, err := models.GetDaybookCountByUserID(diary.AddUserID)
	if err != nil {
		return 0, errors.ParseInputParamErr
	}
	config := settings.GetConfig()
	if daybookNum > config.MaxDaybookLimit {
		return 0, errors.MaxDaybookLimitErr
	}

	//写入 daybook
	diaryID, err := diary.Insert()
	if err != nil {
		return 0, err
	}

	return diaryID, nil
}
