package models

import (
	"database/sql"
)

// space 结构
type Space struct {
	SpaceID            int64  `form:"space_id" json:"space_id" db:"space_id"`
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
	IsShowComment      int8   `form:"is_show_comment" json:"is_show_comment" db:"is_show_comment"`
	IsShowLike         int8   `form:"is_show_like" json:"is_show_like" db:"is_show_like"`
	IsShowView         int8   `form:"is_show_view" json:"is_show_view" db:"is_show_view"`
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
func (space *Space) Insert() (int64, error) {
	sql := `
			INSERT INTO br_space
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
			is_show_comment,
			is_show_like,
			is_show_view,
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
			?,
			?,
			?)
		`

	res, err := db.Exec(sql,
		space.Title, space.Remark, space.AccessURL,
		space.CoverURL, space.CoverMidURL, space.LikeNum,
		space.MemberNum, space.ViewNum,
		space.ViewSetting, space.AccessPwd, space.IsShowComment,
		space.IsShowLike, space.IsShowView, space.AddUTCTime,
		space.AddUserID, space.AddUserName, space.Status, space.InfoNum,
		space.LastActionUTCTime, space.LastActionUserID, space.LastActionUserName, space.ThemeID)

	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()

	return id, err
}

// 根据 userID 列出所有空间
func ListSpacesByUserID(userID int64) ([]Space, error) {
	spaceArray := []Space{}
	sqlstr := "SELECT * FROM br_space WHERE add_user_id = ? ORDER BY space_id DESC"
	err := db.Select(&spaceArray, sqlstr, userID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return spaceArray, err
}

// 通过空间 ID 把空间信息数量加一，并更新最后动作时间,用户ID,用户名
func IncrementInfoNum(spaceId int64, userID int64, userName string, actionUTCTime int64) error {
	sqlstr := `
				UPDATE br_space SET info_num = info_num + 1,
				last_action_utc_time = ?,
				last_action_user_id=?, 
				last_action_user_name = ? 
				WHERE space_id = ?
			`
	_, err := db.Exec(sqlstr, actionUTCTime, userID, userName, spaceId)
	return err
}

// 通过空间 ID 把空间信息数量减一，并更新最后动作时间,用户ID,用户名
func DecrementInfoNum(spaceId int64, userID int64, userName string, actionUTCTime int64) error {
	sqlstr := `
				UPDATE br_space SET info_num = info_num - 1,
				last_action_utc_time = ?,
				last_action_user_id=?, 
				last_action_user_name = ? 
				WHERE space_id = ?
			`
	_, err := db.Exec(sqlstr, actionUTCTime, userID, userName, spaceId)
	return err
}

// 根据ID获取space对象
func GetSpaceBySpaceID(spaceID int64) (*Space, error) {
	space := &Space{}
	sqlstr := "SELECT * FROM br_space WHERE space_id = ?"
	err := db.Get(space, sqlstr, spaceID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return space, err
}

// 根据用户ID获取空间列表，依据最后更新，时间排序日记本
func ListSpacesForDiaryByUserID(userID int64) ([]Space, error) {
	spaceArray := []Space{}
	sqlstr := `SELECT * FROM br_space WHERE add_user_id = ? 
				ORDER BY last_action_utc_time DESC, space_id DESC
				`
	err := db.Select(&spaceArray, sqlstr, userID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return spaceArray, err
}
