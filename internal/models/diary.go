package models

import (
	"database/sql"
)

// space 结构
type Diary struct {
	DiaryID            int64  `form:"diary_id" json:"diary_id" db:"diary_id"`
	Title              string `form:"title" json:"title" db:"title"`
	Remark             string `form:"remark" json:"remark" db:"remark"`
	AccessURL          string `form:"access_url" json:"access_url" db:"access_url"`
	CoverURL           string `form:"cover_url" json:"cover_url" db:"cover_url"`
	CoverMidURL        string `form:"cover_mid_url" json:"cover_mid_url" db:"cover_mid_url"`
	LikeNum            int32  `form:"like_num" json:"like_num" db:"like_num"`
	MemberNum          int32  `form:"member_num" json:"member_num" db:"member_num"`
	ViewNum            int32  `form:"view_num" json:"view_num" db:"view_num"`
	ViewSetting        int8   `form:"view_setting" json:"view_setting" db:"view_setting"`
	AccessPwd          string `form:"access_pwd" json:"access_pwd" db:"access_pwd"`
	SocialSetting      int8   `form:"social_setting" json:"social_setting" db:"social_setting"`
	AddUTCTime         int64  `form:"add_utc_time" json:"add_utc_time" db:"add_utc_time"`
	AddUserID          int64  `form:"add_user_id" json:"add_user_id" db:"add_user_id"`
	AddUserName        string `form:"add_user_name" json:"add_user_name" db:"add_user_name"`
	Status             int8   `form:"status" json:"status" db:"status"`
	InfoNum            int32  `form:"info_num" json:"info_num" db:"info_num"`
	LastActionUTCTime  int64  `form:"last_action_utc_time" json:"last_action_utc_time" db:"last_action_utc_time"`
	LastActionUserID   int64  `form:"last_action_user_id" json:"last_action_user_id" db:"last_action_user_id"`
	LastActionUserName string `form:"last_action_user_name" json:"last_action_user_name" db:"last_action_user_name"`
	ThemeID            int8   `form:"theme_id" json:"theme_id" db:"theme_id"`
}

// 新增space
func (diary *Diary) Insert() (int64, error) {
	sql := `
			INSERT INTO br_diary
			(
			title,
			remark,
			access_url,
			cover_url,
			cover_mid_url,
			like_num,
			member_num,
			view_num,
			view_setting,
			access_pwd,
			social_setting,
			add_utc_time,
			add_user_id,
			add_user_name,
			status,
			info_num,
			last_action_utc_time,
			last_action_user_id,
			last_action_user_name,
			theme_id
			)
			VALUES
			(?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?)
		`

	res, err := db.Exec(sql,
		diary.Title, diary.Remark, diary.AccessURL,
		diary.CoverURL, diary.CoverMidURL, diary.LikeNum,
		diary.MemberNum, diary.ViewNum,
		diary.ViewSetting, diary.AccessPwd, diary.SocialSetting,
		diary.AddUTCTime,
		diary.AddUserID, diary.AddUserName, diary.Status, diary.InfoNum,
		diary.LastActionUTCTime, diary.LastActionUserID, diary.LastActionUserName, diary.ThemeID)

	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()

	return id, err
}

// 更新日记本
func (diary *Diary) Update() error {
	sql := `
		UPDATE br_diary SET 
		title = ?,
		remark = ?,
		view_setting = ?,
		social_setting = ?,
		theme_id = ?,
		last_action_utc_time = ?,
		last_action_user_id = ?,
		last_action_user_name = ?
		WHERE diary_id = ?
	`

	_, err := db.Exec(sql,
		diary.Title, diary.Remark,
		diary.ViewSetting, diary.SocialSetting, diary.ThemeID,
		diary.LastActionUTCTime, diary.LastActionUserID, diary.LastActionUserName,
		diary.DiaryID)

	return err
}

// 通过空间 ID 把空间信息数量加一，并更新最后动作时间,用户ID,用户名
func IncrementPostNum(diaryID int64, userID int64, userName string, actionUTCTime int64) error {
	sqlstr := `
				UPDATE br_diary SET info_num = info_num + 1,
				last_action_utc_time = ?,
				last_action_user_id=?, 
				last_action_user_name = ? 
				WHERE diary_id = ?
			`
	_, err := db.Exec(sqlstr, actionUTCTime, userID, userName, diaryID)
	return err
}

// 通过空间 ID 把空间信息数量减一，并更新最后动作时间,用户ID,用户名
func DecrementPostNum(diaryID int64, userID int64, userName string, actionUTCTime int64) error {
	sqlstr := `
				UPDATE br_diary SET info_num = info_num - 1,
				last_action_utc_time = ?,
				last_action_user_id=?, 
				last_action_user_name = ? 
				WHERE diary_id = ?
			`
	_, err := db.Exec(sqlstr, actionUTCTime, userID, userName, diaryID)
	return err
}

// 根据ID获取space对象
func GetDiaryByDiaryID(diaryID int64) (*Diary, error) {
	diary := &Diary{}
	sqlstr := "SELECT * FROM br_diary WHERE diary_id = ?"
	err := db.Get(diary, sqlstr, diaryID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return diary, err
}

// 根据用户ID获取空间列表，依据最后更新，时间排序日记本
// 用在写日记本界面
func ListDiariesForWriteByUserID(userID int64) ([]Diary, error) {
	diaryArray := []Diary{}
	sqlstr := `SELECT * FROM br_diary WHERE add_user_id = ? 
				ORDER BY last_action_utc_time DESC, diary_id DESC
				`
	err := db.Select(&diaryArray, sqlstr, userID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return diaryArray, err
}

// ListSpaceByUserID 根据特定用户ID检索分页的空间列表。
//
// 参数:
// - userID: 要检索其空间的用户ID。
// - pageIndex: 要检索的页索引（从0开始）。
// - pageSize: 每页的项目数。
//
// 返回值:
// - 一个 diary 对象的切片，表示空间。
// - 一个整数，表示可用空间的总数。
// - 如果在检索过程中有任何问题，则返回错误。
func ListDiariesByUserID(userID int64, pageIndex, pageSize int) ([]Diary, int, error) {
	infoArray := []Diary{}
	totalRecords := 0

	// 计算偏移量，注意：页码通常从1开始，所以需要减1
	offset := (pageIndex - 1) * pageSize

	sqlstr := `SELECT * FROM br_diary WHERE add_user_id = ? 
				ORDER BY diary_id DESC LIMIT ? offset ? `

	err := db.Select(&infoArray, sqlstr, userID, pageSize, offset)
	if err == sql.ErrNoRows {
		return nil, 0, nil
	}

	//获取总记录数
	totalRecords, err = GetDiariesCountByUserID(userID)

	return infoArray, totalRecords, err
}

// 根据用户ID获取空间数量
func GetDiariesCountByUserID(userID int64) (int, error) {
	var count int
	sqlstr := "SELECT COUNT(*) FROM br_diary WHERE add_user_id = ?"

	err := db.Get(&count, sqlstr, userID)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	return count, err
}
