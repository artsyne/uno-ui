package services

import (
	"github.com/biuren/hebe/errors"
	"github.com/biuren/hebe/models"
	"github.com/biuren/hebe/settings"
)

// 创建一个新的日记本
func CreateNewDaybook(daybook models.Daybook) (int64, error) {

	// 校验参数，标题必须填写，设置访问与社交类型必须为1或者2，主题id必须在1-16之间
	if len(daybook.Title) < 1 || (daybook.ViewSetting != 1 &&
		daybook.ViewSetting != 2 || daybook.SocialSetting != 1 &&
		daybook.SocialSetting != 2 || daybook.ThemeID < 1 || daybook.ThemeID > 16) {
		return 0, errors.ParseInputParamErr
	}

	//判断用户增加日记本数量是否超过限制
	daybookNum, err := models.GetDaybookCountByUserID(daybook.AddUserID)
	if err != nil {
		return 0, errors.ParseInputParamErr
	}
	config := settings.GetConfig()
	if daybookNum > config.MaxDaybookLimit {
		return 0, errors.MaxDaybookLimitErr
	}

	//写入 daybook
	daybookId, err := daybook.Insert()
	if err != nil {
		return 0, err
	}

	return daybookId, nil
}
