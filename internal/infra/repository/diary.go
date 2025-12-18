package repository

import (
	"database/sql"
	"fmt"

	"github.com/biuren/hebe/internal/core/database"
	"github.com/biuren/hebe/internal/domain/diary"
	"github.com/biuren/hebe/internal/utils/errors"
	"github.com/jmoiron/sqlx"
)

// DiaryRepository 日记本数据访问层实现
// 负责处理所有与日记本相关的数据库操作
type DiaryRepository struct {
	db *sqlx.DB
}

// NewDiaryRepository 创建新的日记本数据访问层实例
func NewDiaryRepository() diary.RepositoryInterface {
	db := database.GetDB()
	if db == nil {
		panic("database connection is nil")
	}
	return &DiaryRepository{
		db: db,
	}
}

// Insert 新增日记本到数据库
// 参数: diary - 要插入的日记本对象
// 返回: 新创建的日记本ID和可能的错误
func (r *DiaryRepository) Insert(d *diary.Diary) (int64, error) {
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
		(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)
	`

	res, err := r.db.Exec(sql,
		d.Title, d.Remark, d.AccessURL,
		d.CoverURL, d.CoverMidURL, d.LikeNum,
		d.MemberNum, d.ViewNum,
		d.ViewSetting, d.AccessPwd, d.SocialSetting,
		d.AddUTCTime,
		d.AddUserID, d.AddUserName, d.Status, d.InfoNum,
		d.LastActionUTCTime, d.LastActionUserID, d.LastActionUserName, d.ThemeID)

	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()

	return id, err
}

// Update 更新日记本信息
// 参数: diary - 要更新的日记本对象
// 返回: 可能的错误
func (r *DiaryRepository) Update(d *diary.Diary) error {
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

	result, err := r.db.Exec(sql,
		d.Title, d.Remark,
		d.ViewSetting, d.SocialSetting, d.ThemeID,
		d.LastActionUTCTime, d.LastActionUserID, d.LastActionUserName,
		d.DiaryID)

	if err != nil {
		return err
	}

	// 检查是否有行被更新
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.DiaryNotFound.WithDetails(fmt.Sprintf("diary with ID %d not found", d.DiaryID))
	}

	return nil
}

// Delete 删除日记本
// 参数: diaryID - 日记本ID
// 返回: 可能的错误
func (r *DiaryRepository) Delete(diaryID int64) error {
	sql := "DELETE FROM br_diary WHERE diary_id = ?"
	_, err := r.db.Exec(sql, diaryID)
	return err
}

// IncrementPostNum 增加日记本的日记数量
// 通过空间 ID 把空间信息数量加一，并更新最后动作时间,用户ID,用户名
// 参数: diaryID - 日记本ID, userID - 用户ID, userName - 用户名, actionUTCTime - 操作时间
// 返回: 可能的错误
func (r *DiaryRepository) IncrementPostNum(diaryID int64, userID int64, userName string, actionUTCTime int64) error {
	sqlstr := `
		UPDATE br_diary SET info_num = info_num + 1,
		last_action_utc_time = ?,
		last_action_user_id=?, 
		last_action_user_name = ? 
		WHERE diary_id = ?
	`
	_, err := r.db.Exec(sqlstr, actionUTCTime, userID, userName, diaryID)
	return err
}

// DecrementPostNum 减少日记本的日记数量
// 通过空间 ID 把空间信息数量减一，并更新最后动作时间,用户ID,用户名
// 参数: diaryID - 日记本ID, userID - 用户ID, userName - 用户名, actionUTCTime - 操作时间
// 返回: 可能的错误
func (r *DiaryRepository) DecrementPostNum(diaryID int64, userID int64, userName string, actionUTCTime int64) error {
	sqlstr := `
		UPDATE br_diary SET info_num = info_num - 1,
		last_action_utc_time = ?,
		last_action_user_id=?, 
		last_action_user_name = ? 
		WHERE diary_id = ?
	`
	_, err := r.db.Exec(sqlstr, actionUTCTime, userID, userName, diaryID)
	return err
}

// GetByDiaryID 根据ID获取日记本对象
// 参数: diaryID - 日记本ID
// 返回: 日记本对象和可能的错误
func (r *DiaryRepository) GetByDiaryID(diaryID int64) (*diary.Diary, error) {
	d := &diary.Diary{}
	sqlstr := "SELECT * FROM br_diary WHERE diary_id = ?"
	err := r.db.Get(d, sqlstr, diaryID)
	if err == sql.ErrNoRows {
		return nil, errors.DiaryNotFound.WithDetails("日记本不存在")
	}
	return d, err
}

// ListForWriteByUserID 根据用户ID获取日记本列表，依据最后更新时间排序
// 用在写日记本界面
// 参数: userID - 用户ID
// 返回: 日记本列表和可能的错误
func (r *DiaryRepository) ListForWriteByUserID(userID int64) ([]diary.Diary, error) {
	diaryArray := []diary.Diary{}
	sqlstr := `SELECT * FROM br_diary WHERE add_user_id = ? 
		ORDER BY last_action_utc_time DESC, diary_id DESC`
	err := r.db.Select(&diaryArray, sqlstr, userID)
	if err == sql.ErrNoRows {
		return []diary.Diary{}, nil
	}
	return diaryArray, err
}

// ListByUserID 根据特定用户ID检索分页的日记本列表
// 参数: userID - 用户ID, pageIndex - 页索引, pageSize - 每页大小
// 返回: 日记本列表、总数和可能的错误
func (r *DiaryRepository) ListByUserID(userID int64, pageIndex, pageSize int) ([]diary.Diary, int, error) {
	diaryArray := []diary.Diary{}
	offset := (pageIndex - 1) * pageSize
	sqlstr := `SELECT * FROM br_diary WHERE add_user_id = ?
		ORDER BY last_action_utc_time DESC, diary_id DESC LIMIT ? OFFSET ?`
	err := r.db.Select(&diaryArray, sqlstr, userID, pageSize, offset)
	if err == sql.ErrNoRows {
		return []diary.Diary{}, 0, nil
	}
	if err != nil {
		return nil, 0, err
	}

	count, err := r.GetCountByUserID(userID)
	if err != nil {
		return nil, 0, err
	}

	return diaryArray, count, nil
}

// GetCountByUserID 获取用户的日记本数量
// 参数: userID - 用户ID
// 返回: 日记本数量和可能的错误
func (r *DiaryRepository) GetCountByUserID(userID int64) (int, error) {
	var count int
	sqlstr := "SELECT count(*) FROM br_diary WHERE add_user_id = ?"
	err := r.db.Get(&count, sqlstr, userID)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return count, nil
}

// IsDiaryOwner 验证用户是否是日记本的拥有者
func (r *DiaryRepository) IsDiaryOwner(diaryID, userID int64) (bool, error) {
	var ownerID int64
	sqlstr := "SELECT add_user_id FROM br_diary WHERE diary_id = ?"
	err := r.db.Get(&ownerID, sqlstr, diaryID)
	if err == sql.ErrNoRows {
		return false, errors.DiaryNotFound.WithDetails("日记本不存在")
	}
	if err != nil {
		return false, err
	}
	return ownerID == userID, nil
}

// IncrementLikeNum 增加日记本的点赞数量
func (r *DiaryRepository) IncrementLikeNum(diaryID int64) error {
	sqlstr := "UPDATE br_diary SET like_num = like_num + 1 WHERE diary_id = ?"
	_, err := r.db.Exec(sqlstr, diaryID)
	return err
}

// DecrementLikeNum 减少日记本的点赞数量
func (r *DiaryRepository) DecrementLikeNum(diaryID int64) error {
	sqlstr := "UPDATE br_diary SET like_num = like_num - 1 WHERE diary_id = ? AND like_num > 0"
	_, err := r.db.Exec(sqlstr, diaryID)
	return err
}