package diary

import (
	"database/sql"
	"fmt"

	"github.com/biuren/hebe/internal/core/database"
	errors "github.com/biuren/hebe/internal/utils/error"
	"github.com/jmoiron/sqlx"
)

// RepositoryInterface 日记本数据访问层接口
// 定义所有与日记本相关的数据库操作
type RepositoryInterface interface {
	// Insert 新增日记本到数据库
	Insert(diary *Diary) (int64, error)

	// Update 更新日记本信息
	Update(diary *Diary) error

	// IncrementPostNum 增加日记本的日记数量
	IncrementPostNum(diaryID int64, userID int64, userName string, actionUTCTime int64) error

	// DecrementPostNum 减少日记本的日记数量
	DecrementPostNum(diaryID int64, userID int64, userName string, actionUTCTime int64) error

	// GetByDiaryID 根据ID获取日记本对象
	GetByDiaryID(diaryID int64) (*Diary, error)

	// ListForWriteByUserID 根据用户ID获取日记本列表，依据最后更新时间排序
	ListForWriteByUserID(userID int64) ([]Diary, error)

	// ListByUserID 根据特定用户ID检索分页的日记本列表
	ListByUserID(userID int64, pageIndex, pageSize int) ([]Diary, int, error)

	// GetCountByUserID 获取用户的日记本数量
	GetCountByUserID(userID int64) (int, error)
}

// Repository 日记本数据访问层实现
// 负责处理所有与日记本相关的数据库操作
type Repository struct {
	db *sqlx.DB
}

// NewRepository 创建新的日记本数据访问层实例
func NewRepository() RepositoryInterface {
	db := database.GetDB()
	if db == nil {
		panic("database connection is nil")
	}

	return &Repository{
		db: db,
	}
}

// Insert 新增日记本到数据库
// 参数: diary - 要插入的日记本对象
// 返回: 新创建的日记本ID和可能的错误
func (r *Repository) Insert(diary *Diary) (int64, error) {
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

// Update 更新日记本信息
// 参数: diary - 要更新的日记本对象
// 返回: 可能的错误
func (r *Repository) Update(diary *Diary) error {
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
		diary.Title, diary.Remark,
		diary.ViewSetting, diary.SocialSetting, diary.ThemeID,
		diary.LastActionUTCTime, diary.LastActionUserID, diary.LastActionUserName,
		diary.DiaryID)

	if err != nil {
		return err
	}

	// 检查是否有行被更新
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.DiaryNotFound.WithDetails(fmt.Sprintf("diary with ID %d not found", diary.DiaryID))
	}

	return nil
}

// IncrementPostNum 增加日记本的日记数量
// 通过空间 ID 把空间信息数量加一，并更新最后动作时间,用户ID,用户名
// 参数: diaryID - 日记本ID, userID - 用户ID, userName - 用户名, actionUTCTime - 操作时间
// 返回: 可能的错误
func (r *Repository) IncrementPostNum(diaryID int64, userID int64, userName string, actionUTCTime int64) error {
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
func (r *Repository) DecrementPostNum(diaryID int64, userID int64, userName string, actionUTCTime int64) error {
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
func (r *Repository) GetByDiaryID(diaryID int64) (*Diary, error) {
	diary := &Diary{}
	sqlstr := "SELECT * FROM br_diary WHERE diary_id = ?"
	err := r.db.Get(diary, sqlstr, diaryID)
	if err == sql.ErrNoRows {
		return nil, errors.DiaryNotFound.WithDetails("日记本不存在")
	}
	return diary, err
}

// ListForWriteByUserID 根据用户ID获取日记本列表，依据最后更新时间排序
// 用在写日记本界面
// 参数: userID - 用户ID
// 返回: 日记本列表和可能的错误
func (r *Repository) ListForWriteByUserID(userID int64) ([]Diary, error) {
	diaryArray := []Diary{}
	sqlstr := `SELECT * FROM br_diary WHERE add_user_id = ? 
		ORDER BY last_action_utc_time DESC, diary_id DESC`
	err := r.db.Select(&diaryArray, sqlstr, userID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return diaryArray, err
}

// ListByUserID 根据特定用户ID检索分页的日记本列表
// 参数: userID - 要检索其日记本的用户ID, pageIndex - 页索引（从1开始）, pageSize - 每页项目数
// 返回: 日记本对象切片、总数和可能的错误
func (r *Repository) ListByUserID(userID int64, pageIndex, pageSize int) ([]Diary, int, error) {
	infoArray := []Diary{}
	totalRecords := 0

	// 计算偏移量，注意：页码通常从1开始，所以需要减1
	offset := (pageIndex - 1) * pageSize

	sqlstr := `SELECT * FROM br_diary WHERE add_user_id = ? 
		ORDER BY diary_id DESC LIMIT ? offset ?`

	err := r.db.Select(&infoArray, sqlstr, userID, pageSize, offset)
	if err == sql.ErrNoRows {
		return nil, 0, nil
	}

	//获取总记录数
	totalRecords, err = r.GetCountByUserID(userID)
	if err != nil {
		return nil, 0, err
	}

	return infoArray, totalRecords, err
}

// GetCountByUserID 获取用户的日记本数量
// 参数: userID - 用户ID
// 返回: 日记本数量和可能的错误
func (r *Repository) GetCountByUserID(userID int64) (int, error) {
	var count int
	sqlstr := "SELECT COUNT(*) FROM br_diary WHERE add_user_id = ?"
	err := r.db.Get(&count, sqlstr, userID)
	return count, err
}